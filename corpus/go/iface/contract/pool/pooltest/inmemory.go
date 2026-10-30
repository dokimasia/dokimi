// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package pooltest contains the double and the suite of Contract. InMemory is
// the subject of the suite.
package pooltest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/pool"
)

// ErrExhausted is the error of Get on a pool without a free value.
var ErrExhausted = errors.New("pooltest: the pool has no free value")

// InMemory implements [pool.Contract] in memory. The corpus applies the suite
// of Contract to it as a subject.
//
//dokimi:subject pool.Contract new=NewInMemory
type InMemory struct {
	free  []pool.Value
	inUse int
	gets  int
	puts  int
	mu    sync.Mutex
}

// The compiler checks that InMemory implements [pool.Contract].
var _ pool.Contract = (*InMemory)(nil)

// NewInMemory returns a pool of two free values.
func NewInMemory() *InMemory {
	return NewInMemoryOf(
		pool.Value{Key: "a", Body: "first"},
		pool.Value{Key: "b", Body: "second"},
	)
}

// NewInMemoryOf returns a pool of the free values.
func NewInMemoryOf(values ...pool.Value) *InMemory { return &InMemory{free: values} }

// Stats implements [pool.Contract].
func (s *InMemory) Stats(ctx context.Context) (pool.Stats, error) {
	if err := contextErr(ctx); err != nil {
		return pool.Stats{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return pool.Stats{Gets: s.gets, Puts: s.puts, Outstanding: s.inUse}, nil
}

// Get implements [pool.Contract].
func (s *InMemory) Get(ctx context.Context) (pool.Value, error) {
	if err := contextErr(ctx); err != nil {
		return pool.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.free) == 0 {
		return pool.Value{}, ErrExhausted
	}
	v := s.free[len(s.free)-1]
	s.free = s.free[:len(s.free)-1]
	s.inUse++
	s.gets++
	return v, nil
}

// Put implements [pool.Contract]. It counts a put only while a value is
// outstanding.
func (s *InMemory) Put(ctx context.Context, v pool.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.free = append(s.free, v)
	if s.inUse > 0 {
		s.inUse--
		s.puts++
	}
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("pooltest: nil context")
	}
	return ctx.Err()
}
