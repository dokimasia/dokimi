// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package timeawaretest contains the double and the suite of Mixed. InMemory is
// the control of the suite.
package timeawaretest

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/timeaware"
)

// ErrUnseen is the error of AgeOf for a key that no call of Touch touched.
var ErrUnseen = errors.New("timeawaretest: key never touched")

// InMemory implements [timeaware.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject timeaware.Mixed new=NewInMemoryFromHandle
type InMemory struct {
	clk  assert.Clock
	seen map[string]time.Time
	mu   sync.Mutex
}

// The compiler checks that InMemory implements [timeaware.Mixed].
var _ timeaware.Mixed = (*InMemory)(nil)

// NewInMemoryOn returns an InMemory that reads the time from clk.
func NewInMemoryOn(clk assert.Clock) *InMemory {
	return &InMemory{clk: clk, seen: map[string]time.Time{}}
}

// NewInMemory returns an InMemory on a controlled clock at the Unix epoch.
func NewInMemory() *InMemory {
	return NewInMemoryOn(assert.NewControlled(time.Unix(0, 0).UTC()))
}

// NewInMemoryFromHandle returns an InMemory on the clock of the handle h, so
// a check that advances the clock moves the time of the subject.
func NewInMemoryFromHandle(h *MixedHandle) *InMemory { return NewInMemoryOn(h.Clock()) }

// Touch implements [timeaware.Mixed].
func (s *InMemory) Touch(ctx context.Context, key string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.seen[key] = s.clk.Now()
	return nil
}

// AgeOf implements [timeaware.Mixed].
func (s *InMemory) AgeOf(ctx context.Context, key string) (int64, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	at, ok := s.seen[key]
	if !ok {
		return 0, ErrUnseen
	}
	return s.clk.Now().Sub(at).Nanoseconds(), nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("timeawaretest: nil context")
	}
	return ctx.Err()
}
