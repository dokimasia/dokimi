// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package compositewritertest contains the double and the suite of
// CompositeWriter. InMemory is the control of the suite.
package compositewritertest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/compositewriter"
)

// InMemory implements [compositewriter.CompositeWriter] in memory. The
// corpus applies the suite of CompositeWriter to it as a control.
//
//dokimi:subject compositewriter.CompositeWriter new=NewInMemory
type InMemory struct {
	values map[string]compositewriter.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [compositewriter.CompositeWriter].
var _ compositewriter.CompositeWriter = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{values: map[string]compositewriter.Value{}}
}

// Set implements [compositewriter.CompositeWriter].
func (s *InMemory) Set(ctx context.Context, key string, v compositewriter.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[key] = v
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("compositewritertest: nil context")
	}
	return ctx.Err()
}
