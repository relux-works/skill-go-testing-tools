# STORY-260210-3j76j0: model-test-harness

## Description
Send(), SendAndCollect(), ExecCmds() - send message sequences to bubbletea models, collect results

## Scope
Implement Send[M], SendAndCollect[M], and ExecCmds functions in tuitestkit/harness.go. Send feeds a sequence of tea.Msg through a model's Update(), preserving concrete type via generics. SendAndCollect does the same but also collects returned tea.Cmds. ExecCmds executes tea.Cmd functions synchronously, collecting resulting messages, with recursive unpacking of tea.BatchMsg.

## Acceptance Criteria
- Send[M tea.Model](model M, msgs ...tea.Msg) M correctly applies messages sequentially and returns final model state with preserved concrete type
- SendAndCollect[M tea.Model](model M, msgs ...tea.Msg) (M, []tea.Cmd) collects non-nil Cmds returned by Update
- ExecCmds(cmds ...tea.Cmd) []tea.Msg executes Cmds and collects non-nil messages
- ExecCmds recursively unpacks tea.BatchMsg (Batch results)
- All functions handle nil/empty inputs gracefully without panics
- All functions have comprehensive test coverage in harness_test.go
- No new dependencies beyond existing bubbletea
