// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package memory is the Store of the fixture generic in memory.
package memory

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/generic"
)

// ErrNotFound reports a key without a value.
var ErrNotFound = errors.New("memory: the store has no value under the key")

// Store is the store of values by key in memory.
//
//dokimi:subject generic.Store new=NewInts
type Store[K comparable, V any] struct {
	values map[K]V
	mu     sync.Mutex
}

// NewInts returns an empty store of integers by integer, the witnesses of the
// suite of generic.Store.
func NewInts() *Store[int, int] {
	return &Store[int, int]{values: map[int]int{}}
}

// Get returns the value under the key. It returns ErrNotFound for a key
// without a value.
func (s *Store[K, V]) Get(ctx context.Context, key K) (V, error) {
	var zero V
	if err := ctx.Err(); err != nil {
		return zero, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	value, found := s.values[key]
	if !found {
		return zero, ErrNotFound
	}
	return value, nil
}

// Put stores the value under the key.
func (s *Store[K, V]) Put(ctx context.Context, key K, value V) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = value
	return nil
}

// Compile-time proof that Store at the witnesses satisfies generic.Store.
var _ generic.Store[int, int] = (*Store[int, int])(nil)
