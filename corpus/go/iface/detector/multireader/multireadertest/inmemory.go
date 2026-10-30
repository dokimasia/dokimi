// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package multireadertest contains the double and the suite of MultiReader.
// InMemory is the subject of the suite.
package multireadertest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/multireader"
)

// InMemory implements [multireader.MultiReader] in memory. The corpus
// applies the suite of MultiReader to it as a subject.
//
//dokimi:subject multireader.MultiReader new=NewInMemory
type InMemory struct {
	values map[string]multireader.Value
	metas  map[string]multireader.Meta
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [multireader.MultiReader].
var _ multireader.MultiReader = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{
		values: map[string]multireader.Value{},
		metas:  map[string]multireader.Meta{},
	}
}

// GetWithMeta implements [multireader.MultiReader].
func (s *InMemory) GetWithMeta(
	ctx context.Context, key string,
) (multireader.Value, multireader.Meta, error) {
	if err := contextErr(ctx); err != nil {
		return multireader.Value{}, multireader.Meta{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.values[key]
	if !ok {
		return multireader.Value{}, multireader.Meta{}, multireader.ErrNotFound
	}
	return v, s.metas[key], nil
}

// Put stores v and its metadata m under the key of v. A test seeds the reader
// through it.
func (s *InMemory) Put(v multireader.Value, m multireader.Meta) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = v
	s.metas[v.Key] = m
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("multireadertest: nil context")
	}
	return ctx.Err()
}
