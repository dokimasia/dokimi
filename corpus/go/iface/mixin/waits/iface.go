// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package waits is the mixin-axis fixture for the mixin waits, which declares
// that a call can wait until another call changes the state of its subject.
//
// Acquire waits until the caller holds the lock. A benchmark repeats the calls
// of a callable on one subject, so its second Acquire would wait for a Release
// that the benchmark never makes. The mixin keeps Acquire out of each
// generated benchmark and budget, and the benchmarks measure Release and Held.
package waits

import "context"

// Lock grants a resource to one holder at a time.
//
//dokimi:suite
//dokimi:out path=waitstest/
type Lock interface {
	// Acquire waits until the caller holds the lock, or until ctx ends.
	//dokimi:mixin waits
	Acquire(ctx context.Context) error

	// Release releases the lock that the caller holds.
	Release()

	// Held reports whether a caller holds the lock.
	Held() bool
}
