# Mock Patterns for Go Testing

## CLI Executor Interface Extraction

Direct `exec.Command` inside `Update` makes testing impossible. Extract behind an interface:

```go
// Interface: narrow, only what this model needs
type BoardExecutor interface {
    List() ([]byte, error)
    Show(id string) ([]byte, error)
    UpdateStatus(id string, status string) error
}

// Production implementation
type CLIExecutor struct{}
func (e *CLIExecutor) List() ([]byte, error) {
    return exec.Command("task-board", "list", "--json").Output()
}

// Model accepts the interface
type Model struct {
    state    State
    executor BoardExecutor
}
```

The model uses the executor only inside `tea.Cmd` closures. The reducer stays pure.

## When to Mock vs When to Test Pure Functions

| Situation | Approach |
|-----------|----------|
| Reducer logic (state transitions) | Direct test, no mock |
| View rendering | Direct test with `ViewContains`, no mock |
| Key-to-action mapping | Direct test, no mock |
| CLI command execution | Mock the executor |
| File reading/writing | Mock the file interface |
| Timer/tick behavior | Test the tick message handler directly |

Rule of thumb: if the signature is `func(State, Action) State`, test directly. Mock only at the I/O boundary.

## Test Doubles Taxonomy

**Stub** -- returns canned data, doesn't track calls. Good for rendering tests:

```go
type StubExecutor struct{}
func (s *StubExecutor) List() ([]byte, error) {
    return []byte(`[{"id":"TASK-1","title":"Fix bug"}]`), nil
}
```

**Mock** -- records calls + provides canned responses. Good for interaction verification:

```go
type MockBoardExecutor struct {
    tuitestkit.MockCallRecorder
    Responses *tuitestkit.MockResponseMap
}
```

**Fake** -- has real (simplified) logic, maintains state. Good for integration tests:

```go
type FakeExecutor struct{ items map[string]Item }
func (f *FakeExecutor) UpdateStatus(id, status string) error {
    item := f.items[id]; item.Status = status; f.items[id] = item; return nil
}
```

**Choose:** need data? Stub. Need to verify interactions? Mock. Need realistic multi-step? Fake.

## MockCallRecorder + MockResponseMap Composition

tuitestkit provides two embeddable building blocks:

### Building a Mock

```go
type TaskBoardMock struct {
    tuitestkit.MockCallRecorder                    // embedded: Record(), CallCount(), Calls()
    Responses *tuitestkit.MockResponseMap           // configurable canned responses
}

func NewTaskBoardMock() *TaskBoardMock {
    return &TaskBoardMock{Responses: tuitestkit.NewMockResponseMap()}
}

// Each method: Record + Get. Two lines per mock method.
func (m *TaskBoardMock) TreeJSON() ([]byte, error) {
    m.Record("TreeJSON")
    return m.Responses.Get("TreeJSON")
}

func (m *TaskBoardMock) Execute(cmd string, args ...string) ([]byte, error) {
    m.Record("Execute", cmd, args)
    return m.Responses.Get("Execute:" + cmd)
}
```

### Configuring Responses

```go
mock := NewTaskBoardMock()
mock.Responses.Set("TreeJSON", []byte(`{"epics":[]}`), nil)          // success
mock.Responses.SetError("TreeJSON", fmt.Errorf("not found"))          // error
mock.Responses.Set("Execute:ls", []byte("file1\nfile2"), nil)        // keyed by arg
// Get returns (nil, nil) for unknown keys -- silent fallback
```

### Asserting on Calls

```go
tuitestkit.AssertCalled(t, &mock.MockCallRecorder, "TreeJSON")
tuitestkit.AssertCalledN(t, &mock.MockCallRecorder, "TreeJSON", 1)
tuitestkit.AssertCalledWith(t, &mock.MockCallRecorder, "Execute", "ls", []string{"-la"})
tuitestkit.AssertNotCalled(t, &mock.MockCallRecorder, "Delete")

// Detailed inspection
calls := mock.CallsFor("Execute")
for _, c := range calls { fmt.Println(c.Method, c.Args) }

// Reset between test cases
mock.Reset()
```

## Callback Capture Pattern

When components communicate via callbacks instead of through the parent reducer:

```go
type CallbackCapture[T any] struct {
    mu     sync.Mutex
    values []T
}

func (c *CallbackCapture[T]) Capture() func(T) {
    return func(v T) {
        c.mu.Lock()
        defer c.mu.Unlock()
        c.values = append(c.values, v)
    }
}

func (c *CallbackCapture[T]) Values() []T {
    c.mu.Lock()
    defer c.mu.Unlock()
    out := make([]T, len(c.values))
    copy(out, c.values)
    return out
}
```

Usage:

```go
capture := &CallbackCapture[string]{}
model := NewModel(WithOnSelect(capture.Capture()))
model = tuitestkit.Send(model, tuitestkit.Key("enter"))

if vals := capture.Values(); len(vals) != 1 || vals[0] != "TASK-1" {
    t.Errorf("expected callback with TASK-1, got %v", vals)
}
```

## Full Example: Mock Lifecycle

```go
func TestBoardScreen_RefreshLoadsData(t *testing.T) {
    // 1. Create and configure mock
    mock := NewTaskBoardMock()
    mock.Responses.Set("TreeJSON", []byte(`{"epics":[{"id":"EPIC-1","title":"Feature X"}]}`), nil)

    // 2. Create model with mock
    model := NewBoardScreen(mock)

    // 3. Send messages, collect commands
    model, cmds := tuitestkit.SendAndCollect(model,
        tuitestkit.WindowSize(80, 24),
        tuitestkit.Key("r"),
    )

    // 4. Execute commands (triggers mock.TreeJSON())
    msgs := tuitestkit.ExecCmds(cmds...)
    model = tuitestkit.Send(model, msgs...)

    // 5. Assert view and interactions
    tuitestkit.ViewContains(t, model, "EPIC-1")
    tuitestkit.AssertCalled(t, &mock.MockCallRecorder, "TreeJSON")
}

func TestBoardScreen_HandlesError(t *testing.T) {
    mock := NewTaskBoardMock()
    mock.Responses.SetError("TreeJSON", fmt.Errorf("board not found"))

    model := NewBoardScreen(mock)
    model, cmds := tuitestkit.SendAndCollect(model, tuitestkit.WindowSize(80, 24), tuitestkit.Key("r"))
    msgs := tuitestkit.ExecCmds(cmds...)
    model = tuitestkit.Send(model, msgs...)

    tuitestkit.ViewContains(t, model, "board not found")
}
```

## Key Takeaways

1. **Extract interfaces first.** Narrow, per-model. Don't create god interfaces.
2. **Compose, don't inherit.** Embed `MockCallRecorder` + use `MockResponseMap`.
3. **Two lines per mock method:** `m.Record(...)` then `return m.Responses.Get(...)`.
4. **Use assertion helpers** (`AssertCalled`, `AssertCalledWith`) -- don't manually inspect call slices.
5. **Stubs for rendering, mocks for interactions, fakes for integration.** Pick the simplest double.
