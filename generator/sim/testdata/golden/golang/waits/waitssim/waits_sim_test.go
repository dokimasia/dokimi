package waitssim_test

import (
	"testing"

	"example.com/shop/waits/waitssim"
	"go.dokimi.dev/sim"
)

// TestSimulation runs the simulation of Leader over random runs and the
// crash points of each scenario. In short mode it runs 20 random runs and
// each scenario once.
func TestSimulation(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		sim.TestShort(t, waitssim.LeaderSimulation())
		return
	}
	sim.Test(t, waitssim.LeaderSimulation())
}

// TestSimulationFalsify proves each rule of a runtime and each invariant
// of the simulation of Leader on its planted defect. It skips in short
// mode.
func TestSimulationFalsify(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping the proofs of the planted defects in short mode")
	}
	t.Parallel()

	sim.Falsify(t, waitssim.LeaderSimulation())
}

// FuzzSimulation takes the choices of a run of the simulation of Leader
// from the input of the fuzz target.
func FuzzSimulation(f *testing.F) {
	sim.Fuzz(f, waitssim.LeaderSimulation())
}
