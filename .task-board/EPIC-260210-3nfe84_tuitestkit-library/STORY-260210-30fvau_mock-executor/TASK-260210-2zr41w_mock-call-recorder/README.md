# TASK-260210-2zr41w: mock-call-recorder

## Description
Implement MockCall struct and MockCallRecorder with thread-safe call recording. MockCall holds Method string and Args []any. MockCallRecorder provides Record(), CallCount(), CallsFor(), AllCalls(), Reset() methods. Uses sync.Mutex for goroutine safety (bubbletea Cmds run as goroutines). Designed for struct embedding so user mocks get recording for free.

## Scope
(define task scope)

## Acceptance Criteria
(define acceptance criteria)
