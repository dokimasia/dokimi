// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package circuitbreakertest contains the double and the suite of Contract.
// InMemory is the subject of the suite.
package circuitbreakertest

import (
	"context"
	"errors"
	"sync"

	circuitbreaker "go.dokimi.dev/dokimi/corpus/go/iface/contract/circuit-breaker"
)

// ErrOpen is the error of Run for a key whose breaker has tripped.
var ErrOpen = errors.New("circuitbreakertest: breaker is open")

// ErrDownstream is the error of the guarded call of UnwellKey.
var ErrDownstream = errors.New("circuitbreakertest: the downstream fails")

// UnwellKey is the key of a downstream whose every call fails.
const UnwellKey = "unwell"

// threshold is the number of consecutive failures of one key after which the
// breaker of NewInMemory trips.
const threshold = 3

// InMemory implements [circuitbreaker.Contract] in memory. The corpus applies
// the suite of Contract to it as a subject. It counts the consecutive failures
// of each key on its own.
//
//dokimi:subject circuitbreaker.Contract new=NewInMemory
type InMemory struct {
	failures  map[string]int
	threshold int
	mu        sync.Mutex
}

// The compiler checks that InMemory implements [circuitbreaker.Contract].
var _ circuitbreaker.Contract = (*InMemory)(nil)

// NewInMemory returns a closed breaker that trips after three consecutive
// failures of one key.
func NewInMemory() *InMemory { return NewInMemoryWith(threshold) }

// NewInMemoryWith returns a closed breaker that trips after threshold
// consecutive failures of one key.
func NewInMemoryWith(threshold int) *InMemory {
	return &InMemory{threshold: threshold, failures: map[string]int{}}
}

// Run implements [circuitbreaker.Contract]. A success resets the count of
// failures of its key.
func (s *InMemory) Run(ctx context.Context, key string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.failures[key] >= s.threshold {
		return ErrOpen
	}
	if key == UnwellKey {
		s.failures[key]++
		return ErrDownstream
	}
	s.failures[key] = 0
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("circuitbreakertest: nil context")
	}
	return ctx.Err()
}
