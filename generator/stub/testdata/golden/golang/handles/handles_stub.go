package handles

import (
	"context"
	rand "math/rand/v2"
	"time"

	"example.com/shop/other"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// ArchiveSizeCall records one call of Archive.Size.
type ArchiveSizeCall struct {
	Result int
}

// ArchiveOpenLogCall records one call of Archive.OpenLog.
type ArchiveOpenLogCall struct {
	Name   string
	Result Log
	Err    error
}

// ArchiveSizeStub configures Size and records its calls. Size returns no
// error, so a fault armed through it fails the test.
type ArchiveSizeStub struct {
	*stub.MethodStub[ArchiveSizeCall]

	fn       func() int
	fallback *ArchiveSizeReturn
}

// ArchiveSizeReturn contains the fixed results of Size.
type ArchiveSizeReturn struct {
	Result int
}

// Returns sets the results of every call of Size. A function takes
// precedence over them.
func (s *ArchiveSizeStub) Returns(result int) *ArchiveSizeStub {
	s.fallback = &ArchiveSizeReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Size from its
// arguments.
func (s *ArchiveSizeStub) Func(fn func() int) *ArchiveSizeStub {
	s.fn = fn
	return s
}

// ArchiveOpenLogStub configures OpenLog and records its calls. A fault
// armed through it takes precedence over the function and the fixed
// results.
type ArchiveOpenLogStub struct {
	*stub.MethodStub[ArchiveOpenLogCall]

	fn       func(string) (Log, error)
	fallback *ArchiveOpenLogReturn
}

// ArchiveOpenLogReturn contains the fixed results of OpenLog.
type ArchiveOpenLogReturn struct {
	Result Log
	Err    error
}

// Returns sets the results of every call of OpenLog. A function and an
// injected fault take precedence over them.
func (s *ArchiveOpenLogStub) Returns(result Log, err error) *ArchiveOpenLogStub {
	s.fallback = &ArchiveOpenLogReturn{Result: result, Err: err}
	return s
}

// Func sets the function that computes the results of OpenLog from its
// arguments. An injected fault takes precedence over it.
func (s *ArchiveOpenLogStub) Func(fn func(string) (Log, error)) *ArchiveOpenLogStub {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// ArchiveStubOption configures a [ArchiveStub] at its construction.
type ArchiveStubOption func(*ArchiveStub)

// ArchiveStubStrict returns the option that makes an unconfigured method
// fail the test instead of returning its zero value.
func ArchiveStubStrict() ArchiveStubOption {
	return func(s *ArchiveStub) { s.strict = true }
}

// ArchiveStubDelegateTo returns the option that passes each call to impl.
// The double still records the call.
func ArchiveStubDelegateTo(impl Archive) ArchiveStubOption {
	return func(s *ArchiveStub) {
		s.OnSize.Func(impl.Size)
		s.OnOpenLog.Func(impl.OpenLog)
	}
}

// ArchiveStubWithClock returns the option that sets the clock of the
// double to clk. The double takes the latency, the windows of the faults
// and the instants of the calls from clk. A latency on a controlled clock
// costs no real time.
func ArchiveStubWithClock(clk assert.Clock) ArchiveStubOption {
	return func(s *ArchiveStub) { s.group.SetClock(clk) }
}

// ArchiveStubWithRandSource returns the option that sets the random source
// of the probabilistic faults to src. A seeded source repeats the faults
// of a run.
func ArchiveStubWithRandSource(src rand.Source) ArchiveStubOption {
	return func(s *ArchiveStub) { s.group.SetRandSource(src) }
}

// ArchiveStubBenchMode returns the option that turns the recording of
// calls off. The methods still return their results.
func ArchiveStubBenchMode() ArchiveStubOption {
	return func(s *ArchiveStub) { s.group.BenchMode() }
}

// ArchiveStubWithLatency returns the option that makes every method sleep
// d before it returns. Without ArchiveStubWithClock, d passes in real
// time.
func ArchiveStubWithLatency(d time.Duration) ArchiveStubOption {
	return func(s *ArchiveStub) { s.group.SetLatency(d) }
}

// ArchiveStubFailing returns the option that arms err on each method that
// can fail. Each call of such a method then fails with err.
func ArchiveStubFailing(err error) ArchiveStubOption {
	return func(s *ArchiveStub) {
		s.OnOpenLog.Faults(err, 1)
	}
}

// WithArchiveSize returns the option that sets the function of Size.
func WithArchiveSize(fn func() int) ArchiveStubOption {
	return func(s *ArchiveStub) { s.OnSize.Func(fn) }
}

// WithArchiveOpenLog returns the option that sets the function of OpenLog.
func WithArchiveOpenLog(fn func(string) (Log, error)) ArchiveStubOption {
	return func(s *ArchiveStub) { s.OnOpenLog.Func(fn) }
}

// ArchiveStub is a recording test double of Archive. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type ArchiveStub struct {
	OnSize *ArchiveSizeStub
	// OnOpenLog configures OpenLog, which Opener declares.
	OnOpenLog *ArchiveOpenLogStub

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// The compiler checks that ArchiveStub implements Archive.
var _ Archive = (*ArchiveStub)(nil)

// NewArchiveStub returns a double of Archive that reports its failures to
// tb and applies opts. The double checks the expected counts of calls at
// the end of the test. With a nil tb, the double reports nothing.
func NewArchiveStub(tb assert.TB, opts ...ArchiveStubOption) *ArchiveStub {
	s := &ArchiveStub{
		OnSize: &ArchiveSizeStub{
			MethodStub: stub.NewMethodStub[ArchiveSizeCall](tb, "Archive.Size"),
		},
		OnOpenLog: &ArchiveOpenLogStub{
			MethodStub: stub.NewMethodStub[ArchiveOpenLogCall](tb, "Archive.OpenLog"),
		},
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

// NewArchiveStubDelayed returns a double whose methods sleep d on clk
// before they return. It applies opts after the clock and the latency.
func NewArchiveStubDelayed(
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...ArchiveStubOption,
) *ArchiveStub {
	return NewArchiveStub(tb, append([]ArchiveStubOption{
		ArchiveStubWithClock(clk),
		ArchiveStubWithLatency(d),
	}, opts...)...)
}

// NewArchiveStubFailing returns a double on which each call of a method
// that can fail fails with err. It applies opts after the failure.
func NewArchiveStubFailing(tb assert.TB, err error, opts ...ArchiveStubOption) *ArchiveStub {
	return NewArchiveStub(tb, append([]ArchiveStubOption{ArchiveStubFailing(err)}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *ArchiveStub) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *ArchiveSizeStub) invoke() func() ArchiveSizeReturn {
	if s.fn == nil {
		return nil
	}
	return func() ArchiveSizeReturn {
		r0 := s.fn()
		return ArchiveSizeReturn{Result: r0}
	}
}

// Size records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *ArchiveStub) Size() int {
	call := ArchiveSizeCall{}
	r := stub.Answer(s.OnSize.MethodStub, call, stub.Arms[ArchiveSizeCall, ArchiveSizeReturn]{
		Invoke:   s.OnSize.invoke(),
		Fallback: s.OnSize.fallback,
		Stamp: func(c ArchiveSizeCall, r ArchiveSizeReturn) ArchiveSizeCall {
			c.Result = r.Result
			return c
		},
	})
	return r.Result
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *ArchiveOpenLogStub) invoke(name string) func() ArchiveOpenLogReturn {
	if s.fn == nil {
		return nil
	}
	return func() ArchiveOpenLogReturn {
		r0, r1 := s.fn(name)
		return ArchiveOpenLogReturn{Result: r0, Err: r1}
	}
}

// OpenLog records the call. An armed fault comes first, then the function
// of Func, then the results of Returns, and then the zero values.
func (s *ArchiveStub) OpenLog(name string) (Log, error) {
	call := ArchiveOpenLogCall{Name: name}
	r := stub.Answer(
		s.OnOpenLog.MethodStub,
		call,
		stub.Arms[ArchiveOpenLogCall, ArchiveOpenLogReturn]{
			Invoke:   s.OnOpenLog.invoke(name),
			Fallback: s.OnOpenLog.fallback,
			Fault:    func(err error) ArchiveOpenLogReturn { return ArchiveOpenLogReturn{Err: err} },
			Stamp: func(c ArchiveOpenLogCall, r ArchiveOpenLogReturn) ArchiveOpenLogCall {
				c.Result = r.Result
				c.Err = r.Err
				return c
			},
		},
	)
	return r.Result, r.Err
}

// Log is a handle that a method of Archive returns.

// LogAppendCall records one call of Log.Append.
type LogAppendCall struct {
	Line string
	Err  error
}

// LogAppendStub configures Append and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type LogAppendStub struct {
	*stub.MethodStub[LogAppendCall]

	fn       func(string) error
	fallback *LogAppendReturn
}

// LogAppendReturn contains the fixed results of Append.
type LogAppendReturn struct {
	Err error
}

// Returns sets the results of every call of Append. A function and an
// injected fault take precedence over them.
func (s *LogAppendStub) Returns(err error) *LogAppendStub {
	s.fallback = &LogAppendReturn{Err: err}
	return s
}

// Func sets the function that computes the results of Append from its
// arguments. An injected fault takes precedence over it.
func (s *LogAppendStub) Func(fn func(string) error) *LogAppendStub {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// LogStubOption configures a [LogStub] at its construction.
type LogStubOption func(*LogStub)

// LogStubStrict returns the option that makes an unconfigured method fail
// the test instead of returning its zero value.
func LogStubStrict() LogStubOption {
	return func(s *LogStub) { s.strict = true }
}

// LogStubDelegateTo returns the option that passes each call to impl. The
// double still records the call.
func LogStubDelegateTo(impl Log) LogStubOption {
	return func(s *LogStub) {
		s.OnAppend.Func(impl.Append)
	}
}

// LogStubWithClock returns the option that sets the clock of the double to
// clk. The double takes the latency, the windows of the faults and the
// instants of the calls from clk. A latency on a controlled clock costs no
// real time.
func LogStubWithClock(clk assert.Clock) LogStubOption {
	return func(s *LogStub) { s.group.SetClock(clk) }
}

// LogStubWithRandSource returns the option that sets the random source of
// the probabilistic faults to src. A seeded source repeats the faults of a
// run.
func LogStubWithRandSource(src rand.Source) LogStubOption {
	return func(s *LogStub) { s.group.SetRandSource(src) }
}

// LogStubBenchMode returns the option that turns the recording of calls
// off. The methods still return their results.
func LogStubBenchMode() LogStubOption {
	return func(s *LogStub) { s.group.BenchMode() }
}

// LogStubWithLatency returns the option that makes every method sleep d
// before it returns. Without LogStubWithClock, d passes in real time.
func LogStubWithLatency(d time.Duration) LogStubOption {
	return func(s *LogStub) { s.group.SetLatency(d) }
}

// LogStubFailing returns the option that arms err on each method that can
// fail. Each call of such a method then fails with err.
func LogStubFailing(err error) LogStubOption {
	return func(s *LogStub) {
		s.OnAppend.Faults(err, 1)
	}
}

// WithLogAppend returns the option that sets the function of Append.
func WithLogAppend(fn func(string) error) LogStubOption {
	return func(s *LogStub) { s.OnAppend.Func(fn) }
}

// LogStub is a recording test double of Log. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type LogStub struct {
	OnAppend *LogAppendStub

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// The compiler checks that LogStub implements Log.
var _ Log = (*LogStub)(nil)

// NewLogStub returns a double of Log that reports its failures to tb and
// applies opts. The double checks the expected counts of calls at the end
// of the test. With a nil tb, the double reports nothing.
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

// NewLogStubDelayed returns a double whose methods sleep d on clk before
// they return. It applies opts after the clock and the latency.
func NewLogStubDelayed(
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...LogStubOption,
) *LogStub {
	return NewLogStub(tb, append([]LogStubOption{
		LogStubWithClock(clk),
		LogStubWithLatency(d),
	}, opts...)...)
}

// NewLogStubFailing returns a double on which each call of a method that
// can fail fails with err. It applies opts after the failure.
func NewLogStubFailing(tb assert.TB, err error, opts ...LogStubOption) *LogStub {
	return NewLogStub(tb, append([]LogStubOption{LogStubFailing(err)}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *LogStub) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *LogAppendStub) invoke(line string) func() LogAppendReturn {
	if s.fn == nil {
		return nil
	}
	return func() LogAppendReturn {
		r0 := s.fn(line)
		return LogAppendReturn{Err: r0}
	}
}

// Append records the call. An armed fault comes first, then the function
// of Func, then the results of Returns, and then the zero values.
func (s *LogStub) Append(line string) error {
	call := LogAppendCall{Line: line}
	r := stub.Answer(s.OnAppend.MethodStub, call, stub.Arms[LogAppendCall, LogAppendReturn]{
		Invoke:   s.OnAppend.invoke(line),
		Fallback: s.OnAppend.fallback,
		Fault:    func(err error) LogAppendReturn { return LogAppendReturn{Err: err} },
		Stamp: func(c LogAppendCall, r LogAppendReturn) LogAppendCall {
			c.Err = r.Err
			return c
		},
	})
	return r.Err
}

// ProducerOpenCall records one call of Producer.Open.
type ProducerOpenCall struct {
	Ctx    context.Context
	Result Cursor
	Err    error
}

// ProducerResumeCall records one call of Producer.Resume.
type ProducerResumeCall struct {
	Token  string
	Result Cursor
}

// ProducerWatchCall records one call of Producer.Watch.
type ProducerWatchCall struct {
	Result Subscription
}

// ProducerNothingCall records one call of Producer.Nothing.
type ProducerNothingCall struct {
	Result Empty
}

// ProducerElsewhereCall records one call of Producer.Elsewhere.
type ProducerElsewhereCall struct {
	Result other.Thing
}

// ProducerOpenStub configures Open and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type ProducerOpenStub struct {
	*stub.MethodStub[ProducerOpenCall]

	fn       func(context.Context) (Cursor, error)
	fallback *ProducerOpenReturn
}

// ProducerOpenReturn contains the fixed results of Open.
type ProducerOpenReturn struct {
	Result Cursor
	Err    error
}

// Returns sets the results of every call of Open. A function and an
// injected fault take precedence over them.
func (s *ProducerOpenStub) Returns(result Cursor, err error) *ProducerOpenStub {
	s.fallback = &ProducerOpenReturn{Result: result, Err: err}
	return s
}

// Func sets the function that computes the results of Open from its
// arguments. An injected fault takes precedence over it.
func (s *ProducerOpenStub) Func(fn func(context.Context) (Cursor, error)) *ProducerOpenStub {
	s.fn = fn
	return s
}

// ProducerResumeStub configures Resume and records its calls. Resume
// returns no error, so a fault armed through it fails the test.
type ProducerResumeStub struct {
	*stub.MethodStub[ProducerResumeCall]

	fn       func(string) Cursor
	fallback *ProducerResumeReturn
}

// ProducerResumeReturn contains the fixed results of Resume.
type ProducerResumeReturn struct {
	Result Cursor
}

// Returns sets the results of every call of Resume. A function takes
// precedence over them.
func (s *ProducerResumeStub) Returns(result Cursor) *ProducerResumeStub {
	s.fallback = &ProducerResumeReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Resume from its
// arguments.
func (s *ProducerResumeStub) Func(fn func(string) Cursor) *ProducerResumeStub {
	s.fn = fn
	return s
}

// ProducerWatchStub configures Watch and records its calls. Watch returns
// no error, so a fault armed through it fails the test.
type ProducerWatchStub struct {
	*stub.MethodStub[ProducerWatchCall]

	fn       func() Subscription
	fallback *ProducerWatchReturn
}

// ProducerWatchReturn contains the fixed results of Watch.
type ProducerWatchReturn struct {
	Result Subscription
}

// Returns sets the results of every call of Watch. A function takes
// precedence over them.
func (s *ProducerWatchStub) Returns(result Subscription) *ProducerWatchStub {
	s.fallback = &ProducerWatchReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Watch from its
// arguments.
func (s *ProducerWatchStub) Func(fn func() Subscription) *ProducerWatchStub {
	s.fn = fn
	return s
}

// ProducerNothingStub configures Nothing and records its calls. Nothing
// returns no error, so a fault armed through it fails the test.
type ProducerNothingStub struct {
	*stub.MethodStub[ProducerNothingCall]

	fn       func() Empty
	fallback *ProducerNothingReturn
}

// ProducerNothingReturn contains the fixed results of Nothing.
type ProducerNothingReturn struct {
	Result Empty
}

// Returns sets the results of every call of Nothing. A function takes
// precedence over them.
func (s *ProducerNothingStub) Returns(result Empty) *ProducerNothingStub {
	s.fallback = &ProducerNothingReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Nothing from its
// arguments.
func (s *ProducerNothingStub) Func(fn func() Empty) *ProducerNothingStub {
	s.fn = fn
	return s
}

// ProducerElsewhereStub configures Elsewhere and records its calls.
// Elsewhere returns no error, so a fault armed through it fails the test.
type ProducerElsewhereStub struct {
	*stub.MethodStub[ProducerElsewhereCall]

	fn       func() other.Thing
	fallback *ProducerElsewhereReturn
}

// ProducerElsewhereReturn contains the fixed results of Elsewhere.
type ProducerElsewhereReturn struct {
	Result other.Thing
}

// Returns sets the results of every call of Elsewhere. A function takes
// precedence over them.
func (s *ProducerElsewhereStub) Returns(result other.Thing) *ProducerElsewhereStub {
	s.fallback = &ProducerElsewhereReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Elsewhere from its
// arguments.
func (s *ProducerElsewhereStub) Func(fn func() other.Thing) *ProducerElsewhereStub {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// ProducerStubOption configures a [ProducerStub] at its construction.
type ProducerStubOption func(*ProducerStub)

// ProducerStubStrict returns the option that makes an unconfigured method
// fail the test instead of returning its zero value.
func ProducerStubStrict() ProducerStubOption {
	return func(s *ProducerStub) { s.strict = true }
}

// ProducerStubDelegateTo returns the option that passes each call to impl.
// The double still records the call.
func ProducerStubDelegateTo(impl Producer) ProducerStubOption {
	return func(s *ProducerStub) {
		s.OnOpen.Func(impl.Open)
		s.OnResume.Func(impl.Resume)
		s.OnWatch.Func(impl.Watch)
		s.OnNothing.Func(impl.Nothing)
		s.OnElsewhere.Func(impl.Elsewhere)
	}
}

// ProducerStubWithClock returns the option that sets the clock of the
// double to clk. The double takes the latency, the windows of the faults
// and the instants of the calls from clk. A latency on a controlled clock
// costs no real time.
func ProducerStubWithClock(clk assert.Clock) ProducerStubOption {
	return func(s *ProducerStub) { s.group.SetClock(clk) }
}

// ProducerStubWithRandSource returns the option that sets the random
// source of the probabilistic faults to src. A seeded source repeats the
// faults of a run.
func ProducerStubWithRandSource(src rand.Source) ProducerStubOption {
	return func(s *ProducerStub) { s.group.SetRandSource(src) }
}

// ProducerStubBenchMode returns the option that turns the recording of
// calls off. The methods still return their results.
func ProducerStubBenchMode() ProducerStubOption {
	return func(s *ProducerStub) { s.group.BenchMode() }
}

// ProducerStubWithLatency returns the option that makes every method sleep
// d before it returns. Without ProducerStubWithClock, d passes in real
// time.
func ProducerStubWithLatency(d time.Duration) ProducerStubOption {
	return func(s *ProducerStub) { s.group.SetLatency(d) }
}

// ProducerStubFailing returns the option that arms err on each method that
// can fail. Each call of such a method then fails with err.
func ProducerStubFailing(err error) ProducerStubOption {
	return func(s *ProducerStub) {
		s.OnOpen.Faults(err, 1)
	}
}

// WithProducerOpen returns the option that sets the function of Open.
func WithProducerOpen(fn func(context.Context) (Cursor, error)) ProducerStubOption {
	return func(s *ProducerStub) { s.OnOpen.Func(fn) }
}

// WithProducerResume returns the option that sets the function of Resume.
func WithProducerResume(fn func(string) Cursor) ProducerStubOption {
	return func(s *ProducerStub) { s.OnResume.Func(fn) }
}

// WithProducerWatch returns the option that sets the function of Watch.
func WithProducerWatch(fn func() Subscription) ProducerStubOption {
	return func(s *ProducerStub) { s.OnWatch.Func(fn) }
}

// WithProducerNothing returns the option that sets the function of
// Nothing.
func WithProducerNothing(fn func() Empty) ProducerStubOption {
	return func(s *ProducerStub) { s.OnNothing.Func(fn) }
}

// WithProducerElsewhere returns the option that sets the function of
// Elsewhere.
func WithProducerElsewhere(fn func() other.Thing) ProducerStubOption {
	return func(s *ProducerStub) { s.OnElsewhere.Func(fn) }
}

// ProducerStub is a recording test double of Producer. Each field
// On<Method> configures its method. An unconfigured method records the
// call and returns the zero values of its results.
type ProducerStub struct {
	OnOpen      *ProducerOpenStub
	OnResume    *ProducerResumeStub
	OnWatch     *ProducerWatchStub
	OnNothing   *ProducerNothingStub
	OnElsewhere *ProducerElsewhereStub

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// The compiler checks that ProducerStub implements Producer.
var _ Producer = (*ProducerStub)(nil)

// NewProducerStub returns a double of Producer that reports its failures
// to tb and applies opts. The double checks the expected counts of calls
// at the end of the test. With a nil tb, the double reports nothing.
func NewProducerStub(tb assert.TB, opts ...ProducerStubOption) *ProducerStub {
	s := &ProducerStub{
		OnOpen: &ProducerOpenStub{
			MethodStub: stub.NewMethodStub[ProducerOpenCall](tb, "Producer.Open"),
		},
		OnResume: &ProducerResumeStub{
			MethodStub: stub.NewMethodStub[ProducerResumeCall](tb, "Producer.Resume"),
		},
		OnWatch: &ProducerWatchStub{
			MethodStub: stub.NewMethodStub[ProducerWatchCall](tb, "Producer.Watch"),
		},
		OnNothing: &ProducerNothingStub{
			MethodStub: stub.NewMethodStub[ProducerNothingCall](tb, "Producer.Nothing"),
		},
		OnElsewhere: &ProducerElsewhereStub{
			MethodStub: stub.NewMethodStub[ProducerElsewhereCall](tb, "Producer.Elsewhere"),
		},
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

// NewProducerStubDelayed returns a double whose methods sleep d on clk
// before they return. It applies opts after the clock and the latency.
func NewProducerStubDelayed(
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...ProducerStubOption,
) *ProducerStub {
	return NewProducerStub(tb, append([]ProducerStubOption{
		ProducerStubWithClock(clk),
		ProducerStubWithLatency(d),
	}, opts...)...)
}

// NewProducerStubFailing returns a double on which each call of a method
// that can fail fails with err. It applies opts after the failure.
func NewProducerStubFailing(tb assert.TB, err error, opts ...ProducerStubOption) *ProducerStub {
	return NewProducerStub(tb, append([]ProducerStubOption{ProducerStubFailing(err)}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *ProducerStub) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *ProducerOpenStub) invoke(ctx context.Context) func() ProducerOpenReturn {
	if s.fn == nil {
		return nil
	}
	return func() ProducerOpenReturn {
		r0, r1 := s.fn(ctx)
		return ProducerOpenReturn{Result: r0, Err: r1}
	}
}

// Open records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *ProducerStub) Open(ctx context.Context) (Cursor, error) {
	call := ProducerOpenCall{Ctx: ctx}
	r := stub.Answer(s.OnOpen.MethodStub, call, stub.Arms[ProducerOpenCall, ProducerOpenReturn]{
		Invoke:   s.OnOpen.invoke(ctx),
		Fallback: s.OnOpen.fallback,
		Fault:    func(err error) ProducerOpenReturn { return ProducerOpenReturn{Err: err} },
		Stamp: func(c ProducerOpenCall, r ProducerOpenReturn) ProducerOpenCall {
			c.Result = r.Result
			c.Err = r.Err
			return c
		},
	})
	return r.Result, r.Err
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *ProducerResumeStub) invoke(token string) func() ProducerResumeReturn {
	if s.fn == nil {
		return nil
	}
	return func() ProducerResumeReturn {
		r0 := s.fn(token)
		return ProducerResumeReturn{Result: r0}
	}
}

// Resume records the call. An armed fault comes first, then the function
// of Func, then the results of Returns, and then the zero values.
func (s *ProducerStub) Resume(token string) Cursor {
	call := ProducerResumeCall{Token: token}
	r := stub.Answer(
		s.OnResume.MethodStub,
		call,
		stub.Arms[ProducerResumeCall, ProducerResumeReturn]{
			Invoke:   s.OnResume.invoke(token),
			Fallback: s.OnResume.fallback,
			Stamp: func(c ProducerResumeCall, r ProducerResumeReturn) ProducerResumeCall {
				c.Result = r.Result
				return c
			},
		},
	)
	return r.Result
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *ProducerWatchStub) invoke() func() ProducerWatchReturn {
	if s.fn == nil {
		return nil
	}
	return func() ProducerWatchReturn {
		r0 := s.fn()
		return ProducerWatchReturn{Result: r0}
	}
}

// Watch records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *ProducerStub) Watch() Subscription {
	call := ProducerWatchCall{}
	r := stub.Answer(s.OnWatch.MethodStub, call, stub.Arms[ProducerWatchCall, ProducerWatchReturn]{
		Invoke:   s.OnWatch.invoke(),
		Fallback: s.OnWatch.fallback,
		Stamp: func(c ProducerWatchCall, r ProducerWatchReturn) ProducerWatchCall {
			c.Result = r.Result
			return c
		},
	})
	return r.Result
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *ProducerNothingStub) invoke() func() ProducerNothingReturn {
	if s.fn == nil {
		return nil
	}
	return func() ProducerNothingReturn {
		r0 := s.fn()
		return ProducerNothingReturn{Result: r0}
	}
}

// Nothing records the call. An armed fault comes first, then the function
// of Func, then the results of Returns, and then the zero values.
func (s *ProducerStub) Nothing() Empty {
	call := ProducerNothingCall{}
	r := stub.Answer(
		s.OnNothing.MethodStub,
		call,
		stub.Arms[ProducerNothingCall, ProducerNothingReturn]{
			Invoke:   s.OnNothing.invoke(),
			Fallback: s.OnNothing.fallback,
			Stamp: func(c ProducerNothingCall, r ProducerNothingReturn) ProducerNothingCall {
				c.Result = r.Result
				return c
			},
		},
	)
	return r.Result
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *ProducerElsewhereStub) invoke() func() ProducerElsewhereReturn {
	if s.fn == nil {
		return nil
	}
	return func() ProducerElsewhereReturn {
		r0 := s.fn()
		return ProducerElsewhereReturn{Result: r0}
	}
}

// Elsewhere records the call. An armed fault comes first, then the
// function of Func, then the results of Returns, and then the zero values.
func (s *ProducerStub) Elsewhere() other.Thing {
	call := ProducerElsewhereCall{}
	r := stub.Answer(
		s.OnElsewhere.MethodStub,
		call,
		stub.Arms[ProducerElsewhereCall, ProducerElsewhereReturn]{
			Invoke:   s.OnElsewhere.invoke(),
			Fallback: s.OnElsewhere.fallback,
			Stamp: func(c ProducerElsewhereCall, r ProducerElsewhereReturn) ProducerElsewhereCall {
				c.Result = r.Result
				return c
			},
		},
	)
	return r.Result
}

// Cursor is a handle that a method of Producer returns.

// CursorNextCall records one call of Cursor.Next.
type CursorNextCall struct {
	Ctx     context.Context
	Result0 string
	Result1 bool
	Err     error
}

// CursorCloseCall records one call of Cursor.Close.
type CursorCloseCall struct {
	Ctx context.Context
	Err error
}

// CursorNextStub configures Next and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type CursorNextStub struct {
	*stub.MethodStub[CursorNextCall]

	fn       func(context.Context) (string, bool, error)
	fallback *CursorNextReturn
}

// CursorNextReturn contains the fixed results of Next.
type CursorNextReturn struct {
	Result0 string
	Result1 bool
	Err     error
}

// Returns sets the results of every call of Next. A function and an
// injected fault take precedence over them.
func (s *CursorNextStub) Returns(result0 string, result1 bool, err error) *CursorNextStub {
	s.fallback = &CursorNextReturn{Result0: result0, Result1: result1, Err: err}
	return s
}

// Func sets the function that computes the results of Next from its
// arguments. An injected fault takes precedence over it.
func (s *CursorNextStub) Func(fn func(context.Context) (string, bool, error)) *CursorNextStub {
	s.fn = fn
	return s
}

// CursorCloseStub configures Close and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type CursorCloseStub struct {
	*stub.MethodStub[CursorCloseCall]

	fn       func(context.Context) error
	fallback *CursorCloseReturn
}

// CursorCloseReturn contains the fixed results of Close.
type CursorCloseReturn struct {
	Err error
}

// Returns sets the results of every call of Close. A function and an
// injected fault take precedence over them.
func (s *CursorCloseStub) Returns(err error) *CursorCloseStub {
	s.fallback = &CursorCloseReturn{Err: err}
	return s
}

// Func sets the function that computes the results of Close from its
// arguments. An injected fault takes precedence over it.
func (s *CursorCloseStub) Func(fn func(context.Context) error) *CursorCloseStub {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// CursorStubOption configures a [CursorStub] at its construction.
type CursorStubOption func(*CursorStub)

// CursorStubStrict returns the option that makes an unconfigured method
// fail the test instead of returning its zero value.
func CursorStubStrict() CursorStubOption {
	return func(s *CursorStub) { s.strict = true }
}

// CursorStubDelegateTo returns the option that passes each call to impl.
// The double still records the call.
func CursorStubDelegateTo(impl Cursor) CursorStubOption {
	return func(s *CursorStub) {
		s.OnNext.Func(impl.Next)
		s.OnClose.Func(impl.Close)
	}
}

// CursorStubWithClock returns the option that sets the clock of the double
// to clk. The double takes the latency, the windows of the faults and the
// instants of the calls from clk. A latency on a controlled clock costs no
// real time.
func CursorStubWithClock(clk assert.Clock) CursorStubOption {
	return func(s *CursorStub) { s.group.SetClock(clk) }
}

// CursorStubWithRandSource returns the option that sets the random source
// of the probabilistic faults to src. A seeded source repeats the faults
// of a run.
func CursorStubWithRandSource(src rand.Source) CursorStubOption {
	return func(s *CursorStub) { s.group.SetRandSource(src) }
}

// CursorStubBenchMode returns the option that turns the recording of calls
// off. The methods still return their results.
func CursorStubBenchMode() CursorStubOption {
	return func(s *CursorStub) { s.group.BenchMode() }
}

// CursorStubWithLatency returns the option that makes every method sleep d
// before it returns. Without CursorStubWithClock, d passes in real time.
func CursorStubWithLatency(d time.Duration) CursorStubOption {
	return func(s *CursorStub) { s.group.SetLatency(d) }
}

// CursorStubFailing returns the option that arms err on each method that
// can fail. Each call of such a method then fails with err.
func CursorStubFailing(err error) CursorStubOption {
	return func(s *CursorStub) {
		s.OnNext.Faults(err, 1)
		s.OnClose.Faults(err, 1)
	}
}

// WithCursorNext returns the option that sets the function of Next.
func WithCursorNext(fn func(context.Context) (string, bool, error)) CursorStubOption {
	return func(s *CursorStub) { s.OnNext.Func(fn) }
}

// WithCursorClose returns the option that sets the function of Close.
func WithCursorClose(fn func(context.Context) error) CursorStubOption {
	return func(s *CursorStub) { s.OnClose.Func(fn) }
}

// CursorStub is a recording test double of Cursor. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type CursorStub struct {
	OnNext  *CursorNextStub
	OnClose *CursorCloseStub

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// The compiler checks that CursorStub implements Cursor.
var _ Cursor = (*CursorStub)(nil)

// NewCursorStub returns a double of Cursor that reports its failures to tb
// and applies opts. The double checks the expected counts of calls at the
// end of the test. With a nil tb, the double reports nothing.
func NewCursorStub(tb assert.TB, opts ...CursorStubOption) *CursorStub {
	s := &CursorStub{
		OnNext: &CursorNextStub{MethodStub: stub.NewMethodStub[CursorNextCall](tb, "Cursor.Next")},
		OnClose: &CursorCloseStub{
			MethodStub: stub.NewMethodStub[CursorCloseCall](tb, "Cursor.Close"),
		},
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

// NewCursorStubDelayed returns a double whose methods sleep d on clk
// before they return. It applies opts after the clock and the latency.
func NewCursorStubDelayed(
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...CursorStubOption,
) *CursorStub {
	return NewCursorStub(tb, append([]CursorStubOption{
		CursorStubWithClock(clk),
		CursorStubWithLatency(d),
	}, opts...)...)
}

// NewCursorStubFailing returns a double on which each call of a method
// that can fail fails with err. It applies opts after the failure.
func NewCursorStubFailing(tb assert.TB, err error, opts ...CursorStubOption) *CursorStub {
	return NewCursorStub(tb, append([]CursorStubOption{CursorStubFailing(err)}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *CursorStub) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *CursorNextStub) invoke(ctx context.Context) func() CursorNextReturn {
	if s.fn == nil {
		return nil
	}
	return func() CursorNextReturn {
		r0, r1, r2 := s.fn(ctx)
		return CursorNextReturn{Result0: r0, Result1: r1, Err: r2}
	}
}

// Next records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *CursorStub) Next(ctx context.Context) (string, bool, error) {
	call := CursorNextCall{Ctx: ctx}
	r := stub.Answer(s.OnNext.MethodStub, call, stub.Arms[CursorNextCall, CursorNextReturn]{
		Invoke:   s.OnNext.invoke(ctx),
		Fallback: s.OnNext.fallback,
		Fault:    func(err error) CursorNextReturn { return CursorNextReturn{Err: err} },
		Stamp: func(c CursorNextCall, r CursorNextReturn) CursorNextCall {
			c.Result0 = r.Result0
			c.Result1 = r.Result1
			c.Err = r.Err
			return c
		},
	})
	return r.Result0, r.Result1, r.Err
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *CursorCloseStub) invoke(ctx context.Context) func() CursorCloseReturn {
	if s.fn == nil {
		return nil
	}
	return func() CursorCloseReturn {
		r0 := s.fn(ctx)
		return CursorCloseReturn{Err: r0}
	}
}

// Close records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *CursorStub) Close(ctx context.Context) error {
	call := CursorCloseCall{Ctx: ctx}
	r := stub.Answer(s.OnClose.MethodStub, call, stub.Arms[CursorCloseCall, CursorCloseReturn]{
		Invoke:   s.OnClose.invoke(ctx),
		Fallback: s.OnClose.fallback,
		Fault:    func(err error) CursorCloseReturn { return CursorCloseReturn{Err: err} },
		Stamp: func(c CursorCloseCall, r CursorCloseReturn) CursorCloseCall {
			c.Err = r.Err
			return c
		},
	})
	return r.Err
}

// ReaderReadCall records one call of Reader.Read.
type ReaderReadCall struct {
	Result Cursor
}

// ReaderReadStub configures Read and records its calls. Read returns no
// error, so a fault armed through it fails the test.
type ReaderReadStub struct {
	*stub.MethodStub[ReaderReadCall]

	fn       func() Cursor
	fallback *ReaderReadReturn
}

// ReaderReadReturn contains the fixed results of Read.
type ReaderReadReturn struct {
	Result Cursor
}

// Returns sets the results of every call of Read. A function takes
// precedence over them.
func (s *ReaderReadStub) Returns(result Cursor) *ReaderReadStub {
	s.fallback = &ReaderReadReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Read from its
// arguments.
func (s *ReaderReadStub) Func(fn func() Cursor) *ReaderReadStub {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// ReaderStubOption configures a [ReaderStub] at its construction.
type ReaderStubOption func(*ReaderStub)

// ReaderStubStrict returns the option that makes an unconfigured method
// fail the test instead of returning its zero value.
func ReaderStubStrict() ReaderStubOption {
	return func(s *ReaderStub) { s.strict = true }
}

// ReaderStubDelegateTo returns the option that passes each call to impl.
// The double still records the call.
func ReaderStubDelegateTo(impl Reader) ReaderStubOption {
	return func(s *ReaderStub) {
		s.OnRead.Func(impl.Read)
	}
}

// ReaderStubWithClock returns the option that sets the clock of the double
// to clk. The double takes the latency, the windows of the faults and the
// instants of the calls from clk. A latency on a controlled clock costs no
// real time.
func ReaderStubWithClock(clk assert.Clock) ReaderStubOption {
	return func(s *ReaderStub) { s.group.SetClock(clk) }
}

// ReaderStubWithRandSource returns the option that sets the random source
// of the probabilistic faults to src. A seeded source repeats the faults
// of a run.
func ReaderStubWithRandSource(src rand.Source) ReaderStubOption {
	return func(s *ReaderStub) { s.group.SetRandSource(src) }
}

// ReaderStubBenchMode returns the option that turns the recording of calls
// off. The methods still return their results.
func ReaderStubBenchMode() ReaderStubOption {
	return func(s *ReaderStub) { s.group.BenchMode() }
}

// ReaderStubWithLatency returns the option that makes every method sleep d
// before it returns. Without ReaderStubWithClock, d passes in real time.
func ReaderStubWithLatency(d time.Duration) ReaderStubOption {
	return func(s *ReaderStub) { s.group.SetLatency(d) }
}

// WithReaderRead returns the option that sets the function of Read.
func WithReaderRead(fn func() Cursor) ReaderStubOption {
	return func(s *ReaderStub) { s.OnRead.Func(fn) }
}

// ReaderStub is a recording test double of Reader. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type ReaderStub struct {
	OnRead *ReaderReadStub

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// The compiler checks that ReaderStub implements Reader.
var _ Reader = (*ReaderStub)(nil)

// NewReaderStub returns a double of Reader that reports its failures to tb
// and applies opts. The double checks the expected counts of calls at the
// end of the test. With a nil tb, the double reports nothing.
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

// NewReaderStubDelayed returns a double whose methods sleep d on clk
// before they return. It applies opts after the clock and the latency.
func NewReaderStubDelayed(
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...ReaderStubOption,
) *ReaderStub {
	return NewReaderStub(tb, append([]ReaderStubOption{
		ReaderStubWithClock(clk),
		ReaderStubWithLatency(d),
	}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *ReaderStub) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *ReaderReadStub) invoke() func() ReaderReadReturn {
	if s.fn == nil {
		return nil
	}
	return func() ReaderReadReturn {
		r0 := s.fn()
		return ReaderReadReturn{Result: r0}
	}
}

// Read records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *ReaderStub) Read() Cursor {
	call := ReaderReadCall{}
	r := stub.Answer(s.OnRead.MethodStub, call, stub.Arms[ReaderReadCall, ReaderReadReturn]{
		Invoke:   s.OnRead.invoke(),
		Fallback: s.OnRead.fallback,
		Stamp: func(c ReaderReadCall, r ReaderReadReturn) ReaderReadCall {
			c.Result = r.Result
			return c
		},
	})
	return r.Result
}

// SubscriptionCancelCall records one call of Subscription.Cancel.
type SubscriptionCancelCall struct{}

// SubscriptionCancelStub configures Cancel and records its calls. Cancel
// returns no error, so a fault armed through it fails the test.
type SubscriptionCancelStub struct {
	*stub.MethodStub[SubscriptionCancelCall]

	fn func()
}

// Func sets the function that computes the results of Cancel from its
// arguments.
func (s *SubscriptionCancelStub) Func(fn func()) *SubscriptionCancelStub {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// SubscriptionStubOption configures a [SubscriptionStub] at its
// construction.
type SubscriptionStubOption func(*SubscriptionStub)

// SubscriptionStubStrict returns the option that makes an unconfigured
// method fail the test instead of returning its zero value.
func SubscriptionStubStrict() SubscriptionStubOption {
	return func(s *SubscriptionStub) { s.strict = true }
}

// SubscriptionStubDelegateTo returns the option that passes each call to
// impl. The double still records the call.
func SubscriptionStubDelegateTo(impl Subscription) SubscriptionStubOption {
	return func(s *SubscriptionStub) {
		s.OnCancel.Func(impl.Cancel)
	}
}

// SubscriptionStubWithClock returns the option that sets the clock of the
// double to clk. The double takes the latency, the windows of the faults
// and the instants of the calls from clk. A latency on a controlled clock
// costs no real time.
func SubscriptionStubWithClock(clk assert.Clock) SubscriptionStubOption {
	return func(s *SubscriptionStub) { s.group.SetClock(clk) }
}

// SubscriptionStubWithRandSource returns the option that sets the random
// source of the probabilistic faults to src. A seeded source repeats the
// faults of a run.
func SubscriptionStubWithRandSource(src rand.Source) SubscriptionStubOption {
	return func(s *SubscriptionStub) { s.group.SetRandSource(src) }
}

// SubscriptionStubBenchMode returns the option that turns the recording of
// calls off. The methods still return their results.
func SubscriptionStubBenchMode() SubscriptionStubOption {
	return func(s *SubscriptionStub) { s.group.BenchMode() }
}

// SubscriptionStubWithLatency returns the option that makes every method
// sleep d before it returns. Without SubscriptionStubWithClock, d passes
// in real time.
func SubscriptionStubWithLatency(d time.Duration) SubscriptionStubOption {
	return func(s *SubscriptionStub) { s.group.SetLatency(d) }
}

// WithSubscriptionCancel returns the option that sets the function of
// Cancel.
func WithSubscriptionCancel(fn func()) SubscriptionStubOption {
	return func(s *SubscriptionStub) { s.OnCancel.Func(fn) }
}

// SubscriptionStub is a recording test double of Subscription. Each field
// On<Method> configures its method. An unconfigured method records the
// call and returns the zero values of its results.
type SubscriptionStub struct {
	OnCancel *SubscriptionCancelStub

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// The compiler checks that SubscriptionStub implements Subscription.
var _ Subscription = (*SubscriptionStub)(nil)

// NewSubscriptionStub returns a double of Subscription that reports its
// failures to tb and applies opts. The double checks the expected counts
// of calls at the end of the test. With a nil tb, the double reports
// nothing.
func NewSubscriptionStub(tb assert.TB, opts ...SubscriptionStubOption) *SubscriptionStub {
	s := &SubscriptionStub{
		OnCancel: &SubscriptionCancelStub{
			MethodStub: stub.NewMethodStub[SubscriptionCancelCall](tb, "Subscription.Cancel"),
		},
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

// NewSubscriptionStubDelayed returns a double whose methods sleep d on clk
// before they return. It applies opts after the clock and the latency.
func NewSubscriptionStubDelayed(
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...SubscriptionStubOption,
) *SubscriptionStub {
	return NewSubscriptionStub(tb, append([]SubscriptionStubOption{
		SubscriptionStubWithClock(clk),
		SubscriptionStubWithLatency(d),
	}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *SubscriptionStub) ResetCalls() {
	s.group.Reset()
}

// invokeVoid returns a call of the function of Func with the arguments, and
// nil where no function is set. The empty struct is the result of a method
// without results.
func (s *SubscriptionCancelStub) invokeVoid() func() struct{} {
	if s.fn == nil {
		return nil
	}
	return func() struct{} {
		s.fn()
		return struct{}{}
	}
}

// Cancel records the call. An armed fault comes first, then the function
// of Func, then the results of Returns, and then the zero values.
func (s *SubscriptionStub) Cancel() {
	call := SubscriptionCancelCall{}
	stub.Answer(s.OnCancel.MethodStub, call, stub.Arms[SubscriptionCancelCall, struct{}]{
		Invoke: s.OnCancel.invokeVoid(),
	})
}
