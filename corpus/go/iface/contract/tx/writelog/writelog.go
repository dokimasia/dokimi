// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package writelog contains a control of the suite of Contract: a store that
// appends each write to one log and reads the last write of a committed
// transaction.
package writelog

import (
	"context"
	"fmt"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/tx"
)

// write is one Put of a transaction.
type write struct {
	key   string
	value tx.Value
	tx    int64
}

// Store implements [tx.Contract] over a log of every write. commits maps each
// committed transaction to the position of its commit, and open contains each
// transaction that Begin opened and no Commit or Rollback settled. Get returns
// the write of the key whose transaction committed last, and the last such
// write of that transaction. The corpus applies the suite of Contract to it as
// a control of the transactional declarations.
//
//dokimi:subject tx.Contract new=New
type Store struct {
	open    map[int64]bool
	commits map[int64]int
	log     []write
	last    int64
	mu      sync.Mutex
}

// The compiler checks that Store implements [tx.Contract].
var _ tx.Contract = (*Store)(nil)

// New returns an empty Store.
func New() *Store { return &Store{open: map[int64]bool{}, commits: map[int64]int{}} }

// Begin implements [tx.Contract].
func (s *Store) Begin(ctx context.Context) (tx.Tx, error) {
	if err := ctx.Err(); err != nil {
		return tx.Tx{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.last++
	s.open[s.last] = true
	return tx.Tx{ID: s.last}, nil
}

// Put implements [tx.Contract].
func (s *Store) Put(ctx context.Context, h tx.Tx, key string, v tx.Value) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.open[h.ID] {
		return fmt.Errorf("writelog: transaction %d: %w", h.ID, tx.ErrTxClosed)
	}
	s.log = append(s.log, write{key: key, value: v, tx: h.ID})
	return nil
}

// Commit implements [tx.Contract].
func (s *Store) Commit(ctx context.Context, h tx.Tx) error {
	return s.settle(ctx, h, true)
}

// Rollback implements [tx.Contract].
func (s *Store) Rollback(ctx context.Context, h tx.Tx) error {
	return s.settle(ctx, h, false)
}

// Get implements [tx.Contract].
func (s *Store) Get(ctx context.Context, key string) (tx.Value, error) {
	if err := ctx.Err(); err != nil {
		return tx.Value{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	found, latest := -1, -1
	for i, w := range s.log {
		if at, committed := s.commits[w.tx]; committed && w.key == key && at >= latest {
			found, latest = i, at
		}
	}
	if found < 0 {
		return tx.Value{}, fmt.Errorf("writelog: key %q: %w", key, tx.ErrNotFound)
	}
	return s.log[found].value, nil
}

// settle closes the open transaction h, and records its commit where commit
// is true.
func (s *Store) settle(ctx context.Context, h tx.Tx, commit bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.open[h.ID] {
		return fmt.Errorf("writelog: transaction %d: %w", h.ID, tx.ErrTxClosed)
	}
	delete(s.open, h.ID)
	if commit {
		s.commits[h.ID] = len(s.commits)
	}
	return nil
}
