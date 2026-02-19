---
name: go-testing-tools
description: Go testing toolkit with tuitestkit library. Closed-loop agent development cycle for Go applications. Covers test helpers, reducer tests, snapshot/golden file testing, mock executor patterns, view assertions. Triggers: Go test, bubbletea test, tuitestkit, go test helpers, reducer test, snapshot test go, golden file go, тесты го, го тест хелперы, тесты баблти, тест утилиты го, замкнутый цикл тестирования.
---

# Go Testing Tools

Testing toolkit for Go TUI applications built with bubbletea. Provides a Go library (`tuitestkit`), reference documentation, and templates that enable a closed-loop development cycle: write code, write tests, run tests, validate, fix, repeat -- all without human intervention.

Core philosophy: **Elm architecture makes TUI testable.** The reducer is a pure function -- given the same state and action, it always returns the same result. No I/O, no randomness, no terminal. Test everything as function calls.

---

## Prerequisites

- **Go 1.21+**
- **bubbletea** (`github.com/charmbracelet/bubbletea`)
- **lipgloss** (`github.com/charmbracelet/lipgloss`) -- used by most bubbletea apps for styling

Install the library in your project:

```bash
go get github.com/relux-works/skill-go-testing-tools/tuitestkit
```

---

## Quick Start

Minimal test showing the core workflow: build a message, send it to a model, assert on the view.

```go
package myscreen_test

import (
    "testing"

    "github.com/relux-works/skill-go-testing-tools/tuitestkit"
)

func TestBasicNavigation(t *testing.T) {
    m := NewMyModel(testData)
    m = tuitestkit.Send(m, tuitestkit.WindowSize(80, 24))

    // Press down twice
    m = tuitestkit.Send(m, tuitestkit.Key("down"), tuitestkit.Key("down"))

    // Assert the view contains expected content
    tuitestkit.ViewContains(t, m, "third item")
    tuitestkit.ViewNotContains(t, m, "ERROR")
}
```

---

## Testing Pyramid

Five levels, from fastest/most numerous to slowest/fewest:

| Level | What | Speed | Investment |
|-------|------|-------|------------|
| 1. Pure Reducer | `func(State, Action) State` | Nanoseconds | ~60% |
| 2. Component | `Update` + `View` | Microseconds | ~20% |
| 3. Integration | Multi-component with mocks | Milliseconds | ~10% |
| 4. Snapshot | Golden file comparison | Milliseconds | ~5% |
| 5. Behavioral | Full user workflows | Milliseconds | ~5% |

When time is limited: **Level 1 > Level 2 > Level 3**. Pure reducer tests catch the most bugs per line of test code.

Full details: `references/testing-pyramid.md`

---

## Library API Reference

All functions are in the `tuitestkit` package.

### Message Builders (`messages.go`)

Build `tea.Msg` values for driving model updates.

```go
// Key builds a tea.KeyMsg from a human-readable string.
// Supports: "enter", "tab", "esc", "space", "backspace", "up", "down",
// "left", "right", "home", "end", "pgup", "pgdown", "delete", "insert",
// "f1"-"f20", "ctrl+c", "ctrl+a"..."ctrl+z", "alt+h", "alt+enter",
// single runes: "a", "/", "1", etc.
func Key(k string) tea.KeyMsg

// Keys builds a slice of tea.Msg from multiple key strings.
func Keys(keys ...string) []tea.Msg

// WindowSize builds a tea.WindowSizeMsg with the given dimensions.
func WindowSize(w, h int) tea.WindowSizeMsg

// MouseClick builds a tea.MouseMsg for a left-button click at (x, y).
func MouseClick(x, y int) tea.MouseMsg

// MouseClickRight builds a tea.MouseMsg for a right-button click at (x, y).
func MouseClickRight(x, y int) tea.MouseMsg

// MouseScroll builds a tea.MouseMsg for a scroll event.
// Direction constants: ScrollUp, ScrollDown, ScrollLeft, ScrollRight.
func MouseScroll(dir ScrollDir) tea.MouseMsg

// MouseRelease builds a tea.MouseMsg for a button release at (x, y).
func MouseRelease(x, y int) tea.MouseMsg
```

### Model Harness (`harness.go`)

Send messages to bubbletea models and handle commands.

```go
// Send sends messages to a model sequentially, returns the final model.
// Preserves concrete type via generics -- no type assertion needed.
// Panics if Update returns a different type (broken implementation).
// Init() is NOT called.
func Send[M tea.Model](model M, msgs ...tea.Msg) M

// SendAndCollect sends messages, collects all non-nil Cmds returned by Update.
// Returns final model + collected commands.
func SendAndCollect[M tea.Model](model M, msgs ...tea.Msg) (M, []tea.Cmd)

// ExecCmds executes tea.Cmd functions synchronously, collects resulting messages.
// Recursively handles tea.BatchMsg. Nil cmds/messages are skipped.
func ExecCmds(cmds ...tea.Cmd) []tea.Msg
```

**Common pattern -- command pipeline:**

```go
m, cmds := tuitestkit.SendAndCollect(m, tuitestkit.Key("r"))
msgs := tuitestkit.ExecCmds(cmds...)
m = tuitestkit.Send(m, msgs...)
```

### Reducer Harness (`reducer.go`)

Table-driven testing for pure reducer functions.

```go
// ReducerTest defines a single test case for a pure reducer.
// S = state type, A = action type.
type ReducerTest[S any, A any] struct {
    Name    string
    Initial S
    Action  A
    Assert  func(t *testing.T, got S)
}

// Step defines one step in a multi-action reducer sequence.
// Assert is optional -- if nil, applied without per-step validation.
type Step[S any, A any] struct {
    Name   string
    Action A
    Assert func(t *testing.T, got S)
}

// ReducerSequence defines a multi-step test scenario.
// Actions applied sequentially from Initial. Final runs on end state.
type ReducerSequence[S, A any] struct {
    Name    string
    Initial S
    Steps   []Step[S, A]
    Final   func(t *testing.T, got S)
}

// RunReducerTests executes table-driven reducer tests as subtests.
func RunReducerTests[S, A any](t *testing.T, reduce func(S, A) S, tests []ReducerTest[S, A])

// RunReducerSequences executes multi-step sequence tests.
// Per-step Assert runs after each step (if non-nil). Final runs on end state.
func RunReducerSequences[S, A any](t *testing.T, reduce func(S, A) S, sequences []ReducerSequence[S, A])
```

**Invariant checking:**

```go
// Invariant defines a property that must hold for any state.
type Invariant[S any] struct {
    Name  string
    Check func(s S) error
}

// InvariantChecker holds invariants and validates state against all of them.
type InvariantChecker[S any] struct { /* ... */ }

// NewInvariantChecker creates a checker with the given invariants.
func NewInvariantChecker[S any](invariants ...Invariant[S]) *InvariantChecker[S]

// Check validates state against all invariants. Returns combined error or nil.
func (ic *InvariantChecker[S]) Check(s S) error

// WrapWithInvariants wraps a reducer with invariant checking.
// After every reduce call, all invariants are validated.
// t.Fatalf on violation.
func WrapWithInvariants[S, A any](t *testing.T, reduce func(S, A) S, checker *InvariantChecker[S]) func(S, A) S
```

**Example:**

```go
checker := tuitestkit.NewInvariantChecker(
    tuitestkit.Invariant[AppState]{
        Name: "cursor in bounds",
        Check: func(s AppState) error {
            if s.Cursor < 0 || s.Cursor >= len(s.Items) {
                return fmt.Errorf("cursor %d out of [0, %d)", s.Cursor, len(s.Items))
            }
            return nil
        },
    },
)
safeReduce := tuitestkit.WrapWithInvariants(t, Reduce, checker)
tuitestkit.RunReducerTests(t, safeReduce, tests)
```

### Mock Building Blocks (`mock.go`)

Composable primitives for building project-specific mocks.

```go
// MockCall represents a single recorded method invocation.
type MockCall struct {
    Method string
    Args   []any
}

// MockCallRecorder provides thread-safe recording of method calls.
// Embed into project-specific mock structs.
type MockCallRecorder struct { /* ... */ }

func (r *MockCallRecorder) Record(method string, args ...any)
func (r *MockCallRecorder) CallCount(method string) int
func (r *MockCallRecorder) Calls() []MockCall
func (r *MockCallRecorder) CallsFor(method string) []MockCall
func (r *MockCallRecorder) Reset()

// MockResponseMap provides thread-safe canned response storage.
type MockResponseMap struct { /* ... */ }

func NewMockResponseMap() *MockResponseMap
func (m *MockResponseMap) Set(key string, data []byte, err error)
func (m *MockResponseMap) Get(key string) ([]byte, error)
func (m *MockResponseMap) SetError(key string, err error)

// Assertion helpers
func AssertCalled(t testing.TB, r *MockCallRecorder, method string)
func AssertNotCalled(t testing.TB, r *MockCallRecorder, method string)
func AssertCalledN(t testing.TB, r *MockCallRecorder, method string, n int)
func AssertCalledWith(t testing.TB, r *MockCallRecorder, method string, args ...any)
```

**Building a mock (two lines per method):**

```go
type TaskBoardMock struct {
    tuitestkit.MockCallRecorder
    Responses *tuitestkit.MockResponseMap
}

func NewTaskBoardMock() *TaskBoardMock {
    return &TaskBoardMock{Responses: tuitestkit.NewMockResponseMap()}
}

func (m *TaskBoardMock) TreeJSON() ([]byte, error) {
    m.Record("TreeJSON")
    return m.Responses.Get("TreeJSON")
}

func (m *TaskBoardMock) Execute(cmd string, args ...string) ([]byte, error) {
    m.Record("Execute", cmd, args)
    return m.Responses.Get("Execute:" + cmd)
}
```

### View Assertions (`view.go`)

Assert on rendered output. All functions strip ANSI escape codes before comparison.

```go
// StripANSI removes all ANSI escape sequences from a string.
func StripANSI(s string) string

// --- Model-based (call model.View() internally) ---

func ViewContains(t testing.TB, model tea.Model, text string)
func ViewNotContains(t testing.TB, model tea.Model, text string)
func ViewLines(model tea.Model) []string
func ViewLineContains(t testing.TB, model tea.Model, lineIdx int, text string)
func ViewLineEquals(t testing.TB, model tea.Model, lineIdx int, text string)
func ViewMatchesRegex(t testing.TB, model tea.Model, pattern string)

// --- String-based (when you already have the view string) ---

func ContainsStr(t testing.TB, view string, text string)
func NotContainsStr(t testing.TB, view string, text string)
func LinesFromStr(view string) []string
func MatchesRegexStr(t testing.TB, view string, pattern string)
```

### Snapshot Testing (`snapshot.go`)

Golden file comparison for visual regression testing.

```go
// UpdateSnapshots controls whether snapshots are written (true) or compared (false).
// Set via UPDATE_SNAPSHOTS=1 environment variable.
var UpdateSnapshots bool

// SnapshotView captures model.View(), strips ANSI, compares against golden file.
func SnapshotView(t *testing.T, model tea.Model, name string)

// SnapshotViewRaw captures model.View() with raw ANSI codes intact.
func SnapshotViewRaw(t *testing.T, model tea.Model, name string)

// SnapshotStr compares a pre-rendered string (ANSI-stripped) against golden file.
func SnapshotStr(t *testing.T, view string, name string)

// SnapshotStrRaw compares a pre-rendered string (raw ANSI) against golden file.
func SnapshotStrRaw(t *testing.T, view string, name string)
```

Golden files stored at `testdata/snapshots/<name>.golden` relative to the test file.

```bash
# Create or update golden files
UPDATE_SNAPSHOTS=1 go test ./...

# Update for a single package
UPDATE_SNAPSHOTS=1 go test ./internal/ui/screens/board/ -run TestSnapshot
```

Mismatches produce a line-by-line diff with `-`/`+` markers.

---

## Closed-Loop Workflow

The agent development cycle for writing tests:

```
1. UNDERSTAND: Read screen/component code
2. IDENTIFY:   What behaviors need tests?
3. EXTRACT:    If exec.Command present, extract to interface (use mock template)
4. WRITE:      Tests using tuitestkit helpers
5. RUN:        go test ./... -v
6. VALIDATE:   All pass? Check coverage with go test -cover
7. FIX:        If failures, read error, fix code or test, goto 5
8. SNAPSHOT:   Optionally capture view golden files
9. DONE:       Tests green, coverage acceptable
```

### Step Details

**UNDERSTAND:** Read the component source. Identify:
- What state does it manage?
- What messages does it handle in Update?
- Does it have a pure reducer (`func(State, Action) State`)?
- Does it shell out to external commands (exec.Command)?
- What does View() render?

**IDENTIFY:** Map behaviors to test levels:
- Pure state logic -> Level 1 (reducer tests)
- Key bindings, view rendering -> Level 2 (component tests)
- Data loading, command execution -> Level 3 (integration tests with mocks)
- Visual appearance -> Level 4 (snapshot tests)

**EXTRACT:** If the component uses `exec.Command` or any I/O directly:
1. Define an interface with only the methods this component needs
2. Move `exec.Command` calls to a production implementation
3. Build a mock using `MockCallRecorder` + `MockResponseMap`
4. See `assets/templates/executor_interface.go.tmpl` and `assets/templates/mock_executor.go.tmpl`

**WRITE:** Use the appropriate tuitestkit helpers per level:
- Level 1: `RunReducerTests`, `RunReducerSequences`, `WrapWithInvariants`
- Level 2: `Send`, `ViewContains`, `ViewLines`, `Key`, `MouseClick`
- Level 3: `SendAndCollect`, `ExecCmds`, `AssertCalled`, `AssertCalledWith`
- Level 4: `SnapshotView`, `SnapshotStr`

**RUN:** `go test ./... -v` -- all tests must run in under 1 second (no I/O, no network).

**VALIDATE:** `go test ./... -cover` -- target 80%+ for code under test.

**FIX:** Read error messages carefully. Common issues:
- Type assertion panic -> `Update` returns wrong type (broken implementation)
- View assertion failure -> check ANSI stripping, check exact text
- Snapshot mismatch -> run with `UPDATE_SNAPSHOTS=1` if the change is intentional
- Mock assertion failure -> check response key matches (`"Execute:list"` not `"Execute"`)

---

## Templates

Scaffolding templates in `assets/templates/`:

| Template | Purpose |
|----------|---------|
| `reducer_test.go.tmpl` | Table-driven reducer tests with invariants |
| `component_test.go.tmpl` | Component model tests (keyboard, mouse, commands, view) |
| `snapshot_test.go.tmpl` | Snapshot golden file tests |
| `executor_interface.go.tmpl` | Executor interface + production stub |
| `mock_executor.go.tmpl` | Mock executor using MockCallRecorder + MockResponseMap |
| `PROJECT_STRUCTURE.md` | Recommended directory layout for test files |

Copy templates into your project, replace `TODO` markers with your actual types and import paths.

---

## References

In-depth documentation in `references/`:

| Document | Content |
|----------|---------|
| `elm-architecture.md` | Why Elm makes TUI testable, state/action/reducer/view cycle, nested reducers, dependency injection, the executor pattern |
| `mock-patterns.md` | Interface extraction, when to mock vs test pure functions, stubs vs mocks vs fakes, callback capture, full mock lifecycle example |
| `testing-pyramid.md` | All 5 test levels with code examples, decision table for which level to use, coverage strategy |

---

## Setup

Run the setup script to create symlinks:

```bash
./scripts/setup.sh
```

This links the skill to `~/.agents/skills/go-testing-tools`, `~/.claude/skills/go-testing-tools`, and `~/.codex/skills/go-testing-tools`.

To remove:

```bash
./scripts/deinit.sh
```
