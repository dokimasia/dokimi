// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package waitstest contains the suite of Lock. InMemory is the control of
// the suite.
package waitstest

import (
	"context"
	"sync"
	"time"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/waits"
)

// lease is the time for which a holder holds the lock. The lock releases
// itself when the lease ends.
const lease = time.Minute

// InMemory implements [waits.Lock] in memory. A holder holds the lock until
// Release, or until its lease ends. The corpus applies the suite of Lock to it
// as a control.
//
//dokimi:subject waits.Lock new=NewInMemory
type InMemory struct {
	until    time.Time
	released chan struct{}
	mu       sync.Mutex
	held     bool
}

// The compiler checks that InMemory implements [waits.Lock].
var _ waits.Lock = (*InMemory)(nil)

// NewInMemory returns a lock that no caller holds.
func NewInMemory() *InMemory {
	return &InMemory{released: make(chan struct{})}
}

// Acquire implements [waits.Lock]. It waits until the holder releases the lock
// or its lease ends, or until ctx ends. It returns the error of a context that
// ended before the call.
func (l *InMemory) Acquire(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	for {
		l.mu.Lock()
		now := time.Now()
		if !l.held || !now.Before(l.until) {
			l.held, l.until = true, now.Add(lease)
			l.mu.Unlock()
			return nil
		}
		released, wait := l.released, l.until.Sub(now)
		l.mu.Unlock()
		timer := time.NewTimer(wait)
		select {
		case <-released:
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		}
		timer.Stop()
	}
}

// Release implements [waits.Lock]. It does nothing for a lock that no caller
// holds.
func (l *InMemory) Release() {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.held {
		return
	}
	l.held = false
	close(l.released)
	l.released = make(chan struct{})
}

// Held implements [waits.Lock]. A holder whose lease ended holds the lock no
// more.
func (l *InMemory) Held() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.held && time.Now().Before(l.until)
}
