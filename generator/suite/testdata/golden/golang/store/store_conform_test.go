package store

import (
	"testing"

	"go.dokimi.dev/conform"
)

// TestStoreSuite proves that each check of the suite of Store fails on its
// planted defect.
func TestStoreSuite(t *testing.T) {
	t.Parallel()

	conform.Falsify(t, StoreSuite())
}
