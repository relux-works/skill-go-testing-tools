package tuitestkit

import (
	"fmt"
	"testing"
)

// --- Example domain for testing ---

// counterState is a simple state for testing reducers.
type counterState struct {
	Count int
	Min   int
	Max   int
}

// counterAction represents actions on a counter.
type counterAction int

const (
	actionIncrement counterAction = iota
	actionDecrement
	actionReset
	actionDouble
)

// counterReduce is a pure reducer for the counter domain.
func counterReduce(s counterState, a counterAction) counterState {
	switch a {
	case actionIncrement:
		s.Count++
	case actionDecrement:
		s.Count--
	case actionReset:
		s.Count = 0
	case actionDouble:
		s.Count *= 2
	}
	// Clamp to min/max
	if s.Count < s.Min {
		s.Count = s.Min
	}
	if s.Count > s.Max {
		s.Count = s.Max
	}
	return s
}

// --- RunReducerTests ---

func TestRunReducerTests_SingleActions(t *testing.T) {
	tests := []ReducerTest[counterState, counterAction]{
		{
			Name:    "increment from zero",
			Initial: counterState{Count: 0, Min: -10, Max: 10},
			Action:  actionIncrement,
			Assert: func(t *testing.T, got counterState) {
				t.Helper()
				if got.Count != 1 {
					t.Errorf("expected Count=1, got %d", got.Count)
				}
			},
		},
		{
			Name:    "decrement from zero",
			Initial: counterState{Count: 0, Min: -10, Max: 10},
			Action:  actionDecrement,
			Assert: func(t *testing.T, got counterState) {
				t.Helper()
				if got.Count != -1 {
					t.Errorf("expected Count=-1, got %d", got.Count)
				}
			},
		},
		{
			Name:    "reset from 5",
			Initial: counterState{Count: 5, Min: -10, Max: 10},
			Action:  actionReset,
			Assert: func(t *testing.T, got counterState) {
				t.Helper()
				if got.Count != 0 {
					t.Errorf("expected Count=0, got %d", got.Count)
				}
			},
		},
		{
			Name:    "double from 3",
			Initial: counterState{Count: 3, Min: -10, Max: 10},
			Action:  actionDouble,
			Assert: func(t *testing.T, got counterState) {
				t.Helper()
				if got.Count != 6 {
					t.Errorf("expected Count=6, got %d", got.Count)
				}
			},
		},
		{
			Name:    "increment clamped at max",
			Initial: counterState{Count: 10, Min: -10, Max: 10},
			Action:  actionIncrement,
			Assert: func(t *testing.T, got counterState) {
				t.Helper()
				if got.Count != 10 {
					t.Errorf("expected Count=10 (clamped), got %d", got.Count)
				}
			},
		},
		{
			Name:    "decrement clamped at min",
			Initial: counterState{Count: -10, Min: -10, Max: 10},
			Action:  actionDecrement,
			Assert: func(t *testing.T, got counterState) {
				t.Helper()
				if got.Count != -10 {
					t.Errorf("expected Count=-10 (clamped), got %d", got.Count)
				}
			},
		},
	}

	RunReducerTests(t, counterReduce, tests)
}

func TestRunReducerTests_EmptySlice(t *testing.T) {
	// Should not panic with empty tests slice
	RunReducerTests(t, counterReduce, []ReducerTest[counterState, counterAction]{})
}

// --- RunReducerSequences ---

func TestRunReducerSequences_MultiStep(t *testing.T) {
	sequences := []ReducerSequence[counterState, counterAction]{
		{
			Name:    "increment three times",
			Initial: counterState{Count: 0, Min: -10, Max: 10},
			Steps: []Step[counterState, counterAction]{
				{Name: "first increment", Action: actionIncrement, Assert: func(t *testing.T, got counterState) {
					t.Helper()
					if got.Count != 1 {
						t.Errorf("after first increment: expected 1, got %d", got.Count)
					}
				}},
				{Name: "second increment", Action: actionIncrement, Assert: func(t *testing.T, got counterState) {
					t.Helper()
					if got.Count != 2 {
						t.Errorf("after second increment: expected 2, got %d", got.Count)
					}
				}},
				{Name: "third increment", Action: actionIncrement},
			},
			Final: func(t *testing.T, got counterState) {
				t.Helper()
				if got.Count != 3 {
					t.Errorf("final state: expected Count=3, got %d", got.Count)
				}
			},
		},
		{
			Name:    "increment then reset",
			Initial: counterState{Count: 0, Min: -10, Max: 10},
			Steps: []Step[counterState, counterAction]{
				{Name: "bump", Action: actionIncrement},
				{Name: "bump again", Action: actionIncrement},
				{Name: "reset", Action: actionReset},
			},
			Final: func(t *testing.T, got counterState) {
				t.Helper()
				if got.Count != 0 {
					t.Errorf("expected Count=0 after reset, got %d", got.Count)
				}
			},
		},
		{
			Name:    "double overflow clamped",
			Initial: counterState{Count: 8, Min: -10, Max: 10},
			Steps: []Step[counterState, counterAction]{
				{Name: "double", Action: actionDouble, Assert: func(t *testing.T, got counterState) {
					t.Helper()
					// 8*2=16 but clamped to 10
					if got.Count != 10 {
						t.Errorf("expected Count=10 (clamped), got %d", got.Count)
					}
				}},
			},
			Final: func(t *testing.T, got counterState) {
				t.Helper()
				if got.Count != 10 {
					t.Errorf("final state: expected Count=10, got %d", got.Count)
				}
			},
		},
	}

	RunReducerSequences(t, counterReduce, sequences)
}

func TestRunReducerSequences_StepsWithoutAssert(t *testing.T) {
	seq := []ReducerSequence[counterState, counterAction]{
		{
			Name:    "all steps without per-step assert",
			Initial: counterState{Count: 0, Min: -100, Max: 100},
			Steps: []Step[counterState, counterAction]{
				{Action: actionIncrement},
				{Action: actionIncrement},
				{Action: actionDouble},
			},
			Final: func(t *testing.T, got counterState) {
				t.Helper()
				// 0 -> 1 -> 2 -> 4
				if got.Count != 4 {
					t.Errorf("expected Count=4, got %d", got.Count)
				}
			},
		},
	}

	RunReducerSequences(t, counterReduce, seq)
}

func TestRunReducerSequences_EmptySteps(t *testing.T) {
	seq := []ReducerSequence[counterState, counterAction]{
		{
			Name:    "no steps",
			Initial: counterState{Count: 42, Min: 0, Max: 100},
			Steps:   nil,
			Final: func(t *testing.T, got counterState) {
				t.Helper()
				if got.Count != 42 {
					t.Errorf("expected initial state preserved, got Count=%d", got.Count)
				}
			},
		},
	}

	RunReducerSequences(t, counterReduce, seq)
}

func TestRunReducerSequences_NilFinal(t *testing.T) {
	// Should not panic when Final is nil
	seq := []ReducerSequence[counterState, counterAction]{
		{
			Name:    "nil final assertion",
			Initial: counterState{Count: 0, Min: -10, Max: 10},
			Steps: []Step[counterState, counterAction]{
				{Action: actionIncrement},
			},
			Final: nil,
		},
	}

	RunReducerSequences(t, counterReduce, seq)
}

func TestRunReducerSequences_EmptySlice(t *testing.T) {
	RunReducerSequences(t, counterReduce, []ReducerSequence[counterState, counterAction]{})
}

// --- Step name fallback ---

func TestRunReducerSequences_StepNameFallback(t *testing.T) {
	// When step has no name, it should use "step-N" as subtest name
	seq := []ReducerSequence[counterState, counterAction]{
		{
			Name:    "unnamed steps fallback",
			Initial: counterState{Count: 0, Min: -10, Max: 10},
			Steps: []Step[counterState, counterAction]{
				{Action: actionIncrement, Assert: func(t *testing.T, got counterState) {
					t.Helper()
					if got.Count != 1 {
						t.Errorf("expected 1, got %d", got.Count)
					}
				}},
			},
			Final: func(t *testing.T, got counterState) {
				t.Helper()
				if got.Count != 1 {
					t.Errorf("expected 1, got %d", got.Count)
				}
			},
		},
	}

	RunReducerSequences(t, counterReduce, seq)
}

// --- InvariantChecker ---

func TestInvariantChecker_PassingInvariants(t *testing.T) {
	checker := NewInvariantChecker(
		Invariant[counterState]{
			Name: "count within bounds",
			Check: func(s counterState) error {
				if s.Count < s.Min || s.Count > s.Max {
					return fmt.Errorf("count %d out of bounds [%d, %d]", s.Count, s.Min, s.Max)
				}
				return nil
			},
		},
	)

	err := checker.Check(counterState{Count: 5, Min: 0, Max: 10})
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}
}

func TestInvariantChecker_FailingInvariant(t *testing.T) {
	checker := NewInvariantChecker(
		Invariant[counterState]{
			Name: "count within bounds",
			Check: func(s counterState) error {
				if s.Count < s.Min || s.Count > s.Max {
					return fmt.Errorf("count %d out of bounds [%d, %d]", s.Count, s.Min, s.Max)
				}
				return nil
			},
		},
	)

	err := checker.Check(counterState{Count: 15, Min: 0, Max: 10})
	if err == nil {
		t.Error("expected invariant violation error, got nil")
	}
}

func TestInvariantChecker_MultipleInvariants(t *testing.T) {
	checker := NewInvariantChecker(
		Invariant[counterState]{
			Name: "non-negative",
			Check: func(s counterState) error {
				if s.Count < 0 {
					return fmt.Errorf("count is negative: %d", s.Count)
				}
				return nil
			},
		},
		Invariant[counterState]{
			Name: "under max",
			Check: func(s counterState) error {
				if s.Count > s.Max {
					return fmt.Errorf("count %d exceeds max %d", s.Count, s.Max)
				}
				return nil
			},
		},
	)

	// Both pass
	err := checker.Check(counterState{Count: 5, Min: 0, Max: 10})
	if err != nil {
		t.Errorf("expected no error, got: %v", err)
	}

	// First fails
	err = checker.Check(counterState{Count: -1, Min: 0, Max: 10})
	if err == nil {
		t.Error("expected error for negative count, got nil")
	}

	// Second fails
	err = checker.Check(counterState{Count: 20, Min: 0, Max: 10})
	if err == nil {
		t.Error("expected error for exceeding max, got nil")
	}
}

func TestInvariantChecker_NoInvariants(t *testing.T) {
	checker := NewInvariantChecker[counterState]()
	err := checker.Check(counterState{Count: 999, Min: 0, Max: 1})
	if err != nil {
		t.Errorf("expected no error with no invariants, got: %v", err)
	}
}

// --- WrapWithInvariants ---

func TestWrapWithInvariants_PassingInvariant(t *testing.T) {
	checker := NewInvariantChecker(
		Invariant[counterState]{
			Name: "count within bounds",
			Check: func(s counterState) error {
				if s.Count < s.Min || s.Count > s.Max {
					return fmt.Errorf("count %d out of bounds [%d, %d]", s.Count, s.Min, s.Max)
				}
				return nil
			},
		},
	)

	wrapped := WrapWithInvariants(t, counterReduce, checker)

	// The clamped reducer should never violate bounds
	state := counterState{Count: 0, Min: -10, Max: 10}
	state = wrapped(state, actionIncrement)
	if state.Count != 1 {
		t.Errorf("expected Count=1, got %d", state.Count)
	}

	state = wrapped(state, actionDecrement)
	if state.Count != 0 {
		t.Errorf("expected Count=0, got %d", state.Count)
	}
}

func TestWrapWithInvariants_UsedWithRunReducerTests(t *testing.T) {
	checker := NewInvariantChecker(
		Invariant[counterState]{
			Name: "count within bounds",
			Check: func(s counterState) error {
				if s.Count < s.Min || s.Count > s.Max {
					return fmt.Errorf("count %d out of bounds [%d, %d]", s.Count, s.Min, s.Max)
				}
				return nil
			},
		},
	)

	wrapped := WrapWithInvariants(t, counterReduce, checker)

	tests := []ReducerTest[counterState, counterAction]{
		{
			Name:    "increment with invariant check",
			Initial: counterState{Count: 9, Min: 0, Max: 10},
			Action:  actionIncrement,
			Assert: func(t *testing.T, got counterState) {
				t.Helper()
				if got.Count != 10 {
					t.Errorf("expected Count=10, got %d", got.Count)
				}
			},
		},
		{
			Name:    "double clamped with invariant check",
			Initial: counterState{Count: 8, Min: 0, Max: 10},
			Action:  actionDouble,
			Assert: func(t *testing.T, got counterState) {
				t.Helper()
				// 8*2=16, clamped to 10
				if got.Count != 10 {
					t.Errorf("expected Count=10 (clamped), got %d", got.Count)
				}
			},
		},
	}

	RunReducerTests(t, wrapped, tests)
}

func TestWrapWithInvariants_UsedWithSequences(t *testing.T) {
	checker := NewInvariantChecker(
		Invariant[counterState]{
			Name: "count within bounds",
			Check: func(s counterState) error {
				if s.Count < s.Min || s.Count > s.Max {
					return fmt.Errorf("count %d out of bounds [%d, %d]", s.Count, s.Min, s.Max)
				}
				return nil
			},
		},
	)

	wrapped := WrapWithInvariants(t, counterReduce, checker)

	sequences := []ReducerSequence[counterState, counterAction]{
		{
			Name:    "sequence with invariants",
			Initial: counterState{Count: 0, Min: -5, Max: 5},
			Steps: []Step[counterState, counterAction]{
				{Name: "inc", Action: actionIncrement},
				{Name: "inc", Action: actionIncrement},
				{Name: "inc", Action: actionIncrement},
				{Name: "double", Action: actionDouble, Assert: func(t *testing.T, got counterState) {
					t.Helper()
					// 3*2=6, clamped to 5
					if got.Count != 5 {
						t.Errorf("expected Count=5 (clamped), got %d", got.Count)
					}
				}},
				{Name: "reset", Action: actionReset},
			},
			Final: func(t *testing.T, got counterState) {
				t.Helper()
				if got.Count != 0 {
					t.Errorf("expected Count=0 after reset, got %d", got.Count)
				}
			},
		},
	}

	RunReducerSequences(t, wrapped, sequences)
}

// --- String action type (different generic instantiation) ---

type listState struct {
	Items []string
}

type listAction struct {
	Kind string
	Item string
}

func listReduce(s listState, a listAction) listState {
	switch a.Kind {
	case "add":
		s.Items = append(append([]string{}, s.Items...), a.Item)
	case "clear":
		s.Items = nil
	}
	return s
}

func TestRunReducerTests_DifferentTypes(t *testing.T) {
	tests := []ReducerTest[listState, listAction]{
		{
			Name:    "add item",
			Initial: listState{Items: nil},
			Action:  listAction{Kind: "add", Item: "hello"},
			Assert: func(t *testing.T, got listState) {
				t.Helper()
				if len(got.Items) != 1 || got.Items[0] != "hello" {
					t.Errorf("expected [hello], got %v", got.Items)
				}
			},
		},
		{
			Name:    "clear items",
			Initial: listState{Items: []string{"a", "b", "c"}},
			Action:  listAction{Kind: "clear"},
			Assert: func(t *testing.T, got listState) {
				t.Helper()
				if got.Items != nil {
					t.Errorf("expected nil items, got %v", got.Items)
				}
			},
		},
	}

	RunReducerTests(t, listReduce, tests)
}

func TestRunReducerSequences_DifferentTypes(t *testing.T) {
	sequences := []ReducerSequence[listState, listAction]{
		{
			Name:    "build a list then clear",
			Initial: listState{},
			Steps: []Step[listState, listAction]{
				{Name: "add first", Action: listAction{Kind: "add", Item: "a"}},
				{Name: "add second", Action: listAction{Kind: "add", Item: "b"}},
				{Name: "clear", Action: listAction{Kind: "clear"}},
			},
			Final: func(t *testing.T, got listState) {
				t.Helper()
				if got.Items != nil {
					t.Errorf("expected nil after clear, got %v", got.Items)
				}
			},
		},
	}

	RunReducerSequences(t, listReduce, sequences)
}

func TestWrapWithInvariants_DifferentTypes(t *testing.T) {
	checker := NewInvariantChecker(
		Invariant[listState]{
			Name: "max 5 items",
			Check: func(s listState) error {
				if len(s.Items) > 5 {
					return fmt.Errorf("too many items: %d", len(s.Items))
				}
				return nil
			},
		},
	)

	wrapped := WrapWithInvariants(t, listReduce, checker)

	state := listState{}
	for i := range 5 {
		state = wrapped(state, listAction{Kind: "add", Item: fmt.Sprintf("item-%d", i)})
	}
	if len(state.Items) != 5 {
		t.Errorf("expected 5 items, got %d", len(state.Items))
	}
}

// --- Edge case: state immutability ---

func TestReducerDoesNotMutateOriginal(t *testing.T) {
	original := counterState{Count: 5, Min: 0, Max: 10}
	result := counterReduce(original, actionIncrement)

	if original.Count != 5 {
		t.Errorf("original state was mutated: Count=%d", original.Count)
	}
	if result.Count != 6 {
		t.Errorf("expected result Count=6, got %d", result.Count)
	}
}
