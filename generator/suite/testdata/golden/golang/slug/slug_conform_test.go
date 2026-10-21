package slug

import (
	"testing"

	"go.dokimi.dev/conform"
)

// TestSlugSuite proves that each check of the suite of Slug fails on its
// planted defect.
func TestSlugSuite(t *testing.T) {
	t.Parallel()

	conform.Falsify(t, SlugSuite())
}
