// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package queuetest contains the suite of Queue, the check that a consumer
// adds to the suite inside a testing/synctest bubble, and the control of the
// suite.
package queuetest

import (
	"context"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/queue"
)

// room is the number of items that the queue of NewInMemory holds.
const room = 2

// InMemory implements [queue.Queue] in memory. The corpus applies the suite of
// Queue to it as a control.
//
//dokimi:subject queue.Queue new=NewInMemory
type InMemory struct {
	items  chan string
	closed chan struct{}
	once   sync.Once
}

// The compiler checks that InMemory implements [queue.Queue].
var _ queue.Queue = (*InMemory)(nil)

// NewInMemory returns an empty queue with room for two items.
func NewInMemory() *InMemory {
	return &InMemory{items: make(chan string, room), closed: make(chan struct{})}
}

// Write implements [queue.Queue].
func (q *InMemory) Write(ctx context.Context, item string) error {
	select {
	case <-q.closed:
		return queue.ErrClosed
	default:
	}
	select {
	case q.items <- item:
		return nil
	case <-q.closed:
		return queue.ErrClosed
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Read implements [queue.Queue].
func (q *InMemory) Read(ctx context.Context) (string, error) {
	select {
	case item := <-q.items:
		return item, nil
	case <-q.closed:
		select {
		case item := <-q.items:
			return item, nil
		default:
			return "", queue.ErrClosed
		}
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// Close implements [queue.Queue].
func (q *InMemory) Close() error {
	err := queue.ErrClosed
	q.once.Do(func() {
		close(q.closed)
		err = nil
	})
	return err
}
