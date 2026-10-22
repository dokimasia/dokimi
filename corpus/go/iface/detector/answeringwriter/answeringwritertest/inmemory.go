// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package answeringwritertest contains the double and the suite of
// AnsweringWriter. InMemory is the control of the suite.
package answeringwritertest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/answeringwriter"
)

// InMemory implements [answeringwriter.AnsweringWriter] in memory. The
// corpus applies the suite of AnsweringWriter to it as a control.
//
//dokimi:subject answeringwriter.AnsweringWriter new=NewInMemory
type InMemory struct {
	values map[string]answeringwriter.Value
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [answeringwriter.AnsweringWriter].
var _ answeringwriter.AnsweringWriter = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{values: map[string]answeringwriter.Value{}}
}

// Put implements [answeringwriter.AnsweringWriter].
func (s *InMemory) Put(ctx context.Context, v answeringwriter.Value) (answeringwriter.Value, error) {
	if err := contextErr(ctx); err != nil {
		return answeringwriter.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.values[v.Key] = v
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("answeringwritertest: nil context")
	}
	return ctx.Err()
}
