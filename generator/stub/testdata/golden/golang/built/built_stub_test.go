package built

import (
	"context"
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"example.com/shop/billing"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// storeStubGetSubject binds Get into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func storeStubGetSubject(tb assert.TB) stub.Subject[StoreGetCall, StoreGetReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreGetCall, StoreGetReturn]{
		Stub: s.OnGet.MethodStub,
		Call: func() {
			var a0 context.Context
			var a1 string
			a1 = "Get-a"
			_, _ = s.Get(a0, a1)
		},
		Result: func() StoreGetReturn {
			var a0 context.Context
			var a1 string
			a1 = "Get-a"
			got0, got1 := s.Get(a0, a1)
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
			var a0 context.Context
			var a1 string
			a1 = "Get-a"
			_, r1 := s.Get(a0, a1)
			return r1
		},
	}
}

// TestStoreStubGet pins how Get answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestStoreStubGet(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Get", storeStubGetSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 Account
		var want1 error
		want0 = Account{ID: "ID-a"}
		assert.NotEqual(t, want0, storeStubGetSubject(t).Default().Result,
			"the sample of Result must differ from its default result")
		s.OnGet.Returns(want0, want1)
		var a0 context.Context
		var a1 string
		a1 = "Get-a"
		got0, _ := s.Get(a0, a1)
		assert.Equal(t, got0, want0, "Get must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var a0 context.Context
		var a1 string
		a1 = "Get-a"
		_, _ = s.Get(a0, a1)
		got := s.OnGet.AssertCalledOnce(t, "Get must record the call")
		assert.Equal(t, got.Ctx, a0, "the recorded call carries Ctx")
		assert.Equal(t, got.ID, a1, "the recorded call carries ID")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var seen []StoreGetCall
		s.OnGet.OnRecord(func(c StoreGetCall) { seen = append(seen, c) })
		var a0 context.Context
		var a1 string
		a1 = "Get-a"
		_, _ = s.Get(a0, a1)
		_, _ = s.Get(a0, a1)
		assert.Length(t, seen, 2, "OnRecord must fire once per Get call")
	})

	t.Run("wires WithStoreGet at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewStoreStub(t, WithStoreGet(func(_ context.Context, _ string) (Account, error) {
			called = true
			var z0 Account
			var z1 error
			return z0, z1
		}))
		var a0 context.Context
		var a1 string
		a1 = "Get-a"
		_, _ = s.Get(a0, a1)
		assert.True(t, called, "WithStoreGet must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("Get-failing")
		s := NewStoreStubFailing(t, want)
		var a0 context.Context
		var a1 string
		a1 = "Get-a"
		_, r1 := s.Get(a0, a1)
		assert.ErrorIs(t, r1, want, "Get must answer the armed error")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 Account
		var want1 error
		want0 = Account{ID: "ID-a"}
		s.OnGet.Returns(want0, want1)
		var a0 context.Context
		var a1 string
		a1 = "Get-a"
		_, _ = s.Get(a0, a1)
		s.ResetCalls()
		got0, _ := s.Get(a0, a1)
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// storeStubFindSubject binds Find into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func storeStubFindSubject(tb assert.TB) stub.Subject[StoreFindCall, StoreFindReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreFindCall, StoreFindReturn]{
		Stub: s.OnFind.MethodStub,
		Call: func() {
			var a0 context.Context
			var a1 string
			a1 = "Find-a"
			_, _ = s.Find(a0, a1)
		},
		Result: func() StoreFindReturn {
			var a0 context.Context
			var a1 string
			a1 = "Find-a"
			got0, got1 := s.Find(a0, a1)
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
			var a0 context.Context
			var a1 string
			a1 = "Find-a"
			_, r1 := s.Find(a0, a1)
			return r1
		},
	}
}

// TestStoreStubFind pins how Find answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestStoreStubFind(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Find", storeStubFindSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 *Account
		var want1 error
		want0 = &Account{ID: "ID-a"}
		assert.NotEqual(t, want0, storeStubFindSubject(t).Default().Result,
			"the sample of Result must differ from its default result")
		s.OnFind.Returns(want0, want1)
		var a0 context.Context
		var a1 string
		a1 = "Find-a"
		got0, _ := s.Find(a0, a1)
		assert.Equal(t, got0, want0, "Find must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var a0 context.Context
		var a1 string
		a1 = "Find-a"
		_, _ = s.Find(a0, a1)
		got := s.OnFind.AssertCalledOnce(t, "Find must record the call")
		assert.Equal(t, got.Ctx, a0, "the recorded call carries Ctx")
		assert.Equal(t, got.ID, a1, "the recorded call carries ID")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var seen []StoreFindCall
		s.OnFind.OnRecord(func(c StoreFindCall) { seen = append(seen, c) })
		var a0 context.Context
		var a1 string
		a1 = "Find-a"
		_, _ = s.Find(a0, a1)
		_, _ = s.Find(a0, a1)
		assert.Length(t, seen, 2, "OnRecord must fire once per Find call")
	})

	t.Run("wires WithStoreFind at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewStoreStub(t, WithStoreFind(func(_ context.Context, _ string) (*Account, error) {
			called = true
			var z0 *Account
			var z1 error
			return z0, z1
		}))
		var a0 context.Context
		var a1 string
		a1 = "Find-a"
		_, _ = s.Find(a0, a1)
		assert.True(t, called, "WithStoreFind must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("Find-failing")
		s := NewStoreStubFailing(t, want)
		var a0 context.Context
		var a1 string
		a1 = "Find-a"
		_, r1 := s.Find(a0, a1)
		assert.ErrorIs(t, r1, want, "Find must answer the armed error")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 *Account
		var want1 error
		want0 = &Account{ID: "ID-a"}
		s.OnFind.Returns(want0, want1)
		var a0 context.Context
		var a1 string
		a1 = "Find-a"
		_, _ = s.Find(a0, a1)
		s.ResetCalls()
		got0, _ := s.Find(a0, a1)
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// storeStubLookupSubject binds Lookup into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func storeStubLookupSubject(tb assert.TB) stub.Subject[StoreLookupCall, StoreLookupReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreLookupCall, StoreLookupReturn]{
		Stub: s.OnLookup.MethodStub,
		Call: func() {
			var a0 string
			a0 = "Lookup-a"
			_, _ = s.Lookup(a0)
		},
		Result: func() StoreLookupReturn {
			var a0 string
			a0 = "Lookup-a"
			got0, got1 := s.Lookup(a0)
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

// TestStoreStubLookup pins how Lookup answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestStoreStubLookup(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Lookup", storeStubLookupSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 Account
		var want1 bool
		want0 = Account{ID: "ID-a"}
		want1 = true
		assert.NotEqual(t, want0, storeStubLookupSubject(t).Default().Result0,
			"the sample of Result0 must differ from its default result")
		s.OnLookup.Returns(want0, want1)
		var a0 string
		a0 = "Lookup-a"
		got0, got1 := s.Lookup(a0)
		assert.Equal(t, got0, want0, "Lookup must answer with what Returns pinned")
		assert.Equal(t, got1, want1, "Lookup must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var a0 string
		a0 = "Lookup-a"
		_, _ = s.Lookup(a0)
		got := s.OnLookup.AssertCalledOnce(t, "Lookup must record the call")
		assert.Equal(t, got.ID, a0, "the recorded call carries ID")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var seen []StoreLookupCall
		s.OnLookup.OnRecord(func(c StoreLookupCall) { seen = append(seen, c) })
		var a0 string
		a0 = "Lookup-a"
		_, _ = s.Lookup(a0)
		_, _ = s.Lookup(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Lookup call")
	})

	t.Run("wires WithStoreLookup at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewStoreStub(t, WithStoreLookup(func(_ string) (Account, bool) {
			called = true
			var z0 Account
			var z1 bool
			return z0, z1
		}))
		var a0 string
		a0 = "Lookup-a"
		_, _ = s.Lookup(a0)
		assert.True(t, called, "WithStoreLookup must install the override")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 Account
		var want1 bool
		want0 = Account{ID: "ID-a"}
		want1 = true
		s.OnLookup.Returns(want0, want1)
		var a0 string
		a0 = "Lookup-a"
		_, _ = s.Lookup(a0)
		s.ResetCalls()
		got0, got1 := s.Lookup(a0)
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
		assert.Equal(t, got1, want1, "a reset must keep what Returns pinned")
	})
}

// storeStubListSubject binds List into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func storeStubListSubject(tb assert.TB) stub.Subject[StoreListCall, StoreListReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreListCall, StoreListReturn]{
		Stub: s.OnList.MethodStub,
		Call: func() {
			var a0 context.Context
			_, _ = s.List(a0)
		},
		Result: func() StoreListReturn {
			var a0 context.Context
			got0, got1 := s.List(a0)
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
			var a0 context.Context
			_, r1 := s.List(a0)
			return r1
		},
	}
}

// TestStoreStubList pins how List answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestStoreStubList(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "List", storeStubListSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 Page[Account]
		var want1 error
		want0 = Page[Account]{Items: []Account{Account{ID: "ID-a"}}}
		assert.NotEqual(t, want0, storeStubListSubject(t).Default().Result,
			"the sample of Result must differ from its default result")
		s.OnList.Returns(want0, want1)
		var a0 context.Context
		got0, _ := s.List(a0)
		assert.Equal(t, got0, want0, "List must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var a0 context.Context
		_, _ = s.List(a0)
		got := s.OnList.AssertCalledOnce(t, "List must record the call")
		assert.Equal(t, got.Ctx, a0, "the recorded call carries Ctx")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var seen []StoreListCall
		s.OnList.OnRecord(func(c StoreListCall) { seen = append(seen, c) })
		var a0 context.Context
		_, _ = s.List(a0)
		_, _ = s.List(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per List call")
	})

	t.Run("wires WithStoreList at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewStoreStub(t, WithStoreList(func(_ context.Context) (Page[Account], error) {
			called = true
			var z0 Page[Account]
			var z1 error
			return z0, z1
		}))
		var a0 context.Context
		_, _ = s.List(a0)
		assert.True(t, called, "WithStoreList must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("List-failing")
		s := NewStoreStubFailing(t, want)
		var a0 context.Context
		_, r1 := s.List(a0)
		assert.ErrorIs(t, r1, want, "List must answer the armed error")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 Page[Account]
		var want1 error
		want0 = Page[Account]{Items: []Account{Account{ID: "ID-a"}}}
		s.OnList.Returns(want0, want1)
		var a0 context.Context
		_, _ = s.List(a0)
		s.ResetCalls()
		got0, _ := s.List(a0)
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// storeStubNoteSubject binds Note into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func storeStubNoteSubject(tb assert.TB) stub.Subject[StoreNoteCall, StoreNoteReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreNoteCall, StoreNoteReturn]{
		Stub: s.OnNote.MethodStub,
		Call: func() {
			var a0 string
			a0 = "Note-a"
			_ = s.Note(a0)
		},
		Result: func() StoreNoteReturn {
			var a0 string
			a0 = "Note-a"
			got0 := s.Note(a0)
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

// TestStoreStubNote pins how Note answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestStoreStubNote(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Note", storeStubNoteSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 Note
		want0 = Note{Text: "Text-a"}
		s.OnNote.Returns(want0)
		var a0 string
		a0 = "Note-a"
		got0 := s.Note(a0)
		assert.Equal(t, got0, want0, "Note must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var a0 string
		a0 = "Note-a"
		_ = s.Note(a0)
		got := s.OnNote.AssertCalledOnce(t, "Note must record the call")
		assert.Equal(t, got.ID, a0, "the recorded call carries ID")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var seen []StoreNoteCall
		s.OnNote.OnRecord(func(c StoreNoteCall) { seen = append(seen, c) })
		var a0 string
		a0 = "Note-a"
		_ = s.Note(a0)
		_ = s.Note(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Note call")
	})

	t.Run("wires WithStoreNote at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewStoreStub(t, WithStoreNote(func(_ string) Note {
			called = true
			var z0 Note
			return z0
		}))
		var a0 string
		a0 = "Note-a"
		_ = s.Note(a0)
		assert.True(t, called, "WithStoreNote must install the override")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 Note
		want0 = Note{Text: "Text-a"}
		s.OnNote.Returns(want0)
		var a0 string
		a0 = "Note-a"
		_ = s.Note(a0)
		s.ResetCalls()
		got0 := s.Note(a0)
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// storeStubBillSubject binds Bill into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func storeStubBillSubject(tb assert.TB) stub.Subject[StoreBillCall, StoreBillReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreBillCall, StoreBillReturn]{
		Stub: s.OnBill.MethodStub,
		Call: func() {
			var a0 string
			a0 = "Bill-a"
			_ = s.Bill(a0)
		},
		Result: func() StoreBillReturn {
			var a0 string
			a0 = "Bill-a"
			got0 := s.Bill(a0)
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

// TestStoreStubBill pins how Bill answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestStoreStubBill(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Bill", storeStubBillSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 billing.Invoice
		want0 = billing.Invoice{Number: "Number-a"}
		s.OnBill.Returns(want0)
		var a0 string
		a0 = "Bill-a"
		got0 := s.Bill(a0)
		assert.Equal(t, got0, want0, "Bill must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var a0 string
		a0 = "Bill-a"
		_ = s.Bill(a0)
		got := s.OnBill.AssertCalledOnce(t, "Bill must record the call")
		assert.Equal(t, got.ID, a0, "the recorded call carries ID")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var seen []StoreBillCall
		s.OnBill.OnRecord(func(c StoreBillCall) { seen = append(seen, c) })
		var a0 string
		a0 = "Bill-a"
		_ = s.Bill(a0)
		_ = s.Bill(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Bill call")
	})

	t.Run("wires WithStoreBill at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewStoreStub(t, WithStoreBill(func(_ string) billing.Invoice {
			called = true
			var z0 billing.Invoice
			return z0
		}))
		var a0 string
		a0 = "Bill-a"
		_ = s.Bill(a0)
		assert.True(t, called, "WithStoreBill must install the override")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 billing.Invoice
		want0 = billing.Invoice{Number: "Number-a"}
		s.OnBill.Returns(want0)
		var a0 string
		a0 = "Bill-a"
		_ = s.Bill(a0)
		s.ResetCalls()
		got0 := s.Bill(a0)
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// storeStubDouble describes how to build a StoreStub under each option
// whose effect is the same whatever a method's signature.
//
// Get stands in for the double as a whole: these checks assert only that
// an option reached it.
func storeStubDouble() stub.Double[StoreGetCall] {
	instance := func(s *StoreStub) stub.Instance[StoreGetCall] {
		return stub.Instance[StoreGetCall]{
			Stub: s.OnGet.MethodStub,
			Call: func() {
				var a0 context.Context
				var a1 string
				a1 = "Get-a"
				_, _ = s.Get(a0, a1)
			},
			Reset: s.ResetCalls,
		}
	}
	return stub.Double[StoreGetCall]{
		New: func(tb assert.TB) stub.Instance[StoreGetCall] {
			return instance(NewStoreStub(tb))
		},
		WithClock: func(tb assert.TB, clk assert.Clock) stub.Instance[StoreGetCall] {
			return instance(NewStoreStub(tb, StoreStubWithClock(clk)))
		},
		WithLatency: func(tb assert.TB, d time.Duration, clk assert.Clock) stub.Instance[StoreGetCall] {
			return instance(NewStoreStubDelayed(tb, d, clk))
		},
		WithRandSource: func(tb assert.TB, src rand.Source) stub.Instance[StoreGetCall] {
			return instance(NewStoreStub(tb, StoreStubWithRandSource(src)))
		},
		BenchMode: func(tb assert.TB) stub.Instance[StoreGetCall] {
			return instance(NewStoreStub(tb, StoreStubBenchMode()))
		},
		Strict: func(tb assert.TB) stub.Instance[StoreGetCall] {
			return instance(NewStoreStub(tb, StoreStubStrict()))
		},
	}
}

// TestStoreStubOptions pins the settings that apply to the double rather
// than to one method: bench mode, the injected clock and random source,
// and what a reset clears.
func TestStoreStubOptions(t *testing.T) {
	t.Parallel()

	stub.DoubleBehaviour(t, storeStubDouble())
}

// TestStoreStubDelegateTo pins that a wrapped implementation is called
// through, and recorded on the way past.
func TestStoreStubDelegateTo(t *testing.T) {
	t.Parallel()

	inner := NewStoreStub(t)
	s := NewStoreStub(t, StoreStubDelegateTo(inner))

	t.Run("forwards Get to the wrapped implementation", func(t *testing.T) {
		var a0 context.Context
		var a1 string
		a1 = "Get-a"
		_, _ = s.Get(a0, a1)
		inner.OnGet.AssertCalledOnce(t, "Get must reach the wrapped implementation")
	})

	t.Run("surfaces what Get answered", func(t *testing.T) {
		want := errors.New("Get-delegate")
		inner.OnGet.FaultsFor(time.Hour, want)
		var a0 context.Context
		var a1 string
		a1 = "Get-a"
		_, r1 := s.Get(a0, a1)
		assert.ErrorIs(t, r1, want, "Get must surface the wrapped answer")
	})

	t.Run("forwards Find to the wrapped implementation", func(t *testing.T) {
		var a0 context.Context
		var a1 string
		a1 = "Find-a"
		_, _ = s.Find(a0, a1)
		inner.OnFind.AssertCalledOnce(t, "Find must reach the wrapped implementation")
	})

	t.Run("surfaces what Find answered", func(t *testing.T) {
		want := errors.New("Find-delegate")
		inner.OnFind.FaultsFor(time.Hour, want)
		var a0 context.Context
		var a1 string
		a1 = "Find-a"
		_, r1 := s.Find(a0, a1)
		assert.ErrorIs(t, r1, want, "Find must surface the wrapped answer")
	})

	t.Run("forwards Lookup to the wrapped implementation", func(t *testing.T) {
		var a0 string
		a0 = "Lookup-a"
		_, _ = s.Lookup(a0)
		inner.OnLookup.AssertCalledOnce(t, "Lookup must reach the wrapped implementation")
	})

	t.Run("forwards List to the wrapped implementation", func(t *testing.T) {
		var a0 context.Context
		_, _ = s.List(a0)
		inner.OnList.AssertCalledOnce(t, "List must reach the wrapped implementation")
	})

	t.Run("surfaces what List answered", func(t *testing.T) {
		want := errors.New("List-delegate")
		inner.OnList.FaultsFor(time.Hour, want)
		var a0 context.Context
		_, r1 := s.List(a0)
		assert.ErrorIs(t, r1, want, "List must surface the wrapped answer")
	})

	t.Run("forwards Note to the wrapped implementation", func(t *testing.T) {
		var a0 string
		a0 = "Note-a"
		_ = s.Note(a0)
		inner.OnNote.AssertCalledOnce(t, "Note must reach the wrapped implementation")
	})

	t.Run("forwards Bill to the wrapped implementation", func(t *testing.T) {
		var a0 string
		a0 = "Bill-a"
		_ = s.Bill(a0)
		inner.OnBill.AssertCalledOnce(t, "Bill must reach the wrapped implementation")
	})
}
