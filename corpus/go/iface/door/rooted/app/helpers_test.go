// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package app_test

import (
	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/rooted"
	"go.dokimi.dev/dokimi/corpus/go/iface/door/rooted/app"
	"go.dokimi.dev/dokimi/corpus/go/iface/door/rooted/rootedtest"
)

// tallyOptions returns the options of the application of the suite of Tally
// to the tally of the program: its reference, a second tally of the program.
func tallyOptions() []rootedtest.TallyOption {
	return []rootedtest.TallyOption{
		conform.Reference(func(*rootedtest.TallyHandle) rooted.Tally { return app.NewTally() }),
	}
}
