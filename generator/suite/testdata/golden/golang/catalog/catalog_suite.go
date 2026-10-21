package catalog

import (
	"testing"

	"go.dokimi.dev/dokimi-go/conform"
)

// The suite refers to the compatibility witness of its runtime.
var _ = conform.CompatV1

// The names that a test of Catalog writes.
type (
	CatalogHandle = conform.Handle[CatalogMethods]
	CatalogCheck  = conform.Check[CatalogMethods]
	CatalogOption = conform.Option[CatalogMethods]
)

// CatalogMethods is the interface of the exported methods of Catalog. Each
// subject of the suite of Catalog implements it.
type CatalogMethods interface {
	Find(id string) (string, bool)
	Add(id string, name string)
}

// CatalogChecks contains the name of each callable of Catalog, with the
// identities of its checks.
var CatalogChecks = struct {
	Find catalogFind
	Add  catalogAdd
}{
	Find: catalogFind{"Catalog.Find"},
	Add:  catalogAdd{"Catalog.Add"},
}

// catalogFind is the name of Catalog.Find, with the identities of its
// checks.
type catalogFind struct {
	conform.Callable[CatalogMethods]
}

// catalogAdd is the name of Catalog.Add, with the identities of its
// checks.
type catalogAdd struct {
	conform.Callable[CatalogMethods]
}

// CatalogSuite returns the checks of Catalog.
func CatalogSuite() conform.Suite[CatalogMethods] {
	return conform.Suite[CatalogMethods]{Declaration: "Catalog"}
}

// RunCatalog applies the checks of Catalog, and the checks that opts add,
// to a fresh subject of newSubject for each check.
func RunCatalog[T CatalogMethods](t *testing.T, newSubject func() T, opts ...CatalogOption) {
	t.Helper()
	conform.Verify(t, CatalogSuite(), func(*CatalogHandle) CatalogMethods { return newSubject() }, opts...)
}

// RunCatalogOn applies the same checks to a subject that newSubject builds
// on the handle of each check.
func RunCatalogOn[T CatalogMethods](t *testing.T, newSubject func(h *CatalogHandle) T, opts ...CatalogOption) {
	t.Helper()
	conform.Verify(t, CatalogSuite(), func(h *CatalogHandle) CatalogMethods { return newSubject(h) }, opts...)
}

// ControlCatalog applies the same checks to a fresh control of newControl
// for each check. A control is a correct implementation that differs from
// the subject.
func ControlCatalog[T CatalogMethods](t *testing.T, newControl func() T, opts ...CatalogOption) {
	t.Helper()
	conform.Control(t, CatalogSuite(), func(*CatalogHandle) CatalogMethods { return newControl() }, opts...)
}

// ControlCatalogOn applies the same checks to a control that newControl
// builds on the handle of each check.
func ControlCatalogOn[T CatalogMethods](t *testing.T, newControl func(h *CatalogHandle) T, opts ...CatalogOption) {
	t.Helper()
	conform.Control(t, CatalogSuite(), func(h *CatalogHandle) CatalogMethods { return newControl(h) }, opts...)
}
