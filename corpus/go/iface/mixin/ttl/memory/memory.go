// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package memory is the Mixed of the fixture ttl in memory.
package memory

import (
	"context"
	"sync"
	"time"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/ttl"
)

// Lifetime is how long Read returns a value after Put stored it. It is the
// duration of the mixin ttl of Read.
const Lifetime = time.Minute

// Clock is the clock on which a store measures the lifetime of its values.
type Clock interface {
	// Now returns the current instant.
	Now() time.Time
}

// Store keeps each value for its lifetime on a clock.
type Store struct {
	clock  Clock
	values map[string]entry
	mu     sync.Mutex
}

// entry is a stored value and the instant at which its lifetime started.
type entry struct {
	at    time.Time
	value ttl.Value
}

// New returns an empty store on clock.
func New(clock Clock) *Store {
	return &Store{clock: clock, values: map[string]entry{}}
}

// Put stores the value and starts its lifetime.
func (s *Store) Put(ctx context.Context, v ttl.Value) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = entry{value: v, at: s.clock.Now()}
	return nil
}

// Read returns the value under the key until its lifetime ends. It returns
// ErrExpired for a key without a value and for a value whose lifetime has
// ended.
func (s *Store) Read(ctx context.Context, key string) (ttl.Value, error) {
	if err := ctx.Err(); err != nil {
		return ttl.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, found := s.values[key]
	if !found || s.clock.Now().Sub(e.at) >= Lifetime {
		return ttl.Value{}, ttl.ErrExpired
	}
	return e.value, nil
}
