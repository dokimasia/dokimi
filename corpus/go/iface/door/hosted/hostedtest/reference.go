// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package hostedtest

import (
	"context"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/hosted"
)

// reference is the reference of Store: the values by key in memory, which no
// crash loses.
type reference struct {
	values map[string]string
	mu     sync.Mutex
}

// newReference returns an empty reference.
func newReference() *reference { return &reference{values: map[string]string{}} }

// Put stores the value under the key.
func (r *reference) Put(ctx context.Context, key, value string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.values[key] = value
	return nil
}

// Get returns the value under the key, and hosted.ErrNotFound for a key
// without a value.
func (r *reference) Get(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	value, found := r.values[key]
	if !found {
		return "", hosted.ErrNotFound
	}
	return value, nil
}
