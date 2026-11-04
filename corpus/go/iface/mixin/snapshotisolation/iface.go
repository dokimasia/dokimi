// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package snapshotisolation is the mixin-axis fixture for the
// snapshot-isolation mixin, which declares that each transaction reads a
// snapshot of the store at its begin, and that two concurrent transactions do
// not both commit a write of one key.
//
// Begin returns a handle whose methods write, read, commit and roll back, as
// database/sql returns a transaction. Set writes the register of a key, so the
// check of the mixin records a history of registers.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression.
package snapshotisolation

import (
	"context"
	"errors"
)

// ErrClosed is the error of a call of a closed transaction.
var ErrClosed = errors.New("snapshotisolation: closed")

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=snapshotisolationtest/
type Mixed interface {
	// Begin opens a transaction and returns its handle.
	//
	//dokimi:contract tx role=begin closed=ErrClosed
	Begin(ctx context.Context) (Txn, error)
}

// Txn is the handle of a transaction of Mixed.
type Txn interface {
	// Commit applies the writes of the transaction. It fails where another
	// transaction committed a write of the same key after the begin.
	//
	//dokimi:contract tx role=commit
	Commit(ctx context.Context) error

	// Rollback discards the writes of the transaction.
	//
	//dokimi:contract tx role=rollback
	Rollback(ctx context.Context) error

	// Set writes v under key in the transaction.
	//
	//dokimi:mixin snapshot-isolation read=Load
	Set(ctx context.Context, key string, v int64) error

	// Load returns the number that the transaction wrote under key. For a key
	// that the transaction did not write, it returns the number of the
	// snapshot, which is 0 for a key that no transaction wrote.
	Load(ctx context.Context, key string) (int64, error)
}
