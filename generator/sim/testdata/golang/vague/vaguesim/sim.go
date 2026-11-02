// Package vaguesim contains the simulation of the vague machine.
package vaguesim

import (
	"go.dokimi.dev/sim"

	"example.com/shop/vague"
)

// system spawns no process.
//
//dokimi:system
func system(*sim.World[vague.Event, vague.Effect]) {}
