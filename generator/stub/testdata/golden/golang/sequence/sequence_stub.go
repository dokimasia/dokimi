package sequence

import (
	"iter"
	rand "math/rand/v2"
	"time"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/stub"
)

// TableNamesCall records one invocation of Table.Names.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type TableNamesCall struct {
	Result iter.Seq[string]
}

// TableRowsCall records one invocation of Table.Rows.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type TableRowsCall struct {
	Limit  int
	Result iter.Seq2[Item, error]
}

// TableCountsCall records one invocation of Table.Counts.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type TableCountsCall struct {
	Result iter.Seq2[string, int]
}

// TableTailCall records one invocation of Table.Tail.
//
// Fields take their names from the source signature; a slot the source
// left unnamed or blank falls back to a positional name.
type TableTailCall struct {
	Result Lines
}

// --- Per-method configuration ---

// TableNamesStub controls how the double answers Names and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type TableNamesStub struct {
	*stub.MethodStub[TableNamesCall]

	fn       func() iter.Seq[string]
	fallback *TableNamesReturn
}

// TableNamesReturn holds the fixed answer configured through Returns.
type TableNamesReturn struct {
	Result iter.Seq[string]
}

// Returns pins a fixed result for every call to Names. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Names answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *TableNamesStub) Returns(result iter.Seq[string]) *TableNamesStub {
	s.fallback = &TableNamesReturn{Result: result}
	return s
}

// Func supplies a body for Names, for when the answer depends on the
// arguments. Nothing takes precedence over it: Names answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *TableNamesStub) Func(fn func() iter.Seq[string]) *TableNamesStub {
	s.fn = fn
	return s
}

// Yields makes Names return a sequence over items.
//
// Building a range-over-func by hand is a closure with a break-handling
// branch in it, and a test that gets that branch wrong reports a bug in
// itself rather than in the subject.
func (s *TableNamesStub) Yields(items ...string) *TableNamesStub {
	s.fn = func() iter.Seq[string] {
		return func(yield func(string) bool) {
			for _, v := range items {
				if !yield(v) {
					return
				}
			}
		}
	}
	return s
}

// TableRowsStub controls how the double answers Rows and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type TableRowsStub struct {
	*stub.MethodStub[TableRowsCall]

	fn       func(int) iter.Seq2[Item, error]
	fallback *TableRowsReturn
}

// TableRowsReturn holds the fixed answer configured through Returns.
type TableRowsReturn struct {
	Result iter.Seq2[Item, error]
}

// Returns pins a fixed result for every call to Rows. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Rows answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *TableRowsStub) Returns(result iter.Seq2[Item, error]) *TableRowsStub {
	s.fallback = &TableRowsReturn{Result: result}
	return s
}

// Func supplies a body for Rows, for when the answer depends on the
// arguments. Nothing takes precedence over it: Rows answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *TableRowsStub) Func(fn func(int) iter.Seq2[Item, error]) *TableRowsStub {
	s.fn = fn
	return s
}

// Yields makes Rows return a sequence over items.
//
// Building a range-over-func by hand is a closure with a break-handling
// branch in it, and a test that gets that branch wrong reports a bug in
// itself rather than in the subject.
func (s *TableRowsStub) Yields(items ...Item) *TableRowsStub {
	s.fn = func(_ int) iter.Seq2[Item, error] {
		return func(yield func(Item, error) bool) {
			for _, v := range items {
				if !yield(v, nil) {
					return
				}
			}
		}
	}
	return s
}

// YieldsError makes Rows return a sequence that produces items and then
// fails.
//
// Failing partway through is the case a consumer's loop most often gets
// wrong, and it cannot be reached by a sequence that only ever succeeds.
func (s *TableRowsStub) YieldsError(items []Item, err error) *TableRowsStub {
	s.fn = func(_ int) iter.Seq2[Item, error] {
		return func(yield func(Item, error) bool) {
			for _, v := range items {
				if !yield(v, nil) {
					return
				}
			}
			var zero Item
			yield(zero, err)
		}
	}
	return s
}

// TableCountsStub controls how the double answers Counts and records what
// it was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type TableCountsStub struct {
	*stub.MethodStub[TableCountsCall]

	fn       func() iter.Seq2[string, int]
	fallback *TableCountsReturn
}

// TableCountsReturn holds the fixed answer configured through Returns.
type TableCountsReturn struct {
	Result iter.Seq2[string, int]
}

// Returns pins a fixed result for every call to Counts. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Counts answers no error, so the
// double is built with no fault arm and a fault armed through Faults fails
// the test.
func (s *TableCountsStub) Returns(result iter.Seq2[string, int]) *TableCountsStub {
	s.fallback = &TableCountsReturn{Result: result}
	return s
}

// Func supplies a body for Counts, for when the answer depends on the
// arguments. Nothing takes precedence over it: Counts answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *TableCountsStub) Func(fn func() iter.Seq2[string, int]) *TableCountsStub {
	s.fn = fn
	return s
}

// Yields makes Counts return a sequence over items.
//
// Building a range-over-func by hand is a closure with a break-handling
// branch in it, and a test that gets that branch wrong reports a bug in
// itself rather than in the subject.
func (s *TableCountsStub) Yields(items ...string) *TableCountsStub {
	s.fn = func() iter.Seq2[string, int] {
		return func(yield func(string, int) bool) {
			var second int
			for _, v := range items {
				if !yield(v, second) {
					return
				}
			}
		}
	}
	return s
}

// TableTailStub controls how the double answers Tail and records what it
// was asked.
//
// The embedded MethodStub supplies the machinery every method shares: call
// recording, fault injection, latency against a virtual clock, gates,
// call-count expectations, and strict mode.
type TableTailStub struct {
	*stub.MethodStub[TableTailCall]

	fn       func() Lines
	fallback *TableTailReturn
}

// TableTailReturn holds the fixed answer configured through Returns.
type TableTailReturn struct {
	Result Lines
}

// Returns pins a fixed result for every call to Tail. A Func override
// takes precedence over it.
//
// No fault can take precedence here: Tail answers no error, so the double
// is built with no fault arm and a fault armed through Faults fails the
// test.
func (s *TableTailStub) Returns(result Lines) *TableTailStub {
	s.fallback = &TableTailReturn{Result: result}
	return s
}

// Func supplies a body for Tail, for when the answer depends on the
// arguments. Nothing takes precedence over it: Tail answers no error, so
// the double carries no fault arm for one to be reported through.
func (s *TableTailStub) Func(fn func() Lines) *TableTailStub {
	s.fn = fn
	return s
}

// Yields makes Tail return a sequence over items.
//
// Building a range-over-func by hand is a closure with a break-handling
// branch in it, and a test that gets that branch wrong reports a bug in
// itself rather than in the subject.
func (s *TableTailStub) Yields(items ...string) *TableTailStub {
	s.fn = func() Lines {
		return func(yield func(string) bool) {
			for _, v := range items {
				if !yield(v) {
					return
				}
			}
		}
	}
	return s
}

// --- TableStub ---

// The witness of the major version of the runtime of doubles.
var _ = stub.CompatV1

// TableStubOption configures a [TableStub] at construction time.
type TableStubOption func(*TableStub)

// TableStubStrict makes an unconfigured method fail the test rather than
// return its zero value.
func TableStubStrict() TableStubOption {
	return func(s *TableStub) { s.strict = true }
}

// TableStubDelegateTo wraps a real implementation: every method forwards
// to impl and is recorded on the double on the way through.
func TableStubDelegateTo(impl Table) TableStubOption {
	return func(s *TableStub) {
		s.OnNames.Func(impl.Names)
		s.OnRows.Func(impl.Rows)
		s.OnCounts.Func(impl.Counts)
		s.OnTail.Func(impl.Tail)
	}
}

// TableStubWithClock drives latency and time-windowed faults from clk
// rather than from wall time, so a test asserting on a five-second timeout
// does not take five seconds.
func TableStubWithClock(clk assert.Clock) TableStubOption {
	return func(s *TableStub) { s.group.SetClock(clk) }
}

// TableStubWithRandSource makes probabilistic fault injection
// reproducible. A failure nobody can replay is not a finding.
func TableStubWithRandSource(src rand.Source) TableStubOption {
	return func(s *TableStub) { s.group.SetRandSource(src) }
}

// TableStubBenchMode disables call recording. Dispatch still works; only
// the accumulating call log is dropped.
func TableStubBenchMode() TableStubOption {
	return func(s *TableStub) { s.group.BenchMode() }
}

// TableStubWithLatency makes every method sleep d before answering, for a
// test standing in for a slow dependency rather than one slow method. Pair
// it with TableStubWithClock, or d elapses in wall time.
func TableStubWithLatency(d time.Duration) TableStubOption {
	return func(s *TableStub) { s.group.SetLatency(d) }
}

// WithTableNames sets Names's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithTableNames(fn func() iter.Seq[string]) TableStubOption {
	return func(s *TableStub) { s.OnNames.Func(fn) }
}

// WithTableRows sets Rows's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithTableRows(fn func(int) iter.Seq2[Item, error]) TableStubOption {
	return func(s *TableStub) { s.OnRows.Func(fn) }
}

// WithTableCounts sets Counts's body at construction time, for the common
// case of configuring one method and taking the defaults for the rest.
func WithTableCounts(fn func() iter.Seq2[string, int]) TableStubOption {
	return func(s *TableStub) { s.OnCounts.Func(fn) }
}

// WithTableTail sets Tail's body at construction time, for the common case
// of configuring one method and taking the defaults for the rest.
func WithTableTail(fn func() Lines) TableStubOption {
	return func(s *TableStub) { s.OnTail.Func(fn) }
}

// TableStub is a recording test double for Table.
//
// Each On<Method> field is that method's configuration point. Left alone,
// a method returns its zero value and records the call.
type TableStub struct {
	OnNames  *TableNamesStub
	OnRows   *TableRowsStub
	OnCounts *TableCountsStub
	OnTail   *TableTailStub

	// group is every method stub above, viewed through the surface that
	// does not depend on a signature.
	group  *stub.Group
	strict bool
}

// Compile-time proof that TableStub satisfies Table.
var _ Table = (*TableStub)(nil)

// NewTableStub returns a double bound to tb.
//
// Passing tb registers a cleanup that verifies every method's call-count
// expectations when the test ends, so an unmet Times is reported without
// the caller remembering to check. A nil tb skips that, which is what
// benchmarks and non-test callers want.
func NewTableStub(tb assert.TB, opts ...TableStubOption) *TableStub {
	s := &TableStub{
		OnNames:  &TableNamesStub{MethodStub: stub.NewMethodStub[TableNamesCall](tb, "Table.Names")},
		OnRows:   &TableRowsStub{MethodStub: stub.NewMethodStub[TableRowsCall](tb, "Table.Rows")},
		OnCounts: &TableCountsStub{MethodStub: stub.NewMethodStub[TableCountsCall](tb, "Table.Counts")},
		OnTail:   &TableTailStub{MethodStub: stub.NewMethodStub[TableTailCall](tb, "Table.Tail")},
	}
	s.group = stub.NewGroup(
		s.OnNames.MethodStub,
		s.OnRows.MethodStub,
		s.OnCounts.MethodStub,
		s.OnTail.MethodStub,
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

// NewTableStubDelayed returns a double whose every method sleeps d against
// clk before answering — the clock and the latency a slow-dependency test
// wants, in one call rather than two options. Further options apply after
// the pair, so a caller can still add Strict or a method override.
func NewTableStubDelayed(tb assert.TB, d time.Duration, clk assert.Clock, opts ...TableStubOption) *TableStub {
	return NewTableStub(tb, append([]TableStubOption{
		TableStubWithClock(clk),
		TableStubWithLatency(d),
	}, opts...)...)
}

// ResetCalls clears recorded calls, fault counters, and call-count
// expectations. Func and Returns configuration is deliberately preserved, so
// a test can reuse one double across phases without rebuilding it.
func (s *TableStub) ResetCalls() {
	s.group.Reset()
}

// --- Interface methods ---

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *TableNamesStub) invoke() func() TableNamesReturn {
	if s.fn == nil {
		return nil
	}
	return func() TableNamesReturn {
		r0 := s.fn()
		return TableNamesReturn{Result: r0}
	}
}

// Names records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *TableStub) Names() iter.Seq[string] {
	call := TableNamesCall{}
	r := stub.Answer(s.OnNames.MethodStub, &call, stub.Arms[TableNamesCall, TableNamesReturn]{
		Invoke:   s.OnNames.invoke(),
		Fallback: s.OnNames.fallback,
		Stamp: func(c *TableNamesCall, r TableNamesReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *TableRowsStub) invoke(limit int) func() TableRowsReturn {
	if s.fn == nil {
		return nil
	}
	return func() TableRowsReturn {
		r0 := s.fn(limit)
		return TableRowsReturn{Result: r0}
	}
}

// Rows records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *TableStub) Rows(limit int) iter.Seq2[Item, error] {
	call := TableRowsCall{Limit: limit}
	r := stub.Answer(s.OnRows.MethodStub, &call, stub.Arms[TableRowsCall, TableRowsReturn]{
		Invoke:   s.OnRows.invoke(limit),
		Fallback: s.OnRows.fallback,
		Stamp: func(c *TableRowsCall, r TableRowsReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *TableCountsStub) invoke() func() TableCountsReturn {
	if s.fn == nil {
		return nil
	}
	return func() TableCountsReturn {
		r0 := s.fn()
		return TableCountsReturn{Result: r0}
	}
}

// Counts records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *TableStub) Counts() iter.Seq2[string, int] {
	call := TableCountsCall{}
	r := stub.Answer(s.OnCounts.MethodStub, &call, stub.Arms[TableCountsCall, TableCountsReturn]{
		Invoke:   s.OnCounts.invoke(),
		Fallback: s.OnCounts.fallback,
		Stamp: func(c *TableCountsCall, r TableCountsReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}

// invoke adapts the Func override to the shape [stub.Answer] consumes, or
// returns nil when no override is set — which is how Answer tells "no
// override" from "an override that returns zero".
func (s *TableTailStub) invoke() func() TableTailReturn {
	if s.fn == nil {
		return nil
	}
	return func() TableTailReturn {
		r0 := s.fn()
		return TableTailReturn{Result: r0}
	}
}

// Tail records the call and answers it.
//
// Which arm answers — injected fault, Func override, Returns fallback, or
// the zero value — is [stub.Answer]'s to decide, so every generated double
// resolves a call the same way and the ordering is tested once rather than
// restated per method.
func (s *TableStub) Tail() Lines {
	call := TableTailCall{}
	r := stub.Answer(s.OnTail.MethodStub, &call, stub.Arms[TableTailCall, TableTailReturn]{
		Invoke:   s.OnTail.invoke(),
		Fallback: s.OnTail.fallback,
		Stamp: func(c *TableTailCall, r TableTailReturn) {
			c.Result = r.Result
		},
	})
	return r.Result
}
