// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package poisonable is the mixin-axis fixture for the poisonable mixin, which
// declares that a subject driven into a failed state reports it consistently
// rather than intermittently.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression (docs/adr/0016).
package poisonable

import (
	"context"
)

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=poisonabletest/
type Mixed interface {
	// Fail drives the subject into the failed state the probe reports.
	Fail(ctx context.Context) error

	// Probe reports the state. It takes nothing, which is the shape a
	// stored-error accessor has: the answer is state rather than a lookup.
	//dokimi:mixin poisonable induce=Fail
	Probe() error
}
