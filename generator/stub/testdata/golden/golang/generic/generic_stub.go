package generic

import (
	rand "math/rand/v2"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// RepoGetCall records one call of Repo.Get.
type RepoGetCall[T any] struct {
	ID     string
	Result T
	Err    error
}

// RepoPutCall records one call of Repo.Put.
type RepoPutCall[T any] struct {
	V T
}

// RepoListCall records one call of Repo.List.
type RepoListCall[T any] struct {
	Result []T
	Err    error
}

// RepoBoxedCall records one call of Repo.Boxed.
type RepoBoxedCall[T any] struct {
	Result Box[T]
}

// RepoGetStub configures Get and records its calls. A fault armed through
// it takes precedence over the function and the fixed results.
type RepoGetStub[T any] struct {
	*stub.MethodStub[RepoGetCall[T]]

	fn       func(string) (T, error)
	fallback *RepoGetReturn[T]
}

// RepoGetReturn contains the fixed results of Get.
type RepoGetReturn[T any] struct {
	Result T
	Err    error
}

// Returns sets the results of every call of Get. A function and an
// injected fault take precedence over them.
func (s *RepoGetStub[T]) Returns(result T, err error) *RepoGetStub[T] {
	s.fallback = &RepoGetReturn[T]{Result: result, Err: err}
	return s
}

// Func sets the function that computes the results of Get from its
// arguments. An injected fault takes precedence over it.
func (s *RepoGetStub[T]) Func(fn func(string) (T, error)) *RepoGetStub[T] {
	s.fn = fn
	return s
}

// RepoPutStub configures Put and records its calls. Put returns no error,
// so a fault armed through it fails the test.
type RepoPutStub[T any] struct {
	*stub.MethodStub[RepoPutCall[T]]

	fn func(T)
}

// Func sets the function that computes the results of Put from its
// arguments.
func (s *RepoPutStub[T]) Func(fn func(T)) *RepoPutStub[T] {
	s.fn = fn
	return s
}

// RepoListStub configures List and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type RepoListStub[T any] struct {
	*stub.MethodStub[RepoListCall[T]]

	fn       func() ([]T, error)
	fallback *RepoListReturn[T]
}

// RepoListReturn contains the fixed results of List.
type RepoListReturn[T any] struct {
	Result []T
	Err    error
}

// Returns sets the results of every call of List. A function and an
// injected fault take precedence over them.
func (s *RepoListStub[T]) Returns(result []T, err error) *RepoListStub[T] {
	s.fallback = &RepoListReturn[T]{Result: result, Err: err}
	return s
}

// Func sets the function that computes the results of List from its
// arguments. An injected fault takes precedence over it.
func (s *RepoListStub[T]) Func(fn func() ([]T, error)) *RepoListStub[T] {
	s.fn = fn
	return s
}

// RepoBoxedStub configures Boxed and records its calls. Boxed returns no
// error, so a fault armed through it fails the test.
type RepoBoxedStub[T any] struct {
	*stub.MethodStub[RepoBoxedCall[T]]

	fn       func() Box[T]
	fallback *RepoBoxedReturn[T]
}

// RepoBoxedReturn contains the fixed results of Boxed.
type RepoBoxedReturn[T any] struct {
	Result Box[T]
}

// Returns sets the results of every call of Boxed. A function takes
// precedence over them.
func (s *RepoBoxedStub[T]) Returns(result Box[T]) *RepoBoxedStub[T] {
	s.fallback = &RepoBoxedReturn[T]{Result: result}
	return s
}

// Func sets the function that computes the results of Boxed from its
// arguments.
func (s *RepoBoxedStub[T]) Func(fn func() Box[T]) *RepoBoxedStub[T] {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// RepoStubOption configures a [RepoStub] at its construction.
type RepoStubOption[T any] func(*RepoStub[T])

// RepoStubStrict returns the option that makes an unconfigured method fail
// the test instead of returning its zero value.
func RepoStubStrict[T any]() RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.strict = true }
}

// RepoStubDelegateTo returns the option that passes each call to impl. The
// double still records the call.
func RepoStubDelegateTo[T any](impl Repo[T]) RepoStubOption[T] {
	return func(s *RepoStub[T]) {
		s.OnGet.Func(impl.Get)
		s.OnPut.Func(impl.Put)
		s.OnList.Func(impl.List)
		s.OnBoxed.Func(impl.Boxed)
	}
}

// RepoStubWithClock returns the option that sets the clock of the double
// to clk. The double takes the latency, the windows of the faults and the
// instants of the calls from clk. A latency on a controlled clock costs no
// real time.
func RepoStubWithClock[T any](clk assert.Clock) RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.group.SetClock(clk) }
}

// RepoStubWithRandSource returns the option that sets the random source of
// the probabilistic faults to src. A seeded source repeats the faults of a
// run.
func RepoStubWithRandSource[T any](src rand.Source) RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.group.SetRandSource(src) }
}

// RepoStubBenchMode returns the option that turns the recording of calls
// off. The methods still return their results.
func RepoStubBenchMode[T any]() RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.group.BenchMode() }
}

// RepoStubWithLatency returns the option that makes every method sleep d
// before it returns. Without RepoStubWithClock, d passes in real time.
func RepoStubWithLatency[T any](d time.Duration) RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.group.SetLatency(d) }
}

// RepoStubFailing returns the option that arms err on each method that can
// fail. Each call of such a method then fails with err.
func RepoStubFailing[T any](err error) RepoStubOption[T] {
	return func(s *RepoStub[T]) {
		s.OnGet.Faults(err, 1)
		s.OnList.Faults(err, 1)
	}
}

// WithRepoGet returns the option that sets the function of Get.
func WithRepoGet[T any](fn func(string) (T, error)) RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.OnGet.Func(fn) }
}

// WithRepoPut returns the option that sets the function of Put.
func WithRepoPut[T any](fn func(T)) RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.OnPut.Func(fn) }
}

// WithRepoList returns the option that sets the function of List.
func WithRepoList[T any](fn func() ([]T, error)) RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.OnList.Func(fn) }
}

// WithRepoBoxed returns the option that sets the function of Boxed.
func WithRepoBoxed[T any](fn func() Box[T]) RepoStubOption[T] {
	return func(s *RepoStub[T]) { s.OnBoxed.Func(fn) }
}

// RepoStub is a recording test double of Repo. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type RepoStub[T any] struct {
	OnGet   *RepoGetStub[T]
	OnPut   *RepoPutStub[T]
	OnList  *RepoListStub[T]
	OnBoxed *RepoBoxedStub[T]

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// The compiler checks that RepoStub implements Repo at the witnesses of
// its type parameters. The double does not branch on its type parameters,
// so it implements Repo at every instance.
var _ Repo[int] = (*RepoStub[int])(nil)

// NewRepoStub returns a double of Repo that reports its failures to tb and
// applies opts. The double checks the expected counts of calls at the end
// of the test. With a nil tb, the double reports nothing.
func NewRepoStub[T any](tb assert.TB, opts ...RepoStubOption[T]) *RepoStub[T] {
	s := &RepoStub[T]{
		OnGet:  &RepoGetStub[T]{MethodStub: stub.NewMethodStub[RepoGetCall[T]](tb, "Repo.Get")},
		OnPut:  &RepoPutStub[T]{MethodStub: stub.NewMethodStub[RepoPutCall[T]](tb, "Repo.Put")},
		OnList: &RepoListStub[T]{MethodStub: stub.NewMethodStub[RepoListCall[T]](tb, "Repo.List")},
		OnBoxed: &RepoBoxedStub[T]{
			MethodStub: stub.NewMethodStub[RepoBoxedCall[T]](tb, "Repo.Boxed"),
		},
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

// NewRepoStubDelayed returns a double whose methods sleep d on clk before
// they return. It applies opts after the clock and the latency.
func NewRepoStubDelayed[T any](
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...RepoStubOption[T],
) *RepoStub[T] {
	return NewRepoStub[T](tb, append([]RepoStubOption[T]{
		RepoStubWithClock[T](clk),
		RepoStubWithLatency[T](d),
	}, opts...)...)
}

// NewRepoStubFailing returns a double on which each call of a method that
// can fail fails with err. It applies opts after the failure.
func NewRepoStubFailing[T any](tb assert.TB, err error, opts ...RepoStubOption[T]) *RepoStub[T] {
	return NewRepoStub[T](tb, append([]RepoStubOption[T]{RepoStubFailing[T](err)}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *RepoStub[T]) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *RepoGetStub[T]) invoke(id string) func() RepoGetReturn[T] {
	if s.fn == nil {
		return nil
	}
	return func() RepoGetReturn[T] {
		r0, r1 := s.fn(id)
		return RepoGetReturn[T]{Result: r0, Err: r1}
	}
}

// Get records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// invokeVoid returns a call of the function of Func with the arguments, and
// nil where no function is set. The empty struct is the result of a method
// without results.
func (s *RepoPutStub[T]) invokeVoid(v T) func() struct{} {
	if s.fn == nil {
		return nil
	}
	return func() struct{} {
		s.fn(v)
		return struct{}{}
	}
}

// Put records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *RepoStub[T]) Put(v T) {
	call := RepoPutCall[T]{V: v}
	stub.Answer(s.OnPut.MethodStub, &call, stub.Arms[RepoPutCall[T], struct{}]{
		Invoke: s.OnPut.invokeVoid(v),
	})
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *RepoListStub[T]) invoke() func() RepoListReturn[T] {
	if s.fn == nil {
		return nil
	}
	return func() RepoListReturn[T] {
		r0, r1 := s.fn()
		return RepoListReturn[T]{Result: r0, Err: r1}
	}
}

// List records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *RepoBoxedStub[T]) invoke() func() RepoBoxedReturn[T] {
	if s.fn == nil {
		return nil
	}
	return func() RepoBoxedReturn[T] {
		r0 := s.fn()
		return RepoBoxedReturn[T]{Result: r0}
	}
}

// Boxed records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// PairsAllCall records one call of Pairs.All.
type PairsAllCall[K1 comparable, V_ any] struct {
	Result map[K1]V_
}

// PairsAllStub configures All and records its calls. All returns no error,
// so a fault armed through it fails the test.
type PairsAllStub[K1 comparable, V_ any] struct {
	*stub.MethodStub[PairsAllCall[K1, V_]]

	fn       func() map[K1]V_
	fallback *PairsAllReturn[K1, V_]
}

// PairsAllReturn contains the fixed results of All.
type PairsAllReturn[K1 comparable, V_ any] struct {
	Result map[K1]V_
}

// Returns sets the results of every call of All. A function takes
// precedence over them.
func (s *PairsAllStub[K1, V_]) Returns(result map[K1]V_) *PairsAllStub[K1, V_] {
	s.fallback = &PairsAllReturn[K1, V_]{Result: result}
	return s
}

// Func sets the function that computes the results of All from its
// arguments.
func (s *PairsAllStub[K1, V_]) Func(fn func() map[K1]V_) *PairsAllStub[K1, V_] {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// PairsStubOption configures a [PairsStub] at its construction.
type PairsStubOption[K1 comparable, V_ any] func(*PairsStub[K1, V_])

// PairsStubStrict returns the option that makes an unconfigured method
// fail the test instead of returning its zero value.
func PairsStubStrict[K1 comparable, V_ any]() PairsStubOption[K1, V_] {
	return func(s *PairsStub[K1, V_]) { s.strict = true }
}

// PairsStubDelegateTo returns the option that passes each call to impl.
// The double still records the call.
func PairsStubDelegateTo[K1 comparable, V_ any](impl Pairs[K1, V_]) PairsStubOption[K1, V_] {
	return func(s *PairsStub[K1, V_]) {
		s.OnAll.Func(impl.All)
	}
}

// PairsStubWithClock returns the option that sets the clock of the double
// to clk. The double takes the latency, the windows of the faults and the
// instants of the calls from clk. A latency on a controlled clock costs no
// real time.
func PairsStubWithClock[K1 comparable, V_ any](clk assert.Clock) PairsStubOption[K1, V_] {
	return func(s *PairsStub[K1, V_]) { s.group.SetClock(clk) }
}

// PairsStubWithRandSource returns the option that sets the random source
// of the probabilistic faults to src. A seeded source repeats the faults
// of a run.
func PairsStubWithRandSource[K1 comparable, V_ any](src rand.Source) PairsStubOption[K1, V_] {
	return func(s *PairsStub[K1, V_]) { s.group.SetRandSource(src) }
}

// PairsStubBenchMode returns the option that turns the recording of calls
// off. The methods still return their results.
func PairsStubBenchMode[K1 comparable, V_ any]() PairsStubOption[K1, V_] {
	return func(s *PairsStub[K1, V_]) { s.group.BenchMode() }
}

// PairsStubWithLatency returns the option that makes every method sleep d
// before it returns. Without PairsStubWithClock, d passes in real time.
func PairsStubWithLatency[K1 comparable, V_ any](d time.Duration) PairsStubOption[K1, V_] {
	return func(s *PairsStub[K1, V_]) { s.group.SetLatency(d) }
}

// WithPairsAll returns the option that sets the function of All.
func WithPairsAll[K1 comparable, V_ any](fn func() map[K1]V_) PairsStubOption[K1, V_] {
	return func(s *PairsStub[K1, V_]) { s.OnAll.Func(fn) }
}

// PairsStub is a recording test double of Pairs. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type PairsStub[K1 comparable, V_ any] struct {
	OnAll *PairsAllStub[K1, V_]

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// The compiler checks that PairsStub implements Pairs at the witnesses of
// its type parameters. The double does not branch on its type parameters,
// so it implements Pairs at every instance.
var _ Pairs[string, int] = (*PairsStub[string, int])(nil)

// NewPairsStub returns a double of Pairs that reports its failures to tb
// and applies opts. The double checks the expected counts of calls at the
// end of the test. With a nil tb, the double reports nothing.
func NewPairsStub[K1 comparable, V_ any](
	tb assert.TB,
	opts ...PairsStubOption[K1, V_],
) *PairsStub[K1, V_] {
	s := &PairsStub[K1, V_]{
		OnAll: &PairsAllStub[K1, V_]{
			MethodStub: stub.NewMethodStub[PairsAllCall[K1, V_]](tb, "Pairs.All"),
		},
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

// NewPairsStubDelayed returns a double whose methods sleep d on clk before
// they return. It applies opts after the clock and the latency.
func NewPairsStubDelayed[K1 comparable, V_ any](
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...PairsStubOption[K1, V_],
) *PairsStub[K1, V_] {
	return NewPairsStub[K1, V_](tb, append([]PairsStubOption[K1, V_]{
		PairsStubWithClock[K1, V_](clk),
		PairsStubWithLatency[K1, V_](d),
	}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *PairsStub[K1, V_]) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *PairsAllStub[K1, V_]) invoke() func() PairsAllReturn[K1, V_] {
	if s.fn == nil {
		return nil
	}
	return func() PairsAllReturn[K1, V_] {
		r0 := s.fn()
		return PairsAllReturn[K1, V_]{Result: r0}
	}
}

// All records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *PairsStub[K1, V_]) All() map[K1]V_ {
	call := PairsAllCall[K1, V_]{}
	r := stub.Answer(
		s.OnAll.MethodStub,
		&call,
		stub.Arms[PairsAllCall[K1, V_], PairsAllReturn[K1, V_]]{
			Invoke:   s.OnAll.invoke(),
			Fallback: s.OnAll.fallback,
			Stamp: func(c *PairsAllCall[K1, V_], r PairsAllReturn[K1, V_]) {
				c.Result = r.Result
			},
		},
	)
	return r.Result
}

// BagAddCall records one call of Bag.Add.
type BagAddCall[K Keyable] struct {
	K K
}

// BagAddStub configures Add and records its calls. Add returns no error,
// so a fault armed through it fails the test.
type BagAddStub[K Keyable] struct {
	*stub.MethodStub[BagAddCall[K]]

	fn func(K)
}

// Func sets the function that computes the results of Add from its
// arguments.
func (s *BagAddStub[K]) Func(fn func(K)) *BagAddStub[K] {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// BagStubOption configures a [BagStub] at its construction.
type BagStubOption[K Keyable] func(*BagStub[K])

// BagStubStrict returns the option that makes an unconfigured method fail
// the test instead of returning its zero value.
func BagStubStrict[K Keyable]() BagStubOption[K] {
	return func(s *BagStub[K]) { s.strict = true }
}

// BagStubDelegateTo returns the option that passes each call to impl. The
// double still records the call.
func BagStubDelegateTo[K Keyable](impl Bag[K]) BagStubOption[K] {
	return func(s *BagStub[K]) {
		s.OnAdd.Func(impl.Add)
	}
}

// BagStubWithClock returns the option that sets the clock of the double to
// clk. The double takes the latency, the windows of the faults and the
// instants of the calls from clk. A latency on a controlled clock costs no
// real time.
func BagStubWithClock[K Keyable](clk assert.Clock) BagStubOption[K] {
	return func(s *BagStub[K]) { s.group.SetClock(clk) }
}

// BagStubWithRandSource returns the option that sets the random source of
// the probabilistic faults to src. A seeded source repeats the faults of a
// run.
func BagStubWithRandSource[K Keyable](src rand.Source) BagStubOption[K] {
	return func(s *BagStub[K]) { s.group.SetRandSource(src) }
}

// BagStubBenchMode returns the option that turns the recording of calls
// off. The methods still return their results.
func BagStubBenchMode[K Keyable]() BagStubOption[K] {
	return func(s *BagStub[K]) { s.group.BenchMode() }
}

// BagStubWithLatency returns the option that makes every method sleep d
// before it returns. Without BagStubWithClock, d passes in real time.
func BagStubWithLatency[K Keyable](d time.Duration) BagStubOption[K] {
	return func(s *BagStub[K]) { s.group.SetLatency(d) }
}

// WithBagAdd returns the option that sets the function of Add.
func WithBagAdd[K Keyable](fn func(K)) BagStubOption[K] {
	return func(s *BagStub[K]) { s.OnAdd.Func(fn) }
}

// BagStub is a recording test double of Bag. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type BagStub[K Keyable] struct {
	OnAdd *BagAddStub[K]

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// NewBagStub returns a double of Bag that reports its failures to tb and
// applies opts. The double checks the expected counts of calls at the end
// of the test. With a nil tb, the double reports nothing.
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

// NewBagStubDelayed returns a double whose methods sleep d on clk before
// they return. It applies opts after the clock and the latency.
func NewBagStubDelayed[K Keyable](
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...BagStubOption[K],
) *BagStub[K] {
	return NewBagStub[K](tb, append([]BagStubOption[K]{
		BagStubWithClock[K](clk),
		BagStubWithLatency[K](d),
	}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *BagStub[K]) ResetCalls() {
	s.group.Reset()
}

// invokeVoid returns a call of the function of Func with the arguments, and
// nil where no function is set. The empty struct is the result of a method
// without results.
func (s *BagAddStub[K]) invokeVoid(k K) func() struct{} {
	if s.fn == nil {
		return nil
	}
	return func() struct{} {
		s.fn(k)
		return struct{}{}
	}
}

// Add records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *BagStub[K]) Add(k K) {
	call := BagAddCall[K]{K: k}
	stub.Answer(s.OnAdd.MethodStub, &call, stub.Arms[BagAddCall[K], struct{}]{
		Invoke: s.OnAdd.invokeVoid(k),
	})
}
