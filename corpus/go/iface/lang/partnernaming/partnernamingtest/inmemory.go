// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package partnernamingtest contains the double and the suite of Store.
// InMemory is the control of the suite.
package partnernamingtest

import (
	"context"
	"errors"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/partnernaming"
)

// InMemory implements [partnernaming.Store] in memory. The corpus applies
// the suite of Store to it as a control.
//
//dokimi:subject partnernaming.Store new=NewInMemory
type InMemory struct {
	weights map[string]int
	at      map[string]int
	buckets map[int]int
	recs    map[string]partnernaming.Rec
	mu      sync.Mutex
	rev     int64
}

// The compiler checks that InMemory implements [partnernaming.Store].
var _ partnernaming.Store = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{
		weights: map[string]int{},
		at:      map[string]int{},
		buckets: map[int]int{},
		recs:    map[string]partnernaming.Rec{},
	}
}

// Touch implements [partnernaming.Store].
func (s *InMemory) Touch(ctx context.Context, key string, weight int) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.weights[key] += weight
	return nil
}

// Seen implements [partnernaming.Store].
func (s *InMemory) Seen(ctx context.Context, k string) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.weights[k], nil
}

// Move implements [partnernaming.Store].
func (s *InMemory) Move(ctx context.Context, from, to string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.at[from]--
	s.at[to]++
	return nil
}

// At implements [partnernaming.Store].
func (s *InMemory) At(ctx context.Context, where string) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.at[where], nil
}

// Emit implements [partnernaming.Store].
func (s *InMemory) Emit(ctx context.Context, id string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.buckets[len(id)]++
	return nil
}

// Count implements [partnernaming.Store].
func (s *InMemory) Count(ctx context.Context, bucket int) (int, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buckets[bucket], nil
}

// Version implements [partnernaming.Store]. It returns the zero record for a
// key without a record.
func (s *InMemory) Version(ctx context.Context, key string) (partnernaming.Rec, error) {
	if err := contextErr(ctx); err != nil {
		return partnernaming.Rec{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.recs[key], nil
}

// Put implements [partnernaming.Store].
func (s *InMemory) Put(ctx context.Context, r partnernaming.Rec) error {
	return s.store(ctx, r)
}

// Replace implements [partnernaming.Store].
func (s *InMemory) Replace(ctx context.Context, r partnernaming.Rec) error {
	return s.store(ctx, r)
}

// store stores the record r under its key with the next revision.
func (s *InMemory) store(ctx context.Context, r partnernaming.Rec) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rev++
	r.Rev = s.rev
	s.recs[r.Key] = r
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("partnernamingtest: nil context")
	}
	return ctx.Err()
}
