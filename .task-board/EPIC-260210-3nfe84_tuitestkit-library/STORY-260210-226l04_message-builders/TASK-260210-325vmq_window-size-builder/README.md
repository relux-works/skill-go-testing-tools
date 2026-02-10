# TASK-260210-325vmq: window-size-builder

## Description
Implement WindowSize(w, h int) tea.WindowSizeMsg builder. Trivial wrapper that returns tea.WindowSizeMsg{Width: w, Height: h}.

## Scope
(define task scope)

## Acceptance Criteria
- WindowSize(80, 24) returns tea.WindowSizeMsg{Width: 80, Height: 24}\n- WindowSize(0, 0) works without panic\n- Return type is tea.WindowSizeMsg (usable directly as tea.Msg)
