// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package stableordertest contains the double and the suite of Mixed. InMemory
// is the control of the suite.
package stableordertest

import (
	"context"
	"errors"
	"slices"
	"strings"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/stableorder"
)

// InMemory implements [stableorder.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject stableorder.Mixed new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	values []stableorder.Value
}

// The compiler checks that InMemory implements [stableorder.Mixed].
var _ stableorder.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{}
}

// Add implements [stableorder.Mixed].
func (s *InMemory) Add(ctx context.Context, v stableorder.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values = append(s.values, v)
	return nil
}

// Items implements [stableorder.Mixed].
func (s *InMemory) Items(ctx context.Context) ([]stableorder.Value, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]stableorder.Value, len(s.values))
	copy(out, s.values)
	slices.SortStableFunc(out, func(a, b stableorder.Value) int {
		return strings.Compare(a.Key, b.Key)
	})
	return out, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("stableordertest: nil context")
	}
	return ctx.Err()
}
