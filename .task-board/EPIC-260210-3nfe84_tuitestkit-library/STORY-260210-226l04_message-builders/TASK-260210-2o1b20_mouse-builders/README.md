# TASK-260210-2o1b20: mouse-builders

## Description
Implement mouse message builders: MouseClick(x, y int) tea.MouseMsg for left-click, MouseClickRight(x, y int) tea.MouseMsg for right-click, MouseScroll(dir ScrollDir) tea.MouseMsg for scroll wheel events, MouseRelease(x, y int) tea.MouseMsg for button release. Define ScrollDir type with ScrollUp, ScrollDown, ScrollLeft, ScrollRight constants. All builders must populate both the new Button/Action fields AND the deprecated Type field for backward compatibility.

## Scope
(define task scope)

## Acceptance Criteria
- MouseClick(5, 10) returns tea.MouseMsg{X: 5, Y: 10, Action: tea.MouseActionPress, Button: tea.MouseButtonLeft, Type: tea.MouseLeft}\n- MouseClickRight(5, 10) returns right-click with Button: tea.MouseButtonRight, Type: tea.MouseRight\n- MouseScroll(ScrollUp) returns tea.MouseMsg{Button: tea.MouseButtonWheelUp, Type: tea.MouseWheelUp}\n- MouseScroll(ScrollDown/Left/Right) returns corresponding wheel buttons\n- MouseRelease(5, 10) returns tea.MouseMsg{X: 5, Y: 10, Action: tea.MouseActionRelease, Button: tea.MouseButtonNone, Type: tea.MouseRelease}\n- ScrollDir type is exported with four directional constants\n- Deprecated Type field is populated for backward compat in all builders
