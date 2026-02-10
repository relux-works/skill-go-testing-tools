# TASK-260210-3hnv2a: mouse-and-window-tests

## Description
Test suite for MouseClick(), MouseClickRight(), MouseScroll(), MouseRelease(), and WindowSize() builders. Verify all struct fields are correctly populated including deprecated Type field. Test all ScrollDir values. Test coordinate handling. Test that builders produce identical output to manual construction patterns found in board-tui tests.

## Scope
(define task scope)

## Acceptance Criteria
- MouseClick(5, 10) fields: X=5, Y=10, Action=Press, Button=Left, Type=MouseLeft\n- MouseClickRight fields validated\n- MouseScroll(ScrollUp/Down/Left/Right) all produce correct Button and deprecated Type\n- MouseRelease fields: Action=Release, Button=None, Type=MouseRelease\n- WindowSize(80, 24) fields: Width=80, Height=24\n- WindowSize(0, 0) does not panic\n- Equivalence test: MouseMsg{Button: tea.MouseButtonWheelRight} == MouseScroll(ScrollRight) (matching board-tui arkanoid_test.go pattern)\n- Test coverage > 90% for mouse.go and windowsize.go source files
