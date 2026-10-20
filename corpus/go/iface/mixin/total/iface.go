// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package total is the mixin-axis fixture for the total mixin, which
// declares that a callable answers for every input in a named domain rather
// than failing on some of them.
//
// Two methods, differing only in whether they can report a failure. The claim
// is about what a single call returns for a hostile input rather than about
// how two calls relate, so one would state it — but the law's closure threads
// an error, and a method with no error to thread is the shape that made the
// binder emit a call the compiler refuses. A total function that cannot fail
// is the strongest form of the claim, and the fixture the tier had none of.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression (docs/adr/0016).
package total

import (
	"context"
)

// Edge is the input the directive names: the empty string, which is where
// the domain of strings ends.
//
// Declared rather than derived. A domain is prose for a reader — "strings"
// says nothing a generator can build — and a drawn sample sits comfortably
// inside it, so a subject that gives way at the boundary passes on one.
var Edge = ""

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=totaltest/
type Mixed interface {
	// Classify answers for any input in the declared domain.
	//dokimi:mixin total domain=strings edge=Edge
	Classify(ctx context.Context, in string) (string, error)

	// Normalize answers for any input and has no way to report otherwise.
	//dokimi:mixin total domain=strings edge=Edge
	Normalize(in string) string
}
