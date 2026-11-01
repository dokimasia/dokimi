// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package referencetest

import (
	"context"
	"sync"

	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/reference"
)

// RingQueue implements [reference.Fifo] over a ring buffer that doubles when
// it is full. head is the position of the oldest item, and count is the
// number of items. The corpus applies the suite of Fifo to it with a slice as
// its reference.
//
//dokimi:subject reference.Fifo new=NewRingQueue options=RingQueueOptions
type RingQueue struct {
	items []string
	head  int
	count int
	mu    sync.Mutex
}

// The compiler checks that RingQueue implements [reference.Fifo].
var _ reference.Fifo = (*RingQueue)(nil)

// NewRingQueue returns an empty RingQueue with room for two items.
func NewRingQueue() *RingQueue { return &RingQueue{items: make([]string, 2)} }

// RingQueueOptions returns the reference of a RingQueue: a queue over a slice.
func RingQueueOptions() []FifoOption {
	return []FifoOption{conform.Reference(func(*FifoHandle) reference.Fifo { return &sliceQueue{} })}
}

// Push implements [reference.Fifo].
func (q *RingQueue) Push(ctx context.Context, item string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.count == len(q.items) {
		grown := make([]string, 2*len(q.items))
		for i := range q.count {
			grown[i] = q.items[(q.head+i)%len(q.items)]
		}
		q.items, q.head = grown, 0
	}
	q.items[(q.head+q.count)%len(q.items)] = item
	q.count++
	return nil
}

// Pop implements [reference.Fifo].
func (q *RingQueue) Pop(ctx context.Context) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	if q.count == 0 {
		return "", reference.ErrEmpty
	}
	item := q.items[q.head]
	q.head = (q.head + 1) % len(q.items)
	q.count--
	return item, nil
}

// Len implements [reference.Fifo].
func (q *RingQueue) Len(ctx context.Context) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.count, nil
}

// sliceQueue is the reference of a RingQueue: the items in a slice, the
// oldest first.
type sliceQueue struct {
	items []string
}

// Push implements [reference.Fifo].
func (q *sliceQueue) Push(_ context.Context, item string) error {
	q.items = append(q.items, item)
	return nil
}

// Pop implements [reference.Fifo].
func (q *sliceQueue) Pop(context.Context) (string, error) {
	if len(q.items) == 0 {
		return "", reference.ErrEmpty
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, nil
}

// Len implements [reference.Fifo].
func (q *sliceQueue) Len(context.Context) (int, error) { return len(q.items), nil }
