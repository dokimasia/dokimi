// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package readernoerrortest contains the double and the suite of ReaderNoError.
// InMemory is the control of the suite.
package readernoerrortest

import (
	"context"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/readernoerror"
)

// InMemory implements [readernoerror.ReaderNoError] in memory. The corpus
// applies the suite of ReaderNoError to it as a control.
//
//dokimi:subject readernoerror.ReaderNoError new=NewInMemory
type InMemory struct {
	values map[string]readernoerror.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [readernoerror.ReaderNoError].
var _ readernoerror.ReaderNoError = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{values: map[string]readernoerror.Value{}}
}

// Lookup implements [readernoerror.ReaderNoError].
func (s *InMemory) Lookup(ctx context.Context, key string) readernoerror.Value {
	if ctx == nil || ctx.Err() != nil {
		return readernoerror.Value{}
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.values[key]
}

// Put stores v under its key. A test seeds the reader through it.
func (s *InMemory) Put(v readernoerror.Value) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = v
}
