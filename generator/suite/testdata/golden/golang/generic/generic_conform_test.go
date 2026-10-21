package generic

import (
	"testing"

	"go.dokimi.dev/conform"
)

// TestCacheSuite proves that each check of the suite of Cache fails on its
// planted defect.
func TestCacheSuite(t *testing.T) {
	t.Parallel()

	conform.Falsify(t, CacheSuite())
}
