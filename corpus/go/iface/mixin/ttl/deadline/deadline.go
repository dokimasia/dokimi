// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package deadline contains a control of the suite of Mixed: a store that
// keeps the deadline of each value in place of the instant of its write.
package deadline

import (
	"context"
	"sync"
	"time"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/ttl"
	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/ttl/ttltest"
)

// lifetime is the duration of the mixin ttl of Read.
const lifetime = time.Minute

// Store implements [ttl.Mixed] over the clock of the handle. Each entry keeps
// the instant at which its value expires. The corpus applies the suite of
// Mixed to it as a control of the clocked declarations.
//
//dokimi:subject ttl.Mixed new=NewFromHandle
type Store struct {
	clk     assert.Clock
	entries map[string]entry
	mu      sync.Mutex
}

// entry is a value with the instant at which it expires.
type entry struct {
	until time.Time
	value ttl.Value
}

// The compiler checks that Store implements [ttl.Mixed].
var _ ttl.Mixed = (*Store)(nil)

// NewFromHandle returns an empty Store on the clock of h.
func NewFromHandle(h *ttltest.MixedHandle) *Store {
	return &Store{clk: h.Clock(), entries: map[string]entry{}}
}

// Put implements [ttl.Mixed].
func (s *Store) Put(ctx context.Context, v ttl.Value) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[v.Key] = entry{until: s.clk.Now().Add(lifetime), value: v}
	return nil
}

// Read implements [ttl.Mixed].
func (s *Store) Read(ctx context.Context, key string) (ttl.Value, error) {
	if err := ctx.Err(); err != nil {
		return ttl.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[key]
	if !ok || !s.clk.Now().Before(e.until) {
		return ttl.Value{}, ttl.ErrExpired
	}
	return e.value, nil
}
