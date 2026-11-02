// Package remotesim contains the simulation of the nodes of the remote
// register.
package remotesim

import (
	"context"
	"sync"

	"go.dokimi.dev/assert/prop"
	"go.dokimi.dev/sim"

	"example.com/shop/remote"
	"example.com/shop/remote/store"
)

// memory is the writer and the reader of the values of a map, the
// dependencies of each node.
type memory struct {
	values map[string]string
	mu     sync.Mutex
}

// Put stores value under key.
func (m *memory) Put(_ context.Context, key, value string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[key] = value
	return nil
}

// Get returns the value under key, and the empty value for a key without one.
func (m *memory) Get(_ context.Context, key string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.values[key], nil
}

// Read returns the value under key, and the empty value for a key without
// one.
func (m *memory) Read(ctx context.Context, key string) (store.Value, error) {
	value, err := m.Get(ctx, key)
	return store.Value(value), err
}

// system spawns the node a with a writer and a reader of its own.
//
//dokimi:system
func system(w *sim.World[remote.Event, remote.Effect]) {
	m := &memory{values: map[string]string{}}
	w.Spawn("a", newNode, remote.Writer(m), store.Reader(m))
}

// newNode returns a node.
func newNode(*sim.Clock) sim.Machine[remote.Event, remote.Effect] { return remote.New() }

// watch waits on the node a for its stop.
//
//dokimi:action call
func watch(*prop.Case) func(w *sim.World[remote.Event, remote.Effect], client int) {
	return func(w *sim.World[remote.Event, remote.Effect], client int) {
		w.Call(client, "a", "watch", nil, nil, func(r uint64) remote.Event {
			return remote.Watch{Reply: remote.ReplyToken(r)}
		})
	}
}

// write writes a value that the case chooses under the key k of the node a.
//
//dokimi:action call
func write(c *prop.Case) func(w *sim.World[remote.Event, remote.Effect], client int) {
	value := c.Draw(prop.String(), "value")
	return func(w *sim.World[remote.Event, remote.Effect], client int) {
		w.Call(client, "a", "write", []any{value}, []any{"k"}, func(r uint64) remote.Event {
			return remote.Write{Key: "k", Value: value, Reply: remote.ReplyToken(r)}
		})
	}
}
