package catalog

import (
	"testing"

	"go.dokimi.dev/conform"
)

// TestCatalogSuite proves that each check of the suite of Catalog fails on
// its planted defect.
func TestCatalogSuite(t *testing.T) {
	t.Parallel()

	conform.Falsify(t, CatalogSuite())
}
