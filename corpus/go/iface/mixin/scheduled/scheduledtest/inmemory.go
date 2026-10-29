// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package scheduledtest contains the double and the suite of Mixed. InMemory is
// the control of the suite.
package scheduledtest

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/scheduled"
)

// InMemory implements [scheduled.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject scheduled.Mixed new=NewInMemoryFromHandle
type InMemory struct {
	clk   assert.Clock
	tasks []time.Time
	mu    sync.Mutex
}

// The compiler checks that InMemory implements [scheduled.Mixed].
var _ scheduled.Mixed = (*InMemory)(nil)

// NewInMemoryOn returns an InMemory that reads the time from clk.
func NewInMemoryOn(clk assert.Clock) *InMemory {
	return &InMemory{clk: clk}
}

// NewInMemory returns an InMemory on a controlled clock at the Unix epoch.
func NewInMemory() *InMemory {
	return NewInMemoryOn(assert.NewControlled(time.Unix(0, 0).UTC()))
}

// NewInMemoryFromHandle returns an InMemory on the clock of the handle h, so
// a check that advances the clock runs the tasks that are due.
func NewInMemoryFromHandle(h *MixedHandle) *InMemory { return NewInMemoryOn(h.Clock()) }

// At implements [scheduled.Mixed].
func (s *InMemory) At(ctx context.Context, after time.Duration) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = append(s.tasks, s.clk.Now().Add(after))
	return nil
}

// Fired implements [scheduled.Mixed].
func (s *InMemory) Fired(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	now := s.clk.Now()
	var n int
	for _, at := range s.tasks {
		if !at.After(now) {
			n++
		}
	}
	return n, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("scheduledtest: nil context")
	}
	return ctx.Err()
}
