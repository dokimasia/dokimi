// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package stickytest contains the double and the suite of Mixed. InMemory is
// the subject of the suite.
package stickytest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/sticky"
)

// ErrNotFound is the error of a read of a key without a value.
var ErrNotFound = errors.New("stickytest: not found")

// InMemory implements [sticky.Mixed] in memory. The corpus applies the suite
// of Mixed to it as a subject.
//
//dokimi:subject sticky.Mixed new=NewInMemory
type InMemory struct {
	values   map[string]sticky.Value
	resolved map[string]sticky.Value
	mu       sync.Mutex
}

// The compiler checks that InMemory implements [sticky.Mixed].
var _ sticky.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{
		values:   map[string]sticky.Value{},
		resolved: map[string]sticky.Value{},
	}
}

// Store implements [sticky.Mixed].
func (s *InMemory) Store(ctx context.Context, v sticky.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = v
	return nil
}

// Get implements [sticky.Mixed].
func (s *InMemory) Get(ctx context.Context, key string) (sticky.Value, error) {
	if err := contextErr(ctx); err != nil {
		return sticky.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, pinned := s.resolved[key]; pinned {
		return v, nil
	}
	v, ok := s.values[key]
	if !ok {
		return sticky.Value{}, ErrNotFound
	}
	s.resolved[key] = v
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("stickytest: nil context")
	}
	return ctx.Err()
}
