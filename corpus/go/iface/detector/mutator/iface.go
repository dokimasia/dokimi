// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package mutator is the detector-axis fixture for the mutator shape:
// state change with nothing returned. With no error there is no in-band
// failure channel, so the only observable effect is on a later read.
//
// No directive appears here: the classification comes from the signature
// alone, so a detector that misfires shows up as wrong generated output
// rather than as a directive that was misread.
package mutator

import (
	"context"
)

// Mutator is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=mutatortest/
type Mutator interface {
	Touch(ctx context.Context, key string)
}
