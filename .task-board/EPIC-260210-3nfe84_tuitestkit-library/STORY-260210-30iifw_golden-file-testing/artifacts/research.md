# Golden File Testing Research

## Go Community Conventions

### testdata/ vs .snapshots/

Go has a strong convention of using `testdata/` for test fixtures and golden files. The `go build` tool automatically ignores `testdata/` directories. However, SPEC.md specifically says `.snapshots/` convention.

**Decision:** Use `testdata/.snapshots/` — combines Go convention with the SPEC requirement. The `.snapshots/` subdirectory inside `testdata/` keeps golden files organized separately from other fixtures.

Actually, for a **library** that users consume, the snapshot directory should be configurable but default to `testdata/snapshots/` (no dot — Go convention, visible in Finder). The `.snapshots/` from SPEC is the user-facing convention; the library should use a reasonable default.

**Final decision:** Default directory `testdata/snapshots/`, configurable via option.

### Update Mechanism: flag.Bool vs Environment Variable

Two mainstream approaches:

1. **`flag.Bool("-update", ...)`** — used by goldie, gotest.tools. Requires `TestMain` + `flag.Parse()`. Problem: since tuitestkit is a library, we can't define flags in the library itself — that would conflict with the consuming project's flags.

2. **Environment variable (`UPDATE_SNAPSHOTS=1`)** — used by BTBurke/snapshot, go-snapshot-testing. No flag registration needed. Works perfectly for libraries.

**Decision:** Use `os.Getenv("UPDATE_SNAPSHOTS")` checked at call time. Any truthy value (`1`, `true`, `yes`) triggers update mode. This is library-friendly — no flag registration, no TestMain requirement, no conflicts.

Additionally, expose a package-level `var UpdateSnapshots bool` that users can set manually if they prefer `flag.Bool` in their own TestMain.

### Golden File Format

- **Stripped (no ANSI)** — better for readability, git diffs, manual review
- **Raw (with ANSI)** — exact matching, catches style regressions

**Decision:** `SnapshotView()` stores ANSI-stripped content (primary use case). Provide `SnapshotViewRaw()` variant for exact ANSI matching when needed.

File extension: `.golden` (community standard).

### Diff Output on Mismatch

goldie provides three engines: Classic, Colored, Simple.
Most Go golden file libs show line-by-line diff with +/- markers.

**Decision:** Implement a simple unified-diff-style output:
- Show line numbers
- Mark added lines with `+`
- Mark removed lines with `-`
- Show context lines around changes
- Keep it in the library (no external diff dependency)

### Auto-creation of Directories

On first snapshot update, `.snapshots/` directory should be auto-created via `os.MkdirAll`. Standard practice.

## API Design

```go
// Package-level control
var UpdateSnapshots bool

func init() {
    env := os.Getenv("UPDATE_SNAPSHOTS")
    UpdateSnapshots = env == "1" || env == "true" || env == "yes"
}

// Core API
func SnapshotView(t testing.TB, model tea.Model, name string)
func SnapshotViewRaw(t testing.TB, model tea.Model, name string)

// String-based variants (when you already have the view string)
func SnapshotStr(t testing.TB, view string, name string)
func SnapshotStrRaw(t testing.TB, view string, name string)

// Options for customization
type SnapshotOption func(*snapshotConfig)
func WithSnapshotDir(dir string) SnapshotOption
func SnapshotViewWith(t testing.TB, model tea.Model, name string, opts ...SnapshotOption)
```

### Snapshot Directory Resolution

The directory must be relative to the test file. Use `runtime.Caller(1)` to find the caller's file path, then resolve `testdata/snapshots/` relative to it.

### File Naming

`testdata/snapshots/<name>.golden` — name is provided by the user, sanitized to be filesystem-safe.

## Sources

- [goldie](https://github.com/sebdah/goldie) — testdata/ default, -update flag, diff engines
- [xorcare/golden](https://github.com/xorcare/golden) — minimal golden file testing
- [gotest.tools/v3/golden](https://pkg.go.dev/gotest.tools/v3/golden) — -update flag
- [Go golden files blog](https://ieftimov.com/posts/testing-in-go-golden-files/) — flag.Bool pattern, testdata/ convention
- [Golden Files — Why you should use them](https://jarifibrahim.github.io/blog/golden-files-why-you-should-use-them/) — flag.Bool pattern
- [BTBurke/snapshot](https://github.com/BTBurke/snapshot) — UPDATE_SNAPSHOTS env var
- [go-golden](https://pkg.go.dev/github.com/jimeh/go-golden) — GOLDEN_UPDATE env var
