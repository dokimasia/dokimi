// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package boundedcachetest contains the double and the suite of Mixed. InMemory
// is the subject of the suite.
package boundedcachetest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/boundedcache"
)

// capacity is the limit of the mixin bounded of Put.
const capacity = 2

// InMemory implements [boundedcache.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a subject. At its capacity, Put evicts the oldest
// key.
//
//dokimi:subject boundedcache.Mixed new=NewInMemory
type InMemory struct {
	entries  map[boundedcache.Key]boundedcache.Value
	order    []boundedcache.Key
	capacity int
	mu       sync.Mutex
}

// The compiler checks that InMemory implements [boundedcache.Mixed].
var _ boundedcache.Mixed = (*InMemory)(nil)

// NewInMemory returns an empty cache at the capacity of the mixin bounded.
func NewInMemory() *InMemory { return NewInMemoryWith(capacity) }

// NewInMemoryWith returns an empty cache of the capacity.
func NewInMemoryWith(capacity int) *InMemory {
	return &InMemory{
		capacity: capacity,
		entries:  map[boundedcache.Key]boundedcache.Value{},
	}
}

// Put implements [boundedcache.Mixed].
func (s *InMemory) Put(ctx context.Context, key boundedcache.Key, value boundedcache.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, stored := s.entries[key]; !stored {
		if len(s.order) >= s.capacity {
			oldest := s.order[0]
			s.order = s.order[1:]
			delete(s.entries, oldest)
		}
		s.order = append(s.order, key)
	}
	s.entries[key] = value
	return nil
}

// Get implements [boundedcache.Mixed].
func (s *InMemory) Get(ctx context.Context, key boundedcache.Key) (boundedcache.Value, bool) {
	if err := contextErr(ctx); err != nil {
		return boundedcache.Value{}, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, stored := s.entries[key]
	return v, stored
}

// Len implements [boundedcache.Mixed].
func (s *InMemory) Len(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.entries), nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("boundedcachetest: nil context")
	}
	return ctx.Err()
}
