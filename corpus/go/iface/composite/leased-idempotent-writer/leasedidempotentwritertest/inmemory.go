// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package leasedidempotentwritertest contains the double and the suite of
// LeasedWriter. InMemory is the control of the suite.
package leasedidempotentwritertest

import (
	"context"
	"errors"
	"sync"

	leasedidempotentwriter "go.dokimi.dev/dokimi/corpus/go/iface/composite/leased-idempotent-writer"
)

// Registry is the set of leases that the holders of one registry contend
// for. A holder is refused a lease only while another holder of the registry
// has it.
type Registry struct {
	mu     sync.Mutex
	leases map[string]*InMemory
}

// NewRegistry returns a registry without leases.
func NewRegistry() *Registry { return &Registry{leases: map[string]*InMemory{}} }

// Holder returns a new holder that contends for the leases of the registry.
func (r *Registry) Holder() *InMemory { return &InMemory{registry: r} }

// InMemory implements [leasedidempotentwriter.LeasedWriter] in memory. The
// corpus applies the suite of LeasedWriter to it as a control.
//
//dokimi:subject leasedidempotentwriter.LeasedWriter new=NewInMemory
type InMemory struct{ registry *Registry }

// The compiler checks that InMemory implements [leasedidempotentwriter.LeasedWriter].
var _ leasedidempotentwriter.LeasedWriter = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return NewRegistry().Holder() }

// Acquire implements [leasedidempotentwriter.LeasedWriter].
func (s *InMemory) Acquire(ctx context.Context, key string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()
	if holder, taken := s.registry.leases[key]; taken && holder != s {
		return leasedidempotentwriter.ErrHeld
	}
	s.registry.leases[key] = s
	return nil
}

// Release implements [leasedidempotentwriter.LeasedWriter].
func (s *InMemory) Release(ctx context.Context, key string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.registry.mu.Lock()
	defer s.registry.mu.Unlock()
	if holder, taken := s.registry.leases[key]; taken && holder == s {
		delete(s.registry.leases, key)
	}
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("leasedidempotentwritertest: nil context")
	}
	return ctx.Err()
}
