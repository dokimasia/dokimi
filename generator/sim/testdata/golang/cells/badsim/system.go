// Package badsim contains the functions of a system of components that sim
// refuses.
package badsim

import (
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"example.com/shop/cells"
)

// bounded spawns the cell, and its param tolerate has a value that is no
// fault and no mode.
//
//dokimi:system tolerate=[crash:1, threads]
func bounded(w *sim.World[sim.None, sim.None]) {
	sim.SpawnComponent(w, "cell", cells.Open)
}

// second is a second system of the package.
//
//dokimi:system
func second(*sim.World[sim.None, sim.None]) {}

// stray is an action whose run takes a world of other events and effects.
//
//dokimi:action
func stray(*prop.Case) func(w *sim.World[int, int]) { return func(*sim.World[int, int]) {} }
