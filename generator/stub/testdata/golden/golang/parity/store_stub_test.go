package parity

import (
	"context"
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

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
func storeStubGet(tb assert.TB, s *StoreStub) (string, error) {
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
func storeStubList(tb assert.TB, s *StoreStub) ([]string, error) {
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
func storeStubClose(tb assert.TB, s *StoreStub) {
	tb.Helper()
	s.Close()
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
			_, err := storeStubGet(tb, s)
			return err
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
		Override: func(mark func()) {
			s.OnList.Func(func() ([]string, error) {
				mark()
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

// storeStubCloseSubject returns the subject of [stub.Behaviour] for Close.
// The subject calls a new double that reports to tb.
func storeStubCloseSubject(tb assert.TB) stub.Subject[StoreCloseCall, struct{}] {
	tb.Helper()
	s := NewStoreStub(tb)
	return stub.Subject[StoreCloseCall, struct{}]{
		Stub: s.OnClose.MethodStub,
		Call: func() {
			storeStubClose(tb, s)
		},
		Override: func(mark func()) {
			s.OnClose.Func(func() {
				mark()
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

	t.Run("List", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "List", storeStubListSubject)

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

	t.Run("Close", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Close", storeStubCloseSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			var seen []StoreCloseCall
			s.OnClose.OnRecord(func(c StoreCloseCall) { seen = append(seen, c) })
			storeStubClose(t, s)
			storeStubClose(t, s)
			assert.Length(t, seen, 2, "OnRecord runs once for each call of Close")
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
			expect.Equal(t, got0, want0, "ResetCalls keeps the Item of Get")
			expect.ErrorIs(t, got1, want1, "ResetCalls keeps the error of Get")
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

		t.Run("clears the recorded calls of Close", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			storeStubClose(t, s)
			s.ResetCalls()
			assert.Equal(t, s.OnClose.CallCount(), 0, "ResetCalls clears the calls of Close")
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
			_, err := storeStubGet(t, s)
			assert.ErrorIs(t, err, want, "Get fails with the error of NewStoreStubFailing")
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
			_, err := storeStubGet(t, s)
			assert.ErrorIs(t, err, want, "Get fails with the error of the implementation")
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

		t.Run("passes Close to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := NewStoreStub(t)
			s := NewStoreStub(t, StoreStubDelegateTo(inner))
			storeStubClose(t, s)
			inner.OnClose.AssertCalledOnce(t, "Close calls the implementation")
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
			_, err := storeStubGet(t, s)
			assert.ErrorIs(t, err, want, "Get fails with the error of StoreStubFailing")
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
			s := NewStoreStub(t, WithStoreGet(func(_ context.Context, _ string) (string, error) {
				called = true
				var z0 string
				var z1 error
				return z0, z1
			}))
			_, _ = storeStubGet(t, s)
			assert.True(t, called, "WithStoreGet sets the function of Get")
		})
	})

	t.Run("WithStoreList", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of List", func(t *testing.T) {
			t.Parallel()
			called := false
			s := NewStoreStub(t, WithStoreList(func() ([]string, error) {
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
			s := NewStoreStub(t, WithStoreClose(func() {
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

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of Get", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			want0, want1 := storeStubGetPins()
			s.OnGet.Returns(want0, want1)
			got0, got1 := storeStubGet(t, s)
			expect.Equal(t, got0, want0, "Get returns the Item that Returns sets")
			expect.ErrorIs(t, got1, want1, "Get returns the error that Returns sets")
		})

		t.Run("configures Get for strict mode", func(t *testing.T) {
			t.Parallel()
			r := assert.NewRecorder()
			s := NewStoreStub(r)
			s.OnGet.Strict()
			var want0 string
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
			s.OnGet.Func(func(p0 context.Context, p1 string) (string, error) {
				got0, got1 = p0, p1
				var z0 string
				var z1 error
				return z0, z1
			})
			_, _ = storeStubGet(t, s)
			a0, a1 := storeStubGetArgs(t)
			expect.Equal(t, got0, a0, "the function receives Ctx", expect.ByIdentity())
			expect.Equal(t, got1, a1, "the function receives ID")
		})
	})

	t.Run("FaultNotFound", func(t *testing.T) {
		t.Parallel()

		t.Run("fails the next call of Get with ErrNotFound", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			s.OnGet.FaultNotFound()
			_, err := storeStubGet(t, s)
			assert.ErrorIs(t, err, ErrNotFound, "the next call of Get fails with ErrNotFound")
		})

		t.Run("spares the call after the failed call", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
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
			s := NewStoreStub(t)
			s.OnGet.FaultGone()
			_, err := storeStubGet(t, s)
			assert.ErrorIs(t, err, ErrGone, "the next call of Get fails with ErrGone")
		})

		t.Run("spares the call after the failed call", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
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
			s := NewStoreStub(t)
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
			s := NewStoreStub(t)
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

	t.Run("Returns", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the results of each call of List", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			want0, want1 := storeStubListPins()
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
			var want0 []string
			var want1 error
			s.OnList.Returns(want0, want1)
			_, _ = storeStubList(t, s)
			assert.False(t, r.Failed(), "strict mode accepts the call of List")
		})
	})

	t.Run("RetrySchedule", func(t *testing.T) {
		t.Parallel()

		t.Run("fails each call of List before call 3", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			want := errors.New("List-retry")
			s.OnList.RetrySchedule(want)
			for attempt := 1; attempt < 3; attempt++ {
				_, r1 := storeStubList(t, s)
				assert.ErrorIs(t, r1, want, "a call of List before call 3 fails")
			}
		})

		t.Run("spares call 3 of List", func(t *testing.T) {
			t.Parallel()
			s := NewStoreStub(t)
			s.OnList.RetrySchedule(errors.New("List-retry"))
			for attempt := 1; attempt < 3; attempt++ {
				_, _ = storeStubList(t, s)
			}
			_, r1 := storeStubList(t, s)
			assert.NoError(t, r1, "call 3 of List succeeds")
		})
	})
}
