package parity

import (
	"context"
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/dokimi-go/stub"
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
			return StoreGetReturn{Item: got0, Err: got1}
		},
		Override: func(mark func()) {
			s.OnGet.Func(func(_ context.Context, _ string) (string, error) {
				mark()
				var z0 string
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			var a0 context.Context
			var a1 string
			a1 = "Get-a"
			_, err := s.Get(a0, a1)
			return err
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
		var want0 string
		var want1 error
		want0 = "Get-a"
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
		s := NewStoreStub(t, WithStoreGet(func(_ context.Context, _ string) (string, error) {
			called = true
			var z0 string
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
		_, err := s.Get(a0, a1)
		assert.ErrorIs(t, err, want, "Get must answer the armed error")
	})

	t.Run("answers ErrNotFound once FaultNotFound arms it", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		s.OnGet.FaultNotFound()
		var a0 context.Context
		var a1 string
		a1 = "Get-a"
		_, err := s.Get(a0, a1)
		assert.ErrorIs(t, err, ErrNotFound, "FaultNotFound must answer with ErrNotFound")
		_, err = s.Get(a0, a1)
		assert.NoError(t, err, "the fault is one-shot and must clear")
	})

	t.Run("answers ErrGone once FaultGone arms it", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		s.OnGet.FaultGone()
		var a0 context.Context
		var a1 string
		a1 = "Get-a"
		_, err := s.Get(a0, a1)
		assert.ErrorIs(t, err, ErrGone, "FaultGone must answer with ErrGone")
		_, err = s.Get(a0, a1)
		assert.NoError(t, err, "the fault is one-shot and must clear")
	})

	t.Run("faults only the targeted partition", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var a0 context.Context
		var a1 string
		a1 = "Get-a"
		key := a1
		want := errors.New("Get-partition")
		s.OnGet.FaultForPartition(key, want, 1)
		_, err := s.Get(a0, a1)
		assert.ErrorIs(t, err, want, "a call in the targeted partition must fail")
	})

	t.Run("spares the partition FaultForOtherPartitions names", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var a0 context.Context
		var a1 string
		a1 = "Get-a"
		key := a1
		s.OnGet.FaultForOtherPartitions(key, errors.New("Get-partition"), 1)
		_, err := s.Get(a0, a1)
		assert.NoError(t, err, "a call in the named partition must be spared")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 string
		var want1 error
		want0 = "Get-a"
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
			_, _ = s.List()
		},
		Result: func() StoreListReturn {
			got0, got1 := s.List()
			return StoreListReturn{Result: got0, Err: got1}
		},
		Override: func(mark func()) {
			s.OnList.Func(func() ([]string, error) {
				mark()
				var z0 []string
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			_, r1 := s.List()
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
		var want0 []string
		var want1 error
		want0 = []string{"List-a"}
		s.OnList.Returns(want0, want1)
		got0, _ := s.List()
		assert.Equal(t, got0, want0, "List must answer with what Returns pinned")
	})

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		_, _ = s.List()
		s.OnList.AssertCalledOnce(t, "List must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var seen []StoreListCall
		s.OnList.OnRecord(func(c StoreListCall) { seen = append(seen, c) })
		_, _ = s.List()
		_, _ = s.List()
		assert.Length(t, seen, 2, "OnRecord must fire once per List call")
	})

	t.Run("wires WithStoreList at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewStoreStub(t, WithStoreList(func() ([]string, error) {
			called = true
			var z0 []string
			var z1 error
			return z0, z1
		}))
		_, _ = s.List()
		assert.True(t, called, "WithStoreList must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("List-failing")
		s := NewStoreStubFailing(t, want)
		_, r1 := s.List()
		assert.ErrorIs(t, r1, want, "List must answer the armed error")
	})

	t.Run("fails until attempt 3 then succeeds", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		want := errors.New("List-retry")
		s.OnList.RetrySchedule(want)
		for attempt := 1; attempt < 3; attempt++ {
			_, r1 := s.List()
			assert.ErrorIs(t, r1, want, "a call before attempt 3 must fail")
		}
		_, r1 := s.List()
		assert.NoError(t, r1, "attempt 3 must succeed")
	})

	t.Run("keeps the Returns configuration across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var want0 []string
		var want1 error
		want0 = []string{"List-a"}
		s.OnList.Returns(want0, want1)
		_, _ = s.List()
		s.ResetCalls()
		got0, _ := s.List()
		assert.Equal(t, got0, want0, "a reset must keep what Returns pinned")
	})
}

// storeStubCloseSubject binds Close into the shape [stub.Behaviour]
// drives: how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func storeStubCloseSubject(tb assert.TB) stub.Subject[StoreCloseCall, struct{}] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreCloseCall, struct{}]{
		Stub: s.OnClose.MethodStub,
		Call: func() {
			s.Close()
		},
		Override: func(mark func()) {
			s.OnClose.Func(func() {
				mark()
			})
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
		s.Close()
		s.OnClose.AssertCalledOnce(t, "Close must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewStoreStub(t)
		var seen []StoreCloseCall
		s.OnClose.OnRecord(func(c StoreCloseCall) { seen = append(seen, c) })
		s.Close()
		s.Close()
		assert.Length(t, seen, 2, "OnRecord must fire once per Close call")
	})

	t.Run("wires WithStoreClose at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewStoreStub(t, WithStoreClose(func() {
			called = true
		}))
		s.Close()
		assert.True(t, called, "WithStoreClose must install the override")
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
		_, err := s.Get(a0, a1)
		assert.ErrorIs(t, err, want, "Get must surface the wrapped answer")
	})

	t.Run("forwards List to the wrapped implementation", func(t *testing.T) {
		_, _ = s.List()
		inner.OnList.AssertCalledOnce(t, "List must reach the wrapped implementation")
	})

	t.Run("surfaces what List answered", func(t *testing.T) {
		want := errors.New("List-delegate")
		inner.OnList.FaultsFor(time.Hour, want)
		_, r1 := s.List()
		assert.ErrorIs(t, r1, want, "List must surface the wrapped answer")
	})

	t.Run("forwards Close to the wrapped implementation", func(t *testing.T) {
		s.Close()
		inner.OnClose.AssertCalledOnce(t, "Close must reach the wrapped implementation")
	})
}
