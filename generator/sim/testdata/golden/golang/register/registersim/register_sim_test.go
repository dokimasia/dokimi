package registersim_test

import (
	"testing"

	"example.com/shop/register/registersim"
	"go.dokimi.dev/sim"
)

// TestSimulation runs the simulation of Machine over random runs and the
// crash points of each scenario.
func TestSimulation(t *testing.T) {
	sim.Test(t, registersim.MachineSimulation())
}

// TestSimulationFalsify proves each rule of a runtime and each invariant
// of the simulation of Machine on its planted defect.
func TestSimulationFalsify(t *testing.T) {
	sim.Falsify(t, registersim.MachineSimulation())
}

// FuzzSimulation takes the choices of a run of the simulation of Machine
// from the input of the fuzz target.
func FuzzSimulation(f *testing.F) {
	sim.Fuzz(f, registersim.MachineSimulation())
}
