// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package mailbox contains a control of the suite of Contract: a publisher
// that gives each subscriber an unbounded mailbox and a goroutine that
// delivers it.
package mailbox

import (
	"context"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/publisher"
)

// Publisher implements [publisher.Contract] with a mailbox for each
// subscriber. Publish appends the value to each mailbox and never refuses
// one. The goroutine of a subscriber sends the values of its mailbox in order
// on an unbuffered channel, and ends with the context of its Subscribe. The
// corpus applies the suite of Contract to it as a control of the publishers.
//
//dokimi:subject publisher.Contract new=New
type Publisher struct {
	boxes []*box
	mu    sync.Mutex
}

// box is the mailbox of one subscriber. ready has one element while values
// is not empty.
type box struct {
	ready  chan struct{}
	values []publisher.Value
	mu     sync.Mutex
}

// The compiler checks that Publisher implements [publisher.Contract].
var _ publisher.Contract = (*Publisher)(nil)

// New returns a Publisher without subscribers.
func New() *Publisher { return &Publisher{} }

// Publish implements [publisher.Contract].
func (p *Publisher) Publish(ctx context.Context, v publisher.Value) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	for _, b := range p.boxes {
		b.mu.Lock()
		b.values = append(b.values, v)
		b.mu.Unlock()
		select {
		case b.ready <- struct{}{}:
		default:
		}
	}
	return nil
}

// Subscribe implements [publisher.Contract].
func (p *Publisher) Subscribe(ctx context.Context) (<-chan publisher.Value, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	b := &box{ready: make(chan struct{}, 1)}
	p.mu.Lock()
	p.boxes = append(p.boxes, b)
	p.mu.Unlock()
	out := make(chan publisher.Value)
	go b.deliver(ctx, out)
	return out, nil
}

// deliver sends the values of the mailbox on out in order until ctx ends.
func (b *box) deliver(ctx context.Context, out chan<- publisher.Value) {
	for {
		b.mu.Lock()
		if len(b.values) == 0 {
			b.mu.Unlock()
			select {
			case <-b.ready:
				continue
			case <-ctx.Done():
				return
			}
		}
		v := b.values[0]
		b.values = b.values[1:]
		b.mu.Unlock()
		select {
		case out <- v:
		case <-ctx.Done():
			return
		}
	}
}
