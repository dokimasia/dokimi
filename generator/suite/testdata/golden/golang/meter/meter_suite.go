package meter

import (
	"context"
	"testing"
	"time"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Meter:
//
// Meter.Count/smoke: Count returns without a panic for derived inputs.
//
// Meter.Count/cancel: Count returns an error for a cancelled context.
//
// Meter.Count/zero-on-error: Count returns the zero value beside an error.
//
// Meter.Record/smoke: Record returns without a panic for derived inputs.
//
// Meter.Record/cancel: Record returns an error for a cancelled context.
//
// Meter.Merge/smoke: Merge returns without a panic for derived inputs.
//
// Meter.Merge/nil-safe: Merge returns without a panic for a nil argument.
//
// Meter.Total/smoke: Total returns without a panic for derived inputs.
//
// Meter.Total/deterministic: Total returns the same value for the same
// inputs.

// The derived checks of Meter that witness their laws at the shallow end:
//
// Meter.Count/smoke
//
// Meter.Record/smoke
//
// Meter.Merge/smoke
//
// Meter.Total/smoke
//
// Meter.Total/deterministic

// The unclaimed classifications of Meter:
//
// The mixin:concurrent of Meter.Count has the stance needs. The claims of
// the mixin need concurrent callers.
//
// The shape:writer of Meter.Record has the stance silent. A write is an
// input of the other checks, and its result is an error alone.
//
// The shape:mutator of Meter.Merge has the stance silent. A mutation is an
// input of the other checks, and its result is an error alone.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Meter writes.
type (
	MeterHandle = conform.Handle[Meter]
	MeterCheck  = conform.Check[Meter]
	MeterOption = conform.Option[Meter]
)

// MeterChecks contains the name of each callable of Meter, with the
// identities of its checks.
var MeterChecks = struct {
	Count  meterCount
	Record meterRecord
	Merge  meterMerge
	Total  meterTotal
}{
	Count:  meterCount{"Meter.Count"},
	Record: meterRecord{"Meter.Record"},
	Merge:  meterMerge{"Meter.Merge"},
	Total:  meterTotal{"Meter.Total"},
}

// meterCount is the name of Meter.Count, with the identities of its
// checks.
type meterCount struct {
	conform.Callable[Meter]
}

// Smoke returns the identity of the check that Count returns without a
// panic for derived inputs.
func (c meterCount) Smoke() conform.ID[Meter] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Count returns an error for
// a cancelled context.
func (c meterCount) Cancel() conform.ID[Meter] {
	return c.Check("cancel")
}

// ZeroOnError returns the identity of the check that Count returns the
// zero value beside an error.
func (c meterCount) ZeroOnError() conform.ID[Meter] {
	return c.Check("zero-on-error")
}

// AllocationBudget returns the identity of the check
// Meter.Count/allocation-budget.
func (c meterCount) AllocationBudget() conform.ID[Meter] {
	return c.Check("allocation-budget")
}

// ByteBudget returns the identity of the check Meter.Count/byte-budget.
func (c meterCount) ByteBudget() conform.ID[Meter] {
	return c.Check("byte-budget")
}

// LatencyBudget returns the identity of the check
// Meter.Count/latency-budget.
func (c meterCount) LatencyBudget() conform.ID[Meter] {
	return c.Check("latency-budget")
}

// MeanLatencyBudget returns the identity of the check
// Meter.Count/mean-latency-budget.
func (c meterCount) MeanLatencyBudget() conform.ID[Meter] {
	return c.Check("mean-latency-budget")
}

// meterRecord is the name of Meter.Record, with the identities of its
// checks.
type meterRecord struct {
	conform.Callable[Meter]
}

// Smoke returns the identity of the check that Record returns without a
// panic for derived inputs.
func (c meterRecord) Smoke() conform.ID[Meter] {
	return c.Check("smoke")
}

// Cancel returns the identity of the check that Record returns an error
// for a cancelled context.
func (c meterRecord) Cancel() conform.ID[Meter] {
	return c.Check("cancel")
}

// LatencyBudget returns the identity of the check
// Meter.Record/latency-budget.
func (c meterRecord) LatencyBudget() conform.ID[Meter] {
	return c.Check("latency-budget")
}

// meterMerge is the name of Meter.Merge, with the identities of its
// checks.
type meterMerge struct {
	conform.Callable[Meter]
}

// Smoke returns the identity of the check that Merge returns without a
// panic for derived inputs.
func (c meterMerge) Smoke() conform.ID[Meter] {
	return c.Check("smoke")
}

// NilSafe returns the identity of the check that Merge returns without a
// panic for a nil argument.
func (c meterMerge) NilSafe() conform.ID[Meter] {
	return c.Check("nil-safe")
}

// meterTotal is the name of Meter.Total, with the identities of its
// checks.
type meterTotal struct {
	conform.Callable[Meter]
}

// Smoke returns the identity of the check that Total returns without a
// panic for derived inputs.
func (c meterTotal) Smoke() conform.ID[Meter] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that Total returns the
// same value for the same inputs.
func (c meterTotal) Deterministic() conform.ID[Meter] {
	return c.Check("deterministic")
}

// MeterSuite returns the checks of Meter, and the calls and the budgets of
// its callables.
func MeterSuite() conform.Suite[Meter] {
	return conform.Suite[Meter]{
		Declaration: "Meter",
		Checks: []conform.Check[Meter]{
			{
				ID:          MeterChecks.Count.Smoke(),
				Statement:   "Count returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *MeterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx, Unit(1))
						return err
					})
				},
				Defect: meterCountPanics,
			},
			{
				ID:          MeterChecks.Count.Cancel(),
				Statement:   "Count returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *MeterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						_, err := h.Subject().Count(ctx, Unit(1))
						return err
					})
				},
				Defect: meterCountIgnoresCancel,
			},
			{
				ID:          MeterChecks.Count.ZeroOnError(),
				Statement:   "Count returns the zero value beside an error",
				Claim:       spec.ZeroOnError,
				Probe:       spec.ZeroBesideError,
				Reason:      spec.ValueBesideError,
				Law:         spec.NoValueBesideAnError,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *MeterHandle) {
					probe.ZeroBesideError(h, func(ctx context.Context) (int, error) {
						return h.Subject().Count(ctx, Unit(1))
					}, int(0))
				},
				Defect: meterCountValueBeside,
			},
			{
				ID:          MeterChecks.Record.Smoke(),
				Statement:   "Record returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *MeterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						return h.Subject().Record(ctx, Sample{ID: "ID-a"})
					})
				},
				Defect: meterRecordPanics,
			},
			{
				ID:          MeterChecks.Record.Cancel(),
				Statement:   "Record returns an error for a cancelled context",
				Claim:       spec.Cancel,
				Probe:       spec.UnderCancelled,
				Reason:      spec.IgnoredCancel,
				Law:         spec.CancelledContextIsReported,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *MeterHandle) {
					probe.UnderCancelled(h, func(ctx context.Context) error {
						return h.Subject().Record(ctx, Sample{ID: "ID-a"})
					})
				},
				Defect: meterRecordIgnoresCancel,
			},
			{
				ID:          MeterChecks.Merge.Smoke(),
				Statement:   "Merge returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *MeterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Merge(&Sample{ID: "ID-a"})
						return nil
					})
				},
				Defect: meterMergePanics,
			},
			{
				ID:          MeterChecks.Merge.NilSafe(),
				Statement:   "Merge returns without a panic for a nil argument",
				Claim:       spec.NilSafe,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NilArgumentIsAccepted,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *MeterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Merge(&Sample{ID: "ID-a"})
						return nil
					})
				},
				Defect: meterMergePanics,
			},
			{
				ID:          MeterChecks.Total.Smoke(),
				Statement:   "Total returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *MeterHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Total()
						return nil
					})
				},
				Defect: meterTotalPanics,
			},
			{
				ID:          MeterChecks.Total.Deterministic(),
				Statement:   "Total returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[Meter]{MeterChecks.Total.Smoke()},
				Apply: func(h *MeterHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) (int, error) {
						r0 := h.Subject().Total()
						return r0, nil
					})
				},
				Defect: meterTotalDrifts,
			},
		},
		Calls: []conform.Call[Meter]{
			{
				Calls:    meterCountCalls,
				Over:     newMeterOver,
				Callable: MeterChecks.Count.Callable,
				Budget: conform.Budget{
					Allocs:  new(uint64(1)),
					Bytes:   new(uint64(96)),
					Latency: 2 * time.Microsecond,
					Mean:    time.Microsecond,
				},
				Concurrent: true,
			},
			{
				Calls:    meterRecordCalls,
				Over:     newMeterOver,
				Callable: MeterChecks.Record.Callable,
				Budget: conform.Budget{
					Latency: time.Millisecond,
					Warmup:  new(500),
				},
			},
			{
				Calls:    meterMergeCalls,
				Over:     newMeterOver,
				Callable: MeterChecks.Merge.Callable,
			},
			{
				Calls:    meterTotalCalls,
				Over:     newMeterOver,
				Callable: MeterChecks.Total.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.SampledFrom(Grams, Litres)),
		},
	}
}

// RunMeterSuite applies the checks of Meter, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunMeterSuite[T Meter](t *testing.T, newSubject func() T, opts ...MeterOption) {
	t.Helper()
	RunMeterSuiteFromHandle(t, func(*MeterHandle) T { return newSubject() }, opts...)
}

// RunMeterSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunMeterSuiteFromHandle[T Meter](
	t *testing.T,
	newSubject func(h *MeterHandle) T,
	opts ...MeterOption,
) {
	t.Helper()
	conform.Run(t, MeterSuite(), meterSubjects(newSubject), opts...)
}

// FuzzMeterSuite registers on f the fuzz target of the property checks of
// callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzMeterSuite[T Meter](
	f *testing.F,
	callable conform.Callable[Meter],
	newSubject func(h *MeterHandle) T,
	opts ...MeterOption,
) {
	f.Helper()
	conform.Fuzz(f, MeterSuite(), callable, meterSubjects(newSubject), opts...)
}

// RunMeterBudgets applies the allocation ceilings of Meter and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to a fresh subject of newSubject for each case.
func RunMeterBudgets[T Meter](
	t *testing.T,
	newSubject func(h *MeterHandle) T,
	opts ...MeterOption,
) {
	t.Helper()
	conform.RunBudgets(t, MeterSuite(), meterSubjects(newSubject), opts...)
}

// BenchMeterSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchMeterSuite[T Meter](
	b conform.B,
	callable conform.Callable[Meter],
	newSubject func(h *MeterHandle) T,
	opts ...MeterOption,
) {
	b.Helper()
	conform.Bench(b, MeterSuite(), callable, meterSubjects(newSubject), opts...)
}

// BenchMeterSuiteParallel measures callable on a subject of newSubject on
// GOMAXPROCS goroutines at once, under the ceilings of its budget.
func BenchMeterSuiteParallel[T Meter](
	b *testing.B,
	callable conform.Callable[Meter],
	newSubject func(h *MeterHandle) T,
	opts ...MeterOption,
) {
	b.Helper()
	conform.BenchParallel(b, MeterSuite(), callable, meterSubjects(newSubject), opts...)
}

// meterSubjects returns the factory of the subjects of Meter that calls
// newSubject with the handle of each check.
func meterSubjects[T Meter](newSubject func(h *MeterHandle) T) func(h *MeterHandle) Meter {
	return func(h *MeterHandle) Meter { return newSubject(h) }
}

// meterCountCalls generates the calls of Meter.Count with the inputs that
// each case draws.
var meterCountCalls = prop.Composite(func(c *prop.Case) func(context.Context, Meter) error {
	unit := c.Draw(prop.Of[Unit](), "unit")
	return func(ctx context.Context, subject Meter) error {
		_, err := subject.Count(ctx, unit)
		return err
	}
})

// meterRecordCalls generates the calls of Meter.Record with the inputs
// that each case draws.
var meterRecordCalls = prop.Composite(func(c *prop.Case) func(context.Context, Meter) error {
	samples := c.Draw(prop.List(prop.Of[Sample]()), "samples")
	return func(ctx context.Context, subject Meter) error {
		return subject.Record(ctx, samples...)
	}
})

// meterMergeCalls generates the calls of Meter.Merge with the inputs that
// each case draws.
var meterMergeCalls = prop.Composite(func(c *prop.Case) func(context.Context, Meter) error {
	s := c.Draw(prop.Of[*Sample](), "s")
	return func(_ context.Context, subject Meter) error {
		subject.Merge(s)
		return nil
	}
})

// meterTotalCalls generates the calls of Meter.Total with the inputs that
// each case draws.
var meterTotalCalls = prop.Composite(func(*prop.Case) func(context.Context, Meter) error {
	return func(_ context.Context, subject Meter) error {
		_ = subject.Total()
		return nil
	}
})

// meterOverEmbedded is Meter under the name of the field that meterOver
// embeds.
type meterOverEmbedded = Meter

// meterOver is a subject of Meter whose callables call spend and return
// their zero results. The planted defect of a ceiling passes a spend that
// exceeds the ceiling.
type meterOver struct {
	meterOverEmbedded
	spend func()
}

// Count calls spend and returns its zero results.
func (o meterOver) Count(context.Context, Unit) (r0 int, r1 error) {
	o.spend()
	return r0, r1
}

// Record calls spend and returns its zero results.
func (o meterOver) Record(context.Context, ...Sample) (r0 error) {
	o.spend()
	return r0
}

// Merge calls spend and returns its zero results.
func (o meterOver) Merge(*Sample) {
	o.spend()
}

// Total calls spend and returns its zero results.
func (o meterOver) Total() (r0 int) {
	o.spend()
	return r0
}

// newMeterOver returns the planted defect of a ceiling of Meter, whose
// callables call spend.
func newMeterOver(_ *MeterHandle, spend func()) Meter {
	return meterOver{spend: spend}
}

// meterCountPanicsDefect is the planted defect of a Meter whose Count
// panics.
type meterCountPanicsDefect struct {
	meterOverEmbedded
}

// Count is the override of the planted defect.
func (*meterCountPanicsDefect) Count(p0 context.Context, p1 Unit) (r0 int, r1 error) {
	panic("the planted defect of a Meter whose Count panics")
}

// meterCountPanics returns the planted defect of a Meter whose Count
// panics over a fresh subject.
func meterCountPanics(h *MeterHandle) Meter {
	return &meterCountPanicsDefect{meterOverEmbedded: h.Fresh()}
}

// meterCountIgnoresCancelDefect is the planted defect of a Meter whose
// Count ignores the context.
type meterCountIgnoresCancelDefect struct {
	meterOverEmbedded
}

// Count is the override of the planted defect.
func (d *meterCountIgnoresCancelDefect) Count(p0 context.Context, p1 Unit) (r0 int, r1 error) {
	if p0.Err() != nil {
		return r0, r1
	}
	return d.meterOverEmbedded.Count(p0, p1)
}

// meterCountIgnoresCancel returns the planted defect of a Meter whose
// Count ignores the context over a fresh subject.
func meterCountIgnoresCancel(h *MeterHandle) Meter {
	return &meterCountIgnoresCancelDefect{meterOverEmbedded: h.Fresh()}
}

// meterCountValueBesideDefect is the planted defect of a Meter whose Count
// returns a value beside an error.
type meterCountValueBesideDefect struct {
	meterOverEmbedded
}

// Count is the override of the planted defect.
func (d *meterCountValueBesideDefect) Count(p0 context.Context, p1 Unit) (r0 int, r1 error) {
	r0, r1 = d.meterOverEmbedded.Count(p0, p1)
	if r1 != nil {
		r0 = int(7)
	}
	return r0, r1
}

// meterCountValueBeside returns the planted defect of a Meter whose Count
// returns a value beside an error over a fresh subject.
func meterCountValueBeside(h *MeterHandle) Meter {
	return &meterCountValueBesideDefect{meterOverEmbedded: h.Fresh()}
}

// meterRecordPanicsDefect is the planted defect of a Meter whose Record
// panics.
type meterRecordPanicsDefect struct {
	meterOverEmbedded
}

// Record is the override of the planted defect.
func (*meterRecordPanicsDefect) Record(p0 context.Context, p1 ...Sample) (r0 error) {
	panic("the planted defect of a Meter whose Record panics")
}

// meterRecordPanics returns the planted defect of a Meter whose Record
// panics over a fresh subject.
func meterRecordPanics(h *MeterHandle) Meter {
	return &meterRecordPanicsDefect{meterOverEmbedded: h.Fresh()}
}

// meterRecordIgnoresCancelDefect is the planted defect of a Meter whose
// Record ignores the context.
type meterRecordIgnoresCancelDefect struct {
	meterOverEmbedded
}

// Record is the override of the planted defect.
func (d *meterRecordIgnoresCancelDefect) Record(p0 context.Context, p1 ...Sample) (r0 error) {
	if p0.Err() != nil {
		return r0
	}
	return d.meterOverEmbedded.Record(p0, p1...)
}

// meterRecordIgnoresCancel returns the planted defect of a Meter whose
// Record ignores the context over a fresh subject.
func meterRecordIgnoresCancel(h *MeterHandle) Meter {
	return &meterRecordIgnoresCancelDefect{meterOverEmbedded: h.Fresh()}
}

// meterMergePanicsDefect is the planted defect of a Meter whose Merge
// panics.
type meterMergePanicsDefect struct {
	meterOverEmbedded
}

// Merge is the override of the planted defect.
func (*meterMergePanicsDefect) Merge(p0 *Sample) {
	panic("the planted defect of a Meter whose Merge panics")
}

// meterMergePanics returns the planted defect of a Meter whose Merge
// panics over a fresh subject.
func meterMergePanics(h *MeterHandle) Meter {
	return &meterMergePanicsDefect{meterOverEmbedded: h.Fresh()}
}

// meterTotalPanicsDefect is the planted defect of a Meter whose Total
// panics.
type meterTotalPanicsDefect struct {
	meterOverEmbedded
}

// Total is the override of the planted defect.
func (*meterTotalPanicsDefect) Total() (r0 int) {
	panic("the planted defect of a Meter whose Total panics")
}

// meterTotalPanics returns the planted defect of a Meter whose Total
// panics over a fresh subject.
func meterTotalPanics(h *MeterHandle) Meter {
	return &meterTotalPanicsDefect{meterOverEmbedded: h.Fresh()}
}

// meterTotalDriftsDefect is the planted defect of a Meter whose Total
// returns another value on the repeat.
type meterTotalDriftsDefect struct {
	meterOverEmbedded
	calls int
}

// Total is the override of the planted defect.
func (d *meterTotalDriftsDefect) Total() (r0 int) {
	d.calls++
	if d.calls > 1 {
		r0 = int(7)
	}
	return r0
}

// meterTotalDrifts returns the planted defect of a Meter whose Total
// returns another value on the repeat over a fresh subject.
func meterTotalDrifts(h *MeterHandle) Meter {
	return &meterTotalDriftsDefect{meterOverEmbedded: h.Fresh()}
}

// The derived checks of Weigh:
//
// Weigh/smoke: Weigh returns without a panic for derived inputs.
//
// Weigh/deterministic: Weigh returns the same value for the same inputs.

// The derived checks of Weigh that witness their laws at the shallow end:
//
// Weigh/smoke
//
// Weigh/deterministic

// The unclaimed classifications of Weigh:
//
// The mixin:concurrent of Weigh has the stance needs. The claims of the
// mixin need concurrent callers.

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Weigh writes.
type (
	WeighHandle = conform.Handle[WeighFunc]
	WeighCheck  = conform.Check[WeighFunc]
	WeighOption = conform.Option[WeighFunc]
)

// WeighFunc is the type of the function Weigh. Each subject of the suite
// of Weigh has it.
type WeighFunc func(s Sample, density float64) float64

// WeighChecks contains the name of each callable of Weigh, with the
// identities of its checks.
var WeighChecks = weighCallable{"Weigh"}

// weighCallable is the name of Weigh, with the identities of its checks.
type weighCallable struct {
	conform.Callable[WeighFunc]
}

// Smoke returns the identity of the check that Weigh returns without a
// panic for derived inputs.
func (c weighCallable) Smoke() conform.ID[WeighFunc] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that Weigh returns the
// same value for the same inputs.
func (c weighCallable) Deterministic() conform.ID[WeighFunc] {
	return c.Check("deterministic")
}

// AllocationBudget returns the identity of the check
// Weigh/allocation-budget.
func (c weighCallable) AllocationBudget() conform.ID[WeighFunc] {
	return c.Check("allocation-budget")
}

// LatencyBudget returns the identity of the check Weigh/latency-budget.
func (c weighCallable) LatencyBudget() conform.ID[WeighFunc] {
	return c.Check("latency-budget")
}

// WeighSuite returns the checks of Weigh, and the calls and the budgets of
// its callables.
func WeighSuite() conform.Suite[WeighFunc] {
	return conform.Suite[WeighFunc]{
		Declaration: "Weigh",
		Checks: []conform.Check[WeighFunc]{
			{
				ID:          WeighChecks.Smoke(),
				Statement:   "Weigh returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *WeighHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject()(Sample{ID: "ID-a"}, float64(1.5))
						return nil
					})
				},
				Defect: weighPanics,
			},
			{
				ID:          WeighChecks.Deterministic(),
				Statement:   "Weigh returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[WeighFunc]{WeighChecks.Smoke()},
				Apply: func(h *WeighHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) (float64, error) {
						r0 := h.Subject()(Sample{ID: "ID-a"}, float64(1.5))
						return r0, nil
					})
				},
				Defect: weighDrifts,
			},
		},
		Calls: []conform.Call[WeighFunc]{
			{
				Calls:    weighCalls,
				Over:     newWeighOver,
				Callable: WeighChecks.Callable,
				Budget: conform.Budget{
					Allocs:  new(uint64(0)),
					Latency: 100 * time.Nanosecond,
				},
				Concurrent: true,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.SampledFrom(Grams, Litres)),
		},
	}
}

// RunWeighSuite applies the checks of Weigh, and the checks that opts add,
// to subject.
func RunWeighSuite(t *testing.T, subject WeighFunc, opts ...WeighOption) {
	t.Helper()
	RunWeighSuiteFromHandle(t, weighSubjects(subject), opts...)
}

// RunWeighSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunWeighSuiteFromHandle(
	t *testing.T,
	newSubject func(h *WeighHandle) WeighFunc,
	opts ...WeighOption,
) {
	t.Helper()
	conform.Run(t, WeighSuite(), newSubject, opts...)
}

// FuzzWeighSuite registers on f the fuzz target of the property checks of
// Weigh on subject.
func FuzzWeighSuite(f *testing.F, subject WeighFunc, opts ...WeighOption) {
	f.Helper()
	conform.Fuzz(f, WeighSuite(), WeighChecks.Callable, weighSubjects(subject), opts...)
}

// RunWeighBudgets applies the allocation ceiling of Weigh and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to subject.
func RunWeighBudgets(t *testing.T, subject WeighFunc, opts ...WeighOption) {
	t.Helper()
	conform.RunBudgets(t, WeighSuite(), weighSubjects(subject), opts...)
}

// BenchWeighSuite measures subject under the ceilings of the budget of
// Weigh.
func BenchWeighSuite(b conform.B, subject WeighFunc, opts ...WeighOption) {
	b.Helper()
	conform.Bench(b, WeighSuite(), WeighChecks.Callable, weighSubjects(subject), opts...)
}

// BenchWeighSuiteParallel measures subject on GOMAXPROCS goroutines at
// once, under the ceilings of the budget of Weigh.
func BenchWeighSuiteParallel(b *testing.B, subject WeighFunc, opts ...WeighOption) {
	b.Helper()
	conform.BenchParallel(b, WeighSuite(), WeighChecks.Callable, weighSubjects(subject), opts...)
}

// weighSubjects returns the factory of the subjects of Weigh that returns
// subject for each handle.
func weighSubjects(subject WeighFunc) func(*WeighHandle) WeighFunc {
	return func(*WeighHandle) WeighFunc { return subject }
}

// weighCalls generates the calls of Weigh with the inputs that each case
// draws.
var weighCalls = prop.Composite(func(c *prop.Case) func(context.Context, WeighFunc) error {
	s := c.Draw(prop.Of[Sample](), "s")
	density := c.Draw(prop.Of[float64](), "density")
	return func(_ context.Context, subject WeighFunc) error {
		_ = subject(s, density)
		return nil
	}
})

// newWeighOver returns the planted defect of a ceiling of Weigh, which
// calls spend and returns its zero results.
func newWeighOver(_ *WeighHandle, spend func()) WeighFunc {
	return func(Sample, float64) (r0 float64) {
		spend()
		return r0
	}
}

// weighPanics returns the planted defect of a Weigh whose Weigh panics
// over the fresh function.
func weighPanics(h *WeighHandle) WeighFunc {
	return func(p0 Sample, p1 float64) (r0 float64) {
		panic("the planted defect of a Weigh whose Weigh panics")
	}
}

// weighDrifts returns the planted defect of a Weigh whose Weigh returns
// another value on the repeat over the fresh function.
func weighDrifts(h *WeighHandle) WeighFunc {
	calls := 0
	return func(p0 Sample, p1 float64) (r0 float64) {
		calls++
		if calls > 1 {
			r0 = float64(2.5)
		}
		return r0
	}
}
