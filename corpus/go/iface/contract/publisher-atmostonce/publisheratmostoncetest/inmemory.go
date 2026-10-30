// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package publisheratmostoncetest contains the double and the suite of
// Contract. InMemory is the subject of the suite.
package publisheratmostoncetest

import (
	"context"
	"errors"
	"sync"

	publisheratmostonce "go.dokimi.dev/dokimi/corpus/go/iface/contract/publisher-atmostonce"
)

// subscriberBuffer is the number of messages that the channel of a subscriber
// buffers.
const subscriberBuffer = 8

// ErrFull is the error of Publish for a subscriber whose buffer is full.
var ErrFull = errors.New("publisheratmostoncetest: a subscriber is too far behind to take the message")

// InMemory implements [publisheratmostonce.Contract] in memory. The corpus
// applies the suite of Contract to it as a subject.
//
//dokimi:subject publisheratmostonce.Contract new=NewInMemory
type InMemory struct {
	subscribers []chan publisheratmostonce.Value
	mu          sync.Mutex
}

// The compiler checks that InMemory implements [publisheratmostonce.Contract].
var _ publisheratmostonce.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Publish implements [publisheratmostonce.Contract].
func (s *InMemory) Publish(ctx context.Context, v publisheratmostonce.Value) error {
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

// Subscribe implements [publisheratmostonce.Contract].
func (s *InMemory) Subscribe(ctx context.Context) (<-chan publisheratmostonce.Value, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	ch := make(chan publisheratmostonce.Value, subscriberBuffer)
	s.subscribers = append(s.subscribers, ch)
	return ch, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("publisheratmostoncetest: nil context")
	}
	return ctx.Err()
}
