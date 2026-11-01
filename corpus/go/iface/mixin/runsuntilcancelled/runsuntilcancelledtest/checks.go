// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package runsuntilcancelledtest contains the suite of Worker, a check that a
// consumer adds to the suite for the end of the loop, and the subject of the
// suite.
package runsuntilcancelledtest

import (
	"context"
	"testing"
	"testing/synctest"

	"go.dokimi.dev/assert"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/runsuntilcancelled"
)

// StopsWhenItsContextEnds checks that Run returns nil once its context ends.
//
//dokimi:check runsuntilcancelled.Worker.Run defect=runsPastItsContext around=bubble
func StopsWhenItsContextEnds(h *WorkerHandle) {
	ctx, cancel := context.WithCancel(h.Context())
	ended := make(chan error, 1)
	go func() { ended <- h.Subject().Run(ctx) }()
	synctest.Wait()
	cancel()
	synctest.Wait()
	stopped := false
	var err error
	select {
	case err = <-ended:
		stopped = true
	default:
	}
	assert.True(h, stopped && err == nil, h.Statement())
}

// pastItsContext is a worker whose Run runs on after its context ends.
type pastItsContext struct {
	*InMemory
}

// Run runs the loop of the worker with a context that never ends.
func (w pastItsContext) Run(ctx context.Context) error {
	return w.InMemory.Run(context.WithoutCancel(ctx))
}

// runsPastItsContext returns a worker whose Run does not stop when its
// context ends.
func runsPastItsContext(*WorkerHandle) runsuntilcancelled.Worker {
	return pastItsContext{InMemory: NewInMemory()}
}

// bubble runs each run of a check in a testing/synctest bubble of its own.
func bubble(t *testing.T, run func(*testing.T)) {
	t.Helper()
	synctest.Test(t, run)
}
