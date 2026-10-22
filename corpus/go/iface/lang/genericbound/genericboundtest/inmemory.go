// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package genericboundtest contains the double and the suite of Ranked.
// InMemory is the control of the suite.
package genericboundtest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/genericbound"
)

// ErrNotFound is the error of Rank for a key without a rank.
var ErrNotFound = errors.New("genericboundtest: not found")

// InMemory implements [genericbound.Ranked] in memory. The corpus applies the
// suite of Ranked to it as a control.
//
//dokimi:subject genericbound.Ranked new=NewInMemory
type InMemory[K genericbound.Ordered, V any] struct {
	ranks map[K]V
	mu    sync.Mutex
}

// The compiler checks that InMemory implements [genericbound.Ranked].
var _ genericbound.Ranked[string, int] = (*InMemory[string, int])(nil)

// NewInMemory returns an empty InMemory at the witnesses of the suite of
// Ranked.
func NewInMemory() *InMemory[int, genericbound.Score] {
	return NewInMemoryOf[int, genericbound.Score]()
}

// NewInMemoryOf returns an empty InMemory.
func NewInMemoryOf[K genericbound.Ordered, V any]() *InMemory[K, V] {
	return &InMemory[K, V]{ranks: map[K]V{}}
}

// Rank implements [genericbound.Ranked].
func (s *InMemory[K, V]) Rank(ctx context.Context, key K) (V, error) {
	var zero V
	if err := contextErr(ctx); err != nil {
		return zero, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.ranks[key]
	if !ok {
		return zero, ErrNotFound
	}
	return v, nil
}

// Reset implements [genericbound.Ranked].
func (s *InMemory[K, V]) Reset(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	clear(s.ranks)
	return nil
}

// Set stores value as the rank of key.
func (s *InMemory[K, V]) Set(key K, value V) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.ranks[key] = value
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("genericboundtest: nil context")
	}
	return ctx.Err()
}
