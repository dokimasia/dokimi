package leadertest_test

import (
	"go.dokimi.dev/conform"

	"example.com/shop/leader"
)

// falsifyOptions returns the options of the falsification of the suite of
// Primary.
func falsifyOptions() []conform.Option[leader.Primary] { return nil }
