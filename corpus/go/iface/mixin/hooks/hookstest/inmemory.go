// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package hookstest contains the double and the suite of Mixed. InMemory is the
// control of the suite.
package hookstest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/hooks"
)

// InMemory implements [hooks.Mixed] in memory. The corpus applies the suite
// of Mixed to it as a control.
//
//dokimi:subject hooks.Mixed new=NewInMemory
type InMemory struct {
	mu       sync.Mutex
	handlers []func(event string)
}

// The compiler checks that InMemory implements [hooks.Mixed].
var _ hooks.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// OnEvent implements [hooks.Mixed].
func (s *InMemory) OnEvent(fn func(event string)) {
	if fn == nil {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.handlers = append(s.handlers, fn)
}

// Fire implements [hooks.Mixed].
func (s *InMemory) Fire(ctx context.Context, event string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	snapshot := make([]func(event string), len(s.handlers))
	copy(snapshot, s.handlers)
	s.mu.Unlock()

	for _, fn := range snapshot {
		fn(event)
	}
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("hookstest: nil context")
	}
	return ctx.Err()
}
