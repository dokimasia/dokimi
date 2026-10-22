// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package ttlperwritetest contains the double and the suite of Mixed. InMemory
// is the control of the suite.
package ttlperwritetest

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/ttlperwrite"
)

// InMemory implements [ttlperwrite.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject ttlperwrite.Mixed new=NewInMemory
type InMemory struct {
	clk     assert.Clock
	entries map[string]stored
	mu      sync.Mutex
}

// stored is an entry with the instant at which it expires. An entry without a
// lifetime has the zero instant and does not expire.
type stored struct {
	until time.Time
	entry ttlperwrite.Entry
}

// The compiler checks that InMemory implements [ttlperwrite.Mixed].
var _ ttlperwrite.Mixed = (*InMemory)(nil)

// NewInMemory returns an InMemory on a controlled clock at the Unix epoch.
func NewInMemory() *InMemory { return NewInMemoryOn(assert.NewControlled(time.Unix(0, 0).UTC())) }

// NewInMemoryOn returns an InMemory that reads the time from clk.
func NewInMemoryOn(clk assert.Clock) *InMemory {
	return &InMemory{clk: clk, entries: map[string]stored{}}
}

// Put implements [ttlperwrite.Mixed].
func (s *InMemory) Put(ctx context.Context, e ttlperwrite.Entry) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	h := stored{entry: e}
	if e.Lifetime > 0 {
		h.until = s.clk.Now().Add(e.Lifetime)
	}
	s.entries[e.Key] = h
	return nil
}

// Read implements [ttlperwrite.Mixed].
func (s *InMemory) Read(ctx context.Context, key string) (ttlperwrite.Entry, error) {
	if err := contextErr(ctx); err != nil {
		return ttlperwrite.Entry{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	h, present := s.entries[key]
	switch {
	case !present:
		return ttlperwrite.Entry{}, ttlperwrite.ErrExpired
	case !h.until.IsZero() && !s.clk.Now().Before(h.until):
		delete(s.entries, key)
		return ttlperwrite.Entry{}, ttlperwrite.ErrExpired
	}
	return h.entry, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("ttlperwritetest: nil context")
	}
	return ctx.Err()
}

// Reopen returns an InMemory with the clock and the entries of prior. A restart
// after a crash keeps every entry that Put wrote. Do not use prior after the
// call. The two stores share the entries but not the lock.
func Reopen(prior *InMemory) *InMemory {
	return &InMemory{clk: prior.clk, entries: prior.entries}
}
