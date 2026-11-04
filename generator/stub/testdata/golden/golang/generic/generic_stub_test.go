package generic_test

import (
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"example.com/shop/generic"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// repoStubGetArgs returns the samples of the arguments of Get.
func repoStubGetArgs[T any](tb assert.TB) string {
	tb.Helper()
	var a0 string
	a0 = "Get-a"
	return a0
}

// repoStubGet calls Get of s with the samples of its arguments.
func repoStubGet[T any](tb assert.TB, s *generic.RepoStub[T]) (T, error) {
	tb.Helper()
	a0 := repoStubGetArgs[T](tb)
	return s.Get(a0)
}

// repoStubGetPins returns the samples of the results of Get, with a new
// error as its error result.
func repoStubGetPins[T any]() (T, error) {
	var want0 T
	want1 := errors.New("Get-returns")
	return want0, want1
}

// repoStubPutArgs returns the samples of the arguments of Put.
func repoStubPutArgs[T any](tb assert.TB) T {
	tb.Helper()
	var a0 T
	return a0
}

// repoStubPut calls Put of s with the samples of its arguments.
func repoStubPut[T any](tb assert.TB, s *generic.RepoStub[T]) {
	tb.Helper()
	a0 := repoStubPutArgs[T](tb)
	s.Put(a0)
}

// repoStubList calls List of s with the samples of its arguments.
func repoStubList[T any](tb assert.TB, s *generic.RepoStub[T]) ([]T, error) {
	tb.Helper()
	return s.List()
}

// repoStubListPins returns the samples of the results of List, with a new
// error as its error result.
func repoStubListPins[T any]() ([]T, error) {
	var want0 []T
	want1 := errors.New("List-returns")
	return want0, want1
}

// repoStubBoxed calls Boxed of s with the samples of its arguments.
func repoStubBoxed[T any](tb assert.TB, s *generic.RepoStub[T]) generic.Box[T] {
	tb.Helper()
	return s.Boxed()
}

// repoStubGetSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Get. The subject calls a new double that
// reports to tb.
func repoStubGetSubject[T any](
	tb assert.TB,
) stub.Subject[generic.RepoGetCall[T], generic.RepoGetReturn[T]] {
	tb.Helper()
	s := generic.NewRepoStub[T](tb)
	return stub.Subject[generic.RepoGetCall[T], generic.RepoGetReturn[T]]{
		Stub:  s.OnGet.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = repoStubGet[T](tb, s)
		},
		Arguments: func(tb assert.TB, c generic.RepoGetCall[T]) {
			tb.Helper()
			a0 := repoStubGetArgs[T](tb)
			assert.Equal(tb, c.ID, a0, "the call records ID")
		},
		Result: func() generic.RepoGetReturn[T] {
			got0, got1 := repoStubGet[T](tb, s)
			return generic.RepoGetReturn[T]{Result: got0, Err: got1}
		},
		Returns: func() generic.RepoGetReturn[T] {
			want0, want1 := repoStubGetPins[T]()
			s.OnGet.Returns(want0, want1)
			return generic.RepoGetReturn[T]{Result: want0, Err: want1}
		},
		Override: func(seen func(generic.RepoGetCall[T])) {
			s.OnGet.Func(func(p0 string) (T, error) {
				seen(generic.RepoGetCall[T]{ID: p0})
				var z0 T
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			_, r1 := repoStubGet[T](tb, s)
			return r1
		},
	}
}

// repoStubPutSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Put. The subject calls a new double that
// reports to tb.
func repoStubPutSubject[T any](tb assert.TB) stub.Subject[generic.RepoPutCall[T], struct{}] {
	tb.Helper()
	s := generic.NewRepoStub[T](tb)
	return stub.Subject[generic.RepoPutCall[T], struct{}]{
		Stub:  s.OnPut.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			repoStubPut[T](tb, s)
		},
		Arguments: func(tb assert.TB, c generic.RepoPutCall[T]) {
			tb.Helper()
			a0 := repoStubPutArgs[T](tb)
			assert.Equal(tb, c.V, a0, "the call records V")
		},
		Override: func(seen func(generic.RepoPutCall[T])) {
			s.OnPut.Func(func(p0 T) {
				seen(generic.RepoPutCall[T]{V: p0})
			})
		},
	}
}

// repoStubListSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for List. The subject calls a new double that
// reports to tb.
func repoStubListSubject[T any](
	tb assert.TB,
) stub.Subject[generic.RepoListCall[T], generic.RepoListReturn[T]] {
	tb.Helper()
	s := generic.NewRepoStub[T](tb)
	return stub.Subject[generic.RepoListCall[T], generic.RepoListReturn[T]]{
		Stub:  s.OnList.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_, _ = repoStubList[T](tb, s)
		},
		Result: func() generic.RepoListReturn[T] {
			got0, got1 := repoStubList[T](tb, s)
			return generic.RepoListReturn[T]{Result: got0, Err: got1}
		},
		Returns: func() generic.RepoListReturn[T] {
			want0, want1 := repoStubListPins[T]()
			s.OnList.Returns(want0, want1)
			return generic.RepoListReturn[T]{Result: want0, Err: want1}
		},
		Override: func(seen func(generic.RepoListCall[T])) {
			s.OnList.Func(func() ([]T, error) {
				seen(generic.RepoListCall[T]{})
				var z0 []T
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			_, r1 := repoStubList[T](tb, s)
			return r1
		},
	}
}

// repoStubBoxedSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Boxed. The subject calls a new double that
// reports to tb.
func repoStubBoxedSubject[T any](
	tb assert.TB,
) stub.Subject[generic.RepoBoxedCall[T], generic.RepoBoxedReturn[T]] {
	tb.Helper()
	s := generic.NewRepoStub[T](tb)
	return stub.Subject[generic.RepoBoxedCall[T], generic.RepoBoxedReturn[T]]{
		Stub:  s.OnBoxed.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = repoStubBoxed[T](tb, s)
		},
		Result: func() generic.RepoBoxedReturn[T] {
			got0 := repoStubBoxed[T](tb, s)
			return generic.RepoBoxedReturn[T]{Result: got0}
		},
		Override: func(seen func(generic.RepoBoxedCall[T])) {
			s.OnBoxed.Func(func() generic.Box[T] {
				seen(generic.RepoBoxedCall[T]{})
				var z0 generic.Box[T]
				return z0
			})
		},
	}
}

// repoStubInstance is the instance of a generic.RepoStub that the checks
// of the constructors and the options drive.
type repoStubInstance[T any] = stub.Instance[generic.RepoGetCall[T]]

// repoStubInstanceOf returns the instance of s. The options apply to the
// whole double, so the checks call Get alone.
func repoStubInstanceOf[T any](tb assert.TB, s *generic.RepoStub[T]) repoStubInstance[T] {
	tb.Helper()
	return repoStubInstance[T]{
		Stub: s.OnGet.MethodStub,
		Call: func() {
			_, _ = repoStubGet[T](tb, s)
		},
	}
}

// checkRepoStub checks the methods, the constructors and the options of
// the double RepoStub at its type parameters.
func checkRepoStub[T any](t *testing.T) {
	t.Helper()

	t.Run("Get", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Get", repoStubGetSubject[T])
	})

	t.Run("Put", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Put", repoStubPutSubject[T])
	})

	t.Run("List", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "List", repoStubListSubject[T])
	})

	t.Run("Boxed", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "Boxed", repoStubBoxedSubject[T])
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()

		t.Run("Get", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Get", repoStubGetSubject[T])
		})

		t.Run("Put", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Put", repoStubPutSubject[T])
		})

		t.Run("List", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "List", repoStubListSubject[T])
		})

		t.Run("Boxed", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "Boxed", repoStubBoxedSubject[T])
		})
	})

	t.Run("RepoStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) repoStubInstance[T] {
			return repoStubInstanceOf[T](
				tb,
				generic.NewRepoStub[T](tb, generic.RepoStubStrict[T]()),
			)
		})
	})

	t.Run("RepoStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes Get to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := generic.NewRepoStub[T](t)
			s := generic.NewRepoStub[T](t, generic.RepoStubDelegateTo[T](inner))
			_, _ = repoStubGet[T](t, s)
			inner.OnGet.AssertCalledOnce(t, "Get calls the implementation")
		})

		t.Run("fails Get with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Get-delegate")
			inner := generic.NewRepoStub[T](t)
			inner.OnGet.Faults(want, 1)
			s := generic.NewRepoStub[T](t, generic.RepoStubDelegateTo[T](inner))
			_, r1 := repoStubGet[T](t, s)
			assert.ErrorIs(t, r1, want, "Get fails with the error of the implementation")
		})

		t.Run("passes Put to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := generic.NewRepoStub[T](t)
			s := generic.NewRepoStub[T](t, generic.RepoStubDelegateTo[T](inner))
			repoStubPut[T](t, s)
			inner.OnPut.AssertCalledOnce(t, "Put calls the implementation")
		})

		t.Run("passes List to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := generic.NewRepoStub[T](t)
			s := generic.NewRepoStub[T](t, generic.RepoStubDelegateTo[T](inner))
			_, _ = repoStubList[T](t, s)
			inner.OnList.AssertCalledOnce(t, "List calls the implementation")
		})

		t.Run("fails List with the error of the implementation", func(t *testing.T) {
			t.Parallel()
			want := errors.New("List-delegate")
			inner := generic.NewRepoStub[T](t)
			inner.OnList.Faults(want, 1)
			s := generic.NewRepoStub[T](t, generic.RepoStubDelegateTo[T](inner))
			_, r1 := repoStubList[T](t, s)
			assert.ErrorIs(t, r1, want, "List fails with the error of the implementation")
		})

		t.Run("passes Boxed to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := generic.NewRepoStub[T](t)
			s := generic.NewRepoStub[T](t, generic.RepoStubDelegateTo[T](inner))
			_ = repoStubBoxed[T](t, s)
			inner.OnBoxed.AssertCalledOnce(t, "Boxed calls the implementation")
		})
	})

	t.Run("RepoStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) repoStubInstance[T] {
			return repoStubInstanceOf[T](
				tb,
				generic.NewRepoStub[T](tb, generic.RepoStubWithClock[T](clk)),
			)
		})
	})

	t.Run("RepoStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) repoStubInstance[T] {
			return repoStubInstanceOf[T](
				tb,
				generic.NewRepoStub[T](tb, generic.RepoStubWithRandSource[T](src)),
			)
		})
	})

	t.Run("RepoStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) repoStubInstance[T] {
			return repoStubInstanceOf[T](
				tb,
				generic.NewRepoStub[T](tb, generic.RepoStubBenchMode[T]()),
			)
		})
	})

	t.Run("RepoStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) repoStubInstance[T] {
				return repoStubInstanceOf[T](
					tb,
					generic.NewRepoStub[T](
						tb,
						generic.RepoStubWithClock[T](clk),
						generic.RepoStubWithLatency[T](d),
					),
				)
			},
		)
	})

	t.Run("RepoStubFailing", func(t *testing.T) {
		t.Parallel()

		t.Run("fails Get with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("Get-failing")
			s := generic.NewRepoStub[T](t, generic.RepoStubFailing[T](want))
			_, r1 := repoStubGet[T](t, s)
			assert.ErrorIs(t, r1, want, "Get fails with the error of RepoStubFailing")
		})

		t.Run("fails List with its error", func(t *testing.T) {
			t.Parallel()
			want := errors.New("List-failing")
			s := generic.NewRepoStub[T](t, generic.RepoStubFailing[T](want))
			_, r1 := repoStubList[T](t, s)
			assert.ErrorIs(t, r1, want, "List fails with the error of RepoStubFailing")
		})
	})

	t.Run("WithRepoGet", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Get", func(t *testing.T) {
			t.Parallel()
			called := false
			s := generic.NewRepoStub(t, generic.WithRepoGet[T](func(_ string) (T, error) {
				called = true
				var z0 T
				var z1 error
				return z0, z1
			}))
			_, _ = repoStubGet[T](t, s)
			assert.True(t, called, "WithRepoGet sets the function of Get")
		})
	})

	t.Run("WithRepoPut", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Put", func(t *testing.T) {
			t.Parallel()
			called := false
			s := generic.NewRepoStub(t, generic.WithRepoPut[T](func(_ T) {
				called = true
			}))
			repoStubPut[T](t, s)
			assert.True(t, called, "WithRepoPut sets the function of Put")
		})
	})

	t.Run("WithRepoList", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of List", func(t *testing.T) {
			t.Parallel()
			called := false
			s := generic.NewRepoStub(t, generic.WithRepoList[T](func() ([]T, error) {
				called = true
				var z0 []T
				var z1 error
				return z0, z1
			}))
			_, _ = repoStubList[T](t, s)
			assert.True(t, called, "WithRepoList sets the function of List")
		})
	})

	t.Run("WithRepoBoxed", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of Boxed", func(t *testing.T) {
			t.Parallel()
			called := false
			s := generic.NewRepoStub(t, generic.WithRepoBoxed[T](func() generic.Box[T] {
				called = true
				var z0 generic.Box[T]
				return z0
			}))
			_ = repoStubBoxed[T](t, s)
			assert.True(t, called, "WithRepoBoxed sets the function of Boxed")
		})
	})
}

// TestRepoStub runs the checks of RepoStub at the witnesses of its type
// parameters.
func TestRepoStub(t *testing.T) {
	t.Parallel()

	checkRepoStub[int](t)
}

// pairsStubAll calls All of s with the samples of its arguments.
func pairsStubAll[K1 comparable, V_ any](tb assert.TB, s *generic.PairsStub[K1, V_]) map[K1]V_ {
	tb.Helper()
	return s.All()
}

// pairsStubAllSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for All. The subject calls a new double that
// reports to tb.
func pairsStubAllSubject[K1 comparable, V_ any](
	tb assert.TB,
) stub.Subject[generic.PairsAllCall[K1, V_], generic.PairsAllReturn[K1, V_]] {
	tb.Helper()
	s := generic.NewPairsStub[K1, V_](tb)
	return stub.Subject[generic.PairsAllCall[K1, V_], generic.PairsAllReturn[K1, V_]]{
		Stub:  s.OnAll.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = pairsStubAll[K1, V_](tb, s)
		},
		Result: func() generic.PairsAllReturn[K1, V_] {
			got0 := pairsStubAll[K1, V_](tb, s)
			return generic.PairsAllReturn[K1, V_]{Result: got0}
		},
		Override: func(seen func(generic.PairsAllCall[K1, V_])) {
			s.OnAll.Func(func() map[K1]V_ {
				seen(generic.PairsAllCall[K1, V_]{})
				var z0 map[K1]V_
				return z0
			})
		},
	}
}

// pairsStubInstance is the instance of a generic.PairsStub that the checks
// of the constructors and the options drive.
type pairsStubInstance[K1 comparable, V_ any] = stub.Instance[generic.PairsAllCall[K1, V_]]

// pairsStubInstanceOf returns the instance of s. The options apply to the
// whole double, so the checks call All alone.
func pairsStubInstanceOf[K1 comparable, V_ any](
	tb assert.TB,
	s *generic.PairsStub[K1, V_],
) pairsStubInstance[K1, V_] {
	tb.Helper()
	return pairsStubInstance[K1, V_]{
		Stub: s.OnAll.MethodStub,
		Call: func() {
			_ = pairsStubAll[K1, V_](tb, s)
		},
	}
}

// checkPairsStub checks the methods, the constructors and the options of
// the double PairsStub at its type parameters.
func checkPairsStub[K1 comparable, V_ any](t *testing.T) {
	t.Helper()

	t.Run("All", func(t *testing.T) {
		t.Parallel()
		stub.Behaviour(t, "All", pairsStubAllSubject[K1, V_])
	})

	t.Run("ResetCalls", func(t *testing.T) {
		t.Parallel()

		t.Run("All", func(t *testing.T) {
			t.Parallel()
			stub.ResetBehaviour(t, "All", pairsStubAllSubject[K1, V_])
		})
	})

	t.Run("PairsStubStrict", func(t *testing.T) {
		t.Parallel()
		stub.StrictBehaviour(t, func(tb assert.TB) pairsStubInstance[K1, V_] {
			return pairsStubInstanceOf[K1, V_](
				tb,
				generic.NewPairsStub[K1, V_](tb, generic.PairsStubStrict[K1, V_]()),
			)
		})
	})

	t.Run("PairsStubDelegateTo", func(t *testing.T) {
		t.Parallel()

		t.Run("passes All to the implementation", func(t *testing.T) {
			t.Parallel()
			inner := generic.NewPairsStub[K1, V_](t)
			s := generic.NewPairsStub[K1, V_](t, generic.PairsStubDelegateTo[K1, V_](inner))
			_ = pairsStubAll[K1, V_](t, s)
			inner.OnAll.AssertCalledOnce(t, "All calls the implementation")
		})
	})

	t.Run("PairsStubWithClock", func(t *testing.T) {
		t.Parallel()
		stub.ClockBehaviour(t, func(tb assert.TB, clk assert.Clock) pairsStubInstance[K1, V_] {
			return pairsStubInstanceOf[K1, V_](
				tb,
				generic.NewPairsStub[K1, V_](tb, generic.PairsStubWithClock[K1, V_](clk)),
			)
		})
	})

	t.Run("PairsStubWithRandSource", func(t *testing.T) {
		t.Parallel()
		stub.RandSourceBehaviour(t, func(tb assert.TB, src rand.Source) pairsStubInstance[K1, V_] {
			return pairsStubInstanceOf[K1, V_](
				tb,
				generic.NewPairsStub[K1, V_](tb, generic.PairsStubWithRandSource[K1, V_](src)),
			)
		})
	})

	t.Run("PairsStubBenchMode", func(t *testing.T) {
		t.Parallel()
		stub.BenchModeBehaviour(t, func(tb assert.TB) pairsStubInstance[K1, V_] {
			return pairsStubInstanceOf[K1, V_](
				tb,
				generic.NewPairsStub[K1, V_](tb, generic.PairsStubBenchMode[K1, V_]()),
			)
		})
	})

	t.Run("PairsStubWithLatency", func(t *testing.T) {
		t.Parallel()
		stub.LatencyBehaviour(
			t,
			func(tb assert.TB, d time.Duration, clk assert.Clock) pairsStubInstance[K1, V_] {
				return pairsStubInstanceOf[K1, V_](
					tb,
					generic.NewPairsStub[K1, V_](
						tb,
						generic.PairsStubWithClock[K1, V_](clk),
						generic.PairsStubWithLatency[K1, V_](d),
					),
				)
			},
		)
	})

	t.Run("WithPairsAll", func(t *testing.T) {
		t.Parallel()

		t.Run("sets the function of All", func(t *testing.T) {
			t.Parallel()
			called := false
			s := generic.NewPairsStub(t, generic.WithPairsAll[K1, V_](func() map[K1]V_ {
				called = true
				var z0 map[K1]V_
				return z0
			}))
			_ = pairsStubAll[K1, V_](t, s)
			assert.True(t, called, "WithPairsAll sets the function of All")
		})
	})
}

// TestPairsStub runs the checks of PairsStub at the witnesses of its type
// parameters.
func TestPairsStub(t *testing.T) {
	t.Parallel()

	checkPairsStub[string, int](t)
}

// BagStub is generic, and its type parameters have no witnesses that a
// test outside its package can refer to, so the file has no checks of it.
// A test in the package of a consumer instantiates it.
