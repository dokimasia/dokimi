// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package crdtmergetest contains the double and the suite of Mixed. InMemory is
// the control of the suite.
package crdtmergetest

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/crdtmerge"
)

// InMemory implements [crdtmerge.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject crdtmerge.Mixed new=NewInMemory
type InMemory struct {
	mu    sync.Mutex
	items map[string]struct{}
}

// The compiler checks that InMemory implements [crdtmerge.Mixed] and
// [crdtmerge.Replica].
var (
	_ crdtmerge.Mixed   = (*InMemory)(nil)
	_ crdtmerge.Replica = (*InMemory)(nil)
)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{items: map[string]struct{}{}} }

// Add implements [crdtmerge.Mixed].
func (s *InMemory) Add(ctx context.Context, item string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items[item] = struct{}{}
	return nil
}

// Merge implements [crdtmerge.Mixed].
func (s *InMemory) Merge(ctx context.Context, peer crdtmerge.Replica) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	if peer == nil {
		return nil
	}
	theirs, err := peer.Items(ctx)
	if err != nil {
		return fmt.Errorf("crdtmergetest: read the items of the peer: %w", err)
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, item := range theirs {
		s.items[item] = struct{}{}
	}
	return nil
}

// Items implements [crdtmerge.Mixed].
func (s *InMemory) Items(ctx context.Context) ([]string, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]string, 0, len(s.items))
	for item := range s.items {
		out = append(out, item)
	}
	return out, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("crdtmergetest: nil context")
	}
	return ctx.Err()
}
