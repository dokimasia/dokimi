// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package ttltest contains the suite and the double of Mixed. InMemory is the
// subject of the suite.
package ttltest

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/ttl"
)

// lifetime is the duration of the mixin ttl of Read.
const lifetime = time.Minute

// StaleKey is a key whose value Put stores with an age of one lifetime. Read
// reports the value as expired at once.
const StaleKey = "ttltest-stale"

// InMemory implements [ttl.Mixed] in memory. The corpus applies the suite of
// Mixed to it as a subject, on the clock of each check. A write of the
// subject allocates nothing after a warm-up of 1000 writes, which a budget of
// the subject alone states.
//
//dokimi:subject ttl.Mixed new=NewInMemoryFromHandle
//dokimi:budget Mixed.Put allocs=0 p99=1ms warmup=1000
type InMemory struct {
	clk    assert.Clock
	values map[string]entry
	mu     sync.Mutex
}

// entry is a stored value with the instant of its write.
type entry struct {
	at    time.Time
	value ttl.Value
}

// The compiler checks that InMemory implements [ttl.Mixed].
var _ ttl.Mixed = (*InMemory)(nil)

// NewInMemoryOn returns an empty InMemory that reads the time from clk.
func NewInMemoryOn(clk assert.Clock) *InMemory {
	return &InMemory{clk: clk, values: map[string]entry{}}
}

// NewInMemory returns an empty InMemory on a controlled clock at the Unix
// epoch.
func NewInMemory() *InMemory {
	return NewInMemoryOn(assert.NewControlled(time.Unix(0, 0).UTC()))
}

// NewInMemoryFromHandle returns an empty InMemory on the clock of the handle
// h, so a check that advances the clock ends the lifetime of a value.
func NewInMemoryFromHandle(h *MixedHandle) *InMemory { return NewInMemoryOn(h.Clock()) }

// Put implements [ttl.Mixed].
func (s *InMemory) Put(ctx context.Context, v ttl.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	at := s.clk.Now()
	if v.Key == StaleKey {
		at = at.Add(-lifetime)
	}
	s.values[v.Key] = entry{value: v, at: at}
	return nil
}

// Read implements [ttl.Mixed].
func (s *InMemory) Read(ctx context.Context, key string) (ttl.Value, error) {
	if err := contextErr(ctx); err != nil {
		return ttl.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.values[key]
	if !ok || s.clk.Now().Sub(e.at) >= lifetime {
		return ttl.Value{}, ttl.ErrExpired
	}
	return e.value, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("ttltest: nil context")
	}
	return ctx.Err()
}

// Reopen returns an InMemory with the clock and the values of prior. A restart
// after a crash keeps every value that Put wrote. Do not use prior after the
// call. The two stores share the values but not the lock.
func Reopen(prior *InMemory) *InMemory {
	return &InMemory{clk: prior.clk, values: prior.values}
}
