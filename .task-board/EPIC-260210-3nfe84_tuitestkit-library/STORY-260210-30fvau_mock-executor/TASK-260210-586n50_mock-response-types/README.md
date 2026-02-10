# TASK-260210-586n50: mock-response-types

## Description
Implement MockResponse struct (Data []byte, Error error) and MockResponseMap for parameterized methods. MockResponseMap maps string keys to MockResponse with Set(), Get(), SetFallback() methods. Thread-safe with sync.RWMutex. Get() returns fallback when key not found (fallback defaults to error). These are composable building blocks for user-defined mocks.

## Scope
(define task scope)

## Acceptance Criteria
(define acceptance criteria)
