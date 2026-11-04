// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package serializable is the mixin-axis fixture for the serializable mixin,
// which declares that concurrent transactions are equivalent to some serial
// order.
//
// Snapshot isolation excludes dirty writes and dirty, intermediate and
// circular reads. It allows write skew, where two transactions each read a
// key that the other writes and both commit. Serializability excludes write
// skew too.
//
// Each callable of a transaction takes the handle that Begin returns. Append
// appends to the list of a key, so the check of the mixin records a history
// of list appends.
//
// There is no negated form here. eidos declares the mixin directive
// DenyNegation, because a mixin is opt-in and deleting the directive is
// the suppression.
package serializable

import (
	"context"
	"errors"
)

// ErrClosed is the error of a call with the handle of a closed transaction.
var ErrClosed = errors.New("serializable: closed")

// Tx is the handle of a transaction of Mixed.
type Tx struct{ ID int64 }

// Mixed is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=serializabletest/
type Mixed interface {
	// Begin opens a transaction and returns its handle.
	//
	//dokimi:contract tx role=begin closed=ErrClosed
	Begin(ctx context.Context) (Tx, error)

	// Commit applies the appends of the transaction. It fails where another
	// transaction committed a change of a list that the transaction read.
	//
	//dokimi:contract tx role=commit
	Commit(ctx context.Context, tx Tx) error

	// Rollback discards the appends of the transaction.
	//
	//dokimi:contract tx role=rollback
	Rollback(ctx context.Context, tx Tx) error

	// Append appends v to the list of key in the transaction.
	//
	//dokimi:mixin serializable read=Get
	Append(ctx context.Context, tx Tx, key string, v int) error

	// Get returns the list of key in the transaction, with the appends of the
	// transaction at its end.
	Get(ctx context.Context, tx Tx, key string) ([]int, error)
}
