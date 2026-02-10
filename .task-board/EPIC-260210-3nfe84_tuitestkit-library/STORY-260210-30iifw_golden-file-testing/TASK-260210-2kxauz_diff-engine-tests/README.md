# TASK-260210-2kxauz: diff-engine-tests

## Description
Tests in snapshot_diff_test.go for the unified diff engine: (1) Identical strings produce empty diff, (2) Single line addition shows + marker, (3) Single line removal shows - marker, (4) Single line change shows both - and + markers, (5) Multi-line diff with context lines, (6) Empty expected vs non-empty actual, (7) Non-empty expected vs empty actual, (8) Both empty produces empty diff. Verify line numbers in output. Test edge cases: trailing newlines, very long lines.

## Scope
(define task scope)

## Acceptance Criteria
(define acceptance criteria)
