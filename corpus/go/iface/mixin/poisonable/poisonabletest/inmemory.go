// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package poisonabletest contains the double and the suite of Mixed. InMemory
// is the subject of the suite.
package poisonabletest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/poisonable"
)

// ErrPoisoned is the error of Probe after a call of Fail.
var ErrPoisoned = errors.New("poisonabletest: poisoned")

// InMemory implements [poisonable.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a subject.
//
//dokimi:subject poisonable.Mixed new=NewInMemory
type InMemory struct {
	mu       sync.Mutex
	poisoned bool
}

// The compiler checks that InMemory implements [poisonable.Mixed].
var _ poisonable.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Fail implements [poisonable.Mixed].
func (s *InMemory) Fail(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.poisoned = true
	return nil
}

// Probe implements [poisonable.Mixed].
func (s *InMemory) Probe() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.poisoned {
		return ErrPoisoned
	}
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("poisonabletest: nil context")
	}
	return ctx.Err()
}
