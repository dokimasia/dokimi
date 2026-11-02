// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package replicatedsim contains the simulation of the nodes of the
// register: its system, its actions, its spec and its oracles.
package replicatedsim

import (
	"maps"
	"time"

	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/history"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"go.dokimi.dev/dokimi/corpus/go/sim/replicated"
)

// world is the world of the nodes of the register.
type world = sim.World[replicated.Event, replicated.Effect]

// The keys and the values of the calls of the clients.
var (
	keys   = []string{"x", "y"}
	values = []string{"1", "2", "3"}
)

// The operations of the calls of the clients in the history.
const (
	putOperation   = "put"
	getOperation   = "get"
	watchOperation = "watch"
)

// system spawns the primary and the backup, each with a disk of its own.
//
//dokimi:system
func system(w *world) {
	primary, backup := replicated.Disk{}, replicated.Disk{}
	w.Spawn("primary", func(*sim.Clock) sim.Machine[replicated.Event, replicated.Effect] {
		return replicated.NewPrimary(primary)
	})
	w.Spawn("backup", func(*sim.Clock) sim.Machine[replicated.Event, replicated.Effect] {
		return replicated.NewBackup(backup)
	})
}

// hasty spawns a primary that applies and acknowledges each write before the
// backup has it, and the backup.
func hasty(w *world) {
	primary, backup := replicated.Disk{}, replicated.Disk{}
	w.Spawn("primary", func(*sim.Clock) sim.Machine[replicated.Event, replicated.Effect] {
		return &hastyPrimary{disk: primary}
	})
	w.Spawn("backup", func(*sim.Clock) sim.Machine[replicated.Event, replicated.Effect] {
		return replicated.NewBackup(backup)
	})
}

// hastyPrimary is a primary that applies and acknowledges each write at
// once, and copies it to the backup afterwards.
type hastyPrimary struct {
	disk   replicated.Disk
	tokens replicated.Token
}

// Step applies ev and appends its effects to out.
func (h *hastyPrimary) Step(ev replicated.Event, out []replicated.Effect) []replicated.Effect {
	switch ev := ev.(type) {
	case replicated.Started:
		return append(out, replicated.Reply{To: ev.Reply})
	case replicated.Put:
		h.tokens++
		h.disk[ev.Key] = ev.Value
		return append(out, replicated.Reply{To: ev.Reply},
			replicated.Send{Key: ev.Key, Value: ev.Value, Token: h.tokens, Timeout: time.Second})
	case replicated.Get:
		return append(out, replicated.Reply{Value: h.disk[ev.Key], To: ev.Reply})
	default:
		return out
	}
}

// Err returns nil, because the primary never stops itself.
func (*hastyPrimary) Err() error { return nil }

// send copies the value of the effect f to the backup, as a call whose reply
// is the result of f.
func send(
	_ *world, _ *sim.Process[replicated.Event, replicated.Effect], f replicated.Send,
) sim.Outcome[replicated.Event] {
	return sim.Outcome[replicated.Event]{Message: &sim.Message[replicated.Event]{
		To: "backup",
		Call: func(reply uint64) replicated.Event {
			return replicated.Copy{Key: f.Key, Value: f.Value, Reply: replicated.ReplyToken(reply)}
		},
		Result: func(_ any, err error) replicated.Event { return replicated.Copied{Err: err, Token: f.Token} },
	}}
}

// call returns the run of a call of a client to the primary, which records
// the operation with the args under the key in the history, and builds the
// event of the call from its reply token.
func call(operation string, args []any, key string, build func(replicated.ReplyToken) replicated.Event,
) func(w *world, client int) {
	return func(w *world, client int) {
		w.Call(client, "primary", operation, args, []any{key}, func(reply uint64) replicated.Event {
			return build(replicated.ReplyToken(reply))
		})
	}
}

// put writes a value that the case chooses under a key that it chooses.
//
//dokimi:action call
func put(c *prop.Case) func(w *world, client int) {
	key := keys[c.Draw(prop.Integer(0, len(keys)-1), "key")]
	value := values[c.Draw(prop.Integer(0, len(values)-1), "value")]
	return call(putOperation, []any{key, value}, key, func(r replicated.ReplyToken) replicated.Event {
		return replicated.Put{Key: key, Value: value, Reply: r}
	})
}

// get reads a key that the case chooses.
//
//dokimi:action call
func get(c *prop.Case) func(w *world, client int) {
	key := keys[c.Draw(prop.Integer(0, len(keys)-1), "key")]
	return call(getOperation, []any{key}, key, func(r replicated.ReplyToken) replicated.Event {
		return replicated.Get{Key: key, Reply: r}
	})
}

// watch waits for the next write of a key that the case chooses.
//
//dokimi:action call
func watch(c *prop.Case) func(w *world, client int) {
	key := keys[c.Draw(prop.Integer(0, len(keys)-1), "key")]
	return call(watchOperation, []any{key}, key, func(r replicated.ReplyToken) replicated.Event {
		return replicated.Watch{Key: key, Reply: r}
	})
}

// readX reads the key x once the faults have ended.
//
//dokimi:action call final
func readX(*prop.Case) func(w *world, client int) {
	return call(getOperation, []any{"x"}, "x", func(r replicated.ReplyToken) replicated.Event {
		return replicated.Get{Key: "x", Reply: r}
	})
}

// spec is the spec of a register of keys: a put writes its value and
// returns nothing, a get returns the value of its key, and a watch returns
// the value of the write of its key that ends it.
//
//dokimi:spec
func spec() history.Spec[map[string]string] {
	return history.Spec[map[string]string]{
		Initial: func() map[string]string { return map[string]string{} },
		Next: func(state map[string]string, op history.Operation) []map[string]string {
			key, _ := op.Args[0].(string)
			if op.Name != putOperation {
				if !op.Returned(state[key]) {
					return nil
				}
				return []map[string]string{state}
			}
			if !op.Returned(nil) {
				return nil
			}
			value, _ := op.Args[1].(string)
			next := maps.Clone(state)
			next[key] = value
			return []map[string]string{next}
		},
	}
}

// The backup holds the value of each acknowledged write. The run checks it
// after every step.
//
//dokimi:invariant defect=hasty
func durable(h *sim.Check, w *world) {
	copied := map[[2]string]bool{}
	for _, e := range w.Log() {
		if c, isCopy := e.Event.(replicated.Copy); isCopy && e.Process == "backup" && !e.Dropped {
			copied[[2]string{c.Key, c.Value}] = true
		}
	}
	events := w.History().Events()
	for i := range events {
		invoked := &events[events[i].Call]
		if events[i].Kind != history.OK || invoked.Operation != putOperation {
			continue
		}
		key, _ := invoked.Args[0].(string)
		value, _ := invoked.Args[1].(string)
		expect.True(h, copied[[2]string{key, value}], h.Statement())
	}
}

// Each write ends once the faults have ended.
//
//dokimi:converges within=settled
func ended(h *sim.Check, w *world) {
	open := map[int]bool{}
	events := w.History().Events()
	for i := range events {
		switch e := &events[i]; {
		case e.Kind == history.Invoke && e.Operation == putOperation:
			open[e.Index] = true
		case e.Kind != history.Invoke:
			delete(open, e.Call)
		}
	}
	expect.Empty(h, open, h.Statement())
}

// settled returns the bound within which each write ends: the timeout of a
// copy on the slowest clock, twice.
func settled() time.Duration { return 4 * time.Second }

// fingerprint returns the number of the steps of the log.
//
//dokimi:fingerprint
func fingerprint(w *world) uint64 { return uint64(len(w.Log())) }

// putThenGet writes 1 under x on the primary, and then reads x.
//
//dokimi:scenario
func putThenGet(w *world) {
	call(putOperation, []any{"x", "1"}, "x", func(r replicated.ReplyToken) replicated.Event {
		return replicated.Put{Key: "x", Value: "1", Reply: r}
	})(w, 0)
	drain(w)
	call(getOperation, []any{"x"}, "x", func(r replicated.ReplyToken) replicated.Event {
		return replicated.Get{Key: "x", Reply: r}
	})(w, 1)
	drain(w)
}

// drain steps each process with a pending event until no process has one.
func drain(w *world) {
	for pending := w.Pending(); len(pending) > 0; pending = w.Pending() {
		w.Step(pending[0])
	}
}
