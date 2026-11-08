package placedsim_test

import (
	"testing"

	"example.com/shop/placed/placedsim"
	"go.dokimi.dev/sim"
)

// TestSimulation runs the simulation of Machine over random runs and the
// crash points of each scenario.
func TestSimulation(t *testing.T) {
	t.Parallel()

	sim.Test(t, placedsim.MachineSimulation())
}

// TestSimulationFalsify proves each rule of a runtime and each invariant
// of the simulation of Machine on its planted defect.
func TestSimulationFalsify(t *testing.T) {
	t.Parallel()

	sim.Falsify(t, placedsim.MachineSimulation())
}

// FuzzSimulation takes the choices of a run of the simulation of Machine
// from the input of the fuzz target.
func FuzzSimulation(f *testing.F) {
	sim.Fuzz(f, placedsim.MachineSimulation())
}
