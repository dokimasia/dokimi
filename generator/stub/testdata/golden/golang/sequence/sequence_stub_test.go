package sequence_test

import (
	"errors"
	"iter"
	rand "math/rand/v2"
	"testing"
	"time"

	"example.com/shop/sequence"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// tableStubNames calls Names of s with the samples of its arguments.
func tableStubNames(tb assert.TB, s *sequence.TableStub) iter.Seq[string] {
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
func tableStubRows(tb assert.TB, s *sequence.TableStub) iter.Seq2[sequence.Item, error] {
	tb.Helper()
	a0 := tableStubRowsArgs(tb)
	return s.Rows(a0)
}

// tableStubRowsPins returns the samples of the results of Rows.
func tableStubRowsPins() iter.Seq2[sequence.Item, error] {
	var want0 iter.Seq2[sequence.Item, error]
	want0 = iter.Seq2[sequence.Item, error](func(func(sequence.Item, error) bool) {})
	return want0
}

// tableStubCounts calls Counts of s with the samples of its arguments.
func tableStubCounts(tb assert.TB, s *sequence.TableStub) iter.Seq2[string, int] {
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
func tableStubTail(tb assert.TB, s *sequence.TableStub) sequence.Lines {
	tb.Helper()
	return s.Tail()
}

// tableStubTailPins returns the samples of the results of Tail.
func tableStubTailPins() sequence.Lines {
	var want0 sequence.Lines
	want0 = iter.Seq[string](func(func(string) bool) {})
	return want0
}

// tableStubNamesSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Names. The subject calls a new double that
// reports to tb.
func tableStubNamesSubject(
	tb assert.TB,
) stub.Subject[sequence.TableNamesCall, sequence.TableNamesReturn] {
	tb.Helper()
	s := sequence.NewTableStub(tb)
	return stub.Subject[sequence.TableNamesCall, sequence.TableNamesReturn]{
		Stub:  s.OnNames.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = tableStubNames(tb, s)
		},
		Result: func() sequence.TableNamesReturn {
			got0 := tableStubNames(tb, s)
			return sequence.TableNamesReturn{Result: got0}
		},
		Returns: func() sequence.TableNamesReturn {
			want0 := tableStubNamesPins()
			s.OnNames.Returns(want0)
			return sequence.TableNamesReturn{Result: want0}
		},
		Override: func(seen func(sequence.TableNamesCall)) {
			s.OnNames.Func(func() iter.Seq[string] {
				seen(sequence.TableNamesCall{})
				var z0 iter.Seq[string]
				return z0
			})
		},
	}
}

// tableStubRowsSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Rows. The subject calls a new double that
// reports to tb.
func tableStubRowsSubject(
	tb assert.TB,
) stub.Subject[sequence.TableRowsCall, sequence.TableRowsReturn] {
	tb.Helper()
	s := sequence.NewTableStub(tb)
	return stub.Subject[sequence.TableRowsCall, sequence.TableRowsReturn]{
		Stub:  s.OnRows.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = tableStubRows(tb, s)
		},
		Arguments: func(tb assert.TB, c sequence.TableRowsCall) {
			tb.Helper()
			a0 := tableStubRowsArgs(tb)
			assert.Equal(tb, c.Limit, a0, "the call records Limit")
		},
		Result: func() sequence.TableRowsReturn {
			got0 := tableStubRows(tb, s)
			return sequence.TableRowsReturn{Result: got0}
		},
		Returns: func() sequence.TableRowsReturn {
			want0 := tableStubRowsPins()
			s.OnRows.Returns(want0)
			return sequence.TableRowsReturn{Result: want0}
		},
		Override: func(seen func(sequence.TableRowsCall)) {
			s.OnRows.Func(func(p0 int) iter.Seq2[sequence.Item, error] {
				seen(sequence.TableRowsCall{Limit: p0})
				var z0 iter.Seq2[sequence.Item, error]
				return z0
			})
		},
	}
}

// tableStubCountsSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Counts. The subject calls a new double that
// reports to tb.
func tableStubCountsSubject(
	tb assert.TB,
) stub.Subject[sequence.TableCountsCall, sequence.TableCountsReturn] {
	tb.Helper()
	s := sequence.NewTableStub(tb)
	return stub.Subject[sequence.TableCountsCall, sequence.TableCountsReturn]{
		Stub:  s.OnCounts.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = tableStubCounts(tb, s)
		},
		Result: func() sequence.TableCountsReturn {
			got0 := tableStubCounts(tb, s)
			return sequence.TableCountsReturn{Result: got0}
		},
		Returns: func() sequence.TableCountsReturn {
			want0 := tableStubCountsPins()
			s.OnCounts.Returns(want0)
			return sequence.TableCountsReturn{Result: want0}
		},
		Override: func(seen func(sequence.TableCountsCall)) {
			s.OnCounts.Func(func() iter.Seq2[string, int] {
				seen(sequence.TableCountsCall{})
				var z0 iter.Seq2[string, int]
				return z0
			})
		},
	}
}

// tableStubTailSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Tail. The subject calls a new double that
// reports to tb.
func tableStubTailSubject(
	tb assert.TB,
) stub.Subject[sequence.TableTailCall, sequence.TableTailReturn] {
	tb.Helper()
	s := sequence.NewTableStub(tb)
	return stub.Subject[sequence.TableTailCall, sequence.TableTailReturn]{
		Stub:  s.OnTail.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = tableStubTail(tb, s)
		},
		Result: func() sequence.TableTailReturn {
			got0 := tableStubTail(tb, s)
			return sequence.TableTailReturn{Result: got0}
		},
		Returns: func() sequence.TableTailReturn {
			want0 := tableStubTailPins()
			s.OnTail.Returns(want0)
			return sequence.TableTailReturn{Result: want0}
		},
		Override: func(seen func(sequence.TableTailCall)) {
			s.OnTail.Func(func() sequence.Lines {
				seen(sequence.TableTailCall{})
				var z0 sequence.Lines
				return z0
			})
		},
	}
}

// tableStubInstance is the instance of a sequence.TableStub that the
// checks of the constructors and the options drive.
type tableStubInstance = stub.Instance[sequence.TableNamesCall]

// tableStubInstanceOf returns the instance of s. The options apply to the
// whole double, so the checks call Names alone.
func tableStubInstanceOf(tb assert.TB, s *sequence.TableStub) tableStubInstance {
	tb.Helper()
	return tableStubInstance{
		Stub: s.OnNames.MethodStub,
		Call: func() {
			_ = tableStubNames(tb, s)
		},
	}
}

// TestTableStub checks the methods, the constructors and the options of
// the double TableStub.
func TestTableStub(t *testing.T) {
	t.Parallel()

	t.Run("Names", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Names", tableStubNamesSubject)
	})

	t.Run("Rows", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Rows", tableStubRowsSubject)
	})

	t.Run("Counts", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Counts", tableStubCountsSubject)
	})

	t.Run("Tail", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Tail", tableStubTailSubject)
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()

		t.Run("Names", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Names", tableStubNamesSubject)
		})

		t.Run("Rows", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Rows", tableStubRowsSubject)
		})

		t.Run("Counts", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Counts", tableStubCountsSubject)
		})

		t.Run("Tail", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Tail", tableStubTailSubject)
		})
	})

	t.Run("TableStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) tableStubInstance {
			return tableStubInstanceOf(tb, sequence.NewTableStub(tb, sequence.TableStubStrict()))
		})
	})

	t.Run("TableStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Names to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := sequence.NewTableStub(t)
			s := sequence.NewTableStub(t, sequence.TableStubDelegateTo(inner))
			_ = tableStubNames(t, s)
			inner.OnNames.AssertCalledOnce(t, "Names calls the implementation")
		})

		t.Run("passes Rows to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := sequence.NewTableStub(t)
			s := sequence.NewTableStub(t, sequence.TableStubDelegateTo(inner))
			_ = tableStubRows(t, s)
			inner.OnRows.AssertCalledOnce(t, "Rows calls the implementation")
		})

		t.Run("passes Counts to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := sequence.NewTableStub(t)
			s := sequence.NewTableStub(t, sequence.TableStubDelegateTo(inner))
			_ = tableStubCounts(t, s)
			inner.OnCounts.AssertCalledOnce(t, "Counts calls the implementation")
		})

		t.Run("passes Tail to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := sequence.NewTableStub(t)
			s := sequence.NewTableStub(t, sequence.TableStubDelegateTo(inner))
			_ = tableStubTail(t, s)
			inner.OnTail.AssertCalledOnce(t, "Tail calls the implementation")
		})
	})

	t.Run("TableStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) tableStubInstance {
			return tableStubInstanceOf(
				tb,
				sequence.NewTableStub(tb, sequence.TableStubWithClock(clk)),
			)
		})
	})

	t.Run("TableStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) tableStubInstance {
			return tableStubInstanceOf(
				tb,
				sequence.NewTableStub(tb, sequence.TableStubWithRandSource(src)),
			)
		})
	})

	t.Run("TableStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) tableStubInstance {
			return tableStubInstanceOf(tb, sequence.NewTableStub(tb, sequence.TableStubBenchMode()))
		})
	})

	t.Run("TableStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) tableStubInstance {
				return tableStubInstanceOf(
					tb,
					sequence.NewTableStub(
						tb,
						sequence.TableStubWithClock(clk),
						sequence.TableStubWithLatency(d),
					),
				)
			},
		)
	})

	t.Run("WithTableNames", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Names", func(t *testing.T) {
			t.Parallel()
			called := false
			s := sequence.NewTableStub(t, sequence.WithTableNames(func() iter.Seq[string] {
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
			s := sequence.NewTableStub(
				t,
				sequence.WithTableRows(func(_ int) iter.Seq2[sequence.Item, error] {
					called = true
					var z0 iter.Seq2[sequence.Item, error]
					return z0
				}),
			)
			_ = tableStubRows(t, s)
			assert.True(t, called, "WithTableRows sets the function of Rows")
		})
	})

	t.Run("WithTableCounts", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Counts", func(t *testing.T) {
			t.Parallel()
			called := false
			s := sequence.NewTableStub(t, sequence.WithTableCounts(func() iter.Seq2[string, int] {
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
			s := sequence.NewTableStub(t, sequence.WithTableTail(func() sequence.Lines {
				called = true
				var z0 sequence.Lines
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

	t.Run("Yields", func(t *testing.T) {
		t.Parallel()

		t.Run("yields each item", func(t *testing.T) {
			t.Parallel()
			s := sequence.NewTableStub(t)
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
			s := sequence.NewTableStub(t)
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

	t.Run("Yields", func(t *testing.T) {
		t.Parallel()

		t.Run("yields each item", func(t *testing.T) {
			t.Parallel()
			s := sequence.NewTableStub(t)
			var want sequence.Item
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
			s := sequence.NewTableStub(t)
			var want sequence.Item
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
			s := sequence.NewTableStub(t)
			var item sequence.Item
			want := errors.New("Rows-yields")
			s.OnRows.YieldsError([]sequence.Item{item, item}, want)
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
			s := sequence.NewTableStub(t)
			var item sequence.Item
			s.OnRows.YieldsError([]sequence.Item{item, item, item}, errors.New("Rows-yields"))
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

	t.Run("Yields", func(t *testing.T) {
		t.Parallel()

		t.Run("yields each item", func(t *testing.T) {
			t.Parallel()
			s := sequence.NewTableStub(t)
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
			s := sequence.NewTableStub(t)
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

	t.Run("Yields", func(t *testing.T) {
		t.Parallel()

		t.Run("yields each item", func(t *testing.T) {
			t.Parallel()
			s := sequence.NewTableStub(t)
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
			s := sequence.NewTableStub(t)
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
