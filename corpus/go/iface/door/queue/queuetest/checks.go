// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package queuetest

import (
	"context"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/queue"
)

// CloseEndsAWaitingWrite checks that Close ends a write that waits for room
// with ErrClosed.
//
//dokimi:check queue.Queue.Close defect=keepsWritesWaiting around=bubble
func CloseEndsAWaitingWrite(h *QueueHandle) {
	q, ctx := h.Subject(), h.Context()
	ended := make(chan error, 1)
	go func() {
		for {
			if err := q.Write(ctx, "item"); err != nil {
				ended <- err
				return
			}
		}
	}()
	synctest.Wait()
	assert.NoError(h, q.Close(), "Close closes the queue")
	synctest.Wait()
	var err error
	select {
	case err = <-ended:
	default:
	}
	assert.ErrorIs(h, err, queue.ErrClosed, h.Statement())
}

// bubble runs each run of a check in a testing/synctest bubble of its own.
func bubble(t *testing.T, run func(*testing.T)) {
	t.Helper()
	synctest.Test(t, run)
}

// deaf is a queue whose Close closes the queue for later writes and leaves
// each waiting write waiting.
type deaf struct {
	items  chan string
	closed atomic.Bool
}

// keepsWritesWaiting returns a deaf queue with room for two items.
func keepsWritesWaiting(*QueueHandle) queue.Queue {
	return &deaf{items: make(chan string, room)}
}

// Write appends the item, and waits while the queue has no room.
func (q *deaf) Write(ctx context.Context, item string) error {
	if q.closed.Load() {
		return queue.ErrClosed
	}
	select {
	case q.items <- item:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Read removes and returns the oldest item, and waits while the queue has no
// item.
func (q *deaf) Read(ctx context.Context) (string, error) {
	select {
	case item := <-q.items:
		return item, nil
	case <-ctx.Done():
		return "", ctx.Err()
	}
}

// Close closes the queue for later writes.
func (q *deaf) Close() error {
	q.closed.Store(true)
	return nil
}
