// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package overmatchtest contains the double and the suite of Mixed. InMemory is
// the subject of the suite.
package overmatchtest

import (
	"context"
	"errors"
	"slices"
	"strings"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/overmatch"
)

// InMemory implements [overmatch.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a subject.
//
//dokimi:subject overmatch.Mixed new=NewInMemory
type InMemory struct {
	items []overmatch.Value
	mu    sync.Mutex
}

// The compiler checks that InMemory implements [overmatch.Mixed].
var _ overmatch.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Add implements [overmatch.Mixed].
func (s *InMemory) Add(ctx context.Context, v overmatch.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, v)
	return nil
}

// Items implements [overmatch.Mixed].
func (s *InMemory) Items(ctx context.Context) ([]overmatch.Value, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	out := slices.Clone(s.items)
	slices.SortFunc(out, func(a, b overmatch.Value) int {
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
		return errors.New("overmatchtest: nil context")
	}
	return ctx.Err()
}
