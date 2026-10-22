// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package writertest contains the double and the suite of Writer. InMemory is
// the control of the suite.
package writertest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/writer"
)

// InMemory implements [writer.Writer] in memory. The corpus applies the
// suite of Writer to it as a control.
//
//dokimi:subject writer.Writer new=NewInMemory
type InMemory struct {
	values map[string]writer.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [writer.Writer].
var _ writer.Writer = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{values: map[string]writer.Value{}} }

// Put implements [writer.Writer].
func (s *InMemory) Put(ctx context.Context, v writer.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = v
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("writertest: nil context")
	}
	return ctx.Err()
}
