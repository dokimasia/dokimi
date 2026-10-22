// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package publisherexactlyoncetest contains the double and the suite of
// Contract. InMemory is the control of the suite.
package publisherexactlyoncetest

import (
	"context"
	"errors"
	"sync"

	publisherexactlyonce "go.dokimi.dev/dokimi/corpus/go/iface/contract/publisher-exactlyonce"
)

// subscriberBuffer is the number of messages that the channel of a subscriber
// buffers.
const subscriberBuffer = 8

// ErrFull is the error of Publish and Replay for a subscriber whose buffer is
// full.
var ErrFull = errors.New("publisherexactlyoncetest: a subscriber is too far behind to take the message")

// InMemory implements [publisherexactlyonce.Contract] in memory. The corpus
// applies the suite of Contract to it as a control.
//
//dokimi:subject publisherexactlyonce.Contract new=NewInMemory
type InMemory struct {
	delivered   map[publisherexactlyonce.Value]bool
	subscribers []chan publisherexactlyonce.Value
	mu          sync.Mutex
}

// The compiler checks that InMemory implements [publisherexactlyonce.Contract].
var _ publisherexactlyonce.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{delivered: map[publisherexactlyonce.Value]bool{}}
}

// Publish implements [publisherexactlyonce.Contract].
func (s *InMemory) Publish(ctx context.Context, v publisherexactlyonce.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.deliver(v)
}

// Replay implements [publisherexactlyonce.Contract].
func (s *InMemory) Replay(ctx context.Context, v publisherexactlyonce.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.delivered[v] {
		return nil
	}
	return s.deliver(v)
}

// Subscribe implements [publisherexactlyonce.Contract].
func (s *InMemory) Subscribe(ctx context.Context) (<-chan publisherexactlyonce.Value, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	ch := make(chan publisherexactlyonce.Value, subscriberBuffer)
	s.subscribers = append(s.subscribers, ch)
	return ch, nil
}

// deliver sends v to each subscriber and records v as delivered. The caller
// has the lock.
func (s *InMemory) deliver(v publisherexactlyonce.Value) error {
	for _, ch := range s.subscribers {
		select {
		case ch <- v:
		default:
			return ErrFull
		}
	}
	s.delivered[v] = true
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("publisherexactlyoncetest: nil context")
	}
	return ctx.Err()
}
