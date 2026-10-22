package built

import (
	"context"
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"example.com/shop/billing"
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
func storeStubGet(tb assert.TB, s *StoreStub) (Account, error) {
	tb.Helper()
	a0, a1 := storeStubGetArgs(tb)
	return s.Get(a0, a1)
}

// storeStubGetPins returns the samples of the results of Get, with a new
// error as its error result.
func storeStubGetPins() (Account, error) {
	var want0 Account
	want1 := errors.New("Get-returns")
	want0 = Account{ID: "ID-a"}
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
func storeStubFind(tb assert.TB, s *StoreStub) (*Account, error) {
	tb.Helper()
	a0, a1 := storeStubFindArgs(tb)
	return s.Find(a0, a1)
}

// storeStubFindPins returns the samples of the results of Find, with a new
// error as its error result.
func storeStubFindPins() (*Account, error) {
	var want0 *Account
	want1 := errors.New("Find-returns")
	want0 = &Account{ID: "ID-a"}
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
func storeStubLookup(tb assert.TB, s *StoreStub) (Account, bool) {
	tb.Helper()
	a0 := storeStubLookupArgs(tb)
	return s.Lookup(a0)
}

// storeStubLookupPins returns the samples of the results of Lookup.
func storeStubLookupPins() (Account, bool) {
	var want0 Account
	var want1 bool
	want0 = Account{ID: "ID-a"}
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
func storeStubList(tb assert.TB, s *StoreStub) (Page[Account], error) {
	tb.Helper()
	a0 := storeStubListArgs(tb)
	return s.List(a0)
}

// storeStubListPins returns the samples of the results of List, with a new
// error as its error result.
func storeStubListPins() (Page[Account], error) {
	var want0 Page[Account]
	want1 := errors.New("List-returns")
	want0 = Page[Account]{Items: []Account{{ID: "ID-a"}}}
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
func storeStubNote(tb assert.TB, s *StoreStub) Note {
	tb.Helper()
	a0 := storeStubNoteArgs(tb)
	return s.Note(a0)
}

// storeStubNotePins returns the samples of the results of Note.
func storeStubNotePins() Note {
	var want0 Note
	want0 = Note{Text: "Text-a"}
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
func storeStubBill(tb assert.TB, s *StoreStub) billing.Invoice {
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

// storeStubGetSubject returns the subject of [stub.Behaviour] for Get. The
// subject calls a new double that reports to tb.
func storeStubGetSubject(tb assert.TB) stub.Subject[StoreGetCall, StoreGetReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreGetCall, StoreGetReturn]{
		Stub: s.OnGet.MethodStub,
		Call: func() {
			_, _ = storeStubGet(tb, s)
		},
		Result: func() StoreGetReturn {
			got0, got1 := storeStubGet(tb, s)
			return StoreGetReturn{Result: got0, Err: got1}
		},
		Default: func() StoreGetReturn {
			r0 := NewAccountBuilder().Build()
			return StoreGetReturn{Result: r0}
		},
		Override: func(mark func()) {
			s.OnGet.Func(func(_ context.Context, _ string) (Account, error) {
				mark()
				var z0 Account
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

// storeStubFindSubject returns the subject of [stub.Behaviour] for Find.
// The subject calls a new double that reports to tb.
func storeStubFindSubject(tb assert.TB) stub.Subject[StoreFindCall, StoreFindReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreFindCall, StoreFindReturn]{
		Stub: s.OnFind.MethodStub,
		Call: func() {
			_, _ = storeStubFind(tb, s)
		},
		Result: func() StoreFindReturn {
			got0, got1 := storeStubFind(tb, s)
			return StoreFindReturn{Result: got0, Err: got1}
		},
		Default: func() StoreFindReturn {
			r0 := NewAccountBuilder().Build()
			return StoreFindReturn{Result: &r0}
		},
		Override: func(mark func()) {
			s.OnFind.Func(func(_ context.Context, _ string) (*Account, error) {
				mark()
				var z0 *Account
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

// storeStubLookupSubject returns the subject of [stub.Behaviour] for
// Lookup. The subject calls a new double that reports to tb.
func storeStubLookupSubject(tb assert.TB) stub.Subject[StoreLookupCall, StoreLookupReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreLookupCall, StoreLookupReturn]{
		Stub: s.OnLookup.MethodStub,
		Call: func() {
			_, _ = storeStubLookup(tb, s)
		},
		Result: func() StoreLookupReturn {
			got0, got1 := storeStubLookup(tb, s)
			return StoreLookupReturn{Result0: got0, Result1: got1}
		},
		Default: func() StoreLookupReturn {
			r0 := NewAccountBuilder().Build()
			return StoreLookupReturn{Result0: r0}
		},
		Override: func(mark func()) {
			s.OnLookup.Func(func(_ string) (Account, bool) {
				mark()
				var z0 Account
				var z1 bool
				return z0, z1
			})
		},
	}
}

// storeStubListSubject returns the subject of [stub.Behaviour] for List.
// The subject calls a new double that reports to tb.
func storeStubListSubject(tb assert.TB) stub.Subject[StoreListCall, StoreListReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreListCall, StoreListReturn]{
		Stub: s.OnList.MethodStub,
		Call: func() {
			_, _ = storeStubList(tb, s)
		},
		Result: func() StoreListReturn {
			got0, got1 := storeStubList(tb, s)
			return StoreListReturn{Result: got0, Err: got1}
		},
		Default: func() StoreListReturn {
			r0 := NewPageBuilder[Account]().Build()
			return StoreListReturn{Result: r0}
		},
		Override: func(mark func()) {
			s.OnList.Func(func(_ context.Context) (Page[Account], error) {
				mark()
				var z0 Page[Account]
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

// storeStubNoteSubject returns the subject of [stub.Behaviour] for Note.
// The subject calls a new double that reports to tb.
func storeStubNoteSubject(tb assert.TB) stub.Subject[StoreNoteCall, StoreNoteReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreNoteCall, StoreNoteReturn]{
		Stub: s.OnNote.MethodStub,
		Call: func() {
			_ = storeStubNote(tb, s)
		},
		Result: func() StoreNoteReturn {
			got0 := storeStubNote(tb, s)
			return StoreNoteReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnNote.Func(func(_ string) Note {
				mark()
				var z0 Note
				return z0
			})
		},
	}
}

// storeStubBillSubject returns the subject of [stub.Behaviour] for Bill.
// The subject calls a new double that reports to tb.
func storeStubBillSubject(tb assert.TB) stub.Subject[StoreBillCall, StoreBillReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreBillCall, StoreBillReturn]{
		Stub: s.OnBill.MethodStub,
		Call: func() {
			_ = storeStubBill(tb, s)
		},
		Result: func() StoreBillReturn {
			got0 := storeStubBill(tb, s)
			return StoreBillReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnBill.Func(func(_ string) billing.Invoice {
				mark()
				var z0 billing.Invoice
				return z0
			})
		},
	}
}

// storeStubInstance is the instance of a StoreStub that the checks of the
// constructors and the options drive.
type storeStubInstance = stub.Instance[StoreGetCall]

// storeStubInstanceOf returns the instance of s. The options apply to the
// whole double, so the checks call Get alone.
func storeStubInstanceOf(tb assert.TB, s *StoreStub) storeStubInstance {
	tb.Helper()
	return storeStubInstance{
		Stub: s.OnGet.MethodStub,
		Call: func() {
			_, _ = storeStubGet(tb, s)
		},
		Reset: s.ResetCalls,
	}
}

// TestStoreStub checks the methods, the constructors and the options of
// the double StoreStub.
func TestStoreStub(t *testing.T) {
	t.Parallel()

	t.Run("Get", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Get", storeStubGetSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			_, _ = storeStubGet(t, s)
			got := s.OnGet.AssertCalledOnce(t, "Get records the call")
			a0, a1 := storeStubGetArgs(t)
			expect.Equal(t, got.Ctx, a0, "the call records Ctx", expect.ByIdentity())
			expect.Equal(t, got.ID, a1, "the call records ID")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			var seen []StoreGetCall
			s.OnGet.OnRecord(func(c StoreGetCall) { seen = append(seen, c) })
			_, _ = storeStubGet(t, s)
			_, _ = storeStubGet(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Get")
		})
	})

	t.Run("Find", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Find", storeStubFindSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			_, _ = storeStubFind(t, s)
			got := s.OnFind.AssertCalledOnce(t, "Find records the call")
			a0, a1 := storeStubFindArgs(t)
			expect.Equal(t, got.Ctx, a0, "the call records Ctx", expect.ByIdentity())
			expect.Equal(t, got.ID, a1, "the call records ID")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			var seen []StoreFindCall
			s.OnFind.OnRecord(func(c StoreFindCall) { seen = append(seen, c) })
			_, _ = storeStubFind(t, s)
			_, _ = storeStubFind(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Find")
		})
	})

	t.Run("Lookup", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Lookup", storeStubLookupSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			_, _ = storeStubLookup(t, s)
			got := s.OnLookup.AssertCalledOnce(t, "Lookup records the call")
			a0 := storeStubLookupArgs(t)
			assert.Equal(t, got.ID, a0, "the call records ID")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			var seen []StoreLookupCall
			s.OnLookup.OnRecord(func(c StoreLookupCall) { seen = append(seen, c) })
			_, _ = storeStubLookup(t, s)
			_, _ = storeStubLookup(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Lookup")
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "List", storeStubListSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			_, _ = storeStubList(t, s)
			got := s.OnList.AssertCalledOnce(t, "List records the call")
			a0 := storeStubListArgs(t)
			assert.Equal(t, got.Ctx, a0, "the call records Ctx", assert.ByIdentity())
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			var seen []StoreListCall
			s.OnList.OnRecord(func(c StoreListCall) { seen = append(seen, c) })
			_, _ = storeStubList(t, s)
			_, _ = storeStubList(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of List")
		})
	})

	t.Run("Note", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Note", storeStubNoteSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			_ = storeStubNote(t, s)
			got := s.OnNote.AssertCalledOnce(t, "Note records the call")
			a0 := storeStubNoteArgs(t)
			assert.Equal(t, got.ID, a0, "the call records ID")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			var seen []StoreNoteCall
			s.OnNote.OnRecord(func(c StoreNoteCall) { seen = append(seen, c) })
			_ = storeStubNote(t, s)
			_ = storeStubNote(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Note")
		})
	})

	t.Run("Bill", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Bill", storeStubBillSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			_ = storeStubBill(t, s)
			got := s.OnBill.AssertCalledOnce(t, "Bill records the call")
			a0 := storeStubBillArgs(t)
			assert.Equal(t, got.ID, a0, "the call records ID")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			var seen []StoreBillCall
			s.OnBill.OnRecord(func(c StoreBillCall) { seen = append(seen, c) })
			_ = storeStubBill(t, s)
			_ = storeStubBill(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Bill")
		})
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()
		stub.ResetBehaviour(t, func(tb assert.TB) storeStubInstance {
			return storeStubInstanceOf(tb, NewStoreStub(tb))
		})

		t.Run("clears the recorded calls of Get", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			_, _ = storeStubGet(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnGet.CallCount(), 0, "ResetCalls clears the calls of Get")
		})

		t.Run("keeps the results of Returns of Get", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			want0, want1 := storeStubGetPins()
			s.OnGet.Returns(want0, want1)
			s.ResetCalls()
			got0, got1 := storeStubGet(t, s)
			expect.Equal(t, got0, want0, "ResetCalls keeps the Result of Get")
			expect.ErrorIs(t, got1, want1, "ResetCalls keeps the error of Get")
		})

		t.Run("clears the recorded calls of Find", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			_, _ = storeStubFind(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnFind.CallCount(), 0, "ResetCalls clears the calls of Find")
		})

		t.Run("keeps the results of Returns of Find", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			want0, want1 := storeStubFindPins()
			s.OnFind.Returns(want0, want1)
			s.ResetCalls()
			got0, got1 := storeStubFind(t, s)
			expect.Equal(t, got0, want0, "ResetCalls keeps the Result of Find")
			expect.ErrorIs(t, got1, want1, "ResetCalls keeps the error of Find")
		})

		t.Run("clears the recorded calls of Lookup", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			_, _ = storeStubLookup(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnLookup.CallCount(), 0, "ResetCalls clears the calls of Lookup")
		})

		t.Run("keeps the results of Returns of Lookup", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			want0, want1 := storeStubLookupPins()
			s.OnLookup.Returns(want0, want1)
			s.ResetCalls()
			got0, got1 := storeStubLookup(t, s)
			expect.Equal(t, got0, want0, "ResetCalls keeps the Result0 of Lookup")
			expect.Equal(t, got1, want1, "ResetCalls keeps the Result1 of Lookup")
		})

		t.Run("clears the recorded calls of List", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			_, _ = storeStubList(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnList.CallCount(), 0, "ResetCalls clears the calls of List")
		})

		t.Run("keeps the results of Returns of List", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			want0, want1 := storeStubListPins()
			s.OnList.Returns(want0, want1)
			s.ResetCalls()
			got0, got1 := storeStubList(t, s)
			expect.Equal(t, got0, want0, "ResetCalls keeps the Result of List")
			expect.ErrorIs(t, got1, want1, "ResetCalls keeps the error of List")
		})

		t.Run("clears the recorded calls of Note", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			_ = storeStubNote(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnNote.CallCount(), 0, "ResetCalls clears the calls of Note")
		})

		t.Run("keeps the results of Returns of Note", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			want0 := storeStubNotePins()
			s.OnNote.Returns(want0)
			s.ResetCalls()
			got0 := storeStubNote(t, s)
			assert.Equal(t, got0, want0, "ResetCalls keeps the Result of Note")
		})

		t.Run("clears the recorded calls of Bill", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			_ = storeStubBill(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnBill.CallCount(), 0, "ResetCalls clears the calls of Bill")
		})

		t.Run("keeps the results of Returns of Bill", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			want0 := storeStubBillPins()
			s.OnBill.Returns(want0)
			s.ResetCalls()
			got0 := storeStubBill(t, s)
			assert.Equal(t, got0, want0, "ResetCalls keeps the Result of Bill")
		})
	})

	t.Run("NewStoreStubDelayed", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) storeStubInstance {
				return storeStubInstanceOf(tb, NewStoreStubDelayed(tb, d, clk))
			},
		)
	})

	t.Run("NewStoreStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Get with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Get-failing")
			s := NewStoreStubFailing(t, want)
			_, r1 := storeStubGet(t, s)
			assert.ErrorIs(t, r1, want, "Get fails with the error of NewStoreStubFailing")
		})

		t.Run("fails Find with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Find-failing")
			s := NewStoreStubFailing(t, want)
			_, r1 := storeStubFind(t, s)
			assert.ErrorIs(t, r1, want, "Find fails with the error of NewStoreStubFailing")
		})

		t.Run("fails List with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("List-failing")
			s := NewStoreStubFailing(t, want)
			_, r1 := storeStubList(t, s)
			assert.ErrorIs(t, r1, want, "List fails with the error of NewStoreStubFailing")
		})
	})

	t.Run("StoreStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) storeStubInstance {
			return storeStubInstanceOf(tb, NewStoreStub(tb, StoreStubStrict()))
		})
	})

	t.Run("StoreStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Get to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewStoreStub(t)
			s := NewStoreStub(t, StoreStubDelegateTo(inner))
			_, _ = storeStubGet(t, s)
			inner.OnGet.AssertCalledOnce(t, "Get calls the implementation")
		})

		t.Run("fails Get with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Get-delegate")
			inner := NewStoreStub(t)
			inner.OnGet.Faults(want, 1)
			s := NewStoreStub(t, StoreStubDelegateTo(inner))
			_, r1 := storeStubGet(t, s)
			assert.ErrorIs(t, r1, want, "Get fails with the error of the implementation")
		})

		t.Run("passes Find to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewStoreStub(t)
			s := NewStoreStub(t, StoreStubDelegateTo(inner))
			_, _ = storeStubFind(t, s)
			inner.OnFind.AssertCalledOnce(t, "Find calls the implementation")
		})

		t.Run("fails Find with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Find-delegate")
			inner := NewStoreStub(t)
			inner.OnFind.Faults(want, 1)
			s := NewStoreStub(t, StoreStubDelegateTo(inner))
			_, r1 := storeStubFind(t, s)
			assert.ErrorIs(t, r1, want, "Find fails with the error of the implementation")
		})

		t.Run("passes Lookup to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewStoreStub(t)
			s := NewStoreStub(t, StoreStubDelegateTo(inner))
			_, _ = storeStubLookup(t, s)
			inner.OnLookup.AssertCalledOnce(t, "Lookup calls the implementation")
		})

		t.Run("passes List to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewStoreStub(t)
			s := NewStoreStub(t, StoreStubDelegateTo(inner))
			_, _ = storeStubList(t, s)
			inner.OnList.AssertCalledOnce(t, "List calls the implementation")
		})

		t.Run("fails List with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("List-delegate")
			inner := NewStoreStub(t)
			inner.OnList.Faults(want, 1)
			s := NewStoreStub(t, StoreStubDelegateTo(inner))
			_, r1 := storeStubList(t, s)
			assert.ErrorIs(t, r1, want, "List fails with the error of the implementation")
		})

		t.Run("passes Note to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewStoreStub(t)
			s := NewStoreStub(t, StoreStubDelegateTo(inner))
			_ = storeStubNote(t, s)
			inner.OnNote.AssertCalledOnce(t, "Note calls the implementation")
		})

		t.Run("passes Bill to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewStoreStub(t)
			s := NewStoreStub(t, StoreStubDelegateTo(inner))
			_ = storeStubBill(t, s)
			inner.OnBill.AssertCalledOnce(t, "Bill calls the implementation")
		})
	})

	t.Run("StoreStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) storeStubInstance {
			return storeStubInstanceOf(tb, NewStoreStub(tb, StoreStubWithClock(clk)))
		})
	})

	t.Run("StoreStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) storeStubInstance {
			return storeStubInstanceOf(tb, NewStoreStub(tb, StoreStubWithRandSource(src)))
		})
	})

	t.Run("StoreStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) storeStubInstance {
			return storeStubInstanceOf(tb, NewStoreStub(tb, StoreStubBenchMode()))
		})
	})

	t.Run("StoreStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) storeStubInstance {
				return storeStubInstanceOf(
					tb,
					NewStoreStub(tb, StoreStubWithClock(clk), StoreStubWithLatency(d)),
				)
			},
		)
	})

	t.Run("StoreStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Get with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Get-failing")
			s := NewStoreStub(t, StoreStubFailing(want))
			_, r1 := storeStubGet(t, s)
			assert.ErrorIs(t, r1, want, "Get fails with the error of StoreStubFailing")
		})

		t.Run("fails Find with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Find-failing")
			s := NewStoreStub(t, StoreStubFailing(want))
			_, r1 := storeStubFind(t, s)
			assert.ErrorIs(t, r1, want, "Find fails with the error of StoreStubFailing")
		})

		t.Run("fails List with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("List-failing")
			s := NewStoreStub(t, StoreStubFailing(want))
			_, r1 := storeStubList(t, s)
			assert.ErrorIs(t, r1, want, "List fails with the error of StoreStubFailing")
		})
	})

	t.Run("WithStoreGet", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Get", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewStoreStub(t, WithStoreGet(func(_ context.Context, _ string) (Account, error) {
				called = true
				var z0 Account
				var z1 error
				return z0, z1
			}))
			_, _ = storeStubGet(t, s)
			assert.True(t, called, "WithStoreGet sets the function of Get")
		})
	})

	t.Run("WithStoreFind", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Find", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewStoreStub(t, WithStoreFind(func(_ context.Context, _ string) (*Account, error) {
				called = true
				var z0 *Account
				var z1 error
				return z0, z1
			}))
			_, _ = storeStubFind(t, s)
			assert.True(t, called, "WithStoreFind sets the function of Find")
		})
	})

	t.Run("WithStoreLookup", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Lookup", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewStoreStub(t, WithStoreLookup(func(_ string) (Account, bool) {
				called = true
				var z0 Account
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
			s := NewStoreStub(t, WithStoreList(func(_ context.Context) (Page[Account], error) {
				called = true
				var z0 Page[Account]
				var z1 error
				return z0, z1
			}))
			_, _ = storeStubList(t, s)
			assert.True(t, called, "WithStoreList sets the function of List")
		})
	})

	t.Run("WithStoreNote", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Note", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewStoreStub(t, WithStoreNote(func(_ string) Note {
				called = true
				var z0 Note
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
			s := NewStoreStub(t, WithStoreBill(func(_ string) billing.Invoice {
				called = true
				var z0 billing.Invoice
				return z0
			}))
			_ = storeStubBill(t, s)
			assert.True(t, called, "WithStoreBill sets the function of Bill")
		})
	})
}

// TestStoreGetStub checks the methods of StoreGetStub, which configures
// Get.
func TestStoreGetStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Get", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			want0, want1 := storeStubGetPins()
			assert.NotEqual(t, want0, storeStubGetSubject(t).Default().Result,
				"the sample of Result differs from its default result")
			s.OnGet.Returns(want0, want1)
			got0, got1 := storeStubGet(t, s)
			expect.Equal(t, got0, want0, "Get returns the Result that Returns sets")
			expect.ErrorIs(t, got1, want1, "Get returns the error that Returns sets")
		})

		t.Run("configures Get for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewStoreStub(r)
			s.OnGet.Strict()
			var want0 Account
			var want1 error
			s.OnGet.Returns(want0, want1)
			_, _ = storeStubGet(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Get")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Get to the function", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			var got0 context.Context
			var got1 string
			s.OnGet.Func(func(p0 context.Context, p1 string) (Account, error) {
				got0, got1 = p0, p1
				var z0 Account
				var z1 error
				return z0, z1
			})
			_, _ = storeStubGet(t, s)
			a0, a1 := storeStubGetArgs(t)
			expect.Equal(t, got0, a0, "the function receives Ctx", expect.ByIdentity())
			expect.Equal(t, got1, a1, "the function receives ID")
		})
	})
}

// TestStoreFindStub checks the methods of StoreFindStub, which configures
// Find.
func TestStoreFindStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Find", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			want0, want1 := storeStubFindPins()
			assert.NotEqual(t, want0, storeStubFindSubject(t).Default().Result,
				"the sample of Result differs from its default result")
			s.OnFind.Returns(want0, want1)
			got0, got1 := storeStubFind(t, s)
			expect.Equal(t, got0, want0, "Find returns the Result that Returns sets")
			expect.ErrorIs(t, got1, want1, "Find returns the error that Returns sets")
		})

		t.Run("configures Find for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewStoreStub(r)
			s.OnFind.Strict()
			var want0 *Account
			var want1 error
			s.OnFind.Returns(want0, want1)
			_, _ = storeStubFind(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Find")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Find to the function", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			var got0 context.Context
			var got1 string
			s.OnFind.Func(func(p0 context.Context, p1 string) (*Account, error) {
				got0, got1 = p0, p1
				var z0 *Account
				var z1 error
				return z0, z1
			})
			_, _ = storeStubFind(t, s)
			a0, a1 := storeStubFindArgs(t)
			expect.Equal(t, got0, a0, "the function receives Ctx", expect.ByIdentity())
			expect.Equal(t, got1, a1, "the function receives ID")
		})
	})
}

// TestStoreLookupStub checks the methods of StoreLookupStub, which
// configures Lookup.
func TestStoreLookupStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Lookup", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			want0, want1 := storeStubLookupPins()
			assert.NotEqual(t, want0, storeStubLookupSubject(t).Default().Result0,
				"the sample of Result0 differs from its default result")
			s.OnLookup.Returns(want0, want1)
			got0, got1 := storeStubLookup(t, s)
			expect.Equal(t, got0, want0, "Lookup returns the Result0 that Returns sets")
			expect.Equal(t, got1, want1, "Lookup returns the Result1 that Returns sets")
		})

		t.Run("configures Lookup for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewStoreStub(r)
			s.OnLookup.Strict()
			var want0 Account
			var want1 bool
			s.OnLookup.Returns(want0, want1)
			_, _ = storeStubLookup(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Lookup")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Lookup to the function", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			var got0 string
			s.OnLookup.Func(func(p0 string) (Account, bool) {
				got0 = p0
				var z0 Account
				var z1 bool
				return z0, z1
			})
			_, _ = storeStubLookup(t, s)
			a0 := storeStubLookupArgs(t)
			assert.Equal(t, got0, a0, "the function receives ID")
		})
	})
}

// TestStoreListStub checks the methods of StoreListStub, which configures
// List.
func TestStoreListStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of List", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			want0, want1 := storeStubListPins()
			assert.NotEqual(t, want0, storeStubListSubject(t).Default().Result,
				"the sample of Result differs from its default result")
			s.OnList.Returns(want0, want1)
			got0, got1 := storeStubList(t, s)
			expect.Equal(t, got0, want0, "List returns the Result that Returns sets")
			expect.ErrorIs(t, got1, want1, "List returns the error that Returns sets")
		})

		t.Run("configures List for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewStoreStub(r)
			s.OnList.Strict()
			var want0 Page[Account]
			var want1 error
			s.OnList.Returns(want0, want1)
			_, _ = storeStubList(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of List")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of List to the function", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			var got0 context.Context
			s.OnList.Func(func(p0 context.Context) (Page[Account], error) {
				got0 = p0
				var z0 Page[Account]
				var z1 error
				return z0, z1
			})
			_, _ = storeStubList(t, s)
			a0 := storeStubListArgs(t)
			assert.Equal(t, got0, a0, "the function receives Ctx", assert.ByIdentity())
		})
	})
}

// TestStoreNoteStub checks the methods of StoreNoteStub, which configures
// Note.
func TestStoreNoteStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Note", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			want0 := storeStubNotePins()
			s.OnNote.Returns(want0)
			got0 := storeStubNote(t, s)
			assert.Equal(t, got0, want0, "Note returns the Result that Returns sets")
		})

		t.Run("configures Note for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewStoreStub(r)
			s.OnNote.Strict()
			var want0 Note
			s.OnNote.Returns(want0)
			_ = storeStubNote(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Note")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Note to the function", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			var got0 string
			s.OnNote.Func(func(p0 string) Note {
				got0 = p0
				var z0 Note
				return z0
			})
			_ = storeStubNote(t, s)
			a0 := storeStubNoteArgs(t)
			assert.Equal(t, got0, a0, "the function receives ID")
		})
	})
}

// TestStoreBillStub checks the methods of StoreBillStub, which configures
// Bill.
func TestStoreBillStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Bill", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			want0 := storeStubBillPins()
			s.OnBill.Returns(want0)
			got0 := storeStubBill(t, s)
			assert.Equal(t, got0, want0, "Bill returns the Result that Returns sets")
		})

		t.Run("configures Bill for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewStoreStub(r)
			s.OnBill.Strict()
			var want0 billing.Invoice
			s.OnBill.Returns(want0)
			_ = storeStubBill(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Bill")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Bill to the function", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			var got0 string
			s.OnBill.Func(func(p0 string) billing.Invoice {
				got0 = p0
				var z0 billing.Invoice
				return z0
			})
			_ = storeStubBill(t, s)
			a0 := storeStubBillArgs(t)
			assert.Equal(t, got0, a0, "the function receives ID")
		})
	})
}
