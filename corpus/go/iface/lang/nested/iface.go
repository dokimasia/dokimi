// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package nested is the language-axis fixture for inputs whose fields have
// function types and interface types, at each depth of a request.
package nested

import (
	"context"
	"errors"
)

// ErrUndelivered reports a dispatch whose request did not take its parcel,
// because Done or Notify of the request returned an error.
var ErrUndelivered = errors.New("nested: the request did not take the parcel")

// Parcel is the parcel of a dispatch.
type Parcel interface {
	// Weight returns the weight of the parcel in grams.
	Weight() int
}

// Notify notifies the sender of the parcel with the identity id.
type Notify func(id string) error

// DispatchRequest is the request of a dispatch. Parcel has an interface type,
// Done a function type and Notify a defined function type.
type DispatchRequest struct {
	Parcel Parcel
	Done   func() error
	Notify Notify
	ID     string
}

// Batch is a batch of dispatches.
type Batch struct {
	Requests []DispatchRequest
}

// Dispatcher dispatches the parcels of requests.
//
//dokimi:suite
//dokimi:out path=nestedtest/
type Dispatcher interface {
	// Dispatch dispatches the parcel of the request. It returns
	// ErrUndelivered where the request does not take the parcel.
	//
	//dokimi:budget allocs=0
	Dispatch(ctx context.Context, req *DispatchRequest) error
	// DispatchAll dispatches the parcel of each request of the batch, and
	// stops at the first request that returns ErrUndelivered.
	//
	//dokimi:budget allocs=0
	DispatchAll(ctx context.Context, b Batch) error
	// Weight returns the total weight of the dispatched parcels in grams.
	Weight() int
}
