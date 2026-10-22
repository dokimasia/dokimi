// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package permutationtest contains the double and the suite of Mixed. InMemory
// is the control of the suite.
package permutationtest

import (
	"context"
	"errors"
	"slices"
	"strings"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/permutation"
)

// InMemory implements [permutation.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject permutation.Mixed new=NewInMemory
type InMemory struct {
	mu    sync.Mutex
	items []permutation.Value
}

// The compiler checks that InMemory implements [permutation.Mixed].
var _ permutation.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Add implements [permutation.Mixed].
func (s *InMemory) Add(ctx context.Context, v permutation.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, v)
	return nil
}

// Items implements [permutation.Mixed].
func (s *InMemory) Items(ctx context.Context) ([]permutation.Value, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	out := slices.Clone(s.items)
	slices.SortFunc(out, func(a, b permutation.Value) int {
		if k := strings.Compare(a.Key, b.Key); k != 0 {
			return k
		}
		return strings.Compare(a.Body, b.Body)
	})
	return out, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("permutationtest: nil context")
	}
	return ctx.Err()
}
