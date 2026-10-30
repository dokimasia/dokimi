// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package writesfollowreadstest contains the double and the suite of Mixed.
// InMemory is the subject of the suite.
package writesfollowreadstest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/writesfollowreads"
)

// ErrNotFound is the error of Get for a key that no call of Store wrote.
var ErrNotFound = errors.New("writesfollowreadstest: not found")

// InMemory implements [writesfollowreads.Mixed] in memory. The corpus
// applies the suite of Mixed to it as a subject.
//
//dokimi:subject writesfollowreads.Mixed new=NewInMemory
type InMemory struct {
	revs   map[string]int64
	values map[string]writesfollowreads.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [writesfollowreads.Mixed].
var _ writesfollowreads.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{
		revs:   map[string]int64{},
		values: map[string]writesfollowreads.Value{},
	}
}

// Store implements [writesfollowreads.Mixed].
func (s *InMemory) Store(ctx context.Context, v writesfollowreads.Value) (writesfollowreads.Value, error) {
	if err := contextErr(ctx); err != nil {
		return writesfollowreads.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.revs[v.Key]++
	v.Rev = s.revs[v.Key]
	s.values[v.Key] = v
	return v, nil
}

// Get implements [writesfollowreads.Mixed].
func (s *InMemory) Get(ctx context.Context, key string) (writesfollowreads.Value, error) {
	if err := contextErr(ctx); err != nil {
		return writesfollowreads.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.values[key]
	if !ok {
		return writesfollowreads.Value{}, ErrNotFound
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("writesfollowreadstest: nil context")
	}
	return ctx.Err()
}
