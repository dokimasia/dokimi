// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package retrysucceeds is the mixin-axis fixture for the retrysucceeds mixin, which
// declares that a transient failure succeeds on a later attempt.
//
// The interface carries whatever methods the law needs to be stateable.
// A mixin whose law spans two calls cannot be hosted by a single method:
// there would be nothing to compare against, and the generated subtest
// would pass by having nothing to check.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression (docs/adr/0016).
package retrysucceeds

import (
	"context"
)

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=retrysucceedstest/
type Mixed interface {
	// Call fails transiently before succeeding. The law drives it repeatedly,
	// so a method that either always works or never does cannot host it.
	//
	// The fault directive carries only a key — no sentinel — which is the
	// shape a retry-loop test wants: RetrySchedule takes whatever error the
	// test supplies, and forcing a named sentinel here would invent a helper
	// nobody asked for.
	//dokimi:mixin retry-succeeds attempts=3
	//dokimi:fault retry=3
	Call(ctx context.Context, key string) error

	// Attempts reports how many tries the subject has seen, which is how the
	// law distinguishes a retry from a first success.
	Attempts(ctx context.Context) (int, error)
}
