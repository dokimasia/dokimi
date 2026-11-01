package wrap

import (
	"context"
	"sync/atomic"
	"testing"

	"example.com/shop/lib/tool"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/assert/stateful"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Reach:
//
// Reach.Store/smoke#generated: Store returns without a panic for derived
// inputs.
//
// Reach.Store/cancel#generated: Store returns an error for a cancelled
// context.
//
// Reach.Current/smoke: Current returns without a panic for derived inputs.
//
// Reach.Current/cancel: Current returns an error for a cancelled context.
//
// Reach.Current/zero-on-error: Current returns the zero value beside an
// error.
//
// Reach.Current/read-after-write#generated: Current returns what Store
// stored before it.
//
// Reach.Current/reads-each-last-write: Current returns what Store wrote
// right before it, over any sequence of calls.

// The derived checks of Reach that witness their laws at the shallow end:
//
// Reach.Current/smoke
//
// Reach.Current/read-after-write#generated

// The refused claims of Reach:
//
// Reach.Store/smoke: a sample of the check refers to a declaration that
// the file of the suite cannot import. The remedy is to author a sample of
// each input of Store with the directive sample.
//
// Reach.Store/cancel: a sample of the check refers to a declaration that
// the file of the suite cannot import. The remedy is to author a sample of
// each input of Store with the directive sample.
//
// Reach.Current/read-after-write: a sample of the check refers to a
// declaration that the file of the suite cannot import. The remedy is to
// author a sample of each input of Current with the directive sample.
//
// Reach.Gear/smoke: a sample of the check refers to a declaration that the
// file of the suite cannot import. The remedy is to author a sample of
// each input of Gear with the directive sample.
//
// Reach.Gear/cancel: a sample of the check refers to a declaration that
// the file of the suite cannot import. The remedy is to author a sample of
// each input of Gear with the directive sample.
//
// Reach.Gear/zero-on-error: a sample of the check refers to a declaration
// that the file of the suite cannot import. The remedy is to author a
// sample of each input of Gear with the directive sample.

// The unclaimed classifications of Reach:
//
// The shape:writer of Reach.Store has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Reach writes.
type (
	ReachHandle = conform.Handle[Reach]
	ReachCheck  = conform.Check[Reach]
	ReachOption = conform.Option[Reach]
)

// ReachChecks contains the name of each callable of Reach, with the
// identities of its checks.
var ReachChecks = struct {
	Store   reachStore
	Current reachCurrent
	Gear    reachGear
}{
	Store:   reachStore{"Reach.Store"},
	Current: reachCurrent{"Reach.Current"},
	Gear:    reachGear{"Reach.Gear"},
}

// reachStore is the name of Reach.Store, with the identities of its
// checks.
type reachStore struct {
	conform.Callable[Reach]
}

// SmokeGenerated returns the identity of the check that Store returns
// without a panic for derived inputs.
func (c reachStore) SmokeGenerated() conform.ID[Reach] {
	return c.Check("smoke#generated")
}

// CancelGenerated returns the identity of the check that Store returns an
// error for a cancelled context.
func (c reachStore) CancelGenerated() conform.ID[Reach] {
	return c.Check("cancel#generated")
}

// reachCurrent is the name of Reach.Current, with the identities of its
// checks.
type reachCurrent struct {
	conform.Callable[Reach]
}

// Smoke returns the identity of the check that Current returns without a
// panic for derived inputs.
func (c reachCurrent) Smoke() conform.ID[Reach] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Current returns an error
// for a cancelled context.
func (c reachCurrent) Cancel() conform.ID[Reach] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Current returns the
// zero value beside an error.
func (c reachCurrent) ZeroOnError() conform.ID[Reach] {
	return c.Check("zero-on-error")
}

// ReadAfterWriteGenerated returns the identity of the check that Current
// returns what Store stored before it.
func (c reachCurrent) ReadAfterWriteGenerated() conform.ID[Reach] {
	return c.Check("read-after-write#generated")
}

// ReadsEachLastWrite returns the identity of the check that Current
// returns what Store wrote right before it, over any sequence of calls.
func (c reachCurrent) ReadsEachLastWrite() conform.ID[Reach] {
	return c.Check("reads-each-last-write")
}

// reachGear is the name of Reach.Gear, with the identities of its checks.
type reachGear struct {
	conform.Callable[Reach]
}

// ReachSuite returns the checks of Reach, and the calls and the budgets of
// its callables.
func ReachSuite() conform.Suite[Reach] {
	return conform.Suite[Reach]{
		Declaration: "Reach",
		Package:     "example.com/shop/wrap",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Reach]{
			{
				ID:          ReachChecks.Store.SmokeGenerated(),
				Statement:   "Store returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ReachHandle, c *prop.Case) {
					drawnP := c.Draw(prop.Of[tool.Pitch](), "p")
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, drawnP)
					})
				},
				Defect: reachStorePanics,
			},
			{
				ID:          ReachChecks.Store.CancelGenerated(),
				Statement:   "Store returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *ReachHandle, c *prop.Case) {
					drawnP := c.Draw(prop.Of[tool.Pitch](), "p")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Store(ctx, drawnP)
					})
				},
				Defect: reachStoreIgnoresCancel,
			},
			{
				ID:          ReachChecks.Current.Smoke(),
				Statement:   "Current returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *ReachHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Current(ctx)
						return err
					})
				},
				Defect: reachCurrentPanics,
			},
			{
				ID:          ReachChecks.Current.Cancel(),
				Statement:   "Current returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ReachHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Current(ctx)
						return err
					})
				},
				Defect: reachCurrentIgnoresCancel,
			},
			{
				ID:          ReachChecks.Current.ZeroOnError(),
				Statement:   "Current returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *ReachHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (tool.Pitch, error) {
						return h.Subject().Current(ctx)
					}, tool.Pitch(0))
				},
				Defect: reachCurrentValueBeside,
			},
			{
				ID:          ReachChecks.Current.ReadAfterWriteGenerated(),
				Statement:   "Current returns what Store stored before it",
				Claim:       spec.ReadAfterWrite,
				Probe:       spec.WriteThenRead,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Shallow:     true,
				Property: func(h *ReachHandle, c *prop.Case) {
					drawnP := c.Draw(prop.Of[tool.Pitch](), "p")
					probe.WriteThenRead(h, func(ctx context.Context) (tool.Pitch, error) {
						return h.Subject().Current(ctx)
					}, func(ctx context.Context) (tool.Pitch, error) {
						w := drawnP
						err := h.Subject().Store(ctx, w)
						return w, err
					})
				},
				Defect: reachCurrentReadsAnotherValue,
			},
			{
				ID:          ReachChecks.Current.ReadsEachLastWrite(),
				Statement:   "Current returns what Store wrote right before it, over any sequence of calls",
				Claim:       spec.ReadsEachLastWrite,
				Probe:       spec.ReadsReturnTheLastWrite,
				Reason:      spec.UnobservedWrite,
				Law:         spec.ReadReturnsTheWriteBeforeIt,
				Machineries: []spec.MachineryID{spec.GeneratedSequences},
				Entails:     []conform.ID[Reach]{ReachChecks.Current.ReadAfterWrite()},
				Property: func(h *ReachHandle, c *prop.Case) {
					probe.ReadsReturnTheLastWrite(
						h,
						c,
						reachActions,
						"Store",
						func(ctx context.Context, args []any) conform.Output {
							values, err := reachCurrentCall(ctx, h.Subject(), []any{})
							return conform.OutputOf(err, values...)
						},
						func(args []any) tool.Pitch {
							v, _ := args[0].(tool.Pitch)
							return v
						},
					)
				},
				Defect: reachCurrentReadsAStaleValueAfterThree,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Smoke,
				Callable: "Reach.Store",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Store with the directive sample",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Reach.Store",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Store with the directive sample",
			},
			{
				Claim:    spec.ReadAfterWrite,
				Callable: "Reach.Current",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Current with the directive sample",
			},
			{
				Claim:    spec.Smoke,
				Callable: "Reach.Gear",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Gear with the directive sample",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Reach.Gear",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Gear with the directive sample",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Reach.Gear",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Gear with the directive sample",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Reach.Store",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Reach]{
			{
				Calls:    reachStoreCalls,
				Over:     newReachOver,
				Callable: ReachChecks.Store.Callable,
			},
			{
				Calls:    reachCurrentCalls,
				Over:     newReachOver,
				Callable: ReachChecks.Current.Callable,
			},
			{
				Calls:    reachGearCalls,
				Over:     newReachOver,
				Callable: ReachChecks.Gear.Callable,
			},
		},
	}
}

// RunReachSuite applies the checks of Reach, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunReachSuite[T Reach](t *testing.T, newSubject func() T, opts ...ReachOption) {
	t.Helper()
	RunReachSuiteFromHandle(t, func(*ReachHandle) T { return newSubject() }, opts...)
}

// RunReachSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunReachSuiteFromHandle[T Reach](
	t *testing.T,
	newSubject func(h *ReachHandle) T,
	opts ...ReachOption,
) {
	t.Helper()
	conform.Run(t, ReachSuite(), reachSubjects(newSubject), opts...)
}

// FuzzReachSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzReachSuite[T Reach](
	f *testing.F,
	callable conform.Callable[Reach],
	newSubject func(h *ReachHandle) T,
	opts ...ReachOption,
) {
	f.Helper()
	conform.Fuzz(f, ReachSuite(), callable, reachSubjects(newSubject), opts...)
}

// RunReachBudgets applies the allocation ceilings of Reach and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunReachBudgets[T Reach](
	t *testing.T,
	newSubject func(h *ReachHandle) T,
	opts ...ReachOption,
) {
	t.Helper()
	conform.RunBudgets(t, ReachSuite(), reachSubjects(newSubject), opts...)
}

// BenchReachSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchReachSuite[T Reach](
	b conform.B,
	callable conform.Callable[Reach],
	newSubject func(h *ReachHandle) T,
	opts ...ReachOption,
) {
	b.Helper()
	conform.Bench(b, ReachSuite(), callable, reachSubjects(newSubject), opts...)
}

// reachSubjects returns the factory of the subjects of Reach that calls
// newSubject with the handle of each check.
func reachSubjects[T Reach](newSubject func(h *ReachHandle) T) func(h *ReachHandle) Reach {
	return func(h *ReachHandle) Reach { return newSubject(h) }
}

// reachStoreCalls generates the calls of Reach.Store with the inputs that
// each case draws.
var reachStoreCalls = prop.Composite(func(c *prop.Case) func(context.Context, Reach) error {
	p := c.Draw(prop.Of[tool.Pitch](), "p")
	return func(ctx context.Context, subject Reach) error {
		return subject.Store(ctx, p)
	}
})

// reachCurrentCalls generates the calls of Reach.Current with the inputs
// that each case draws.
var reachCurrentCalls = prop.Composite(func(*prop.Case) func(context.Context, Reach) error {
	return func(ctx context.Context, subject Reach) error {
		_, err := subject.Current(ctx)
		return err
	}
})

// reachGearCalls generates the calls of Reach.Gear with the inputs that
// each case draws.
var reachGearCalls = prop.Composite(func(*prop.Case) func(context.Context, Reach) error {
	return func(ctx context.Context, subject Reach) error {
		_, err := subject.Gear(ctx)
		return err
	}
})

// reachOverEmbedded is Reach under the name of the field that reachOver
// embeds.
type reachOverEmbedded = Reach

// reachOver is a subject of Reach whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type reachOver struct {
	reachOverEmbedded
	spend func()
}

// Store calls spend and returns its zero results.
func (o reachOver) Store(context.Context, tool.Pitch) (r0 error) {
	o.spend()
	return r0
}

// Current calls spend and returns its zero results.
func (o reachOver) Current(context.Context) (r0 tool.Pitch, r1 error) {
	o.spend()
	return r0, r1
}

// Gear calls spend and returns its zero results.
func (o reachOver) Gear(context.Context) (r0 tool.Gear, r1 error) {
	o.spend()
	return r0, r1
}

// newReachOver returns the planted defect of a ceiling of Reach, whose
// callables call spend.
func newReachOver(_ *ReachHandle, spend func()) Reach {
	return reachOver{spend: spend}
}

// reachStorePanicsDefect is the planted defect of a Reach whose Store
// panics.
type reachStorePanicsDefect struct {
	reachOverEmbedded
}

// Store is the override of the planted defect.
func (*reachStorePanicsDefect) Store(p0 context.Context, p1 tool.Pitch) (r0 error) {
	panic("the planted defect of a Reach whose Store panics")
}

// reachStorePanics returns the planted defect of a Reach whose Store
// panics over a fresh subject.
func reachStorePanics(h *ReachHandle) Reach {
	return &reachStorePanicsDefect{reachOverEmbedded: h.Fresh()}
}

// reachStoreIgnoresCancelDefect is the planted defect of a Reach whose
// Store ignores the context.
type reachStoreIgnoresCancelDefect struct {
	reachOverEmbedded
}

// Store is the override of the planted defect.
func (d *reachStoreIgnoresCancelDefect) Store(p0 context.Context, p1 tool.Pitch) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.reachOverEmbedded.Store(p0, p1)
}

// reachStoreIgnoresCancel returns the planted defect of a Reach whose
// Store ignores the context over a fresh subject.
func reachStoreIgnoresCancel(h *ReachHandle) Reach {
	return &reachStoreIgnoresCancelDefect{reachOverEmbedded: h.Fresh()}
}

// reachCurrentPanicsDefect is the planted defect of a Reach whose Current
// panics.
type reachCurrentPanicsDefect struct {
	reachOverEmbedded
}

// Current is the override of the planted defect.
func (*reachCurrentPanicsDefect) Current(p0 context.Context) (r0 tool.Pitch, r1 error) {
	panic("the planted defect of a Reach whose Current panics")
}

// reachCurrentPanics returns the planted defect of a Reach whose Current
// panics over a fresh subject.
func reachCurrentPanics(h *ReachHandle) Reach {
	return &reachCurrentPanicsDefect{reachOverEmbedded: h.Fresh()}
}

// reachCurrentIgnoresCancelDefect is the planted defect of a Reach whose
// Current ignores the context.
type reachCurrentIgnoresCancelDefect struct {
	reachOverEmbedded
}

// Current is the override of the planted defect.
func (d *reachCurrentIgnoresCancelDefect) Current(p0 context.Context) (r0 tool.Pitch, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.reachOverEmbedded.Current(p0)
}

// reachCurrentIgnoresCancel returns the planted defect of a Reach whose
// Current ignores the context over a fresh subject.
func reachCurrentIgnoresCancel(h *ReachHandle) Reach {
	return &reachCurrentIgnoresCancelDefect{reachOverEmbedded: h.Fresh()}
}

// reachCurrentValueBesideDefect is the planted defect of a Reach whose
// Current returns a value beside an error.
type reachCurrentValueBesideDefect struct {
	reachOverEmbedded
}

// Current is the override of the planted defect.
func (d *reachCurrentValueBesideDefect) Current(p0 context.Context) (r0 tool.Pitch, r1 error) {
	r0, r1 = d.reachOverEmbedded.Current(p0)
	if r1 != nil {
		r0 = tool.Pitch(7)
	}
	return r0, r1
}

// reachCurrentValueBeside returns the planted defect of a Reach whose
// Current returns a value beside an error over a fresh subject.
func reachCurrentValueBeside(h *ReachHandle) Reach {
	return &reachCurrentValueBesideDefect{reachOverEmbedded: h.Fresh()}
}

// reachCurrentReadsAnotherValueDefect is the planted defect of a Reach
// whose Current returns a value that Store did not write.
type reachCurrentReadsAnotherValueDefect struct {
	reachOverEmbedded
}

// Current is the override of the planted defect.
func (d *reachCurrentReadsAnotherValueDefect) Current(
	p0 context.Context,
) (r0 tool.Pitch, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	r0, r1 = d.reachOverEmbedded.Current(p0)
	if r1 == nil {
		r0 = tool.Pitch(7)
	}
	return r0, r1
}

// reachCurrentReadsAnotherValue returns the planted defect of a Reach
// whose Current returns a value that Store did not write over a fresh
// subject.
func reachCurrentReadsAnotherValue(h *ReachHandle) Reach {
	return &reachCurrentReadsAnotherValueDefect{reachOverEmbedded: h.Fresh()}
}

// reachCurrentReadsAStaleValueAfterThreeDefect is the planted defect of a
// Reach whose Current returns its first value again from its fourth call
// on.
type reachCurrentReadsAStaleValueAfterThreeDefect struct {
	reachOverEmbedded
	calls int
	first *tool.Pitch
}

// Current is the override of the planted defect.
func (d *reachCurrentReadsAStaleValueAfterThreeDefect) Current(
	p0 context.Context,
) (r0 tool.Pitch, r1 error) {
	if err := p0.Err(); err != nil {
		r1 = err
		return r0, r1
	}
	d.calls++
	r0, r1 = d.reachOverEmbedded.Current(p0)
	if d.first == nil {
		d.first = &r0
	}
	if d.calls > 3 {
		r0 = *d.first
		r1 = nil
	}
	return r0, r1
}

// reachCurrentReadsAStaleValueAfterThree returns the planted defect of a
// Reach whose Current returns its first value again from its fourth call
// on over a fresh subject.
func reachCurrentReadsAStaleValueAfterThree(h *ReachHandle) Reach {
	return &reachCurrentReadsAStaleValueAfterThreeDefect{reachOverEmbedded: h.Fresh()}
}

// reachActions returns the actions of the generated sequences of Reach
// over the subject that subject returns. The actions of a case take their
// keys from the pools that it draws.
func reachActions(c *prop.Case, subject func() Reach) []stateful.Action[struct{}] {
	return []stateful.Action[struct{}]{
		{
			Name: "Store",
			Input: func(c *prop.Case, _ struct{}) any {
				p := c.Draw(prop.Of[tool.Pitch](), "store.p")
				return []any{p}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Store", args)
				values, err := reachStoreCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Current",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Current", args)
				values, err := reachCurrentCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Gear",
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Gear", args)
				values, err := reachGearCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// reachStoreCall calls Store on subject with the arguments args, and
// returns the values of the call beside its error.
func reachStoreCall(ctx context.Context, subject Reach, args []any) ([]any, error) {
	a0, _ := args[0].(tool.Pitch)
	return nil, subject.Store(ctx, a0)
}

// reachCurrentCall calls Current on subject with the arguments args, and
// returns the values of the call beside its error.
func reachCurrentCall(ctx context.Context, subject Reach, args []any) ([]any, error) {
	r0, err := subject.Current(ctx)
	return []any{r0}, err
}

// reachGearCall calls Gear on subject with the arguments args, and returns
// the values of the call beside its error.
func reachGearCall(ctx context.Context, subject Reach, args []any) ([]any, error) {
	r0, err := subject.Gear(ctx)
	return []any{r0}, err
}

// The refused claims of Tuned:
//
// Tuned.Tune/smoke: a sample of the check refers to a declaration that the
// file of the suite cannot import. The remedy is to author a sample of
// each input of Tune with the directive sample.

// The unclaimed classifications of Tuned:
//
// The shape:writer of Tuned.Tune has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Tuned.Watch has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Tuned.Turn has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Tuned writes.
type (
	TunedHandle = conform.Handle[Tuned]
	TunedCheck  = conform.Check[Tuned]
	TunedOption = conform.Option[Tuned]
)

// TunedChecks contains the name of each callable of Tuned, with the
// identities of its checks.
var TunedChecks = struct {
	Tune  tunedTune
	Watch tunedWatch
	Turn  tunedTurn
}{
	Tune:  tunedTune{"Tuned.Tune"},
	Watch: tunedWatch{"Tuned.Watch"},
	Turn:  tunedTurn{"Tuned.Turn"},
}

// tunedTune is the name of Tuned.Tune, with the identities of its checks.
type tunedTune struct {
	conform.Callable[Tuned]
}

// tunedWatch is the name of Tuned.Watch, with the identities of its
// checks.
type tunedWatch struct {
	conform.Callable[Tuned]
}

// tunedTurn is the name of Tuned.Turn, with the identities of its checks.
type tunedTurn struct {
	conform.Callable[Tuned]
}

// TunedSuite returns the checks of Tuned, and the calls and the budgets of
// its callables.
func TunedSuite() conform.Suite[Tuned] {
	return conform.Suite[Tuned]{
		Declaration: "Tuned",
		Package:     "example.com/shop/wrap",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Smoke,
				Callable: "Tuned.Tune",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Tune with the directive sample",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:writer",
				Callable:       "Tuned.Tune",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Tuned.Watch",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Tuned.Turn",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Tuned]{
			{
				Calls:    tunedTuneCalls,
				Over:     newTunedOver,
				Callable: TunedChecks.Tune.Callable,
			},
			{
				Calls:    tunedWatchCalls,
				Over:     newTunedOver,
				Callable: TunedChecks.Watch.Callable,
			},
			{
				Calls:    tunedTurnCalls,
				Over:     newTunedOver,
				Callable: TunedChecks.Turn.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.SampledFrom(tool.Percent, tool.Steps)),
			prop.Using(tunedSettingValues),
			prop.Using(prop.Composite(func(c *prop.Case) tool.Note {
				noteWriteResults := c.Draw(
					prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
					"note-write",
				)
				var noteWriteAt atomic.Int64
				return &tunedNoteInput{
					onWrite: func(string) (r0 error) {
						i := int(noteWriteAt.Add(1) - 1)
						if i >= len(noteWriteResults) {
							r0 = conform.ErrGenerated
							return r0
						}
						return noteWriteResults[i]
					},
				}
			})),
		},
	}
}

// RunTunedSuite applies the checks of Tuned, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunTunedSuite[T Tuned](t *testing.T, newSubject func() T, opts ...TunedOption) {
	t.Helper()
	RunTunedSuiteFromHandle(t, func(*TunedHandle) T { return newSubject() }, opts...)
}

// RunTunedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunTunedSuiteFromHandle[T Tuned](
	t *testing.T,
	newSubject func(h *TunedHandle) T,
	opts ...TunedOption,
) {
	t.Helper()
	conform.Run(t, TunedSuite(), tunedSubjects(newSubject), opts...)
}

// FuzzTunedSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzTunedSuite[T Tuned](
	f *testing.F,
	callable conform.Callable[Tuned],
	newSubject func(h *TunedHandle) T,
	opts ...TunedOption,
) {
	f.Helper()
	conform.Fuzz(f, TunedSuite(), callable, tunedSubjects(newSubject), opts...)
}

// RunTunedBudgets applies the allocation ceilings of Tuned and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunTunedBudgets[T Tuned](
	t *testing.T,
	newSubject func(h *TunedHandle) T,
	opts ...TunedOption,
) {
	t.Helper()
	conform.RunBudgets(t, TunedSuite(), tunedSubjects(newSubject), opts...)
}

// BenchTunedSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchTunedSuite[T Tuned](
	b conform.B,
	callable conform.Callable[Tuned],
	newSubject func(h *TunedHandle) T,
	opts ...TunedOption,
) {
	b.Helper()
	conform.Bench(b, TunedSuite(), callable, tunedSubjects(newSubject), opts...)
}

// tunedSubjects returns the factory of the subjects of Tuned that calls
// newSubject with the handle of each check.
func tunedSubjects[T Tuned](newSubject func(h *TunedHandle) T) func(h *TunedHandle) Tuned {
	return func(h *TunedHandle) Tuned { return newSubject(h) }
}

// tunedTuneCalls generates the calls of Tuned.Tune with the inputs that
// each case draws.
var tunedTuneCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tuned) error {
	r := c.Draw(prop.Of[tool.Request](), "r")
	return func(_ context.Context, subject Tuned) error {
		return subject.Tune(r)
	}
})

// tunedWatchCalls generates the calls of Tuned.Watch with the inputs that
// each case draws.
var tunedWatchCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tuned) error {
	g := c.Draw(prop.Of[tool.Gauge](), "g")
	return func(_ context.Context, subject Tuned) error {
		subject.Watch(g)
		return nil
	}
})

// tunedTurnCalls generates the calls of Tuned.Turn with the inputs that
// each case draws.
var tunedTurnCalls = prop.Composite(func(c *prop.Case) func(context.Context, Tuned) error {
	d := c.Draw(prop.Of[tool.Dial](), "d")
	return func(_ context.Context, subject Tuned) error {
		subject.Turn(d)
		return nil
	}
})

// tunedSettingValues generates the values of Setting that keep or change
// each field of its starting value.
var tunedSettingValues = prop.Composite(func(c *prop.Case) tool.Setting {
	changed := c.Draw(prop.Of[tool.Setting](), "setting")
	value := tool.Setting{Value: 5}
	if c.Draw(prop.Boolean(), "setting.value") {
		value.Value = changed.Value
	}
	return value
})

// tunedNoteInputEmbedded is Note under the name of the field that
// tunedNoteInput embeds.
type tunedNoteInputEmbedded = tool.Note

// tunedNoteInput implements Note. Its methods return the results that the
// generators of the suite of Tuned draw for them, in turn.
type tunedNoteInput struct {
	tunedNoteInputEmbedded
	onWrite func(a0 string) error
}

// Write calls the function of Write.
func (in tunedNoteInput) Write(a0 string) error {
	return in.onWrite(a0)
}

// tunedOverEmbedded is Tuned under the name of the field that tunedOver
// embeds.
type tunedOverEmbedded = Tuned

// tunedOver is a subject of Tuned whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type tunedOver struct {
	tunedOverEmbedded
	spend func()
}

// Tune calls spend and returns its zero results.
func (o tunedOver) Tune(tool.Request) (r0 error) {
	o.spend()
	return r0
}

// Watch calls spend and returns its zero results.
func (o tunedOver) Watch(tool.Gauge) {
	o.spend()
}

// Turn calls spend and returns its zero results.
func (o tunedOver) Turn(tool.Dial) {
	o.spend()
}

// newTunedOver returns the planted defect of a ceiling of Tuned, whose
// callables call spend.
func newTunedOver(_ *TunedHandle, spend func()) Tuned {
	return tunedOver{spend: spend}
}

// The derived checks of Limited:
//
// Limited.Limit/smoke: Limit returns without a panic for derived inputs.
//
// Limited.Limit/smoke#generated: Limit returns without a panic for derived
// inputs.
//
// Limited.Grip/smoke#generated: Grip returns without a panic for derived
// inputs.
//
// Limited.Span/smoke#generated: Span returns without a panic for derived
// inputs.
//
// Limited.Pitch/smoke#generated: Pitch returns without a panic for derived
// inputs.
//
// Limited.Measure/smoke: Measure returns without a panic for derived
// inputs.
//
// Limited.Measure/smoke#generated: Measure returns without a panic for
// derived inputs.
//
// Limited.Phase/smoke: Phase returns without a panic for derived inputs.
//
// Limited.Phase/smoke#generated: Phase returns without a panic for derived
// inputs.

// The derived checks of Limited that witness their laws at the shallow
// end:
//
// Limited.Limit/smoke
//
// Limited.Measure/smoke
//
// Limited.Phase/smoke

// The refused claims of Limited:
//
// Limited.Grip/smoke: a sample of the check refers to a declaration that
// the file of the suite cannot import. The remedy is to author a sample of
// each input of Grip with the directive sample.
//
// Limited.Span/smoke: a sample of the check refers to a declaration that
// the file of the suite cannot import. The remedy is to author a sample of
// each input of Span with the directive sample.
//
// Limited.Pitch/smoke: a sample of the check refers to a declaration that
// the file of the suite cannot import. The remedy is to author a sample of
// each input of Pitch with the directive sample.

// The unclaimed classifications of Limited:
//
// The shape:mutator of Limited.Limit has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Limited.Grip has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Limited.Clamp has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Limited.Frame has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Limited.Mount has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Limited.Span has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Limited.Pitch has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Limited.Measure has the stance silent. A mutation
// is an input of the other checks, and its result is an error alone.
//
// The shape:mutator of Limited.Phase has the stance silent. A mutation is
// an input of the other checks, and its result is an error alone.
//
// The shape:writer of Limited.Sense has the stance silent. A write is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Limited writes.
type (
	LimitedHandle = conform.Handle[Limited]
	LimitedCheck  = conform.Check[Limited]
	LimitedOption = conform.Option[Limited]
)

// LimitedChecks contains the name of each callable of Limited, with the
// identities of its checks.
var LimitedChecks = struct {
	Limit   limitedLimit
	Grip    limitedGrip
	Clamp   limitedClamp
	Frame   limitedFrame
	Mount   limitedMount
	Span    limitedSpan
	Pitch   limitedPitch
	Measure limitedMeasure
	Phase   limitedPhase
	Sense   limitedSense
}{
	Limit:   limitedLimit{"Limited.Limit"},
	Grip:    limitedGrip{"Limited.Grip"},
	Clamp:   limitedClamp{"Limited.Clamp"},
	Frame:   limitedFrame{"Limited.Frame"},
	Mount:   limitedMount{"Limited.Mount"},
	Span:    limitedSpan{"Limited.Span"},
	Pitch:   limitedPitch{"Limited.Pitch"},
	Measure: limitedMeasure{"Limited.Measure"},
	Phase:   limitedPhase{"Limited.Phase"},
	Sense:   limitedSense{"Limited.Sense"},
}

// limitedLimit is the name of Limited.Limit, with the identities of its
// checks.
type limitedLimit struct {
	conform.Callable[Limited]
}

// Smoke returns the identity of the check that Limit returns without a
// panic for derived inputs.
func (c limitedLimit) Smoke() conform.ID[Limited] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Limit returns
// without a panic for derived inputs.
func (c limitedLimit) SmokeGenerated() conform.ID[Limited] {
	return c.Check("smoke#generated")
}

// limitedGrip is the name of Limited.Grip, with the identities of its
// checks.
type limitedGrip struct {
	conform.Callable[Limited]
}

// SmokeGenerated returns the identity of the check that Grip returns
// without a panic for derived inputs.
func (c limitedGrip) SmokeGenerated() conform.ID[Limited] {
	return c.Check("smoke#generated")
}

// limitedClamp is the name of Limited.Clamp, with the identities of its
// checks.
type limitedClamp struct {
	conform.Callable[Limited]
}

// limitedFrame is the name of Limited.Frame, with the identities of its
// checks.
type limitedFrame struct {
	conform.Callable[Limited]
}

// limitedMount is the name of Limited.Mount, with the identities of its
// checks.
type limitedMount struct {
	conform.Callable[Limited]
}

// limitedSpan is the name of Limited.Span, with the identities of its
// checks.
type limitedSpan struct {
	conform.Callable[Limited]
}

// SmokeGenerated returns the identity of the check that Span returns
// without a panic for derived inputs.
func (c limitedSpan) SmokeGenerated() conform.ID[Limited] {
	return c.Check("smoke#generated")
}

// limitedPitch is the name of Limited.Pitch, with the identities of its
// checks.
type limitedPitch struct {
	conform.Callable[Limited]
}

// SmokeGenerated returns the identity of the check that Pitch returns
// without a panic for derived inputs.
func (c limitedPitch) SmokeGenerated() conform.ID[Limited] {
	return c.Check("smoke#generated")
}

// limitedMeasure is the name of Limited.Measure, with the identities of
// its checks.
type limitedMeasure struct {
	conform.Callable[Limited]
}

// Smoke returns the identity of the check that Measure returns without a
// panic for derived inputs.
func (c limitedMeasure) Smoke() conform.ID[Limited] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Measure returns
// without a panic for derived inputs.
func (c limitedMeasure) SmokeGenerated() conform.ID[Limited] {
	return c.Check("smoke#generated")
}

// limitedPhase is the name of Limited.Phase, with the identities of its
// checks.
type limitedPhase struct {
	conform.Callable[Limited]
}

// Smoke returns the identity of the check that Phase returns without a
// panic for derived inputs.
func (c limitedPhase) Smoke() conform.ID[Limited] {
	return c.Check("smoke")
}

// SmokeGenerated returns the identity of the check that Phase returns
// without a panic for derived inputs.
func (c limitedPhase) SmokeGenerated() conform.ID[Limited] {
	return c.Check("smoke#generated")
}

// limitedSense is the name of Limited.Sense, with the identities of its
// checks.
type limitedSense struct {
	conform.Callable[Limited]
}

// LimitedSuite returns the checks of Limited, and the calls and the
// budgets of its callables.
func LimitedSuite() conform.Suite[Limited] {
	return conform.Suite[Limited]{
		Declaration: "Limited",
		Package:     "example.com/shop/wrap",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Limited]{
			{
				ID:          LimitedChecks.Limit.Smoke(),
				Statement:   "Limit returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LimitedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Limit(tool.Limits{Count: 2})
						return nil
					})
				},
				Defect: limitedLimitPanics,
			},
			{
				ID:          LimitedChecks.Limit.SmokeGenerated(),
				Statement:   "Limit returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LimitedHandle, c *prop.Case) {
					drawnL := c.Draw(prop.Of[tool.Limits](), "l")
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Limit(drawnL)
						return nil
					})
				},
				Defect: limitedLimitPanics,
			},
			{
				ID:          LimitedChecks.Grip.SmokeGenerated(),
				Statement:   "Grip returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LimitedHandle, c *prop.Case) {
					drawnG := c.Draw(prop.Of[tool.Grip](), "g")
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Grip(drawnG)
						return nil
					})
				},
				Defect: limitedGripPanics,
			},
			{
				ID:          LimitedChecks.Span.SmokeGenerated(),
				Statement:   "Span returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LimitedHandle, c *prop.Case) {
					drawnS := c.Draw(prop.Of[tool.Span](), "s")
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Span(drawnS)
						return nil
					})
				},
				Defect: limitedSpanPanics,
			},
			{
				ID:          LimitedChecks.Pitch.SmokeGenerated(),
				Statement:   "Pitch returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LimitedHandle, c *prop.Case) {
					drawnP := c.Draw(prop.Of[tool.Pitch](), "p")
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Pitch(drawnP)
						return nil
					})
				},
				Defect: limitedPitchPanics,
			},
			{
				ID:          LimitedChecks.Measure.Smoke(),
				Statement:   "Measure returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LimitedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Measure(tool.Unit(1))
						return nil
					})
				},
				Defect: limitedMeasurePanics,
			},
			{
				ID:          LimitedChecks.Measure.SmokeGenerated(),
				Statement:   "Measure returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LimitedHandle, c *prop.Case) {
					drawnU := c.Draw(prop.Of[tool.Unit](), "u")
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Measure(drawnU)
						return nil
					})
				},
				Defect: limitedMeasurePanics,
			},
			{
				ID:          LimitedChecks.Phase.Smoke(),
				Statement:   "Phase returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *LimitedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Phase(tool.Phase(1))
						return nil
					})
				},
				Defect: limitedPhasePanics,
			},
			{
				ID:          LimitedChecks.Phase.SmokeGenerated(),
				Statement:   "Phase returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *LimitedHandle, c *prop.Case) {
					drawnP := c.Draw(prop.Of[tool.Phase](), "p")
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Phase(drawnP)
						return nil
					})
				},
				Defect: limitedPhasePanics,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Smoke,
				Callable: "Limited.Grip",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Grip with the directive sample",
			},
			{
				Claim:    spec.Smoke,
				Callable: "Limited.Span",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Span with the directive sample",
			},
			{
				Claim:    spec.Smoke,
				Callable: "Limited.Pitch",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Pitch with the directive sample",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:mutator",
				Callable:       "Limited.Limit",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Limited.Grip",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Limited.Clamp",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Limited.Frame",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Limited.Mount",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Limited.Span",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Limited.Pitch",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Limited.Measure",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:mutator",
				Callable:       "Limited.Phase",
				Stance:         spec.Silent,
				Note:           "A mutation is an input of the other checks, and its result is an error alone.",
			},
			{
				Classification: "shape:writer",
				Callable:       "Limited.Sense",
				Stance:         spec.Silent,
				Note:           "A write is an input of the other checks, and its result is an error alone.",
			},
		},
		Calls: []conform.Call[Limited]{
			{
				Calls:    limitedLimitCalls,
				Over:     newLimitedOver,
				Callable: LimitedChecks.Limit.Callable,
			},
			{
				Calls:    limitedGripCalls,
				Over:     newLimitedOver,
				Callable: LimitedChecks.Grip.Callable,
			},
			{
				Calls:    limitedClampCalls,
				Over:     newLimitedOver,
				Callable: LimitedChecks.Clamp.Callable,
			},
			{
				Calls:    limitedFrameCalls,
				Over:     newLimitedOver,
				Callable: LimitedChecks.Frame.Callable,
			},
			{
				Calls:    limitedMountCalls,
				Over:     newLimitedOver,
				Callable: LimitedChecks.Mount.Callable,
			},
			{
				Calls:    limitedSpanCalls,
				Over:     newLimitedOver,
				Callable: LimitedChecks.Span.Callable,
			},
			{
				Calls:    limitedPitchCalls,
				Over:     newLimitedOver,
				Callable: LimitedChecks.Pitch.Callable,
			},
			{
				Calls:    limitedMeasureCalls,
				Over:     newLimitedOver,
				Callable: LimitedChecks.Measure.Callable,
			},
			{
				Calls:    limitedPhaseCalls,
				Over:     newLimitedOver,
				Callable: LimitedChecks.Phase.Callable,
			},
			{
				Calls:    limitedSenseCalls,
				Over:     newLimitedOver,
				Callable: LimitedChecks.Sense.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.SampledFrom(tool.Percent, tool.Steps)),
			prop.Using(limitedLimitsValues),
			prop.Using(limitedGripValues),
			prop.Using(limitedClampValues),
			prop.Using(limitedFrameValues),
			prop.Using(limitedSpanValues),
		},
	}
}

// RunLimitedSuite applies the checks of Limited, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunLimitedSuite[T Limited](t *testing.T, newSubject func() T, opts ...LimitedOption) {
	t.Helper()
	RunLimitedSuiteFromHandle(t, func(*LimitedHandle) T { return newSubject() }, opts...)
}

// RunLimitedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunLimitedSuiteFromHandle[T Limited](
	t *testing.T,
	newSubject func(h *LimitedHandle) T,
	opts ...LimitedOption,
) {
	t.Helper()
	conform.Run(t, LimitedSuite(), limitedSubjects(newSubject), opts...)
}

// FuzzLimitedSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzLimitedSuite[T Limited](
	f *testing.F,
	callable conform.Callable[Limited],
	newSubject func(h *LimitedHandle) T,
	opts ...LimitedOption,
) {
	f.Helper()
	conform.Fuzz(f, LimitedSuite(), callable, limitedSubjects(newSubject), opts...)
}

// RunLimitedBudgets applies the allocation ceilings of Limited and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunLimitedBudgets[T Limited](
	t *testing.T,
	newSubject func(h *LimitedHandle) T,
	opts ...LimitedOption,
) {
	t.Helper()
	conform.RunBudgets(t, LimitedSuite(), limitedSubjects(newSubject), opts...)
}

// BenchLimitedSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchLimitedSuite[T Limited](
	b conform.B,
	callable conform.Callable[Limited],
	newSubject func(h *LimitedHandle) T,
	opts ...LimitedOption,
) {
	b.Helper()
	conform.Bench(b, LimitedSuite(), callable, limitedSubjects(newSubject), opts...)
}

// limitedSubjects returns the factory of the subjects of Limited that
// calls newSubject with the handle of each check.
func limitedSubjects[T Limited](
	newSubject func(h *LimitedHandle) T,
) func(h *LimitedHandle) Limited {
	return func(h *LimitedHandle) Limited { return newSubject(h) }
}

// limitedLimitCalls generates the calls of Limited.Limit with the inputs
// that each case draws.
var limitedLimitCalls = prop.Composite(func(c *prop.Case) func(context.Context, Limited) error {
	l := c.Draw(prop.Of[tool.Limits](), "l")
	return func(_ context.Context, subject Limited) error {
		subject.Limit(l)
		return nil
	}
})

// limitedGripCalls generates the calls of Limited.Grip with the inputs
// that each case draws.
var limitedGripCalls = prop.Composite(func(c *prop.Case) func(context.Context, Limited) error {
	g := c.Draw(prop.Of[tool.Grip](), "g")
	return func(_ context.Context, subject Limited) error {
		subject.Grip(g)
		return nil
	}
})

// limitedClampCalls generates the calls of Limited.Clamp with the inputs
// that each case draws.
var limitedClampCalls = prop.Composite(func(c *prop.Case) func(context.Context, Limited) error {
	c0 := c.Draw(prop.Of[tool.Clamp](), "c")
	return func(_ context.Context, subject Limited) error {
		subject.Clamp(c0)
		return nil
	}
})

// limitedFrameCalls generates the calls of Limited.Frame with the inputs
// that each case draws.
var limitedFrameCalls = prop.Composite(func(c *prop.Case) func(context.Context, Limited) error {
	f := c.Draw(prop.Of[tool.Frame](), "f")
	return func(_ context.Context, subject Limited) error {
		subject.Frame(f)
		return nil
	}
})

// limitedMountCalls generates the calls of Limited.Mount with the inputs
// that each case draws.
var limitedMountCalls = prop.Composite(func(c *prop.Case) func(context.Context, Limited) error {
	m := c.Draw(prop.Of[tool.Mount](), "m")
	return func(_ context.Context, subject Limited) error {
		subject.Mount(m)
		return nil
	}
})

// limitedSpanCalls generates the calls of Limited.Span with the inputs
// that each case draws.
var limitedSpanCalls = prop.Composite(func(c *prop.Case) func(context.Context, Limited) error {
	s := c.Draw(prop.Of[tool.Span](), "s")
	return func(_ context.Context, subject Limited) error {
		subject.Span(s)
		return nil
	}
})

// limitedPitchCalls generates the calls of Limited.Pitch with the inputs
// that each case draws.
var limitedPitchCalls = prop.Composite(func(c *prop.Case) func(context.Context, Limited) error {
	p := c.Draw(prop.Of[tool.Pitch](), "p")
	return func(_ context.Context, subject Limited) error {
		subject.Pitch(p)
		return nil
	}
})

// limitedMeasureCalls generates the calls of Limited.Measure with the
// inputs that each case draws.
var limitedMeasureCalls = prop.Composite(func(c *prop.Case) func(context.Context, Limited) error {
	u := c.Draw(prop.Of[tool.Unit](), "u")
	return func(_ context.Context, subject Limited) error {
		subject.Measure(u)
		return nil
	}
})

// limitedPhaseCalls generates the calls of Limited.Phase with the inputs
// that each case draws.
var limitedPhaseCalls = prop.Composite(func(c *prop.Case) func(context.Context, Limited) error {
	p := c.Draw(prop.Of[tool.Phase](), "p")
	return func(_ context.Context, subject Limited) error {
		subject.Phase(p)
		return nil
	}
})

// limitedSenseCalls generates the calls of Limited.Sense with the inputs
// that each case draws.
var limitedSenseCalls = prop.Composite(func(c *prop.Case) func(context.Context, Limited) error {
	sSenseResults := c.Draw(
		prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
		"s-sense",
	)
	var sSenseAt atomic.Int64
	sProbeResults := c.Draw(prop.List(prop.SampledFrom[error](nil, tool.ErrJammed)), "s-probe")
	var sProbeAt atomic.Int64
	sReadResults0 := c.Draw(prop.List(prop.Of[int]()), "s-read-0")
	sReadResults1 := c.Draw(
		prop.List(
			prop.SampledFrom[error](nil, conform.ErrGenerated),
			prop.MinSize(len(sReadResults0)),
			prop.MaxSize(len(sReadResults0)),
		),
		"s-read-1",
	)
	var sReadAt atomic.Int64
	s := &limitedSensorInput{
		onSense: func() (r0 error) {
			i := int(sSenseAt.Add(1) - 1)
			if i >= len(sSenseResults) {
				r0 = conform.ErrGenerated
				return r0
			}
			return sSenseResults[i]
		},
		onProbe: func() (r0 error) {
			i := int(sProbeAt.Add(1) - 1)
			if i >= len(sProbeResults) {
				r0 = tool.ErrJammed
				return r0
			}
			return sProbeResults[i]
		},
		onRead: func(string) (r0 int, r1 error) {
			i := int(sReadAt.Add(1) - 1)
			if i >= len(sReadResults0) {
				r1 = conform.ErrGenerated
				return r0, r1
			}
			return sReadResults0[i],
				sReadResults1[i]
		},
	}
	return func(_ context.Context, subject Limited) error {
		return subject.Sense(s)
	}
})

// limitedLimitsValues generates the values of tool.Limits that keep or
// change each field of its starting value.
var limitedLimitsValues = prop.Composite(func(c *prop.Case) tool.Limits {
	changed := c.Draw(prop.Of[tool.Limits](), "limits")
	value := tool.Limits{Count: 2}
	if c.Draw(prop.Boolean(), "limits.level") {
		value.Level = changed.Level
	}
	if c.Draw(prop.Boolean(), "limits.count") {
		value.Count = changed.Count
	}
	return value
})

// limitedGripValues generates the values of tool.Grip that keep or change
// each field of its starting value.
var limitedGripValues = prop.Composite(func(c *prop.Case) tool.Grip {
	changed := c.Draw(prop.Of[tool.Grip](), "grip")
	value := tool.Grip{Force: 4}
	if c.Draw(prop.Boolean(), "grip.force") {
		value.Force = changed.Force
	}
	return value
})

// limitedClampValues generates the starting value of tool.Clamp, which has
// no field to change.
var limitedClampValues = prop.Composite(func(*prop.Case) tool.Clamp {
	return tool.Clamp{}
})

// limitedFrameValues generates the starting value of tool.Frame, which has
// no field to change.
var limitedFrameValues = prop.Composite(func(*prop.Case) tool.Frame {
	return tool.Frame{}
})

// limitedSpanValues generates the values of tool.Span that keep or change
// its starting value.
var limitedSpanValues = prop.Composite(func(c *prop.Case) tool.Span {
	changed := c.Draw(prop.Of[tool.Span](), "span")
	if c.Draw(prop.Boolean(), "span.changed") {
		return changed
	}
	return 8
})

// limitedSensorInputEmbedded is tool.Sensor under the name of the field
// that limitedSensorInput embeds.
type limitedSensorInputEmbedded = tool.Sensor

// limitedSensorInput implements tool.Sensor. Its methods return the
// results that the generators of the suite of Limited draw for them, in
// turn.
type limitedSensorInput struct {
	limitedSensorInputEmbedded
	onSense func() error
	onProbe func() error
	onRead  func(a0 string) (int, error)
}

// Sense calls the function of Sense.
func (in limitedSensorInput) Sense() error {
	return in.onSense()
}

// Probe calls the function of Probe.
func (in limitedSensorInput) Probe() error {
	return in.onProbe()
}

// Read calls the function of Read.
func (in limitedSensorInput) Read(a0 string) (int, error) {
	return in.onRead(a0)
}

// limitedOverEmbedded is Limited under the name of the field that
// limitedOver embeds.
type limitedOverEmbedded = Limited

// limitedOver is a subject of Limited whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type limitedOver struct {
	limitedOverEmbedded
	spend func()
}

// Limit calls spend and returns its zero results.
func (o limitedOver) Limit(tool.Limits) {
	o.spend()
}

// Grip calls spend and returns its zero results.
func (o limitedOver) Grip(tool.Grip) {
	o.spend()
}

// Clamp calls spend and returns its zero results.
func (o limitedOver) Clamp(tool.Clamp) {
	o.spend()
}

// Frame calls spend and returns its zero results.
func (o limitedOver) Frame(tool.Frame) {
	o.spend()
}

// Mount calls spend and returns its zero results.
func (o limitedOver) Mount(tool.Mount) {
	o.spend()
}

// Span calls spend and returns its zero results.
func (o limitedOver) Span(tool.Span) {
	o.spend()
}

// Pitch calls spend and returns its zero results.
func (o limitedOver) Pitch(tool.Pitch) {
	o.spend()
}

// Measure calls spend and returns its zero results.
func (o limitedOver) Measure(tool.Unit) {
	o.spend()
}

// Phase calls spend and returns its zero results.
func (o limitedOver) Phase(tool.Phase) {
	o.spend()
}

// Sense calls spend and returns its zero results.
func (o limitedOver) Sense(tool.Sensor) (r0 error) {
	o.spend()
	return r0
}

// newLimitedOver returns the planted defect of a ceiling of Limited, whose
// callables call spend.
func newLimitedOver(_ *LimitedHandle, spend func()) Limited {
	return limitedOver{spend: spend}
}

// limitedLimitPanicsDefect is the planted defect of a Limited whose Limit
// panics.
type limitedLimitPanicsDefect struct {
	limitedOverEmbedded
}

// Limit is the override of the planted defect.
func (*limitedLimitPanicsDefect) Limit(p0 tool.Limits) {
	panic("the planted defect of a Limited whose Limit panics")
}

// limitedLimitPanics returns the planted defect of a Limited whose Limit
// panics over a fresh subject.
func limitedLimitPanics(h *LimitedHandle) Limited {
	return &limitedLimitPanicsDefect{limitedOverEmbedded: h.Fresh()}
}

// limitedGripPanicsDefect is the planted defect of a Limited whose Grip
// panics.
type limitedGripPanicsDefect struct {
	limitedOverEmbedded
}

// Grip is the override of the planted defect.
func (*limitedGripPanicsDefect) Grip(p0 tool.Grip) {
	panic("the planted defect of a Limited whose Grip panics")
}

// limitedGripPanics returns the planted defect of a Limited whose Grip
// panics over a fresh subject.
func limitedGripPanics(h *LimitedHandle) Limited {
	return &limitedGripPanicsDefect{limitedOverEmbedded: h.Fresh()}
}

// limitedSpanPanicsDefect is the planted defect of a Limited whose Span
// panics.
type limitedSpanPanicsDefect struct {
	limitedOverEmbedded
}

// Span is the override of the planted defect.
func (*limitedSpanPanicsDefect) Span(p0 tool.Span) {
	panic("the planted defect of a Limited whose Span panics")
}

// limitedSpanPanics returns the planted defect of a Limited whose Span
// panics over a fresh subject.
func limitedSpanPanics(h *LimitedHandle) Limited {
	return &limitedSpanPanicsDefect{limitedOverEmbedded: h.Fresh()}
}

// limitedPitchPanicsDefect is the planted defect of a Limited whose Pitch
// panics.
type limitedPitchPanicsDefect struct {
	limitedOverEmbedded
}

// Pitch is the override of the planted defect.
func (*limitedPitchPanicsDefect) Pitch(p0 tool.Pitch) {
	panic("the planted defect of a Limited whose Pitch panics")
}

// limitedPitchPanics returns the planted defect of a Limited whose Pitch
// panics over a fresh subject.
func limitedPitchPanics(h *LimitedHandle) Limited {
	return &limitedPitchPanicsDefect{limitedOverEmbedded: h.Fresh()}
}

// limitedMeasurePanicsDefect is the planted defect of a Limited whose
// Measure panics.
type limitedMeasurePanicsDefect struct {
	limitedOverEmbedded
}

// Measure is the override of the planted defect.
func (*limitedMeasurePanicsDefect) Measure(p0 tool.Unit) {
	panic("the planted defect of a Limited whose Measure panics")
}

// limitedMeasurePanics returns the planted defect of a Limited whose
// Measure panics over a fresh subject.
func limitedMeasurePanics(h *LimitedHandle) Limited {
	return &limitedMeasurePanicsDefect{limitedOverEmbedded: h.Fresh()}
}

// limitedPhasePanicsDefect is the planted defect of a Limited whose Phase
// panics.
type limitedPhasePanicsDefect struct {
	limitedOverEmbedded
}

// Phase is the override of the planted defect.
func (*limitedPhasePanicsDefect) Phase(p0 tool.Phase) {
	panic("the planted defect of a Limited whose Phase panics")
}

// limitedPhasePanics returns the planted defect of a Limited whose Phase
// panics over a fresh subject.
func limitedPhasePanics(h *LimitedHandle) Limited {
	return &limitedPhasePanicsDefect{limitedOverEmbedded: h.Fresh()}
}

// The derived checks of Played:
//
// Played.Play/smoke: Play returns without a panic for derived inputs.
//
// Played.Play/cancel: Play returns an error for a cancelled context.
//
// Played.StateAt/smoke#generated: StateAt returns without a panic for
// derived inputs.
//
// Played.StateAt/cancel#generated: StateAt returns an error for a
// cancelled context.
//
// Played.StateAt/zero-on-error#generated: StateAt returns the zero value
// beside an error.

// The derived checks of Played that witness their laws at the shallow end:
//
// Played.Play/smoke

// The refused claims of Played:
//
// Played.Play/follows-its-graph: a sample of the check refers to a
// declaration that the file of the suite cannot import. The remedy is to
// author a sample of each input of Play with the directive sample.
//
// Played.StateAt/smoke: a sample of the check refers to a declaration that
// the file of the suite cannot import. The remedy is to author a sample of
// each input of StateAt with the directive sample.
//
// Played.StateAt/cancel: a sample of the check refers to a declaration
// that the file of the suite cannot import. The remedy is to author a
// sample of each input of StateAt with the directive sample.
//
// Played.StateAt/zero-on-error: a sample of the check refers to a
// declaration that the file of the suite cannot import. The remedy is to
// author a sample of each input of StateAt with the directive sample.

// The unclaimed classifications of Played:
//
// The shape:lifecycle of Played.Play has the stance silent. A lifecycle
// call returns nothing that a check could compare.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Played writes.
type (
	PlayedHandle = conform.Handle[Played]
	PlayedCheck  = conform.Check[Played]
	PlayedOption = conform.Option[Played]
)

// PlayedChecks contains the name of each callable of Played, with the
// identities of its checks.
var PlayedChecks = struct {
	Play    playedPlay
	StateAt playedStateAt
}{
	Play:    playedPlay{"Played.Play"},
	StateAt: playedStateAt{"Played.StateAt"},
}

// playedPlay is the name of Played.Play, with the identities of its
// checks.
type playedPlay struct {
	conform.Callable[Played]
}

// Smoke returns the identity of the check that Play returns without a
// panic for derived inputs.
func (c playedPlay) Smoke() conform.ID[Played] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Play returns an error for
// a cancelled context.
func (c playedPlay) Cancel() conform.ID[Played] {
	return c.Check("cancel")
}

// playedStateAt is the name of Played.StateAt, with the identities of its
// checks.
type playedStateAt struct {
	conform.Callable[Played]
}

// SmokeGenerated returns the identity of the check that StateAt returns
// without a panic for derived inputs.
func (c playedStateAt) SmokeGenerated() conform.ID[Played] {
	return c.Check("smoke#generated")
}

// CancelGenerated returns the identity of the check that StateAt returns
// an error for a cancelled context.
func (c playedStateAt) CancelGenerated() conform.ID[Played] {
	return c.Check("cancel#generated")
}

// ZeroOnErrorGenerated returns the identity of the check that StateAt
// returns the zero value beside an error.
func (c playedStateAt) ZeroOnErrorGenerated() conform.ID[Played] {
	return c.Check("zero-on-error#generated")
}

// PlayedSuite returns the checks of Played, and the calls and the budgets
// of its callables.
func PlayedSuite() conform.Suite[Played] {
	return conform.Suite[Played]{
		Declaration: "Played",
		Package:     "example.com/shop/wrap",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Played]{
			{
				ID:          PlayedChecks.Play.Smoke(),
				Statement:   "Play returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PlayedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Play(ctx)
					})
				},
				Defect: playedPlayPanics,
			},
			{
				ID:          PlayedChecks.Play.Cancel(),
				Statement:   "Play returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *PlayedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Play(ctx)
					})
				},
				Defect: playedPlayIgnoresCancel,
			},
			{
				ID:          PlayedChecks.StateAt.SmokeGenerated(),
				Statement:   "StateAt returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *PlayedHandle, c *prop.Case) {
					drawnP := c.Draw(prop.Of[tool.Pitch](), "p")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().StateAt(ctx, drawnP)
						return err
					})
				},
				Defect: playedStateAtPanics,
			},
			{
				ID:          PlayedChecks.StateAt.CancelGenerated(),
				Statement:   "StateAt returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *PlayedHandle, c *prop.Case) {
					drawnP := c.Draw(prop.Of[tool.Pitch](), "p")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().StateAt(ctx, drawnP)
						return err
					})
				},
				Defect: playedStateAtIgnoresCancel,
			},
			{
				ID:          PlayedChecks.StateAt.ZeroOnErrorGenerated(),
				Statement:   "StateAt returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *PlayedHandle, c *prop.Case) {
					drawnP := c.Draw(prop.Of[tool.Pitch](), "p")
					probe.ZeroBesideError(h, func(ctx context.Context) (string, error) {
						return h.Subject().StateAt(ctx, drawnP)
					}, string(""))
				},
				Defect: playedStateAtValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.FollowsItsGraph,
				Callable: "Played.Play",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Play with the directive sample",
			},
			{
				Claim:    spec.Smoke,
				Callable: "Played.StateAt",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of StateAt with the directive sample",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Played.StateAt",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of StateAt with the directive sample",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Played.StateAt",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of StateAt with the directive sample",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:lifecycle",
				Callable:       "Played.Play",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
		},
		Calls: []conform.Call[Played]{
			{
				Calls:    playedPlayCalls,
				Over:     newPlayedOver,
				Callable: PlayedChecks.Play.Callable,
			},
			{
				Calls:    playedStateAtCalls,
				Over:     newPlayedOver,
				Callable: PlayedChecks.StateAt.Callable,
			},
		},
	}
}

// RunPlayedSuite applies the checks of Played, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunPlayedSuite[T Played](t *testing.T, newSubject func() T, opts ...PlayedOption) {
	t.Helper()
	RunPlayedSuiteFromHandle(t, func(*PlayedHandle) T { return newSubject() }, opts...)
}

// RunPlayedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunPlayedSuiteFromHandle[T Played](
	t *testing.T,
	newSubject func(h *PlayedHandle) T,
	opts ...PlayedOption,
) {
	t.Helper()
	conform.Run(t, PlayedSuite(), playedSubjects(newSubject), opts...)
}

// FuzzPlayedSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzPlayedSuite[T Played](
	f *testing.F,
	callable conform.Callable[Played],
	newSubject func(h *PlayedHandle) T,
	opts ...PlayedOption,
) {
	f.Helper()
	conform.Fuzz(f, PlayedSuite(), callable, playedSubjects(newSubject), opts...)
}

// RunPlayedBudgets applies the allocation ceilings of Played and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunPlayedBudgets[T Played](
	t *testing.T,
	newSubject func(h *PlayedHandle) T,
	opts ...PlayedOption,
) {
	t.Helper()
	conform.RunBudgets(t, PlayedSuite(), playedSubjects(newSubject), opts...)
}

// BenchPlayedSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchPlayedSuite[T Played](
	b conform.B,
	callable conform.Callable[Played],
	newSubject func(h *PlayedHandle) T,
	opts ...PlayedOption,
) {
	b.Helper()
	conform.Bench(b, PlayedSuite(), callable, playedSubjects(newSubject), opts...)
}

// playedSubjects returns the factory of the subjects of Played that calls
// newSubject with the handle of each check.
func playedSubjects[T Played](newSubject func(h *PlayedHandle) T) func(h *PlayedHandle) Played {
	return func(h *PlayedHandle) Played { return newSubject(h) }
}

// playedPlayCalls generates the calls of Played.Play with the inputs that
// each case draws.
var playedPlayCalls = prop.Composite(func(*prop.Case) func(context.Context, Played) error {
	return func(ctx context.Context, subject Played) error {
		return subject.Play(ctx)
	}
})

// playedStateAtCalls generates the calls of Played.StateAt with the inputs
// that each case draws.
var playedStateAtCalls = prop.Composite(func(c *prop.Case) func(context.Context, Played) error {
	p := c.Draw(prop.Of[tool.Pitch](), "p")
	return func(ctx context.Context, subject Played) error {
		_, err := subject.StateAt(ctx, p)
		return err
	}
})

// playedOverEmbedded is Played under the name of the field that playedOver
// embeds.
type playedOverEmbedded = Played

// playedOver is a subject of Played whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type playedOver struct {
	playedOverEmbedded
	spend func()
}

// Play calls spend and returns its zero results.
func (o playedOver) Play(context.Context) (r0 error) {
	o.spend()
	return r0
}

// StateAt calls spend and returns its zero results.
func (o playedOver) StateAt(context.Context, tool.Pitch) (r0 string, r1 error) {
	o.spend()
	return r0, r1
}

// newPlayedOver returns the planted defect of a ceiling of Played, whose
// callables call spend.
func newPlayedOver(_ *PlayedHandle, spend func()) Played {
	return playedOver{spend: spend}
}

// playedPlayPanicsDefect is the planted defect of a Played whose Play
// panics.
type playedPlayPanicsDefect struct {
	playedOverEmbedded
}

// Play is the override of the planted defect.
func (*playedPlayPanicsDefect) Play(p0 context.Context) (r0 error) {
	panic("the planted defect of a Played whose Play panics")
}

// playedPlayPanics returns the planted defect of a Played whose Play
// panics over a fresh subject.
func playedPlayPanics(h *PlayedHandle) Played {
	return &playedPlayPanicsDefect{playedOverEmbedded: h.Fresh()}
}

// playedPlayIgnoresCancelDefect is the planted defect of a Played whose
// Play ignores the context.
type playedPlayIgnoresCancelDefect struct {
	playedOverEmbedded
}

// Play is the override of the planted defect.
func (d *playedPlayIgnoresCancelDefect) Play(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.playedOverEmbedded.Play(p0)
}

// playedPlayIgnoresCancel returns the planted defect of a Played whose
// Play ignores the context over a fresh subject.
func playedPlayIgnoresCancel(h *PlayedHandle) Played {
	return &playedPlayIgnoresCancelDefect{playedOverEmbedded: h.Fresh()}
}

// playedStateAtPanicsDefect is the planted defect of a Played whose
// StateAt panics.
type playedStateAtPanicsDefect struct {
	playedOverEmbedded
}

// StateAt is the override of the planted defect.
func (*playedStateAtPanicsDefect) StateAt(p0 context.Context, p1 tool.Pitch) (r0 string, r1 error) {
	panic("the planted defect of a Played whose StateAt panics")
}

// playedStateAtPanics returns the planted defect of a Played whose StateAt
// panics over a fresh subject.
func playedStateAtPanics(h *PlayedHandle) Played {
	return &playedStateAtPanicsDefect{playedOverEmbedded: h.Fresh()}
}

// playedStateAtIgnoresCancelDefect is the planted defect of a Played whose
// StateAt ignores the context.
type playedStateAtIgnoresCancelDefect struct {
	playedOverEmbedded
}

// StateAt is the override of the planted defect.
func (d *playedStateAtIgnoresCancelDefect) StateAt(
	p0 context.Context,
	p1 tool.Pitch,
) (r0 string, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.playedOverEmbedded.StateAt(p0, p1)
}

// playedStateAtIgnoresCancel returns the planted defect of a Played whose
// StateAt ignores the context over a fresh subject.
func playedStateAtIgnoresCancel(h *PlayedHandle) Played {
	return &playedStateAtIgnoresCancelDefect{playedOverEmbedded: h.Fresh()}
}

// playedStateAtValueBesideDefect is the planted defect of a Played whose
// StateAt returns a value beside an error.
type playedStateAtValueBesideDefect struct {
	playedOverEmbedded
}

// StateAt is the override of the planted defect.
func (d *playedStateAtValueBesideDefect) StateAt(
	p0 context.Context,
	p1 tool.Pitch,
) (r0 string, r1 error) {
	r0, r1 = d.playedOverEmbedded.StateAt(p0, p1)
	if r1 != nil {
		r0 = string("sample-b")
	}
	return r0, r1
}

// playedStateAtValueBeside returns the planted defect of a Played whose
// StateAt returns a value beside an error over a fresh subject.
func playedStateAtValueBeside(h *PlayedHandle) Played {
	return &playedStateAtValueBesideDefect{playedOverEmbedded: h.Fresh()}
}

// The derived checks of Gauged:
//
// Gauged.Fill/smoke: Fill returns without a panic for derived inputs.
//
// Gauged.Fill/cancel: Fill returns an error for a cancelled context.
//
// Gauged.LevelAt/smoke#generated: LevelAt returns without a panic for
// derived inputs.
//
// Gauged.LevelAt/cancel#generated: LevelAt returns an error for a
// cancelled context.
//
// Gauged.LevelAt/zero-on-error#generated: LevelAt returns the zero value
// beside an error.

// The derived checks of Gauged that witness their laws at the shallow end:
//
// Gauged.Fill/smoke

// The refused claims of Gauged:
//
// Gauged.Fill/keeps-within-its-bounds: a sample of the check refers to a
// declaration that the file of the suite cannot import. The remedy is to
// author a sample of each input of Fill with the directive sample.
//
// Gauged.LevelAt/smoke: a sample of the check refers to a declaration that
// the file of the suite cannot import. The remedy is to author a sample of
// each input of LevelAt with the directive sample.
//
// Gauged.LevelAt/cancel: a sample of the check refers to a declaration
// that the file of the suite cannot import. The remedy is to author a
// sample of each input of LevelAt with the directive sample.
//
// Gauged.LevelAt/zero-on-error: a sample of the check refers to a
// declaration that the file of the suite cannot import. The remedy is to
// author a sample of each input of LevelAt with the directive sample.

// The unclaimed classifications of Gauged:
//
// The shape:lifecycle of Gauged.Fill has the stance silent. A lifecycle
// call returns nothing that a check could compare.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Gauged writes.
type (
	GaugedHandle = conform.Handle[Gauged]
	GaugedCheck  = conform.Check[Gauged]
	GaugedOption = conform.Option[Gauged]
)

// GaugedChecks contains the name of each callable of Gauged, with the
// identities of its checks.
var GaugedChecks = struct {
	Fill    gaugedFill
	LevelAt gaugedLevelAt
}{
	Fill:    gaugedFill{"Gauged.Fill"},
	LevelAt: gaugedLevelAt{"Gauged.LevelAt"},
}

// gaugedFill is the name of Gauged.Fill, with the identities of its
// checks.
type gaugedFill struct {
	conform.Callable[Gauged]
}

// Smoke returns the identity of the check that Fill returns without a
// panic for derived inputs.
func (c gaugedFill) Smoke() conform.ID[Gauged] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Fill returns an error for
// a cancelled context.
func (c gaugedFill) Cancel() conform.ID[Gauged] {
	return c.Check("cancel")
}

// gaugedLevelAt is the name of Gauged.LevelAt, with the identities of its
// checks.
type gaugedLevelAt struct {
	conform.Callable[Gauged]
}

// SmokeGenerated returns the identity of the check that LevelAt returns
// without a panic for derived inputs.
func (c gaugedLevelAt) SmokeGenerated() conform.ID[Gauged] {
	return c.Check("smoke#generated")
}

// CancelGenerated returns the identity of the check that LevelAt returns
// an error for a cancelled context.
func (c gaugedLevelAt) CancelGenerated() conform.ID[Gauged] {
	return c.Check("cancel#generated")
}

// ZeroOnErrorGenerated returns the identity of the check that LevelAt
// returns the zero value beside an error.
func (c gaugedLevelAt) ZeroOnErrorGenerated() conform.ID[Gauged] {
	return c.Check("zero-on-error#generated")
}

// GaugedSuite returns the checks of Gauged, and the calls and the budgets
// of its callables.
func GaugedSuite() conform.Suite[Gauged] {
	return conform.Suite[Gauged]{
		Declaration: "Gauged",
		Package:     "example.com/shop/wrap",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[Gauged]{
			{
				ID:          GaugedChecks.Fill.Smoke(),
				Statement:   "Fill returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *GaugedHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Fill(ctx)
					})
				},
				Defect: gaugedFillPanics,
			},
			{
				ID:          GaugedChecks.Fill.Cancel(),
				Statement:   "Fill returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *GaugedHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Fill(ctx)
					})
				},
				Defect: gaugedFillIgnoresCancel,
			},
			{
				ID:          GaugedChecks.LevelAt.SmokeGenerated(),
				Statement:   "LevelAt returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *GaugedHandle, c *prop.Case) {
					drawnP := c.Draw(prop.Of[tool.Pitch](), "p")
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().LevelAt(ctx, drawnP)
						return err
					})
				},
				Defect: gaugedLevelAtPanics,
			},
			{
				ID:          GaugedChecks.LevelAt.CancelGenerated(),
				Statement:   "LevelAt returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *GaugedHandle, c *prop.Case) {
					drawnP := c.Draw(prop.Of[tool.Pitch](), "p")
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().LevelAt(ctx, drawnP)
						return err
					})
				},
				Defect: gaugedLevelAtIgnoresCancel,
			},
			{
				ID:          GaugedChecks.LevelAt.ZeroOnErrorGenerated(),
				Statement:   "LevelAt returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
				Property: func(h *GaugedHandle, c *prop.Case) {
					drawnP := c.Draw(prop.Of[tool.Pitch](), "p")
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject().LevelAt(ctx, drawnP)
					}, int(0))
				},
				Defect: gaugedLevelAtValueBeside,
			},
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.KeepsWithinItsBounds,
				Callable: "Gauged.Fill",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Fill with the directive sample",
			},
			{
				Claim:    spec.Smoke,
				Callable: "Gauged.LevelAt",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of LevelAt with the directive sample",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Gauged.LevelAt",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of LevelAt with the directive sample",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Gauged.LevelAt",
				Gate:     spec.ImportableSamples,
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of LevelAt with the directive sample",
			},
		},
		Unclaimed: []conform.Unclaimed{
			{
				Classification: "shape:lifecycle",
				Callable:       "Gauged.Fill",
				Stance:         spec.Silent,
				Note:           "A lifecycle call returns nothing that a check could compare.",
			},
		},
		Calls: []conform.Call[Gauged]{
			{
				Calls:    gaugedFillCalls,
				Over:     newGaugedOver,
				Callable: GaugedChecks.Fill.Callable,
			},
			{
				Calls:    gaugedLevelAtCalls,
				Over:     newGaugedOver,
				Callable: GaugedChecks.LevelAt.Callable,
			},
		},
	}
}

// RunGaugedSuite applies the checks of Gauged, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunGaugedSuite[T Gauged](t *testing.T, newSubject func() T, opts ...GaugedOption) {
	t.Helper()
	RunGaugedSuiteFromHandle(t, func(*GaugedHandle) T { return newSubject() }, opts...)
}

// RunGaugedSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunGaugedSuiteFromHandle[T Gauged](
	t *testing.T,
	newSubject func(h *GaugedHandle) T,
	opts ...GaugedOption,
) {
	t.Helper()
	conform.Run(t, GaugedSuite(), gaugedSubjects(newSubject), opts...)
}

// FuzzGaugedSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzGaugedSuite[T Gauged](
	f *testing.F,
	callable conform.Callable[Gauged],
	newSubject func(h *GaugedHandle) T,
	opts ...GaugedOption,
) {
	f.Helper()
	conform.Fuzz(f, GaugedSuite(), callable, gaugedSubjects(newSubject), opts...)
}

// RunGaugedBudgets applies the allocation ceilings of Gauged and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunGaugedBudgets[T Gauged](
	t *testing.T,
	newSubject func(h *GaugedHandle) T,
	opts ...GaugedOption,
) {
	t.Helper()
	conform.RunBudgets(t, GaugedSuite(), gaugedSubjects(newSubject), opts...)
}

// BenchGaugedSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchGaugedSuite[T Gauged](
	b conform.B,
	callable conform.Callable[Gauged],
	newSubject func(h *GaugedHandle) T,
	opts ...GaugedOption,
) {
	b.Helper()
	conform.Bench(b, GaugedSuite(), callable, gaugedSubjects(newSubject), opts...)
}

// gaugedSubjects returns the factory of the subjects of Gauged that calls
// newSubject with the handle of each check.
func gaugedSubjects[T Gauged](newSubject func(h *GaugedHandle) T) func(h *GaugedHandle) Gauged {
	return func(h *GaugedHandle) Gauged { return newSubject(h) }
}

// gaugedFillCalls generates the calls of Gauged.Fill with the inputs that
// each case draws.
var gaugedFillCalls = prop.Composite(func(*prop.Case) func(context.Context, Gauged) error {
	return func(ctx context.Context, subject Gauged) error {
		return subject.Fill(ctx)
	}
})

// gaugedLevelAtCalls generates the calls of Gauged.LevelAt with the inputs
// that each case draws.
var gaugedLevelAtCalls = prop.Composite(func(c *prop.Case) func(context.Context, Gauged) error {
	p := c.Draw(prop.Of[tool.Pitch](), "p")
	return func(ctx context.Context, subject Gauged) error {
		_, err := subject.LevelAt(ctx, p)
		return err
	}
})

// gaugedOverEmbedded is Gauged under the name of the field that gaugedOver
// embeds.
type gaugedOverEmbedded = Gauged

// gaugedOver is a subject of Gauged whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type gaugedOver struct {
	gaugedOverEmbedded
	spend func()
}

// Fill calls spend and returns its zero results.
func (o gaugedOver) Fill(context.Context) (r0 error) {
	o.spend()
	return r0
}

// LevelAt calls spend and returns its zero results.
func (o gaugedOver) LevelAt(context.Context, tool.Pitch) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// newGaugedOver returns the planted defect of a ceiling of Gauged, whose
// callables call spend.
func newGaugedOver(_ *GaugedHandle, spend func()) Gauged {
	return gaugedOver{spend: spend}
}

// gaugedFillPanicsDefect is the planted defect of a Gauged whose Fill
// panics.
type gaugedFillPanicsDefect struct {
	gaugedOverEmbedded
}

// Fill is the override of the planted defect.
func (*gaugedFillPanicsDefect) Fill(p0 context.Context) (r0 error) {
	panic("the planted defect of a Gauged whose Fill panics")
}

// gaugedFillPanics returns the planted defect of a Gauged whose Fill
// panics over a fresh subject.
func gaugedFillPanics(h *GaugedHandle) Gauged {
	return &gaugedFillPanicsDefect{gaugedOverEmbedded: h.Fresh()}
}

// gaugedFillIgnoresCancelDefect is the planted defect of a Gauged whose
// Fill ignores the context.
type gaugedFillIgnoresCancelDefect struct {
	gaugedOverEmbedded
}

// Fill is the override of the planted defect.
func (d *gaugedFillIgnoresCancelDefect) Fill(p0 context.Context) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.gaugedOverEmbedded.Fill(p0)
}

// gaugedFillIgnoresCancel returns the planted defect of a Gauged whose
// Fill ignores the context over a fresh subject.
func gaugedFillIgnoresCancel(h *GaugedHandle) Gauged {
	return &gaugedFillIgnoresCancelDefect{gaugedOverEmbedded: h.Fresh()}
}

// gaugedLevelAtPanicsDefect is the planted defect of a Gauged whose
// LevelAt panics.
type gaugedLevelAtPanicsDefect struct {
	gaugedOverEmbedded
}

// LevelAt is the override of the planted defect.
func (*gaugedLevelAtPanicsDefect) LevelAt(p0 context.Context, p1 tool.Pitch) (r0 int, r1 error) {
	panic("the planted defect of a Gauged whose LevelAt panics")
}

// gaugedLevelAtPanics returns the planted defect of a Gauged whose LevelAt
// panics over a fresh subject.
func gaugedLevelAtPanics(h *GaugedHandle) Gauged {
	return &gaugedLevelAtPanicsDefect{gaugedOverEmbedded: h.Fresh()}
}

// gaugedLevelAtIgnoresCancelDefect is the planted defect of a Gauged whose
// LevelAt ignores the context.
type gaugedLevelAtIgnoresCancelDefect struct {
	gaugedOverEmbedded
}

// LevelAt is the override of the planted defect.
func (d *gaugedLevelAtIgnoresCancelDefect) LevelAt(
	p0 context.Context,
	p1 tool.Pitch,
) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.gaugedOverEmbedded.LevelAt(p0, p1)
}

// gaugedLevelAtIgnoresCancel returns the planted defect of a Gauged whose
// LevelAt ignores the context over a fresh subject.
func gaugedLevelAtIgnoresCancel(h *GaugedHandle) Gauged {
	return &gaugedLevelAtIgnoresCancelDefect{gaugedOverEmbedded: h.Fresh()}
}

// gaugedLevelAtValueBesideDefect is the planted defect of a Gauged whose
// LevelAt returns a value beside an error.
type gaugedLevelAtValueBesideDefect struct {
	gaugedOverEmbedded
}

// LevelAt is the override of the planted defect.
func (d *gaugedLevelAtValueBesideDefect) LevelAt(
	p0 context.Context,
	p1 tool.Pitch,
) (r0 int, r1 error) {
	r0, r1 = d.gaugedOverEmbedded.LevelAt(p0, p1)
	if r1 != nil {
		r0 = int(7)
	}
	return r0, r1
}

// gaugedLevelAtValueBeside returns the planted defect of a Gauged whose
// LevelAt returns a value beside an error over a fresh subject.
func gaugedLevelAtValueBeside(h *GaugedHandle) Gauged {
	return &gaugedLevelAtValueBesideDefect{gaugedOverEmbedded: h.Fresh()}
}

// The derived checks of Plain:
//
// Plain.Name/smoke: Name returns without a panic for derived inputs.
//
// Plain.Name/deterministic: Name returns the same value for the same
// inputs.

// The derived checks of Plain that witness their laws at the shallow end:
//
// Plain.Name/smoke
//
// Plain.Name/deterministic

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Plain writes.
type (
	PlainHandle = conform.Handle[PlainMethods]
	PlainCheck  = conform.Check[PlainMethods]
	PlainOption = conform.Option[PlainMethods]
)

// PlainMethods is the interface of the exported methods of Plain. Each
// subject of the suite of Plain implements it.
type PlainMethods interface {
	Name() string
}

// PlainChecks contains the name of each callable of Plain, with the
// identities of its checks.
var PlainChecks = struct {
	Name plainName
}{
	Name: plainName{"Plain.Name"},
}

// plainName is the name of Plain.Name, with the identities of its checks.
type plainName struct {
	conform.Callable[PlainMethods]
}

// Smoke returns the identity of the check that Name returns without a
// panic for derived inputs.
func (c plainName) Smoke() conform.ID[PlainMethods] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that Name returns the
// same value for the same inputs.
func (c plainName) Deterministic() conform.ID[PlainMethods] {
	return c.Check("deterministic")
}

// PlainSuite returns the checks of Plain, and the calls and the budgets of
// its callables.
func PlainSuite() conform.Suite[PlainMethods] {
	return conform.Suite[PlainMethods]{
		Declaration: "Plain",
		Package:     "example.com/shop/wrap",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[PlainMethods]{
			{
				ID:          PlainChecks.Name.Smoke(),
				Statement:   "Name returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *PlainHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Name()
						return nil
					})
				},
				Defect: plainNamePanics,
			},
			{
				ID:          PlainChecks.Name.Deterministic(),
				Statement:   "Name returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[PlainMethods]{PlainChecks.Name.Smoke()},
				Apply: func(h *PlainHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) (string, error) {
						r0 := h.Subject().Name()
						return r0, nil
					})
				},
				Defect: plainNameDrifts,
			},
		},
		Calls: []conform.Call[PlainMethods]{
			{
				Calls:    plainNameCalls,
				Over:     newPlainOver,
				Callable: PlainChecks.Name.Callable,
			},
		},
	}
}

// RunPlainSuite applies the checks of Plain, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunPlainSuite[T PlainMethods](t *testing.T, newSubject func() T, opts ...PlainOption) {
	t.Helper()
	RunPlainSuiteFromHandle(t, func(*PlainHandle) T { return newSubject() }, opts...)
}

// RunPlainSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunPlainSuiteFromHandle[T PlainMethods](
	t *testing.T,
	newSubject func(h *PlainHandle) T,
	opts ...PlainOption,
) {
	t.Helper()
	conform.Run(t, PlainSuite(), plainSubjects(newSubject), opts...)
}

// FuzzPlainSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzPlainSuite[T PlainMethods](
	f *testing.F,
	callable conform.Callable[PlainMethods],
	newSubject func(h *PlainHandle) T,
	opts ...PlainOption,
) {
	f.Helper()
	conform.Fuzz(f, PlainSuite(), callable, plainSubjects(newSubject), opts...)
}

// RunPlainBudgets applies the allocation ceilings of Plain and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunPlainBudgets[T PlainMethods](
	t *testing.T,
	newSubject func(h *PlainHandle) T,
	opts ...PlainOption,
) {
	t.Helper()
	conform.RunBudgets(t, PlainSuite(), plainSubjects(newSubject), opts...)
}

// BenchPlainSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchPlainSuite[T PlainMethods](
	b conform.B,
	callable conform.Callable[PlainMethods],
	newSubject func(h *PlainHandle) T,
	opts ...PlainOption,
) {
	b.Helper()
	conform.Bench(b, PlainSuite(), callable, plainSubjects(newSubject), opts...)
}

// plainSubjects returns the factory of the subjects of Plain that calls
// newSubject with the handle of each check.
func plainSubjects[T PlainMethods](
	newSubject func(h *PlainHandle) T,
) func(h *PlainHandle) PlainMethods {
	return func(h *PlainHandle) PlainMethods { return newSubject(h) }
}

// plainNameCalls generates the calls of Plain.Name with the inputs that
// each case draws.
var plainNameCalls = prop.Composite(func(*prop.Case) func(context.Context, PlainMethods) error {
	return func(_ context.Context, subject PlainMethods) error {
		_ = subject.Name()
		return nil
	}
})

// plainOverEmbedded is PlainMethods under the name of the field that
// plainOver embeds.
type plainOverEmbedded = PlainMethods

// plainOver is a subject of Plain whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type plainOver struct {
	plainOverEmbedded
	spend func()
}

// Name calls spend and returns its zero results.
func (o plainOver) Name() (r0 string) {
	o.spend()
	return r0
}

// newPlainOver returns the planted defect of a ceiling of Plain, whose
// callables call spend.
func newPlainOver(_ *PlainHandle, spend func()) PlainMethods {
	return plainOver{spend: spend}
}

// plainNamePanicsDefect is the planted defect of a Plain whose Name
// panics.
type plainNamePanicsDefect struct {
	plainOverEmbedded
}

// Name is the override of the planted defect.
func (*plainNamePanicsDefect) Name() (r0 string) {
	panic("the planted defect of a Plain whose Name panics")
}

// plainNamePanics returns the planted defect of a Plain whose Name panics
// over a fresh subject.
func plainNamePanics(h *PlainHandle) PlainMethods {
	return &plainNamePanicsDefect{plainOverEmbedded: h.Fresh()}
}

// plainNameDriftsDefect is the planted defect of a Plain whose Name
// returns another value on the repeat.
type plainNameDriftsDefect struct {
	plainOverEmbedded
	calls int
}

// Name is the override of the planted defect.
func (d *plainNameDriftsDefect) Name() (r0 string) {
	d.calls++
	if d.calls > 1 {
		r0 = string("sample-b")
	}
	return r0
}

// plainNameDrifts returns the planted defect of a Plain whose Name returns
// another value on the repeat over a fresh subject.
func plainNameDrifts(h *PlainHandle) PlainMethods {
	return &plainNameDriftsDefect{plainOverEmbedded: h.Fresh()}
}
