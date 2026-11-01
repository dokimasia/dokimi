// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package ledger is the fixture of a transaction over a store whose calls fail
// on demand. Run stages the writes of its body, and commits them to the store
// in one call once the body succeeds. A commit that fails leaves the store as
// it was.
package ledger

import (
	"context"
	"errors"
)

// ErrMissing is the error of a read of a key without a committed value.
var ErrMissing = errors.New("ledger: missing")

// Ledger runs bodies in transactions over a store.
//
//dokimi:suite machineries=[faults]
//dokimi:out path=ledgertest/
type Ledger interface {
	// Run runs the body in a transaction. It commits what the body wrote
	// through Put once the body succeeds, and discards it when the body or the
	// commit fails.
	//
	//dokimi:contract transaction role=fn not-found=ErrMissing write=Put read=Get
	Run(ctx context.Context, body func(ctx context.Context) error) error

	// Put writes the value under the key: into the transaction of the context
	// inside Run, and into the store outside Run.
	Put(ctx context.Context, key, value string) error

	// Get returns the committed value under the key, and ErrMissing for a key
	// without one.
	Get(ctx context.Context, key string) (string, error)
}
