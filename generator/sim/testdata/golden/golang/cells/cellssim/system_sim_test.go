package cellssim_test

import (
	"testing"

	"example.com/shop/cells/cellssim"
	"go.dokimi.dev/sim"
)

// TestSimulation runs the simulation of the system cluster over random
// runs and the crash points of each scenario.
func TestSimulation(t *testing.T) {
	t.Parallel()

	sim.Test(t, cellssim.ClusterSimulation())
}

// TestSimulationFalsify proves each invariant of the simulation of the
// system cluster on its planted defect.
func TestSimulationFalsify(t *testing.T) {
	t.Parallel()

	sim.Falsify(t, cellssim.ClusterSimulation())
}

// FuzzSimulation takes the choices of a run of the simulation of the
// system cluster from the input of the fuzz target.
func FuzzSimulation(f *testing.F) {
	sim.Fuzz(f, cellssim.ClusterSimulation())
}
