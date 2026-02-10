package tuitestkit

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// --- Test model ---

type incMsg struct{}
type decMsg struct{}

// cmdMsg triggers the model to return a cmd from Update.
type cmdMsg struct {
	cmd tea.Cmd
}

type counterModel struct {
	count int
}

func (m counterModel) Init() tea.Cmd { return nil }

func (m counterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case incMsg:
		m.count++
	case decMsg:
		m.count--
	case cmdMsg:
		return m, msg.cmd
	}
	return m, nil
}

func (m counterModel) View() string {
	return fmt.Sprintf("count: %d", m.count)
}

// --- Send tests ---

func TestSend_NoMessages(t *testing.T) {
	m := counterModel{count: 42}
	got := Send(m)
	if got.count != 42 {
		t.Errorf("Send with no messages: count = %d, want 42", got.count)
	}
}

func TestSend_SingleMessage(t *testing.T) {
	m := counterModel{count: 0}
	got := Send(m, incMsg{})
	if got.count != 1 {
		t.Errorf("Send single inc: count = %d, want 1", got.count)
	}
}

func TestSend_MultipleMessages(t *testing.T) {
	m := counterModel{count: 0}
	got := Send(m, incMsg{}, incMsg{}, incMsg{}, decMsg{})
	if got.count != 2 {
		t.Errorf("Send multiple: count = %d, want 2", got.count)
	}
}

func TestSend_PreservesConcreteType(t *testing.T) {
	m := counterModel{count: 5}
	got := Send(m, incMsg{})

	// The returned value should be counterModel, not tea.Model.
	// If this compiles, the type is preserved. We also verify the value.
	var _ counterModel = got
	if got.count != 6 {
		t.Errorf("Send preserves type: count = %d, want 6", got.count)
	}
}

func TestSend_DoesNotMutateOriginal(t *testing.T) {
	m := counterModel{count: 10}
	got := Send(m, incMsg{})
	if m.count != 10 {
		t.Errorf("original model mutated: count = %d, want 10", m.count)
	}
	if got.count != 11 {
		t.Errorf("returned model: count = %d, want 11", got.count)
	}
}

func TestSend_WithKeyMsgs(t *testing.T) {
	// Verify Send works with tea.Msg interface (not just concrete types).
	// Key() returns tea.KeyMsg which implements tea.Msg.
	m := counterModel{count: 0}
	// Keys that the counter doesn't handle — count should remain 0.
	got := Send(m, Key("a"), Key("enter"))
	if got.count != 0 {
		t.Errorf("Send with unhandled keys: count = %d, want 0", got.count)
	}
}

// --- SendAndCollect tests ---

func TestSendAndCollect_NoMessages(t *testing.T) {
	m := counterModel{count: 5}
	got, cmds := SendAndCollect(m)
	if got.count != 5 {
		t.Errorf("SendAndCollect no msgs: count = %d, want 5", got.count)
	}
	if len(cmds) != 0 {
		t.Errorf("SendAndCollect no msgs: got %d cmds, want 0", len(cmds))
	}
}

func TestSendAndCollect_CollectsCmds(t *testing.T) {
	dummyCmd := func() tea.Msg { return incMsg{} }
	m := counterModel{count: 0}

	got, cmds := SendAndCollect(m, cmdMsg{cmd: dummyCmd})
	if got.count != 0 {
		t.Errorf("SendAndCollect: count = %d, want 0", got.count)
	}
	if len(cmds) != 1 {
		t.Fatalf("SendAndCollect: got %d cmds, want 1", len(cmds))
	}
	// Execute the collected cmd and verify it produces the right message.
	msg := cmds[0]()
	if _, ok := msg.(incMsg); !ok {
		t.Errorf("collected cmd returned %T, want incMsg", msg)
	}
}

func TestSendAndCollect_SkipsNilCmds(t *testing.T) {
	m := counterModel{count: 0}
	// incMsg returns nil cmd, cmdMsg returns a real cmd.
	dummyCmd := func() tea.Msg { return decMsg{} }

	got, cmds := SendAndCollect(m, incMsg{}, cmdMsg{cmd: dummyCmd}, incMsg{})
	if got.count != 2 {
		t.Errorf("SendAndCollect: count = %d, want 2", got.count)
	}
	// Only cmdMsg produced a non-nil cmd.
	if len(cmds) != 1 {
		t.Errorf("SendAndCollect: got %d cmds, want 1", len(cmds))
	}
}

func TestSendAndCollect_MultipleCmds(t *testing.T) {
	cmd1 := func() tea.Msg { return incMsg{} }
	cmd2 := func() tea.Msg { return decMsg{} }
	m := counterModel{count: 0}

	_, cmds := SendAndCollect(m, cmdMsg{cmd: cmd1}, cmdMsg{cmd: cmd2})
	if len(cmds) != 2 {
		t.Fatalf("SendAndCollect: got %d cmds, want 2", len(cmds))
	}

	// Verify first cmd.
	msg1 := cmds[0]()
	if _, ok := msg1.(incMsg); !ok {
		t.Errorf("cmd[0] returned %T, want incMsg", msg1)
	}

	// Verify second cmd.
	msg2 := cmds[1]()
	if _, ok := msg2.(decMsg); !ok {
		t.Errorf("cmd[1] returned %T, want decMsg", msg2)
	}
}

func TestSendAndCollect_PreservesConcreteType(t *testing.T) {
	m := counterModel{count: 3}
	got, _ := SendAndCollect(m, incMsg{})
	var _ counterModel = got
	if got.count != 4 {
		t.Errorf("SendAndCollect preserves type: count = %d, want 4", got.count)
	}
}

// --- ExecCmds tests ---

func TestExecCmds_NilInput(t *testing.T) {
	msgs := ExecCmds()
	if len(msgs) != 0 {
		t.Errorf("ExecCmds(): got %d msgs, want 0", len(msgs))
	}
}

func TestExecCmds_NilCmdInSlice(t *testing.T) {
	msgs := ExecCmds(nil, nil, nil)
	if len(msgs) != 0 {
		t.Errorf("ExecCmds(nil, nil, nil): got %d msgs, want 0", len(msgs))
	}
}

func TestExecCmds_RegularCmds(t *testing.T) {
	cmd1 := func() tea.Msg { return incMsg{} }
	cmd2 := func() tea.Msg { return decMsg{} }

	msgs := ExecCmds(cmd1, cmd2)
	if len(msgs) != 2 {
		t.Fatalf("ExecCmds: got %d msgs, want 2", len(msgs))
	}
	if _, ok := msgs[0].(incMsg); !ok {
		t.Errorf("msgs[0] = %T, want incMsg", msgs[0])
	}
	if _, ok := msgs[1].(decMsg); !ok {
		t.Errorf("msgs[1] = %T, want decMsg", msgs[1])
	}
}

func TestExecCmds_CmdReturningNilMsg(t *testing.T) {
	nilCmd := func() tea.Msg { return nil }
	realCmd := func() tea.Msg { return incMsg{} }

	msgs := ExecCmds(nilCmd, realCmd)
	if len(msgs) != 1 {
		t.Fatalf("ExecCmds: got %d msgs, want 1", len(msgs))
	}
	if _, ok := msgs[0].(incMsg); !ok {
		t.Errorf("msgs[0] = %T, want incMsg", msgs[0])
	}
}

func TestExecCmds_BatchMsgRecursion(t *testing.T) {
	innerCmd1 := func() tea.Msg { return incMsg{} }
	innerCmd2 := func() tea.Msg { return decMsg{} }

	batchCmd := func() tea.Msg {
		return tea.BatchMsg{innerCmd1, innerCmd2}
	}

	msgs := ExecCmds(batchCmd)
	if len(msgs) != 2 {
		t.Fatalf("ExecCmds batch: got %d msgs, want 2", len(msgs))
	}
	if _, ok := msgs[0].(incMsg); !ok {
		t.Errorf("msgs[0] = %T, want incMsg", msgs[0])
	}
	if _, ok := msgs[1].(decMsg); !ok {
		t.Errorf("msgs[1] = %T, want decMsg", msgs[1])
	}
}

func TestExecCmds_NestedBatchMsg(t *testing.T) {
	leafCmd := func() tea.Msg { return incMsg{} }

	innerBatch := func() tea.Msg {
		return tea.BatchMsg{leafCmd, leafCmd}
	}

	outerBatch := func() tea.Msg {
		return tea.BatchMsg{innerBatch, leafCmd}
	}

	msgs := ExecCmds(outerBatch)
	// outerBatch -> [innerBatch, leafCmd]
	// innerBatch -> [leafCmd, leafCmd]
	// Total: 3 incMsg
	if len(msgs) != 3 {
		t.Fatalf("ExecCmds nested batch: got %d msgs, want 3", len(msgs))
	}
	for i, msg := range msgs {
		if _, ok := msg.(incMsg); !ok {
			t.Errorf("msgs[%d] = %T, want incMsg", i, msg)
		}
	}
}

func TestExecCmds_BatchWithNilCmd(t *testing.T) {
	realCmd := func() tea.Msg { return incMsg{} }

	batchCmd := func() tea.Msg {
		return tea.BatchMsg{nil, realCmd, nil}
	}

	msgs := ExecCmds(batchCmd)
	if len(msgs) != 1 {
		t.Fatalf("ExecCmds batch with nils: got %d msgs, want 1", len(msgs))
	}
	if _, ok := msgs[0].(incMsg); !ok {
		t.Errorf("msgs[0] = %T, want incMsg", msgs[0])
	}
}

func TestExecCmds_MixedRegularAndBatch(t *testing.T) {
	regularCmd := func() tea.Msg { return decMsg{} }
	innerCmd := func() tea.Msg { return incMsg{} }
	batchCmd := func() tea.Msg {
		return tea.BatchMsg{innerCmd}
	}

	msgs := ExecCmds(regularCmd, batchCmd)
	if len(msgs) != 2 {
		t.Fatalf("ExecCmds mixed: got %d msgs, want 2", len(msgs))
	}
	if _, ok := msgs[0].(decMsg); !ok {
		t.Errorf("msgs[0] = %T, want decMsg", msgs[0])
	}
	if _, ok := msgs[1].(incMsg); !ok {
		t.Errorf("msgs[1] = %T, want incMsg", msgs[1])
	}
}

// --- Integration: SendAndCollect + ExecCmds + Send ---

func TestIntegration_SendCollectExecSend(t *testing.T) {
	// Model that returns a cmd on cmdMsg, then we exec and feed back.
	cmd := func() tea.Msg { return incMsg{} }
	m := counterModel{count: 0}

	// Send a cmdMsg that will produce a cmd.
	m, cmds := SendAndCollect(m, cmdMsg{cmd: cmd})
	if m.count != 0 {
		t.Errorf("after SendAndCollect: count = %d, want 0", m.count)
	}

	// Execute the cmds to get messages.
	msgs := ExecCmds(cmds...)
	if len(msgs) != 1 {
		t.Fatalf("ExecCmds: got %d msgs, want 1", len(msgs))
	}

	// Feed the resulting messages back to the model.
	m = Send(m, msgs...)
	if m.count != 1 {
		t.Errorf("after feeding exec results: count = %d, want 1", m.count)
	}
}
