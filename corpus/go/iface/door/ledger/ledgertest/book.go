// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package ledgertest contains the suite of Ledger and the subject of the
// suite.
package ledgertest

import (
	"context"
	"fmt"
	"maps"
	"sync"

	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/ledger"
)

// Book is the Ledger over a store in memory. The corpus applies the suite of
// Ledger to it as a subject, over a store whose calls fail on demand.
//
//dokimi:subject ledger.Ledger new=NewBook options=BookOptions
type Book struct {
	store *store
}

// The compiler checks that Book implements [ledger.Ledger].
var _ ledger.Ledger = (*Book)(nil)

// NewBook returns a book over an empty store whose calls never fail.
func NewBook() *Book { return &Book{store: &store{values: map[string]string{}}} }

// staged contains the writes of the body of one transaction.
type staged struct {
	writes map[string]string
	mu     sync.Mutex
}

// stagedKey is the key of the transaction of a context.
type stagedKey struct{}

// Run implements [ledger.Ledger].
func (b *Book) Run(ctx context.Context, body func(ctx context.Context) error) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	tx := &staged{writes: map[string]string{}}
	if err := body(context.WithValue(ctx, stagedKey{}, tx)); err != nil {
		return err
	}
	tx.mu.Lock()
	defer tx.mu.Unlock()
	return b.store.commit(tx.writes)
}

// Put implements [ledger.Ledger].
func (b *Book) Put(ctx context.Context, key, value string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if tx, inside := ctx.Value(stagedKey{}).(*staged); inside {
		tx.mu.Lock()
		defer tx.mu.Unlock()
		tx.writes[key] = value
		return nil
	}
	return b.store.commit(map[string]string{key: value})
}

// Get implements [ledger.Ledger].
func (b *Book) Get(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	return b.store.get(key)
}

// store contains the committed values of a book. Each of its calls asks its
// dependency first, and a call that the dependency fails changes nothing.
type store struct {
	dependency *conform.Dependency
	values     map[string]string
	mu         sync.Mutex
}

// ask returns the error of the call of the callable of the dependency of the
// store, and nil for a store without a dependency.
func (s *store) ask(callable string) error {
	if s.dependency == nil {
		return nil
	}
	if err := s.dependency.Call(callable); err != nil {
		return fmt.Errorf("ledgertest: %s: %w", callable, err)
	}
	return nil
}

// commit writes each value of writes under its key in one call.
func (s *store) commit(writes map[string]string) error {
	if err := s.ask("Commit"); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	maps.Copy(s.values, writes)
	return nil
}

// get returns the value under the key, and ledger.ErrMissing for a key
// without one.
func (s *store) get(key string) (string, error) {
	if err := s.ask("Get"); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	value, found := s.values[key]
	if !found {
		return "", ledger.ErrMissing
	}
	return value, nil
}

// BookOptions returns the door of a book over a store whose calls fail on
// demand.
func BookOptions() []LedgerOption {
	return []LedgerOption{conform.Faultable(faulty)}
}

// faulty returns a book over an empty store whose calls ask a dependency of
// the injector first.
func faulty(_ *LedgerHandle, inj *conform.Injector) ledger.Ledger {
	return &Book{store: &store{dependency: inj.Dependency(), values: map[string]string{}}}
}
