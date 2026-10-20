package generic

import (
	rand "math/rand/v2"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/dokimi-go/stub"
)

// RepoGetCall records one invocation of Repo.Get.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type RepoGetCall[T any] struct {
	ID     string
	Result T
	Err    error
}

// RepoPutCall records one invocation of Repo.Put.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type RepoPutCall[T any] struct {
	V T
}

// RepoListCall records one invocation of Repo.List.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type RepoListCall[T any] struct {
	Result []T
	Err    error
}

// RepoBoxedCall records one invocation of Repo.Boxed.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type RepoBoxedCall[T any] struct {
	Result Box[T]
}

// --- Per-method configuration ---

// RepoGetStub controls how the double answers Get and records what it was
// asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type RepoGetStub[T any] struct {
	*stub.MethodStub[RepoGetCall[T]]

	fn       func(string) (T, error)
	fallback *RepoGetReturn[T]
}

// RepoGetReturn holds the fixed answer configured through Returns.
type RepoGetReturn[T any] struct {
	Result T
	Err    error
}

// Returns pins a fixed result for every call to Get. A Func override takes
// precedence over it, and so does an injected fault.
func (s *RepoGetStub[T]) Returns(result T, err error) *RepoGetStub[T] {
	s.fallback = &RepoGetReturn[T]{Result: result, Err: err}
	return s
}

// Func supplies a body for Get, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *RepoGetStub[T]) Func(fn func(string) (T, error)) *RepoGetStub[T] {
	s.fn = fn
	return s
}

// RepoPutStub controls how the double answers Put and records what it was
// asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type RepoPutStub[T any] struct {
	*stub.MethodStub[RepoPutCall[T]]

	fn func(T)
}

// Func supplies a body for Put, for when the answer depends on the
// arguments. Nothing takes precedence over it: Put answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *RepoPutStub[T]) Func(fn func(T)) *RepoPutStub[T] {
	s.fn = fn
	return s
}

// RepoListStub controls how the double answers List and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type RepoListStub[T any] struct {
	*stub.MethodStub[RepoListCall[T]]

	fn       func() ([]T, error)
	fallback *RepoListReturn[T]
}

// RepoListReturn holds the fixed answer configured through Returns.
type RepoListReturn[T any] struct {
	Result []T
	Err    error
}

// Returns pins a fixed result for every call to List. A Func override
// takes precedence over it, and so does an injected fault.
func (s *RepoListStub[T]) Returns(result []T, err error) *RepoListStub[T] {
	s.fallback = &RepoListReturn[T]{Result: result, Err: err}
	return s
}

// Func supplies a body for List, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *RepoListStub[T]) Func(fn func() ([]T, error)) *RepoListStub[T] {
	s.fn = fn
	return s
}

// RepoBoxedStub controls how the double answers Boxed and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type RepoBoxedStub[T any] struct {
	*stub.MethodStub[RepoBoxedCall[T]]

	fn       func() Box[T]
	fallback *RepoBoxedReturn[T]
}

// RepoBoxedReturn holds the fixed answer configured through Returns.
type RepoBoxedReturn[T any] struct {
	Result Box[T]
}

// Returns pins a fixed result for every call to Boxed. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Boxed answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *RepoBoxedStub[T]) Returns(result Box[T]) *RepoBoxedStub[T] {
	s.fallback = &RepoBoxedReturn[T]{Result: result}
	return s
}

// Func supplies a body for Boxed, for when the answer depends on the
// arguments. Nothing takes precedence over it: Boxed answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *RepoBoxedStub[T]) Func(fn func() Box[T]) *RepoBoxedStub[T] {
	s.fn = fn
	return s
}

// --- RepoStub ---

// The double-runtime witness for the surface this double rides.
var _ = stub.CompatV1

// RepoStubOption configures a [RepoStub] at construction time.
type RepoStubOption[T any] func(*RepoStub[T])

// RepoStubStrict makes an unconfigured method fail the test rather than
// return its zero value.
func RepoStubStrict[T any]() RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.strict = true }
}

// RepoStubDelegateTo wraps a real implementation: every method forwards to
// impl and is recorded on the double on the way through.
func RepoStubDelegateTo[T any](impl Repo[T]) RepoStubOption[T] {
	return func(s *RepoStub[T]) {
		s.OnGet.Func(impl.Get)
		s.OnPut.Func(impl.Put)
		s.OnList.Func(impl.List)
		s.OnBoxed.Func(impl.Boxed)
	}
}

// RepoStubWithClock drives latency and time-windowed faults from clk
// rather than from wall time, so a test asserting on a five-second timeout
// does not take five seconds.
func RepoStubWithClock[T any](clk assert.Clock) RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.group.SetClock(clk) }
}

// RepoStubWithRandSource makes probabilistic fault injection reproducible.
// A failure nobody can replay is not a finding.
func RepoStubWithRandSource[T any](src rand.Source) RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.group.SetRandSource(src) }
}

// RepoStubBenchMode disables call recording. Dispatch still works; only
// the accumulating call log is dropped.
func RepoStubBenchMode[T any]() RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.group.BenchMode() }
}

// RepoStubWithLatency makes every method sleep d before answering, for a
// test standing in for a slow dependency rather than one slow method. Pair
// it with RepoStubWithClock, or d elapses in wall time.
func RepoStubWithLatency[T any](d time.Duration) RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.group.SetLatency(d) }
}

// RepoStubFailing arms every method that can fail with err, so a caller's
// error path is one option rather than one Faults call per method. A
// method answering no error is left alone.
func RepoStubFailing[T any](err error) RepoStubOption[T] {
	return func(s *RepoStub[T]) {
		s.OnGet.Faults(err, 1)
		s.OnList.Faults(err, 1)
	}
}

// WithRepoGet sets Get's body at construction time, for the common case of
// configuring one method and taking the defaults for the rest.
func WithRepoGet[T any](fn func(string) (T, error)) RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.OnGet.Func(fn) }
}

// WithRepoPut sets Put's body at construction time, for the common case of
// configuring one method and taking the defaults for the rest.
func WithRepoPut[T any](fn func(T)) RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.OnPut.Func(fn) }
}

// WithRepoList sets List's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithRepoList[T any](fn func() ([]T, error)) RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.OnList.Func(fn) }
}

// WithRepoBoxed sets Boxed's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithRepoBoxed[T any](fn func() Box[T]) RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.OnBoxed.Func(fn) }
}

// RepoStub is a recording test double for Repo.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method returns its zero value and records the call.
type RepoStub[T any] struct {
	OnGet   *RepoGetStub[T]
	OnPut   *RepoPutStub[T]
	OnList  *RepoListStub[T]
	OnBoxed *RepoBoxedStub[T]

	// group is every method stub above, viewed through the surface that
	// does not depend on a signature.
	group  *stub.Group
	strict bool
}

// Compile-time proof that RepoStub satisfies Repo.
//
// A generic double is asserted at one instantiation. It never branches on
// its type parameters, so satisfying the interface at these types means
// satisfying it at every type the constraints admit.
var _ Repo[int] = (*RepoStub[int])(nil)

// NewRepoStub returns a double bound to tb.
//
// Passing tb registers a cleanup that verifies every method's call-count
// expectations when the test ends, so an unmet Times is reported without
// the caller remembering to check. A nil tb skips that, which is what
// benchmarks and non-test callers want.
func NewRepoStub[T any](tb assert.TB, opts ...RepoStubOption[T]) *RepoStub[T] {
	s := &RepoStub[T]{
		OnGet:   &RepoGetStub[T]{MethodStub: stub.NewMethodStub[RepoGetCall[T]](tb, "Repo.Get")},
		OnPut:   &RepoPutStub[T]{MethodStub: stub.NewMethodStub[RepoPutCall[T]](tb, "Repo.Put")},
		OnList:  &RepoListStub[T]{MethodStub: stub.NewMethodStub[RepoListCall[T]](tb, "Repo.List")},
		OnBoxed: &RepoBoxedStub[T]{MethodStub: stub.NewMethodStub[RepoBoxedCall[T]](tb, "Repo.Boxed")},
	}
	s.group = stub.NewGroup(
		s.OnGet.MethodStub,
		s.OnPut.MethodStub,
		s.OnList.MethodStub,
		s.OnBoxed.MethodStub,
	)
	for _, opt := range opts {
		opt(s)
	}
	if s.strict {
		s.group.StrictAll()
	}
	s.group.Bind(tb)
	return s
}

// NewRepoStubDelayed returns a double whose every method sleeps d against
// clk before answering — the clock and the latency a slow-dependency test
// wants, in one call rather than two options. Further options apply after
// the pair, so a caller can still add Strict or a method override.
func NewRepoStubDelayed[T any](tb assert.TB, d time.Duration, clk assert.Clock, opts ...RepoStubOption[T]) *RepoStub[T] {
	return NewRepoStub[T](tb, append([]RepoStubOption[T]{
		RepoStubWithClock[T](clk),
		RepoStubWithLatency[T](d),
	}, opts...)...)
}

// NewRepoStubFailing returns a double whose every fallible method answers
// err — [RepoStubFailing] applied at construction, for the test whose
// whole subject is the caller's error path. Further options apply after
// it.
func NewRepoStubFailing[T any](tb assert.TB, err error, opts ...RepoStubOption[T]) *RepoStub[T] {
	return NewRepoStub[T](tb, append([]RepoStubOption[T]{RepoStubFailing[T](err)}, opts...)...)
}

// ResetCalls clears recorded calls, fault counters, and call-count
// expectations. Func and Returns configuration is deliberately preserved, so
// a test can reuse one double across phases without rebuilding it.
func (s *RepoStub[T]) ResetCalls() {
	s.group.Reset()
}

// --- Interface methods ---

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *RepoGetStub[T]) invoke(id string) func() RepoGetReturn[T] {
	if s.fn == nil {
		return nil
	}
	return func() RepoGetReturn[T] {
		r0, r1 := s.fn(id)
		return RepoGetReturn[T]{Result: r0, Err: r1}
	}
}

// Get records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *RepoStub[T]) Get(id string) (T, error) {
	call := RepoGetCall[T]{ID: id}
	r := stub.Answer(s.OnGet.MethodStub, &call, stub.Arms[RepoGetCall[T], RepoGetReturn[T]]{
		Invoke:   s.OnGet.invoke(id),
		Fallback: s.OnGet.fallback,
		Fault:    func(err error) RepoGetReturn[T] { return RepoGetReturn[T]{Err: err} },
		Stamp: func(c *RepoGetCall[T], r RepoGetReturn[T]) {
			c.Result = r.Result
			c.Err = r.Err
		},
	})
	return r.Result, r.Err
}

// invokeVoid adapts the Func override to the shape [stub.Answer] consumes.
// A method returning nothing still has an override worth dispatching to and
// a call worth recording; the empty struct is the return tuple of no returns.
func (s *RepoPutStub[T]) invokeVoid(v T) func() struct{} {
	if s.fn == nil {
		return nil
	}
	return func() struct{} {
		s.fn(v)
		return struct{}{}
	}
}

// Put records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *RepoStub[T]) Put(v T) {
	call := RepoPutCall[T]{V: v}
	stub.Answer(s.OnPut.MethodStub, &call, stub.Arms[RepoPutCall[T], struct{}]{
		Invoke: s.OnPut.invokeVoid(v),
	})
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *RepoListStub[T]) invoke() func() RepoListReturn[T] {
	if s.fn == nil {
		return nil
	}
	return func() RepoListReturn[T] {
		r0, r1 := s.fn()
		return RepoListReturn[T]{Result: r0, Err: r1}
	}
}

// List records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *RepoStub[T]) List() ([]T, error) {
	call := RepoListCall[T]{}
	r := stub.Answer(s.OnList.MethodStub, &call, stub.Arms[RepoListCall[T], RepoListReturn[T]]{
		Invoke:   s.OnList.invoke(),
		Fallback: s.OnList.fallback,
		Fault:    func(err error) RepoListReturn[T] { return RepoListReturn[T]{Err: err} },
		Stamp: func(c *RepoListCall[T], r RepoListReturn[T]) {
			c.Result = r.Result
			c.Err = r.Err
		},
	})
	return r.Result, r.Err
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *RepoBoxedStub[T]) invoke() func() RepoBoxedReturn[T] {
	if s.fn == nil {
		return nil
	}
	return func() RepoBoxedReturn[T] {
		r0 := s.fn()
		return RepoBoxedReturn[T]{Result: r0}
	}
}

// Boxed records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *RepoStub[T]) Boxed() Box[T] {
	call := RepoBoxedCall[T]{}
	r := stub.Answer(s.OnBoxed.MethodStub, &call, stub.Arms[RepoBoxedCall[T], RepoBoxedReturn[T]]{
		Invoke:   s.OnBoxed.invoke(),
		Fallback: s.OnBoxed.fallback,
		Stamp: func(c *RepoBoxedCall[T], r RepoBoxedReturn[T]) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// PairsAllCall records one invocation of Pairs.All.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type PairsAllCall[K1 comparable, V_ any] struct {
	Result map[K1]V_
}

// --- Per-method configuration ---

// PairsAllStub controls how the double answers All and records what it was
// asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type PairsAllStub[K1 comparable, V_ any] struct {
	*stub.MethodStub[PairsAllCall[K1, V_]]

	fn       func() map[K1]V_
	fallback *PairsAllReturn[K1, V_]
}

// PairsAllReturn holds the fixed answer configured through Returns.
type PairsAllReturn[K1 comparable, V_ any] struct {
	Result map[K1]V_
}

// Returns pins a fixed result for every call to All. A Func override takes
// precedence over it.
//
// No fault can take precedence here: All answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *PairsAllStub[K1, V_]) Returns(result map[K1]V_) *PairsAllStub[K1, V_] {
	s.fallback = &PairsAllReturn[K1, V_]{Result: result}
	return s
}

// Func supplies a body for All, for when the answer depends on the
// arguments. Nothing takes precedence over it: All answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *PairsAllStub[K1, V_]) Func(fn func() map[K1]V_) *PairsAllStub[K1, V_] {
	s.fn = fn
	return s
}

// --- PairsStub ---

// The double-runtime witness for the surface this double rides.
var _ = stub.CompatV1

// PairsStubOption configures a [PairsStub] at construction time.
type PairsStubOption[K1 comparable, V_ any] func(*PairsStub[K1, V_])

// PairsStubStrict makes an unconfigured method fail the test rather than
// return its zero value.
func PairsStubStrict[K1 comparable, V_ any]() PairsStubOption[K1, V_] {
	return func(s *PairsStub[K1, V_]) { s.strict = true }
}

// PairsStubDelegateTo wraps a real implementation: every method forwards
// to impl and is recorded on the double on the way through.
func PairsStubDelegateTo[K1 comparable, V_ any](impl Pairs[K1, V_]) PairsStubOption[K1, V_] {
	return func(s *PairsStub[K1, V_]) {
		s.OnAll.Func(impl.All)
	}
}

// PairsStubWithClock drives latency and time-windowed faults from clk
// rather than from wall time, so a test asserting on a five-second timeout
// does not take five seconds.
func PairsStubWithClock[K1 comparable, V_ any](clk assert.Clock) PairsStubOption[K1, V_] {
	return func(s *PairsStub[K1, V_]) { s.group.SetClock(clk) }
}

// PairsStubWithRandSource makes probabilistic fault injection
// reproducible. A failure nobody can replay is not a finding.
func PairsStubWithRandSource[K1 comparable, V_ any](src rand.Source) PairsStubOption[K1, V_] {
	return func(s *PairsStub[K1, V_]) { s.group.SetRandSource(src) }
}

// PairsStubBenchMode disables call recording. Dispatch still works; only
// the accumulating call log is dropped.
func PairsStubBenchMode[K1 comparable, V_ any]() PairsStubOption[K1, V_] {
	return func(s *PairsStub[K1, V_]) { s.group.BenchMode() }
}

// PairsStubWithLatency makes every method sleep d before answering, for a
// test standing in for a slow dependency rather than one slow method. Pair
// it with PairsStubWithClock, or d elapses in wall time.
func PairsStubWithLatency[K1 comparable, V_ any](d time.Duration) PairsStubOption[K1, V_] {
	return func(s *PairsStub[K1, V_]) { s.group.SetLatency(d) }
}

// WithPairsAll sets All's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithPairsAll[K1 comparable, V_ any](fn func() map[K1]V_) PairsStubOption[K1, V_] {
	return func(s *PairsStub[K1, V_]) { s.OnAll.Func(fn) }
}

// PairsStub is a recording test double for Pairs.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method returns its zero value and records the call.
type PairsStub[K1 comparable, V_ any] struct {
	OnAll *PairsAllStub[K1, V_]

	// group is every method stub above, viewed through the surface that
	// does not depend on a signature.
	group  *stub.Group
	strict bool
}

// Compile-time proof that PairsStub satisfies Pairs.
//
// A generic double is asserted at one instantiation. It never branches on
// its type parameters, so satisfying the interface at these types means
// satisfying it at every type the constraints admit.
var _ Pairs[string, int] = (*PairsStub[string, int])(nil)

// NewPairsStub returns a double bound to tb.
//
// Passing tb registers a cleanup that verifies every method's call-count
// expectations when the test ends, so an unmet Times is reported without
// the caller remembering to check. A nil tb skips that, which is what
// benchmarks and non-test callers want.
func NewPairsStub[K1 comparable, V_ any](tb assert.TB, opts ...PairsStubOption[K1, V_]) *PairsStub[K1, V_] {
	s := &PairsStub[K1, V_]{
		OnAll: &PairsAllStub[K1, V_]{MethodStub: stub.NewMethodStub[PairsAllCall[K1, V_]](tb, "Pairs.All")},
	}
	s.group = stub.NewGroup(
		s.OnAll.MethodStub,
	)
	for _, opt := range opts {
		opt(s)
	}
	if s.strict {
		s.group.StrictAll()
	}
	s.group.Bind(tb)
	return s
}

// NewPairsStubDelayed returns a double whose every method sleeps d against
// clk before answering — the clock and the latency a slow-dependency test
// wants, in one call rather than two options. Further options apply after
// the pair, so a caller can still add Strict or a method override.
func NewPairsStubDelayed[K1 comparable, V_ any](tb assert.TB, d time.Duration, clk assert.Clock, opts ...PairsStubOption[K1, V_]) *PairsStub[K1, V_] {
	return NewPairsStub[K1, V_](tb, append([]PairsStubOption[K1, V_]{
		PairsStubWithClock[K1, V_](clk),
		PairsStubWithLatency[K1, V_](d),
	}, opts...)...)
}

// ResetCalls clears recorded calls, fault counters, and call-count
// expectations. Func and Returns configuration is deliberately preserved, so
// a test can reuse one double across phases without rebuilding it.
func (s *PairsStub[K1, V_]) ResetCalls() {
	s.group.Reset()
}

// --- Interface methods ---

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *PairsAllStub[K1, V_]) invoke() func() PairsAllReturn[K1, V_] {
	if s.fn == nil {
		return nil
	}
	return func() PairsAllReturn[K1, V_] {
		r0 := s.fn()
		return PairsAllReturn[K1, V_]{Result: r0}
	}
}

// All records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *PairsStub[K1, V_]) All() map[K1]V_ {
	call := PairsAllCall[K1, V_]{}
	r := stub.Answer(s.OnAll.MethodStub, &call, stub.Arms[PairsAllCall[K1, V_], PairsAllReturn[K1, V_]]{
		Invoke:   s.OnAll.invoke(),
		Fallback: s.OnAll.fallback,
		Stamp: func(c *PairsAllCall[K1, V_], r PairsAllReturn[K1, V_]) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// BagAddCall records one invocation of Bag.Add.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type BagAddCall[K Keyable] struct {
	K K
}

// --- Per-method configuration ---

// BagAddStub controls how the double answers Add and records what it was
// asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type BagAddStub[K Keyable] struct {
	*stub.MethodStub[BagAddCall[K]]

	fn func(K)
}

// Func supplies a body for Add, for when the answer depends on the
// arguments. Nothing takes precedence over it: Add answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *BagAddStub[K]) Func(fn func(K)) *BagAddStub[K] {
	s.fn = fn
	return s
}

// --- BagStub ---

// The double-runtime witness for the surface this double rides.
var _ = stub.CompatV1

// BagStubOption configures a [BagStub] at construction time.
type BagStubOption[K Keyable] func(*BagStub[K])

// BagStubStrict makes an unconfigured method fail the test rather than
// return its zero value.
func BagStubStrict[K Keyable]() BagStubOption[K] {
	return func(s *BagStub[K]) { s.strict = true }
}

// BagStubDelegateTo wraps a real implementation: every method forwards to
// impl and is recorded on the double on the way through.
func BagStubDelegateTo[K Keyable](impl Bag[K]) BagStubOption[K] {
	return func(s *BagStub[K]) {
		s.OnAdd.Func(impl.Add)
	}
}

// BagStubWithClock drives latency and time-windowed faults from clk rather
// than from wall time, so a test asserting on a five-second timeout does
// not take five seconds.
func BagStubWithClock[K Keyable](clk assert.Clock) BagStubOption[K] {
	return func(s *BagStub[K]) { s.group.SetClock(clk) }
}

// BagStubWithRandSource makes probabilistic fault injection reproducible.
// A failure nobody can replay is not a finding.
func BagStubWithRandSource[K Keyable](src rand.Source) BagStubOption[K] {
	return func(s *BagStub[K]) { s.group.SetRandSource(src) }
}

// BagStubBenchMode disables call recording. Dispatch still works; only the
// accumulating call log is dropped.
func BagStubBenchMode[K Keyable]() BagStubOption[K] {
	return func(s *BagStub[K]) { s.group.BenchMode() }
}

// BagStubWithLatency makes every method sleep d before answering, for a
// test standing in for a slow dependency rather than one slow method. Pair
// it with BagStubWithClock, or d elapses in wall time.
func BagStubWithLatency[K Keyable](d time.Duration) BagStubOption[K] {
	return func(s *BagStub[K]) { s.group.SetLatency(d) }
}

// WithBagAdd sets Add's body at construction time, for the common case of
// configuring one method and taking the defaults for the rest.
func WithBagAdd[K Keyable](fn func(K)) BagStubOption[K] {
	return func(s *BagStub[K]) { s.OnAdd.Func(fn) }
}

// BagStub is a recording test double for Bag.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method returns its zero value and records the call.
type BagStub[K Keyable] struct {
	OnAdd *BagAddStub[K]

	// group is every method stub above, viewed through the surface that
	// does not depend on a signature.
	group  *stub.Group
	strict bool
}

// NewBagStub returns a double bound to tb.
//
// Passing tb registers a cleanup that verifies every method's call-count
// expectations when the test ends, so an unmet Times is reported without
// the caller remembering to check. A nil tb skips that, which is what
// benchmarks and non-test callers want.
func NewBagStub[K Keyable](tb assert.TB, opts ...BagStubOption[K]) *BagStub[K] {
	s := &BagStub[K]{
		OnAdd: &BagAddStub[K]{MethodStub: stub.NewMethodStub[BagAddCall[K]](tb, "Bag.Add")},
	}
	s.group = stub.NewGroup(
		s.OnAdd.MethodStub,
	)
	for _, opt := range opts {
		opt(s)
	}
	if s.strict {
		s.group.StrictAll()
	}
	s.group.Bind(tb)
	return s
}

// NewBagStubDelayed returns a double whose every method sleeps d against
// clk before answering — the clock and the latency a slow-dependency test
// wants, in one call rather than two options. Further options apply after
// the pair, so a caller can still add Strict or a method override.
func NewBagStubDelayed[K Keyable](tb assert.TB, d time.Duration, clk assert.Clock, opts ...BagStubOption[K]) *BagStub[K] {
	return NewBagStub[K](tb, append([]BagStubOption[K]{
		BagStubWithClock[K](clk),
		BagStubWithLatency[K](d),
	}, opts...)...)
}

// ResetCalls clears recorded calls, fault counters, and call-count
// expectations. Func and Returns configuration is deliberately preserved, so
// a test can reuse one double across phases without rebuilding it.
func (s *BagStub[K]) ResetCalls() {
	s.group.Reset()
}

// --- Interface methods ---

// invokeVoid adapts the Func override to the shape [stub.Answer] consumes.
// A method returning nothing still has an override worth dispatching to and
// a call worth recording; the empty struct is the return tuple of no returns.
func (s *BagAddStub[K]) invokeVoid(k K) func() struct{} {
	if s.fn == nil {
		return nil
	}
	return func() struct{} {
		s.fn(k)
		return struct{}{}
	}
}

// Add records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *BagStub[K]) Add(k K) {
	call := BagAddCall[K]{K: k}
	stub.Answer(s.OnAdd.MethodStub, &call, stub.Arms[BagAddCall[K], struct{}]{
		Invoke: s.OnAdd.invokeVoid(k),
	})
}
