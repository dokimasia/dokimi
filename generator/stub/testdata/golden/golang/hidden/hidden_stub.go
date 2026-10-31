package hidden

import (
	"context"
	rand "math/rand/v2"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// HostsFindCall records one call of Hosts.Find.
type HostsFindCall struct {
	Ctx    context.Context
	Name   string
	Result Host
	Err    error
}

// HostsFindStub configures Find and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type HostsFindStub struct {
	*stub.MethodStub[HostsFindCall]

	fn       func(context.Context, string) (Host, error)
	fallback *HostsFindReturn
}

// HostsFindReturn contains the fixed results of Find.
type HostsFindReturn struct {
	Result Host
	Err    error
}

// Returns sets the results of every call of Find. A function and an
// injected fault take precedence over them.
func (s *HostsFindStub) Returns(result Host, err error) *HostsFindStub {
	s.fallback = &HostsFindReturn{Result: result, Err: err}
	return s
}

// Func sets the function that computes the results of Find from its
// arguments. An injected fault takes precedence over it.
func (s *HostsFindStub) Func(fn func(context.Context, string) (Host, error)) *HostsFindStub {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// HostsStubOption configures a [HostsStub] at its construction.
type HostsStubOption func(*HostsStub)

// HostsStubStrict returns the option that makes an unconfigured method
// fail the test instead of returning its zero value.
func HostsStubStrict() HostsStubOption {
	return func(s *HostsStub) { s.strict = true }
}

// HostsStubDelegateTo returns the option that passes each call to impl.
// The double still records the call.
func HostsStubDelegateTo(impl Hosts) HostsStubOption {
	return func(s *HostsStub) {
		s.OnFind.Func(impl.Find)
	}
}

// HostsStubWithClock returns the option that sets the clock of the double
// to clk. The double takes the latency, the windows of the faults and the
// instants of the calls from clk. A latency on a controlled clock costs no
// real time.
func HostsStubWithClock(clk assert.Clock) HostsStubOption {
	return func(s *HostsStub) { s.group.SetClock(clk) }
}

// HostsStubWithRandSource returns the option that sets the random source
// of the probabilistic faults to src. A seeded source repeats the faults
// of a run.
func HostsStubWithRandSource(src rand.Source) HostsStubOption {
	return func(s *HostsStub) { s.group.SetRandSource(src) }
}

// HostsStubBenchMode returns the option that turns the recording of calls
// off. The methods still return their results.
func HostsStubBenchMode() HostsStubOption {
	return func(s *HostsStub) { s.group.BenchMode() }
}

// HostsStubWithLatency returns the option that makes every method sleep d
// before it returns. Without HostsStubWithClock, d passes in real time.
func HostsStubWithLatency(d time.Duration) HostsStubOption {
	return func(s *HostsStub) { s.group.SetLatency(d) }
}

// HostsStubFailing returns the option that arms err on each method that
// can fail. Each call of such a method then fails with err.
func HostsStubFailing(err error) HostsStubOption {
	return func(s *HostsStub) {
		s.OnFind.Faults(err, 1)
	}
}

// WithHostsFind returns the option that sets the function of Find.
func WithHostsFind(fn func(context.Context, string) (Host, error)) HostsStubOption {
	return func(s *HostsStub) { s.OnFind.Func(fn) }
}

// HostsStub is a recording test double of Hosts. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type HostsStub struct {
	OnFind *HostsFindStub

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// The compiler checks that HostsStub implements Hosts.
var _ Hosts = (*HostsStub)(nil)

// NewHostsStub returns a double of Hosts that reports its failures to tb
// and applies opts. The double checks the expected counts of calls at the
// end of the test. With a nil tb, the double reports nothing.
func NewHostsStub(tb assert.TB, opts ...HostsStubOption) *HostsStub {
	s := &HostsStub{
		OnFind: &HostsFindStub{MethodStub: stub.NewMethodStub[HostsFindCall](tb, "Hosts.Find")},
	}
	s.group = stub.NewGroup(
		s.OnFind.MethodStub,
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

// NewHostsStubDelayed returns a double whose methods sleep d on clk before
// they return. It applies opts after the clock and the latency.
func NewHostsStubDelayed(
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...HostsStubOption,
) *HostsStub {
	return NewHostsStub(tb, append([]HostsStubOption{
		HostsStubWithClock(clk),
		HostsStubWithLatency(d),
	}, opts...)...)
}

// NewHostsStubFailing returns a double on which each call of a method that
// can fail fails with err. It applies opts after the failure.
func NewHostsStubFailing(tb assert.TB, err error, opts ...HostsStubOption) *HostsStub {
	return NewHostsStub(tb, append([]HostsStubOption{HostsStubFailing(err)}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *HostsStub) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *HostsFindStub) invoke(ctx context.Context, name string) func() HostsFindReturn {
	if s.fn == nil {
		return nil
	}
	return func() HostsFindReturn {
		r0, r1 := s.fn(ctx, name)
		return HostsFindReturn{Result: r0, Err: r1}
	}
}

// Find records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *HostsStub) Find(ctx context.Context, name string) (Host, error) {
	call := HostsFindCall{Ctx: ctx, Name: name}
	r := stub.Answer(s.OnFind.MethodStub, call, stub.Arms[HostsFindCall, HostsFindReturn]{
		Invoke:   s.OnFind.invoke(ctx, name),
		Fallback: s.OnFind.fallback,
		Fault:    func(err error) HostsFindReturn { return HostsFindReturn{Err: err} },
		Stamp: func(c HostsFindCall, r HostsFindReturn) HostsFindCall {
			c.Result = r.Result
			c.Err = r.Err
			return c
		},
	})
	return r.Result, r.Err
}

// storeGetCall records one call of store.Get.
type storeGetCall struct {
	Ctx    context.Context
	Key    string
	Result string
	Err    error
}

// storeOpenCall records one call of store.Open.
type storeOpenCall struct {
	Ctx    context.Context
	Result session
	Err    error
}

// storeGetStub configures Get and records its calls. A fault armed through
// it takes precedence over the function and the fixed results.
type storeGetStub struct {
	*stub.MethodStub[storeGetCall]

	fn       func(context.Context, string) (string, error)
	fallback *storeGetReturn
}

// storeGetReturn contains the fixed results of Get.
type storeGetReturn struct {
	Result string
	Err    error
}

// Returns sets the results of every call of Get. A function and an
// injected fault take precedence over them.
func (s *storeGetStub) Returns(result string, err error) *storeGetStub {
	s.fallback = &storeGetReturn{Result: result, Err: err}
	return s
}

// Func sets the function that computes the results of Get from its
// arguments. An injected fault takes precedence over it.
func (s *storeGetStub) Func(fn func(context.Context, string) (string, error)) *storeGetStub {
	s.fn = fn
	return s
}

// storeOpenStub configures Open and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type storeOpenStub struct {
	*stub.MethodStub[storeOpenCall]

	fn       func(context.Context) (session, error)
	fallback *storeOpenReturn
}

// storeOpenReturn contains the fixed results of Open.
type storeOpenReturn struct {
	Result session
	Err    error
}

// Returns sets the results of every call of Open. A function and an
// injected fault take precedence over them.
func (s *storeOpenStub) Returns(result session, err error) *storeOpenStub {
	s.fallback = &storeOpenReturn{Result: result, Err: err}
	return s
}

// Func sets the function that computes the results of Open from its
// arguments. An injected fault takes precedence over it.
func (s *storeOpenStub) Func(fn func(context.Context) (session, error)) *storeOpenStub {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// storeStubOption configures a [storeStub] at its construction.
type storeStubOption func(*storeStub)

// storeStubStrict returns the option that makes an unconfigured method
// fail the test instead of returning its zero value.
func storeStubStrict() storeStubOption {
	return func(s *storeStub) { s.strict = true }
}

// storeStubDelegateTo returns the option that passes each call to impl.
// The double still records the call.
func storeStubDelegateTo(impl store) storeStubOption {
	return func(s *storeStub) {
		s.OnGet.Func(impl.Get)
		s.OnOpen.Func(impl.Open)
	}
}

// storeStubWithClock returns the option that sets the clock of the double
// to clk. The double takes the latency, the windows of the faults and the
// instants of the calls from clk. A latency on a controlled clock costs no
// real time.
func storeStubWithClock(clk assert.Clock) storeStubOption {
	return func(s *storeStub) { s.group.SetClock(clk) }
}

// storeStubWithRandSource returns the option that sets the random source
// of the probabilistic faults to src. A seeded source repeats the faults
// of a run.
func storeStubWithRandSource(src rand.Source) storeStubOption {
	return func(s *storeStub) { s.group.SetRandSource(src) }
}

// storeStubBenchMode returns the option that turns the recording of calls
// off. The methods still return their results.
func storeStubBenchMode() storeStubOption {
	return func(s *storeStub) { s.group.BenchMode() }
}

// storeStubWithLatency returns the option that makes every method sleep d
// before it returns. Without storeStubWithClock, d passes in real time.
func storeStubWithLatency(d time.Duration) storeStubOption {
	return func(s *storeStub) { s.group.SetLatency(d) }
}

// storeStubFailing returns the option that arms err on each method that
// can fail. Each call of such a method then fails with err.
func storeStubFailing(err error) storeStubOption {
	return func(s *storeStub) {
		s.OnGet.Faults(err, 1)
		s.OnOpen.Faults(err, 1)
	}
}

// withStoreGet returns the option that sets the function of Get.
func withStoreGet(fn func(context.Context, string) (string, error)) storeStubOption {
	return func(s *storeStub) { s.OnGet.Func(fn) }
}

// withStoreOpen returns the option that sets the function of Open.
func withStoreOpen(fn func(context.Context) (session, error)) storeStubOption {
	return func(s *storeStub) { s.OnOpen.Func(fn) }
}

// storeStub is a recording test double of store. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type storeStub struct {
	OnGet  *storeGetStub
	OnOpen *storeOpenStub

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// The compiler checks that storeStub implements store.
var _ store = (*storeStub)(nil)

// newStoreStub returns a double of store that reports its failures to tb
// and applies opts. The double checks the expected counts of calls at the
// end of the test. With a nil tb, the double reports nothing.
func newStoreStub(tb assert.TB, opts ...storeStubOption) *storeStub {
	s := &storeStub{
		OnGet:  &storeGetStub{MethodStub: stub.NewMethodStub[storeGetCall](tb, "store.Get")},
		OnOpen: &storeOpenStub{MethodStub: stub.NewMethodStub[storeOpenCall](tb, "store.Open")},
	}
	s.group = stub.NewGroup(
		s.OnGet.MethodStub,
		s.OnOpen.MethodStub,
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

// newStoreStubDelayed returns a double whose methods sleep d on clk before
// they return. It applies opts after the clock and the latency.
func newStoreStubDelayed(
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...storeStubOption,
) *storeStub {
	return newStoreStub(tb, append([]storeStubOption{
		storeStubWithClock(clk),
		storeStubWithLatency(d),
	}, opts...)...)
}

// newStoreStubFailing returns a double on which each call of a method that
// can fail fails with err. It applies opts after the failure.
func newStoreStubFailing(tb assert.TB, err error, opts ...storeStubOption) *storeStub {
	return newStoreStub(tb, append([]storeStubOption{storeStubFailing(err)}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *storeStub) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *storeGetStub) invoke(ctx context.Context, key string) func() storeGetReturn {
	if s.fn == nil {
		return nil
	}
	return func() storeGetReturn {
		r0, r1 := s.fn(ctx, key)
		return storeGetReturn{Result: r0, Err: r1}
	}
}

// Get records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *storeStub) Get(ctx context.Context, key string) (string, error) {
	call := storeGetCall{Ctx: ctx, Key: key}
	r := stub.Answer(s.OnGet.MethodStub, call, stub.Arms[storeGetCall, storeGetReturn]{
		Invoke:   s.OnGet.invoke(ctx, key),
		Fallback: s.OnGet.fallback,
		Fault:    func(err error) storeGetReturn { return storeGetReturn{Err: err} },
		Stamp: func(c storeGetCall, r storeGetReturn) storeGetCall {
			c.Result = r.Result
			c.Err = r.Err
			return c
		},
	})
	return r.Result, r.Err
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *storeOpenStub) invoke(ctx context.Context) func() storeOpenReturn {
	if s.fn == nil {
		return nil
	}
	return func() storeOpenReturn {
		r0, r1 := s.fn(ctx)
		return storeOpenReturn{Result: r0, Err: r1}
	}
}

// Open records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *storeStub) Open(ctx context.Context) (session, error) {
	call := storeOpenCall{Ctx: ctx}
	r := stub.Answer(s.OnOpen.MethodStub, call, stub.Arms[storeOpenCall, storeOpenReturn]{
		Invoke:   s.OnOpen.invoke(ctx),
		Fallback: s.OnOpen.fallback,
		Fault:    func(err error) storeOpenReturn { return storeOpenReturn{Err: err} },
		Stamp: func(c storeOpenCall, r storeOpenReturn) storeOpenCall {
			c.Result = r.Result
			c.Err = r.Err
			return c
		},
	})
	return r.Result, r.Err
}

// session is a handle that a method of store returns.

// sessionCloseCall records one call of session.Close.
type sessionCloseCall struct {
	Err error
}

// sessionCloseStub configures Close and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type sessionCloseStub struct {
	*stub.MethodStub[sessionCloseCall]

	fn       func() error
	fallback *sessionCloseReturn
}

// sessionCloseReturn contains the fixed results of Close.
type sessionCloseReturn struct {
	Err error
}

// Returns sets the results of every call of Close. A function and an
// injected fault take precedence over them.
func (s *sessionCloseStub) Returns(err error) *sessionCloseStub {
	s.fallback = &sessionCloseReturn{Err: err}
	return s
}

// Func sets the function that computes the results of Close from its
// arguments. An injected fault takes precedence over it.
func (s *sessionCloseStub) Func(fn func() error) *sessionCloseStub {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// sessionStubOption configures a [sessionStub] at its construction.
type sessionStubOption func(*sessionStub)

// sessionStubStrict returns the option that makes an unconfigured method
// fail the test instead of returning its zero value.
func sessionStubStrict() sessionStubOption {
	return func(s *sessionStub) { s.strict = true }
}

// sessionStubDelegateTo returns the option that passes each call to impl.
// The double still records the call.
func sessionStubDelegateTo(impl session) sessionStubOption {
	return func(s *sessionStub) {
		s.OnClose.Func(impl.Close)
	}
}

// sessionStubWithClock returns the option that sets the clock of the
// double to clk. The double takes the latency, the windows of the faults
// and the instants of the calls from clk. A latency on a controlled clock
// costs no real time.
func sessionStubWithClock(clk assert.Clock) sessionStubOption {
	return func(s *sessionStub) { s.group.SetClock(clk) }
}

// sessionStubWithRandSource returns the option that sets the random source
// of the probabilistic faults to src. A seeded source repeats the faults
// of a run.
func sessionStubWithRandSource(src rand.Source) sessionStubOption {
	return func(s *sessionStub) { s.group.SetRandSource(src) }
}

// sessionStubBenchMode returns the option that turns the recording of
// calls off. The methods still return their results.
func sessionStubBenchMode() sessionStubOption {
	return func(s *sessionStub) { s.group.BenchMode() }
}

// sessionStubWithLatency returns the option that makes every method sleep
// d before it returns. Without sessionStubWithClock, d passes in real
// time.
func sessionStubWithLatency(d time.Duration) sessionStubOption {
	return func(s *sessionStub) { s.group.SetLatency(d) }
}

// sessionStubFailing returns the option that arms err on each method that
// can fail. Each call of such a method then fails with err.
func sessionStubFailing(err error) sessionStubOption {
	return func(s *sessionStub) {
		s.OnClose.Faults(err, 1)
	}
}

// withSessionClose returns the option that sets the function of Close.
func withSessionClose(fn func() error) sessionStubOption {
	return func(s *sessionStub) { s.OnClose.Func(fn) }
}

// sessionStub is a recording test double of session. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type sessionStub struct {
	OnClose *sessionCloseStub

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// The compiler checks that sessionStub implements session.
var _ session = (*sessionStub)(nil)

// newSessionStub returns a double of session that reports its failures to
// tb and applies opts. The double checks the expected counts of calls at
// the end of the test. With a nil tb, the double reports nothing.
func newSessionStub(tb assert.TB, opts ...sessionStubOption) *sessionStub {
	s := &sessionStub{
		OnClose: &sessionCloseStub{
			MethodStub: stub.NewMethodStub[sessionCloseCall](tb, "session.Close"),
		},
	}
	s.group = stub.NewGroup(
		s.OnClose.MethodStub,
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

// newSessionStubDelayed returns a double whose methods sleep d on clk
// before they return. It applies opts after the clock and the latency.
func newSessionStubDelayed(
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...sessionStubOption,
) *sessionStub {
	return newSessionStub(tb, append([]sessionStubOption{
		sessionStubWithClock(clk),
		sessionStubWithLatency(d),
	}, opts...)...)
}

// newSessionStubFailing returns a double on which each call of a method
// that can fail fails with err. It applies opts after the failure.
func newSessionStubFailing(tb assert.TB, err error, opts ...sessionStubOption) *sessionStub {
	return newSessionStub(tb, append([]sessionStubOption{sessionStubFailing(err)}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *sessionStub) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *sessionCloseStub) invoke() func() sessionCloseReturn {
	if s.fn == nil {
		return nil
	}
	return func() sessionCloseReturn {
		r0 := s.fn()
		return sessionCloseReturn{Err: r0}
	}
}

// Close records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *sessionStub) Close() error {
	call := sessionCloseCall{}
	r := stub.Answer(s.OnClose.MethodStub, call, stub.Arms[sessionCloseCall, sessionCloseReturn]{
		Invoke:   s.OnClose.invoke(),
		Fallback: s.OnClose.fallback,
		Fault:    func(err error) sessionCloseReturn { return sessionCloseReturn{Err: err} },
		Stamp: func(c sessionCloseCall, r sessionCloseReturn) sessionCloseCall {
			c.Err = r.Err
			return c
		},
	})
	return r.Err
}

// FinderFindCall records one call of Finder.Find.
type FinderFindCall[T any] struct {
	Ctx    context.Context
	Name   string
	Result T
	Err    error
}

// FinderFindStub configures Find and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type FinderFindStub[T any] struct {
	*stub.MethodStub[FinderFindCall[T]]

	fn       func(context.Context, string) (T, error)
	fallback *FinderFindReturn[T]
}

// FinderFindReturn contains the fixed results of Find.
type FinderFindReturn[T any] struct {
	Result T
	Err    error
}

// Returns sets the results of every call of Find. A function and an
// injected fault take precedence over them.
func (s *FinderFindStub[T]) Returns(result T, err error) *FinderFindStub[T] {
	s.fallback = &FinderFindReturn[T]{Result: result, Err: err}
	return s
}

// Func sets the function that computes the results of Find from its
// arguments. An injected fault takes precedence over it.
func (s *FinderFindStub[T]) Func(fn func(context.Context, string) (T, error)) *FinderFindStub[T] {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// FinderStubOption configures a [FinderStub] at its construction.
type FinderStubOption[T any] func(*FinderStub[T])

// FinderStubStrict returns the option that makes an unconfigured method
// fail the test instead of returning its zero value.
func FinderStubStrict[T any]() FinderStubOption[T] {
	return func(s *FinderStub[T]) { s.strict = true }
}

// FinderStubDelegateTo returns the option that passes each call to impl.
// The double still records the call.
func FinderStubDelegateTo[T any](impl Finder[T]) FinderStubOption[T] {
	return func(s *FinderStub[T]) {
		s.OnFind.Func(impl.Find)
	}
}

// FinderStubWithClock returns the option that sets the clock of the double
// to clk. The double takes the latency, the windows of the faults and the
// instants of the calls from clk. A latency on a controlled clock costs no
// real time.
func FinderStubWithClock[T any](clk assert.Clock) FinderStubOption[T] {
	return func(s *FinderStub[T]) { s.group.SetClock(clk) }
}

// FinderStubWithRandSource returns the option that sets the random source
// of the probabilistic faults to src. A seeded source repeats the faults
// of a run.
func FinderStubWithRandSource[T any](src rand.Source) FinderStubOption[T] {
	return func(s *FinderStub[T]) { s.group.SetRandSource(src) }
}

// FinderStubBenchMode returns the option that turns the recording of calls
// off. The methods still return their results.
func FinderStubBenchMode[T any]() FinderStubOption[T] {
	return func(s *FinderStub[T]) { s.group.BenchMode() }
}

// FinderStubWithLatency returns the option that makes every method sleep d
// before it returns. Without FinderStubWithClock, d passes in real time.
func FinderStubWithLatency[T any](d time.Duration) FinderStubOption[T] {
	return func(s *FinderStub[T]) { s.group.SetLatency(d) }
}

// FinderStubFailing returns the option that arms err on each method that
// can fail. Each call of such a method then fails with err.
func FinderStubFailing[T any](err error) FinderStubOption[T] {
	return func(s *FinderStub[T]) {
		s.OnFind.Faults(err, 1)
	}
}

// WithFinderFind returns the option that sets the function of Find.
func WithFinderFind[T any](fn func(context.Context, string) (T, error)) FinderStubOption[T] {
	return func(s *FinderStub[T]) { s.OnFind.Func(fn) }
}

// FinderStub is a recording test double of Finder. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type FinderStub[T any] struct {
	OnFind *FinderFindStub[T]

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// NewFinderStub returns a double of Finder that reports its failures to tb
// and applies opts. The double checks the expected counts of calls at the
// end of the test. With a nil tb, the double reports nothing.
func NewFinderStub[T any](tb assert.TB, opts ...FinderStubOption[T]) *FinderStub[T] {
	s := &FinderStub[T]{
		OnFind: &FinderFindStub[T]{
			MethodStub: stub.NewMethodStub[FinderFindCall[T]](tb, "Finder.Find"),
		},
	}
	s.group = stub.NewGroup(
		s.OnFind.MethodStub,
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

// NewFinderStubDelayed returns a double whose methods sleep d on clk
// before they return. It applies opts after the clock and the latency.
func NewFinderStubDelayed[T any](
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...FinderStubOption[T],
) *FinderStub[T] {
	return NewFinderStub[T](tb, append([]FinderStubOption[T]{
		FinderStubWithClock[T](clk),
		FinderStubWithLatency[T](d),
	}, opts...)...)
}

// NewFinderStubFailing returns a double on which each call of a method
// that can fail fails with err. It applies opts after the failure.
func NewFinderStubFailing[T any](
	tb assert.TB,
	err error,
	opts ...FinderStubOption[T],
) *FinderStub[T] {
	return NewFinderStub[T](
		tb,
		append([]FinderStubOption[T]{FinderStubFailing[T](err)}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *FinderStub[T]) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *FinderFindStub[T]) invoke(ctx context.Context, name string) func() FinderFindReturn[T] {
	if s.fn == nil {
		return nil
	}
	return func() FinderFindReturn[T] {
		r0, r1 := s.fn(ctx, name)
		return FinderFindReturn[T]{Result: r0, Err: r1}
	}
}

// Find records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *FinderStub[T]) Find(ctx context.Context, name string) (T, error) {
	call := FinderFindCall[T]{Ctx: ctx, Name: name}
	r := stub.Answer(s.OnFind.MethodStub, call, stub.Arms[FinderFindCall[T], FinderFindReturn[T]]{
		Invoke:   s.OnFind.invoke(ctx, name),
		Fallback: s.OnFind.fallback,
		Fault:    func(err error) FinderFindReturn[T] { return FinderFindReturn[T]{Err: err} },
		Stamp: func(c FinderFindCall[T], r FinderFindReturn[T]) FinderFindCall[T] {
			c.Result = r.Result
			c.Err = r.Err
			return c
		},
	})
	return r.Result, r.Err
}
