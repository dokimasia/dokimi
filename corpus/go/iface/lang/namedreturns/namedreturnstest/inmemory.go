// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package namedreturnstest contains the double and the suite of Service.
// InMemory is the subject of the suite.
package namedreturnstest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/namedreturns"
)

// ErrNotFound is the error of each reader for an identity without a value.
var ErrNotFound = errors.New("namedreturnstest: not found")

// InMemory implements [namedreturns.Service] in memory. The corpus applies
// the suite of Service to it as a subject.
//
//dokimi:subject namedreturns.Service new=NewInMemory
type InMemory struct {
	items map[string]string
	mu    sync.Mutex
}

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{items: map[string]string{}} }

// Put stores value under id. A test seeds the service through it, because the
// interface declares no writer.
func (s *InMemory) Put(id, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[id] = value
}

// Named implements [namedreturns.Service].
func (s *InMemory) Named(ctx context.Context, id string) (item string, err error) {
	return s.lookup(ctx, id)
}

// Unnamed implements [namedreturns.Service].
func (s *InMemory) Unnamed(ctx context.Context, id string) (string, error) {
	return s.lookup(ctx, id)
}

// PartiallyNamed implements [namedreturns.Service].
func (s *InMemory) PartiallyNamed(ctx context.Context, id string) (item string, _ error) {
	return s.lookup(ctx, id)
}

// lookup returns the value under id. The three readers share it, so they
// differ in their declarations alone.
func (s *InMemory) lookup(ctx context.Context, id string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.items[id]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("namedreturnstest: nil context")
	}
	return ctx.Err()
}

// The compiler checks that InMemory implements [namedreturns.Service].
var _ namedreturns.Service = (*InMemory)(nil)
