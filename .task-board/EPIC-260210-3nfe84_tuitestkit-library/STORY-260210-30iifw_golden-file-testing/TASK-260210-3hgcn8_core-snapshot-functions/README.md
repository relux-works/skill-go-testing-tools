# TASK-260210-3hgcn8: core-snapshot-functions

## Description
Implement the core snapshot read/write/compare logic in snapshot.go: (1) snapshotPath(name, cfg) — resolves full path to golden file. (2) readSnapshot(path) — reads golden file, returns content and bool for exists. (3) writeSnapshot(path, content) — writes golden file with os.MkdirAll for auto-creating directories. (4) compareSnapshot(t, name, actual, strip bool, opts) — the main engine: if UpdateSnapshots, write and return; if golden file missing, fail with helpful message; if mismatch, fail with diff output. Public API: SnapshotView(t, model, name), SnapshotViewRaw(t, model, name), SnapshotStr(t, view, name), SnapshotStrRaw(t, view, name). All accept testing.TB.

## Scope
(define task scope)

## Acceptance Criteria
(define acceptance criteria)
