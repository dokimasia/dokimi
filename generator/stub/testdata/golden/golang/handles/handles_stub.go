package handles

import (
	"context"
	rand "math/rand/v2"
	"time"

	"example.com/shop/other"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/dokimi-go/stub"
)

// ArchiveSizeCall records one invocation of Archive.Size.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type ArchiveSizeCall struct {
	Result int
}

// ArchiveOpenLogCall records one invocation of Archive.OpenLog.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type ArchiveOpenLogCall struct {
	Name   string
	Result Log
	Err    error
}

// --- Per-method configuration ---

// ArchiveSizeStub controls how the double answers Size and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type ArchiveSizeStub struct {
	*stub.MethodStub[ArchiveSizeCall]

	fn       func() int
	fallback *ArchiveSizeReturn
}

// ArchiveSizeReturn holds the fixed answer configured through Returns.
type ArchiveSizeReturn struct {
	Result int
}

// Returns pins a fixed result for every call to Size. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Size answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *ArchiveSizeStub) Returns(result int) *ArchiveSizeStub {
	s.fallback = &ArchiveSizeReturn{Result: result}
	return s
}

// Func supplies a body for Size, for when the answer depends on the
// arguments. Nothing takes precedence over it: Size answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *ArchiveSizeStub) Func(fn func() int) *ArchiveSizeStub {
	s.fn = fn
	return s
}

// ArchiveOpenLogStub controls how the double answers OpenLog and records
// what it was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type ArchiveOpenLogStub struct {
	*stub.MethodStub[ArchiveOpenLogCall]

	fn       func(string) (Log, error)
	fallback *ArchiveOpenLogReturn
}

// ArchiveOpenLogReturn holds the fixed answer configured through Returns.
type ArchiveOpenLogReturn struct {
	Result Log
	Err    error
}

// Returns pins a fixed result for every call to OpenLog. A Func override
// takes precedence over it, and so does an injected fault.
func (s *ArchiveOpenLogStub) Returns(result Log, err error) *ArchiveOpenLogStub {
	s.fallback = &ArchiveOpenLogReturn{Result: result, Err: err}
	return s
}

// Func supplies a body for OpenLog, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *ArchiveOpenLogStub) Func(fn func(string) (Log, error)) *ArchiveOpenLogStub {
	s.fn = fn
	return s
}

// --- ArchiveStub ---

// The double-runtime witness for the surface this double rides.
var _ = stub.CompatV1

// ArchiveStubOption configures a [ArchiveStub] at construction time.
type ArchiveStubOption func(*ArchiveStub)

// ArchiveStubStrict makes an unconfigured method fail the test rather than
// return its zero value.
func ArchiveStubStrict() ArchiveStubOption {
	return func(s *ArchiveStub) { s.strict = true }
}

// ArchiveStubDelegateTo wraps a real implementation: every method forwards
// to impl and is recorded on the double on the way through.
func ArchiveStubDelegateTo(impl Archive) ArchiveStubOption {
	return func(s *ArchiveStub) {
		s.OnSize.Func(impl.Size)
		s.OnOpenLog.Func(impl.OpenLog)
	}
}

// ArchiveStubWithClock drives latency and time-windowed faults from clk
// rather than from wall time, so a test asserting on a five-second timeout
// does not take five seconds.
func ArchiveStubWithClock(clk assert.Clock) ArchiveStubOption {
	return func(s *ArchiveStub) { s.group.SetClock(clk) }
}

// ArchiveStubWithRandSource makes probabilistic fault injection
// reproducible. A failure nobody can replay is not a finding.
func ArchiveStubWithRandSource(src rand.Source) ArchiveStubOption {
	return func(s *ArchiveStub) { s.group.SetRandSource(src) }
}

// ArchiveStubBenchMode disables call recording. Dispatch still works; only
// the accumulating call log is dropped.
func ArchiveStubBenchMode() ArchiveStubOption {
	return func(s *ArchiveStub) { s.group.BenchMode() }
}

// ArchiveStubWithLatency makes every method sleep d before answering, for
// a test standing in for a slow dependency rather than one slow method.
// Pair it with ArchiveStubWithClock, or d elapses in wall time.
func ArchiveStubWithLatency(d time.Duration) ArchiveStubOption {
	return func(s *ArchiveStub) { s.group.SetLatency(d) }
}

// ArchiveStubFailing arms every method that can fail with err, so a
// caller's error path is one option rather than one Faults call per
// method. A method answering no error is left alone.
func ArchiveStubFailing(err error) ArchiveStubOption {
	return func(s *ArchiveStub) {
		s.OnOpenLog.Faults(err, 1)
	}
}

// WithArchiveSize sets Size's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithArchiveSize(fn func() int) ArchiveStubOption {
	return func(s *ArchiveStub) { s.OnSize.Func(fn) }
}

// WithArchiveOpenLog sets OpenLog's body at construction time, for the
// common case of configuring one method and taking the defaults for the
// rest.
func WithArchiveOpenLog(fn func(string) (Log, error)) ArchiveStubOption {
	return func(s *ArchiveStub) { s.OnOpenLog.Func(fn) }
}

// ArchiveStub is a recording test double for Archive.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method returns its zero value and records the call.
type ArchiveStub struct {
	OnSize *ArchiveSizeStub
	// OnOpenLog configures OpenLog, contributed by Opener.
	OnOpenLog *ArchiveOpenLogStub

	// group is every method stub above, viewed through the surface that
	// does not depend on a signature.
	group  *stub.Group
	strict bool
}

// Compile-time proof that ArchiveStub satisfies Archive.
var _ Archive = (*ArchiveStub)(nil)

// NewArchiveStub returns a double bound to tb.
//
// Passing tb registers a cleanup that verifies every method's call-count
// expectations when the test ends, so an unmet Times is reported without
// the caller remembering to check. A nil tb skips that, which is what
// benchmarks and non-test callers want.
func NewArchiveStub(tb assert.TB, opts ...ArchiveStubOption) *ArchiveStub {
	s := &ArchiveStub{
		OnSize:    &ArchiveSizeStub{MethodStub: stub.NewMethodStub[ArchiveSizeCall](tb, "Archive.Size")},
		OnOpenLog: &ArchiveOpenLogStub{MethodStub: stub.NewMethodStub[ArchiveOpenLogCall](tb, "Archive.OpenLog")},
	}
	s.group = stub.NewGroup(
		s.OnSize.MethodStub,
		s.OnOpenLog.MethodStub,
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

// NewArchiveStubDelayed returns a double whose every method sleeps d
// against clk before answering — the clock and the latency a
// slow-dependency test wants, in one call rather than two options. Further
// options apply after the pair, so a caller can still add Strict or a
// method override.
func NewArchiveStubDelayed(tb assert.TB, d time.Duration, clk assert.Clock, opts ...ArchiveStubOption) *ArchiveStub {
	return NewArchiveStub(tb, append([]ArchiveStubOption{
		ArchiveStubWithClock(clk),
		ArchiveStubWithLatency(d),
	}, opts...)...)
}

// NewArchiveStubFailing returns a double whose every fallible method
// answers err — [ArchiveStubFailing] applied at construction, for the test
// whose whole subject is the caller's error path. Further options apply
// after it.
func NewArchiveStubFailing(tb assert.TB, err error, opts ...ArchiveStubOption) *ArchiveStub {
	return NewArchiveStub(tb, append([]ArchiveStubOption{ArchiveStubFailing(err)}, opts...)...)
}

// ResetCalls clears recorded calls, fault counters, and call-count
// expectations. Func and Returns configuration is deliberately preserved, so
// a test can reuse one double across phases without rebuilding it.
func (s *ArchiveStub) ResetCalls() {
	s.group.Reset()
}

// --- Interface methods ---

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *ArchiveSizeStub) invoke() func() ArchiveSizeReturn {
	if s.fn == nil {
		return nil
	}
	return func() ArchiveSizeReturn {
		r0 := s.fn()
		return ArchiveSizeReturn{Result: r0}
	}
}

// Size records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *ArchiveStub) Size() int {
	call := ArchiveSizeCall{}
	r := stub.Answer(s.OnSize.MethodStub, &call, stub.Arms[ArchiveSizeCall, ArchiveSizeReturn]{
		Invoke:   s.OnSize.invoke(),
		Fallback: s.OnSize.fallback,
		Stamp: func(c *ArchiveSizeCall, r ArchiveSizeReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *ArchiveOpenLogStub) invoke(name string) func() ArchiveOpenLogReturn {
	if s.fn == nil {
		return nil
	}
	return func() ArchiveOpenLogReturn {
		r0, r1 := s.fn(name)
		return ArchiveOpenLogReturn{Result: r0, Err: r1}
	}
}

// OpenLog records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *ArchiveStub) OpenLog(name string) (Log, error) {
	call := ArchiveOpenLogCall{Name: name}
	r := stub.Answer(s.OnOpenLog.MethodStub, &call, stub.Arms[ArchiveOpenLogCall, ArchiveOpenLogReturn]{
		Invoke:   s.OnOpenLog.invoke(name),
		Fallback: s.OnOpenLog.fallback,
		Fault:    func(err error) ArchiveOpenLogReturn { return ArchiveOpenLogReturn{Err: err} },
		Stamp: func(c *ArchiveOpenLogCall, r ArchiveOpenLogReturn) {
			c.Result = r.Result
			c.Err = r.Err
		},
	})
	return r.Result, r.Err
}

// --- Log, the handle Archive's methods answer ---

// LogAppendCall records one invocation of Log.Append.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type LogAppendCall struct {
	Line string
	Err  error
}

// --- Per-method configuration ---

// LogAppendStub controls how the double answers Append and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type LogAppendStub struct {
	*stub.MethodStub[LogAppendCall]

	fn       func(string) error
	fallback *LogAppendReturn
}

// LogAppendReturn holds the fixed answer configured through Returns.
type LogAppendReturn struct {
	Err error
}

// Returns pins a fixed result for every call to Append. A Func override
// takes precedence over it, and so does an injected fault.
func (s *LogAppendStub) Returns(err error) *LogAppendStub {
	s.fallback = &LogAppendReturn{Err: err}
	return s
}

// Func supplies a body for Append, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *LogAppendStub) Func(fn func(string) error) *LogAppendStub {
	s.fn = fn
	return s
}

// --- LogStub ---

// The double-runtime witness for the surface this double rides.
var _ = stub.CompatV1

// LogStubOption configures a [LogStub] at construction time.
type LogStubOption func(*LogStub)

// LogStubStrict makes an unconfigured method fail the test rather than
// return its zero value.
func LogStubStrict() LogStubOption {
	return func(s *LogStub) { s.strict = true }
}

// LogStubDelegateTo wraps a real implementation: every method forwards to
// impl and is recorded on the double on the way through.
func LogStubDelegateTo(impl Log) LogStubOption {
	return func(s *LogStub) {
		s.OnAppend.Func(impl.Append)
	}
}

// LogStubWithClock drives latency and time-windowed faults from clk rather
// than from wall time, so a test asserting on a five-second timeout does
// not take five seconds.
func LogStubWithClock(clk assert.Clock) LogStubOption {
	return func(s *LogStub) { s.group.SetClock(clk) }
}

// LogStubWithRandSource makes probabilistic fault injection reproducible.
// A failure nobody can replay is not a finding.
func LogStubWithRandSource(src rand.Source) LogStubOption {
	return func(s *LogStub) { s.group.SetRandSource(src) }
}

// LogStubBenchMode disables call recording. Dispatch still works; only the
// accumulating call log is dropped.
func LogStubBenchMode() LogStubOption {
	return func(s *LogStub) { s.group.BenchMode() }
}

// LogStubWithLatency makes every method sleep d before answering, for a
// test standing in for a slow dependency rather than one slow method. Pair
// it with LogStubWithClock, or d elapses in wall time.
func LogStubWithLatency(d time.Duration) LogStubOption {
	return func(s *LogStub) { s.group.SetLatency(d) }
}

// LogStubFailing arms every method that can fail with err, so a caller's
// error path is one option rather than one Faults call per method. A
// method answering no error is left alone.
func LogStubFailing(err error) LogStubOption {
	return func(s *LogStub) {
		s.OnAppend.Faults(err, 1)
	}
}

// WithLogAppend sets Append's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithLogAppend(fn func(string) error) LogStubOption {
	return func(s *LogStub) { s.OnAppend.Func(fn) }
}

// LogStub is a recording test double for Log.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method returns its zero value and records the call.
type LogStub struct {
	OnAppend *LogAppendStub

	// group is every method stub above, viewed through the surface that
	// does not depend on a signature.
	group  *stub.Group
	strict bool
}

// Compile-time proof that LogStub satisfies Log.
var _ Log = (*LogStub)(nil)

// NewLogStub returns a double bound to tb.
//
// Passing tb registers a cleanup that verifies every method's call-count
// expectations when the test ends, so an unmet Times is reported without
// the caller remembering to check. A nil tb skips that, which is what
// benchmarks and non-test callers want.
func NewLogStub(tb assert.TB, opts ...LogStubOption) *LogStub {
	s := &LogStub{
		OnAppend: &LogAppendStub{MethodStub: stub.NewMethodStub[LogAppendCall](tb, "Log.Append")},
	}
	s.group = stub.NewGroup(
		s.OnAppend.MethodStub,
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

// NewLogStubDelayed returns a double whose every method sleeps d against
// clk before answering — the clock and the latency a slow-dependency test
// wants, in one call rather than two options. Further options apply after
// the pair, so a caller can still add Strict or a method override.
func NewLogStubDelayed(tb assert.TB, d time.Duration, clk assert.Clock, opts ...LogStubOption) *LogStub {
	return NewLogStub(tb, append([]LogStubOption{
		LogStubWithClock(clk),
		LogStubWithLatency(d),
	}, opts...)...)
}

// NewLogStubFailing returns a double whose every fallible method answers
// err — [LogStubFailing] applied at construction, for the test whose whole
// subject is the caller's error path. Further options apply after it.
func NewLogStubFailing(tb assert.TB, err error, opts ...LogStubOption) *LogStub {
	return NewLogStub(tb, append([]LogStubOption{LogStubFailing(err)}, opts...)...)
}

// ResetCalls clears recorded calls, fault counters, and call-count
// expectations. Func and Returns configuration is deliberately preserved, so
// a test can reuse one double across phases without rebuilding it.
func (s *LogStub) ResetCalls() {
	s.group.Reset()
}

// --- Interface methods ---

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *LogAppendStub) invoke(line string) func() LogAppendReturn {
	if s.fn == nil {
		return nil
	}
	return func() LogAppendReturn {
		r0 := s.fn(line)
		return LogAppendReturn{Err: r0}
	}
}

// Append records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *LogStub) Append(line string) error {
	call := LogAppendCall{Line: line}
	r := stub.Answer(s.OnAppend.MethodStub, &call, stub.Arms[LogAppendCall, LogAppendReturn]{
		Invoke:   s.OnAppend.invoke(line),
		Fallback: s.OnAppend.fallback,
		Fault:    func(err error) LogAppendReturn { return LogAppendReturn{Err: err} },
		Stamp: func(c *LogAppendCall, r LogAppendReturn) {
			c.Err = r.Err
		},
	})
	return r.Err
}

// ProducerOpenCall records one invocation of Producer.Open.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type ProducerOpenCall struct {
	Ctx    context.Context
	Result Cursor
	Err    error
}

// ProducerResumeCall records one invocation of Producer.Resume.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type ProducerResumeCall struct {
	Token  string
	Result Cursor
}

// ProducerWatchCall records one invocation of Producer.Watch.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type ProducerWatchCall struct {
	Result Subscription
}

// ProducerNothingCall records one invocation of Producer.Nothing.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type ProducerNothingCall struct {
	Result Empty
}

// ProducerElsewhereCall records one invocation of Producer.Elsewhere.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type ProducerElsewhereCall struct {
	Result other.Thing
}

// --- Per-method configuration ---

// ProducerOpenStub controls how the double answers Open and records what
// it was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type ProducerOpenStub struct {
	*stub.MethodStub[ProducerOpenCall]

	fn       func(context.Context) (Cursor, error)
	fallback *ProducerOpenReturn
}

// ProducerOpenReturn holds the fixed answer configured through Returns.
type ProducerOpenReturn struct {
	Result Cursor
	Err    error
}

// Returns pins a fixed result for every call to Open. A Func override
// takes precedence over it, and so does an injected fault.
func (s *ProducerOpenStub) Returns(result Cursor, err error) *ProducerOpenStub {
	s.fallback = &ProducerOpenReturn{Result: result, Err: err}
	return s
}

// Func supplies a body for Open, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *ProducerOpenStub) Func(fn func(context.Context) (Cursor, error)) *ProducerOpenStub {
	s.fn = fn
	return s
}

// ProducerResumeStub controls how the double answers Resume and records
// what it was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type ProducerResumeStub struct {
	*stub.MethodStub[ProducerResumeCall]

	fn       func(string) Cursor
	fallback *ProducerResumeReturn
}

// ProducerResumeReturn holds the fixed answer configured through Returns.
type ProducerResumeReturn struct {
	Result Cursor
}

// Returns pins a fixed result for every call to Resume. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Resume answers no error, so the
// double is built with no fault arm and a fault armed through Faults fails
// the test.
func (s *ProducerResumeStub) Returns(result Cursor) *ProducerResumeStub {
	s.fallback = &ProducerResumeReturn{Result: result}
	return s
}

// Func supplies a body for Resume, for when the answer depends on the
// arguments. Nothing takes precedence over it: Resume answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *ProducerResumeStub) Func(fn func(string) Cursor) *ProducerResumeStub {
	s.fn = fn
	return s
}

// ProducerWatchStub controls how the double answers Watch and records what
// it was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type ProducerWatchStub struct {
	*stub.MethodStub[ProducerWatchCall]

	fn       func() Subscription
	fallback *ProducerWatchReturn
}

// ProducerWatchReturn holds the fixed answer configured through Returns.
type ProducerWatchReturn struct {
	Result Subscription
}

// Returns pins a fixed result for every call to Watch. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Watch answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *ProducerWatchStub) Returns(result Subscription) *ProducerWatchStub {
	s.fallback = &ProducerWatchReturn{Result: result}
	return s
}

// Func supplies a body for Watch, for when the answer depends on the
// arguments. Nothing takes precedence over it: Watch answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *ProducerWatchStub) Func(fn func() Subscription) *ProducerWatchStub {
	s.fn = fn
	return s
}

// ProducerNothingStub controls how the double answers Nothing and records
// what it was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type ProducerNothingStub struct {
	*stub.MethodStub[ProducerNothingCall]

	fn       func() Empty
	fallback *ProducerNothingReturn
}

// ProducerNothingReturn holds the fixed answer configured through Returns.
type ProducerNothingReturn struct {
	Result Empty
}

// Returns pins a fixed result for every call to Nothing. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Nothing answers no error, so the
// double is built with no fault arm and a fault armed through Faults fails
// the test.
func (s *ProducerNothingStub) Returns(result Empty) *ProducerNothingStub {
	s.fallback = &ProducerNothingReturn{Result: result}
	return s
}

// Func supplies a body for Nothing, for when the answer depends on the
// arguments. Nothing takes precedence over it: Nothing answers no error,
// so the double carries no fault arm for one to be reported through.
func (s *ProducerNothingStub) Func(fn func() Empty) *ProducerNothingStub {
	s.fn = fn
	return s
}

// ProducerElsewhereStub controls how the double answers Elsewhere and
// records what it was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type ProducerElsewhereStub struct {
	*stub.MethodStub[ProducerElsewhereCall]

	fn       func() other.Thing
	fallback *ProducerElsewhereReturn
}

// ProducerElsewhereReturn holds the fixed answer configured through
// Returns.
type ProducerElsewhereReturn struct {
	Result other.Thing
}

// Returns pins a fixed result for every call to Elsewhere. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Elsewhere answers no error, so the
// double is built with no fault arm and a fault armed through Faults fails
// the test.
func (s *ProducerElsewhereStub) Returns(result other.Thing) *ProducerElsewhereStub {
	s.fallback = &ProducerElsewhereReturn{Result: result}
	return s
}

// Func supplies a body for Elsewhere, for when the answer depends on the
// arguments. Nothing takes precedence over it: Elsewhere answers no error,
// so the double carries no fault arm for one to be reported through.
func (s *ProducerElsewhereStub) Func(fn func() other.Thing) *ProducerElsewhereStub {
	s.fn = fn
	return s
}

// --- ProducerStub ---

// The double-runtime witness for the surface this double rides.
var _ = stub.CompatV1

// ProducerStubOption configures a [ProducerStub] at construction time.
type ProducerStubOption func(*ProducerStub)

// ProducerStubStrict makes an unconfigured method fail the test rather
// than return its zero value.
func ProducerStubStrict() ProducerStubOption {
	return func(s *ProducerStub) { s.strict = true }
}

// ProducerStubDelegateTo wraps a real implementation: every method
// forwards to impl and is recorded on the double on the way through.
func ProducerStubDelegateTo(impl Producer) ProducerStubOption {
	return func(s *ProducerStub) {
		s.OnOpen.Func(impl.Open)
		s.OnResume.Func(impl.Resume)
		s.OnWatch.Func(impl.Watch)
		s.OnNothing.Func(impl.Nothing)
		s.OnElsewhere.Func(impl.Elsewhere)
	}
}

// ProducerStubWithClock drives latency and time-windowed faults from clk
// rather than from wall time, so a test asserting on a five-second timeout
// does not take five seconds.
func ProducerStubWithClock(clk assert.Clock) ProducerStubOption {
	return func(s *ProducerStub) { s.group.SetClock(clk) }
}

// ProducerStubWithRandSource makes probabilistic fault injection
// reproducible. A failure nobody can replay is not a finding.
func ProducerStubWithRandSource(src rand.Source) ProducerStubOption {
	return func(s *ProducerStub) { s.group.SetRandSource(src) }
}

// ProducerStubBenchMode disables call recording. Dispatch still works;
// only the accumulating call log is dropped.
func ProducerStubBenchMode() ProducerStubOption {
	return func(s *ProducerStub) { s.group.BenchMode() }
}

// ProducerStubWithLatency makes every method sleep d before answering, for
// a test standing in for a slow dependency rather than one slow method.
// Pair it with ProducerStubWithClock, or d elapses in wall time.
func ProducerStubWithLatency(d time.Duration) ProducerStubOption {
	return func(s *ProducerStub) { s.group.SetLatency(d) }
}

// ProducerStubFailing arms every method that can fail with err, so a
// caller's error path is one option rather than one Faults call per
// method. A method answering no error is left alone.
func ProducerStubFailing(err error) ProducerStubOption {
	return func(s *ProducerStub) {
		s.OnOpen.Faults(err, 1)
	}
}

// WithProducerOpen sets Open's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithProducerOpen(fn func(context.Context) (Cursor, error)) ProducerStubOption {
	return func(s *ProducerStub) { s.OnOpen.Func(fn) }
}

// WithProducerResume sets Resume's body at construction time, for the
// common case of configuring one method and taking the defaults for the
// rest.
func WithProducerResume(fn func(string) Cursor) ProducerStubOption {
	return func(s *ProducerStub) { s.OnResume.Func(fn) }
}

// WithProducerWatch sets Watch's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithProducerWatch(fn func() Subscription) ProducerStubOption {
	return func(s *ProducerStub) { s.OnWatch.Func(fn) }
}

// WithProducerNothing sets Nothing's body at construction time, for the
// common case of configuring one method and taking the defaults for the
// rest.
func WithProducerNothing(fn func() Empty) ProducerStubOption {
	return func(s *ProducerStub) { s.OnNothing.Func(fn) }
}

// WithProducerElsewhere sets Elsewhere's body at construction time, for
// the common case of configuring one method and taking the defaults for
// the rest.
func WithProducerElsewhere(fn func() other.Thing) ProducerStubOption {
	return func(s *ProducerStub) { s.OnElsewhere.Func(fn) }
}

// ProducerStub is a recording test double for Producer.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method returns its zero value and records the call.
type ProducerStub struct {
	OnOpen      *ProducerOpenStub
	OnResume    *ProducerResumeStub
	OnWatch     *ProducerWatchStub
	OnNothing   *ProducerNothingStub
	OnElsewhere *ProducerElsewhereStub

	// group is every method stub above, viewed through the surface that
	// does not depend on a signature.
	group  *stub.Group
	strict bool
}

// Compile-time proof that ProducerStub satisfies Producer.
var _ Producer = (*ProducerStub)(nil)

// NewProducerStub returns a double bound to tb.
//
// Passing tb registers a cleanup that verifies every method's call-count
// expectations when the test ends, so an unmet Times is reported without
// the caller remembering to check. A nil tb skips that, which is what
// benchmarks and non-test callers want.
func NewProducerStub(tb assert.TB, opts ...ProducerStubOption) *ProducerStub {
	s := &ProducerStub{
		OnOpen:      &ProducerOpenStub{MethodStub: stub.NewMethodStub[ProducerOpenCall](tb, "Producer.Open")},
		OnResume:    &ProducerResumeStub{MethodStub: stub.NewMethodStub[ProducerResumeCall](tb, "Producer.Resume")},
		OnWatch:     &ProducerWatchStub{MethodStub: stub.NewMethodStub[ProducerWatchCall](tb, "Producer.Watch")},
		OnNothing:   &ProducerNothingStub{MethodStub: stub.NewMethodStub[ProducerNothingCall](tb, "Producer.Nothing")},
		OnElsewhere: &ProducerElsewhereStub{MethodStub: stub.NewMethodStub[ProducerElsewhereCall](tb, "Producer.Elsewhere")},
	}
	s.group = stub.NewGroup(
		s.OnOpen.MethodStub,
		s.OnResume.MethodStub,
		s.OnWatch.MethodStub,
		s.OnNothing.MethodStub,
		s.OnElsewhere.MethodStub,
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

// NewProducerStubDelayed returns a double whose every method sleeps d
// against clk before answering — the clock and the latency a
// slow-dependency test wants, in one call rather than two options. Further
// options apply after the pair, so a caller can still add Strict or a
// method override.
func NewProducerStubDelayed(tb assert.TB, d time.Duration, clk assert.Clock, opts ...ProducerStubOption) *ProducerStub {
	return NewProducerStub(tb, append([]ProducerStubOption{
		ProducerStubWithClock(clk),
		ProducerStubWithLatency(d),
	}, opts...)...)
}

// NewProducerStubFailing returns a double whose every fallible method
// answers err — [ProducerStubFailing] applied at construction, for the
// test whose whole subject is the caller's error path. Further options
// apply after it.
func NewProducerStubFailing(tb assert.TB, err error, opts ...ProducerStubOption) *ProducerStub {
	return NewProducerStub(tb, append([]ProducerStubOption{ProducerStubFailing(err)}, opts...)...)
}

// ResetCalls clears recorded calls, fault counters, and call-count
// expectations. Func and Returns configuration is deliberately preserved, so
// a test can reuse one double across phases without rebuilding it.
func (s *ProducerStub) ResetCalls() {
	s.group.Reset()
}

// --- Interface methods ---

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *ProducerOpenStub) invoke(ctx context.Context) func() ProducerOpenReturn {
	if s.fn == nil {
		return nil
	}
	return func() ProducerOpenReturn {
		r0, r1 := s.fn(ctx)
		return ProducerOpenReturn{Result: r0, Err: r1}
	}
}

// Open records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *ProducerStub) Open(ctx context.Context) (Cursor, error) {
	call := ProducerOpenCall{Ctx: ctx}
	r := stub.Answer(s.OnOpen.MethodStub, &call, stub.Arms[ProducerOpenCall, ProducerOpenReturn]{
		Invoke:   s.OnOpen.invoke(ctx),
		Fallback: s.OnOpen.fallback,
		Fault:    func(err error) ProducerOpenReturn { return ProducerOpenReturn{Err: err} },
		Stamp: func(c *ProducerOpenCall, r ProducerOpenReturn) {
			c.Result = r.Result
			c.Err = r.Err
		},
	})
	return r.Result, r.Err
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *ProducerResumeStub) invoke(token string) func() ProducerResumeReturn {
	if s.fn == nil {
		return nil
	}
	return func() ProducerResumeReturn {
		r0 := s.fn(token)
		return ProducerResumeReturn{Result: r0}
	}
}

// Resume records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *ProducerStub) Resume(token string) Cursor {
	call := ProducerResumeCall{Token: token}
	r := stub.Answer(s.OnResume.MethodStub, &call, stub.Arms[ProducerResumeCall, ProducerResumeReturn]{
		Invoke:   s.OnResume.invoke(token),
		Fallback: s.OnResume.fallback,
		Stamp: func(c *ProducerResumeCall, r ProducerResumeReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *ProducerWatchStub) invoke() func() ProducerWatchReturn {
	if s.fn == nil {
		return nil
	}
	return func() ProducerWatchReturn {
		r0 := s.fn()
		return ProducerWatchReturn{Result: r0}
	}
}

// Watch records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *ProducerStub) Watch() Subscription {
	call := ProducerWatchCall{}
	r := stub.Answer(s.OnWatch.MethodStub, &call, stub.Arms[ProducerWatchCall, ProducerWatchReturn]{
		Invoke:   s.OnWatch.invoke(),
		Fallback: s.OnWatch.fallback,
		Stamp: func(c *ProducerWatchCall, r ProducerWatchReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *ProducerNothingStub) invoke() func() ProducerNothingReturn {
	if s.fn == nil {
		return nil
	}
	return func() ProducerNothingReturn {
		r0 := s.fn()
		return ProducerNothingReturn{Result: r0}
	}
}

// Nothing records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *ProducerStub) Nothing() Empty {
	call := ProducerNothingCall{}
	r := stub.Answer(s.OnNothing.MethodStub, &call, stub.Arms[ProducerNothingCall, ProducerNothingReturn]{
		Invoke:   s.OnNothing.invoke(),
		Fallback: s.OnNothing.fallback,
		Stamp: func(c *ProducerNothingCall, r ProducerNothingReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *ProducerElsewhereStub) invoke() func() ProducerElsewhereReturn {
	if s.fn == nil {
		return nil
	}
	return func() ProducerElsewhereReturn {
		r0 := s.fn()
		return ProducerElsewhereReturn{Result: r0}
	}
}

// Elsewhere records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *ProducerStub) Elsewhere() other.Thing {
	call := ProducerElsewhereCall{}
	r := stub.Answer(s.OnElsewhere.MethodStub, &call, stub.Arms[ProducerElsewhereCall, ProducerElsewhereReturn]{
		Invoke:   s.OnElsewhere.invoke(),
		Fallback: s.OnElsewhere.fallback,
		Stamp: func(c *ProducerElsewhereCall, r ProducerElsewhereReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// --- Cursor, the handle Producer's methods answer ---

// CursorNextCall records one invocation of Cursor.Next.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type CursorNextCall struct {
	Ctx     context.Context
	Result0 string
	Result1 bool
	Err     error
}

// CursorCloseCall records one invocation of Cursor.Close.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type CursorCloseCall struct {
	Ctx context.Context
	Err error
}

// --- Per-method configuration ---

// CursorNextStub controls how the double answers Next and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type CursorNextStub struct {
	*stub.MethodStub[CursorNextCall]

	fn       func(context.Context) (string, bool, error)
	fallback *CursorNextReturn
}

// CursorNextReturn holds the fixed answer configured through Returns.
type CursorNextReturn struct {
	Result0 string
	Result1 bool
	Err     error
}

// Returns pins a fixed result for every call to Next. A Func override
// takes precedence over it, and so does an injected fault.
func (s *CursorNextStub) Returns(result0 string, result1 bool, err error) *CursorNextStub {
	s.fallback = &CursorNextReturn{Result0: result0, Result1: result1, Err: err}
	return s
}

// Func supplies a body for Next, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *CursorNextStub) Func(fn func(context.Context) (string, bool, error)) *CursorNextStub {
	s.fn = fn
	return s
}

// CursorCloseStub controls how the double answers Close and records what
// it was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type CursorCloseStub struct {
	*stub.MethodStub[CursorCloseCall]

	fn       func(context.Context) error
	fallback *CursorCloseReturn
}

// CursorCloseReturn holds the fixed answer configured through Returns.
type CursorCloseReturn struct {
	Err error
}

// Returns pins a fixed result for every call to Close. A Func override
// takes precedence over it, and so does an injected fault.
func (s *CursorCloseStub) Returns(err error) *CursorCloseStub {
	s.fallback = &CursorCloseReturn{Err: err}
	return s
}

// Func supplies a body for Close, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *CursorCloseStub) Func(fn func(context.Context) error) *CursorCloseStub {
	s.fn = fn
	return s
}

// --- CursorStub ---

// The double-runtime witness for the surface this double rides.
var _ = stub.CompatV1

// CursorStubOption configures a [CursorStub] at construction time.
type CursorStubOption func(*CursorStub)

// CursorStubStrict makes an unconfigured method fail the test rather than
// return its zero value.
func CursorStubStrict() CursorStubOption {
	return func(s *CursorStub) { s.strict = true }
}

// CursorStubDelegateTo wraps a real implementation: every method forwards
// to impl and is recorded on the double on the way through.
func CursorStubDelegateTo(impl Cursor) CursorStubOption {
	return func(s *CursorStub) {
		s.OnNext.Func(impl.Next)
		s.OnClose.Func(impl.Close)
	}
}

// CursorStubWithClock drives latency and time-windowed faults from clk
// rather than from wall time, so a test asserting on a five-second timeout
// does not take five seconds.
func CursorStubWithClock(clk assert.Clock) CursorStubOption {
	return func(s *CursorStub) { s.group.SetClock(clk) }
}

// CursorStubWithRandSource makes probabilistic fault injection
// reproducible. A failure nobody can replay is not a finding.
func CursorStubWithRandSource(src rand.Source) CursorStubOption {
	return func(s *CursorStub) { s.group.SetRandSource(src) }
}

// CursorStubBenchMode disables call recording. Dispatch still works; only
// the accumulating call log is dropped.
func CursorStubBenchMode() CursorStubOption {
	return func(s *CursorStub) { s.group.BenchMode() }
}

// CursorStubWithLatency makes every method sleep d before answering, for a
// test standing in for a slow dependency rather than one slow method. Pair
// it with CursorStubWithClock, or d elapses in wall time.
func CursorStubWithLatency(d time.Duration) CursorStubOption {
	return func(s *CursorStub) { s.group.SetLatency(d) }
}

// CursorStubFailing arms every method that can fail with err, so a
// caller's error path is one option rather than one Faults call per
// method. A method answering no error is left alone.
func CursorStubFailing(err error) CursorStubOption {
	return func(s *CursorStub) {
		s.OnNext.Faults(err, 1)
		s.OnClose.Faults(err, 1)
	}
}

// WithCursorNext sets Next's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithCursorNext(fn func(context.Context) (string, bool, error)) CursorStubOption {
	return func(s *CursorStub) { s.OnNext.Func(fn) }
}

// WithCursorClose sets Close's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithCursorClose(fn func(context.Context) error) CursorStubOption {
	return func(s *CursorStub) { s.OnClose.Func(fn) }
}

// CursorStub is a recording test double for Cursor.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method returns its zero value and records the call.
type CursorStub struct {
	OnNext  *CursorNextStub
	OnClose *CursorCloseStub

	// group is every method stub above, viewed through the surface that
	// does not depend on a signature.
	group  *stub.Group
	strict bool
}

// Compile-time proof that CursorStub satisfies Cursor.
var _ Cursor = (*CursorStub)(nil)

// NewCursorStub returns a double bound to tb.
//
// Passing tb registers a cleanup that verifies every method's call-count
// expectations when the test ends, so an unmet Times is reported without
// the caller remembering to check. A nil tb skips that, which is what
// benchmarks and non-test callers want.
func NewCursorStub(tb assert.TB, opts ...CursorStubOption) *CursorStub {
	s := &CursorStub{
		OnNext:  &CursorNextStub{MethodStub: stub.NewMethodStub[CursorNextCall](tb, "Cursor.Next")},
		OnClose: &CursorCloseStub{MethodStub: stub.NewMethodStub[CursorCloseCall](tb, "Cursor.Close")},
	}
	s.group = stub.NewGroup(
		s.OnNext.MethodStub,
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

// NewCursorStubDelayed returns a double whose every method sleeps d
// against clk before answering — the clock and the latency a
// slow-dependency test wants, in one call rather than two options. Further
// options apply after the pair, so a caller can still add Strict or a
// method override.
func NewCursorStubDelayed(tb assert.TB, d time.Duration, clk assert.Clock, opts ...CursorStubOption) *CursorStub {
	return NewCursorStub(tb, append([]CursorStubOption{
		CursorStubWithClock(clk),
		CursorStubWithLatency(d),
	}, opts...)...)
}

// NewCursorStubFailing returns a double whose every fallible method
// answers err — [CursorStubFailing] applied at construction, for the test
// whose whole subject is the caller's error path. Further options apply
// after it.
func NewCursorStubFailing(tb assert.TB, err error, opts ...CursorStubOption) *CursorStub {
	return NewCursorStub(tb, append([]CursorStubOption{CursorStubFailing(err)}, opts...)...)
}

// ResetCalls clears recorded calls, fault counters, and call-count
// expectations. Func and Returns configuration is deliberately preserved, so
// a test can reuse one double across phases without rebuilding it.
func (s *CursorStub) ResetCalls() {
	s.group.Reset()
}

// --- Interface methods ---

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *CursorNextStub) invoke(ctx context.Context) func() CursorNextReturn {
	if s.fn == nil {
		return nil
	}
	return func() CursorNextReturn {
		r0, r1, r2 := s.fn(ctx)
		return CursorNextReturn{Result0: r0, Result1: r1, Err: r2}
	}
}

// Next records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *CursorStub) Next(ctx context.Context) (string, bool, error) {
	call := CursorNextCall{Ctx: ctx}
	r := stub.Answer(s.OnNext.MethodStub, &call, stub.Arms[CursorNextCall, CursorNextReturn]{
		Invoke:   s.OnNext.invoke(ctx),
		Fallback: s.OnNext.fallback,
		Fault:    func(err error) CursorNextReturn { return CursorNextReturn{Err: err} },
		Stamp: func(c *CursorNextCall, r CursorNextReturn) {
			c.Result0 = r.Result0
			c.Result1 = r.Result1
			c.Err = r.Err
		},
	})
	return r.Result0, r.Result1, r.Err
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *CursorCloseStub) invoke(ctx context.Context) func() CursorCloseReturn {
	if s.fn == nil {
		return nil
	}
	return func() CursorCloseReturn {
		r0 := s.fn(ctx)
		return CursorCloseReturn{Err: r0}
	}
}

// Close records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *CursorStub) Close(ctx context.Context) error {
	call := CursorCloseCall{Ctx: ctx}
	r := stub.Answer(s.OnClose.MethodStub, &call, stub.Arms[CursorCloseCall, CursorCloseReturn]{
		Invoke:   s.OnClose.invoke(ctx),
		Fallback: s.OnClose.fallback,
		Fault:    func(err error) CursorCloseReturn { return CursorCloseReturn{Err: err} },
		Stamp: func(c *CursorCloseCall, r CursorCloseReturn) {
			c.Err = r.Err
		},
	})
	return r.Err
}

// ReaderReadCall records one invocation of Reader.Read.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type ReaderReadCall struct {
	Result Cursor
}

// --- Per-method configuration ---

// ReaderReadStub controls how the double answers Read and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type ReaderReadStub struct {
	*stub.MethodStub[ReaderReadCall]

	fn       func() Cursor
	fallback *ReaderReadReturn
}

// ReaderReadReturn holds the fixed answer configured through Returns.
type ReaderReadReturn struct {
	Result Cursor
}

// Returns pins a fixed result for every call to Read. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Read answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *ReaderReadStub) Returns(result Cursor) *ReaderReadStub {
	s.fallback = &ReaderReadReturn{Result: result}
	return s
}

// Func supplies a body for Read, for when the answer depends on the
// arguments. Nothing takes precedence over it: Read answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *ReaderReadStub) Func(fn func() Cursor) *ReaderReadStub {
	s.fn = fn
	return s
}

// --- ReaderStub ---

// The double-runtime witness for the surface this double rides.
var _ = stub.CompatV1

// ReaderStubOption configures a [ReaderStub] at construction time.
type ReaderStubOption func(*ReaderStub)

// ReaderStubStrict makes an unconfigured method fail the test rather than
// return its zero value.
func ReaderStubStrict() ReaderStubOption {
	return func(s *ReaderStub) { s.strict = true }
}

// ReaderStubDelegateTo wraps a real implementation: every method forwards
// to impl and is recorded on the double on the way through.
func ReaderStubDelegateTo(impl Reader) ReaderStubOption {
	return func(s *ReaderStub) {
		s.OnRead.Func(impl.Read)
	}
}

// ReaderStubWithClock drives latency and time-windowed faults from clk
// rather than from wall time, so a test asserting on a five-second timeout
// does not take five seconds.
func ReaderStubWithClock(clk assert.Clock) ReaderStubOption {
	return func(s *ReaderStub) { s.group.SetClock(clk) }
}

// ReaderStubWithRandSource makes probabilistic fault injection
// reproducible. A failure nobody can replay is not a finding.
func ReaderStubWithRandSource(src rand.Source) ReaderStubOption {
	return func(s *ReaderStub) { s.group.SetRandSource(src) }
}

// ReaderStubBenchMode disables call recording. Dispatch still works; only
// the accumulating call log is dropped.
func ReaderStubBenchMode() ReaderStubOption {
	return func(s *ReaderStub) { s.group.BenchMode() }
}

// ReaderStubWithLatency makes every method sleep d before answering, for a
// test standing in for a slow dependency rather than one slow method. Pair
// it with ReaderStubWithClock, or d elapses in wall time.
func ReaderStubWithLatency(d time.Duration) ReaderStubOption {
	return func(s *ReaderStub) { s.group.SetLatency(d) }
}

// WithReaderRead sets Read's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithReaderRead(fn func() Cursor) ReaderStubOption {
	return func(s *ReaderStub) { s.OnRead.Func(fn) }
}

// ReaderStub is a recording test double for Reader.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method returns its zero value and records the call.
type ReaderStub struct {
	OnRead *ReaderReadStub

	// group is every method stub above, viewed through the surface that
	// does not depend on a signature.
	group  *stub.Group
	strict bool
}

// Compile-time proof that ReaderStub satisfies Reader.
var _ Reader = (*ReaderStub)(nil)

// NewReaderStub returns a double bound to tb.
//
// Passing tb registers a cleanup that verifies every method's call-count
// expectations when the test ends, so an unmet Times is reported without
// the caller remembering to check. A nil tb skips that, which is what
// benchmarks and non-test callers want.
func NewReaderStub(tb assert.TB, opts ...ReaderStubOption) *ReaderStub {
	s := &ReaderStub{
		OnRead: &ReaderReadStub{MethodStub: stub.NewMethodStub[ReaderReadCall](tb, "Reader.Read")},
	}
	s.group = stub.NewGroup(
		s.OnRead.MethodStub,
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

// NewReaderStubDelayed returns a double whose every method sleeps d
// against clk before answering — the clock and the latency a
// slow-dependency test wants, in one call rather than two options. Further
// options apply after the pair, so a caller can still add Strict or a
// method override.
func NewReaderStubDelayed(tb assert.TB, d time.Duration, clk assert.Clock, opts ...ReaderStubOption) *ReaderStub {
	return NewReaderStub(tb, append([]ReaderStubOption{
		ReaderStubWithClock(clk),
		ReaderStubWithLatency(d),
	}, opts...)...)
}

// ResetCalls clears recorded calls, fault counters, and call-count
// expectations. Func and Returns configuration is deliberately preserved, so
// a test can reuse one double across phases without rebuilding it.
func (s *ReaderStub) ResetCalls() {
	s.group.Reset()
}

// --- Interface methods ---

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *ReaderReadStub) invoke() func() ReaderReadReturn {
	if s.fn == nil {
		return nil
	}
	return func() ReaderReadReturn {
		r0 := s.fn()
		return ReaderReadReturn{Result: r0}
	}
}

// Read records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *ReaderStub) Read() Cursor {
	call := ReaderReadCall{}
	r := stub.Answer(s.OnRead.MethodStub, &call, stub.Arms[ReaderReadCall, ReaderReadReturn]{
		Invoke:   s.OnRead.invoke(),
		Fallback: s.OnRead.fallback,
		Stamp: func(c *ReaderReadCall, r ReaderReadReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// SubscriptionCancelCall records one invocation of Subscription.Cancel.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type SubscriptionCancelCall struct {
}

// --- Per-method configuration ---

// SubscriptionCancelStub controls how the double answers Cancel and
// records what it was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type SubscriptionCancelStub struct {
	*stub.MethodStub[SubscriptionCancelCall]

	fn func()
}

// Func supplies a body for Cancel, for when the answer depends on the
// arguments. Nothing takes precedence over it: Cancel answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *SubscriptionCancelStub) Func(fn func()) *SubscriptionCancelStub {
	s.fn = fn
	return s
}

// --- SubscriptionStub ---

// The double-runtime witness for the surface this double rides.
var _ = stub.CompatV1

// SubscriptionStubOption configures a [SubscriptionStub] at construction
// time.
type SubscriptionStubOption func(*SubscriptionStub)

// SubscriptionStubStrict makes an unconfigured method fail the test rather
// than return its zero value.
func SubscriptionStubStrict() SubscriptionStubOption {
	return func(s *SubscriptionStub) { s.strict = true }
}

// SubscriptionStubDelegateTo wraps a real implementation: every method
// forwards to impl and is recorded on the double on the way through.
func SubscriptionStubDelegateTo(impl Subscription) SubscriptionStubOption {
	return func(s *SubscriptionStub) {
		s.OnCancel.Func(impl.Cancel)
	}
}

// SubscriptionStubWithClock drives latency and time-windowed faults from
// clk rather than from wall time, so a test asserting on a five-second
// timeout does not take five seconds.
func SubscriptionStubWithClock(clk assert.Clock) SubscriptionStubOption {
	return func(s *SubscriptionStub) { s.group.SetClock(clk) }
}

// SubscriptionStubWithRandSource makes probabilistic fault injection
// reproducible. A failure nobody can replay is not a finding.
func SubscriptionStubWithRandSource(src rand.Source) SubscriptionStubOption {
	return func(s *SubscriptionStub) { s.group.SetRandSource(src) }
}

// SubscriptionStubBenchMode disables call recording. Dispatch still works;
// only the accumulating call log is dropped.
func SubscriptionStubBenchMode() SubscriptionStubOption {
	return func(s *SubscriptionStub) { s.group.BenchMode() }
}

// SubscriptionStubWithLatency makes every method sleep d before answering,
// for a test standing in for a slow dependency rather than one slow
// method. Pair it with SubscriptionStubWithClock, or d elapses in wall
// time.
func SubscriptionStubWithLatency(d time.Duration) SubscriptionStubOption {
	return func(s *SubscriptionStub) { s.group.SetLatency(d) }
}

// WithSubscriptionCancel sets Cancel's body at construction time, for the
// common case of configuring one method and taking the defaults for the
// rest.
func WithSubscriptionCancel(fn func()) SubscriptionStubOption {
	return func(s *SubscriptionStub) { s.OnCancel.Func(fn) }
}

// SubscriptionStub is a recording test double for Subscription.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method returns its zero value and records the call.
type SubscriptionStub struct {
	OnCancel *SubscriptionCancelStub

	// group is every method stub above, viewed through the surface that
	// does not depend on a signature.
	group  *stub.Group
	strict bool
}

// Compile-time proof that SubscriptionStub satisfies Subscription.
var _ Subscription = (*SubscriptionStub)(nil)

// NewSubscriptionStub returns a double bound to tb.
//
// Passing tb registers a cleanup that verifies every method's call-count
// expectations when the test ends, so an unmet Times is reported without
// the caller remembering to check. A nil tb skips that, which is what
// benchmarks and non-test callers want.
func NewSubscriptionStub(tb assert.TB, opts ...SubscriptionStubOption) *SubscriptionStub {
	s := &SubscriptionStub{
		OnCancel: &SubscriptionCancelStub{MethodStub: stub.NewMethodStub[SubscriptionCancelCall](tb, "Subscription.Cancel")},
	}
	s.group = stub.NewGroup(
		s.OnCancel.MethodStub,
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

// NewSubscriptionStubDelayed returns a double whose every method sleeps d
// against clk before answering — the clock and the latency a
// slow-dependency test wants, in one call rather than two options. Further
// options apply after the pair, so a caller can still add Strict or a
// method override.
func NewSubscriptionStubDelayed(tb assert.TB, d time.Duration, clk assert.Clock, opts ...SubscriptionStubOption) *SubscriptionStub {
	return NewSubscriptionStub(tb, append([]SubscriptionStubOption{
		SubscriptionStubWithClock(clk),
		SubscriptionStubWithLatency(d),
	}, opts...)...)
}

// ResetCalls clears recorded calls, fault counters, and call-count
// expectations. Func and Returns configuration is deliberately preserved, so
// a test can reuse one double across phases without rebuilding it.
func (s *SubscriptionStub) ResetCalls() {
	s.group.Reset()
}

// --- Interface methods ---

// invokeVoid adapts the Func override to the shape [stub.Answer] consumes.
// A method returning nothing still has an override worth dispatching to and
// a call worth recording; the empty struct is the return tuple of no returns.
func (s *SubscriptionCancelStub) invokeVoid() func() struct{} {
	if s.fn == nil {
		return nil
	}
	return func() struct{} {
		s.fn()
		return struct{}{}
	}
}

// Cancel records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *SubscriptionStub) Cancel() {
	call := SubscriptionCancelCall{}
	stub.Answer(s.OnCancel.MethodStub, &call, stub.Arms[SubscriptionCancelCall, struct{}]{
		Invoke: s.OnCancel.invokeVoid(),
	})
}
