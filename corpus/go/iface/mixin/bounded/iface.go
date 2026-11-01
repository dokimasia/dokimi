// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package bounded is the fixture of the mixin bounded, which declares that
// the use of a callable remains between a floor and a ceiling.
//
// The check reads the use through the callable that the param observe names
// after each call of a sequence. The interface therefore declares the call
// that grows the use and the call that reads it.
//
// The fixture has no negated form. A mixin is opt-in, so eidos refuses the
// negation of its directive.
package bounded

import (
	"context"
)

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=boundedtest/
type Mixed interface {
	// Add adds an item to the collection. Its use is the length of the list
	// that List returns, and the limit is five items, so a sequence of six
	// calls of Add crosses it.
	//dokimi:mixin bounded limit=5 observe=List
	Add(ctx context.Context, item string) error

	// List returns the items of the collection.
	List(ctx context.Context) ([]string, error)
}
