// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package queue is the fixture of the environment of a check. Its queue has
// room for a fixed number of items, a write waits for room, and Close ends
// each write that waits. The check of the suite waits with synctest.Wait until
// a write blocks, so it runs in a testing/synctest bubble. Its falsification on
// the planted defect runs in a bubble too.
package queue

import (
	"context"
	"errors"
)

// ErrClosed is the error of a call that Close ends, and of each write after
// Close.
var ErrClosed = errors.New("queue: closed")

// Queue is a queue of items with room for a fixed number of them.
//
//dokimi:suite
//dokimi:out path=queuetest/
type Queue interface {
	// Write appends the item. It waits while the queue has no room, and
	// returns ErrClosed for a closed queue and when Close ends its wait.
	//dokimi:mixin waits
	Write(ctx context.Context, item string) error

	// Read removes and returns the oldest item. It waits while the queue has
	// no item, and returns ErrClosed for a closed queue without items and when
	// Close ends its wait.
	//dokimi:mixin waits
	Read(ctx context.Context) (string, error)

	// Close closes the queue and ends each call that waits. It returns
	// ErrClosed for a closed queue.
	Close() error
}
