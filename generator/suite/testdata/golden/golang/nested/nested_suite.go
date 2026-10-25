package nested

import (
	"context"
	"sync/atomic"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
)

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Courier writes.
type (
	CourierHandle = conform.Handle[Courier]
	CourierCheck  = conform.Check[Courier]
	CourierOption = conform.Option[Courier]
)

// CourierChecks contains the name of each callable of Courier, with the
// identities of its checks.
var CourierChecks = struct {
	Ship courierShip
	Send courierSend
}{
	Ship: courierShip{"Courier.Ship"},
	Send: courierSend{"Courier.Send"},
}

// courierShip is the name of Courier.Ship, with the identities of its
// checks.
type courierShip struct {
	conform.Callable[Courier]
}

// courierSend is the name of Courier.Send, with the identities of its
// checks.
type courierSend struct {
	conform.Callable[Courier]
}

// CourierSuite returns the checks of Courier, and the calls and the
// budgets of its callables.
func CourierSuite() conform.Suite[Courier] {
	return conform.Suite[Courier]{
		Declaration: "Courier",
		Calls: []conform.Call[Courier]{
			{
				Calls:    courierShipCalls,
				Over:     newCourierOver,
				Callable: CourierChecks.Ship.Callable,
			},
			{
				Calls:    courierSendCalls,
				Over:     newCourierOver,
				Callable: CourierChecks.Send.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.Composite(func(c *prop.Case) Labeler {
				labelerLabelResults0 := c.Draw(prop.List(prop.Of[Receipt]()), "labeler-label-0")
				labelerLabelResults1 := c.Draw(
					prop.List(
						prop.SampledFrom[error](nil, ErrTorn),
						prop.MinSize(len(labelerLabelResults0)),
						prop.MaxSize(len(labelerLabelResults0)),
					),
					"labeler-label-1",
				)
				var labelerLabelAt atomic.Int64
				return courierLabelerInput{
					onLabel: func(string) (r0 Receipt, r1 error) {
						i := int(labelerLabelAt.Add(1) - 1)
						if i >= len(labelerLabelResults0) {
							r1 = ErrTorn
							return r0, r1
						}
						return labelerLabelResults0[i],
							labelerLabelResults1[i]
					},
				}
			})),
			prop.Using(prop.Composite(func(c *prop.Case) func() error {
				funcResults := c.Draw(
					prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
					"func",
				)
				var funcAt atomic.Int64
				return func() (r0 error) {
					i := int(funcAt.Add(1) - 1)
					if i >= len(funcResults) {
						r0 = conform.ErrGenerated
						return r0
					}
					return funcResults[i]
				}
			})),
			prop.Using(prop.Composite(func(c *prop.Case) Ack {
				ackResults := c.Draw(
					prop.List(prop.SampledFrom[error](nil, conform.ErrGenerated)),
					"ack",
				)
				var ackAt atomic.Int64
				return func(string) (r0 error) {
					i := int(ackAt.Add(1) - 1)
					if i >= len(ackResults) {
						r0 = conform.ErrGenerated
						return r0
					}
					return ackResults[i]
				}
			})),
			prop.Using(prop.Composite(func(c *prop.Case) Scale[int] {
				scaleWeighResults := c.Draw(prop.List(prop.Of[int]()), "scale-weigh")
				var scaleWeighAt atomic.Int64
				return courierScaleInput{
					onWeigh: func() (r0 int) {
						i := int(scaleWeighAt.Add(1) - 1)
						if i >= len(scaleWeighResults) {
							return r0
						}
						return scaleWeighResults[i]
					},
				}
			})),
			prop.Using(prop.Composite(func(*prop.Case) func() {
				return func() {}
			})),
		},
	}
}

// RunCourierSuite applies the checks of Courier, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunCourierSuite[T Courier](t *testing.T, newSubject func() T, opts ...CourierOption) {
	t.Helper()
	RunCourierSuiteFromHandle(t, func(*CourierHandle) T { return newSubject() }, opts...)
}

// RunCourierSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunCourierSuiteFromHandle[T Courier](
	t *testing.T,
	newSubject func(h *CourierHandle) T,
	opts ...CourierOption,
) {
	t.Helper()
	conform.Run(t, CourierSuite(), courierSubjects(newSubject), opts...)
}

// FuzzCourierSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzCourierSuite[T Courier](
	f *testing.F,
	callable conform.Callable[Courier],
	newSubject func(h *CourierHandle) T,
	opts ...CourierOption,
) {
	f.Helper()
	conform.Fuzz(f, CourierSuite(), callable, courierSubjects(newSubject), opts...)
}

// RunCourierBudgets applies the allocation ceilings of Courier and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunCourierBudgets[T Courier](
	t *testing.T,
	newSubject func(h *CourierHandle) T,
	opts ...CourierOption,
) {
	t.Helper()
	conform.RunBudgets(t, CourierSuite(), courierSubjects(newSubject), opts...)
}

// BenchCourierSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchCourierSuite[T Courier](
	b conform.B,
	callable conform.Callable[Courier],
	newSubject func(h *CourierHandle) T,
	opts ...CourierOption,
) {
	b.Helper()
	conform.Bench(b, CourierSuite(), callable, courierSubjects(newSubject), opts...)
}

// courierSubjects returns the factory of the subjects of Courier that
// calls newSubject with the handle of each check.
func courierSubjects[T Courier](
	newSubject func(h *CourierHandle) T,
) func(h *CourierHandle) Courier {
	return func(h *CourierHandle) Courier { return newSubject(h) }
}

// courierShipCalls generates the calls of Courier.Ship with the inputs
// that each case draws.
var courierShipCalls = prop.Composite(func(c *prop.Case) func(context.Context, Courier) error {
	req := c.Draw(prop.Of[ShipRequest](), "req")
	return func(ctx context.Context, subject Courier) error {
		return subject.Ship(ctx, &req)
	}
})

// courierSendCalls generates the calls of Courier.Send with the inputs
// that each case draws.
var courierSendCalls = prop.Composite(func(c *prop.Case) func(context.Context, Courier) error {
	b := c.Draw(prop.Of[Batch](), "b")
	return func(_ context.Context, subject Courier) error {
		return subject.Send(b)
	}
})

// courierLabelerInputEmbedded is Labeler under the name of the field that
// courierLabelerInput embeds.
type courierLabelerInputEmbedded = Labeler

// courierLabelerInput is a Labeler whose methods return the results that
// the generators of the suite of Courier draw for them, in turn.
type courierLabelerInput struct {
	courierLabelerInputEmbedded
	onLabel func(a0 string) (Receipt, error)
}

// Label calls the function of Label.
func (in courierLabelerInput) Label(a0 string) (Receipt, error) {
	return in.onLabel(a0)
}

// courierScaleInputEmbedded is Scale under the name of the field that
// courierScaleInput embeds.
type courierScaleInputEmbedded = Scale[int]

// courierScaleInput is a Scale whose methods return the results that the
// generators of the suite of Courier draw for them, in turn.
type courierScaleInput struct {
	courierScaleInputEmbedded
	onWeigh func() int
}

// Weigh calls the function of Weigh.
func (in courierScaleInput) Weigh() int {
	return in.onWeigh()
}

// courierOverEmbedded is Courier under the name of the field that
// courierOver embeds.
type courierOverEmbedded = Courier

// courierOver is a subject of Courier whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type courierOver struct {
	courierOverEmbedded
	spend func()
}

// Ship calls spend and returns its zero results.
func (o courierOver) Ship(context.Context, *ShipRequest) (r0 error) {
	o.spend()
	return r0
}

// Send calls spend and returns its zero results.
func (o courierOver) Send(Batch) (r0 error) {
	o.spend()
	return r0
}

// newCourierOver returns the planted defect of a ceiling of Courier, whose
// callables call spend.
func newCourierOver(_ *CourierHandle, spend func()) Courier {
	return courierOver{spend: spend}
}

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Hire writes.
type (
	HireHandle = conform.Handle[HireFunc]
	HireCheck  = conform.Check[HireFunc]
	HireOption = conform.Option[HireFunc]
)

// HireFunc is the type of the function Hire. Each subject of the suite of
// Hire has it.
type HireFunc func(hire func() Labeler) int

// HireChecks contains the name of each callable of Hire, with the
// identities of its checks.
var HireChecks = hireCallable{"Hire"}

// hireCallable is the name of Hire, with the identities of its checks.
type hireCallable struct {
	conform.Callable[HireFunc]
}

// HireSuite returns the checks of Hire, and the calls and the budgets of
// its callables.
func HireSuite() conform.Suite[HireFunc] {
	return conform.Suite[HireFunc]{
		Declaration: "Hire",
		Calls: []conform.Call[HireFunc]{
			{
				Calls:    hireCalls,
				Over:     newHireOver,
				Callable: HireChecks.Callable,
			},
		},
		Using: []prop.Option{
			prop.Using(prop.Composite(func(c *prop.Case) Labeler {
				labelerLabelResults0 := c.Draw(prop.List(prop.Of[Receipt]()), "labeler-label-0")
				labelerLabelResults1 := c.Draw(
					prop.List(
						prop.SampledFrom[error](nil, ErrTorn),
						prop.MinSize(len(labelerLabelResults0)),
						prop.MaxSize(len(labelerLabelResults0)),
					),
					"labeler-label-1",
				)
				var labelerLabelAt atomic.Int64
				return hireLabelerInput{
					onLabel: func(string) (r0 Receipt, r1 error) {
						i := int(labelerLabelAt.Add(1) - 1)
						if i >= len(labelerLabelResults0) {
							r1 = ErrTorn
							return r0, r1
						}
						return labelerLabelResults0[i],
							labelerLabelResults1[i]
					},
				}
			})),
			prop.Using(prop.Composite(func(*prop.Case) func() {
				return func() {}
			})),
		},
	}
}

// RunHireSuite applies the checks of Hire, and the checks that opts add,
// to subject.
func RunHireSuite(t *testing.T, subject HireFunc, opts ...HireOption) {
	t.Helper()
	RunHireSuiteFromHandle(t, hireSubjects(subject), opts...)
}

// RunHireSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunHireSuiteFromHandle(
	t *testing.T,
	newSubject func(h *HireHandle) HireFunc,
	opts ...HireOption,
) {
	t.Helper()
	conform.Run(t, HireSuite(), newSubject, opts...)
}

// FuzzHireSuite registers on f the fuzz target of the property checks of
// Hire on subject.
func FuzzHireSuite(f *testing.F, subject HireFunc, opts ...HireOption) {
	f.Helper()
	conform.Fuzz(f, HireSuite(), HireChecks.Callable, hireSubjects(subject), opts...)
}

// RunHireBudgets applies the allocation ceiling of Hire and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to subject.
func RunHireBudgets(t *testing.T, subject HireFunc, opts ...HireOption) {
	t.Helper()
	conform.RunBudgets(t, HireSuite(), hireSubjects(subject), opts...)
}

// BenchHireSuite measures subject under the ceilings of the budget of
// Hire.
func BenchHireSuite(b conform.B, subject HireFunc, opts ...HireOption) {
	b.Helper()
	conform.Bench(b, HireSuite(), HireChecks.Callable, hireSubjects(subject), opts...)
}

// hireSubjects returns the factory of the subjects of Hire that returns
// subject for each handle.
func hireSubjects(subject HireFunc) func(*HireHandle) HireFunc {
	return func(*HireHandle) HireFunc { return subject }
}

// hireCalls generates the calls of Hire with the inputs that each case
// draws.
var hireCalls = prop.Composite(func(c *prop.Case) func(context.Context, HireFunc) error {
	hireResults := c.Draw(prop.List(prop.Of[Labeler]()), "hire")
	var hireAt atomic.Int64
	hire := func() (r0 Labeler) {
		i := int(hireAt.Add(1) - 1)
		if i >= len(hireResults) {
			return r0
		}
		return hireResults[i]
	}
	return func(_ context.Context, subject HireFunc) error {
		_ = subject(hire)
		return nil
	}
})

// hireLabelerInputEmbedded is Labeler under the name of the field that
// hireLabelerInput embeds.
type hireLabelerInputEmbedded = Labeler

// hireLabelerInput is a Labeler whose methods return the results that the
// generators of the suite of Hire draw for them, in turn.
type hireLabelerInput struct {
	hireLabelerInputEmbedded
	onLabel func(a0 string) (Receipt, error)
}

// Label calls the function of Label.
func (in hireLabelerInput) Label(a0 string) (Receipt, error) {
	return in.onLabel(a0)
}

// newHireOver returns the planted defect of a ceiling of Hire, which calls
// spend and returns its zero results.
func newHireOver(_ *HireHandle, spend func()) HireFunc {
	return func(func() Labeler) (r0 int) {
		spend()
		return r0
	}
}
