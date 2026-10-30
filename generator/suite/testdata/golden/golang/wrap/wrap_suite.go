package wrap

import (
	"context"
	"sync/atomic"
	"testing"

	"example.com/shop/lib/tool"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Reach:
//
// Reach.Current/smoke: Current returns without a panic for derived inputs.
//
// Reach.Current/cancel: Current returns an error for a cancelled context.
//
// Reach.Current/zero-on-error: Current returns the zero value beside an
// error.

// The derived checks of Reach that witness their laws at the shallow end:
//
// Reach.Current/smoke

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
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Smoke,
				Callable: "Reach.Store",
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Store with the directive sample",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Reach.Store",
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Store with the directive sample",
			},
			{
				Claim:    spec.ReadAfterWrite,
				Callable: "Reach.Current",
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Current with the directive sample",
			},
			{
				Claim:    spec.Smoke,
				Callable: "Reach.Gear",
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Gear with the directive sample",
			},
			{
				Claim:    spec.Cancel,
				Callable: "Reach.Gear",
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Gear with the directive sample",
			},
			{
				Claim:    spec.ZeroOnError,
				Callable: "Reach.Gear",
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
// Limited.Measure/smoke: Measure returns without a panic for derived
// inputs.
//
// Limited.Phase/smoke: Phase returns without a panic for derived inputs.

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

// limitedGrip is the name of Limited.Grip, with the identities of its
// checks.
type limitedGrip struct {
	conform.Callable[Limited]
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

// limitedPitch is the name of Limited.Pitch, with the identities of its
// checks.
type limitedPitch struct {
	conform.Callable[Limited]
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
		},
		Refusals: []conform.Refusal{
			{
				Claim:    spec.Smoke,
				Callable: "Limited.Grip",
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Grip with the directive sample",
			},
			{
				Claim:    spec.Smoke,
				Callable: "Limited.Span",
				Why:      "a sample of the check refers to a declaration that the file of the suite cannot import",
				Remedy:   "author a sample of each input of Span with the directive sample",
			},
			{
				Claim:    spec.Smoke,
				Callable: "Limited.Pitch",
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
