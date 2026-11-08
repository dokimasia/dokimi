// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package waitssim contains the simulation of the leaders of waits: its
// system, its performers, its actions and its oracles.
package waitssim

import (
	"context"
	"slices"
	"strconv"

	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"go.dokimi.dev/dokimi/corpus/go/sim/waits"
)

// world is the world of the leaders of waits, process one of its processes,
// and outcome the outcome of one of their effects.
type (
	world   = sim.World[waits.Event, waits.Effect]
	process = sim.Process[waits.Event, waits.Effect]
	outcome = sim.Outcome[waits.Event, waits.Effect]
)

// leaders are the leaders of the system, and scope the scope of their terms.
var leaders = []string{"a", "b"}

const scope = "s"

// terms is the store of the terms of the scope. Its line contains the
// candidacies in the order of their campaigns, and the first of them has the
// term of the scope. A store of the defect keeps a candidacy at its resign.
type terms struct {
	line  []*candidacy
	term  uint64
	keeps bool
}

// candidacy is a candidacy of the session of a leader in the store, with its
// term, 0 before the grant.
type candidacy struct {
	store   *terms
	session string
	term    uint64
}

// Campaign enters a candidacy of the session at the end of the line.
func (t *terms) Campaign(_ context.Context, _, session string) (waits.Candidacy, error) {
	c := &candidacy{store: t, session: session}
	t.line = append(t.line, c)
	t.grant()
	return c, nil
}

// Term returns the term of the first candidacy of the line, and 0 for an
// empty line.
func (t *terms) Term(string) uint64 {
	if len(t.line) == 0 {
		return 0
	}
	return t.line[0].term
}

// grant grants the next term to the first candidacy of the line where that
// candidacy has none.
func (t *terms) grant() {
	if len(t.line) > 0 && t.line[0].term == 0 {
		t.term++
		t.line[0].term = t.term
	}
}

// Term returns the term of the candidacy.
func (c *candidacy) Term() uint64 { return c.term }

// Resign removes the candidacy from the line and grants the next term. A
// store of the defect keeps the candidacy.
func (c *candidacy) Resign() {
	t := c.store
	if t.keeps {
		return
	}
	t.line = slices.DeleteFunc(t.line, func(o *candidacy) bool { return o == c })
	t.grant()
}

// spawned returns the system of the leaders over the store t.
func spawned(t *terms) func(w *world) {
	return func(w *world) {
		for _, name := range leaders {
			w.Spawn(name, func(*sim.Clock) sim.Machine[waits.Event, waits.Effect] {
				return waits.NewLeader(scope)
			}, waits.Store(t))
		}
	}
}

// system spawns the leaders over one store.
//
//dokimi:system
func system(w *world) { spawned(&terms{})(w) }

// leaky spawns the leaders over a store that keeps the candidacy of each
// campaign whose wait ends.
func leaky(w *world) { spawned(&terms{keeps: true})(w) }

// campaign campaigns for the term of the scope of the effect f in the store,
// for the session of the leader p: its name and its epoch. Its wait delivers
// the grant of the term once the store grants it. After the cancel of the
// campaign it delivers the end of the term, or the failure of a campaign
// without a term.
func campaign(w *world, p *process, f waits.Campaign) outcome {
	session := p.Name() + "#" + strconv.FormatUint(p.Epoch(), 10)
	c, err := sim.Dependency[waits.Store](w, p).Campaign(w.Context(), f.Scope, session)
	if err != nil {
		return outcome{Event: waits.Campaigned{Err: err, Token: f.Token}}
	}
	var granted uint64
	cancelled := false
	return outcome{
		Next: func() (waits.Event, bool, bool) {
			switch {
			case cancelled && granted != 0:
				return waits.TermEnded{Token: f.Token}, true, true
			case cancelled:
				return waits.Campaigned{Err: context.Canceled, Token: f.Token}, true, true
			case granted == 0 && c.Term() != 0:
				granted = c.Term()
				return waits.Campaigned{Token: f.Token, Term: granted}, true, false
			default:
				return nil, false, false
			}
		},
		Cancel: func() {
			cancelled = true
			c.Resign()
		},
	}
}

// watch watches the terms of the scope of the effect f in the store of the
// leader p. Its wait delivers each change of the term, and ends without an
// event at its cancel.
func watch(w *world, p *process, f waits.Watch) outcome {
	store := sim.Dependency[waits.Store](w, p)
	var last uint64
	cancelled := false
	return outcome{
		Next: func() (waits.Event, bool, bool) {
			if cancelled {
				return nil, false, true
			}
			if term := store.Term(f.Scope); term != last {
				last = term
				return waits.Watched{Token: f.Token, Term: term}, true, false
			}
			return nil, false, false
		},
		Cancel: func() { cancelled = true },
	}
}

// leaderOf returns a leader that the case chooses.
func leaderOf(c *prop.Case) string {
	return leaders[c.Draw(prop.Integer(0, len(leaders)-1), "leader")]
}

// stepDown makes a leader that the case chooses step down.
//
//dokimi:action call
func stepDown(c *prop.Case) func(w *world, client int) {
	leader := leaderOf(c)
	return func(w *world, client int) {
		w.Call(client, leader, "step-down", nil, nil, func(r uint64) waits.Event {
			return waits.StepDown{Reply: waits.ReplyToken(r)}
		})
	}
}

// rejoin makes a leader that the case chooses campaign and watch again.
//
//dokimi:action call
func rejoin(c *prop.Case) func(w *world, client int) {
	leader := leaderOf(c)
	return func(w *world, client int) {
		w.Call(client, leader, "rejoin", nil, nil, func(r uint64) waits.Event {
			return waits.Rejoin{Reply: waits.ReplyToken(r)}
		})
	}
}

// lead waits on a leader that the case chooses until the leader has a term.
//
//dokimi:action call
func lead(c *prop.Case) func(w *world, client int) {
	leader := leaderOf(c)
	return func(w *world, client int) {
		w.Call(client, leader, "lead", nil, nil, func(r uint64) waits.Event {
			return waits.Lead{Reply: waits.ReplyToken(r)}
		})
	}
}

// The store has at most one candidacy of each session. The run checks it after
// every step.
//
//dokimi:invariant defect=leaky
func single(h *sim.Check, w *world) {
	t, _ := sim.Dependency[waits.Store](w, w.Process(leaders[0])).(*terms)
	seen := map[string]bool{}
	for _, c := range t.line {
		expect.False(h, seen[c.session], h.Statement())
		seen[c.session] = true
	}
}

// stepDownThenRejoin makes the leader a step down and then rejoin, once the
// leaders have started.
//
//dokimi:scenario
func stepDownThenRejoin(w *world) {
	drain(w)
	w.Call(0, leaders[0], "step-down", nil, nil, func(r uint64) waits.Event {
		return waits.StepDown{Reply: waits.ReplyToken(r)}
	})
	drain(w)
	w.Call(0, leaders[0], "rejoin", nil, nil, func(r uint64) waits.Event {
		return waits.Rejoin{Reply: waits.ReplyToken(r)}
	})
	drain(w)
}

// drain steps each process with a pending event until no process has one.
func drain(w *world) {
	for pending := w.Pending(); len(pending) > 0; pending = w.Pending() {
		w.Step(pending[0])
	}
}
