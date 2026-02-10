package tuitestkit

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// --- test helpers ---

// stubModel is a minimal bubbletea model for snapshot tests.
type stubModel struct {
	view string
}

func (m stubModel) Init() tea.Cmd                           { return nil }
func (m stubModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) { return m, nil }
func (m stubModel) View() string                            { return m.view }

// withSnapshotDir sets snapshotBaseDir for the duration of the test and
// restores the original value afterwards. Also saves/restores UpdateSnapshots.
func withSnapshotDir(t *testing.T, dir string) {
	t.Helper()
	origBase := snapshotBaseDir
	origUpdate := UpdateSnapshots
	snapshotBaseDir = dir
	t.Cleanup(func() {
		snapshotBaseDir = origBase
		UpdateSnapshots = origUpdate
	})
}

// fakeT intercepts Helper/Errorf/Fatalf so we can inspect failure behavior
// without aborting the real test. Fatalf panics with a sentinel so the caller
// can recover and inspect the error message.
type fakeT struct {
	failed  bool
	fataled bool
	lastErr string
}

// fatalSentinel is the panic value used by fakeT.Fatalf.
type fatalSentinel struct{}

func (f *fakeT) Helper() {}

func (f *fakeT) Errorf(format string, args ...any) {
	f.failed = true
	f.lastErr = fmt.Sprintf(format, args...)
}

func (f *fakeT) Fatalf(format string, args ...any) {
	f.fataled = true
	f.failed = true
	f.lastErr = fmt.Sprintf(format, args...)
	panic(fatalSentinel{})
}

// runSnapshot calls snapshot() on fakeT, recovering from Fatalf panics.
func runSnapshot(ft *fakeT, content, name string, callerSkip int) {
	defer func() {
		if r := recover(); r != nil {
			if _, ok := r.(fatalSentinel); !ok {
				panic(r) // re-panic if not our sentinel
			}
		}
	}()
	snapshot(ft, content, name, callerSkip)
}

// --- unifiedDiff tests (TASK-260210-2kxauz) ---

func TestUnifiedDiff_Identical(t *testing.T) {
	diff := unifiedDiff("hello\nworld", "hello\nworld")
	if !strings.Contains(diff, "--- expected") {
		t.Fatal("diff should contain header")
	}
	// No + or - lines (only space-prefixed context lines).
	for _, line := range strings.Split(diff, "\n") {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			t.Errorf("unexpected addition line in identical diff: %q", line)
		}
		if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			t.Errorf("unexpected removal line in identical diff: %q", line)
		}
	}
}

func TestUnifiedDiff_Addition(t *testing.T) {
	diff := unifiedDiff("hello", "hello\nworld")
	// Should show "world" as an added line.
	found := false
	for _, l := range strings.Split(diff, "\n") {
		if strings.HasPrefix(l, "+") && !strings.HasPrefix(l, "+++") && strings.Contains(l, "world") {
			found = true
		}
	}
	if !found {
		t.Errorf("diff should show addition of 'world':\n%s", diff)
	}
}

func TestUnifiedDiff_Removal(t *testing.T) {
	diff := unifiedDiff("hello\nworld", "hello")
	found := false
	for _, l := range strings.Split(diff, "\n") {
		if strings.HasPrefix(l, "-") && !strings.HasPrefix(l, "---") && strings.Contains(l, "world") {
			found = true
		}
	}
	if !found {
		t.Errorf("diff should show removal of 'world':\n%s", diff)
	}
}

func TestUnifiedDiff_Change(t *testing.T) {
	diff := unifiedDiff("hello\nworld", "hello\nearth")
	if !strings.Contains(diff, "world") || !strings.Contains(diff, "earth") {
		t.Errorf("diff should show both old and new lines:\n%s", diff)
	}
	// "world" should be a removal, "earth" an addition.
	hasRemoval := false
	hasAddition := false
	for _, l := range strings.Split(diff, "\n") {
		if strings.HasPrefix(l, "-") && !strings.HasPrefix(l, "---") && strings.Contains(l, "world") {
			hasRemoval = true
		}
		if strings.HasPrefix(l, "+") && !strings.HasPrefix(l, "+++") && strings.Contains(l, "earth") {
			hasAddition = true
		}
	}
	if !hasRemoval {
		t.Error("diff should show 'world' as removed")
	}
	if !hasAddition {
		t.Error("diff should show 'earth' as added")
	}
}

func TestUnifiedDiff_MultiLine(t *testing.T) {
	expected := "line1\nline2\nline3\nline4"
	actual := "line1\nchanged2\nline3\nnew4\nline5"
	diff := unifiedDiff(expected, actual)

	if !strings.Contains(diff, "line2") {
		t.Error("diff should mention removed 'line2'")
	}
	if !strings.Contains(diff, "changed2") {
		t.Error("diff should mention added 'changed2'")
	}
}

func TestUnifiedDiff_Empty_Both(t *testing.T) {
	diff := unifiedDiff("", "")
	// Only headers and the single empty-line context.
	for _, l := range strings.Split(diff, "\n") {
		if strings.HasPrefix(l, "+") && !strings.HasPrefix(l, "+++") {
			t.Errorf("empty-vs-empty should have no additions: %q", l)
		}
		if strings.HasPrefix(l, "-") && !strings.HasPrefix(l, "---") {
			t.Errorf("empty-vs-empty should have no removals: %q", l)
		}
	}
}

func TestUnifiedDiff_Empty_ToNonEmpty(t *testing.T) {
	diff := unifiedDiff("", "hello")
	if !strings.Contains(diff, "hello") {
		t.Error("should show 'hello' as addition")
	}
}

func TestUnifiedDiff_NonEmpty_ToEmpty(t *testing.T) {
	diff := unifiedDiff("hello", "")
	found := false
	for _, l := range strings.Split(diff, "\n") {
		if strings.HasPrefix(l, "-") && !strings.HasPrefix(l, "---") && strings.Contains(l, "hello") {
			found = true
		}
	}
	if !found {
		t.Error("should show 'hello' as removal")
	}
}

func TestUnifiedDiff_LineNumbers(t *testing.T) {
	diff := unifiedDiff("aaa\nbbb", "aaa\nccc")
	// Line number "2" should appear for the removal of "bbb".
	if !strings.Contains(diff, "2") {
		t.Errorf("diff should contain line numbers:\n%s", diff)
	}
}

func TestUnifiedDiff_Headers(t *testing.T) {
	diff := unifiedDiff("a", "b")
	if !strings.HasPrefix(diff, "--- expected\n+++ actual\n") {
		t.Errorf("diff should start with standard headers:\n%s", diff)
	}
}

// --- Snapshot config tests (TASK-260210-yyzpq3) ---

func TestUpdateSnapshotsVar(t *testing.T) {
	orig := UpdateSnapshots
	defer func() { UpdateSnapshots = orig }()

	UpdateSnapshots = true
	if !UpdateSnapshots {
		t.Fatal("UpdateSnapshots should be true")
	}
	UpdateSnapshots = false
	if UpdateSnapshots {
		t.Fatal("UpdateSnapshots should be false")
	}
}

func TestSnapshotPath_WithBaseDir(t *testing.T) {
	orig := snapshotBaseDir
	defer func() { snapshotBaseDir = orig }()

	snapshotBaseDir = "/tmp/test-snapshots"
	got := snapshotPath("my-view", 1)
	want := "/tmp/test-snapshots/my-view.golden"
	if got != want {
		t.Errorf("snapshotPath = %q, want %q", got, want)
	}
}

func TestSnapshotPath_GoldenExtension(t *testing.T) {
	orig := snapshotBaseDir
	defer func() { snapshotBaseDir = orig }()

	snapshotBaseDir = "/tmp/x"
	got := snapshotPath("foo", 1)
	if !strings.HasSuffix(got, ".golden") {
		t.Errorf("snapshotPath should end with .golden: %q", got)
	}
}

func TestSnapshotPath_DefaultUsesRuntimeCaller(t *testing.T) {
	orig := snapshotBaseDir
	defer func() { snapshotBaseDir = orig }()

	snapshotBaseDir = ""
	got := snapshotPath("test-name", 1)
	if !strings.HasSuffix(got, filepath.Join("testdata", "snapshots", "test-name.golden")) {
		t.Errorf("snapshotPath = %q, expected suffix testdata/snapshots/test-name.golden", got)
	}
}

// --- Snapshot function tests (TASK-260210-3hgcn8 / TASK-260210-15lvy5) ---

func TestSnapshotStr_CreateAndMatch(t *testing.T) {
	dir := t.TempDir()
	withSnapshotDir(t, dir)
	UpdateSnapshots = true

	SnapshotStr(t, "hello world", "basic")

	path := filepath.Join(dir, "basic.golden")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("golden file not created: %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("golden content = %q, want %q", string(data), "hello world")
	}

	// Compare mode should pass.
	UpdateSnapshots = false
	SnapshotStr(t, "hello world", "basic")
}

func TestSnapshotStr_Mismatch(t *testing.T) {
	dir := t.TempDir()
	withSnapshotDir(t, dir)

	if err := os.WriteFile(filepath.Join(dir, "mm.golden"), []byte("expected"), 0o644); err != nil {
		t.Fatal(err)
	}

	UpdateSnapshots = false
	ft := &fakeT{}
	snapshot(ft, "actual", "mm", 1)

	if !ft.failed {
		t.Error("expected snapshot to report failure on mismatch")
	}
	if !strings.Contains(ft.lastErr, "mm") {
		t.Errorf("error should mention snapshot name, got: %s", ft.lastErr)
	}
	if !strings.Contains(ft.lastErr, "expected") || !strings.Contains(ft.lastErr, "actual") {
		t.Errorf("error should contain diff with both values, got: %s", ft.lastErr)
	}
}

func TestSnapshotStr_MissingGoldenFile(t *testing.T) {
	dir := t.TempDir()
	withSnapshotDir(t, dir)
	UpdateSnapshots = false

	ft := &fakeT{}
	runSnapshot(ft, "content", "nonexistent", 1)

	if !ft.fataled {
		t.Error("expected fatal on missing golden file")
	}
	if !strings.Contains(ft.lastErr, "UPDATE_SNAPSHOTS=1") {
		t.Errorf("error should mention UPDATE_SNAPSHOTS=1, got: %s", ft.lastErr)
	}
}

func TestSnapshotStr_StripsANSI(t *testing.T) {
	dir := t.TempDir()
	withSnapshotDir(t, dir)
	UpdateSnapshots = true

	ansiView := "\x1b[31mred text\x1b[0m"
	SnapshotStr(t, ansiView, "ansi-stripped")

	data, err := os.ReadFile(filepath.Join(dir, "ansi-stripped.golden"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "red text" {
		t.Errorf("expected stripped %q, got %q", "red text", string(data))
	}
}

func TestSnapshotStrRaw_PreservesANSI(t *testing.T) {
	dir := t.TempDir()
	withSnapshotDir(t, dir)
	UpdateSnapshots = true

	ansiView := "\x1b[31mred text\x1b[0m"
	SnapshotStrRaw(t, ansiView, "ansi-raw")

	data, err := os.ReadFile(filepath.Join(dir, "ansi-raw.golden"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != ansiView {
		t.Errorf("expected raw %q, got %q", ansiView, string(data))
	}
}

func TestSnapshotView_Model(t *testing.T) {
	dir := t.TempDir()
	withSnapshotDir(t, dir)
	UpdateSnapshots = true

	model := stubModel{view: "\x1b[32mgreen\x1b[0m line"}
	SnapshotView(t, model, "model-stripped")

	data, err := os.ReadFile(filepath.Join(dir, "model-stripped.golden"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "green line" {
		t.Errorf("expected %q, got %q", "green line", string(data))
	}

	// Verify comparison passes.
	UpdateSnapshots = false
	SnapshotView(t, model, "model-stripped")
}

func TestSnapshotViewRaw_Model(t *testing.T) {
	dir := t.TempDir()
	withSnapshotDir(t, dir)
	UpdateSnapshots = true

	rawView := "\x1b[32mgreen\x1b[0m line"
	model := stubModel{view: rawView}
	SnapshotViewRaw(t, model, "model-raw")

	data, err := os.ReadFile(filepath.Join(dir, "model-raw.golden"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != rawView {
		t.Errorf("expected %q, got %q", rawView, string(data))
	}

	// Verify comparison passes.
	UpdateSnapshots = false
	SnapshotViewRaw(t, model, "model-raw")
}

func TestSnapshot_DirectoryAutoCreation(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "nested", "deep")
	withSnapshotDir(t, dir)
	UpdateSnapshots = true

	SnapshotStr(t, "auto-created", "nested-test")

	data, err := os.ReadFile(filepath.Join(dir, "nested-test.golden"))
	if err != nil {
		t.Fatalf("file should have been created in nested dir: %v", err)
	}
	if string(data) != "auto-created" {
		t.Errorf("content = %q, want %q", string(data), "auto-created")
	}
}

func TestSnapshot_UpdateOverwrites(t *testing.T) {
	dir := t.TempDir()
	withSnapshotDir(t, dir)
	UpdateSnapshots = true

	SnapshotStr(t, "version 1", "overwrite")
	SnapshotStr(t, "version 2", "overwrite")

	data, err := os.ReadFile(filepath.Join(dir, "overwrite.golden"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "version 2" {
		t.Errorf("expected overwritten %q, got %q", "version 2", string(data))
	}
}

func TestSnapshot_MultilineContent(t *testing.T) {
	dir := t.TempDir()
	withSnapshotDir(t, dir)
	UpdateSnapshots = true

	content := "line 1\nline 2\nline 3\n"
	SnapshotStr(t, content, "multiline")

	UpdateSnapshots = false
	SnapshotStr(t, content, "multiline")
}

func TestSnapshot_EmptyContent(t *testing.T) {
	dir := t.TempDir()
	withSnapshotDir(t, dir)
	UpdateSnapshots = true

	SnapshotStr(t, "", "empty")

	data, err := os.ReadFile(filepath.Join(dir, "empty.golden"))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "" {
		t.Errorf("expected empty content, got %q", string(data))
	}

	UpdateSnapshots = false
	SnapshotStr(t, "", "empty")
}

func TestSnapshot_MismatchShowsDiff(t *testing.T) {
	dir := t.TempDir()
	withSnapshotDir(t, dir)

	golden := "line 1\nline 2\nline 3"
	if err := os.WriteFile(filepath.Join(dir, "diffcheck.golden"), []byte(golden), 0o644); err != nil {
		t.Fatal(err)
	}

	UpdateSnapshots = false
	ft := &fakeT{}
	snapshot(ft, "line 1\nCHANGED\nline 3", "diffcheck", 1)

	if !ft.failed {
		t.Error("expected failure")
	}
	// Diff should contain the removed and added lines.
	if !strings.Contains(ft.lastErr, "line 2") {
		t.Error("diff should contain removed 'line 2'")
	}
	if !strings.Contains(ft.lastErr, "CHANGED") {
		t.Error("diff should contain added 'CHANGED'")
	}
}
