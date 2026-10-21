package sequence

import (
	"errors"
	"iter"
	rand "math/rand/v2"
	"testing"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// tableStubNamesSubject binds Names into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func tableStubNamesSubject(tb assert.TB) stub.Subject[TableNamesCall, TableNamesReturn] {
	tb.Helper()
	s := NewTableStub(tb)
	return stub.Subject[TableNamesCall, TableNamesReturn]{
		Stub: s.OnNames.MethodStub,
		Call: func() {
			_ = s.Names()
		},
		Result: func() TableNamesReturn {
			got0 := s.Names()
			return TableNamesReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnNames.Func(func() iter.Seq[string] {
				mark()
				var z0 iter.Seq[string]
				return z0
			})
		},
	}
}

// TestTableStubNames pins how Names answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestTableStubNames(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Names", tableStubNamesSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want0 iter.Seq[string]
		want0 = iter.Seq[string](func(func(string) bool) {})
		s.OnNames.Returns(want0)
		got0 := s.Names()
		assert.Equal(t, got0, want0, "Names must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		_ = s.Names()
		s.OnNames.AssertCalledOnce(t, "Names must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var seen []TableNamesCall
		s.OnNames.OnRecord(func(c TableNamesCall) { seen = append(seen, c) })
		_ = s.Names()
		_ = s.Names()
		assert.Length(t, seen, 2, "OnRecord must fire once per Names call")
	})

	t.Run("wires WithTableNames at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewTableStub(t, WithTableNames(func() iter.Seq[string] {
			called = true
			var z0 iter.Seq[string]
			return z0
		}))
		_ = s.Names()
		assert.True(t, called, "WithTableNames must install the override")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want0 iter.Seq[string]
		want0 = iter.Seq[string](func(func(string) bool) {})
		s.OnNames.Returns(want0)
		_ = s.Names()
		s.ResetCalls()
		got0 := s.Names()
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})

	t.Run("yields the items Yields configured", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want string
		s.OnNames.Yields(want, want)
		got0 := s.Names()
		seen := 0
		for range got0 {
			seen++
		}
		assert.Equal(t, seen, 2, "Names must yield every configured item")
	})

	t.Run("stops yielding when the caller breaks", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want string
		s.OnNames.Yields(want, want, want)
		got0 := s.Names()
		seen := 0
		for range got0 {
			seen++
			break
		}
		assert.Equal(t, seen, 1, "Names must stop when the caller breaks")
	})
}

// tableStubRowsSubject binds Rows into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func tableStubRowsSubject(tb assert.TB) stub.Subject[TableRowsCall, TableRowsReturn] {
	tb.Helper()
	s := NewTableStub(tb)
	return stub.Subject[TableRowsCall, TableRowsReturn]{
		Stub: s.OnRows.MethodStub,
		Call: func() {
			var a0 int
			a0 = 42
			_ = s.Rows(a0)
		},
		Result: func() TableRowsReturn {
			var a0 int
			a0 = 42
			got0 := s.Rows(a0)
			return TableRowsReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnRows.Func(func(_ int) iter.Seq2[Item, error] {
				mark()
				var z0 iter.Seq2[Item, error]
				return z0
			})
		},
	}
}

// TestTableStubRows pins how Rows answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestTableStubRows(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Rows", tableStubRowsSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want0 iter.Seq2[Item, error]
		want0 = iter.Seq2[Item, error](func(func(Item, error) bool) {})
		s.OnRows.Returns(want0)
		var a0 int
		a0 = 42
		got0 := s.Rows(a0)
		assert.Equal(t, got0, want0, "Rows must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var a0 int
		a0 = 42
		_ = s.Rows(a0)
		got := s.OnRows.AssertCalledOnce(t, "Rows must record the call")
		assert.Equal(t, got.Limit, a0, "the recorded call carries Limit")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var seen []TableRowsCall
		s.OnRows.OnRecord(func(c TableRowsCall) { seen = append(seen, c) })
		var a0 int
		a0 = 42
		_ = s.Rows(a0)
		_ = s.Rows(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Rows call")
	})

	t.Run("wires WithTableRows at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewTableStub(t, WithTableRows(func(_ int) iter.Seq2[Item, error] {
			called = true
			var z0 iter.Seq2[Item, error]
			return z0
		}))
		var a0 int
		a0 = 42
		_ = s.Rows(a0)
		assert.True(t, called, "WithTableRows must install the override")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want0 iter.Seq2[Item, error]
		want0 = iter.Seq2[Item, error](func(func(Item, error) bool) {})
		s.OnRows.Returns(want0)
		var a0 int
		a0 = 42
		_ = s.Rows(a0)
		s.ResetCalls()
		got0 := s.Rows(a0)
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})

	t.Run("yields the items Yields configured", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want Item
		s.OnRows.Yields(want, want)
		var a0 int
		a0 = 42
		got0 := s.Rows(a0)
		seen := 0
		for range got0 {
			seen++
		}
		assert.Equal(t, seen, 2, "Rows must yield every configured item")
	})

	t.Run("stops yielding when the caller breaks", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want Item
		s.OnRows.Yields(want, want, want)
		var a0 int
		a0 = 42
		got0 := s.Rows(a0)
		seen := 0
		for range got0 {
			seen++
			break
		}
		assert.Equal(t, seen, 1, "Rows must stop when the caller breaks")
	})

	t.Run("yields the items YieldsError configured and then fails", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var item Item
		want := errors.New("Rows-yields")
		s.OnRows.YieldsError([]Item{item, item}, want)
		var a0 int
		a0 = 42
		got0 := s.Rows(a0)
		seen := 0
		var last error
		for _, err := range got0 {
			if err != nil {
				last = err
				break
			}
			seen++
		}
		assert.Equal(t, seen, 2, "Rows must yield every item before failing")
		assert.ErrorIs(t, last, want, "Rows must end with the configured error")
	})

	t.Run("stops a failing sequence when the caller breaks", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var item Item
		s.OnRows.YieldsError([]Item{item, item, item}, errors.New("Rows-yields"))
		var a0 int
		a0 = 42
		got0 := s.Rows(a0)
		seen := 0
		for range got0 {
			seen++
			break
		}
		assert.Equal(t, seen, 1, "Rows must stop when the caller breaks")
	})
}

// tableStubCountsSubject binds Counts into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func tableStubCountsSubject(tb assert.TB) stub.Subject[TableCountsCall, TableCountsReturn] {
	tb.Helper()
	s := NewTableStub(tb)
	return stub.Subject[TableCountsCall, TableCountsReturn]{
		Stub: s.OnCounts.MethodStub,
		Call: func() {
			_ = s.Counts()
		},
		Result: func() TableCountsReturn {
			got0 := s.Counts()
			return TableCountsReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnCounts.Func(func() iter.Seq2[string, int] {
				mark()
				var z0 iter.Seq2[string, int]
				return z0
			})
		},
	}
}

// TestTableStubCounts pins how Counts answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestTableStubCounts(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Counts", tableStubCountsSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want0 iter.Seq2[string, int]
		want0 = iter.Seq2[string, int](func(func(string, int) bool) {})
		s.OnCounts.Returns(want0)
		got0 := s.Counts()
		assert.Equal(t, got0, want0, "Counts must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		_ = s.Counts()
		s.OnCounts.AssertCalledOnce(t, "Counts must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var seen []TableCountsCall
		s.OnCounts.OnRecord(func(c TableCountsCall) { seen = append(seen, c) })
		_ = s.Counts()
		_ = s.Counts()
		assert.Length(t, seen, 2, "OnRecord must fire once per Counts call")
	})

	t.Run("wires WithTableCounts at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewTableStub(t, WithTableCounts(func() iter.Seq2[string, int] {
			called = true
			var z0 iter.Seq2[string, int]
			return z0
		}))
		_ = s.Counts()
		assert.True(t, called, "WithTableCounts must install the override")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want0 iter.Seq2[string, int]
		want0 = iter.Seq2[string, int](func(func(string, int) bool) {})
		s.OnCounts.Returns(want0)
		_ = s.Counts()
		s.ResetCalls()
		got0 := s.Counts()
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})

	t.Run("yields the items Yields configured", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want string
		s.OnCounts.Yields(want, want)
		got0 := s.Counts()
		seen := 0
		for range got0 {
			seen++
		}
		assert.Equal(t, seen, 2, "Counts must yield every configured item")
	})

	t.Run("stops yielding when the caller breaks", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want string
		s.OnCounts.Yields(want, want, want)
		got0 := s.Counts()
		seen := 0
		for range got0 {
			seen++
			break
		}
		assert.Equal(t, seen, 1, "Counts must stop when the caller breaks")
	})
}

// tableStubTailSubject binds Tail into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func tableStubTailSubject(tb assert.TB) stub.Subject[TableTailCall, TableTailReturn] {
	tb.Helper()
	s := NewTableStub(tb)
	return stub.Subject[TableTailCall, TableTailReturn]{
		Stub: s.OnTail.MethodStub,
		Call: func() {
			_ = s.Tail()
		},
		Result: func() TableTailReturn {
			got0 := s.Tail()
			return TableTailReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnTail.Func(func() Lines {
				mark()
				var z0 Lines
				return z0
			})
		},
	}
}

// TestTableStubTail pins how Tail answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestTableStubTail(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Tail", tableStubTailSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want0 Lines
		want0 = iter.Seq[string](func(func(string) bool) {})
		s.OnTail.Returns(want0)
		got0 := s.Tail()
		assert.Equal(t, got0, want0, "Tail must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		_ = s.Tail()
		s.OnTail.AssertCalledOnce(t, "Tail must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var seen []TableTailCall
		s.OnTail.OnRecord(func(c TableTailCall) { seen = append(seen, c) })
		_ = s.Tail()
		_ = s.Tail()
		assert.Length(t, seen, 2, "OnRecord must fire once per Tail call")
	})

	t.Run("wires WithTableTail at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewTableStub(t, WithTableTail(func() Lines {
			called = true
			var z0 Lines
			return z0
		}))
		_ = s.Tail()
		assert.True(t, called, "WithTableTail must install the override")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want0 Lines
		want0 = iter.Seq[string](func(func(string) bool) {})
		s.OnTail.Returns(want0)
		_ = s.Tail()
		s.ResetCalls()
		got0 := s.Tail()
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})

	t.Run("yields the items Yields configured", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want string
		s.OnTail.Yields(want, want)
		got0 := s.Tail()
		seen := 0
		for range got0 {
			seen++
		}
		assert.Equal(t, seen, 2, "Tail must yield every configured item")
	})

	t.Run("stops yielding when the caller breaks", func(t *testing.T) {
		t.Parallel()
		s := NewTableStub(t)
		var want string
		s.OnTail.Yields(want, want, want)
		got0 := s.Tail()
		seen := 0
		for range got0 {
			seen++
			break
		}
		assert.Equal(t, seen, 1, "Tail must stop when the caller breaks")
	})
}

// tableStubDouble describes how to build a TableStub under each option
// whose effect is the same whatever a method's signature.
//
// Names stands in for the double as a whole: these checks assert only that
// an option reached it.
func tableStubDouble() stub.Double[TableNamesCall] {
	instance := func(s *TableStub) stub.Instance[TableNamesCall] {
		return stub.Instance[TableNamesCall]{
			Stub: s.OnNames.MethodStub,
			Call: func() {
				_ = s.Names()
			},
			Reset: s.ResetCalls,
		}
	}
	return stub.Double[TableNamesCall]{
		New: func(tb assert.TB) stub.Instance[TableNamesCall] {
			return instance(NewTableStub(tb))
		},
		WithClock: func(tb assert.TB, clk assert.Clock) stub.Instance[TableNamesCall] {
			return instance(NewTableStub(tb, TableStubWithClock(clk)))
		},
		WithLatency: func(tb assert.TB, d time.Duration, clk assert.Clock) stub.Instance[TableNamesCall] {
			return instance(NewTableStubDelayed(tb, d, clk))
		},
		WithRandSource: func(tb assert.TB, src rand.Source) stub.Instance[TableNamesCall] {
			return instance(NewTableStub(tb, TableStubWithRandSource(src)))
		},
		BenchMode: func(tb assert.TB) stub.Instance[TableNamesCall] {
			return instance(NewTableStub(tb, TableStubBenchMode()))
		},
		Strict: func(tb assert.TB) stub.Instance[TableNamesCall] {
			return instance(NewTableStub(tb, TableStubStrict()))
		},
	}
}

// TestTableStubOptions pins the settings that apply to the double rather
// than to one method: bench mode, the injected clock and random source,
// and what a reset clears.
func TestTableStubOptions(t *testing.T) {
	t.Parallel()

	stub.DoubleBehaviour(t, tableStubDouble())
}

// TestTableStubDelegateTo pins that a wrapped implementation is called
// through, and recorded on the way past.
func TestTableStubDelegateTo(t *testing.T) {
	t.Parallel()

	inner := NewTableStub(t)
	s := NewTableStub(t, TableStubDelegateTo(inner))

	t.Run("forwards Names to the wrapped implementation", func(t *testing.T) {
		_ = s.Names()
		inner.OnNames.AssertCalledOnce(t, "Names must reach the wrapped implementation")
	})

	t.Run("forwards Rows to the wrapped implementation", func(t *testing.T) {
		var a0 int
		a0 = 42
		_ = s.Rows(a0)
		inner.OnRows.AssertCalledOnce(t, "Rows must reach the wrapped implementation")
	})

	t.Run("forwards Counts to the wrapped implementation", func(t *testing.T) {
		_ = s.Counts()
		inner.OnCounts.AssertCalledOnce(t, "Counts must reach the wrapped implementation")
	})

	t.Run("forwards Tail to the wrapped implementation", func(t *testing.T) {
		_ = s.Tail()
		inner.OnTail.AssertCalledOnce(t, "Tail must reach the wrapped implementation")
	})
}
