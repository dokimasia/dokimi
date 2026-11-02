// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package ledgersim_test

import (
	"slices"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"go.dokimi.dev/dokimi/corpus/go/sim/ledger"
	"go.dokimi.dev/dokimi/corpus/go/sim/ledger/ledgersim"
)

// world is the world of the node of the ledger.
type world = sim.World[ledger.Event, ledger.Effect]

// skewed spawns a node that reads the lists of a transaction when it takes
// the transaction, and appends to them once it has persisted it.
func skewed(w *world) {
	disk := ledger.Disk{}
	w.Spawn("ledger", func(*sim.Clock) sim.Machine[ledger.Event, ledger.Effect] { return &skewedNode{disk: disk} })
}

// skewedNode reads the lists of each transaction from a snapshot of the
// disk at the start of the transaction, and applies its appends once it has
// persisted it, so two transactions that overlap each miss the appends of
// the other.
type skewedNode struct {
	disk    ledger.Disk
	pending []skewedTxn
	tokens  ledger.Token
}

// skewedTxn is a transaction that waits for its persist: its operations with
// their reads, the token of the persist and the reply token of its call.
type skewedTxn struct {
	ops   []any
	token ledger.Token
	reply ledger.ReplyToken
}

// Step applies ev and appends its effects to out.
func (n *skewedNode) Step(ev ledger.Event, out []ledger.Effect) []ledger.Effect {
	switch ev := ev.(type) {
	case ledger.Started:
		return append(out, ledger.Reply{To: ev.Reply})
	case ledger.Txn:
		n.tokens++
		snapshot := ledger.Disk{}
		for key, list := range n.disk {
			snapshot[key] = slices.Clone(list)
		}
		n.pending = append(n.pending, skewedTxn{ops: ledger.Apply(snapshot, ev.Ops), token: n.tokens, reply: ev.Reply})
		return append(out, ledger.Persist{Token: n.tokens})
	case ledger.Persisted:
		at := slices.IndexFunc(n.pending, func(p skewedTxn) bool { return p.token == ev.Token })
		if at < 0 {
			return out
		}
		p := n.pending[at]
		n.pending = slices.Delete(n.pending, at, at+1)
		for _, op := range p.ops {
			if micro, _ := op.([]any); micro[0] == "append" {
				key, _ := micro[1].(string)
				n.disk[key] = append(n.disk[key], micro[2])
			}
		}
		return append(out, ledger.Reply{Value: p.ops, To: p.reply})
	default:
		return out
	}
}

// Err returns nil, because the node never stops itself.
func (*skewedNode) Err() error { return nil }

// The simulation of the ledger fails on a node whose transactions miss the
// appends of the transactions that overlap them.
func TestDefects(t *testing.T) {
	t.Parallel()

	t.Run("fails the serializability of a node that skews its writes", func(t *testing.T) {
		t.Parallel()

		s := ledgersim.LedgerSimulation()
		s.System, s.Rules = skewed, nil
		got := assert.Rejects(t, "the simulation of a skewed node fails at the seed 1", func(tb assert.TB) {
			prop.ForAll(tb, "the simulation finds no failure", func(c *prop.Case) { sim.Run(c, s) },
				prop.Seed(1), prop.Store(""))
		})
		inner, isCase := got[0].CaseFailure()
		assert.True(t, isCase, "the failure of the simulation contains the failure of its case")
		expect.Equal(t, inner.Assertion, "serializable",
			"two transactions that overlap each read a list without the append of the other")
	})
}
