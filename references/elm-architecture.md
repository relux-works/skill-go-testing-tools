# Elm Architecture for Testable TUIs

## Why Elm Architecture Makes TUIs Testable

The Elm architecture separates every application into four pure concepts:

```
State  -->  View(State)  -->  User Input  -->  Action  -->  Reducer(State, Action) --> new State
```

The key insight: **the reducer is a pure function**. Given the same state and the same action, it always produces the same result. No I/O, no randomness, no time dependency. This means you can test every state transition with a simple function call and an assertion.

Compare this to testing through a running UI event loop: you'd need to start a terminal, simulate keystrokes, wait for renders, and scrape output. With Elm, you skip all of that. The reducer is the source of truth for behavior, and it runs in nanoseconds.

## The Cycle: State, Action, Reducer, View

### State

A plain struct holding everything the application needs to render and respond to input:

```go
type AppState struct {
    Items       []string
    Cursor      int
    FilterText  string
    ShowHelp    bool
}
```

State is **inert data**. It doesn't know how to render itself or handle input. It's just a container.

### Action

A value that describes what happened. In Go, use a sum type via interfaces or distinct structs:

```go
type Action interface{}

type MoveCursor struct{ Delta int }
type SetFilter  struct{ Text string }
type ToggleHelp struct{}
type SelectItem struct{}
```

Actions carry the minimum information needed. `MoveCursor{Delta: 1}` says "move down one" without knowing the current position or list length. That's the reducer's job.

### Reducer

A pure function `func(State, Action) State`. This is where all logic lives:

```go
func Reduce(s AppState, a Action) AppState {
    switch act := a.(type) {
    case MoveCursor:
        s.Cursor = clamp(s.Cursor+act.Delta, 0, len(s.Items)-1)
    case SetFilter:
        s.FilterText = act.Text
        s.Cursor = 0 // reset cursor when filter changes
    case ToggleHelp:
        s.ShowHelp = !s.ShowHelp
    case SelectItem:
        // mark current item as selected, etc.
    }
    return s
}
```

The reducer takes old state + action, returns new state. It never calls `exec.Command`, never reads files, never touches the network. That constraint is what makes it testable.

### View

A pure function `func(State) string`. Takes state, returns a string for the terminal:

```go
func View(s AppState) string {
    var b strings.Builder
    for i, item := range s.Items {
        if i == s.Cursor {
            b.WriteString("> " + item + "\n")
        } else {
            b.WriteString("  " + item + "\n")
        }
    }
    return b.String()
}
```

The view is also a pure function. Given the same state, it always produces the same output. This enables snapshot testing.

## How Bubbletea Maps to Elm

Bubbletea uses different names for the same concepts:

| Elm Concept | Bubbletea Equivalent | Role |
|-------------|---------------------|------|
| State | `Model` (the struct) | Data container |
| Action | `tea.Msg` | Event/command description |
| Reducer | `Update(msg) (Model, Cmd)` | State transition logic |
| View | `View() string` | Render state to string |

The mapping is almost 1:1, with one addition: bubbletea's `Update` returns a `tea.Cmd` alongside the new model. Commands represent side effects (fetch data, run a process, quit). The Elm architecture handles this by treating commands as **descriptions of effects**, not the effects themselves. The runtime executes them; the reducer only describes what should happen.

### The Cmd Complication

In pure Elm, the reducer returns `(State, Cmd)` where `Cmd` is a description. Bubbletea follows this: `tea.Cmd` is a `func() tea.Msg`. The runtime calls it asynchronously and feeds the result back as a new message.

For testing, this means:
- **Reducer logic** (state transitions) can be tested as pure functions
- **Commands** need `ExecCmds()` to execute synchronously in tests
- **Side effects** (CLI calls, file I/O) need mocking via the executor pattern

## Composition: Nested Reducers

Real applications have multiple screens or components. Each gets its own state + reducer:

```go
// Top-level state composes sub-states
type AppState struct {
    Screen    ScreenType
    Board     BoardState
    Settings  SettingsState
    Filter    FilterState
}

// Top-level reducer delegates to sub-reducers
func Reduce(s AppState, a Action) AppState {
    switch act := a.(type) {
    case BoardAction:
        s.Board = ReduceBoard(s.Board, act)
    case SettingsAction:
        s.Settings = ReduceSettings(s.Settings, act)
    case FilterAction:
        s.Filter = ReduceFilter(s.Filter, act)
    case SwitchScreen:
        s.Screen = act.Target
    }
    return s
}
```

### Action Wrapping

When a child component produces an action that the parent needs to handle, wrap it:

```go
type BoardAction interface{ boardAction() }
type SettingsAction interface{ settingsAction() }

// In bubbletea terms, the parent's Update examines the msg type
// and routes to the appropriate child's Update
```

This keeps each sub-reducer isolated. You can test `ReduceBoard` independently with only `BoardState` and `BoardAction`, without constructing the entire application state.

### Testing Composed Reducers

Test each sub-reducer in isolation:

```go
tuitestkit.RunReducerTests(t, ReduceBoard, []tuitestkit.ReducerTest[BoardState, BoardAction]{
    {
        Name:    "move cursor down",
        Initial: BoardState{Cursor: 0, Items: items},
        Action:  MoveCursorDown{},
        Assert: func(t *testing.T, got BoardState) {
            if got.Cursor != 1 {
                t.Errorf("cursor = %d, want 1", got.Cursor)
            }
        },
    },
})
```

Then test the top-level reducer for routing correctness (does `BoardAction` reach `ReduceBoard`?) without re-testing the sub-reducer's internal logic.

## Dependency Injection for Side Effects

The reducer must stay pure. But real apps need to call CLI tools, read files, or make HTTP requests. The solution: **extract side effects behind an interface and inject them**.

### The Executor Pattern

```go
// Define an interface for all external operations
type Executor interface {
    RunCommand(cmd string, args ...string) ([]byte, error)
    ReadConfig() (Config, error)
}

// The bubbletea model holds the executor
type Model struct {
    state    AppState
    executor Executor
}

// Update uses the executor to create Cmds, not in the reducer itself
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        // Pure state transition
        m.state = Reduce(m.state, keyToAction(msg))

        // Side effect as a Cmd (description, not execution)
        if needsRefresh(m.state) {
            return m, m.fetchData()
        }
    case DataLoaded:
        // Pure state transition with loaded data
        m.state = Reduce(m.state, SetData{Items: msg.Items})
    }
    return m, nil
}

func (m Model) fetchData() tea.Cmd {
    return func() tea.Msg {
        data, err := m.executor.RunCommand("task-board", "list")
        if err != nil {
            return DataError{Err: err}
        }
        return DataLoaded{Items: parseItems(data)}
    }
}
```

### Testing with Mocks

In tests, inject a mock executor using tuitestkit's building blocks:

```go
type MockExec struct {
    tuitestkit.MockCallRecorder
    Responses *tuitestkit.MockResponseMap
}

func NewMockExec() *MockExec {
    return &MockExec{Responses: tuitestkit.NewMockResponseMap()}
}

func (m *MockExec) RunCommand(cmd string, args ...string) ([]byte, error) {
    m.Record("RunCommand", cmd, args)
    return m.Responses.Get("RunCommand:" + cmd)
}

// In test:
mock := NewMockExec()
mock.Responses.Set("RunCommand:task-board", []byte(`[{"id":"TASK-1"}]`), nil)

model := NewModel(mock)
model, cmds := tuitestkit.SendAndCollect(model, tuitestkit.Key("r"))
msgs := tuitestkit.ExecCmds(cmds...)
model = tuitestkit.Send(model, msgs...)

tuitestkit.ViewContains(t, model, "TASK-1")
tuitestkit.AssertCalled(t, &mock.MockCallRecorder, "RunCommand")
```

The reducer never touches the executor. The `Update` method uses the executor only inside `tea.Cmd` closures. Tests execute those closures synchronously via `ExecCmds`, keeping everything deterministic.

## Summary

| Principle | What It Means for Testing |
|-----------|--------------------------|
| Pure reducer | Test state transitions as simple function calls |
| Actions as values | Construct any scenario by creating action values |
| View as pure function | Snapshot test any visual state by constructing the state directly |
| Cmd as description | Execute commands synchronously in tests with `ExecCmds` |
| Dependency injection | Mock all external operations, verify calls with `AssertCalled` |
| Nested reducers | Test each component in isolation, test routing at the top level |
