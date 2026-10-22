package slug

import (
	"testing"

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

// SlugSuite returns the checks of Slug.
func SlugSuite() conform.Suite[SlugFunc] {
	return conform.Suite[SlugFunc]{Declaration: "Slug", Checks: []conform.Check[SlugFunc]{
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
	}}
}

// RunSlugSuite applies the checks of Slug, and the checks that opts add,
// to subject.
func RunSlugSuite(t *testing.T, subject SlugFunc, opts ...SlugOption) {
	t.Helper()
	RunSlugSuiteFromHandle(t, func(*SlugHandle) SlugFunc { return subject }, opts...)
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
	conform.Fuzz(
		f,
		SlugSuite(),
		SlugChecks.Callable,
		func(*SlugHandle) SlugFunc { return subject },
		opts...)
}
