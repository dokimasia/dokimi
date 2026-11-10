// Package defectlesssim contains a simulation whose invariant and convergence
// have no planted defect.
package defectlesssim

import (
	"go.dokimi.dev/sim"

	"example.com/shop/defectless"
)

// system spawns no process.
//
//dokimi:system
func system(*sim.World[defectless.Event, defectless.Effect]) {}

// The world has no process.
//
//dokimi:invariant
func empty(*sim.Check, *sim.World[defectless.Event, defectless.Effect]) {}

// The world has no process once the faults end.
//
//dokimi:converges
func emptied(*sim.Check, *sim.World[defectless.Event, defectless.Effect]) {}
