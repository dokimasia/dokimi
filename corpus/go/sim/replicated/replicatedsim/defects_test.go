// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package replicatedsim_test

import (
	"strings"
	"testing"

	"go.dokimi.dev/assert"
	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/golden"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"go.dokimi.dev/dokimi/corpus/go/sim/replicated"
	"go.dokimi.dev/dokimi/corpus/go/sim/replicated/replicatedsim"
)

// world is the world of the nodes of the register.
type world = sim.World[replicated.Event, replicated.Effect]

// scattered spawns a primary that replies to its start once for each entry
// of a map, in the order in which Go iterates the map.
func scattered(w *world) {
	w.Spawn("primary", func(*sim.Clock) sim.Machine[replicated.Event, replicated.Effect] {
		return scatteredPrimary{}
	})
}

// scatteredPrimary replies to its start once for each entry of a map, in the
// order of the iteration of the map.
type scatteredPrimary struct{}

// Step applies ev and appends its effects to out.
func (scatteredPrimary) Step(ev replicated.Event, out []replicated.Effect) []replicated.Effect {
	if s, started := ev.(replicated.Started); started {
		entries := map[int]bool{}
		for i := range 64 {
			entries[i] = true
		}
		for i := range entries {
			out = append(out, replicated.Reply{Value: i, To: s.Reply})
		}
	}
	return out
}

// Err returns nil, because the primary never stops itself.
func (scatteredPrimary) Err() error { return nil }

// hastySimulation returns the simulation of the register with the defect of
// its invariant, a primary that acknowledges each write before the backup has
// it, and with an action that reads a key from the backup.
func hastySimulation() sim.Simulation[replicated.Event, replicated.Effect, map[string]string] {
	s := replicatedsim.NodeSimulation()
	s.System, s.Rules, s.Invariants, s.Scenarios = s.Invariants[0].Defect, nil, nil, nil
	s.Actions = append(s.Actions, sim.Action[replicated.Event, replicated.Effect]{
		Name: "get from the backup",
		Call: true,
		Input: func(c *prop.Case) any {
			return []string{"x", "y"}[c.Draw(prop.Integer(0, 1), "key")]
		},
		Run: func(w *world, client int, in any) {
			key, _ := in.(string)
			w.Call(client, "backup", "get", []any{key}, []any{key}, func(r uint64) replicated.Event {
				return replicated.Get{Key: key, Reply: replicated.ReplyToken(r)}
			})
		},
	})
	return s
}

// drain steps each process with a pending event until no process has one.
func drain(w *world) {
	for pending := w.Pending(); len(pending) > 0; pending = w.Pending() {
		w.Step(pending[0])
	}
}

// lostUnderPartition writes 1 under x on the primary while a partition cuts
// the primary from the backup, repairs the partition, and reads x from the
// backup.
func lostUnderPartition(w *world) {
	w.Partition([]string{"primary"}, []string{"backup"})
	w.Call(0, "primary", "put", []any{"x", "1"}, []any{"x"}, func(r uint64) replicated.Event {
		return replicated.Put{Key: "x", Value: "1", Reply: replicated.ReplyToken(r)}
	})
	drain(w)
	w.Repair()
	w.Call(1, "backup", "get", []any{"x"}, []any{"x"}, func(r uint64) replicated.Event {
		return replicated.Get{Key: "x", Reply: replicated.ReplyToken(r)}
	})
	drain(w)
}

// The simulation of the register fails on the defects of its machines that
// the checks of a history and of determinism find.
func TestDefects(t *testing.T) {
	t.Parallel()

	t.Run("fails the history of a primary that acknowledges a write before the backup has it", func(t *testing.T) {
		t.Parallel()

		s := hastySimulation()
		got := assert.Rejects(t, "the simulation of a hasty primary fails at the seed 1", func(tb assert.TB) {
			prop.ForAll(tb, "the simulation finds no failure", func(c *prop.Case) { sim.Run(c, s) },
				prop.Seed(1), prop.Store(""))
		})
		inner, isCase := got[0].CaseFailure()
		assert.True(t, isCase, "the failure of the simulation contains the failure of its case")
		expect.Equal(t, inner.Assertion, "linearizable",
			"a read of the backup misses a write that the primary acknowledged")
	})

	t.Run("notes the log and the history of the shrunk case of its seed", func(t *testing.T) {
		t.Parallel()

		s := hastySimulation()
		r := assert.NewRecorder()
		prop.ForAll(r, "the simulation finds no failure", func(c *prop.Case) { sim.Run(c, s) },
			prop.Seed(1), prop.Store(""))
		assert.True(t, r.Failed(), "the simulation of a hasty primary fails at the seed 1")
		golden.Match(t, "hasty.log", []byte(strings.Join(r.Notes(), "\n")), golden.ShouldUpdate())
	})

	t.Run("fails the history of a primary that loses an acknowledged write under a partition", func(t *testing.T) {
		t.Parallel()

		s := replicatedsim.NodeSimulation()
		s.System, s.Actions, s.Rules, s.Invariants = s.Invariants[0].Defect, nil, nil, nil
		s.Scenarios = []sim.Scenario[replicated.Event, replicated.Effect]{
			{Name: "lostUnderPartition", Run: lostUnderPartition},
		}
		got := assert.Rejects(t, "the scenario of a partition fails", func(tb assert.TB) { sim.Test(tb, s) })
		expect.Equal(t, got[0].Contract, "the scenario lostUnderPartition finds no failure",
			"the scenario fails without a crash")
		inner, isCase := got[0].CaseFailure()
		assert.True(t, isCase, "the failure of the scenario contains the failure of its case")
		expect.Equal(t, inner.Contract, "the history of the calls of the clients is linearizable",
			"a read of the backup after the repair misses the write that the primary acknowledged")
	})

	t.Run("fails the determinism of a machine that iterates a map into its effects", func(t *testing.T) {
		t.Parallel()

		s := replicatedsim.NodeSimulation()
		s.System = func(*prop.Case) func(w *world) { return scattered }
		s.Rules, s.Invariants, s.Converges, s.Scenarios = nil, nil, nil, nil
		got := assert.Rejects(t, "the simulation of a scattered primary fails", func(tb assert.TB) { sim.Test(tb, s) })
		expect.Equal(t, got[0].Detail["outcome"], any(prop.Flaky),
			"the replay of a passing case makes other steps than the case made")
		divergence, _ := got[0].Detail["divergence"].(*prop.Divergence)
		assert.NotNil(t, divergence, "the record of the run states the divergence")
		expect.Equal(t, divergence.What, prop.FingerprintDifference, "the checksum of a step differs")
		expect.NotNil(t, divergence.Step, "the divergence names the first step at which the two runs differ")
	})

	t.Run("runs the invariants after a crash at the first crash point of a scenario", func(t *testing.T) {
		t.Parallel()

		s := replicatedsim.NodeSimulation()
		s.Tolerate, s.Require = nil, nil
		s.Invariants = []sim.Invariant[replicated.Event, replicated.Effect]{{
			Statement: "each node keeps its first epoch",
			Check: func(h *sim.Check, w *world) {
				for _, name := range []string{"primary", "backup"} {
					expect.Equal(h, w.Process(name).Epoch(), uint64(1), h.Statement())
				}
			},
		}}
		got := assert.Rejects(t, "the crash points of the scenario fail", func(tb assert.TB) { sim.Test(tb, s) })
		expect.Equal(t, got[0].Contract, "the scenario putThenGet finds no failure with a crash before step 1",
			"the crash before the first step of the scenario fails")
		inner, isCase := got[0].CaseFailure()
		assert.True(t, isCase, "the failure of the scenario contains the failure of its case")
		expect.Equal(t, inner.Contract, "each node keeps its first epoch", "the invariant runs after the crash")
	})
}
