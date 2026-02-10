# TASK-260210-15lvy5: snapshot-tests

## Description
Comprehensive tests in snapshot_test.go covering: (1) SnapshotStr creates golden file when UpdateSnapshots=true, (2) SnapshotStr passes when content matches golden file, (3) SnapshotStr fails with diff when content mismatches, (4) SnapshotView strips ANSI before comparison, (5) SnapshotViewRaw preserves ANSI in golden file, (6) Missing golden file produces helpful error message, (7) Directory auto-creation on first write, (8) Name sanitization for filesystem safety, (9) WithSnapshotDir option overrides default directory. Use t.TempDir() for isolated test directories. Test diff output format correctness.

## Scope
(define task scope)

## Acceptance Criteria
(define acceptance criteria)
