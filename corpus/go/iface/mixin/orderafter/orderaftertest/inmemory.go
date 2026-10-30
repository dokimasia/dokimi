// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package orderaftertest contains the double and the suite of Mixed. InMemory
// is the subject of the suite.
package orderaftertest

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/orderafter"
)

// ErrNotPrepared is the error of a commit before its prepare. It wraps
// [orderafter.ErrNotReady].
var ErrNotPrepared = fmt.Errorf("orderaftertest: commit before prepare: %w", orderafter.ErrNotReady)

// InMemory implements [orderafter.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a subject.
//
//dokimi:subject orderafter.Mixed new=NewInMemory
type InMemory struct {
	mu       sync.Mutex
	prepared bool
	commits  int
}

// The compiler checks that InMemory implements [orderafter.Mixed].
var _ orderafter.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Prepare implements [orderafter.Mixed].
func (s *InMemory) Prepare(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.prepared = true
	return nil
}

// Commit implements [orderafter.Mixed].
func (s *InMemory) Commit(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.prepared {
		return ErrNotPrepared
	}
	s.commits++
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("orderaftertest: nil context")
	}
	return ctx.Err()
}
