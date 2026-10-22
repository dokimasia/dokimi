// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package batchreadertest contains the double and the suite of BatchReader.
// InMemory is the control of the suite.
package batchreadertest

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/batchreader"
)

// ErrNotFound is the error of GetAll for a key without a value. GetAll wraps
// it with the key.
var ErrNotFound = errors.New("batchreadertest: not found")

// InMemory implements [batchreader.BatchReader] in memory. The corpus
// applies the suite of BatchReader to it as a control.
//
//dokimi:subject batchreader.BatchReader new=NewInMemory
type InMemory struct {
	mu     sync.Mutex
	values map[string]batchreader.Value
}

// The compiler checks that InMemory implements [batchreader.BatchReader].
var _ batchreader.BatchReader = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{values: map[string]batchreader.Value{}}
}

// GetAll implements [batchreader.BatchReader].
func (s *InMemory) GetAll(ctx context.Context, keys ...string) ([]batchreader.Value, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make([]batchreader.Value, 0, len(keys))
	for _, k := range keys {
		v, ok := s.values[k]
		if !ok {
			return nil, fmt.Errorf("%w: %s", ErrNotFound, k)
		}
		out = append(out, v)
	}
	return out, nil
}

// Put stores v under its key. A test seeds the reader through it.
func (s *InMemory) Put(v batchreader.Value) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = v
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("batchreadertest: nil context")
	}
	return ctx.Err()
}
