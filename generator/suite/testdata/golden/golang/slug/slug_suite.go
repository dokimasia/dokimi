package slug

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Slug:
//
// Slug/smoke: Slug returns without a panic for derived inputs.
//
// Slug/deterministic: Slug returns the same value for the same inputs.

// The derived checks of Slug that witness their laws at the shallow end:
//
// Slug/smoke
//
// Slug/deterministic

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Slug writes.
type (
	SlugHandle = conform.Handle[SlugFunc]
	SlugCheck  = conform.Check[SlugFunc]
	SlugOption = conform.Option[SlugFunc]
)

// SlugFunc is the type of the function Slug. Each subject of the suite of
// Slug has it.
type SlugFunc func(title string, keep ...string) string

// SlugChecks contains the name of each callable of Slug, with the
// identities of its checks.
var SlugChecks = slugCallable{"Slug"}

// slugCallable is the name of Slug, with the identities of its checks.
type slugCallable struct {
	conform.Callable[SlugFunc]
}

// KeepsTheOrder returns the identity of the check that Slug keeps the
// words of the title in their order.
func (c slugCallable) KeepsTheOrder() conform.ID[SlugFunc] {
	return c.Check("keeps-the-order")
}

// KeepsTheWords returns the identity of the check that Slug keeps each
// word of the title.
func (c slugCallable) KeepsTheWords() conform.ID[SlugFunc] {
	return c.Check("keeps-the-words")
}

// Smoke returns the identity of the check that Slug returns without a
// panic for derived inputs.
func (c slugCallable) Smoke() conform.ID[SlugFunc] {
	return c.Check("smoke")
}

// Deterministic returns the identity of the check that Slug returns the
// same value for the same inputs.
func (c slugCallable) Deterministic() conform.ID[SlugFunc] {
	return c.Check("deterministic")
}

// SlugSuite returns the checks of Slug, and the calls and the budgets of
// its callables.
func SlugSuite() conform.Suite[SlugFunc] {
	return conform.Suite[SlugFunc]{
		Declaration: "Slug",
		Package:     "example.com/shop/slug",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[SlugFunc]{
			{
				ID:          SlugChecks.Smoke(),
				Statement:   "Slug returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *SlugHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject()(string("title-a"), string("keep-a"))
						return nil
					})
				},
				Defect: slugPanics,
			},
			{
				ID:          SlugChecks.Deterministic(),
				Statement:   "Slug returns the same value for the same inputs",
				Claim:       spec.Deterministic,
				Probe:       spec.RepeatAgrees,
				Reason:      spec.UnstableRepeat,
				Law:         spec.RepeatReturnsAlike,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Entails:     []conform.ID[SlugFunc]{SlugChecks.Smoke()},
				Apply: func(h *SlugHandle) {
					probe.RepeatAgrees(h, func(ctx context.Context) (string, error) {
						r0 := h.Subject()(string("title-a"), string("keep-a"))
						return r0, nil
					})
				},
				Defect: slugDrifts,
			},
			{
				ID:        SlugChecks.KeepsTheOrder(),
				Statement: "Slug keeps the words of the title in their order",
				Apply:     KeepsTheOrder,
				Defect:    dropsWords,
			},
			{
				ID:        SlugChecks.KeepsTheWords(),
				Statement: "Slug keeps each word of the title",
				Apply:     KeepsTheWords,
				Defect:    dropsWords,
			},
		},
		Calls: []conform.Call[SlugFunc]{
			{
				Calls:    slugCalls,
				Over:     newSlugOver,
				Callable: SlugChecks.Callable,
			},
		},
	}
}

// RunSlugSuite applies the checks of Slug, and the checks that opts add,
// to subject.
func RunSlugSuite(t *testing.T, subject SlugFunc, opts ...SlugOption) {
	t.Helper()
	RunSlugSuiteFromHandle(t, slugSubjects(subject), opts...)
}

// RunSlugSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunSlugSuiteFromHandle(
	t *testing.T,
	newSubject func(h *SlugHandle) SlugFunc,
	opts ...SlugOption,
) {
	t.Helper()
	conform.Run(t, SlugSuite(), newSubject, opts...)
}

// FuzzSlugSuite registers on f the fuzz target of the property checks of
// Slug on subject.
func FuzzSlugSuite(f *testing.F, subject SlugFunc, opts ...SlugOption) {
	f.Helper()
	conform.Fuzz(f, SlugSuite(), SlugChecks.Callable, slugSubjects(subject), opts...)
}

// RunSlugBudgets applies the allocation ceiling of Slug and its checks
// that count allocations, and the ceilings and the checks that opts add,
// to subject.
func RunSlugBudgets(t *testing.T, subject SlugFunc, opts ...SlugOption) {
	t.Helper()
	conform.RunBudgets(t, SlugSuite(), slugSubjects(subject), opts...)
}

// BenchSlugSuite measures subject under the ceilings of the budget of
// Slug.
func BenchSlugSuite(b conform.B, subject SlugFunc, opts ...SlugOption) {
	b.Helper()
	conform.Bench(b, SlugSuite(), SlugChecks.Callable, slugSubjects(subject), opts...)
}

// slugSubjects returns the factory of the subjects of Slug that returns
// subject for each handle.
func slugSubjects(subject SlugFunc) func(*SlugHandle) SlugFunc {
	return func(*SlugHandle) SlugFunc { return subject }
}

// slugCalls generates the calls of Slug with the inputs that each case
// draws.
var slugCalls = prop.Composite(func(c *prop.Case) func(context.Context, SlugFunc) error {
	title := c.Draw(prop.Of[string](), "title")
	keep := c.Draw(prop.List(prop.Of[string]()), "keep")
	return func(_ context.Context, subject SlugFunc) error {
		_ = subject(title, keep...)
		return nil
	}
})

// newSlugOver returns the planted defect of a ceiling of Slug, which calls
// spend and returns its zero results.
func newSlugOver(_ *SlugHandle, spend func()) SlugFunc {
	return func(string, ...string) (r0 string) {
		spend()
		return r0
	}
}

// slugPanics returns the planted defect of a Slug whose Slug panics over
// the fresh function.
func slugPanics(h *SlugHandle) SlugFunc {
	return func(p0 string, p1 ...string) (r0 string) {
		panic("the planted defect of a Slug whose Slug panics")
	}
}

// slugDrifts returns the planted defect of a Slug whose Slug returns
// another value on the repeat over the fresh function.
func slugDrifts(h *SlugHandle) SlugFunc {
	calls := 0
	return func(p0 string, p1 ...string) (r0 string) {
		calls++
		if calls > 1 {
			r0 = string("sample-b")
		}
		return r0
	}
}
