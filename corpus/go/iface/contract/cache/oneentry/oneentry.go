// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package oneentry contains a control of the suite of Contract: a cache of one
// entry in front of a backing store.
package oneentry

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/cache"
)

// ErrNotFound is the error of a key without a value in the backing store.
var ErrNotFound = errors.New("oneentry: no value under that key")

// Cache implements [cache.Contract] with a cache of one entry. Each Lookup that
// misses replaces the entry with the value of the backing store, and each
// Store empties the entry. The corpus applies the suite of Contract to it as a
// control of the caches.
//
//dokimi:subject cache.Contract new=New
type Cache struct {
	backing map[string]cache.Value
	key     string
	entry   cache.Value
	mu      sync.Mutex
	full    bool
}

// The compiler checks that Cache implements [cache.Contract].
var _ cache.Contract = (*Cache)(nil)

// New returns a Cache over an empty backing store.
func New() *Cache { return &Cache{backing: map[string]cache.Value{}} }

// Store implements [cache.Contract].
func (c *Cache) Store(ctx context.Context, key string, v cache.Value) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.backing[key] = v
	c.full = false
	return nil
}

// Lookup implements [cache.Contract].
func (c *Cache) Lookup(ctx context.Context, key string) (cache.Value, error) {
	if err := ctx.Err(); err != nil {
		return cache.Value{}, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.full && c.key == key {
		return c.entry, nil
	}
	v, present := c.backing[key]
	if !present {
		return cache.Value{}, ErrNotFound
	}
	c.key, c.entry, c.full = key, v, true
	return v, nil
}

// Fetch implements [cache.Contract].
func (c *Cache) Fetch(ctx context.Context, key string) (cache.Value, error) {
	if err := ctx.Err(); err != nil {
		return cache.Value{}, err
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	v, present := c.backing[key]
	if !present {
		return cache.Value{}, ErrNotFound
	}
	return v, nil
}
