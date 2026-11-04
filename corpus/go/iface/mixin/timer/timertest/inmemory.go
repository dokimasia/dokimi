// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package timertest contains the double and the suite of Mixed. InMemory is
// the subject of the suite.
package timertest

import (
	"sync/atomic"

	"go.dokimi.dev/dokimi/corpus/go/iface/mixin/timer"
)

// InMemory implements [timer.Mixed] with a flag of an active timer. The
// corpus applies the suite of Mixed to it as a subject.
//
//dokimi:subject timer.Mixed new=NewInMemory
type InMemory struct {
	active atomic.Bool
}

// The compiler checks that InMemory implements [timer.Mixed].
var _ timer.Mixed = (*InMemory)(nil)

// NewInMemory returns an active timer.
func NewInMemory() *InMemory {
	t := &InMemory{}
	t.active.Store(true)
	return t
}

// Stop implements [timer.Mixed].
func (t *InMemory) Stop() bool { return t.active.Swap(false) }
