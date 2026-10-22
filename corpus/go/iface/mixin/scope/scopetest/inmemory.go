// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package scopetest contains the double and the suite of Mixed. InMemory is the
// control of the suite.
package scopetest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/scope"
)

// ErrNotFound is the error of a read of a key without a value.
var ErrNotFound = errors.New("scopetest: not found")

// InMemory implements [scope.Mixed] in memory. The corpus applies the suite
// of Mixed to it as a control.
//
//dokimi:subject scope.Mixed new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	scopes map[string]map[string]string
}

// The compiler checks that InMemory implements [scope.Mixed].
var _ scope.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{scopes: map[string]map[string]string{}}
}

// Set implements [scope.Mixed].
func (s *InMemory) Set(ctx context.Context, sc, key, value string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.scopes[sc] == nil {
		s.scopes[sc] = map[string]string{}
	}
	s.scopes[sc][key] = value
	return nil
}

// Get implements [scope.Mixed].
func (s *InMemory) Get(ctx context.Context, sc, key string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.scopes[sc][key]
	if !ok {
		return "", ErrNotFound
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("scopetest: nil context")
	}
	return ctx.Err()
}
