package scoped_test

import (
	errors2 "errors"
	rand2 "math/rand/v2"
	testing2 "testing"
	time2 "time"

	"example.com/shop/scoped"
	assert2 "go.dokimi.dev/assert"
	expect2 "go.dokimi.dev/assert/expect"
	stub2 "go.dokimi.dev/stub"
)

// storeStubGetArgs returns the samples of the arguments of Get.
func storeStubGetArgs(tb assert2.TB) string {
	tb.Helper()
	var a0 string
	a0 = "Get-a"
	return a0
}

// storeStubGet calls Get of s with the samples of its arguments.
func storeStubGet(tb assert2.TB, s *scoped.StoreStub) (string, error) {
	tb.Helper()
	a0 := storeStubGetArgs(tb)
	return s.Get(a0)
}

// storeStubGetPins returns the samples of the results of Get, with a new
// error as its error result.
func storeStubGetPins() (string, error) {
	var want0 string
	want1 := errors2.New("Get-returns")
	want0 = "Get-a"
	return want0, want1
}

// storeStubList calls List of s with the samples of its arguments.
func storeStubList(tb assert2.TB, s *scoped.StoreStub) ([]string, error) {
	tb.Helper()
	return s.List()
}

// storeStubListPins returns the samples of the results of List, with a new
// error as its error result.
func storeStubListPins() ([]string, error) {
	var want0 []string
	want1 := errors2.New("List-returns")
	want0 = []string{"List-a"}
	return want0, want1
}

// storeStubGetSubject returns the subject of [stub2.Behaviour] for Get.
// The subject calls a new double that reports to tb.
func storeStubGetSubject(tb assert2.TB) stub2.Subject[scoped.StoreGetCall, scoped.StoreGetReturn] {
	tb.Helper()
	s := scoped.NewStoreStub(tb)
	return stub2.Subject[scoped.StoreGetCall, scoped.StoreGetReturn]{
		Stub: s.OnGet.MethodStub,
		Call: func() {
			_, _ = storeStubGet(tb, s)
		},
		Result: func() scoped.StoreGetReturn {
			got0, got1 := storeStubGet(tb, s)
			return scoped.StoreGetReturn{Result: got0, Err: got1}
		},
		Override: func(mark func()) {
			s.OnGet.Func(func(_ string) (string, error) {
				mark()
				var z0 string
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

// storeStubListSubject returns the subject of [stub2.Behaviour] for List.
// The subject calls a new double that reports to tb.
func storeStubListSubject(
	tb assert2.TB,
) stub2.Subject[scoped.StoreListCall, scoped.StoreListReturn] {
	tb.Helper()
	s := scoped.NewStoreStub(tb)
	return stub2.Subject[scoped.StoreListCall, scoped.StoreListReturn]{
		Stub: s.OnList.MethodStub,
		Call: func() {
			_, _ = storeStubList(tb, s)
		},
		Result: func() scoped.StoreListReturn {
			got0, got1 := storeStubList(tb, s)
			return scoped.StoreListReturn{Result: got0, Err: got1}
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

// storeStubInstance is the instance of a scoped.StoreStub that the checks
// of the constructors and the options drive.
type storeStubInstance = stub2.Instance[scoped.StoreGetCall]

// storeStubInstanceOf returns the instance of s. The options apply to the
// whole double, so the checks call Get alone.
func storeStubInstanceOf(tb assert2.TB, s *scoped.StoreStub) storeStubInstance {
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
func TestStoreStub(t *testing2.T) {
	t.Parallel()

	t.Run("Get", func(t *testing2.T) {
		t.Parallel()
		stub2.Behaviour(t, "Get", storeStubGetSubject)

		t.Run("records the arguments of the call", func(t *testing2.T) {
			t.Parallel()
			s := scoped.NewStoreStub(t)
			_, _ = storeStubGet(t, s)
			got := s.OnGet.AssertCalledOnce(t, "Get records the call")
			a0 := storeStubGetArgs(t)
			assert2.Equal(t, got.ID, a0, "the call records ID")
		})

		t.Run("runs the OnRecord hook once for each call", func(t *testing2.T) {
			t.Parallel()
			s := scoped.NewStoreStub(t)
			var seen []scoped.StoreGetCall
			s.OnGet.OnRecord(func(c scoped.StoreGetCall) { seen = append(seen, c) })
			_, _ = storeStubGet(t, s)
			_, _ = storeStubGet(t, s)
			assert2.Length(t, seen, 2, "OnRecord runs once for each call of Get")
		})
	})

	t.Run("List", func(t *testing2.T) {
		t.Parallel()
		stub2.Behaviour(t, "List", storeStubListSubject)

		t.Run("runs the OnRecord hook once for each call", func(t *testing2.T) {
			t.Parallel()
			s := scoped.NewStoreStub(t)
			var seen []scoped.StoreListCall
			s.OnList.OnRecord(func(c scoped.StoreListCall) { seen = append(seen, c) })
			_, _ = storeStubList(t, s)
			_, _ = storeStubList(t, s)
			assert2.Length(t, seen, 2, "OnRecord runs once for each call of List")
		})
	})

	t.Run("ResetCalls", func(t *testing2.T) {
		t.Parallel()
		stub2.ResetBehaviour(t, func(tb assert2.TB) storeStubInstance {
			return storeStubInstanceOf(tb, scoped.NewStoreStub(tb))
		})

		t.Run("clears the recorded calls of Get", func(t *testing2.T) {
			t.Parallel()
			s := scoped.NewStoreStub(t)
			_, _ = storeStubGet(t, s)
			s.ResetCalls()
			assert2.Equal(t, s.OnGet.CallCount(), 0, "ResetCalls clears the calls of Get")
		})

		t.Run("keeps the results of Returns of Get", func(t *testing2.T) {
			t.Parallel()
			s := scoped.NewStoreStub(t)
			want0, want1 := storeStubGetPins()
			s.OnGet.Returns(want0, want1)
			s.ResetCalls()
			got0, got1 := storeStubGet(t, s)
			expect2.Equal(t, got0, want0, "ResetCalls keeps the Result of Get")
			expect2.ErrorIs(t, got1, want1, "ResetCalls keeps the error of Get")
		})

		t.Run("clears the recorded calls of List", func(t *testing2.T) {
			t.Parallel()
			s := scoped.NewStoreStub(t)
			_, _ = storeStubList(t, s)
			s.ResetCalls()
			assert2.Equal(t, s.OnList.CallCount(), 0, "ResetCalls clears the calls of List")
		})

		t.Run("keeps the results of Returns of List", func(t *testing2.T) {
			t.Parallel()
			s := scoped.NewStoreStub(t)
			want0, want1 := storeStubListPins()
			s.OnList.Returns(want0, want1)
			s.ResetCalls()
			got0, got1 := storeStubList(t, s)
			expect2.Equal(t, got0, want0, "ResetCalls keeps the Result of List")
			expect2.ErrorIs(t, got1, want1, "ResetCalls keeps the error of List")
		})
	})

	t.Run("NewStoreStubDelayed", func(t *testing2.T) {
		t.Parallel()
		stub2.LatencyBehaviour(
			t,
			func(tb assert2.TB, d time2.Duration, clk assert2.Clock) storeStubInstance {
				return storeStubInstanceOf(tb, scoped.NewStoreStubDelayed(tb, d, clk))
			},
		)
	})

	t.Run("NewStoreStubFailing", func(t *testing2.T) {
		t.Parallel()

		t.Run("fails Get with its error", func(t *testing2.T) {
			t.Parallel()
			want := errors2.New("Get-failing")
			s := scoped.NewStoreStubFailing(t, want)
			_, r1 := storeStubGet(t, s)
			assert2.ErrorIs(t, r1, want, "Get fails with the error of NewStoreStubFailing")
		})

		t.Run("fails List with its error", func(t *testing2.T) {
			t.Parallel()
			want := errors2.New("List-failing")
			s := scoped.NewStoreStubFailing(t, want)
			_, r1 := storeStubList(t, s)
			assert2.ErrorIs(t, r1, want, "List fails with the error of NewStoreStubFailing")
		})
	})

	t.Run("StoreStubStrict", func(t *testing2.T) {
		t.Parallel()
		stub2.StrictBehaviour(t, func(tb assert2.TB) storeStubInstance {
			return storeStubInstanceOf(tb, scoped.NewStoreStub(tb, scoped.StoreStubStrict()))
		})
	})

	t.Run("StoreStubDelegateTo", func(t *testing2.T) {
		t.Parallel()

		t.Run("passes Get to the implementation", func(t *testing2.T) {
			t.Parallel()
			inner := scoped.NewStoreStub(t)
			s := scoped.NewStoreStub(t, scoped.StoreStubDelegateTo(inner))
			_, _ = storeStubGet(t, s)
			inner.OnGet.AssertCalledOnce(t, "Get calls the implementation")
		})

		t.Run("fails Get with the error of the implementation", func(t *testing2.T) {
			t.Parallel()
			want := errors2.New("Get-delegate")
			inner := scoped.NewStoreStub(t)
			inner.OnGet.Faults(want, 1)
			s := scoped.NewStoreStub(t, scoped.StoreStubDelegateTo(inner))
			_, r1 := storeStubGet(t, s)
			assert2.ErrorIs(t, r1, want, "Get fails with the error of the implementation")
		})

		t.Run("passes List to the implementation", func(t *testing2.T) {
			t.Parallel()
			inner := scoped.NewStoreStub(t)
			s := scoped.NewStoreStub(t, scoped.StoreStubDelegateTo(inner))
			_, _ = storeStubList(t, s)
			inner.OnList.AssertCalledOnce(t, "List calls the implementation")
		})

		t.Run("fails List with the error of the implementation", func(t *testing2.T) {
			t.Parallel()
			want := errors2.New("List-delegate")
			inner := scoped.NewStoreStub(t)
			inner.OnList.Faults(want, 1)
			s := scoped.NewStoreStub(t, scoped.StoreStubDelegateTo(inner))
			_, r1 := storeStubList(t, s)
			assert2.ErrorIs(t, r1, want, "List fails with the error of the implementation")
		})
	})

	t.Run("StoreStubWithClock", func(t *testing2.T) {
		t.Parallel()
		stub2.ClockBehaviour(t, func(tb assert2.TB, clk assert2.Clock) storeStubInstance {
			return storeStubInstanceOf(tb, scoped.NewStoreStub(tb, scoped.StoreStubWithClock(clk)))
		})
	})

	t.Run("StoreStubWithRandSource", func(t *testing2.T) {
		t.Parallel()
		stub2.RandSourceBehaviour(t, func(tb assert2.TB, src rand2.Source) storeStubInstance {
			return storeStubInstanceOf(
				tb,
				scoped.NewStoreStub(tb, scoped.StoreStubWithRandSource(src)),
			)
		})
	})

	t.Run("StoreStubBenchMode", func(t *testing2.T) {
		t.Parallel()
		stub2.BenchModeBehaviour(t, func(tb assert2.TB) storeStubInstance {
			return storeStubInstanceOf(tb, scoped.NewStoreStub(tb, scoped.StoreStubBenchMode()))
		})
	})

	t.Run("StoreStubWithLatency", func(t *testing2.T) {
		t.Parallel()
		stub2.LatencyBehaviour(
			t,
			func(tb assert2.TB, d time2.Duration, clk assert2.Clock) storeStubInstance {
				return storeStubInstanceOf(
					tb,
					scoped.NewStoreStub(
						tb,
						scoped.StoreStubWithClock(clk),
						scoped.StoreStubWithLatency(d),
					),
				)
			},
		)
	})

	t.Run("StoreStubFailing", func(t *testing2.T) {
		t.Parallel()

		t.Run("fails Get with its error", func(t *testing2.T) {
			t.Parallel()
			want := errors2.New("Get-failing")
			s := scoped.NewStoreStub(t, scoped.StoreStubFailing(want))
			_, r1 := storeStubGet(t, s)
			assert2.ErrorIs(t, r1, want, "Get fails with the error of StoreStubFailing")
		})

		t.Run("fails List with its error", func(t *testing2.T) {
			t.Parallel()
			want := errors2.New("List-failing")
			s := scoped.NewStoreStub(t, scoped.StoreStubFailing(want))
			_, r1 := storeStubList(t, s)
			assert2.ErrorIs(t, r1, want, "List fails with the error of StoreStubFailing")
		})
	})

	t.Run("WithStoreGet", func(t *testing2.T) {
		t.Parallel()

		t.Run("sets the function of Get", func(t *testing2.T) {
			t.Parallel()
			called := false
			s := scoped.NewStoreStub(t, scoped.WithStoreGet(func(_ string) (string, error) {
				called = true
				var z0 string
				var z1 error
				return z0, z1
			}))
			_, _ = storeStubGet(t, s)
			assert2.True(t, called, "WithStoreGet sets the function of Get")
		})
	})

	t.Run("WithStoreList", func(t *testing2.T) {
		t.Parallel()

		t.Run("sets the function of List", func(t *testing2.T) {
			t.Parallel()
			called := false
			s := scoped.NewStoreStub(t, scoped.WithStoreList(func() ([]string, error) {
				called = true
				var z0 []string
				var z1 error
				return z0, z1
			}))
			_, _ = storeStubList(t, s)
			assert2.True(t, called, "WithStoreList sets the function of List")
		})
	})
}

// TestStoreGetStub checks the methods of StoreGetStub, which configures
// Get.
func TestStoreGetStub(t *testing2.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing2.T) {
		t.Parallel()

		t.Run("sets the results of each call of Get", func(t *testing2.T) {
			t.Parallel()
			s := scoped.NewStoreStub(t)
			want0, want1 := storeStubGetPins()
			s.OnGet.Returns(want0, want1)
			got0, got1 := storeStubGet(t, s)
			expect2.Equal(t, got0, want0, "Get returns the Result that Returns sets")
			expect2.ErrorIs(t, got1, want1, "Get returns the error that Returns sets")
		})

		t.Run("configures Get for strict mode", func(t *testing2.T) {
			t.Parallel()
			r := assert2.NewRecorder()
			s := scoped.NewStoreStub(r)
			s.OnGet.Strict()
			var want0 string
			var want1 error
			s.OnGet.Returns(want0, want1)
			_, _ = storeStubGet(t, s)
			assert2.False(t, r.Failed(), "strict mode accepts the call of Get")
		})
	})

	t.Run("Func", func(t *testing2.T) {
		t.Parallel()

		t.Run("passes the arguments of each call of Get to the function", func(t *testing2.T) {
			t.Parallel()
			s := scoped.NewStoreStub(t)
			var got0 string
			s.OnGet.Func(func(p0 string) (string, error) {
				got0 = p0
				var z0 string
				var z1 error
				return z0, z1
			})
			_, _ = storeStubGet(t, s)
			a0 := storeStubGetArgs(t)
			assert2.Equal(t, got0, a0, "the function receives ID")
		})
	})

	t.Run("FaultForPartition", func(t *testing2.T) {
		t.Parallel()

		t.Run("fails a call of Get in the partition of the key", func(t *testing2.T) {
			t.Parallel()
			s := scoped.NewStoreStub(t)
			key := storeStubGetArgs(t)
			want := errors2.New("Get-partition")
			s.OnGet.FaultForPartition(key, want, 1)
			_, r1 := storeStubGet(t, s)
			assert2.ErrorIs(t, r1, want, "a call in the partition of the key fails")
		})
	})

	t.Run("FaultForOtherPartitions", func(t *testing2.T) {
		t.Parallel()

		t.Run("spares a call of Get in the partition of the key", func(t *testing2.T) {
			t.Parallel()
			s := scoped.NewStoreStub(t)
			key := storeStubGetArgs(t)
			s.OnGet.FaultForOtherPartitions(key, errors2.New("Get-partition"), 1)
			_, r1 := storeStubGet(t, s)
			assert2.NoError(t, r1, "a call in the partition of the key succeeds")
		})
	})
}

// TestStoreListStub checks the methods of StoreListStub, which configures
// List.
func TestStoreListStub(t *testing2.T) {
	t.Parallel()

	t.Run("Returns", func(t *testing2.T) {
		t.Parallel()

		t.Run("sets the results of each call of List", func(t *testing2.T) {
			t.Parallel()
			s := scoped.NewStoreStub(t)
			want0, want1 := storeStubListPins()
			s.OnList.Returns(want0, want1)
			got0, got1 := storeStubList(t, s)
			expect2.Equal(t, got0, want0, "List returns the Result that Returns sets")
			expect2.ErrorIs(t, got1, want1, "List returns the error that Returns sets")
		})

		t.Run("configures List for strict mode", func(t *testing2.T) {
			t.Parallel()
			r := assert2.NewRecorder()
			s := scoped.NewStoreStub(r)
			s.OnList.Strict()
			var want0 []string
			var want1 error
			s.OnList.Returns(want0, want1)
			_, _ = storeStubList(t, s)
			assert2.False(t, r.Failed(), "strict mode accepts the call of List")
		})
	})

	t.Run("RetrySchedule", func(t *testing2.T) {
		t.Parallel()

		t.Run("fails each call of List before call 3", func(t *testing2.T) {
			t.Parallel()
			s := scoped.NewStoreStub(t)
			want := errors2.New("List-retry")
			s.OnList.RetrySchedule(want)
			for attempt := 1; attempt < 3; attempt++ {
				_, r1 := storeStubList(t, s)
				assert2.ErrorIs(t, r1, want, "a call of List before call 3 fails")
			}
		})

		t.Run("spares call 3 of List", func(t *testing2.T) {
			t.Parallel()
			s := scoped.NewStoreStub(t)
			s.OnList.RetrySchedule(errors2.New("List-retry"))
			for attempt := 1; attempt < 3; attempt++ {
				_, _ = storeStubList(t, s)
			}
			_, r1 := storeStubList(t, s)
			assert2.NoError(t, r1, "call 3 of List succeeds")
		})
	})
}
