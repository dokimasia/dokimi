// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package poisonaccessor is the detector-axis fixture for the poisonaccessor shape:
// a latched failure state. Once poisoned every later call reports the same
// error, which is what makes the accessor worth having: the failure outlives the call that caused it.
//
// No directive appears here: the classification comes from the signature
// alone, so a detector that misfires shows up as wrong generated output
// rather than as a directive that was misread.
package poisonaccessor

import (
	"errors"
)

// ErrPoisoned is the latched failure this fixture reports once tripped.
var ErrPoisoned = errors.New("poisonaccessor: poisoned")

// PoisonAccessor is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=poisonaccessortest/
type PoisonAccessor interface {
	Err() error
}
