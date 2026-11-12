// Package hostedsim contains the simulation of the members of hosted.
package hostedsim

import (
	"context"
	"sync"

	"go.dokimi.dev/assert/expect"
	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"example.com/shop/hosted"
)

// world is the world of the members of hosted, and outcome the outcome of
// their effects.
type (
	world   = sim.World[hosted.Event, hosted.Effect]
	outcome = sim.Outcome[hosted.Event, hosted.Effect]
)

// disk is the storage of the values of one member.
type disk struct {
	values map[string]string
	mu     sync.Mutex
}

// Write writes value under key.
func (d *disk) Write(_ context.Context, key, value string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.values[key] = value
	return nil
}

// Read returns the value under key, and the empty value for a key without
// one.
func (d *disk) Read(_ context.Context, key string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.values[key], nil
}

// directory is the registry of the addresses of the members.
type directory map[string]string

// Lookup returns the address of the member named name, and the empty address
// for a name without one.
func (d directory) Lookup(_ context.Context, name string) (string, error) {
	return d[name], nil
}

// members returns the system of the members names. Each member has a
// storage of its own, and every member shares one registry.
func members(names ...string) func(w *world) {
	return func(w *world) {
		registry := directory{}
		for _, name := range names {
			registry[name] = name + ":7000"
		}
		for _, name := range names {
			storage := &disk{values: map[string]string{}}
			w.Spawn(name, func(*sim.Clock) sim.Machine[hosted.Event, hosted.Effect] {
				return hosted.New()
			}, hosted.Storage(storage), hosted.Registry(registry))
		}
	}
}

// system spawns the members a and b.
//
//dokimi:system
func system(w *world) { members("a", "b")(w) }

// lonely spawns the member a alone.
func lonely(w *world) { members("a")(w) }

// persist writes the value of the effect f under its key through the storage
// of the process p.
func persist(w *world, p *sim.Process[hosted.Event, hosted.Effect], f hosted.Persist) outcome {
	err := sim.Dependency[hosted.Storage](w, p).Write(w.Context(), f.Key, f.Value)
	return outcome{Event: hosted.Persisted{Err: err, Token: f.Token}}
}

// put writes a value that the case chooses under the key k on the member a.
//
//dokimi:action call
func put(c *prop.Case) func(w *world, client int) {
	value := c.Draw(prop.String(), "value")
	return func(w *world, client int) {
		w.Call(client, "a", "put", []any{value}, nil, func(r uint64) hosted.Event {
			return hosted.Put{Key: "k", Value: value, Reply: hosted.ReplyToken(r)}
		})
	}
}

// get reads the value under the key k on the member a.
//
//dokimi:action call
func get(*prop.Case) func(w *world, client int) {
	return func(w *world, client int) {
		w.Call(client, "a", "get", nil, nil, func(r uint64) hosted.Event {
			return hosted.Get{Key: "k", Reply: hosted.ReplyToken(r)}
		})
	}
}

// await waits on the member a for its next put.
//
//dokimi:action call
func await(*prop.Case) func(w *world, client int) {
	return func(w *world, client int) {
		w.Call(client, "a", "await", nil, nil, func(r uint64) hosted.Event {
			return hosted.Await{Reply: hosted.ReplyToken(r)}
		})
	}
}

// find reads the address of the member b on the member a.
//
//dokimi:action call
func find(*prop.Case) func(w *world, client int) {
	return func(w *world, client int) {
		w.Call(client, "a", "find", nil, nil, func(r uint64) hosted.Event {
			return hosted.Find{Name: "b", Reply: hosted.ReplyToken(r)}
		})
	}
}

// The world has the member b. The run checks it after every step.
//
//dokimi:invariant defect=lonely
func spawned(h *sim.Check, w *world) {
	expect.NotNil(h, w.Process("b"), h.Statement())
}
