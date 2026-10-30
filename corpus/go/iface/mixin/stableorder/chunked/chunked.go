// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package chunked contains a control of the suite of Mixed: a list that grows
// by blocks of a fixed size.
package chunked

import (
	"cmp"
	"context"
	"slices"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/stableorder"
)

// blockSize is the number of values of each block.
const blockSize = 4

// List implements [stableorder.Mixed] over blocks of at most blockSize values.
// Each block before the last is full. Items returns the values sorted by key,
// in the order of their adds within one key. The corpus applies the suite of
// Mixed to it as a control of the appenders and logs.
//
//dokimi:subject stableorder.Mixed new=New
type List struct {
	blocks [][]stableorder.Value
	mu     sync.Mutex
}

// The compiler checks that List implements [stableorder.Mixed].
var _ stableorder.Mixed = (*List)(nil)

// New returns an empty List.
func New() *List { return &List{} }

// Add implements [stableorder.Mixed].
func (l *List) Add(ctx context.Context, v stableorder.Value) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	if n := len(l.blocks); n == 0 || len(l.blocks[n-1]) == blockSize {
		l.blocks = append(l.blocks, make([]stableorder.Value, 0, blockSize))
	}
	last := len(l.blocks) - 1
	l.blocks[last] = append(l.blocks[last], v)
	return nil
}

// Items implements [stableorder.Mixed].
func (l *List) Items(ctx context.Context) ([]stableorder.Value, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	out := slices.Concat(l.blocks...)
	slices.SortStableFunc(out, func(a, b stableorder.Value) int { return cmp.Compare(a.Key, b.Key) })
	return out, nil
}
