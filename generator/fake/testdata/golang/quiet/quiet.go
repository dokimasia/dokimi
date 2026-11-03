// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package quiet keeps values with a writer that returns nothing, the fixture
// of a store that takes a context and returns no error.
package quiet

import "context"

// Value is a value of the store.
type Value struct {
	Key  string
	Body string
}

// Quiet keeps values under their keys, and its writer returns nothing.
//
//dokimi:suite
//dokimi:fake
//dokimi:out path=quiettest/
type Quiet interface {
	// Put stores v.
	//
	//dokimi:contract persister role=writer
	Put(ctx context.Context, v Value)

	// Get returns the value under key.
	//
	//dokimi:contract persister role=reader
	Get(ctx context.Context, key string) (Value, error)
}
