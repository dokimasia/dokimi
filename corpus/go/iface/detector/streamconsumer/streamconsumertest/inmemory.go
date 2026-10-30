// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package streamconsumertest contains the double and the suite of
// StreamConsumer. InMemory is the subject of the suite.
package streamconsumertest

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/streamconsumer"
)

// ErrNoSource is the error of Ingest for a nil source.
var ErrNoSource = errors.New("streamconsumertest: nil source")

// ErrSourceFailed is the error of a source that fails, which a test passes to
// Ingest.
var ErrSourceFailed = errors.New("streamconsumertest: source failed")

// InMemory implements [streamconsumer.StreamConsumer] in memory. The corpus
// applies the suite of StreamConsumer to it as a subject.
//
//dokimi:subject streamconsumer.StreamConsumer new=NewInMemory
type InMemory struct {
	ingested []streamconsumer.Value
	mu       sync.Mutex
}

// The compiler checks that InMemory implements [streamconsumer.StreamConsumer].
var _ streamconsumer.StreamConsumer = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory { return &InMemory{} }

// Ingest implements [streamconsumer.StreamConsumer].
func (s *InMemory) Ingest(ctx context.Context, src streamconsumer.Source) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	if src == nil {
		return 0, ErrNoSource
	}

	var taken []streamconsumer.Value
	for {
		v, ok, err := src.Next(ctx)
		if err != nil {
			return 0, fmt.Errorf("streamconsumertest: read the source: %w", err)
		}
		if !ok {
			break
		}
		taken = append(taken, v)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	s.ingested = append(s.ingested, taken...)
	return len(taken), nil
}

// SliceSource is a [streamconsumer.Source] of the values of a slice, in their
// order.
type SliceSource struct {
	values []streamconsumer.Value
	at     int
	mu     sync.Mutex
}

// The compiler checks that SliceSource implements [streamconsumer.Source].
var _ streamconsumer.Source = (*SliceSource)(nil)

// NewSliceSource returns a source of the values.
func NewSliceSource(values ...streamconsumer.Value) *SliceSource {
	return &SliceSource{values: values}
}

// Next implements [streamconsumer.Source].
func (s *SliceSource) Next(ctx context.Context) (streamconsumer.Value, bool, error) {
	if err := contextErr(ctx); err != nil {
		return streamconsumer.Value{}, false, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.at >= len(s.values) {
		return streamconsumer.Value{}, false, nil
	}
	v := s.values[s.at]
	s.at++
	return v, true, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("streamconsumertest: nil context")
	}
	return ctx.Err()
}
