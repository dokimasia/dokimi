// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package castest contains the double and the suite of Contract. InMemory is
// the control of the suite.
package castest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/cas"
)

// ErrEmpty is the error of Get for a cell that no Put wrote.
var ErrEmpty = errors.New("castest: no write set the cell yet")

// InMemory implements [cas.Contract] in memory. The corpus applies the suite
// of Contract to it as a control.
//
//dokimi:subject cas.Contract new=NewInMemory
type InMemory struct {
	value   cas.Value
	current int64
	mu      sync.Mutex
	present bool
}

// The compiler checks that InMemory implements [cas.Contract].
var _ cas.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Put implements [cas.Contract].
func (s *InMemory) Put(ctx context.Context, v cas.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if v.Version != s.current {
		return cas.ErrMismatch
	}
	s.value = v
	s.current++
	s.present = true
	return nil
}

// Get implements [cas.Contract].
func (s *InMemory) Get(ctx context.Context) (cas.Value, error) {
	if err := contextErr(ctx); err != nil {
		return cas.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.present {
		return cas.Value{}, ErrEmpty
	}
	return s.value, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("castest: nil context")
	}
	return ctx.Err()
}
