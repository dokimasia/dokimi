package mixed

import (
	"context"
	rand "math/rand/v2"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// SessionPrepareCall records one call of Session.Prepare.
type SessionPrepareCall struct {
	Ctx context.Context
	Err error
}

// SessionCommitCall records one call of Session.Commit.
type SessionCommitCall struct {
	Ctx context.Context
	Err error
}

// SessionOldCall records one call of Session.Old.
type SessionOldCall struct {
	Key    string
	Result string
	Err    error
}

// SessionPrepareStub configures Prepare and records its calls. A fault
// armed through it takes precedence over the function and the fixed
// results.
type SessionPrepareStub struct {
	*stub.MethodStub[SessionPrepareCall]

	fn       func(context.Context) error
	fallback *SessionPrepareReturn
}

// SessionPrepareReturn contains the fixed results of Prepare.
type SessionPrepareReturn struct {
	Err error
}

// Returns sets the results of every call of Prepare. A function and an
// injected fault take precedence over them.
func (s *SessionPrepareStub) Returns(err error) *SessionPrepareStub {
	s.fallback = &SessionPrepareReturn{Err: err}
	return s
}

// Func sets the function that computes the results of Prepare from its
// arguments. An injected fault takes precedence over it.
func (s *SessionPrepareStub) Func(fn func(context.Context) error) *SessionPrepareStub {
	s.fn = fn
	return s
}

// SessionCommitStub configures Commit and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type SessionCommitStub struct {
	*stub.MethodStub[SessionCommitCall]

	fn       func(context.Context) error
	fallback *SessionCommitReturn
}

// SessionCommitReturn contains the fixed results of Commit.
type SessionCommitReturn struct {
	Err error
}

// Returns sets the results of every call of Commit. A function and an
// injected fault take precedence over them.
func (s *SessionCommitStub) Returns(err error) *SessionCommitStub {
	s.fallback = &SessionCommitReturn{Err: err}
	return s
}

// Func sets the function that computes the results of Commit from its
// arguments. An injected fault takes precedence over it.
func (s *SessionCommitStub) Func(fn func(context.Context) error) *SessionCommitStub {
	s.fn = fn
	return s
}

// SessionOldStub configures Old and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type SessionOldStub struct {
	*stub.MethodStub[SessionOldCall]

	fn       func(string) (string, error)
	fallback *SessionOldReturn
}

// SessionOldReturn contains the fixed results of Old.
type SessionOldReturn struct {
	Result string
	Err    error
}

// Returns sets the results of every call of Old. A function and an
// injected fault take precedence over them.
func (s *SessionOldStub) Returns(result string, err error) *SessionOldStub {
	s.fallback = &SessionOldReturn{Result: result, Err: err}
	return s
}

// Func sets the function that computes the results of Old from its
// arguments. An injected fault takes precedence over it.
func (s *SessionOldStub) Func(fn func(string) (string, error)) *SessionOldStub {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// SessionStubOption configures a [SessionStub] at its construction.
type SessionStubOption func(*SessionStub)

// SessionStubStrict returns the option that makes an unconfigured method
// fail the test instead of returning its zero value.
func SessionStubStrict() SessionStubOption {
	return func(s *SessionStub) { s.strict = true }
}

// SessionStubDelegateTo returns the option that passes each call to impl.
// The double still records the call.
func SessionStubDelegateTo(impl Session) SessionStubOption {
	return func(s *SessionStub) {
		s.OnPrepare.Func(impl.Prepare)
		s.OnCommit.Func(impl.Commit)
		s.OnOld.Func(impl.Old)
	}
}

// SessionStubWithClock returns the option that sets the clock of the
// double to clk. The double takes the latency, the windows of the faults
// and the instants of the calls from clk. A latency on a controlled clock
// costs no real time.
func SessionStubWithClock(clk assert.Clock) SessionStubOption {
	return func(s *SessionStub) { s.group.SetClock(clk) }
}

// SessionStubWithRandSource returns the option that sets the random source
// of the probabilistic faults to src. A seeded source repeats the faults
// of a run.
func SessionStubWithRandSource(src rand.Source) SessionStubOption {
	return func(s *SessionStub) { s.group.SetRandSource(src) }
}

// SessionStubBenchMode returns the option that turns the recording of
// calls off. The methods still return their results.
func SessionStubBenchMode() SessionStubOption {
	return func(s *SessionStub) { s.group.BenchMode() }
}

// SessionStubWithLatency returns the option that makes every method sleep
// d before it returns. Without SessionStubWithClock, d passes in real
// time.
func SessionStubWithLatency(d time.Duration) SessionStubOption {
	return func(s *SessionStub) { s.group.SetLatency(d) }
}

// SessionStubFailing returns the option that arms err on each method that
// can fail. Each call of such a method then fails with err.
func SessionStubFailing(err error) SessionStubOption {
	return func(s *SessionStub) {
		s.OnPrepare.Faults(err, 1)
		s.OnCommit.Faults(err, 1)
		s.OnOld.Faults(err, 1)
	}
}

// WithSessionPrepare returns the option that sets the function of Prepare.
func WithSessionPrepare(fn func(context.Context) error) SessionStubOption {
	return func(s *SessionStub) { s.OnPrepare.Func(fn) }
}

// WithSessionCommit returns the option that sets the function of Commit.
func WithSessionCommit(fn func(context.Context) error) SessionStubOption {
	return func(s *SessionStub) { s.OnCommit.Func(fn) }
}

// WithSessionOld returns the option that sets the function of Old.
func WithSessionOld(fn func(string) (string, error)) SessionStubOption {
	return func(s *SessionStub) { s.OnOld.Func(fn) }
}

// SessionStub is a recording test double of Session. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type SessionStub struct {
	OnPrepare *SessionPrepareStub
	OnCommit  *SessionCommitStub
	OnOld     *SessionOldStub

	// group contains the stub of each method.
	group *stub.Group

	order  *stub.OrderTracker
	strict bool
}

// The compiler checks that SessionStub implements Session.
var _ Session = (*SessionStub)(nil)

// NewSessionStub returns a double of Session that reports its failures to
// tb and applies opts. The double checks the expected counts of calls at
// the end of the test. With a nil tb, the double reports nothing.
func NewSessionStub(tb assert.TB, opts ...SessionStubOption) *SessionStub {
	s := &SessionStub{
		OnPrepare: &SessionPrepareStub{
			MethodStub: stub.NewMethodStub[SessionPrepareCall](tb, "Session.Prepare"),
		},
		OnCommit: &SessionCommitStub{
			MethodStub: stub.NewMethodStub[SessionCommitCall](tb, "Session.Commit"),
		},
		OnOld: &SessionOldStub{
			MethodStub: stub.NewMethodStub[SessionOldCall](tb, "Session.Old"),
		},
	}
	s.group = stub.NewGroup(
		s.OnPrepare.MethodStub,
		s.OnCommit.MethodStub,
		s.OnOld.MethodStub,
	)
	for _, opt := range opts {
		opt(s)
	}
	if s.strict {
		s.group.StrictAll()
	}
	// The tracker fails the test in strict mode alone, so it follows the
	// options.
	s.order = stub.NewOrderTracker(tb, s.strict)
	s.group.Bind(tb)
	return s
}

// NewSessionStubDelayed returns a double whose methods sleep d on clk
// before they return. It applies opts after the clock and the latency.
func NewSessionStubDelayed(
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...SessionStubOption,
) *SessionStub {
	return NewSessionStub(tb, append([]SessionStubOption{
		SessionStubWithClock(clk),
		SessionStubWithLatency(d),
	}, opts...)...)
}

// NewSessionStubFailing returns a double on which each call of a method
// that can fail fails with err. It applies opts after the failure.
func NewSessionStubFailing(tb assert.TB, err error, opts ...SessionStubOption) *SessionStub {
	return NewSessionStub(tb, append([]SessionStubOption{SessionStubFailing(err)}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *SessionStub) ResetCalls() {
	s.group.Reset()
	s.order.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *SessionPrepareStub) invoke(ctx context.Context) func() SessionPrepareReturn {
	if s.fn == nil {
		return nil
	}
	return func() SessionPrepareReturn {
		r0 := s.fn(ctx)
		return SessionPrepareReturn{Err: r0}
	}
}

// Prepare records the call. An armed fault comes first, then the function
// of Func, then the results of Returns, and then the zero values.
func (s *SessionStub) Prepare(ctx context.Context) error {
	defer s.order.Record("Prepare")
	call := SessionPrepareCall{Ctx: ctx}
	r := stub.Answer(
		s.OnPrepare.MethodStub,
		call,
		stub.Arms[SessionPrepareCall, SessionPrepareReturn]{
			Invoke:   s.OnPrepare.invoke(ctx),
			Fallback: s.OnPrepare.fallback,
			Fault:    func(err error) SessionPrepareReturn { return SessionPrepareReturn{Err: err} },
			Stamp: func(c SessionPrepareCall, r SessionPrepareReturn) SessionPrepareCall {
				c.Err = r.Err
				return c
			},
		},
	)
	return r.Err
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *SessionCommitStub) invoke(ctx context.Context) func() SessionCommitReturn {
	if s.fn == nil {
		return nil
	}
	return func() SessionCommitReturn {
		r0 := s.fn(ctx)
		return SessionCommitReturn{Err: r0}
	}
}

// Commit records the call. An armed fault comes first, then the function
// of Func, then the results of Returns, and then the zero values.
func (s *SessionStub) Commit(ctx context.Context) error {
	s.order.AssertAfter("Commit", "Prepare")
	defer s.order.Record("Commit")
	call := SessionCommitCall{Ctx: ctx}
	r := stub.Answer(s.OnCommit.MethodStub, call, stub.Arms[SessionCommitCall, SessionCommitReturn]{
		Invoke:   s.OnCommit.invoke(ctx),
		Fallback: s.OnCommit.fallback,
		Fault:    func(err error) SessionCommitReturn { return SessionCommitReturn{Err: err} },
		Stamp: func(c SessionCommitCall, r SessionCommitReturn) SessionCommitCall {
			c.Err = r.Err
			return c
		},
	})
	return r.Err
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *SessionOldStub) invoke(key string) func() SessionOldReturn {
	if s.fn == nil {
		return nil
	}
	return func() SessionOldReturn {
		r0, r1 := s.fn(key)
		return SessionOldReturn{Result: r0, Err: r1}
	}
}

// Old records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *SessionStub) Old(key string) (string, error) {
	// Bench mode logs nothing, because a log allocates on each call.
	if tb, logs := s.OnOld.TB().(interface{ Logf(string, ...any) }); logs &&
		!s.OnOld.IsBenchMode() {
		tb.Logf("Session.Old is deprecated")
	}
	defer s.order.Record("Old")
	call := SessionOldCall{Key: key}
	r := stub.Answer(s.OnOld.MethodStub, call, stub.Arms[SessionOldCall, SessionOldReturn]{
		Invoke:   s.OnOld.invoke(key),
		Fallback: s.OnOld.fallback,
		Fault:    func(err error) SessionOldReturn { return SessionOldReturn{Err: err} },
		Stamp: func(c SessionOldCall, r SessionOldReturn) SessionOldCall {
			c.Result = r.Result
			c.Err = r.Err
			return c
		},
	})
	return r.Result, r.Err
}
