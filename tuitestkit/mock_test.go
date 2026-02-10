package tuitestkit

import (
	"errors"
	"fmt"
	"sync"
	"testing"
)

// --- MockCallRecorder tests ---

func TestMockCallRecorder_Record(t *testing.T) {
	var r MockCallRecorder
	r.Record("Foo")
	r.Record("Bar", 1, "hello")

	calls := r.Calls()
	if len(calls) != 2 {
		t.Fatalf("expected 2 calls, got %d", len(calls))
	}
	if calls[0].Method != "Foo" {
		t.Errorf("expected method Foo, got %s", calls[0].Method)
	}
	if calls[0].Args != nil {
		t.Errorf("expected nil args for Foo, got %v", calls[0].Args)
	}
	if calls[1].Method != "Bar" {
		t.Errorf("expected method Bar, got %s", calls[1].Method)
	}
	if len(calls[1].Args) != 2 || calls[1].Args[0] != 1 || calls[1].Args[1] != "hello" {
		t.Errorf("unexpected args for Bar: %v", calls[1].Args)
	}
}

func TestMockCallRecorder_CallCount(t *testing.T) {
	var r MockCallRecorder
	if r.CallCount("Foo") != 0 {
		t.Error("expected 0 calls before any recording")
	}
	r.Record("Foo")
	r.Record("Bar")
	r.Record("Foo")
	r.Record("Foo")
	if got := r.CallCount("Foo"); got != 3 {
		t.Errorf("expected Foo called 3 times, got %d", got)
	}
	if got := r.CallCount("Bar"); got != 1 {
		t.Errorf("expected Bar called 1 time, got %d", got)
	}
	if got := r.CallCount("Missing"); got != 0 {
		t.Errorf("expected Missing called 0 times, got %d", got)
	}
}

func TestMockCallRecorder_Calls_ReturnsCopy(t *testing.T) {
	var r MockCallRecorder
	r.Record("A")
	calls := r.Calls()
	calls[0].Method = "MUTATED"
	// Original should be unchanged.
	original := r.Calls()
	if original[0].Method != "A" {
		t.Error("Calls() did not return a copy — internal state was mutated")
	}
}

func TestMockCallRecorder_CallsFor(t *testing.T) {
	var r MockCallRecorder
	r.Record("Exec", "ls")
	r.Record("Open", "/tmp")
	r.Record("Exec", "pwd")
	r.Record("Exec", "cat", "file.txt")

	execCalls := r.CallsFor("Exec")
	if len(execCalls) != 3 {
		t.Fatalf("expected 3 Exec calls, got %d", len(execCalls))
	}
	if execCalls[0].Args[0] != "ls" {
		t.Errorf("expected first Exec arg 'ls', got %v", execCalls[0].Args[0])
	}
	if execCalls[1].Args[0] != "pwd" {
		t.Errorf("expected second Exec arg 'pwd', got %v", execCalls[1].Args[0])
	}

	openCalls := r.CallsFor("Open")
	if len(openCalls) != 1 {
		t.Fatalf("expected 1 Open call, got %d", len(openCalls))
	}

	missing := r.CallsFor("Missing")
	if len(missing) != 0 {
		t.Errorf("expected 0 calls for Missing, got %d", len(missing))
	}
}

func TestMockCallRecorder_Reset(t *testing.T) {
	var r MockCallRecorder
	r.Record("A")
	r.Record("B")
	r.Reset()
	if got := len(r.Calls()); got != 0 {
		t.Errorf("expected 0 calls after Reset, got %d", got)
	}
	if got := r.CallCount("A"); got != 0 {
		t.Errorf("expected 0 for A after Reset, got %d", got)
	}
}

func TestMockCallRecorder_Concurrent(t *testing.T) {
	var r MockCallRecorder
	const goroutines = 100
	const callsPerGoroutine = 50

	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			for j := 0; j < callsPerGoroutine; j++ {
				r.Record(fmt.Sprintf("method-%d", id), j)
			}
		}(i)
	}
	wg.Wait()

	total := len(r.Calls())
	expected := goroutines * callsPerGoroutine
	if total != expected {
		t.Errorf("expected %d total calls, got %d", expected, total)
	}
}

// --- MockResponseMap tests ---

func TestNewMockResponseMap(t *testing.T) {
	m := NewMockResponseMap()
	if m == nil {
		t.Fatal("NewMockResponseMap returned nil")
	}
}

func TestMockResponseMap_Set_Get(t *testing.T) {
	m := NewMockResponseMap()
	data := []byte(`{"items":[]}`)
	m.Set("TreeJSON", data, nil)

	got, err := m.Get("TreeJSON")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if string(got) != string(data) {
		t.Errorf("expected %s, got %s", data, got)
	}
}

func TestMockResponseMap_Get_NotFound(t *testing.T) {
	m := NewMockResponseMap()
	data, err := m.Get("missing")
	if data != nil {
		t.Errorf("expected nil data for missing key, got %v", data)
	}
	if err != nil {
		t.Errorf("expected nil error for missing key, got %v", err)
	}
}

func TestMockResponseMap_Set_WithError(t *testing.T) {
	m := NewMockResponseMap()
	expectedErr := errors.New("connection refused")
	m.Set("Fetch", []byte("partial"), expectedErr)

	data, err := m.Get("Fetch")
	if string(data) != "partial" {
		t.Errorf("expected 'partial', got %s", data)
	}
	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}

func TestMockResponseMap_SetError(t *testing.T) {
	m := NewMockResponseMap()
	expectedErr := errors.New("timeout")
	m.SetError("Slow", expectedErr)

	data, err := m.Get("Slow")
	if data != nil {
		t.Errorf("expected nil data from SetError, got %v", data)
	}
	if err != expectedErr {
		t.Errorf("expected %v, got %v", expectedErr, err)
	}
}

func TestMockResponseMap_Overwrite(t *testing.T) {
	m := NewMockResponseMap()
	m.Set("Key", []byte("first"), nil)
	m.Set("Key", []byte("second"), nil)

	data, _ := m.Get("Key")
	if string(data) != "second" {
		t.Errorf("expected 'second' after overwrite, got %s", data)
	}
}

func TestMockResponseMap_Concurrent(t *testing.T) {
	m := NewMockResponseMap()
	const goroutines = 100

	var wg sync.WaitGroup
	wg.Add(goroutines * 2)

	// Writers
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", id)
			m.Set(key, []byte(fmt.Sprintf("data-%d", id)), nil)
		}(i)
	}

	// Readers
	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			key := fmt.Sprintf("key-%d", id)
			m.Get(key) // just ensure no panic
		}(i)
	}

	wg.Wait()
}

// --- Assertion helper tests ---

// mockTB is a fake testing.TB for verifying assertion helper behavior.
type mockTB struct {
	testing.TB // embed interface for unimplemented methods
	failed     bool
	logs       []string
}

func (m *mockTB) Helper() {}
func (m *mockTB) Errorf(format string, args ...any) {
	m.failed = true
	m.logs = append(m.logs, fmt.Sprintf(format, args...))
}
func (m *mockTB) Log(args ...any) {
	m.logs = append(m.logs, fmt.Sprint(args...))
}

func TestAssertCalled_Pass(t *testing.T) {
	var r MockCallRecorder
	r.Record("Do")
	tb := &mockTB{}
	AssertCalled(tb, &r, "Do")
	if tb.failed {
		t.Error("AssertCalled should pass when method was called")
	}
}

func TestAssertCalled_Fail(t *testing.T) {
	var r MockCallRecorder
	tb := &mockTB{}
	AssertCalled(tb, &r, "Do")
	if !tb.failed {
		t.Error("AssertCalled should fail when method was not called")
	}
}

func TestAssertNotCalled_Pass(t *testing.T) {
	var r MockCallRecorder
	tb := &mockTB{}
	AssertNotCalled(tb, &r, "Do")
	if tb.failed {
		t.Error("AssertNotCalled should pass when method was not called")
	}
}

func TestAssertNotCalled_Fail(t *testing.T) {
	var r MockCallRecorder
	r.Record("Do")
	tb := &mockTB{}
	AssertNotCalled(tb, &r, "Do")
	if !tb.failed {
		t.Error("AssertNotCalled should fail when method was called")
	}
}

func TestAssertCalledN_Pass(t *testing.T) {
	var r MockCallRecorder
	r.Record("Do")
	r.Record("Do")
	r.Record("Do")
	tb := &mockTB{}
	AssertCalledN(tb, &r, "Do", 3)
	if tb.failed {
		t.Error("AssertCalledN should pass when count matches")
	}
}

func TestAssertCalledN_Fail(t *testing.T) {
	var r MockCallRecorder
	r.Record("Do")
	tb := &mockTB{}
	AssertCalledN(tb, &r, "Do", 5)
	if !tb.failed {
		t.Error("AssertCalledN should fail when count doesn't match")
	}
}

func TestAssertCalledN_Zero(t *testing.T) {
	var r MockCallRecorder
	tb := &mockTB{}
	AssertCalledN(tb, &r, "Do", 0)
	if tb.failed {
		t.Error("AssertCalledN with n=0 should pass when method was never called")
	}
}

func TestAssertCalledWith_Pass(t *testing.T) {
	var r MockCallRecorder
	r.Record("Exec", "ls", "-la")
	tb := &mockTB{}
	AssertCalledWith(tb, &r, "Exec", "ls", "-la")
	if tb.failed {
		t.Errorf("AssertCalledWith should pass when args match; logs: %v", tb.logs)
	}
}

func TestAssertCalledWith_Fail_WrongArgs(t *testing.T) {
	var r MockCallRecorder
	r.Record("Exec", "ls")
	tb := &mockTB{}
	AssertCalledWith(tb, &r, "Exec", "pwd")
	if !tb.failed {
		t.Error("AssertCalledWith should fail when no call matches the args")
	}
}

func TestAssertCalledWith_Fail_NeverCalled(t *testing.T) {
	var r MockCallRecorder
	tb := &mockTB{}
	AssertCalledWith(tb, &r, "Exec", "ls")
	if !tb.failed {
		t.Error("AssertCalledWith should fail when method was never called")
	}
}

func TestAssertCalledWith_MatchesAmongMultipleCalls(t *testing.T) {
	var r MockCallRecorder
	r.Record("Exec", "ls")
	r.Record("Exec", "pwd")
	r.Record("Exec", "cat", "file.txt")
	tb := &mockTB{}
	AssertCalledWith(tb, &r, "Exec", "pwd")
	if tb.failed {
		t.Error("AssertCalledWith should find match among multiple calls")
	}
}

func TestAssertCalledWith_DeepEqual(t *testing.T) {
	var r MockCallRecorder
	r.Record("Process", []string{"a", "b"}, map[string]int{"x": 1})
	tb := &mockTB{}
	AssertCalledWith(tb, &r, "Process", []string{"a", "b"}, map[string]int{"x": 1})
	if tb.failed {
		t.Errorf("AssertCalledWith should use DeepEqual for complex types; logs: %v", tb.logs)
	}
}

func TestAssertCalledWith_NoArgs(t *testing.T) {
	var r MockCallRecorder
	r.Record("Ping")
	tb := &mockTB{}
	// Record with no args creates nil Args slice.
	// AssertCalledWith with no variadic args passes nil.
	// We need to handle this: both should be "no args".
	AssertCalledWith(tb, &r, "Ping")
	if tb.failed {
		t.Errorf("AssertCalledWith with no args should match Record with no args; logs: %v", tb.logs)
	}
}

// --- Composition / integration test ---

func TestMock_Composition(t *testing.T) {
	// Demonstrates how users compose MockCallRecorder + MockResponseMap.
	type myMock struct {
		MockCallRecorder
		responses *MockResponseMap
	}

	m := &myMock{responses: NewMockResponseMap()}
	m.responses.Set("TreeJSON", []byte(`{"root":{}}`), nil)

	// Simulate a method call.
	m.Record("TreeJSON")
	data, err := m.responses.Get("TreeJSON")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != `{"root":{}}` {
		t.Errorf("unexpected data: %s", data)
	}
	AssertCalled(t, &m.MockCallRecorder, "TreeJSON")
	AssertCalledN(t, &m.MockCallRecorder, "TreeJSON", 1)
}

func TestMock_Composition_ErrorPath(t *testing.T) {
	type myMock struct {
		MockCallRecorder
		responses *MockResponseMap
	}

	m := &myMock{responses: NewMockResponseMap()}
	m.responses.SetError("Fail", errors.New("boom"))

	m.Record("Fail")
	data, err := m.responses.Get("Fail")

	if data != nil {
		t.Errorf("expected nil data, got %v", data)
	}
	if err == nil || err.Error() != "boom" {
		t.Errorf("expected error 'boom', got %v", err)
	}
	AssertCalled(t, &m.MockCallRecorder, "Fail")
}
