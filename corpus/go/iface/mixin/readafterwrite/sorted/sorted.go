// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package sorted contains a control of the suite of Mixed: a store that keeps
// its keys in a sorted slice and finds a key by binary search.
package sorted

import (
	"context"
	"errors"
	"slices"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/readafterwrite"
)

// ErrNotFound is the error of a read of a key without a value.
var ErrNotFound = errors.New("sorted: not found")

// Store implements [readafterwrite.Mixed] over two slices. keys is sorted,
// and values[i] is the value of keys[i]. The corpus applies the suite of Mixed
// to it as a control of the stores and readers.
//
//dokimi:subject readafterwrite.Mixed new=New
type Store struct {
	keys   []string
	values []string
	mu     sync.Mutex
}

// The compiler checks that Store implements [readafterwrite.Mixed].
var _ readafterwrite.Mixed = (*Store)(nil)

// New returns an empty Store.
func New() *Store { return &Store{} }

// Write implements [readafterwrite.Mixed].
func (s *Store) Write(ctx context.Context, key, value string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	at, found := slices.BinarySearch(s.keys, key)
	if found {
		s.values[at] = value
		return nil
	}
	s.keys = slices.Insert(s.keys, at, key)
	s.values = slices.Insert(s.values, at, value)
	return nil
}

// Read implements [readafterwrite.Mixed].
func (s *Store) Read(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	at, found := slices.BinarySearch(s.keys, key)
	if !found {
		return "", ErrNotFound
	}
	return s.values[at], nil
}
