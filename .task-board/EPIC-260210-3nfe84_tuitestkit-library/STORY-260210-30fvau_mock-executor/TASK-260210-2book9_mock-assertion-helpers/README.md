# TASK-260210-2book9: mock-assertion-helpers

## Description
Implement test assertion functions: AssertCalled(t, recorder, method), AssertCalledN(t, recorder, method, n), AssertCalledWith(t, recorder, method, args...), AssertNotCalled(t, recorder, method). All accept testing.TB for compatibility with both testing.T and testing.B. Use t.Helper() for clean error reporting. Compare args with reflect.DeepEqual for AssertCalledWith.

## Scope
(define task scope)

## Acceptance Criteria
(define acceptance criteria)
