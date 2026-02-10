# Testing Pyramid for Go TUI Applications

## Overview

```
           /\
          /  \        Level 5: Behavioral (full workflows)
         /----\       Level 4: Snapshot (visual regression)
        /------\      Level 3: Integration (multi-component)
       /--------\     Level 2: Component (Update + View)
      /----------\    Level 1: Pure reducer (fastest, most numerous)
```

**Time investment:** ~60% Level 1, ~20% Level 2, ~10% Level 3, ~5% Level 4, ~5% Level 5.

---

## Level 1: Pure Reducer Tests

**What:** Test `func(State, Action) State` in isolation. **Speed:** Nanoseconds.

### ReducerTest (Single Action)

```go
tuitestkit.RunReducerTests(t, Reduce, []tuitestkit.ReducerTest[AppState, Action]{
    {
        Name:    "move down from top",
        Initial: AppState{Items: items, Cursor: 0},
        Action:  MoveCursor{Delta: 1},
        Assert:  func(t *testing.T, got AppState) {
            if got.Cursor != 1 { t.Errorf("cursor = %d, want 1", got.Cursor) }
        },
    },
    {
        Name:    "clamp at bottom",
        Initial: AppState{Items: items, Cursor: 2},
        Action:  MoveCursor{Delta: 1},
        Assert:  func(t *testing.T, got AppState) {
            if got.Cursor != 2 { t.Errorf("cursor = %d, want 2", got.Cursor) }
        },
    },
})
```

### ReducerSequence (Multi-Step)

Chain actions from the same initial state. `Assert` per step is optional; `Final` runs on end state:

```go
tuitestkit.RunReducerSequences(t, Reduce, []tuitestkit.ReducerSequence[AppState, Action]{
    {
        Name:    "filter then select",
        Initial: AppState{Items: allItems},
        Steps: []tuitestkit.Step[AppState, Action]{
            {Name: "apply filter", Action: SetFilter{Text: "bug"}, Assert: func(t *testing.T, got AppState) {
                if len(got.FilteredItems()) > len(allItems) { t.Error("filter should reduce items") }
            }},
            {Name: "move down", Action: MoveCursor{Delta: 1}},
            {Name: "select", Action: SelectItem{}},
        },
        Final: func(t *testing.T, got AppState) {
            if got.Selected == "" { t.Error("expected item selected") }
        },
    },
})
```

### Invariant Checking

Enforce properties after every reducer call with `WrapWithInvariants`:

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

**Test at Level 1:** every action type, edge cases (empty lists, boundaries), invalid input, state invariants.

---

## Level 2: Component Model Tests (Update + View)

**What:** Test bubbletea Model's `Update` and `View`. **Speed:** Microseconds.

### Send + View Assertions

```go
m := NewBoardModel(testItems)
m = tuitestkit.Send(m, tuitestkit.WindowSize(80, 24))
m = tuitestkit.Send(m, tuitestkit.Key("down"), tuitestkit.Key("down"))

tuitestkit.ViewContains(t, m, "> TASK-3")
tuitestkit.ViewNotContains(t, m, "ERROR")
```

### Line-Level and Regex Assertions

```go
tuitestkit.ViewLineContains(t, m, 0, "Task Board")     // header check
tuitestkit.ViewMatchesRegex(t, m, `TASK-\d+`)          // pattern match

lines := tuitestkit.ViewLines(m)  // []string, ANSI-stripped, trailing blanks trimmed
```

### Mouse Input

```go
m = tuitestkit.Send(m, tuitestkit.MouseClick(10, 5))         // left click
m = tuitestkit.Send(m, tuitestkit.MouseScroll(tuitestkit.ScrollDown))
m = tuitestkit.Send(m, tuitestkit.MouseClickRight(10, 5))    // right click
```

**Test at Level 2:** key bindings, view content correctness, window size adaptation, mouse interactions, empty states.

---

## Level 3: Integration Tests (Multi-Component)

**What:** Test component interaction through commands. Mocks required. **Speed:** Milliseconds.

### SendAndCollect + ExecCmds

The core pattern: send messages, collect returned commands, execute them synchronously, feed results back:

```go
mock := NewBoardMock()
mock.Responses.Set("List", []byte(`[{"id":"TASK-1","title":"Bug fix"}]`), nil)

m := NewBoardModel(mock)
m, cmds := tuitestkit.SendAndCollect(m, tuitestkit.WindowSize(80, 24), tuitestkit.Key("r"))
msgs := tuitestkit.ExecCmds(cmds...)
m = tuitestkit.Send(m, msgs...)

tuitestkit.ViewContains(t, m, "TASK-1")
tuitestkit.AssertCalled(t, &mock.MockCallRecorder, "List")
```

### Error Paths

```go
mock.Responses.SetError("List", fmt.Errorf("permission denied"))
// ... same Send/Collect/Exec cycle ...
tuitestkit.ViewContains(t, m, "permission denied")
```

### Multi-Round Chains

When Init or subsequent messages trigger more commands, repeat the collect/exec cycle:

```go
cmd := m.Init()
msgs := tuitestkit.ExecCmds(cmd)
m, cmds := tuitestkit.SendAndCollect(m, msgs...)
if len(cmds) > 0 {
    msgs = tuitestkit.ExecCmds(cmds...)
    m = tuitestkit.Send(m, msgs...)
}
```

**Test at Level 3:** data loading (success + error), cross-component communication, command chaining.

---

## Level 4: Snapshot Tests (Visual Regression)

**What:** Compare View() output against golden files. **Speed:** Milliseconds.

```go
tuitestkit.SnapshotView(t, m, "board-default")        // ANSI-stripped
tuitestkit.SnapshotViewRaw(t, m, "board-default-raw")  // raw ANSI
tuitestkit.SnapshotStr(t, view, "from-string")         // string variant
```

Golden files stored at `testdata/snapshots/<name>.golden` relative to test file. Create/update:

```bash
UPDATE_SNAPSHOTS=1 go test ./...
```

Use **deterministic data** (no timestamps, no randomness). Always set `WindowSize` explicitly.

**Test at Level 4:** default appearance, empty state, different terminal sizes, error views.

---

## Level 5: Behavioral Tests (Full Workflow)

**What:** Simulate multi-step user journeys. **Speed:** Milliseconds.

```go
m := NewAppModel(mock)
m = tuitestkit.Send(m, tuitestkit.WindowSize(120, 40))

// Load data
msgs := tuitestkit.ExecCmds(m.Init())
m = tuitestkit.Send(m, msgs...)

// Open filter, type, apply
m = tuitestkit.Send(m, tuitestkit.Key("/"))
m = tuitestkit.Send(m, tuitestkit.Key("b"), tuitestkit.Key("u"), tuitestkit.Key("g"))
m = tuitestkit.Send(m, tuitestkit.Key("enter"))
tuitestkit.ViewContains(t, m, "Fix bug")
tuitestkit.ViewNotContains(t, m, "Add feature")

// Navigate and select
m, cmds := tuitestkit.SendAndCollect(m, tuitestkit.Key("down"), tuitestkit.Key("enter"))
msgs = tuitestkit.ExecCmds(cmds...)
m = tuitestkit.Send(m, msgs...)
tuitestkit.ViewContains(t, m, "TASK-3")

// Go back, filter preserved
m = tuitestkit.Send(m, tuitestkit.Key("esc"))
tuitestkit.ViewContains(t, m, "Filter: bug")
```

Combine with `WrapWithInvariants` to catch state corruption across long action sequences.

**Test at Level 5:** critical user journeys (3-5 max), navigation flows, error recovery paths.

---

## Quick Reference: Which Level for What

| Situation | Level | Why |
|-----------|-------|-----|
| New action type | 1 (Reducer) | Fast, all edge cases |
| Cursor/filter logic | 1 (Reducer) | Pure state, table-driven |
| Key binding works | 2 (Component) | Tests Update dispatching |
| View shows data | 2 (Component) | Tests View rendering |
| Mouse click/scroll | 2 (Component) | Update + View for mouse |
| CLI data loads | 3 (Integration) | Mock + command execution |
| Error displays | 3 (Integration) | Mock error + view check |
| Layout looks right | 4 (Snapshot) | Full visual capture |
| User workflow | 5 (Behavioral) | Multi-step interaction |

## Coverage Strategy

1. **Level 1 first.** Every reducer action, every edge case. 60%+ of your logic.
2. **Level 2 for every screen.** Renders without crash, shows content, handles input.
3. **Level 3 for every data path.** Success + error for each CLI operation.
4. **Level 4 for visual stability.** One snapshot per screen per important state.
5. **Level 5 sparingly.** The 3-5 workflows that matter most.

When time is limited: **Level 1 > Level 2 > Level 3**. Pure reducer tests catch the most bugs per line of test code.
