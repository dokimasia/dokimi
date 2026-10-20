// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package integrationonly is the mixin-axis fixture for the integrationonly mixin, which
// declares that the method needs external infrastructure and is skipped by default.
//
// The interface carries whatever methods the law needs to be stateable.
// A mixin whose law spans two calls cannot be hosted by a single method:
// there would be nothing to compare against, and the generated subtest
// would pass by having nothing to check.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression (docs/adr/0016).
package integrationonly

import (
	"context"
)

// Mixed is the fixture interface.
//
// Mixed is the fixture interface. The double carries Connect like any other
// method: the mixin is a law for the suite and model tiers about what a test
// runs against, not an instruction to the stub generator, and a double
// omitting the method would not satisfy the interface at all.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=integrationonlytest/
type Mixed interface {
	// Connect reaches outside the process, which is what the mixin
	// declares. What it does NOT do yet is act on that: no rule derives a
	// check from it and nothing gates the generated file behind a build
	// tag, so the four signature checks below run wherever the package
	// does. The generated header carries the gap by name.
	//
	// The stamp is here anyway, and that is the fixture's whole point: a
	// classification eidos registers and testkit has no rule for still
	// has to reach the census, or it reads as covered by being absent
	// from every count.
	//dokimi:mixin integration-only
	Connect(ctx context.Context, dsn string) error
}
