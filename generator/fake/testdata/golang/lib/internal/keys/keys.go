// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package keys is internal to the tree lib. The directive out moves the fake
// of Keys out of that tree, so the file of the fake cannot import Keys.
package keys

import "context"

// Value is a value of the store.
type Value struct {
	Key  string
	Body string
}

// Keys keeps values under their keys.
//
//dokimi:suite
//dokimi:fake
//dokimi:out path=../../../keystest/
type Keys interface {
	// Put stores v.
	//
	//dokimi:contract persister role=writer
	Put(ctx context.Context, v Value) (Value, error)

	// Get returns the value under key.
	//
	//dokimi:contract persister role=reader
	Get(ctx context.Context, key string) (Value, error)
}
