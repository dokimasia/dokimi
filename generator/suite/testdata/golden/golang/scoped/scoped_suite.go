package scoped

import (
	context2 "context"
	testing2 "testing"
	time2 "time"

	prop2 "go.dokimi.dev/assert/prop"
	conform2 "go.dokimi.dev/conform"
	probe2 "go.dokimi.dev/conform/probe"
	spec2 "go.dokimi.dev/spec"
)

// The suite refers to the compatibility witness of its runtime.
var _ = conform2.CompatV1

// The names that a test of Hyphenate writes.
type (
	HyphenateHandle = conform2.Handle[HyphenateFunc]
	HyphenateCheck  = conform2.Check[HyphenateFunc]
	HyphenateOption = conform2.Option[HyphenateFunc]
)

// HyphenateFunc is the type of the function Hyphenate. Each subject of the
// suite of Hyphenate has it.
type HyphenateFunc func(title string) string

// HyphenateChecks contains the name of each callable of Hyphenate, with
// the identities of its checks.
var HyphenateChecks = hyphenateCallable{"Hyphenate"}

// hyphenateCallable is the name of Hyphenate, with the identities of its
// checks.
type hyphenateCallable struct {
	conform2.Callable[HyphenateFunc]
}

// Smoke returns the identity of the check that Hyphenate returns without a
// panic for derived inputs.
func (c hyphenateCallable) Smoke() conform2.ID[HyphenateFunc] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that Hyphenate returns
// the same value for the same inputs.
func (c hyphenateCallable) Deterministic() conform2.ID[HyphenateFunc] {
	return c.Check("deterministic")
}

// LatencyBudget returns the identity of the check
// Hyphenate/latency-budget.
func (c hyphenateCallable) LatencyBudget() conform2.ID[HyphenateFunc] {
	return c.Check("latency-budget")
}

// MeanLatencyBudget returns the identity of the check
// Hyphenate/mean-latency-budget.
func (c hyphenateCallable) MeanLatencyBudget() conform2.ID[HyphenateFunc] {
	return c.Check("mean-latency-budget")
}

// HyphenateSuite returns the checks of Hyphenate, and the calls and the
// budgets of its callables.
func HyphenateSuite() conform2.Suite[HyphenateFunc] {
	return conform2.Suite[HyphenateFunc]{
		Declaration: "Hyphenate",
		Package:     "example.com/shop/scoped",
		Generator:   conform2.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform2.Check[HyphenateFunc]{
			{
				ID:          HyphenateChecks.Smoke(),
				Statement:   "Hyphenate returns without a panic for derived inputs",
				Claim:       spec2.Smoke,
				Machineries: []spec2.MachineryID{spec2.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *HyphenateHandle) {
					probe2.Survives(h, func(ctx context2.Context) error {
						h.Subject()(string("title-a"))
						return nil
					})
				},
				Defect: hyphenatePanics,
			},
			{
				ID:          HyphenateChecks.Smoke().Generated(),
				Statement:   "Hyphenate returns without a panic for derived inputs",
				Claim:       spec2.Smoke,
				Machineries: []spec2.MachineryID{spec2.GeneratedInputs, spec2.NothingBeyondACaller},
				Property: func(h *HyphenateHandle, c *prop2.Case) {
					drawnTitle := c.Draw(prop2.Of[string](), "title")
					probe2.Survives(h, func(ctx context2.Context) error {
						h.Subject()(drawnTitle)
						return nil
					})
				},
				Defect: hyphenatePanics,
			},
			{
				ID:          HyphenateChecks.Deterministic(),
				Statement:   "Hyphenate returns the same value for the same inputs",
				Claim:       spec2.Deterministic,
				Machineries: []spec2.MachineryID{spec2.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform2.ID[HyphenateFunc]{HyphenateChecks.Smoke()},
				Apply: func(h *HyphenateHandle) {
					probe2.RepeatAgrees(h, func(ctx context2.Context) (string, error) {
						r0 := h.Subject()(string("title-a"))
						return r0, nil
					})
				},
				Defect: hyphenateDrifts,
			},
			{
				ID:          HyphenateChecks.Deterministic().Generated(),
				Statement:   "Hyphenate returns the same value for the same inputs",
				Claim:       spec2.Deterministic,
				Machineries: []spec2.MachineryID{spec2.GeneratedInputs, spec2.NothingBeyondACaller},
				Entails:     []conform2.ID[HyphenateFunc]{HyphenateChecks.Smoke().Generated()},
				Property: func(h *HyphenateHandle, c *prop2.Case) {
					drawnTitle := c.Draw(prop2.Of[string](), "title")
					probe2.RepeatAgrees(h, func(ctx context2.Context) (string, error) {
						r0 := h.Subject()(drawnTitle)
						return r0, nil
					})
				},
				Defect: hyphenateDrifts,
			},
		},
		Calls: []conform2.Call[HyphenateFunc]{
			{
				Calls:    hyphenateCalls,
				Over:     newHyphenateOver,
				Callable: HyphenateChecks.Callable,
				Budget: conform2.Budget{
					Latency: 2 * time2.Millisecond,
					Mean:    time2.Millisecond,
				},
			},
		},
	}
}

// RunHyphenateSuite applies the checks of Hyphenate, and the checks that
// opts add, to subject.
func RunHyphenateSuite(t *testing2.T, subject HyphenateFunc, opts ...HyphenateOption) {
	t.Helper()
	RunHyphenateSuiteFromHandle(t, hyphenateSubjects(subject), opts...)
}

// RunHyphenateSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunHyphenateSuiteFromHandle(
	t *testing2.T,
	newSubject func(h *HyphenateHandle) HyphenateFunc,
	opts ...HyphenateOption,
) {
	t.Helper()
	conform2.Run(t, HyphenateSuite(), newSubject, opts...)
}

// FuzzHyphenateSuite registers on f the fuzz target of the property checks
// of Hyphenate on subject.
func FuzzHyphenateSuite(f *testing2.F, subject HyphenateFunc, opts ...HyphenateOption) {
	f.Helper()
	conform2.Fuzz(
		f,
		HyphenateSuite(),
		HyphenateChecks.Callable,
		hyphenateSubjects(subject),
		opts...)
}

// RunHyphenateBudgets applies the allocation ceiling of Hyphenate and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to subject.
func RunHyphenateBudgets(t *testing2.T, subject HyphenateFunc, opts ...HyphenateOption) {
	t.Helper()
	conform2.RunBudgets(t, HyphenateSuite(), hyphenateSubjects(subject), opts...)
}

// BenchHyphenateSuite measures subject under the ceilings of the budget of
// Hyphenate.
func BenchHyphenateSuite(b conform2.B, subject HyphenateFunc, opts ...HyphenateOption) {
	b.Helper()
	conform2.Bench(
		b,
		HyphenateSuite(),
		HyphenateChecks.Callable,
		hyphenateSubjects(subject),
		opts...)
}

// hyphenateSubjects returns the factory of the subjects of Hyphenate that
// returns subject for each handle.
func hyphenateSubjects(subject HyphenateFunc) func(*HyphenateHandle) HyphenateFunc {
	return func(*HyphenateHandle) HyphenateFunc { return subject }
}

// hyphenateCalls generates the calls of Hyphenate with the inputs that
// each case draws.
var hyphenateCalls = prop2.Composite(
	func(c *prop2.Case) func(context2.Context, HyphenateFunc) error {
		title := c.Draw(prop2.Of[string](), "title")
		return func(_ context2.Context, subject HyphenateFunc) error {
			_ = subject(title)
			return nil
		}
	},
)

// newHyphenateOver returns the planted defect of a ceiling of Hyphenate,
// which calls spend and returns its zero results.
func newHyphenateOver(_ *HyphenateHandle, spend func()) HyphenateFunc {
	return func(string) (r0 string) {
		spend()
		return r0
	}
}

// hyphenatePanics returns the planted defect of a Hyphenate whose
// Hyphenate panics over the fresh function.
func hyphenatePanics(h *HyphenateHandle) HyphenateFunc {
	return func(p0 string) (r0 string) {
		panic("the planted defect of a Hyphenate whose Hyphenate panics")
	}
}

// hyphenateDrifts returns the planted defect of a Hyphenate whose
// Hyphenate returns another value on the repeat over the fresh function.
func hyphenateDrifts(h *HyphenateHandle) HyphenateFunc {
	calls := 0
	return func(p0 string) (r0 string) {
		calls++
		if calls > 1 {
			r0 = string("sample-b")
		}
		return r0
	}
}
