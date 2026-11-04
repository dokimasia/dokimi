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

// storeStubGetSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Get. The subject calls a new double that
// reports to tb.
func storeStubGetSubject(tb assert.TB) stub.Subject[catalog.StoreGetCall, catalog.StoreGetReturn] {
	tb.Helper()
	s := catalog.NewStoreStub(tb)
	return stub.Subject[catalog.StoreGetCall, catalog.StoreGetReturn]{
		Stub:  s.OnGet.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = storeStubGet(tb, s)
		},
		Arguments: func(tb assert.TB, c catalog.StoreGetCall) {
			tb.Helper()
			a0, a1 := storeStubGetArgs(tb)
			expect.Equal(tb, c.Ctx, a0, "the call records Ctx", expect.ByIdentity())
			expect.Equal(tb, c.ID, a1, "the call records ID")
		},
		Result: func() catalog.StoreGetReturn {
			got0, got1 := storeStubGet(tb, s)
			return catalog.StoreGetReturn{Result: got0, Err: got1}
		},
		Returns: func() catalog.StoreGetReturn {
			want0, want1 := storeStubGetPins()
			s.OnGet.Returns(want0, want1)
			return catalog.StoreGetReturn{Result: want0, Err: want1}
		},
		Override: func(seen func(catalog.StoreGetCall)) {
			s.OnGet.Func(func(p0 context.Context, p1 string) (catalog.Item, error) {
				seen(catalog.StoreGetCall{Ctx: p0, ID: p1})
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

// storeStubPutSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Put. The subject calls a new double that
// reports to tb.
func storeStubPutSubject(tb assert.TB) stub.Subject[catalog.StorePutCall, catalog.StorePutReturn] {
	tb.Helper()
	s := catalog.NewStoreStub(tb)
	return stub.Subject[catalog.StorePutCall, catalog.StorePutReturn]{
		Stub:  s.OnPut.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = storeStubPut(tb, s)
		},
		Arguments: func(tb assert.TB, c catalog.StorePutCall) {
			tb.Helper()
			a0, a1 := storeStubPutArgs(tb)
			expect.Equal(tb, c.A0, a0, "the call records A0", expect.ByIdentity())
			expect.Equal(tb, c.A1, a1, "the call records A1")
		},
		Result: func() catalog.StorePutReturn {
			got0 := storeStubPut(tb, s)
			return catalog.StorePutReturn{Err: got0}
		},
		Returns: func() catalog.StorePutReturn {
			want0 := storeStubPutPins()
			s.OnPut.Returns(want0)
			return catalog.StorePutReturn{Err: want0}
		},
		Override: func(seen func(catalog.StorePutCall)) {
			s.OnPut.Func(func(p0 context.Context, p1 catalog.Item) error {
				seen(catalog.StorePutCall{A0: p0, A1: p1})
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

// storeStubLenSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Len. The subject calls a new double that
// reports to tb.
func storeStubLenSubject(tb assert.TB) stub.Subject[catalog.StoreLenCall, catalog.StoreLenReturn] {
	tb.Helper()
	s := catalog.NewStoreStub(tb)
	return stub.Subject[catalog.StoreLenCall, catalog.StoreLenReturn]{
		Stub:  s.OnLen.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = storeStubLen(tb, s)
		},
		Result: func() catalog.StoreLenReturn {
			got0 := storeStubLen(tb, s)
			return catalog.StoreLenReturn{Result: got0}
		},
		Returns: func() catalog.StoreLenReturn {
			want0 := storeStubLenPins()
			s.OnLen.Returns(want0)
			return catalog.StoreLenReturn{Result: want0}
		},
		Override: func(seen func(catalog.StoreLenCall)) {
			s.OnLen.Func(func() int {
				seen(catalog.StoreLenCall{})
				var z0 int
				return z0
			})
		},
	}
}

// storeStubLookupSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Lookup. The subject calls a new double that
// reports to tb.
func storeStubLookupSubject(
	tb assert.TB,
) stub.Subject[catalog.StoreLookupCall, catalog.StoreLookupReturn] {
	tb.Helper()
	s := catalog.NewStoreStub(tb)
	return stub.Subject[catalog.StoreLookupCall, catalog.StoreLookupReturn]{
		Stub:  s.OnLookup.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = storeStubLookup(tb, s)
		},
		Arguments: func(tb assert.TB, c catalog.StoreLookupCall) {
			tb.Helper()
			a0 := storeStubLookupArgs(tb)
			assert.Equal(tb, c.Key, a0, "the call records Key")
		},
		Result: func() catalog.StoreLookupReturn {
			got0, got1 := storeStubLookup(tb, s)
			return catalog.StoreLookupReturn{Item: got0, Ok: got1}
		},
		Returns: func() catalog.StoreLookupReturn {
			want0, want1 := storeStubLookupPins()
			s.OnLookup.Returns(want0, want1)
			return catalog.StoreLookupReturn{Item: want0, Ok: want1}
		},
		Override: func(seen func(catalog.StoreLookupCall)) {
			s.OnLookup.Func(func(p0 string) (catalog.Item, bool) {
				seen(catalog.StoreLookupCall{Key: p0})
				var z0 catalog.Item
				var z1 bool
				return z0, z1
			})
		},
	}
}

// storeStubLogSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Log. The subject calls a new double that
// reports to tb.
func storeStubLogSubject(tb assert.TB) stub.Subject[catalog.StoreLogCall, struct{}] {
	tb.Helper()
	s := catalog.NewStoreStub(tb)
	return stub.Subject[catalog.StoreLogCall, struct{}]{
		Stub:  s.OnLog.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			storeStubLog(tb, s)
		},
		Arguments: func(tb assert.TB, c catalog.StoreLogCall) {
			tb.Helper()
			a0, a1 := storeStubLogArgs(tb)
			expect.Equal(tb, c.Format, a0, "the call records Format")
			expect.Equal(tb, c.Args, a1, "the call records Args")
		},
		Override: func(seen func(catalog.StoreLogCall)) {
			s.OnLog.Func(func(p0 string, p1 ...string) {
				seen(catalog.StoreLogCall{Format: p0, Args: p1})
			})
		},
	}
}

// storeStubSplitSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Split. The subject calls a new double that
// reports to tb.
func storeStubSplitSubject(
	tb assert.TB,
) stub.Subject[catalog.StoreSplitCall, catalog.StoreSplitReturn] {
	tb.Helper()
	s := catalog.NewStoreStub(tb)
	return stub.Subject[catalog.StoreSplitCall, catalog.StoreSplitReturn]{
		Stub:  s.OnSplit.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _, _ = storeStubSplit(tb, s)
		},
		Arguments: func(tb assert.TB, c catalog.StoreSplitCall) {
			tb.Helper()
			a0 := storeStubSplitArgs(tb)
			assert.Equal(tb, c.S, a0, "the call records S")
		},
		Result: func() catalog.StoreSplitReturn {
			got0, got1, got2 := storeStubSplit(tb, s)
			return catalog.StoreSplitReturn{Result0: got0, Result1: got1, Err: got2}
		},
		Returns: func() catalog.StoreSplitReturn {
			want0, want1, want2 := storeStubSplitPins()
			s.OnSplit.Returns(want0, want1, want2)
			return catalog.StoreSplitReturn{Result0: want0, Result1: want1, Err: want2}
		},
		Override: func(seen func(catalog.StoreSplitCall)) {
			s.OnSplit.Func(func(p0 string) (string, string, error) {
				seen(catalog.StoreSplitCall{S: p0})
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

// storeStubCloseSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Close. The subject calls a new double that
// reports to tb.
func storeStubCloseSubject(
	tb assert.TB,
) stub.Subject[catalog.StoreCloseCall, catalog.StoreCloseReturn] {
	tb.Helper()
	s := catalog.NewStoreStub(tb)
	return stub.Subject[catalog.StoreCloseCall, catalog.StoreCloseReturn]{
		Stub:  s.OnClose.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = storeStubClose(tb, s)
		},
		Result: func() catalog.StoreCloseReturn {
			got0 := storeStubClose(tb, s)
			return catalog.StoreCloseReturn{Err: got0}
		},
		Returns: func() catalog.StoreCloseReturn {
			want0 := storeStubClosePins()
			s.OnClose.Returns(want0)
			return catalog.StoreCloseReturn{Err: want0}
		},
		Override: func(seen func(catalog.StoreCloseCall)) {
			s.OnClose.Func(func() error {
				seen(catalog.StoreCloseCall{})
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

	t.Run("Put", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Put", storeStubPutSubject)
	})

	t.Run("Len", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Len", storeStubLenSubject)
	})

	t.Run("Lookup", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Lookup", storeStubLookupSubject)
	})

	t.Run("Log", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Log", storeStubLogSubject)
	})

	t.Run("Split", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Split", storeStubSplitSubject)
	})

	t.Run("Close", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Close", storeStubCloseSubject)
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()

		t.Run("Get", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Get", storeStubGetSubject)
		})

		t.Run("Put", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Put", storeStubPutSubject)
		})

		t.Run("Len", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Len", storeStubLenSubject)
		})

		t.Run("Lookup", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Lookup", storeStubLookupSubject)
		})

		t.Run("Log", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Log", storeStubLogSubject)
		})

		t.Run("Split", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Split", storeStubSplitSubject)
		})

		t.Run("Close", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Close", storeStubCloseSubject)
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

// clockStubNowSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Now. The subject calls a new double that
// reports to tb.
func clockStubNowSubject(tb assert.TB) stub.Subject[catalog.ClockNowCall, catalog.ClockNowReturn] {
	tb.Helper()
	s := catalog.NewClockStub(tb)
	return stub.Subject[catalog.ClockNowCall, catalog.ClockNowReturn]{
		Stub:  s.OnNow.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = clockStubNow(tb, s)
		},
		Result: func() catalog.ClockNowReturn {
			got0 := clockStubNow(tb, s)
			return catalog.ClockNowReturn{Result: got0}
		},
		Returns: func() catalog.ClockNowReturn {
			want0 := clockStubNowPins()
			s.OnNow.Returns(want0)
			return catalog.ClockNowReturn{Result: want0}
		},
		Override: func(seen func(catalog.ClockNowCall)) {
			s.OnNow.Func(func() time.Time {
				seen(catalog.ClockNowCall{})
				var z0 time.Time
				return z0
			})
		},
	}
}

// clockStubZoneSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Zone. The subject calls a new double that
// reports to tb.
func clockStubZoneSubject(
	tb assert.TB,
) stub.Subject[catalog.ClockZoneCall, catalog.ClockZoneReturn] {
	tb.Helper()
	s := catalog.NewClockStub(tb)
	return stub.Subject[catalog.ClockZoneCall, catalog.ClockZoneReturn]{
		Stub:  s.OnZone.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = clockStubZone(tb, s)
		},
		Result: func() catalog.ClockZoneReturn {
			got0 := clockStubZone(tb, s)
			return catalog.ClockZoneReturn{Result: got0}
		},
		Returns: func() catalog.ClockZoneReturn {
			want0 := clockStubZonePins()
			s.OnZone.Returns(want0)
			return catalog.ClockZoneReturn{Result: want0}
		},
		Override: func(seen func(catalog.ClockZoneCall)) {
			s.OnZone.Func(func() string {
				seen(catalog.ClockZoneCall{})
				var z0 string
				return z0
			})
		},
	}
}

// clockStubEverySubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Every. The subject calls a new double that
// reports to tb.
func clockStubEverySubject(
	tb assert.TB,
) stub.Subject[catalog.ClockEveryCall, catalog.ClockEveryReturn] {
	tb.Helper()
	s := catalog.NewClockStub(tb)
	return stub.Subject[catalog.ClockEveryCall, catalog.ClockEveryReturn]{
		Stub:  s.OnEvery.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = clockStubEvery(tb, s)
		},
		Arguments: func(tb assert.TB, c catalog.ClockEveryCall) {
			tb.Helper()
			a0, a1 := clockStubEveryArgs(tb)
			expect.Equal(tb, c.D, a0, "the call records D")
			expect.Equal(tb, c.Fn, a1, "the call records Fn")
		},
		Result: func() catalog.ClockEveryReturn {
			got0 := clockStubEvery(tb, s)
			return catalog.ClockEveryReturn{Stop: got0}
		},
		Returns: func() catalog.ClockEveryReturn {
			want0 := clockStubEveryPins()
			s.OnEvery.Returns(want0)
			return catalog.ClockEveryReturn{Stop: want0}
		},
		Override: func(seen func(catalog.ClockEveryCall)) {
			s.OnEvery.Func(func(p0 time.Duration, p1 func()) func() {
				seen(catalog.ClockEveryCall{D: p0, Fn: p1})
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
	}
}

// TestClockStub checks the methods, the constructors and the options of
// the double ClockStub.
func TestClockStub(t *testing.T) {
	t.Parallel()

	t.Run("Now", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Now", clockStubNowSubject)
	})

	t.Run("Zone", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Zone", clockStubZoneSubject)
	})

	t.Run("Every", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Every", clockStubEverySubject)
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()

		t.Run("Now", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Now", clockStubNowSubject)
		})

		t.Run("Zone", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Zone", clockStubZoneSubject)
		})

		t.Run("Every", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Every", clockStubEverySubject)
		})
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
