// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package sortedkeys contains a control of the suite of Contract: a lease
// table that keeps the held keys in a sorted slice.
package sortedkeys

import (
	"context"
	"slices"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/contract/lease"
)

// Table implements [lease.Contract] over a sorted slice of the held keys. A
// lease ends with Release, or once the context of its Acquire ends, through
// context.AfterFunc. The corpus applies the suite of Contract to it as a
// control of the leases.
//
//dokimi:subject lease.Contract new=New
type Table struct {
	held []string
	mu   sync.Mutex
}

// The compiler checks that Table implements [lease.Contract].
var _ lease.Contract = (*Table)(nil)

// New returns a Table without leases.
func New() *Table { return &Table{} }

// Acquire implements [lease.Contract].
func (t *Table) Acquire(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	at, held := slices.BinarySearch(t.held, key)
	if held {
		return lease.ErrHeld
	}
	t.held = slices.Insert(t.held, at, key)
	context.AfterFunc(ctx, func() { t.drop(key) })
	return nil
}

// Release implements [lease.Contract].
func (t *Table) Release(ctx context.Context, key string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	t.drop(key)
	return nil
}

// drop ends the lease of key.
func (t *Table) drop(key string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if at, held := slices.BinarySearch(t.held, key); held {
		t.held = slices.Delete(t.held, at, at+1)
	}
}
