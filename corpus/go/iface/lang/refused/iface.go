// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package refused is the language-axis fixture for the directive refuse,
// which withdraws a derived claim of a callable for every subject of its
// suite, with a reason.
//
// Update fails only for an instant beyond the max offset of the clock, and no
// sample of an instant is beyond it. A check of the claim zero-on-error on the
// samples would never engage, so Update refuses the claim.
package refused

import "errors"

// ErrOffset is the error of Update for an instant beyond the max offset of the
// clock.
var ErrOffset = errors.New("refused: the instant is beyond the max offset")

// Instant is an instant of a hybrid logical clock: a wall time in nanoseconds
// since the Unix epoch, and a logical count of the events at that wall time.
type Instant struct {
	Wall    int64
	Logical int32
}

// Clock is a hybrid logical clock.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=refusedtest/
type Clock interface {
	// Update merges the instant of a peer into the clock, and returns the
	// instant of the clock after the merge. It returns ErrOffset for an
	// instant beyond the max offset of the clock.
	//
	//dokimi:refuse zero-on-error reason="Update fails only for an instant beyond the max offset of the clock"
	Update(observed Instant) (Instant, error)
}
