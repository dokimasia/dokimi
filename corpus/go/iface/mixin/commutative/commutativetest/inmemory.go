// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package commutativetest contains the double and the suite of Mixed. InMemory
// is the subject of the suite.
package commutativetest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/commutative"
)

// InMemory implements [commutative.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a subject.
//
//dokimi:subject commutative.Mixed new=NewInMemory
type InMemory struct {
	mu    sync.Mutex
	total int
}

// The compiler checks that InMemory implements [commutative.Mixed].
var _ commutative.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Apply implements [commutative.Mixed].
func (s *InMemory) Apply(ctx context.Context, d commutative.Delta) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.total += d.Amount
	return nil
}

// Total implements [commutative.Mixed].
func (s *InMemory) Total(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.total, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("commutativetest: nil context")
	}
	return ctx.Err()
}
