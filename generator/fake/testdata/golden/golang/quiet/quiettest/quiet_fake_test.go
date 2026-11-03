package quiettest_test

import (
	"testing"

	"example.com/shop/quiet/quiettest"
)

// TestQuietFake applies the checks of the suite of Quiet to QuietFake.
func TestQuietFake(t *testing.T) {
	t.Parallel()

	quiettest.RunQuietSuite(t, quiettest.NewQuietFake)
}
