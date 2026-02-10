# TASK-260210-1of26b: implement-send-generic

## Description
Implement Send[M tea.Model](model M, msgs ...tea.Msg) M in tuitestkit/harness.go. Iterates through msgs, calls model.Update(msg) for each, type-asserts result back to M, returns final model. Panics on type assertion failure. Handles empty msgs (returns model as-is).

## Scope
(define task scope)

## Acceptance Criteria
(define acceptance criteria)
