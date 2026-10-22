// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package outboxtest contains the double and the suite of Contract. InMemory is
// the control of the suite.
package outboxtest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/outbox"
)

// subscriberBuffer is the number of records that the channel of a subscriber
// buffers.
const subscriberBuffer = 8

// ErrFull is the error of Append and Subscribe for a subscriber whose buffer
// is full. The outbox reports the record instead of dropping it.
var ErrFull = errors.New("outboxtest: a subscriber is too far behind to take the record")

// InMemory implements [outbox.Contract] in memory. The corpus applies the
// suite of Contract to it as a control.
//
//dokimi:subject outbox.Contract new=NewInMemory
type InMemory struct {
	mu          sync.Mutex
	records     []outbox.Value
	subscribers []chan outbox.Value
}

// The compiler checks that InMemory implements [outbox.Contract].
var _ outbox.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Append implements [outbox.Contract].
func (s *InMemory) Append(ctx context.Context, v outbox.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = append(s.records, v)
	for _, ch := range s.subscribers {
		select {
		case ch <- v:
		default:
			return ErrFull
		}
	}
	return nil
}

// Subscribe implements [outbox.Contract].
func (s *InMemory) Subscribe(ctx context.Context) (<-chan outbox.Value, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	ch := make(chan outbox.Value, subscriberBuffer)
	for _, v := range s.records {
		select {
		case ch <- v:
		default:
			return nil, ErrFull
		}
	}
	s.subscribers = append(s.subscribers, ch)
	return ch, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("outboxtest: nil context")
	}
	return ctx.Err()
}
