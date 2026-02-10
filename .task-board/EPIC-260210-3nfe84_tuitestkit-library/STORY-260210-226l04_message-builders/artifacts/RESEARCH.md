# Research: Message Builders for tuitestkit

## Bubbletea Version

`github.com/charmbracelet/bubbletea v1.3.10`

---

## Key Types and Structures

### tea.KeyMsg / tea.Key

```go
type KeyMsg Key

type Key struct {
    Type  KeyType
    Runes []rune
    Alt   bool
    Paste bool
}
```

**KeyType** is an `int`. Two categories:

#### Control key types (positive values, from C0/C1 codes)
These map to `ctrl+<letter>` strings, plus aliases:
- `KeyNull` (0) = `keyNUL` = `ctrl+@`
- `KeyBreak` (3) = `keyETX` = `ctrl+c`
- `KeyEnter` (13) = `keyCR` = `"enter"`
- `KeyBackspace` (127) = `keyDEL` = `"backspace"`
- `KeyTab` (9) = `keyHT` = `"tab"`
- `KeyEsc` / `KeyEscape` (27) = `keyESC` = `"esc"`
- `KeyCtrlA` through `KeyCtrlZ` (1-26, mapping to ctrl+a through ctrl+z)
- Special aliases: `KeyCtrlOpenBracket` = esc, `KeyCtrlBackslash`, `KeyCtrlCloseBracket`, `KeyCtrlCaret`, `KeyCtrlUnderscore`, `KeyCtrlQuestionMark`

#### Other key types (negative iota values)
- `KeyRunes` (-1) = `"runes"` -- for regular character input
- `KeyUp`, `KeyDown`, `KeyRight`, `KeyLeft`
- `KeyShiftTab` = `"shift+tab"`
- `KeyHome`, `KeyEnd`
- `KeyPgUp`, `KeyPgDown`
- `KeyCtrlPgUp`, `KeyCtrlPgDown`
- `KeyDelete`, `KeyInsert`
- `KeySpace` = `" "`
- `KeyCtrlUp/Down/Right/Left`
- `KeyCtrlHome`, `KeyCtrlEnd`
- `KeyShiftUp/Down/Right/Left`
- `KeyShiftHome`, `KeyShiftEnd`
- `KeyCtrlShiftUp/Down/Left/Right`
- `KeyCtrlShiftHome`, `KeyCtrlShiftEnd`
- `KeyF1` through `KeyF20`

#### String representations (from `keyNames` map)
Used by `Key.String()` for comparison. Exact string values:
- `"ctrl+@"`, `"ctrl+a"` .. `"ctrl+z"`, `"ctrl+\\"`, `"ctrl+]"`, `"ctrl+^"`, `"ctrl+_"`
- `"tab"` (ctrl+i), `"enter"` (ctrl+m), `"esc"`, `"backspace"`
- `"up"`, `"down"`, `"right"`, `"left"`, `" "` (space)
- `"shift+tab"`, `"home"`, `"end"`, `"ctrl+home"`, `"ctrl+end"`, `"shift+home"`, `"shift+end"`, `"ctrl+shift+home"`, `"ctrl+shift+end"`
- `"pgup"`, `"pgdown"`, `"ctrl+pgup"`, `"ctrl+pgdown"`
- `"delete"`, `"insert"`
- `"ctrl+up"`, `"ctrl+down"`, `"ctrl+right"`, `"ctrl+left"`
- `"shift+up"`, `"shift+down"`, `"shift+right"`, `"shift+left"`
- `"ctrl+shift+up"`, `"ctrl+shift+down"`, `"ctrl+shift+left"`, `"ctrl+shift+right"`
- `"f1"` through `"f20"`

Alt modifier prepends `"alt+"` via `Key.String()`.

For `KeyRunes`, the string is the runes themselves (e.g. `"h"`, `"j"`, etc.), or with alt: `"alt+h"`.

---

### tea.WindowSizeMsg

```go
type WindowSizeMsg struct {
    Width  int
    Height int
}
```

Simple struct, trivial builder.

---

### tea.MouseMsg / tea.MouseEvent

```go
type MouseMsg MouseEvent

type MouseEvent struct {
    X      int
    Y      int
    Shift  bool
    Alt    bool
    Ctrl   bool
    Action MouseAction
    Button MouseButton
    Type   MouseEventType  // Deprecated
}
```

#### MouseAction
```go
const (
    MouseActionPress   MouseAction = iota  // 0
    MouseActionRelease                      // 1
    MouseActionMotion                       // 2
)
```

#### MouseButton
```go
const (
    MouseButtonNone       MouseButton = iota  // 0
    MouseButtonLeft                            // 1
    MouseButtonMiddle                          // 2
    MouseButtonRight                           // 3
    MouseButtonWheelUp                         // 4
    MouseButtonWheelDown                       // 5
    MouseButtonWheelLeft                       // 6
    MouseButtonWheelRight                      // 7
    MouseButtonBackward                        // 8
    MouseButtonForward                         // 9
    MouseButton10                              // 10
    MouseButton11                              // 11
)
```

#### MouseEventType (deprecated, but still populated for backward compat)
```go
const (
    MouseUnknown MouseEventType = iota
    MouseLeft
    MouseRight
    MouseMiddle
    MouseRelease
    MouseWheelUp
    MouseWheelDown
    MouseWheelLeft
    MouseWheelRight
    MouseBackward
    MouseForward
    MouseMotion
)
```

---

### tea.FocusMsg / tea.BlurMsg

```go
type FocusMsg struct{}
type BlurMsg struct{}
```

Simple zero-value structs. Could add builders but they are trivially constructible.

---

### Other Msg types (internal/less common)

- `QuitMsg struct{}` -- sent by `Quit()`
- `SuspendMsg struct{}` -- sent by `Suspend()`
- `ResumeMsg struct{}`
- `InterruptMsg struct{}`
- `BatchMsg []Cmd`
- Various internal messages (clearScreen, enterAltScreen, etc.) -- unexported, not for user testing

---

## Patterns Found in Real TUI Test Code

### board-tui/settings_test.go patterns
```go
tea.KeyMsg{Type: tea.KeyTab}
tea.KeyMsg{Type: tea.KeyShiftTab}
tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'l'}}
tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
tea.KeyMsg{Type: tea.KeyDown}
tea.KeyMsg{Type: tea.KeyEnter}
tea.KeyMsg{Type: tea.KeyEsc}
```

### board-tui/arkanoid_test.go patterns
```go
tea.KeyMsg{Type: tea.KeyRight}
tea.KeyMsg{Type: tea.KeyLeft}
tea.MouseMsg{Button: tea.MouseButtonWheelRight}
tea.MouseMsg{Button: tea.MouseButtonWheelLeft}
```

### Common pain points
1. Rune keys require `tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}}` -- 3 fields for a single character
2. No way to express `"ctrl+c"` as a string to get `tea.KeyMsg{Type: tea.KeyCtrlC}` without knowing the constant name
3. Mouse events require setting both `Button` and `Action` fields, plus deprecated `Type` for backward compat

---

## Proposed API Design

### Key(k string) tea.KeyMsg

Accepts the **same string format that `Key.String()` returns**. This is the canonical representation.

```go
// Single runes
Key("h")       // → tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}}
Key(" ")       // → tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}}

// Named keys
Key("enter")   // → tea.KeyMsg{Type: tea.KeyEnter}
Key("esc")     // → tea.KeyMsg{Type: tea.KeyEsc}
Key("tab")     // → tea.KeyMsg{Type: tea.KeyTab}
Key("backspace") // → tea.KeyMsg{Type: tea.KeyBackspace}
Key("up")      // → tea.KeyMsg{Type: tea.KeyUp}
Key("down")    // → tea.KeyMsg{Type: tea.KeyDown}
Key("left")    // → tea.KeyMsg{Type: tea.KeyLeft}
Key("right")   // → tea.KeyMsg{Type: tea.KeyRight}
Key("pgup")    // → tea.KeyMsg{Type: tea.KeyPgUp}
Key("pgdown")  // → tea.KeyMsg{Type: tea.KeyPgDown}
Key("home")    // → tea.KeyMsg{Type: tea.KeyHome}
Key("end")     // → tea.KeyMsg{Type: tea.KeyEnd}
Key("delete")  // → tea.KeyMsg{Type: tea.KeyDelete}
Key("insert")  // → tea.KeyMsg{Type: tea.KeyInsert}
Key("space")   // alternative for " "

// Ctrl combinations
Key("ctrl+c")  // → tea.KeyMsg{Type: tea.KeyCtrlC}
Key("ctrl+a")  // → tea.KeyMsg{Type: tea.KeyCtrlA}
Key("ctrl+z")  // → tea.KeyMsg{Type: tea.KeyCtrlZ}

// Shift combinations
Key("shift+tab")   // → tea.KeyMsg{Type: tea.KeyShiftTab}
Key("shift+up")    // → tea.KeyMsg{Type: tea.KeyShiftUp}
Key("shift+down")  // → tea.KeyMsg{Type: tea.KeyShiftDown}

// Ctrl+Shift combinations
Key("ctrl+shift+up")    // → tea.KeyMsg{Type: tea.KeyCtrlShiftUp}

// Alt modifier
Key("alt+h")   // → tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'h'}, Alt: true}
Key("alt+enter") // → tea.KeyMsg{Type: tea.KeyEnter, Alt: true}

// Function keys
Key("f1")      // → tea.KeyMsg{Type: tea.KeyF1}
Key("f12")     // → tea.KeyMsg{Type: tea.KeyF12}
```

**Implementation approach**: Build a reverse map from `keyNames` string values to `KeyType` constants. Parse the string: strip `"alt+"` prefix if present (set `Alt: true`), look up remaining string in reverse map. If single rune and not in map, create `KeyRunes` msg.

### Keys(keys ...string) []tea.Msg

Simple batch helper:

```go
Keys("h", "j", "enter")
// → []tea.Msg{Key("h"), Key("j"), Key("enter")}
```

Returns `[]tea.Msg` (not `[]tea.KeyMsg`) for direct use with test harness.

### WindowSize(w, h int) tea.WindowSizeMsg

Trivial:

```go
WindowSize(80, 24) // → tea.WindowSizeMsg{Width: 80, Height: 24}
```

### Mouse builders

```go
type ScrollDir int
const (
    ScrollUp    ScrollDir = iota
    ScrollDown
    ScrollLeft
    ScrollRight
)

MouseClick(x, y int) tea.MouseMsg
// → tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionPress, Button: tea.MouseButtonLeft, Type: tea.MouseLeft}

MouseClickRight(x, y int) tea.MouseMsg
// → right-click at position

MouseScroll(dir ScrollDir) tea.MouseMsg
// → tea.MouseMsg with appropriate WheelUp/Down/Left/Right button

MouseRelease(x, y int) tea.MouseMsg
// → tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionRelease, Button: tea.MouseButtonNone, Type: tea.MouseRelease}
```

**Note**: The deprecated `Type` field should be populated for backward compatibility, since existing code may check either `Button`+`Action` or `Type`.

---

## Task Breakdown

1. **key-builder** -- `Key()` function with string-to-KeyMsg parsing and reverse keyNames map
2. **keys-builder** -- `Keys()` batch helper
3. **window-size-builder** -- `WindowSize()` builder
4. **mouse-builders** -- `MouseClick()`, `MouseClickRight()`, `MouseScroll()`, `MouseRelease()`, `ScrollDir` type
5. **key-builder-tests** -- comprehensive tests for Key() covering all key types, alt modifier, edge cases
6. **mouse-builder-tests** -- tests for all mouse builders
7. **integration-validation** -- verify builders produce identical output to manual construction patterns from board-tui tests
