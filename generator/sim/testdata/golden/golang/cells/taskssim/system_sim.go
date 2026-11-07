package taskssim

import (
	"go.dokimi.dev/assert/history"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/sim"
	"go.dokimi.dev/spec"
)

// ScheduledSimulation returns the simulation of the system of the
// components that scheduled spawns, with its checks.
func ScheduledSimulation() sim.Simulation[sim.None, sim.None, []history.Operation] {
	return sim.Simulation[sim.None, sim.None, []history.Operation]{
		System: scheduled,
		Spec:   history.SpecFrom(register),
		Tolerate: map[sim.Fault]int{
			sim.FaultCrash:   1,
			sim.FaultRestart: 1,
		},
		Actions: []sim.Action[sim.None, sim.None]{
			{
				Name:  "get",
				Call:  true,
				Input: func(c *prop.Case) any { return get(c) },
				Run: func(w *sim.World[sim.None, sim.None], client int, in any) {
					run, _ := in.(func(*sim.World[sim.None, sim.None], int))
					run(w, client)
				},
			},
			{
				Name:  "set",
				Call:  true,
				Input: func(c *prop.Case) any { return set(c) },
				Run: func(w *sim.World[sim.None, sim.None], client int, in any) {
					run, _ := in.(func(*sim.World[sim.None, sim.None], int))
					run(w, client)
				},
			},
		},
		Around: conform.Bubble,
		Mode:   spec.ModeTasks,
	}
}
