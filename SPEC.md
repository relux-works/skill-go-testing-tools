# SPEC: Go Testing Tools

## Vision

**Closed-loop agent development cycle for Go TUI applications.** The agent writes TUI code, writes tests, runs them, validates results, and fixes issues — all autonomously, without human intervention.

Built on the insight that **Elm architecture (state → action → reducer → view)** makes TUI applications uniquely testable through pure function testing.

---

## Problem

Go TUI apps built with bubbletea are hard to test because:

1. **No established patterns** — bubbletea community has minimal test guidance
2. **exec.Command coupling** — screens that shell out to CLI tools can't be unit tested
3. **View rendering** — lipgloss/ANSI output is hard to assert on
4. **Agent can't validate** — without tests, the agent has no feedback loop to verify its own work
5. **Boilerplate** — each project reinvents mock patterns, test helpers, message simulation

---

## Solution

A skill + Go library that gives agents everything they need for the closed-loop:

```
Agent writes code → Agent writes tests → Agent runs tests → Tests pass/fail → Agent fixes → Repeat
```

### Three Deliverables

1. **SKILL.md** — agent workflow guide (like ios-testing-tools but for Go TUI)
2. **Go library (`tuitestkit`)** — reusable test utilities as a Go module
3. **Reference docs & templates** — patterns, architecture guides, copy-paste scaffolding

---

## Requirements

### R1: Skill Documentation (SKILL.md)

Complete agent guide covering:

- **Elm architecture testing philosophy** — why pure reducers are the foundation
- **Testing pyramid for TUI** — unit (reducers) → component (Update/View) → integration (multi-component) → snapshot (visual)
- **Closed-loop workflow** — write → test → run → validate → fix cycle
- **When to use which test type** — decision tree for agents
- **Prerequisites** — Go 1.21+, bubbletea, lipgloss

### R2: Go Library (`tuitestkit`)

Go module at `github.com/relux-works/skill-go-testing-tools/tuitestkit` (or similar):

#### R2.1: Message Builders
- `Key(k string) tea.KeyMsg` — build key messages (`Key("enter")`, `Key("ctrl+c")`)
- `Keys(keys ...string) []tea.Msg` — batch key sequences
- `WindowSize(w, h int) tea.WindowSizeMsg`
- `MouseClick(x, y int) tea.MouseMsg`
- `MouseScroll(dir ScrollDir) tea.MouseMsg`

#### R2.2: Model Test Harness
- `Send(model, msgs...) model` — send sequence of messages, return final state
- `SendAndCollect(model, msgs...) (model, []tea.Cmd)` — also collect commands
- `ExecCmds(cmds ...tea.Cmd) []tea.Msg` — execute commands, collect messages

#### R2.3: Reducer Test Harness
- `ReducerTest[S, A]` — table-driven test struct for pure reducers
- Assert helpers: `StateEquals`, `FieldEquals`, `FieldChanged`
- Invariant checker: register invariants, auto-check after every Reduce()

#### R2.4: Mock Executor
- `MockExecutor` — configurable mock for CLI executor pattern
- Canned responses by command
- Call recording (what was called, with what args)
- Error simulation

#### R2.5: View Assertions
- `ViewContains(model, text)` — check rendered output contains text
- `ViewNotContains(model, text)` — negative check
- `ViewLines(model) []string` — split view into lines for line-by-line assertion
- ANSI stripping for clean text comparison

#### R2.6: Golden File / Snapshot Testing
- `SnapshotView(t, model, name)` — save/compare view output against golden file
- `UpdateSnapshots` flag for regenerating (`go test -update-snapshots`)
- `.snapshots/` directory convention
- ANSI-aware diff on mismatch

### R3: Reference Documentation

#### R3.1: Architecture Patterns (`references/elm-architecture.md`)
- Why Elm makes TUI testable
- State → Action → Reducer → View cycle
- Composition: nested reducers, action wrapping
- Dependency injection for side effects

#### R3.2: Mock Patterns (`references/mock-patterns.md`)
- CLI executor interface extraction
- When to mock vs when to test pure functions
- Test doubles: mocks vs stubs vs fakes
- Callback capture pattern

#### R3.3: Testing Pyramid (`references/testing-pyramid.md`)
- Level 1: Pure reducer tests (fastest, most numerous)
- Level 2: Component model tests (Update + View)
- Level 3: Integration tests (multi-component interaction)
- Level 4: Snapshot tests (visual regression)
- Level 5: Behavioral tests (full workflow simulation)

### R4: Asset Templates

#### R4.1: Test File Templates
- `_test.go` template for reducer testing
- `_test.go` template for component model testing
- `_test.go` template for snapshot testing

#### R4.2: Mock Templates
- CLI executor interface + mock implementation
- Configurable response map

#### R4.3: Project Scaffolding
- Recommended directory structure for test code
- `testdata/` conventions for fixtures

### R5: Scripts & Tooling

- `scripts/setup.sh` — install skill, set up symlinks
- `scripts/check-tools.sh` — verify Go version, bubbletea installed
- `scripts/deinit.sh` — remove symlinks

### R6: Closed-Loop Agent Workflow

The core cycle that the skill teaches agents:

```
1. UNDERSTAND: Read screen/component code
2. IDENTIFY: What behaviors need tests?
3. EXTRACT: If exec.Command present → extract to interface (use mock template)
4. WRITE: Tests using tuitestkit helpers
5. RUN: `go test ./... -v`
6. VALIDATE: All pass? Check coverage with `go test -cover`
7. FIX: If failures → read error → fix code or test → goto 5
8. SNAPSHOT: Optionally capture view golden files
9. DONE: Tests green, coverage acceptable
```

The agent must be able to execute this entire cycle without human input.

---

## Non-Goals

- Visual screenshot testing (terminal doesn't have screenshots like iOS)
- Browser-based testing
- Performance benchmarking framework (use standard Go benchmarks)
- CI/CD pipeline templates (project-specific)

---

## Success Criteria

- Agent can write comprehensive tests for any bubbletea TUI screen using the skill
- Tests run in < 1 second (no I/O, no network)
- Library has zero non-test dependencies beyond bubbletea/lipgloss
- Skill documentation is self-contained — agent doesn't need to search for patterns
- Closed-loop works: agent detects failures and fixes them without human help
