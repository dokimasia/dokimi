// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package txtest contains the double and the suite of Contract. InMemory is the
// control of the suite.
package txtest

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/tx"
)

// InMemory implements [tx.Contract] in memory. The corpus applies the suite
// of Contract to it as a control.
//
//dokimi:subject tx.Contract new=NewInMemory
type InMemory struct {
	staged    map[int64]map[string]tx.Value
	committed map[string]tx.Value
	nextID    int64
	mu        sync.Mutex
}

// The compiler checks that InMemory implements [tx.Contract].
var _ tx.Contract = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{
		staged:    map[int64]map[string]tx.Value{},
		committed: map[string]tx.Value{},
	}
}

// Begin implements [tx.Contract].
func (s *InMemory) Begin(ctx context.Context) (tx.Tx, error) {
	if err := contextErr(ctx); err != nil {
		return tx.Tx{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nextID++
	s.staged[s.nextID] = map[string]tx.Value{}
	return tx.Tx{ID: s.nextID}, nil
}

// Put implements [tx.Contract].
func (s *InMemory) Put(ctx context.Context, h tx.Tx, key string, v tx.Value) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	stage, open := s.staged[h.ID]
	if !open {
		return fmt.Errorf("txtest: transaction %d: %w", h.ID, tx.ErrTxClosed)
	}
	stage[key] = v
	return nil
}

// Commit implements [tx.Contract].
func (s *InMemory) Commit(ctx context.Context, h tx.Tx) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	stage, open := s.staged[h.ID]
	if !open {
		return fmt.Errorf("txtest: transaction %d: %w", h.ID, tx.ErrTxClosed)
	}
	maps.Copy(s.committed, stage)
	delete(s.staged, h.ID)
	return nil
}

// Rollback implements [tx.Contract].
func (s *InMemory) Rollback(ctx context.Context, h tx.Tx) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, open := s.staged[h.ID]; !open {
		return fmt.Errorf("txtest: transaction %d: %w", h.ID, tx.ErrTxClosed)
	}
	delete(s.staged, h.ID)
	return nil
}

// Get implements [tx.Contract].
func (s *InMemory) Get(ctx context.Context, key string) (tx.Value, error) {
	if err := contextErr(ctx); err != nil {
		return tx.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	v, present := s.committed[key]
	if !present {
		return tx.Value{}, fmt.Errorf("txtest: key %q: %w", key, tx.ErrNotFound)
	}
	return v, nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("txtest: nil context")
	}
	return ctx.Err()
}
