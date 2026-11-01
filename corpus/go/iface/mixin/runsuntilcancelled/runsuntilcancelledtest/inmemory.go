// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package runsuntilcancelledtest

import (
	"context"
	"sync/atomic"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/runsuntilcancelled"
)

// InMemory implements [runsuntilcancelled.Worker] in memory. Its loop waits
// for the end of the context of Run. The corpus applies the suite of Worker to
// it as a subject.
//
//dokimi:subject runsuntilcancelled.Worker new=NewInMemory
type InMemory struct {
	running atomic.Int32
}

// The compiler checks that InMemory implements [runsuntilcancelled.Worker].
var _ runsuntilcancelled.Worker = (*InMemory)(nil)

// NewInMemory returns a worker whose loop does not run.
func NewInMemory() *InMemory { return &InMemory{} }

// Run implements [runsuntilcancelled.Worker]. It runs the loop until ctx
// ends, and returns nil.
func (w *InMemory) Run(ctx context.Context) error {
	w.running.Add(1)
	defer w.running.Add(-1)
	<-ctx.Done()
	return nil
}

// Running implements [runsuntilcancelled.Worker].
func (w *InMemory) Running() bool { return w.running.Load() > 0 }
