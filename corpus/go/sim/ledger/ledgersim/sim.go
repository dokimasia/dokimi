// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package ledgersim contains the simulation of the node of the ledger: its
// system and the transactions of its clients.
package ledgersim

import (
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"go.dokimi.dev/dokimi/corpus/go/sim/ledger"
)

// world is the world of the node of the ledger.
type world = sim.World[ledger.Event, ledger.Effect]

// keys are the lists of the transactions.
var keys = []string{"x", "y"}

// system spawns the node with a disk.
//
//dokimi:system
func system(w *world) {
	disk := ledger.Disk{}
	w.Spawn("ledger", func(*sim.Clock) sim.Machine[ledger.Event, ledger.Effect] { return ledger.NewNode(disk) })
}

// persist persists a transaction at once.
func persist(
	_ *world, _ *sim.Process[ledger.Event, ledger.Effect], f ledger.Persist,
) sim.Outcome[ledger.Event, ledger.Effect] {
	return sim.Outcome[ledger.Event, ledger.Effect]{Event: ledger.Persisted(f)}
}

// txn applies a transaction of one to three operations that the case
// chooses. Each append writes a value that no other append of the run
// writes.
//
//dokimi:action call
func txn(c *prop.Case) func(w *world, client int) {
	n := c.Draw(prop.Integer(1, 3), "operations")
	appends, lists := make([]bool, n), make([]string, n)
	for i := range n {
		appends[i] = c.Draw(prop.Boolean(), "append")
		lists[i] = keys[c.Draw(prop.Integer(0, len(keys)-1), "list")]
	}
	return func(w *world, client int) {
		base := 10 * len(w.History().Events())
		ops := make([]any, n)
		for i := range n {
			ops[i] = []any{"read", lists[i], nil}
			if appends[i] {
				ops[i] = []any{"append", lists[i], base + i}
			}
		}
		w.Call(client, "ledger", "txn", ops, []any{"x", "y"}, func(reply uint64) ledger.Event {
			return ledger.Txn{Ops: ops, Reply: ledger.ReplyToken(reply)}
		})
	}
}

// await waits for the stop of the node, as a transaction without
// operations.
//
//dokimi:action call
func await(*prop.Case) func(w *world, client int) {
	return func(w *world, client int) {
		w.Call(client, "ledger", "txn", []any{}, []any{"x", "y"}, func(reply uint64) ledger.Event {
			return ledger.Txn{Ops: []any{}, Reply: ledger.ReplyToken(reply)}
		})
	}
}
