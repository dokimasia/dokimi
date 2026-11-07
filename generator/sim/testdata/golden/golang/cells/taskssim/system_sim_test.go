package taskssim_test

import (
	"testing"

	"example.com/shop/cells/taskssim"
	"go.dokimi.dev/sim"
)

// TestSimulation runs the simulation of the system scheduled over random
// runs and the crash points of each scenario.
func TestSimulation(t *testing.T) {
	t.Parallel()

	sim.Test(t, taskssim.ScheduledSimulation())
}

// TestSimulationFalsify proves each invariant of the simulation of the
// system scheduled on its planted defect.
func TestSimulationFalsify(t *testing.T) {
	t.Parallel()

	sim.Falsify(t, taskssim.ScheduledSimulation())
}

// FuzzSimulation takes the choices of a run of the simulation of the
// system scheduled from the input of the fuzz target.
func FuzzSimulation(f *testing.F) {
	sim.Fuzz(f, taskssim.ScheduledSimulation())
}
