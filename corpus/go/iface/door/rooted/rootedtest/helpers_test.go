// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package rootedtest_test

import (
	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/rooted"
	"go.dokimi.dev/dokimi/corpus/go/iface/door/rooted/app"
	"go.dokimi.dev/dokimi/corpus/go/iface/door/rooted/rootedtest"
)

// falsifyOptions returns the options of the falsification of the suite of
// Tally: the reference of the planted defects, a tally of the composition
// root.
func falsifyOptions() []rootedtest.TallyOption {
	return []rootedtest.TallyOption{
		conform.Reference(func(*rootedtest.TallyHandle) rooted.Tally { return app.NewTally() }),
	}
}
