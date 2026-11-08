// Package taskssim contains the simulation of a system of one cell in Tasks
// mode.
package taskssim

import (
	"context"

	"go.dokimi.dev/assert/history"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"example.com/shop/cells"
)

// cellName is the name of the process of the cell.
const cellName = "cell"

// scheduled spawns the cell, and a spare cell where the case chooses one.
//
//dokimi:system tolerate=[crash:1, restart:1, tasks]
func scheduled(c *prop.Case) func(w *sim.World[sim.None, sim.None]) {
	spare := c.Draw(prop.Boolean(), "spare")
	return func(w *sim.World[sim.None, sim.None]) {
		sim.SpawnComponent(w, cellName, cells.Open)
		if spare {
			sim.SpawnComponent(w, "spare", cells.Open)
		}
	}
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
