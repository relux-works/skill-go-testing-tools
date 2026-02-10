# TASK-260210-3h6o1x: key-builder-function

## Description
Implement Key(k string) tea.KeyMsg builder that parses human-readable key names (matching bubbletea keyNames format) into tea.KeyMsg structs. Build reverse map from keyNames strings to KeyType constants. Handle: single runes, named keys (enter, esc, tab, etc.), ctrl+letter combos, shift+key, ctrl+shift+key, alt+modifier prefix, function keys f1-f20, space key.

## Scope
(define task scope)

## Acceptance Criteria
- Key("h") returns tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}\n- Key("enter") returns tea.KeyMsg{Type: tea.KeyEnter}\n- Key("ctrl+c") returns tea.KeyMsg{Type: tea.KeyCtrlC}\n- Key("alt+h") returns tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}, Alt: true}\n- Key("shift+tab") returns tea.KeyMsg{Type: tea.KeyShiftTab}\n- Key("ctrl+shift+up") returns tea.KeyMsg{Type: tea.KeyCtrlShiftUp}\n- Key("f1") through Key("f20") return correct function key types\n- Key(" ") and Key("space") both return tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}}\n- Key("backspace") returns tea.KeyMsg{Type: tea.KeyBackspace}\n- All ctrl+a through ctrl+z work correctly\n- Reverse map built from bubbletea keyNames constant (not hardcoded duplicates)
