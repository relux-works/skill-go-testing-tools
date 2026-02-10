package tuitestkit

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// styledModel is a minimal tea.Model for testing view assertions.
type styledModel struct{ content string }

func (m styledModel) Init() tea.Cmd                           { return nil }
func (m styledModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m styledModel) View() string                            { return m.content }

// ansiWrap applies a lipgloss style to produce real ANSI output.
func ansiWrap(s string) string {
	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("9"))
	return style.Render(s)
}

// --- StripANSI tests ---

func TestStripANSI_Plain(t *testing.T) {
	input := "hello world"
	got := StripANSI(input)
	if got != "hello world" {
		t.Errorf("StripANSI plain: got %q, want %q", got, "hello world")
	}
}

func TestStripANSI_WithEscapes(t *testing.T) {
	input := "\033[1;31mhello\033[0m"
	got := StripANSI(input)
	if got != "hello" {
		t.Errorf("StripANSI escaped: got %q, want %q", got, "hello")
	}
}

func TestStripANSI_Lipgloss(t *testing.T) {
	styled := ansiWrap("bold red")
	got := StripANSI(styled)
	if got != "bold red" {
		t.Errorf("StripANSI lipgloss: got %q, want %q", got, "bold red")
	}
}

func TestStripANSI_Empty(t *testing.T) {
	got := StripANSI("")
	if got != "" {
		t.Errorf("StripANSI empty: got %q, want %q", got, "")
	}
}

// --- ContainsStr / NotContainsStr tests ---

func TestContainsStr_PlainMatch(t *testing.T) {
	ContainsStr(t, "hello world", "world")
}

func TestContainsStr_ANSIMatch(t *testing.T) {
	view := ansiWrap("hello") + " world"
	ContainsStr(t, view, "hello world")
}

func TestContainsStr_Fail(t *testing.T) {
	fake := &testing.T{}
	ContainsStr(fake, "hello world", "missing")
	if !fake.Failed() {
		t.Error("ContainsStr should have failed for missing text")
	}
}

func TestNotContainsStr_PlainAbsent(t *testing.T) {
	NotContainsStr(t, "hello world", "missing")
}

func TestNotContainsStr_ANSIAbsent(t *testing.T) {
	view := ansiWrap("hello") + " world"
	NotContainsStr(t, view, "missing")
}

func TestNotContainsStr_Fail(t *testing.T) {
	fake := &testing.T{}
	NotContainsStr(fake, "hello world", "world")
	if !fake.Failed() {
		t.Error("NotContainsStr should have failed for present text")
	}
}

// --- ViewContains / ViewNotContains tests ---

func TestViewContains_PlainModel(t *testing.T) {
	m := styledModel{content: "Status: OK"}
	ViewContains(t, m, "OK")
}

func TestViewContains_StyledModel(t *testing.T) {
	m := styledModel{content: ansiWrap("Status") + ": OK"}
	ViewContains(t, m, "Status: OK")
}

func TestViewContains_Fail(t *testing.T) {
	fake := &testing.T{}
	m := styledModel{content: "Status: OK"}
	ViewContains(fake, m, "ERROR")
	if !fake.Failed() {
		t.Error("ViewContains should have failed")
	}
}

func TestViewNotContains_PlainModel(t *testing.T) {
	m := styledModel{content: "Status: OK"}
	ViewNotContains(t, m, "ERROR")
}

func TestViewNotContains_Fail(t *testing.T) {
	fake := &testing.T{}
	m := styledModel{content: "Status: OK"}
	ViewNotContains(fake, m, "OK")
	if !fake.Failed() {
		t.Error("ViewNotContains should have failed")
	}
}

// --- LinesFromStr / ViewLines tests ---

func TestLinesFromStr_Plain(t *testing.T) {
	lines := LinesFromStr("line one\nline two\nline three")
	if len(lines) != 3 {
		t.Fatalf("LinesFromStr: got %d lines, want 3", len(lines))
	}
	if lines[0] != "line one" {
		t.Errorf("LinesFromStr[0] = %q, want %q", lines[0], "line one")
	}
	if lines[2] != "line three" {
		t.Errorf("LinesFromStr[2] = %q, want %q", lines[2], "line three")
	}
}

func TestLinesFromStr_StripsANSI(t *testing.T) {
	view := ansiWrap("styled") + "\nplain"
	lines := LinesFromStr(view)
	if len(lines) != 2 {
		t.Fatalf("LinesFromStr ANSI: got %d lines, want 2", len(lines))
	}
	if lines[0] != "styled" {
		t.Errorf("LinesFromStr ANSI [0] = %q, want %q", lines[0], "styled")
	}
	if lines[1] != "plain" {
		t.Errorf("LinesFromStr ANSI [1] = %q, want %q", lines[1], "plain")
	}
}

func TestLinesFromStr_TrimsTrailingEmpty(t *testing.T) {
	lines := LinesFromStr("hello\n\n\n")
	if len(lines) != 1 {
		t.Fatalf("LinesFromStr trailing: got %d lines, want 1", len(lines))
	}
	if lines[0] != "hello" {
		t.Errorf("LinesFromStr trailing[0] = %q, want %q", lines[0], "hello")
	}
}

func TestLinesFromStr_AllEmpty(t *testing.T) {
	lines := LinesFromStr("\n\n\n")
	if len(lines) != 0 {
		t.Errorf("LinesFromStr all empty: got %d lines, want 0", len(lines))
	}
}

func TestLinesFromStr_SingleLine(t *testing.T) {
	lines := LinesFromStr("single")
	if len(lines) != 1 {
		t.Fatalf("LinesFromStr single: got %d lines, want 1", len(lines))
	}
	if lines[0] != "single" {
		t.Errorf("LinesFromStr single[0] = %q, want %q", lines[0], "single")
	}
}

func TestViewLines_Model(t *testing.T) {
	m := styledModel{content: ansiWrap("Header") + "\nBody\nFooter\n"}
	lines := ViewLines(m)
	if len(lines) != 3 {
		t.Fatalf("ViewLines: got %d lines, want 3", len(lines))
	}
	if lines[0] != "Header" {
		t.Errorf("ViewLines[0] = %q, want %q", lines[0], "Header")
	}
}

// --- ViewLineContains / ViewLineEquals tests ---

func TestViewLineContains_Match(t *testing.T) {
	m := styledModel{content: "first\nsecond contains target\nthird"}
	ViewLineContains(t, m, 1, "target")
}

func TestViewLineContains_Styled(t *testing.T) {
	m := styledModel{content: "first\n" + ansiWrap("styled target") + "\nthird"}
	ViewLineContains(t, m, 1, "styled target")
}

func TestViewLineContains_Fail_NotFound(t *testing.T) {
	fake := &testing.T{}
	m := styledModel{content: "first\nsecond\nthird"}
	ViewLineContains(fake, m, 1, "missing")
	if !fake.Failed() {
		t.Error("ViewLineContains should have failed for missing text")
	}
}

func TestViewLineContains_Fail_OutOfBounds(t *testing.T) {
	fake := &testing.T{}
	m := styledModel{content: "only one line"}
	ViewLineContains(fake, m, 5, "text")
	if !fake.Failed() {
		t.Error("ViewLineContains should have failed for out-of-bounds index")
	}
}

func TestViewLineContains_Fail_NegativeIndex(t *testing.T) {
	fake := &testing.T{}
	m := styledModel{content: "line"}
	ViewLineContains(fake, m, -1, "line")
	if !fake.Failed() {
		t.Error("ViewLineContains should have failed for negative index")
	}
}

func TestViewLineEquals_Match(t *testing.T) {
	m := styledModel{content: "exact match\nsecond"}
	ViewLineEquals(t, m, 0, "exact match")
}

func TestViewLineEquals_Styled(t *testing.T) {
	m := styledModel{content: ansiWrap("exact match") + "\nsecond"}
	ViewLineEquals(t, m, 0, "exact match")
}

func TestViewLineEquals_Fail_Mismatch(t *testing.T) {
	fake := &testing.T{}
	m := styledModel{content: "actual text\nsecond"}
	ViewLineEquals(fake, m, 0, "different text")
	if !fake.Failed() {
		t.Error("ViewLineEquals should have failed for mismatched text")
	}
}

func TestViewLineEquals_Fail_OutOfBounds(t *testing.T) {
	fake := &testing.T{}
	m := styledModel{content: "one line"}
	ViewLineEquals(fake, m, 10, "text")
	if !fake.Failed() {
		t.Error("ViewLineEquals should have failed for out-of-bounds index")
	}
}

// --- MatchesRegexStr / ViewMatchesRegex tests ---

func TestMatchesRegexStr_PlainMatch(t *testing.T) {
	MatchesRegexStr(t, "Status: 200 OK", `Status: \d+ OK`)
}

func TestMatchesRegexStr_ANSIMatch(t *testing.T) {
	view := ansiWrap("Status") + ": 200 OK"
	MatchesRegexStr(t, view, `Status: \d+ OK`)
}

func TestMatchesRegexStr_Fail_NoMatch(t *testing.T) {
	fake := &testing.T{}
	MatchesRegexStr(fake, "hello world", `^\d+$`)
	if !fake.Failed() {
		t.Error("MatchesRegexStr should have failed for non-matching pattern")
	}
}

func TestMatchesRegexStr_Fail_BadRegex(t *testing.T) {
	fake := &testing.T{}
	MatchesRegexStr(fake, "hello", `[invalid`)
	if !fake.Failed() {
		t.Error("MatchesRegexStr should have failed for invalid regex")
	}
}

func TestViewMatchesRegex_Model(t *testing.T) {
	m := styledModel{content: "Items: 42 total"}
	ViewMatchesRegex(t, m, `Items: \d+ total`)
}

func TestViewMatchesRegex_StyledModel(t *testing.T) {
	m := styledModel{content: ansiWrap("Items") + ": 42 total"}
	ViewMatchesRegex(t, m, `Items: \d+ total`)
}

func TestViewMatchesRegex_Fail(t *testing.T) {
	fake := &testing.T{}
	m := styledModel{content: "no numbers here"}
	ViewMatchesRegex(fake, m, `\d+`)
	if !fake.Failed() {
		t.Error("ViewMatchesRegex should have failed")
	}
}
