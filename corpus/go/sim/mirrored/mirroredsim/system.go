// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package mirroredsim contains the simulation of the mirrored register: its
// system, its actions, its spec and its oracles.
package mirroredsim

import (
	"context"
	"slices"
	"time"

	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/history"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"go.dokimi.dev/dokimi/corpus/go/sim/mirrored"
)

// world is the world of the components of the register.
type world = sim.World[sim.None, sim.None]

// The names of the processes of the register.
const (
	primaryName = "primary"
	backupName  = "backup"
)

// values are the values that the clients write, beside the empty value of a
// node without one.
var values = []string{"1", "2"}

// system spawns the primary and the backup, each on a host of its own.
//
//dokimi:system tolerate=[crash:1, restart:1, partition:1, reset:1] require=[crash:0.1, partition:0.1]
func system(w *world) {
	sim.SpawnComponent(w, primaryName, mirrored.Primary(backupName))
	sim.SpawnComponent(w, backupName, mirrored.Backup)
}

// junk spawns the register with a primary whose disk starts with a value that
// no client writes, so the invariant clean fails.
func junk(w *world) {
	sim.SpawnComponent(w, primaryName, mirrored.Seeded("junk", backupName))
	sim.SpawnComponent(w, backupName, mirrored.Backup)
}

// stale spawns the register with a primary that mirrors to a host that does
// not exist and starts with the value 1, so the convergence converged fails.
func stale(w *world) {
	sim.SpawnComponent(w, primaryName, mirrored.Seeded("1", "nowhere"))
	sim.SpawnComponent(w, backupName, mirrored.Backup)
}

// set writes a value that the case chooses on the primary.
//
//dokimi:action call
func set(c *prop.Case) func(w *world, client int) {
	v := c.Draw(prop.SampledFrom(values...), "value")
	return func(w *world, client int) {
		primary := sim.ComponentOf[*mirrored.Node](w, primaryName)
		w.Invoke(client, primaryName, "Set", []any{v}, nil, func(ctx context.Context) (any, error) {
			return nil, primary.Value().Set(ctx, v)
		})
	}
}

// get reads the value of the primary.
//
//dokimi:action call
func get(*prop.Case) func(w *world, client int) {
	return func(w *world, client int) {
		primary := sim.ComponentOf[*mirrored.Node](w, primaryName)
		w.Invoke(client, primaryName, "Get", nil, nil, func(ctx context.Context) (any, error) {
			return primary.Value().Get(ctx)
		})
	}
}

// register is the spec of the calls of the clients: a register of one value.
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

// Each node that runs keeps the empty value or a value that a client writes.
// The run checks it after every step.
//
//dokimi:invariant defect=junk
func clean(h *sim.Check, w *world) {
	for _, name := range []string{primaryName, backupName} {
		node := sim.ComponentOf[*mirrored.Node](w, name)
		if !node.Running() {
			continue
		}
		got, err := node.Value().Get(w.Context())
		expect.NoError(h, err, "the node "+name+" reads its value")
		expect.True(h, got == "" || slices.Contains(values, got), h.Statement())
	}
}

// Once the faults end, the backup has the value of the primary.
//
//dokimi:converges within=settled defect=stale
func converged(h *sim.Check, w *world) {
	primary, err := sim.ComponentOf[*mirrored.Node](w, primaryName).Value().Get(w.Context())
	expect.NoError(h, err, "the primary reads its value")
	backup, err := sim.ComponentOf[*mirrored.Node](w, backupName).Value().Get(w.Context())
	expect.NoError(h, err, "the backup reads its value")
	expect.Equal(h, backup, primary, h.Statement())
}

// settled returns the bound of the convergence converged: ten periods of the
// sends of the primary, which covers the rate error of its host.
func settled() time.Duration { return 10 * mirrored.Period }
