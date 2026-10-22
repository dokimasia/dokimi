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

// tableStubNames calls Names of s with the samples of its arguments.
func tableStubNames(tb assert.TB, s *TableStub) iter.Seq[string] {
	tb.Helper()
	return s.Names()
}

// tableStubNamesPins returns the samples of the results of Names.
func tableStubNamesPins() iter.Seq[string] {
	var want0 iter.Seq[string]
	want0 = iter.Seq[string](func(func(string) bool) {})
	return want0
}

// tableStubRowsArgs returns the samples of the arguments of Rows.
func tableStubRowsArgs(tb assert.TB) int {
	tb.Helper()
	var a0 int
	a0 = 42
	return a0
}

// tableStubRows calls Rows of s with the samples of its arguments.
func tableStubRows(tb assert.TB, s *TableStub) iter.Seq2[Item, error] {
	tb.Helper()
	a0 := tableStubRowsArgs(tb)
	return s.Rows(a0)
}

// tableStubRowsPins returns the samples of the results of Rows.
func tableStubRowsPins() iter.Seq2[Item, error] {
	var want0 iter.Seq2[Item, error]
	want0 = iter.Seq2[Item, error](func(func(Item, error) bool) {})
	return want0
}

// tableStubCounts calls Counts of s with the samples of its arguments.
func tableStubCounts(tb assert.TB, s *TableStub) iter.Seq2[string, int] {
	tb.Helper()
	return s.Counts()
}

// tableStubCountsPins returns the samples of the results of Counts.
func tableStubCountsPins() iter.Seq2[string, int] {
	var want0 iter.Seq2[string, int]
	want0 = iter.Seq2[string, int](func(func(string, int) bool) {})
	return want0
}

// tableStubTail calls Tail of s with the samples of its arguments.
func tableStubTail(tb assert.TB, s *TableStub) Lines {
	tb.Helper()
	return s.Tail()
}

// tableStubTailPins returns the samples of the results of Tail.
func tableStubTailPins() Lines {
	var want0 Lines
	want0 = iter.Seq[string](func(func(string) bool) {})
	return want0
}

// tableStubNamesSubject returns the subject of [stub.Behaviour] for Names.
// The subject calls a new double that reports to tb.
func tableStubNamesSubject(tb assert.TB) stub.Subject[TableNamesCall, TableNamesReturn] {
	tb.Helper()
	s := NewTableStub(tb)
	return stub.Subject[TableNamesCall, TableNamesReturn]{
		Stub: s.OnNames.MethodStub,
		Call: func() {
			_ = tableStubNames(tb, s)
		},
		Result: func() TableNamesReturn {
			got0 := tableStubNames(tb, s)
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

// tableStubRowsSubject returns the subject of [stub.Behaviour] for Rows.
// The subject calls a new double that reports to tb.
func tableStubRowsSubject(tb assert.TB) stub.Subject[TableRowsCall, TableRowsReturn] {
	tb.Helper()
	s := NewTableStub(tb)
	return stub.Subject[TableRowsCall, TableRowsReturn]{
		Stub: s.OnRows.MethodStub,
		Call: func() {
			_ = tableStubRows(tb, s)
		},
		Result: func() TableRowsReturn {
			got0 := tableStubRows(tb, s)
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

// tableStubCountsSubject returns the subject of [stub.Behaviour] for
// Counts. The subject calls a new double that reports to tb.
func tableStubCountsSubject(tb assert.TB) stub.Subject[TableCountsCall, TableCountsReturn] {
	tb.Helper()
	s := NewTableStub(tb)
	return stub.Subject[TableCountsCall, TableCountsReturn]{
		Stub: s.OnCounts.MethodStub,
		Call: func() {
			_ = tableStubCounts(tb, s)
		},
		Result: func() TableCountsReturn {
			got0 := tableStubCounts(tb, s)
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

// tableStubTailSubject returns the subject of [stub.Behaviour] for Tail.
// The subject calls a new double that reports to tb.
func tableStubTailSubject(tb assert.TB) stub.Subject[TableTailCall, TableTailReturn] {
	tb.Helper()
	s := NewTableStub(tb)
	return stub.Subject[TableTailCall, TableTailReturn]{
		Stub: s.OnTail.MethodStub,
		Call: func() {
			_ = tableStubTail(tb, s)
		},
		Result: func() TableTailReturn {
			got0 := tableStubTail(tb, s)
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

// tableStubInstance is the instance of a TableStub that the checks of the
// constructors and the options drive.
type tableStubInstance = stub.Instance[TableNamesCall]

// tableStubInstanceOf returns the instance of s. The options apply to the
// whole double, so the checks call Names alone.
func tableStubInstanceOf(tb assert.TB, s *TableStub) tableStubInstance {
	tb.Helper()
	return tableStubInstance{
		Stub: s.OnNames.MethodStub,
		Call: func() {
			_ = tableStubNames(tb, s)
		},
		Reset: s.ResetCalls,
	}
}

// TestTableStub checks the methods, the constructors and the options of
// the double TableStub.
func TestTableStub(t *testing.T) {
	t.Parallel()

	t.Run("Names", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Names", tableStubNamesSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			var seen []TableNamesCall
			s.OnNames.OnRecord(func(c TableNamesCall) { seen = append(seen, c) })
			_ = tableStubNames(t, s)
			_ = tableStubNames(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Names")
		})
	})

	t.Run("Rows", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Rows", tableStubRowsSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			_ = tableStubRows(t, s)
			got := s.OnRows.AssertCalledOnce(t, "Rows records the call")
			a0 := tableStubRowsArgs(t)
			assert.Equal(t, got.Limit, a0, "the call records Limit")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			var seen []TableRowsCall
			s.OnRows.OnRecord(func(c TableRowsCall) { seen = append(seen, c) })
			_ = tableStubRows(t, s)
			_ = tableStubRows(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Rows")
		})
	})

	t.Run("Counts", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Counts", tableStubCountsSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			var seen []TableCountsCall
			s.OnCounts.OnRecord(func(c TableCountsCall) { seen = append(seen, c) })
			_ = tableStubCounts(t, s)
			_ = tableStubCounts(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Counts")
		})
	})

	t.Run("Tail", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Tail", tableStubTailSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			var seen []TableTailCall
			s.OnTail.OnRecord(func(c TableTailCall) { seen = append(seen, c) })
			_ = tableStubTail(t, s)
			_ = tableStubTail(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Tail")
		})
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()
		stub.ResetBehaviour(t, func(tb assert.TB) tableStubInstance {
			return tableStubInstanceOf(tb, NewTableStub(tb))
		})

		t.Run("clears the recorded calls of Names", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			_ = tableStubNames(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnNames.CallCount(), 0, "ResetCalls clears the calls of Names")
		})

		t.Run("keeps the results of Returns of Names", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			want0 := tableStubNamesPins()
			s.OnNames.Returns(want0)
			s.ResetCalls()
			got0 := tableStubNames(t, s)
			assert.Equal(t, got0, want0, "ResetCalls keeps the Result of Names")
		})

		t.Run("clears the recorded calls of Rows", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			_ = tableStubRows(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnRows.CallCount(), 0, "ResetCalls clears the calls of Rows")
		})

		t.Run("keeps the results of Returns of Rows", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			want0 := tableStubRowsPins()
			s.OnRows.Returns(want0)
			s.ResetCalls()
			got0 := tableStubRows(t, s)
			assert.Equal(t, got0, want0, "ResetCalls keeps the Result of Rows")
		})

		t.Run("clears the recorded calls of Counts", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			_ = tableStubCounts(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnCounts.CallCount(), 0, "ResetCalls clears the calls of Counts")
		})

		t.Run("keeps the results of Returns of Counts", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			want0 := tableStubCountsPins()
			s.OnCounts.Returns(want0)
			s.ResetCalls()
			got0 := tableStubCounts(t, s)
			assert.Equal(t, got0, want0, "ResetCalls keeps the Result of Counts")
		})

		t.Run("clears the recorded calls of Tail", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			_ = tableStubTail(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnTail.CallCount(), 0, "ResetCalls clears the calls of Tail")
		})

		t.Run("keeps the results of Returns of Tail", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			want0 := tableStubTailPins()
			s.OnTail.Returns(want0)
			s.ResetCalls()
			got0 := tableStubTail(t, s)
			assert.Equal(t, got0, want0, "ResetCalls keeps the Result of Tail")
		})
	})

	t.Run("NewTableStubDelayed", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) tableStubInstance {
				return tableStubInstanceOf(tb, NewTableStubDelayed(tb, d, clk))
			},
		)
	})

	t.Run("TableStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) tableStubInstance {
			return tableStubInstanceOf(tb, NewTableStub(tb, TableStubStrict()))
		})
	})

	t.Run("TableStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Names to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewTableStub(t)
			s := NewTableStub(t, TableStubDelegateTo(inner))
			_ = tableStubNames(t, s)
			inner.OnNames.AssertCalledOnce(t, "Names calls the implementation")
		})

		t.Run("passes Rows to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewTableStub(t)
			s := NewTableStub(t, TableStubDelegateTo(inner))
			_ = tableStubRows(t, s)
			inner.OnRows.AssertCalledOnce(t, "Rows calls the implementation")
		})

		t.Run("passes Counts to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewTableStub(t)
			s := NewTableStub(t, TableStubDelegateTo(inner))
			_ = tableStubCounts(t, s)
			inner.OnCounts.AssertCalledOnce(t, "Counts calls the implementation")
		})

		t.Run("passes Tail to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewTableStub(t)
			s := NewTableStub(t, TableStubDelegateTo(inner))
			_ = tableStubTail(t, s)
			inner.OnTail.AssertCalledOnce(t, "Tail calls the implementation")
		})
	})

	t.Run("TableStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) tableStubInstance {
			return tableStubInstanceOf(tb, NewTableStub(tb, TableStubWithClock(clk)))
		})
	})

	t.Run("TableStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) tableStubInstance {
			return tableStubInstanceOf(tb, NewTableStub(tb, TableStubWithRandSource(src)))
		})
	})

	t.Run("TableStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) tableStubInstance {
			return tableStubInstanceOf(tb, NewTableStub(tb, TableStubBenchMode()))
		})
	})

	t.Run("TableStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) tableStubInstance {
				return tableStubInstanceOf(
					tb,
					NewTableStub(tb, TableStubWithClock(clk), TableStubWithLatency(d)),
				)
			},
		)
	})

	t.Run("WithTableNames", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Names", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewTableStub(t, WithTableNames(func() iter.Seq[string] {
				called = true
				var z0 iter.Seq[string]
				return z0
			}))
			_ = tableStubNames(t, s)
			assert.True(t, called, "WithTableNames sets the function of Names")
		})
	})

	t.Run("WithTableRows", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Rows", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewTableStub(t, WithTableRows(func(_ int) iter.Seq2[Item, error] {
				called = true
				var z0 iter.Seq2[Item, error]
				return z0
			}))
			_ = tableStubRows(t, s)
			assert.True(t, called, "WithTableRows sets the function of Rows")
		})
	})

	t.Run("WithTableCounts", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Counts", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewTableStub(t, WithTableCounts(func() iter.Seq2[string, int] {
				called = true
				var z0 iter.Seq2[string, int]
				return z0
			}))
			_ = tableStubCounts(t, s)
			assert.True(t, called, "WithTableCounts sets the function of Counts")
		})
	})

	t.Run("WithTableTail", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Tail", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewTableStub(t, WithTableTail(func() Lines {
				called = true
				var z0 Lines
				return z0
			}))
			_ = tableStubTail(t, s)
			assert.True(t, called, "WithTableTail sets the function of Tail")
		})
	})
}

// TestTableNamesStub checks the methods of TableNamesStub, which
// configures Names.
func TestTableNamesStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Names", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			want0 := tableStubNamesPins()
			s.OnNames.Returns(want0)
			got0 := tableStubNames(t, s)
			assert.Equal(t, got0, want0, "Names returns the Result that Returns sets")
		})

		t.Run("configures Names for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewTableStub(r)
			s.OnNames.Strict()
			var want0 iter.Seq[string]
			s.OnNames.Returns(want0)
			_ = tableStubNames(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Names")
		})
	})

	t.Run("Yields", func(t *testing.T) {
		t.Parallel()

		t.Run("yields each item", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			var want string
			s.OnNames.Yields(want, want)
			got0 := tableStubNames(t, s)
			seen := 0
			for range got0 {
				seen++
			}
			assert.Equal(t, seen, 2, "Names yields each item")
		})

		t.Run("stops yielding when the caller breaks", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			var want string
			s.OnNames.Yields(want, want, want)
			got0 := tableStubNames(t, s)
			seen := 0
			for range got0 {
				seen++
				break
			}
			assert.Equal(t, seen, 1, "Names stops when the caller breaks")
		})
	})
}

// TestTableRowsStub checks the methods of TableRowsStub, which configures
// Rows.
func TestTableRowsStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Rows", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			want0 := tableStubRowsPins()
			s.OnRows.Returns(want0)
			got0 := tableStubRows(t, s)
			assert.Equal(t, got0, want0, "Rows returns the Result that Returns sets")
		})

		t.Run("configures Rows for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewTableStub(r)
			s.OnRows.Strict()
			var want0 iter.Seq2[Item, error]
			s.OnRows.Returns(want0)
			_ = tableStubRows(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Rows")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Rows to the function", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			var got0 int
			s.OnRows.Func(func(p0 int) iter.Seq2[Item, error] {
				got0 = p0
				var z0 iter.Seq2[Item, error]
				return z0
			})
			_ = tableStubRows(t, s)
			a0 := tableStubRowsArgs(t)
			assert.Equal(t, got0, a0, "the function receives Limit")
		})
	})

	t.Run("Yields", func(t *testing.T) {
		t.Parallel()

		t.Run("yields each item", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			var want Item
			s.OnRows.Yields(want, want)
			got0 := tableStubRows(t, s)
			seen := 0
			for range got0 {
				seen++
			}
			assert.Equal(t, seen, 2, "Rows yields each item")
		})

		t.Run("stops yielding when the caller breaks", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			var want Item
			s.OnRows.Yields(want, want, want)
			got0 := tableStubRows(t, s)
			seen := 0
			for range got0 {
				seen++
				break
			}
			assert.Equal(t, seen, 1, "Rows stops when the caller breaks")
		})
	})

	t.Run("YieldsError", func(t *testing.T) {
		t.Parallel()

		t.Run("yields the error after the items", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			var item Item
			want := errors.New("Rows-yields")
			s.OnRows.YieldsError([]Item{item, item}, want)
			got0 := tableStubRows(t, s)
			seen := 0
			var last error
			for _, err := range got0 {
				if err != nil {
					last = err
					break
				}
				seen++
			}
			assert.Equal(t, seen, 2, "Rows yields each item before the error")
			assert.ErrorIs(t, last, want, "Rows ends with the error of YieldsError")
		})

		t.Run("stops the sequence when the caller breaks", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			var item Item
			s.OnRows.YieldsError([]Item{item, item, item}, errors.New("Rows-yields"))
			got0 := tableStubRows(t, s)
			seen := 0
			for range got0 {
				seen++
				break
			}
			assert.Equal(t, seen, 1, "Rows stops when the caller breaks")
		})
	})
}

// TestTableCountsStub checks the methods of TableCountsStub, which
// configures Counts.
func TestTableCountsStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Counts", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			want0 := tableStubCountsPins()
			s.OnCounts.Returns(want0)
			got0 := tableStubCounts(t, s)
			assert.Equal(t, got0, want0, "Counts returns the Result that Returns sets")
		})

		t.Run("configures Counts for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewTableStub(r)
			s.OnCounts.Strict()
			var want0 iter.Seq2[string, int]
			s.OnCounts.Returns(want0)
			_ = tableStubCounts(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Counts")
		})
	})

	t.Run("Yields", func(t *testing.T) {
		t.Parallel()

		t.Run("yields each item", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			var want string
			s.OnCounts.Yields(want, want)
			got0 := tableStubCounts(t, s)
			seen := 0
			for range got0 {
				seen++
			}
			assert.Equal(t, seen, 2, "Counts yields each item")
		})

		t.Run("stops yielding when the caller breaks", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			var want string
			s.OnCounts.Yields(want, want, want)
			got0 := tableStubCounts(t, s)
			seen := 0
			for range got0 {
				seen++
				break
			}
			assert.Equal(t, seen, 1, "Counts stops when the caller breaks")
		})
	})
}

// TestTableTailStub checks the methods of TableTailStub, which configures
// Tail.
func TestTableTailStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Tail", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			want0 := tableStubTailPins()
			s.OnTail.Returns(want0)
			got0 := tableStubTail(t, s)
			assert.Equal(t, got0, want0, "Tail returns the Result that Returns sets")
		})

		t.Run("configures Tail for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewTableStub(r)
			s.OnTail.Strict()
			var want0 Lines
			s.OnTail.Returns(want0)
			_ = tableStubTail(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Tail")
		})
	})

	t.Run("Yields", func(t *testing.T) {
		t.Parallel()

		t.Run("yields each item", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			var want string
			s.OnTail.Yields(want, want)
			got0 := tableStubTail(t, s)
			seen := 0
			for range got0 {
				seen++
			}
			assert.Equal(t, seen, 2, "Tail yields each item")
		})

		t.Run("stops yielding when the caller breaks", func(t *testing.T) {
			t.Parallel()
			s := NewTableStub(t)
			var want string
			s.OnTail.Yields(want, want, want)
			got0 := tableStubTail(t, s)
			seen := 0
			for range got0 {
				seen++
				break
			}
			assert.Equal(t, seen, 1, "Tail stops when the caller breaks")
		})
	})
}
