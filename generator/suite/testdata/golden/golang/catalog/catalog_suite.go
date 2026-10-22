package catalog

import (
	"testing"

	"go.dokimi.dev/conform"
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
	Add(id, name string)
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

// FindsTheAddedProduct returns the identity of the check that Find returns
// the product that Add listed.
func (c catalogFind) FindsTheAddedProduct() conform.ID[CatalogMethods] {
	return c.Check("finds-the-added-product")
}

// catalogAdd is the name of Catalog.Add, with the identities of its
// checks.
type catalogAdd struct {
	conform.Callable[CatalogMethods]
}

// CatalogSuite returns the checks of Catalog.
func CatalogSuite() conform.Suite[CatalogMethods] {
	return conform.Suite[CatalogMethods]{
		Declaration: "Catalog",
		Checks: []conform.Check[CatalogMethods]{
			{
				ID:        CatalogChecks.Find.FindsTheAddedProduct(),
				Statement: "Find returns the product that Add listed",
				Apply:     FindsTheAddedProduct,
				Defect:    amnesiac,
			},
		},
	}
}

// RunCatalogSuite applies the checks of Catalog, and the checks that opts
// add, to a fresh subject of newSubject for each check.
func RunCatalogSuite[T CatalogMethods](t *testing.T, newSubject func() T, opts ...CatalogOption) {
	t.Helper()
	RunCatalogSuiteFromHandle(t, func(*CatalogHandle) T { return newSubject() }, opts...)
}

// RunCatalogSuiteFromHandle applies the same checks to a subject that
// newSubject builds from the handle of each check.
func RunCatalogSuiteFromHandle[T CatalogMethods](
	t *testing.T,
	newSubject func(h *CatalogHandle) T,
	opts ...CatalogOption,
) {
	t.Helper()
	conform.Run(
		t,
		CatalogSuite(),
		func(h *CatalogHandle) CatalogMethods { return newSubject(h) },
		opts...)
}

// FuzzCatalogSuite registers on f the fuzz target of the property checks
// of callable, on a subject that newSubject builds from the handle of each
// case.
func FuzzCatalogSuite[T CatalogMethods](
	f *testing.F,
	callable conform.Callable[CatalogMethods],
	newSubject func(h *CatalogHandle) T,
	opts ...CatalogOption,
) {
	f.Helper()
	conform.Fuzz(
		f,
		CatalogSuite(),
		callable,
		func(h *CatalogHandle) CatalogMethods { return newSubject(h) },
		opts...)
}
