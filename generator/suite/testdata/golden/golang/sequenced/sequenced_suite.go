package sequenced

import (
	"context"
	"iter"
	"slices"
	"strings"
	"sync/atomic"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/assert/stateful"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Counter:
//
// Counter.Inc/keeps-within-its-bounds: Inc keeps the use that Value reads
// within its bounds, over any sequence of calls.
//
// Counter.Value/smoke: Value returns without a panic for derived inputs.
//
// Counter.Value/cancel: Value returns an error for a cancelled context.
//
// Counter.Value/zero-on-error: Value returns the zero value beside an
// error.
//
// Counter.Value/monotonic: Value never returns less than it returned
// before.
//
// Counter.Value/never-goes-back: Value never returns less than a result
// that it returned before, over any sequence of calls.
//
// Counter.Grow/smoke: Grow returns without a panic for derived inputs.
//
// Counter.Grow/cancel: Grow returns an error for a cancelled context.
//
// Counter.Grow/keeps-within-its-bounds: Grow keeps the use that Size reads
// within its bounds, over any sequence of calls.
//
// Counter.Size/smoke: Size returns without a panic for derived inputs.
//
// Counter.Size/cancel: Size returns an error for a cancelled context.
//
// Counter.Size/zero-on-error: Size returns the zero value beside an error.
//
// Counter.Throttle/smoke: Throttle returns without a panic for derived
// inputs.
//
// Counter.Throttle/cancel: Throttle returns an error for a cancelled
// context.
//
// Counter.Load/smoke: Load returns without a panic for derived inputs.
//
// Counter.Load/cancel: Load returns an error for a cancelled context.
//
// Counter.Load/zero-on-error: Load returns the zero value beside an error.

// The derived checks of Counter that witness their laws at the shallow
// end:
//
// Counter.Value/smoke
//
// Counter.Value/monotonic
//
// Counter.Grow/smoke
//
// Counter.Size/smoke
//
// Counter.Throttle/smoke
//
// Counter.Load/smoke

// The refused claims of Counter:
//
// Counter.Inc/cancel: eidos derives no sample of the parameter opts of
// Inc. The remedy is to author a sample of the parameter opts with the
// directive sample.
//
// Counter.Throttle/keeps-within-its-bounds: the reader that observe names
// for Throttle returns no number and no length. The remedy is to name with
// observe a callable that returns the use of Throttle as a number or a
// list.

// The unclaimed classifications of Counter:
//
// The shape:writer of Counter.Inc has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:lifecycle of Counter.Grow has the stance silent. A lifecycle
// call returns nothing that a check could compare.
//
// The shape:lifecycle of Counter.Throttle has the stance silent. A
// lifecycle call returns nothing that a check could compare.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Counter writes.
type (
	CounterHandle = conform.Handle[Counter]
	CounterCheck  = conform.Check[Counter]
	CounterOption = conform.Option[Counter]
)

// CounterChecks contains the name of each callable of Counter, with the
// identities of its checks.
var CounterChecks = struct {
	Inc      counterInc
	Value    counterValue
	Grow     counterGrow
	Size     counterSize
	Throttle counterThrottle
	Load     counterLoad
}{
	Inc:      counterInc{"Counter.Inc"},
	Value:    counterValue{"Counter.Value"},
	Grow:     counterGrow{"Counter.Grow"},
	Size:     counterSize{"Counter.Size"},
	Throttle: counterThrottle{"Counter.Throttle"},
	Load:     counterLoad{"Counter.Load"},
}

// counterInc is the name of Counter.Inc, with the identities of its
// checks.
type counterInc struct {
	conform.Callable[Counter]
}

// KeepsWithinItsBounds returns the identity of the check that Inc keeps
// the use that Value reads within its bounds, over any sequence of calls.
func (c counterInc) KeepsWithinItsBounds() conform.ID[Counter] {
	return c.Check("keeps-within-its-bounds")
}

// counterValue is the name of Counter.Value, with the identities of its
// checks.
type counterValue struct {
	conform.Callable[Counter]
}

// Smoke returns the identity of the check that Value returns without a
// panic for derived inputs.
func (c counterValue) Smoke() conform.ID[Counter] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Value returns an error for
// a cancelled context.
func (c counterValue) Cancel() conform.ID[Counter] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Value returns the
// zero value beside an error.
func (c counterValue) ZeroOnError() conform.ID[Counter] {
	return c.Check("zero-on-error")
}

// Monotonic returns the identity of the check that Value never returns
// less than it returned before.
func (c counterValue) Monotonic() conform.ID[Counter] {
	return c.Check("monotonic")
}

// NeverGoesBack returns the identity of the check that Value never returns
// less than a result that it returned before, over any sequence of calls.
func (c counterValue) NeverGoesBack() conform.ID[Counter] {
	return c.Check("never-goes-back")
}

// counterGrow is the name of Counter.Grow, with the identities of its
// checks.
type counterGrow struct {
	conform.Callable[Counter]
}

// Smoke returns the identity of the check that Grow returns without a
// panic for derived inputs.
func (c counterGrow) Smoke() conform.ID[Counter] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Grow returns an error for
// a cancelled context.
func (c counterGrow) Cancel() conform.ID[Counter] {
	return c.Check("cancel")
}

// KeepsWithinItsBounds returns the identity of the check that Grow keeps
// the use that Size reads within its bounds, over any sequence of calls.
func (c counterGrow) KeepsWithinItsBounds() conform.ID[Counter] {
	return c.Check("keeps-within-its-bounds")
}

// counterSize is the name of Counter.Size, with the identities of its
// checks.
type counterSize struct {
	conform.Callable[Counter]
}

// Smoke returns the identity of the check that Size returns without a
// panic for derived inputs.
func (c counterSize) Smoke() conform.ID[Counter] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Size returns an error for
// a cancelled context.
func (c counterSize) Cancel() conform.ID[Counter] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Size returns the zero
// value beside an error.
func (c counterSize) ZeroOnError() conform.ID[Counter] {
	return c.Check("zero-on-error")
}

// counterThrottle is the name of Counter.Throttle, with the identities of
// its checks.
type counterThrottle struct {
	conform.Callable[Counter]
}

// Smoke returns the identity of the check that Throttle returns without a
// panic for derived inputs.
func (c counterThrottle) Smoke() conform.ID[Counter] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Throttle returns an error
// for a cancelled context.
func (c counterThrottle) Cancel() conform.ID[Counter] {
	return c.Check("cancel")
}

// counterLoad is the name of Counter.Load, with the identities of its
// checks.
type counterLoad struct {
	conform.Callable[Counter]
}

// Smoke returns the identity of the check that Load returns without a
// panic for derived inputs.
func (c counterLoad) Smoke() conform.ID[Counter] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Load returns an error for
// a cancelled context.
func (c counterLoad) Cancel() conform.ID[Counter] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Load returns the zero
// value beside an error.
func (c counterLoad) ZeroOnError() conform.ID[Counter] {
	return c.Check("zero-on-error")
}

// CounterSuite returns the checks of Counter, and the calls and the
// budgets of its callables.
func CounterSuite() conform.Suite[Counter] {
	return conform.Suite[Counter]{
		Declaration: "Counter",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Counter]{
			{
				ID:          CounterChecks.Inc.KeepsWithinItsBounds(),
				Statement:   "Inc keeps the use that Value reads within its bounds, over any sequence of calls",
				Claim:       spec.KeepsWithinItsBounds,
				Probe:       spec.UseIsWithinTheBounds,
				Reason:      spec.ExceededBound,
				Law:         spec.UseIsWithinItsBounds,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *CounterHandle, c *prop.Case) {
					probe.UseIsWithinTheBounds(
						h,
						c,
						counterActions,
						func(ctx context.Context) (int, error) {
							values, err := counterValueCall(ctx, h.Subject(), []any{})
							v, _ := values[0].(int)
							return int(v), err
						},
						0,
						3,
					)
				},
				Defect: counterValueGrowsPastTheLimitAfterThree,
			},
			{
				ID:          CounterChecks.Value.Smoke(),
				Statement:   "Value returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CounterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Value(ctx)
						return err
					})
				},
				Defect: counterValuePanics,
			},
			{
				ID:          CounterChecks.Value.Cancel(),
				Statement:   "Value returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Value(ctx)
						return err
					})
				},
				Defect: counterValueIgnoresCancel,
			},
			{
				ID:          CounterChecks.Value.ZeroOnError(),
				Statement:   "Value returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject().Value(ctx)
					}, int(0))
				},
				Defect: counterValueValueBeside,
			},
			{
				ID:          CounterChecks.Value.Monotonic(),
				Statement:   "Value never returns less than it returned before",
				Claim:       spec.Monotonic,
				Probe:       spec.NeverDecreases,
				Reason:      spec.Regressed,
				Law:         spec.ResultsNeverDecrease,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CounterHandle) {
					probe.NeverDecreases(h, func(ctx context.Context) (int, error) {
						return h.Subject().Value(ctx)
					})
				},
				Defect: counterValueDescends,
			},
			{
				ID:          CounterChecks.Value.NeverGoesBack(),
				Statement:   "Value never returns less than a result that it returned before, over any sequence of calls",
				Claim:       spec.NeverGoesBack,
				Probe:       spec.ResultsNeverFall,
				Reason:      spec.Regressed,
				Law:         spec.ResultsNeverDecrease,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Counter]{CounterChecks.Value.Monotonic()},
				Property: func(h *CounterHandle, c *prop.Case) {
					probe.ResultsNeverFall[Counter, int](h, c, counterActions, "Value")
				},
				Defect: counterValueGoesBackAfterThree,
			},
			{
				ID:          CounterChecks.Grow.Smoke(),
				Statement:   "Grow returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CounterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Grow(ctx)
					})
				},
				Defect: counterGrowPanics,
			},
			{
				ID:          CounterChecks.Grow.Cancel(),
				Statement:   "Grow returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Grow(ctx)
					})
				},
				Defect: counterGrowIgnoresCancel,
			},
			{
				ID:          CounterChecks.Grow.KeepsWithinItsBounds(),
				Statement:   "Grow keeps the use that Size reads within its bounds, over any sequence of calls",
				Claim:       spec.KeepsWithinItsBounds,
				Probe:       spec.UseIsWithinTheBounds,
				Reason:      spec.ExceededBound,
				Law:         spec.UseIsWithinItsBounds,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *CounterHandle, c *prop.Case) {
					probe.UseIsWithinTheBounds(
						h,
						c,
						counterActions,
						func(ctx context.Context) (int, error) {
							values, err := counterSizeCall(ctx, h.Subject(), []any{})
							v, _ := values[0].(uint64)
							return int(v), err
						},
						0,
						64,
					)
				},
				Defect: counterSizeGrowsPastTheLimitAfterThree,
			},
			{
				ID:          CounterChecks.Size.Smoke(),
				Statement:   "Size returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CounterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Size(ctx)
						return err
					})
				},
				Defect: counterSizePanics,
			},
			{
				ID:          CounterChecks.Size.Cancel(),
				Statement:   "Size returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Size(ctx)
						return err
					})
				},
				Defect: counterSizeIgnoresCancel,
			},
			{
				ID:          CounterChecks.Size.ZeroOnError(),
				Statement:   "Size returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (uint64, error) {
						return h.Subject().Size(ctx)
					}, uint64(0))
				},
				Defect: counterSizeValueBeside,
			},
			{
				ID:          CounterChecks.Throttle.Smoke(),
				Statement:   "Throttle returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CounterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Throttle(ctx)
					})
				},
				Defect: counterThrottlePanics,
			},
			{
				ID:          CounterChecks.Throttle.Cancel(),
				Statement:   "Throttle returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Throttle(ctx)
					})
				},
				Defect: counterThrottleIgnoresCancel,
			},
			{
				ID:          CounterChecks.Load.Smoke(),
				Statement:   "Load returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CounterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Load(ctx)
						return err
					})
				},
				Defect: counterLoadPanics,
			},
			{
				ID:          CounterChecks.Load.Cancel(),
				Statement:   "Load returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Load(ctx)
						return err
					})
				},
				Defect: counterLoadIgnoresCancel,
			},
			{
				ID:          CounterChecks.Load.ZeroOnError(),
				Statement:   "Load returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CounterHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (float64, error) {
						return h.Subject().Load(ctx)
					}, float64(0))
				},
				Defect: counterLoadValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "Counter.Inc",
				Gate:     spec.InputSamples,
				Why:      "eidos derives no sample of the parameter opts of Inc",
				Remedy:   "author a sample of the parameter opts with the directive sample",
			},
			{
				Claim:    spec.KeepsWithinItsBounds,
				Callable: "Counter.Throttle",
				Gate:     spec.CountableUse,
				Why:      "the reader that observe names for Throttle returns no number and no length",
				Remedy:   "name with observe a callable that returns the use of Throttle as a number or a list",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Counter.Inc",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Counter.Grow",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Counter.Throttle",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
		},
		Calls: []conform.Call[Counter]{
			{
				Calls:    counterIncCalls,
				Over:     newCounterOver,
				Callable: CounterChecks.Inc.Callable,
			},
			{
				Calls:    counterValueCalls,
				Over:     newCounterOver,
				Callable: CounterChecks.Value.Callable,
			},
			{
				Calls:    counterGrowCalls,
				Over:     newCounterOver,
				Callable: CounterChecks.Grow.Callable,
			},
			{
				Calls:    counterSizeCalls,
				Over:     newCounterOver,
				Callable: CounterChecks.Size.Callable,
			},
			{
				Calls:    counterThrottleCalls,
				Over:     newCounterOver,
				Callable: CounterChecks.Throttle.Callable,
			},
			{
				Calls:    counterLoadCalls,
				Over:     newCounterOver,
				Callable: CounterChecks.Load.Callable,
			},
		},
	}
}

// RunCounterSuite applies the checks of Counter, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunCounterSuite[T Counter](t *testing.T, newSubject func() T, opts ...CounterOption) {
	t.Helper()
	RunCounterSuiteFromHandle(t, func(*CounterHandle) T { return newSubject() }, opts...)
}

// RunCounterSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunCounterSuiteFromHandle[T Counter](
	t *testing.T,
	newSubject func(h *CounterHandle) T,
	opts ...CounterOption,
) {
	t.Helper()
	conform.Run(t, CounterSuite(), counterSubjects(newSubject), opts...)
}

// FuzzCounterSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzCounterSuite[T Counter](
	f *testing.F,
	callable conform.Callable[Counter],
	newSubject func(h *CounterHandle) T,
	opts ...CounterOption,
) {
	f.Helper()
	conform.Fuzz(f, CounterSuite(), callable, counterSubjects(newSubject), opts...)
}

// RunCounterBudgets applies the allocation ceilings of Counter and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunCounterBudgets[T Counter](
	t *testing.T,
	newSubject func(h *CounterHandle) T,
	opts ...CounterOption,
) {
	t.Helper()
	conform.RunBudgets(t, CounterSuite(), counterSubjects(newSubject), opts...)
}

// BenchCounterSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchCounterSuite[T Counter](
	b conform.B,
	callable conform.Callable[Counter],
	newSubject func(h *CounterHandle) T,
	opts ...CounterOption,
) {
	b.Helper()
	conform.Bench(b, CounterSuite(), callable, counterSubjects(newSubject), opts...)
}

// counterSubjects returns the factory of the subjects of Counter that
// calls newSubject with the handle of each check.
func counterSubjects[T Counter](
	newSubject func(h *CounterHandle) T,
) func(h *CounterHandle) Counter {
	return func(h *CounterHandle) Counter { return newSubject(h) }
}

// counterIncCalls generates the calls of Counter.Inc with the inputs that
// each case draws.
var counterIncCalls = prop.Composite(func(*prop.Case) func(context.Context, Counter) error {
	return func(ctx context.Context, subject Counter) error {
		return subject.Inc(ctx)
	}
})

// counterValueCalls generates the calls of Counter.Value with the inputs
// that each case draws.
var counterValueCalls = prop.Composite(func(*prop.Case) func(context.Context, Counter) error {
	return func(ctx context.Context, subject Counter) error {
		_, err := subject.Value(ctx)
		return err
	}
})

// counterGrowCalls generates the calls of Counter.Grow with the inputs
// that each case draws.
var counterGrowCalls = prop.Composite(func(*prop.Case) func(context.Context, Counter) error {
	return func(ctx context.Context, subject Counter) error {
		return subject.Grow(ctx)
	}
})

// counterSizeCalls generates the calls of Counter.Size with the inputs
// that each case draws.
var counterSizeCalls = prop.Composite(func(*prop.Case) func(context.Context, Counter) error {
	return func(ctx context.Context, subject Counter) error {
		_, err := subject.Size(ctx)
		return err
	}
})

// counterThrottleCalls generates the calls of Counter.Throttle with the
// inputs that each case draws.
var counterThrottleCalls = prop.Composite(func(*prop.Case) func(context.Context, Counter) error {
	return func(ctx context.Context, subject Counter) error {
		return subject.Throttle(ctx)
	}
})

// counterLoadCalls generates the calls of Counter.Load with the inputs
// that each case draws.
var counterLoadCalls = prop.Composite(func(*prop.Case) func(context.Context, Counter) error {
	return func(ctx context.Context, subject Counter) error {
		_, err := subject.Load(ctx)
		return err
	}
})

// counterOverEmbedded is Counter under the name of the field that
// counterOver embeds.
type counterOverEmbedded = Counter

// counterOver is a subject of Counter whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type counterOver struct {
	counterOverEmbedded
	spend func()
}

// Inc calls spend and returns its zero results.
func (o counterOver) Inc(context.Context, ...Option) (r0 error) {
	o.spend()
	return r0
}

// Value calls spend and returns its zero results.
func (o counterOver) Value(context.Context) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// Grow calls spend and returns its zero results.
func (o counterOver) Grow(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Size calls spend and returns its zero results.
func (o counterOver) Size(context.Context) (r0 uint64, r1 error) {
	o.spend()
	return r0, r1
}

// Throttle calls spend and returns its zero results.
func (o counterOver) Throttle(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Load calls spend and returns its zero results.
func (o counterOver) Load(context.Context) (r0 float64, r1 error) {
	o.spend()
	return r0, r1
}

// newCounterOver returns the planted defect of a ceiling of Counter, whose
// callables call spend.
func newCounterOver(_ *CounterHandle, spend func()) Counter {
	return counterOver{spend: spend}
}

// counterValueGrowsPastTheLimitAfterThreeDefect is the planted defect of a
// Counter whose Value reports a use above the limit from its fourth call
// on.
type counterValueGrowsPastTheLimitAfterThreeDefect struct {
	counterOverEmbedded
	calls int
}

// Value is the override of the planted defect.
func (d *counterValueGrowsPastTheLimitAfterThreeDefect) Value(
	p0 context.Context,
) (r0 int, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.counterOverEmbedded.Value(p0)
	d.calls++
	if d.calls > 3 {
		r0 = 3 + 1
	}
	return r0, r1
}

// counterValueGrowsPastTheLimitAfterThree returns the planted defect of a
// Counter whose Value reports a use above the limit from its fourth call
// on over a fresh subject.
func counterValueGrowsPastTheLimitAfterThree(h *CounterHandle) Counter {
	return &counterValueGrowsPastTheLimitAfterThreeDefect{counterOverEmbedded: h.Fresh()}
}

// counterValuePanicsDefect is the planted defect of a Counter whose Value
// panics.
type counterValuePanicsDefect struct {
	counterOverEmbedded
}

// Value is the override of the planted defect.
func (*counterValuePanicsDefect) Value(p0 context.Context) (r0 int, r1 error) {
	panic("the planted defect of a Counter whose Value panics")
}

// counterValuePanics returns the planted defect of a Counter whose Value
// panics over a fresh subject.
func counterValuePanics(h *CounterHandle) Counter {
	return &counterValuePanicsDefect{counterOverEmbedded: h.Fresh()}
}

// counterValueIgnoresCancelDefect is the planted defect of a Counter whose
// Value ignores the context.
type counterValueIgnoresCancelDefect struct {
	counterOverEmbedded
}

// Value is the override of the planted defect.
func (d *counterValueIgnoresCancelDefect) Value(p0 context.Context) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.counterOverEmbedded.Value(p0)
}

// counterValueIgnoresCancel returns the planted defect of a Counter whose
// Value ignores the context over a fresh subject.
func counterValueIgnoresCancel(h *CounterHandle) Counter {
	return &counterValueIgnoresCancelDefect{counterOverEmbedded: h.Fresh()}
}

// counterValueValueBesideDefect is the planted defect of a Counter whose
// Value returns a value beside an error.
type counterValueValueBesideDefect struct {
	counterOverEmbedded
}

// Value is the override of the planted defect.
func (d *counterValueValueBesideDefect) Value(p0 context.Context) (r0 int, r1 error) {
	r0, r1 = d.counterOverEmbedded.Value(p0)
	if r1 != nil {
		r0 = int(7)
	}
	return r0, r1
}

// counterValueValueBeside returns the planted defect of a Counter whose
// Value returns a value beside an error over a fresh subject.
func counterValueValueBeside(h *CounterHandle) Counter {
	return &counterValueValueBesideDefect{counterOverEmbedded: h.Fresh()}
}

// counterValueDescendsDefect is the planted defect of a Counter whose
// Value returns less than its previous result.
type counterValueDescendsDefect struct {
	counterOverEmbedded
	calls int
}

// Value is the override of the planted defect.
func (d *counterValueDescendsDefect) Value(p0 context.Context) (r0 int, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	if d.calls == 1 {
		r0 = int(7)
	}
	return r0, r1
}

// counterValueDescends returns the planted defect of a Counter whose Value
// returns less than its previous result over a fresh subject.
func counterValueDescends(h *CounterHandle) Counter {
	return &counterValueDescendsDefect{counterOverEmbedded: h.Fresh()}
}

// counterValueGoesBackAfterThreeDefect is the planted defect of a Counter
// whose Value returns its first result again from its fourth call on.
type counterValueGoesBackAfterThreeDefect struct {
	counterOverEmbedded
	calls int
	first *int
}

// Value is the override of the planted defect.
func (d *counterValueGoesBackAfterThreeDefect) Value(p0 context.Context) (r0 int, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.counterOverEmbedded.Value(p0)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// counterValueGoesBackAfterThree returns the planted defect of a Counter
// whose Value returns its first result again from its fourth call on over
// a fresh subject.
func counterValueGoesBackAfterThree(h *CounterHandle) Counter {
	return &counterValueGoesBackAfterThreeDefect{counterOverEmbedded: h.Fresh()}
}

// counterGrowPanicsDefect is the planted defect of a Counter whose Grow
// panics.
type counterGrowPanicsDefect struct {
	counterOverEmbedded
}

// Grow is the override of the planted defect.
func (*counterGrowPanicsDefect) Grow(p0 context.Context) (r0 error) {
	panic("the planted defect of a Counter whose Grow panics")
}

// counterGrowPanics returns the planted defect of a Counter whose Grow
// panics over a fresh subject.
func counterGrowPanics(h *CounterHandle) Counter {
	return &counterGrowPanicsDefect{counterOverEmbedded: h.Fresh()}
}

// counterGrowIgnoresCancelDefect is the planted defect of a Counter whose
// Grow ignores the context.
type counterGrowIgnoresCancelDefect struct {
	counterOverEmbedded
}

// Grow is the override of the planted defect.
func (d *counterGrowIgnoresCancelDefect) Grow(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.counterOverEmbedded.Grow(p0)
}

// counterGrowIgnoresCancel returns the planted defect of a Counter whose
// Grow ignores the context over a fresh subject.
func counterGrowIgnoresCancel(h *CounterHandle) Counter {
	return &counterGrowIgnoresCancelDefect{counterOverEmbedded: h.Fresh()}
}

// counterSizeGrowsPastTheLimitAfterThreeDefect is the planted defect of a
// Counter whose Size reports a use above the limit from its fourth call
// on.
type counterSizeGrowsPastTheLimitAfterThreeDefect struct {
	counterOverEmbedded
	calls int
}

// Size is the override of the planted defect.
func (d *counterSizeGrowsPastTheLimitAfterThreeDefect) Size(
	p0 context.Context,
) (r0 uint64, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.counterOverEmbedded.Size(p0)
	d.calls++
	if d.calls > 3 {
		r0 = 64 + 1
	}
	return r0, r1
}

// counterSizeGrowsPastTheLimitAfterThree returns the planted defect of a
// Counter whose Size reports a use above the limit from its fourth call on
// over a fresh subject.
func counterSizeGrowsPastTheLimitAfterThree(h *CounterHandle) Counter {
	return &counterSizeGrowsPastTheLimitAfterThreeDefect{counterOverEmbedded: h.Fresh()}
}

// counterSizePanicsDefect is the planted defect of a Counter whose Size
// panics.
type counterSizePanicsDefect struct {
	counterOverEmbedded
}

// Size is the override of the planted defect.
func (*counterSizePanicsDefect) Size(p0 context.Context) (r0 uint64, r1 error) {
	panic("the planted defect of a Counter whose Size panics")
}

// counterSizePanics returns the planted defect of a Counter whose Size
// panics over a fresh subject.
func counterSizePanics(h *CounterHandle) Counter {
	return &counterSizePanicsDefect{counterOverEmbedded: h.Fresh()}
}

// counterSizeIgnoresCancelDefect is the planted defect of a Counter whose
// Size ignores the context.
type counterSizeIgnoresCancelDefect struct {
	counterOverEmbedded
}

// Size is the override of the planted defect.
func (d *counterSizeIgnoresCancelDefect) Size(p0 context.Context) (r0 uint64, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.counterOverEmbedded.Size(p0)
}

// counterSizeIgnoresCancel returns the planted defect of a Counter whose
// Size ignores the context over a fresh subject.
func counterSizeIgnoresCancel(h *CounterHandle) Counter {
	return &counterSizeIgnoresCancelDefect{counterOverEmbedded: h.Fresh()}
}

// counterSizeValueBesideDefect is the planted defect of a Counter whose
// Size returns a value beside an error.
type counterSizeValueBesideDefect struct {
	counterOverEmbedded
}

// Size is the override of the planted defect.
func (d *counterSizeValueBesideDefect) Size(p0 context.Context) (r0 uint64, r1 error) {
	r0, r1 = d.counterOverEmbedded.Size(p0)
	if r1 != nil {
		r0 = uint64(7)
	}
	return r0, r1
}

// counterSizeValueBeside returns the planted defect of a Counter whose
// Size returns a value beside an error over a fresh subject.
func counterSizeValueBeside(h *CounterHandle) Counter {
	return &counterSizeValueBesideDefect{counterOverEmbedded: h.Fresh()}
}

// counterThrottlePanicsDefect is the planted defect of a Counter whose
// Throttle panics.
type counterThrottlePanicsDefect struct {
	counterOverEmbedded
}

// Throttle is the override of the planted defect.
func (*counterThrottlePanicsDefect) Throttle(p0 context.Context) (r0 error) {
	panic("the planted defect of a Counter whose Throttle panics")
}

// counterThrottlePanics returns the planted defect of a Counter whose
// Throttle panics over a fresh subject.
func counterThrottlePanics(h *CounterHandle) Counter {
	return &counterThrottlePanicsDefect{counterOverEmbedded: h.Fresh()}
}

// counterThrottleIgnoresCancelDefect is the planted defect of a Counter
// whose Throttle ignores the context.
type counterThrottleIgnoresCancelDefect struct {
	counterOverEmbedded
}

// Throttle is the override of the planted defect.
func (d *counterThrottleIgnoresCancelDefect) Throttle(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.counterOverEmbedded.Throttle(p0)
}

// counterThrottleIgnoresCancel returns the planted defect of a Counter
// whose Throttle ignores the context over a fresh subject.
func counterThrottleIgnoresCancel(h *CounterHandle) Counter {
	return &counterThrottleIgnoresCancelDefect{counterOverEmbedded: h.Fresh()}
}

// counterLoadPanicsDefect is the planted defect of a Counter whose Load
// panics.
type counterLoadPanicsDefect struct {
	counterOverEmbedded
}

// Load is the override of the planted defect.
func (*counterLoadPanicsDefect) Load(p0 context.Context) (r0 float64, r1 error) {
	panic("the planted defect of a Counter whose Load panics")
}

// counterLoadPanics returns the planted defect of a Counter whose Load
// panics over a fresh subject.
func counterLoadPanics(h *CounterHandle) Counter {
	return &counterLoadPanicsDefect{counterOverEmbedded: h.Fresh()}
}

// counterLoadIgnoresCancelDefect is the planted defect of a Counter whose
// Load ignores the context.
type counterLoadIgnoresCancelDefect struct {
	counterOverEmbedded
}

// Load is the override of the planted defect.
func (d *counterLoadIgnoresCancelDefect) Load(p0 context.Context) (r0 float64, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.counterOverEmbedded.Load(p0)
}

// counterLoadIgnoresCancel returns the planted defect of a Counter whose
// Load ignores the context over a fresh subject.
func counterLoadIgnoresCancel(h *CounterHandle) Counter {
	return &counterLoadIgnoresCancelDefect{counterOverEmbedded: h.Fresh()}
}

// counterLoadValueBesideDefect is the planted defect of a Counter whose
// Load returns a value beside an error.
type counterLoadValueBesideDefect struct {
	counterOverEmbedded
}

// Load is the override of the planted defect.
func (d *counterLoadValueBesideDefect) Load(p0 context.Context) (r0 float64, r1 error) {
	r0, r1 = d.counterOverEmbedded.Load(p0)
	if r1 != nil {
		r0 = float64(2.5)
	}
	return r0, r1
}

// counterLoadValueBeside returns the planted defect of a Counter whose
// Load returns a value beside an error over a fresh subject.
func counterLoadValueBeside(h *CounterHandle) Counter {
	return &counterLoadValueBesideDefect{counterOverEmbedded: h.Fresh()}
}

// counterActions returns the actions of the generated sequences of Counter
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func counterActions(c *prop.Case, subject func() Counter) []stateful.Action[struct{}] {
	return []stateful.Action[struct{}]{
		{
			Name: "Inc",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Inc", args)
				values, err := counterIncCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Value",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Value", args)
				values, err := counterValueCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Grow",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Grow", args)
				values, err := counterGrowCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Size",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Size", args)
				values, err := counterSizeCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Throttle",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Throttle", args)
				values, err := counterThrottleCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Load",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Load", args)
				values, err := counterLoadCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// counterIncCall calls Inc on subject with the arguments args, and returns
// the values of the call beside its error.
func counterIncCall(ctx context.Context, subject Counter, args []any) ([]any, error) {
	return nil, subject.Inc(ctx)
}

// counterValueCall calls Value on subject with the arguments args, and
// returns the values of the call beside its error.
func counterValueCall(ctx context.Context, subject Counter, args []any) ([]any, error) {
	r0, err := subject.Value(ctx)
	return []any{r0}, err
}

// counterGrowCall calls Grow on subject with the arguments args, and
// returns the values of the call beside its error.
func counterGrowCall(ctx context.Context, subject Counter, args []any) ([]any, error) {
	return nil, subject.Grow(ctx)
}

// counterSizeCall calls Size on subject with the arguments args, and
// returns the values of the call beside its error.
func counterSizeCall(ctx context.Context, subject Counter, args []any) ([]any, error) {
	r0, err := subject.Size(ctx)
	return []any{r0}, err
}

// counterThrottleCall calls Throttle on subject with the arguments args,
// and returns the values of the call beside its error.
func counterThrottleCall(ctx context.Context, subject Counter, args []any) ([]any, error) {
	return nil, subject.Throttle(ctx)
}

// counterLoadCall calls Load on subject with the arguments args, and
// returns the values of the call beside its error.
func counterLoadCall(ctx context.Context, subject Counter, args []any) ([]any, error) {
	r0, err := subject.Load(ctx)
	return []any{r0}, err
}

// The derived checks of Notes:
//
// Notes.Push/smoke: Push returns without a panic for derived inputs.
//
// Notes.Push/smoke#generated: Push returns without a panic for derived
// inputs.
//
// Notes.Push/cancel: Push returns an error for a cancelled context.
//
// Notes.Push/cancel#generated: Push returns an error for a cancelled
// context.
//
// Notes.Push/keeps-within-its-bounds: Push keeps the use that Items reads
// within its bounds, over any sequence of calls.
//
// Notes.Pop/smoke: Pop returns without a panic for derived inputs.
//
// Notes.Pop/cancel: Pop returns an error for a cancelled context.
//
// Notes.Items/smoke: Items returns without a panic for derived inputs.
//
// Notes.Items/smoke#generated: Items returns without a panic for derived
// inputs.
//
// Notes.Items/cancel: Items returns an error for a cancelled context.
//
// Notes.Items/cancel#generated: Items returns an error for a cancelled
// context.
//
// Notes.Items/zero-on-error: Items returns the zero value beside an error.
//
// Notes.Items/zero-on-error#generated: Items returns the zero value beside
// an error.
//
// Notes.Append/smoke: Append returns without a panic for derived inputs.
//
// Notes.Append/smoke#generated: Append returns without a panic for derived
// inputs.
//
// Notes.Append/cancel: Append returns an error for a cancelled context.
//
// Notes.Append/cancel#generated: Append returns an error for a cancelled
// context.
//
// Notes.Append/keeps-within-its-bounds: Append keeps the use that Text
// reads within its bounds, over any sequence of calls.
//
// Notes.Text/smoke: Text returns without a panic for derived inputs.
//
// Notes.Text/cancel: Text returns an error for a cancelled context.
//
// Notes.Text/zero-on-error: Text returns the zero value beside an error.
//
// Notes.Pin/smoke: Pin returns without a panic for derived inputs.
//
// Notes.Pin/smoke#generated: Pin returns without a panic for derived
// inputs.
//
// Notes.Pin/cancel: Pin returns an error for a cancelled context.
//
// Notes.Pin/cancel#generated: Pin returns an error for a cancelled
// context.
//
// Notes.Head/smoke: Head returns without a panic for derived inputs.
//
// Notes.Head/cancel: Head returns an error for a cancelled context.
//
// Notes.Head/zero-on-error: Head returns the zero value beside an error.

// The derived checks of Notes that witness their laws at the shallow end:
//
// Notes.Push/smoke
//
// Notes.Pop/smoke
//
// Notes.Items/smoke
//
// Notes.Append/smoke
//
// Notes.Text/smoke
//
// Notes.Pin/smoke
//
// Notes.Head/smoke

// The refused claims of Notes:
//
// Notes.Pop/keeps-within-its-bounds: the reader that observe names for Pop
// returns no number and no length. The remedy is to name with observe a
// callable that returns the use of Pop as a number or a list.
//
// Notes.Pin/keeps-within-its-bounds: the reader that observe names for Pin
// returns no number and no length. The remedy is to name with observe a
// callable that returns the use of Pin as a number or a list.

// The unclaimed classifications of Notes:
//
// The shape:writer of Notes.Push has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:lifecycle of Notes.Pop has the stance silent. A lifecycle call
// returns nothing that a check could compare.
//
// The shape:writer of Notes.Append has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Notes.Pin has the stance silent. A write is an input
// of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Notes writes.
type (
	NotesHandle = conform.Handle[Notes]
	NotesCheck  = conform.Check[Notes]
	NotesOption = conform.Option[Notes]
)

// NotesChecks contains the name of each callable of Notes, with the
// identities of its checks.
var NotesChecks = struct {
	Push   notesPush
	Pop    notesPop
	Items  notesItems
	Append notesAppend
	Text   notesText
	Pin    notesPin
	Head   notesHead
}{
	Push:   notesPush{"Notes.Push"},
	Pop:    notesPop{"Notes.Pop"},
	Items:  notesItems{"Notes.Items"},
	Append: notesAppend{"Notes.Append"},
	Text:   notesText{"Notes.Text"},
	Pin:    notesPin{"Notes.Pin"},
	Head:   notesHead{"Notes.Head"},
}

// notesPush is the name of Notes.Push, with the identities of its checks.
type notesPush struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Push returns without a
// panic for derived inputs.
func (c notesPush) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Push returns
// without a panic for derived inputs.
func (c notesPush) SmokeGenerated() conform.ID[Notes] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Push returns an error for
// a cancelled context.
func (c notesPush) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Push returns an
// error for a cancelled context.
func (c notesPush) CancelGenerated() conform.ID[Notes] {
	return c.Check("cancel#generated")
}

// KeepsWithinItsBounds returns the identity of the check that Push keeps
// the use that Items reads within its bounds, over any sequence of calls.
func (c notesPush) KeepsWithinItsBounds() conform.ID[Notes] {
	return c.Check("keeps-within-its-bounds")
}

// notesPop is the name of Notes.Pop, with the identities of its checks.
type notesPop struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Pop returns without a panic
// for derived inputs.
func (c notesPop) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Pop returns an error for a
// cancelled context.
func (c notesPop) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// notesItems is the name of Notes.Items, with the identities of its
// checks.
type notesItems struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Items returns without a
// panic for derived inputs.
func (c notesItems) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Items returns
// without a panic for derived inputs.
func (c notesItems) SmokeGenerated() conform.ID[Notes] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Items returns an error for
// a cancelled context.
func (c notesItems) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Items returns an
// error for a cancelled context.
func (c notesItems) CancelGenerated() conform.ID[Notes] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Items returns the
// zero value beside an error.
func (c notesItems) ZeroOnError() conform.ID[Notes] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Items
// returns the zero value beside an error.
func (c notesItems) ZeroOnErrorGenerated() conform.ID[Notes] {
	return c.Check("zero-on-error#generated")
}

// notesAppend is the name of Notes.Append, with the identities of its
// checks.
type notesAppend struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Append returns without a
// panic for derived inputs.
func (c notesAppend) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Append returns
// without a panic for derived inputs.
func (c notesAppend) SmokeGenerated() conform.ID[Notes] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Append returns an error
// for a cancelled context.
func (c notesAppend) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Append returns an
// error for a cancelled context.
func (c notesAppend) CancelGenerated() conform.ID[Notes] {
	return c.Check("cancel#generated")
}

// KeepsWithinItsBounds returns the identity of the check that Append keeps
// the use that Text reads within its bounds, over any sequence of calls.
func (c notesAppend) KeepsWithinItsBounds() conform.ID[Notes] {
	return c.Check("keeps-within-its-bounds")
}

// notesText is the name of Notes.Text, with the identities of its checks.
type notesText struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Text returns without a
// panic for derived inputs.
func (c notesText) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Text returns an error for
// a cancelled context.
func (c notesText) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Text returns the zero
// value beside an error.
func (c notesText) ZeroOnError() conform.ID[Notes] {
	return c.Check("zero-on-error")
}

// notesPin is the name of Notes.Pin, with the identities of its checks.
type notesPin struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Pin returns without a panic
// for derived inputs.
func (c notesPin) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Pin returns
// without a panic for derived inputs.
func (c notesPin) SmokeGenerated() conform.ID[Notes] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Pin returns an error for a
// cancelled context.
func (c notesPin) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Pin returns an
// error for a cancelled context.
func (c notesPin) CancelGenerated() conform.ID[Notes] {
	return c.Check("cancel#generated")
}

// notesHead is the name of Notes.Head, with the identities of its checks.
type notesHead struct {
	conform.Callable[Notes]
}

// Smoke returns the identity of the check that Head returns without a
// panic for derived inputs.
func (c notesHead) Smoke() conform.ID[Notes] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Head returns an error for
// a cancelled context.
func (c notesHead) Cancel() conform.ID[Notes] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Head returns the zero
// value beside an error.
func (c notesHead) ZeroOnError() conform.ID[Notes] {
	return c.Check("zero-on-error")
}

// NotesSuite returns the checks of Notes, and the calls and the budgets of
// its callables.
func NotesSuite() conform.Suite[Notes] {
	return conform.Suite[Notes]{
		Declaration: "Notes",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Notes]{
			{
				ID:          NotesChecks.Push.Smoke(),
				Statement:   "Push returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Push(ctx, string("item-a"))
					})
				},
				Defect: notesPushPanics,
			},
			{
				ID:          NotesChecks.Push.SmokeGenerated(),
				Statement:   "Push returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Push(ctx, drawnItem)
					})
				},
				Defect: notesPushPanics,
			},
			{
				ID:          NotesChecks.Push.Cancel(),
				Statement:   "Push returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Push(ctx, string("item-a"))
					})
				},
				Defect: notesPushIgnoresCancel,
			},
			{
				ID:          NotesChecks.Push.CancelGenerated(),
				Statement:   "Push returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Push(ctx, drawnItem)
					})
				},
				Defect: notesPushIgnoresCancel,
			},
			{
				ID:          NotesChecks.Push.KeepsWithinItsBounds(),
				Statement:   "Push keeps the use that Items reads within its bounds, over any sequence of calls",
				Claim:       spec.KeepsWithinItsBounds,
				Probe:       spec.UseIsWithinTheBounds,
				Reason:      spec.ExceededBound,
				Law:         spec.UseIsWithinItsBounds,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *NotesHandle, c *prop.Case) {
					probe.UseIsWithinTheBounds(
						h,
						c,
						notesActions,
						func(ctx context.Context) (int, error) {
							values, err := notesItemsCall(
								ctx,
								h.Subject(),
								[]any{string("prefix-a")},
							)
							v, _ := values[0].([]string)
							return len(v), err
						},
						0,
						4,
					)
				},
				Defect: notesItemsGrowsPastTheLimitAfterThree,
			},
			{
				ID:          NotesChecks.Pop.Smoke(),
				Statement:   "Pop returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Pop(ctx)
					})
				},
				Defect: notesPopPanics,
			},
			{
				ID:          NotesChecks.Pop.Cancel(),
				Statement:   "Pop returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Pop(ctx)
					})
				},
				Defect: notesPopIgnoresCancel,
			},
			{
				ID:          NotesChecks.Items.Smoke(),
				Statement:   "Items returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Items(ctx, string("prefix-a"))
						return err
					})
				},
				Defect: notesItemsPanics,
			},
			{
				ID:          NotesChecks.Items.SmokeGenerated(),
				Statement:   "Items returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnPrefix := c.Draw(prop.Of[string](), "prefix")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Items(ctx, drawnPrefix)
						return err
					})
				},
				Defect: notesItemsPanics,
			},
			{
				ID:          NotesChecks.Items.Cancel(),
				Statement:   "Items returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Items(ctx, string("prefix-a"))
						return err
					})
				},
				Defect: notesItemsIgnoresCancel,
			},
			{
				ID:          NotesChecks.Items.CancelGenerated(),
				Statement:   "Items returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnPrefix := c.Draw(prop.Of[string](), "prefix")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Items(ctx, drawnPrefix)
						return err
					})
				},
				Defect: notesItemsIgnoresCancel,
			},
			{
				ID:          NotesChecks.Items.ZeroOnError(),
				Statement:   "Items returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) ([]string, error) {
						return h.Subject().Items(ctx, string("prefix-a"))
					}, []string(nil))
				},
				Defect: notesItemsValueBeside,
			},
			{
				ID:          NotesChecks.Items.ZeroOnErrorGenerated(),
				Statement:   "Items returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnPrefix := c.Draw(prop.Of[string](), "prefix")
					probe.ZeroBesideError(h, func(ctx context.Context) ([]string, error) {
						return h.Subject().Items(ctx, drawnPrefix)
					}, []string(nil))
				},
				Defect: notesItemsValueBeside,
			},
			{
				ID:          NotesChecks.Append.Smoke(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, string("text-a"))
					})
				},
				Defect: notesAppendPanics,
			},
			{
				ID:          NotesChecks.Append.SmokeGenerated(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnText := c.Draw(prop.Of[string](), "text")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, drawnText)
					})
				},
				Defect: notesAppendPanics,
			},
			{
				ID:          NotesChecks.Append.Cancel(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, string("text-a"))
					})
				},
				Defect: notesAppendIgnoresCancel,
			},
			{
				ID:          NotesChecks.Append.CancelGenerated(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnText := c.Draw(prop.Of[string](), "text")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Append(ctx, drawnText)
					})
				},
				Defect: notesAppendIgnoresCancel,
			},
			{
				ID:          NotesChecks.Append.KeepsWithinItsBounds(),
				Statement:   "Append keeps the use that Text reads within its bounds, over any sequence of calls",
				Claim:       spec.KeepsWithinItsBounds,
				Probe:       spec.UseIsWithinTheBounds,
				Reason:      spec.ExceededBound,
				Law:         spec.UseIsWithinItsBounds,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *NotesHandle, c *prop.Case) {
					probe.UseIsWithinTheBounds(
						h,
						c,
						notesActions,
						func(ctx context.Context) (int, error) {
							values, err := notesTextCall(ctx, h.Subject(), []any{})
							v, _ := values[0].(string)
							return len(v), err
						},
						0,
						8,
					)
				},
				Defect: notesTextGrowsPastTheLimitAfterThree,
			},
			{
				ID:          NotesChecks.Text.Smoke(),
				Statement:   "Text returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Text(ctx)
						return err
					})
				},
				Defect: notesTextPanics,
			},
			{
				ID:          NotesChecks.Text.Cancel(),
				Statement:   "Text returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Text(ctx)
						return err
					})
				},
				Defect: notesTextIgnoresCancel,
			},
			{
				ID:          NotesChecks.Text.ZeroOnError(),
				Statement:   "Text returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Text(ctx)
					}, string(""))
				},
				Defect: notesTextValueBeside,
			},
			{
				ID:          NotesChecks.Pin.Smoke(),
				Statement:   "Pin returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Pin(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: notesPinPanics,
			},
			{
				ID:          NotesChecks.Pin.SmokeGenerated(),
				Statement:   "Pin returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Pin(ctx, drawnIt)
					})
				},
				Defect: notesPinPanics,
			},
			{
				ID:          NotesChecks.Pin.Cancel(),
				Statement:   "Pin returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Pin(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: notesPinIgnoresCancel,
			},
			{
				ID:          NotesChecks.Pin.CancelGenerated(),
				Statement:   "Pin returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *NotesHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Pin(ctx, drawnIt)
					})
				},
				Defect: notesPinIgnoresCancel,
			},
			{
				ID:          NotesChecks.Head.Smoke(),
				Statement:   "Head returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *NotesHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Head(ctx)
						return err
					})
				},
				Defect: notesHeadPanics,
			},
			{
				ID:          NotesChecks.Head.Cancel(),
				Statement:   "Head returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Head(ctx)
						return err
					})
				},
				Defect: notesHeadIgnoresCancel,
			},
			{
				ID:          NotesChecks.Head.ZeroOnError(),
				Statement:   "Head returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *NotesHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Head(ctx)
					}, Item{})
				},
				Defect: notesHeadValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.KeepsWithinItsBounds,
				Callable: "Notes.Pop",
				Gate:     spec.CountableUse,
				Why:      "the reader that observe names for Pop returns no number and no length",
				Remedy:   "name with observe a callable that returns the use of Pop as a number or a list",
			},
			{
				Claim:    spec.KeepsWithinItsBounds,
				Callable: "Notes.Pin",
				Gate:     spec.CountableUse,
				Why:      "the reader that observe names for Pin returns no number and no length",
				Remedy:   "name with observe a callable that returns the use of Pin as a number or a list",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Notes.Push",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Notes.Pop",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Notes.Append",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Notes.Pin",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Notes]{
			{
				Calls:    notesPushCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Push.Callable,
			},
			{
				Calls:    notesPopCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Pop.Callable,
			},
			{
				Calls:    notesItemsCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Items.Callable,
			},
			{
				Calls:    notesAppendCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Append.Callable,
			},
			{
				Calls:    notesTextCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Text.Callable,
			},
			{
				Calls:    notesPinCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Pin.Callable,
			},
			{
				Calls:    notesHeadCalls,
				Over:     newNotesOver,
				Callable: NotesChecks.Head.Callable,
			},
		},
	}
}

// RunNotesSuite applies the checks of Notes, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunNotesSuite[T Notes](t *testing.T, newSubject func() T, opts ...NotesOption) {
	t.Helper()
	RunNotesSuiteFromHandle(t, func(*NotesHandle) T { return newSubject() }, opts...)
}

// RunNotesSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunNotesSuiteFromHandle[T Notes](
	t *testing.T,
	newSubject func(h *NotesHandle) T,
	opts ...NotesOption,
) {
	t.Helper()
	conform.Run(t, NotesSuite(), notesSubjects(newSubject), opts...)
}

// FuzzNotesSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzNotesSuite[T Notes](
	f *testing.F,
	callable conform.Callable[Notes],
	newSubject func(h *NotesHandle) T,
	opts ...NotesOption,
) {
	f.Helper()
	conform.Fuzz(f, NotesSuite(), callable, notesSubjects(newSubject), opts...)
}

// RunNotesBudgets applies the allocation ceilings of Notes and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunNotesBudgets[T Notes](
	t *testing.T,
	newSubject func(h *NotesHandle) T,
	opts ...NotesOption,
) {
	t.Helper()
	conform.RunBudgets(t, NotesSuite(), notesSubjects(newSubject), opts...)
}

// BenchNotesSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchNotesSuite[T Notes](
	b conform.B,
	callable conform.Callable[Notes],
	newSubject func(h *NotesHandle) T,
	opts ...NotesOption,
) {
	b.Helper()
	conform.Bench(b, NotesSuite(), callable, notesSubjects(newSubject), opts...)
}

// notesSubjects returns the factory of the subjects of Notes that calls
// newSubject with the handle of each check.
func notesSubjects[T Notes](newSubject func(h *NotesHandle) T) func(h *NotesHandle) Notes {
	return func(h *NotesHandle) Notes { return newSubject(h) }
}

// notesPushCalls generates the calls of Notes.Push with the inputs that
// each case draws.
var notesPushCalls = prop.Composite(func(c *prop.Case) func(context.Context, Notes) error {
	item := c.Draw(prop.Of[string](), "item")
	return func(ctx context.Context, subject Notes) error {
		return subject.Push(ctx, item)
	}
})

// notesPopCalls generates the calls of Notes.Pop with the inputs that each
// case draws.
var notesPopCalls = prop.Composite(func(*prop.Case) func(context.Context, Notes) error {
	return func(ctx context.Context, subject Notes) error {
		return subject.Pop(ctx)
	}
})

// notesItemsCalls generates the calls of Notes.Items with the inputs that
// each case draws.
var notesItemsCalls = prop.Composite(func(c *prop.Case) func(context.Context, Notes) error {
	prefix := c.Draw(prop.Of[string](), "prefix")
	return func(ctx context.Context, subject Notes) error {
		_, err := subject.Items(ctx, prefix)
		return err
	}
})

// notesAppendCalls generates the calls of Notes.Append with the inputs
// that each case draws.
var notesAppendCalls = prop.Composite(func(c *prop.Case) func(context.Context, Notes) error {
	text := c.Draw(prop.Of[string](), "text")
	return func(ctx context.Context, subject Notes) error {
		return subject.Append(ctx, text)
	}
})

// notesTextCalls generates the calls of Notes.Text with the inputs that
// each case draws.
var notesTextCalls = prop.Composite(func(*prop.Case) func(context.Context, Notes) error {
	return func(ctx context.Context, subject Notes) error {
		_, err := subject.Text(ctx)
		return err
	}
})

// notesPinCalls generates the calls of Notes.Pin with the inputs that each
// case draws.
var notesPinCalls = prop.Composite(func(c *prop.Case) func(context.Context, Notes) error {
	it := c.Draw(prop.Of[Item](), "it")
	return func(ctx context.Context, subject Notes) error {
		return subject.Pin(ctx, it)
	}
})

// notesHeadCalls generates the calls of Notes.Head with the inputs that
// each case draws.
var notesHeadCalls = prop.Composite(func(*prop.Case) func(context.Context, Notes) error {
	return func(ctx context.Context, subject Notes) error {
		_, err := subject.Head(ctx)
		return err
	}
})

// notesOverEmbedded is Notes under the name of the field that notesOver
// embeds.
type notesOverEmbedded = Notes

// notesOver is a subject of Notes whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type notesOver struct {
	notesOverEmbedded
	spend func()
}

// Push calls spend and returns its zero results.
func (o notesOver) Push(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Pop calls spend and returns its zero results.
func (o notesOver) Pop(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Items calls spend and returns its zero results.
func (o notesOver) Items(context.Context, string) (r0 []string, r1 error) {
	o.spend()
	return r0, r1
}

// Append calls spend and returns its zero results.
func (o notesOver) Append(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Text calls spend and returns its zero results.
func (o notesOver) Text(context.Context) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// Pin calls spend and returns its zero results.
func (o notesOver) Pin(context.Context, Item) (r0 error) {
	o.spend()
	return r0
}

// Head calls spend and returns its zero results.
func (o notesOver) Head(context.Context) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// newNotesOver returns the planted defect of a ceiling of Notes, whose
// callables call spend.
func newNotesOver(_ *NotesHandle, spend func()) Notes {
	return notesOver{spend: spend}
}

// notesPushPanicsDefect is the planted defect of a Notes whose Push
// panics.
type notesPushPanicsDefect struct {
	notesOverEmbedded
}

// Push is the override of the planted defect.
func (*notesPushPanicsDefect) Push(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Notes whose Push panics")
}

// notesPushPanics returns the planted defect of a Notes whose Push panics
// over a fresh subject.
func notesPushPanics(h *NotesHandle) Notes {
	return &notesPushPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesPushIgnoresCancelDefect is the planted defect of a Notes whose Push
// ignores the context.
type notesPushIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Push is the override of the planted defect.
func (d *notesPushIgnoresCancelDefect) Push(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.notesOverEmbedded.Push(p0, p1)
}

// notesPushIgnoresCancel returns the planted defect of a Notes whose Push
// ignores the context over a fresh subject.
func notesPushIgnoresCancel(h *NotesHandle) Notes {
	return &notesPushIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesItemsGrowsPastTheLimitAfterThreeDefect is the planted defect of a
// Notes whose Items reports a use above the limit from its fourth call on.
type notesItemsGrowsPastTheLimitAfterThreeDefect struct {
	notesOverEmbedded
	calls int
}

// Items is the override of the planted defect.
func (d *notesItemsGrowsPastTheLimitAfterThreeDefect) Items(
	p0 context.Context,
	p1 string,
) (r0 []string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.notesOverEmbedded.Items(p0, p1)
	d.calls++
	if d.calls > 3 {
		r0 = make([]string, 4+1)
	}
	return r0, r1
}

// notesItemsGrowsPastTheLimitAfterThree returns the planted defect of a
// Notes whose Items reports a use above the limit from its fourth call on
// over a fresh subject.
func notesItemsGrowsPastTheLimitAfterThree(h *NotesHandle) Notes {
	return &notesItemsGrowsPastTheLimitAfterThreeDefect{notesOverEmbedded: h.Fresh()}
}

// notesPopPanicsDefect is the planted defect of a Notes whose Pop panics.
type notesPopPanicsDefect struct {
	notesOverEmbedded
}

// Pop is the override of the planted defect.
func (*notesPopPanicsDefect) Pop(p0 context.Context) (r0 error) {
	panic("the planted defect of a Notes whose Pop panics")
}

// notesPopPanics returns the planted defect of a Notes whose Pop panics
// over a fresh subject.
func notesPopPanics(h *NotesHandle) Notes {
	return &notesPopPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesPopIgnoresCancelDefect is the planted defect of a Notes whose Pop
// ignores the context.
type notesPopIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Pop is the override of the planted defect.
func (d *notesPopIgnoresCancelDefect) Pop(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.notesOverEmbedded.Pop(p0)
}

// notesPopIgnoresCancel returns the planted defect of a Notes whose Pop
// ignores the context over a fresh subject.
func notesPopIgnoresCancel(h *NotesHandle) Notes {
	return &notesPopIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesItemsPanicsDefect is the planted defect of a Notes whose Items
// panics.
type notesItemsPanicsDefect struct {
	notesOverEmbedded
}

// Items is the override of the planted defect.
func (*notesItemsPanicsDefect) Items(p0 context.Context, p1 string) (r0 []string, r1 error) {
	panic("the planted defect of a Notes whose Items panics")
}

// notesItemsPanics returns the planted defect of a Notes whose Items
// panics over a fresh subject.
func notesItemsPanics(h *NotesHandle) Notes {
	return &notesItemsPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesItemsIgnoresCancelDefect is the planted defect of a Notes whose
// Items ignores the context.
type notesItemsIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Items is the override of the planted defect.
func (d *notesItemsIgnoresCancelDefect) Items(
	p0 context.Context,
	p1 string,
) (r0 []string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.notesOverEmbedded.Items(p0, p1)
}

// notesItemsIgnoresCancel returns the planted defect of a Notes whose
// Items ignores the context over a fresh subject.
func notesItemsIgnoresCancel(h *NotesHandle) Notes {
	return &notesItemsIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesItemsValueBesideDefect is the planted defect of a Notes whose Items
// returns a value beside an error.
type notesItemsValueBesideDefect struct {
	notesOverEmbedded
}

// Items is the override of the planted defect.
func (d *notesItemsValueBesideDefect) Items(p0 context.Context, p1 string) (r0 []string, r1 error) {
	r0, r1 = d.notesOverEmbedded.Items(p0, p1)
	if r1 != nil {
		r0 = []string{"sample-b"}
	}
	return r0, r1
}

// notesItemsValueBeside returns the planted defect of a Notes whose Items
// returns a value beside an error over a fresh subject.
func notesItemsValueBeside(h *NotesHandle) Notes {
	return &notesItemsValueBesideDefect{notesOverEmbedded: h.Fresh()}
}

// notesAppendPanicsDefect is the planted defect of a Notes whose Append
// panics.
type notesAppendPanicsDefect struct {
	notesOverEmbedded
}

// Append is the override of the planted defect.
func (*notesAppendPanicsDefect) Append(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Notes whose Append panics")
}

// notesAppendPanics returns the planted defect of a Notes whose Append
// panics over a fresh subject.
func notesAppendPanics(h *NotesHandle) Notes {
	return &notesAppendPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesAppendIgnoresCancelDefect is the planted defect of a Notes whose
// Append ignores the context.
type notesAppendIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Append is the override of the planted defect.
func (d *notesAppendIgnoresCancelDefect) Append(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.notesOverEmbedded.Append(p0, p1)
}

// notesAppendIgnoresCancel returns the planted defect of a Notes whose
// Append ignores the context over a fresh subject.
func notesAppendIgnoresCancel(h *NotesHandle) Notes {
	return &notesAppendIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesTextGrowsPastTheLimitAfterThreeDefect is the planted defect of a
// Notes whose Text reports a use above the limit from its fourth call on.
type notesTextGrowsPastTheLimitAfterThreeDefect struct {
	notesOverEmbedded
	calls int
}

// Text is the override of the planted defect.
func (d *notesTextGrowsPastTheLimitAfterThreeDefect) Text(
	p0 context.Context,
) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.notesOverEmbedded.Text(p0)
	d.calls++
	if d.calls > 3 {
		r0 = string(strings.Repeat("x", 8+1))
	}
	return r0, r1
}

// notesTextGrowsPastTheLimitAfterThree returns the planted defect of a
// Notes whose Text reports a use above the limit from its fourth call on
// over a fresh subject.
func notesTextGrowsPastTheLimitAfterThree(h *NotesHandle) Notes {
	return &notesTextGrowsPastTheLimitAfterThreeDefect{notesOverEmbedded: h.Fresh()}
}

// notesTextPanicsDefect is the planted defect of a Notes whose Text
// panics.
type notesTextPanicsDefect struct {
	notesOverEmbedded
}

// Text is the override of the planted defect.
func (*notesTextPanicsDefect) Text(p0 context.Context) (r0 string, r1 error) {
	panic("the planted defect of a Notes whose Text panics")
}

// notesTextPanics returns the planted defect of a Notes whose Text panics
// over a fresh subject.
func notesTextPanics(h *NotesHandle) Notes {
	return &notesTextPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesTextIgnoresCancelDefect is the planted defect of a Notes whose Text
// ignores the context.
type notesTextIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Text is the override of the planted defect.
func (d *notesTextIgnoresCancelDefect) Text(p0 context.Context) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.notesOverEmbedded.Text(p0)
}

// notesTextIgnoresCancel returns the planted defect of a Notes whose Text
// ignores the context over a fresh subject.
func notesTextIgnoresCancel(h *NotesHandle) Notes {
	return &notesTextIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesTextValueBesideDefect is the planted defect of a Notes whose Text
// returns a value beside an error.
type notesTextValueBesideDefect struct {
	notesOverEmbedded
}

// Text is the override of the planted defect.
func (d *notesTextValueBesideDefect) Text(p0 context.Context) (r0 string, r1 error) {
	r0, r1 = d.notesOverEmbedded.Text(p0)
	if r1 != nil {
		r0 = string("sample-b")
	}
	return r0, r1
}

// notesTextValueBeside returns the planted defect of a Notes whose Text
// returns a value beside an error over a fresh subject.
func notesTextValueBeside(h *NotesHandle) Notes {
	return &notesTextValueBesideDefect{notesOverEmbedded: h.Fresh()}
}

// notesPinPanicsDefect is the planted defect of a Notes whose Pin panics.
type notesPinPanicsDefect struct {
	notesOverEmbedded
}

// Pin is the override of the planted defect.
func (*notesPinPanicsDefect) Pin(p0 context.Context, p1 Item) (r0 error) {
	panic("the planted defect of a Notes whose Pin panics")
}

// notesPinPanics returns the planted defect of a Notes whose Pin panics
// over a fresh subject.
func notesPinPanics(h *NotesHandle) Notes {
	return &notesPinPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesPinIgnoresCancelDefect is the planted defect of a Notes whose Pin
// ignores the context.
type notesPinIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Pin is the override of the planted defect.
func (d *notesPinIgnoresCancelDefect) Pin(p0 context.Context, p1 Item) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.notesOverEmbedded.Pin(p0, p1)
}

// notesPinIgnoresCancel returns the planted defect of a Notes whose Pin
// ignores the context over a fresh subject.
func notesPinIgnoresCancel(h *NotesHandle) Notes {
	return &notesPinIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesHeadPanicsDefect is the planted defect of a Notes whose Head
// panics.
type notesHeadPanicsDefect struct {
	notesOverEmbedded
}

// Head is the override of the planted defect.
func (*notesHeadPanicsDefect) Head(p0 context.Context) (r0 Item, r1 error) {
	panic("the planted defect of a Notes whose Head panics")
}

// notesHeadPanics returns the planted defect of a Notes whose Head panics
// over a fresh subject.
func notesHeadPanics(h *NotesHandle) Notes {
	return &notesHeadPanicsDefect{notesOverEmbedded: h.Fresh()}
}

// notesHeadIgnoresCancelDefect is the planted defect of a Notes whose Head
// ignores the context.
type notesHeadIgnoresCancelDefect struct {
	notesOverEmbedded
}

// Head is the override of the planted defect.
func (d *notesHeadIgnoresCancelDefect) Head(p0 context.Context) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.notesOverEmbedded.Head(p0)
}

// notesHeadIgnoresCancel returns the planted defect of a Notes whose Head
// ignores the context over a fresh subject.
func notesHeadIgnoresCancel(h *NotesHandle) Notes {
	return &notesHeadIgnoresCancelDefect{notesOverEmbedded: h.Fresh()}
}

// notesHeadValueBesideDefect is the planted defect of a Notes whose Head
// returns a value beside an error.
type notesHeadValueBesideDefect struct {
	notesOverEmbedded
}

// Head is the override of the planted defect.
func (d *notesHeadValueBesideDefect) Head(p0 context.Context) (r0 Item, r1 error) {
	r0, r1 = d.notesOverEmbedded.Head(p0)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// notesHeadValueBeside returns the planted defect of a Notes whose Head
// returns a value beside an error over a fresh subject.
func notesHeadValueBeside(h *NotesHandle) Notes {
	return &notesHeadValueBesideDefect{notesOverEmbedded: h.Fresh()}
}

// notesActions returns the actions of the generated sequences of Notes
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func notesActions(c *prop.Case, subject func() Notes) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Push",
			Input: func(c *prop.Case, _ struct{}) any {
				item := c.Draw(prop.Of[string](), "push.item")
				return []any{item}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Push", args)
				values, err := notesPushCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Pop",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Pop", args)
				values, err := notesPopCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Items",
			Input: func(c *prop.Case, _ struct{}) any {
				prefix := c.Draw(prop.SampledFrom(keysString...), "items.prefix")
				return []any{prefix}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Items", args, args[0])
				values, err := notesItemsCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Append",
			Input: func(c *prop.Case, _ struct{}) any {
				text := c.Draw(prop.Of[string](), "append.text")
				return []any{text}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Append", args)
				values, err := notesAppendCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Text",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Text", args)
				values, err := notesTextCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Pin",
			Input: func(c *prop.Case, _ struct{}) any {
				it := c.Draw(prop.Of[Item](), "pin.it")
				return []any{it}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Pin", args)
				values, err := notesPinCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Head",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Head", args)
				values, err := notesHeadCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// notesPushCall calls Push on subject with the arguments args, and returns
// the values of the call beside its error.
func notesPushCall(ctx context.Context, subject Notes, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Push(ctx, a0)
}

// notesPopCall calls Pop on subject with the arguments args, and returns
// the values of the call beside its error.
func notesPopCall(ctx context.Context, subject Notes, args []any) ([]any, error) {
	return nil, subject.Pop(ctx)
}

// notesItemsCall calls Items on subject with the arguments args, and
// returns the values of the call beside its error.
func notesItemsCall(ctx context.Context, subject Notes, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Items(ctx, a0)
	return []any{r0}, err
}

// notesAppendCall calls Append on subject with the arguments args, and
// returns the values of the call beside its error.
func notesAppendCall(ctx context.Context, subject Notes, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Append(ctx, a0)
}

// notesTextCall calls Text on subject with the arguments args, and returns
// the values of the call beside its error.
func notesTextCall(ctx context.Context, subject Notes, args []any) ([]any, error) {
	r0, err := subject.Text(ctx)
	return []any{r0}, err
}

// notesPinCall calls Pin on subject with the arguments args, and returns
// the values of the call beside its error.
func notesPinCall(ctx context.Context, subject Notes, args []any) ([]any, error) {
	a0, _ := args[0].(Item)
	return nil, subject.Pin(ctx, a0)
}

// notesHeadCall calls Head on subject with the arguments args, and returns
// the values of the call beside its error.
func notesHeadCall(ctx context.Context, subject Notes, args []any) ([]any, error) {
	r0, err := subject.Head(ctx)
	return []any{r0}, err
}

// The derived checks of Set:
//
// Set.Merge/merges-in-any-order: two replicas that Add made diverge read
// the same state through Items after Merge merges them in either order.
//
// Set.Join/merges-in-any-order: two replicas that Add made diverge read
// the same state through Items after Join merges them in either order.
//
// Set.Add/smoke: Add returns without a panic for derived inputs.
//
// Set.Add/smoke#generated: Add returns without a panic for derived inputs.
//
// Set.Add/cancel: Add returns an error for a cancelled context.
//
// Set.Add/cancel#generated: Add returns an error for a cancelled context.
//
// Set.Items/smoke: Items returns without a panic for derived inputs.
//
// Set.Items/cancel: Items returns an error for a cancelled context.
//
// Set.Items/zero-on-error: Items returns the zero value beside an error.

// The derived checks of Set that witness their laws at the shallow end:
//
// Set.Add/smoke
//
// Set.Items/smoke

// The refused claims of Set:
//
// Set.Merge/cancel: eidos derives no sample of the parameter peer of
// Merge. The remedy is to author a sample of the parameter peer with the
// directive sample.
//
// Set.Join/cancel: eidos derives no sample of the parameter other of Join.
// The remedy is to author a sample of the parameter other with the
// directive sample.
//
// Set.Weigh/cancel: eidos derives no sample of the parameter peer of
// Weigh. The remedy is to author a sample of the parameter peer with the
// directive sample.
//
// Set.Weigh/merges-in-any-order: Weigh takes no input that another replica
// of Set can be. The remedy is to give Weigh an input of an interface that
// Set implements.

// The unclaimed classifications of Set:
//
// The shape:writer of Set.Merge has the stance silent. A write is an input
// of the other checks, and its result is an error alone.
//
// The shape:writer of Set.Join has the stance silent. A write is an input
// of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Set.Weigh has the stance silent. A write
// is an input of the other checks, and its result is an error alone.
//
// The shape:writer of Set.Add has the stance silent. A write is an input
// of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Set writes.
type (
	SetHandle = conform.Handle[Set]
	SetCheck  = conform.Check[Set]
	SetOption = conform.Option[Set]
)

// SetChecks contains the name of each callable of Set, with the identities
// of its checks.
var SetChecks = struct {
	Merge setMerge
	Join  setJoin
	Weigh setWeigh
	Add   setAdd
	Items setItems
}{
	Merge: setMerge{"Set.Merge"},
	Join:  setJoin{"Set.Join"},
	Weigh: setWeigh{"Set.Weigh"},
	Add:   setAdd{"Set.Add"},
	Items: setItems{"Set.Items"},
}

// setMerge is the name of Set.Merge, with the identities of its checks.
type setMerge struct {
	conform.Callable[Set]
}

// MergesInAnyOrder returns the identity of the check that two replicas
// that Add made diverge read the same state through Items after Merge
// merges them in either order.
func (c setMerge) MergesInAnyOrder() conform.ID[Set] {
	return c.Check("merges-in-any-order")
}

// setJoin is the name of Set.Join, with the identities of its checks.
type setJoin struct {
	conform.Callable[Set]
}

// MergesInAnyOrder returns the identity of the check that two replicas
// that Add made diverge read the same state through Items after Join
// merges them in either order.
func (c setJoin) MergesInAnyOrder() conform.ID[Set] {
	return c.Check("merges-in-any-order")
}

// setWeigh is the name of Set.Weigh, with the identities of its checks.
type setWeigh struct {
	conform.Callable[Set]
}

// setAdd is the name of Set.Add, with the identities of its checks.
type setAdd struct {
	conform.Callable[Set]
}

// Smoke returns the identity of the check that Add returns without a panic
// for derived inputs.
func (c setAdd) Smoke() conform.ID[Set] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Add returns
// without a panic for derived inputs.
func (c setAdd) SmokeGenerated() conform.ID[Set] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Add returns an error for a
// cancelled context.
func (c setAdd) Cancel() conform.ID[Set] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Add returns an
// error for a cancelled context.
func (c setAdd) CancelGenerated() conform.ID[Set] {
	return c.Check("cancel#generated")
}

// setItems is the name of Set.Items, with the identities of its checks.
type setItems struct {
	conform.Callable[Set]
}

// Smoke returns the identity of the check that Items returns without a
// panic for derived inputs.
func (c setItems) Smoke() conform.ID[Set] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Items returns an error for
// a cancelled context.
func (c setItems) Cancel() conform.ID[Set] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Items returns the
// zero value beside an error.
func (c setItems) ZeroOnError() conform.ID[Set] {
	return c.Check("zero-on-error")
}

// SetSuite returns the checks of Set, and the calls and the budgets of its
// callables.
func SetSuite() conform.Suite[Set] {
	return conform.Suite[Set]{
		Declaration: "Set",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Set]{
			{
				ID:          SetChecks.Merge.MergesInAnyOrder(),
				Statement:   "two replicas that Add made diverge read the same state through Items after Merge merges them in either order",
				Claim:       spec.MergesInAnyOrder,
				Probe:       spec.MergesCommute,
				Reason:      spec.OrderDependentMerge,
				Law:         spec.MergeOrderDoesNotMatter,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *SetHandle, c *prop.Case) {
					probe.MergesCommute(h, c, setActions, setReference, "Merge", "Add", "Items")
				},
				Defect: setMergeIgnoresTheOtherReplica,
			},
			{
				ID:          SetChecks.Join.MergesInAnyOrder(),
				Statement:   "two replicas that Add made diverge read the same state through Items after Join merges them in either order",
				Claim:       spec.MergesInAnyOrder,
				Probe:       spec.MergesCommute,
				Reason:      spec.OrderDependentMerge,
				Law:         spec.MergeOrderDoesNotMatter,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *SetHandle, c *prop.Case) {
					probe.MergesCommute(h, c, setActions, setReference, "Join", "Add", "Items")
				},
				Defect: setJoinIgnoresTheOtherReplica,
			},
			{
				ID:          SetChecks.Add.Smoke(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *SetHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, string("item-a"))
					})
				},
				Defect: setAddPanics,
			},
			{
				ID:          SetChecks.Add.SmokeGenerated(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *SetHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, drawnItem)
					})
				},
				Defect: setAddPanics,
			},
			{
				ID:          SetChecks.Add.Cancel(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *SetHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, string("item-a"))
					})
				},
				Defect: setAddIgnoresCancel,
			},
			{
				ID:          SetChecks.Add.CancelGenerated(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *SetHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx, drawnItem)
					})
				},
				Defect: setAddIgnoresCancel,
			},
			{
				ID:          SetChecks.Items.Smoke(),
				Statement:   "Items returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *SetHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Items(ctx)
						return err
					})
				},
				Defect: setItemsPanics,
			},
			{
				ID:          SetChecks.Items.Cancel(),
				Statement:   "Items returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *SetHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Items(ctx)
						return err
					})
				},
				Defect: setItemsIgnoresCancel,
			},
			{
				ID:          SetChecks.Items.ZeroOnError(),
				Statement:   "Items returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *SetHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) ([]string, error) {
						return h.Subject().Items(ctx)
					}, []string(nil))
				},
				Defect: setItemsValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "Set.Merge",
				Gate:     spec.InputSamples,
				Why:      "eidos derives no sample of the parameter peer of Merge",
				Remedy:   "author a sample of the parameter peer with the directive sample",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Set.Join",
				Gate:     spec.InputSamples,
				Why:      "eidos derives no sample of the parameter other of Join",
				Remedy:   "author a sample of the parameter other with the directive sample",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Set.Weigh",
				Gate:     spec.InputSamples,
				Why:      "eidos derives no sample of the parameter peer of Weigh",
				Remedy:   "author a sample of the parameter peer with the directive sample",
			},
			{
				Claim:    spec.MergesInAnyOrder,
				Callable: "Set.Weigh",
				Gate:     spec.TakesAReplica,
				Why:      "Weigh takes no input that another replica of Set can be",
				Remedy:   "give Weigh an input of an interface that Set implements",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Set.Merge",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Set.Join",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Set.Weigh",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Set.Add",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Set]{
			{
				Calls:    setMergeCalls,
				Over:     newSetOver,
				Callable: SetChecks.Merge.Callable,
			},
			{
				Calls:    setJoinCalls,
				Over:     newSetOver,
				Callable: SetChecks.Join.Callable,
			},
			{
				Calls:    setWeighCalls,
				Over:     newSetOver,
				Callable: SetChecks.Weigh.Callable,
			},
			{
				Calls:    setAddCalls,
				Over:     newSetOver,
				Callable: SetChecks.Add.Callable,
			},
			{
				Calls:    setItemsCalls,
				Over:     newSetOver,
				Callable: SetChecks.Items.Callable,
			},
		},
	}
}

// RunSetSuite applies the checks of Set, and the checks that opts add, to
// a fresh subject of newSubject for each check.
func RunSetSuite[T Set](t *testing.T, newSubject func() T, opts ...SetOption) {
	t.Helper()
	RunSetSuiteFromHandle(t, func(*SetHandle) T { return newSubject() }, opts...)
}

// RunSetSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunSetSuiteFromHandle[T Set](
	t *testing.T,
	newSubject func(h *SetHandle) T,
	opts ...SetOption,
) {
	t.Helper()
	conform.Run(t, SetSuite(), setSubjects(newSubject), opts...)
}

// FuzzSetSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzSetSuite[T Set](
	f *testing.F,
	callable conform.Callable[Set],
	newSubject func(h *SetHandle) T,
	opts ...SetOption,
) {
	f.Helper()
	conform.Fuzz(f, SetSuite(), callable, setSubjects(newSubject), opts...)
}

// RunSetBudgets applies the allocation ceilings of Set and its checks that
// count allocations, and the ceilings and the checks that opts add, to a
// fresh subject of newSubject for each case.
func RunSetBudgets[T Set](t *testing.T, newSubject func(h *SetHandle) T, opts ...SetOption) {
	t.Helper()
	conform.RunBudgets(t, SetSuite(), setSubjects(newSubject), opts...)
}

// BenchSetSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchSetSuite[T Set](
	b conform.B,
	callable conform.Callable[Set],
	newSubject func(h *SetHandle) T,
	opts ...SetOption,
) {
	b.Helper()
	conform.Bench(b, SetSuite(), callable, setSubjects(newSubject), opts...)
}

// setSubjects returns the factory of the subjects of Set that calls
// newSubject with the handle of each check.
func setSubjects[T Set](newSubject func(h *SetHandle) T) func(h *SetHandle) Set {
	return func(h *SetHandle) Set { return newSubject(h) }
}

// setMergeCalls generates the calls of Set.Merge with the inputs that each
// case draws.
var setMergeCalls = prop.Composite(func(c *prop.Case) func(context.Context, Set) error {
	peerItemsResults0 := c.Draw(prop.List(prop.Of[[]string]()), "peer-items-0")
	peerItemsResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, conform.ErrGenerated),
			prop.MinSize(len(peerItemsResults0)),
			prop.MaxSize(len(peerItemsResults0)),
		),
		"peer-items-1",
	)
	var peerItemsAt atomic.Int64
	peer := &setReplicaInput{
		onItems: func(context.Context) (r0 []string, r1 error) {
			i := int(peerItemsAt.Add(1) - 1)
			if i >= len(peerItemsResults0) {
				r1 = conform.ErrGenerated
				return r0, r1
			}
			return peerItemsResults0[i],
				peerItemsResults1[i]
		},
	}
	return func(ctx context.Context, subject Set) error {
		return subject.Merge(ctx, peer)
	}
})

// setJoinCalls generates the calls of Set.Join with the inputs that each
// case draws.
var setJoinCalls = prop.Composite(func(c *prop.Case) func(context.Context, Set) error {
	otherMergeResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"other-merge",
	)
	var otherMergeAt atomic.Int64
	otherJoinResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"other-join",
	)
	var otherJoinAt atomic.Int64
	otherWeighResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"other-weigh",
	)
	var otherWeighAt atomic.Int64
	otherAddResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"other-add",
	)
	var otherAddAt atomic.Int64
	otherItemsResults0 := c.Draw(prop.List(prop.Of[[]string]()), "other-items-0")
	otherItemsResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, conform.ErrGenerated),
			prop.MinSize(len(otherItemsResults0)),
			prop.MaxSize(len(otherItemsResults0)),
		),
		"other-items-1",
	)
	var otherItemsAt atomic.Int64
	other := &setSetInput{
		onMerge: func(context.Context, Replica) (r0 error) {
			i := int(otherMergeAt.Add(1) - 1)
			if i >= len(otherMergeResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return otherMergeResults[i]
		},
		onJoin: func(context.Context, Set) (r0 error) {
			i := int(otherJoinAt.Add(1) - 1)
			if i >= len(otherJoinResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return otherJoinResults[i]
		},
		onWeigh: func(context.Context, Replica, int) (r0 error) {
			i := int(otherWeighAt.Add(1) - 1)
			if i >= len(otherWeighResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return otherWeighResults[i]
		},
		onAdd: func(context.Context, string) (r0 error) {
			i := int(otherAddAt.Add(1) - 1)
			if i >= len(otherAddResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return otherAddResults[i]
		},
		onItems: func(context.Context) (r0 []string, r1 error) {
			i := int(otherItemsAt.Add(1) - 1)
			if i >= len(otherItemsResults0) {
				r1 = conform.ErrGenerated
				return r0, r1
			}
			return otherItemsResults0[i],
				otherItemsResults1[i]
		},
	}
	return func(ctx context.Context, subject Set) error {
		return subject.Join(ctx, other)
	}
})

// setWeighCalls generates the calls of Set.Weigh with the inputs that each
// case draws.
var setWeighCalls = prop.Composite(func(c *prop.Case) func(context.Context, Set) error {
	peerItemsResults0 := c.Draw(prop.List(prop.Of[[]string]()), "peer-items-0")
	peerItemsResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, conform.ErrGenerated),
			prop.MinSize(len(peerItemsResults0)),
			prop.MaxSize(len(peerItemsResults0)),
		),
		"peer-items-1",
	)
	var peerItemsAt atomic.Int64
	peer := &setReplicaInput{
		onItems: func(context.Context) (r0 []string, r1 error) {
			i := int(peerItemsAt.Add(1) - 1)
			if i >= len(peerItemsResults0) {
				r1 = conform.ErrGenerated
				return r0, r1
			}
			return peerItemsResults0[i],
				peerItemsResults1[i]
		},
	}
	weight := c.Draw(prop.Of[int](), "weight")
	return func(ctx context.Context, subject Set) error {
		return subject.Weigh(ctx, peer, weight)
	}
})

// setAddCalls generates the calls of Set.Add with the inputs that each
// case draws.
var setAddCalls = prop.Composite(func(c *prop.Case) func(context.Context, Set) error {
	item := c.Draw(prop.Of[string](), "item")
	return func(ctx context.Context, subject Set) error {
		return subject.Add(ctx, item)
	}
})

// setItemsCalls generates the calls of Set.Items with the inputs that each
// case draws.
var setItemsCalls = prop.Composite(func(*prop.Case) func(context.Context, Set) error {
	return func(ctx context.Context, subject Set) error {
		_, err := subject.Items(ctx)
		return err
	}
})

// setReplicaInputEmbedded is Replica under the name of the field that
// setReplicaInput embeds.
type setReplicaInputEmbedded = Replica

// setReplicaInput implements Replica. Its methods return the results that
// the generators of the suite of Set draw for them, in turn.
type setReplicaInput struct {
	setReplicaInputEmbedded
	onItems func(a0 context.Context) ([]string, error)
}

// Items calls the function of Items.
func (in setReplicaInput) Items(a0 context.Context) ([]string, error) {
	return in.onItems(a0)
}

// setSetInputEmbedded is Set under the name of the field that setSetInput
// embeds.
type setSetInputEmbedded = Set

// setSetInput implements Set. Its methods return the results that the
// generators of the suite of Set draw for them, in turn.
type setSetInput struct {
	setSetInputEmbedded
	onMerge func(a0 context.Context, a1 Replica) error
	onJoin  func(a0 context.Context, a1 Set) error
	onWeigh func(a0 context.Context, a1 Replica, a2 int) error
	onAdd   func(a0 context.Context, a1 string) error
	onItems func(a0 context.Context) ([]string, error)
}

// Merge calls the function of Merge.
func (in setSetInput) Merge(a0 context.Context, a1 Replica) error {
	return in.onMerge(a0, a1)
}

// Join calls the function of Join.
func (in setSetInput) Join(a0 context.Context, a1 Set) error {
	return in.onJoin(a0, a1)
}

// Weigh calls the function of Weigh.
func (in setSetInput) Weigh(a0 context.Context, a1 Replica, a2 int) error {
	return in.onWeigh(a0, a1, a2)
}

// Add calls the function of Add.
func (in setSetInput) Add(a0 context.Context, a1 string) error {
	return in.onAdd(a0, a1)
}

// Items calls the function of Items.
func (in setSetInput) Items(a0 context.Context) ([]string, error) {
	return in.onItems(a0)
}

// setOverEmbedded is Set under the name of the field that setOver embeds.
type setOverEmbedded = Set

// setOver is a subject of Set whose callables call spend and return their
// zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type setOver struct {
	setOverEmbedded
	spend func()
}

// Merge calls spend and returns its zero results.
func (o setOver) Merge(context.Context, Replica) (r0 error) {
	o.spend()
	return r0
}

// Join calls spend and returns its zero results.
func (o setOver) Join(context.Context, Set) (r0 error) {
	o.spend()
	return r0
}

// Weigh calls spend and returns its zero results.
func (o setOver) Weigh(context.Context, Replica, int) (r0 error) {
	o.spend()
	return r0
}

// Add calls spend and returns its zero results.
func (o setOver) Add(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Items calls spend and returns its zero results.
func (o setOver) Items(context.Context) (r0 []string, r1 error) {
	o.spend()
	return r0, r1
}

// newSetOver returns the planted defect of a ceiling of Set, whose
// callables call spend.
func newSetOver(_ *SetHandle, spend func()) Set {
	return setOver{spend: spend}
}

// setMergeIgnoresTheOtherReplicaDefect is the planted defect of a Set
// whose Merge ignores the other replica.
type setMergeIgnoresTheOtherReplicaDefect struct {
	setOverEmbedded
}

// Merge is the override of the planted defect.
func (*setMergeIgnoresTheOtherReplicaDefect) Merge(p0 context.Context, p1 Replica) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// setMergeIgnoresTheOtherReplica returns the planted defect of a Set whose
// Merge ignores the other replica over a fresh subject.
func setMergeIgnoresTheOtherReplica(h *SetHandle) Set {
	return &setMergeIgnoresTheOtherReplicaDefect{setOverEmbedded: h.Fresh()}
}

// setJoinIgnoresTheOtherReplicaDefect is the planted defect of a Set whose
// Join ignores the other replica.
type setJoinIgnoresTheOtherReplicaDefect struct {
	setOverEmbedded
}

// Join is the override of the planted defect.
func (*setJoinIgnoresTheOtherReplicaDefect) Join(p0 context.Context, p1 Set) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// setJoinIgnoresTheOtherReplica returns the planted defect of a Set whose
// Join ignores the other replica over a fresh subject.
func setJoinIgnoresTheOtherReplica(h *SetHandle) Set {
	return &setJoinIgnoresTheOtherReplicaDefect{setOverEmbedded: h.Fresh()}
}

// setAddPanicsDefect is the planted defect of a Set whose Add panics.
type setAddPanicsDefect struct {
	setOverEmbedded
}

// Add is the override of the planted defect.
func (*setAddPanicsDefect) Add(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Set whose Add panics")
}

// setAddPanics returns the planted defect of a Set whose Add panics over a
// fresh subject.
func setAddPanics(h *SetHandle) Set {
	return &setAddPanicsDefect{setOverEmbedded: h.Fresh()}
}

// setAddIgnoresCancelDefect is the planted defect of a Set whose Add
// ignores the context.
type setAddIgnoresCancelDefect struct {
	setOverEmbedded
}

// Add is the override of the planted defect.
func (d *setAddIgnoresCancelDefect) Add(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.setOverEmbedded.Add(p0, p1)
}

// setAddIgnoresCancel returns the planted defect of a Set whose Add
// ignores the context over a fresh subject.
func setAddIgnoresCancel(h *SetHandle) Set {
	return &setAddIgnoresCancelDefect{setOverEmbedded: h.Fresh()}
}

// setItemsPanicsDefect is the planted defect of a Set whose Items panics.
type setItemsPanicsDefect struct {
	setOverEmbedded
}

// Items is the override of the planted defect.
func (*setItemsPanicsDefect) Items(p0 context.Context) (r0 []string, r1 error) {
	panic("the planted defect of a Set whose Items panics")
}

// setItemsPanics returns the planted defect of a Set whose Items panics
// over a fresh subject.
func setItemsPanics(h *SetHandle) Set {
	return &setItemsPanicsDefect{setOverEmbedded: h.Fresh()}
}

// setItemsIgnoresCancelDefect is the planted defect of a Set whose Items
// ignores the context.
type setItemsIgnoresCancelDefect struct {
	setOverEmbedded
}

// Items is the override of the planted defect.
func (d *setItemsIgnoresCancelDefect) Items(p0 context.Context) (r0 []string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.setOverEmbedded.Items(p0)
}

// setItemsIgnoresCancel returns the planted defect of a Set whose Items
// ignores the context over a fresh subject.
func setItemsIgnoresCancel(h *SetHandle) Set {
	return &setItemsIgnoresCancelDefect{setOverEmbedded: h.Fresh()}
}

// setItemsValueBesideDefect is the planted defect of a Set whose Items
// returns a value beside an error.
type setItemsValueBesideDefect struct {
	setOverEmbedded
}

// Items is the override of the planted defect.
func (d *setItemsValueBesideDefect) Items(p0 context.Context) (r0 []string, r1 error) {
	r0, r1 = d.setOverEmbedded.Items(p0)
	if r1 != nil {
		r0 = []string{"sample-b"}
	}
	return r0, r1
}

// setItemsValueBeside returns the planted defect of a Set whose Items
// returns a value beside an error over a fresh subject.
func setItemsValueBeside(h *SetHandle) Set {
	return &setItemsValueBesideDefect{setOverEmbedded: h.Fresh()}
}

// setActions returns the actions of the generated sequences of Set over
// the subject that subject returns. The actions of a case take their keys
// from the pools that it draws.
func setActions(c *prop.Case, subject func() Set) []stateful.Action[struct{}] {
	keysReplica := c.Draw(
		prop.List(prop.Of[Replica](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.replica",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Merge",
			Input: func(c *prop.Case, _ struct{}) any {
				peerItemsResults0 := c.Draw(prop.List(prop.Of[[]string]()), "peer-items-0")
				peerItemsResults1 := c.Draw(
					prop.List(
						prop.SampledFrom[error](nil, conform.ErrGenerated),
						prop.MinSize(len(peerItemsResults0)),
						prop.MaxSize(len(peerItemsResults0)),
					),
					"peer-items-1",
				)
				var peerItemsAt atomic.Int64
				peer := &setReplicaInput{
					onItems: func(context.Context) (r0 []string, r1 error) {
						i := int(peerItemsAt.Add(1) - 1)
						if i >= len(peerItemsResults0) {
							r1 = conform.ErrGenerated
							return r0, r1
						}
						return peerItemsResults0[i],
							peerItemsResults1[i]
					},
				}
				return []any{peer}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Merge", args)
				values, err := setMergeCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Join",
			Input: func(c *prop.Case, _ struct{}) any {
				otherMergeResults := c.Draw(
					prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
					"other-merge",
				)
				var otherMergeAt atomic.Int64
				otherJoinResults := c.Draw(
					prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
					"other-join",
				)
				var otherJoinAt atomic.Int64
				otherWeighResults := c.Draw(
					prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
					"other-weigh",
				)
				var otherWeighAt atomic.Int64
				otherAddResults := c.Draw(
					prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
					"other-add",
				)
				var otherAddAt atomic.Int64
				otherItemsResults0 := c.Draw(prop.List(prop.Of[[]string]()), "other-items-0")
				otherItemsResults1 := c.Draw(
					prop.List(
						prop.SampledFrom[error](nil, conform.ErrGenerated),
						prop.MinSize(len(otherItemsResults0)),
						prop.MaxSize(len(otherItemsResults0)),
					),
					"other-items-1",
				)
				var otherItemsAt atomic.Int64
				other := &setSetInput{
					onMerge: func(context.Context, Replica) (r0 error) {
						i := int(otherMergeAt.Add(1) - 1)
						if i >= len(otherMergeResults) {
							r0 = conform.ErrGenerated
							return r0
						}
						return otherMergeResults[i]
					},
					onJoin: func(context.Context, Set) (r0 error) {
						i := int(otherJoinAt.Add(1) - 1)
						if i >= len(otherJoinResults) {
							r0 = conform.ErrGenerated
							return r0
						}
						return otherJoinResults[i]
					},
					onWeigh: func(context.Context, Replica, int) (r0 error) {
						i := int(otherWeighAt.Add(1) - 1)
						if i >= len(otherWeighResults) {
							r0 = conform.ErrGenerated
							return r0
						}
						return otherWeighResults[i]
					},
					onAdd: func(context.Context, string) (r0 error) {
						i := int(otherAddAt.Add(1) - 1)
						if i >= len(otherAddResults) {
							r0 = conform.ErrGenerated
							return r0
						}
						return otherAddResults[i]
					},
					onItems: func(context.Context) (r0 []string, r1 error) {
						i := int(otherItemsAt.Add(1) - 1)
						if i >= len(otherItemsResults0) {
							r1 = conform.ErrGenerated
							return r0, r1
						}
						return otherItemsResults0[i],
							otherItemsResults1[i]
					},
				}
				return []any{other}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Join", args)
				values, err := setJoinCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Weigh",
			Input: func(c *prop.Case, _ struct{}) any {
				peer := c.Draw(prop.SampledFrom(keysReplica...), "weigh.peer")
				weight := c.Draw(prop.Of[int](), "weigh.weight")
				return []any{peer, weight}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Weigh", args, args[0])
				values, err := setWeighCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Add",
			Input: func(c *prop.Case, _ struct{}) any {
				item := c.Draw(prop.Of[string](), "add.item")
				return []any{item}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Add", args)
				values, err := setAddCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Items",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Items", args)
				values, err := setItemsCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// setMergeCall calls Merge on subject with the arguments args, and returns
// the values of the call beside its error.
func setMergeCall(ctx context.Context, subject Set, args []any) ([]any, error) {
	a0, _ := args[0].(Replica)
	return nil, subject.Merge(ctx, a0)
}

// setJoinCall calls Join on subject with the arguments args, and returns
// the values of the call beside its error.
func setJoinCall(ctx context.Context, subject Set, args []any) ([]any, error) {
	a0, _ := args[0].(Set)
	return nil, subject.Join(ctx, a0)
}

// setWeighCall calls Weigh on subject with the arguments args, and returns
// the values of the call beside its error.
func setWeighCall(ctx context.Context, subject Set, args []any) ([]any, error) {
	a0, _ := args[0].(Replica)
	a1, _ := args[1].(int)
	return nil, subject.Weigh(ctx, a0, a1)
}

// setAddCall calls Add on subject with the arguments args, and returns the
// values of the call beside its error.
func setAddCall(ctx context.Context, subject Set, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Add(ctx, a0)
}

// setItemsCall calls Items on subject with the arguments args, and returns
// the values of the call beside its error.
func setItemsCall(ctx context.Context, subject Set, args []any) ([]any, error) {
	r0, err := subject.Items(ctx)
	return []any{r0}, err
}

// setReference returns the function that applies an operation of a history
// to the reference r by the name of its callable.
func setReference(r Set) func(context.Context, string, []any) conform.Output {
	return func(ctx context.Context, operation string, args []any) conform.Output {
		var values []any
		var err error
		switch operation {
		case "Merge":
			values, err = setMergeCall(ctx, r, args)
		case "Join":
			values, err = setJoinCall(ctx, r, args)
		case "Weigh":
			values, err = setWeighCall(ctx, r, args)
		case "Add":
			values, err = setAddCall(ctx, r, args)
		case "Items":
			values, err = setItemsCall(ctx, r, args)
		}
		return conform.OutputOf(err, values...)
	}
}

// The derived checks of Tally:
//
// Tally.Merge/merges-in-any-order: two replicas that Add made diverge read
// the same state through Count after Merge merges them in either order.
//
// Tally.Add/smoke: Add returns without a panic for derived inputs.
//
// Tally.Add/cancel: Add returns an error for a cancelled context.
//
// Tally.Count/smoke: Count returns without a panic for derived inputs.
//
// Tally.Count/cancel: Count returns an error for a cancelled context.
//
// Tally.Count/zero-on-error: Count returns the zero value beside an error.

// The derived checks of Tally that witness their laws at the shallow end:
//
// Tally.Add/smoke
//
// Tally.Count/smoke

// The refused claims of Tally:
//
// Tally.Merge/cancel: eidos derives no sample of the parameter other of
// Merge. The remedy is to author a sample of the parameter other with the
// directive sample.

// The unclaimed classifications of Tally:
//
// The shape:writer of Tally.Merge has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:lifecycle of Tally.Add has the stance silent. A lifecycle call
// returns nothing that a check could compare.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Tally writes.
type (
	TallyHandle = conform.Handle[TallyMethods]
	TallyCheck  = conform.Check[TallyMethods]
	TallyOption = conform.Option[TallyMethods]
)

// TallyMethods is the interface of the exported methods of Tally. Each
// subject of the suite of Tally implements it.
type TallyMethods interface {
	Merge(ctx context.Context, other *Tally) error
	Add(ctx context.Context) error
	Count(ctx context.Context) (int, error)
}

// TallyChecks contains the name of each callable of Tally, with the
// identities of its checks.
var TallyChecks = struct {
	Merge tallyMerge
	Add   tallyAdd
	Count tallyCount
}{
	Merge: tallyMerge{"Tally.Merge"},
	Add:   tallyAdd{"Tally.Add"},
	Count: tallyCount{"Tally.Count"},
}

// tallyMerge is the name of Tally.Merge, with the identities of its
// checks.
type tallyMerge struct {
	conform.Callable[TallyMethods]
}

// MergesInAnyOrder returns the identity of the check that two replicas
// that Add made diverge read the same state through Count after Merge
// merges them in either order.
func (c tallyMerge) MergesInAnyOrder() conform.ID[TallyMethods] {
	return c.Check("merges-in-any-order")
}

// tallyAdd is the name of Tally.Add, with the identities of its checks.
type tallyAdd struct {
	conform.Callable[TallyMethods]
}

// Smoke returns the identity of the check that Add returns without a panic
// for derived inputs.
func (c tallyAdd) Smoke() conform.ID[TallyMethods] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Add returns an error for a
// cancelled context.
func (c tallyAdd) Cancel() conform.ID[TallyMethods] {
	return c.Check("cancel")
}

// tallyCount is the name of Tally.Count, with the identities of its
// checks.
type tallyCount struct {
	conform.Callable[TallyMethods]
}

// Smoke returns the identity of the check that Count returns without a
// panic for derived inputs.
func (c tallyCount) Smoke() conform.ID[TallyMethods] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Count returns an error for
// a cancelled context.
func (c tallyCount) Cancel() conform.ID[TallyMethods] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Count returns the
// zero value beside an error.
func (c tallyCount) ZeroOnError() conform.ID[TallyMethods] {
	return c.Check("zero-on-error")
}

// TallySuite returns the checks of Tally, and the calls and the budgets of
// its callables.
func TallySuite() conform.Suite[TallyMethods] {
	return conform.Suite[TallyMethods]{
		Declaration: "Tally",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[TallyMethods]{
			{
				ID:          TallyChecks.Merge.MergesInAnyOrder(),
				Statement:   "two replicas that Add made diverge read the same state through Count after Merge merges them in either order",
				Claim:       spec.MergesInAnyOrder,
				Probe:       spec.MergesCommute,
				Reason:      spec.OrderDependentMerge,
				Law:         spec.MergeOrderDoesNotMatter,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *TallyHandle, c *prop.Case) {
					probe.MergesCommute(h, c, tallyActions, tallyReference, "Merge", "Add", "Count")
				},
				Defect: tallyMergeIgnoresTheOtherReplica,
			},
			{
				ID:          TallyChecks.Add.Smoke(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TallyHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx)
					})
				},
				Defect: tallyAddPanics,
			},
			{
				ID:          TallyChecks.Add.Cancel(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TallyHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx)
					})
				},
				Defect: tallyAddIgnoresCancel,
			},
			{
				ID:          TallyChecks.Count.Smoke(),
				Statement:   "Count returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *TallyHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx)
						return err
					})
				},
				Defect: tallyCountPanics,
			},
			{
				ID:          TallyChecks.Count.Cancel(),
				Statement:   "Count returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TallyHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx)
						return err
					})
				},
				Defect: tallyCountIgnoresCancel,
			},
			{
				ID:          TallyChecks.Count.ZeroOnError(),
				Statement:   "Count returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *TallyHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject().Count(ctx)
					}, int(0))
				},
				Defect: tallyCountValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "Tally.Merge",
				Gate:     spec.InputSamples,
				Why:      "eidos derives no sample of the parameter other of Merge",
				Remedy:   "author a sample of the parameter other with the directive sample",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Tally.Merge",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Tally.Add",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
		},
		Calls: []conform.Call[TallyMethods]{
			{
				Calls:    tallyMergeCalls,
				Over:     newTallyOver,
				Callable: TallyChecks.Merge.Callable,
			},
			{
				Calls:    tallyAddCalls,
				Over:     newTallyOver,
				Callable: TallyChecks.Add.Callable,
			},
			{
				Calls:    tallyCountCalls,
				Over:     newTallyOver,
				Callable: TallyChecks.Count.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(tallyTallyValues),
		},
	}
}

// RunTallySuite applies the checks of Tally, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunTallySuite[T TallyMethods](t *testing.T, newSubject func() T, opts ...TallyOption) {
	t.Helper()
	RunTallySuiteFromHandle(t, func(*TallyHandle) T { return newSubject() }, opts...)
}

// RunTallySuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunTallySuiteFromHandle[T TallyMethods](
	t *testing.T,
	newSubject func(h *TallyHandle) T,
	opts ...TallyOption,
) {
	t.Helper()
	conform.Run(t, TallySuite(), tallySubjects(newSubject), opts...)
}

// FuzzTallySuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzTallySuite[T TallyMethods](
	f *testing.F,
	callable conform.Callable[TallyMethods],
	newSubject func(h *TallyHandle) T,
	opts ...TallyOption,
) {
	f.Helper()
	conform.Fuzz(f, TallySuite(), callable, tallySubjects(newSubject), opts...)
}

// RunTallyBudgets applies the allocation ceilings of Tally and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunTallyBudgets[T TallyMethods](
	t *testing.T,
	newSubject func(h *TallyHandle) T,
	opts ...TallyOption,
) {
	t.Helper()
	conform.RunBudgets(t, TallySuite(), tallySubjects(newSubject), opts...)
}

// BenchTallySuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchTallySuite[T TallyMethods](
	b conform.B,
	callable conform.Callable[TallyMethods],
	newSubject func(h *TallyHandle) T,
	opts ...TallyOption,
) {
	b.Helper()
	conform.Bench(b, TallySuite(), callable, tallySubjects(newSubject), opts...)
}

// tallySubjects returns the factory of the subjects of Tally that calls
// newSubject with the handle of each check.
func tallySubjects[T TallyMethods](
	newSubject func(h *TallyHandle) T,
) func(h *TallyHandle) TallyMethods {
	return func(h *TallyHandle) TallyMethods { return newSubject(h) }
}

// tallyMergeCalls generates the calls of Tally.Merge with the inputs that
// each case draws.
var tallyMergeCalls = prop.Composite(func(c *prop.Case) func(context.Context, TallyMethods) error {
	other := c.Draw(prop.Of[Tally](), "other")
	return func(ctx context.Context, subject TallyMethods) error {
		return subject.Merge(ctx, &other)
	}
})

// tallyAddCalls generates the calls of Tally.Add with the inputs that each
// case draws.
var tallyAddCalls = prop.Composite(func(*prop.Case) func(context.Context, TallyMethods) error {
	return func(ctx context.Context, subject TallyMethods) error {
		return subject.Add(ctx)
	}
})

// tallyCountCalls generates the calls of Tally.Count with the inputs that
// each case draws.
var tallyCountCalls = prop.Composite(func(*prop.Case) func(context.Context, TallyMethods) error {
	return func(ctx context.Context, subject TallyMethods) error {
		_, err := subject.Count(ctx)
		return err
	}
})

// tallyTallyValues generates the starting value of Tally, which has no
// field to change.
var tallyTallyValues = prop.Composite(func(*prop.Case) Tally {
	return Tally{}
})

// tallyOverEmbedded is TallyMethods under the name of the field that
// tallyOver embeds.
type tallyOverEmbedded = TallyMethods

// tallyOver is a subject of Tally whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type tallyOver struct {
	tallyOverEmbedded
	spend func()
}

// Merge calls spend and returns its zero results.
func (o tallyOver) Merge(context.Context, *Tally) (r0 error) {
	o.spend()
	return r0
}

// Add calls spend and returns its zero results.
func (o tallyOver) Add(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Count calls spend and returns its zero results.
func (o tallyOver) Count(context.Context) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// newTallyOver returns the planted defect of a ceiling of Tally, whose
// callables call spend.
func newTallyOver(_ *TallyHandle, spend func()) TallyMethods {
	return tallyOver{spend: spend}
}

// tallyMergeIgnoresTheOtherReplicaDefect is the planted defect of a Tally
// whose Merge ignores the other replica.
type tallyMergeIgnoresTheOtherReplicaDefect struct {
	tallyOverEmbedded
}

// Merge is the override of the planted defect.
func (*tallyMergeIgnoresTheOtherReplicaDefect) Merge(p0 context.Context, p1 *Tally) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	return r0
}

// tallyMergeIgnoresTheOtherReplica returns the planted defect of a Tally
// whose Merge ignores the other replica over a fresh subject.
func tallyMergeIgnoresTheOtherReplica(h *TallyHandle) TallyMethods {
	return &tallyMergeIgnoresTheOtherReplicaDefect{tallyOverEmbedded: h.Fresh()}
}

// tallyAddPanicsDefect is the planted defect of a Tally whose Add panics.
type tallyAddPanicsDefect struct {
	tallyOverEmbedded
}

// Add is the override of the planted defect.
func (*tallyAddPanicsDefect) Add(p0 context.Context) (r0 error) {
	panic("the planted defect of a Tally whose Add panics")
}

// tallyAddPanics returns the planted defect of a Tally whose Add panics
// over a fresh subject.
func tallyAddPanics(h *TallyHandle) TallyMethods {
	return &tallyAddPanicsDefect{tallyOverEmbedded: h.Fresh()}
}

// tallyAddIgnoresCancelDefect is the planted defect of a Tally whose Add
// ignores the context.
type tallyAddIgnoresCancelDefect struct {
	tallyOverEmbedded
}

// Add is the override of the planted defect.
func (d *tallyAddIgnoresCancelDefect) Add(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.tallyOverEmbedded.Add(p0)
}

// tallyAddIgnoresCancel returns the planted defect of a Tally whose Add
// ignores the context over a fresh subject.
func tallyAddIgnoresCancel(h *TallyHandle) TallyMethods {
	return &tallyAddIgnoresCancelDefect{tallyOverEmbedded: h.Fresh()}
}

// tallyCountPanicsDefect is the planted defect of a Tally whose Count
// panics.
type tallyCountPanicsDefect struct {
	tallyOverEmbedded
}

// Count is the override of the planted defect.
func (*tallyCountPanicsDefect) Count(p0 context.Context) (r0 int, r1 error) {
	panic("the planted defect of a Tally whose Count panics")
}

// tallyCountPanics returns the planted defect of a Tally whose Count
// panics over a fresh subject.
func tallyCountPanics(h *TallyHandle) TallyMethods {
	return &tallyCountPanicsDefect{tallyOverEmbedded: h.Fresh()}
}

// tallyCountIgnoresCancelDefect is the planted defect of a Tally whose
// Count ignores the context.
type tallyCountIgnoresCancelDefect struct {
	tallyOverEmbedded
}

// Count is the override of the planted defect.
func (d *tallyCountIgnoresCancelDefect) Count(p0 context.Context) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.tallyOverEmbedded.Count(p0)
}

// tallyCountIgnoresCancel returns the planted defect of a Tally whose
// Count ignores the context over a fresh subject.
func tallyCountIgnoresCancel(h *TallyHandle) TallyMethods {
	return &tallyCountIgnoresCancelDefect{tallyOverEmbedded: h.Fresh()}
}

// tallyCountValueBesideDefect is the planted defect of a Tally whose Count
// returns a value beside an error.
type tallyCountValueBesideDefect struct {
	tallyOverEmbedded
}

// Count is the override of the planted defect.
func (d *tallyCountValueBesideDefect) Count(p0 context.Context) (r0 int, r1 error) {
	r0, r1 = d.tallyOverEmbedded.Count(p0)
	if r1 != nil {
		r0 = int(7)
	}
	return r0, r1
}

// tallyCountValueBeside returns the planted defect of a Tally whose Count
// returns a value beside an error over a fresh subject.
func tallyCountValueBeside(h *TallyHandle) TallyMethods {
	return &tallyCountValueBesideDefect{tallyOverEmbedded: h.Fresh()}
}

// tallyActions returns the actions of the generated sequences of Tally
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func tallyActions(c *prop.Case, subject func() TallyMethods) []stateful.Action[struct{}] {
	return []stateful.Action[struct{}]{
		{
			Name: "Merge",
			Input: func(c *prop.Case, _ struct{}) any {
				other := c.Draw(prop.Of[Tally](), "merge.other")
				return []any{&other}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Merge", args)
				values, err := tallyMergeCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Add",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Add", args)
				values, err := tallyAddCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Count",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Count", args)
				values, err := tallyCountCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// tallyMergeCall calls Merge on subject with the arguments args, and
// returns the values of the call beside its error.
func tallyMergeCall(ctx context.Context, subject TallyMethods, args []any) ([]any, error) {
	a0, _ := args[0].(*Tally)
	return nil, subject.Merge(ctx, a0)
}

// tallyAddCall calls Add on subject with the arguments args, and returns
// the values of the call beside its error.
func tallyAddCall(ctx context.Context, subject TallyMethods, args []any) ([]any, error) {
	return nil, subject.Add(ctx)
}

// tallyCountCall calls Count on subject with the arguments args, and
// returns the values of the call beside its error.
func tallyCountCall(ctx context.Context, subject TallyMethods, args []any) ([]any, error) {
	r0, err := subject.Count(ctx)
	return []any{r0}, err
}

// tallyReference returns the function that applies an operation of a
// history to the reference r by the name of its callable.
func tallyReference(r TallyMethods) func(context.Context, string, []any) conform.Output {
	return func(ctx context.Context, operation string, args []any) conform.Output {
		var values []any
		var err error
		switch operation {
		case "Merge":
			values, err = tallyMergeCall(ctx, r, args)
		case "Add":
			values, err = tallyAddCall(ctx, r, args)
		case "Count":
			values, err = tallyCountCall(ctx, r, args)
		}
		return conform.OutputOf(err, values...)
	}
}

// The derived checks of Pile:
//
// Pile.Absorb/smoke: Absorb returns without a panic for derived inputs.
//
// Pile.Absorb/smoke#generated: Absorb returns without a panic for derived
// inputs.
//
// Pile.Absorb/cancel: Absorb returns an error for a cancelled context.
//
// Pile.Absorb/cancel#generated: Absorb returns an error for a cancelled
// context.
//
// Pile.Add/smoke: Add returns without a panic for derived inputs.
//
// Pile.Add/cancel: Add returns an error for a cancelled context.
//
// Pile.Count/smoke: Count returns without a panic for derived inputs.
//
// Pile.Count/cancel: Count returns an error for a cancelled context.
//
// Pile.Count/zero-on-error: Count returns the zero value beside an error.

// The derived checks of Pile that witness their laws at the shallow end:
//
// Pile.Absorb/smoke
//
// Pile.Add/smoke
//
// Pile.Count/smoke

// The refused claims of Pile:
//
// Pile.Pour/cancel: eidos derives no sample of the parameter into of Pour.
// The remedy is to author a sample of the parameter into with the
// directive sample.
//
// Pile.Absorb/merges-in-any-order: Absorb takes no input that another
// replica of Pile can be. The remedy is to give Absorb an input of an
// interface that Pile implements.
//
// Pile.Pour/merges-in-any-order: Pour takes no input that another replica
// of Pile can be. The remedy is to give Pour an input of an interface that
// Pile implements.

// The unclaimed classifications of Pile:
//
// The shape:writer of Pile.Absorb has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Pile.Pour has the stance silent. A write is an input
// of the other checks, and its result is an error alone.
//
// The shape:lifecycle of Pile.Add has the stance silent. A lifecycle call
// returns nothing that a check could compare.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Pile writes.
type (
	PileHandle = conform.Handle[Pile]
	PileCheck  = conform.Check[Pile]
	PileOption = conform.Option[Pile]
)

// PileChecks contains the name of each callable of Pile, with the
// identities of its checks.
var PileChecks = struct {
	Absorb pileAbsorb
	Pour   pilePour
	Add    pileAdd
	Count  pileCount
}{
	Absorb: pileAbsorb{"Pile.Absorb"},
	Pour:   pilePour{"Pile.Pour"},
	Add:    pileAdd{"Pile.Add"},
	Count:  pileCount{"Pile.Count"},
}

// pileAbsorb is the name of Pile.Absorb, with the identities of its
// checks.
type pileAbsorb struct {
	conform.Callable[Pile]
}

// Smoke returns the identity of the check that Absorb returns without a
// panic for derived inputs.
func (c pileAbsorb) Smoke() conform.ID[Pile] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Absorb returns
// without a panic for derived inputs.
func (c pileAbsorb) SmokeGenerated() conform.ID[Pile] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Absorb returns an error
// for a cancelled context.
func (c pileAbsorb) Cancel() conform.ID[Pile] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Absorb returns an
// error for a cancelled context.
func (c pileAbsorb) CancelGenerated() conform.ID[Pile] {
	return c.Check("cancel#generated")
}

// pilePour is the name of Pile.Pour, with the identities of its checks.
type pilePour struct {
	conform.Callable[Pile]
}

// pileAdd is the name of Pile.Add, with the identities of its checks.
type pileAdd struct {
	conform.Callable[Pile]
}

// Smoke returns the identity of the check that Add returns without a panic
// for derived inputs.
func (c pileAdd) Smoke() conform.ID[Pile] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Add returns an error for a
// cancelled context.
func (c pileAdd) Cancel() conform.ID[Pile] {
	return c.Check("cancel")
}

// pileCount is the name of Pile.Count, with the identities of its checks.
type pileCount struct {
	conform.Callable[Pile]
}

// Smoke returns the identity of the check that Count returns without a
// panic for derived inputs.
func (c pileCount) Smoke() conform.ID[Pile] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Count returns an error for
// a cancelled context.
func (c pileCount) Cancel() conform.ID[Pile] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Count returns the
// zero value beside an error.
func (c pileCount) ZeroOnError() conform.ID[Pile] {
	return c.Check("zero-on-error")
}

// PileSuite returns the checks of Pile, and the calls and the budgets of
// its callables.
func PileSuite() conform.Suite[Pile] {
	return conform.Suite[Pile]{
		Declaration: "Pile",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Pile]{
			{
				ID:          PileChecks.Absorb.Smoke(),
				Statement:   "Absorb returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PileHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Absorb(ctx, int(42))
					})
				},
				Defect: pileAbsorbPanics,
			},
			{
				ID:          PileChecks.Absorb.SmokeGenerated(),
				Statement:   "Absorb returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *PileHandle, c *prop.Case) {
					drawnN := c.Draw(prop.Of[int](), "n")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Absorb(ctx, drawnN)
					})
				},
				Defect: pileAbsorbPanics,
			},
			{
				ID:          PileChecks.Absorb.Cancel(),
				Statement:   "Absorb returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PileHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Absorb(ctx, int(42))
					})
				},
				Defect: pileAbsorbIgnoresCancel,
			},
			{
				ID:          PileChecks.Absorb.CancelGenerated(),
				Statement:   "Absorb returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *PileHandle, c *prop.Case) {
					drawnN := c.Draw(prop.Of[int](), "n")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Absorb(ctx, drawnN)
					})
				},
				Defect: pileAbsorbIgnoresCancel,
			},
			{
				ID:          PileChecks.Add.Smoke(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PileHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx)
					})
				},
				Defect: pileAddPanics,
			},
			{
				ID:          PileChecks.Add.Cancel(),
				Statement:   "Add returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PileHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Add(ctx)
					})
				},
				Defect: pileAddIgnoresCancel,
			},
			{
				ID:          PileChecks.Count.Smoke(),
				Statement:   "Count returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PileHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx)
						return err
					})
				},
				Defect: pileCountPanics,
			},
			{
				ID:          PileChecks.Count.Cancel(),
				Statement:   "Count returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PileHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx)
						return err
					})
				},
				Defect: pileCountIgnoresCancel,
			},
			{
				ID:          PileChecks.Count.ZeroOnError(),
				Statement:   "Count returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PileHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject().Count(ctx)
					}, int(0))
				},
				Defect: pileCountValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "Pile.Pour",
				Gate:     spec.InputSamples,
				Why:      "eidos derives no sample of the parameter into of Pour",
				Remedy:   "author a sample of the parameter into with the directive sample",
			},
			{
				Claim:    spec.MergesInAnyOrder,
				Callable: "Pile.Absorb",
				Gate:     spec.TakesAReplica,
				Why:      "Absorb takes no input that another replica of Pile can be",
				Remedy:   "give Absorb an input of an interface that Pile implements",
			},
			{
				Claim:    spec.MergesInAnyOrder,
				Callable: "Pile.Pour",
				Gate:     spec.TakesAReplica,
				Why:      "Pour takes no input that another replica of Pile can be",
				Remedy:   "give Pour an input of an interface that Pile implements",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Pile.Absorb",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Pile.Pour",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:lifecycle",
				Callable:       "Pile.Add",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
		},
		Calls: []conform.Call[Pile]{
			{
				Calls:    pileAbsorbCalls,
				Over:     newPileOver,
				Callable: PileChecks.Absorb.Callable,
			},
			{
				Calls:    pilePourCalls,
				Over:     newPileOver,
				Callable: PileChecks.Pour.Callable,
			},
			{
				Calls:    pileAddCalls,
				Over:     newPileOver,
				Callable: PileChecks.Add.Callable,
			},
			{
				Calls:    pileCountCalls,
				Over:     newPileOver,
				Callable: PileChecks.Count.Callable,
			},
		},
	}
}

// RunPileSuite applies the checks of Pile, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunPileSuite[T Pile](t *testing.T, newSubject func() T, opts ...PileOption) {
	t.Helper()
	RunPileSuiteFromHandle(t, func(*PileHandle) T { return newSubject() }, opts...)
}

// RunPileSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunPileSuiteFromHandle[T Pile](
	t *testing.T,
	newSubject func(h *PileHandle) T,
	opts ...PileOption,
) {
	t.Helper()
	conform.Run(t, PileSuite(), pileSubjects(newSubject), opts...)
}

// FuzzPileSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzPileSuite[T Pile](
	f *testing.F,
	callable conform.Callable[Pile],
	newSubject func(h *PileHandle) T,
	opts ...PileOption,
) {
	f.Helper()
	conform.Fuzz(f, PileSuite(), callable, pileSubjects(newSubject), opts...)
}

// RunPileBudgets applies the allocation ceilings of Pile and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunPileBudgets[T Pile](t *testing.T, newSubject func(h *PileHandle) T, opts ...PileOption) {
	t.Helper()
	conform.RunBudgets(t, PileSuite(), pileSubjects(newSubject), opts...)
}

// BenchPileSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchPileSuite[T Pile](
	b conform.B,
	callable conform.Callable[Pile],
	newSubject func(h *PileHandle) T,
	opts ...PileOption,
) {
	b.Helper()
	conform.Bench(b, PileSuite(), callable, pileSubjects(newSubject), opts...)
}

// pileSubjects returns the factory of the subjects of Pile that calls
// newSubject with the handle of each check.
func pileSubjects[T Pile](newSubject func(h *PileHandle) T) func(h *PileHandle) Pile {
	return func(h *PileHandle) Pile { return newSubject(h) }
}

// pileAbsorbCalls generates the calls of Pile.Absorb with the inputs that
// each case draws.
var pileAbsorbCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pile) error {
	n := c.Draw(prop.Of[int](), "n")
	return func(ctx context.Context, subject Pile) error {
		return subject.Absorb(ctx, n)
	}
})

// pilePourCalls generates the calls of Pile.Pour with the inputs that each
// case draws.
var pilePourCalls = prop.Composite(func(c *prop.Case) func(context.Context, Pile) error {
	intoDrainResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"into-drain",
	)
	var intoDrainAt atomic.Int64
	into := &pileSinkInput{
		onDrain: func(context.Context) (r0 error) {
			i := int(intoDrainAt.Add(1) - 1)
			if i >= len(intoDrainResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return intoDrainResults[i]
		},
	}
	return func(ctx context.Context, subject Pile) error {
		return subject.Pour(ctx, into)
	}
})

// pileAddCalls generates the calls of Pile.Add with the inputs that each
// case draws.
var pileAddCalls = prop.Composite(func(*prop.Case) func(context.Context, Pile) error {
	return func(ctx context.Context, subject Pile) error {
		return subject.Add(ctx)
	}
})

// pileCountCalls generates the calls of Pile.Count with the inputs that
// each case draws.
var pileCountCalls = prop.Composite(func(*prop.Case) func(context.Context, Pile) error {
	return func(ctx context.Context, subject Pile) error {
		_, err := subject.Count(ctx)
		return err
	}
})

// pileSinkInputEmbedded is Sink under the name of the field that
// pileSinkInput embeds.
type pileSinkInputEmbedded = Sink

// pileSinkInput implements Sink. Its methods return the results that the
// generators of the suite of Pile draw for them, in turn.
type pileSinkInput struct {
	pileSinkInputEmbedded
	onDrain func(a0 context.Context) error
}

// Drain calls the function of Drain.
func (in pileSinkInput) Drain(a0 context.Context) error {
	return in.onDrain(a0)
}

// pileOverEmbedded is Pile under the name of the field that pileOver
// embeds.
type pileOverEmbedded = Pile

// pileOver is a subject of Pile whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type pileOver struct {
	pileOverEmbedded
	spend func()
}

// Absorb calls spend and returns its zero results.
func (o pileOver) Absorb(context.Context, int) (r0 error) {
	o.spend()
	return r0
}

// Pour calls spend and returns its zero results.
func (o pileOver) Pour(context.Context, Sink) (r0 error) {
	o.spend()
	return r0
}

// Add calls spend and returns its zero results.
func (o pileOver) Add(context.Context) (r0 error) {
	o.spend()
	return r0
}

// Count calls spend and returns its zero results.
func (o pileOver) Count(context.Context) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// newPileOver returns the planted defect of a ceiling of Pile, whose
// callables call spend.
func newPileOver(_ *PileHandle, spend func()) Pile {
	return pileOver{spend: spend}
}

// pileAbsorbPanicsDefect is the planted defect of a Pile whose Absorb
// panics.
type pileAbsorbPanicsDefect struct {
	pileOverEmbedded
}

// Absorb is the override of the planted defect.
func (*pileAbsorbPanicsDefect) Absorb(p0 context.Context, p1 int) (r0 error) {
	panic("the planted defect of a Pile whose Absorb panics")
}

// pileAbsorbPanics returns the planted defect of a Pile whose Absorb
// panics over a fresh subject.
func pileAbsorbPanics(h *PileHandle) Pile {
	return &pileAbsorbPanicsDefect{pileOverEmbedded: h.Fresh()}
}

// pileAbsorbIgnoresCancelDefect is the planted defect of a Pile whose
// Absorb ignores the context.
type pileAbsorbIgnoresCancelDefect struct {
	pileOverEmbedded
}

// Absorb is the override of the planted defect.
func (d *pileAbsorbIgnoresCancelDefect) Absorb(p0 context.Context, p1 int) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.pileOverEmbedded.Absorb(p0, p1)
}

// pileAbsorbIgnoresCancel returns the planted defect of a Pile whose
// Absorb ignores the context over a fresh subject.
func pileAbsorbIgnoresCancel(h *PileHandle) Pile {
	return &pileAbsorbIgnoresCancelDefect{pileOverEmbedded: h.Fresh()}
}

// pileAddPanicsDefect is the planted defect of a Pile whose Add panics.
type pileAddPanicsDefect struct {
	pileOverEmbedded
}

// Add is the override of the planted defect.
func (*pileAddPanicsDefect) Add(p0 context.Context) (r0 error) {
	panic("the planted defect of a Pile whose Add panics")
}

// pileAddPanics returns the planted defect of a Pile whose Add panics over
// a fresh subject.
func pileAddPanics(h *PileHandle) Pile {
	return &pileAddPanicsDefect{pileOverEmbedded: h.Fresh()}
}

// pileAddIgnoresCancelDefect is the planted defect of a Pile whose Add
// ignores the context.
type pileAddIgnoresCancelDefect struct {
	pileOverEmbedded
}

// Add is the override of the planted defect.
func (d *pileAddIgnoresCancelDefect) Add(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.pileOverEmbedded.Add(p0)
}

// pileAddIgnoresCancel returns the planted defect of a Pile whose Add
// ignores the context over a fresh subject.
func pileAddIgnoresCancel(h *PileHandle) Pile {
	return &pileAddIgnoresCancelDefect{pileOverEmbedded: h.Fresh()}
}

// pileCountPanicsDefect is the planted defect of a Pile whose Count
// panics.
type pileCountPanicsDefect struct {
	pileOverEmbedded
}

// Count is the override of the planted defect.
func (*pileCountPanicsDefect) Count(p0 context.Context) (r0 int, r1 error) {
	panic("the planted defect of a Pile whose Count panics")
}

// pileCountPanics returns the planted defect of a Pile whose Count panics
// over a fresh subject.
func pileCountPanics(h *PileHandle) Pile {
	return &pileCountPanicsDefect{pileOverEmbedded: h.Fresh()}
}

// pileCountIgnoresCancelDefect is the planted defect of a Pile whose Count
// ignores the context.
type pileCountIgnoresCancelDefect struct {
	pileOverEmbedded
}

// Count is the override of the planted defect.
func (d *pileCountIgnoresCancelDefect) Count(p0 context.Context) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.pileOverEmbedded.Count(p0)
}

// pileCountIgnoresCancel returns the planted defect of a Pile whose Count
// ignores the context over a fresh subject.
func pileCountIgnoresCancel(h *PileHandle) Pile {
	return &pileCountIgnoresCancelDefect{pileOverEmbedded: h.Fresh()}
}

// pileCountValueBesideDefect is the planted defect of a Pile whose Count
// returns a value beside an error.
type pileCountValueBesideDefect struct {
	pileOverEmbedded
}

// Count is the override of the planted defect.
func (d *pileCountValueBesideDefect) Count(p0 context.Context) (r0 int, r1 error) {
	r0, r1 = d.pileOverEmbedded.Count(p0)
	if r1 != nil {
		r0 = int(7)
	}
	return r0, r1
}

// pileCountValueBeside returns the planted defect of a Pile whose Count
// returns a value beside an error over a fresh subject.
func pileCountValueBeside(h *PileHandle) Pile {
	return &pileCountValueBesideDefect{pileOverEmbedded: h.Fresh()}
}

// The derived checks of Run:
//
// Run.Advance/smoke: Advance returns without a panic for derived inputs.
//
// Run.Advance/smoke#generated: Advance returns without a panic for derived
// inputs.
//
// Run.Advance/cancel: Advance returns an error for a cancelled context.
//
// Run.Advance/cancel#generated: Advance returns an error for a cancelled
// context.
//
// Run.Advance/follows-its-graph: each call of Advance leaves the state
// that Status reads at the state before it or at one of its successors,
// over any sequence of calls.
//
// Run.Retry/smoke: Retry returns without a panic for derived inputs.
//
// Run.Retry/smoke#generated: Retry returns without a panic for derived
// inputs.
//
// Run.Retry/cancel: Retry returns an error for a cancelled context.
//
// Run.Retry/cancel#generated: Retry returns an error for a cancelled
// context.
//
// Run.Retry/follows-its-graph: each call of Retry leaves the state that
// Status reads at the state before it or at one of its successors, over
// any sequence of calls.
//
// Run.Restart/smoke: Restart returns without a panic for derived inputs.
//
// Run.Restart/smoke#generated: Restart returns without a panic for derived
// inputs.
//
// Run.Restart/cancel: Restart returns an error for a cancelled context.
//
// Run.Restart/cancel#generated: Restart returns an error for a cancelled
// context.
//
// Run.Toggle/smoke: Toggle returns without a panic for derived inputs.
//
// Run.Toggle/smoke#generated: Toggle returns without a panic for derived
// inputs.
//
// Run.Toggle/cancel: Toggle returns an error for a cancelled context.
//
// Run.Toggle/cancel#generated: Toggle returns an error for a cancelled
// context.
//
// Run.Toggle/follows-its-graph: each call of Toggle leaves the state that
// Status reads at the state before it or at one of its successors, over
// any sequence of calls.
//
// Run.Status/smoke: Status returns without a panic for derived inputs.
//
// Run.Status/smoke#generated: Status returns without a panic for derived
// inputs.
//
// Run.Status/cancel: Status returns an error for a cancelled context.
//
// Run.Status/cancel#generated: Status returns an error for a cancelled
// context.

// The derived checks of Run that witness their laws at the shallow end:
//
// Run.Advance/smoke
//
// Run.Retry/smoke
//
// Run.Restart/smoke
//
// Run.Toggle/smoke
//
// Run.Status/smoke

// The refused claims of Run:
//
// Run.Restart/follows-its-graph: the variable that transitions names for
// Restart is no graph of the states that observe reads. The remedy is to
// declare transitions as a map from each state to the list of its
// successors.

// The unclaimed classifications of Run:
//
// The shape:writer of Run.Advance has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:writer of Run.Retry has the stance silent. A write is an input
// of the other checks, and its result is an error alone.
//
// The shape:writer of Run.Restart has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:composite-writer of Run.Toggle has the stance silent. A write
// is an input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Run writes.
type (
	RunHandle = conform.Handle[Run]
	RunCheck  = conform.Check[Run]
	RunOption = conform.Option[Run]
)

// RunChecks contains the name of each callable of Run, with the identities
// of its checks.
var RunChecks = struct {
	Advance runAdvance
	Retry   runRetry
	Restart runRestart
	Toggle  runToggle
	Status  runStatus
}{
	Advance: runAdvance{"Run.Advance"},
	Retry:   runRetry{"Run.Retry"},
	Restart: runRestart{"Run.Restart"},
	Toggle:  runToggle{"Run.Toggle"},
	Status:  runStatus{"Run.Status"},
}

// runAdvance is the name of Run.Advance, with the identities of its
// checks.
type runAdvance struct {
	conform.Callable[Run]
}

// Smoke returns the identity of the check that Advance returns without a
// panic for derived inputs.
func (c runAdvance) Smoke() conform.ID[Run] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Advance returns
// without a panic for derived inputs.
func (c runAdvance) SmokeGenerated() conform.ID[Run] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Advance returns an error
// for a cancelled context.
func (c runAdvance) Cancel() conform.ID[Run] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Advance returns
// an error for a cancelled context.
func (c runAdvance) CancelGenerated() conform.ID[Run] {
	return c.Check("cancel#generated")
}

// FollowsItsGraph returns the identity of the check that each call of
// Advance leaves the state that Status reads at the state before it or at
// one of its successors, over any sequence of calls.
func (c runAdvance) FollowsItsGraph() conform.ID[Run] {
	return c.Check("follows-its-graph")
}

// runRetry is the name of Run.Retry, with the identities of its checks.
type runRetry struct {
	conform.Callable[Run]
}

// Smoke returns the identity of the check that Retry returns without a
// panic for derived inputs.
func (c runRetry) Smoke() conform.ID[Run] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Retry returns
// without a panic for derived inputs.
func (c runRetry) SmokeGenerated() conform.ID[Run] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Retry returns an error for
// a cancelled context.
func (c runRetry) Cancel() conform.ID[Run] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Retry returns an
// error for a cancelled context.
func (c runRetry) CancelGenerated() conform.ID[Run] {
	return c.Check("cancel#generated")
}

// FollowsItsGraph returns the identity of the check that each call of
// Retry leaves the state that Status reads at the state before it or at
// one of its successors, over any sequence of calls.
func (c runRetry) FollowsItsGraph() conform.ID[Run] {
	return c.Check("follows-its-graph")
}

// runRestart is the name of Run.Restart, with the identities of its
// checks.
type runRestart struct {
	conform.Callable[Run]
}

// Smoke returns the identity of the check that Restart returns without a
// panic for derived inputs.
func (c runRestart) Smoke() conform.ID[Run] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Restart returns
// without a panic for derived inputs.
func (c runRestart) SmokeGenerated() conform.ID[Run] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Restart returns an error
// for a cancelled context.
func (c runRestart) Cancel() conform.ID[Run] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Restart returns
// an error for a cancelled context.
func (c runRestart) CancelGenerated() conform.ID[Run] {
	return c.Check("cancel#generated")
}

// runToggle is the name of Run.Toggle, with the identities of its checks.
type runToggle struct {
	conform.Callable[Run]
}

// Smoke returns the identity of the check that Toggle returns without a
// panic for derived inputs.
func (c runToggle) Smoke() conform.ID[Run] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Toggle returns
// without a panic for derived inputs.
func (c runToggle) SmokeGenerated() conform.ID[Run] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Toggle returns an error
// for a cancelled context.
func (c runToggle) Cancel() conform.ID[Run] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Toggle returns an
// error for a cancelled context.
func (c runToggle) CancelGenerated() conform.ID[Run] {
	return c.Check("cancel#generated")
}

// FollowsItsGraph returns the identity of the check that each call of
// Toggle leaves the state that Status reads at the state before it or at
// one of its successors, over any sequence of calls.
func (c runToggle) FollowsItsGraph() conform.ID[Run] {
	return c.Check("follows-its-graph")
}

// runStatus is the name of Run.Status, with the identities of its checks.
type runStatus struct {
	conform.Callable[Run]
}

// Smoke returns the identity of the check that Status returns without a
// panic for derived inputs.
func (c runStatus) Smoke() conform.ID[Run] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Status returns
// without a panic for derived inputs.
func (c runStatus) SmokeGenerated() conform.ID[Run] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Status returns an error
// for a cancelled context.
func (c runStatus) Cancel() conform.ID[Run] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Status returns an
// error for a cancelled context.
func (c runStatus) CancelGenerated() conform.ID[Run] {
	return c.Check("cancel#generated")
}

// RunSuite returns the checks of Run, and the calls and the budgets of its
// callables.
func RunSuite() conform.Suite[Run] {
	return conform.Suite[Run]{
		Declaration: "Run",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Run]{
			{
				ID:          RunChecks.Advance.Smoke(),
				Statement:   "Advance returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RunHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Advance(ctx, string("key-a"))
					})
				},
				Defect: runAdvancePanics,
			},
			{
				ID:          RunChecks.Advance.SmokeGenerated(),
				Statement:   "Advance returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Advance(ctx, drawnKey)
					})
				},
				Defect: runAdvancePanics,
			},
			{
				ID:          RunChecks.Advance.Cancel(),
				Statement:   "Advance returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RunHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Advance(ctx, string("key-a"))
					})
				},
				Defect: runAdvanceIgnoresCancel,
			},
			{
				ID:          RunChecks.Advance.CancelGenerated(),
				Statement:   "Advance returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Advance(ctx, drawnKey)
					})
				},
				Defect: runAdvanceIgnoresCancel,
			},
			{
				ID:          RunChecks.Advance.FollowsItsGraph(),
				Statement:   "each call of Advance leaves the state that Status reads at the state before it or at one of its successors, over any sequence of calls",
				Claim:       spec.FollowsItsGraph,
				Probe:       spec.StatesMoveAlongTheGraph,
				Reason:      spec.ForbiddenTransition,
				Law:         spec.StatesFollowTheGraph,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *RunHandle, c *prop.Case) {
					probe.StatesMoveAlongTheGraph(
						h,
						c,
						runActions,
						"Advance",
						func(ctx context.Context, args []any) conform.Output {
							values, err := runStatusCall(
								ctx,
								h.Subject(),
								[]any{args[0], bool(true)},
							)
							return conform.OutputOf(err, values...)
						},
						Graph,
					)
				},
				Defect: runAdvanceAdmitsAnyTransitionAfterThree,
			},
			{
				ID:          RunChecks.Retry.Smoke(),
				Statement:   "Retry returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RunHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Retry(ctx, string("key-a"))
					})
				},
				Defect: runRetryPanics,
			},
			{
				ID:          RunChecks.Retry.SmokeGenerated(),
				Statement:   "Retry returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Retry(ctx, drawnKey)
					})
				},
				Defect: runRetryPanics,
			},
			{
				ID:          RunChecks.Retry.Cancel(),
				Statement:   "Retry returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RunHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Retry(ctx, string("key-a"))
					})
				},
				Defect: runRetryIgnoresCancel,
			},
			{
				ID:          RunChecks.Retry.CancelGenerated(),
				Statement:   "Retry returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Retry(ctx, drawnKey)
					})
				},
				Defect: runRetryIgnoresCancel,
			},
			{
				ID:          RunChecks.Retry.FollowsItsGraph(),
				Statement:   "each call of Retry leaves the state that Status reads at the state before it or at one of its successors, over any sequence of calls",
				Claim:       spec.FollowsItsGraph,
				Probe:       spec.StatesMoveAlongTheGraph,
				Reason:      spec.ForbiddenTransition,
				Law:         spec.StatesFollowTheGraph,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *RunHandle, c *prop.Case) {
					probe.StatesMoveAlongTheGraph(
						h,
						c,
						runActions,
						"Retry",
						func(ctx context.Context, args []any) conform.Output {
							values, err := runStatusCall(
								ctx,
								h.Subject(),
								[]any{args[0], bool(true)},
							)
							return conform.OutputOf(err, values...)
						},
						Retries,
					)
				},
				Defect: runRetryAdmitsAnyTransitionAfterThree,
			},
			{
				ID:          RunChecks.Restart.Smoke(),
				Statement:   "Restart returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RunHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Restart(ctx, string("key-a"))
					})
				},
				Defect: runRestartPanics,
			},
			{
				ID:          RunChecks.Restart.SmokeGenerated(),
				Statement:   "Restart returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Restart(ctx, drawnKey)
					})
				},
				Defect: runRestartPanics,
			},
			{
				ID:          RunChecks.Restart.Cancel(),
				Statement:   "Restart returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RunHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Restart(ctx, string("key-a"))
					})
				},
				Defect: runRestartIgnoresCancel,
			},
			{
				ID:          RunChecks.Restart.CancelGenerated(),
				Statement:   "Restart returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Restart(ctx, drawnKey)
					})
				},
				Defect: runRestartIgnoresCancel,
			},
			{
				ID:          RunChecks.Toggle.Smoke(),
				Statement:   "Toggle returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RunHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().
							Toggle(ctx, string("key-a"), func() bool { return false })
					})
				},
				Defect: runTogglePanics,
			},
			{
				ID:          RunChecks.Toggle.SmokeGenerated(),
				Statement:   "Toggle returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					verboseResults := c.Draw(prop.List(prop.Of[bool]()), "verbose")
					drawnVerbose := func() func() bool {
						var verboseAt atomic.Int64
						return func() (r0 bool) {
							i := int(verboseAt.Add(1) - 1)
							if i >= len(verboseResults) {
								return r0
							}
							return verboseResults[i]
						}
					}
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Toggle(ctx, drawnKey, drawnVerbose())
					})
				},
				Defect: runTogglePanics,
			},
			{
				ID:          RunChecks.Toggle.Cancel(),
				Statement:   "Toggle returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RunHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().
							Toggle(ctx, string("key-a"), func() bool { return false })
					})
				},
				Defect: runToggleIgnoresCancel,
			},
			{
				ID:          RunChecks.Toggle.CancelGenerated(),
				Statement:   "Toggle returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					verboseResults := c.Draw(prop.List(prop.Of[bool]()), "verbose")
					drawnVerbose := func() func() bool {
						var verboseAt atomic.Int64
						return func() (r0 bool) {
							i := int(verboseAt.Add(1) - 1)
							if i >= len(verboseResults) {
								return r0
							}
							return verboseResults[i]
						}
					}
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Toggle(ctx, drawnKey, drawnVerbose())
					})
				},
				Defect: runToggleIgnoresCancel,
			},
			{
				ID:          RunChecks.Toggle.FollowsItsGraph(),
				Statement:   "each call of Toggle leaves the state that Status reads at the state before it or at one of its successors, over any sequence of calls",
				Claim:       spec.FollowsItsGraph,
				Probe:       spec.StatesMoveAlongTheGraph,
				Reason:      spec.ForbiddenTransition,
				Law:         spec.StatesFollowTheGraph,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *RunHandle, c *prop.Case) {
					probe.StatesMoveAlongTheGraph(
						h,
						c,
						runActions,
						"Toggle",
						func(ctx context.Context, args []any) conform.Output {
							values, err := runStatusCall(
								ctx,
								h.Subject(),
								[]any{args[0], bool(true)},
							)
							return conform.OutputOf(err, values...)
						},
						Graph,
					)
				},
				Defect: runToggleAdmitsAnyTransitionAfterThree,
			},
			{
				ID:          RunChecks.Status.Smoke(),
				Statement:   "Status returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RunHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Status(ctx, string("key-a"), bool(true))
						return err
					})
				},
				Defect: runStatusPanics,
			},
			{
				ID:          RunChecks.Status.SmokeGenerated(),
				Statement:   "Status returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnVerbose := c.Draw(prop.Of[bool](), "verbose")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Status(ctx, drawnKey, drawnVerbose)
						return err
					})
				},
				Defect: runStatusPanics,
			},
			{
				ID:          RunChecks.Status.Cancel(),
				Statement:   "Status returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RunHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Status(ctx, string("key-a"), bool(true))
						return err
					})
				},
				Defect: runStatusIgnoresCancel,
			},
			{
				ID:          RunChecks.Status.CancelGenerated(),
				Statement:   "Status returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RunHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnVerbose := c.Draw(prop.Of[bool](), "verbose")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Status(ctx, drawnKey, drawnVerbose)
						return err
					})
				},
				Defect: runStatusIgnoresCancel,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.FollowsItsGraph,
				Callable: "Run.Restart",
				Gate:     spec.GraphOfStates,
				Why:      "the variable that transitions names for Restart is no graph of the states that observe reads",
				Remedy:   "declare transitions as a map from each state to the list of its successors",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Run.Advance",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Run.Retry",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Run.Restart",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:composite-writer",
				Callable:       "Run.Toggle",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Run]{
			{
				Calls:    runAdvanceCalls,
				Over:     newRunOver,
				Callable: RunChecks.Advance.Callable,
			},
			{
				Calls:    runRetryCalls,
				Over:     newRunOver,
				Callable: RunChecks.Retry.Callable,
			},
			{
				Calls:    runRestartCalls,
				Over:     newRunOver,
				Callable: RunChecks.Restart.Callable,
			},
			{
				Calls:    runToggleCalls,
				Over:     newRunOver,
				Callable: RunChecks.Toggle.Callable,
			},
			{
				Calls:    runStatusCalls,
				Over:     newRunOver,
				Callable: RunChecks.Status.Callable,
			},
		},
	}
}

// RunRunSuite applies the checks of Run, and the checks that opts add, to
// a fresh subject of newSubject for each check.
func RunRunSuite[T Run](t *testing.T, newSubject func() T, opts ...RunOption) {
	t.Helper()
	RunRunSuiteFromHandle(t, func(*RunHandle) T { return newSubject() }, opts...)
}

// RunRunSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunRunSuiteFromHandle[T Run](
	t *testing.T,
	newSubject func(h *RunHandle) T,
	opts ...RunOption,
) {
	t.Helper()
	conform.Run(t, RunSuite(), runSubjects(newSubject), opts...)
}

// FuzzRunSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzRunSuite[T Run](
	f *testing.F,
	callable conform.Callable[Run],
	newSubject func(h *RunHandle) T,
	opts ...RunOption,
) {
	f.Helper()
	conform.Fuzz(f, RunSuite(), callable, runSubjects(newSubject), opts...)
}

// RunRunBudgets applies the allocation ceilings of Run and its checks that
// count allocations, and the ceilings and the checks that opts add, to a
// fresh subject of newSubject for each case.
func RunRunBudgets[T Run](t *testing.T, newSubject func(h *RunHandle) T, opts ...RunOption) {
	t.Helper()
	conform.RunBudgets(t, RunSuite(), runSubjects(newSubject), opts...)
}

// BenchRunSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchRunSuite[T Run](
	b conform.B,
	callable conform.Callable[Run],
	newSubject func(h *RunHandle) T,
	opts ...RunOption,
) {
	b.Helper()
	conform.Bench(b, RunSuite(), callable, runSubjects(newSubject), opts...)
}

// runSubjects returns the factory of the subjects of Run that calls
// newSubject with the handle of each check.
func runSubjects[T Run](newSubject func(h *RunHandle) T) func(h *RunHandle) Run {
	return func(h *RunHandle) Run { return newSubject(h) }
}

// runAdvanceCalls generates the calls of Run.Advance with the inputs that
// each case draws.
var runAdvanceCalls = prop.Composite(func(c *prop.Case) func(context.Context, Run) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Run) error {
		return subject.Advance(ctx, key)
	}
})

// runRetryCalls generates the calls of Run.Retry with the inputs that each
// case draws.
var runRetryCalls = prop.Composite(func(c *prop.Case) func(context.Context, Run) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Run) error {
		return subject.Retry(ctx, key)
	}
})

// runRestartCalls generates the calls of Run.Restart with the inputs that
// each case draws.
var runRestartCalls = prop.Composite(func(c *prop.Case) func(context.Context, Run) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Run) error {
		return subject.Restart(ctx, key)
	}
})

// runToggleCalls generates the calls of Run.Toggle with the inputs that
// each case draws.
var runToggleCalls = prop.Composite(func(c *prop.Case) func(context.Context, Run) error {
	key := c.Draw(prop.Of[string](), "key")
	verboseResults := c.Draw(prop.List(prop.Of[bool]()), "verbose")
	var verboseAt atomic.Int64
	verbose := func() (r0 bool) {
		i := int(verboseAt.Add(1) - 1)
		if i >= len(verboseResults) {
			return r0
		}
		return verboseResults[i]
	}
	return func(ctx context.Context, subject Run) error {
		return subject.Toggle(ctx, key, verbose)
	}
})

// runStatusCalls generates the calls of Run.Status with the inputs that
// each case draws.
var runStatusCalls = prop.Composite(func(c *prop.Case) func(context.Context, Run) error {
	key := c.Draw(prop.Of[string](), "key")
	verbose := c.Draw(prop.Of[bool](), "verbose")
	return func(ctx context.Context, subject Run) error {
		_, err := subject.Status(ctx, key, verbose)
		return err
	}
})

// runOverEmbedded is Run under the name of the field that runOver embeds.
type runOverEmbedded = Run

// runOver is a subject of Run whose callables call spend and return their
// zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type runOver struct {
	runOverEmbedded
	spend func()
}

// Advance calls spend and returns its zero results.
func (o runOver) Advance(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Retry calls spend and returns its zero results.
func (o runOver) Retry(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Restart calls spend and returns its zero results.
func (o runOver) Restart(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Toggle calls spend and returns its zero results.
func (o runOver) Toggle(context.Context, string, func() bool) (r0 error) {
	o.spend()
	return r0
}

// Status calls spend and returns its zero results.
func (o runOver) Status(context.Context, string, bool) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newRunOver returns the planted defect of a ceiling of Run, whose
// callables call spend.
func newRunOver(_ *RunHandle, spend func()) Run {
	return runOver{spend: spend}
}

// runAdvancePanicsDefect is the planted defect of a Run whose Advance
// panics.
type runAdvancePanicsDefect struct {
	runOverEmbedded
}

// Advance is the override of the planted defect.
func (*runAdvancePanicsDefect) Advance(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Run whose Advance panics")
}

// runAdvancePanics returns the planted defect of a Run whose Advance
// panics over a fresh subject.
func runAdvancePanics(h *RunHandle) Run {
	return &runAdvancePanicsDefect{runOverEmbedded: h.Fresh()}
}

// runAdvanceIgnoresCancelDefect is the planted defect of a Run whose
// Advance ignores the context.
type runAdvanceIgnoresCancelDefect struct {
	runOverEmbedded
}

// Advance is the override of the planted defect.
func (d *runAdvanceIgnoresCancelDefect) Advance(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.runOverEmbedded.Advance(p0, p1)
}

// runAdvanceIgnoresCancel returns the planted defect of a Run whose
// Advance ignores the context over a fresh subject.
func runAdvanceIgnoresCancel(h *RunHandle) Run {
	return &runAdvanceIgnoresCancelDefect{runOverEmbedded: h.Fresh()}
}

// runAdvanceAdmitsAnyTransitionAfterThreeDefect is the planted defect of a
// Run whose Advance makes its transition twice from its fourth call on.
type runAdvanceAdmitsAnyTransitionAfterThreeDefect struct {
	runOverEmbedded
	calls int
}

// Advance is the override of the planted defect.
func (d *runAdvanceAdmitsAnyTransitionAfterThreeDefect) Advance(
	p0 context.Context,
	p1 string,
) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.calls++
	if d.calls > 3 {
		r0 = d.runOverEmbedded.Advance(p0, p1)
	}
	return d.runOverEmbedded.Advance(p0, p1)
}

// runAdvanceAdmitsAnyTransitionAfterThree returns the planted defect of a
// Run whose Advance makes its transition twice from its fourth call on
// over a fresh subject.
func runAdvanceAdmitsAnyTransitionAfterThree(h *RunHandle) Run {
	return &runAdvanceAdmitsAnyTransitionAfterThreeDefect{runOverEmbedded: h.Fresh()}
}

// runRetryPanicsDefect is the planted defect of a Run whose Retry panics.
type runRetryPanicsDefect struct {
	runOverEmbedded
}

// Retry is the override of the planted defect.
func (*runRetryPanicsDefect) Retry(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Run whose Retry panics")
}

// runRetryPanics returns the planted defect of a Run whose Retry panics
// over a fresh subject.
func runRetryPanics(h *RunHandle) Run {
	return &runRetryPanicsDefect{runOverEmbedded: h.Fresh()}
}

// runRetryIgnoresCancelDefect is the planted defect of a Run whose Retry
// ignores the context.
type runRetryIgnoresCancelDefect struct {
	runOverEmbedded
}

// Retry is the override of the planted defect.
func (d *runRetryIgnoresCancelDefect) Retry(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.runOverEmbedded.Retry(p0, p1)
}

// runRetryIgnoresCancel returns the planted defect of a Run whose Retry
// ignores the context over a fresh subject.
func runRetryIgnoresCancel(h *RunHandle) Run {
	return &runRetryIgnoresCancelDefect{runOverEmbedded: h.Fresh()}
}

// runRetryAdmitsAnyTransitionAfterThreeDefect is the planted defect of a
// Run whose Retry makes its transition twice from its fourth call on.
type runRetryAdmitsAnyTransitionAfterThreeDefect struct {
	runOverEmbedded
	calls int
}

// Retry is the override of the planted defect.
func (d *runRetryAdmitsAnyTransitionAfterThreeDefect) Retry(
	p0 context.Context,
	p1 string,
) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.calls++
	if d.calls > 3 {
		r0 = d.runOverEmbedded.Retry(p0, p1)
	}
	return d.runOverEmbedded.Retry(p0, p1)
}

// runRetryAdmitsAnyTransitionAfterThree returns the planted defect of a
// Run whose Retry makes its transition twice from its fourth call on over
// a fresh subject.
func runRetryAdmitsAnyTransitionAfterThree(h *RunHandle) Run {
	return &runRetryAdmitsAnyTransitionAfterThreeDefect{runOverEmbedded: h.Fresh()}
}

// runRestartPanicsDefect is the planted defect of a Run whose Restart
// panics.
type runRestartPanicsDefect struct {
	runOverEmbedded
}

// Restart is the override of the planted defect.
func (*runRestartPanicsDefect) Restart(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Run whose Restart panics")
}

// runRestartPanics returns the planted defect of a Run whose Restart
// panics over a fresh subject.
func runRestartPanics(h *RunHandle) Run {
	return &runRestartPanicsDefect{runOverEmbedded: h.Fresh()}
}

// runRestartIgnoresCancelDefect is the planted defect of a Run whose
// Restart ignores the context.
type runRestartIgnoresCancelDefect struct {
	runOverEmbedded
}

// Restart is the override of the planted defect.
func (d *runRestartIgnoresCancelDefect) Restart(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.runOverEmbedded.Restart(p0, p1)
}

// runRestartIgnoresCancel returns the planted defect of a Run whose
// Restart ignores the context over a fresh subject.
func runRestartIgnoresCancel(h *RunHandle) Run {
	return &runRestartIgnoresCancelDefect{runOverEmbedded: h.Fresh()}
}

// runTogglePanicsDefect is the planted defect of a Run whose Toggle
// panics.
type runTogglePanicsDefect struct {
	runOverEmbedded
}

// Toggle is the override of the planted defect.
func (*runTogglePanicsDefect) Toggle(p0 context.Context, p1 string, p2 func() bool) (r0 error) {
	panic("the planted defect of a Run whose Toggle panics")
}

// runTogglePanics returns the planted defect of a Run whose Toggle panics
// over a fresh subject.
func runTogglePanics(h *RunHandle) Run {
	return &runTogglePanicsDefect{runOverEmbedded: h.Fresh()}
}

// runToggleIgnoresCancelDefect is the planted defect of a Run whose Toggle
// ignores the context.
type runToggleIgnoresCancelDefect struct {
	runOverEmbedded
}

// Toggle is the override of the planted defect.
func (d *runToggleIgnoresCancelDefect) Toggle(
	p0 context.Context,
	p1 string,
	p2 func() bool,
) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.runOverEmbedded.Toggle(p0, p1, p2)
}

// runToggleIgnoresCancel returns the planted defect of a Run whose Toggle
// ignores the context over a fresh subject.
func runToggleIgnoresCancel(h *RunHandle) Run {
	return &runToggleIgnoresCancelDefect{runOverEmbedded: h.Fresh()}
}

// runToggleAdmitsAnyTransitionAfterThreeDefect is the planted defect of a
// Run whose Toggle makes its transition twice from its fourth call on.
type runToggleAdmitsAnyTransitionAfterThreeDefect struct {
	runOverEmbedded
	calls int
}

// Toggle is the override of the planted defect.
func (d *runToggleAdmitsAnyTransitionAfterThreeDefect) Toggle(
	p0 context.Context,
	p1 string,
	p2 func() bool,
) (r0 error) {
	if err := p0.Err(); err != nil {
		r0 = err
		return r0
	}
	d.calls++
	if d.calls > 3 {
		r0 = d.runOverEmbedded.Toggle(p0, p1, p2)
	}
	return d.runOverEmbedded.Toggle(p0, p1, p2)
}

// runToggleAdmitsAnyTransitionAfterThree returns the planted defect of a
// Run whose Toggle makes its transition twice from its fourth call on over
// a fresh subject.
func runToggleAdmitsAnyTransitionAfterThree(h *RunHandle) Run {
	return &runToggleAdmitsAnyTransitionAfterThreeDefect{runOverEmbedded: h.Fresh()}
}

// runStatusPanicsDefect is the planted defect of a Run whose Status
// panics.
type runStatusPanicsDefect struct {
	runOverEmbedded
}

// Status is the override of the planted defect.
func (*runStatusPanicsDefect) Status(p0 context.Context, p1 string, p2 bool) (r0 string, r1 error) {
	panic("the planted defect of a Run whose Status panics")
}

// runStatusPanics returns the planted defect of a Run whose Status panics
// over a fresh subject.
func runStatusPanics(h *RunHandle) Run {
	return &runStatusPanicsDefect{runOverEmbedded: h.Fresh()}
}

// runStatusIgnoresCancelDefect is the planted defect of a Run whose Status
// ignores the context.
type runStatusIgnoresCancelDefect struct {
	runOverEmbedded
}

// Status is the override of the planted defect.
func (d *runStatusIgnoresCancelDefect) Status(
	p0 context.Context,
	p1 string,
	p2 bool,
) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.runOverEmbedded.Status(p0, p1, p2)
}

// runStatusIgnoresCancel returns the planted defect of a Run whose Status
// ignores the context over a fresh subject.
func runStatusIgnoresCancel(h *RunHandle) Run {
	return &runStatusIgnoresCancelDefect{runOverEmbedded: h.Fresh()}
}

// runActions returns the actions of the generated sequences of Run over
// the subject that subject returns. The actions of a case take their keys
// from the pools that it draws.
func runActions(c *prop.Case, subject func() Run) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Advance",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.Of[string](), "advance.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Advance", args)
				values, err := runAdvanceCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Retry",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.Of[string](), "retry.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Retry", args)
				values, err := runRetryCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Restart",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.Of[string](), "restart.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Restart", args)
				values, err := runRestartCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Toggle",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "toggle.key")
				verboseResults := c.Draw(prop.List(prop.Of[bool]()), "verbose")
				var verboseAt atomic.Int64
				verbose := func() (r0 bool) {
					i := int(verboseAt.Add(1) - 1)
					if i >= len(verboseResults) {
						return r0
					}
					return verboseResults[i]
				}
				return []any{key, verbose}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Toggle", args, args[0])
				values, err := runToggleCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Status",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "status.key")
				verbose := c.Draw(prop.Of[bool](), "status.verbose")
				return []any{key, verbose}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Status", args, args[0])
				values, err := runStatusCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// runAdvanceCall calls Advance on subject with the arguments args, and
// returns the values of the call beside its error.
func runAdvanceCall(ctx context.Context, subject Run, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Advance(ctx, a0)
}

// runRetryCall calls Retry on subject with the arguments args, and returns
// the values of the call beside its error.
func runRetryCall(ctx context.Context, subject Run, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Retry(ctx, a0)
}

// runRestartCall calls Restart on subject with the arguments args, and
// returns the values of the call beside its error.
func runRestartCall(ctx context.Context, subject Run, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Restart(ctx, a0)
}

// runToggleCall calls Toggle on subject with the arguments args, and
// returns the values of the call beside its error.
func runToggleCall(ctx context.Context, subject Run, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(func() bool)
	return nil, subject.Toggle(ctx, a0, a1)
}

// runStatusCall calls Status on subject with the arguments args, and
// returns the values of the call beside its error.
func runStatusCall(ctx context.Context, subject Run, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(bool)
	r0, err := subject.Status(ctx, a0, a1)
	return []any{r0}, err
}

// The derived checks of Log:
//
// Log.Append/smoke: Append returns without a panic for derived inputs.
//
// Log.Append/smoke#generated: Append returns without a panic for derived
// inputs.
//
// Log.Append/cancel: Append returns an error for a cancelled context.
//
// Log.Append/cancel#generated: Append returns an error for a cancelled
// context.
//
// Log.Append/zero-on-error: Append returns the zero value beside an error.
//
// Log.Append/zero-on-error#generated: Append returns the zero value beside
// an error.
//
// Log.Append/appender: Append returns a position above the position of the
// append before it.
//
// Log.Append/appender#generated: Append returns a position above the
// position of the append before it.
//
// Log.Append/appends-keep-ascending: Append returns a position above the
// one before it, over any sequence of calls.

// The derived checks of Log that witness their laws at the shallow end:
//
// Log.Append/smoke
//
// Log.Append/appender
//
// Log.Append/appender#generated

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Log writes.
type (
	LogHandle = conform.Handle[Log]
	LogCheck  = conform.Check[Log]
	LogOption = conform.Option[Log]
)

// LogChecks contains the name of each callable of Log, with the identities
// of its checks.
var LogChecks = struct {
	Append logAppend
}{
	Append: logAppend{"Log.Append"},
}

// logAppend is the name of Log.Append, with the identities of its checks.
type logAppend struct {
	conform.Callable[Log]
}

// Smoke returns the identity of the check that Append returns without a
// panic for derived inputs.
func (c logAppend) Smoke() conform.ID[Log] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Append returns
// without a panic for derived inputs.
func (c logAppend) SmokeGenerated() conform.ID[Log] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Append returns an error
// for a cancelled context.
func (c logAppend) Cancel() conform.ID[Log] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Append returns an
// error for a cancelled context.
func (c logAppend) CancelGenerated() conform.ID[Log] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Append returns the
// zero value beside an error.
func (c logAppend) ZeroOnError() conform.ID[Log] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Append
// returns the zero value beside an error.
func (c logAppend) ZeroOnErrorGenerated() conform.ID[Log] {
	return c.Check("zero-on-error#generated")
}

// Appender returns the identity of the check that Append returns a
// position above the position of the append before it.
func (c logAppend) Appender() conform.ID[Log] {
	return c.Check("appender")
}

// AppenderGenerated returns the identity of the check that Append returns
// a position above the position of the append before it.
func (c logAppend) AppenderGenerated() conform.ID[Log] {
	return c.Check("appender#generated")
}

// AppendsKeepAscending returns the identity of the check that Append
// returns a position above the one before it, over any sequence of calls.
func (c logAppend) AppendsKeepAscending() conform.ID[Log] {
	return c.Check("appends-keep-ascending")
}

// LogSuite returns the checks of Log, and the calls and the budgets of its
// callables.
func LogSuite() conform.Suite[Log] {
	return conform.Suite[Log]{
		Declaration: "Log",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Log]{
			{
				ID:          LogChecks.Append.Smoke(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LogHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Append(ctx, string("entry-a"))
						return err
					})
				},
				Defect: logAppendPanics,
			},
			{
				ID:          LogChecks.Append.SmokeGenerated(),
				Statement:   "Append returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LogHandle, c *prop.Case) {
					drawnEntry := c.Draw(prop.Of[string](), "entry")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Append(ctx, drawnEntry)
						return err
					})
				},
				Defect: logAppendPanics,
			},
			{
				ID:          LogChecks.Append.Cancel(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LogHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Append(ctx, string("entry-a"))
						return err
					})
				},
				Defect: logAppendIgnoresCancel,
			},
			{
				ID:          LogChecks.Append.CancelGenerated(),
				Statement:   "Append returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LogHandle, c *prop.Case) {
					drawnEntry := c.Draw(prop.Of[string](), "entry")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Append(ctx, drawnEntry)
						return err
					})
				},
				Defect: logAppendIgnoresCancel,
			},
			{
				ID:          LogChecks.Append.ZeroOnError(),
				Statement:   "Append returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *LogHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int64, error) {
						return h.Subject().Append(ctx, string("entry-a"))
					}, int64(0))
				},
				Defect: logAppendValueBeside,
			},
			{
				ID:          LogChecks.Append.ZeroOnErrorGenerated(),
				Statement:   "Append returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LogHandle, c *prop.Case) {
					drawnEntry := c.Draw(prop.Of[string](), "entry")
					probe.ZeroBesideError(h, func(ctx context.Context) (int64, error) {
						return h.Subject().Append(ctx, drawnEntry)
					}, int64(0))
				},
				Defect: logAppendValueBeside,
			},
			{
				ID:          LogChecks.Append.Appender(),
				Statement:   "Append returns a position above the position of the append before it",
				Claim:       spec.Appender,
				Probe:       spec.RepeatAscends,
				Reason:      spec.NotAscending,
				Law:         spec.AppendReturnsAboveTheLast,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LogHandle) {
					probe.RepeatAscends(h, func(ctx context.Context) (int64, error) {
						return h.Subject().Append(ctx, string("entry-a"))
					})
				},
				Defect: logAppendDescends,
			},
			{
				ID:          LogChecks.Append.AppenderGenerated(),
				Statement:   "Append returns a position above the position of the append before it",
				Claim:       spec.Appender,
				Probe:       spec.RepeatAscends,
				Reason:      spec.NotAscending,
				Law:         spec.AppendReturnsAboveTheLast,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *LogHandle, c *prop.Case) {
					drawnEntry := c.Draw(prop.Of[string](), "entry")
					probe.RepeatAscends(h, func(ctx context.Context) (int64, error) {
						return h.Subject().Append(ctx, drawnEntry)
					})
				},
				Defect: logAppendDescends,
			},
			{
				ID:          LogChecks.Append.AppendsKeepAscending(),
				Statement:   "Append returns a position above the one before it, over any sequence of calls",
				Claim:       spec.AppendsKeepAscending,
				Probe:       spec.PositionsRise,
				Reason:      spec.NotAscending,
				Law:         spec.AppendReturnsAboveTheLast,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Log]{LogChecks.Append.Appender()},
				Property: func(h *LogHandle, c *prop.Case) {
					probe.PositionsRise[Log, int64](h, c, logActions, "Append")
				},
				Defect: logAppendDescendsAfterThree,
			},
		},
		Calls: []conform.Call[Log]{
			{
				Calls:    logAppendCalls,
				Over:     newLogOver,
				Callable: LogChecks.Append.Callable,
			},
		},
	}
}

// RunLogSuite applies the checks of Log, and the checks that opts add, to
// a fresh subject of newSubject for each check.
func RunLogSuite[T Log](t *testing.T, newSubject func() T, opts ...LogOption) {
	t.Helper()
	RunLogSuiteFromHandle(t, func(*LogHandle) T { return newSubject() }, opts...)
}

// RunLogSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunLogSuiteFromHandle[T Log](
	t *testing.T,
	newSubject func(h *LogHandle) T,
	opts ...LogOption,
) {
	t.Helper()
	conform.Run(t, LogSuite(), logSubjects(newSubject), opts...)
}

// FuzzLogSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzLogSuite[T Log](
	f *testing.F,
	callable conform.Callable[Log],
	newSubject func(h *LogHandle) T,
	opts ...LogOption,
) {
	f.Helper()
	conform.Fuzz(f, LogSuite(), callable, logSubjects(newSubject), opts...)
}

// RunLogBudgets applies the allocation ceilings of Log and its checks that
// count allocations, and the ceilings and the checks that opts add, to a
// fresh subject of newSubject for each case.
func RunLogBudgets[T Log](t *testing.T, newSubject func(h *LogHandle) T, opts ...LogOption) {
	t.Helper()
	conform.RunBudgets(t, LogSuite(), logSubjects(newSubject), opts...)
}

// BenchLogSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchLogSuite[T Log](
	b conform.B,
	callable conform.Callable[Log],
	newSubject func(h *LogHandle) T,
	opts ...LogOption,
) {
	b.Helper()
	conform.Bench(b, LogSuite(), callable, logSubjects(newSubject), opts...)
}

// logSubjects returns the factory of the subjects of Log that calls
// newSubject with the handle of each check.
func logSubjects[T Log](newSubject func(h *LogHandle) T) func(h *LogHandle) Log {
	return func(h *LogHandle) Log { return newSubject(h) }
}

// logAppendCalls generates the calls of Log.Append with the inputs that
// each case draws.
var logAppendCalls = prop.Composite(func(c *prop.Case) func(context.Context, Log) error {
	entry := c.Draw(prop.Of[string](), "entry")
	return func(ctx context.Context, subject Log) error {
		_, err := subject.Append(ctx, entry)
		return err
	}
})

// logOverEmbedded is Log under the name of the field that logOver embeds.
type logOverEmbedded = Log

// logOver is a subject of Log whose callables call spend and return their
// zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type logOver struct {
	logOverEmbedded
	spend func()
}

// Append calls spend and returns its zero results.
func (o logOver) Append(context.Context, string) (r0 int64, r1 error) {
	o.spend()
	return r0, r1
}

// newLogOver returns the planted defect of a ceiling of Log, whose
// callables call spend.
func newLogOver(_ *LogHandle, spend func()) Log {
	return logOver{spend: spend}
}

// logAppendPanicsDefect is the planted defect of a Log whose Append
// panics.
type logAppendPanicsDefect struct {
	logOverEmbedded
}

// Append is the override of the planted defect.
func (*logAppendPanicsDefect) Append(p0 context.Context, p1 string) (r0 int64, r1 error) {
	panic("the planted defect of a Log whose Append panics")
}

// logAppendPanics returns the planted defect of a Log whose Append panics
// over a fresh subject.
func logAppendPanics(h *LogHandle) Log {
	return &logAppendPanicsDefect{logOverEmbedded: h.Fresh()}
}

// logAppendIgnoresCancelDefect is the planted defect of a Log whose Append
// ignores the context.
type logAppendIgnoresCancelDefect struct {
	logOverEmbedded
}

// Append is the override of the planted defect.
func (d *logAppendIgnoresCancelDefect) Append(p0 context.Context, p1 string) (r0 int64, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.logOverEmbedded.Append(p0, p1)
}

// logAppendIgnoresCancel returns the planted defect of a Log whose Append
// ignores the context over a fresh subject.
func logAppendIgnoresCancel(h *LogHandle) Log {
	return &logAppendIgnoresCancelDefect{logOverEmbedded: h.Fresh()}
}

// logAppendValueBesideDefect is the planted defect of a Log whose Append
// returns a value beside an error.
type logAppendValueBesideDefect struct {
	logOverEmbedded
}

// Append is the override of the planted defect.
func (d *logAppendValueBesideDefect) Append(p0 context.Context, p1 string) (r0 int64, r1 error) {
	r0, r1 = d.logOverEmbedded.Append(p0, p1)
	if r1 != nil {
		r0 = int64(7)
	}
	return r0, r1
}

// logAppendValueBeside returns the planted defect of a Log whose Append
// returns a value beside an error over a fresh subject.
func logAppendValueBeside(h *LogHandle) Log {
	return &logAppendValueBesideDefect{logOverEmbedded: h.Fresh()}
}

// logAppendDescendsDefect is the planted defect of a Log whose Append
// returns a position that is not above the last.
type logAppendDescendsDefect struct {
	logOverEmbedded
	calls int
}

// Append is the override of the planted defect.
func (d *logAppendDescendsDefect) Append(p0 context.Context, p1 string) (r0 int64, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	if d.calls == 1 {
		r0 = int64(7)
	}
	return r0, r1
}

// logAppendDescends returns the planted defect of a Log whose Append
// returns a position that is not above the last over a fresh subject.
func logAppendDescends(h *LogHandle) Log {
	return &logAppendDescendsDefect{logOverEmbedded: h.Fresh()}
}

// logAppendDescendsAfterThreeDefect is the planted defect of a Log whose
// Append returns the position 0 from its fourth call on.
type logAppendDescendsAfterThreeDefect struct {
	logOverEmbedded
	calls int
}

// Append is the override of the planted defect.
func (d *logAppendDescendsAfterThreeDefect) Append(
	p0 context.Context,
	p1 string,
) (r0 int64, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	if d.calls > 3 {
		return r0, r1
	}
	return d.logOverEmbedded.Append(p0, p1)
}

// logAppendDescendsAfterThree returns the planted defect of a Log whose
// Append returns the position 0 from its fourth call on over a fresh
// subject.
func logAppendDescendsAfterThree(h *LogHandle) Log {
	return &logAppendDescendsAfterThreeDefect{logOverEmbedded: h.Fresh()}
}

// logActions returns the actions of the generated sequences of Log over
// the subject that subject returns. The actions of a case take their keys
// from the pools that it draws.
func logActions(c *prop.Case, subject func() Log) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Append",
			Input: func(c *prop.Case, _ struct{}) any {
				entry := c.Draw(prop.SampledFrom(keysString...), "append.entry")
				return []any{entry}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Append", args, args[0])
				values, err := logAppendCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// logAppendCall calls Append on subject with the arguments args, and
// returns the values of the call beside its error.
func logAppendCall(ctx context.Context, subject Log, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Append(ctx, a0)
	return []any{r0}, err
}

// The derived checks of Feed:
//
// Feed.Put/smoke: Put returns without a panic for derived inputs.
//
// Feed.Put/smoke#generated: Put returns without a panic for derived
// inputs.
//
// Feed.Put/cancel: Put returns an error for a cancelled context.
//
// Feed.Put/cancel#generated: Put returns an error for a cancelled context.
//
// Feed.All/smoke: All returns without a panic for derived inputs.
//
// Feed.All/no-duplicates: no element appears twice in one All.
//
// Feed.All/no-duplicates#generated: no element appears twice in one All.
//
// Feed.All/streams-what-was-stored: All yields each value that Put stored.
//
// Feed.All/streams-what-was-stored#generated: All yields each value that
// Put stored.
//
// Feed.All/yields-each-once: All never yields an element twice, over any
// sequence of calls.
//
// Feed.Each/smoke: Each returns without a panic for derived inputs.
//
// Feed.Each/no-duplicates: no element appears twice in one Each.
//
// Feed.Each/no-duplicates#generated: no element appears twice in one Each.
//
// Feed.Each/streams-what-was-stored: Each yields each value that Put
// stored.
//
// Feed.Each/streams-what-was-stored#generated: Each yields each value that
// Put stored.
//
// Feed.Each/yields-each-once: Each never yields an element twice, over any
// sequence of calls.
//
// Feed.Pairs/smoke: Pairs returns without a panic for derived inputs.
//
// Feed.Pairs/no-duplicates: no element appears twice in one Pairs.
//
// Feed.Pairs/no-duplicates#generated: no element appears twice in one
// Pairs.
//
// Feed.Pairs/streams-what-was-stored: Pairs yields each value that Put
// stored.
//
// Feed.Pairs/streams-what-was-stored#generated: Pairs yields each value
// that Put stored.
//
// Feed.Pairs/yields-each-once: Pairs never yields an element twice, over
// any sequence of calls.

// The derived checks of Feed that witness their laws at the shallow end:
//
// Feed.Put/smoke
//
// Feed.All/smoke
//
// Feed.All/no-duplicates
//
// Feed.All/no-duplicates#generated
//
// Feed.All/streams-what-was-stored
//
// Feed.All/streams-what-was-stored#generated
//
// Feed.Each/smoke
//
// Feed.Each/no-duplicates
//
// Feed.Each/no-duplicates#generated
//
// Feed.Each/streams-what-was-stored
//
// Feed.Each/streams-what-was-stored#generated
//
// Feed.Pairs/smoke
//
// Feed.Pairs/no-duplicates
//
// Feed.Pairs/no-duplicates#generated
//
// Feed.Pairs/streams-what-was-stored
//
// Feed.Pairs/streams-what-was-stored#generated

// The refused claims of Feed:
//
// Feed.All/cancel: All cannot fail. The remedy is to give All an error
// result.
//
// Feed.Each/cancel: Each cannot fail. The remedy is to give Each an error
// result.
//
// Feed.Pairs/cancel: Pairs cannot fail. The remedy is to give Pairs an
// error result.

// The unclaimed classifications of Feed:
//
// The shape:writer of Feed.Put has the stance silent. A write is an input
// of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Feed writes.
type (
	FeedHandle = conform.Handle[Feed]
	FeedCheck  = conform.Check[Feed]
	FeedOption = conform.Option[Feed]
)

// FeedChecks contains the name of each callable of Feed, with the
// identities of its checks.
var FeedChecks = struct {
	Put   feedPut
	All   feedAll
	Each  feedEach
	Pairs feedPairs
}{
	Put:   feedPut{"Feed.Put"},
	All:   feedAll{"Feed.All"},
	Each:  feedEach{"Feed.Each"},
	Pairs: feedPairs{"Feed.Pairs"},
}

// feedPut is the name of Feed.Put, with the identities of its checks.
type feedPut struct {
	conform.Callable[Feed]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c feedPut) Smoke() conform.ID[Feed] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Put returns
// without a panic for derived inputs.
func (c feedPut) SmokeGenerated() conform.ID[Feed] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c feedPut) Cancel() conform.ID[Feed] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Put returns an
// error for a cancelled context.
func (c feedPut) CancelGenerated() conform.ID[Feed] {
	return c.Check("cancel#generated")
}

// feedAll is the name of Feed.All, with the identities of its checks.
type feedAll struct {
	conform.Callable[Feed]
}

// Smoke returns the identity of the check that All returns without a panic
// for derived inputs.
func (c feedAll) Smoke() conform.ID[Feed] {
	return c.Check("smoke")
}

// NoDuplicates returns the identity of the check that no element appears
// twice in one All.
func (c feedAll) NoDuplicates() conform.ID[Feed] {
	return c.Check("no-duplicates")
}

// NoDuplicatesGenerated returns the identity of the check that no element
// appears twice in one All.
func (c feedAll) NoDuplicatesGenerated() conform.ID[Feed] {
	return c.Check("no-duplicates#generated")
}

// StreamsWhatWasStored returns the identity of the check that All yields
// each value that Put stored.
func (c feedAll) StreamsWhatWasStored() conform.ID[Feed] {
	return c.Check("streams-what-was-stored")
}

// StreamsWhatWasStoredGenerated returns the identity of the check that All
// yields each value that Put stored.
func (c feedAll) StreamsWhatWasStoredGenerated() conform.ID[Feed] {
	return c.Check("streams-what-was-stored#generated")
}

// YieldsEachOnce returns the identity of the check that All never yields
// an element twice, over any sequence of calls.
func (c feedAll) YieldsEachOnce() conform.ID[Feed] {
	return c.Check("yields-each-once")
}

// feedEach is the name of Feed.Each, with the identities of its checks.
type feedEach struct {
	conform.Callable[Feed]
}

// Smoke returns the identity of the check that Each returns without a
// panic for derived inputs.
func (c feedEach) Smoke() conform.ID[Feed] {
	return c.Check("smoke")
}

// NoDuplicates returns the identity of the check that no element appears
// twice in one Each.
func (c feedEach) NoDuplicates() conform.ID[Feed] {
	return c.Check("no-duplicates")
}

// NoDuplicatesGenerated returns the identity of the check that no element
// appears twice in one Each.
func (c feedEach) NoDuplicatesGenerated() conform.ID[Feed] {
	return c.Check("no-duplicates#generated")
}

// StreamsWhatWasStored returns the identity of the check that Each yields
// each value that Put stored.
func (c feedEach) StreamsWhatWasStored() conform.ID[Feed] {
	return c.Check("streams-what-was-stored")
}

// StreamsWhatWasStoredGenerated returns the identity of the check that
// Each yields each value that Put stored.
func (c feedEach) StreamsWhatWasStoredGenerated() conform.ID[Feed] {
	return c.Check("streams-what-was-stored#generated")
}

// YieldsEachOnce returns the identity of the check that Each never yields
// an element twice, over any sequence of calls.
func (c feedEach) YieldsEachOnce() conform.ID[Feed] {
	return c.Check("yields-each-once")
}

// feedPairs is the name of Feed.Pairs, with the identities of its checks.
type feedPairs struct {
	conform.Callable[Feed]
}

// Smoke returns the identity of the check that Pairs returns without a
// panic for derived inputs.
func (c feedPairs) Smoke() conform.ID[Feed] {
	return c.Check("smoke")
}

// NoDuplicates returns the identity of the check that no element appears
// twice in one Pairs.
func (c feedPairs) NoDuplicates() conform.ID[Feed] {
	return c.Check("no-duplicates")
}

// NoDuplicatesGenerated returns the identity of the check that no element
// appears twice in one Pairs.
func (c feedPairs) NoDuplicatesGenerated() conform.ID[Feed] {
	return c.Check("no-duplicates#generated")
}

// StreamsWhatWasStored returns the identity of the check that Pairs yields
// each value that Put stored.
func (c feedPairs) StreamsWhatWasStored() conform.ID[Feed] {
	return c.Check("streams-what-was-stored")
}

// StreamsWhatWasStoredGenerated returns the identity of the check that
// Pairs yields each value that Put stored.
func (c feedPairs) StreamsWhatWasStoredGenerated() conform.ID[Feed] {
	return c.Check("streams-what-was-stored#generated")
}

// YieldsEachOnce returns the identity of the check that Pairs never yields
// an element twice, over any sequence of calls.
func (c feedPairs) YieldsEachOnce() conform.ID[Feed] {
	return c.Check("yields-each-once")
}

// FeedSuite returns the checks of Feed, and the calls and the budgets of
// its callables.
func FeedSuite() conform.Suite[Feed] {
	return conform.Suite[Feed]{
		Declaration: "Feed",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Feed]{
			{
				ID:          FeedChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("item-a"))
					})
				},
				Defect: feedPutPanics,
			},
			{
				ID:          FeedChecks.Put.SmokeGenerated(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnItem)
					})
				},
				Defect: feedPutPanics,
			},
			{
				ID:          FeedChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *FeedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, string("item-a"))
					})
				},
				Defect: feedPutIgnoresCancel,
			},
			{
				ID:          FeedChecks.Put.CancelGenerated(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(prop.Of[string](), "item")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnItem)
					})
				},
				Defect: feedPutIgnoresCancel,
			},
			{
				ID:          FeedChecks.All.Smoke(),
				Statement:   "All returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().All(ctx)
						return nil
					})
				},
				Defect: feedAllPanics,
			},
			{
				ID:          FeedChecks.All.NoDuplicates(),
				Statement:   "no element appears twice in one All",
				Claim:       spec.NoDuplicates,
				Probe:       spec.NoElementTwice,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.NoElementTwice(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().All(ctx)
						var drained []string
						for e, err := range r0 {
							if err != nil {
								return drained, err
							}
							drained = append(drained, e)
						}
						return drained, nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := string("item-a")
							if calls > 1 {
								w = string("item-b")
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedAllDrainsTwice,
			},
			{
				ID:          FeedChecks.All.NoDuplicatesGenerated(),
				Statement:   "no element appears twice in one All",
				Claim:       spec.NoDuplicates,
				Probe:       spec.NoElementTwice,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.NoElementTwice(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().All(ctx)
						var drained []string
						for e, err := range r0 {
							if err != nil {
								return drained, err
							}
							drained = append(drained, e)
						}
						return drained, nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnItem[0]
							if calls > 1 {
								w = drawnItem[1]
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedAllDrainsTwice,
			},
			{
				ID:          FeedChecks.All.StreamsWhatWasStored(),
				Statement:   "All yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().All(ctx)
						var drained []string
						for e, err := range r0 {
							if err != nil {
								return drained, err
							}
							drained = append(drained, e)
						}
						return drained, nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := string("item-a")
							if calls > 1 {
								w = string("item-b")
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedAllDrainsNothing,
			},
			{
				ID:          FeedChecks.All.StreamsWhatWasStoredGenerated(),
				Statement:   "All yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().All(ctx)
						var drained []string
						for e, err := range r0 {
							if err != nil {
								return drained, err
							}
							drained = append(drained, e)
						}
						return drained, nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnItem[0]
							if calls > 1 {
								w = drawnItem[1]
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedAllDrainsNothing,
			},
			{
				ID:          FeedChecks.All.YieldsEachOnce(),
				Statement:   "All never yields an element twice, over any sequence of calls",
				Claim:       spec.YieldsEachOnce,
				Probe:       spec.DrainsWithoutRepeats,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Feed]{FeedChecks.All.NoDuplicates()},
				Property: func(h *FeedHandle, c *prop.Case) {
					probe.DrainsWithoutRepeats[Feed, string](h, c, feedActions, "All")
				},
				Defect: feedAllRepeatsAfterThree,
			},
			{
				ID:          FeedChecks.Each.Smoke(),
				Statement:   "Each returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Each(ctx)
						return nil
					})
				},
				Defect: feedEachPanics,
			},
			{
				ID:          FeedChecks.Each.NoDuplicates(),
				Statement:   "no element appears twice in one Each",
				Claim:       spec.NoDuplicates,
				Probe:       spec.NoElementTwice,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.NoElementTwice(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().Each(ctx)
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := string("item-a")
							if calls > 1 {
								w = string("item-b")
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedEachDrainsTwice,
			},
			{
				ID:          FeedChecks.Each.NoDuplicatesGenerated(),
				Statement:   "no element appears twice in one Each",
				Claim:       spec.NoDuplicates,
				Probe:       spec.NoElementTwice,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.NoElementTwice(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().Each(ctx)
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnItem[0]
							if calls > 1 {
								w = drawnItem[1]
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedEachDrainsTwice,
			},
			{
				ID:          FeedChecks.Each.StreamsWhatWasStored(),
				Statement:   "Each yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().Each(ctx)
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := string("item-a")
							if calls > 1 {
								w = string("item-b")
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedEachDrainsNothing,
			},
			{
				ID:          FeedChecks.Each.StreamsWhatWasStoredGenerated(),
				Statement:   "Each yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.KeepsWhatWentIn(h, func(ctx context.Context) ([]string, error) {
						r0 := h.Subject().Each(ctx)
						return slices.Collect(r0), nil
					}, func() func(ctx context.Context) error {
						calls := 0
						return func(ctx context.Context) error {
							calls++
							w := drawnItem[0]
							if calls > 1 {
								w = drawnItem[1]
							}
							return h.Subject().Put(ctx, w)
						}
					}())
				},
				Defect: feedEachDrainsNothing,
			},
			{
				ID:          FeedChecks.Each.YieldsEachOnce(),
				Statement:   "Each never yields an element twice, over any sequence of calls",
				Claim:       spec.YieldsEachOnce,
				Probe:       spec.DrainsWithoutRepeats,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Feed]{FeedChecks.Each.NoDuplicates()},
				Property: func(h *FeedHandle, c *prop.Case) {
					probe.DrainsWithoutRepeats[Feed, string](h, c, feedActions, "Each")
				},
				Defect: feedEachRepeatsAfterThree,
			},
			{
				ID:          FeedChecks.Pairs.Smoke(),
				Statement:   "Pairs returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Pairs(ctx)
						return nil
					})
				},
				Defect: feedPairsPanics,
			},
			{
				ID:          FeedChecks.Pairs.NoDuplicates(),
				Statement:   "no element appears twice in one Pairs",
				Claim:       spec.NoDuplicates,
				Probe:       spec.NoElementTwice,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.NoElementTwice(
						h,
						func(ctx context.Context) ([]probe.Pair[string, int], error) {
							r0 := h.Subject().Pairs(ctx)
							var drained []probe.Pair[string, int]
							for k, v := range r0 {
								drained = append(
									drained,
									probe.Pair[string, int]{First: k, Second: v},
								)
							}
							return drained, nil
						},
						func() func(ctx context.Context) error {
							calls := 0
							return func(ctx context.Context) error {
								calls++
								w := string("item-a")
								if calls > 1 {
									w = string("item-b")
								}
								return h.Subject().Put(ctx, w)
							}
						}(),
					)
				},
				Defect: feedPairsDrainsTwice,
			},
			{
				ID:          FeedChecks.Pairs.NoDuplicatesGenerated(),
				Statement:   "no element appears twice in one Pairs",
				Claim:       spec.NoDuplicates,
				Probe:       spec.NoElementTwice,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.NoElementTwice(
						h,
						func(ctx context.Context) ([]probe.Pair[string, int], error) {
							r0 := h.Subject().Pairs(ctx)
							var drained []probe.Pair[string, int]
							for k, v := range r0 {
								drained = append(
									drained,
									probe.Pair[string, int]{First: k, Second: v},
								)
							}
							return drained, nil
						},
						func() func(ctx context.Context) error {
							calls := 0
							return func(ctx context.Context) error {
								calls++
								w := drawnItem[0]
								if calls > 1 {
									w = drawnItem[1]
								}
								return h.Subject().Put(ctx, w)
							}
						}(),
					)
				},
				Defect: feedPairsDrainsTwice,
			},
			{
				ID:          FeedChecks.Pairs.StreamsWhatWasStored(),
				Statement:   "Pairs yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *FeedHandle) {
					probe.KeepsWhatWentIn(
						h,
						func(ctx context.Context) ([]probe.Pair[string, int], error) {
							r0 := h.Subject().Pairs(ctx)
							var drained []probe.Pair[string, int]
							for k, v := range r0 {
								drained = append(
									drained,
									probe.Pair[string, int]{First: k, Second: v},
								)
							}
							return drained, nil
						},
						func() func(ctx context.Context) error {
							calls := 0
							return func(ctx context.Context) error {
								calls++
								w := string("item-a")
								if calls > 1 {
									w = string("item-b")
								}
								return h.Subject().Put(ctx, w)
							}
						}(),
					)
				},
				Defect: feedPairsDrainsNothing,
			},
			{
				ID:          FeedChecks.Pairs.StreamsWhatWasStoredGenerated(),
				Statement:   "Pairs yields each value that Put stored",
				Claim:       spec.StreamsWhatWasStored,
				Probe:       spec.KeepsWhatWentIn,
				Reason:      spec.ShrunkDrain,
				Law:         spec.DrainKeepsWhatWentIn,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *FeedHandle, c *prop.Case) {
					drawnItem := c.Draw(
						prop.List(
							prop.Of[string](),
							prop.Unique(),
							prop.MinSize(2),
							prop.MaxSize(2),
						),
						"item",
					)
					probe.KeepsWhatWentIn(
						h,
						func(ctx context.Context) ([]probe.Pair[string, int], error) {
							r0 := h.Subject().Pairs(ctx)
							var drained []probe.Pair[string, int]
							for k, v := range r0 {
								drained = append(
									drained,
									probe.Pair[string, int]{First: k, Second: v},
								)
							}
							return drained, nil
						},
						func() func(ctx context.Context) error {
							calls := 0
							return func(ctx context.Context) error {
								calls++
								w := drawnItem[0]
								if calls > 1 {
									w = drawnItem[1]
								}
								return h.Subject().Put(ctx, w)
							}
						}(),
					)
				},
				Defect: feedPairsDrainsNothing,
			},
			{
				ID:          FeedChecks.Pairs.YieldsEachOnce(),
				Statement:   "Pairs never yields an element twice, over any sequence of calls",
				Claim:       spec.YieldsEachOnce,
				Probe:       spec.DrainsWithoutRepeats,
				Reason:      spec.RepeatedElement,
				Law:         spec.NoElementIsDrainedTwice,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Feed]{FeedChecks.Pairs.NoDuplicates()},
				Property: func(h *FeedHandle, c *prop.Case) {
					probe.DrainsWithoutRepeats[Feed, probe.Pair[string, int]](
						h,
						c,
						feedActions,
						"Pairs",
					)
				},
				Defect: feedPairsRepeatsAfterThree,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cancel,
				Callable: "Feed.All",
				Gate:     spec.Failure,
				Why:      "All cannot fail",
				Remedy:   "give All an error result",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Feed.Each",
				Gate:     spec.Failure,
				Why:      "Each cannot fail",
				Remedy:   "give Each an error result",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Feed.Pairs",
				Gate:     spec.Failure,
				Why:      "Pairs cannot fail",
				Remedy:   "give Pairs an error result",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Feed.Put",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Feed]{
			{
				Calls:    feedPutCalls,
				Over:     newFeedOver,
				Callable: FeedChecks.Put.Callable,
			},
			{
				Calls:    feedAllCalls,
				Over:     newFeedOver,
				Callable: FeedChecks.All.Callable,
			},
			{
				Calls:    feedEachCalls,
				Over:     newFeedOver,
				Callable: FeedChecks.Each.Callable,
			},
			{
				Calls:    feedPairsCalls,
				Over:     newFeedOver,
				Callable: FeedChecks.Pairs.Callable,
			},
		},
	}
}

// RunFeedSuite applies the checks of Feed, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunFeedSuite[T Feed](t *testing.T, newSubject func() T, opts ...FeedOption) {
	t.Helper()
	RunFeedSuiteFromHandle(t, func(*FeedHandle) T { return newSubject() }, opts...)
}

// RunFeedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunFeedSuiteFromHandle[T Feed](
	t *testing.T,
	newSubject func(h *FeedHandle) T,
	opts ...FeedOption,
) {
	t.Helper()
	conform.Run(t, FeedSuite(), feedSubjects(newSubject), opts...)
}

// FuzzFeedSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzFeedSuite[T Feed](
	f *testing.F,
	callable conform.Callable[Feed],
	newSubject func(h *FeedHandle) T,
	opts ...FeedOption,
) {
	f.Helper()
	conform.Fuzz(f, FeedSuite(), callable, feedSubjects(newSubject), opts...)
}

// RunFeedBudgets applies the allocation ceilings of Feed and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunFeedBudgets[T Feed](t *testing.T, newSubject func(h *FeedHandle) T, opts ...FeedOption) {
	t.Helper()
	conform.RunBudgets(t, FeedSuite(), feedSubjects(newSubject), opts...)
}

// BenchFeedSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchFeedSuite[T Feed](
	b conform.B,
	callable conform.Callable[Feed],
	newSubject func(h *FeedHandle) T,
	opts ...FeedOption,
) {
	b.Helper()
	conform.Bench(b, FeedSuite(), callable, feedSubjects(newSubject), opts...)
}

// feedSubjects returns the factory of the subjects of Feed that calls
// newSubject with the handle of each check.
func feedSubjects[T Feed](newSubject func(h *FeedHandle) T) func(h *FeedHandle) Feed {
	return func(h *FeedHandle) Feed { return newSubject(h) }
}

// feedPutCalls generates the calls of Feed.Put with the inputs that each
// case draws.
var feedPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Feed) error {
	item := c.Draw(prop.Of[string](), "item")
	return func(ctx context.Context, subject Feed) error {
		return subject.Put(ctx, item)
	}
})

// feedAllCalls generates the calls of Feed.All with the inputs that each
// case draws.
var feedAllCalls = prop.Composite(func(*prop.Case) func(context.Context, Feed) error {
	return func(ctx context.Context, subject Feed) error {
		_ = subject.All(ctx)
		return nil
	}
})

// feedEachCalls generates the calls of Feed.Each with the inputs that each
// case draws.
var feedEachCalls = prop.Composite(func(*prop.Case) func(context.Context, Feed) error {
	return func(ctx context.Context, subject Feed) error {
		_ = subject.Each(ctx)
		return nil
	}
})

// feedPairsCalls generates the calls of Feed.Pairs with the inputs that
// each case draws.
var feedPairsCalls = prop.Composite(func(*prop.Case) func(context.Context, Feed) error {
	return func(ctx context.Context, subject Feed) error {
		_ = subject.Pairs(ctx)
		return nil
	}
})

// feedOverEmbedded is Feed under the name of the field that feedOver
// embeds.
type feedOverEmbedded = Feed

// feedOver is a subject of Feed whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type feedOver struct {
	feedOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o feedOver) Put(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// All calls spend and returns its zero results.
func (o feedOver) All(context.Context) (r0 iter.Seq2[string, error]) {
	o.spend()
	return r0
}

// Each calls spend and returns its zero results.
func (o feedOver) Each(context.Context) (r0 iter.Seq[string]) {
	o.spend()
	return r0
}

// Pairs calls spend and returns its zero results.
func (o feedOver) Pairs(context.Context) (r0 iter.Seq2[string, int]) {
	o.spend()
	return r0
}

// newFeedOver returns the planted defect of a ceiling of Feed, whose
// callables call spend.
func newFeedOver(_ *FeedHandle, spend func()) Feed {
	return feedOver{spend: spend}
}

// feedPutPanicsDefect is the planted defect of a Feed whose Put panics.
type feedPutPanicsDefect struct {
	feedOverEmbedded
}

// Put is the override of the planted defect.
func (*feedPutPanicsDefect) Put(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Feed whose Put panics")
}

// feedPutPanics returns the planted defect of a Feed whose Put panics over
// a fresh subject.
func feedPutPanics(h *FeedHandle) Feed {
	return &feedPutPanicsDefect{feedOverEmbedded: h.Fresh()}
}

// feedPutIgnoresCancelDefect is the planted defect of a Feed whose Put
// ignores the context.
type feedPutIgnoresCancelDefect struct {
	feedOverEmbedded
}

// Put is the override of the planted defect.
func (d *feedPutIgnoresCancelDefect) Put(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.feedOverEmbedded.Put(p0, p1)
}

// feedPutIgnoresCancel returns the planted defect of a Feed whose Put
// ignores the context over a fresh subject.
func feedPutIgnoresCancel(h *FeedHandle) Feed {
	return &feedPutIgnoresCancelDefect{feedOverEmbedded: h.Fresh()}
}

// feedAllPanicsDefect is the planted defect of a Feed whose All panics.
type feedAllPanicsDefect struct {
	feedOverEmbedded
}

// All is the override of the planted defect.
func (*feedAllPanicsDefect) All(p0 context.Context) (r0 iter.Seq2[string, error]) {
	panic("the planted defect of a Feed whose All panics")
}

// feedAllPanics returns the planted defect of a Feed whose All panics over
// a fresh subject.
func feedAllPanics(h *FeedHandle) Feed {
	return &feedAllPanicsDefect{feedOverEmbedded: h.Fresh()}
}

// feedAllDrainsTwiceDefect is the planted defect of a Feed whose All
// returns one element twice.
type feedAllDrainsTwiceDefect struct {
	feedOverEmbedded
}

// All is the override of the planted defect.
func (d *feedAllDrainsTwiceDefect) All(p0 context.Context) (r0 iter.Seq2[string, error]) {
	r0 = d.feedOverEmbedded.All(p0)
	if len(r0) > 0 {
		r0 = append(r0[:1:1], r0...)
	}
	return r0
}

// feedAllDrainsTwice returns the planted defect of a Feed whose All
// returns one element twice over a fresh subject.
func feedAllDrainsTwice(h *FeedHandle) Feed {
	return &feedAllDrainsTwiceDefect{feedOverEmbedded: h.Fresh()}
}

// feedAllDrainsNothingDefect is the planted defect of a Feed whose All
// loses what Put stored.
type feedAllDrainsNothingDefect struct {
	feedOverEmbedded
}

// All is the override of the planted defect.
func (*feedAllDrainsNothingDefect) All(p0 context.Context) (r0 iter.Seq2[string, error]) {
	r0 = func(func(string, error) bool) {}
	return r0
}

// feedAllDrainsNothing returns the planted defect of a Feed whose All
// loses what Put stored over a fresh subject.
func feedAllDrainsNothing(h *FeedHandle) Feed {
	return &feedAllDrainsNothingDefect{feedOverEmbedded: h.Fresh()}
}

// feedAllRepeatsAfterThreeDefect is the planted defect of a Feed whose All
// yields its first element twice from its fourth call on.
type feedAllRepeatsAfterThreeDefect struct {
	feedOverEmbedded
	calls int
}

// All is the override of the planted defect.
func (d *feedAllRepeatsAfterThreeDefect) All(p0 context.Context) (r0 iter.Seq2[string, error]) {
	r0 = d.feedOverEmbedded.All(p0)
	d.calls++
	if d.calls > 3 && len(r0) > 0 {
		r0 = append(r0[:1:1], r0...)
	}
	return r0
}

// feedAllRepeatsAfterThree returns the planted defect of a Feed whose All
// yields its first element twice from its fourth call on over a fresh
// subject.
func feedAllRepeatsAfterThree(h *FeedHandle) Feed {
	return &feedAllRepeatsAfterThreeDefect{feedOverEmbedded: h.Fresh()}
}

// feedEachPanicsDefect is the planted defect of a Feed whose Each panics.
type feedEachPanicsDefect struct {
	feedOverEmbedded
}

// Each is the override of the planted defect.
func (*feedEachPanicsDefect) Each(p0 context.Context) (r0 iter.Seq[string]) {
	panic("the planted defect of a Feed whose Each panics")
}

// feedEachPanics returns the planted defect of a Feed whose Each panics
// over a fresh subject.
func feedEachPanics(h *FeedHandle) Feed {
	return &feedEachPanicsDefect{feedOverEmbedded: h.Fresh()}
}

// feedEachDrainsTwiceDefect is the planted defect of a Feed whose Each
// returns one element twice.
type feedEachDrainsTwiceDefect struct {
	feedOverEmbedded
}

// Each is the override of the planted defect.
func (d *feedEachDrainsTwiceDefect) Each(p0 context.Context) (r0 iter.Seq[string]) {
	r0 = d.feedOverEmbedded.Each(p0)
	if len(r0) > 0 {
		r0 = append(r0[:1:1], r0...)
	}
	return r0
}

// feedEachDrainsTwice returns the planted defect of a Feed whose Each
// returns one element twice over a fresh subject.
func feedEachDrainsTwice(h *FeedHandle) Feed {
	return &feedEachDrainsTwiceDefect{feedOverEmbedded: h.Fresh()}
}

// feedEachDrainsNothingDefect is the planted defect of a Feed whose Each
// loses what Put stored.
type feedEachDrainsNothingDefect struct {
	feedOverEmbedded
}

// Each is the override of the planted defect.
func (*feedEachDrainsNothingDefect) Each(p0 context.Context) (r0 iter.Seq[string]) {
	r0 = func(func(string) bool) {}
	return r0
}

// feedEachDrainsNothing returns the planted defect of a Feed whose Each
// loses what Put stored over a fresh subject.
func feedEachDrainsNothing(h *FeedHandle) Feed {
	return &feedEachDrainsNothingDefect{feedOverEmbedded: h.Fresh()}
}

// feedEachRepeatsAfterThreeDefect is the planted defect of a Feed whose
// Each yields its first element twice from its fourth call on.
type feedEachRepeatsAfterThreeDefect struct {
	feedOverEmbedded
	calls int
}

// Each is the override of the planted defect.
func (d *feedEachRepeatsAfterThreeDefect) Each(p0 context.Context) (r0 iter.Seq[string]) {
	r0 = d.feedOverEmbedded.Each(p0)
	d.calls++
	if d.calls > 3 && len(r0) > 0 {
		r0 = append(r0[:1:1], r0...)
	}
	return r0
}

// feedEachRepeatsAfterThree returns the planted defect of a Feed whose
// Each yields its first element twice from its fourth call on over a fresh
// subject.
func feedEachRepeatsAfterThree(h *FeedHandle) Feed {
	return &feedEachRepeatsAfterThreeDefect{feedOverEmbedded: h.Fresh()}
}

// feedPairsPanicsDefect is the planted defect of a Feed whose Pairs
// panics.
type feedPairsPanicsDefect struct {
	feedOverEmbedded
}

// Pairs is the override of the planted defect.
func (*feedPairsPanicsDefect) Pairs(p0 context.Context) (r0 iter.Seq2[string, int]) {
	panic("the planted defect of a Feed whose Pairs panics")
}

// feedPairsPanics returns the planted defect of a Feed whose Pairs panics
// over a fresh subject.
func feedPairsPanics(h *FeedHandle) Feed {
	return &feedPairsPanicsDefect{feedOverEmbedded: h.Fresh()}
}

// feedPairsDrainsTwiceDefect is the planted defect of a Feed whose Pairs
// returns one element twice.
type feedPairsDrainsTwiceDefect struct {
	feedOverEmbedded
}

// Pairs is the override of the planted defect.
func (d *feedPairsDrainsTwiceDefect) Pairs(p0 context.Context) (r0 iter.Seq2[string, int]) {
	r0 = d.feedOverEmbedded.Pairs(p0)
	if len(r0) > 0 {
		r0 = append(r0[:1:1], r0...)
	}
	return r0
}

// feedPairsDrainsTwice returns the planted defect of a Feed whose Pairs
// returns one element twice over a fresh subject.
func feedPairsDrainsTwice(h *FeedHandle) Feed {
	return &feedPairsDrainsTwiceDefect{feedOverEmbedded: h.Fresh()}
}

// feedPairsDrainsNothingDefect is the planted defect of a Feed whose Pairs
// loses what Put stored.
type feedPairsDrainsNothingDefect struct {
	feedOverEmbedded
}

// Pairs is the override of the planted defect.
func (*feedPairsDrainsNothingDefect) Pairs(p0 context.Context) (r0 iter.Seq2[string, int]) {
	r0 = func(func(string, int) bool) {}
	return r0
}

// feedPairsDrainsNothing returns the planted defect of a Feed whose Pairs
// loses what Put stored over a fresh subject.
func feedPairsDrainsNothing(h *FeedHandle) Feed {
	return &feedPairsDrainsNothingDefect{feedOverEmbedded: h.Fresh()}
}

// feedPairsRepeatsAfterThreeDefect is the planted defect of a Feed whose
// Pairs yields its first element twice from its fourth call on.
type feedPairsRepeatsAfterThreeDefect struct {
	feedOverEmbedded
	calls int
}

// Pairs is the override of the planted defect.
func (d *feedPairsRepeatsAfterThreeDefect) Pairs(p0 context.Context) (r0 iter.Seq2[string, int]) {
	r0 = d.feedOverEmbedded.Pairs(p0)
	d.calls++
	if d.calls > 3 && len(r0) > 0 {
		r0 = append(r0[:1:1], r0...)
	}
	return r0
}

// feedPairsRepeatsAfterThree returns the planted defect of a Feed whose
// Pairs yields its first element twice from its fourth call on over a
// fresh subject.
func feedPairsRepeatsAfterThree(h *FeedHandle) Feed {
	return &feedPairsRepeatsAfterThreeDefect{feedOverEmbedded: h.Fresh()}
}

// feedActions returns the actions of the generated sequences of Feed over
// the subject that subject returns. The actions of a case take their keys
// from the pools that it draws.
func feedActions(c *prop.Case, subject func() Feed) []stateful.Action[struct{}] {
	return []stateful.Action[struct{}]{
		{
			Name: "Put",
			Input: func(c *prop.Case, _ struct{}) any {
				item := c.Draw(prop.Of[string](), "put.item")
				return []any{item}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Put", args)
				values, err := feedPutCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "All",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "All", args)
				values, err := feedAllCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Each",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Each", args)
				values, err := feedEachCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Pairs",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Pairs", args)
				values, err := feedPairsCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// feedPutCall calls Put on subject with the arguments args, and returns
// the values of the call beside its error.
func feedPutCall(ctx context.Context, subject Feed, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Put(ctx, a0)
}

// feedAllCall calls All on subject with the arguments args, and returns
// the values of the call beside its error.
func feedAllCall(ctx context.Context, subject Feed, args []any) ([]any, error) {
	r0 := subject.All(ctx)
	var drained []string
	for e, err := range r0 {
		if err != nil {
			return []any{drained}, err
		}
		drained = append(drained, e)
	}
	return []any{drained}, nil
}

// feedEachCall calls Each on subject with the arguments args, and returns
// the values of the call beside its error.
func feedEachCall(ctx context.Context, subject Feed, args []any) ([]any, error) {
	r0 := subject.Each(ctx)
	return []any{slices.Collect(r0)}, nil
}

// feedPairsCall calls Pairs on subject with the arguments args, and
// returns the values of the call beside its error.
func feedPairsCall(ctx context.Context, subject Feed, args []any) ([]any, error) {
	r0 := subject.Pairs(ctx)
	var drained []probe.Pair[string, int]
	for k, v := range r0 {
		drained = append(drained, probe.Pair[string, int]{First: k, Second: v})
	}
	return []any{drained}, nil
}

// The derived checks of Register:
//
// Register.Set/smoke: Set returns without a panic for derived inputs.
//
// Register.Set/smoke#generated: Set returns without a panic for derived
// inputs.
//
// Register.Set/cancel: Set returns an error for a cancelled context.
//
// Register.Set/cancel#generated: Set returns an error for a cancelled
// context.
//
// Register.Get/smoke: Get returns without a panic for derived inputs.
//
// Register.Get/cancel: Get returns an error for a cancelled context.
//
// Register.Get/zero-on-error: Get returns the zero value beside an error.
//
// Register.Get/agrees-with-the-reference: every sequence of calls returns
// on Register what it returns on the reference.

// The derived checks of Register that witness their laws at the shallow
// end:
//
// Register.Set/smoke
//
// Register.Get/smoke

// The unclaimed classifications of Register:
//
// The shape:writer of Register.Set has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Register writes.
type (
	RegisterHandle = conform.Handle[Register]
	RegisterCheck  = conform.Check[Register]
	RegisterOption = conform.Option[Register]
)

// RegisterChecks contains the name of each callable of Register, with the
// identities of its checks.
var RegisterChecks = struct {
	Set registerSet
	Get registerGet
}{
	Set: registerSet{"Register.Set"},
	Get: registerGet{"Register.Get"},
}

// registerSet is the name of Register.Set, with the identities of its
// checks.
type registerSet struct {
	conform.Callable[Register]
}

// Smoke returns the identity of the check that Set returns without a panic
// for derived inputs.
func (c registerSet) Smoke() conform.ID[Register] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Set returns
// without a panic for derived inputs.
func (c registerSet) SmokeGenerated() conform.ID[Register] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Set returns an error for a
// cancelled context.
func (c registerSet) Cancel() conform.ID[Register] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Set returns an
// error for a cancelled context.
func (c registerSet) CancelGenerated() conform.ID[Register] {
	return c.Check("cancel#generated")
}

// registerGet is the name of Register.Get, with the identities of its
// checks.
type registerGet struct {
	conform.Callable[Register]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c registerGet) Smoke() conform.ID[Register] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c registerGet) Cancel() conform.ID[Register] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c registerGet) ZeroOnError() conform.ID[Register] {
	return c.Check("zero-on-error")
}

// AgreesWithTheReference returns the identity of the check that every
// sequence of calls returns on Register what it returns on the reference.
func (c registerGet) AgreesWithTheReference() conform.ID[Register] {
	return "Register/agrees-with-the-reference"
}

// RegisterSuite returns the checks of Register, and the calls and the
// budgets of its callables.
func RegisterSuite() conform.Suite[Register] {
	return conform.Suite[Register]{
		Declaration: "Register",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Register]{
			{
				ID:          RegisterChecks.Set.Smoke(),
				Statement:   "Set returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RegisterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, string("v-a"))
					})
				},
				Defect: registerSetPanics,
			},
			{
				ID:          RegisterChecks.Set.SmokeGenerated(),
				Statement:   "Set returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegisterHandle, c *prop.Case) {
					drawnV := c.Draw(prop.Of[string](), "v")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, drawnV)
					})
				},
				Defect: registerSetPanics,
			},
			{
				ID:          RegisterChecks.Set.Cancel(),
				Statement:   "Set returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegisterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, string("v-a"))
					})
				},
				Defect: registerSetIgnoresCancel,
			},
			{
				ID:          RegisterChecks.Set.CancelGenerated(),
				Statement:   "Set returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *RegisterHandle, c *prop.Case) {
					drawnV := c.Draw(prop.Of[string](), "v")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Set(ctx, drawnV)
					})
				},
				Defect: registerSetIgnoresCancel,
			},
			{
				ID:          RegisterChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *RegisterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx)
						return err
					})
				},
				Defect: registerGetPanics,
			},
			{
				ID:          RegisterChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegisterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx)
						return err
					})
				},
				Defect: registerGetIgnoresCancel,
			},
			{
				ID:          RegisterChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *RegisterHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().Get(ctx)
					}, string(""))
				},
				Defect: registerGetValueBeside,
			},
			{
				ID:        RegisterChecks.Get.AgreesWithTheReference(),
				Statement: "every sequence of calls returns on Register what it returns on the reference",
				Claim:     spec.AgreesWithTheReference,
				Probe:     spec.SequencesAgreeWithTheReference,
				Reason:    spec.DivergesFromTheReference,
				Law:       spec.EverySequenceAgreesWithTheReference,
				Machineries: []spec.MachineryID{
					spec.GeneratedSequences,
					spec.ReferenceImplementation,
				},
				Property: func(h *RegisterHandle, c *prop.Case) {
					probe.SequencesAgreeWithTheReference(h, c, registerActions, registerReference)
				},
				Defect: registerGetDivergesAfterThree,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Register.Set",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Register]{
			{
				Calls:    registerSetCalls,
				Over:     newRegisterOver,
				Callable: RegisterChecks.Set.Callable,
			},
			{
				Calls:    registerGetCalls,
				Over:     newRegisterOver,
				Callable: RegisterChecks.Get.Callable,
			},
		},
	}
}

// RunRegisterSuite applies the checks of Register, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunRegisterSuite[T Register](t *testing.T, newSubject func() T, opts ...RegisterOption) {
	t.Helper()
	RunRegisterSuiteFromHandle(t, func(*RegisterHandle) T { return newSubject() }, opts...)
}

// RunRegisterSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunRegisterSuiteFromHandle[T Register](
	t *testing.T,
	newSubject func(h *RegisterHandle) T,
	opts ...RegisterOption,
) {
	t.Helper()
	conform.Run(t, RegisterSuite(), registerSubjects(newSubject), opts...)
}

// FuzzRegisterSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzRegisterSuite[T Register](
	f *testing.F,
	callable conform.Callable[Register],
	newSubject func(h *RegisterHandle) T,
	opts ...RegisterOption,
) {
	f.Helper()
	conform.Fuzz(f, RegisterSuite(), callable, registerSubjects(newSubject), opts...)
}

// RunRegisterBudgets applies the allocation ceilings of Register and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunRegisterBudgets[T Register](
	t *testing.T,
	newSubject func(h *RegisterHandle) T,
	opts ...RegisterOption,
) {
	t.Helper()
	conform.RunBudgets(t, RegisterSuite(), registerSubjects(newSubject), opts...)
}

// BenchRegisterSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchRegisterSuite[T Register](
	b conform.B,
	callable conform.Callable[Register],
	newSubject func(h *RegisterHandle) T,
	opts ...RegisterOption,
) {
	b.Helper()
	conform.Bench(b, RegisterSuite(), callable, registerSubjects(newSubject), opts...)
}

// registerSubjects returns the factory of the subjects of Register that
// calls newSubject with the handle of each check.
func registerSubjects[T Register](
	newSubject func(h *RegisterHandle) T,
) func(h *RegisterHandle) Register {
	return func(h *RegisterHandle) Register { return newSubject(h) }
}

// registerSetCalls generates the calls of Register.Set with the inputs
// that each case draws.
var registerSetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Register) error {
	v := c.Draw(prop.Of[string](), "v")
	return func(ctx context.Context, subject Register) error {
		return subject.Set(ctx, v)
	}
})

// registerGetCalls generates the calls of Register.Get with the inputs
// that each case draws.
var registerGetCalls = prop.Composite(func(*prop.Case) func(context.Context, Register) error {
	return func(ctx context.Context, subject Register) error {
		_, err := subject.Get(ctx)
		return err
	}
})

// registerOverEmbedded is Register under the name of the field that
// registerOver embeds.
type registerOverEmbedded = Register

// registerOver is a subject of Register whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type registerOver struct {
	registerOverEmbedded
	spend func()
}

// Set calls spend and returns its zero results.
func (o registerOver) Set(context.Context, string) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o registerOver) Get(context.Context) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newRegisterOver returns the planted defect of a ceiling of Register,
// whose callables call spend.
func newRegisterOver(_ *RegisterHandle, spend func()) Register {
	return registerOver{spend: spend}
}

// registerSetPanicsDefect is the planted defect of a Register whose Set
// panics.
type registerSetPanicsDefect struct {
	registerOverEmbedded
}

// Set is the override of the planted defect.
func (*registerSetPanicsDefect) Set(p0 context.Context, p1 string) (r0 error) {
	panic("the planted defect of a Register whose Set panics")
}

// registerSetPanics returns the planted defect of a Register whose Set
// panics over a fresh subject.
func registerSetPanics(h *RegisterHandle) Register {
	return &registerSetPanicsDefect{registerOverEmbedded: h.Fresh()}
}

// registerSetIgnoresCancelDefect is the planted defect of a Register whose
// Set ignores the context.
type registerSetIgnoresCancelDefect struct {
	registerOverEmbedded
}

// Set is the override of the planted defect.
func (d *registerSetIgnoresCancelDefect) Set(p0 context.Context, p1 string) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.registerOverEmbedded.Set(p0, p1)
}

// registerSetIgnoresCancel returns the planted defect of a Register whose
// Set ignores the context over a fresh subject.
func registerSetIgnoresCancel(h *RegisterHandle) Register {
	return &registerSetIgnoresCancelDefect{registerOverEmbedded: h.Fresh()}
}

// registerGetPanicsDefect is the planted defect of a Register whose Get
// panics.
type registerGetPanicsDefect struct {
	registerOverEmbedded
}

// Get is the override of the planted defect.
func (*registerGetPanicsDefect) Get(p0 context.Context) (r0 string, r1 error) {
	panic("the planted defect of a Register whose Get panics")
}

// registerGetPanics returns the planted defect of a Register whose Get
// panics over a fresh subject.
func registerGetPanics(h *RegisterHandle) Register {
	return &registerGetPanicsDefect{registerOverEmbedded: h.Fresh()}
}

// registerGetIgnoresCancelDefect is the planted defect of a Register whose
// Get ignores the context.
type registerGetIgnoresCancelDefect struct {
	registerOverEmbedded
}

// Get is the override of the planted defect.
func (d *registerGetIgnoresCancelDefect) Get(p0 context.Context) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.registerOverEmbedded.Get(p0)
}

// registerGetIgnoresCancel returns the planted defect of a Register whose
// Get ignores the context over a fresh subject.
func registerGetIgnoresCancel(h *RegisterHandle) Register {
	return &registerGetIgnoresCancelDefect{registerOverEmbedded: h.Fresh()}
}

// registerGetValueBesideDefect is the planted defect of a Register whose
// Get returns a value beside an error.
type registerGetValueBesideDefect struct {
	registerOverEmbedded
}

// Get is the override of the planted defect.
func (d *registerGetValueBesideDefect) Get(p0 context.Context) (r0 string, r1 error) {
	r0, r1 = d.registerOverEmbedded.Get(p0)
	if r1 != nil {
		r0 = string("sample-b")
	}
	return r0, r1
}

// registerGetValueBeside returns the planted defect of a Register whose
// Get returns a value beside an error over a fresh subject.
func registerGetValueBeside(h *RegisterHandle) Register {
	return &registerGetValueBesideDefect{registerOverEmbedded: h.Fresh()}
}

// registerGetDivergesAfterThreeDefect is the planted defect of a Register
// whose first callable with a value result returns another value from its
// fourth call on.
type registerGetDivergesAfterThreeDefect struct {
	registerOverEmbedded
	calls int
}

// Get is the override of the planted defect.
func (d *registerGetDivergesAfterThreeDefect) Get(p0 context.Context) (r0 string, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.registerOverEmbedded.Get(p0)
	d.calls++
	if d.calls > 3 {
		r0 = string("sample-b")
	}
	return r0, r1
}

// registerGetDivergesAfterThree returns the planted defect of a Register
// whose first callable with a value result returns another value from its
// fourth call on over a fresh subject.
func registerGetDivergesAfterThree(h *RegisterHandle) Register {
	return &registerGetDivergesAfterThreeDefect{registerOverEmbedded: h.Fresh()}
}

// registerActions returns the actions of the generated sequences of
// Register over the subject that subject returns. The actions of a case
// take their keys from the pools that it draws.
func registerActions(c *prop.Case, subject func() Register) []stateful.Action[struct{}] {
	return []stateful.Action[struct{}]{
		{
			Name: "Set",
			Input: func(c *prop.Case, _ struct{}) any {
				v := c.Draw(prop.Of[string](), "set.v")
				return []any{v}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Set", args)
				values, err := registerSetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Get",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Get", args)
				values, err := registerGetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// registerSetCall calls Set on subject with the arguments args, and
// returns the values of the call beside its error.
func registerSetCall(ctx context.Context, subject Register, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	return nil, subject.Set(ctx, a0)
}

// registerGetCall calls Get on subject with the arguments args, and
// returns the values of the call beside its error.
func registerGetCall(ctx context.Context, subject Register, args []any) ([]any, error) {
	r0, err := subject.Get(ctx)
	return Unstamped("Get", []any{r0}), err
}

// registerReference returns the function that applies an operation of a
// history to the reference r by the name of its callable.
func registerReference(r Register) func(context.Context, string, []any) conform.Output {
	return func(ctx context.Context, operation string, args []any) conform.Output {
		var values []any
		var err error
		switch operation {
		case "Set":
			values, err = registerSetCall(ctx, r, args)
		case "Get":
			values, err = registerGetCall(ctx, r, args)
		}
		return conform.OutputOf(err, values...)
	}
}

// The derived checks of Double:
//
// Double/smoke: Double returns without a panic for derived inputs.
//
// Double/smoke#generated: Double returns without a panic for derived
// inputs.
//
// Double/cancel: Double returns an error for a cancelled context.
//
// Double/cancel#generated: Double returns an error for a cancelled
// context.
//
// Double/zero-on-error: Double returns the zero value beside an error.
//
// Double/zero-on-error#generated: Double returns the zero value beside an
// error.
//
// Double/agrees-with-the-reference: every sequence of calls returns on
// Double what it returns on the reference.

// The derived checks of Double that witness their laws at the shallow end:
//
// Double/smoke

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Double writes.
type (
	DoubleHandle = conform.Handle[DoubleFunc]
	DoubleCheck  = conform.Check[DoubleFunc]
	DoubleOption = conform.Option[DoubleFunc]
)

// DoubleFunc is the type of the function Double. Each subject of the suite
// of Double has it.
type DoubleFunc func(ctx context.Context, n int) (int, error)

// DoubleChecks contains the name of each callable of Double, with the
// identities of its checks.
var DoubleChecks = doubleCallable{"Double"}

// doubleCallable is the name of Double, with the identities of its checks.
type doubleCallable struct {
	conform.Callable[DoubleFunc]
}

// Smoke returns the identity of the check that Double returns without a
// panic for derived inputs.
func (c doubleCallable) Smoke() conform.ID[DoubleFunc] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Double returns
// without a panic for derived inputs.
func (c doubleCallable) SmokeGenerated() conform.ID[DoubleFunc] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Double returns an error
// for a cancelled context.
func (c doubleCallable) Cancel() conform.ID[DoubleFunc] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Double returns an
// error for a cancelled context.
func (c doubleCallable) CancelGenerated() conform.ID[DoubleFunc] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Double returns the
// zero value beside an error.
func (c doubleCallable) ZeroOnError() conform.ID[DoubleFunc] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Double
// returns the zero value beside an error.
func (c doubleCallable) ZeroOnErrorGenerated() conform.ID[DoubleFunc] {
	return c.Check("zero-on-error#generated")
}

// AgreesWithTheReference returns the identity of the check that every
// sequence of calls returns on Double what it returns on the reference.
func (c doubleCallable) AgreesWithTheReference() conform.ID[DoubleFunc] {
	return "Double/agrees-with-the-reference"
}

// DoubleSuite returns the checks of Double, and the calls and the budgets
// of its callables.
func DoubleSuite() conform.Suite[DoubleFunc] {
	return conform.Suite[DoubleFunc]{
		Declaration: "Double",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[DoubleFunc]{
			{
				ID:          DoubleChecks.Smoke(),
				Statement:   "Double returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *DoubleHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, int(42))
						return err
					})
				},
				Defect: doublePanics,
			},
			{
				ID:          DoubleChecks.SmokeGenerated(),
				Statement:   "Double returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *DoubleHandle, c *prop.Case) {
					drawnN := c.Draw(prop.Of[int](), "n")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, drawnN)
						return err
					})
				},
				Defect: doublePanics,
			},
			{
				ID:          DoubleChecks.Cancel(),
				Statement:   "Double returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *DoubleHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, int(42))
						return err
					})
				},
				Defect: doubleIgnoresCancel,
			},
			{
				ID:          DoubleChecks.CancelGenerated(),
				Statement:   "Double returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *DoubleHandle, c *prop.Case) {
					drawnN := c.Draw(prop.Of[int](), "n")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject()(ctx, drawnN)
						return err
					})
				},
				Defect: doubleIgnoresCancel,
			},
			{
				ID:          DoubleChecks.ZeroOnError(),
				Statement:   "Double returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *DoubleHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject()(ctx, int(42))
					}, int(0))
				},
				Defect: doubleValueBeside,
			},
			{
				ID:          DoubleChecks.ZeroOnErrorGenerated(),
				Statement:   "Double returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *DoubleHandle, c *prop.Case) {
					drawnN := c.Draw(prop.Of[int](), "n")
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject()(ctx, drawnN)
					}, int(0))
				},
				Defect: doubleValueBeside,
			},
			{
				ID:        DoubleChecks.AgreesWithTheReference(),
				Statement: "every sequence of calls returns on Double what it returns on the reference",
				Claim:     spec.AgreesWithTheReference,
				Probe:     spec.SequencesAgreeWithTheReference,
				Reason:    spec.DivergesFromTheReference,
				Law:       spec.EverySequenceAgreesWithTheReference,
				Machineries: []spec.MachineryID{
					spec.GeneratedSequences,
					spec.ReferenceImplementation,
				},
				Property: func(h *DoubleHandle, c *prop.Case) {
					probe.SequencesAgreeWithTheReference(h, c, doubleActions, doubleReference)
				},
				Defect: doubleDivergesAfterThree,
			},
		},
		Calls: []conform.Call[DoubleFunc]{
			{
				Calls:    doubleCalls,
				Over:     newDoubleOver,
				Callable: DoubleChecks.Callable,
			},
		},
	}
}

// RunDoubleSuite applies the checks of Double, and the checks that opts
// add, to subject.
func RunDoubleSuite(t *testing.T, subject DoubleFunc, opts ...DoubleOption) {
	t.Helper()
	RunDoubleSuiteFromHandle(t, doubleSubjects(subject), opts...)
}

// RunDoubleSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunDoubleSuiteFromHandle(
	t *testing.T,
	newSubject func(h *DoubleHandle) DoubleFunc,
	opts ...DoubleOption,
) {
	t.Helper()
	conform.Run(t, DoubleSuite(), newSubject, opts...)
}

// FuzzDoubleSuite registers on f the fuzz target of the property checks of
// Double on subject.
func FuzzDoubleSuite(f *testing.F, subject DoubleFunc, opts ...DoubleOption) {
	f.Helper()
	conform.Fuzz(f, DoubleSuite(), DoubleChecks.Callable, doubleSubjects(subject), opts...)
}

// RunDoubleBudgets applies the allocation ceiling of Double and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to subject.
func RunDoubleBudgets(t *testing.T, subject DoubleFunc, opts ...DoubleOption) {
	t.Helper()
	conform.RunBudgets(t, DoubleSuite(), doubleSubjects(subject), opts...)
}

// BenchDoubleSuite measures subject under the ceilings of the budget of
// Double.
func BenchDoubleSuite(b conform.B, subject DoubleFunc, opts ...DoubleOption) {
	b.Helper()
	conform.Bench(b, DoubleSuite(), DoubleChecks.Callable, doubleSubjects(subject), opts...)
}

// doubleSubjects returns the factory of the subjects of Double that
// returns subject for each handle.
func doubleSubjects(subject DoubleFunc) func(*DoubleHandle) DoubleFunc {
	return func(*DoubleHandle) DoubleFunc { return subject }
}

// doubleCalls generates the calls of Double with the inputs that each case
// draws.
var doubleCalls = prop.Composite(func(c *prop.Case) func(context.Context, DoubleFunc) error {
	n := c.Draw(prop.Of[int](), "n")
	return func(ctx context.Context, subject DoubleFunc) error {
		_, err := subject(ctx, n)
		return err
	}
})

// newDoubleOver returns the planted defect of a ceiling of Double, which
// calls spend and returns its zero results.
func newDoubleOver(_ *DoubleHandle, spend func()) DoubleFunc {
	return func(context.Context, int) (r0 int, r1 error) {
		spend()
		return r0, r1
	}
}

// doublePanics returns the planted defect of a Double whose Double panics
// over the fresh function.
func doublePanics(h *DoubleHandle) DoubleFunc {
	return func(p0 context.Context, p1 int) (r0 int, r1 error) {
		panic("the planted defect of a Double whose Double panics")
	}
}

// doubleIgnoresCancel returns the planted defect of a Double whose Double
// ignores the context over the fresh function.
func doubleIgnoresCancel(h *DoubleHandle) DoubleFunc {
	fresh := h.Fresh()
	return func(p0 context.Context, p1 int) (r0 int, r1 error) {
		if p0.Err() != nil {
			return r0, r1
		}
		return fresh(p0, p1)
	}
}

// doubleValueBeside returns the planted defect of a Double whose Double
// returns a value beside an error over the fresh function.
func doubleValueBeside(h *DoubleHandle) DoubleFunc {
	fresh := h.Fresh()
	return func(p0 context.Context, p1 int) (r0 int, r1 error) {
		r0, r1 = fresh(p0, p1)
		if r1 != nil {
			r0 = int(7)
		}
		return r0, r1
	}
}

// doubleDivergesAfterThree returns the planted defect of a Double whose
// first callable with a value result returns another value from its fourth
// call on over the fresh function.
func doubleDivergesAfterThree(h *DoubleHandle) DoubleFunc {
	fresh := h.Fresh()
	calls := 0
	return func(p0 context.Context, p1 int) (r0 int, r1 error) {
		if err := p0.Err(); err != nil {
			r1 = err
			return r0, r1
		}
		r0, r1 = fresh(p0, p1)
		calls++
		if calls > 3 {
			r0 = int(7)
		}
		return r0, r1
	}
}

// doubleActions returns the actions of the generated sequences of Double
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func doubleActions(c *prop.Case, subject func() DoubleFunc) []stateful.Action[struct{}] {
	keysInt := c.Draw(
		prop.List(prop.Of[int](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.int",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Double",
			Input: func(c *prop.Case, _ struct{}) any {
				n := c.Draw(prop.SampledFrom(keysInt...), "double.n")
				return []any{n}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Double", args, args[0])
				values, err := doubleCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// doubleCall calls Double on subject with the arguments args, and returns
// the values of the call beside its error.
func doubleCall(ctx context.Context, subject DoubleFunc, args []any) ([]any, error) {
	a0, _ := args[0].(int)
	r0, err := subject(ctx, a0)
	return []any{r0}, err
}

// doubleReference returns the function that applies an operation of a
// history to the reference r by the name of its callable.
func doubleReference(r DoubleFunc) func(context.Context, string, []any) conform.Output {
	return func(ctx context.Context, operation string, args []any) conform.Output {
		var values []any
		var err error
		switch operation {
		case "Double":
			values, err = doubleCall(ctx, r, args)
		}
		return conform.OutputOf(err, values...)
	}
}

// The derived checks of Cached:
//
// Cached.Get/smoke: Get returns without a panic for derived inputs.
//
// Cached.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Cached.Get/cancel: Get returns an error for a cancelled context.
//
// Cached.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Cached.Get/zero-on-error: Get returns the zero value beside an error.
//
// Cached.Get/zero-on-error#generated: Get returns the zero value beside an
// error.
//
// Cached.Get/cache-agrees-over-sequences: Get returns what Load returns
// for its arguments, over any sequence of calls.
//
// Cached.Load/smoke: Load returns without a panic for derived inputs.
//
// Cached.Load/smoke#generated: Load returns without a panic for derived
// inputs.
//
// Cached.Load/cancel: Load returns an error for a cancelled context.
//
// Cached.Load/cancel#generated: Load returns an error for a cancelled
// context.
//
// Cached.Load/zero-on-error: Load returns the zero value beside an error.
//
// Cached.Load/zero-on-error#generated: Load returns the zero value beside
// an error.
//
// Cached.Peek/smoke: Peek returns without a panic for derived inputs.
//
// Cached.Peek/smoke#generated: Peek returns without a panic for derived
// inputs.
//
// Cached.Peek/cancel: Peek returns an error for a cancelled context.
//
// Cached.Peek/cancel#generated: Peek returns an error for a cancelled
// context.
//
// Cached.Peek/zero-on-error: Peek returns the zero value beside an error.
//
// Cached.Peek/zero-on-error#generated: Peek returns the zero value beside
// an error.
//
// Cached.Peek/cache-agrees-over-sequences: Peek returns what Fetch returns
// for its arguments, over any sequence of calls.
//
// Cached.Fetch/smoke: Fetch returns without a panic for derived inputs.
//
// Cached.Fetch/smoke#generated: Fetch returns without a panic for derived
// inputs.
//
// Cached.Fetch/cancel: Fetch returns an error for a cancelled context.
//
// Cached.Fetch/cancel#generated: Fetch returns an error for a cancelled
// context.
//
// Cached.Fetch/zero-on-error: Fetch returns the zero value beside an
// error.
//
// Cached.Fetch/zero-on-error#generated: Fetch returns the zero value
// beside an error.

// The derived checks of Cached that witness their laws at the shallow end:
//
// Cached.Get/smoke
//
// Cached.Load/smoke
//
// Cached.Peek/smoke
//
// Cached.Fetch/smoke

// The refused claims of Cached:
//
// Cached.Get/cache: no callable of the declaration is the partner put of
// Get. The remedy is to name the partner of Get with the param of the
// classification.
//
// Cached.Peek/cache: no callable of the declaration is the partner put of
// Peek. The remedy is to name the partner of Peek with the param of the
// classification.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Cached writes.
type (
	CachedHandle = conform.Handle[Cached]
	CachedCheck  = conform.Check[Cached]
	CachedOption = conform.Option[Cached]
)

// CachedChecks contains the name of each callable of Cached, with the
// identities of its checks.
var CachedChecks = struct {
	Get   cachedGet
	Load  cachedLoad
	Peek  cachedPeek
	Fetch cachedFetch
}{
	Get:   cachedGet{"Cached.Get"},
	Load:  cachedLoad{"Cached.Load"},
	Peek:  cachedPeek{"Cached.Peek"},
	Fetch: cachedFetch{"Cached.Fetch"},
}

// cachedGet is the name of Cached.Get, with the identities of its checks.
type cachedGet struct {
	conform.Callable[Cached]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c cachedGet) Smoke() conform.ID[Cached] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c cachedGet) SmokeGenerated() conform.ID[Cached] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c cachedGet) Cancel() conform.ID[Cached] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c cachedGet) CancelGenerated() conform.ID[Cached] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c cachedGet) ZeroOnError() conform.ID[Cached] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c cachedGet) ZeroOnErrorGenerated() conform.ID[Cached] {
	return c.Check("zero-on-error#generated")
}

// CacheAgreesOverSequences returns the identity of the check that Get
// returns what Load returns for its arguments, over any sequence of calls.
func (c cachedGet) CacheAgreesOverSequences() conform.ID[Cached] {
	return c.Check("cache-agrees-over-sequences")
}

// cachedLoad is the name of Cached.Load, with the identities of its
// checks.
type cachedLoad struct {
	conform.Callable[Cached]
}

// Smoke returns the identity of the check that Load returns without a
// panic for derived inputs.
func (c cachedLoad) Smoke() conform.ID[Cached] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Load returns
// without a panic for derived inputs.
func (c cachedLoad) SmokeGenerated() conform.ID[Cached] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Load returns an error for
// a cancelled context.
func (c cachedLoad) Cancel() conform.ID[Cached] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Load returns an
// error for a cancelled context.
func (c cachedLoad) CancelGenerated() conform.ID[Cached] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Load returns the zero
// value beside an error.
func (c cachedLoad) ZeroOnError() conform.ID[Cached] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Load returns
// the zero value beside an error.
func (c cachedLoad) ZeroOnErrorGenerated() conform.ID[Cached] {
	return c.Check("zero-on-error#generated")
}

// cachedPeek is the name of Cached.Peek, with the identities of its
// checks.
type cachedPeek struct {
	conform.Callable[Cached]
}

// Smoke returns the identity of the check that Peek returns without a
// panic for derived inputs.
func (c cachedPeek) Smoke() conform.ID[Cached] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Peek returns
// without a panic for derived inputs.
func (c cachedPeek) SmokeGenerated() conform.ID[Cached] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Peek returns an error for
// a cancelled context.
func (c cachedPeek) Cancel() conform.ID[Cached] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Peek returns an
// error for a cancelled context.
func (c cachedPeek) CancelGenerated() conform.ID[Cached] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Peek returns the zero
// value beside an error.
func (c cachedPeek) ZeroOnError() conform.ID[Cached] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Peek returns
// the zero value beside an error.
func (c cachedPeek) ZeroOnErrorGenerated() conform.ID[Cached] {
	return c.Check("zero-on-error#generated")
}

// CacheAgreesOverSequences returns the identity of the check that Peek
// returns what Fetch returns for its arguments, over any sequence of
// calls.
func (c cachedPeek) CacheAgreesOverSequences() conform.ID[Cached] {
	return c.Check("cache-agrees-over-sequences")
}

// cachedFetch is the name of Cached.Fetch, with the identities of its
// checks.
type cachedFetch struct {
	conform.Callable[Cached]
}

// Smoke returns the identity of the check that Fetch returns without a
// panic for derived inputs.
func (c cachedFetch) Smoke() conform.ID[Cached] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Fetch returns
// without a panic for derived inputs.
func (c cachedFetch) SmokeGenerated() conform.ID[Cached] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Fetch returns an error for
// a cancelled context.
func (c cachedFetch) Cancel() conform.ID[Cached] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Fetch returns an
// error for a cancelled context.
func (c cachedFetch) CancelGenerated() conform.ID[Cached] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Fetch returns the
// zero value beside an error.
func (c cachedFetch) ZeroOnError() conform.ID[Cached] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Fetch
// returns the zero value beside an error.
func (c cachedFetch) ZeroOnErrorGenerated() conform.ID[Cached] {
	return c.Check("zero-on-error#generated")
}

// CachedSuite returns the checks of Cached, and the calls and the budgets
// of its callables.
func CachedSuite() conform.Suite[Cached] {
	return conform.Suite[Cached]{
		Declaration: "Cached",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Cached]{
			{
				ID:          CachedChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CachedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: cachedGetPanics,
			},
			{
				ID:          CachedChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: cachedGetPanics,
			},
			{
				ID:          CachedChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: cachedGetIgnoresCancel,
			},
			{
				ID:          CachedChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: cachedGetIgnoresCancel,
			},
			{
				ID:          CachedChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Item{})
				},
				Defect: cachedGetValueBeside,
			},
			{
				ID:          CachedChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, Item{})
				},
				Defect: cachedGetValueBeside,
			},
			{
				ID:          CachedChecks.Get.CacheAgreesOverSequences(),
				Statement:   "Get returns what Load returns for its arguments, over any sequence of calls",
				Claim:       spec.CacheAgreesOverSequences,
				Probe:       spec.ReadsAgreeWithTheBacking,
				Reason:      spec.CacheDisagrees,
				Law:         spec.CacheAgreesWithItsBacking,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *CachedHandle, c *prop.Case) {
					probe.ReadsAgreeWithTheBacking(
						h,
						c,
						cachedActions,
						"Get",
						func(ctx context.Context, args []any) conform.Output {
							values, err := cachedLoadCall(ctx, h.Subject(), []any{string("id-a")})
							return conform.OutputOf(err, values...)
						},
					)
				},
				Defect: cachedGetServesAStaleValueAfterThree,
			},
			{
				ID:          CachedChecks.Load.Smoke(),
				Statement:   "Load returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CachedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Load(ctx, string("id-a"))
						return err
					})
				},
				Defect: cachedLoadPanics,
			},
			{
				ID:          CachedChecks.Load.SmokeGenerated(),
				Statement:   "Load returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnID := c.Draw(prop.Of[string](), "id")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Load(ctx, drawnID)
						return err
					})
				},
				Defect: cachedLoadPanics,
			},
			{
				ID:          CachedChecks.Load.Cancel(),
				Statement:   "Load returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Load(ctx, string("id-a"))
						return err
					})
				},
				Defect: cachedLoadIgnoresCancel,
			},
			{
				ID:          CachedChecks.Load.CancelGenerated(),
				Statement:   "Load returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnID := c.Draw(prop.Of[string](), "id")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Load(ctx, drawnID)
						return err
					})
				},
				Defect: cachedLoadIgnoresCancel,
			},
			{
				ID:          CachedChecks.Load.ZeroOnError(),
				Statement:   "Load returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Load(ctx, string("id-a"))
					}, Item{})
				},
				Defect: cachedLoadValueBeside,
			},
			{
				ID:          CachedChecks.Load.ZeroOnErrorGenerated(),
				Statement:   "Load returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnID := c.Draw(prop.Of[string](), "id")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Load(ctx, drawnID)
					}, Item{})
				},
				Defect: cachedLoadValueBeside,
			},
			{
				ID:          CachedChecks.Peek.Smoke(),
				Statement:   "Peek returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CachedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Peek(ctx, string("key-a"))
						return err
					})
				},
				Defect: cachedPeekPanics,
			},
			{
				ID:          CachedChecks.Peek.SmokeGenerated(),
				Statement:   "Peek returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Peek(ctx, drawnKey)
						return err
					})
				},
				Defect: cachedPeekPanics,
			},
			{
				ID:          CachedChecks.Peek.Cancel(),
				Statement:   "Peek returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Peek(ctx, string("key-a"))
						return err
					})
				},
				Defect: cachedPeekIgnoresCancel,
			},
			{
				ID:          CachedChecks.Peek.CancelGenerated(),
				Statement:   "Peek returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Peek(ctx, drawnKey)
						return err
					})
				},
				Defect: cachedPeekIgnoresCancel,
			},
			{
				ID:          CachedChecks.Peek.ZeroOnError(),
				Statement:   "Peek returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Peek(ctx, string("key-a"))
					}, Item{})
				},
				Defect: cachedPeekValueBeside,
			},
			{
				ID:          CachedChecks.Peek.ZeroOnErrorGenerated(),
				Statement:   "Peek returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Peek(ctx, drawnKey)
					}, Item{})
				},
				Defect: cachedPeekValueBeside,
			},
			{
				ID:          CachedChecks.Peek.CacheAgreesOverSequences(),
				Statement:   "Peek returns what Fetch returns for its arguments, over any sequence of calls",
				Claim:       spec.CacheAgreesOverSequences,
				Probe:       spec.ReadsAgreeWithTheBacking,
				Reason:      spec.CacheDisagrees,
				Law:         spec.CacheAgreesWithItsBacking,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Property: func(h *CachedHandle, c *prop.Case) {
					probe.ReadsAgreeWithTheBacking(
						h,
						c,
						cachedActions,
						"Peek",
						func(ctx context.Context, args []any) conform.Output {
							values, err := cachedFetchCall(ctx, h.Subject(), []any{int(42)})
							return conform.OutputOf(err, values...)
						},
					)
				},
				Defect: cachedPeekServesAStaleValueAfterThree,
			},
			{
				ID:          CachedChecks.Fetch.Smoke(),
				Statement:   "Fetch returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CachedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Fetch(ctx, int(42))
						return err
					})
				},
				Defect: cachedFetchPanics,
			},
			{
				ID:          CachedChecks.Fetch.SmokeGenerated(),
				Statement:   "Fetch returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[int](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Fetch(ctx, drawnKey)
						return err
					})
				},
				Defect: cachedFetchPanics,
			},
			{
				ID:          CachedChecks.Fetch.Cancel(),
				Statement:   "Fetch returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Fetch(ctx, int(42))
						return err
					})
				},
				Defect: cachedFetchIgnoresCancel,
			},
			{
				ID:          CachedChecks.Fetch.CancelGenerated(),
				Statement:   "Fetch returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[int](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Fetch(ctx, drawnKey)
						return err
					})
				},
				Defect: cachedFetchIgnoresCancel,
			},
			{
				ID:          CachedChecks.Fetch.ZeroOnError(),
				Statement:   "Fetch returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CachedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Fetch(ctx, int(42))
					}, Item{})
				},
				Defect: cachedFetchValueBeside,
			},
			{
				ID:          CachedChecks.Fetch.ZeroOnErrorGenerated(),
				Statement:   "Fetch returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *CachedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[int](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Fetch(ctx, drawnKey)
					}, Item{})
				},
				Defect: cachedFetchValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Cache,
				Callable: "Cached.Get",
				Gate:     spec.PartnerFound,
				Why:      "no callable of the declaration is the partner put of Get",
				Remedy:   "name the partner of Get with the param of the classification",
			},
			{
				Claim:    spec.Cache,
				Callable: "Cached.Peek",
				Gate:     spec.PartnerFound,
				Why:      "no callable of the declaration is the partner put of Peek",
				Remedy:   "name the partner of Peek with the param of the classification",
			},
		},
		Calls: []conform.Call[Cached]{
			{
				Calls:    cachedGetCalls,
				Over:     newCachedOver,
				Callable: CachedChecks.Get.Callable,
			},
			{
				Calls:    cachedLoadCalls,
				Over:     newCachedOver,
				Callable: CachedChecks.Load.Callable,
			},
			{
				Calls:    cachedPeekCalls,
				Over:     newCachedOver,
				Callable: CachedChecks.Peek.Callable,
			},
			{
				Calls:    cachedFetchCalls,
				Over:     newCachedOver,
				Callable: CachedChecks.Fetch.Callable,
			},
		},
	}
}

// RunCachedSuite applies the checks of Cached, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunCachedSuite[T Cached](t *testing.T, newSubject func() T, opts ...CachedOption) {
	t.Helper()
	RunCachedSuiteFromHandle(t, func(*CachedHandle) T { return newSubject() }, opts...)
}

// RunCachedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunCachedSuiteFromHandle[T Cached](
	t *testing.T,
	newSubject func(h *CachedHandle) T,
	opts ...CachedOption,
) {
	t.Helper()
	conform.Run(t, CachedSuite(), cachedSubjects(newSubject), opts...)
}

// FuzzCachedSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzCachedSuite[T Cached](
	f *testing.F,
	callable conform.Callable[Cached],
	newSubject func(h *CachedHandle) T,
	opts ...CachedOption,
) {
	f.Helper()
	conform.Fuzz(f, CachedSuite(), callable, cachedSubjects(newSubject), opts...)
}

// RunCachedBudgets applies the allocation ceilings of Cached and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunCachedBudgets[T Cached](
	t *testing.T,
	newSubject func(h *CachedHandle) T,
	opts ...CachedOption,
) {
	t.Helper()
	conform.RunBudgets(t, CachedSuite(), cachedSubjects(newSubject), opts...)
}

// BenchCachedSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchCachedSuite[T Cached](
	b conform.B,
	callable conform.Callable[Cached],
	newSubject func(h *CachedHandle) T,
	opts ...CachedOption,
) {
	b.Helper()
	conform.Bench(b, CachedSuite(), callable, cachedSubjects(newSubject), opts...)
}

// cachedSubjects returns the factory of the subjects of Cached that calls
// newSubject with the handle of each check.
func cachedSubjects[T Cached](newSubject func(h *CachedHandle) T) func(h *CachedHandle) Cached {
	return func(h *CachedHandle) Cached { return newSubject(h) }
}

// cachedGetCalls generates the calls of Cached.Get with the inputs that
// each case draws.
var cachedGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Cached) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Cached) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// cachedLoadCalls generates the calls of Cached.Load with the inputs that
// each case draws.
var cachedLoadCalls = prop.Composite(func(c *prop.Case) func(context.Context, Cached) error {
	id := c.Draw(prop.Of[string](), "id")
	return func(ctx context.Context, subject Cached) error {
		_, err := subject.Load(ctx, id)
		return err
	}
})

// cachedPeekCalls generates the calls of Cached.Peek with the inputs that
// each case draws.
var cachedPeekCalls = prop.Composite(func(c *prop.Case) func(context.Context, Cached) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Cached) error {
		_, err := subject.Peek(ctx, key)
		return err
	}
})

// cachedFetchCalls generates the calls of Cached.Fetch with the inputs
// that each case draws.
var cachedFetchCalls = prop.Composite(func(c *prop.Case) func(context.Context, Cached) error {
	key := c.Draw(prop.Of[int](), "key")
	return func(ctx context.Context, subject Cached) error {
		_, err := subject.Fetch(ctx, key)
		return err
	}
})

// cachedOverEmbedded is Cached under the name of the field that cachedOver
// embeds.
type cachedOverEmbedded = Cached

// cachedOver is a subject of Cached whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type cachedOver struct {
	cachedOverEmbedded
	spend func()
}

// Get calls spend and returns its zero results.
func (o cachedOver) Get(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// Load calls spend and returns its zero results.
func (o cachedOver) Load(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// Peek calls spend and returns its zero results.
func (o cachedOver) Peek(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// Fetch calls spend and returns its zero results.
func (o cachedOver) Fetch(context.Context, int) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// newCachedOver returns the planted defect of a ceiling of Cached, whose
// callables call spend.
func newCachedOver(_ *CachedHandle, spend func()) Cached {
	return cachedOver{spend: spend}
}

// cachedGetPanicsDefect is the planted defect of a Cached whose Get
// panics.
type cachedGetPanicsDefect struct {
	cachedOverEmbedded
}

// Get is the override of the planted defect.
func (*cachedGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	panic("the planted defect of a Cached whose Get panics")
}

// cachedGetPanics returns the planted defect of a Cached whose Get panics
// over a fresh subject.
func cachedGetPanics(h *CachedHandle) Cached {
	return &cachedGetPanicsDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedGetIgnoresCancelDefect is the planted defect of a Cached whose Get
// ignores the context.
type cachedGetIgnoresCancelDefect struct {
	cachedOverEmbedded
}

// Get is the override of the planted defect.
func (d *cachedGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.cachedOverEmbedded.Get(p0, p1)
}

// cachedGetIgnoresCancel returns the planted defect of a Cached whose Get
// ignores the context over a fresh subject.
func cachedGetIgnoresCancel(h *CachedHandle) Cached {
	return &cachedGetIgnoresCancelDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedGetValueBesideDefect is the planted defect of a Cached whose Get
// returns a value beside an error.
type cachedGetValueBesideDefect struct {
	cachedOverEmbedded
}

// Get is the override of the planted defect.
func (d *cachedGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	r0, r1 = d.cachedOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// cachedGetValueBeside returns the planted defect of a Cached whose Get
// returns a value beside an error over a fresh subject.
func cachedGetValueBeside(h *CachedHandle) Cached {
	return &cachedGetValueBesideDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedGetServesAStaleValueAfterThreeDefect is the planted defect of a
// Cached whose Get returns its first value again from its fourth call on.
type cachedGetServesAStaleValueAfterThreeDefect struct {
	cachedOverEmbedded
	calls int
	first *Item
}

// Get is the override of the planted defect.
func (d *cachedGetServesAStaleValueAfterThreeDefect) Get(
	p0 context.Context,
	p1 string,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.cachedOverEmbedded.Get(p0, p1)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// cachedGetServesAStaleValueAfterThree returns the planted defect of a
// Cached whose Get returns its first value again from its fourth call on
// over a fresh subject.
func cachedGetServesAStaleValueAfterThree(h *CachedHandle) Cached {
	return &cachedGetServesAStaleValueAfterThreeDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedLoadPanicsDefect is the planted defect of a Cached whose Load
// panics.
type cachedLoadPanicsDefect struct {
	cachedOverEmbedded
}

// Load is the override of the planted defect.
func (*cachedLoadPanicsDefect) Load(p0 context.Context, p1 string) (r0 Item, r1 error) {
	panic("the planted defect of a Cached whose Load panics")
}

// cachedLoadPanics returns the planted defect of a Cached whose Load
// panics over a fresh subject.
func cachedLoadPanics(h *CachedHandle) Cached {
	return &cachedLoadPanicsDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedLoadIgnoresCancelDefect is the planted defect of a Cached whose
// Load ignores the context.
type cachedLoadIgnoresCancelDefect struct {
	cachedOverEmbedded
}

// Load is the override of the planted defect.
func (d *cachedLoadIgnoresCancelDefect) Load(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.cachedOverEmbedded.Load(p0, p1)
}

// cachedLoadIgnoresCancel returns the planted defect of a Cached whose
// Load ignores the context over a fresh subject.
func cachedLoadIgnoresCancel(h *CachedHandle) Cached {
	return &cachedLoadIgnoresCancelDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedLoadValueBesideDefect is the planted defect of a Cached whose Load
// returns a value beside an error.
type cachedLoadValueBesideDefect struct {
	cachedOverEmbedded
}

// Load is the override of the planted defect.
func (d *cachedLoadValueBesideDefect) Load(p0 context.Context, p1 string) (r0 Item, r1 error) {
	r0, r1 = d.cachedOverEmbedded.Load(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// cachedLoadValueBeside returns the planted defect of a Cached whose Load
// returns a value beside an error over a fresh subject.
func cachedLoadValueBeside(h *CachedHandle) Cached {
	return &cachedLoadValueBesideDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedPeekPanicsDefect is the planted defect of a Cached whose Peek
// panics.
type cachedPeekPanicsDefect struct {
	cachedOverEmbedded
}

// Peek is the override of the planted defect.
func (*cachedPeekPanicsDefect) Peek(p0 context.Context, p1 string) (r0 Item, r1 error) {
	panic("the planted defect of a Cached whose Peek panics")
}

// cachedPeekPanics returns the planted defect of a Cached whose Peek
// panics over a fresh subject.
func cachedPeekPanics(h *CachedHandle) Cached {
	return &cachedPeekPanicsDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedPeekIgnoresCancelDefect is the planted defect of a Cached whose
// Peek ignores the context.
type cachedPeekIgnoresCancelDefect struct {
	cachedOverEmbedded
}

// Peek is the override of the planted defect.
func (d *cachedPeekIgnoresCancelDefect) Peek(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.cachedOverEmbedded.Peek(p0, p1)
}

// cachedPeekIgnoresCancel returns the planted defect of a Cached whose
// Peek ignores the context over a fresh subject.
func cachedPeekIgnoresCancel(h *CachedHandle) Cached {
	return &cachedPeekIgnoresCancelDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedPeekValueBesideDefect is the planted defect of a Cached whose Peek
// returns a value beside an error.
type cachedPeekValueBesideDefect struct {
	cachedOverEmbedded
}

// Peek is the override of the planted defect.
func (d *cachedPeekValueBesideDefect) Peek(p0 context.Context, p1 string) (r0 Item, r1 error) {
	r0, r1 = d.cachedOverEmbedded.Peek(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// cachedPeekValueBeside returns the planted defect of a Cached whose Peek
// returns a value beside an error over a fresh subject.
func cachedPeekValueBeside(h *CachedHandle) Cached {
	return &cachedPeekValueBesideDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedPeekServesAStaleValueAfterThreeDefect is the planted defect of a
// Cached whose Peek returns its first value again from its fourth call on.
type cachedPeekServesAStaleValueAfterThreeDefect struct {
	cachedOverEmbedded
	calls int
	first *Item
}

// Peek is the override of the planted defect.
func (d *cachedPeekServesAStaleValueAfterThreeDefect) Peek(
	p0 context.Context,
	p1 string,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.cachedOverEmbedded.Peek(p0, p1)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// cachedPeekServesAStaleValueAfterThree returns the planted defect of a
// Cached whose Peek returns its first value again from its fourth call on
// over a fresh subject.
func cachedPeekServesAStaleValueAfterThree(h *CachedHandle) Cached {
	return &cachedPeekServesAStaleValueAfterThreeDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedFetchPanicsDefect is the planted defect of a Cached whose Fetch
// panics.
type cachedFetchPanicsDefect struct {
	cachedOverEmbedded
}

// Fetch is the override of the planted defect.
func (*cachedFetchPanicsDefect) Fetch(p0 context.Context, p1 int) (r0 Item, r1 error) {
	panic("the planted defect of a Cached whose Fetch panics")
}

// cachedFetchPanics returns the planted defect of a Cached whose Fetch
// panics over a fresh subject.
func cachedFetchPanics(h *CachedHandle) Cached {
	return &cachedFetchPanicsDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedFetchIgnoresCancelDefect is the planted defect of a Cached whose
// Fetch ignores the context.
type cachedFetchIgnoresCancelDefect struct {
	cachedOverEmbedded
}

// Fetch is the override of the planted defect.
func (d *cachedFetchIgnoresCancelDefect) Fetch(p0 context.Context, p1 int) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.cachedOverEmbedded.Fetch(p0, p1)
}

// cachedFetchIgnoresCancel returns the planted defect of a Cached whose
// Fetch ignores the context over a fresh subject.
func cachedFetchIgnoresCancel(h *CachedHandle) Cached {
	return &cachedFetchIgnoresCancelDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedFetchValueBesideDefect is the planted defect of a Cached whose
// Fetch returns a value beside an error.
type cachedFetchValueBesideDefect struct {
	cachedOverEmbedded
}

// Fetch is the override of the planted defect.
func (d *cachedFetchValueBesideDefect) Fetch(p0 context.Context, p1 int) (r0 Item, r1 error) {
	r0, r1 = d.cachedOverEmbedded.Fetch(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// cachedFetchValueBeside returns the planted defect of a Cached whose
// Fetch returns a value beside an error over a fresh subject.
func cachedFetchValueBeside(h *CachedHandle) Cached {
	return &cachedFetchValueBesideDefect{cachedOverEmbedded: h.Fresh()}
}

// cachedActions returns the actions of the generated sequences of Cached
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func cachedActions(c *prop.Case, subject func() Cached) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	keysInt := c.Draw(
		prop.List(prop.Of[int](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.int",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Get",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "get.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Get", args, args[0])
				values, err := cachedGetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Load",
			Input: func(c *prop.Case, _ struct{}) any {
				id := c.Draw(prop.SampledFrom(keysString...), "load.id")
				return []any{id}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Load", args, args[0])
				values, err := cachedLoadCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Peek",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "peek.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Peek", args, args[0])
				values, err := cachedPeekCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Fetch",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysInt...), "fetch.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Fetch", args, args[0])
				values, err := cachedFetchCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// cachedGetCall calls Get on subject with the arguments args, and returns
// the values of the call beside its error.
func cachedGetCall(ctx context.Context, subject Cached, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Get(ctx, a0)
	return []any{r0}, err
}

// cachedLoadCall calls Load on subject with the arguments args, and
// returns the values of the call beside its error.
func cachedLoadCall(ctx context.Context, subject Cached, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Load(ctx, a0)
	return []any{r0}, err
}

// cachedPeekCall calls Peek on subject with the arguments args, and
// returns the values of the call beside its error.
func cachedPeekCall(ctx context.Context, subject Cached, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Peek(ctx, a0)
	return []any{r0}, err
}

// cachedFetchCall calls Fetch on subject with the arguments args, and
// returns the values of the call beside its error.
func cachedFetchCall(ctx context.Context, subject Cached, args []any) ([]any, error) {
	a0, _ := args[0].(int)
	r0, err := subject.Fetch(ctx, a0)
	return []any{r0}, err
}

// The derived checks of Index:
//
// Index.Put/smoke: Put returns without a panic for derived inputs.
//
// Index.Put/smoke#generated: Put returns without a panic for derived
// inputs.
//
// Index.Put/cancel: Put returns an error for a cancelled context.
//
// Index.Put/cancel#generated: Put returns an error for a cancelled
// context.
//
// Index.Get/smoke: Get returns without a panic for derived inputs.
//
// Index.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Index.Get/cancel: Get returns an error for a cancelled context.
//
// Index.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Index.Get/zero-on-error: Get returns the zero value beside an error.
//
// Index.Get/zero-on-error#generated: Get returns the zero value beside an
// error.
//
// Index.Get/read-after-write: Get returns what Put stored before it.
//
// Index.Get/read-after-write#generated: Get returns what Put stored before
// it.
//
// Index.Get/reads-each-last-write: Get returns what Put wrote right before
// it, over any sequence of calls.
//
// Index.Find/smoke: Find returns without a panic for derived inputs.
//
// Index.Find/smoke#generated: Find returns without a panic for derived
// inputs.
//
// Index.Find/cancel: Find returns an error for a cancelled context.
//
// Index.Find/cancel#generated: Find returns an error for a cancelled
// context.
//
// Index.Find/read-after-write: Find returns what Put stored before it.
//
// Index.Find/read-after-write#generated: Find returns what Put stored
// before it.
//
// Index.Find/reads-each-last-write: Find returns what Put wrote right
// before it, over any sequence of calls.
//
// Index.ByRev/smoke: ByRev returns without a panic for derived inputs.
//
// Index.ByRev/smoke#generated: ByRev returns without a panic for derived
// inputs.
//
// Index.ByRev/cancel: ByRev returns an error for a cancelled context.
//
// Index.ByRev/cancel#generated: ByRev returns an error for a cancelled
// context.
//
// Index.ByRev/zero-on-error: ByRev returns the zero value beside an error.
//
// Index.ByRev/zero-on-error#generated: ByRev returns the zero value beside
// an error.
//
// Index.ByRev/read-after-write: ByRev returns what Put stored before it.
//
// Index.ByRev/read-after-write#generated: ByRev returns what Put stored
// before it.
//
// Index.ByRev/reads-each-last-write: ByRev returns what Put wrote right
// before it, over any sequence of calls.
//
// Index.Lookup/smoke: Lookup returns without a panic for derived inputs.
//
// Index.Lookup/smoke#generated: Lookup returns without a panic for derived
// inputs.
//
// Index.Lookup/cancel: Lookup returns an error for a cancelled context.
//
// Index.Lookup/cancel#generated: Lookup returns an error for a cancelled
// context.
//
// Index.Lookup/read-after-write: Lookup returns what Put stored before it.
//
// Index.Lookup/read-after-write#generated: Lookup returns what Put stored
// before it.
//
// Index.Lookup/reads-each-last-write: Lookup returns what Put wrote right
// before it, over any sequence of calls.

// The derived checks of Index that witness their laws at the shallow end:
//
// Index.Put/smoke
//
// Index.Get/smoke
//
// Index.Get/read-after-write
//
// Index.Get/read-after-write#generated
//
// Index.Find/smoke
//
// Index.Find/read-after-write
//
// Index.Find/read-after-write#generated
//
// Index.ByRev/smoke
//
// Index.ByRev/read-after-write
//
// Index.ByRev/read-after-write#generated
//
// Index.Lookup/smoke
//
// Index.Lookup/read-after-write
//
// Index.Lookup/read-after-write#generated

// The unclaimed classifications of Index:
//
// The shape:writer of Index.Put has the stance silent. A write is an input
// of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Index writes.
type (
	IndexHandle = conform.Handle[Index]
	IndexCheck  = conform.Check[Index]
	IndexOption = conform.Option[Index]
)

// IndexChecks contains the name of each callable of Index, with the
// identities of its checks.
var IndexChecks = struct {
	Put    indexPut
	Get    indexGet
	Find   indexFind
	ByRev  indexByRev
	Lookup indexLookup
}{
	Put:    indexPut{"Index.Put"},
	Get:    indexGet{"Index.Get"},
	Find:   indexFind{"Index.Find"},
	ByRev:  indexByRev{"Index.ByRev"},
	Lookup: indexLookup{"Index.Lookup"},
}

// indexPut is the name of Index.Put, with the identities of its checks.
type indexPut struct {
	conform.Callable[Index]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c indexPut) Smoke() conform.ID[Index] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Put returns
// without a panic for derived inputs.
func (c indexPut) SmokeGenerated() conform.ID[Index] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c indexPut) Cancel() conform.ID[Index] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Put returns an
// error for a cancelled context.
func (c indexPut) CancelGenerated() conform.ID[Index] {
	return c.Check("cancel#generated")
}

// indexGet is the name of Index.Get, with the identities of its checks.
type indexGet struct {
	conform.Callable[Index]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c indexGet) Smoke() conform.ID[Index] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c indexGet) SmokeGenerated() conform.ID[Index] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c indexGet) Cancel() conform.ID[Index] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c indexGet) CancelGenerated() conform.ID[Index] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c indexGet) ZeroOnError() conform.ID[Index] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c indexGet) ZeroOnErrorGenerated() conform.ID[Index] {
	return c.Check("zero-on-error#generated")
}

// ReadAfterWrite returns the identity of the check that Get returns what
// Put stored before it.
func (c indexGet) ReadAfterWrite() conform.ID[Index] {
	return c.Check("read-after-write")
}

// ReadAfterWriteGenerated returns the identity of the check that Get
// returns what Put stored before it.
func (c indexGet) ReadAfterWriteGenerated() conform.ID[Index] {
	return c.Check("read-after-write#generated")
}

// ReadsEachLastWrite returns the identity of the check that Get returns
// what Put wrote right before it, over any sequence of calls.
func (c indexGet) ReadsEachLastWrite() conform.ID[Index] {
	return c.Check("reads-each-last-write")
}

// indexFind is the name of Index.Find, with the identities of its checks.
type indexFind struct {
	conform.Callable[Index]
}

// Smoke returns the identity of the check that Find returns without a
// panic for derived inputs.
func (c indexFind) Smoke() conform.ID[Index] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Find returns
// without a panic for derived inputs.
func (c indexFind) SmokeGenerated() conform.ID[Index] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Find returns an error for
// a cancelled context.
func (c indexFind) Cancel() conform.ID[Index] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Find returns an
// error for a cancelled context.
func (c indexFind) CancelGenerated() conform.ID[Index] {
	return c.Check("cancel#generated")
}

// ReadAfterWrite returns the identity of the check that Find returns what
// Put stored before it.
func (c indexFind) ReadAfterWrite() conform.ID[Index] {
	return c.Check("read-after-write")
}

// ReadAfterWriteGenerated returns the identity of the check that Find
// returns what Put stored before it.
func (c indexFind) ReadAfterWriteGenerated() conform.ID[Index] {
	return c.Check("read-after-write#generated")
}

// ReadsEachLastWrite returns the identity of the check that Find returns
// what Put wrote right before it, over any sequence of calls.
func (c indexFind) ReadsEachLastWrite() conform.ID[Index] {
	return c.Check("reads-each-last-write")
}

// indexByRev is the name of Index.ByRev, with the identities of its
// checks.
type indexByRev struct {
	conform.Callable[Index]
}

// Smoke returns the identity of the check that ByRev returns without a
// panic for derived inputs.
func (c indexByRev) Smoke() conform.ID[Index] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that ByRev returns
// without a panic for derived inputs.
func (c indexByRev) SmokeGenerated() conform.ID[Index] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that ByRev returns an error for
// a cancelled context.
func (c indexByRev) Cancel() conform.ID[Index] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that ByRev returns an
// error for a cancelled context.
func (c indexByRev) CancelGenerated() conform.ID[Index] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that ByRev returns the
// zero value beside an error.
func (c indexByRev) ZeroOnError() conform.ID[Index] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that ByRev
// returns the zero value beside an error.
func (c indexByRev) ZeroOnErrorGenerated() conform.ID[Index] {
	return c.Check("zero-on-error#generated")
}

// ReadAfterWrite returns the identity of the check that ByRev returns what
// Put stored before it.
func (c indexByRev) ReadAfterWrite() conform.ID[Index] {
	return c.Check("read-after-write")
}

// ReadAfterWriteGenerated returns the identity of the check that ByRev
// returns what Put stored before it.
func (c indexByRev) ReadAfterWriteGenerated() conform.ID[Index] {
	return c.Check("read-after-write#generated")
}

// ReadsEachLastWrite returns the identity of the check that ByRev returns
// what Put wrote right before it, over any sequence of calls.
func (c indexByRev) ReadsEachLastWrite() conform.ID[Index] {
	return c.Check("reads-each-last-write")
}

// indexLookup is the name of Index.Lookup, with the identities of its
// checks.
type indexLookup struct {
	conform.Callable[Index]
}

// Smoke returns the identity of the check that Lookup returns without a
// panic for derived inputs.
func (c indexLookup) Smoke() conform.ID[Index] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Lookup returns
// without a panic for derived inputs.
func (c indexLookup) SmokeGenerated() conform.ID[Index] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Lookup returns an error
// for a cancelled context.
func (c indexLookup) Cancel() conform.ID[Index] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Lookup returns an
// error for a cancelled context.
func (c indexLookup) CancelGenerated() conform.ID[Index] {
	return c.Check("cancel#generated")
}

// ReadAfterWrite returns the identity of the check that Lookup returns
// what Put stored before it.
func (c indexLookup) ReadAfterWrite() conform.ID[Index] {
	return c.Check("read-after-write")
}

// ReadAfterWriteGenerated returns the identity of the check that Lookup
// returns what Put stored before it.
func (c indexLookup) ReadAfterWriteGenerated() conform.ID[Index] {
	return c.Check("read-after-write#generated")
}

// ReadsEachLastWrite returns the identity of the check that Lookup returns
// what Put wrote right before it, over any sequence of calls.
func (c indexLookup) ReadsEachLastWrite() conform.ID[Index] {
	return c.Check("reads-each-last-write")
}

// IndexSuite returns the checks of Index, and the calls and the budgets of
// its callables.
func IndexSuite() conform.Suite[Index] {
	return conform.Suite[Index]{
		Declaration: "Index",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Index]{
			{
				ID:          IndexChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: indexPutPanics,
			},
			{
				ID:          IndexChecks.Put.SmokeGenerated(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnIt)
					})
				},
				Defect: indexPutPanics,
			},
			{
				ID:          IndexChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *IndexHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: indexPutIgnoresCancel,
			},
			{
				ID:          IndexChecks.Put.CancelGenerated(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnIt)
					})
				},
				Defect: indexPutIgnoresCancel,
			},
			{
				ID:          IndexChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: indexGetPanics,
			},
			{
				ID:          IndexChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: indexGetPanics,
			},
			{
				ID:          IndexChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *IndexHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: indexGetIgnoresCancel,
			},
			{
				ID:          IndexChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: indexGetIgnoresCancel,
			},
			{
				ID:          IndexChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *IndexHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Item{})
				},
				Defect: indexGetValueBeside,
			},
			{
				ID:          IndexChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, Item{})
				},
				Defect: indexGetValueBeside,
			},
			{
				ID:          IndexChecks.Get.ReadAfterWrite(),
				Statement:   "Get returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, (Item{Key: "Key-a"}).Key)
					}, func(ctx context.Context) (Item, error) {
						w := Item{Key: "Key-a"}
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexGetReadsAnotherValue,
			},
			{
				ID:          IndexChecks.Get.ReadAfterWriteGenerated(),
				Statement:   "Get returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, drawnIt.Key)
					}, func(ctx context.Context) (Item, error) {
						w := drawnIt
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexGetReadsAnotherValue,
			},
			{
				ID:          IndexChecks.Get.ReadsEachLastWrite(),
				Statement:   "Get returns what Put wrote right before it, over any sequence of calls",
				Claim:       spec.ReadsEachLastWrite,
				Probe:       spec.ReadsReturnTheLastWrite,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Index]{IndexChecks.Get.ReadAfterWrite()},
				Property: func(h *IndexHandle, c *prop.Case) {
					probe.ReadsReturnTheLastWrite(
						h,
						c,
						indexActions,
						"Put",
						func(ctx context.Context, args []any) conform.Output {
							w, _ := args[0].(Item)
							values, err := indexGetCall(ctx, h.Subject(), []any{w.Key})
							return conform.OutputOf(err, values...)
						},
						func(args []any) Item {
							v, _ := args[0].(Item)
							return v
						},
					)
				},
				Defect: indexGetReadsAStaleValueAfterThree,
			},
			{
				ID:          IndexChecks.Find.Smoke(),
				Statement:   "Find returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().
							Find(ctx, string("key-a"), func(Item) bool { return false })
						return err
					})
				},
				Defect: indexFindPanics,
			},
			{
				ID:          IndexChecks.Find.SmokeGenerated(),
				Statement:   "Find returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					keepResults := c.Draw(prop.List(prop.Of[bool]()), "keep")
					drawnKeep := func() func(Item) bool {
						var keepAt atomic.Int64
						return func(Item) (r0 bool) {
							i := int(keepAt.Add(1) - 1)
							if i >= len(keepResults) {
								return r0
							}
							return keepResults[i]
						}
					}
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Find(ctx, drawnKey, drawnKeep())
						return err
					})
				},
				Defect: indexFindPanics,
			},
			{
				ID:          IndexChecks.Find.Cancel(),
				Statement:   "Find returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *IndexHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().
							Find(ctx, string("key-a"), func(Item) bool { return false })
						return err
					})
				},
				Defect: indexFindIgnoresCancel,
			},
			{
				ID:          IndexChecks.Find.CancelGenerated(),
				Statement:   "Find returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					keepResults := c.Draw(prop.List(prop.Of[bool]()), "keep")
					drawnKeep := func() func(Item) bool {
						var keepAt atomic.Int64
						return func(Item) (r0 bool) {
							i := int(keepAt.Add(1) - 1)
							if i >= len(keepResults) {
								return r0
							}
							return keepResults[i]
						}
					}
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Find(ctx, drawnKey, drawnKeep())
						return err
					})
				},
				Defect: indexFindIgnoresCancel,
			},
			{
				ID:          IndexChecks.Find.ReadAfterWrite(),
				Statement:   "Find returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().
							Find(ctx, (Item{Key: "Key-a"}).Key, func(Item) bool { return false })
					}, func(ctx context.Context) (Item, error) {
						w := Item{Key: "Key-a"}
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexFindReadsAnotherValue,
			},
			{
				ID:          IndexChecks.Find.ReadAfterWriteGenerated(),
				Statement:   "Find returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					keepResults := c.Draw(prop.List(prop.Of[bool]()), "keep")
					drawnKeep := func() func(Item) bool {
						var keepAt atomic.Int64
						return func(Item) (r0 bool) {
							i := int(keepAt.Add(1) - 1)
							if i >= len(keepResults) {
								return r0
							}
							return keepResults[i]
						}
					}
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Find(ctx, drawnIt.Key, drawnKeep())
					}, func(ctx context.Context) (Item, error) {
						w := drawnIt
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexFindReadsAnotherValue,
			},
			{
				ID:          IndexChecks.Find.ReadsEachLastWrite(),
				Statement:   "Find returns what Put wrote right before it, over any sequence of calls",
				Claim:       spec.ReadsEachLastWrite,
				Probe:       spec.ReadsReturnTheLastWrite,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Index]{IndexChecks.Find.ReadAfterWrite()},
				Property: func(h *IndexHandle, c *prop.Case) {
					probe.ReadsReturnTheLastWrite(
						h,
						c,
						indexActions,
						"Put",
						func(ctx context.Context, args []any) conform.Output {
							w, _ := args[0].(Item)
							values, err := indexFindCall(
								ctx,
								h.Subject(),
								[]any{w.Key, func(Item) bool { return false }},
							)
							return conform.OutputOf(err, values...)
						},
						func(args []any) Item {
							v, _ := args[0].(Item)
							return v
						},
					)
				},
				Defect: indexFindReadsAStaleValueAfterThree,
			},
			{
				ID:          IndexChecks.ByRev.Smoke(),
				Statement:   "ByRev returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().ByRev(ctx, int(42))
						return err
					})
				},
				Defect: indexByRevPanics,
			},
			{
				ID:          IndexChecks.ByRev.SmokeGenerated(),
				Statement:   "ByRev returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnRev := c.Draw(prop.Of[int](), "rev")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().ByRev(ctx, drawnRev)
						return err
					})
				},
				Defect: indexByRevPanics,
			},
			{
				ID:          IndexChecks.ByRev.Cancel(),
				Statement:   "ByRev returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *IndexHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().ByRev(ctx, int(42))
						return err
					})
				},
				Defect: indexByRevIgnoresCancel,
			},
			{
				ID:          IndexChecks.ByRev.CancelGenerated(),
				Statement:   "ByRev returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnRev := c.Draw(prop.Of[int](), "rev")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().ByRev(ctx, drawnRev)
						return err
					})
				},
				Defect: indexByRevIgnoresCancel,
			},
			{
				ID:          IndexChecks.ByRev.ZeroOnError(),
				Statement:   "ByRev returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *IndexHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().ByRev(ctx, int(42))
					}, Item{})
				},
				Defect: indexByRevValueBeside,
			},
			{
				ID:          IndexChecks.ByRev.ZeroOnErrorGenerated(),
				Statement:   "ByRev returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnRev := c.Draw(prop.Of[int](), "rev")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().ByRev(ctx, drawnRev)
					}, Item{})
				},
				Defect: indexByRevValueBeside,
			},
			{
				ID:          IndexChecks.ByRev.ReadAfterWrite(),
				Statement:   "ByRev returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().ByRev(ctx, (Item{Key: "Key-a"}).Rev)
					}, func(ctx context.Context) (Item, error) {
						w := Item{Key: "Key-a"}
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexByRevReadsAnotherValue,
			},
			{
				ID:          IndexChecks.ByRev.ReadAfterWriteGenerated(),
				Statement:   "ByRev returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().ByRev(ctx, drawnIt.Rev)
					}, func(ctx context.Context) (Item, error) {
						w := drawnIt
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexByRevReadsAnotherValue,
			},
			{
				ID:          IndexChecks.ByRev.ReadsEachLastWrite(),
				Statement:   "ByRev returns what Put wrote right before it, over any sequence of calls",
				Claim:       spec.ReadsEachLastWrite,
				Probe:       spec.ReadsReturnTheLastWrite,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Index]{IndexChecks.ByRev.ReadAfterWrite()},
				Property: func(h *IndexHandle, c *prop.Case) {
					probe.ReadsReturnTheLastWrite(
						h,
						c,
						indexActions,
						"Put",
						func(ctx context.Context, args []any) conform.Output {
							w, _ := args[0].(Item)
							values, err := indexByRevCall(ctx, h.Subject(), []any{w.Rev})
							return conform.OutputOf(err, values...)
						},
						func(args []any) Item {
							v, _ := args[0].(Item)
							return v
						},
					)
				},
				Defect: indexByRevReadsAStaleValueAfterThree,
			},
			{
				ID:          IndexChecks.Lookup.Smoke(),
				Statement:   "Lookup returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Lookup(ctx, string("key-a"), int(42))
						return err
					})
				},
				Defect: indexLookupPanics,
			},
			{
				ID:          IndexChecks.Lookup.SmokeGenerated(),
				Statement:   "Lookup returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnRev := c.Draw(prop.Of[int](), "rev")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Lookup(ctx, drawnKey, drawnRev)
						return err
					})
				},
				Defect: indexLookupPanics,
			},
			{
				ID:          IndexChecks.Lookup.Cancel(),
				Statement:   "Lookup returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *IndexHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Lookup(ctx, string("key-a"), int(42))
						return err
					})
				},
				Defect: indexLookupIgnoresCancel,
			},
			{
				ID:          IndexChecks.Lookup.CancelGenerated(),
				Statement:   "Lookup returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					drawnRev := c.Draw(prop.Of[int](), "rev")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Lookup(ctx, drawnKey, drawnRev)
						return err
					})
				},
				Defect: indexLookupIgnoresCancel,
			},
			{
				ID:          IndexChecks.Lookup.ReadAfterWrite(),
				Statement:   "Lookup returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *IndexHandle) {
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().
							Lookup(ctx, (Item{Key: "Key-a"}).Key, (Item{Key: "Key-a"}).Rev)
					}, func(ctx context.Context) (Item, error) {
						w := Item{Key: "Key-a"}
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexLookupReadsAnotherValue,
			},
			{
				ID:          IndexChecks.Lookup.ReadAfterWriteGenerated(),
				Statement:   "Lookup returns what Put stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *IndexHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.WriteThenRead(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Lookup(ctx, drawnIt.Key, drawnIt.Rev)
					}, func(ctx context.Context) (Item, error) {
						w := drawnIt
						err := h.Subject().Put(ctx, w)
						return w, err
					})
				},
				Defect: indexLookupReadsAnotherValue,
			},
			{
				ID:          IndexChecks.Lookup.ReadsEachLastWrite(),
				Statement:   "Lookup returns what Put wrote right before it, over any sequence of calls",
				Claim:       spec.ReadsEachLastWrite,
				Probe:       spec.ReadsReturnTheLastWrite,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Index]{IndexChecks.Lookup.ReadAfterWrite()},
				Property: func(h *IndexHandle, c *prop.Case) {
					probe.ReadsReturnTheLastWrite(
						h,
						c,
						indexActions,
						"Put",
						func(ctx context.Context, args []any) conform.Output {
							w, _ := args[0].(Item)
							values, err := indexLookupCall(ctx, h.Subject(), []any{w.Key, w.Rev})
							return conform.OutputOf(err, values...)
						},
						func(args []any) Item {
							v, _ := args[0].(Item)
							return v
						},
					)
				},
				Defect: indexLookupReadsAStaleValueAfterThree,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Index.Put",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Index]{
			{
				Calls:    indexPutCalls,
				Over:     newIndexOver,
				Callable: IndexChecks.Put.Callable,
			},
			{
				Calls:    indexGetCalls,
				Over:     newIndexOver,
				Callable: IndexChecks.Get.Callable,
			},
			{
				Calls:    indexFindCalls,
				Over:     newIndexOver,
				Callable: IndexChecks.Find.Callable,
			},
			{
				Calls:    indexByRevCalls,
				Over:     newIndexOver,
				Callable: IndexChecks.ByRev.Callable,
			},
			{
				Calls:    indexLookupCalls,
				Over:     newIndexOver,
				Callable: IndexChecks.Lookup.Callable,
			},
		},
	}
}

// RunIndexSuite applies the checks of Index, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunIndexSuite[T Index](t *testing.T, newSubject func() T, opts ...IndexOption) {
	t.Helper()
	RunIndexSuiteFromHandle(t, func(*IndexHandle) T { return newSubject() }, opts...)
}

// RunIndexSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunIndexSuiteFromHandle[T Index](
	t *testing.T,
	newSubject func(h *IndexHandle) T,
	opts ...IndexOption,
) {
	t.Helper()
	conform.Run(t, IndexSuite(), indexSubjects(newSubject), opts...)
}

// FuzzIndexSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzIndexSuite[T Index](
	f *testing.F,
	callable conform.Callable[Index],
	newSubject func(h *IndexHandle) T,
	opts ...IndexOption,
) {
	f.Helper()
	conform.Fuzz(f, IndexSuite(), callable, indexSubjects(newSubject), opts...)
}

// RunIndexBudgets applies the allocation ceilings of Index and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunIndexBudgets[T Index](
	t *testing.T,
	newSubject func(h *IndexHandle) T,
	opts ...IndexOption,
) {
	t.Helper()
	conform.RunBudgets(t, IndexSuite(), indexSubjects(newSubject), opts...)
}

// BenchIndexSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchIndexSuite[T Index](
	b conform.B,
	callable conform.Callable[Index],
	newSubject func(h *IndexHandle) T,
	opts ...IndexOption,
) {
	b.Helper()
	conform.Bench(b, IndexSuite(), callable, indexSubjects(newSubject), opts...)
}

// indexSubjects returns the factory of the subjects of Index that calls
// newSubject with the handle of each check.
func indexSubjects[T Index](newSubject func(h *IndexHandle) T) func(h *IndexHandle) Index {
	return func(h *IndexHandle) Index { return newSubject(h) }
}

// indexPutCalls generates the calls of Index.Put with the inputs that each
// case draws.
var indexPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Index) error {
	it := c.Draw(prop.Of[Item](), "it")
	return func(ctx context.Context, subject Index) error {
		return subject.Put(ctx, it)
	}
})

// indexGetCalls generates the calls of Index.Get with the inputs that each
// case draws.
var indexGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Index) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Index) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// indexFindCalls generates the calls of Index.Find with the inputs that
// each case draws.
var indexFindCalls = prop.Composite(func(c *prop.Case) func(context.Context, Index) error {
	key := c.Draw(prop.Of[string](), "key")
	keepResults := c.Draw(prop.List(prop.Of[bool]()), "keep")
	var keepAt atomic.Int64
	keep := func(Item) (r0 bool) {
		i := int(keepAt.Add(1) - 1)
		if i >= len(keepResults) {
			return r0
		}
		return keepResults[i]
	}
	return func(ctx context.Context, subject Index) error {
		_, err := subject.Find(ctx, key, keep)
		return err
	}
})

// indexByRevCalls generates the calls of Index.ByRev with the inputs that
// each case draws.
var indexByRevCalls = prop.Composite(func(c *prop.Case) func(context.Context, Index) error {
	rev := c.Draw(prop.Of[int](), "rev")
	return func(ctx context.Context, subject Index) error {
		_, err := subject.ByRev(ctx, rev)
		return err
	}
})

// indexLookupCalls generates the calls of Index.Lookup with the inputs
// that each case draws.
var indexLookupCalls = prop.Composite(func(c *prop.Case) func(context.Context, Index) error {
	key := c.Draw(prop.Of[string](), "key")
	rev := c.Draw(prop.Of[int](), "rev")
	return func(ctx context.Context, subject Index) error {
		_, err := subject.Lookup(ctx, key, rev)
		return err
	}
})

// indexOverEmbedded is Index under the name of the field that indexOver
// embeds.
type indexOverEmbedded = Index

// indexOver is a subject of Index whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type indexOver struct {
	indexOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o indexOver) Put(context.Context, Item) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o indexOver) Get(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// Find calls spend and returns its zero results.
func (o indexOver) Find(context.Context, string, func(Item) bool) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// ByRev calls spend and returns its zero results.
func (o indexOver) ByRev(context.Context, int) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// Lookup calls spend and returns its zero results.
func (o indexOver) Lookup(context.Context, string, int) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// newIndexOver returns the planted defect of a ceiling of Index, whose
// callables call spend.
func newIndexOver(_ *IndexHandle, spend func()) Index {
	return indexOver{spend: spend}
}

// indexPutPanicsDefect is the planted defect of a Index whose Put panics.
type indexPutPanicsDefect struct {
	indexOverEmbedded
}

// Put is the override of the planted defect.
func (*indexPutPanicsDefect) Put(p0 context.Context, p1 Item) (r0 error) {
	panic("the planted defect of a Index whose Put panics")
}

// indexPutPanics returns the planted defect of a Index whose Put panics
// over a fresh subject.
func indexPutPanics(h *IndexHandle) Index {
	return &indexPutPanicsDefect{indexOverEmbedded: h.Fresh()}
}

// indexPutIgnoresCancelDefect is the planted defect of a Index whose Put
// ignores the context.
type indexPutIgnoresCancelDefect struct {
	indexOverEmbedded
}

// Put is the override of the planted defect.
func (d *indexPutIgnoresCancelDefect) Put(p0 context.Context, p1 Item) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.indexOverEmbedded.Put(p0, p1)
}

// indexPutIgnoresCancel returns the planted defect of a Index whose Put
// ignores the context over a fresh subject.
func indexPutIgnoresCancel(h *IndexHandle) Index {
	return &indexPutIgnoresCancelDefect{indexOverEmbedded: h.Fresh()}
}

// indexGetPanicsDefect is the planted defect of a Index whose Get panics.
type indexGetPanicsDefect struct {
	indexOverEmbedded
}

// Get is the override of the planted defect.
func (*indexGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	panic("the planted defect of a Index whose Get panics")
}

// indexGetPanics returns the planted defect of a Index whose Get panics
// over a fresh subject.
func indexGetPanics(h *IndexHandle) Index {
	return &indexGetPanicsDefect{indexOverEmbedded: h.Fresh()}
}

// indexGetIgnoresCancelDefect is the planted defect of a Index whose Get
// ignores the context.
type indexGetIgnoresCancelDefect struct {
	indexOverEmbedded
}

// Get is the override of the planted defect.
func (d *indexGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.indexOverEmbedded.Get(p0, p1)
}

// indexGetIgnoresCancel returns the planted defect of a Index whose Get
// ignores the context over a fresh subject.
func indexGetIgnoresCancel(h *IndexHandle) Index {
	return &indexGetIgnoresCancelDefect{indexOverEmbedded: h.Fresh()}
}

// indexGetValueBesideDefect is the planted defect of a Index whose Get
// returns a value beside an error.
type indexGetValueBesideDefect struct {
	indexOverEmbedded
}

// Get is the override of the planted defect.
func (d *indexGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	r0, r1 = d.indexOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// indexGetValueBeside returns the planted defect of a Index whose Get
// returns a value beside an error over a fresh subject.
func indexGetValueBeside(h *IndexHandle) Index {
	return &indexGetValueBesideDefect{indexOverEmbedded: h.Fresh()}
}

// indexGetReadsAnotherValueDefect is the planted defect of a Index whose
// Get returns a value that Put did not write.
type indexGetReadsAnotherValueDefect struct {
	indexOverEmbedded
}

// Get is the override of the planted defect.
func (d *indexGetReadsAnotherValueDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.indexOverEmbedded.Get(p0, p1)
	if r1 == nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// indexGetReadsAnotherValue returns the planted defect of a Index whose
// Get returns a value that Put did not write over a fresh subject.
func indexGetReadsAnotherValue(h *IndexHandle) Index {
	return &indexGetReadsAnotherValueDefect{indexOverEmbedded: h.Fresh()}
}

// indexGetReadsAStaleValueAfterThreeDefect is the planted defect of a
// Index whose Get returns its first value again from its fourth call on.
type indexGetReadsAStaleValueAfterThreeDefect struct {
	indexOverEmbedded
	calls int
	first *Item
}

// Get is the override of the planted defect.
func (d *indexGetReadsAStaleValueAfterThreeDefect) Get(
	p0 context.Context,
	p1 string,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.indexOverEmbedded.Get(p0, p1)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// indexGetReadsAStaleValueAfterThree returns the planted defect of a Index
// whose Get returns its first value again from its fourth call on over a
// fresh subject.
func indexGetReadsAStaleValueAfterThree(h *IndexHandle) Index {
	return &indexGetReadsAStaleValueAfterThreeDefect{indexOverEmbedded: h.Fresh()}
}

// indexFindPanicsDefect is the planted defect of a Index whose Find
// panics.
type indexFindPanicsDefect struct {
	indexOverEmbedded
}

// Find is the override of the planted defect.
func (*indexFindPanicsDefect) Find(
	p0 context.Context,
	p1 string,
	p2 func(Item) bool,
) (r0 Item, r1 error) {
	panic("the planted defect of a Index whose Find panics")
}

// indexFindPanics returns the planted defect of a Index whose Find panics
// over a fresh subject.
func indexFindPanics(h *IndexHandle) Index {
	return &indexFindPanicsDefect{indexOverEmbedded: h.Fresh()}
}

// indexFindIgnoresCancelDefect is the planted defect of a Index whose Find
// ignores the context.
type indexFindIgnoresCancelDefect struct {
	indexOverEmbedded
}

// Find is the override of the planted defect.
func (d *indexFindIgnoresCancelDefect) Find(
	p0 context.Context,
	p1 string,
	p2 func(Item) bool,
) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.indexOverEmbedded.Find(p0, p1, p2)
}

// indexFindIgnoresCancel returns the planted defect of a Index whose Find
// ignores the context over a fresh subject.
func indexFindIgnoresCancel(h *IndexHandle) Index {
	return &indexFindIgnoresCancelDefect{indexOverEmbedded: h.Fresh()}
}

// indexFindReadsAnotherValueDefect is the planted defect of a Index whose
// Find returns a value that Put did not write.
type indexFindReadsAnotherValueDefect struct {
	indexOverEmbedded
}

// Find is the override of the planted defect.
func (d *indexFindReadsAnotherValueDefect) Find(
	p0 context.Context,
	p1 string,
	p2 func(Item) bool,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.indexOverEmbedded.Find(p0, p1, p2)
	if r1 == nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// indexFindReadsAnotherValue returns the planted defect of a Index whose
// Find returns a value that Put did not write over a fresh subject.
func indexFindReadsAnotherValue(h *IndexHandle) Index {
	return &indexFindReadsAnotherValueDefect{indexOverEmbedded: h.Fresh()}
}

// indexFindReadsAStaleValueAfterThreeDefect is the planted defect of a
// Index whose Find returns its first value again from its fourth call on.
type indexFindReadsAStaleValueAfterThreeDefect struct {
	indexOverEmbedded
	calls int
	first *Item
}

// Find is the override of the planted defect.
func (d *indexFindReadsAStaleValueAfterThreeDefect) Find(
	p0 context.Context,
	p1 string,
	p2 func(Item) bool,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.indexOverEmbedded.Find(p0, p1, p2)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// indexFindReadsAStaleValueAfterThree returns the planted defect of a
// Index whose Find returns its first value again from its fourth call on
// over a fresh subject.
func indexFindReadsAStaleValueAfterThree(h *IndexHandle) Index {
	return &indexFindReadsAStaleValueAfterThreeDefect{indexOverEmbedded: h.Fresh()}
}

// indexByRevPanicsDefect is the planted defect of a Index whose ByRev
// panics.
type indexByRevPanicsDefect struct {
	indexOverEmbedded
}

// ByRev is the override of the planted defect.
func (*indexByRevPanicsDefect) ByRev(p0 context.Context, p1 int) (r0 Item, r1 error) {
	panic("the planted defect of a Index whose ByRev panics")
}

// indexByRevPanics returns the planted defect of a Index whose ByRev
// panics over a fresh subject.
func indexByRevPanics(h *IndexHandle) Index {
	return &indexByRevPanicsDefect{indexOverEmbedded: h.Fresh()}
}

// indexByRevIgnoresCancelDefect is the planted defect of a Index whose
// ByRev ignores the context.
type indexByRevIgnoresCancelDefect struct {
	indexOverEmbedded
}

// ByRev is the override of the planted defect.
func (d *indexByRevIgnoresCancelDefect) ByRev(p0 context.Context, p1 int) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.indexOverEmbedded.ByRev(p0, p1)
}

// indexByRevIgnoresCancel returns the planted defect of a Index whose
// ByRev ignores the context over a fresh subject.
func indexByRevIgnoresCancel(h *IndexHandle) Index {
	return &indexByRevIgnoresCancelDefect{indexOverEmbedded: h.Fresh()}
}

// indexByRevValueBesideDefect is the planted defect of a Index whose ByRev
// returns a value beside an error.
type indexByRevValueBesideDefect struct {
	indexOverEmbedded
}

// ByRev is the override of the planted defect.
func (d *indexByRevValueBesideDefect) ByRev(p0 context.Context, p1 int) (r0 Item, r1 error) {
	r0, r1 = d.indexOverEmbedded.ByRev(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// indexByRevValueBeside returns the planted defect of a Index whose ByRev
// returns a value beside an error over a fresh subject.
func indexByRevValueBeside(h *IndexHandle) Index {
	return &indexByRevValueBesideDefect{indexOverEmbedded: h.Fresh()}
}

// indexByRevReadsAnotherValueDefect is the planted defect of a Index whose
// ByRev returns a value that Put did not write.
type indexByRevReadsAnotherValueDefect struct {
	indexOverEmbedded
}

// ByRev is the override of the planted defect.
func (d *indexByRevReadsAnotherValueDefect) ByRev(p0 context.Context, p1 int) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.indexOverEmbedded.ByRev(p0, p1)
	if r1 == nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// indexByRevReadsAnotherValue returns the planted defect of a Index whose
// ByRev returns a value that Put did not write over a fresh subject.
func indexByRevReadsAnotherValue(h *IndexHandle) Index {
	return &indexByRevReadsAnotherValueDefect{indexOverEmbedded: h.Fresh()}
}

// indexByRevReadsAStaleValueAfterThreeDefect is the planted defect of a
// Index whose ByRev returns its first value again from its fourth call on.
type indexByRevReadsAStaleValueAfterThreeDefect struct {
	indexOverEmbedded
	calls int
	first *Item
}

// ByRev is the override of the planted defect.
func (d *indexByRevReadsAStaleValueAfterThreeDefect) ByRev(
	p0 context.Context,
	p1 int,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.indexOverEmbedded.ByRev(p0, p1)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// indexByRevReadsAStaleValueAfterThree returns the planted defect of a
// Index whose ByRev returns its first value again from its fourth call on
// over a fresh subject.
func indexByRevReadsAStaleValueAfterThree(h *IndexHandle) Index {
	return &indexByRevReadsAStaleValueAfterThreeDefect{indexOverEmbedded: h.Fresh()}
}

// indexLookupPanicsDefect is the planted defect of a Index whose Lookup
// panics.
type indexLookupPanicsDefect struct {
	indexOverEmbedded
}

// Lookup is the override of the planted defect.
func (*indexLookupPanicsDefect) Lookup(p0 context.Context, p1 string, p2 int) (r0 Item, r1 error) {
	panic("the planted defect of a Index whose Lookup panics")
}

// indexLookupPanics returns the planted defect of a Index whose Lookup
// panics over a fresh subject.
func indexLookupPanics(h *IndexHandle) Index {
	return &indexLookupPanicsDefect{indexOverEmbedded: h.Fresh()}
}

// indexLookupIgnoresCancelDefect is the planted defect of a Index whose
// Lookup ignores the context.
type indexLookupIgnoresCancelDefect struct {
	indexOverEmbedded
}

// Lookup is the override of the planted defect.
func (d *indexLookupIgnoresCancelDefect) Lookup(
	p0 context.Context,
	p1 string,
	p2 int,
) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.indexOverEmbedded.Lookup(p0, p1, p2)
}

// indexLookupIgnoresCancel returns the planted defect of a Index whose
// Lookup ignores the context over a fresh subject.
func indexLookupIgnoresCancel(h *IndexHandle) Index {
	return &indexLookupIgnoresCancelDefect{indexOverEmbedded: h.Fresh()}
}

// indexLookupReadsAnotherValueDefect is the planted defect of a Index
// whose Lookup returns a value that Put did not write.
type indexLookupReadsAnotherValueDefect struct {
	indexOverEmbedded
}

// Lookup is the override of the planted defect.
func (d *indexLookupReadsAnotherValueDefect) Lookup(
	p0 context.Context,
	p1 string,
	p2 int,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.indexOverEmbedded.Lookup(p0, p1, p2)
	if r1 == nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// indexLookupReadsAnotherValue returns the planted defect of a Index whose
// Lookup returns a value that Put did not write over a fresh subject.
func indexLookupReadsAnotherValue(h *IndexHandle) Index {
	return &indexLookupReadsAnotherValueDefect{indexOverEmbedded: h.Fresh()}
}

// indexLookupReadsAStaleValueAfterThreeDefect is the planted defect of a
// Index whose Lookup returns its first value again from its fourth call
// on.
type indexLookupReadsAStaleValueAfterThreeDefect struct {
	indexOverEmbedded
	calls int
	first *Item
}

// Lookup is the override of the planted defect.
func (d *indexLookupReadsAStaleValueAfterThreeDefect) Lookup(
	p0 context.Context,
	p1 string,
	p2 int,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.indexOverEmbedded.Lookup(p0, p1, p2)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// indexLookupReadsAStaleValueAfterThree returns the planted defect of a
// Index whose Lookup returns its first value again from its fourth call on
// over a fresh subject.
func indexLookupReadsAStaleValueAfterThree(h *IndexHandle) Index {
	return &indexLookupReadsAStaleValueAfterThreeDefect{indexOverEmbedded: h.Fresh()}
}

// indexActions returns the actions of the generated sequences of Index
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func indexActions(c *prop.Case, subject func() Index) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	keysInt := c.Draw(
		prop.List(prop.Of[int](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.int",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Put",
			Input: func(c *prop.Case, _ struct{}) any {
				it := c.Draw(prop.Of[Item](), "put.it")
				it.Key = c.Draw(prop.SampledFrom(keysString...), "put.it.key")
				it.Rev = c.Draw(prop.SampledFrom(keysInt...), "put.it.rev")
				return []any{it}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				key0, _ := args[0].(Item)
				key1, _ := args[0].(Item)
				call := c.History().Invoke(client, "Put", args, key0.Key, key1.Rev)
				values, err := indexPutCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Get",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "get.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Get", args, args[0])
				values, err := indexGetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Find",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.Of[string](), "find.key")
				keepResults := c.Draw(prop.List(prop.Of[bool]()), "keep")
				var keepAt atomic.Int64
				keep := func(Item) (r0 bool) {
					i := int(keepAt.Add(1) - 1)
					if i >= len(keepResults) {
						return r0
					}
					return keepResults[i]
				}
				return []any{key, keep}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Find", args)
				values, err := indexFindCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "ByRev",
			Input: func(c *prop.Case, _ struct{}) any {
				rev := c.Draw(prop.SampledFrom(keysInt...), "by-rev.rev")
				return []any{rev}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "ByRev", args, args[0])
				values, err := indexByRevCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Lookup",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.Of[string](), "lookup.key")
				rev := c.Draw(prop.Of[int](), "lookup.rev")
				return []any{key, rev}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Lookup", args)
				values, err := indexLookupCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// indexPutCall calls Put on subject with the arguments args, and returns
// the values of the call beside its error.
func indexPutCall(ctx context.Context, subject Index, args []any) ([]any, error) {
	a0, _ := args[0].(Item)
	return nil, subject.Put(ctx, a0)
}

// indexGetCall calls Get on subject with the arguments args, and returns
// the values of the call beside its error.
func indexGetCall(ctx context.Context, subject Index, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Get(ctx, a0)
	return []any{r0}, err
}

// indexFindCall calls Find on subject with the arguments args, and returns
// the values of the call beside its error.
func indexFindCall(ctx context.Context, subject Index, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(func(Item) bool)
	r0, err := subject.Find(ctx, a0, a1)
	return []any{r0}, err
}

// indexByRevCall calls ByRev on subject with the arguments args, and
// returns the values of the call beside its error.
func indexByRevCall(ctx context.Context, subject Index, args []any) ([]any, error) {
	a0, _ := args[0].(int)
	r0, err := subject.ByRev(ctx, a0)
	return []any{r0}, err
}

// indexLookupCall calls Lookup on subject with the arguments args, and
// returns the values of the call beside its error.
func indexLookupCall(ctx context.Context, subject Index, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	a1, _ := args[1].(int)
	r0, err := subject.Lookup(ctx, a0, a1)
	return []any{r0}, err
}

// The derived checks of Stamped:
//
// Stamped.Store/smoke: Store returns without a panic for derived inputs.
//
// Stamped.Store/smoke#generated: Store returns without a panic for derived
// inputs.
//
// Stamped.Store/cancel: Store returns an error for a cancelled context.
//
// Stamped.Store/cancel#generated: Store returns an error for a cancelled
// context.
//
// Stamped.Get/smoke: Get returns without a panic for derived inputs.
//
// Stamped.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Stamped.Get/cancel: Get returns an error for a cancelled context.
//
// Stamped.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Stamped.Get/zero-on-error: Get returns the zero value beside an error.
//
// Stamped.Get/zero-on-error#generated: Get returns the zero value beside
// an error.
//
// Stamped.Get/monotonic-reads: Get never returns an older version than one
// that it returned before.
//
// Stamped.Get/monotonic-reads#generated: Get never returns an older
// version than one that it returned before.
//
// Stamped.Get/reads-never-go-back: Get never returns a version below one
// that it returned before, over any sequence of calls.

// The derived checks of Stamped that witness their laws at the shallow
// end:
//
// Stamped.Store/smoke
//
// Stamped.Get/smoke
//
// Stamped.Get/monotonic-reads
//
// Stamped.Get/monotonic-reads#generated

// The unclaimed classifications of Stamped:
//
// The shape:writer of Stamped.Store has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Stamped writes.
type (
	StampedHandle = conform.Handle[Stamped]
	StampedCheck  = conform.Check[Stamped]
	StampedOption = conform.Option[Stamped]
)

// StampedChecks contains the name of each callable of Stamped, with the
// identities of its checks.
var StampedChecks = struct {
	Store stampedStore
	Get   stampedGet
}{
	Store: stampedStore{"Stamped.Store"},
	Get:   stampedGet{"Stamped.Get"},
}

// stampedStore is the name of Stamped.Store, with the identities of its
// checks.
type stampedStore struct {
	conform.Callable[Stamped]
}

// Smoke returns the identity of the check that Store returns without a
// panic for derived inputs.
func (c stampedStore) Smoke() conform.ID[Stamped] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Store returns
// without a panic for derived inputs.
func (c stampedStore) SmokeGenerated() conform.ID[Stamped] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Store returns an error for
// a cancelled context.
func (c stampedStore) Cancel() conform.ID[Stamped] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Store returns an
// error for a cancelled context.
func (c stampedStore) CancelGenerated() conform.ID[Stamped] {
	return c.Check("cancel#generated")
}

// stampedGet is the name of Stamped.Get, with the identities of its
// checks.
type stampedGet struct {
	conform.Callable[Stamped]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c stampedGet) Smoke() conform.ID[Stamped] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c stampedGet) SmokeGenerated() conform.ID[Stamped] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c stampedGet) Cancel() conform.ID[Stamped] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c stampedGet) CancelGenerated() conform.ID[Stamped] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c stampedGet) ZeroOnError() conform.ID[Stamped] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c stampedGet) ZeroOnErrorGenerated() conform.ID[Stamped] {
	return c.Check("zero-on-error#generated")
}

// MonotonicReads returns the identity of the check that Get never returns
// an older version than one that it returned before.
func (c stampedGet) MonotonicReads() conform.ID[Stamped] {
	return c.Check("monotonic-reads")
}

// MonotonicReadsGenerated returns the identity of the check that Get never
// returns an older version than one that it returned before.
func (c stampedGet) MonotonicReadsGenerated() conform.ID[Stamped] {
	return c.Check("monotonic-reads#generated")
}

// ReadsNeverGoBack returns the identity of the check that Get never
// returns a version below one that it returned before, over any sequence
// of calls.
func (c stampedGet) ReadsNeverGoBack() conform.ID[Stamped] {
	return c.Check("reads-never-go-back")
}

// StampedSuite returns the checks of Stamped, and the calls and the
// budgets of its callables.
func StampedSuite() conform.Suite[Stamped] {
	return conform.Suite[Stamped]{
		Declaration: "Stamped",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Stamped]{
			{
				ID:          StampedChecks.Store.Smoke(),
				Statement:   "Store returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StampedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: stampedStorePanics,
			},
			{
				ID:          StampedChecks.Store.SmokeGenerated(),
				Statement:   "Store returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *StampedHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, drawnIt)
					})
				},
				Defect: stampedStorePanics,
			},
			{
				ID:          StampedChecks.Store.Cancel(),
				Statement:   "Store returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StampedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: stampedStoreIgnoresCancel,
			},
			{
				ID:          StampedChecks.Store.CancelGenerated(),
				Statement:   "Store returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *StampedHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, drawnIt)
					})
				},
				Defect: stampedStoreIgnoresCancel,
			},
			{
				ID:          StampedChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StampedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: stampedGetPanics,
			},
			{
				ID:          StampedChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *StampedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: stampedGetPanics,
			},
			{
				ID:          StampedChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StampedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: stampedGetIgnoresCancel,
			},
			{
				ID:          StampedChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *StampedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: stampedGetIgnoresCancel,
			},
			{
				ID:          StampedChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *StampedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Item{})
				},
				Defect: stampedGetValueBeside,
			},
			{
				ID:          StampedChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *StampedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, Item{})
				},
				Defect: stampedGetValueBeside,
			},
			{
				ID:          StampedChecks.Get.MonotonicReads(),
				Statement:   "Get never returns an older version than one that it returned before",
				Claim:       spec.MonotonicReads,
				Probe:       spec.ReadsNeverRegress,
				Reason:      spec.Regressed,
				Law:         spec.ReadsNeverGoBackwards,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *StampedHandle) {
					probe.ReadsNeverRegress(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, (Item{Key: "Key-a"}).Key)
					}, func(ctx context.Context) error {
						return h.Subject().Store(ctx, Item{Key: "Key-a"})
					}, func(v Item) int { return v.Rev })
				},
				Defect: stampedGetVersionDescends,
			},
			{
				ID:          StampedChecks.Get.MonotonicReadsGenerated(),
				Statement:   "Get never returns an older version than one that it returned before",
				Claim:       spec.MonotonicReads,
				Probe:       spec.ReadsNeverRegress,
				Reason:      spec.Regressed,
				Law:         spec.ReadsNeverGoBackwards,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *StampedHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.ReadsNeverRegress(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, drawnIt.Key)
					}, func(ctx context.Context) error {
						return h.Subject().Store(ctx, drawnIt)
					}, func(v Item) int { return v.Rev })
				},
				Defect: stampedGetVersionDescends,
			},
			{
				ID:          StampedChecks.Get.ReadsNeverGoBack(),
				Statement:   "Get never returns a version below one that it returned before, over any sequence of calls",
				Claim:       spec.ReadsNeverGoBack,
				Probe:       spec.ReadVersionsNeverFall,
				Reason:      spec.Regressed,
				Law:         spec.ReadsNeverGoBackwards,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Stamped]{StampedChecks.Get.MonotonicReads()},
				Property: func(h *StampedHandle, c *prop.Case) {
					probe.ReadVersionsNeverFall(
						h,
						c,
						stampedActions,
						"Get",
						func(v Item) int { return v.Rev },
					)
				},
				Defect: stampedGetReadsThePastAfterThree,
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Stamped.Store",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Stamped]{
			{
				Calls:    stampedStoreCalls,
				Over:     newStampedOver,
				Callable: StampedChecks.Store.Callable,
			},
			{
				Calls:    stampedGetCalls,
				Over:     newStampedOver,
				Callable: StampedChecks.Get.Callable,
			},
		},
	}
}

// RunStampedSuite applies the checks of Stamped, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunStampedSuite[T Stamped](t *testing.T, newSubject func() T, opts ...StampedOption) {
	t.Helper()
	RunStampedSuiteFromHandle(t, func(*StampedHandle) T { return newSubject() }, opts...)
}

// RunStampedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunStampedSuiteFromHandle[T Stamped](
	t *testing.T,
	newSubject func(h *StampedHandle) T,
	opts ...StampedOption,
) {
	t.Helper()
	conform.Run(t, StampedSuite(), stampedSubjects(newSubject), opts...)
}

// FuzzStampedSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzStampedSuite[T Stamped](
	f *testing.F,
	callable conform.Callable[Stamped],
	newSubject func(h *StampedHandle) T,
	opts ...StampedOption,
) {
	f.Helper()
	conform.Fuzz(f, StampedSuite(), callable, stampedSubjects(newSubject), opts...)
}

// RunStampedBudgets applies the allocation ceilings of Stamped and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunStampedBudgets[T Stamped](
	t *testing.T,
	newSubject func(h *StampedHandle) T,
	opts ...StampedOption,
) {
	t.Helper()
	conform.RunBudgets(t, StampedSuite(), stampedSubjects(newSubject), opts...)
}

// BenchStampedSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchStampedSuite[T Stamped](
	b conform.B,
	callable conform.Callable[Stamped],
	newSubject func(h *StampedHandle) T,
	opts ...StampedOption,
) {
	b.Helper()
	conform.Bench(b, StampedSuite(), callable, stampedSubjects(newSubject), opts...)
}

// stampedSubjects returns the factory of the subjects of Stamped that
// calls newSubject with the handle of each check.
func stampedSubjects[T Stamped](
	newSubject func(h *StampedHandle) T,
) func(h *StampedHandle) Stamped {
	return func(h *StampedHandle) Stamped { return newSubject(h) }
}

// stampedStoreCalls generates the calls of Stamped.Store with the inputs
// that each case draws.
var stampedStoreCalls = prop.Composite(func(c *prop.Case) func(context.Context, Stamped) error {
	it := c.Draw(prop.Of[Item](), "it")
	return func(ctx context.Context, subject Stamped) error {
		return subject.Store(ctx, it)
	}
})

// stampedGetCalls generates the calls of Stamped.Get with the inputs that
// each case draws.
var stampedGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Stamped) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Stamped) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// stampedOverEmbedded is Stamped under the name of the field that
// stampedOver embeds.
type stampedOverEmbedded = Stamped

// stampedOver is a subject of Stamped whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type stampedOver struct {
	stampedOverEmbedded
	spend func()
}

// Store calls spend and returns its zero results.
func (o stampedOver) Store(context.Context, Item) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o stampedOver) Get(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// newStampedOver returns the planted defect of a ceiling of Stamped, whose
// callables call spend.
func newStampedOver(_ *StampedHandle, spend func()) Stamped {
	return stampedOver{spend: spend}
}

// stampedStorePanicsDefect is the planted defect of a Stamped whose Store
// panics.
type stampedStorePanicsDefect struct {
	stampedOverEmbedded
}

// Store is the override of the planted defect.
func (*stampedStorePanicsDefect) Store(p0 context.Context, p1 Item) (r0 error) {
	panic("the planted defect of a Stamped whose Store panics")
}

// stampedStorePanics returns the planted defect of a Stamped whose Store
// panics over a fresh subject.
func stampedStorePanics(h *StampedHandle) Stamped {
	return &stampedStorePanicsDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedStoreIgnoresCancelDefect is the planted defect of a Stamped whose
// Store ignores the context.
type stampedStoreIgnoresCancelDefect struct {
	stampedOverEmbedded
}

// Store is the override of the planted defect.
func (d *stampedStoreIgnoresCancelDefect) Store(p0 context.Context, p1 Item) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.stampedOverEmbedded.Store(p0, p1)
}

// stampedStoreIgnoresCancel returns the planted defect of a Stamped whose
// Store ignores the context over a fresh subject.
func stampedStoreIgnoresCancel(h *StampedHandle) Stamped {
	return &stampedStoreIgnoresCancelDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedGetPanicsDefect is the planted defect of a Stamped whose Get
// panics.
type stampedGetPanicsDefect struct {
	stampedOverEmbedded
}

// Get is the override of the planted defect.
func (*stampedGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	panic("the planted defect of a Stamped whose Get panics")
}

// stampedGetPanics returns the planted defect of a Stamped whose Get
// panics over a fresh subject.
func stampedGetPanics(h *StampedHandle) Stamped {
	return &stampedGetPanicsDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedGetIgnoresCancelDefect is the planted defect of a Stamped whose
// Get ignores the context.
type stampedGetIgnoresCancelDefect struct {
	stampedOverEmbedded
}

// Get is the override of the planted defect.
func (d *stampedGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.stampedOverEmbedded.Get(p0, p1)
}

// stampedGetIgnoresCancel returns the planted defect of a Stamped whose
// Get ignores the context over a fresh subject.
func stampedGetIgnoresCancel(h *StampedHandle) Stamped {
	return &stampedGetIgnoresCancelDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedGetValueBesideDefect is the planted defect of a Stamped whose Get
// returns a value beside an error.
type stampedGetValueBesideDefect struct {
	stampedOverEmbedded
}

// Get is the override of the planted defect.
func (d *stampedGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	r0, r1 = d.stampedOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// stampedGetValueBeside returns the planted defect of a Stamped whose Get
// returns a value beside an error over a fresh subject.
func stampedGetValueBeside(h *StampedHandle) Stamped {
	return &stampedGetValueBesideDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedGetVersionDescendsDefect is the planted defect of a Stamped whose
// Get returns an older version on the second read.
type stampedGetVersionDescendsDefect struct {
	stampedOverEmbedded
	calls int
}

// Get is the override of the planted defect.
func (d *stampedGetVersionDescendsDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0 = Item{Key: "Key-b"}
	if d.calls == 1 {
		r0.Rev = 1
	} else {
		r0.Rev = 0
	}
	return r0, r1
}

// stampedGetVersionDescends returns the planted defect of a Stamped whose
// Get returns an older version on the second read over a fresh subject.
func stampedGetVersionDescends(h *StampedHandle) Stamped {
	return &stampedGetVersionDescendsDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedGetReadsThePastAfterThreeDefect is the planted defect of a
// Stamped whose Get returns its first value again from its fourth call on.
type stampedGetReadsThePastAfterThreeDefect struct {
	stampedOverEmbedded
	calls int
	first *Item
}

// Get is the override of the planted defect.
func (d *stampedGetReadsThePastAfterThreeDefect) Get(
	p0 context.Context,
	p1 string,
) (r0 Item, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.stampedOverEmbedded.Get(p0, p1)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// stampedGetReadsThePastAfterThree returns the planted defect of a Stamped
// whose Get returns its first value again from its fourth call on over a
// fresh subject.
func stampedGetReadsThePastAfterThree(h *StampedHandle) Stamped {
	return &stampedGetReadsThePastAfterThreeDefect{stampedOverEmbedded: h.Fresh()}
}

// stampedActions returns the actions of the generated sequences of Stamped
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func stampedActions(c *prop.Case, subject func() Stamped) []stateful.Action[struct{}] {
	keysString := c.Draw(
		prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string",
	)
	return []stateful.Action[struct{}]{
		{
			Name: "Store",
			Input: func(c *prop.Case, _ struct{}) any {
				it := c.Draw(prop.Of[Item](), "store.it")
				it.Key = c.Draw(prop.SampledFrom(keysString...), "store.it.key")
				return []any{it}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				key0, _ := args[0].(Item)
				call := c.History().Invoke(client, "Store", args, key0.Key)
				values, err := stampedStoreCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Get",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "get.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Get", args, args[0])
				values, err := stampedGetCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// stampedStoreCall calls Store on subject with the arguments args, and
// returns the values of the call beside its error.
func stampedStoreCall(ctx context.Context, subject Stamped, args []any) ([]any, error) {
	a0, _ := args[0].(Item)
	return nil, subject.Store(ctx, a0)
}

// stampedGetCall calls Get on subject with the arguments args, and returns
// the values of the call beside its error.
func stampedGetCall(ctx context.Context, subject Stamped, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Get(ctx, a0)
	return []any{r0}, err
}

// The derived checks of Versioned:
//
// Versioned.Put/smoke: Put returns without a panic for derived inputs.
//
// Versioned.Put/smoke#generated: Put returns without a panic for derived
// inputs.
//
// Versioned.Put/cancel: Put returns an error for a cancelled context.
//
// Versioned.Put/cancel#generated: Put returns an error for a cancelled
// context.
//
// Versioned.Get/smoke: Get returns without a panic for derived inputs.
//
// Versioned.Get/smoke#generated: Get returns without a panic for derived
// inputs.
//
// Versioned.Get/cancel: Get returns an error for a cancelled context.
//
// Versioned.Get/cancel#generated: Get returns an error for a cancelled
// context.
//
// Versioned.Get/zero-on-error: Get returns the zero value beside an error.
//
// Versioned.Get/zero-on-error#generated: Get returns the zero value beside
// an error.

// The derived checks of Versioned that witness their laws at the shallow
// end:
//
// Versioned.Put/smoke
//
// Versioned.Get/smoke

// The refused claims of Versioned:
//
// Versioned.Get/writes-follow-reads: the partner write of Get returns no
// value of the type that Get returns. The remedy is to return the written
// value from the partner write of Get.
//
// Versioned.Get/writes-follow-their-reads: the partner write of Get
// returns no value of the type that Get returns. The remedy is to return
// the written value from the partner write of Get.

// The unclaimed classifications of Versioned:
//
// The shape:writer of Versioned.Put has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Versioned writes.
type (
	VersionedHandle = conform.Handle[Versioned]
	VersionedCheck  = conform.Check[Versioned]
	VersionedOption = conform.Option[Versioned]
)

// VersionedChecks contains the name of each callable of Versioned, with
// the identities of its checks.
var VersionedChecks = struct {
	Put versionedPut
	Get versionedGet
}{
	Put: versionedPut{"Versioned.Put"},
	Get: versionedGet{"Versioned.Get"},
}

// versionedPut is the name of Versioned.Put, with the identities of its
// checks.
type versionedPut struct {
	conform.Callable[Versioned]
}

// Smoke returns the identity of the check that Put returns without a panic
// for derived inputs.
func (c versionedPut) Smoke() conform.ID[Versioned] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Put returns
// without a panic for derived inputs.
func (c versionedPut) SmokeGenerated() conform.ID[Versioned] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Put returns an error for a
// cancelled context.
func (c versionedPut) Cancel() conform.ID[Versioned] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Put returns an
// error for a cancelled context.
func (c versionedPut) CancelGenerated() conform.ID[Versioned] {
	return c.Check("cancel#generated")
}

// versionedGet is the name of Versioned.Get, with the identities of its
// checks.
type versionedGet struct {
	conform.Callable[Versioned]
}

// Smoke returns the identity of the check that Get returns without a panic
// for derived inputs.
func (c versionedGet) Smoke() conform.ID[Versioned] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Get returns
// without a panic for derived inputs.
func (c versionedGet) SmokeGenerated() conform.ID[Versioned] {
	return c.Check("smoke#generated")
}

// Cancel returns the identity of the check that Get returns an error for a
// cancelled context.
func (c versionedGet) Cancel() conform.ID[Versioned] {
	return c.Check("cancel")
}

// CancelGenerated returns the identity of the check that Get returns an
// error for a cancelled context.
func (c versionedGet) CancelGenerated() conform.ID[Versioned] {
	return c.Check("cancel#generated")
}

// ZeroOnError returns the identity of the check that Get returns the zero
// value beside an error.
func (c versionedGet) ZeroOnError() conform.ID[Versioned] {
	return c.Check("zero-on-error")
}

// ZeroOnErrorGenerated returns the identity of the check that Get returns
// the zero value beside an error.
func (c versionedGet) ZeroOnErrorGenerated() conform.ID[Versioned] {
	return c.Check("zero-on-error#generated")
}

// VersionedSuite returns the checks of Versioned, and the calls and the
// budgets of its callables.
func VersionedSuite() conform.Suite[Versioned] {
	return conform.Suite[Versioned]{
		Declaration: "Versioned",
		Package:     "example.com/shop/sequenced",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Versioned]{
			{
				ID:          VersionedChecks.Put.Smoke(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *VersionedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: versionedPutPanics,
			},
			{
				ID:          VersionedChecks.Put.SmokeGenerated(),
				Statement:   "Put returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VersionedHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnIt)
					})
				},
				Defect: versionedPutPanics,
			},
			{
				ID:          VersionedChecks.Put.Cancel(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VersionedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, Item{Key: "Key-a"})
					})
				},
				Defect: versionedPutIgnoresCancel,
			},
			{
				ID:          VersionedChecks.Put.CancelGenerated(),
				Statement:   "Put returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VersionedHandle, c *prop.Case) {
					drawnIt := c.Draw(prop.Of[Item](), "it")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Put(ctx, drawnIt)
					})
				},
				Defect: versionedPutIgnoresCancel,
			},
			{
				ID:          VersionedChecks.Get.Smoke(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *VersionedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: versionedGetPanics,
			},
			{
				ID:          VersionedChecks.Get.SmokeGenerated(),
				Statement:   "Get returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VersionedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: versionedGetPanics,
			},
			{
				ID:          VersionedChecks.Get.Cancel(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VersionedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, string("key-a"))
						return err
					})
				},
				Defect: versionedGetIgnoresCancel,
			},
			{
				ID:          VersionedChecks.Get.CancelGenerated(),
				Statement:   "Get returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VersionedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Get(ctx, drawnKey)
						return err
					})
				},
				Defect: versionedGetIgnoresCancel,
			},
			{
				ID:          VersionedChecks.Get.ZeroOnError(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *VersionedHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, string("key-a"))
					}, Item{})
				},
				Defect: versionedGetValueBeside,
			},
			{
				ID:          VersionedChecks.Get.ZeroOnErrorGenerated(),
				Statement:   "Get returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *VersionedHandle, c *prop.Case) {
					drawnKey := c.Draw(prop.Of[string](), "key")
					probe.ZeroBesideError(h, func(ctx context.Context) (Item, error) {
						return h.Subject().Get(ctx, drawnKey)
					}, Item{})
				},
				Defect: versionedGetValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.WritesFollowReads,
				Callable: "Versioned.Get",
				Gate:     spec.WriteReturnsTheValue,
				Why:      "the partner write of Get returns no value of the type that Get returns",
				Remedy:   "return the written value from the partner write of Get",
			},
			{
				Claim:    spec.WritesFollowTheirReads,
				Callable: "Versioned.Get",
				Gate:     spec.WriteReturnsTheValue,
				Why:      "the partner write of Get returns no value of the type that Get returns",
				Remedy:   "return the written value from the partner write of Get",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Versioned.Put",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Versioned]{
			{
				Calls:    versionedPutCalls,
				Over:     newVersionedOver,
				Callable: VersionedChecks.Put.Callable,
			},
			{
				Calls:    versionedGetCalls,
				Over:     newVersionedOver,
				Callable: VersionedChecks.Get.Callable,
			},
		},
	}
}

// RunVersionedSuite applies the checks of Versioned, and the checks that
// opts add, to a fresh subject of newSubject for each check.
func RunVersionedSuite[T Versioned](t *testing.T, newSubject func() T, opts ...VersionedOption) {
	t.Helper()
	RunVersionedSuiteFromHandle(t, func(*VersionedHandle) T { return newSubject() }, opts...)
}

// RunVersionedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunVersionedSuiteFromHandle[T Versioned](
	t *testing.T,
	newSubject func(h *VersionedHandle) T,
	opts ...VersionedOption,
) {
	t.Helper()
	conform.Run(t, VersionedSuite(), versionedSubjects(newSubject), opts...)
}

// FuzzVersionedSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzVersionedSuite[T Versioned](
	f *testing.F,
	callable conform.Callable[Versioned],
	newSubject func(h *VersionedHandle) T,
	opts ...VersionedOption,
) {
	f.Helper()
	conform.Fuzz(f, VersionedSuite(), callable, versionedSubjects(newSubject), opts...)
}

// RunVersionedBudgets applies the allocation ceilings of Versioned and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunVersionedBudgets[T Versioned](
	t *testing.T,
	newSubject func(h *VersionedHandle) T,
	opts ...VersionedOption,
) {
	t.Helper()
	conform.RunBudgets(t, VersionedSuite(), versionedSubjects(newSubject), opts...)
}

// BenchVersionedSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchVersionedSuite[T Versioned](
	b conform.B,
	callable conform.Callable[Versioned],
	newSubject func(h *VersionedHandle) T,
	opts ...VersionedOption,
) {
	b.Helper()
	conform.Bench(b, VersionedSuite(), callable, versionedSubjects(newSubject), opts...)
}

// versionedSubjects returns the factory of the subjects of Versioned that
// calls newSubject with the handle of each check.
func versionedSubjects[T Versioned](
	newSubject func(h *VersionedHandle) T,
) func(h *VersionedHandle) Versioned {
	return func(h *VersionedHandle) Versioned { return newSubject(h) }
}

// versionedPutCalls generates the calls of Versioned.Put with the inputs
// that each case draws.
var versionedPutCalls = prop.Composite(func(c *prop.Case) func(context.Context, Versioned) error {
	it := c.Draw(prop.Of[Item](), "it")
	return func(ctx context.Context, subject Versioned) error {
		return subject.Put(ctx, it)
	}
})

// versionedGetCalls generates the calls of Versioned.Get with the inputs
// that each case draws.
var versionedGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, Versioned) error {
	key := c.Draw(prop.Of[string](), "key")
	return func(ctx context.Context, subject Versioned) error {
		_, err := subject.Get(ctx, key)
		return err
	}
})

// versionedOverEmbedded is Versioned under the name of the field that
// versionedOver embeds.
type versionedOverEmbedded = Versioned

// versionedOver is a subject of Versioned whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type versionedOver struct {
	versionedOverEmbedded
	spend func()
}

// Put calls spend and returns its zero results.
func (o versionedOver) Put(context.Context, Item) (r0 error) {
	o.spend()
	return r0
}

// Get calls spend and returns its zero results.
func (o versionedOver) Get(context.Context, string) (r0 Item, r1 error) {
	o.spend()
	return r0, r1
}

// newVersionedOver returns the planted defect of a ceiling of Versioned,
// whose callables call spend.
func newVersionedOver(_ *VersionedHandle, spend func()) Versioned {
	return versionedOver{spend: spend}
}

// versionedPutPanicsDefect is the planted defect of a Versioned whose Put
// panics.
type versionedPutPanicsDefect struct {
	versionedOverEmbedded
}

// Put is the override of the planted defect.
func (*versionedPutPanicsDefect) Put(p0 context.Context, p1 Item) (r0 error) {
	panic("the planted defect of a Versioned whose Put panics")
}

// versionedPutPanics returns the planted defect of a Versioned whose Put
// panics over a fresh subject.
func versionedPutPanics(h *VersionedHandle) Versioned {
	return &versionedPutPanicsDefect{versionedOverEmbedded: h.Fresh()}
}

// versionedPutIgnoresCancelDefect is the planted defect of a Versioned
// whose Put ignores the context.
type versionedPutIgnoresCancelDefect struct {
	versionedOverEmbedded
}

// Put is the override of the planted defect.
func (d *versionedPutIgnoresCancelDefect) Put(p0 context.Context, p1 Item) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.versionedOverEmbedded.Put(p0, p1)
}

// versionedPutIgnoresCancel returns the planted defect of a Versioned
// whose Put ignores the context over a fresh subject.
func versionedPutIgnoresCancel(h *VersionedHandle) Versioned {
	return &versionedPutIgnoresCancelDefect{versionedOverEmbedded: h.Fresh()}
}

// versionedGetPanicsDefect is the planted defect of a Versioned whose Get
// panics.
type versionedGetPanicsDefect struct {
	versionedOverEmbedded
}

// Get is the override of the planted defect.
func (*versionedGetPanicsDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	panic("the planted defect of a Versioned whose Get panics")
}

// versionedGetPanics returns the planted defect of a Versioned whose Get
// panics over a fresh subject.
func versionedGetPanics(h *VersionedHandle) Versioned {
	return &versionedGetPanicsDefect{versionedOverEmbedded: h.Fresh()}
}

// versionedGetIgnoresCancelDefect is the planted defect of a Versioned
// whose Get ignores the context.
type versionedGetIgnoresCancelDefect struct {
	versionedOverEmbedded
}

// Get is the override of the planted defect.
func (d *versionedGetIgnoresCancelDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.versionedOverEmbedded.Get(p0, p1)
}

// versionedGetIgnoresCancel returns the planted defect of a Versioned
// whose Get ignores the context over a fresh subject.
func versionedGetIgnoresCancel(h *VersionedHandle) Versioned {
	return &versionedGetIgnoresCancelDefect{versionedOverEmbedded: h.Fresh()}
}

// versionedGetValueBesideDefect is the planted defect of a Versioned whose
// Get returns a value beside an error.
type versionedGetValueBesideDefect struct {
	versionedOverEmbedded
}

// Get is the override of the planted defect.
func (d *versionedGetValueBesideDefect) Get(p0 context.Context, p1 string) (r0 Item, r1 error) {
	r0, r1 = d.versionedOverEmbedded.Get(p0, p1)
	if r1 != nil {
		r0 = Item{Key: "Key-b"}
	}
	return r0, r1
}

// versionedGetValueBeside returns the planted defect of a Versioned whose
// Get returns a value beside an error over a fresh subject.
func versionedGetValueBeside(h *VersionedHandle) Versioned {
	return &versionedGetValueBesideDefect{versionedOverEmbedded: h.Fresh()}
}
