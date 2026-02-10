# TASK-260210-2mxre1: keys-batch-builder

## Description
Implement Keys(keys ...string) []tea.Msg helper that takes variadic key name strings and returns a slice of tea.Msg. Delegates to Key() for each string. Returns []tea.Msg (not []tea.KeyMsg) for direct use with test harness Send() functions.

## Scope
(define task scope)

## Acceptance Criteria
- Keys("h", "j", "enter") returns []tea.Msg with 3 elements\n- Each element matches what Key() would return for that string\n- Return type is []tea.Msg (not []tea.KeyMsg)\n- Empty call Keys() returns empty slice (not nil)\n- Single key Keys("enter") returns single-element slice
