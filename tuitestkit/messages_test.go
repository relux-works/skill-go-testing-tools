package tuitestkit

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// --- Key() tests ---

func TestKey_SpecialKeys(t *testing.T) {
	tests := []struct {
		input    string
		wantType tea.KeyType
	}{
		{"enter", tea.KeyEnter},
		{"tab", tea.KeyTab},
		{"esc", tea.KeyEsc},
		{"escape", tea.KeyEscape},
		{"backspace", tea.KeyBackspace},
		{"up", tea.KeyUp},
		{"down", tea.KeyDown},
		{"left", tea.KeyLeft},
		{"right", tea.KeyRight},
		{"home", tea.KeyHome},
		{"end", tea.KeyEnd},
		{"pgup", tea.KeyPgUp},
		{"pgdown", tea.KeyPgDown},
		{"delete", tea.KeyDelete},
		{"insert", tea.KeyInsert},
		{"shift+tab", tea.KeyShiftTab},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			msg := Key(tt.input)
			if msg.Type != tt.wantType {
				t.Errorf("Key(%q).Type = %v, want %v", tt.input, msg.Type, tt.wantType)
			}
		})
	}
}

func TestKey_Space(t *testing.T) {
	msg := Key("space")
	if msg.Type != tea.KeySpace {
		t.Errorf("Key(\"space\").Type = %v, want KeySpace (%v)", msg.Type, tea.KeySpace)
	}
	if len(msg.Runes) != 1 || msg.Runes[0] != ' ' {
		t.Errorf("Key(\"space\").Runes = %v, want [' ']", msg.Runes)
	}
}

func TestKey_FKeys(t *testing.T) {
	fkeys := map[string]tea.KeyType{
		"f1":  tea.KeyF1,
		"f2":  tea.KeyF2,
		"f3":  tea.KeyF3,
		"f4":  tea.KeyF4,
		"f5":  tea.KeyF5,
		"f6":  tea.KeyF6,
		"f7":  tea.KeyF7,
		"f8":  tea.KeyF8,
		"f9":  tea.KeyF9,
		"f10": tea.KeyF10,
		"f11": tea.KeyF11,
		"f12": tea.KeyF12,
		"f13": tea.KeyF13,
		"f14": tea.KeyF14,
		"f15": tea.KeyF15,
		"f16": tea.KeyF16,
		"f17": tea.KeyF17,
		"f18": tea.KeyF18,
		"f19": tea.KeyF19,
		"f20": tea.KeyF20,
	}

	for name, wantType := range fkeys {
		t.Run(name, func(t *testing.T) {
			msg := Key(name)
			if msg.Type != wantType {
				t.Errorf("Key(%q).Type = %v, want %v", name, msg.Type, wantType)
			}
		})
	}
}

func TestKey_CtrlCombos(t *testing.T) {
	tests := []struct {
		input    string
		wantType tea.KeyType
	}{
		{"ctrl+a", tea.KeyCtrlA},
		{"ctrl+b", tea.KeyCtrlB},
		{"ctrl+c", tea.KeyCtrlC},
		{"ctrl+d", tea.KeyCtrlD},
		{"ctrl+e", tea.KeyCtrlE},
		{"ctrl+f", tea.KeyCtrlF},
		{"ctrl+g", tea.KeyCtrlG},
		{"ctrl+h", tea.KeyCtrlH},
		{"ctrl+i", tea.KeyCtrlI},
		{"ctrl+j", tea.KeyCtrlJ},
		{"ctrl+k", tea.KeyCtrlK},
		{"ctrl+l", tea.KeyCtrlL},
		{"ctrl+m", tea.KeyCtrlM},
		{"ctrl+n", tea.KeyCtrlN},
		{"ctrl+o", tea.KeyCtrlO},
		{"ctrl+p", tea.KeyCtrlP},
		{"ctrl+q", tea.KeyCtrlQ},
		{"ctrl+r", tea.KeyCtrlR},
		{"ctrl+s", tea.KeyCtrlS},
		{"ctrl+t", tea.KeyCtrlT},
		{"ctrl+u", tea.KeyCtrlU},
		{"ctrl+v", tea.KeyCtrlV},
		{"ctrl+w", tea.KeyCtrlW},
		{"ctrl+x", tea.KeyCtrlX},
		{"ctrl+y", tea.KeyCtrlY},
		{"ctrl+z", tea.KeyCtrlZ},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			msg := Key(tt.input)
			if msg.Type != tt.wantType {
				t.Errorf("Key(%q).Type = %v, want %v", tt.input, msg.Type, tt.wantType)
			}
			if msg.Alt {
				t.Errorf("Key(%q).Alt = true, want false", tt.input)
			}
		})
	}
}

func TestKey_CtrlSpecialSymbols(t *testing.T) {
	tests := []struct {
		input    string
		wantType tea.KeyType
	}{
		{"ctrl+@", tea.KeyCtrlAt},
		{"ctrl+[", tea.KeyCtrlOpenBracket},
		{"ctrl+\\", tea.KeyCtrlBackslash},
		{"ctrl+]", tea.KeyCtrlCloseBracket},
		{"ctrl+^", tea.KeyCtrlCaret},
		{"ctrl+_", tea.KeyCtrlUnderscore},
		{"ctrl+?", tea.KeyCtrlQuestionMark},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			msg := Key(tt.input)
			if msg.Type != tt.wantType {
				t.Errorf("Key(%q).Type = %v, want %v", tt.input, msg.Type, tt.wantType)
			}
		})
	}
}

func TestKey_SingleRune(t *testing.T) {
	tests := []struct {
		input    string
		wantRune rune
	}{
		{"a", 'a'},
		{"b", 'b'},
		{"z", 'z'},
		{"1", '1'},
		{"0", '0'},
		{"/", '/'},
		{".", '.'},
		{"-", '-'},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			msg := Key(tt.input)
			if msg.Type != tea.KeyRunes {
				t.Errorf("Key(%q).Type = %v, want KeyRunes", tt.input, msg.Type)
			}
			if len(msg.Runes) != 1 || msg.Runes[0] != tt.wantRune {
				t.Errorf("Key(%q).Runes = %v, want [%c]", tt.input, msg.Runes, tt.wantRune)
			}
			if msg.Alt {
				t.Errorf("Key(%q).Alt = true, want false", tt.input)
			}
		})
	}
}

func TestKey_AltCombos(t *testing.T) {
	t.Run("alt+rune", func(t *testing.T) {
		msg := Key("alt+h")
		if msg.Type != tea.KeyRunes {
			t.Errorf("Key(\"alt+h\").Type = %v, want KeyRunes", msg.Type)
		}
		if len(msg.Runes) != 1 || msg.Runes[0] != 'h' {
			t.Errorf("Key(\"alt+h\").Runes = %v, want ['h']", msg.Runes)
		}
		if !msg.Alt {
			t.Errorf("Key(\"alt+h\").Alt = false, want true")
		}
	})

	t.Run("alt+enter", func(t *testing.T) {
		msg := Key("alt+enter")
		if msg.Type != tea.KeyEnter {
			t.Errorf("Key(\"alt+enter\").Type = %v, want KeyEnter", msg.Type)
		}
		if !msg.Alt {
			t.Errorf("Key(\"alt+enter\").Alt = false, want true")
		}
	})

	t.Run("alt+space", func(t *testing.T) {
		msg := Key("alt+space")
		if msg.Type != tea.KeySpace {
			t.Errorf("Key(\"alt+space\").Type = %v, want KeySpace", msg.Type)
		}
		if !msg.Alt {
			t.Errorf("Key(\"alt+space\").Alt = false, want true")
		}
	})

	t.Run("alt+a", func(t *testing.T) {
		msg := Key("alt+a")
		if msg.Type != tea.KeyRunes {
			t.Errorf("Key(\"alt+a\").Type = %v, want KeyRunes", msg.Type)
		}
		if len(msg.Runes) != 1 || msg.Runes[0] != 'a' {
			t.Errorf("Key(\"alt+a\").Runes = %v, want ['a']", msg.Runes)
		}
		if !msg.Alt {
			t.Errorf("Key(\"alt+a\").Alt = false, want true")
		}
	})

	t.Run("alt+f1", func(t *testing.T) {
		msg := Key("alt+f1")
		if msg.Type != tea.KeyF1 {
			t.Errorf("Key(\"alt+f1\").Type = %v, want KeyF1", msg.Type)
		}
		if !msg.Alt {
			t.Errorf("Key(\"alt+f1\").Alt = false, want true")
		}
	})
}

func TestKey_CaseInsensitive(t *testing.T) {
	tests := []struct {
		input    string
		wantType tea.KeyType
	}{
		{"Enter", tea.KeyEnter},
		{"ENTER", tea.KeyEnter},
		{"Tab", tea.KeyTab},
		{"ESC", tea.KeyEsc},
		{"Ctrl+C", tea.KeyCtrlC},
		{"CTRL+C", tea.KeyCtrlC},
		{"Alt+Enter", tea.KeyEnter},
		{"F1", tea.KeyF1},
		{"F12", tea.KeyF12},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			msg := Key(tt.input)
			if msg.Type != tt.wantType {
				t.Errorf("Key(%q).Type = %v, want %v", tt.input, msg.Type, tt.wantType)
			}
		})
	}
}

func TestKey_StringRoundtrip(t *testing.T) {
	// Verify that Key().String() produces the expected bubbletea string representation
	tests := []struct {
		input   string
		wantStr string
	}{
		{"enter", "enter"},
		{"tab", "tab"},
		{"esc", "esc"},
		{"backspace", "backspace"},
		{"ctrl+c", "ctrl+c"},
		{"ctrl+a", "ctrl+a"},
		{"a", "a"},
		{"z", "z"},
		{"alt+h", "alt+h"},
		{"alt+enter", "alt+enter"},
		{"up", "up"},
		{"down", "down"},
		{"f1", "f1"},
		{"f12", "f12"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			msg := Key(tt.input)
			got := msg.String()
			if got != tt.wantStr {
				t.Errorf("Key(%q).String() = %q, want %q", tt.input, got, tt.wantStr)
			}
		})
	}
}

// --- Keys() tests ---

func TestKeys_Empty(t *testing.T) {
	msgs := Keys()
	if len(msgs) != 0 {
		t.Errorf("Keys() returned %d messages, want 0", len(msgs))
	}
}

func TestKeys_Single(t *testing.T) {
	msgs := Keys("enter")
	if len(msgs) != 1 {
		t.Fatalf("Keys(\"enter\") returned %d messages, want 1", len(msgs))
	}
	km, ok := msgs[0].(tea.KeyMsg)
	if !ok {
		t.Fatalf("Keys(\"enter\")[0] type = %T, want tea.KeyMsg", msgs[0])
	}
	if km.Type != tea.KeyEnter {
		t.Errorf("Keys(\"enter\")[0].Type = %v, want KeyEnter", km.Type)
	}
}

func TestKeys_Multiple(t *testing.T) {
	msgs := Keys("h", "e", "l", "l", "o", "enter")
	if len(msgs) != 6 {
		t.Fatalf("Keys() returned %d messages, want 6", len(msgs))
	}

	// Check the rune keys
	expected := []rune{'h', 'e', 'l', 'l', 'o'}
	for i, r := range expected {
		km, ok := msgs[i].(tea.KeyMsg)
		if !ok {
			t.Fatalf("msgs[%d] type = %T, want tea.KeyMsg", i, msgs[i])
		}
		if km.Type != tea.KeyRunes {
			t.Errorf("msgs[%d].Type = %v, want KeyRunes", i, km.Type)
		}
		if len(km.Runes) != 1 || km.Runes[0] != r {
			t.Errorf("msgs[%d].Runes = %v, want [%c]", i, km.Runes, r)
		}
	}

	// Check the enter key
	km, ok := msgs[5].(tea.KeyMsg)
	if !ok {
		t.Fatalf("msgs[5] type = %T, want tea.KeyMsg", msgs[5])
	}
	if km.Type != tea.KeyEnter {
		t.Errorf("msgs[5].Type = %v, want KeyEnter", km.Type)
	}
}

func TestKeys_MixedTypes(t *testing.T) {
	msgs := Keys("ctrl+a", "x", "alt+enter", "f5")
	if len(msgs) != 4 {
		t.Fatalf("Keys() returned %d messages, want 4", len(msgs))
	}

	assertKeyMsg(t, msgs[0], tea.KeyCtrlA, false)
	assertKeyMsgRune(t, msgs[1], 'x')
	assertKeyMsg(t, msgs[2], tea.KeyEnter, true)
	assertKeyMsg(t, msgs[3], tea.KeyF5, false)
}

// --- WindowSize() tests ---

func TestWindowSize(t *testing.T) {
	msg := WindowSize(80, 24)
	if msg.Width != 80 {
		t.Errorf("WindowSize(80, 24).Width = %d, want 80", msg.Width)
	}
	if msg.Height != 24 {
		t.Errorf("WindowSize(80, 24).Height = %d, want 24", msg.Height)
	}
}

func TestWindowSize_Zero(t *testing.T) {
	msg := WindowSize(0, 0)
	if msg.Width != 0 {
		t.Errorf("WindowSize(0, 0).Width = %d, want 0", msg.Width)
	}
	if msg.Height != 0 {
		t.Errorf("WindowSize(0, 0).Height = %d, want 0", msg.Height)
	}
}

func TestWindowSize_Large(t *testing.T) {
	msg := WindowSize(1920, 1080)
	if msg.Width != 1920 {
		t.Errorf("WindowSize(1920, 1080).Width = %d, want 1920", msg.Width)
	}
	if msg.Height != 1080 {
		t.Errorf("WindowSize(1920, 1080).Height = %d, want 1080", msg.Height)
	}
}

// --- MouseClick() tests ---

func TestMouseClick(t *testing.T) {
	msg := MouseClick(10, 5)
	if msg.X != 10 {
		t.Errorf("MouseClick(10, 5).X = %d, want 10", msg.X)
	}
	if msg.Y != 5 {
		t.Errorf("MouseClick(10, 5).Y = %d, want 5", msg.Y)
	}
	if msg.Button != tea.MouseButtonLeft {
		t.Errorf("MouseClick(10, 5).Button = %v, want MouseButtonLeft", msg.Button)
	}
	if msg.Action != tea.MouseActionPress {
		t.Errorf("MouseClick(10, 5).Action = %v, want MouseActionPress", msg.Action)
	}
	if msg.Type != tea.MouseLeft {
		t.Errorf("MouseClick(10, 5).Type = %v, want MouseLeft", msg.Type)
	}
}

func TestMouseClick_Origin(t *testing.T) {
	msg := MouseClick(0, 0)
	if msg.X != 0 || msg.Y != 0 {
		t.Errorf("MouseClick(0, 0) position = (%d, %d), want (0, 0)", msg.X, msg.Y)
	}
}

// --- MouseClickRight() tests ---

func TestMouseClickRight(t *testing.T) {
	msg := MouseClickRight(15, 20)
	if msg.X != 15 {
		t.Errorf("MouseClickRight(15, 20).X = %d, want 15", msg.X)
	}
	if msg.Y != 20 {
		t.Errorf("MouseClickRight(15, 20).Y = %d, want 20", msg.Y)
	}
	if msg.Button != tea.MouseButtonRight {
		t.Errorf("MouseClickRight(15, 20).Button = %v, want MouseButtonRight", msg.Button)
	}
	if msg.Action != tea.MouseActionPress {
		t.Errorf("MouseClickRight(15, 20).Action = %v, want MouseActionPress", msg.Action)
	}
	if msg.Type != tea.MouseRight {
		t.Errorf("MouseClickRight(15, 20).Type = %v, want MouseRight", msg.Type)
	}
}

// --- MouseScroll() tests ---

func TestMouseScroll_Up(t *testing.T) {
	msg := MouseScroll(ScrollUp)
	if msg.Button != tea.MouseButtonWheelUp {
		t.Errorf("MouseScroll(ScrollUp).Button = %v, want MouseButtonWheelUp", msg.Button)
	}
	if msg.Action != tea.MouseActionPress {
		t.Errorf("MouseScroll(ScrollUp).Action = %v, want MouseActionPress", msg.Action)
	}
	if msg.Type != tea.MouseWheelUp {
		t.Errorf("MouseScroll(ScrollUp).Type = %v, want MouseWheelUp", msg.Type)
	}
}

func TestMouseScroll_Down(t *testing.T) {
	msg := MouseScroll(ScrollDown)
	if msg.Button != tea.MouseButtonWheelDown {
		t.Errorf("MouseScroll(ScrollDown).Button = %v, want MouseButtonWheelDown", msg.Button)
	}
	if msg.Type != tea.MouseWheelDown {
		t.Errorf("MouseScroll(ScrollDown).Type = %v, want MouseWheelDown", msg.Type)
	}
}

func TestMouseScroll_Left(t *testing.T) {
	msg := MouseScroll(ScrollLeft)
	if msg.Button != tea.MouseButtonWheelLeft {
		t.Errorf("MouseScroll(ScrollLeft).Button = %v, want MouseButtonWheelLeft", msg.Button)
	}
	if msg.Type != tea.MouseWheelLeft {
		t.Errorf("MouseScroll(ScrollLeft).Type = %v, want MouseWheelLeft", msg.Type)
	}
}

func TestMouseScroll_Right(t *testing.T) {
	msg := MouseScroll(ScrollRight)
	if msg.Button != tea.MouseButtonWheelRight {
		t.Errorf("MouseScroll(ScrollRight).Button = %v, want MouseButtonWheelRight", msg.Button)
	}
	if msg.Type != tea.MouseWheelRight {
		t.Errorf("MouseScroll(ScrollRight).Type = %v, want MouseWheelRight", msg.Type)
	}
}

func TestMouseScroll_DefaultPosition(t *testing.T) {
	msg := MouseScroll(ScrollUp)
	if msg.X != 0 || msg.Y != 0 {
		t.Errorf("MouseScroll(ScrollUp) position = (%d, %d), want (0, 0)", msg.X, msg.Y)
	}
}

// --- MouseRelease() tests ---

func TestMouseRelease(t *testing.T) {
	msg := MouseRelease(10, 5)
	if msg.X != 10 {
		t.Errorf("MouseRelease(10, 5).X = %d, want 10", msg.X)
	}
	if msg.Y != 5 {
		t.Errorf("MouseRelease(10, 5).Y = %d, want 5", msg.Y)
	}
	if msg.Action != tea.MouseActionRelease {
		t.Errorf("MouseRelease(10, 5).Action = %v, want MouseActionRelease", msg.Action)
	}
	if msg.Button != tea.MouseButtonNone {
		t.Errorf("MouseRelease(10, 5).Button = %v, want MouseButtonNone", msg.Button)
	}
	if msg.Type != tea.MouseRelease {
		t.Errorf("MouseRelease(10, 5).Type = %v, want MouseRelease", msg.Type)
	}
}

// --- Helpers ---

func assertKeyMsg(t *testing.T, msg tea.Msg, wantType tea.KeyType, wantAlt bool) {
	t.Helper()
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		t.Fatalf("msg type = %T, want tea.KeyMsg", msg)
	}
	if km.Type != wantType {
		t.Errorf("msg.Type = %v, want %v", km.Type, wantType)
	}
	if km.Alt != wantAlt {
		t.Errorf("msg.Alt = %v, want %v", km.Alt, wantAlt)
	}
}

func assertKeyMsgRune(t *testing.T, msg tea.Msg, wantRune rune) {
	t.Helper()
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		t.Fatalf("msg type = %T, want tea.KeyMsg", msg)
	}
	if km.Type != tea.KeyRunes {
		t.Errorf("msg.Type = %v, want KeyRunes", km.Type)
	}
	if len(km.Runes) != 1 || km.Runes[0] != wantRune {
		t.Errorf("msg.Runes = %v, want [%c]", km.Runes, wantRune)
	}
}
