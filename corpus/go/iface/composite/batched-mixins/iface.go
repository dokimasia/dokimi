// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package batchedmixins is the fixture for several mixins on one line.
//
// A comment line can contain several directives. Each directive after the
// first opens with the mark dokimi: after a blank. A param belongs to the
// directive that it follows. Put and List declare three mixins on one line
// each.
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
	//dokimi:mixin idempotent dokimi:mixin concurrent dokimi:mixin side-effect observe=List
	Put(ctx context.Context, key, value string) error

	// Read has one mixin with a param.
	//dokimi:mixin read-after-write write=Put
	Read(ctx context.Context, key string) (string, error)

	// List has two mixins without params and one mixin with a param.
	//dokimi:mixin cacheable dokimi:mixin pure dokimi:mixin bounded limit=50
	List(ctx context.Context) ([]string, error)
}
