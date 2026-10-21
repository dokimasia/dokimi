// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package batchedmixinstest

import (
	"context"
	"errors"
	"maps"
	"slices"
	"sync"

	batchedmixins "go.dokimi.dev/dokimi/corpus/go/iface/composite/batched-mixins"
)

// ErrNotFound reports a key that no write stored.
var ErrNotFound = errors.New("batchedmixinstest: no write stored the key")

// listLimit is the most keys that List returns, the limit of the mixin
// bounded of List.
const listLimit = 50

// InMemory is the Batched in memory.
//
//dokimi:subject batchedmixins.Batched new=NewInMemory
type InMemory struct {
	values map[string]string
	mu     sync.Mutex
}

// NewInMemory returns an empty InMemory.
func NewInMemory() *InMemory {
	return &InMemory{values: map[string]string{}}
}

// Put stores the value under the key. A repeat of a write stores the same
// value again.
func (m *InMemory) Put(ctx context.Context, key, value string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.values[key] = value
	return nil
}

// Read returns the value under the key. It returns ErrNotFound for a key that
// no write stored.
func (m *InMemory) Read(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	value, found := m.values[key]
	if !found {
		return "", ErrNotFound
	}
	return value, nil
}

// List returns the stored keys in ascending order, at most listLimit of them.
func (m *InMemory) List(ctx context.Context) ([]string, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	keys := slices.Sorted(maps.Keys(m.values))
	return keys[:min(len(keys), listLimit)], nil
}

// Compile-time proof that InMemory satisfies Batched.
var _ batchedmixins.Batched = (*InMemory)(nil)
