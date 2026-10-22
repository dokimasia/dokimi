// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package pointerreadertest contains the double and the suite of PointerReader.
// InMemory is the control of the suite.
package pointerreadertest

import (
	"context"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/pointerreader"
)

// InMemory implements [pointerreader.PointerReader] in memory. The corpus
// applies the suite of PointerReader to it as a control.
//
//dokimi:subject pointerreader.PointerReader new=NewInMemory
type InMemory struct {
	values map[string]pointerreader.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [pointerreader.PointerReader].
var _ pointerreader.PointerReader = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{values: map[string]pointerreader.Value{}}
}

// Find implements [pointerreader.PointerReader].
func (s *InMemory) Find(ctx context.Context, key string) *pointerreader.Value {
	if ctx == nil || ctx.Err() != nil {
		return nil
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.values[key]
	if !ok {
		return nil
	}
	return &v
}

// Put stores v under its key. A test seeds the reader through it.
func (s *InMemory) Put(v pointerreader.Value) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = v
}
