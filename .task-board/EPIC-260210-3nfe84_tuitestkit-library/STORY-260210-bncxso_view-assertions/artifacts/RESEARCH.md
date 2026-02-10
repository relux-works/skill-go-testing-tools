# Research: View Assertions

## Current State of View Testing in board-tui

### settings_test.go (the only file that tests View())

The test at `tools/board-tui/settings_test.go:144-169` calls `m.View()` and checks for substrings using hand-rolled helpers:

```go
func containsString(s, substr string) bool {
    return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstring(s, substr))
}

func containsSubstring(s, substr string) bool {
    for i := 0; i <= len(s)-len(substr); i++ {
        if s[i:i+len(substr)] == substr {
            return true
        }
    }
    return false
}
```

This works because lipgloss injects ANSI escape codes **around** styled text, not inside words. So `"Settings"` still appears as a contiguous byte sequence in the output, even though it's wrapped with color codes like `\x1b[38;5;212mSettings\x1b[0m`. The raw `strings.Contains` or byte scanning approach works for simple cases, but breaks when:

- Text is split across multiple ANSI segments (e.g., per-character styling)
- Testing for multi-word phrases where styling changes mid-phrase
- Assertions need to be on specific lines (line splitting is affected by ANSI)

### Other test files

No other `*_test.go` files call `.View()`. They test state/reducers (pure functions), which is the correct Elm-architecture pattern. View assertions are the exception, not the norm.

## ANSI Stripping Options in Go

### Option A: `github.com/charmbracelet/x/ansi` (RECOMMENDED)

**`ansi.Strip(s string) string`** — removes all ANSI escape sequences from a string.

Pros:
- Already a transitive dependency of bubbletea/lipgloss (confirmed in board-tui go.mod: `github.com/charmbracelet/x/ansi v0.11.5`)
- Maintained by the same team that builds lipgloss/bubbletea — guaranteed to handle all sequences lipgloss generates
- Uses a proper state-machine parser (not regex), which handles edge cases better
- Also provides `StringWidth()` for display-width-aware operations

Cons:
- None significant. This is the canonical choice.

### Option B: `github.com/acarl005/stripansi`

Uses regex: `\x1b\[[0-9;]*m` (approximate). Simple and proven but:
- Regex-based — may miss some exotic sequences (OSC, DCS, etc.)
- Extra dependency when we already have charmbracelet/x/ansi transitively
- Less maintained (last commit years ago)

### Option C: `github.com/muesli/reflow/ansi`

Part of muesli/reflow (also a transitive dependency). Has ANSI-aware text operations but no direct `Strip()` function. Not suitable as a standalone stripping utility.

### Option D: Roll our own regex

Standard regex: `\x1b\[[0-9;]*[a-zA-Z]` or `\x1b\[[\d;]*m`. Simple to implement but:
- Doesn't handle all escape sequence types (OSC 8 hyperlinks, etc.)
- Reinventing what charmbracelet/x/ansi already does
- Maintenance burden

**Decision: Use `charmbracelet/x/ansi.Strip()`** — zero new dependencies, battle-tested, canonical.

## What lipgloss Actually Outputs

lipgloss uses `termenv` under the hood for color profiles. Output format depends on terminal capabilities:

- **TrueColor**: `\x1b[38;2;R;G;Bm` (24-bit RGB)
- **256 Color**: `\x1b[38;5;Nm` (256 palette)
- **16 Color**: `\x1b[31m` etc. (basic ANSI)
- **Ascii**: No escape codes at all

Plus bold (`\x1b[1m`), italic (`\x1b[3m`), underline (`\x1b[4m`), reset (`\x1b[0m`), etc.

In tests, the profile depends on the terminal running `go test`. The Charm team recommends `lipgloss.SetColorProfile(termenv.Ascii)` for golden file tests. But for our view assertions, we take the opposite approach: **let lipgloss render with full styling, then strip ANSI for text comparison**. This tests the real rendering path.

## API Design

### Core: StripANSI utility

```go
// StripANSI removes all ANSI escape sequences from s.
// Delegates to charmbracelet/x/ansi.Strip().
func StripANSI(s string) string
```

### Assertion helpers (accept testing.TB for both testing.T and testing.B)

```go
// ViewContains asserts model.View() contains text after stripping ANSI.
func ViewContains(t testing.TB, model tea.Model, text string)

// ViewNotContains asserts model.View() does NOT contain text after stripping ANSI.
func ViewNotContains(t testing.TB, model tea.Model, text string)

// ViewLines returns model.View() split into lines, each line ANSI-stripped.
func ViewLines(model tea.Model) []string

// ViewLineContains asserts that a specific line (0-indexed) contains text.
func ViewLineContains(t testing.TB, model tea.Model, lineIdx int, text string)

// ViewLineEquals asserts that a specific line (0-indexed) exactly equals text.
func ViewLineEquals(t testing.TB, model tea.Model, lineIdx int, text string)

// ViewMatchesRegex asserts model.View() (ANSI-stripped) matches the regex pattern.
func ViewMatchesRegex(t testing.TB, model tea.Model, pattern string)
```

### Design decisions

1. **Accept `tea.Model` not `string`**: Helpers call `.View()` themselves. This reads better and matches the pattern `ViewContains(t, model, "Settings")` from the spec.
2. **Use `testing.TB`**: Works with both `*testing.T` and `*testing.B`, more flexible than just `*testing.T`.
3. **`t.Helper()`**: All assertion functions call `t.Helper()` so failures report at the caller's line.
4. **`t.Errorf` not `t.Fatalf`**: Non-fatal by default — let all assertions run. Users can wrap in `if` for early termination.
5. **No raw string overloads initially**: The spec says `ViewContains(model, text)` — we can add `StringContains(t, s, text)` for raw string testing if needed, but the primary API is model-based.
6. **StripANSI is public**: Useful as a standalone utility for custom assertions.

## Task Breakdown

1. **StripANSI function + tests** — core utility wrapping `charmbracelet/x/ansi.Strip()`, with tests that verify stripping of various ANSI patterns (colors, bold, reset, 256-color, truecolor, combined).
2. **ViewContains/ViewNotContains + tests** — the two most common assertion helpers.
3. **ViewLines + ViewLineContains/ViewLineEquals + tests** — line-based assertions.
4. **ViewMatchesRegex + tests** — regex-based assertion for more flexible matching.

Dependencies: Task 1 is the foundation. Tasks 2-4 depend on Task 1. Tasks 2-4 are independent of each other.
