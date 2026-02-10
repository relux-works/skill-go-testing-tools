# TASK-260210-2bzqcq: implement-exec-cmds

## Description
Implement ExecCmds(cmds ...tea.Cmd) []tea.Msg in tuitestkit/harness.go. Executes each Cmd by calling it synchronously. Skips nil Cmds. Collects non-nil resulting messages. Recursively unpacks tea.BatchMsg: when a Cmd returns a BatchMsg ([]Cmd), executes all contained Cmds and collects their results recursively.

## Scope
(define task scope)

## Acceptance Criteria
(define acceptance criteria)
