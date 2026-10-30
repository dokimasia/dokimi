// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package defaultonerrortest contains the double and the suite of Mixed.
// InMemory is the subject of the suite.
package defaultonerrortest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/defaultonerror"
)

// ErrNotFound is the error of Get for a key without a body. Get returns
// [defaultonerror.DefaultBody] beside it.
var ErrNotFound = errors.New("defaultonerrortest: not found")

// InMemory implements [defaultonerror.Mixed] in memory. The corpus applies
// the suite of Mixed to it as a subject.
//
//dokimi:subject defaultonerror.Mixed new=NewInMemory
type InMemory struct {
	values map[string]defaultonerror.Body
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [defaultonerror.Mixed].
var _ defaultonerror.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{values: map[string]defaultonerror.Body{}}
}

// Store implements [defaultonerror.Mixed].
func (s *InMemory) Store(ctx context.Context, key string, body defaultonerror.Body) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = body
	return nil
}

// Get implements [defaultonerror.Mixed].
func (s *InMemory) Get(ctx context.Context, key string) (defaultonerror.Body, error) {
	if err := contextErr(ctx); err != nil {
		return defaultonerror.DefaultBody, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.values[key]
	if !ok {
		return defaultonerror.DefaultBody, ErrNotFound
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("defaultonerrortest: nil context")
	}
	return ctx.Err()
}

// Reopen returns an InMemory with the bodies of prior. A restart after a crash
// keeps every body that Store wrote. Do not use prior after the call. The two
// stores share the bodies but not the lock.
func Reopen(prior *InMemory) *InMemory {
	return &InMemory{values: prior.values}
}
