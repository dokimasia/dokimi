// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package leasetest contains the double and the suite of Contract. InMemory is
// the subject of the suite.
package leasetest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/lease"
)

// InMemory implements [lease.Contract] in memory. The corpus applies the
// suite of Contract to it as a subject.
//
//dokimi:subject lease.Contract new=NewInMemory
type InMemory struct {
	leased map[string]bool
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [lease.Contract].
var _ lease.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{leased: map[string]bool{}} }

// Acquire implements [lease.Contract].
func (s *InMemory) Acquire(ctx context.Context, key string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.leased[key] {
		return lease.ErrHeld
	}
	s.leased[key] = true

	go func() {
		<-ctx.Done()
		s.mu.Lock()
		defer s.mu.Unlock()
		delete(s.leased, key)
	}()
	return nil
}

// Release implements [lease.Contract].
func (s *InMemory) Release(ctx context.Context, key string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.leased, key)
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("leasetest: nil context")
	}
	return ctx.Err()
}
