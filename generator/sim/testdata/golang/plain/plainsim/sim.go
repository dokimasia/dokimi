// Package plainsim contains the simulation of the plain machine.
package plainsim

import (
	"go.dokimi.dev/assert/history"
	"go.dokimi.dev/sim"

	"example.com/shop/plain"
)

// system spawns no process.
//
//dokimi:system
func system(*sim.World[plain.Event, plain.Effect]) {}

// beta performs Beta without a result.
func beta(
	*sim.World[plain.Event, plain.Effect], *sim.Process[plain.Event, plain.Effect], plain.Beta,
) sim.Outcome[plain.Event] {
	return sim.Outcome[plain.Event]{}
}

// subject returns the subject whose sequential behaviour is the spec of the
// calls.
//
//dokimi:spec
func subject() history.Subject { return nil }

// The world settles once the faults end.
//
//dokimi:converges
func settles(*sim.Check, *sim.World[plain.Event, plain.Effect]) {}
