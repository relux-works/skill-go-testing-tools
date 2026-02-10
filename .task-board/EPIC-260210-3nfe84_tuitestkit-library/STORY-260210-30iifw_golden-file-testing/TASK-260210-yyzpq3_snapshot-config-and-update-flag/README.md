# TASK-260210-yyzpq3: snapshot-config-and-update-flag

## Description
Create snapshot.go with core configuration: package-level var UpdateSnapshots bool, init() that reads UPDATE_SNAPSHOTS env var (accepts 1/true/yes), snapshotConfig struct with Dir field, SnapshotOption type, WithSnapshotDir() option, default dir resolution using runtime.Caller to find testdata/snapshots/ relative to test file. Also implement name sanitization (replace non-filesystem-safe chars). This is the foundation all other snapshot functions depend on.

## Scope
(define task scope)

## Acceptance Criteria
(define acceptance criteria)
