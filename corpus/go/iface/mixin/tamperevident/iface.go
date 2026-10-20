// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package tamperevident is the mixin-axis fixture for the tamperevident mixin, which
// declares that a stored value altered outside the interface is detected
// rather than served as if intact.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression (docs/adr/0016).
package tamperevident

import (
	"context"
)

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=tamperevidenttest/
type Mixed interface {
	// Store records a value under its integrity check.
	//dokimi:mixin tamper-evident tamper=Corrupt verify=Verify
	Store(ctx context.Context, body string) error

	// Corrupt alters the stored bytes behind the interface's back, which is
	// the only way to reach the state Verify exists to detect.
	Corrupt(ctx context.Context) error

	// Verify reports whether what is stored is what was stored.
	Verify(ctx context.Context) error
}
