// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package serializabletest contains the double and the suite of Mixed. InMemory
// is the control of the suite.
package serializabletest

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/serializable"
)

// ErrNotFound is the error of Get for a key that no recorded entry has.
var ErrNotFound = errors.New("serializabletest: not found")

// ErrNotAdmissible is the error of Record for an entry that the isolation
// level refuses.
var ErrNotAdmissible = errors.New("serializabletest: the isolation level refuses the entry")

// InMemory implements [serializable.Mixed] in memory. The corpus applies the
// suite of Mixed to it as a control.
//
//dokimi:subject serializable.Mixed new=NewInMemory
type InMemory struct {
	mu      sync.Mutex
	entries []serializable.Entry
	txn     int
	version map[string]int64
	wrote   map[keyed]bool
}

// The compiler checks that InMemory implements [serializable.Mixed].
var _ serializable.Mixed = (*InMemory)(nil)

// NewInMemory returns a new InMemory.
func NewInMemory() *InMemory {
	return &InMemory{version: map[string]int64{}, wrote: map[keyed]bool{}}
}

// Record implements [serializable.Mixed].
func (s *InMemory) Record(ctx context.Context, e serializable.Entry) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.admits(e); err != nil {
		return err
	}
	if e.Txn > s.txn {
		s.txn = e.Txn
	}
	if e.Write {
		s.version[e.Key] = e.Version
		s.wrote[keyed{txn: e.Txn, key: e.Key}] = true
	}
	s.entries = append(s.entries, e)
	return nil
}

// History implements [serializable.Mixed].
func (s *InMemory) History(ctx context.Context) ([]serializable.Entry, error) {
	if err := contextErr(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return slices.Clone(s.entries), nil
}

// Get implements [serializable.Mixed].
func (s *InMemory) Get(ctx context.Context, key string) (serializable.Entry, error) {
	if err := contextErr(ctx); err != nil {
		return serializable.Entry{}, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, e := range slices.Backward(s.entries) {
		if e.Key == key {
			return e, nil
		}
	}
	return serializable.Entry{}, ErrNotFound
}

// keyed is the pair of a transaction and a key that the transaction wrote.
type keyed struct {
	txn int
	key string
}

// admits returns an error that wraps ErrNotAdmissible when the isolation level
// refuses e, and nil otherwise.
func (s *InMemory) admits(e serializable.Entry) error {
	switch {
	case e.Txn <= 0:
		return fmt.Errorf("%w: transaction %d", ErrNotAdmissible, e.Txn)
	case e.Version <= 0:
		return fmt.Errorf("%w: version %d", ErrNotAdmissible, e.Version)
	case e.Txn < s.txn:
		return fmt.Errorf("%w: transaction %d closed at %d", ErrNotAdmissible, e.Txn, s.txn)
	case e.Write && e.Version <= s.version[e.Key]:
		return fmt.Errorf("%w: %q is at version %d", ErrNotAdmissible, e.Key, s.version[e.Key])
	case e.Write && s.wrote[keyed{txn: e.Txn, key: e.Key}]:
		return fmt.Errorf("%w: transaction %d already wrote %q", ErrNotAdmissible, e.Txn, e.Key)
	case !e.Write && e.Version != s.version[e.Key]:
		return fmt.Errorf("%w: %q is at version %d, not %d",
			ErrNotAdmissible, e.Key, s.version[e.Key], e.Version)
	}
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("serializabletest: nil context")
	}
	return ctx.Err()
}
