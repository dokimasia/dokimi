// Package partialsim contains a simulation whose performers do not have the
// signature of a performer.
package partialsim

import (
	"go.dokimi.dev/sim"

	"example.com/shop/partial"
)

// world and process are the world and a process of the machines of partial.
type (
	world   = sim.World[partial.Event, partial.Effect]
	process = sim.Process[partial.Event, partial.Effect]
)

// Outcome is a type of the package with the name of the outcome of an
// effect.
type Outcome[E any] struct{ Event E }

// system spawns no process.
//
//dokimi:system
func system(*sim.World[partial.Event, partial.Effect]) {}

// wrong performs Wrong on a world in place of a process.
func wrong(*world, *world, partial.Wrong) sim.Outcome[partial.Event, partial.Effect] {
	return sim.Outcome[partial.Event, partial.Effect]{}
}

// pair performs the effect Wrong in place of Pair.
func pair(*world, *process, partial.Wrong) sim.Outcome[partial.Event, partial.Effect] {
	return sim.Outcome[partial.Event, partial.Effect]{}
}

// odd performs Odd into the outcome of the package.
func odd(*world, *process, partial.Odd) Outcome[partial.Event] { return Outcome[partial.Event]{} }

// skew performs Skew into an outcome whose events are the effects.
func skew(*world, *process, partial.Skew) sim.Outcome[partial.Effect, partial.Effect] {
	return sim.Outcome[partial.Effect, partial.Effect]{}
}
