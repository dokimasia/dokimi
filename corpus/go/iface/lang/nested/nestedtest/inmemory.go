// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package nestedtest contains the suite of Dispatcher. InMemory is the control
// of the suite.
package nestedtest

import (
	"context"
	"sync"

	"go.dokimi.dev/dokimi/corpus/go/iface/lang/nested"
)

// InMemory implements [nested.Dispatcher] in memory. Each dispatch weighs the
// parcel of its request and calls Done and Notify of the request, so it reads
// each field of a function type and of an interface type. The corpus applies
// the suite of Dispatcher to it as a control.
//
//dokimi:subject nested.Dispatcher new=NewInMemory
type InMemory struct {
	weight int
	mu     sync.Mutex
}

// NewInMemory returns an InMemory without dispatches.
func NewInMemory() *InMemory { return &InMemory{} }

// Dispatch implements [nested.Dispatcher].
func (d *InMemory) Dispatch(_ context.Context, req *nested.DispatchRequest) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.dispatch(req)
}

// DispatchAll implements [nested.Dispatcher].
func (d *InMemory) DispatchAll(_ context.Context, b nested.Batch) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	for i := range b.Requests {
		if err := d.dispatch(&b.Requests[i]); err != nil {
			return err
		}
	}
	return nil
}

// Weight implements [nested.Dispatcher].
func (d *InMemory) Weight() int {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.weight
}

// dispatch weighs the parcel of the request, and calls Done and Notify of the
// request. It returns [nested.ErrUndelivered] where either returns an error.
func (d *InMemory) dispatch(req *nested.DispatchRequest) error {
	d.weight += req.Parcel.Weight()
	if req.Done() != nil {
		return nested.ErrUndelivered
	}
	if req.Notify(req.ID) != nil {
		return nested.ErrUndelivered
	}
	return nil
}

// The compiler checks that InMemory implements [nested.Dispatcher].
var _ nested.Dispatcher = (*InMemory)(nil)
