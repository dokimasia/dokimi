package genericsim_test

import (
	"testing"

	"example.com/shop/generic/genericsim"
	"go.dokimi.dev/sim"
)

// TestSimulation runs the simulation of Replier over random runs and the
// crash points of each scenario.
func TestSimulation(t *testing.T) {
	t.Parallel()

	sim.Test(t, genericsim.ReplierSimulation())
}

// TestSimulationFalsify proves each rule of a runtime and each invariant
// of the simulation of Replier on its planted defect.
func TestSimulationFalsify(t *testing.T) {
	t.Parallel()

	sim.Falsify(t, genericsim.ReplierSimulation())
}

// FuzzSimulation takes the choices of a run of the simulation of Replier
// from the input of the fuzz target.
func FuzzSimulation(f *testing.F) {
	sim.Fuzz(f, genericsim.ReplierSimulation())
}
