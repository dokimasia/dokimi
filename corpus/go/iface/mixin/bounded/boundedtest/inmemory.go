// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package boundedtest contains the double and the suite of Mixed. InMemory is
// the control of the suite.
package boundedtest

import (
	"context"
	"errors"
	"slices"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/bounded"
)

// limit is the limit of the mixin bounded of List.
const limit = 5

// InMemory implements [bounded.Mixed] in memory. The corpus applies the suite
// of Mixed to it as a control.
//
//dokimi:subject bounded.Mixed new=NewInMemory
type InMemory struct {
	items []string
	limit int
	mu    sync.Mutex
}

// The compiler checks that InMemory implements [bounded.Mixed].
var _ bounded.Mixed = (*InMemory)(nil)

// NewInMemory returns an empty InMemory whose List returns at most the limit of
// the mixin bounded.
func NewInMemory() *InMemory { return NewInMemoryWith(limit) }

// NewInMemoryWith returns an empty InMemory whose List returns at most limit
// items.
func NewInMemoryWith(limit int) *InMemory { return &InMemory{limit: limit} }

// Add implements [bounded.Mixed].
func (s *InMemory) Add(ctx context.Context, item string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, item)
	return nil
}

// List implements [bounded.Mixed]. It returns the first items in the order of
// their additions.
func (s *InMemory) List(ctx context.Context) ([]string, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.items[:min(len(s.items), s.limit)]), nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("boundedtest: nil context")
	}
	return ctx.Err()
}
