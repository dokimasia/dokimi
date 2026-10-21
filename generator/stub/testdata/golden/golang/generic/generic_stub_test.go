package generic

import (
	"errors"
	rand "math/rand/v2"
	"testing"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// repoStubGetSubject binds Get into the shape [stub.Behaviour] drives: how
// to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func repoStubGetSubject[T any](tb assert.TB) stub.Subject[RepoGetCall[T], RepoGetReturn[T]] {
	tb.Helper()
	s := NewRepoStub[T](tb)
	return stub.Subject[RepoGetCall[T], RepoGetReturn[T]]{
		Stub: s.OnGet.MethodStub,
		Call: func() {
			var a0 string
			a0 = "Get-a"
			_, _ = s.Get(a0)
		},
		Result: func() RepoGetReturn[T] {
			var a0 string
			a0 = "Get-a"
			got0, got1 := s.Get(a0)
			return RepoGetReturn[T]{Result: got0, Err: got1}
		},
		Override: func(mark func()) {
			s.OnGet.Func(func(_ string) (T, error) {
				mark()
				var z0 T
				var z1 error
				return z0, z1
			})
		},
		Fails: func() error {
			var a0 string
			a0 = "Get-a"
			_, r1 := s.Get(a0)
			return r1
		},
	}
}

// TestRepoStubGet pins how Get answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func repoStubGetChecks[T any](t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Get", repoStubGetSubject[T])

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewRepoStub[T](t)
		var a0 string
		a0 = "Get-a"
		_, _ = s.Get(a0)
		got := s.OnGet.AssertCalledOnce(t, "Get must record the call")
		assert.Equal(t, got.ID, a0, "the recorded call carries ID")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewRepoStub[T](t)
		var seen []RepoGetCall[T]
		s.OnGet.OnRecord(func(c RepoGetCall[T]) { seen = append(seen, c) })
		var a0 string
		a0 = "Get-a"
		_, _ = s.Get(a0)
		_, _ = s.Get(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Get call")
	})

	t.Run("wires WithRepoGet at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewRepoStub(t, WithRepoGet[T](func(_ string) (T, error) {
			called = true
			var z0 T
			var z1 error
			return z0, z1
		}))
		var a0 string
		a0 = "Get-a"
		_, _ = s.Get(a0)
		assert.True(t, called, "WithRepoGet must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("Get-failing")
		s := NewRepoStubFailing[T](t, want)
		var a0 string
		a0 = "Get-a"
		_, r1 := s.Get(a0)
		assert.ErrorIs(t, r1, want, "Get must answer the armed error")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewRepoStub[T](t)
		var want0 T
		var want1 error
		s.OnGet.Returns(want0, want1)
		var a0 string
		a0 = "Get-a"
		_, _ = s.Get(a0)
		s.ResetCalls()
		_, _ = s.Get(a0)
		s.OnGet.AssertCalledOnce(t, "Get must record the call that followed the reset")
	})
}

// repoStubPutSubject binds Put into the shape [stub.Behaviour] drives: how
// to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func repoStubPutSubject[T any](tb assert.TB) stub.Subject[RepoPutCall[T], struct{}] {
	tb.Helper()
	s := NewRepoStub[T](tb)
	return stub.Subject[RepoPutCall[T], struct{}]{
		Stub: s.OnPut.MethodStub,
		Call: func() {
			var a0 T
			s.Put(a0)
		},
		Override: func(mark func()) {
			s.OnPut.Func(func(_ T) {
				mark()
			})
		},
	}
}

// TestRepoStubPut pins how Put answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func repoStubPutChecks[T any](t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Put", repoStubPutSubject[T])

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewRepoStub[T](t)
		var a0 T
		s.Put(a0)
		got := s.OnPut.AssertCalledOnce(t, "Put must record the call")
		assert.Equal(t, got.V, a0, "the recorded call carries V")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewRepoStub[T](t)
		var seen []RepoPutCall[T]
		s.OnPut.OnRecord(func(c RepoPutCall[T]) { seen = append(seen, c) })
		var a0 T
		s.Put(a0)
		s.Put(a0)
		assert.Length(t, seen, 2, "OnRecord must fire once per Put call")
	})

	t.Run("wires WithRepoPut at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewRepoStub(t, WithRepoPut[T](func(_ T) {
			called = true
		}))
		var a0 T
		s.Put(a0)
		assert.True(t, called, "WithRepoPut must install the override")
	})
}

// repoStubListSubject binds List into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func repoStubListSubject[T any](tb assert.TB) stub.Subject[RepoListCall[T], RepoListReturn[T]] {
	tb.Helper()
	s := NewRepoStub[T](tb)
	return stub.Subject[RepoListCall[T], RepoListReturn[T]]{
		Stub: s.OnList.MethodStub,
		Call: func() {
			_, _ = s.List()
		},
		Result: func() RepoListReturn[T] {
			got0, got1 := s.List()
			return RepoListReturn[T]{Result: got0, Err: got1}
		},
		Override: func(mark func()) {
			s.OnList.Func(func() ([]T, error) {
				mark()
				var z0 []T
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

// TestRepoStubList pins how List answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func repoStubListChecks[T any](t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "List", repoStubListSubject[T])

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewRepoStub[T](t)
		_, _ = s.List()
		s.OnList.AssertCalledOnce(t, "List must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewRepoStub[T](t)
		var seen []RepoListCall[T]
		s.OnList.OnRecord(func(c RepoListCall[T]) { seen = append(seen, c) })
		_, _ = s.List()
		_, _ = s.List()
		assert.Length(t, seen, 2, "OnRecord must fire once per List call")
	})

	t.Run("wires WithRepoList at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewRepoStub(t, WithRepoList[T](func() ([]T, error) {
			called = true
			var z0 []T
			var z1 error
			return z0, z1
		}))
		_, _ = s.List()
		assert.True(t, called, "WithRepoList must install the override")
	})

	t.Run("answers the error the Failing preset armed", func(t *testing.T) {
		t.Parallel()
		want := errors.New("List-failing")
		s := NewRepoStubFailing[T](t, want)
		_, r1 := s.List()
		assert.ErrorIs(t, r1, want, "List must answer the armed error")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewRepoStub[T](t)
		var want0 []T
		var want1 error
		s.OnList.Returns(want0, want1)
		_, _ = s.List()
		s.ResetCalls()
		_, _ = s.List()
		s.OnList.AssertCalledOnce(t, "List must record the call that followed the reset")
	})
}

// repoStubBoxedSubject binds Boxed into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func repoStubBoxedSubject[T any](tb assert.TB) stub.Subject[RepoBoxedCall[T], RepoBoxedReturn[T]] {
	tb.Helper()
	s := NewRepoStub[T](tb)
	return stub.Subject[RepoBoxedCall[T], RepoBoxedReturn[T]]{
		Stub: s.OnBoxed.MethodStub,
		Call: func() {
			_ = s.Boxed()
		},
		Result: func() RepoBoxedReturn[T] {
			got0 := s.Boxed()
			return RepoBoxedReturn[T]{Result: got0}
		},
		Override: func(mark func()) {
			s.OnBoxed.Func(func() Box[T] {
				mark()
				var z0 Box[T]
				return z0
			})
		},
	}
}

// TestRepoStubBoxed pins how Boxed answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func repoStubBoxedChecks[T any](t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "Boxed", repoStubBoxedSubject[T])

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewRepoStub[T](t)
		_ = s.Boxed()
		s.OnBoxed.AssertCalledOnce(t, "Boxed must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewRepoStub[T](t)
		var seen []RepoBoxedCall[T]
		s.OnBoxed.OnRecord(func(c RepoBoxedCall[T]) { seen = append(seen, c) })
		_ = s.Boxed()
		_ = s.Boxed()
		assert.Length(t, seen, 2, "OnRecord must fire once per Boxed call")
	})

	t.Run("wires WithRepoBoxed at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewRepoStub(t, WithRepoBoxed[T](func() Box[T] {
			called = true
			var z0 Box[T]
			return z0
		}))
		_ = s.Boxed()
		assert.True(t, called, "WithRepoBoxed must install the override")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewRepoStub[T](t)
		var want0 Box[T]
		s.OnBoxed.Returns(want0)
		_ = s.Boxed()
		s.ResetCalls()
		_ = s.Boxed()
		s.OnBoxed.AssertCalledOnce(t, "Boxed must record the call that followed the reset")
	})
}

// repoStubDouble describes how to build a RepoStub under each option whose
// effect is the same whatever a method's signature.
//
// Get stands in for the double as a whole: these checks assert only that
// an option reached it.
func repoStubDouble[T any]() stub.Double[RepoGetCall[T]] {
	instance := func(s *RepoStub[T]) stub.Instance[RepoGetCall[T]] {
		return stub.Instance[RepoGetCall[T]]{
			Stub: s.OnGet.MethodStub,
			Call: func() {
				var a0 string
				a0 = "Get-a"
				_, _ = s.Get(a0)
			},
			Reset: s.ResetCalls,
		}
	}
	return stub.Double[RepoGetCall[T]]{
		New: func(tb assert.TB) stub.Instance[RepoGetCall[T]] {
			return instance(NewRepoStub[T](tb))
		},
		WithClock: func(tb assert.TB, clk assert.Clock) stub.Instance[RepoGetCall[T]] {
			return instance(NewRepoStub[T](tb, RepoStubWithClock[T](clk)))
		},
		WithLatency: func(tb assert.TB, d time.Duration, clk assert.Clock) stub.Instance[RepoGetCall[T]] {
			return instance(NewRepoStubDelayed[T](tb, d, clk))
		},
		WithRandSource: func(tb assert.TB, src rand.Source) stub.Instance[RepoGetCall[T]] {
			return instance(NewRepoStub[T](tb, RepoStubWithRandSource[T](src)))
		},
		BenchMode: func(tb assert.TB) stub.Instance[RepoGetCall[T]] {
			return instance(NewRepoStub[T](tb, RepoStubBenchMode[T]()))
		},
		Strict: func(tb assert.TB) stub.Instance[RepoGetCall[T]] {
			return instance(NewRepoStub[T](tb, RepoStubStrict[T]()))
		},
	}
}

// TestRepoStubOptions pins the settings that apply to the double rather
// than to one method: bench mode, the injected clock and random source,
// and what a reset clears.
func repoStubOptionChecks[T any](t *testing.T) {
	t.Parallel()

	stub.DoubleBehaviour(t, repoStubDouble[T]())
}

// TestRepoStubDelegateTo pins that a wrapped implementation is called
// through, and recorded on the way past.
func repoStubDelegationChecks[T any](t *testing.T) {
	t.Parallel()

	inner := NewRepoStub[T](t)
	s := NewRepoStub[T](t, RepoStubDelegateTo[T](inner))

	t.Run("forwards Get to the wrapped implementation", func(t *testing.T) {
		var a0 string
		a0 = "Get-a"
		_, _ = s.Get(a0)
		inner.OnGet.AssertCalledOnce(t, "Get must reach the wrapped implementation")
	})

	t.Run("surfaces what Get answered", func(t *testing.T) {
		want := errors.New("Get-delegate")
		inner.OnGet.FaultsFor(time.Hour, want)
		var a0 string
		a0 = "Get-a"
		_, r1 := s.Get(a0)
		assert.ErrorIs(t, r1, want, "Get must surface the wrapped answer")
	})

	t.Run("forwards Put to the wrapped implementation", func(t *testing.T) {
		var a0 T
		s.Put(a0)
		inner.OnPut.AssertCalledOnce(t, "Put must reach the wrapped implementation")
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

	t.Run("forwards Boxed to the wrapped implementation", func(t *testing.T) {
		_ = s.Boxed()
		inner.OnBoxed.AssertCalledOnce(t, "Boxed must reach the wrapped implementation")
	})
}

// TestRepoStubGet runs Get's checks at the witness types its constraints
// admit.
//
// A Go test function cannot take type parameters, so the checks live in a
// generic helper and this instantiates it at one set of witnesses.
func TestRepoStubGet(t *testing.T) { repoStubGetChecks[int](t) }

// TestRepoStubPut runs Put's checks at the witness types its constraints
// admit.
//
// A Go test function cannot take type parameters, so the checks live in a
// generic helper and this instantiates it at one set of witnesses.
func TestRepoStubPut(t *testing.T) { repoStubPutChecks[int](t) }

// TestRepoStubList runs List's checks at the witness types its constraints
// admit.
//
// A Go test function cannot take type parameters, so the checks live in a
// generic helper and this instantiates it at one set of witnesses.
func TestRepoStubList(t *testing.T) { repoStubListChecks[int](t) }

// TestRepoStubBoxed runs Boxed's checks at the witness types its
// constraints admit.
//
// A Go test function cannot take type parameters, so the checks live in a
// generic helper and this instantiates it at one set of witnesses.
func TestRepoStubBoxed(t *testing.T) { repoStubBoxedChecks[int](t) }

// TestRepoStubOptions runs the whole-double checks at the same witnesses.
func TestRepoStubOptions(t *testing.T) { repoStubOptionChecks[int](t) }

// TestRepoStubDelegateTo runs the delegation checks at the same witnesses.
func TestRepoStubDelegateTo(t *testing.T) { repoStubDelegationChecks[int](t) }

// pairsStubAllSubject binds All into the shape [stub.Behaviour] drives:
// how to call it, what it answers with, and how to override it.
//
// Each call builds a fresh double: several checks assert on failure and
// need one bound to a failable TB rather than to the running test.
func pairsStubAllSubject[K1 comparable, V_ any](tb assert.TB) stub.Subject[PairsAllCall[K1, V_], PairsAllReturn[K1, V_]] {
	tb.Helper()
	s := NewPairsStub[K1, V_](tb)
	return stub.Subject[PairsAllCall[K1, V_], PairsAllReturn[K1, V_]]{
		Stub: s.OnAll.MethodStub,
		Call: func() {
			_ = s.All()
		},
		Result: func() PairsAllReturn[K1, V_] {
			got0 := s.All()
			return PairsAllReturn[K1, V_]{Result: got0}
		},
		Override: func(mark func()) {
			s.OnAll.Func(func() map[K1]V_ {
				mark()
				var z0 map[K1]V_
				return z0
			})
		},
	}
}

// TestPairsStubAll pins how All answers.
//
// Recording, resetting, call-count expectations, strict mode, fault
// injection and zero-value dispatch are the same contract for every
// method, so they are asserted once in [stub.Behaviour] rather than
// restated here. What remains below needs a value this method's signature
// can tell apart from a zero one.
func pairsStubAllChecks[K1 comparable, V_ any](t *testing.T) {
	t.Parallel()

	stub.Behaviour(t, "All", pairsStubAllSubject[K1, V_])

	t.Run("records what it was called with", func(t *testing.T) {
		t.Parallel()
		s := NewPairsStub[K1, V_](t)
		_ = s.All()
		s.OnAll.AssertCalledOnce(t, "All must record the call")
	})

	t.Run("fires the OnRecord hook for every call", func(t *testing.T) {
		t.Parallel()
		s := NewPairsStub[K1, V_](t)
		var seen []PairsAllCall[K1, V_]
		s.OnAll.OnRecord(func(c PairsAllCall[K1, V_]) { seen = append(seen, c) })
		_ = s.All()
		_ = s.All()
		assert.Length(t, seen, 2, "OnRecord must fire once per All call")
	})

	t.Run("wires WithPairsAll at construction", func(t *testing.T) {
		t.Parallel()
		called := false
		s := NewPairsStub(t, WithPairsAll[K1, V_](func() map[K1]V_ {
			called = true
			var z0 map[K1]V_
			return z0
		}))
		_ = s.All()
		assert.True(t, called, "WithPairsAll must install the override")
	})

	t.Run("keeps recording across a reset", func(t *testing.T) {
		t.Parallel()
		s := NewPairsStub[K1, V_](t)
		var want0 map[K1]V_
		s.OnAll.Returns(want0)
		_ = s.All()
		s.ResetCalls()
		_ = s.All()
		s.OnAll.AssertCalledOnce(t, "All must record the call that followed the reset")
	})
}

// pairsStubDouble describes how to build a PairsStub under each option
// whose effect is the same whatever a method's signature.
//
// All stands in for the double as a whole: these checks assert only that
// an option reached it.
func pairsStubDouble[K1 comparable, V_ any]() stub.Double[PairsAllCall[K1, V_]] {
	instance := func(s *PairsStub[K1, V_]) stub.Instance[PairsAllCall[K1, V_]] {
		return stub.Instance[PairsAllCall[K1, V_]]{
			Stub: s.OnAll.MethodStub,
			Call: func() {
				_ = s.All()
			},
			Reset: s.ResetCalls,
		}
	}
	return stub.Double[PairsAllCall[K1, V_]]{
		New: func(tb assert.TB) stub.Instance[PairsAllCall[K1, V_]] {
			return instance(NewPairsStub[K1, V_](tb))
		},
		WithClock: func(tb assert.TB, clk assert.Clock) stub.Instance[PairsAllCall[K1, V_]] {
			return instance(NewPairsStub[K1, V_](tb, PairsStubWithClock[K1, V_](clk)))
		},
		WithLatency: func(tb assert.TB, d time.Duration, clk assert.Clock) stub.Instance[PairsAllCall[K1, V_]] {
			return instance(NewPairsStubDelayed[K1, V_](tb, d, clk))
		},
		WithRandSource: func(tb assert.TB, src rand.Source) stub.Instance[PairsAllCall[K1, V_]] {
			return instance(NewPairsStub[K1, V_](tb, PairsStubWithRandSource[K1, V_](src)))
		},
		BenchMode: func(tb assert.TB) stub.Instance[PairsAllCall[K1, V_]] {
			return instance(NewPairsStub[K1, V_](tb, PairsStubBenchMode[K1, V_]()))
		},
		Strict: func(tb assert.TB) stub.Instance[PairsAllCall[K1, V_]] {
			return instance(NewPairsStub[K1, V_](tb, PairsStubStrict[K1, V_]()))
		},
	}
}

// TestPairsStubOptions pins the settings that apply to the double rather
// than to one method: bench mode, the injected clock and random source,
// and what a reset clears.
func pairsStubOptionChecks[K1 comparable, V_ any](t *testing.T) {
	t.Parallel()

	stub.DoubleBehaviour(t, pairsStubDouble[K1, V_]())
}

// TestPairsStubDelegateTo pins that a wrapped implementation is called
// through, and recorded on the way past.
func pairsStubDelegationChecks[K1 comparable, V_ any](t *testing.T) {
	t.Parallel()

	inner := NewPairsStub[K1, V_](t)
	s := NewPairsStub[K1, V_](t, PairsStubDelegateTo[K1, V_](inner))

	t.Run("forwards All to the wrapped implementation", func(t *testing.T) {
		_ = s.All()
		inner.OnAll.AssertCalledOnce(t, "All must reach the wrapped implementation")
	})
}

// TestPairsStubAll runs All's checks at the witness types its constraints
// admit.
//
// A Go test function cannot take type parameters, so the checks live in a
// generic helper and this instantiates it at one set of witnesses.
func TestPairsStubAll(t *testing.T) { pairsStubAllChecks[string, int](t) }

// TestPairsStubOptions runs the whole-double checks at the same witnesses.
func TestPairsStubOptions(t *testing.T) { pairsStubOptionChecks[string, int](t) }

// TestPairsStubDelegateTo runs the delegation checks at the same
// witnesses.
func TestPairsStubDelegateTo(t *testing.T) { pairsStubDelegationChecks[string, int](t) }

// Skipped: BagStub is generic — instantiate it with concrete types and
// write the checks in the consuming package.
