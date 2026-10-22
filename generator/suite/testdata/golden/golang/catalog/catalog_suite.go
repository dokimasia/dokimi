package catalog

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
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

// AllocationBudget returns the identity of the check
// Catalog.Find/allocation-budget.
func (c catalogFind) AllocationBudget() conform.ID[CatalogMethods] {
	return c.Check("allocation-budget")
}

// catalogAdd is the name of Catalog.Add, with the identities of its
// checks.
type catalogAdd struct {
	conform.Callable[CatalogMethods]
}

// CatalogSuite returns the checks of Catalog, and the calls and the
// budgets of its callables.
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
		Calls: []conform.Call[CatalogMethods]{
			{
				Calls:    catalogFindCalls,
				Over:     newCatalogOver,
				Callable: CatalogChecks.Find.Callable,
				Budget: conform.Budget{
					Allocs: new(uint64(0)),
				},
			},
			{
				Calls:    catalogAddCalls,
				Over:     newCatalogOver,
				Callable: CatalogChecks.Add.Callable,
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

// RunCatalogBudgets applies the allocation ceilings of Catalog, and the
// ceilings that opts add, to a fresh subject of newSubject for each case.
func RunCatalogBudgets[T CatalogMethods](
	t *testing.T,
	newSubject func(h *CatalogHandle) T,
	opts ...CatalogOption,
) {
	t.Helper()
	conform.RunBudgets(
		t,
		CatalogSuite(),
		func(h *CatalogHandle) CatalogMethods { return newSubject(h) },
		opts...)
}

// BenchCatalogSuite measures callable on a subject of newSubject, under
// the ceilings of its budget.
func BenchCatalogSuite[T CatalogMethods](
	b conform.B,
	callable conform.Callable[CatalogMethods],
	newSubject func(h *CatalogHandle) T,
	opts ...CatalogOption,
) {
	b.Helper()
	conform.Bench(
		b,
		CatalogSuite(),
		callable,
		func(h *CatalogHandle) CatalogMethods { return newSubject(h) },
		opts...)
}

// BenchCatalogSuiteParallel measures callable on a subject of newSubject
// on GOMAXPROCS goroutines at once, under the ceilings of its budget.
func BenchCatalogSuiteParallel[T CatalogMethods](
	b *testing.B,
	callable conform.Callable[CatalogMethods],
	newSubject func(h *CatalogHandle) T,
	opts ...CatalogOption,
) {
	b.Helper()
	conform.BenchParallel(
		b,
		CatalogSuite(),
		callable,
		func(h *CatalogHandle) CatalogMethods { return newSubject(h) },
		opts...)
}

// catalogFindCalls generates the calls of Catalog.Find with the inputs
// that each case draws.
var catalogFindCalls = prop.Composite(
	func(c *prop.Case) func(context.Context, CatalogMethods) error {
		id := c.Draw(prop.Of[string](), "id")
		return func(_ context.Context, subject CatalogMethods) error {
			_, _ = subject.Find(id)
			return nil
		}
	},
)

// catalogAddCalls generates the calls of Catalog.Add with the inputs that
// each case draws.
var catalogAddCalls = prop.Composite(
	func(c *prop.Case) func(context.Context, CatalogMethods) error {
		id := c.Draw(prop.Of[string](), "id")
		name := c.Draw(prop.Of[string](), "name")
		return func(_ context.Context, subject CatalogMethods) error {
			subject.Add(id, name)
			return nil
		}
	},
)

// catalogOver is a subject of Catalog whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type catalogOver struct {
	CatalogMethods
	spend func()
}

// Find calls spend and returns its zero results.
func (o catalogOver) Find(string) (r0 string, r1 bool) {
	o.spend()
	return r0, r1
}

// Add calls spend and returns its zero results.
func (o catalogOver) Add(string, string) {
	o.spend()
}

// newCatalogOver returns the planted defect of a ceiling of Catalog, whose
// callables call spend.
func newCatalogOver(_ *CatalogHandle, spend func()) CatalogMethods {
	return catalogOver{spend: spend}
}
