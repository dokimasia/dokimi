// Package cellssim contains the simulation of a system of two cells in native
// mode, under the blocking faults.
package cellssim

import (
	"context"

	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/history"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"example.com/shop/cells"
)

// The names of the processes of the cells. The clients call the cell, and
// the spare runs beside it.
const (
	cellName  = "cell"
	spareName = "spare"
)

// cluster spawns the cell and the spare.
//
//dokimi:system tolerate=[crash:1, restart:1, pause:1, blocking] require=[crash:0.2] strategy=uniform mean-steps=60 max-steps=60
func cluster(w *sim.World[sim.None, sim.None]) {
	sim.SpawnComponent(w, cellName, cells.Open)
	sim.SpawnComponent(w, spareName, cells.Open)
}

// lonely spawns the cell alone, so the invariant spawned fails.
func lonely(w *sim.World[sim.None, sim.None]) {
	sim.SpawnComponent(w, cellName, cells.Open)
}

// set sets the cell to a value that the case chooses.
//
//dokimi:action call
func set(c *prop.Case) func(w *sim.World[sim.None, sim.None], client int) {
	v := c.Draw(prop.SampledFrom("a", "b"), "value")
	return func(w *sim.World[sim.None, sim.None], client int) {
		cell := sim.ComponentOf[*cells.Cell](w, cellName)
		w.Invoke(client, cellName, "Set", []any{v}, nil, func(ctx context.Context) (any, error) {
			return nil, cell.Value().Set(ctx, v)
		})
	}
}

// get reads the value of the cell.
//
//dokimi:action call
func get(*prop.Case) func(w *sim.World[sim.None, sim.None], client int) {
	return func(w *sim.World[sim.None, sim.None], client int) {
		cell := sim.ComponentOf[*cells.Cell](w, cellName)
		w.Invoke(client, cellName, "Get", nil, nil, func(ctx context.Context) (any, error) {
			return cell.Value().Get(ctx)
		})
	}
}

// register is the spec of the cell, a register of one value.
//
//dokimi:spec
func register() history.Subject {
	value := ""
	return func(operation string, args []any) any {
		if operation == "Set" {
			value, _ = args[0].(string)
			return nil
		}
		return value
	}
}

// The system has the spare. The run checks it after every step.
//
//dokimi:invariant defect=lonely
func spawned(h *sim.Check, w *sim.World[sim.None, sim.None]) {
	expect.NotNil(h, w.Process(spareName), h.Statement())
}
