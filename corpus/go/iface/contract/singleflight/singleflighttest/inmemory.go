// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package singleflighttest contains the double and the suite of Contract.
// InMemory is the subject of the suite.
package singleflighttest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/singleflight"
)

// call is the computation of one key. done closes once val and err have the
// result that each caller of the key shares.
type call struct {
	done chan struct{}
	err  error
	val  string
}

// InMemory implements [singleflight.Contract] in memory. The corpus applies
// the suite of Contract to it as a subject. Each call of Run passes its
// computation, and InMemory calls no dependency of its own, so it has no
// failure of a dependency to report.
//
//dokimi:subject singleflight.Contract new=NewInMemory options=InMemoryOptions
//dokimi:excuse Contract.Run/failure-is-reported reason="an InMemory calls no dependency, because each call passes its computation"
//dokimi:excuse Contract.Flights/failure-is-reported reason="Flights reads a count in memory and calls no dependency"
type InMemory struct {
	calls   map[string]*call
	flights int
	mu      sync.Mutex
}

// The compiler checks that InMemory implements [singleflight.Contract].
var _ singleflight.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{calls: map[string]*call{}} }

// InMemoryOptions returns the door of an InMemory for the checks whose
// dependencies fail on demand.
func InMemoryOptions() []ContractOption {
	return []ContractOption{conform.Faultable(faulty)}
}

// faulty returns a new InMemory. An InMemory has no dependency of its own,
// because it computes the result of a key through the computation of the
// call.
func faulty(*ContractHandle, *conform.Injector) singleflight.Contract { return NewInMemory() }

// Run implements [singleflight.Contract].
func (s *InMemory) Run(ctx context.Context, key string, compute func() string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}

	s.mu.Lock()
	if c, seen := s.calls[key]; seen {
		s.mu.Unlock()
		<-c.done
		return c.val, nil
	}
	c := &call{done: make(chan struct{})}
	s.calls[key] = c
	s.flights++
	s.mu.Unlock()

	if compute != nil {
		c.val = compute()
	}
	close(c.done)
	return c.val, nil
}

// Flights implements [singleflight.Contract].
func (s *InMemory) Flights(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.flights, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("singleflighttest: nil context")
	}
	return ctx.Err()
}
