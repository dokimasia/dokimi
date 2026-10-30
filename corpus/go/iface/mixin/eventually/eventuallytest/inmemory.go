// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package eventuallytest contains the double and the suite of Mixed. InMemory
// is the subject of the suite.
package eventuallytest

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"slices"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/eventually"
)

// InMemory implements [eventually.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a subject.
//
//dokimi:subject eventually.Mixed new=NewInMemory
type InMemory struct {
	settled map[string]bool
	pending []string
	mu      sync.Mutex
}

// The compiler checks that InMemory implements [eventually.Mixed].
var _ eventually.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{settled: map[string]bool{}} }

// Publish implements [eventually.Mixed].
func (s *InMemory) Publish(ctx context.Context, item string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pending = append(s.pending, item)
	return nil
}

// Settle implements [eventually.Mixed].
func (s *InMemory) Settle(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range s.pending {
		s.settled[item] = true
	}
	s.pending = nil
	return nil
}

// Sync implements [eventually.Mixed].
func (s *InMemory) Sync(ctx context.Context, peer eventually.Mixed) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if peer == nil {
		return errors.New("eventuallytest: nothing to sync with")
	}
	items, err := peer.Items(ctx)
	if err != nil {
		return fmt.Errorf("eventuallytest: read the items of the peer: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range items {
		s.settled[item] = true
	}
	return nil
}

// Items implements [eventually.Mixed].
func (s *InMemory) Items(ctx context.Context) ([]string, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Sorted(maps.Keys(s.settled)), nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("eventuallytest: nil context")
	}
	return ctx.Err()
}
