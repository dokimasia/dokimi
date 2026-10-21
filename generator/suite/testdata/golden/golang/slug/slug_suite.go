package slug

import (
	"testing"

	"go.dokimi.dev/dokimi-go/conform"
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

// SlugSuite returns the checks of Slug.
func SlugSuite() conform.Suite[SlugFunc] {
	return conform.Suite[SlugFunc]{Declaration: "Slug"}
}

// RunSlug applies the checks of Slug, and the checks that opts add, to
// subject.
func RunSlug(t *testing.T, subject SlugFunc, opts ...SlugOption) {
	t.Helper()
	conform.Verify(t, SlugSuite(), func(*SlugHandle) SlugFunc { return subject }, opts...)
}

// RunSlugOn applies the same checks to a subject that newSubject builds on
// the handle of each check.
func RunSlugOn(t *testing.T, newSubject func(h *SlugHandle) SlugFunc, opts ...SlugOption) {
	t.Helper()
	conform.Verify(t, SlugSuite(), newSubject, opts...)
}

// ControlSlug applies the same checks to control. A control is a correct
// implementation that differs from the subject.
func ControlSlug(t *testing.T, control SlugFunc, opts ...SlugOption) {
	t.Helper()
	conform.Control(t, SlugSuite(), func(*SlugHandle) SlugFunc { return control }, opts...)
}

// ControlSlugOn applies the same checks to a control that newControl
// builds on the handle of each check.
func ControlSlugOn(t *testing.T, newControl func(h *SlugHandle) SlugFunc, opts ...SlugOption) {
	t.Helper()
	conform.Control(t, SlugSuite(), newControl, opts...)
}
