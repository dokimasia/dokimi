// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package nestedtest contains the suite of Dispatcher. InMemory is the subject
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
// the suite of Dispatcher to it as a subject.
//
//dokimi:subject nested.Dispatcher new=NewInMemory
type InMemory struct {
	weight int
	mu     sync.Mutex
}

// NewInMemory returns an InMemory without dispatches.
func NewInMemory() *InMemory { return &InMemory{} }

// Dispatch implements [nested.Dispatcher]. It returns the error of a
// cancelled context.
func (d *InMemory) Dispatch(ctx context.Context, req *nested.DispatchRequest) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.dispatch(req)
}

// DispatchAll implements [nested.Dispatcher]. It returns the error of a
// cancelled context.
func (d *InMemory) DispatchAll(ctx context.Context, b nested.Batch) error {
	if err := ctx.Err(); err != nil {
		return err
	}
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
// request. It returns [nested.ErrUndelivered] for a request without a parcel,
// a Done or a Notify, and where Done or Notify returns an error.
func (d *InMemory) dispatch(req *nested.DispatchRequest) error {
	if req.Parcel == nil || req.Done == nil || req.Notify == nil {
		return nested.ErrUndelivered
	}
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
