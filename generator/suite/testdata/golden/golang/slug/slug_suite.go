package slug

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
)

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

// SlugSuite returns the checks of Slug, and the calls and the budgets of
// its callables.
func SlugSuite() conform.Suite[SlugFunc] {
	return conform.Suite[SlugFunc]{
		Declaration: "Slug",
		Checks: []conform.Check[SlugFunc]{
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

// RunSlugBudgets applies the allocation ceiling of Slug, and the ceilings
// that opts add, to subject.
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
