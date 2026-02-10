# TASK-260210-3btfbq: view-line-assertions

## Description
ViewLines(model tea.Model) []string — splits model.View() into lines after ANSI stripping. ViewLineContains(t testing.TB, model tea.Model, lineIdx int, text string) — asserts specific line contains text. ViewLineEquals(t testing.TB, model tea.Model, lineIdx int, text string) — asserts specific line exactly equals text. Must handle: out-of-bounds lineIdx (report error, not panic), empty lines, trailing newlines. Tests with multi-line styled output.

## Scope
(define task scope)

## Acceptance Criteria
(define acceptance criteria)
