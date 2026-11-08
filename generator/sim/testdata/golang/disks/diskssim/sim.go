// Package diskssim contains the simulation of the disks machine.
package diskssim

import (
	"go.dokimi.dev/sim"

	"example.com/shop/disks"
)

// system spawns no process.
//
//dokimi:system
func system(*sim.World[disks.Event, disks.Effect]) {}

// flush performs Flush without a result.
func flush(
	*sim.World[disks.Event, disks.Effect], *sim.Process[disks.Event, disks.Effect], disks.Flush,
) sim.Outcome[disks.Event, disks.Effect] {
	return sim.Outcome[disks.Event, disks.Effect]{}
}
