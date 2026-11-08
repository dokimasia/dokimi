// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package placedsim contains the simulation of the agents and the cells of
// placed: its system, its performers, its actions and its oracles.
package placedsim

import (
	"context"
	"slices"
	"strconv"
	"sync"
	"time"

	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"go.dokimi.dev/dokimi/corpus/go/sim/placed"
)

// world is the world of the machines of placed, process one of its
// processes, and outcome the outcome of one of their effects.
type (
	world   = sim.World[placed.Event, placed.Effect]
	process = sim.Process[placed.Event, placed.Effect]
	outcome = sim.Outcome[placed.Event, placed.Effect]
)

// hosts are the hosts of the system. The agent of each host can start the
// cells 1 to cells of the host.
var hosts = []string{"east", "west"}

const cells = 2

// stamp returns the time of a write of the cell p.
type stamp func(p *process) time.Time

// record is one write of the store: its time, its cell and its value.
type record struct {
	at    time.Time
	cell  string
	value string
}

// memory is the store of the cells. It keeps the writes in the order in
// which the cells write them.
type memory struct {
	records []record
	mu      sync.Mutex
}

// Write appends the write of value for cell at the time at.
func (m *memory) Write(_ context.Context, cell, value string, at time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.records = append(m.records, record{at: at, cell: cell, value: value})
	return nil
}

// agentOf returns the name of the agent of the host.
func agentOf(host string) string { return host + "-agent" }

// cellOf returns the name of the cell n of the host.
func cellOf(host string, n int) string { return host + "-cell-" + strconv.Itoa(n) }

// spawned returns the system of an agent on each host over one store, whose
// writes take their times from the stamp that at returns. The case chooses
// the timeout of a write.
func spawned(c *prop.Case, at func(w *world) stamp) func(w *world) {
	timeout := c.Draw(prop.Duration(time.Millisecond, time.Second), "timeout")
	return func(w *world) {
		store := &memory{}
		for _, host := range hosts {
			w.AddHost(host).Spawn(agentOf(host), func(*sim.Clock) sim.Machine[placed.Event, placed.Effect] {
				return placed.NewAgent(timeout)
			}, placed.Store(store), at(w))
		}
	}
}

// system spawns an agent on each host over one store, which stamps each
// write with the time of the run. The case chooses the timeout of a write.
//
//dokimi:system
func system(c *prop.Case) func(w *world) {
	return spawned(c, func(w *world) stamp { return func(*process) time.Time { return w.Now() } })
}

// skewed spawns an agent on each host over one store, which stamps each write
// with the clock of the host of its cell. The case chooses the timeout of a
// write.
func skewed(c *prop.Case) func(w *world) {
	return spawned(c, func(*world) stamp { return func(p *process) time.Time { return p.Clock().Now() } })
}

// place starts the cell of the effect f on the host of the agent p. The cell
// takes the store and the stamp of the agent, and its reply to its start is
// the result of f.
func place(w *world, p *process, f placed.Place) outcome {
	return outcome{Launch: &sim.Launch[placed.Event, placed.Effect]{
		Machine: func(*sim.Clock) sim.Machine[placed.Event, placed.Effect] { return placed.NewCell(f.Timeout) },
		Call:    func(reply uint64) placed.Event { return placed.Started{Reply: placed.ReplyToken(reply)} },
		Result:  func(_ any, err error) placed.Event { return placed.Placed{Err: err, Token: f.Token} },
		Name:    f.Cell,
		Deps:    []any{sim.Dependency[placed.Store](w, p), sim.Dependency[stamp](w, p)},
	}}
}

// unplace stops the cell of the effect f.
func unplace(_ *world, _ *process, f placed.Unplace) outcome {
	return outcome{Event: placed.Unplaced{Token: f.Token}, Stop: f.Cell}
}

// write writes the value of the effect f of the cell p to the store, at the
// time of the stamp of the cell.
func write(w *world, p *process, f placed.Write) outcome {
	at := sim.Dependency[stamp](w, p)(p)
	err := sim.Dependency[placed.Store](w, p).Write(w.Context(), p.Name(), f.Value, at)
	return outcome{Event: placed.Written{Err: err, Token: f.Token}}
}

// target returns a host and a cell of the host that the case chooses.
func target(c *prop.Case) (host, cell string) {
	host = hosts[c.Draw(prop.Integer(0, len(hosts)-1), "host")]
	return host, cellOf(host, c.Draw(prop.Integer(1, cells), "cell"))
}

// build builds the event of a call from its reply token.
type build func(reply placed.ReplyToken) placed.Event

// order makes the agent of host start or stop the cell as a call of client.
// The history records the operation on the cell.
func order(w *world, client int, host, operation, cell string, event build) {
	w.Call(client, agentOf(host), operation, []any{cell}, []any{cell}, func(reply uint64) placed.Event {
		return event(placed.ReplyToken(reply))
	})
}

// ask calls the cell as a call of client, where an agent has started the
// cell. The history records the operation with the args.
func ask(w *world, client int, cell, operation string, args []any, event build) {
	if w.Process(cell) == nil {
		return
	}
	w.Call(client, cell, operation, args, []any{cell}, func(reply uint64) placed.Event {
		return event(placed.ReplyToken(reply))
	})
}

// add makes the agent of a host that the case chooses start a cell of the
// host.
//
//dokimi:action call
func add(c *prop.Case) func(w *world, client int) {
	host, cell := target(c)
	return func(w *world, client int) {
		order(w, client, host, "add", cell, func(r placed.ReplyToken) placed.Event {
			return placed.Add{Cell: cell, Reply: r}
		})
	}
}

// drop makes the agent of a host that the case chooses stop a cell of the
// host.
//
//dokimi:action call
func drop(c *prop.Case) func(w *world, client int) {
	host, cell := target(c)
	return func(w *world, client int) {
		order(w, client, host, "drop", cell, func(r placed.ReplyToken) placed.Event {
			return placed.Drop{Cell: cell, Reply: r}
		})
	}
}

// set writes a value that the case chooses to a cell that the case chooses.
//
//dokimi:action call
func set(c *prop.Case) func(w *world, client int) {
	_, cell := target(c)
	value := strconv.Itoa(c.Draw(prop.Integer(0, 9), "value"))
	return func(w *world, client int) { setOn(w, client, cell, value) }
}

// setOn writes value to the cell as a call of client, where an agent has
// started the cell.
func setOn(w *world, client int, cell, value string) {
	ask(w, client, cell, "set", []any{value}, func(r placed.ReplyToken) placed.Event {
		return placed.Set{Value: value, Reply: r}
	})
}

// await waits for the next write of a cell that the case chooses.
//
//dokimi:action call
func await(c *prop.Case) func(w *world, client int) {
	_, cell := target(c)
	return func(w *world, client int) {
		ask(w, client, cell, "await", nil, func(r placed.ReplyToken) placed.Event { return placed.Await{Reply: r} })
	}
}

// The store keeps the writes in the order of their times. The run checks it
// after every step.
//
//dokimi:invariant defect=skewed
func ordered(h *sim.Check, w *world) {
	m, _ := sim.Dependency[placed.Store](w, w.Process(agentOf(hosts[0]))).(*memory)
	m.mu.Lock()
	defer m.mu.Unlock()
	expect.True(h, slices.IsSortedFunc(m.records, func(a, b record) int { return a.at.Compare(b.at) }),
		h.Statement())
}

// placeThenSet starts the cell 1 of the host east and writes to it twice.
// The clock of east jumps back by a minute between the two writes.
//
//dokimi:scenario
func placeThenSet(w *world) {
	host, cell := hosts[0], cellOf(hosts[0], 1)
	order(w, 0, host, "add", cell, func(r placed.ReplyToken) placed.Event { return placed.Add{Cell: cell, Reply: r} })
	drain(w)
	setOn(w, 1, cell, "1")
	drain(w)
	w.Jump(host, -time.Minute)
	setOn(w, 1, cell, "2")
	drain(w)
}

// drain steps each process with a pending event until no process has one.
func drain(w *world) {
	for pending := w.Pending(); len(pending) > 0; pending = w.Pending() {
		w.Step(pending[0])
	}
}
