// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package snapshotisolationtest contains the double and the suite of Mixed.
// InMemory is the subject of the suite.
package snapshotisolationtest

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/snapshotisolation"
)

// ErrConflict is the error of a commit of a transaction that wrote a key that
// another transaction committed a write of after its begin.
var ErrConflict = errors.New("snapshotisolationtest: conflict")

// InMemory implements [snapshotisolation.Mixed] in memory. The corpus applies
// the suite of Mixed to it as a subject.
//
// A transaction reads a copy of the committed numbers from its begin, with its
// own writes over it. Its commit fails with ErrConflict where another
// transaction committed a write of a key that it wrote after its begin, and
// applies its writes otherwise. The store allows write skew, so its
// transactions are not serializable. No call keeps the store locked after it
// returns, so the transactions of concurrent clients interleave.
//
//dokimi:subject snapshotisolation.Mixed new=NewInMemory
type InMemory struct {
	numbers map[string]int64
	changed map[string]int64
	commits int64
	mu      sync.Mutex
}

// The compiler checks that InMemory implements [snapshotisolation.Mixed].
var _ snapshotisolation.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{numbers: map[string]int64{}, changed: map[string]int64{}}
}

// Begin implements [snapshotisolation.Mixed].
func (s *InMemory) Begin(ctx context.Context) (snapshotisolation.Txn, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return &txn{store: s, snapshot: maps.Clone(s.numbers), writes: map[string]int64{}, begin: s.commits}, nil
}

// txn implements [snapshotisolation.Txn] for a transaction of InMemory: the
// numbers at its begin, its writes, and the number of commits before its
// begin.
type txn struct {
	store    *InMemory
	snapshot map[string]int64
	writes   map[string]int64
	begin    int64
	closed   bool
}

// Commit implements [snapshotisolation.Txn].
func (t *txn) Commit(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s := t.store
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := t.close(); err != nil {
		return err
	}
	for key := range t.writes {
		if s.changed[key] > t.begin {
			return fmt.Errorf("snapshotisolationtest: key %q: %w", key, ErrConflict)
		}
	}
	s.commits++
	for key, v := range t.writes {
		s.numbers[key] = v
		s.changed[key] = s.commits
	}
	return nil
}

// Rollback implements [snapshotisolation.Txn].
func (t *txn) Rollback(ctx context.Context) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	t.store.mu.Lock()
	defer t.store.mu.Unlock()
	return t.close()
}

// Set implements [snapshotisolation.Txn].
func (t *txn) Set(ctx context.Context, key string, v int64) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	t.store.mu.Lock()
	defer t.store.mu.Unlock()
	if t.closed {
		return errClosed
	}
	t.writes[key] = v
	return nil
}

// Load implements [snapshotisolation.Txn].
func (t *txn) Load(ctx context.Context, key string) (int64, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	t.store.mu.Lock()
	defer t.store.mu.Unlock()
	if t.closed {
		return 0, errClosed
	}
	if v, written := t.writes[key]; written {
		return v, nil
	}
	return t.snapshot[key], nil
}

// errClosed is the error of a call of a closed transaction.
var errClosed = fmt.Errorf("snapshotisolationtest: %w", snapshotisolation.ErrClosed)

// close closes the transaction, and returns errClosed for a transaction that
// is closed already.
func (t *txn) close() error {
	if t.closed {
		return errClosed
	}
	t.closed = true
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("snapshotisolationtest: nil context")
	}
	return ctx.Err()
}
