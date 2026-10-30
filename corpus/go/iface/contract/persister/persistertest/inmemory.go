// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package persistertest contains the double and the suite of Contract, and the
// builder of Value. InMemory is the subject of the suite.
package persistertest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/persister"
)

// ErrNotFound is the error of Get for a key without a value.
var ErrNotFound = errors.New("persistertest: no value under that key")

// InMemory implements [persister.Contract] in memory. The corpus applies the
// suite of Contract to it as a subject.
//
//dokimi:subject persister.Contract new=NewInMemory
type InMemory struct {
	values map[string]persister.Value
	mu     sync.Mutex
	// gone reports that the medium failed. Put then refuses each write and
	// records nothing.
	gone bool
}

// The compiler checks that InMemory implements [persister.Contract].
var _ persister.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{values: map[string]persister.Value{}} }

// Put implements [persister.Contract].
func (s *InMemory) Put(ctx context.Context, v persister.Value) (persister.Value, error) {
	if err := contextErr(ctx); err != nil {
		return persister.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.gone {
		return persister.Value{}, persister.ErrMediumGone
	}
	s.values[v.Key] = v
	return v, nil
}

// LoseTheMedium fails the medium, so each later Put returns
// [persister.ErrMediumGone].
func (s *InMemory) LoseTheMedium() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.gone = true
}

// Get implements [persister.Contract].
func (s *InMemory) Get(ctx context.Context, key string) (persister.Value, error) {
	if err := contextErr(ctx); err != nil {
		return persister.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, present := s.values[key]
	if !present {
		return persister.Value{}, ErrNotFound
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("persistertest: nil context")
	}
	return ctx.Err()
}

// Reopen returns an InMemory with the values of prior and a working medium. A
// restart after a crash keeps every value that Put wrote. Do not use prior
// after the call. The two stores share the values but not the lock.
func Reopen(prior *InMemory) *InMemory {
	return &InMemory{values: prior.values}
}
