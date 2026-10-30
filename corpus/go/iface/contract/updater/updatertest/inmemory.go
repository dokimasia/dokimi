// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package updatertest contains the double and the suite of Contract. InMemory
// is the subject of the suite.
package updatertest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/updater"
)

// ErrNotFound is the error of Get for a key without a value.
var ErrNotFound = errors.New("updatertest: no value under that key")

// InMemory implements [updater.Contract] in memory. The corpus applies the
// suite of Contract to it as a subject.
//
//dokimi:subject updater.Contract new=NewInMemory
type InMemory struct {
	values map[string]updater.Value
	writes map[string]int
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [updater.Contract].
var _ updater.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{
		values: map[string]updater.Value{},
		writes: map[string]int{},
	}
}

// Put implements [updater.Contract].
func (s *InMemory) Put(ctx context.Context, v updater.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = v
	s.writes[v.Key]++
	return nil
}

// Get implements [updater.Contract].
func (s *InMemory) Get(ctx context.Context, key string) (updater.Value, error) {
	if err := contextErr(ctx); err != nil {
		return updater.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, present := s.values[key]
	if !present {
		return updater.Value{}, ErrNotFound
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("updatertest: nil context")
	}
	return ctx.Err()
}

// Reopen returns an InMemory with the values and the write counts of prior. A
// restart after a crash keeps both. Do not use prior after the call. The two
// stores share the values but not the lock.
func Reopen(prior *InMemory) *InMemory {
	return &InMemory{values: prior.values, writes: prior.writes}
}
