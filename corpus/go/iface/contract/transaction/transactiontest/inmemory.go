// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package transactiontest contains the double and the suite of Contract.
// InMemory is the subject of the suite.
package transactiontest

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"strconv"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/transaction"
)

// RunKey is the key of the entry that each committed run writes.
const RunKey = "run"

// InMemory implements [transaction.Contract] in memory. The corpus applies
// the suite of Contract to it as a subject.
//
//dokimi:subject transaction.Contract new=NewInMemory
type InMemory struct {
	entries map[string]transaction.Value
	// staged contains the writes of the open run alone. A commit copies them
	// into the entries key by key, so a Put outside the run keeps its value.
	staged map[string]transaction.Value
	runs   int
	runMu  sync.Mutex
	mu     sync.Mutex
}

// The compiler checks that InMemory implements [transaction.Contract].
var _ transaction.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{entries: map[string]transaction.Value{}}
}

// Run implements [transaction.Contract].
func (s *InMemory) Run(ctx context.Context, body func(ctx context.Context) error) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.runMu.Lock()
	defer s.runMu.Unlock()

	s.mu.Lock()
	runs := s.runs + 1
	s.staged = map[string]transaction.Value{
		RunKey: {Key: RunKey, Body: strconv.Itoa(runs)},
	}
	s.mu.Unlock()

	var bodyErr error
	if body != nil {
		bodyErr = body(ctx)
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	staged := s.staged
	s.staged = nil
	if bodyErr != nil {
		return bodyErr
	}
	maps.Copy(s.entries, staged)
	s.runs = runs
	return nil
}

// Put implements [transaction.Contract].
func (s *InMemory) Put(ctx context.Context, key string, v transaction.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.staged != nil {
		s.staged[key] = v
		return nil
	}
	s.entries[key] = v
	return nil
}

// Get implements [transaction.Contract].
func (s *InMemory) Get(ctx context.Context, key string) (transaction.Value, error) {
	if err := contextErr(ctx); err != nil {
		return transaction.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, present := s.entries[key]
	if !present {
		return transaction.Value{}, fmt.Errorf("transactiontest: key %q: %w", key, transaction.ErrNotFound)
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("transactiontest: nil context")
	}
	return ctx.Err()
}
