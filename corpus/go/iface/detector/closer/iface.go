// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package closer is the detector-axis fixture for the closer shape:
// the bare teardown — a callable taking nothing and returning only an
// error, named for release.
//
// This is the catalogue's one name-gated detector, and the fixture
// holds that: `Close() error` and `Err() error` are the same signature
// with opposite semantics, so the name is what separates a teardown
// from a poison probe. A teardown is expected to be allowed to answer
// differently the second time, which is why no repeat claim is stated
// at the floor.
//
// No directive appears here: the classification comes from the
// signature and the name alone, so a detector that misfires shows up
// as wrong generated output rather than as a directive that was
// misread.
package closer

// Closer is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=closertest/
type Closer interface {
	Close() error
}
