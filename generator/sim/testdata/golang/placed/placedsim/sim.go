// Package placedsim contains the simulation of the agents of placed.
package placedsim

import (
	"context"
	"sync"
	"time"

	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"example.com/shop/placed"
)

// world is the world of the agents of placed, and outcome the outcome of
// their effects.
type (
	world   = sim.World[placed.Event, placed.Effect]
	outcome = sim.Outcome[placed.Event, placed.Effect]
)

// memory is the store of the values of the cells, which the agents share.
type memory struct {
	values map[string]string
	mu     sync.Mutex
}

// Write writes value for cell.
func (m *memory) Write(_ context.Context, cell, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[cell] = value
	return nil
}

// Read returns the value of cell, and the empty value for a cell without one.
func (m *memory) Read(_ context.Context, cell string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.values[cell], nil
}

// agents returns the system of the agents names over one store. The case
// chooses the timeout of their writes.
func agents(c *prop.Case, names ...string) func(w *world) {
	timeout := c.Draw(prop.Duration(time.Millisecond, time.Second), "timeout")
	return func(w *world) {
		m := &memory{values: map[string]string{}}
		for _, name := range names {
			w.Spawn(name, func(*sim.Clock) sim.Machine[placed.Event, placed.Effect] {
				return placed.New(timeout)
			}, placed.Store(m))
		}
	}
}

// system spawns the agents a and b over one store. The case chooses the
// timeout of their writes.
//
//dokimi:system
func system(c *prop.Case) func(w *world) { return agents(c, "a", "b") }

// lonely spawns the agent a alone. The case chooses the timeout of its
// writes.
func lonely(c *prop.Case) func(w *world) { return agents(c, "a") }

// write writes the value of the effect f to the cell of the process p.
func write(w *world, p *sim.Process[placed.Event, placed.Effect], f placed.Write) outcome {
	err := sim.Dependency[placed.Store](w, p).Write(w.Context(), p.Name(), f.Value)
	return outcome{Event: placed.Written{Err: err, Token: f.Token}}
}

// assign writes a value that the case chooses on the agent a.
//
//dokimi:action call
func assign(c *prop.Case) func(w *world, client int) {
	value := c.Draw(prop.String(), "value")
	return func(w *world, client int) {
		w.Call(client, "a", "assign", []any{value}, nil, func(r uint64) placed.Event {
			return placed.Assign{Value: value, Reply: placed.ReplyToken(r)}
		})
	}
}

// await waits on the agent a for its next write.
//
//dokimi:action call
func await(*prop.Case) func(w *world, client int) {
	return func(w *world, client int) {
		w.Call(client, "a", "await", nil, nil, func(r uint64) placed.Event {
			return placed.Await{Reply: placed.ReplyToken(r)}
		})
	}
}

// The world has the agent b. The run checks it after every step.
//
//dokimi:invariant defect=lonely
func spawned(h *sim.Check, w *world) {
	expect.NotNil(h, w.Process("b"), h.Statement())
}
