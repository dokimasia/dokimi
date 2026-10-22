// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package deleteremovestest contains the double and the suite of Mixed.
// InMemory is the control of the suite.
package deleteremovestest

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/deleteremoves"
)

// ErrNotFound is the error of a read of a key without a value. It wraps
// [deleteremoves.ErrGone].
var ErrNotFound = fmt.Errorf("deleteremovestest: %w", deleteremoves.ErrGone)

// InMemory implements [deleteremoves.Mixed] in memory. The corpus applies
// the suite of Mixed to it as a control.
//
//dokimi:subject deleteremoves.Mixed new=NewInMemory
type InMemory struct {
	mu    sync.Mutex
	items map[string]string
}

// The compiler checks that InMemory implements [deleteremoves.Mixed].
var _ deleteremoves.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{items: map[string]string{}} }

// Put implements [deleteremoves.Mixed].
func (s *InMemory) Put(ctx context.Context, key, value string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[key] = value
	return nil
}

// Delete implements [deleteremoves.Mixed].
func (s *InMemory) Delete(ctx context.Context, key string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.items, key)
	return nil
}

// Read implements [deleteremoves.Mixed].
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

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("deleteremovestest: nil context")
	}
	return ctx.Err()
}
