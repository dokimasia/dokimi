// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package meter declares a generic cache whose witness is a type of the
// internal package keys. The directive out moves the fake of Meter out of the
// tree lib, so the files of the fake cannot refer to the witness.
package meter

import (
	"context"

	"example.com/shop/lib/internal/keys"
)

// Reading is the value of a meter, a value of the internal package keys.
type Reading = keys.Value

// Meter reads values of V through a cache.
//
//dokimi:suite
//dokimi:fake
//dokimi:witness V=keys.Value
//dokimi:out path=../../metertest/
type Meter[V any] interface {
	// Store fills the backing with v under key.
	Store(ctx context.Context, key string, v V) error

	// Lookup returns the cached value under key.
	//
	//dokimi:contract cache role=cache
	Lookup(ctx context.Context, key string) (V, error)

	// Fetch returns the value under key from the backing.
	//
	//dokimi:contract cache role=backing
	Fetch(ctx context.Context, key string) (V, error)
}
