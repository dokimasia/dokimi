// Package valuedsim contains a system of components with a signature that sim
// refuses, and systems whose worlds are no worlds of components.
package valuedsim

import (
	"errors"

	"go.dokimi.dev/sim"
)

// errUnknown is an error that leaves the outcome of a call unknown.
var errUnknown = errors.New("valuedsim: the outcome of the call is unknown")

// pair is a generic type of two type parameters that is no world.
type pair[E, F any] struct{}

// valued is a system of components that returns a value.
//
//dokimi:system
func valued(*sim.World[sim.None, sim.None]) int { return 0 }

// idle is a system that takes no world, which the simulation of a machine
// interface binds.
//
//dokimi:system
func idle() {}

// half is a system whose world has events of another type than sim.None.
//
//dokimi:system
func half(*sim.World[int, sim.None]) {}

// twin is a system whose world has effects of another type than sim.None.
//
//dokimi:system
func twin(*sim.World[sim.None, int]) {}

// paired is a system that takes another generic type of sim.None.
//
//dokimi:system
func paired(*pair[sim.None, sim.None]) {}

// uncertain is a system whose world has events of another type than sim.None.
// Its directive states the errors of an unknown outcome, which the directive
// reply of a machine states.
//
//dokimi:system unknown=[errUnknown]
func uncertain(*sim.World[int, sim.None]) {}
