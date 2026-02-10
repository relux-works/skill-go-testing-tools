# Research: Mock Executor for tuitestkit

## Problem

TUI apps shell out to external CLI tools via `exec.Command`. To unit-test screens that depend on CLI output, you need:
1. An executor interface (project-specific -- user defines their own)
2. A mock executor (for tests -- **this is what tuitestkit provides**)

The mock needs to: return canned responses, record calls, simulate errors -- all without generics/reflection magic.

---

## Real-World Example: board-tui's MockCLIExecutor

The board-tui project already has a hand-rolled mock at `tools/board-tui/internal/cli/mock_executor.go`. Key patterns:

### Interface (7 methods)
```go
type CLIExecutor interface {
    TreeJSON() ([]byte, error)
    ListTasksJSON() ([]byte, error)
    ListStoriesJSON() ([]byte, error)
    ShowElementJSON(id string) ([]byte, error)
    AgentsJSON(staleMinutes int) ([]byte, error)
    Validate() error
    SelfUpdate() (stdout string, stderr string, err error)
}
```

### Mock Pattern
```go
type MockCLIExecutor struct {
    mu sync.Mutex
    // Canned response fields per method
    TreeJSONResponse  []byte
    TreeJSONError     error
    // ... etc for each method
    ShowElementJSONResponses map[string][]byte  // keyed by ID for parameterized methods
    Calls []MockCall  // call recording
}

type MockCall struct {
    Method string
    Args   []interface{}
}
```

### What Works Well
- **Simple struct fields** for canned responses -- no builder chain needed
- **Call recording** with method name + args -- useful for assertions
- **`CallCount(method)` and `CallsFor(method)`** helper methods
- **Thread-safe** with `sync.Mutex` (important for bubbletea Cmd goroutines)
- **Compile-time check**: `var _ CLIExecutor = (*MockCLIExecutor)(nil)`

### What's Boilerplate
Every project would rewrite:
- `MockCall` struct
- `record()` method
- `CallCount()` / `CallsFor()` helpers
- Mutex locking pattern

---

## Go Mock Approaches Compared

### 1. testify/mock (github.com/stretchr/testify)
```go
type MockExec struct { mock.Mock }
func (m *MockExec) TreeJSON() ([]byte, error) {
    args := m.Called()
    return args.Get(0).([]byte), args.Error(1)
}
// Usage: mock.On("TreeJSON").Return([]byte(`{}`), nil)
```
- **Pros**: Builder chain, flexible matchers, call expectations
- **Cons**: External dependency (tuitestkit spec says "zero non-test dependencies beyond bubbletea/lipgloss"), runtime type assertions, stringly-typed method names

### 2. gomock (go.uber.org/mock)
- Code generation via mockgen
- **Cons**: Requires separate tool, generated code, heavy

### 3. Hand-rolled (board-tui pattern)
- No dependencies, full control
- **Cons**: Boilerplate per project

### 4. tuitestkit approach: **Reusable building blocks + templates**
The sweet spot -- provide the infrastructure pieces that every mock needs, let the user compose them for their specific interface.

---

## Design Decision: Building Blocks, Not Generic Framework

**Rejected**: Generic mock framework with `On("method").Return(...)` -- requires reflection or code generation, adds complexity, fights Go's type system.

**Chosen**: Provide reusable types and a template pattern:

### Core Types (in tuitestkit package)
```go
// MockCallRecorder tracks method calls with thread safety.
type MockCallRecorder struct {
    mu    sync.Mutex
    calls []MockCall
}

type MockCall struct {
    Method string
    Args   []any
}

func (r *MockCallRecorder) Record(method string, args ...any)
func (r *MockCallRecorder) CallCount(method string) int
func (r *MockCallRecorder) CallsFor(method string) []MockCall
func (r *MockCallRecorder) AllCalls() []MockCall
func (r *MockCallRecorder) Reset()
```

### Response Helpers
```go
// MockResponse is a pre-configured (data, error) pair.
type MockResponse struct {
    Data  []byte
    Error error
}

// MockResponseMap maps string keys to responses (for parameterized methods).
type MockResponseMap struct {
    mu        sync.RWMutex
    responses map[string]MockResponse
    fallback  MockResponse
}

func NewMockResponseMap() *MockResponseMap
func (m *MockResponseMap) Set(key string, data []byte, err error)
func (m *MockResponseMap) Get(key string) ([]byte, error)
func (m *MockResponseMap) SetFallback(data []byte, err error)
```

### User Creates Their Mock Using These
```go
type MockCLIExecutor struct {
    tuitestkit.MockCallRecorder  // embed for call recording
    TreeJSONResp    tuitestkit.MockResponse
    ShowElementResp tuitestkit.MockResponseMap
    // ... project-specific fields
}

func (m *MockCLIExecutor) TreeJSON() ([]byte, error) {
    m.Record("TreeJSON")
    return m.TreeJSONResp.Data, m.TreeJSONResp.Error
}

func (m *MockCLIExecutor) ShowElementJSON(id string) ([]byte, error) {
    m.Record("ShowElementJSON", id)
    return m.ShowElementResp.Get(id)
}
```

### Test Assertion Helpers
```go
// AssertCalled checks that a method was called at least once.
func AssertCalled(t testing.TB, r *MockCallRecorder, method string)

// AssertCalledN checks exact call count.
func AssertCalledN(t testing.TB, r *MockCallRecorder, method string, n int)

// AssertCalledWith checks a method was called with specific args.
func AssertCalledWith(t testing.TB, r *MockCallRecorder, method string, args ...any)

// AssertNotCalled checks a method was never called.
func AssertNotCalled(t testing.TB, r *MockCallRecorder, method string)
```

---

## Task Breakdown

1. **MockCall + MockCallRecorder** -- core call recording with thread safety
2. **MockResponse + MockResponseMap** -- canned response building blocks
3. **Mock assertion helpers** -- AssertCalled, AssertCalledN, AssertCalledWith, AssertNotCalled
4. **Tests for all the above** -- comprehensive, table-driven
5. **Usage example** -- doc test or example_test.go showing the full pattern

---

## Key Design Notes

- **Embed pattern**: `MockCallRecorder` is designed for embedding, so user's mock struct gets Record/CallCount/CallsFor "for free"
- **No generics needed**: `MockResponse` works with `[]byte` because CLI executors return bytes (per the established pattern)
- **Thread safety**: All recorder and response map operations are mutex-protected (bubbletea runs Cmds as goroutines)
- **Zero external deps**: Just Go stdlib + testing
- **Relationship to STORY-260210-7fupsc (mock-templates)**: This story builds the library types. The templates story provides copy-paste scaffolding showing how to compose these types into a complete mock for a given interface.
