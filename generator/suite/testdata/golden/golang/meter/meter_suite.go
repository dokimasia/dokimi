package meter

import (
	"context"
	"testing"
	"time"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
)

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

// meterTotal is the name of Meter.Total, with the identities of its
// checks.
type meterTotal struct {
	conform.Callable[Meter]
}

// MeterSuite returns the checks of Meter, and the calls and the budgets of
// its callables.
func MeterSuite() conform.Suite[Meter] {
	return conform.Suite[Meter]{
		Declaration: "Meter",
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
	conform.Run(t, MeterSuite(), func(h *MeterHandle) Meter { return newSubject(h) }, opts...)
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
	conform.Fuzz(
		f,
		MeterSuite(),
		callable,
		func(h *MeterHandle) Meter { return newSubject(h) },
		opts...)
}

// RunMeterBudgets applies the allocation ceilings of Meter, and the
// ceilings that opts add, to a fresh subject of newSubject for each case.
func RunMeterBudgets[T Meter](
	t *testing.T,
	newSubject func(h *MeterHandle) T,
	opts ...MeterOption,
) {
	t.Helper()
	conform.RunBudgets(
		t,
		MeterSuite(),
		func(h *MeterHandle) Meter { return newSubject(h) },
		opts...)
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
	conform.Bench(
		b,
		MeterSuite(),
		callable,
		func(h *MeterHandle) Meter { return newSubject(h) },
		opts...)
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
	conform.BenchParallel(
		b,
		MeterSuite(),
		callable,
		func(h *MeterHandle) Meter { return newSubject(h) },
		opts...)
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
	RunWeighSuiteFromHandle(t, func(*WeighHandle) WeighFunc { return subject }, opts...)
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
	conform.Fuzz(
		f,
		WeighSuite(),
		WeighChecks.Callable,
		func(*WeighHandle) WeighFunc { return subject },
		opts...)
}

// RunWeighBudgets applies the allocation ceiling of Weigh, and the
// ceilings that opts add, to subject.
func RunWeighBudgets(t *testing.T, subject WeighFunc, opts ...WeighOption) {
	t.Helper()
	conform.RunBudgets(t, WeighSuite(), func(*WeighHandle) WeighFunc { return subject }, opts...)
}

// BenchWeighSuite measures subject under the ceilings of the budget of
// Weigh.
func BenchWeighSuite(b conform.B, subject WeighFunc, opts ...WeighOption) {
	b.Helper()
	conform.Bench(
		b,
		WeighSuite(),
		WeighChecks.Callable,
		func(*WeighHandle) WeighFunc { return subject },
		opts...)
}

// BenchWeighSuiteParallel measures subject on GOMAXPROCS goroutines at
// once, under the ceilings of the budget of Weigh.
func BenchWeighSuiteParallel(b *testing.B, subject WeighFunc, opts ...WeighOption) {
	b.Helper()
	conform.BenchParallel(
		b,
		WeighSuite(),
		WeighChecks.Callable,
		func(*WeighHandle) WeighFunc { return subject },
		opts...)
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
