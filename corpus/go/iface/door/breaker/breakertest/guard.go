// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package breakertest contains the suite of Breaker and the subject of the
// suite.
package breakertest

import (
	"context"
	"sync"

	"go.dokimi.dev/conform"

	"go.dokimi.dev/dokimi/corpus/go/iface/door/breaker"
)

// threshold is the number of failures of the operation in a row that open a
// guard.
const threshold = 3

// Guard is the Breaker over an operation. The corpus applies the suite of
// Breaker to it as a subject, over an operation whose calls fail on demand.
//
//dokimi:subject breaker.Breaker new=NewGuard options=GuardOptions
type Guard struct {
	operation func(ctx context.Context) error
	failures  int
	mu        sync.Mutex
}

// The compiler checks that Guard implements [breaker.Breaker].
var _ breaker.Breaker = (*Guard)(nil)

// NewGuard returns a closed guard over an operation that always succeeds.
func NewGuard() *Guard {
	return &Guard{operation: func(context.Context) error { return nil }}
}

// Do implements [breaker.Breaker].
func (g *Guard) Do(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.failures >= threshold {
		return breaker.ErrOpen
	}
	if err := g.operation(ctx); err != nil {
		g.failures++
		return err
	}
	g.failures = 0
	return nil
}

// GuardOptions returns the door of a guard over an operation whose calls fail
// on demand.
func GuardOptions() []BreakerOption {
	return []BreakerOption{conform.Faultable(faulty)}
}

// faulty returns a closed guard over an operation that asks a dependency of
// the injector.
func faulty(_ *BreakerHandle, inj *conform.Injector) breaker.Breaker {
	operation := inj.Dependency()
	return &Guard{operation: func(context.Context) error { return operation.Call("Do") }}
}
