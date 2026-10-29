// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package windowedtest contains the double and the suite of Mixed. InMemory is
// the control of the suite.
package windowedtest

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/windowed"
)

// window is the span of time before now in which CountIn counts a record.
const window = time.Minute

// InMemory implements [windowed.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject windowed.Mixed new=NewInMemoryFromHandle
type InMemory struct {
	clk  assert.Clock
	seen map[string][]time.Time
	mu   sync.Mutex
}

// The compiler checks that InMemory implements [windowed.Mixed].
var _ windowed.Mixed = (*InMemory)(nil)

// NewInMemoryOn returns an InMemory that reads the time from clk.
func NewInMemoryOn(clk assert.Clock) *InMemory {
	return &InMemory{clk: clk, seen: map[string][]time.Time{}}
}

// NewInMemory returns an InMemory on a controlled clock at the Unix epoch.
func NewInMemory() *InMemory {
	return NewInMemoryOn(assert.NewControlled(time.Unix(0, 0).UTC()))
}

// NewInMemoryFromHandle returns an InMemory on the clock of the handle h, so
// a check that advances the clock moves the time of the subject.
func NewInMemoryFromHandle(h *MixedHandle) *InMemory { return NewInMemoryOn(h.Clock()) }

// Record implements [windowed.Mixed].
func (s *InMemory) Record(ctx context.Context, key string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seen[key] = append(s.seen[key], s.clk.Now())
	return nil
}

// CountIn implements [windowed.Mixed].
func (s *InMemory) CountIn(ctx context.Context, key string) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	cutoff := s.clk.Now().Add(-window)
	var n int
	for _, at := range s.seen[key] {
		if at.After(cutoff) {
			n++
		}
	}
	return n, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("windowedtest: nil context")
	}
	return ctx.Err()
}
