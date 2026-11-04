package parity_test

import (
	"context"
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"example.com/shop/parity"
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
func storeStubGet(tb assert.TB, s *parity.StoreStub) (string, error) {
	tb.Helper()
	a0, a1 := storeStubGetArgs(tb)
	return s.Get(a0, a1)
}

// storeStubGetPins returns the samples of the results of Get, with a new
// error as its error result.
func storeStubGetPins() (string, error) {
	var want0 string
	want1 := errors.New("Get-returns")
	want0 = "Get-a"
	return want0, want1
}

// storeStubList calls List of s with the samples of its arguments.
func storeStubList(tb assert.TB, s *parity.StoreStub) ([]string, error) {
	tb.Helper()
	return s.List()
}

// storeStubListPins returns the samples of the results of List, with a new
// error as its error result.
func storeStubListPins() ([]string, error) {
	var want0 []string
	want1 := errors.New("List-returns")
	want0 = []string{"List-a"}
	return want0, want1
}

// storeStubClose calls Close of s with the samples of its arguments.
func storeStubClose(tb assert.TB, s *parity.StoreStub) {
	tb.Helper()
	s.Close()
}

// storeStubGetSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Get. The subject calls a new double that
// reports to tb.
func storeStubGetSubject(tb assert.TB) stub.Subject[parity.StoreGetCall, parity.StoreGetReturn] {
	tb.Helper()
	s := parity.NewStoreStub(tb)
	return stub.Subject[parity.StoreGetCall, parity.StoreGetReturn]{
		Stub:  s.OnGet.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = storeStubGet(tb, s)
		},
		Arguments: func(tb assert.TB, c parity.StoreGetCall) {
			tb.Helper()
			a0, a1 := storeStubGetArgs(tb)
			expect.Equal(tb, c.Ctx, a0, "the call records Ctx", expect.ByIdentity())
			expect.Equal(tb, c.ID, a1, "the call records ID")
		},
		Result: func() parity.StoreGetReturn {
			got0, got1 := storeStubGet(tb, s)
			return parity.StoreGetReturn{Item: got0, Err: got1}
		},
		Returns: func() parity.StoreGetReturn {
			want0, want1 := storeStubGetPins()
			s.OnGet.Returns(want0, want1)
			return parity.StoreGetReturn{Item: want0, Err: want1}
		},
		Override: func(seen func(parity.StoreGetCall)) {
			s.OnGet.Func(func(p0 context.Context, p1 string) (string, error) {
				seen(parity.StoreGetCall{Ctx: p0, ID: p1})
				var z0 string
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			_, err := storeStubGet(tb, s)
			return err
		},
	}
}

// storeStubListSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for List. The subject calls a new double that
// reports to tb.
func storeStubListSubject(tb assert.TB) stub.Subject[parity.StoreListCall, parity.StoreListReturn] {
	tb.Helper()
	s := parity.NewStoreStub(tb)
	return stub.Subject[parity.StoreListCall, parity.StoreListReturn]{
		Stub:  s.OnList.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = storeStubList(tb, s)
		},
		Result: func() parity.StoreListReturn {
			got0, got1 := storeStubList(tb, s)
			return parity.StoreListReturn{Result: got0, Err: got1}
		},
		Returns: func() parity.StoreListReturn {
			want0, want1 := storeStubListPins()
			s.OnList.Returns(want0, want1)
			return parity.StoreListReturn{Result: want0, Err: want1}
		},
		Override: func(seen func(parity.StoreListCall)) {
			s.OnList.Func(func() ([]string, error) {
				seen(parity.StoreListCall{})
				var z0 []string
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

// storeStubCloseSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Close. The subject calls a new double that
// reports to tb.
func storeStubCloseSubject(tb assert.TB) stub.Subject[parity.StoreCloseCall, struct{}] {
	tb.Helper()
	s := parity.NewStoreStub(tb)
	return stub.Subject[parity.StoreCloseCall, struct{}]{
		Stub:  s.OnClose.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			storeStubClose(tb, s)
		},
		Override: func(seen func(parity.StoreCloseCall)) {
			s.OnClose.Func(func() {
				seen(parity.StoreCloseCall{})
			})
		},
	}
}

// storeStubInstance is the instance of a parity.StoreStub that the checks
// of the constructors and the options drive.
type storeStubInstance = stub.Instance[parity.StoreGetCall]

// storeStubInstanceOf returns the instance of s. The options apply to the
// whole double, so the checks call Get alone.
func storeStubInstanceOf(tb assert.TB, s *parity.StoreStub) storeStubInstance {
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

	t.Run("List", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "List", storeStubListSubject)
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

		t.Run("List", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "List", storeStubListSubject)
		})

		t.Run("Close", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Close", storeStubCloseSubject)
		})
	})

	t.Run("StoreStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) storeStubInstance {
			return storeStubInstanceOf(tb, parity.NewStoreStub(tb, parity.StoreStubStrict()))
		})
	})

	t.Run("StoreStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Get to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := parity.NewStoreStub(t)
			s := parity.NewStoreStub(t, parity.StoreStubDelegateTo(inner))
			_, _ = storeStubGet(t, s)
			inner.OnGet.AssertCalledOnce(t, "Get calls the implementation")
		})

		t.Run("fails Get with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Get-delegate")
			inner := parity.NewStoreStub(t)
			inner.OnGet.Faults(want, 1)
			s := parity.NewStoreStub(t, parity.StoreStubDelegateTo(inner))
			_, err := storeStubGet(t, s)
			assert.ErrorIs(t, err, want, "Get fails with the error of the implementation")
		})

		t.Run("passes List to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := parity.NewStoreStub(t)
			s := parity.NewStoreStub(t, parity.StoreStubDelegateTo(inner))
			_, _ = storeStubList(t, s)
			inner.OnList.AssertCalledOnce(t, "List calls the implementation")
		})

		t.Run("fails List with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("List-delegate")
			inner := parity.NewStoreStub(t)
			inner.OnList.Faults(want, 1)
			s := parity.NewStoreStub(t, parity.StoreStubDelegateTo(inner))
			_, r1 := storeStubList(t, s)
			assert.ErrorIs(t, r1, want, "List fails with the error of the implementation")
		})

		t.Run("passes Close to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := parity.NewStoreStub(t)
			s := parity.NewStoreStub(t, parity.StoreStubDelegateTo(inner))
			storeStubClose(t, s)
			inner.OnClose.AssertCalledOnce(t, "Close calls the implementation")
		})
	})

	t.Run("StoreStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) storeStubInstance {
			return storeStubInstanceOf(tb, parity.NewStoreStub(tb, parity.StoreStubWithClock(clk)))
		})
	})

	t.Run("StoreStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) storeStubInstance {
			return storeStubInstanceOf(
				tb,
				parity.NewStoreStub(tb, parity.StoreStubWithRandSource(src)),
			)
		})
	})

	t.Run("StoreStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) storeStubInstance {
			return storeStubInstanceOf(tb, parity.NewStoreStub(tb, parity.StoreStubBenchMode()))
		})
	})

	t.Run("StoreStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) storeStubInstance {
				return storeStubInstanceOf(
					tb,
					parity.NewStoreStub(
						tb,
						parity.StoreStubWithClock(clk),
						parity.StoreStubWithLatency(d),
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
			s := parity.NewStoreStub(t, parity.StoreStubFailing(want))
			_, err := storeStubGet(t, s)
			assert.ErrorIs(t, err, want, "Get fails with the error of StoreStubFailing")
		})

		t.Run("fails List with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("List-failing")
			s := parity.NewStoreStub(t, parity.StoreStubFailing(want))
			_, r1 := storeStubList(t, s)
			assert.ErrorIs(t, r1, want, "List fails with the error of StoreStubFailing")
		})
	})

	t.Run("WithStoreGet", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Get", func(t *testing.T) {
			t.Parallel()
			called := false
			s := parity.NewStoreStub(
				t,
				parity.WithStoreGet(func(_ context.Context, _ string) (string, error) {
					called = true
					var z0 string
					var z1 error
					return z0, z1
				}),
			)
			_, _ = storeStubGet(t, s)
			assert.True(t, called, "WithStoreGet sets the function of Get")
		})
	})

	t.Run("WithStoreList", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of List", func(t *testing.T) {
			t.Parallel()
			called := false
			s := parity.NewStoreStub(t, parity.WithStoreList(func() ([]string, error) {
				called = true
				var z0 []string
				var z1 error
				return z0, z1
			}))
			_, _ = storeStubList(t, s)
			assert.True(t, called, "WithStoreList sets the function of List")
		})
	})

	t.Run("WithStoreClose", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Close", func(t *testing.T) {
			t.Parallel()
			called := false
			s := parity.NewStoreStub(t, parity.WithStoreClose(func() {
				called = true
			}))
			storeStubClose(t, s)
			assert.True(t, called, "WithStoreClose sets the function of Close")
		})
	})
}

// TestStoreGetStub checks the methods of StoreGetStub, which configures
// Get.
func TestStoreGetStub(t *testing.T) {
	t.Parallel()

	t.Run("FaultNotFound", func(t *testing.T) {
		t.Parallel()

		t.Run("fails the next call of Get with ErrNotFound", func(t *testing.T) {
			t.Parallel()
			s := parity.NewStoreStub(t)
			s.OnGet.FaultNotFound()
			_, err := storeStubGet(t, s)
			assert.ErrorIs(
				t,
				err,
				parity.ErrNotFound,
				"the next call of Get fails with ErrNotFound",
			)
		})

		t.Run("spares the call after the failed call", func(t *testing.T) {
			t.Parallel()
			s := parity.NewStoreStub(t)
			s.OnGet.FaultNotFound()
			_, _ = storeStubGet(t, s)
			_, err := storeStubGet(t, s)
			assert.NoError(t, err, "the call of Get after the failed call succeeds")
		})
	})

	t.Run("FaultGone", func(t *testing.T) {
		t.Parallel()

		t.Run("fails the next call of Get with ErrGone", func(t *testing.T) {
			t.Parallel()
			s := parity.NewStoreStub(t)
			s.OnGet.FaultGone()
			_, err := storeStubGet(t, s)
			assert.ErrorIs(t, err, parity.ErrGone, "the next call of Get fails with ErrGone")
		})

		t.Run("spares the call after the failed call", func(t *testing.T) {
			t.Parallel()
			s := parity.NewStoreStub(t)
			s.OnGet.FaultGone()
			_, _ = storeStubGet(t, s)
			_, err := storeStubGet(t, s)
			assert.NoError(t, err, "the call of Get after the failed call succeeds")
		})
	})

	t.Run("FaultForPartition", func(t *testing.T) {
		t.Parallel()

		t.Run("fails a call of Get in the partition of the key", func(t *testing.T) {
			t.Parallel()
			s := parity.NewStoreStub(t)
			_, key := storeStubGetArgs(t)
			want := errors.New("Get-partition")
			s.OnGet.FaultForPartition(key, want, 1)
			_, err := storeStubGet(t, s)
			assert.ErrorIs(t, err, want, "a call in the partition of the key fails")
		})
	})

	t.Run("FaultForOtherPartitions", func(t *testing.T) {
		t.Parallel()

		t.Run("spares a call of Get in the partition of the key", func(t *testing.T) {
			t.Parallel()
			s := parity.NewStoreStub(t)
			_, key := storeStubGetArgs(t)
			s.OnGet.FaultForOtherPartitions(key, errors.New("Get-partition"), 1)
			_, err := storeStubGet(t, s)
			assert.NoError(t, err, "a call in the partition of the key succeeds")
		})
	})
}

// TestStoreListStub checks the methods of StoreListStub, which configures
// List.
func TestStoreListStub(t *testing.T) {
	t.Parallel()

	t.Run("RetrySchedule", func(t *testing.T) {
		t.Parallel()

		t.Run("fails each call of List before call 3", func(t *testing.T) {
			t.Parallel()
			s := parity.NewStoreStub(t)
			want := errors.New("List-retry")
			s.OnList.RetrySchedule(want)
			for attempt := 1; attempt < 3; attempt++ {
				_, r1 := storeStubList(t, s)
				assert.ErrorIs(t, r1, want, "a call of List before call 3 fails")
			}
		})

		t.Run("spares call 3 of List", func(t *testing.T) {
			t.Parallel()
			s := parity.NewStoreStub(t)
			s.OnList.RetrySchedule(errors.New("List-retry"))
			for attempt := 1; attempt < 3; attempt++ {
				_, _ = storeStubList(t, s)
			}
			_, r1 := storeStubList(t, s)
			assert.NoError(t, r1, "call 3 of List succeeds")
		})
	})
}
