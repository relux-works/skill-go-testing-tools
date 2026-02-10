# TASK-260210-7k7c0j: key-builder-tests

## Description
Comprehensive test suite for Key() and Keys() builders. Table-driven tests covering: all named keys from keyNames map, all ctrl+letter combinations (ctrl+a through ctrl+z), alt modifier for both runes and named keys, shift/ctrl+shift combinations, function keys f1-f20, space handling (both " " and "space"), edge cases (empty string, unknown key name). Verify roundtrip: Key(name).String() == name for all supported keys.

## Scope
(define task scope)

## Acceptance Criteria
- Table-driven tests for all named keys (enter, esc, tab, backspace, up, down, left, right, home, end, pgup, pgdown, delete, insert, space)\n- Tests for all ctrl+a through ctrl+z\n- Tests for alt+rune and alt+named-key\n- Tests for shift+tab, shift+up/down/left/right, shift+home/end\n- Tests for ctrl+up/down/left/right, ctrl+home/end, ctrl+pgup/pgdown\n- Tests for ctrl+shift combos\n- Tests for f1 through f20\n- Roundtrip test: Key(x).String() == x for all valid inputs\n- Keys() returns correct length and element types\n- Test coverage > 90% for key.go source file
