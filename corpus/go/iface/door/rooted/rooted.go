// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package rooted declares a tally whose kit has no implementation of its
// own. The composition root of the program builds the tally, above the kit,
// so the test of the suite takes the reference of its planted defects from
// the options of its falsification.
package rooted

import "context"

// Tally counts the calls of Inc.
//
//dokimi:suite falsify=falsifyOptions
//dokimi:out path=rootedtest/
type Tally interface {
	// Inc adds one to the count.
	Inc(ctx context.Context) error

	// Count returns the count.
	Count(ctx context.Context) (int, error)
}
