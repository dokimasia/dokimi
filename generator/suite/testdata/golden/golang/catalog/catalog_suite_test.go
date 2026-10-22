package catalog

import (
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/conform"
)

// catalogSubject is a value of Catalog that only TestCatalogSuite builds.
type catalogSubject struct{ CatalogMethods }

// TestCatalogSuite checks the suite of Catalog and the two forms of its
// application.
func TestCatalogSuite(t *testing.T) {
	t.Parallel()

	t.Run("CatalogSuite", func(t *testing.T) {
		t.Parallel()
		conform.Falsify(t, CatalogSuite())
	})

	t.Run("RunCatalogSuite", func(t *testing.T) {
		t.Parallel()
		var built *catalogSubject
		RunCatalogSuite(
			t,
			func() *catalogSubject {
				built = &catalogSubject{}
				return built
			},
			conform.With(CatalogCheck{
				ID:        "Catalog/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *CatalogHandle) {
					got := h.Subject()
					assert.Equal(h, got, CatalogMethods(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*CatalogHandle) CatalogMethods { return &catalogSubject{} },
			}),
			conform.Excuse(
				CatalogChecks.Find.FindsTheAddedProduct(),
				"the test applies the check of its option alone",
			),
		)
	})

	t.Run("RunCatalogSuiteFromHandle", func(t *testing.T) {
		t.Parallel()
		var built *catalogSubject
		RunCatalogSuiteFromHandle(
			t,
			func(*CatalogHandle) *catalogSubject {
				built = &catalogSubject{}
				return built
			},
			conform.With(CatalogCheck{
				ID:        "Catalog/builds-the-subject",
				Statement: "the subject of a check is the value that the factory builds",
				Apply: func(h *CatalogHandle) {
					got := h.Subject()
					assert.Equal(h, got, CatalogMethods(built), h.Statement(), assert.ByIdentity())
				},
				Defect: func(*CatalogHandle) CatalogMethods { return &catalogSubject{} },
			}),
			conform.Excuse(
				CatalogChecks.Find.FindsTheAddedProduct(),
				"the test applies the check of its option alone",
			),
		)
	})
}

// TestCatalogSuiteAllocs falsifies the allocation ceilings of the suite of
// Catalog.
func TestCatalogSuiteAllocs(t *testing.T) {
	t.Run("CatalogSuite", func(t *testing.T) {
		conform.FalsifyAllocs(t, CatalogSuite())
	})
}

// TestCatalogConformsToCatalog applies the checks of the suite of Catalog
// to Catalog.
func TestCatalogConformsToCatalog(t *testing.T) {
	t.Parallel()

	RunCatalogSuite(t, New)
}

// TestCatalogAllocs applies the allocation ceilings of the suite of
// Catalog to Catalog.
func TestCatalogAllocs(t *testing.T) {
	RunCatalogBudgets(t, func(*CatalogHandle) *Catalog { return New() })
}

// BenchmarkCatalog measures each callable of the suite of Catalog on
// Catalog.
func BenchmarkCatalog(b *testing.B) {
	b.Run("Find", func(b *testing.B) {
		BenchCatalogSuite(
			b,
			CatalogChecks.Find.Callable,
			func(*CatalogHandle) *Catalog { return New() },
		)
	})
	b.Run("Add", func(b *testing.B) {
		BenchCatalogSuite(
			b,
			CatalogChecks.Add.Callable,
			func(*CatalogHandle) *Catalog { return New() },
		)
	})
}

// FuzzCatalogFind fuzzes the property checks of Catalog.Find on Catalog.
func FuzzCatalogFind(f *testing.F) {
	FuzzCatalogSuite(f, CatalogChecks.Find.Callable, func(*CatalogHandle) *Catalog { return New() })
}

// FuzzCatalogAdd fuzzes the property checks of Catalog.Add on Catalog.
func FuzzCatalogAdd(f *testing.F) {
	FuzzCatalogSuite(f, CatalogChecks.Add.Callable, func(*CatalogHandle) *Catalog { return New() })
}
