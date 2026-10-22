// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package readerwithbooltest contains the double and the suite of
// ReaderWithBool. InMemory is the control of the suite.
package readerwithbooltest

import (
	"context"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/readerwithbool"
)

// InMemory implements [readerwithbool.ReaderWithBool] in memory. The corpus
// applies the suite of ReaderWithBool to it as a control.
//
//dokimi:subject readerwithbool.ReaderWithBool new=NewInMemory
type InMemory struct {
	values map[string]readerwithbool.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [readerwithbool.ReaderWithBool].
var _ readerwithbool.ReaderWithBool = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{values: map[string]readerwithbool.Value{}}
}

// Load implements [readerwithbool.ReaderWithBool].
func (s *InMemory) Load(ctx context.Context, key string) (readerwithbool.Value, bool) {
	if ctx == nil || ctx.Err() != nil {
		return readerwithbool.Value{}, false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, ok := s.values[key]
	if !ok {
		return readerwithbool.Value{}, false
	}
	return v, true
}

// Put stores v under its key. A test seeds the reader through it.
func (s *InMemory) Put(v readerwithbool.Value) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = v
}
