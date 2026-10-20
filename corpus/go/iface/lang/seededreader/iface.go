// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package seededreader is the language-axis fixture for the SEED SEAM:
// an interface nothing can write to.
//
// A run populates a subject through the surface under test, and this
// surface offers no way in. A read against an empty subject cannot tell
// a miss from a bug, so the harness receives the corpus instead of a
// bare constructor — the suite decides what is seeded, from its own
// pools, and the harness decides how loading happens.
//
// Both roles are stamped because the corpus is zipped from them: a key
// pool with no payload pool names half a map.
package seededreader

import "context"

// Key identifies one document.
//
// A directive rather than a tag: both are named types, and a named type
// has no field to carry one.
//
//dokimi:default "test-doc"
type Key string

// Body is what a document holds.
//
//dokimi:default "test-contents"
type Body string

// Catalog reads documents and writes none.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=seededreadertest/
type Catalog interface {
	// Lookup answers for a seeded key and reports a miss for any other.
	Lookup(ctx context.Context, key Key) (Body, error)

	// Len reports how many documents were loaded.
	Len(ctx context.Context) (int, error)
}
