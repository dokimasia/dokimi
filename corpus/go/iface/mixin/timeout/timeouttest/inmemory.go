// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package timeouttest contains the double and the suite of Mixed. InMemory is
// the control of the suite.
package timeouttest

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/timeout"
)

// InMemory implements [timeout.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject timeout.Mixed new=NewInMemory
type InMemory struct {
	clk   assert.Clock
	delay time.Duration
	mu    sync.Mutex
}

// The compiler checks that InMemory implements [timeout.Mixed].
var _ timeout.Mixed = (*InMemory)(nil)

// NewInMemoryOn returns an InMemory on clk whose Slow returns at once.
func NewInMemoryOn(clk assert.Clock) *InMemory {
	return &InMemory{clk: clk}
}

// NewInMemory returns an InMemory on the system clock whose Slow returns at
// once.
func NewInMemory() *InMemory { return NewInMemoryOn(assert.System{}) }

// WithDelay returns an InMemory whose Slow waits d on clk, or until its
// context ends.
func WithDelay(clk assert.Clock, d time.Duration) *InMemory {
	return &InMemory{clk: clk, delay: d}
}

// Slow implements [timeout.Mixed]. A wait that its context ends leaves its
// sleeper on the clock until the clock passes the delay.
func (s *InMemory) Slow(ctx context.Context, _ string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	clk, d := s.clk, s.delay
	s.mu.Unlock()
	if d <= 0 {
		return nil
	}
	slept := make(chan struct{})
	go func() {
		clk.Sleep(d)
		close(slept)
	}()
	select {
	case <-slept:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("timeouttest: nil context")
	}
	return ctx.Err()
}
