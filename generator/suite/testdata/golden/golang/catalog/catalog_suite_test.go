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
