// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package publisheratleastoncetest contains the double and the suite of
// Contract. InMemory is the subject of the suite.
package publisheratleastoncetest

import (
	"context"
	"errors"
	"sync"

	publisheratleastonce "go.dokimi.dev/dokimi/corpus/go/iface/contract/publisher-atleastonce"
)

// subscriberBuffer is the number of messages that the channel of a subscriber
// buffers.
const subscriberBuffer = 8

// ErrFull is the error of Publish for a subscriber whose buffer is full.
var ErrFull = errors.New("publisheratleastoncetest: a subscriber is too far behind to take the message")

// InMemory implements [publisheratleastonce.Contract] in memory. The corpus
// applies the suite of Contract to it as a subject.
//
//dokimi:subject publisheratleastonce.Contract new=NewInMemory
type InMemory struct {
	subscribers []chan publisheratleastonce.Value
	mu          sync.Mutex
}

// The compiler checks that InMemory implements [publisheratleastonce.Contract].
var _ publisheratleastonce.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Publish implements [publisheratleastonce.Contract].
func (s *InMemory) Publish(ctx context.Context, v publisheratleastonce.Value) error {
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

// Subscribe implements [publisheratleastonce.Contract].
func (s *InMemory) Subscribe(ctx context.Context) (<-chan publisheratleastonce.Value, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	ch := make(chan publisheratleastonce.Value, subscriberBuffer)
	s.subscribers = append(s.subscribers, ch)
	return ch, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("publisheratleastoncetest: nil context")
	}
	return ctx.Err()
}
