# Research: Model Test Harness (R2.2)

## Bubbletea Interfaces & Types

### tea.Model
```go
type Model interface {
    Init() Cmd
    Update(Msg) (Model, Cmd)
    View() string
}
```

### tea.Cmd and tea.Msg
```go
type Cmd func() Msg
type Msg interface{}
```

### tea.Batch
```go
func Batch(cmds ...Cmd) Cmd  // returns a Cmd that produces BatchMsg
type BatchMsg []Cmd           // exported, can type-assert on it
```

### tea.Sequence
```go
func Sequence(cmds ...Cmd) Cmd  // returns a Cmd that produces sequenceMsg (unexported)
```
`sequenceMsg` is unexported -- we cannot type-assert on it in `ExecCmds`. This means `ExecCmds` can recursively unpack `BatchMsg` but NOT `sequenceMsg`. This is fine -- `Sequence` is for ordered execution within the runtime, and in tests we'd just use separate calls anyway.

## API Design

### Send[M tea.Model](model M, msgs ...tea.Msg) M

- Generic function preserving concrete model type
- Calls `model.Update(msg)` for each message sequentially
- Type-asserts the returned `tea.Model` back to `M`
- Discards returned `tea.Cmd`s (use `SendAndCollect` if you need them)
- Does NOT call `Init()` -- caller handles that separately if needed
- Panics if type assertion fails (indicates a broken Update implementation)

### SendAndCollect[M tea.Model](model M, msgs ...tea.Msg) (M, []tea.Cmd)

- Same as `Send` but collects all returned `tea.Cmd`s into a slice
- nil Cmds are excluded from the collected slice (keeps it clean)
- Returned slice order matches message processing order

### ExecCmds(cmds ...tea.Cmd) []tea.Msg

- Executes each `Cmd` by calling it (synchronously)
- Skips nil Cmds gracefully
- Collects non-nil `tea.Msg` results
- Recursively handles `tea.BatchMsg`: when a Cmd returns a `BatchMsg`, execute all contained Cmds recursively
- Does NOT handle `sequenceMsg` (unexported) -- just returns it as a regular message
- Single-level recursion is probably sufficient, but full recursion handles nested Batch

## Existing Patterns in the Codebase

- Package: `tuitestkit` (no internal packages)
- Style: exported top-level functions, generic where type preservation matters
- Testing: table-driven tests with helper domains (counterState, listState)
- Dependencies: only `github.com/charmbracelet/bubbletea` and `github.com/charmbracelet/x/ansi`
- No `testing.T` parameter in non-assertion functions (Send/ExecCmds are pure helpers, not assertions)

## File Placement

New file: `tuitestkit/harness.go` -- follows the existing pattern of one file per concern (messages.go, reducer.go, view.go, mock.go).

Test file: `tuitestkit/harness_test.go` -- matching test file.

## Edge Cases

1. Empty msgs slice in Send/SendAndCollect -- return model as-is
2. Update returns a different concrete type -- panic (this indicates a bubbletea anti-pattern)
3. Nil Cmd from Update -- skip in collected slice
4. Nil Cmd in ExecCmds input -- skip
5. Cmd returns nil Msg -- skip in collected messages
6. BatchMsg containing nil Cmds -- skip
7. Nested BatchMsg (Batch returning BatchMsg) -- handle recursively
