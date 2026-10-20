// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package xsssafe is the mixin-axis fixture for the xsssafe mixin, which
// declares that a value rendered into markup cannot close a tag or open a
// script.
//
// One method is enough here: the claim is about what a single call returns
// for a hostile or unusual input, not about how two calls relate.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression (docs/adr/0016).
package xsssafe

import (
	"context"
)

// Unsafe is the value the directive names: markup that closes a tag and
// opens a script, which is what a renderer must not hand back.
//
// Declared rather than derived. A drawn string is "test-in" and nothing
// in it is dangerous, so a check over one passes a renderer that escapes
// nothing at all.
var Unsafe = `<script>alert(1)</script>`

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=xsssafetest/
type Mixed interface {
	// Render escapes a raw value for embedding in markup.
	//dokimi:mixin xss-safe unsafe=Unsafe
	Render(ctx context.Context, in string) (string, error)
}
