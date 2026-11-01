// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package referencetest contains the subjects of the declarations of the
// fixture reference and their reference implementations. Each subject keeps
// its state in another representation than its reference.
package referencetest

import (
	"context"
	"errors"
	"slices"
	"sync"

	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/reference"
)

// pair is one value of a PairStore under its key.
type pair struct {
	key, value string
}

// PairStore implements [reference.Store] over a slice of pairs. A Put of a
// stored key replaces the value of its pair, and the revision of a write is
// the count of the writes of the store. The corpus applies the suite of Store
// to it with a map as its reference.
//
//dokimi:subject reference.Store new=NewPairStore options=PairStoreOptions
type PairStore struct {
	pairs    []pair
	revision int64
	mu       sync.Mutex
}

// The compiler checks that PairStore implements [reference.Store].
var _ reference.Store = (*PairStore)(nil)

// NewPairStore returns an empty PairStore.
func NewPairStore() *PairStore { return &PairStore{} }

// PairStoreOptions returns the reference of a PairStore: a store over a map.
func PairStoreOptions() []StoreOption {
	return []StoreOption{conform.Reference(func(*StoreHandle) reference.Store {
		return &mapStore{values: map[string]string{}}
	})}
}

// Put implements [reference.Store].
func (s *PairStore) Put(ctx context.Context, key, value string) (int64, error) {
	if err := contextErr(ctx); err != nil {
		return 0, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.revision++
	if i := s.index(key); i >= 0 {
		s.pairs[i].value = value
		return s.revision, nil
	}
	s.pairs = append(s.pairs, pair{key: key, value: value})
	return s.revision, nil
}

// Get implements [reference.Store].
func (s *PairStore) Get(ctx context.Context, key string) (string, error) {
	if err := contextErr(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	i := s.index(key)
	if i < 0 {
		return "", reference.ErrNotFound
	}
	return s.pairs[i].value, nil
}

// Delete implements [reference.Store].
func (s *PairStore) Delete(ctx context.Context, key string) error {
	if err := contextErr(ctx); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if i := s.index(key); i >= 0 {
		s.pairs = slices.Delete(s.pairs, i, i+1)
	}
	return nil
}

// index returns the position of the pair of key, and -1 for a key without a
// pair. The caller holds the mutex.
func (s *PairStore) index(key string) int {
	return slices.IndexFunc(s.pairs, func(p pair) bool { return p.key == key })
}

// mapStore is the reference of a PairStore: the values in a map. The revision
// of a write is the count of the keys with a value after it, which differs
// from the revision of a PairStore.
type mapStore struct {
	values map[string]string
}

// Put implements [reference.Store].
func (s *mapStore) Put(_ context.Context, key, value string) (int64, error) {
	s.values[key] = value
	return int64(len(s.values)), nil
}

// Get implements [reference.Store].
func (s *mapStore) Get(_ context.Context, key string) (string, error) {
	value, stored := s.values[key]
	if !stored {
		return "", reference.ErrNotFound
	}
	return value, nil
}

// Delete implements [reference.Store].
func (s *mapStore) Delete(_ context.Context, key string) error {
	delete(s.values, key)
	return nil
}

// contextErr returns the error of a canceled or expired context, and an error
// for a nil context.
func contextErr(ctx context.Context) error {
	if ctx == nil {
		return errors.New("referencetest: nil context")
	}
	return ctx.Err()
}
