// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package streamreflectsmutationstest contains the double and the suite of
// Mixed. InMemory is the subject of the suite.
package streamreflectsmutationstest

import (
	"context"
	"errors"
	"iter"
	"slices"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/streamreflectsmutations"
)

// InMemory implements [streamreflectsmutations.Mixed] in memory as a set of
// items. The corpus applies the suite of Mixed to it as a subject.
//
// Remove marks the entry of an item as removed and keeps its place, so the
// position of a running iteration keeps pointing at the next entry. Add adds
// an entry for an item without a live entry.
//
//dokimi:subject streamreflectsmutations.Mixed new=NewInMemory
type InMemory struct {
	entries []entry
	mu      sync.Mutex
}

// entry is an item of InMemory, and whether Remove removed it.
type entry struct {
	item    string
	removed bool
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
	if s.live(item) < 0 {
		s.entries = append(s.entries, entry{item: item})
	}
	return nil
}

// Remove implements [streamreflectsmutations.Mixed].
func (s *InMemory) Remove(ctx context.Context, item string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if i := s.live(item); i >= 0 {
		s.entries[i].removed = true
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
			if i >= len(s.entries) {
				s.mu.Unlock()
				return
			}
			e := s.entries[i]
			s.mu.Unlock()

			if e.removed {
				continue
			}
			if !yield(e.item, nil) {
				return
			}
		}
	}
}

// live returns the position of the live entry of item, and -1 for an item
// without one. The caller holds the lock.
func (s *InMemory) live(item string) int {
	return slices.IndexFunc(s.entries, func(e entry) bool { return !e.removed && e.item == item })
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("streamreflectsmutationstest: nil context")
	}
	return ctx.Err()
}
