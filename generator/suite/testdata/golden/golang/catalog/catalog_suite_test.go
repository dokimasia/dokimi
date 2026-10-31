package catalog_test

import (
	"testing"

	"example.com/shop/catalog"
	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// catalogSubjectEmbedded is catalog.CatalogMethods under the name of the
// field that catalogSubject embeds.
type catalogSubjectEmbedded = catalog.CatalogMethods

// catalogSubject is a value of Catalog that only TestCatalogSuite builds.
type catalogSubject struct{ catalogSubjectEmbedded }

// TestCatalogSuite checks the suite of Catalog and the two forms of its
// application.
func TestCatalogSuite(t *testing.T) {
	t.Parallel()

	t.Run("CatalogSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, catalog.CatalogSuite())
	})

	t.Run("RunCatalogSuite", func(t *testing.T) {
		t.Parallel()
		var built *catalogSubject
		catalog.RunCatalogSuite(
			t,
			func() *catalogSubject {
				built = &catalogSubject{}
				return built
			},
			conform.With(catalog.CatalogCheck{
				ID:        "Catalog/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *catalog.CatalogHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						catalog.CatalogMethods(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*catalog.CatalogHandle) catalog.CatalogMethods { return &catalogSubject{} },
			}),
			conform.Excuse(
				catalog.CatalogChecks.Find.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				catalog.CatalogChecks.Find.ZeroOnMiss(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				catalog.CatalogChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				catalog.CatalogChecks.Find.FindsTheAddedProduct(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunCatalogSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *catalogSubject
		catalog.RunCatalogSuiteFromHandle(
			t,
			func(*catalog.CatalogHandle) *catalogSubject {
				built = &catalogSubject{}
				return built
			},
			conform.With(catalog.CatalogCheck{
				ID:        "Catalog/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *catalog.CatalogHandle) {
					got := h.Subject()
					assert.Equal(
						h,
						got,
						catalog.CatalogMethods(built),
						h.Statement(),
						assert.ByIdentity(),
					)
				},
				Defect: func(*catalog.CatalogHandle) catalog.CatalogMethods { return &catalogSubject{} },
			}),
			conform.Excuse(
				catalog.CatalogChecks.Find.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				catalog.CatalogChecks.Find.ZeroOnMiss(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				catalog.CatalogChecks.Add.Smoke(),
				"the test applies the check of its option alone",
			),
			conform.Excuse(
				catalog.CatalogChecks.Find.FindsTheAddedProduct(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("CatalogChecks", func(t *testing.T) {
		t.Parallel()
		assert.Equal(
			t,
			string(catalog.CatalogChecks.Find.FindsTheAddedProduct()),
			"Catalog.Find/finds-the-added-product",
			"the identity of the check joins the callable and the segment of the check",
		)
		assert.Equal(
			t,
			string(catalog.CatalogChecks.Find.AllocationBudget()),
			"Catalog.Find/allocation-budget",
			"the identity of the check joins the callable and the segment of the check",
		)
	})
}

// TestCatalogSuiteAllocs falsifies the allocation ceilings and the checks
// that count allocations of the suite of Catalog.
func TestCatalogSuiteAllocs(t *testing.T) {
	t.Run("CatalogSuite", func(t *testing.T) {
		conform.FalsifyAllocs(t, catalog.CatalogSuite())
	})
}

// catalogIdle returns the planted defect of a ceiling of Catalog with a
// spend of nothing, a subject that answers each call at once.
func catalogIdle(h *catalog.CatalogHandle) catalog.CatalogMethods {
	return catalog.CatalogSuite().Calls[0].Over(h, func() {})
}

// TestCatalogSuiteBudgets applies the allocation ceilings of the suite of
// Catalog to catalogIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no
// input of their own.
func TestCatalogSuiteBudgets(t *testing.T) {
	catalog.RunCatalogBudgets(
		t,
		catalogIdle,
		conform.Ceilings(
			catalog.CatalogChecks.Add.Callable,
			conform.Budget{Allocs: new(uint64(0))},
		),
		conform.Excuse(
			catalog.CatalogChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			catalog.CatalogChecks.Find.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			catalog.CatalogChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			catalog.CatalogChecks.Find.FindsTheAddedProduct(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCatalogSuiteFind registers the fuzz target of Catalog.Find on
// catalogIdle, with every check of the suite excused.
func FuzzCatalogSuiteFind(f *testing.F) {
	catalog.FuzzCatalogSuite(
		f,
		catalog.CatalogChecks.Find.Callable,
		catalogIdle,
		conform.Excuse(
			catalog.CatalogChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			catalog.CatalogChecks.Find.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			catalog.CatalogChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			catalog.CatalogChecks.Find.FindsTheAddedProduct(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// FuzzCatalogSuiteAdd registers the fuzz target of Catalog.Add on
// catalogIdle, with every check of the suite excused.
func FuzzCatalogSuiteAdd(f *testing.F) {
	catalog.FuzzCatalogSuite(
		f,
		catalog.CatalogChecks.Add.Callable,
		catalogIdle,
		conform.Excuse(
			catalog.CatalogChecks.Find.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			catalog.CatalogChecks.Find.ZeroOnMiss(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			catalog.CatalogChecks.Add.Smoke(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
		conform.Excuse(
			catalog.CatalogChecks.Find.FindsTheAddedProduct(),
			"the planted defect of a ceiling answers each call with its zero results",
		),
	)
}

// BenchmarkCatalogSuite measures each callable of the suite of Catalog on
// catalogIdle.
func BenchmarkCatalogSuite(b *testing.B) {
	b.Run("Find", func(b *testing.B) {
		catalog.BenchCatalogSuite(b, catalog.CatalogChecks.Find.Callable, catalogIdle)
	})
	b.Run("Add", func(b *testing.B) {
		catalog.BenchCatalogSuite(b, catalog.CatalogChecks.Add.Callable, catalogIdle)
	})
}

// TestCatalogConformsToCatalog applies the checks of the suite of Catalog
// to Catalog.
func TestCatalogConformsToCatalog(t *testing.T) {
	t.Parallel()

	catalog.RunCatalogSuite(t, catalog.New)
}

// TestCatalogConformsToCatalogAllocs applies the allocation ceilings and
// the checks that count allocations of the suite of Catalog to Catalog.
func TestCatalogConformsToCatalogAllocs(t *testing.T) {
	catalog.RunCatalogBudgets(
		t,
		func(*catalog.CatalogHandle) *catalog.Catalog { return catalog.New() },
	)
}

// BenchmarkCatalogConformsToCatalog measures each callable of the suite of
// Catalog on Catalog.
func BenchmarkCatalogConformsToCatalog(b *testing.B) {
	b.Run("Find", func(b *testing.B) {
		catalog.BenchCatalogSuite(
			b,
			catalog.CatalogChecks.Find.Callable,
			func(*catalog.CatalogHandle) *catalog.Catalog { return catalog.New() },
		)
	})
	b.Run("Add", func(b *testing.B) {
		catalog.BenchCatalogSuite(
			b,
			catalog.CatalogChecks.Add.Callable,
			func(*catalog.CatalogHandle) *catalog.Catalog { return catalog.New() },
		)
	})
}

// FuzzCatalogConformsToCatalogFind fuzzes the property checks of
// Catalog.Find on Catalog.
func FuzzCatalogConformsToCatalogFind(f *testing.F) {
	catalog.FuzzCatalogSuite(
		f,
		catalog.CatalogChecks.Find.Callable,
		func(*catalog.CatalogHandle) *catalog.Catalog { return catalog.New() },
	)
}

// FuzzCatalogConformsToCatalogAdd fuzzes the property checks of
// Catalog.Add on Catalog.
func FuzzCatalogConformsToCatalogAdd(f *testing.F) {
	catalog.FuzzCatalogSuite(
		f,
		catalog.CatalogChecks.Add.Callable,
		func(*catalog.CatalogHandle) *catalog.Catalog { return catalog.New() },
	)
}
