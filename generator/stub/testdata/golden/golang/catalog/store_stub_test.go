package catalog

import (
	"context"
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

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
		Override: func(mark func()) {
			s.OnGet.Func(func(_ context.Context, _ string) (Item, error) {
				mark()
				var z0 Item
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
		var want0 Item
		var want1 error
		want0 = Item{ID: "ID-a"}
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
		s := NewStoreStub(t, WithStoreGet(func(_ context.Context, _ string) (Item, error) {
			called = true
			var z0 Item
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
		var want0 Item
		var want1 error
		want0 = Item{ID: "ID-a"}
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

// storeStubPutSubject binds Put into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func storeStubPutSubject(tb assert.TB) stub.Subject[StorePutCall, StorePutReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StorePutCall, StorePutReturn]{
		Stub: s.OnPut.MethodStub,
		Call: func() {
			var a0 context.Context
			var a1 Item
			a1 = Item{ID: "ID-a"}
			_ = s.Put(a0, a1)
		},
		Result: func() StorePutReturn {
			var a0 context.Context
			var a1 Item
			a1 = Item{ID: "ID-a"}
			got0 := s.Put(a0, a1)
			return StorePutReturn{Err: got0}
		},
		Override: func(mark func()) {
			s.OnPut.Func(func(_ context.Context, _ Item) error {
				mark()
				var z0 error
				return z0
			})
		},
		Fails: func() error {
			var a0 context.Context
			var a1 Item
			a1 = Item{ID: "ID-a"}
			r0 := s.Put(a0, a1)
			return r0
		},
	}
}

// TestStoreStubPut pins how Put answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestStoreStubPut(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Put", storeStubPutSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var a0 context.Context
		var a1 Item
		a1 = Item{ID: "ID-a"}
		_ = s.Put(a0, a1)
		got := s.OnPut.AssertCalledOnce(t, "Put must record the call")
		assert.Equal(t, got.A0, a0, "the recorded call carries A0")
		assert.Equal(t, got.A1, a1, "the recorded call carries A1")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var seen []StorePutCall
		s.OnPut.OnRecord(func(c StorePutCall) { seen = append(seen, c) })
		var a0 context.Context
		var a1 Item
		a1 = Item{ID: "ID-a"}
		_ = s.Put(a0, a1)
		_ = s.Put(a0, a1)
		assert.Length(t, seen, 2, "OnRecord must fire once per Put call")
	})

	t.Run("wires WithStorePut at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewStoreStub(t, WithStorePut(func(_ context.Context, _ Item) error {
			called = true
			var z0 error
			return z0
		}))
		var a0 context.Context
		var a1 Item
		a1 = Item{ID: "ID-a"}
		_ = s.Put(a0, a1)
		assert.True(t, called, "WithStorePut must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("Put-failing")
		s := NewStoreStubFailing(t, want)
		var a0 context.Context
		var a1 Item
		a1 = Item{ID: "ID-a"}
		r0 := s.Put(a0, a1)
		assert.ErrorIs(t, r0, want, "Put must answer the armed error")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 error
		s.OnPut.Returns(want0)
		var a0 context.Context
		var a1 Item
		a1 = Item{ID: "ID-a"}
		_ = s.Put(a0, a1)
		s.ResetCalls()
		_ = s.Put(a0, a1)
		s.OnPut.AssertCalledOnce(t, "Put must record the call that followed the reset")
	})
}

// storeStubLenSubject binds Len into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func storeStubLenSubject(tb assert.TB) stub.Subject[StoreLenCall, StoreLenReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreLenCall, StoreLenReturn]{
		Stub: s.OnLen.MethodStub,
		Call: func() {
			_ = s.Len()
		},
		Result: func() StoreLenReturn {
			got0 := s.Len()
			return StoreLenReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnLen.Func(func() int {
				mark()
				var z0 int
				return z0
			})
		},
	}
}

// TestStoreStubLen pins how Len answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestStoreStubLen(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Len", storeStubLenSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 int
		want0 = 42
		s.OnLen.Returns(want0)
		got0 := s.Len()
		assert.Equal(t, got0, want0, "Len must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		_ = s.Len()
		s.OnLen.AssertCalledOnce(t, "Len must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var seen []StoreLenCall
		s.OnLen.OnRecord(func(c StoreLenCall) { seen = append(seen, c) })
		_ = s.Len()
		_ = s.Len()
		assert.Length(t, seen, 2, "OnRecord must fire once per Len call")
	})

	t.Run("wires WithStoreLen at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewStoreStub(t, WithStoreLen(func() int {
			called = true
			var z0 int
			return z0
		}))
		_ = s.Len()
		assert.True(t, called, "WithStoreLen must install the override")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 int
		want0 = 42
		s.OnLen.Returns(want0)
		_ = s.Len()
		s.ResetCalls()
		got0 := s.Len()
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
			return StoreLookupReturn{Item: got0, Ok: got1}
		},
		Override: func(mark func()) {
			s.OnLookup.Func(func(_ string) (Item, bool) {
				mark()
				var z0 Item
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
		var want0 Item
		var want1 bool
		want0 = Item{ID: "ID-a"}
		want1 = true
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
		assert.Equal(t, got.Key, a0, "the recorded call carries Key")
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
		s := NewStoreStub(t, WithStoreLookup(func(_ string) (Item, bool) {
			called = true
			var z0 Item
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
		var want0 Item
		var want1 bool
		want0 = Item{ID: "ID-a"}
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

// storeStubLogSubject binds Log into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func storeStubLogSubject(tb assert.TB) stub.Subject[StoreLogCall, struct{}] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreLogCall, struct{}]{
		Stub: s.OnLog.MethodStub,
		Call: func() {
			var a0 string
			var a1 []string
			a0 = "Log-a"
			s.Log(a0, a1...)
		},
		Override: func(mark func()) {
			s.OnLog.Func(func(_ string, _ ...string) {
				mark()
			})
		},
	}
}

// TestStoreStubLog pins how Log answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestStoreStubLog(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Log", storeStubLogSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var a0 string
		var a1 []string
		a0 = "Log-a"
		s.Log(a0, a1...)
		got := s.OnLog.AssertCalledOnce(t, "Log must record the call")
		assert.Equal(t, got.Format, a0, "the recorded call carries Format")
		assert.Equal(t, got.Args, a1, "the recorded call carries Args")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var seen []StoreLogCall
		s.OnLog.OnRecord(func(c StoreLogCall) { seen = append(seen, c) })
		var a0 string
		var a1 []string
		a0 = "Log-a"
		s.Log(a0, a1...)
		s.Log(a0, a1...)
		assert.Length(t, seen, 2, "OnRecord must fire once per Log call")
	})

	t.Run("wires WithStoreLog at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewStoreStub(t, WithStoreLog(func(_ string, _ ...string) {
			called = true
		}))
		var a0 string
		var a1 []string
		a0 = "Log-a"
		s.Log(a0, a1...)
		assert.True(t, called, "WithStoreLog must install the override")
	})
}

// storeStubSplitSubject binds Split into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func storeStubSplitSubject(tb assert.TB) stub.Subject[StoreSplitCall, StoreSplitReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreSplitCall, StoreSplitReturn]{
		Stub: s.OnSplit.MethodStub,
		Call: func() {
			var a0 string
			a0 = "Split-a"
			_, _, _ = s.Split(a0)
		},
		Result: func() StoreSplitReturn {
			var a0 string
			a0 = "Split-a"
			got0, got1, got2 := s.Split(a0)
			return StoreSplitReturn{Result0: got0, Result1: got1, Err: got2}
		},
		Override: func(mark func()) {
			s.OnSplit.Func(func(_ string) (string, string, error) {
				mark()
				var z0 string
				var z1 string
				var z2 error
				return z0, z1, z2
			})
		},
		Fails: func() error {
			var a0 string
			a0 = "Split-a"
			_, _, r2 := s.Split(a0)
			return r2
		},
	}
}

// TestStoreStubSplit pins how Split answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestStoreStubSplit(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Split", storeStubSplitSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 string
		var want1 string
		var want2 error
		want0 = "Split-a"
		want1 = "Split-a"
		s.OnSplit.Returns(want0, want1, want2)
		var a0 string
		a0 = "Split-a"
		got0, got1, _ := s.Split(a0)
		assert.Equal(t, got0, want0, "Split must answer with what Returns pinned")
		assert.Equal(t, got1, want1, "Split must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var a0 string
		a0 = "Split-a"
		_, _, _ = s.Split(a0)
		got := s.OnSplit.AssertCalledOnce(t, "Split must record the call")
		assert.Equal(t, got.S, a0, "the recorded call carries S")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var seen []StoreSplitCall
		s.OnSplit.OnRecord(func(c StoreSplitCall) { seen = append(seen, c) })
		var a0 string
		a0 = "Split-a"
		_, _, _ = s.Split(a0)
		_, _, _ = s.Split(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Split call")
	})

	t.Run("wires WithStoreSplit at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewStoreStub(t, WithStoreSplit(func(_ string) (string, string, error) {
			called = true
			var z0 string
			var z1 string
			var z2 error
			return z0, z1, z2
		}))
		var a0 string
		a0 = "Split-a"
		_, _, _ = s.Split(a0)
		assert.True(t, called, "WithStoreSplit must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("Split-failing")
		s := NewStoreStubFailing(t, want)
		var a0 string
		a0 = "Split-a"
		_, _, r2 := s.Split(a0)
		assert.ErrorIs(t, r2, want, "Split must answer the armed error")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 string
		var want1 string
		var want2 error
		want0 = "Split-a"
		want1 = "Split-a"
		s.OnSplit.Returns(want0, want1, want2)
		var a0 string
		a0 = "Split-a"
		_, _, _ = s.Split(a0)
		s.ResetCalls()
		got0, got1, _ := s.Split(a0)
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
		assert.Equal(t, got1, want1, "a reset must keep what Returns pinned")
	})
}

// storeStubCloseSubject binds Close into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func storeStubCloseSubject(tb assert.TB) stub.Subject[StoreCloseCall, StoreCloseReturn] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreCloseCall, StoreCloseReturn]{
		Stub: s.OnClose.MethodStub,
		Call: func() {
			_ = s.Close()
		},
		Result: func() StoreCloseReturn {
			got0 := s.Close()
			return StoreCloseReturn{Err: got0}
		},
		Override: func(mark func()) {
			s.OnClose.Func(func() error {
				mark()
				var z0 error
				return z0
			})
		},
		Fails: func() error {
			r0 := s.Close()
			return r0
		},
	}
}

// TestStoreStubClose pins how Close answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestStoreStubClose(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Close", storeStubCloseSubject)

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		_ = s.Close()
		s.OnClose.AssertCalledOnce(t, "Close must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var seen []StoreCloseCall
		s.OnClose.OnRecord(func(c StoreCloseCall) { seen = append(seen, c) })
		_ = s.Close()
		_ = s.Close()
		assert.Length(t, seen, 2, "OnRecord must fire once per Close call")
	})

	t.Run("wires WithStoreClose at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewStoreStub(t, WithStoreClose(func() error {
			called = true
			var z0 error
			return z0
		}))
		_ = s.Close()
		assert.True(t, called, "WithStoreClose must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("Close-failing")
		s := NewStoreStubFailing(t, want)
		r0 := s.Close()
		assert.ErrorIs(t, r0, want, "Close must answer the armed error")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 error
		s.OnClose.Returns(want0)
		_ = s.Close()
		s.ResetCalls()
		_ = s.Close()
		s.OnClose.AssertCalledOnce(t, "Close must record the call that followed the reset")
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

	t.Run("forwards Put to the wrapped implementation", func(t *testing.T) {
		var a0 context.Context
		var a1 Item
		a1 = Item{ID: "ID-a"}
		_ = s.Put(a0, a1)
		inner.OnPut.AssertCalledOnce(t, "Put must reach the wrapped implementation")
	})

	t.Run("surfaces what Put answered", func(t *testing.T) {
		want := errors.New("Put-delegate")
		inner.OnPut.FaultsFor(time.Hour, want)
		var a0 context.Context
		var a1 Item
		a1 = Item{ID: "ID-a"}
		r0 := s.Put(a0, a1)
		assert.ErrorIs(t, r0, want, "Put must surface the wrapped answer")
	})

	t.Run("forwards Len to the wrapped implementation", func(t *testing.T) {
		_ = s.Len()
		inner.OnLen.AssertCalledOnce(t, "Len must reach the wrapped implementation")
	})

	t.Run("forwards Lookup to the wrapped implementation", func(t *testing.T) {
		var a0 string
		a0 = "Lookup-a"
		_, _ = s.Lookup(a0)
		inner.OnLookup.AssertCalledOnce(t, "Lookup must reach the wrapped implementation")
	})

	t.Run("forwards Log to the wrapped implementation", func(t *testing.T) {
		var a0 string
		var a1 []string
		a0 = "Log-a"
		s.Log(a0, a1...)
		inner.OnLog.AssertCalledOnce(t, "Log must reach the wrapped implementation")
	})

	t.Run("forwards Split to the wrapped implementation", func(t *testing.T) {
		var a0 string
		a0 = "Split-a"
		_, _, _ = s.Split(a0)
		inner.OnSplit.AssertCalledOnce(t, "Split must reach the wrapped implementation")
	})

	t.Run("surfaces what Split answered", func(t *testing.T) {
		want := errors.New("Split-delegate")
		inner.OnSplit.FaultsFor(time.Hour, want)
		var a0 string
		a0 = "Split-a"
		_, _, r2 := s.Split(a0)
		assert.ErrorIs(t, r2, want, "Split must surface the wrapped answer")
	})

	t.Run("forwards Close to the wrapped implementation", func(t *testing.T) {
		_ = s.Close()
		inner.OnClose.AssertCalledOnce(t, "Close must reach the wrapped implementation")
	})

	t.Run("surfaces what Close answered", func(t *testing.T) {
		want := errors.New("Close-delegate")
		inner.OnClose.FaultsFor(time.Hour, want)
		r0 := s.Close()
		assert.ErrorIs(t, r0, want, "Close must surface the wrapped answer")
	})
}

// clockStubNowSubject binds Now into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func clockStubNowSubject(tb assert.TB) stub.Subject[ClockNowCall, ClockNowReturn] {
	tb.Helper()
	s := NewClockStub(tb)
	return stub.Subject[ClockNowCall, ClockNowReturn]{
		Stub: s.OnNow.MethodStub,
		Call: func() {
			_ = s.Now()
		},
		Result: func() ClockNowReturn {
			got0 := s.Now()
			return ClockNowReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnNow.Func(func() time.Time {
				mark()
				var z0 time.Time
				return z0
			})
		},
	}
}

// TestClockStubNow pins how Now answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestClockStubNow(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Now", clockStubNowSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewClockStub(t)
		var want0 time.Time
		want0 = time.Unix(1700000000, 0)
		s.OnNow.Returns(want0)
		got0 := s.Now()
		assert.Equal(t, got0, want0, "Now must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewClockStub(t)
		_ = s.Now()
		s.OnNow.AssertCalledOnce(t, "Now must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewClockStub(t)
		var seen []ClockNowCall
		s.OnNow.OnRecord(func(c ClockNowCall) { seen = append(seen, c) })
		_ = s.Now()
		_ = s.Now()
		assert.Length(t, seen, 2, "OnRecord must fire once per Now call")
	})

	t.Run("wires WithClockNow at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewClockStub(t, WithClockNow(func() time.Time {
			called = true
			var z0 time.Time
			return z0
		}))
		_ = s.Now()
		assert.True(t, called, "WithClockNow must install the override")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewClockStub(t)
		var want0 time.Time
		want0 = time.Unix(1700000000, 0)
		s.OnNow.Returns(want0)
		_ = s.Now()
		s.ResetCalls()
		got0 := s.Now()
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// clockStubZoneSubject binds Zone into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func clockStubZoneSubject(tb assert.TB) stub.Subject[ClockZoneCall, ClockZoneReturn] {
	tb.Helper()
	s := NewClockStub(tb)
	return stub.Subject[ClockZoneCall, ClockZoneReturn]{
		Stub: s.OnZone.MethodStub,
		Call: func() {
			_ = s.Zone()
		},
		Result: func() ClockZoneReturn {
			got0 := s.Zone()
			return ClockZoneReturn{Result: got0}
		},
		Override: func(mark func()) {
			s.OnZone.Func(func() string {
				mark()
				var z0 string
				return z0
			})
		},
	}
}

// TestClockStubZone pins how Zone answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestClockStubZone(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Zone", clockStubZoneSubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewClockStub(t)
		var want0 string
		want0 = "Zone-a"
		s.OnZone.Returns(want0)
		got0 := s.Zone()
		assert.Equal(t, got0, want0, "Zone must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewClockStub(t)
		_ = s.Zone()
		s.OnZone.AssertCalledOnce(t, "Zone must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewClockStub(t)
		var seen []ClockZoneCall
		s.OnZone.OnRecord(func(c ClockZoneCall) { seen = append(seen, c) })
		_ = s.Zone()
		_ = s.Zone()
		assert.Length(t, seen, 2, "OnRecord must fire once per Zone call")
	})

	t.Run("wires WithClockZone at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewClockStub(t, WithClockZone(func() string {
			called = true
			var z0 string
			return z0
		}))
		_ = s.Zone()
		assert.True(t, called, "WithClockZone must install the override")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewClockStub(t)
		var want0 string
		want0 = "Zone-a"
		s.OnZone.Returns(want0)
		_ = s.Zone()
		s.ResetCalls()
		got0 := s.Zone()
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// clockStubEverySubject binds Every into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func clockStubEverySubject(tb assert.TB) stub.Subject[ClockEveryCall, ClockEveryReturn] {
	tb.Helper()
	s := NewClockStub(tb)
	return stub.Subject[ClockEveryCall, ClockEveryReturn]{
		Stub: s.OnEvery.MethodStub,
		Call: func() {
			var a0 time.Duration
			var a1 func()
			a0 = time.Duration(1)
			a1 = func() {}
			_ = s.Every(a0, a1)
		},
		Result: func() ClockEveryReturn {
			var a0 time.Duration
			var a1 func()
			a0 = time.Duration(1)
			a1 = func() {}
			got0 := s.Every(a0, a1)
			return ClockEveryReturn{Stop: got0}
		},
		Override: func(mark func()) {
			s.OnEvery.Func(func(_ time.Duration, _ func()) func() {
				mark()
				var z0 func()
				return z0
			})
		},
	}
}

// TestClockStubEvery pins how Every answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func TestClockStubEvery(t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Every", clockStubEverySubject)

	t.Run("answers with the value pinned by Returns", func(t *testing.T) {
		t.Parallel()
		s := NewClockStub(t)
		var want0 func()
		want0 = func() {}
		s.OnEvery.Returns(want0)
		var a0 time.Duration
		var a1 func()
		a0 = time.Duration(1)
		a1 = func() {}
		got0 := s.Every(a0, a1)
		assert.Equal(t, got0, want0, "Every must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewClockStub(t)
		var a0 time.Duration
		var a1 func()
		a0 = time.Duration(1)
		a1 = func() {}
		_ = s.Every(a0, a1)
		got := s.OnEvery.AssertCalledOnce(t, "Every must record the call")
		assert.Equal(t, got.D, a0, "the recorded call carries D")
		assert.Equal(t, got.Fn, a1, "the recorded call carries Fn")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewClockStub(t)
		var seen []ClockEveryCall
		s.OnEvery.OnRecord(func(c ClockEveryCall) { seen = append(seen, c) })
		var a0 time.Duration
		var a1 func()
		a0 = time.Duration(1)
		a1 = func() {}
		_ = s.Every(a0, a1)
		_ = s.Every(a0, a1)
		assert.Length(t, seen, 2, "OnRecord must fire once per Every call")
	})

	t.Run("wires WithClockEvery at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewClockStub(t, WithClockEvery(func(_ time.Duration, _ func()) func() {
			called = true
			var z0 func()
			return z0
		}))
		var a0 time.Duration
		var a1 func()
		a0 = time.Duration(1)
		a1 = func() {}
		_ = s.Every(a0, a1)
		assert.True(t, called, "WithClockEvery must install the override")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewClockStub(t)
		var want0 func()
		want0 = func() {}
		s.OnEvery.Returns(want0)
		var a0 time.Duration
		var a1 func()
		a0 = time.Duration(1)
		a1 = func() {}
		_ = s.Every(a0, a1)
		s.ResetCalls()
		got0 := s.Every(a0, a1)
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// clockStubDouble describes how to build a ClockStub under each option
// whose effect is the same whatever a method's signature.
//
// Now stands in for the double as a whole: these checks assert only that
// an option reached it.
func clockStubDouble() stub.Double[ClockNowCall] {
	instance := func(s *ClockStub) stub.Instance[ClockNowCall] {
		return stub.Instance[ClockNowCall]{
			Stub: s.OnNow.MethodStub,
			Call: func() {
				_ = s.Now()
			},
			Reset: s.ResetCalls,
		}
	}
	return stub.Double[ClockNowCall]{
		New: func(tb assert.TB) stub.Instance[ClockNowCall] {
			return instance(NewClockStub(tb))
		},
		WithClock: func(tb assert.TB, clk assert.Clock) stub.Instance[ClockNowCall] {
			return instance(NewClockStub(tb, ClockStubWithClock(clk)))
		},
		WithLatency: func(tb assert.TB, d time.Duration, clk assert.Clock) stub.Instance[ClockNowCall] {
			return instance(NewClockStubDelayed(tb, d, clk))
		},
		WithRandSource: func(tb assert.TB, src rand.Source) stub.Instance[ClockNowCall] {
			return instance(NewClockStub(tb, ClockStubWithRandSource(src)))
		},
		BenchMode: func(tb assert.TB) stub.Instance[ClockNowCall] {
			return instance(NewClockStub(tb, ClockStubBenchMode()))
		},
		Strict: func(tb assert.TB) stub.Instance[ClockNowCall] {
			return instance(NewClockStub(tb, ClockStubStrict()))
		},
	}
}

// TestClockStubOptions pins the settings that apply to the double rather
// than to one method: bench mode, the injected clock and random source,
// and what a reset clears.
func TestClockStubOptions(t *testing.T) {
	t.Parallel()

	stub.DoubleBehaviour(t, clockStubDouble())
}
