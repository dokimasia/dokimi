package catalog

import (
	"context"
	rand "math/rand/v2"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// StoreGetCall records one invocation of Store.Get.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StoreGetCall struct {
	Ctx    context.Context
	ID     string
	Result Item
	Err    error
}

// StorePutCall records one invocation of Store.Put.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StorePutCall struct {
	A0  context.Context
	A1  Item
	Err error
}

// StoreLenCall records one invocation of Store.Len.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StoreLenCall struct {
	Result int
}

// StoreLookupCall records one invocation of Store.Lookup.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StoreLookupCall struct {
	Key  string
	Item Item
	Ok   bool
}

// StoreLogCall records one invocation of Store.Log.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StoreLogCall struct {
	Format string
	Args   []string
}

// StoreSplitCall records one invocation of Store.Split.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StoreSplitCall struct {
	S       string
	Result0 string
	Result1 string
	Err     error
}

// StoreCloseCall records one invocation of Store.Close.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type StoreCloseCall struct {
	Err error
}

// --- Per-method configuration ---

// StoreGetStub controls how the double answers Get and records what it was
// asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type StoreGetStub struct {
	*stub.MethodStub[StoreGetCall]

	fn       func(context.Context, string) (Item, error)
	fallback *StoreGetReturn
}

// StoreGetReturn holds the fixed answer configured through Returns.
type StoreGetReturn struct {
	Result Item
	Err    error
}

// Returns pins a fixed result for every call to Get. A Func override takes
// precedence over it, and so does an injected fault.
func (s *StoreGetStub) Returns(result Item, err error) *StoreGetStub {
	s.fallback = &StoreGetReturn{Result: result, Err: err}
	return s
}

// Func supplies a body for Get, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *StoreGetStub) Func(fn func(context.Context, string) (Item, error)) *StoreGetStub {
	s.fn = fn
	return s
}

// StorePutStub controls how the double answers Put and records what it was
// asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type StorePutStub struct {
	*stub.MethodStub[StorePutCall]

	fn       func(context.Context, Item) error
	fallback *StorePutReturn
}

// StorePutReturn holds the fixed answer configured through Returns.
type StorePutReturn struct {
	Err error
}

// Returns pins a fixed result for every call to Put. A Func override takes
// precedence over it, and so does an injected fault.
func (s *StorePutStub) Returns(err error) *StorePutStub {
	s.fallback = &StorePutReturn{Err: err}
	return s
}

// Func supplies a body for Put, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *StorePutStub) Func(fn func(context.Context, Item) error) *StorePutStub {
	s.fn = fn
	return s
}

// StoreLenStub controls how the double answers Len and records what it was
// asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type StoreLenStub struct {
	*stub.MethodStub[StoreLenCall]

	fn       func() int
	fallback *StoreLenReturn
}

// StoreLenReturn holds the fixed answer configured through Returns.
type StoreLenReturn struct {
	Result int
}

// Returns pins a fixed result for every call to Len. A Func override takes
// precedence over it.
//
// No fault can take precedence here: Len answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *StoreLenStub) Returns(result int) *StoreLenStub {
	s.fallback = &StoreLenReturn{Result: result}
	return s
}

// Func supplies a body for Len, for when the answer depends on the
// arguments. Nothing takes precedence over it: Len answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *StoreLenStub) Func(fn func() int) *StoreLenStub {
	s.fn = fn
	return s
}

// StoreLookupStub controls how the double answers Lookup and records what
// it was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type StoreLookupStub struct {
	*stub.MethodStub[StoreLookupCall]

	fn       func(string) (Item, bool)
	fallback *StoreLookupReturn
}

// StoreLookupReturn holds the fixed answer configured through Returns.
type StoreLookupReturn struct {
	Item Item
	Ok   bool
}

// Returns pins a fixed result for every call to Lookup. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Lookup answers no error, so the
// double is built with no fault arm and a fault armed through Faults fails
// the test.
func (s *StoreLookupStub) Returns(item Item, ok bool) *StoreLookupStub {
	s.fallback = &StoreLookupReturn{Item: item, Ok: ok}
	return s
}

// Func supplies a body for Lookup, for when the answer depends on the
// arguments. Nothing takes precedence over it: Lookup answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *StoreLookupStub) Func(fn func(string) (Item, bool)) *StoreLookupStub {
	s.fn = fn
	return s
}

// StoreLogStub controls how the double answers Log and records what it was
// asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type StoreLogStub struct {
	*stub.MethodStub[StoreLogCall]

	fn func(string, ...string)
}

// Func supplies a body for Log, for when the answer depends on the
// arguments. Nothing takes precedence over it: Log answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *StoreLogStub) Func(fn func(string, ...string)) *StoreLogStub {
	s.fn = fn
	return s
}

// StoreSplitStub controls how the double answers Split and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type StoreSplitStub struct {
	*stub.MethodStub[StoreSplitCall]

	fn       func(string) (string, string, error)
	fallback *StoreSplitReturn
}

// StoreSplitReturn holds the fixed answer configured through Returns.
type StoreSplitReturn struct {
	Result0 string
	Result1 string
	Err     error
}

// Returns pins a fixed result for every call to Split. A Func override
// takes precedence over it, and so does an injected fault.
func (s *StoreSplitStub) Returns(result0 string, result1 string, err error) *StoreSplitStub {
	s.fallback = &StoreSplitReturn{Result0: result0, Result1: result1, Err: err}
	return s
}

// Func supplies a body for Split, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *StoreSplitStub) Func(fn func(string) (string, string, error)) *StoreSplitStub {
	s.fn = fn
	return s
}

// StoreCloseStub controls how the double answers Close and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type StoreCloseStub struct {
	*stub.MethodStub[StoreCloseCall]

	fn       func() error
	fallback *StoreCloseReturn
}

// StoreCloseReturn holds the fixed answer configured through Returns.
type StoreCloseReturn struct {
	Err error
}

// Returns pins a fixed result for every call to Close. A Func override
// takes precedence over it, and so does an injected fault.
func (s *StoreCloseStub) Returns(err error) *StoreCloseStub {
	s.fallback = &StoreCloseReturn{Err: err}
	return s
}

// Func supplies a body for Close, for when the answer depends on the
// arguments. An injected fault still takes precedence.
func (s *StoreCloseStub) Func(fn func() error) *StoreCloseStub {
	s.fn = fn
	return s
}

// --- StoreStub ---

// The witness of the major version of the runtime of doubles.
var _ = stub.CompatV1

// StoreStubOption configures a [StoreStub] at construction time.
type StoreStubOption func(*StoreStub)

// StoreStubStrict makes an unconfigured method fail the test rather than
// return its zero value.
func StoreStubStrict() StoreStubOption {
	return func(s *StoreStub) { s.strict = true }
}

// StoreStubDelegateTo wraps a real implementation: every method forwards
// to impl and is recorded on the double on the way through.
func StoreStubDelegateTo(impl Store) StoreStubOption {
	return func(s *StoreStub) {
		s.OnGet.Func(impl.Get)
		s.OnPut.Func(impl.Put)
		s.OnLen.Func(impl.Len)
		s.OnLookup.Func(impl.Lookup)
		s.OnLog.Func(impl.Log)
		s.OnSplit.Func(impl.Split)
		s.OnClose.Func(impl.Close)
	}
}

// StoreStubWithClock drives latency and time-windowed faults from clk
// rather than from wall time, so a test asserting on a five-second timeout
// does not take five seconds.
func StoreStubWithClock(clk assert.Clock) StoreStubOption {
	return func(s *StoreStub) { s.group.SetClock(clk) }
}

// StoreStubWithRandSource makes probabilistic fault injection
// reproducible. A failure nobody can replay is not a finding.
func StoreStubWithRandSource(src rand.Source) StoreStubOption {
	return func(s *StoreStub) { s.group.SetRandSource(src) }
}

// StoreStubBenchMode disables call recording. Dispatch still works; only
// the accumulating call log is dropped.
func StoreStubBenchMode() StoreStubOption {
	return func(s *StoreStub) { s.group.BenchMode() }
}

// StoreStubWithLatency makes every method sleep d before answering, for a
// test standing in for a slow dependency rather than one slow method. Pair
// it with StoreStubWithClock, or d elapses in wall time.
func StoreStubWithLatency(d time.Duration) StoreStubOption {
	return func(s *StoreStub) { s.group.SetLatency(d) }
}

// StoreStubFailing arms every method that can fail with err, so a caller's
// error path is one option rather than one Faults call per method. A
// method answering no error is left alone.
func StoreStubFailing(err error) StoreStubOption {
	return func(s *StoreStub) {
		s.OnGet.Faults(err, 1)
		s.OnPut.Faults(err, 1)
		s.OnSplit.Faults(err, 1)
		s.OnClose.Faults(err, 1)
	}
}

// WithStoreGet sets Get's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithStoreGet(fn func(context.Context, string) (Item, error)) StoreStubOption {
	return func(s *StoreStub) { s.OnGet.Func(fn) }
}

// WithStorePut sets Put's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithStorePut(fn func(context.Context, Item) error) StoreStubOption {
	return func(s *StoreStub) { s.OnPut.Func(fn) }
}

// WithStoreLen sets Len's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithStoreLen(fn func() int) StoreStubOption {
	return func(s *StoreStub) { s.OnLen.Func(fn) }
}

// WithStoreLookup sets Lookup's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithStoreLookup(fn func(string) (Item, bool)) StoreStubOption {
	return func(s *StoreStub) { s.OnLookup.Func(fn) }
}

// WithStoreLog sets Log's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithStoreLog(fn func(string, ...string)) StoreStubOption {
	return func(s *StoreStub) { s.OnLog.Func(fn) }
}

// WithStoreSplit sets Split's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithStoreSplit(fn func(string) (string, string, error)) StoreStubOption {
	return func(s *StoreStub) { s.OnSplit.Func(fn) }
}

// WithStoreClose sets Close's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithStoreClose(fn func() error) StoreStubOption {
	return func(s *StoreStub) { s.OnClose.Func(fn) }
}

// StoreStub is a recording test double for Store.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method returns its zero value and records the call.
type StoreStub struct {
	OnGet    *StoreGetStub
	OnPut    *StorePutStub
	OnLen    *StoreLenStub
	OnLookup *StoreLookupStub
	OnLog    *StoreLogStub
	OnSplit  *StoreSplitStub
	// OnClose configures Close, contributed by Closer.
	OnClose *StoreCloseStub

	// group is every method stub above, viewed through the surface that
	// does not depend on a signature.
	group  *stub.Group
	strict bool
}

// Compile-time proof that StoreStub satisfies Store.
var _ Store = (*StoreStub)(nil)

// NewStoreStub returns a double bound to tb.
//
// Passing tb registers a cleanup that verifies every method's call-count
// expectations when the test ends, so an unmet Times is reported without
// the caller remembering to check. A nil tb skips that, which is what
// benchmarks and non-test callers want.
func NewStoreStub(tb assert.TB, opts ...StoreStubOption) *StoreStub {
	s := &StoreStub{
		OnGet:    &StoreGetStub{MethodStub: stub.NewMethodStub[StoreGetCall](tb, "Store.Get")},
		OnPut:    &StorePutStub{MethodStub: stub.NewMethodStub[StorePutCall](tb, "Store.Put")},
		OnLen:    &StoreLenStub{MethodStub: stub.NewMethodStub[StoreLenCall](tb, "Store.Len")},
		OnLookup: &StoreLookupStub{MethodStub: stub.NewMethodStub[StoreLookupCall](tb, "Store.Lookup")},
		OnLog:    &StoreLogStub{MethodStub: stub.NewMethodStub[StoreLogCall](tb, "Store.Log")},
		OnSplit:  &StoreSplitStub{MethodStub: stub.NewMethodStub[StoreSplitCall](tb, "Store.Split")},
		OnClose:  &StoreCloseStub{MethodStub: stub.NewMethodStub[StoreCloseCall](tb, "Store.Close")},
	}
	s.group = stub.NewGroup(
		s.OnGet.MethodStub,
		s.OnPut.MethodStub,
		s.OnLen.MethodStub,
		s.OnLookup.MethodStub,
		s.OnLog.MethodStub,
		s.OnSplit.MethodStub,
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

// NewStoreStubDelayed returns a double whose every method sleeps d against
// clk before answering — the clock and the latency a slow-dependency test
// wants, in one call rather than two options. Further options apply after
// the pair, so a caller can still add Strict or a method override.
func NewStoreStubDelayed(tb assert.TB, d time.Duration, clk assert.Clock, opts ...StoreStubOption) *StoreStub {
	return NewStoreStub(tb, append([]StoreStubOption{
		StoreStubWithClock(clk),
		StoreStubWithLatency(d),
	}, opts...)...)
}

// NewStoreStubFailing returns a double whose every fallible method answers
// err — [StoreStubFailing] applied at construction, for the test whose
// whole subject is the caller's error path. Further options apply after
// it.
func NewStoreStubFailing(tb assert.TB, err error, opts ...StoreStubOption) *StoreStub {
	return NewStoreStub(tb, append([]StoreStubOption{StoreStubFailing(err)}, opts...)...)
}

// ResetCalls clears recorded calls, fault counters, and call-count
// expectations. Func and Returns configuration is deliberately preserved, so
// a test can reuse one double across phases without rebuilding it.
func (s *StoreStub) ResetCalls() {
	s.group.Reset()
}

// --- Interface methods ---

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *StoreGetStub) invoke(ctx context.Context, id string) func() StoreGetReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreGetReturn {
		r0, r1 := s.fn(ctx, id)
		return StoreGetReturn{Result: r0, Err: r1}
	}
}

// Get records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *StoreStub) Get(ctx context.Context, id string) (Item, error) {
	call := StoreGetCall{Ctx: ctx, ID: id}
	r := stub.Answer(s.OnGet.MethodStub, &call, stub.Arms[StoreGetCall, StoreGetReturn]{
		Invoke:   s.OnGet.invoke(ctx, id),
		Fallback: s.OnGet.fallback,
		Fault:    func(err error) StoreGetReturn { return StoreGetReturn{Err: err} },
		Stamp: func(c *StoreGetCall, r StoreGetReturn) {
			c.Result = r.Result
			c.Err = r.Err
		},
	})
	return r.Result, r.Err
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *StorePutStub) invoke(a0 context.Context, a1 Item) func() StorePutReturn {
	if s.fn == nil {
		return nil
	}
	return func() StorePutReturn {
		r0 := s.fn(a0, a1)
		return StorePutReturn{Err: r0}
	}
}

// Put records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *StoreStub) Put(a0 context.Context, a1 Item) error {
	call := StorePutCall{A0: a0, A1: a1}
	r := stub.Answer(s.OnPut.MethodStub, &call, stub.Arms[StorePutCall, StorePutReturn]{
		Invoke:   s.OnPut.invoke(a0, a1),
		Fallback: s.OnPut.fallback,
		Fault:    func(err error) StorePutReturn { return StorePutReturn{Err: err} },
		Stamp: func(c *StorePutCall, r StorePutReturn) {
			c.Err = r.Err
		},
	})
	return r.Err
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *StoreLenStub) invoke() func() StoreLenReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreLenReturn {
		r0 := s.fn()
		return StoreLenReturn{Result: r0}
	}
}

// Len records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *StoreStub) Len() int {
	call := StoreLenCall{}
	r := stub.Answer(s.OnLen.MethodStub, &call, stub.Arms[StoreLenCall, StoreLenReturn]{
		Invoke:   s.OnLen.invoke(),
		Fallback: s.OnLen.fallback,
		Stamp: func(c *StoreLenCall, r StoreLenReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *StoreLookupStub) invoke(key string) func() StoreLookupReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreLookupReturn {
		item, ok := s.fn(key)
		return StoreLookupReturn{Item: item, Ok: ok}
	}
}

// Lookup records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *StoreStub) Lookup(key string) (item Item, ok bool) {
	call := StoreLookupCall{Key: key}
	r := stub.Answer(s.OnLookup.MethodStub, &call, stub.Arms[StoreLookupCall, StoreLookupReturn]{
		Invoke:   s.OnLookup.invoke(key),
		Fallback: s.OnLookup.fallback,
		Stamp: func(c *StoreLookupCall, r StoreLookupReturn) {
			c.Item = r.Item
			c.Ok = r.Ok
		},
	})
	return r.Item, r.Ok
}

// invokeVoid adapts the Func override to the shape [stub.Answer] consumes.
// A method returning nothing still has an override worth dispatching to and
// a call worth recording; the empty struct is the return tuple of no returns.
func (s *StoreLogStub) invokeVoid(format string, args ...string) func() struct{} {
	if s.fn == nil {
		return nil
	}
	return func() struct{} {
		s.fn(format, args...)
		return struct{}{}
	}
}

// Log records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *StoreStub) Log(format string, args ...string) {
	call := StoreLogCall{Format: format, Args: args}
	stub.Answer(s.OnLog.MethodStub, &call, stub.Arms[StoreLogCall, struct{}]{
		Invoke: s.OnLog.invokeVoid(format, args...),
	})
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s0 *StoreSplitStub) invoke(s string) func() StoreSplitReturn {
	if s0.fn == nil {
		return nil
	}
	return func() StoreSplitReturn {
		r0, r1, r2 := s0.fn(s)
		return StoreSplitReturn{Result0: r0, Result1: r1, Err: r2}
	}
}

// Split records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s0 *StoreStub) Split(s string) (string, string, error) {
	call := StoreSplitCall{S: s}
	r := stub.Answer(s0.OnSplit.MethodStub, &call, stub.Arms[StoreSplitCall, StoreSplitReturn]{
		Invoke:   s0.OnSplit.invoke(s),
		Fallback: s0.OnSplit.fallback,
		Fault:    func(err error) StoreSplitReturn { return StoreSplitReturn{Err: err} },
		Stamp: func(c *StoreSplitCall, r StoreSplitReturn) {
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
func (s *StoreCloseStub) invoke() func() StoreCloseReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreCloseReturn {
		r0 := s.fn()
		return StoreCloseReturn{Err: r0}
	}
}

// Close records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *StoreStub) Close() error {
	call := StoreCloseCall{}
	r := stub.Answer(s.OnClose.MethodStub, &call, stub.Arms[StoreCloseCall, StoreCloseReturn]{
		Invoke:   s.OnClose.invoke(),
		Fallback: s.OnClose.fallback,
		Fault:    func(err error) StoreCloseReturn { return StoreCloseReturn{Err: err} },
		Stamp: func(c *StoreCloseCall, r StoreCloseReturn) {
			c.Err = r.Err
		},
	})
	return r.Err
}

// ClockNowCall records one invocation of Clock.Now.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type ClockNowCall struct {
	Result time.Time
}

// ClockZoneCall records one invocation of Clock.Zone.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type ClockZoneCall struct {
	Result string
}

// ClockEveryCall records one invocation of Clock.Every.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type ClockEveryCall struct {
	D    time.Duration
	Fn   func()
	Stop func()
}

// --- Per-method configuration ---

// ClockNowStub controls how the double answers Now and records what it was
// asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type ClockNowStub struct {
	*stub.MethodStub[ClockNowCall]

	fn       func() time.Time
	fallback *ClockNowReturn
}

// ClockNowReturn holds the fixed answer configured through Returns.
type ClockNowReturn struct {
	Result time.Time
}

// Returns pins a fixed result for every call to Now. A Func override takes
// precedence over it.
//
// No fault can take precedence here: Now answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *ClockNowStub) Returns(result time.Time) *ClockNowStub {
	s.fallback = &ClockNowReturn{Result: result}
	return s
}

// Func supplies a body for Now, for when the answer depends on the
// arguments. Nothing takes precedence over it: Now answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *ClockNowStub) Func(fn func() time.Time) *ClockNowStub {
	s.fn = fn
	return s
}

// ClockZoneStub controls how the double answers Zone and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type ClockZoneStub struct {
	*stub.MethodStub[ClockZoneCall]

	fn       func() string
	fallback *ClockZoneReturn
}

// ClockZoneReturn holds the fixed answer configured through Returns.
type ClockZoneReturn struct {
	Result string
}

// Returns pins a fixed result for every call to Zone. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Zone answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *ClockZoneStub) Returns(result string) *ClockZoneStub {
	s.fallback = &ClockZoneReturn{Result: result}
	return s
}

// Func supplies a body for Zone, for when the answer depends on the
// arguments. Nothing takes precedence over it: Zone answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *ClockZoneStub) Func(fn func() string) *ClockZoneStub {
	s.fn = fn
	return s
}

// ClockEveryStub controls how the double answers Every and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type ClockEveryStub struct {
	*stub.MethodStub[ClockEveryCall]

	fn       func(time.Duration, func()) func()
	fallback *ClockEveryReturn
}

// ClockEveryReturn holds the fixed answer configured through Returns.
type ClockEveryReturn struct {
	Stop func()
}

// Returns pins a fixed result for every call to Every. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Every answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *ClockEveryStub) Returns(stop func()) *ClockEveryStub {
	s.fallback = &ClockEveryReturn{Stop: stop}
	return s
}

// Func supplies a body for Every, for when the answer depends on the
// arguments. Nothing takes precedence over it: Every answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *ClockEveryStub) Func(fn func(time.Duration, func()) func()) *ClockEveryStub {
	s.fn = fn
	return s
}

// --- ClockStub ---

// The witness of the major version of the runtime of doubles.
var _ = stub.CompatV1

// ClockStubOption configures a [ClockStub] at construction time.
type ClockStubOption func(*ClockStub)

// ClockStubStrict makes an unconfigured method fail the test rather than
// return its zero value.
func ClockStubStrict() ClockStubOption {
	return func(s *ClockStub) { s.strict = true }
}

// ClockStubDelegateTo wraps a real implementation: every method forwards
// to impl and is recorded on the double on the way through.
func ClockStubDelegateTo(impl *Clock) ClockStubOption {
	return func(s *ClockStub) {
		s.OnNow.Func(impl.Now)
		s.OnZone.Func(impl.Zone)
		s.OnEvery.Func(impl.Every)
	}
}

// ClockStubWithClock drives latency and time-windowed faults from clk
// rather than from wall time, so a test asserting on a five-second timeout
// does not take five seconds.
func ClockStubWithClock(clk assert.Clock) ClockStubOption {
	return func(s *ClockStub) { s.group.SetClock(clk) }
}

// ClockStubWithRandSource makes probabilistic fault injection
// reproducible. A failure nobody can replay is not a finding.
func ClockStubWithRandSource(src rand.Source) ClockStubOption {
	return func(s *ClockStub) { s.group.SetRandSource(src) }
}

// ClockStubBenchMode disables call recording. Dispatch still works; only
// the accumulating call log is dropped.
func ClockStubBenchMode() ClockStubOption {
	return func(s *ClockStub) { s.group.BenchMode() }
}

// ClockStubWithLatency makes every method sleep d before answering, for a
// test standing in for a slow dependency rather than one slow method. Pair
// it with ClockStubWithClock, or d elapses in wall time.
func ClockStubWithLatency(d time.Duration) ClockStubOption {
	return func(s *ClockStub) { s.group.SetLatency(d) }
}

// WithClockNow sets Now's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithClockNow(fn func() time.Time) ClockStubOption {
	return func(s *ClockStub) { s.OnNow.Func(fn) }
}

// WithClockZone sets Zone's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithClockZone(fn func() string) ClockStubOption {
	return func(s *ClockStub) { s.OnZone.Func(fn) }
}

// WithClockEvery sets Every's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithClockEvery(fn func(time.Duration, func()) func()) ClockStubOption {
	return func(s *ClockStub) { s.OnEvery.Func(fn) }
}

// ClockStub is a recording test double for Clock.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method returns its zero value and records the call.
type ClockStub struct {
	OnNow   *ClockNowStub
	OnZone  *ClockZoneStub
	OnEvery *ClockEveryStub

	// group is every method stub above, viewed through the surface that
	// does not depend on a signature.
	group  *stub.Group
	strict bool
}

// NewClockStub returns a double bound to tb.
//
// Passing tb registers a cleanup that verifies every method's call-count
// expectations when the test ends, so an unmet Times is reported without
// the caller remembering to check. A nil tb skips that, which is what
// benchmarks and non-test callers want.
func NewClockStub(tb assert.TB, opts ...ClockStubOption) *ClockStub {
	s := &ClockStub{
		OnNow:   &ClockNowStub{MethodStub: stub.NewMethodStub[ClockNowCall](tb, "Clock.Now")},
		OnZone:  &ClockZoneStub{MethodStub: stub.NewMethodStub[ClockZoneCall](tb, "Clock.Zone")},
		OnEvery: &ClockEveryStub{MethodStub: stub.NewMethodStub[ClockEveryCall](tb, "Clock.Every")},
	}
	s.group = stub.NewGroup(
		s.OnNow.MethodStub,
		s.OnZone.MethodStub,
		s.OnEvery.MethodStub,
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

// NewClockStubDelayed returns a double whose every method sleeps d against
// clk before answering — the clock and the latency a slow-dependency test
// wants, in one call rather than two options. Further options apply after
// the pair, so a caller can still add Strict or a method override.
func NewClockStubDelayed(tb assert.TB, d time.Duration, clk assert.Clock, opts ...ClockStubOption) *ClockStub {
	return NewClockStub(tb, append([]ClockStubOption{
		ClockStubWithClock(clk),
		ClockStubWithLatency(d),
	}, opts...)...)
}

// ResetCalls clears recorded calls, fault counters, and call-count
// expectations. Func and Returns configuration is deliberately preserved, so
// a test can reuse one double across phases without rebuilding it.
func (s *ClockStub) ResetCalls() {
	s.group.Reset()
}

// --- Interface methods ---

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *ClockNowStub) invoke() func() ClockNowReturn {
	if s.fn == nil {
		return nil
	}
	return func() ClockNowReturn {
		r0 := s.fn()
		return ClockNowReturn{Result: r0}
	}
}

// Now records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *ClockStub) Now() time.Time {
	call := ClockNowCall{}
	r := stub.Answer(s.OnNow.MethodStub, &call, stub.Arms[ClockNowCall, ClockNowReturn]{
		Invoke:   s.OnNow.invoke(),
		Fallback: s.OnNow.fallback,
		Stamp: func(c *ClockNowCall, r ClockNowReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *ClockZoneStub) invoke() func() ClockZoneReturn {
	if s.fn == nil {
		return nil
	}
	return func() ClockZoneReturn {
		r0 := s.fn()
		return ClockZoneReturn{Result: r0}
	}
}

// Zone records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *ClockStub) Zone() string {
	call := ClockZoneCall{}
	r := stub.Answer(s.OnZone.MethodStub, &call, stub.Arms[ClockZoneCall, ClockZoneReturn]{
		Invoke:   s.OnZone.invoke(),
		Fallback: s.OnZone.fallback,
		Stamp: func(c *ClockZoneCall, r ClockZoneReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *ClockEveryStub) invoke(d time.Duration, fn func()) func() ClockEveryReturn {
	if s.fn == nil {
		return nil
	}
	return func() ClockEveryReturn {
		stop := s.fn(d, fn)
		return ClockEveryReturn{Stop: stop}
	}
}

// Every records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *ClockStub) Every(d time.Duration, fn func()) (stop func()) {
	call := ClockEveryCall{D: d, Fn: fn}
	r := stub.Answer(s.OnEvery.MethodStub, &call, stub.Arms[ClockEveryCall, ClockEveryReturn]{
		Invoke:   s.OnEvery.invoke(d, fn),
		Fallback: s.OnEvery.fallback,
		Stamp: func(c *ClockEveryCall, r ClockEveryReturn) {
			c.Stop = r.Stop
		},
	})
	return r.Stop
}
