// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package clash declares a store whose parameters take the names of the
// receiver and of the locals of a fake.
package clash

import "context"

// Value is a value of the store.
type Value struct {
	Present string
	Body    string
}

// Shadow keeps values under their keys, with parameters named f, r0 and
// present.
//
//dokimi:suite
//dokimi:fake
//dokimi:out path=clashtest/
type Shadow interface {
	// Put stores r0 and returns it as stored.
	//
	//dokimi:contract persister role=writer
	Put(f context.Context, r0 Value) (Value, error)

	// Get returns the value under present.
	//
	//dokimi:contract persister role=reader
	Get(f context.Context, present string) (Value, error)
}
