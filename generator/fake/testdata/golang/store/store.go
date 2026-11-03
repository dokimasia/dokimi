// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package store keeps values under their keys, the fixture of a fake of the
// contract persister.
package store

import "context"

// Value is a value of the store.
type Value struct {
	Key  string
	Body string
}

// Store keeps values under their keys.
//
//dokimi:suite
//dokimi:fake
//dokimi:out path=storetest/
type Store interface {
	// Put stores v and returns it as stored.
	//
	//dokimi:contract persister role=writer
	Put(ctx context.Context, v Value) (Value, error)

	// Get returns the value under key.
	//
	//dokimi:contract persister role=reader
	Get(ctx context.Context, key string) (Value, error)
}
