package memory

import (
	"testing"
)

// TestHiddenConformsToStore applies the checks of the suite of Store to
// Hidden.
func TestHiddenConformsToStore(t *testing.T) {
	t.Parallel()

	RunStoreSuite(t, newHidden)
}
