// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package ttltest

import (
	"context"
	"errors"
	"sync"
	"time"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/ttl"
)

// StaleKey is a key whose value Put stores with an age of one lifetime. Read
// reports the value as expired at once.
const StaleKey = "ttltest-stale"

// InMemory implements [ttl.Mixed] in memory. The corpus applies the suite of
// Mixed to it as a control. NewInMemory builds it on a clock of its own. The
// check of the lifetime advances the clock of its handle, so the control
// excuses that check. The store of the package memory applies the check on the
// clock of the handle.
//
//dokimi:subject ttl.Mixed new=NewInMemory
//dokimi:excuse Mixed.Read/expires-after-the-lifetime reason="the check does not advance the clock that NewInMemory reads"
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
