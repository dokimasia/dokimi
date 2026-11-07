// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package app is the composition root of the tally: the package above the
// kit rootedtest that builds the tally of the program.
package app

import (
	"context"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/rooted"
)

// Tally is the tally of the program, under a mutex.
//
//dokimi:subject rooted.Tally new=NewTally options=tallyOptions
type Tally struct {
	mu    sync.Mutex
	count int
}

// The compiler checks that Tally implements [rooted.Tally].
var _ rooted.Tally = (*Tally)(nil)

// NewTally returns a tally at zero.
func NewTally() *Tally {
	return &Tally{}
}

// Inc adds one to the count.
func (t *Tally) Inc(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	t.count++
	return nil
}

// Count returns the count.
func (t *Tally) Count(ctx context.Context) (int, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.count, nil
}
