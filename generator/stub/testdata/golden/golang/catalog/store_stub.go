package catalog

import (
	"context"
	rand "math/rand/v2"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// StoreGetCall records one call of Store.Get.
type StoreGetCall struct {
	Ctx    context.Context
	ID     string
	Result Item
	Err    error
}

// StorePutCall records one call of Store.Put.
type StorePutCall struct {
	A0  context.Context
	A1  Item
	Err error
}

// StoreLenCall records one call of Store.Len.
type StoreLenCall struct {
	Result int
}

// StoreLookupCall records one call of Store.Lookup.
type StoreLookupCall struct {
	Key  string
	Item Item
	Ok   bool
}

// StoreLogCall records one call of Store.Log.
type StoreLogCall struct {
	Format string
	Args   []string
}

// StoreSplitCall records one call of Store.Split.
type StoreSplitCall struct {
	S       string
	Result0 string
	Result1 string
	Err     error
}

// StoreCloseCall records one call of Store.Close.
type StoreCloseCall struct {
	Err error
}

// StoreGetStub configures Get and records its calls. A fault armed through
// it takes precedence over the function and the fixed results.
type StoreGetStub struct {
	*stub.MethodStub[StoreGetCall]

	fn       func(context.Context, string) (Item, error)
	fallback *StoreGetReturn
}

// StoreGetReturn contains the fixed results of Get.
type StoreGetReturn struct {
	Result Item
	Err    error
}

// Returns sets the results of every call of Get. A function and an
// injected fault take precedence over them.
func (s *StoreGetStub) Returns(result Item, err error) *StoreGetStub {
	s.fallback = &StoreGetReturn{Result: result, Err: err}
	return s
}

// Func sets the function that computes the results of Get from its
// arguments. An injected fault takes precedence over it.
func (s *StoreGetStub) Func(fn func(context.Context, string) (Item, error)) *StoreGetStub {
	s.fn = fn
	return s
}

// StorePutStub configures Put and records its calls. A fault armed through
// it takes precedence over the function and the fixed results.
type StorePutStub struct {
	*stub.MethodStub[StorePutCall]

	fn       func(context.Context, Item) error
	fallback *StorePutReturn
}

// StorePutReturn contains the fixed results of Put.
type StorePutReturn struct {
	Err error
}

// Returns sets the results of every call of Put. A function and an
// injected fault take precedence over them.
func (s *StorePutStub) Returns(err error) *StorePutStub {
	s.fallback = &StorePutReturn{Err: err}
	return s
}

// Func sets the function that computes the results of Put from its
// arguments. An injected fault takes precedence over it.
func (s *StorePutStub) Func(fn func(context.Context, Item) error) *StorePutStub {
	s.fn = fn
	return s
}

// StoreLenStub configures Len and records its calls. Len returns no error,
// so a fault armed through it fails the test.
type StoreLenStub struct {
	*stub.MethodStub[StoreLenCall]

	fn       func() int
	fallback *StoreLenReturn
}

// StoreLenReturn contains the fixed results of Len.
type StoreLenReturn struct {
	Result int
}

// Returns sets the results of every call of Len. A function takes
// precedence over them.
func (s *StoreLenStub) Returns(result int) *StoreLenStub {
	s.fallback = &StoreLenReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Len from its
// arguments.
func (s *StoreLenStub) Func(fn func() int) *StoreLenStub {
	s.fn = fn
	return s
}

// StoreLookupStub configures Lookup and records its calls. Lookup returns
// no error, so a fault armed through it fails the test.
type StoreLookupStub struct {
	*stub.MethodStub[StoreLookupCall]

	fn       func(string) (Item, bool)
	fallback *StoreLookupReturn
}

// StoreLookupReturn contains the fixed results of Lookup.
type StoreLookupReturn struct {
	Item Item
	Ok   bool
}

// Returns sets the results of every call of Lookup. A function takes
// precedence over them.
func (s *StoreLookupStub) Returns(item Item, ok bool) *StoreLookupStub {
	s.fallback = &StoreLookupReturn{Item: item, Ok: ok}
	return s
}

// Func sets the function that computes the results of Lookup from its
// arguments.
func (s *StoreLookupStub) Func(fn func(string) (Item, bool)) *StoreLookupStub {
	s.fn = fn
	return s
}

// StoreLogStub configures Log and records its calls. Log returns no error,
// so a fault armed through it fails the test.
type StoreLogStub struct {
	*stub.MethodStub[StoreLogCall]

	fn func(string, ...string)
}

// Func sets the function that computes the results of Log from its
// arguments.
func (s *StoreLogStub) Func(fn func(string, ...string)) *StoreLogStub {
	s.fn = fn
	return s
}

// StoreSplitStub configures Split and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type StoreSplitStub struct {
	*stub.MethodStub[StoreSplitCall]

	fn       func(string) (string, string, error)
	fallback *StoreSplitReturn
}

// StoreSplitReturn contains the fixed results of Split.
type StoreSplitReturn struct {
	Result0 string
	Result1 string
	Err     error
}

// Returns sets the results of every call of Split. A function and an
// injected fault take precedence over them.
func (s *StoreSplitStub) Returns(result0, result1 string, err error) *StoreSplitStub {
	s.fallback = &StoreSplitReturn{Result0: result0, Result1: result1, Err: err}
	return s
}

// Func sets the function that computes the results of Split from its
// arguments. An injected fault takes precedence over it.
func (s *StoreSplitStub) Func(fn func(string) (string, string, error)) *StoreSplitStub {
	s.fn = fn
	return s
}

// StoreCloseStub configures Close and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type StoreCloseStub struct {
	*stub.MethodStub[StoreCloseCall]

	fn       func() error
	fallback *StoreCloseReturn
}

// StoreCloseReturn contains the fixed results of Close.
type StoreCloseReturn struct {
	Err error
}

// Returns sets the results of every call of Close. A function and an
// injected fault take precedence over them.
func (s *StoreCloseStub) Returns(err error) *StoreCloseStub {
	s.fallback = &StoreCloseReturn{Err: err}
	return s
}

// Func sets the function that computes the results of Close from its
// arguments. An injected fault takes precedence over it.
func (s *StoreCloseStub) Func(fn func() error) *StoreCloseStub {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// StoreStubOption configures a [StoreStub] at its construction.
type StoreStubOption func(*StoreStub)

// StoreStubStrict returns the option that makes an unconfigured method
// fail the test instead of returning its zero value.
func StoreStubStrict() StoreStubOption {
	return func(s *StoreStub) { s.strict = true }
}

// StoreStubDelegateTo returns the option that passes each call to impl.
// The double still records the call.
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

// StoreStubWithClock returns the option that sets the clock of the double
// to clk. The double takes the latency, the windows of the faults and the
// instants of the calls from clk. A latency on a controlled clock costs no
// real time.
func StoreStubWithClock(clk assert.Clock) StoreStubOption {
	return func(s *StoreStub) { s.group.SetClock(clk) }
}

// StoreStubWithRandSource returns the option that sets the random source
// of the probabilistic faults to src. A seeded source repeats the faults
// of a run.
func StoreStubWithRandSource(src rand.Source) StoreStubOption {
	return func(s *StoreStub) { s.group.SetRandSource(src) }
}

// StoreStubBenchMode returns the option that turns the recording of calls
// off. The methods still return their results.
func StoreStubBenchMode() StoreStubOption {
	return func(s *StoreStub) { s.group.BenchMode() }
}

// StoreStubWithLatency returns the option that makes every method sleep d
// before it returns. Without StoreStubWithClock, d passes in real time.
func StoreStubWithLatency(d time.Duration) StoreStubOption {
	return func(s *StoreStub) { s.group.SetLatency(d) }
}

// StoreStubFailing returns the option that arms err on each method that
// can fail. Each call of such a method then fails with err.
func StoreStubFailing(err error) StoreStubOption {
	return func(s *StoreStub) {
		s.OnGet.Faults(err, 1)
		s.OnPut.Faults(err, 1)
		s.OnSplit.Faults(err, 1)
		s.OnClose.Faults(err, 1)
	}
}

// WithStoreGet returns the option that sets the function of Get.
func WithStoreGet(fn func(context.Context, string) (Item, error)) StoreStubOption {
	return func(s *StoreStub) { s.OnGet.Func(fn) }
}

// WithStorePut returns the option that sets the function of Put.
func WithStorePut(fn func(context.Context, Item) error) StoreStubOption {
	return func(s *StoreStub) { s.OnPut.Func(fn) }
}

// WithStoreLen returns the option that sets the function of Len.
func WithStoreLen(fn func() int) StoreStubOption {
	return func(s *StoreStub) { s.OnLen.Func(fn) }
}

// WithStoreLookup returns the option that sets the function of Lookup.
func WithStoreLookup(fn func(string) (Item, bool)) StoreStubOption {
	return func(s *StoreStub) { s.OnLookup.Func(fn) }
}

// WithStoreLog returns the option that sets the function of Log.
func WithStoreLog(fn func(string, ...string)) StoreStubOption {
	return func(s *StoreStub) { s.OnLog.Func(fn) }
}

// WithStoreSplit returns the option that sets the function of Split.
func WithStoreSplit(fn func(string) (string, string, error)) StoreStubOption {
	return func(s *StoreStub) { s.OnSplit.Func(fn) }
}

// WithStoreClose returns the option that sets the function of Close.
func WithStoreClose(fn func() error) StoreStubOption {
	return func(s *StoreStub) { s.OnClose.Func(fn) }
}

// StoreStub is a recording test double of Store. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type StoreStub struct {
	OnGet    *StoreGetStub
	OnPut    *StorePutStub
	OnLen    *StoreLenStub
	OnLookup *StoreLookupStub
	OnLog    *StoreLogStub
	OnSplit  *StoreSplitStub
	// OnClose configures Close, which Closer declares.
	OnClose *StoreCloseStub

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// The compiler checks that StoreStub implements Store.
var _ Store = (*StoreStub)(nil)

// NewStoreStub returns a double of Store that reports its failures to tb
// and applies opts. The double checks the expected counts of calls at the
// end of the test. With a nil tb, the double reports nothing.
func NewStoreStub(tb assert.TB, opts ...StoreStubOption) *StoreStub {
	s := &StoreStub{
		OnGet: &StoreGetStub{MethodStub: stub.NewMethodStub[StoreGetCall](tb, "Store.Get")},
		OnPut: &StorePutStub{MethodStub: stub.NewMethodStub[StorePutCall](tb, "Store.Put")},
		OnLen: &StoreLenStub{MethodStub: stub.NewMethodStub[StoreLenCall](tb, "Store.Len")},
		OnLookup: &StoreLookupStub{
			MethodStub: stub.NewMethodStub[StoreLookupCall](tb, "Store.Lookup"),
		},
		OnLog: &StoreLogStub{MethodStub: stub.NewMethodStub[StoreLogCall](tb, "Store.Log")},
		OnSplit: &StoreSplitStub{
			MethodStub: stub.NewMethodStub[StoreSplitCall](tb, "Store.Split"),
		},
		OnClose: &StoreCloseStub{
			MethodStub: stub.NewMethodStub[StoreCloseCall](tb, "Store.Close"),
		},
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

// NewStoreStubDelayed returns a double whose methods sleep d on clk before
// they return. It applies opts after the clock and the latency.
func NewStoreStubDelayed(
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...StoreStubOption,
) *StoreStub {
	return NewStoreStub(tb, append([]StoreStubOption{
		StoreStubWithClock(clk),
		StoreStubWithLatency(d),
	}, opts...)...)
}

// NewStoreStubFailing returns a double on which each call of a method that
// can fail fails with err. It applies opts after the failure.
func NewStoreStubFailing(tb assert.TB, err error, opts ...StoreStubOption) *StoreStub {
	return NewStoreStub(tb, append([]StoreStubOption{StoreStubFailing(err)}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *StoreStub) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *StoreGetStub) invoke(ctx context.Context, id string) func() StoreGetReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreGetReturn {
		r0, r1 := s.fn(ctx, id)
		return StoreGetReturn{Result: r0, Err: r1}
	}
}

// Get records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *StorePutStub) invoke(a0 context.Context, a1 Item) func() StorePutReturn {
	if s.fn == nil {
		return nil
	}
	return func() StorePutReturn {
		r0 := s.fn(a0, a1)
		return StorePutReturn{Err: r0}
	}
}

// Put records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *StoreLenStub) invoke() func() StoreLenReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreLenReturn {
		r0 := s.fn()
		return StoreLenReturn{Result: r0}
	}
}

// Len records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *StoreLookupStub) invoke(key string) func() StoreLookupReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreLookupReturn {
		item, ok := s.fn(key)
		return StoreLookupReturn{Item: item, Ok: ok}
	}
}

// Lookup records the call. An armed fault comes first, then the function
// of Func, then the results of Returns, and then the zero values.
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

// invokeVoid returns a call of the function of Func with the arguments, and
// nil where no function is set. The empty struct is the result of a method
// without results.
func (s *StoreLogStub) invokeVoid(format string, args ...string) func() struct{} {
	if s.fn == nil {
		return nil
	}
	return func() struct{} {
		s.fn(format, args...)
		return struct{}{}
	}
}

// Log records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *StoreStub) Log(format string, args ...string) {
	call := StoreLogCall{Format: format, Args: args}
	stub.Answer(s.OnLog.MethodStub, &call, stub.Arms[StoreLogCall, struct{}]{
		Invoke: s.OnLog.invokeVoid(format, args...),
	})
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s0 *StoreSplitStub) invoke(s string) func() StoreSplitReturn {
	if s0.fn == nil {
		return nil
	}
	return func() StoreSplitReturn {
		r0, r1, r2 := s0.fn(s)
		return StoreSplitReturn{Result0: r0, Result1: r1, Err: r2}
	}
}

// Split records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *StoreCloseStub) invoke() func() StoreCloseReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreCloseReturn {
		r0 := s.fn()
		return StoreCloseReturn{Err: r0}
	}
}

// Close records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// ClockNowCall records one call of Clock.Now.
type ClockNowCall struct {
	Result time.Time
}

// ClockZoneCall records one call of Clock.Zone.
type ClockZoneCall struct {
	Result string
}

// ClockEveryCall records one call of Clock.Every.
type ClockEveryCall struct {
	D    time.Duration
	Fn   func()
	Stop func()
}

// ClockNowStub configures Now and records its calls. Now returns no error,
// so a fault armed through it fails the test.
type ClockNowStub struct {
	*stub.MethodStub[ClockNowCall]

	fn       func() time.Time
	fallback *ClockNowReturn
}

// ClockNowReturn contains the fixed results of Now.
type ClockNowReturn struct {
	Result time.Time
}

// Returns sets the results of every call of Now. A function takes
// precedence over them.
func (s *ClockNowStub) Returns(result time.Time) *ClockNowStub {
	s.fallback = &ClockNowReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Now from its
// arguments.
func (s *ClockNowStub) Func(fn func() time.Time) *ClockNowStub {
	s.fn = fn
	return s
}

// ClockZoneStub configures Zone and records its calls. Zone returns no
// error, so a fault armed through it fails the test.
type ClockZoneStub struct {
	*stub.MethodStub[ClockZoneCall]

	fn       func() string
	fallback *ClockZoneReturn
}

// ClockZoneReturn contains the fixed results of Zone.
type ClockZoneReturn struct {
	Result string
}

// Returns sets the results of every call of Zone. A function takes
// precedence over them.
func (s *ClockZoneStub) Returns(result string) *ClockZoneStub {
	s.fallback = &ClockZoneReturn{Result: result}
	return s
}

// Func sets the function that computes the results of Zone from its
// arguments.
func (s *ClockZoneStub) Func(fn func() string) *ClockZoneStub {
	s.fn = fn
	return s
}

// ClockEveryStub configures Every and records its calls. Every returns no
// error, so a fault armed through it fails the test.
type ClockEveryStub struct {
	*stub.MethodStub[ClockEveryCall]

	fn       func(time.Duration, func()) func()
	fallback *ClockEveryReturn
}

// ClockEveryReturn contains the fixed results of Every.
type ClockEveryReturn struct {
	Stop func()
}

// Returns sets the results of every call of Every. A function takes
// precedence over them.
func (s *ClockEveryStub) Returns(stop func()) *ClockEveryStub {
	s.fallback = &ClockEveryReturn{Stop: stop}
	return s
}

// Func sets the function that computes the results of Every from its
// arguments.
func (s *ClockEveryStub) Func(fn func(time.Duration, func()) func()) *ClockEveryStub {
	s.fn = fn
	return s
}

// The reference to stub.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub.CompatV1

// ClockStubOption configures a [ClockStub] at its construction.
type ClockStubOption func(*ClockStub)

// ClockStubStrict returns the option that makes an unconfigured method
// fail the test instead of returning its zero value.
func ClockStubStrict() ClockStubOption {
	return func(s *ClockStub) { s.strict = true }
}

// ClockStubDelegateTo returns the option that passes each call to impl.
// impl has the methods of Clock. Both a pointer to a Clock and a second
// double have them. The double still records the call.
func ClockStubDelegateTo(impl interface {
	Now() time.Time
	Zone() string
	Every(d time.Duration, fn func()) (stop func())
},
) ClockStubOption {
	return func(s *ClockStub) {
		s.OnNow.Func(impl.Now)
		s.OnZone.Func(impl.Zone)
		s.OnEvery.Func(impl.Every)
	}
}

// ClockStubWithClock returns the option that sets the clock of the double
// to clk. The double takes the latency, the windows of the faults and the
// instants of the calls from clk. A latency on a controlled clock costs no
// real time.
func ClockStubWithClock(clk assert.Clock) ClockStubOption {
	return func(s *ClockStub) { s.group.SetClock(clk) }
}

// ClockStubWithRandSource returns the option that sets the random source
// of the probabilistic faults to src. A seeded source repeats the faults
// of a run.
func ClockStubWithRandSource(src rand.Source) ClockStubOption {
	return func(s *ClockStub) { s.group.SetRandSource(src) }
}

// ClockStubBenchMode returns the option that turns the recording of calls
// off. The methods still return their results.
func ClockStubBenchMode() ClockStubOption {
	return func(s *ClockStub) { s.group.BenchMode() }
}

// ClockStubWithLatency returns the option that makes every method sleep d
// before it returns. Without ClockStubWithClock, d passes in real time.
func ClockStubWithLatency(d time.Duration) ClockStubOption {
	return func(s *ClockStub) { s.group.SetLatency(d) }
}

// WithClockNow returns the option that sets the function of Now.
func WithClockNow(fn func() time.Time) ClockStubOption {
	return func(s *ClockStub) { s.OnNow.Func(fn) }
}

// WithClockZone returns the option that sets the function of Zone.
func WithClockZone(fn func() string) ClockStubOption {
	return func(s *ClockStub) { s.OnZone.Func(fn) }
}

// WithClockEvery returns the option that sets the function of Every.
func WithClockEvery(fn func(time.Duration, func()) func()) ClockStubOption {
	return func(s *ClockStub) { s.OnEvery.Func(fn) }
}

// ClockStub is a recording test double of Clock. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type ClockStub struct {
	OnNow   *ClockNowStub
	OnZone  *ClockZoneStub
	OnEvery *ClockEveryStub

	// group contains the stub of each method.
	group  *stub.Group
	strict bool
}

// NewClockStub returns a double of Clock that reports its failures to tb
// and applies opts. The double checks the expected counts of calls at the
// end of the test. With a nil tb, the double reports nothing.
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

// NewClockStubDelayed returns a double whose methods sleep d on clk before
// they return. It applies opts after the clock and the latency.
func NewClockStubDelayed(
	tb assert.TB,
	d time.Duration,
	clk assert.Clock,
	opts ...ClockStubOption,
) *ClockStub {
	return NewClockStub(tb, append([]ClockStubOption{
		ClockStubWithClock(clk),
		ClockStubWithLatency(d),
	}, opts...)...)
}

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *ClockStub) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *ClockNowStub) invoke() func() ClockNowReturn {
	if s.fn == nil {
		return nil
	}
	return func() ClockNowReturn {
		r0 := s.fn()
		return ClockNowReturn{Result: r0}
	}
}

// Now records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *ClockZoneStub) invoke() func() ClockZoneReturn {
	if s.fn == nil {
		return nil
	}
	return func() ClockZoneReturn {
		r0 := s.fn()
		return ClockZoneReturn{Result: r0}
	}
}

// Zone records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *ClockEveryStub) invoke(d time.Duration, fn func()) func() ClockEveryReturn {
	if s.fn == nil {
		return nil
	}
	return func() ClockEveryReturn {
		stop := s.fn(d, fn)
		return ClockEveryReturn{Stop: stop}
	}
}

// Every records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
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
