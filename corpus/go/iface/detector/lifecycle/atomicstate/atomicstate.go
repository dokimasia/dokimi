// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package atomicstate contains a control of the suite of Lifecycle: a
// lifecycle whose state is an atomic integer without a mutex.
package atomicstate

import (
	"context"
	"sync/atomic"

	"go.dokimi.dev/dokimi/corpus/go/iface/detector/lifecycle"
)

// The states of a Lifecycle.
const (
	open int32 = iota
	closed
)

// Lifecycle implements [lifecycle.Lifecycle] with its state in an atomic
// integer. The corpus applies the suite of Lifecycle to it as a control of
// the flags and lifecycles.
//
//dokimi:subject lifecycle.Lifecycle new=New
type Lifecycle struct {
	state atomic.Int32
}

// The compiler checks that Lifecycle implements [lifecycle.Lifecycle].
var _ lifecycle.Lifecycle = (*Lifecycle)(nil)

// New returns an open Lifecycle.
func New() *Lifecycle { return &Lifecycle{} }

// Close implements [lifecycle.Lifecycle]. A second Close succeeds.
func (l *Lifecycle) Close(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	l.state.CompareAndSwap(open, closed)
	return nil
}
