# Research: Reducer Test Harness

## Analyzed Sources

Real-world reducer tests from `board-tui` project (6 test files, ~450 lines total).

### Pattern Analysis

#### 1. Reducer Signatures Found

Two patterns exist in the codebase:

**Value-receiver reducers (components):**
```go
func Reduce(state State, action Action) State          // command, filter, dialog, board
```

**Pointer-receiver reducer (app-level):**
```go
func Reduce(state *AppState, action Action) *AppState  // top-level state
```

The harness must support both via generics. The pointer variant is less common but exists at the app level.

**Decision:** The generic type parameter `S` can be either value or pointer. For the pointer case, users pass `*AppState` as `S`. The `Reduce` function signature in the harness is simply `func(S, A) S` which works for both.

#### 2. Common Test Patterns

**Pattern A: Simple action -> single assertion**
Most common (~70% of tests). One initial state, one action, one or two fields checked.
```go
s := NewAppState()
s = Reduce(s, SetScreen{Screen: ScreenSettings})
if s.Screen != ScreenSettings { t.Error("...") }
```

**Pattern B: Sequential actions (multi-step)**
~20% of tests. Apply multiple actions, assert after each.
```go
s = Reduce(s, SelectDown{})    // assert idx == 1
s = Reduce(s, SelectDown{})    // assert idx == 2
s = Reduce(s, SelectDown{})    // assert idx == 0 (wrap)
```

**Pattern C: State setup + action + deep assertion**
~10% of tests. More complex initial state, action, then multi-field assertion.
```go
s := NewAppState()
s.Shared.Config = &Config{}
statuses := []string{"development", "to-review"}
s = Reduce(s, SetAgentsStatusFilter{Statuses: statuses})
// assert len(s.Shared.AgentsStatusFilter) == 2
// assert s.Shared.Config.AgentsStatusFilter updated too
```

#### 3. Action Type Patterns

Actions use a marker interface pattern:
```go
type Action interface { actionMarker() }
type action struct{}
func (action) actionMarker() {}
```

Each action is a struct embedding the marker. This means `A` in the generic must be an interface, and concrete actions are passed as values.

#### 4. Nested Reducer Delegation

Board reducer delegates to child reducers (filter, command, dialog) via wrapper actions:
```go
case FilterAction:
    state.Filter = filter.Reduce(state.Filter, a.Action)
```

This means tests at different levels use different `(S, A)` pairs. The harness must be generic over any `(S, A)` combination.

#### 5. Assertion Patterns

- Simple field equality checks (`if s.X != expected`)
- Boolean state checks (`if !s.Active`)
- Nil/non-nil checks (`if s.Config != nil`)
- Length checks (`if len(s.Items) != 2`)
- No deep-equal patterns found -- all assertions are manual field checks.

**Decision:** Use `func(S) error` assertion functions. Keeps it simple, no need for matcher DSL. Users write their own checks.

#### 6. Initial State Patterns

- Constructor function: `NewAppState()`, `Initial(commands)`, `Initial()`
- Inline struct literal: `State{Active: true, Input: "old"}`
- Constructor + mutation: `s := NewAppState(); s.Shared.Config = &Config{}`

**Decision:** `Initial S` field in test case. Users can use any of the above.

### Invariant Patterns Identified

Cross-cutting properties that should hold after ANY action:

1. **Index bounds:** `SelectedIdx` should never exceed `len(items)-1` or go negative
2. **Mutual exclusion:** Dialog open + command active should never both be true simultaneously
3. **State consistency:** If `Loading == false` and `Error == nil`, data should be present
4. **Config sync:** When config exists, certain shared fields must match config fields

These are perfect candidates for `InvariantChecker` -- register once, auto-check after every `Reduce()`.

### Existing Go Libraries Reviewed

- **Standard table-driven tests** (Go wiki pattern) -- our baseline, but we add generics + invariants
- **shoenig/test** -- generic assertions, but assertion-focused, not reducer-focused
- **alecthomas/assert** -- simple generic assertions, no reducer framework
- **flyingmutant/rapid** -- property-based/stateful testing, overkill for our use case
- **go-quicktest/qt** -- good generic helpers, but again assertion-focused

**Decision:** No existing library provides reducer-specific table-driven testing with invariant checking. This is a genuine gap. Our library fills it.

### API Design

#### ReducerTest

```go
type ReducerTest[S any, A any] struct {
    Name    string
    Initial S
    Action  A
    Assert  func(t *testing.T, got S)
}

func RunReducerTests[S any, A any](
    t *testing.T,
    reduce func(S, A) S,
    tests []ReducerTest[S, A],
)
```

Uses `t.Run(test.Name, ...)` for sub-test isolation. Each test: apply action to initial, run assert.

#### ReducerSequence (for multi-step)

```go
type Step[S any, A any] struct {
    Name   string           // optional label
    Action A
    Assert func(t *testing.T, got S)  // optional per-step assertion
}

type ReducerSequence[S any, A any] struct {
    Name    string
    Initial S
    Steps   []Step[S, A]
    Final   func(t *testing.T, got S)  // optional final assertion
}

func RunReducerSequences[S any, A any](
    t *testing.T,
    reduce func(S, A) S,
    sequences []ReducerSequence[S, A],
)
```

#### InvariantChecker

```go
type Invariant[S any] struct {
    Name  string
    Check func(S) error
}

type InvariantChecker[S any] struct {
    invariants []Invariant[S]
}

func NewInvariantChecker[S any](invariants ...Invariant[S]) *InvariantChecker[S]

func (c *InvariantChecker[S]) CheckAll(t *testing.T, state S)

// Integration: wrap a reducer so invariants auto-check
func (c *InvariantChecker[S]) Wrap(reduce func(S, A) S) func(S, A) S
```

The `Wrap` method returns a new reducer that runs the original then checks all invariants. This way:
```go
checker := NewInvariantChecker[State](...)
reduce := checker.Wrap(command.Reduce)
// Now every Reduce() call auto-checks invariants
RunReducerTests(t, reduce, tests)
```

**Problem with Wrap:** `InvariantChecker[S]` doesn't know about `A`. Two options:
1. Make `Wrap` a standalone function: `func WrapWithInvariants[S, A any](reduce func(S,A)S, checker *InvariantChecker[S]) func(S,A)S`
2. Make `Wrap` generic on the method level -- not possible in Go (methods can't have their own type params).

**Decision:** Standalone function `WrapWithInvariants`.

### Task Breakdown

1. **Define core types** -- `ReducerTest[S,A]`, `Step[S,A]`, `ReducerSequence[S,A]`
2. **Implement RunReducerTests** -- table-driven runner with t.Run
3. **Implement RunReducerSequences** -- sequential multi-step runner
4. **Implement InvariantChecker** -- invariant registration + CheckAll
5. **Implement WrapWithInvariants** -- reducer wrapping for auto-checking
6. **Tests for the harness itself** -- test RunReducerTests, RunReducerSequences, InvariantChecker using a simple counter reducer
