package cellssim

import (
	"errors"

	"example.com/shop/cells"
	"go.dokimi.dev/assert/history"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/conform"
	"go.dokimi.dev/sim"
	"go.dokimi.dev/spec"
)

// ClusterSimulation returns the simulation of the system of the components
// that cluster spawns, with its checks.
func ClusterSimulation() sim.Simulation[sim.None, sim.None, []history.Operation] {
	return sim.Simulation[sim.None, sim.None, []history.Operation]{
		Anatomy: sim.Anatomy[sim.None, sim.None]{
			Unknown: func(err error) bool {
				return errors.Is(err, cells.ErrAmbiguous)
			},
		},
		System: func(*prop.Case) func(w *sim.World[sim.None, sim.None]) { return cluster },
		Spec:   history.SpecFrom(register),
		Tolerate: map[sim.Fault]int{
			sim.FaultCrash:   1,
			sim.FaultRestart: 1,
			sim.FaultPause:   1,
		},
		Require: map[sim.Fault]float64{
			sim.FaultCrash: 0.2,
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
		Invariants: []sim.Invariant[sim.None, sim.None]{
			{
				Statement: "The system has the spare",
				Check:     spawned,
				Defect:    func(*prop.Case) func(w *sim.World[sim.None, sim.None]) { return lonely },
			},
		},
		Around:    conform.Bubble,
		Mode:      spec.ModeNative,
		MeanSteps: 60,
		MaxSteps:  60,
		Strategy:  sim.Uniform(),
		Blocking:  true,
	}
}
