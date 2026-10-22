// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package upsertertest contains the double and the suite of Contract. InMemory
// is the control of the suite.
package upsertertest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/upserter"
)

// ErrNotFound is the error of Get for a key without a value.
var ErrNotFound = errors.New("upsertertest: no value under that key")

// InMemory implements [upserter.Contract] in memory. The corpus applies the
// suite of Contract to it as a control.
//
//dokimi:subject upserter.Contract new=NewInMemory
type InMemory struct {
	values map[string]upserter.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [upserter.Contract].
var _ upserter.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{values: map[string]upserter.Value{}} }

// Put implements [upserter.Contract].
func (s *InMemory) Put(ctx context.Context, v upserter.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = v
	return nil
}

// Get implements [upserter.Contract].
func (s *InMemory) Get(ctx context.Context, key string) (upserter.Value, error) {
	if err := contextErr(ctx); err != nil {
		return upserter.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, present := s.values[key]
	if !present {
		return upserter.Value{}, ErrNotFound
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("upsertertest: nil context")
	}
	return ctx.Err()
}

// Reopen returns an InMemory with the values of prior. A restart after a crash
// keeps every value that Put wrote. Do not use prior after the call. The two
// stores share the values but not the lock.
func Reopen(prior *InMemory) *InMemory {
	return &InMemory{values: prior.values}
}
