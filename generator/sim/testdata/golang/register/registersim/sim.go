// Package registersim contains the simulation of the nodes of the register.
package registersim

import (
	"time"

	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/history"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"example.com/shop/register"
)

// system spawns the nodes a and b.
//
//dokimi:system
func system(w *sim.World[register.Event, register.Effect]) {
	w.Spawn("a", newNode)
	w.Spawn("b", newNode)
}

// lonely spawns the node a alone, so a write to b fails the run.
func lonely(w *sim.World[register.Event, register.Effect]) {
	w.Spawn("a", newNode)
}

// newNode returns a node.
func newNode(*sim.Clock) sim.Machine[register.Event, register.Effect] { return register.New() }

// commit performs a commit on a disk that never fails.
func commit(
	_ *sim.World[register.Event, register.Effect], _ *sim.Process[register.Event, register.Effect], f register.Commit,
) sim.Outcome[register.Event, register.Effect] {
	return sim.Outcome[register.Event, register.Effect]{Event: register.Committed{Token: f.Token}}
}

// write writes a value that the case chooses on the node a.
//
//dokimi:action call
func write(c *prop.Case) func(w *sim.World[register.Event, register.Effect], client int) {
	value := c.Draw(prop.String(), "value")
	return func(w *sim.World[register.Event, register.Effect], client int) {
		w.Call(client, "a", "write", []any{value}, nil, func(r uint64) register.Event {
			return register.Write{Value: value, Reply: register.ReplyToken(r)}
		})
	}
}

// watch waits on the node a for its stop.
//
//dokimi:action call
func watch(*prop.Case) func(w *sim.World[register.Event, register.Effect], client int) {
	return func(w *sim.World[register.Event, register.Effect], client int) {
		w.Call(client, "a", "watch", nil, nil, func(r uint64) register.Event {
			return register.Watch{Reply: register.ReplyToken(r)}
		})
	}
}

// settle writes the empty value on the node a once the faults end.
//
//dokimi:action call final
func settle(*prop.Case) func(w *sim.World[register.Event, register.Effect], client int) {
	return func(w *sim.World[register.Event, register.Effect], client int) {
		w.Call(client, "a", "write", []any{""}, nil, func(r uint64) register.Event {
			return register.Write{Reply: register.ReplyToken(r)}
		})
	}
}

// pause advances the time of the run by a second.
//
//dokimi:action
func pause(*prop.Case) func(w *sim.World[register.Event, register.Effect]) {
	return func(w *sim.World[register.Event, register.Effect]) { w.Advance(time.Second) }
}

// rest advances the time of the run by a second once the faults end.
//
//dokimi:action final
func rest(*prop.Case) func(w *sim.World[register.Event, register.Effect]) {
	return func(w *sim.World[register.Event, register.Effect]) { w.Advance(time.Second) }
}

// spec is the spec of a register whose write returns the value that it
// writes, and whose watch changes nothing.
//
//dokimi:spec
func spec() history.Spec[string] {
	return history.Spec[string]{
		Initial: func() string { return "" },
		Next: func(state string, op history.Operation) []string {
			if op.Name == "watch" {
				return []string{state}
			}
			v, _ := op.Args[0].(string)
			if !op.Returned(v) {
				return nil
			}
			return []string{v}
		},
	}
}

// The system has the node b. The run checks it after every step.
//
//dokimi:invariant defect=lonely
func spawned(h *sim.Check, w *sim.World[register.Event, register.Effect]) {
	expect.NotNil(h, w.Process("b"), h.Statement())
}

// The log of the run has a step once the faults have ended.
//
//dokimi:converges within=bound
func replied(h *sim.Check, w *sim.World[register.Event, register.Effect]) {
	expect.NotEmpty(h, w.Log(), h.Statement())
}

// bound returns the bound of the convergence replied.
func bound() time.Duration { return time.Minute }

// fingerprint returns the number of the steps of the log.
//
//dokimi:fingerprint
func fingerprint(w *sim.World[register.Event, register.Effect]) uint64 { return uint64(len(w.Log())) }

// writeOnA writes the value x on the node a.
//
//dokimi:scenario
func writeOnA(w *sim.World[register.Event, register.Effect]) {
	w.Call(0, "a", "write", []any{"x"}, nil, func(r uint64) register.Event {
		return register.Write{Value: "x", Reply: register.ReplyToken(r)}
	})
	for pending := w.Pending(); len(pending) > 0; pending = w.Pending() {
		w.Step(pending[0])
	}
}
