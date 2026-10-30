package catalog

import (
	"context"
	"testing"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/conform/probe"
	"go.dokimi.dev/spec"
)

// The derived checks of Catalog:
//
// Catalog.Find/smoke: Find returns without a panic for derived inputs.
//
// Catalog.Find/zero-on-miss: Find returns the zero value beside a false
// presence flag.
//
// Catalog.Add/smoke: Add returns without a panic for derived inputs.

// The derived checks of Catalog that witness their laws at the shallow
// end:
//
// Catalog.Find/smoke
//
// Catalog.Add/smoke

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

// Smoke returns the identity of the check that Find returns without a
// panic for derived inputs.
func (c catalogFind) Smoke() conform.ID[CatalogMethods] {
	return c.Check("smoke")
}

// ZeroOnMiss returns the identity of the check that Find returns the zero
// value beside a false presence flag.
func (c catalogFind) ZeroOnMiss() conform.ID[CatalogMethods] {
	return c.Check("zero-on-miss")
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

// Smoke returns the identity of the check that Add returns without a panic
// for derived inputs.
func (c catalogAdd) Smoke() conform.ID[CatalogMethods] {
	return c.Check("smoke")
}

// CatalogSuite returns the checks of Catalog, and the calls and the
// budgets of its callables.
func CatalogSuite() conform.Suite[CatalogMethods] {
	return conform.Suite[CatalogMethods]{
		Declaration: "Catalog",
		Package:     "example.com/shop/catalog",
		Generator:   conform.Generator{Name: "dokimi", Version: "0.1.0"},
		Checks: []conform.Check[CatalogMethods]{
			{
				ID:          CatalogChecks.Find.Smoke(),
				Statement:   "Find returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CatalogHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Find(string("id-a"))
						return nil
					})
				},
				Defect: catalogFindPanics,
			},
			{
				ID:          CatalogChecks.Find.ZeroOnMiss(),
				Statement:   "Find returns the zero value beside a false presence flag",
				Claim:       spec.ZeroOnMiss,
				Probe:       spec.ZeroBesideFalse,
				Reason:      spec.ValueBesideFalse,
				Law:         spec.NoValueBesideAFalseFlag,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Apply: func(h *CatalogHandle) {
					probe.ZeroBesideFalse(h, func(ctx context.Context) (string, bool) {
						return h.Subject().Find(string("id-a"))
					}, string(""))
				},
				Defect: catalogFindValueBesideAFalseFlag,
			},
			{
				ID:          CatalogChecks.Add.Smoke(),
				Statement:   "Add returns without a panic for derived inputs",
				Claim:       spec.Smoke,
				Probe:       spec.Survives,
				Reason:      spec.Panicked,
				Law:         spec.NoPanicOnAnyInput,
				Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
				Shallow:     true,
				Apply: func(h *CatalogHandle) {
					probe.Survives(h, func(ctx context.Context) error {
						h.Subject().Add(string("id-a"), string("name-a"))
						return nil
					})
				},
				Defect: catalogAddPanics,
			},
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
	conform.Run(t, CatalogSuite(), catalogSubjects(newSubject), opts...)
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
	conform.Fuzz(f, CatalogSuite(), callable, catalogSubjects(newSubject), opts...)
}

// RunCatalogBudgets applies the allocation ceilings of Catalog and its
// checks that count allocations, and the ceilings and the checks that opts
// add, to a fresh subject of newSubject for each case.
func RunCatalogBudgets[T CatalogMethods](
	t *testing.T,
	newSubject func(h *CatalogHandle) T,
	opts ...CatalogOption,
) {
	t.Helper()
	conform.RunBudgets(t, CatalogSuite(), catalogSubjects(newSubject), opts...)
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
	conform.Bench(b, CatalogSuite(), callable, catalogSubjects(newSubject), opts...)
}

// catalogSubjects returns the factory of the subjects of Catalog that
// calls newSubject with the handle of each check.
func catalogSubjects[T CatalogMethods](
	newSubject func(h *CatalogHandle) T,
) func(h *CatalogHandle) CatalogMethods {
	return func(h *CatalogHandle) CatalogMethods { return newSubject(h) }
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

// catalogOverEmbedded is CatalogMethods under the name of the field that
// catalogOver embeds.
type catalogOverEmbedded = CatalogMethods

// catalogOver is a subject of Catalog whose callables call spend and
// return their zero results. The planted defect of a ceiling passes a
// spend that exceeds the ceiling.
type catalogOver struct {
	catalogOverEmbedded
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

// catalogFindPanicsDefect is the planted defect of a Catalog whose Find
// panics.
type catalogFindPanicsDefect struct {
	catalogOverEmbedded
}

// Find is the override of the planted defect.
func (*catalogFindPanicsDefect) Find(p0 string) (r0 string, r1 bool) {
	panic("the planted defect of a Catalog whose Find panics")
}

// catalogFindPanics returns the planted defect of a Catalog whose Find
// panics over a fresh subject.
func catalogFindPanics(h *CatalogHandle) CatalogMethods {
	return &catalogFindPanicsDefect{catalogOverEmbedded: h.Fresh()}
}

// catalogFindValueBesideAFalseFlagDefect is the planted defect of a
// Catalog whose Find returns a value beside a false presence flag.
type catalogFindValueBesideAFalseFlagDefect struct {
	catalogOverEmbedded
}

// Find is the override of the planted defect.
func (d *catalogFindValueBesideAFalseFlagDefect) Find(p0 string) (r0 string, r1 bool) {
	r0, r1 = d.catalogOverEmbedded.Find(p0)
	if !r1 {
		r0 = string("sample-b")
	}
	return r0, r1
}

// catalogFindValueBesideAFalseFlag returns the planted defect of a Catalog
// whose Find returns a value beside a false presence flag over a fresh
// subject.
func catalogFindValueBesideAFalseFlag(h *CatalogHandle) CatalogMethods {
	return &catalogFindValueBesideAFalseFlagDefect{catalogOverEmbedded: h.Fresh()}
}

// catalogAddPanicsDefect is the planted defect of a Catalog whose Add
// panics.
type catalogAddPanicsDefect struct {
	catalogOverEmbedded
}

// Add is the override of the planted defect.
func (*catalogAddPanicsDefect) Add(p0, p1 string) {
	panic("the planted defect of a Catalog whose Add panics")
}

// catalogAddPanics returns the planted defect of a Catalog whose Add
// panics over a fresh subject.
func catalogAddPanics(h *CatalogHandle) CatalogMethods {
	return &catalogAddPanicsDefect{catalogOverEmbedded: h.Fresh()}
}
