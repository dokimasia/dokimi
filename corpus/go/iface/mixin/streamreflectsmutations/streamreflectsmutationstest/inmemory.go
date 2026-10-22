// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package streamreflectsmutationstest contains the double and the suite of
// Mixed. InMemory is the control of the suite.
package streamreflectsmutationstest

import (
	"context"
	"errors"
	"iter"
	"slices"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/streamreflectsmutations"
)

// InMemory implements [streamreflectsmutations.Mixed] in memory. The corpus
// applies the suite of Mixed to it as a control.
//
//dokimi:subject streamreflectsmutations.Mixed new=NewInMemory
type InMemory struct {
	items []string
	mu    sync.Mutex
}

// The compiler checks that InMemory implements [streamreflectsmutations.Mixed].
var _ streamreflectsmutations.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Add implements [streamreflectsmutations.Mixed].
func (s *InMemory) Add(ctx context.Context, item string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.items = append(s.items, item)
	return nil
}

// Remove implements [streamreflectsmutations.Mixed].
func (s *InMemory) Remove(ctx context.Context, item string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if i := slices.Index(s.items, item); i >= 0 {
		s.items = slices.Delete(s.items, i, i+1)
	}
	return nil
}

// Stream implements [streamreflectsmutations.Mixed].
func (s *InMemory) Stream(ctx context.Context) iter.Seq2[string, error] {
	return func(yield func(string, error) bool) {
		if err := contextErr(ctx); err != nil {
			yield("", err)
			return
		}
		for i := 0; ; i++ {
			s.mu.Lock()
			if i >= len(s.items) {
				s.mu.Unlock()
				return
			}
			item := s.items[i]
			s.mu.Unlock()

			if !yield(item, nil) {
				return
			}
		}
	}
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("streamreflectsmutationstest: nil context")
	}
	return ctx.Err()
}
