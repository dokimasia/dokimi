// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package runsuntilcancelled is the mixin-axis fixture for the mixin
// runs-until-cancelled, which declares that a call runs until its context
// ends.
//
// Run runs the loop of a worker until its context ends. A check that calls
// Run with the context of the check never returns, and the end of the context
// completes Run without an error. The mixin keeps Run out of the checks of a
// single caller, of each generated benchmark and of each budget. The suite
// checks Running, and a check of the consumer ends the context of Run.
package runsuntilcancelled

import "context"

// Worker runs the loop of a worker.
//
//dokimi:suite
//dokimi:out path=runsuntilcancelledtest/
type Worker interface {
	// Run runs the loop until ctx ends, and returns nil when ctx ends.
	//dokimi:mixin runs-until-cancelled
	Run(ctx context.Context) error

	// Running reports whether a call of Run runs the loop.
	Running() bool
}
