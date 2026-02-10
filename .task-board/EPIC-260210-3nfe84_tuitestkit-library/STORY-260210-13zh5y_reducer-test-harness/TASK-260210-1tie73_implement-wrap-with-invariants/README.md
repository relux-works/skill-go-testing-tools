# TASK-260210-1tie73: implement-wrap-with-invariants

## Description
Implement WrapWithInvariants[S,A](reduce func(S,A)S, checker *InvariantChecker[S]) func(S,A)S standalone function. Returns a new reducer that calls the original then runs checker.CheckAll. Must be a standalone function (not method) because Go methods cannot have additional type parameters. Place in invariant.go alongside InvariantChecker.

## Scope
(define task scope)

## Acceptance Criteria
(define acceptance criteria)
