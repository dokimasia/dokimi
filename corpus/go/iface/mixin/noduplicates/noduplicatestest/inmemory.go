// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package noduplicatestest contains the double and the suite of Mixed. InMemory
// is the control of the suite.
package noduplicatestest

import (
	"context"
	"errors"
	"slices"
	"strings"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/noduplicates"
)

// InMemory implements [noduplicates.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject noduplicates.Mixed new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	values map[string]noduplicates.Value
}

// The compiler checks that InMemory implements [noduplicates.Mixed].
var _ noduplicates.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{values: map[string]noduplicates.Value{}}
}

// Add implements [noduplicates.Mixed].
func (s *InMemory) Add(ctx context.Context, v noduplicates.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = v
	return nil
}

// Items implements [noduplicates.Mixed].
func (s *InMemory) Items(ctx context.Context) ([]noduplicates.Value, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]noduplicates.Value, 0, len(s.values))
	for _, v := range s.values {
		out = append(out, v)
	}
	slices.SortFunc(out, func(a, b noduplicates.Value) int {
		return strings.Compare(a.Key, b.Key)
	})
	return out, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("noduplicatestest: nil context")
	}
	return ctx.Err()
}
