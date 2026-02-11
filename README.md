# Go Testing Tools

Agent skill + Go library for testing bubbletea TUI applications. Enables a **closed-loop agent development cycle**: write code, write tests, run, validate, fix — all autonomously.

## What's Inside

### `tuitestkit/` — Go Library

Test utilities as a Go module (`github.com/ivalx1s/skill-go-testing-tools/tuitestkit`):

| File | What |
|------|------|
| `messages.go` | Message builders: `Key()`, `Keys()`, `WindowSize()`, `MouseClick()`, `MouseScroll()` |
| `harness.go` | Model harness: `Send[M]()`, `SendAndCollect[M]()`, `ExecCmds()` |
| `reducer.go` | Reducer harness: `RunReducerTests()`, `RunReducerSequences()`, `WrapWithInvariants()` |
| `mock.go` | Mock building blocks: `MockCallRecorder`, `MockResponseMap`, assertion helpers |
| `view.go` | View assertions: `ViewContains()`, `ViewLines()`, `ViewMatchesRegex()`, `StripANSI()` |
| `snapshot.go` | Golden file testing: `SnapshotView()`, `SnapshotStr()`, unified diff engine |

155 tests, zero external dependencies beyond bubbletea.

### `references/` — Architecture Docs

- `elm-architecture.md` — Why Elm makes TUI testable, State/Action/Reducer/View cycle
- `mock-patterns.md` — CLI executor extraction, test doubles, composition patterns
- `testing-pyramid.md` — 5-level pyramid: reducer → component → integration → snapshot → behavioral

### `assets/templates/` — Copy-Paste Templates

- `reducer_test.go.tmpl` — Pure reducer test template
- `component_test.go.tmpl` — Bubbletea screen test template
- `snapshot_test.go.tmpl` — Golden file test template
- `executor_interface.go.tmpl` — CLI executor interface pattern
- `mock_executor.go.tmpl` — Mock executor with call recording
- `PROJECT_STRUCTURE.md` — Recommended directory layout

### `scripts/` — Setup & Tooling

- `setup.sh` — Install skill (symlinks to `~/.agents/skills/`, `~/.claude/skills/`, `~/.codex/skills/`)
- `deinit.sh` — Remove symlinks
- `check-tools.sh` — Verify Go version and bubbletea dependency

## Setup

```bash
./scripts/setup.sh
```

## Quick Start

```go
import (
    "testing"
    kit "github.com/ivalx1s/skill-go-testing-tools/tuitestkit"
)

func TestMyScreen(t *testing.T) {
    m := NewMyModel()

    // Send keys, get final model
    m = kit.Send(m, kit.Key("down"), kit.Key("enter"))

    // Assert on view output
    kit.ViewContains(t, m, "Expected text")
}
```

## License

MIT
