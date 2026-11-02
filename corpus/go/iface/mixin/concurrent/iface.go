// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package concurrent is the mixin-axis fixture for the concurrent mixin, which
// declares that the method is safe to call from several goroutines at once.
//
// The law of the mixin is linearizability: every history of concurrent calls
// has an order in which a reference implementation, called one at a time,
// returns the same results. The interface keeps a value under each key and
// states the machinery reference. The mixin read-after-write pairs its read
// with its write, so each call records its key and the check splits the
// history by key.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression.
package concurrent

import (
	"context"
)

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite machineries=[reference]
//dokimi:out path=concurrenttest/
type Mixed interface {
	// Put stores value under key.
	//dokimi:mixin concurrent
	Put(ctx context.Context, key, value string) error

	// Get returns the value under key, and the empty value for a key without
	// a value.
	//dokimi:mixin concurrent
	//dokimi:mixin read-after-write write=Put
	Get(ctx context.Context, key string) (string, error)
}
