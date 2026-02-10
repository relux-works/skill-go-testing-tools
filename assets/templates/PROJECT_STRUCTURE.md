# Recommended Test Directory Structure

```
myapp/
├── internal/
│   ├── cli/
│   │   ├── executor.go               # Executor interface (see executor_interface.go.tmpl)
│   │   └── executor_real.go           # Production implementation (os/exec)
│   │
│   ├── state/
│   │   ├── state.go                   # State type definition
│   │   ├── actions.go                 # Action types
│   │   ├── reducer.go                 # Pure reducer: func(State, Action) State
│   │   └── reducer_test.go            # Reducer tests (see reducer_test.go.tmpl)
│   │                                  #   - Table-driven with ReducerTest
│   │                                  #   - Multi-step with ReducerSequence
│   │                                  #   - Invariant checking with WrapWithInvariants
│   │
│   ├── ui/
│   │   ├── screens/
│   │   │   ├── board/
│   │   │   │   ├── board.go           # bubbletea Model (Init, Update, View)
│   │   │   │   ├── board_test.go      # Component tests (see component_test.go.tmpl)
│   │   │   │   │                      #   - Send/SendAndCollect for input simulation
│   │   │   │   │                      #   - ViewContains/ViewLines for output assertions
│   │   │   │   │                      #   - ExecCmds for command pipeline testing
│   │   │   │   ├── board_snapshot_test.go  # Snapshot tests (see snapshot_test.go.tmpl)
│   │   │   │   │                      #   - SnapshotView for golden file comparison
│   │   │   │   │                      #   - UPDATE_SNAPSHOTS=1 to regenerate
│   │   │   │   ├── mock_executor_test.go   # Mock for this screen's tests (see mock_executor.go.tmpl)
│   │   │   │   │                      #   - MockCallRecorder for call tracking
│   │   │   │   │                      #   - MockResponseMap for canned responses
│   │   │   │   └── testdata/
│   │   │   │       └── snapshots/     # Golden files (*.golden)
│   │   │   │                          #   - Committed to git
│   │   │   │                          #   - Regenerated with UPDATE_SNAPSHOTS=1
│   │   │   │
│   │   │   ├── detail/
│   │   │   │   ├── detail.go
│   │   │   │   ├── detail_test.go
│   │   │   │   └── testdata/
│   │   │   │       └── snapshots/
│   │   │   │
│   │   │   └── settings/
│   │   │       ├── settings.go
│   │   │       ├── settings_test.go
│   │   │       └── testdata/
│   │   │           └── snapshots/
│   │   │
│   │   └── components/
│   │       ├── filter/
│   │       │   ├── filter.go
│   │       │   └── filter_test.go
│   │       └── dialog/
│   │           ├── dialog.go
│   │           └── dialog_test.go
│   │
│   └── domain/                        # Domain types (shared across packages)
│       └── types.go
│
├── go.mod
├── go.sum
└── main.go
```

## Key Placement Rules

**Reducer tests** (`reducer_test.go`) go next to the reducer implementation.
Pure functions, no mock needed. Use `ReducerTest` for single-action cases,
`ReducerSequence` for multi-step flows, `WrapWithInvariants` for safety nets.

**Component tests** (`*_test.go`) go next to the component.
Each screen/component gets its own test file. Use `Send`/`SendAndCollect`
to simulate user input, `ViewContains`/`ViewLines` to assert output.

**Snapshot tests** (`*_snapshot_test.go`) go next to the component.
Separate file from component tests to keep concerns clear.
Golden files live in `testdata/snapshots/` within the same package directory.

**Mock files** (`mock_*_test.go`) go in the test package that uses them.
Build mocks by composing `MockCallRecorder` + `MockResponseMap`.
One mock per interface, shared across tests in the same package.

## Running Tests

```bash
# Run all tests
go test ./...

# Run with verbose output
go test ./... -v

# Run only reducer tests
go test ./internal/state/ -v -run TestReducer

# Run only snapshot tests
go test ./internal/ui/screens/board/ -v -run TestSnapshot

# Update golden files
UPDATE_SNAPSHOTS=1 go test ./...

# Update golden files for a single package
UPDATE_SNAPSHOTS=1 go test ./internal/ui/screens/board/ -run TestSnapshot
```
