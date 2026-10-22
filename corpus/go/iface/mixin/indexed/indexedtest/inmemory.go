// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package indexedtest contains the double and the suite of Ranked. InMemory is
// the control of the suite.
package indexedtest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/indexed"
)

// ErrOutOfRange is the error of a read at an index outside the items.
var ErrOutOfRange = errors.New("indexedtest: index out of range")

// InMemory implements [indexed.Ranked] in memory. The corpus applies the
// suite of Ranked to it as a control.
//
//dokimi:subject indexed.Ranked new=NewInMemory
type InMemory struct {
	values []indexed.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [indexed.Ranked].
var _ indexed.Ranked = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Add implements [indexed.Ranked].
func (s *InMemory) Add(ctx context.Context, v indexed.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values = append(s.values, v)
	return nil
}

// Len implements [indexed.Ranked].
func (s *InMemory) Len(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.values), nil
}

// At implements [indexed.Ranked].
func (s *InMemory) At(ctx context.Context, i int) (indexed.Value, error) {
	if err := contextErr(ctx); err != nil {
		return indexed.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if i < 0 || i >= len(s.values) {
		return indexed.Value{}, ErrOutOfRange
	}
	return s.values[i], nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("indexedtest: nil context")
	}
	return ctx.Err()
}
