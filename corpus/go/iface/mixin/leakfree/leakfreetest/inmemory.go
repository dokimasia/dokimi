// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package leakfreetest contains the double and the suite of Mixed. InMemory is
// the subject of the suite.
package leakfreetest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/leakfree"
)

// ErrNotAcquired is the error of Release without an outstanding acquisition.
var ErrNotAcquired = errors.New("leakfreetest: nothing to release")

// InMemory implements [leakfree.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a subject.
//
//dokimi:subject leakfree.Mixed new=NewInMemory
type InMemory struct {
	acquired []*acquisition
	mu       sync.Mutex
}

// acquisition is one call of Acquire with the goroutine that it starts.
// Release closes release and waits until the goroutine closes exited.
type acquisition struct {
	release chan struct{}
	exited  chan struct{}
}

// The compiler checks that InMemory implements [leakfree.Mixed].
var _ leakfree.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Acquire implements [leakfree.Mixed].
func (s *InMemory) Acquire(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	h := &acquisition{release: make(chan struct{}), exited: make(chan struct{})}
	go func() {
		<-h.release
		close(h.exited)
	}()
	s.acquired = append(s.acquired, h)
	return nil
}

// Release implements [leakfree.Mixed].
func (s *InMemory) Release(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.acquired) == 0 {
		return ErrNotAcquired
	}
	last := len(s.acquired) - 1
	h := s.acquired[last]
	s.acquired = s.acquired[:last]
	close(h.release)
	<-h.exited
	return nil
}

// Outstanding implements [leakfree.Mixed].
func (s *InMemory) Outstanding(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.acquired), nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("leakfreetest: nil context")
	}
	return ctx.Err()
}
