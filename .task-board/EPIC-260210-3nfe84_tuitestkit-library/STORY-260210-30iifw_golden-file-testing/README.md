# STORY-260210-30iifw: golden-file-testing

## Description
SnapshotView(), UpdateSnapshots flag, .snapshots/ convention, ANSI-aware diff on mismatch

## Scope
Add golden file / snapshot testing to tuitestkit. Implement SnapshotView() and SnapshotStr() functions that save/compare bubbletea model view output against golden files stored in testdata/snapshots/. Support UPDATE_SNAPSHOTS env var and package-level UpdateSnapshots bool for regenerating snapshots. Provide ANSI-stripped (default) and raw variants. Show ANSI-aware unified diff on mismatch. Auto-create snapshot directories. No external dependencies beyond what tuitestkit already has.

## Acceptance Criteria
1. SnapshotView(t, model, name) compares model.View() (ANSI-stripped) against testdata/snapshots/<name>.golden
2. SnapshotViewRaw(t, model, name) compares model.View() with ANSI preserved
3. SnapshotStr(t, view, name) and SnapshotStrRaw(t, view, name) string-based variants
4. UPDATE_SNAPSHOTS=1 env var triggers golden file creation/update
5. Package-level var UpdateSnapshots bool for manual control
6. On mismatch: unified diff output showing expected vs got with line numbers
7. Snapshot directory auto-created on first update
8. Golden files use .golden extension
9. Snapshot dir defaults to testdata/snapshots/ relative to test file (via runtime.Caller)
10. All functions accept testing.TB for compatibility with testing.T and testing.B
11. Tests cover: create snapshot, compare match, compare mismatch with diff, update mode, missing snapshot, stripped vs raw, directory auto-creation
