// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package reader is the detector-axis fixture for the reader shape:
// a single-key fetch that can fail — the commonest shape in any
// store-like interface, and what every miss-sentinel law attaches to.
//
// No directive appears here: the classification comes from the signature
// alone, so a detector that misfires shows up as wrong generated output
// rather than as a directive that was misread.
package reader

import (
	"context"
	"errors"
)

// ErrNotFound is the miss sentinel this fixture's reads report.
var ErrNotFound = errors.New("reader: not found")

// Value is the payload the fixture reads.
type Value struct{ Key, Body string }

// Reader is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=readertest/
type Reader interface {
	// The fault line carries all three settings at once — a sentinel, a
	// retry attempt, and the partition field — because the settings are
	// independent and a projection that dropped one when another was
	// present would still pass every fixture carrying them apart. The
	// double is routed into readertest, so the sentinel reference this
	// renders has to qualify.
	//dokimi:mixin not-found sentinel=ErrNotFound
	//dokimi:fault errors=[ErrNotFound] retry=3 partition=key
	Get(ctx context.Context, key string) (Value, error)
}
