// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package publishertest contains the double and the suite of Contract. InMemory
// is the control of the suite.
package publishertest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/publisher"
)

// subscriberBuffer is the number of messages that the channel of a subscriber
// buffers.
const subscriberBuffer = 8

// ErrFull is the error of Publish for a subscriber whose buffer is full.
var ErrFull = errors.New("publishertest: a subscriber is too far behind to take the message")

// InMemory implements [publisher.Contract] in memory. The corpus applies the
// suite of Contract to it as a control.
//
//dokimi:subject publisher.Contract new=NewInMemory
type InMemory struct {
	subscribers []chan publisher.Value
	mu          sync.Mutex
}

// The compiler checks that InMemory implements [publisher.Contract].
var _ publisher.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Publish implements [publisher.Contract].
func (s *InMemory) Publish(ctx context.Context, v publisher.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, ch := range s.subscribers {
		select {
		case ch <- v:
		default:
			return ErrFull
		}
	}
	return nil
}

// Subscribe implements [publisher.Contract].
func (s *InMemory) Subscribe(ctx context.Context) (<-chan publisher.Value, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	ch := make(chan publisher.Value, subscriberBuffer)
	s.subscribers = append(s.subscribers, ch)
	return ch, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("publishertest: nil context")
	}
	return ctx.Err()
}
