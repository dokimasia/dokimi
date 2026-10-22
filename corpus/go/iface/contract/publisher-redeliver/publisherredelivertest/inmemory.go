// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package publisherredelivertest contains the double and the suite of Contract.
// InMemory is the control of the suite.
package publisherredelivertest

import (
	"context"
	"errors"
	"sync"

	publisherredeliver "go.dokimi.dev/dokimi/corpus/go/iface/contract/publisher-redeliver"
)

// subscriberBuffer is the number of messages that the channel of a subscriber
// buffers.
const subscriberBuffer = 8

// ErrFull is the error of Publish and Republish for a subscriber whose buffer
// is full.
var ErrFull = errors.New("publisherredelivertest: a subscriber is too far behind to take the message")

// InMemory implements [publisherredeliver.Contract] in memory. The corpus
// applies the suite of Contract to it as a control.
//
//dokimi:subject publisherredeliver.Contract new=NewInMemory
type InMemory struct {
	subscribers []chan publisherredeliver.Value
	mu          sync.Mutex
}

// The compiler checks that InMemory implements [publisherredeliver.Contract].
var _ publisherredeliver.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Publish implements [publisherredeliver.Contract].
func (s *InMemory) Publish(ctx context.Context, v publisherredeliver.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	return s.deliver(v)
}

// Republish implements [publisherredeliver.Contract].
func (s *InMemory) Republish(ctx context.Context, v publisherredeliver.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	return s.deliver(v)
}

// Subscribe implements [publisherredeliver.Contract].
func (s *InMemory) Subscribe(ctx context.Context) (<-chan publisherredeliver.Value, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	ch := make(chan publisherredeliver.Value, subscriberBuffer)
	s.subscribers = append(s.subscribers, ch)
	return ch, nil
}

// deliver sends v to each subscriber. It returns ErrFull for the first
// subscriber whose buffer is full.
func (s *InMemory) deliver(v publisherredeliver.Value) error {
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

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("publisherredelivertest: nil context")
	}
	return ctx.Err()
}
