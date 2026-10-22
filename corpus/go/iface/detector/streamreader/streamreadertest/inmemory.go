// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package streamreadertest contains the double and the suite of StreamReader.
// InMemory is the control of the suite.
package streamreadertest

import (
	"context"
	"errors"
	"iter"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/streamreader"
)

// InMemory implements [streamreader.StreamReader] in memory. The corpus
// applies the suite of StreamReader to it as a control.
//
//dokimi:subject streamreader.StreamReader new=NewInMemory
type InMemory struct {
	values []streamreader.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [streamreader.StreamReader].
var _ streamreader.StreamReader = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Add implements [streamreader.StreamReader].
func (s *InMemory) Add(ctx context.Context, v streamreader.Value) error {
	if ctx == nil {
		return errors.New("streamreadertest: nil context")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values = append(s.values, v)
	return nil
}

// List implements [streamreader.StreamReader].
func (s *InMemory) List(ctx context.Context) iter.Seq2[streamreader.Value, error] {
	if ctx == nil || ctx.Err() != nil {
		err := context.Canceled
		if ctx != nil {
			err = ctx.Err()
		}
		return func(yield func(streamreader.Value, error) bool) {
			yield(streamreader.Value{}, err)
		}
	}

	s.mu.Lock()
	snapshot := make([]streamreader.Value, len(s.values))
	copy(snapshot, s.values)
	s.mu.Unlock()

	return func(yield func(streamreader.Value, error) bool) {
		for _, v := range snapshot {
			if !yield(v, nil) {
				return
			}
		}
	}
}

// Put appends v to the stream. A test seeds the reader through it.
func (s *InMemory) Put(v streamreader.Value) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values = append(s.values, v)
}
