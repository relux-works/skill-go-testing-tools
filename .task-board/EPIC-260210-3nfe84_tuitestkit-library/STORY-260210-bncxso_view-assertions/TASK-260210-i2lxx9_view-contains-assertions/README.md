# TASK-260210-i2lxx9: view-contains-assertions

## Description
ViewContains(t testing.TB, model tea.Model, text string) and ViewNotContains(t testing.TB, model tea.Model, text string) assertion helpers. Call model.View(), strip ANSI via StripANSI, then check strings.Contains. Use t.Helper() for clean error reporting. Use t.Errorf (non-fatal) by default. Tests: mock tea.Model with styled View() output, verify positive/negative assertions, verify failure messages.

## Scope
(define task scope)

## Acceptance Criteria
(define acceptance criteria)
