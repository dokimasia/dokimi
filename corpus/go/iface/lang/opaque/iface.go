// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package opaque is the fixture of a request with a field whose type has no
// exported field. The property engine reads no field of such a type, so the
// suite starts each value of it from its zero value.
package opaque

import "context"

// Set is a set of numbers. Its fields are unexported, and its zero value is
// the empty set.
type Set struct {
	items []int
}

// Len returns the number of items in s.
func (s Set) Len() int { return len(s.items) }

// PutRequest is the request of Put.
type PutRequest struct {
	Key string
	Set Set
}

// Store keeps sets.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=opaquetest/
type Store interface {
	// Put stores the set of the request under its key.
	//
	//dokimi:budget allocs=0
	Put(ctx context.Context, req *PutRequest) error
}
