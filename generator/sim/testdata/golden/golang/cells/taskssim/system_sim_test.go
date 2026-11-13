package taskssim_test

import (
	"testing"

	"example.com/shop/cells/taskssim"
	"go.dokimi.dev/sim"
)

// TestSimulation runs the simulation of the system scheduled over random
// runs and the crash points of each scenario. In short mode it runs 20
// random runs and each scenario once.
func TestSimulation(t *testing.T) {
	if testing.Short() {
		sim.TestShort(t, taskssim.ScheduledSimulation())
		return
	}
	sim.Test(t, taskssim.ScheduledSimulation())
}

// TestSimulationFalsify proves each invariant of the simulation of the
// system scheduled on its planted defect. It skips in short mode.
func TestSimulationFalsify(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping the proofs of the planted defects in short mode")
	}
	sim.Falsify(t, taskssim.ScheduledSimulation())
}

// FuzzSimulation takes the choices of a run of the simulation of the
// system scheduled from the input of the fuzz target.
func FuzzSimulation(f *testing.F) {
	sim.Fuzz(f, taskssim.ScheduledSimulation())
}
