// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package restrictedpooltest contains the double and the suite of Store.
// InMemory is the control of the suite.
package restrictedpooltest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/restrictedpool"
)

// InMemory implements [restrictedpool.Store] in memory. The corpus applies
// the suite of Store to it as a control.
//
//dokimi:subject restrictedpool.Store new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	bodies map[restrictedpool.Key]restrictedpool.Body
}

// The compiler checks that InMemory implements [restrictedpool.Store].
var _ restrictedpool.Store = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{bodies: map[restrictedpool.Key]restrictedpool.Body{}}
}

// Put implements [restrictedpool.Store].
func (s *InMemory) Put(ctx context.Context, key restrictedpool.Key, body restrictedpool.Body) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.bodies[key] = body
	return nil
}

// Get implements [restrictedpool.Store].
func (s *InMemory) Get(ctx context.Context, key restrictedpool.Key) (restrictedpool.Body, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	body, ok := s.bodies[key]
	if !ok {
		return "", restrictedpool.ErrNotFound
	}
	return body, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("restrictedpooltest: nil context")
	}
	return ctx.Err()
}
