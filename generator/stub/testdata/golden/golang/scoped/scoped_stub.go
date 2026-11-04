package scoped

import (
	rand2 "math/rand/v2"
	time2 "time"

	assert2 "go.dokimi.dev/assert"
	fault2 "go.dokimi.dev/fault"
	stub2 "go.dokimi.dev/stub"
)

// StoreGetCall records one call of Store.Get.
type StoreGetCall struct {
	ID     string
	Result string
	Err    error
}

// StoreListCall records one call of Store.List.
type StoreListCall struct {
	Result []string
	Err    error
}

// StoreGetStub configures Get and records its calls. A fault armed through
// it takes precedence over the function and the fixed results.
type StoreGetStub struct {
	*stub2.MethodStub[StoreGetCall]

	fn       func(string) (string, error)
	fallback *StoreGetReturn
}

// StoreGetReturn contains the fixed results of Get.
type StoreGetReturn struct {
	Result string
	Err    error
}

// Returns sets the results of every call of Get. A function and an
// injected fault take precedence over them.
func (s *StoreGetStub) Returns(result string, err error) *StoreGetStub {
	s.fallback = &StoreGetReturn{Result: result, Err: err}
	return s
}

// Func sets the function that computes the results of Get from its
// arguments. An injected fault takes precedence over it.
func (s *StoreGetStub) Func(fn func(string) (string, error)) *StoreGetStub {
	s.fn = fn
	return s
}

// FaultForPartition makes Get fail every nth call whose ID matches key.
// The fault leaves the calls of the other partitions alone. A test of the
// isolation of a partition needs a fault that selects its calls.
func (s *StoreGetStub) FaultForPartition(key string, err error, n int) *StoreGetStub {
	s.FaultsWhen(func(c StoreGetCall) bool { return c.ID == key }, err, n)
	return s
}

// FaultForOtherPartitions makes Get fail every nth call whose ID differs
// from key. The fault leaves the calls of the partition of key alone.
func (s *StoreGetStub) FaultForOtherPartitions(key string, err error, n int) *StoreGetStub {
	s.FaultsWhen(func(c StoreGetCall) bool { return c.ID != key }, err, n)
	return s
}

// StoreListStub configures List and records its calls. A fault armed
// through it takes precedence over the function and the fixed results.
type StoreListStub struct {
	*stub2.MethodStub[StoreListCall]

	fn       func() ([]string, error)
	fallback *StoreListReturn
}

// StoreListReturn contains the fixed results of List.
type StoreListReturn struct {
	Result []string
	Err    error
}

// Returns sets the results of every call of List. A function and an
// injected fault take precedence over them.
func (s *StoreListStub) Returns(result []string, err error) *StoreListStub {
	s.fallback = &StoreListReturn{Result: result, Err: err}
	return s
}

// Func sets the function that computes the results of List from its
// arguments. An injected fault takes precedence over it.
func (s *StoreListStub) Func(fn func() ([]string, error)) *StoreListStub {
	s.fn = fn
	return s
}

// RetrySchedule makes List fail with err on every call before call 3. Call
// 3 and every call after it succeed.
func (s *StoreListStub) RetrySchedule(err error) *StoreListStub {
	s.SetFault(fault2.NewRetryFault[StoreListCall](err, 3))
	return s
}

// The reference to stub2.CompatV1 fails the build against another major
// version of the runtime of doubles.
var _ = stub2.CompatV1

// The reference to fault2.CompatV1 fails the build against another major
// version of the fault strategies.
var _ = fault2.CompatV1

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
		s.OnList.Func(impl.List)
	}
}

// StoreStubWithClock returns the option that sets the clock of the double
// to clk. The double takes the latency, the windows of the faults and the
// instants of the calls from clk. A latency on a controlled clock costs no
// real time.
func StoreStubWithClock(clk assert2.Clock) StoreStubOption {
	return func(s *StoreStub) { s.group.SetClock(clk) }
}

// StoreStubWithRandSource returns the option that sets the random source
// of the probabilistic faults to src. A seeded source repeats the faults
// of a run.
func StoreStubWithRandSource(src rand2.Source) StoreStubOption {
	return func(s *StoreStub) { s.group.SetRandSource(src) }
}

// StoreStubBenchMode returns the option that turns the recording of calls
// off. The methods still return their results.
func StoreStubBenchMode() StoreStubOption {
	return func(s *StoreStub) { s.group.BenchMode() }
}

// StoreStubWithLatency returns the option that makes every method sleep d
// before it returns. Without StoreStubWithClock, d passes in real time.
func StoreStubWithLatency(d time2.Duration) StoreStubOption {
	return func(s *StoreStub) { s.group.SetLatency(d) }
}

// StoreStubFailing returns the option that arms err on each method that
// can fail. Each call of such a method then fails with err.
func StoreStubFailing(err error) StoreStubOption {
	return func(s *StoreStub) {
		s.OnGet.Faults(err, 1)
		s.OnList.Faults(err, 1)
	}
}

// WithStoreGet returns the option that sets the function of Get.
func WithStoreGet(fn func(string) (string, error)) StoreStubOption {
	return func(s *StoreStub) { s.OnGet.Func(fn) }
}

// WithStoreList returns the option that sets the function of List.
func WithStoreList(fn func() ([]string, error)) StoreStubOption {
	return func(s *StoreStub) { s.OnList.Func(fn) }
}

// StoreStub is a recording test double of Store. Each field On<Method>
// configures its method. An unconfigured method records the call and
// returns the zero values of its results.
type StoreStub struct {
	OnGet  *StoreGetStub
	OnList *StoreListStub

	// group contains the stub of each method.
	group  *stub2.Group
	strict bool
}

// The compiler checks that StoreStub implements Store.
var _ Store = (*StoreStub)(nil)

// NewStoreStub returns a double of Store that reports its failures to tb
// and applies opts. The double checks the expected counts of calls at the
// end of the test. With a nil tb, the double reports nothing.
func NewStoreStub(tb assert2.TB, opts ...StoreStubOption) *StoreStub {
	s := &StoreStub{
		OnGet:  &StoreGetStub{MethodStub: stub2.NewMethodStub[StoreGetCall](tb, "Store.Get")},
		OnList: &StoreListStub{MethodStub: stub2.NewMethodStub[StoreListCall](tb, "Store.List")},
	}
	s.group = stub2.NewGroup(
		s.OnGet.MethodStub,
		s.OnList.MethodStub,
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

// ResetCalls clears the recorded calls, the counters of the faults and the
// expected counts. It keeps the functions, the fixed results and the faults.
func (s *StoreStub) ResetCalls() {
	s.group.Reset()
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *StoreGetStub) invoke(id string) func() StoreGetReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreGetReturn {
		r0, r1 := s.fn(id)
		return StoreGetReturn{Result: r0, Err: r1}
	}
}

// Get records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *StoreStub) Get(id string) (string, error) {
	call := StoreGetCall{ID: id}
	r := stub2.Answer(s.OnGet.MethodStub, call, stub2.Arms[StoreGetCall, StoreGetReturn]{
		Invoke:   s.OnGet.invoke(id),
		Fallback: s.OnGet.fallback,
		Fault:    func(err error) StoreGetReturn { return StoreGetReturn{Err: err} },
		Stamp: func(c StoreGetCall, r StoreGetReturn) StoreGetCall {
			c.Result = r.Result
			c.Err = r.Err
			return c
		},
	})
	return r.Result, r.Err
}

// invoke returns a call of the function of Func with the arguments, and nil
// where no function is set.
func (s *StoreListStub) invoke() func() StoreListReturn {
	if s.fn == nil {
		return nil
	}
	return func() StoreListReturn {
		r0, r1 := s.fn()
		return StoreListReturn{Result: r0, Err: r1}
	}
}

// List records the call. An armed fault comes first, then the function of
// Func, then the results of Returns, and then the zero values.
func (s *StoreStub) List() ([]string, error) {
	call := StoreListCall{}
	r := stub2.Answer(s.OnList.MethodStub, call, stub2.Arms[StoreListCall, StoreListReturn]{
		Invoke:   s.OnList.invoke(),
		Fallback: s.OnList.fallback,
		Fault:    func(err error) StoreListReturn { return StoreListReturn{Err: err} },
		Stamp: func(c StoreListCall, r StoreListReturn) StoreListCall {
			c.Result = r.Result
			c.Err = r.Err
			return c
		},
	})
	return r.Result, r.Err
}
