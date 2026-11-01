// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package readthroughtest contains the double of Store, the suite of Cache and
// the subject of the suite.
package readthroughtest

import (
	"context"
	"fmt"
	"sync"

	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/readthrough"
)

// attempts is the number of attempts of a write of the store.
const attempts = 3

// Cached is the Cache over a store. The corpus applies the suite of Cache to
// it as a subject, over a double of a store in memory whose methods fail on
// demand.
//
//dokimi:subject readthrough.Cache new=NewMemoryCached options=CachedOptions
type Cached struct {
	store  readthrough.Store
	values map[string]readthrough.Value
	mu     sync.Mutex
}

// The compiler checks that Cached implements [readthrough.Cache].
var _ readthrough.Cache = (*Cached)(nil)

// NewCached returns an empty cache over the store.
func NewCached(store readthrough.Store) *Cached {
	return &Cached{store: store, values: map[string]readthrough.Value{}}
}

// NewMemoryCached returns an empty cache over an empty store in memory.
func NewMemoryCached() *Cached { return NewCached(NewMemory()) }

// Get implements [readthrough.Cache].
func (c *Cached) Get(ctx context.Context, key string) (readthrough.Value, error) {
	if err := ctx.Err(); err != nil {
		return readthrough.Value{}, err
	}
	c.mu.Lock()
	v, kept := c.values[key]
	c.mu.Unlock()
	if kept {
		return v, nil
	}
	v, err := c.store.Get(ctx, key)
	if err != nil {
		return readthrough.Value{}, fmt.Errorf("readthroughtest: get %s: %w", key, err)
	}
	c.keep(key, v)
	return v, nil
}

// Put implements [readthrough.Cache].
func (c *Cached) Put(ctx context.Context, v readthrough.Value) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	var err error
	for range attempts {
		if err = c.store.Put(ctx, v); err == nil {
			c.keep(v.Key, v)
			return nil
		}
	}
	return fmt.Errorf("readthroughtest: put %s: %w", v.Key, err)
}

// keep keeps the value under the key in memory.
func (c *Cached) keep(key string, v readthrough.Value) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[key] = v
}

// Memory is the Store in memory.
type Memory struct {
	values map[string]readthrough.Value
	mu     sync.Mutex
}

// The compiler checks that Memory implements [readthrough.Store].
var _ readthrough.Store = (*Memory)(nil)

// NewMemory returns an empty store.
func NewMemory() *Memory { return &Memory{values: map[string]readthrough.Value{}} }

// Get implements [readthrough.Store].
func (m *Memory) Get(ctx context.Context, key string) (readthrough.Value, error) {
	if err := ctx.Err(); err != nil {
		return readthrough.Value{}, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	v, found := m.values[key]
	if !found {
		return readthrough.Value{}, readthrough.ErrMissing
	}
	return v, nil
}

// Put implements [readthrough.Store].
func (m *Memory) Put(ctx context.Context, v readthrough.Value) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[v.Key] = v
	return nil
}

// CachedOptions returns the door of a cache over a double of a store in
// memory, whose methods fail on demand.
func CachedOptions() []CacheOption {
	return []CacheOption{conform.Faultable(faulty)}
}

// faulty returns an empty cache over a double of an empty store in memory.
// Each call of a method of the double asks the injector first.
func faulty(h *CacheHandle, inj *conform.Injector) readthrough.Cache {
	store := inj.Dependency()
	double := NewStoreStub(h, StoreStubDelegateTo(NewMemory()))
	double.OnGet.SetFault(conform.FaultOf[StoreGetCall](store, "Get"))
	double.OnPut.SetFault(conform.FaultOf[StorePutCall](store, "Put"))
	return NewCached(double)
}
