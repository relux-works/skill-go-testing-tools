# TASK-260210-2auypu: view-regex-assertion

## Description
ViewMatchesRegex(t testing.TB, model tea.Model, pattern string) — asserts model.View() (ANSI-stripped) matches the given regex pattern. Uses regexp.MatchString under the hood. Reports the pattern and cleaned view on failure. Must handle: invalid regex (report compilation error via t.Errorf, not panic), multiline matching. Tests with various regex patterns against styled output.

## Scope
(define task scope)

## Acceptance Criteria
(define acceptance criteria)
