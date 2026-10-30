// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package leaderelectiontest contains the double and the suite of Contract.
// InMemory is the subject of the suite.
package leaderelectiontest

import (
	"context"
	"errors"
	"sync"

	leaderelection "go.dokimi.dev/dokimi/corpus/go/iface/contract/leader-election"
)

// ErrHeld is the error of Campaign while another candidate leads.
var ErrHeld = errors.New("leaderelectiontest: another candidate leads")

// Registry is the group whose candidates contend for one leadership.
type Registry struct {
	leader *InMemory
	mu     sync.Mutex
}

// NewRegistry returns a registry without a leader.
func NewRegistry() *Registry { return &Registry{} }

// Candidate returns a new candidate that contends for the leadership of the
// registry.
func (r *Registry) Candidate() *InMemory { return &InMemory{registry: r} }

// InMemory implements [leaderelection.Contract] in memory. The corpus
// applies the suite of Contract to it as a subject.
//
//dokimi:subject leaderelection.Contract new=NewInMemory
type InMemory struct{ registry *Registry }

// The compiler checks that InMemory implements [leaderelection.Contract].
var _ leaderelection.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return NewRegistry().Candidate() }

// Campaign implements [leaderelection.Contract].
func (s *InMemory) Campaign(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()
	if s.registry.leader != nil && s.registry.leader != s {
		return ErrHeld
	}
	s.registry.leader = s
	return nil
}

// Resign implements [leaderelection.Contract].
func (s *InMemory) Resign(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()
	if s.registry.leader == s {
		s.registry.leader = nil
	}
	return nil
}

// IsLeader implements [leaderelection.Contract].
func (s *InMemory) IsLeader(ctx context.Context) bool {
	if ctx != nil && ctx.Err() != nil {
		return false
	}
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()
	return s.registry.leader == s
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("leaderelectiontest: nil context")
	}
	return ctx.Err()
}
