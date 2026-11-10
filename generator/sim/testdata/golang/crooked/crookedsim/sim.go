// Package crookedsim contains the functions of a simulation, each with the
// signature of another directive than its own.
package crookedsim

import (
	"time"

	"go.dokimi.dev/assert/history"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"example.com/shop/crooked"
)

// world is the world of the machines of crooked.
type world = sim.World[crooked.Event, crooked.Effect]

// pair is a generic type of the event and the effect interfaces that is no
// world.
type pair[E, F any] struct{}

// system spawns no process.
//
//dokimi:system
func system(*sim.World[crooked.Event, crooked.Effect]) {}

// again is a second system, under the spelling of the directive with the
// name of the plugin. Its directive states a param of the simulation, which
// the directive sim of a machine interface takes.
//
//dokimi:sim:system strategy=uniform
func again(*sim.World[crooked.Event, crooked.Effect]) {}

// broken is a system that takes no world.
//
//dokimi:system
func broken(int) {}

// first is a spec of the calls.
//
//dokimi:spec
func first() history.Spec[int] { return history.Spec[int]{} }

// second is a second spec of the calls.
//
//dokimi:spec
func second() history.Spec[int] { return history.Spec[int]{} }

// parameterized is a spec that takes a parameter.
//
//dokimi:spec
func parameterized(int) history.Spec[int] { return history.Spec[int]{} }

// twofold is a spec that also returns an error.
//
//dokimi:spec
func twofold() (history.Spec[int], error) { return history.Spec[int]{}, nil }

// hashed is a fingerprint.
//
//dokimi:fingerprint
func hashed(*sim.World[crooked.Event, crooked.Effect]) uint64 { return 0 }

// rehashed is a second fingerprint.
//
//dokimi:fingerprint
func rehashed(*sim.World[crooked.Event, crooked.Effect]) uint64 { return 0 }

// unseen is a fingerprint that takes no world.
//
//dokimi:fingerprint
func unseen() uint64 { return 0 }

// caseless is an action that takes no case.
//
//dokimi:action
func caseless(int) func(*sim.World[crooked.Event, crooked.Effect]) {
	return func(*sim.World[crooked.Event, crooked.Effect]) {}
}

// counted is a fingerprint that returns a signed count.
//
//dokimi:fingerprint
func counted(*sim.World[crooked.Event, crooked.Effect]) int { return 0 }

// numbered is a spec that returns no spec.
//
//dokimi:spec
func numbered() int { return 0 }

// idle is an action that runs nothing.
//
//dokimi:action
func idle(*prop.Case) func() { return func() {} }

// inverted is an action whose run returns the world in place of taking it.
//
//dokimi:action
func inverted(*prop.Case) func() *world { return func() *world { return nil } }

// doubled is an action that also returns an error.
//
//dokimi:action
func doubled(*prop.Case) (func(*world), error) { return func(*world) {}, nil }

// anonymous is an action call without the client of the call.
//
//dokimi:action call
func anonymous(*prop.Case) func(*sim.World[crooked.Event, crooked.Effect]) {
	return func(*sim.World[crooked.Event, crooked.Effect]) {}
}

// misclient is an action call whose client is a name.
//
//dokimi:action call
func misclient(*prop.Case) func(*world, string) { return func(*world, string) {} }

// refinal is a final change whose directive states final twice.
//
//dokimi:action final final
func refinal(*prop.Case) func(*world) { return func(*world) {} }

// The world is not empty.
//
//dokimi:invariant defect=system
func unchecked(*sim.World[crooked.Event, crooked.Effect]) {}

// The world has a process.
//
//dokimi:invariant defect=system
func fallible(*sim.Check, *world) error { return nil }

// The world has its processes.
//
//dokimi:invariant defect=counted
func misplanted(*sim.Check, *sim.World[crooked.Event, crooked.Effect]) {}

// The world settles.
//
//dokimi:converges within=late defect=system
func unbounded(*sim.Check, *sim.World[crooked.Event, crooked.Effect]) {}

// late returns a bound that is no duration.
func late() int { return int(time.Second) }

// The world settles in time.
//
//dokimi:converges within=lateBy defect=system
func drifting(*sim.Check, *world) {}

// lateBy returns a bound that takes a parameter.
func lateBy(int) time.Duration { return time.Second }

// The world settles soon.
//
//dokimi:converges defect=system
func uncheckedConvergence(*sim.World[crooked.Event, crooked.Effect]) {}

// The world settles at last.
//
//dokimi:converges defect=counted
func misplantedConvergence(*sim.Check, *sim.World[crooked.Event, crooked.Effect]) {}

// twice is a scenario that takes the world twice.
//
//dokimi:scenario
func twice(*sim.World[crooked.Event, crooked.Effect], *sim.World[crooked.Event, crooked.Effect]) {}

// paired is a scenario on a pair in place of a world.
//
//dokimi:scenario
func paired(*pair[crooked.Event, crooked.Effect]) {}

// foreign is a scenario on a world whose event and effect interfaces are
// swapped.
//
//dokimi:scenario
func foreign(*sim.World[crooked.Effect, crooked.Event]) {}

// aliased is a scenario whose world has a name of its own.
//
//dokimi:scenario
func aliased(*world) {}
