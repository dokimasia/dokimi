// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package batchedmixins is the fixture for several mixins on one method.
//
// The mixin directive takes one variant, so a method with several mixins has
// one directive line for each. Put and List stack mixins without params on a
// mixin with params.
package batchedmixins

import "context"

// Batched is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=batchedmixinstest/
type Batched interface {
	// Put has three mixins. side-effect observes the write through List,
	// which returns a value before the write and after it. Read refuses a key
	// that no write stored, so it cannot observe the change.
	//dokimi:mixin idempotent
	//dokimi:mixin concurrent
	//dokimi:mixin side-effect observe=List
	Put(ctx context.Context, key, value string) error

	// Read has one mixin with a param.
	//dokimi:mixin read-after-write write=Put
	Read(ctx context.Context, key string) (string, error)

	// List has two mixins without params and one mixin with a param.
	//dokimi:mixin cacheable
	//dokimi:mixin pure
	//dokimi:mixin bounded limit=50
	List(ctx context.Context) ([]string, error)
}
