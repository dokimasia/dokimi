// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package serializabletest contains the double and the suite of Mixed. InMemory
// is the subject of the suite.
package serializabletest

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/serializable"
)

// ErrConflict is the error of a commit of a transaction that read a list that
// another transaction changed after its begin.
var ErrConflict = errors.New("serializabletest: conflict")

// InMemory implements [serializable.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a subject.
//
// A transaction reads the committed lists with its own appends at their end.
// Its commit fails with ErrConflict where a list that it read changed after
// its begin, and applies its appends otherwise. Each list that a committed
// transaction read is unchanged at its commit, so the order of the commits is
// a serial order. No call keeps the store locked after it returns, so the
// transactions of concurrent clients interleave.
//
//dokimi:subject serializable.Mixed new=NewInMemory
type InMemory struct {
	lists   map[string][]int
	changed map[string]int64
	open    map[int64]*txn
	next    int64
	commits int64
	mu      sync.Mutex
}

// txn is an open transaction with the keys that it read and its appends.
// begin is the number of commits before its begin.
type txn struct {
	read    map[string]bool
	appends []appended
	begin   int64
}

// appended is one append of a transaction.
type appended struct {
	key string
	v   int
}

// The compiler checks that InMemory implements [serializable.Mixed].
var _ serializable.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{lists: map[string][]int{}, changed: map[string]int64{}, open: map[int64]*txn{}}
}

// Begin implements [serializable.Mixed].
func (s *InMemory) Begin(ctx context.Context) (serializable.Tx, error) {
	if err := contextErr(ctx); err != nil {
		return serializable.Tx{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.next++
	s.open[s.next] = &txn{read: map[string]bool{}, begin: s.commits}
	return serializable.Tx{ID: s.next}, nil
}

// Commit implements [serializable.Mixed].
func (s *InMemory) Commit(ctx context.Context, tx serializable.Tx) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	t, err := s.close(tx)
	if err != nil {
		return err
	}
	for key := range t.read {
		if s.changed[key] > t.begin {
			return fmt.Errorf("serializabletest: transaction %d read %q: %w", tx.ID, key, ErrConflict)
		}
	}
	s.commits++
	for _, a := range t.appends {
		s.lists[a.key] = append(s.lists[a.key], a.v)
		s.changed[a.key] = s.commits
	}
	return nil
}

// Rollback implements [serializable.Mixed].
func (s *InMemory) Rollback(ctx context.Context, tx serializable.Tx) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.close(tx)
	return err
}

// Append implements [serializable.Mixed].
func (s *InMemory) Append(ctx context.Context, tx serializable.Tx, key string, v int) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	t, err := s.of(tx)
	if err != nil {
		return err
	}
	t.appends = append(t.appends, appended{key: key, v: v})
	return nil
}

// Get implements [serializable.Mixed].
func (s *InMemory) Get(ctx context.Context, tx serializable.Tx, key string) ([]int, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	t, err := s.of(tx)
	if err != nil {
		return nil, err
	}
	t.read[key] = true
	list := slices.Clone(s.lists[key])
	for _, a := range t.appends {
		if a.key == key {
			list = append(list, a.v)
		}
	}
	return list, nil
}

// of returns the open transaction of tx, and an error that wraps
// [serializable.ErrClosed] for a transaction that is not open.
func (s *InMemory) of(tx serializable.Tx) (*txn, error) {
	t, open := s.open[tx.ID]
	if !open {
		return nil, fmt.Errorf("serializabletest: transaction %d: %w", tx.ID, serializable.ErrClosed)
	}
	return t, nil
}

// close closes the open transaction of tx and returns it.
func (s *InMemory) close(tx serializable.Tx) (*txn, error) {
	t, err := s.of(tx)
	if err == nil {
		delete(s.open, tx.ID)
	}
	return t, err
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("serializabletest: nil context")
	}
	return ctx.Err()
}
