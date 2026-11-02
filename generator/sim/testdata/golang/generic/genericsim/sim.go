// Package genericsim contains the simulation of the generic machine.
package genericsim

import (
	"go.dokimi.dev/sim"

	"example.com/shop/generic"
)

// system spawns the node a.
//
//dokimi:system
func system(w *sim.World[generic.Event, generic.Effect]) {
	w.Spawn("a", func(*sim.Clock) sim.Machine[generic.Event, generic.Effect] { return generic.Node{} })
}
