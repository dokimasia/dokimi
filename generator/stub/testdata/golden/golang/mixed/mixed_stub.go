package mixed

import (
	"context"
	rand "math/rand/v2"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// SessionPrepareCall records one invocation of Session.Prepare.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type SessionPrepareCall struct {
	Ctx context.Context
	Err error
}

// SessionCommitCall records one invocation of Session.Commit.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type SessionCommitCall struct {
	Ctx context.Context
	Err error
}

// SessionOldCall records one invocation of Session.Old.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type SessionOldCall struct {
	Key    string
	Result string
	Err    error
}

// --- Per-method configuration ---

// SessionPrepareStub controls how the double answers Prepare and records
// what it was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type SessionPrepareStub struct {
	*stub.MethodStub[SessionPrepareCall]

	fn       func(context.Context) error
	fallback *SessionPrepareReturn
}

// SessionPrepareReturn holds the fixed answer configured through Returns.
type SessionPrepareReturn struct {
	Err error
}

// Returns pins a fixed result for every call to Prepare. A Func override
// takes precedence over it, and so does an injected fault.
func (s *SessionPrepareStub) Returns(err error) *SessionPrepareStub {
	s.fallback = &SessionPrepareReturn{Err: err}
	return s
}

// Func supplies a body for Prepare, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *SessionPrepareStub) Func(fn func(context.Context) error) *SessionPrepareStub {
	s.fn = fn
	return s
}

// SessionCommitStub controls how the double answers Commit and records
// what it was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type SessionCommitStub struct {
	*stub.MethodStub[SessionCommitCall]

	fn       func(context.Context) error
	fallback *SessionCommitReturn
}

// SessionCommitReturn holds the fixed answer configured through Returns.
type SessionCommitReturn struct {
	Err error
}

// Returns pins a fixed result for every call to Commit. A Func override
// takes precedence over it, and so does an injected fault.
func (s *SessionCommitStub) Returns(err error) *SessionCommitStub {
	s.fallback = &SessionCommitReturn{Err: err}
	return s
}

// Func supplies a body for Commit, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *SessionCommitStub) Func(fn func(context.Context) error) *SessionCommitStub {
	s.fn = fn
	return s
}

// SessionOldStub controls how the double answers Old and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type SessionOldStub struct {
	*stub.MethodStub[SessionOldCall]

	fn       func(string) (string, error)
	fallback *SessionOldReturn
}

// SessionOldReturn holds the fixed answer configured through Returns.
type SessionOldReturn struct {
	Result string
	Err    error
}

// Returns pins a fixed result for every call to Old. A Func override takes
// precedence over it, and so does an injected fault.
func (s *SessionOldStub) Returns(result string, err error) *SessionOldStub {
	s.fallback = &SessionOldReturn{Result: result, Err: err}
	return s
}

// Func supplies a body for Old, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *SessionOldStub) Func(fn func(string) (string, error)) *SessionOldStub {
	s.fn = fn
	return s
}

// --- SessionStub ---

// The witness of the major version of the runtime of doubles.
var _ = stub.CompatV1

// SessionStubOption configures a [SessionStub] at construction time.
type SessionStubOption func(*SessionStub)

// SessionStubStrict makes an unconfigured method fail the test rather than
// return its zero value.
func SessionStubStrict() SessionStubOption {
	return func(s *SessionStub) { s.strict = true }
}

// SessionStubDelegateTo wraps a real implementation: every method forwards
// to impl and is recorded on the double on the way through.
func SessionStubDelegateTo(impl Session) SessionStubOption {
	return func(s *SessionStub) {
		s.OnPrepare.Func(impl.Prepare)
		s.OnCommit.Func(impl.Commit)
		s.OnOld.Func(impl.Old)
	}
}

// SessionStubWithClock drives latency and time-windowed faults from clk
// rather than from wall time, so a test asserting on a five-second timeout
// does not take five seconds.
func SessionStubWithClock(clk assert.Clock) SessionStubOption {
	return func(s *SessionStub) { s.group.SetClock(clk) }
}

// SessionStubWithRandSource makes probabilistic fault injection
// reproducible. A failure nobody can replay is not a finding.
func SessionStubWithRandSource(src rand.Source) SessionStubOption {
	return func(s *SessionStub) { s.group.SetRandSource(src) }
}

// SessionStubBenchMode disables call recording. Dispatch still works; only
// the accumulating call log is dropped.
func SessionStubBenchMode() SessionStubOption {
	return func(s *SessionStub) { s.group.BenchMode() }
}

// SessionStubWithLatency makes every method sleep d before answering, for
// a test standing in for a slow dependency rather than one slow method.
// Pair it with SessionStubWithClock, or d elapses in wall time.
func SessionStubWithLatency(d time.Duration) SessionStubOption {
	return func(s *SessionStub) { s.group.SetLatency(d) }
}

// SessionStubFailing arms every method that can fail with err, so a
// caller's error path is one option rather than one Faults call per
// method. A method answering no error is left alone.
func SessionStubFailing(err error) SessionStubOption {
	return func(s *SessionStub) {
		s.OnPrepare.Faults(err, 1)
		s.OnCommit.Faults(err, 1)
		s.OnOld.Faults(err, 1)
	}
}

// WithSessionPrepare sets Prepare's body at construction time, for the
// common case of configuring one method and taking the defaults for the
// rest.
func WithSessionPrepare(fn func(context.Context) error) SessionStubOption {
	return func(s *SessionStub) { s.OnPrepare.Func(fn) }
}

// WithSessionCommit sets Commit's body at construction time, for the
// common case of configuring one method and taking the defaults for the
// rest.
func WithSessionCommit(fn func(context.Context) error) SessionStubOption {
	return func(s *SessionStub) { s.OnCommit.Func(fn) }
}

// WithSessionOld sets Old's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithSessionOld(fn func(string) (string, error)) SessionStubOption {
	return func(s *SessionStub) { s.OnOld.Func(fn) }
}

// SessionStub is a recording test double for Session.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method returns its zero value and records the call.
type SessionStub struct {
	OnPrepare *SessionPrepareStub
	OnCommit  *SessionCommitStub
	OnOld     *SessionOldStub

	// group is every method stub above, viewed through the surface that
	// does not depend on a signature.
	group *stub.Group

	order  *stub.OrderTracker
	strict bool
}

// Compile-time proof that SessionStub satisfies Session.
var _ Session = (*SessionStub)(nil)

// NewSessionStub returns a double bound to tb.
//
// Passing tb registers a cleanup that verifies every method's call-count
// expectations when the test ends, so an unmet Times is reported without
// the caller remembering to check. A nil tb skips that, which is what
// benchmarks and non-test callers want.
func NewSessionStub(tb assert.TB, opts ...SessionStubOption) *SessionStub {
	s := &SessionStub{
		OnPrepare: &SessionPrepareStub{MethodStub: stub.NewMethodStub[SessionPrepareCall](tb, "Session.Prepare")},
		OnCommit:  &SessionCommitStub{MethodStub: stub.NewMethodStub[SessionCommitCall](tb, "Session.Commit")},
		OnOld:     &SessionOldStub{MethodStub: stub.NewMethodStub[SessionOldCall](tb, "Session.Old")},
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
	// Allocated after the options run: the tracker only fails the test in
	// strict mode, so it has to be told which mode it is in.
	s.order = stub.NewOrderTracker(tb, s.strict)
	s.group.Bind(tb)
	return s
}

// NewSessionStubDelayed returns a double whose every method sleeps d
// against clk before answering — the clock and the latency a
// slow-dependency test wants, in one call rather than two options. Further
// options apply after the pair, so a caller can still add Strict or a
// method override.
func NewSessionStubDelayed(tb assert.TB, d time.Duration, clk assert.Clock, opts ...SessionStubOption) *SessionStub {
	return NewSessionStub(tb, append([]SessionStubOption{
		SessionStubWithClock(clk),
		SessionStubWithLatency(d),
	}, opts...)...)
}

// NewSessionStubFailing returns a double whose every fallible method
// answers err — [SessionStubFailing] applied at construction, for the test
// whose whole subject is the caller's error path. Further options apply
// after it.
func NewSessionStubFailing(tb assert.TB, err error, opts ...SessionStubOption) *SessionStub {
	return NewSessionStub(tb, append([]SessionStubOption{SessionStubFailing(err)}, opts...)...)
}

// ResetCalls clears recorded calls, fault counters, and call-count
// expectations. Func and Returns configuration is deliberately preserved, so
// a test can reuse one double across phases without rebuilding it.
func (s *SessionStub) ResetCalls() {
	s.group.Reset()
	s.order.Reset()
}

// --- Interface methods ---

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *SessionPrepareStub) invoke(ctx context.Context) func() SessionPrepareReturn {
	if s.fn == nil {
		return nil
	}
	return func() SessionPrepareReturn {
		r0 := s.fn(ctx)
		return SessionPrepareReturn{Err: r0}
	}
}

// Prepare records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *SessionStub) Prepare(ctx context.Context) error {
	defer s.order.Record("Prepare")
	call := SessionPrepareCall{Ctx: ctx}
	r := stub.Answer(s.OnPrepare.MethodStub, &call, stub.Arms[SessionPrepareCall, SessionPrepareReturn]{
		Invoke:   s.OnPrepare.invoke(ctx),
		Fallback: s.OnPrepare.fallback,
		Fault:    func(err error) SessionPrepareReturn { return SessionPrepareReturn{Err: err} },
		Stamp: func(c *SessionPrepareCall, r SessionPrepareReturn) {
			c.Err = r.Err
		},
	})
	return r.Err
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *SessionCommitStub) invoke(ctx context.Context) func() SessionCommitReturn {
	if s.fn == nil {
		return nil
	}
	return func() SessionCommitReturn {
		r0 := s.fn(ctx)
		return SessionCommitReturn{Err: r0}
	}
}

// Commit records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *SessionStub) Commit(ctx context.Context) error {
	s.order.AssertAfter("Commit", "Prepare")
	defer s.order.Record("Commit")
	call := SessionCommitCall{Ctx: ctx}
	r := stub.Answer(s.OnCommit.MethodStub, &call, stub.Arms[SessionCommitCall, SessionCommitReturn]{
		Invoke:   s.OnCommit.invoke(ctx),
		Fallback: s.OnCommit.fallback,
		Fault:    func(err error) SessionCommitReturn { return SessionCommitReturn{Err: err} },
		Stamp: func(c *SessionCommitCall, r SessionCommitReturn) {
			c.Err = r.Err
		},
	})
	return r.Err
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *SessionOldStub) invoke(key string) func() SessionOldReturn {
	if s.fn == nil {
		return nil
	}
	return func() SessionOldReturn {
		r0, r1 := s.fn(key)
		return SessionOldReturn{Result: r0, Err: r1}
	}
}

// Old records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *SessionStub) Old(key string) (string, error) {
	// Skipped in bench mode: Logf allocates the formatted message per
	// dispatch, which is what a benchmark through this double would
	// otherwise be measuring.
	if tb, logs := s.OnOld.TB().(interface{ Logf(string, ...any) }); logs && !s.OnOld.IsBenchMode() {
		tb.Logf("Session.Old is deprecated")
	}
	defer s.order.Record("Old")
	call := SessionOldCall{Key: key}
	r := stub.Answer(s.OnOld.MethodStub, &call, stub.Arms[SessionOldCall, SessionOldReturn]{
		Invoke:   s.OnOld.invoke(key),
		Fallback: s.OnOld.fallback,
		Fault:    func(err error) SessionOldReturn { return SessionOldReturn{Err: err} },
		Stamp: func(c *SessionOldCall, r SessionOldReturn) {
			c.Result = r.Result
			c.Err = r.Err
		},
	})
	return r.Result, r.Err
}
