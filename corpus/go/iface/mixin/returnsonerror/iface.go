// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package returnsonerror is the mixin-axis fixture for the returns-on-error
// mixin, which declares that a failed call returns the value of a field of
// its request beside the error.
//
// Read and Fill append a record to the buffer of their request, and return
// the buffer unchanged with every error, so the caller keeps its buffer. A
// failed call does not return the zero value, so the mixin refuses the claim
// zero-on-error with the gate returns-no-field.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression.
package returnsonerror

import (
	"context"
	"errors"
)

// ErrNotFound is the error of a read of a key without a record.
var ErrNotFound = errors.New("returnsonerror: not found")

// Request is a read of the record of a key into a buffer.
type Request struct {
	Key string
	Dst []byte
}

// Result is the record of a read and its size.
type Result struct {
	Record []byte
	Size   int
}

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=returnsonerrortest/
type Mixed interface {
	// Put stores the record under key.
	Put(ctx context.Context, key string, record []byte) error

	// Read appends the record of req.Key to req.Dst, and returns the buffer
	// as the record of its result. It returns req.Dst as the record of its
	// result with every error.
	//
	//dokimi:mixin returns-on-error field=Dst result=Record
	Read(ctx context.Context, req *Request) (Result, error)

	// Fill appends the record of req.Key to req.Dst, and returns the buffer.
	// It returns req.Dst with every error.
	//
	//dokimi:mixin returns-on-error field=Dst
	Fill(ctx context.Context, req Request) ([]byte, error)
}
