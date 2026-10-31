package catalog_test

import (
	"context"
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"example.com/shop/catalog"
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
func storeStubGet(tb assert.TB, s *catalog.StoreStub) (catalog.Item, error) {
	tb.Helper()
	a0, a1 := storeStubGetArgs(tb)
	return s.Get(a0, a1)
}

// storeStubGetPins returns the samples of the results of Get, with a new
// error as its error result.
func storeStubGetPins() (catalog.Item, error) {
	var want0 catalog.Item
	want1 := errors.New("Get-returns")
	want0 = catalog.Item{ID: "ID-a"}
	return want0, want1
}

// storeStubPutArgs returns the samples of the arguments of Put.
func storeStubPutArgs(tb assert.TB) (context.Context, catalog.Item) {
	tb.Helper()
	a0 := assert.Context(tb)
	var a1 catalog.Item
	a1 = catalog.Item{ID: "ID-a"}
	return a0, a1
}

// storeStubPut calls Put of s with the samples of its arguments.
func storeStubPut(tb assert.TB, s *catalog.StoreStub) error {
	tb.Helper()
	a0, a1 := storeStubPutArgs(tb)
	return s.Put(a0, a1)
}

// storeStubPutPins returns the samples of the results of Put, with a new
// error as its error result.
func storeStubPutPins() error {
	want0 := errors.New("Put-returns")
	return want0
}

// storeStubLen calls Len of s with the samples of its arguments.
func storeStubLen(tb assert.TB, s *catalog.StoreStub) int {
	tb.Helper()
	return s.Len()
}

// storeStubLenPins returns the samples of the results of Len.
func storeStubLenPins() int {
	var want0 int
	want0 = 42
	return want0
}

// storeStubLookupArgs returns the samples of the arguments of Lookup.
func storeStubLookupArgs(tb assert.TB) string {
	tb.Helper()
	var a0 string
	a0 = "Lookup-a"
	return a0
}

// storeStubLookup calls Lookup of s with the samples of its arguments.
func storeStubLookup(tb assert.TB, s *catalog.StoreStub) (catalog.Item, bool) {
	tb.Helper()
	a0 := storeStubLookupArgs(tb)
	return s.Lookup(a0)
}

// storeStubLookupPins returns the samples of the results of Lookup.
func storeStubLookupPins() (catalog.Item, bool) {
	var want0 catalog.Item
	var want1 bool
	want0 = catalog.Item{ID: "ID-a"}
	want1 = true
	return want0, want1
}

// storeStubLogArgs returns the samples of the arguments of Log.
func storeStubLogArgs(tb assert.TB) (string, []string) {
	tb.Helper()
	var a0 string
	var a1 []string
	a0 = "Log-a"
	return a0, a1
}

// storeStubLog calls Log of s with the samples of its arguments.
func storeStubLog(tb assert.TB, s *catalog.StoreStub) {
	tb.Helper()
	a0, a1 := storeStubLogArgs(tb)
	s.Log(a0, a1...)
}

// storeStubSplitArgs returns the samples of the arguments of Split.
func storeStubSplitArgs(tb assert.TB) string {
	tb.Helper()
	var a0 string
	a0 = "Split-a"
	return a0
}

// storeStubSplit calls Split of s with the samples of its arguments.
func storeStubSplit(tb assert.TB, s *catalog.StoreStub) (string, string, error) {
	tb.Helper()
	a0 := storeStubSplitArgs(tb)
	return s.Split(a0)
}

// storeStubSplitPins returns the samples of the results of Split, with a
// new error as its error result.
func storeStubSplitPins() (string, string, error) {
	var want0 string
	var want1 string
	want2 := errors.New("Split-returns")
	want0 = "Split-a"
	want1 = "Split-a"
	return want0, want1, want2
}

// storeStubClose calls Close of s with the samples of its arguments.
func storeStubClose(tb assert.TB, s *catalog.StoreStub) error {
	tb.Helper()
	return s.Close()
}

// storeStubClosePins returns the samples of the results of Close, with a
// new error as its error result.
func storeStubClosePins() error {
	want0 := errors.New("Close-returns")
	return want0
}

// storeStubGetSubject returns the subject of [stub.Behaviour] for Get. The
// subject calls a new double that reports to tb.
func storeStubGetSubject(tb assert.TB) stub.Subject[catalog.StoreGetCall, catalog.StoreGetReturn] {
	tb.Helper()
	s := catalog.NewStoreStub(tb)
	return stub.Subject[catalog.StoreGetCall, catalog.StoreGetReturn]{
		Stub: s.OnGet.MethodStub,
		Call: func() {
			_, _ = storeStubGet(tb, s)
		},
		Result: func() catalog.StoreGetReturn {
			got0, got1 := storeStubGet(tb, s)
			return catalog.StoreGetReturn{Result: got0, Err: got1}
		},
		Override: func(mark func()) {
			s.OnGet.Func(func(_ context.Context, _ string) (catalog.Item, error) {
				mark()
				var z0 catalog.Item
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

// storeStubPutSubject returns the subject of [stub.Behaviour] for Put. The
// subject calls a new double that reports to tb.
func storeStubPutSubject(tb assert.TB) stub.Subject[catalog.StorePutCall, catalog.StorePutReturn] {
	tb.Helper()
	s := catalog.NewStoreStub(tb)
	return stub.Subject[catalog.StorePutCall, catalog.StorePutReturn]{
		Stub: s.OnPut.MethodStub,
		Call: func() {
			_ = storeStubPut(tb, s)
		},
		Result: func() catalog.StorePutReturn {
			got0 := storeStubPut(tb, s)
			return catalog.StorePutReturn{Err: got0}
		},
		Override: func(mark func()) {
			s.OnPut.Func(func(_ context.Context, _ catalog.Item) error {
				mark()
				var z0 error
				return z0
			})
		},
		Fails: func() error {
			r0 := storeStubPut(tb, s)
			return r0
		},
	}
}

// storeStubLenSubject returns the subject of [stub.Behaviour] for Len. The
// subject calls a new double that reports to tb.
func storeStubLenSubject(tb assert.TB) stub.Subject[catalog.StoreLenCall, catalog.StoreLenReturn] {
	tb.Helper()
	s := catalog.NewStoreStub(tb)
	return stub.Subject[catalog.StoreLenCall, catalog.StoreLenReturn]{
		Stub: s.OnLen.MethodStub,
		Call: func() {
			_ = storeStubLen(tb, s)
		},
		Result: func() catalog.StoreLenReturn {
			got0 := storeStubLen(tb, s)
			return catalog.StoreLenReturn{Result: got0}
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

// storeStubLookupSubject returns the subject of [stub.Behaviour] for
// Lookup. The subject calls a new double that reports to tb.
func storeStubLookupSubject(
	tb assert.TB,
) stub.Subject[catalog.StoreLookupCall, catalog.StoreLookupReturn] {
	tb.Helper()
	s := catalog.NewStoreStub(tb)
	return stub.Subject[catalog.StoreLookupCall, catalog.StoreLookupReturn]{
		Stub: s.OnLookup.MethodStub,
		Call: func() {
			_, _ = storeStubLookup(tb, s)
		},
		Result: func() catalog.StoreLookupReturn {
			got0, got1 := storeStubLookup(tb, s)
			return catalog.StoreLookupReturn{Item: got0, Ok: got1}
		},
		Override: func(mark func()) {
			s.OnLookup.Func(func(_ string) (catalog.Item, bool) {
				mark()
				var z0 catalog.Item
				var z1 bool
				return z0, z1
			})
		},
	}
}

// storeStubLogSubject returns the subject of [stub.Behaviour] for Log. The
// subject calls a new double that reports to tb.
func storeStubLogSubject(tb assert.TB) stub.Subject[catalog.StoreLogCall, struct{}] {
	tb.Helper()
	s := catalog.NewStoreStub(tb)
	return stub.Subject[catalog.StoreLogCall, struct{}]{
		Stub: s.OnLog.MethodStub,
		Call: func() {
			storeStubLog(tb, s)
		},
		Override: func(mark func()) {
			s.OnLog.Func(func(_ string, _ ...string) {
				mark()
			})
		},
	}
}

// storeStubSplitSubject returns the subject of [stub.Behaviour] for Split.
// The subject calls a new double that reports to tb.
func storeStubSplitSubject(
	tb assert.TB,
) stub.Subject[catalog.StoreSplitCall, catalog.StoreSplitReturn] {
	tb.Helper()
	s := catalog.NewStoreStub(tb)
	return stub.Subject[catalog.StoreSplitCall, catalog.StoreSplitReturn]{
		Stub: s.OnSplit.MethodStub,
		Call: func() {
			_, _, _ = storeStubSplit(tb, s)
		},
		Result: func() catalog.StoreSplitReturn {
			got0, got1, got2 := storeStubSplit(tb, s)
			return catalog.StoreSplitReturn{Result0: got0, Result1: got1, Err: got2}
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
			_, _, r2 := storeStubSplit(tb, s)
			return r2
		},
	}
}

// storeStubCloseSubject returns the subject of [stub.Behaviour] for Close.
// The subject calls a new double that reports to tb.
func storeStubCloseSubject(
	tb assert.TB,
) stub.Subject[catalog.StoreCloseCall, catalog.StoreCloseReturn] {
	tb.Helper()
	s := catalog.NewStoreStub(tb)
	return stub.Subject[catalog.StoreCloseCall, catalog.StoreCloseReturn]{
		Stub: s.OnClose.MethodStub,
		Call: func() {
			_ = storeStubClose(tb, s)
		},
		Result: func() catalog.StoreCloseReturn {
			got0 := storeStubClose(tb, s)
			return catalog.StoreCloseReturn{Err: got0}
		},
		Override: func(mark func()) {
			s.OnClose.Func(func() error {
				mark()
				var z0 error
				return z0
			})
		},
		Fails: func() error {
			r0 := storeStubClose(tb, s)
			return r0
		},
	}
}

// storeStubInstance is the instance of a catalog.StoreStub that the checks
// of the constructors and the options drive.
type storeStubInstance = stub.Instance[catalog.StoreGetCall]

// storeStubInstanceOf returns the instance of s. The options apply to the
// whole double, so the checks call Get alone.
func storeStubInstanceOf(tb assert.TB, s *catalog.StoreStub) storeStubInstance {
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
			s := catalog.NewStoreStub(t)
			_, _ = storeStubGet(t, s)
			got := s.OnGet.AssertCalledOnce(t, "Get records the call")
			a0, a1 := storeStubGetArgs(t)
			expect.Equal(t, got.Ctx, a0, "the call records Ctx", expect.ByIdentity())
			expect.Equal(t, got.ID, a1, "the call records ID")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			var seen []catalog.StoreGetCall
			s.OnGet.OnRecord(func(c catalog.StoreGetCall) { seen = append(seen, c) })
			_, _ = storeStubGet(t, s)
			_, _ = storeStubGet(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Get")
		})
	})

	t.Run("Put", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Put", storeStubPutSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			_ = storeStubPut(t, s)
			got := s.OnPut.AssertCalledOnce(t, "Put records the call")
			a0, a1 := storeStubPutArgs(t)
			expect.Equal(t, got.A0, a0, "the call records A0", expect.ByIdentity())
			expect.Equal(t, got.A1, a1, "the call records A1")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			var seen []catalog.StorePutCall
			s.OnPut.OnRecord(func(c catalog.StorePutCall) { seen = append(seen, c) })
			_ = storeStubPut(t, s)
			_ = storeStubPut(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Put")
		})
	})

	t.Run("Len", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Len", storeStubLenSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			var seen []catalog.StoreLenCall
			s.OnLen.OnRecord(func(c catalog.StoreLenCall) { seen = append(seen, c) })
			_ = storeStubLen(t, s)
			_ = storeStubLen(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Len")
		})
	})

	t.Run("Lookup", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Lookup", storeStubLookupSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			_, _ = storeStubLookup(t, s)
			got := s.OnLookup.AssertCalledOnce(t, "Lookup records the call")
			a0 := storeStubLookupArgs(t)
			assert.Equal(t, got.Key, a0, "the call records Key")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			var seen []catalog.StoreLookupCall
			s.OnLookup.OnRecord(func(c catalog.StoreLookupCall) { seen = append(seen, c) })
			_, _ = storeStubLookup(t, s)
			_, _ = storeStubLookup(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Lookup")
		})
	})

	t.Run("Log", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Log", storeStubLogSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			storeStubLog(t, s)
			got := s.OnLog.AssertCalledOnce(t, "Log records the call")
			a0, a1 := storeStubLogArgs(t)
			expect.Equal(t, got.Format, a0, "the call records Format")
			expect.Equal(t, got.Args, a1, "the call records Args")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			var seen []catalog.StoreLogCall
			s.OnLog.OnRecord(func(c catalog.StoreLogCall) { seen = append(seen, c) })
			storeStubLog(t, s)
			storeStubLog(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Log")
		})
	})

	t.Run("Split", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Split", storeStubSplitSubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			_, _, _ = storeStubSplit(t, s)
			got := s.OnSplit.AssertCalledOnce(t, "Split records the call")
			a0 := storeStubSplitArgs(t)
			assert.Equal(t, got.S, a0, "the call records S")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			var seen []catalog.StoreSplitCall
			s.OnSplit.OnRecord(func(c catalog.StoreSplitCall) { seen = append(seen, c) })
			_, _, _ = storeStubSplit(t, s)
			_, _, _ = storeStubSplit(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Split")
		})
	})

	t.Run("Close", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Close", storeStubCloseSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			var seen []catalog.StoreCloseCall
			s.OnClose.OnRecord(func(c catalog.StoreCloseCall) { seen = append(seen, c) })
			_ = storeStubClose(t, s)
			_ = storeStubClose(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Close")
		})
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()
		stub.ResetBehaviour(t, func(tb assert.TB) storeStubInstance {
			return storeStubInstanceOf(tb, catalog.NewStoreStub(tb))
		})

		t.Run("clears the recorded calls of Get", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			_, _ = storeStubGet(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnGet.CallCount(), 0, "ResetCalls clears the calls of Get")
		})

		t.Run("keeps the results of Returns of Get", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			want0, want1 := storeStubGetPins()
			s.OnGet.Returns(want0, want1)
			s.ResetCalls()
			got0, got1 := storeStubGet(t, s)
			expect.Equal(t, got0, want0, "ResetCalls keeps the Result of Get")
			expect.ErrorIs(t, got1, want1, "ResetCalls keeps the error of Get")
		})

		t.Run("clears the recorded calls of Put", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			_ = storeStubPut(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnPut.CallCount(), 0, "ResetCalls clears the calls of Put")
		})

		t.Run("keeps the results of Returns of Put", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			want0 := storeStubPutPins()
			s.OnPut.Returns(want0)
			s.ResetCalls()
			got0 := storeStubPut(t, s)
			assert.ErrorIs(t, got0, want0, "ResetCalls keeps the error of Put")
		})

		t.Run("clears the recorded calls of Len", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			_ = storeStubLen(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnLen.CallCount(), 0, "ResetCalls clears the calls of Len")
		})

		t.Run("keeps the results of Returns of Len", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			want0 := storeStubLenPins()
			s.OnLen.Returns(want0)
			s.ResetCalls()
			got0 := storeStubLen(t, s)
			assert.Equal(t, got0, want0, "ResetCalls keeps the Result of Len")
		})

		t.Run("clears the recorded calls of Lookup", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			_, _ = storeStubLookup(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnLookup.CallCount(), 0, "ResetCalls clears the calls of Lookup")
		})

		t.Run("keeps the results of Returns of Lookup", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			want0, want1 := storeStubLookupPins()
			s.OnLookup.Returns(want0, want1)
			s.ResetCalls()
			got0, got1 := storeStubLookup(t, s)
			expect.Equal(t, got0, want0, "ResetCalls keeps the Item of Lookup")
			expect.Equal(t, got1, want1, "ResetCalls keeps the Ok of Lookup")
		})

		t.Run("clears the recorded calls of Log", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			storeStubLog(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnLog.CallCount(), 0, "ResetCalls clears the calls of Log")
		})

		t.Run("clears the recorded calls of Split", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			_, _, _ = storeStubSplit(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnSplit.CallCount(), 0, "ResetCalls clears the calls of Split")
		})

		t.Run("keeps the results of Returns of Split", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			want0, want1, want2 := storeStubSplitPins()
			s.OnSplit.Returns(want0, want1, want2)
			s.ResetCalls()
			got0, got1, got2 := storeStubSplit(t, s)
			expect.Equal(t, got0, want0, "ResetCalls keeps the Result0 of Split")
			expect.Equal(t, got1, want1, "ResetCalls keeps the Result1 of Split")
			expect.ErrorIs(t, got2, want2, "ResetCalls keeps the error of Split")
		})

		t.Run("clears the recorded calls of Close", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			_ = storeStubClose(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnClose.CallCount(), 0, "ResetCalls clears the calls of Close")
		})

		t.Run("keeps the results of Returns of Close", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			want0 := storeStubClosePins()
			s.OnClose.Returns(want0)
			s.ResetCalls()
			got0 := storeStubClose(t, s)
			assert.ErrorIs(t, got0, want0, "ResetCalls keeps the error of Close")
		})
	})

	t.Run("NewStoreStubDelayed", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) storeStubInstance {
				return storeStubInstanceOf(tb, catalog.NewStoreStubDelayed(tb, d, clk))
			},
		)
	})

	t.Run("NewStoreStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Get with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Get-failing")
			s := catalog.NewStoreStubFailing(t, want)
			_, r1 := storeStubGet(t, s)
			assert.ErrorIs(t, r1, want, "Get fails with the error of NewStoreStubFailing")
		})

		t.Run("fails Put with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Put-failing")
			s := catalog.NewStoreStubFailing(t, want)
			r0 := storeStubPut(t, s)
			assert.ErrorIs(t, r0, want, "Put fails with the error of NewStoreStubFailing")
		})

		t.Run("fails Split with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Split-failing")
			s := catalog.NewStoreStubFailing(t, want)
			_, _, r2 := storeStubSplit(t, s)
			assert.ErrorIs(t, r2, want, "Split fails with the error of NewStoreStubFailing")
		})

		t.Run("fails Close with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Close-failing")
			s := catalog.NewStoreStubFailing(t, want)
			r0 := storeStubClose(t, s)
			assert.ErrorIs(t, r0, want, "Close fails with the error of NewStoreStubFailing")
		})
	})

	t.Run("StoreStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) storeStubInstance {
			return storeStubInstanceOf(tb, catalog.NewStoreStub(tb, catalog.StoreStubStrict()))
		})
	})

	t.Run("StoreStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Get to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := catalog.NewStoreStub(t)
			s := catalog.NewStoreStub(t, catalog.StoreStubDelegateTo(inner))
			_, _ = storeStubGet(t, s)
			inner.OnGet.AssertCalledOnce(t, "Get calls the implementation")
		})

		t.Run("fails Get with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Get-delegate")
			inner := catalog.NewStoreStub(t)
			inner.OnGet.Faults(want, 1)
			s := catalog.NewStoreStub(t, catalog.StoreStubDelegateTo(inner))
			_, r1 := storeStubGet(t, s)
			assert.ErrorIs(t, r1, want, "Get fails with the error of the implementation")
		})

		t.Run("passes Put to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := catalog.NewStoreStub(t)
			s := catalog.NewStoreStub(t, catalog.StoreStubDelegateTo(inner))
			_ = storeStubPut(t, s)
			inner.OnPut.AssertCalledOnce(t, "Put calls the implementation")
		})

		t.Run("fails Put with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Put-delegate")
			inner := catalog.NewStoreStub(t)
			inner.OnPut.Faults(want, 1)
			s := catalog.NewStoreStub(t, catalog.StoreStubDelegateTo(inner))
			r0 := storeStubPut(t, s)
			assert.ErrorIs(t, r0, want, "Put fails with the error of the implementation")
		})

		t.Run("passes Len to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := catalog.NewStoreStub(t)
			s := catalog.NewStoreStub(t, catalog.StoreStubDelegateTo(inner))
			_ = storeStubLen(t, s)
			inner.OnLen.AssertCalledOnce(t, "Len calls the implementation")
		})

		t.Run("passes Lookup to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := catalog.NewStoreStub(t)
			s := catalog.NewStoreStub(t, catalog.StoreStubDelegateTo(inner))
			_, _ = storeStubLookup(t, s)
			inner.OnLookup.AssertCalledOnce(t, "Lookup calls the implementation")
		})

		t.Run("passes Log to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := catalog.NewStoreStub(t)
			s := catalog.NewStoreStub(t, catalog.StoreStubDelegateTo(inner))
			storeStubLog(t, s)
			inner.OnLog.AssertCalledOnce(t, "Log calls the implementation")
		})

		t.Run("passes Split to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := catalog.NewStoreStub(t)
			s := catalog.NewStoreStub(t, catalog.StoreStubDelegateTo(inner))
			_, _, _ = storeStubSplit(t, s)
			inner.OnSplit.AssertCalledOnce(t, "Split calls the implementation")
		})

		t.Run("fails Split with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Split-delegate")
			inner := catalog.NewStoreStub(t)
			inner.OnSplit.Faults(want, 1)
			s := catalog.NewStoreStub(t, catalog.StoreStubDelegateTo(inner))
			_, _, r2 := storeStubSplit(t, s)
			assert.ErrorIs(t, r2, want, "Split fails with the error of the implementation")
		})

		t.Run("passes Close to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := catalog.NewStoreStub(t)
			s := catalog.NewStoreStub(t, catalog.StoreStubDelegateTo(inner))
			_ = storeStubClose(t, s)
			inner.OnClose.AssertCalledOnce(t, "Close calls the implementation")
		})

		t.Run("fails Close with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Close-delegate")
			inner := catalog.NewStoreStub(t)
			inner.OnClose.Faults(want, 1)
			s := catalog.NewStoreStub(t, catalog.StoreStubDelegateTo(inner))
			r0 := storeStubClose(t, s)
			assert.ErrorIs(t, r0, want, "Close fails with the error of the implementation")
		})
	})

	t.Run("StoreStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) storeStubInstance {
			return storeStubInstanceOf(
				tb,
				catalog.NewStoreStub(tb, catalog.StoreStubWithClock(clk)),
			)
		})
	})

	t.Run("StoreStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) storeStubInstance {
			return storeStubInstanceOf(
				tb,
				catalog.NewStoreStub(tb, catalog.StoreStubWithRandSource(src)),
			)
		})
	})

	t.Run("StoreStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) storeStubInstance {
			return storeStubInstanceOf(tb, catalog.NewStoreStub(tb, catalog.StoreStubBenchMode()))
		})
	})

	t.Run("StoreStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) storeStubInstance {
				return storeStubInstanceOf(
					tb,
					catalog.NewStoreStub(
						tb,
						catalog.StoreStubWithClock(clk),
						catalog.StoreStubWithLatency(d),
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
			s := catalog.NewStoreStub(t, catalog.StoreStubFailing(want))
			_, r1 := storeStubGet(t, s)
			assert.ErrorIs(t, r1, want, "Get fails with the error of StoreStubFailing")
		})

		t.Run("fails Put with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Put-failing")
			s := catalog.NewStoreStub(t, catalog.StoreStubFailing(want))
			r0 := storeStubPut(t, s)
			assert.ErrorIs(t, r0, want, "Put fails with the error of StoreStubFailing")
		})

		t.Run("fails Split with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Split-failing")
			s := catalog.NewStoreStub(t, catalog.StoreStubFailing(want))
			_, _, r2 := storeStubSplit(t, s)
			assert.ErrorIs(t, r2, want, "Split fails with the error of StoreStubFailing")
		})

		t.Run("fails Close with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Close-failing")
			s := catalog.NewStoreStub(t, catalog.StoreStubFailing(want))
			r0 := storeStubClose(t, s)
			assert.ErrorIs(t, r0, want, "Close fails with the error of StoreStubFailing")
		})
	})

	t.Run("WithStoreGet", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Get", func(t *testing.T) {
			t.Parallel()
			called := false
			s := catalog.NewStoreStub(
				t,
				catalog.WithStoreGet(func(_ context.Context, _ string) (catalog.Item, error) {
					called = true
					var z0 catalog.Item
					var z1 error
					return z0, z1
				}),
			)
			_, _ = storeStubGet(t, s)
			assert.True(t, called, "WithStoreGet sets the function of Get")
		})
	})

	t.Run("WithStorePut", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Put", func(t *testing.T) {
			t.Parallel()
			called := false
			s := catalog.NewStoreStub(
				t,
				catalog.WithStorePut(func(_ context.Context, _ catalog.Item) error {
					called = true
					var z0 error
					return z0
				}),
			)
			_ = storeStubPut(t, s)
			assert.True(t, called, "WithStorePut sets the function of Put")
		})
	})

	t.Run("WithStoreLen", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Len", func(t *testing.T) {
			t.Parallel()
			called := false
			s := catalog.NewStoreStub(t, catalog.WithStoreLen(func() int {
				called = true
				var z0 int
				return z0
			}))
			_ = storeStubLen(t, s)
			assert.True(t, called, "WithStoreLen sets the function of Len")
		})
	})

	t.Run("WithStoreLookup", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Lookup", func(t *testing.T) {
			t.Parallel()
			called := false
			s := catalog.NewStoreStub(
				t,
				catalog.WithStoreLookup(func(_ string) (catalog.Item, bool) {
					called = true
					var z0 catalog.Item
					var z1 bool
					return z0, z1
				}),
			)
			_, _ = storeStubLookup(t, s)
			assert.True(t, called, "WithStoreLookup sets the function of Lookup")
		})
	})

	t.Run("WithStoreLog", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Log", func(t *testing.T) {
			t.Parallel()
			called := false
			s := catalog.NewStoreStub(t, catalog.WithStoreLog(func(_ string, _ ...string) {
				called = true
			}))
			storeStubLog(t, s)
			assert.True(t, called, "WithStoreLog sets the function of Log")
		})
	})

	t.Run("WithStoreSplit", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Split", func(t *testing.T) {
			t.Parallel()
			called := false
			s := catalog.NewStoreStub(
				t,
				catalog.WithStoreSplit(func(_ string) (string, string, error) {
					called = true
					var z0 string
					var z1 string
					var z2 error
					return z0, z1, z2
				}),
			)
			_, _, _ = storeStubSplit(t, s)
			assert.True(t, called, "WithStoreSplit sets the function of Split")
		})
	})

	t.Run("WithStoreClose", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Close", func(t *testing.T) {
			t.Parallel()
			called := false
			s := catalog.NewStoreStub(t, catalog.WithStoreClose(func() error {
				called = true
				var z0 error
				return z0
			}))
			_ = storeStubClose(t, s)
			assert.True(t, called, "WithStoreClose sets the function of Close")
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
			s := catalog.NewStoreStub(t)
			want0, want1 := storeStubGetPins()
			s.OnGet.Returns(want0, want1)
			got0, got1 := storeStubGet(t, s)
			expect.Equal(t, got0, want0, "Get returns the Result that Returns sets")
			expect.ErrorIs(t, got1, want1, "Get returns the error that Returns sets")
		})

		t.Run("configures Get for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := catalog.NewStoreStub(r)
			s.OnGet.Strict()
			var want0 catalog.Item
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
			s := catalog.NewStoreStub(t)
			var got0 context.Context
			var got1 string
			s.OnGet.Func(func(p0 context.Context, p1 string) (catalog.Item, error) {
				got0, got1 = p0, p1
				var z0 catalog.Item
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

// TestStorePutStub checks the methods of StorePutStub, which configures
// Put.
func TestStorePutStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Put", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			want0 := storeStubPutPins()
			s.OnPut.Returns(want0)
			got0 := storeStubPut(t, s)
			assert.ErrorIs(t, got0, want0, "Put returns the error that Returns sets")
		})

		t.Run("configures Put for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := catalog.NewStoreStub(r)
			s.OnPut.Strict()
			var want0 error
			s.OnPut.Returns(want0)
			_ = storeStubPut(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Put")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Put to the function", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			var got0 context.Context
			var got1 catalog.Item
			s.OnPut.Func(func(p0 context.Context, p1 catalog.Item) error {
				got0, got1 = p0, p1
				var z0 error
				return z0
			})
			_ = storeStubPut(t, s)
			a0, a1 := storeStubPutArgs(t)
			expect.Equal(t, got0, a0, "the function receives A0", expect.ByIdentity())
			expect.Equal(t, got1, a1, "the function receives A1")
		})
	})
}

// TestStoreLenStub checks the methods of StoreLenStub, which configures
// Len.
func TestStoreLenStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Len", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			want0 := storeStubLenPins()
			s.OnLen.Returns(want0)
			got0 := storeStubLen(t, s)
			assert.Equal(t, got0, want0, "Len returns the Result that Returns sets")
		})

		t.Run("configures Len for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := catalog.NewStoreStub(r)
			s.OnLen.Strict()
			var want0 int
			s.OnLen.Returns(want0)
			_ = storeStubLen(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Len")
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
			s := catalog.NewStoreStub(t)
			want0, want1 := storeStubLookupPins()
			s.OnLookup.Returns(want0, want1)
			got0, got1 := storeStubLookup(t, s)
			expect.Equal(t, got0, want0, "Lookup returns the Item that Returns sets")
			expect.Equal(t, got1, want1, "Lookup returns the Ok that Returns sets")
		})

		t.Run("configures Lookup for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := catalog.NewStoreStub(r)
			s.OnLookup.Strict()
			var want0 catalog.Item
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
			s := catalog.NewStoreStub(t)
			var got0 string
			s.OnLookup.Func(func(p0 string) (catalog.Item, bool) {
				got0 = p0
				var z0 catalog.Item
				var z1 bool
				return z0, z1
			})
			_, _ = storeStubLookup(t, s)
			a0 := storeStubLookupArgs(t)
			assert.Equal(t, got0, a0, "the function receives Key")
		})
	})
}

// TestStoreLogStub checks the methods of StoreLogStub, which configures
// Log.
func TestStoreLogStub(t *testing.T) {
	t.Parallel()

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Log to the function", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			var got0 string
			var got1 []string
			s.OnLog.Func(func(p0 string, p1 ...string) {
				got0, got1 = p0, p1
			})
			storeStubLog(t, s)
			a0, a1 := storeStubLogArgs(t)
			expect.Equal(t, got0, a0, "the function receives Format")
			expect.Equal(t, got1, a1, "the function receives Args")
		})
	})
}

// TestStoreSplitStub checks the methods of StoreSplitStub, which
// configures Split.
func TestStoreSplitStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Split", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			want0, want1, want2 := storeStubSplitPins()
			s.OnSplit.Returns(want0, want1, want2)
			got0, got1, got2 := storeStubSplit(t, s)
			expect.Equal(t, got0, want0, "Split returns the Result0 that Returns sets")
			expect.Equal(t, got1, want1, "Split returns the Result1 that Returns sets")
			expect.ErrorIs(t, got2, want2, "Split returns the error that Returns sets")
		})

		t.Run("configures Split for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := catalog.NewStoreStub(r)
			s.OnSplit.Strict()
			var want0 string
			var want1 string
			var want2 error
			s.OnSplit.Returns(want0, want1, want2)
			_, _, _ = storeStubSplit(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Split")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Split to the function", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			var got0 string
			s.OnSplit.Func(func(p0 string) (string, string, error) {
				got0 = p0
				var z0 string
				var z1 string
				var z2 error
				return z0, z1, z2
			})
			_, _, _ = storeStubSplit(t, s)
			a0 := storeStubSplitArgs(t)
			assert.Equal(t, got0, a0, "the function receives S")
		})
	})
}

// TestStoreCloseStub checks the methods of StoreCloseStub, which
// configures Close.
func TestStoreCloseStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Close", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewStoreStub(t)
			want0 := storeStubClosePins()
			s.OnClose.Returns(want0)
			got0 := storeStubClose(t, s)
			assert.ErrorIs(t, got0, want0, "Close returns the error that Returns sets")
		})

		t.Run("configures Close for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := catalog.NewStoreStub(r)
			s.OnClose.Strict()
			var want0 error
			s.OnClose.Returns(want0)
			_ = storeStubClose(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Close")
		})
	})
}

// clockStubNow calls Now of s with the samples of its arguments.
func clockStubNow(tb assert.TB, s *catalog.ClockStub) time.Time {
	tb.Helper()
	return s.Now()
}

// clockStubNowPins returns the samples of the results of Now.
func clockStubNowPins() time.Time {
	var want0 time.Time
	want0 = time.Unix(1700000000, 0)
	return want0
}

// clockStubZone calls Zone of s with the samples of its arguments.
func clockStubZone(tb assert.TB, s *catalog.ClockStub) string {
	tb.Helper()
	return s.Zone()
}

// clockStubZonePins returns the samples of the results of Zone.
func clockStubZonePins() string {
	var want0 string
	want0 = "Zone-a"
	return want0
}

// clockStubEveryArgs returns the samples of the arguments of Every.
func clockStubEveryArgs(tb assert.TB) (time.Duration, func()) {
	tb.Helper()
	var a0 time.Duration
	var a1 func()
	a0 = time.Duration(1)
	a1 = func() {}
	return a0, a1
}

// clockStubEvery calls Every of s with the samples of its arguments.
func clockStubEvery(tb assert.TB, s *catalog.ClockStub) func() {
	tb.Helper()
	a0, a1 := clockStubEveryArgs(tb)
	return s.Every(a0, a1)
}

// clockStubEveryPins returns the samples of the results of Every.
func clockStubEveryPins() func() {
	var want0 func()
	want0 = func() {}
	return want0
}

// clockStubNowSubject returns the subject of [stub.Behaviour] for Now. The
// subject calls a new double that reports to tb.
func clockStubNowSubject(tb assert.TB) stub.Subject[catalog.ClockNowCall, catalog.ClockNowReturn] {
	tb.Helper()
	s := catalog.NewClockStub(tb)
	return stub.Subject[catalog.ClockNowCall, catalog.ClockNowReturn]{
		Stub: s.OnNow.MethodStub,
		Call: func() {
			_ = clockStubNow(tb, s)
		},
		Result: func() catalog.ClockNowReturn {
			got0 := clockStubNow(tb, s)
			return catalog.ClockNowReturn{Result: got0}
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

// clockStubZoneSubject returns the subject of [stub.Behaviour] for Zone.
// The subject calls a new double that reports to tb.
func clockStubZoneSubject(
	tb assert.TB,
) stub.Subject[catalog.ClockZoneCall, catalog.ClockZoneReturn] {
	tb.Helper()
	s := catalog.NewClockStub(tb)
	return stub.Subject[catalog.ClockZoneCall, catalog.ClockZoneReturn]{
		Stub: s.OnZone.MethodStub,
		Call: func() {
			_ = clockStubZone(tb, s)
		},
		Result: func() catalog.ClockZoneReturn {
			got0 := clockStubZone(tb, s)
			return catalog.ClockZoneReturn{Result: got0}
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

// clockStubEverySubject returns the subject of [stub.Behaviour] for Every.
// The subject calls a new double that reports to tb.
func clockStubEverySubject(
	tb assert.TB,
) stub.Subject[catalog.ClockEveryCall, catalog.ClockEveryReturn] {
	tb.Helper()
	s := catalog.NewClockStub(tb)
	return stub.Subject[catalog.ClockEveryCall, catalog.ClockEveryReturn]{
		Stub: s.OnEvery.MethodStub,
		Call: func() {
			_ = clockStubEvery(tb, s)
		},
		Result: func() catalog.ClockEveryReturn {
			got0 := clockStubEvery(tb, s)
			return catalog.ClockEveryReturn{Stop: got0}
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

// clockStubInstance is the instance of a catalog.ClockStub that the checks
// of the constructors and the options drive.
type clockStubInstance = stub.Instance[catalog.ClockNowCall]

// clockStubInstanceOf returns the instance of s. The options apply to the
// whole double, so the checks call Now alone.
func clockStubInstanceOf(tb assert.TB, s *catalog.ClockStub) clockStubInstance {
	tb.Helper()
	return clockStubInstance{
		Stub: s.OnNow.MethodStub,
		Call: func() {
			_ = clockStubNow(tb, s)
		},
		Reset: s.ResetCalls,
	}
}

// TestClockStub checks the methods, the constructors and the options of
// the double ClockStub.
func TestClockStub(t *testing.T) {
	t.Parallel()

	t.Run("Now", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Now", clockStubNowSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewClockStub(t)
			var seen []catalog.ClockNowCall
			s.OnNow.OnRecord(func(c catalog.ClockNowCall) { seen = append(seen, c) })
			_ = clockStubNow(t, s)
			_ = clockStubNow(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Now")
		})
	})

	t.Run("Zone", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Zone", clockStubZoneSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewClockStub(t)
			var seen []catalog.ClockZoneCall
			s.OnZone.OnRecord(func(c catalog.ClockZoneCall) { seen = append(seen, c) })
			_ = clockStubZone(t, s)
			_ = clockStubZone(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Zone")
		})
	})

	t.Run("Every", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Every", clockStubEverySubject)

		t.Run("records the arguments of the call", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewClockStub(t)
			_ = clockStubEvery(t, s)
			got := s.OnEvery.AssertCalledOnce(t, "Every records the call")
			a0, a1 := clockStubEveryArgs(t)
			expect.Equal(t, got.D, a0, "the call records D")
			expect.Equal(t, got.Fn, a1, "the call records Fn")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewClockStub(t)
			var seen []catalog.ClockEveryCall
			s.OnEvery.OnRecord(func(c catalog.ClockEveryCall) { seen = append(seen, c) })
			_ = clockStubEvery(t, s)
			_ = clockStubEvery(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Every")
		})
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()
		stub.ResetBehaviour(t, func(tb assert.TB) clockStubInstance {
			return clockStubInstanceOf(tb, catalog.NewClockStub(tb))
		})

		t.Run("clears the recorded calls of Now", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewClockStub(t)
			_ = clockStubNow(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnNow.CallCount(), 0, "ResetCalls clears the calls of Now")
		})

		t.Run("keeps the results of Returns of Now", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewClockStub(t)
			want0 := clockStubNowPins()
			s.OnNow.Returns(want0)
			s.ResetCalls()
			got0 := clockStubNow(t, s)
			assert.Equal(t, got0, want0, "ResetCalls keeps the Result of Now")
		})

		t.Run("clears the recorded calls of Zone", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewClockStub(t)
			_ = clockStubZone(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnZone.CallCount(), 0, "ResetCalls clears the calls of Zone")
		})

		t.Run("keeps the results of Returns of Zone", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewClockStub(t)
			want0 := clockStubZonePins()
			s.OnZone.Returns(want0)
			s.ResetCalls()
			got0 := clockStubZone(t, s)
			assert.Equal(t, got0, want0, "ResetCalls keeps the Result of Zone")
		})

		t.Run("clears the recorded calls of Every", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewClockStub(t)
			_ = clockStubEvery(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnEvery.CallCount(), 0, "ResetCalls clears the calls of Every")
		})

		t.Run("keeps the results of Returns of Every", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewClockStub(t)
			want0 := clockStubEveryPins()
			s.OnEvery.Returns(want0)
			s.ResetCalls()
			got0 := clockStubEvery(t, s)
			assert.Equal(t, got0, want0, "ResetCalls keeps the Stop of Every")
		})
	})

	t.Run("NewClockStubDelayed", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) clockStubInstance {
				return clockStubInstanceOf(tb, catalog.NewClockStubDelayed(tb, d, clk))
			},
		)
	})

	t.Run("ClockStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) clockStubInstance {
			return clockStubInstanceOf(tb, catalog.NewClockStub(tb, catalog.ClockStubStrict()))
		})
	})

	t.Run("ClockStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Now to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := catalog.NewClockStub(t)
			s := catalog.NewClockStub(t, catalog.ClockStubDelegateTo(inner))
			_ = clockStubNow(t, s)
			inner.OnNow.AssertCalledOnce(t, "Now calls the implementation")
		})

		t.Run("passes Zone to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := catalog.NewClockStub(t)
			s := catalog.NewClockStub(t, catalog.ClockStubDelegateTo(inner))
			_ = clockStubZone(t, s)
			inner.OnZone.AssertCalledOnce(t, "Zone calls the implementation")
		})

		t.Run("passes Every to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := catalog.NewClockStub(t)
			s := catalog.NewClockStub(t, catalog.ClockStubDelegateTo(inner))
			_ = clockStubEvery(t, s)
			inner.OnEvery.AssertCalledOnce(t, "Every calls the implementation")
		})
	})

	t.Run("ClockStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) clockStubInstance {
			return clockStubInstanceOf(
				tb,
				catalog.NewClockStub(tb, catalog.ClockStubWithClock(clk)),
			)
		})
	})

	t.Run("ClockStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) clockStubInstance {
			return clockStubInstanceOf(
				tb,
				catalog.NewClockStub(tb, catalog.ClockStubWithRandSource(src)),
			)
		})
	})

	t.Run("ClockStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) clockStubInstance {
			return clockStubInstanceOf(tb, catalog.NewClockStub(tb, catalog.ClockStubBenchMode()))
		})
	})

	t.Run("ClockStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) clockStubInstance {
				return clockStubInstanceOf(
					tb,
					catalog.NewClockStub(
						tb,
						catalog.ClockStubWithClock(clk),
						catalog.ClockStubWithLatency(d),
					),
				)
			},
		)
	})

	t.Run("WithClockNow", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Now", func(t *testing.T) {
			t.Parallel()
			called := false
			s := catalog.NewClockStub(t, catalog.WithClockNow(func() time.Time {
				called = true
				var z0 time.Time
				return z0
			}))
			_ = clockStubNow(t, s)
			assert.True(t, called, "WithClockNow sets the function of Now")
		})
	})

	t.Run("WithClockZone", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Zone", func(t *testing.T) {
			t.Parallel()
			called := false
			s := catalog.NewClockStub(t, catalog.WithClockZone(func() string {
				called = true
				var z0 string
				return z0
			}))
			_ = clockStubZone(t, s)
			assert.True(t, called, "WithClockZone sets the function of Zone")
		})
	})

	t.Run("WithClockEvery", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Every", func(t *testing.T) {
			t.Parallel()
			called := false
			s := catalog.NewClockStub(
				t,
				catalog.WithClockEvery(func(_ time.Duration, _ func()) func() {
					called = true
					var z0 func()
					return z0
				}),
			)
			_ = clockStubEvery(t, s)
			assert.True(t, called, "WithClockEvery sets the function of Every")
		})
	})
}

// TestClockNowStub checks the methods of ClockNowStub, which configures
// Now.
func TestClockNowStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Now", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewClockStub(t)
			want0 := clockStubNowPins()
			s.OnNow.Returns(want0)
			got0 := clockStubNow(t, s)
			assert.Equal(t, got0, want0, "Now returns the Result that Returns sets")
		})

		t.Run("configures Now for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := catalog.NewClockStub(r)
			s.OnNow.Strict()
			var want0 time.Time
			s.OnNow.Returns(want0)
			_ = clockStubNow(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Now")
		})
	})
}

// TestClockZoneStub checks the methods of ClockZoneStub, which configures
// Zone.
func TestClockZoneStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Zone", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewClockStub(t)
			want0 := clockStubZonePins()
			s.OnZone.Returns(want0)
			got0 := clockStubZone(t, s)
			assert.Equal(t, got0, want0, "Zone returns the Result that Returns sets")
		})

		t.Run("configures Zone for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := catalog.NewClockStub(r)
			s.OnZone.Strict()
			var want0 string
			s.OnZone.Returns(want0)
			_ = clockStubZone(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Zone")
		})
	})
}

// TestClockEveryStub checks the methods of ClockEveryStub, which
// configures Every.
func TestClockEveryStub(t *testing.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Every", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewClockStub(t)
			want0 := clockStubEveryPins()
			s.OnEvery.Returns(want0)
			got0 := clockStubEvery(t, s)
			assert.Equal(t, got0, want0, "Every returns the Stop that Returns sets")
		})

		t.Run("configures Every for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := catalog.NewClockStub(r)
			s.OnEvery.Strict()
			var want0 func()
			s.OnEvery.Returns(want0)
			_ = clockStubEvery(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of Every")
		})
	})

	t.Run("Func", func(t *testing.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Every to the function", func(t *testing.T) {
			t.Parallel()
			s := catalog.NewClockStub(t)
			var got0 time.Duration
			var got1 func()
			s.OnEvery.Func(func(p0 time.Duration, p1 func()) func() {
				got0, got1 = p0, p1
				var z0 func()
				return z0
			})
			_ = clockStubEvery(t, s)
			a0, a1 := clockStubEveryArgs(t)
			expect.Equal(t, got0, a0, "the function receives D")
			expect.Equal(t, got1, a1, "the function receives Fn")
		})
	})
}
