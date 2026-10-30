// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package ring contains a control of the suite of Queue: a ring buffer of a
// fixed size under a mutex.
package ring

import (
	"context"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/queue"
)

// size is the number of items that a Buffer holds.
const size = 2

// Buffer implements [queue.Queue] over an array of size items. head is the
// position of the oldest item, and count is the number of items. A call that
// waits selects on changed, which each change of the state closes and
// replaces. The corpus applies the suite of Queue to it as a control of the
// queues.
//
//dokimi:subject queue.Queue new=New
type Buffer struct {
	changed chan struct{}
	items   [size]string
	head    int
	count   int
	mu      sync.Mutex
	closed  bool
}

// The compiler checks that Buffer implements [queue.Queue].
var _ queue.Queue = (*Buffer)(nil)

// New returns an empty Buffer.
func New() *Buffer { return &Buffer{changed: make(chan struct{})} }

// Write implements [queue.Queue].
func (b *Buffer) Write(ctx context.Context, item string) error {
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		b.mu.Lock()
		if b.closed {
			b.mu.Unlock()
			return queue.ErrClosed
		}
		if b.count < size {
			b.items[(b.head+b.count)%size] = item
			b.count++
			b.signal()
			b.mu.Unlock()
			return nil
		}
		changed := b.changed
		b.mu.Unlock()
		select {
		case <-changed:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// Read implements [queue.Queue].
func (b *Buffer) Read(ctx context.Context) (string, error) {
	for {
		b.mu.Lock()
		if b.count > 0 {
			item := b.items[b.head]
			b.head = (b.head + 1) % size
			b.count--
			b.signal()
			b.mu.Unlock()
			return item, nil
		}
		if b.closed {
			b.mu.Unlock()
			return "", queue.ErrClosed
		}
		changed := b.changed
		b.mu.Unlock()
		select {
		case <-changed:
		case <-ctx.Done():
			return "", ctx.Err()
		}
	}
}

// Close implements [queue.Queue].
func (b *Buffer) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return queue.ErrClosed
	}
	b.closed = true
	b.signal()
	return nil
}

// signal wakes each call that waits. The caller holds mu.
func (b *Buffer) signal() {
	close(b.changed)
	b.changed = make(chan struct{})
}
