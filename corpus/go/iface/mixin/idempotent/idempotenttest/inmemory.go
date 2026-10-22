// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package idempotenttest contains the double and the suite of Mixed. InMemory
// is the control of the suite.
package idempotenttest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/idempotent"
)

// ErrNotFound is the error of a read of a key without a value.
var ErrNotFound = errors.New("idempotenttest: not found")

// InMemory implements [idempotent.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject idempotent.Mixed new=NewInMemory
type InMemory struct {
	items  map[string]string
	writes int
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [idempotent.Mixed].
var _ idempotent.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{items: map[string]string{}} }

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("idempotenttest: nil context")
	}
	return ctx.Err()
}

// Put implements [idempotent.Mixed].
func (s *InMemory) Put(ctx context.Context, key, value string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.items[key]; ok && existing == value {
		return nil
	}
	s.items[key] = value
	s.writes++
	return nil
}

// Read implements [idempotent.Mixed].
func (s *InMemory) Read(ctx context.Context, key string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.items[key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}
