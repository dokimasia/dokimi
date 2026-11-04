package built_test

import (
	"context"
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"example.com/shop/billing"
	"example.com/shop/built"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/stub"
)

// storeStubGetArgs returns the samples of the arguments of Get.
func storeStubGetArgs(tb assert.TB) (context.Context, string) {
	tb.Helper()
	a0 := assert.Context(tb)
	var a1 string
	a1 = "Get-a"
	return a0, a1
}

// storeStubGet calls Get of s with the samples of its arguments.
func storeStubGet(tb assert.TB, s *built.StoreStub) (built.Account, error) {
	tb.Helper()
	a0, a1 := storeStubGetArgs(tb)
	return s.Get(a0, a1)
}

// storeStubGetPins returns the samples of the results of Get, with a new
// error as its error result.
func storeStubGetPins() (built.Account, error) {
	var want0 built.Account
	want1 := errors.New("Get-returns")
	want0 = built.Account{ID: "ID-a"}
	return want0, want1
}

// storeStubFindArgs returns the samples of the arguments of Find.
func storeStubFindArgs(tb assert.TB) (context.Context, string) {
	tb.Helper()
	a0 := assert.Context(tb)
	var a1 string
	a1 = "Find-a"
	return a0, a1
}

// storeStubFind calls Find of s with the samples of its arguments.
func storeStubFind(tb assert.TB, s *built.StoreStub) (*built.Account, error) {
	tb.Helper()
	a0, a1 := storeStubFindArgs(tb)
	return s.Find(a0, a1)
}

// storeStubFindPins returns the samples of the results of Find, with a new
// error as its error result.
func storeStubFindPins() (*built.Account, error) {
	var want0 *built.Account
	want1 := errors.New("Find-returns")
	want0 = &built.Account{ID: "ID-a"}
	return want0, want1
}

// storeStubLookupArgs returns the samples of the arguments of Lookup.
func storeStubLookupArgs(tb assert.TB) string {
	tb.Helper()
	var a0 string
	a0 = "Lookup-a"
	return a0
}

// storeStubLookup calls Lookup of s with the samples of its arguments.
func storeStubLookup(tb assert.TB, s *built.StoreStub) (built.Account, bool) {
	tb.Helper()
	a0 := storeStubLookupArgs(tb)
	return s.Lookup(a0)
}

// storeStubLookupPins returns the samples of the results of Lookup.
func storeStubLookupPins() (built.Account, bool) {
	var want0 built.Account
	var want1 bool
	want0 = built.Account{ID: "ID-a"}
	want1 = true
	return want0, want1
}

// storeStubListArgs returns the samples of the arguments of List.
func storeStubListArgs(tb assert.TB) context.Context {
	tb.Helper()
	a0 := assert.Context(tb)
	return a0
}

// storeStubList calls List of s with the samples of its arguments.
func storeStubList(tb assert.TB, s *built.StoreStub) (built.Page[built.Account], error) {
	tb.Helper()
	a0 := storeStubListArgs(tb)
	return s.List(a0)
}

// storeStubListPins returns the samples of the results of List, with a new
// error as its error result.
func storeStubListPins() (built.Page[built.Account], error) {
	var want0 built.Page[built.Account]
	want1 := errors.New("List-returns")
	want0 = built.Page[built.Account]{Items: []built.Account{{ID: "ID-a"}}}
	return want0, want1
}

// storeStubNoteArgs returns the samples of the arguments of Note.
func storeStubNoteArgs(tb assert.TB) string {
	tb.Helper()
	var a0 string
	a0 = "Note-a"
	return a0
}

// storeStubNote calls Note of s with the samples of its arguments.
func storeStubNote(tb assert.TB, s *built.StoreStub) built.Note {
	tb.Helper()
	a0 := storeStubNoteArgs(tb)
	return s.Note(a0)
}

// storeStubNotePins returns the samples of the results of Note.
func storeStubNotePins() built.Note {
	var want0 built.Note
	want0 = built.Note{Text: "Text-a"}
	return want0
}

// storeStubBillArgs returns the samples of the arguments of Bill.
func storeStubBillArgs(tb assert.TB) string {
	tb.Helper()
	var a0 string
	a0 = "Bill-a"
	return a0
}

// storeStubBill calls Bill of s with the samples of its arguments.
func storeStubBill(tb assert.TB, s *built.StoreStub) billing.Invoice {
	tb.Helper()
	a0 := storeStubBillArgs(tb)
	return s.Bill(a0)
}

// storeStubBillPins returns the samples of the results of Bill.
func storeStubBillPins() billing.Invoice {
	var want0 billing.Invoice
	want0 = billing.Invoice{Number: "Number-a"}
	return want0
}

// storeStubGetSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Get. The subject calls a new double that
// reports to tb.
func storeStubGetSubject(tb assert.TB) stub.Subject[built.StoreGetCall, built.StoreGetReturn] {
	tb.Helper()
	s := built.NewStoreStub(tb)
	return stub.Subject[built.StoreGetCall, built.StoreGetReturn]{
		Stub:  s.OnGet.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = storeStubGet(tb, s)
		},
		Arguments: func(tb assert.TB, c built.StoreGetCall) {
			tb.Helper()
			a0, a1 := storeStubGetArgs(tb)
			expect.Equal(tb, c.Ctx, a0, "the call records Ctx", expect.ByIdentity())
			expect.Equal(tb, c.ID, a1, "the call records ID")
		},
		Result: func() built.StoreGetReturn {
			got0, got1 := storeStubGet(tb, s)
			return built.StoreGetReturn{Result: got0, Err: got1}
		},
		Default: func() built.StoreGetReturn {
			r0 := built.NewAccountBuilder().Build()
			return built.StoreGetReturn{Result: r0}
		},
		Returns: func() built.StoreGetReturn {
			want0, want1 := storeStubGetPins()
			s.OnGet.Returns(want0, want1)
			return built.StoreGetReturn{Result: want0, Err: want1}
		},
		Override: func(seen func(built.StoreGetCall)) {
			s.OnGet.Func(func(p0 context.Context, p1 string) (built.Account, error) {
				seen(built.StoreGetCall{Ctx: p0, ID: p1})
				var z0 built.Account
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			_, r1 := storeStubGet(tb, s)
			return r1
		},
	}
}

// storeStubFindSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Find. The subject calls a new double that
// reports to tb.
func storeStubFindSubject(tb assert.TB) stub.Subject[built.StoreFindCall, built.StoreFindReturn] {
	tb.Helper()
	s := built.NewStoreStub(tb)
	return stub.Subject[built.StoreFindCall, built.StoreFindReturn]{
		Stub:  s.OnFind.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = storeStubFind(tb, s)
		},
		Arguments: func(tb assert.TB, c built.StoreFindCall) {
			tb.Helper()
			a0, a1 := storeStubFindArgs(tb)
			expect.Equal(tb, c.Ctx, a0, "the call records Ctx", expect.ByIdentity())
			expect.Equal(tb, c.ID, a1, "the call records ID")
		},
		Result: func() built.StoreFindReturn {
			got0, got1 := storeStubFind(tb, s)
			return built.StoreFindReturn{Result: got0, Err: got1}
		},
		Default: func() built.StoreFindReturn {
			r0 := built.NewAccountBuilder().Build()
			return built.StoreFindReturn{Result: &r0}
		},
		Returns: func() built.StoreFindReturn {
			want0, want1 := storeStubFindPins()
			s.OnFind.Returns(want0, want1)
			return built.StoreFindReturn{Result: want0, Err: want1}
		},
		Override: func(seen func(built.StoreFindCall)) {
			s.OnFind.Func(func(p0 context.Context, p1 string) (*built.Account, error) {
				seen(built.StoreFindCall{Ctx: p0, ID: p1})
				var z0 *built.Account
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			_, r1 := storeStubFind(tb, s)
			return r1
		},
	}
}

// storeStubLookupSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Lookup. The subject calls a new double that
// reports to tb.
func storeStubLookupSubject(
	tb assert.TB,
) stub.Subject[built.StoreLookupCall, built.StoreLookupReturn] {
	tb.Helper()
	s := built.NewStoreStub(tb)
	return stub.Subject[built.StoreLookupCall, built.StoreLookupReturn]{
		Stub:  s.OnLookup.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = storeStubLookup(tb, s)
		},
		Arguments: func(tb assert.TB, c built.StoreLookupCall) {
			tb.Helper()
			a0 := storeStubLookupArgs(tb)
			assert.Equal(tb, c.ID, a0, "the call records ID")
		},
		Result: func() built.StoreLookupReturn {
			got0, got1 := storeStubLookup(tb, s)
			return built.StoreLookupReturn{Result0: got0, Result1: got1}
		},
		Default: func() built.StoreLookupReturn {
			r0 := built.NewAccountBuilder().Build()
			return built.StoreLookupReturn{Result0: r0}
		},
		Returns: func() built.StoreLookupReturn {
			want0, want1 := storeStubLookupPins()
			s.OnLookup.Returns(want0, want1)
			return built.StoreLookupReturn{Result0: want0, Result1: want1}
		},
		Override: func(seen func(built.StoreLookupCall)) {
			s.OnLookup.Func(func(p0 string) (built.Account, bool) {
				seen(built.StoreLookupCall{ID: p0})
				var z0 built.Account
				var z1 bool
				return z0, z1
			})
		},
	}
}

// storeStubListSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for List. The subject calls a new double that
// reports to tb.
func storeStubListSubject(tb assert.TB) stub.Subject[built.StoreListCall, built.StoreListReturn] {
	tb.Helper()
	s := built.NewStoreStub(tb)
	return stub.Subject[built.StoreListCall, built.StoreListReturn]{
		Stub:  s.OnList.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = storeStubList(tb, s)
		},
		Arguments: func(tb assert.TB, c built.StoreListCall) {
			tb.Helper()
			a0 := storeStubListArgs(tb)
			assert.Equal(tb, c.Ctx, a0, "the call records Ctx", assert.ByIdentity())
		},
		Result: func() built.StoreListReturn {
			got0, got1 := storeStubList(tb, s)
			return built.StoreListReturn{Result: got0, Err: got1}
		},
		Default: func() built.StoreListReturn {
			r0 := built.NewPageBuilder[built.Account]().Build()
			return built.StoreListReturn{Result: r0}
		},
		Returns: func() built.StoreListReturn {
			want0, want1 := storeStubListPins()
			s.OnList.Returns(want0, want1)
			return built.StoreListReturn{Result: want0, Err: want1}
		},
		Override: func(seen func(built.StoreListCall)) {
			s.OnList.Func(func(p0 context.Context) (built.Page[built.Account], error) {
				seen(built.StoreListCall{Ctx: p0})
				var z0 built.Page[built.Account]
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			_, r1 := storeStubList(tb, s)
			return r1
		},
	}
}

// storeStubNoteSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Note. The subject calls a new double that
// reports to tb.
func storeStubNoteSubject(tb assert.TB) stub.Subject[built.StoreNoteCall, built.StoreNoteReturn] {
	tb.Helper()
	s := built.NewStoreStub(tb)
	return stub.Subject[built.StoreNoteCall, built.StoreNoteReturn]{
		Stub:  s.OnNote.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = storeStubNote(tb, s)
		},
		Arguments: func(tb assert.TB, c built.StoreNoteCall) {
			tb.Helper()
			a0 := storeStubNoteArgs(tb)
			assert.Equal(tb, c.ID, a0, "the call records ID")
		},
		Result: func() built.StoreNoteReturn {
			got0 := storeStubNote(tb, s)
			return built.StoreNoteReturn{Result: got0}
		},
		Returns: func() built.StoreNoteReturn {
			want0 := storeStubNotePins()
			s.OnNote.Returns(want0)
			return built.StoreNoteReturn{Result: want0}
		},
		Override: func(seen func(built.StoreNoteCall)) {
			s.OnNote.Func(func(p0 string) built.Note {
				seen(built.StoreNoteCall{ID: p0})
				var z0 built.Note
				return z0
			})
		},
	}
}

// storeStubBillSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Bill. The subject calls a new double that
// reports to tb.
func storeStubBillSubject(tb assert.TB) stub.Subject[built.StoreBillCall, built.StoreBillReturn] {
	tb.Helper()
	s := built.NewStoreStub(tb)
	return stub.Subject[built.StoreBillCall, built.StoreBillReturn]{
		Stub:  s.OnBill.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = storeStubBill(tb, s)
		},
		Arguments: func(tb assert.TB, c built.StoreBillCall) {
			tb.Helper()
			a0 := storeStubBillArgs(tb)
			assert.Equal(tb, c.ID, a0, "the call records ID")
		},
		Result: func() built.StoreBillReturn {
			got0 := storeStubBill(tb, s)
			return built.StoreBillReturn{Result: got0}
		},
		Returns: func() built.StoreBillReturn {
			want0 := storeStubBillPins()
			s.OnBill.Returns(want0)
			return built.StoreBillReturn{Result: want0}
		},
		Override: func(seen func(built.StoreBillCall)) {
			s.OnBill.Func(func(p0 string) billing.Invoice {
				seen(built.StoreBillCall{ID: p0})
				var z0 billing.Invoice
				return z0
			})
		},
	}
}

// storeStubInstance is the instance of a built.StoreStub that the checks
// of the constructors and the options drive.
type storeStubInstance = stub.Instance[built.StoreGetCall]

// storeStubInstanceOf returns the instance of s. The options apply to the
// whole double, so the checks call Get alone.
func storeStubInstanceOf(tb assert.TB, s *built.StoreStub) storeStubInstance {
	tb.Helper()
	return storeStubInstance{
		Stub: s.OnGet.MethodStub,
		Call: func() {
			_, _ = storeStubGet(tb, s)
		},
	}
}

// TestStoreStub checks the methods, the constructors and the options of
// the double StoreStub.
func TestStoreStub(t *testing.T) {
	t.Parallel()

	t.Run("Get", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Get", storeStubGetSubject)
	})

	t.Run("Find", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Find", storeStubFindSubject)
	})

	t.Run("Lookup", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Lookup", storeStubLookupSubject)
	})

	t.Run("List", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "List", storeStubListSubject)
	})

	t.Run("Note", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Note", storeStubNoteSubject)
	})

	t.Run("Bill", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Bill", storeStubBillSubject)
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()

		t.Run("Get", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Get", storeStubGetSubject)
		})

		t.Run("Find", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Find", storeStubFindSubject)
		})

		t.Run("Lookup", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Lookup", storeStubLookupSubject)
		})

		t.Run("List", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "List", storeStubListSubject)
		})

		t.Run("Note", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Note", storeStubNoteSubject)
		})

		t.Run("Bill", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Bill", storeStubBillSubject)
		})
	})

	t.Run("StoreStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) storeStubInstance {
			return storeStubInstanceOf(tb, built.NewStoreStub(tb, built.StoreStubStrict()))
		})
	})

	t.Run("StoreStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Get to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := built.NewStoreStub(t)
			s := built.NewStoreStub(t, built.StoreStubDelegateTo(inner))
			_, _ = storeStubGet(t, s)
			inner.OnGet.AssertCalledOnce(t, "Get calls the implementation")
		})

		t.Run("fails Get with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Get-delegate")
			inner := built.NewStoreStub(t)
			inner.OnGet.Faults(want, 1)
			s := built.NewStoreStub(t, built.StoreStubDelegateTo(inner))
			_, r1 := storeStubGet(t, s)
			assert.ErrorIs(t, r1, want, "Get fails with the error of the implementation")
		})

		t.Run("passes Find to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := built.NewStoreStub(t)
			s := built.NewStoreStub(t, built.StoreStubDelegateTo(inner))
			_, _ = storeStubFind(t, s)
			inner.OnFind.AssertCalledOnce(t, "Find calls the implementation")
		})

		t.Run("fails Find with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Find-delegate")
			inner := built.NewStoreStub(t)
			inner.OnFind.Faults(want, 1)
			s := built.NewStoreStub(t, built.StoreStubDelegateTo(inner))
			_, r1 := storeStubFind(t, s)
			assert.ErrorIs(t, r1, want, "Find fails with the error of the implementation")
		})

		t.Run("passes Lookup to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := built.NewStoreStub(t)
			s := built.NewStoreStub(t, built.StoreStubDelegateTo(inner))
			_, _ = storeStubLookup(t, s)
			inner.OnLookup.AssertCalledOnce(t, "Lookup calls the implementation")
		})

		t.Run("passes List to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := built.NewStoreStub(t)
			s := built.NewStoreStub(t, built.StoreStubDelegateTo(inner))
			_, _ = storeStubList(t, s)
			inner.OnList.AssertCalledOnce(t, "List calls the implementation")
		})

		t.Run("fails List with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("List-delegate")
			inner := built.NewStoreStub(t)
			inner.OnList.Faults(want, 1)
			s := built.NewStoreStub(t, built.StoreStubDelegateTo(inner))
			_, r1 := storeStubList(t, s)
			assert.ErrorIs(t, r1, want, "List fails with the error of the implementation")
		})

		t.Run("passes Note to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := built.NewStoreStub(t)
			s := built.NewStoreStub(t, built.StoreStubDelegateTo(inner))
			_ = storeStubNote(t, s)
			inner.OnNote.AssertCalledOnce(t, "Note calls the implementation")
		})

		t.Run("passes Bill to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := built.NewStoreStub(t)
			s := built.NewStoreStub(t, built.StoreStubDelegateTo(inner))
			_ = storeStubBill(t, s)
			inner.OnBill.AssertCalledOnce(t, "Bill calls the implementation")
		})
	})

	t.Run("StoreStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) storeStubInstance {
			return storeStubInstanceOf(tb, built.NewStoreStub(tb, built.StoreStubWithClock(clk)))
		})
	})

	t.Run("StoreStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) storeStubInstance {
			return storeStubInstanceOf(
				tb,
				built.NewStoreStub(tb, built.StoreStubWithRandSource(src)),
			)
		})
	})

	t.Run("StoreStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) storeStubInstance {
			return storeStubInstanceOf(tb, built.NewStoreStub(tb, built.StoreStubBenchMode()))
		})
	})

	t.Run("StoreStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) storeStubInstance {
				return storeStubInstanceOf(
					tb,
					built.NewStoreStub(
						tb,
						built.StoreStubWithClock(clk),
						built.StoreStubWithLatency(d),
					),
				)
			},
		)
	})

	t.Run("StoreStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Get with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Get-failing")
			s := built.NewStoreStub(t, built.StoreStubFailing(want))
			_, r1 := storeStubGet(t, s)
			assert.ErrorIs(t, r1, want, "Get fails with the error of StoreStubFailing")
		})

		t.Run("fails Find with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Find-failing")
			s := built.NewStoreStub(t, built.StoreStubFailing(want))
			_, r1 := storeStubFind(t, s)
			assert.ErrorIs(t, r1, want, "Find fails with the error of StoreStubFailing")
		})

		t.Run("fails List with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("List-failing")
			s := built.NewStoreStub(t, built.StoreStubFailing(want))
			_, r1 := storeStubList(t, s)
			assert.ErrorIs(t, r1, want, "List fails with the error of StoreStubFailing")
		})
	})

	t.Run("WithStoreGet", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Get", func(t *testing.T) {
			t.Parallel()
			called := false
			s := built.NewStoreStub(
				t,
				built.WithStoreGet(func(_ context.Context, _ string) (built.Account, error) {
					called = true
					var z0 built.Account
					var z1 error
					return z0, z1
				}),
			)
			_, _ = storeStubGet(t, s)
			assert.True(t, called, "WithStoreGet sets the function of Get")
		})
	})

	t.Run("WithStoreFind", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Find", func(t *testing.T) {
			t.Parallel()
			called := false
			s := built.NewStoreStub(
				t,
				built.WithStoreFind(func(_ context.Context, _ string) (*built.Account, error) {
					called = true
					var z0 *built.Account
					var z1 error
					return z0, z1
				}),
			)
			_, _ = storeStubFind(t, s)
			assert.True(t, called, "WithStoreFind sets the function of Find")
		})
	})

	t.Run("WithStoreLookup", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Lookup", func(t *testing.T) {
			t.Parallel()
			called := false
			s := built.NewStoreStub(t, built.WithStoreLookup(func(_ string) (built.Account, bool) {
				called = true
				var z0 built.Account
				var z1 bool
				return z0, z1
			}))
			_, _ = storeStubLookup(t, s)
			assert.True(t, called, "WithStoreLookup sets the function of Lookup")
		})
	})

	t.Run("WithStoreList", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of List", func(t *testing.T) {
			t.Parallel()
			called := false
			s := built.NewStoreStub(
				t,
				built.WithStoreList(func(_ context.Context) (built.Page[built.Account], error) {
					called = true
					var z0 built.Page[built.Account]
					var z1 error
					return z0, z1
				}),
			)
			_, _ = storeStubList(t, s)
			assert.True(t, called, "WithStoreList sets the function of List")
		})
	})

	t.Run("WithStoreNote", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Note", func(t *testing.T) {
			t.Parallel()
			called := false
			s := built.NewStoreStub(t, built.WithStoreNote(func(_ string) built.Note {
				called = true
				var z0 built.Note
				return z0
			}))
			_ = storeStubNote(t, s)
			assert.True(t, called, "WithStoreNote sets the function of Note")
		})
	})

	t.Run("WithStoreBill", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Bill", func(t *testing.T) {
			t.Parallel()
			called := false
			s := built.NewStoreStub(t, built.WithStoreBill(func(_ string) billing.Invoice {
				called = true
				var z0 billing.Invoice
				return z0
			}))
			_ = storeStubBill(t, s)
			assert.True(t, called, "WithStoreBill sets the function of Bill")
		})
	})
}
