// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package txwithretry stacks the tx contract with the retrysucceeds mixin on
// the commit role.
//
// It earns a composite fixture because retry changes what the contract's
// terminal-state rule means. Tx says a transaction reaches exactly one
// terminal state: once committed it cannot be rolled back, and once rolled
// back it cannot be committed. Retry says a call failing transiently succeeds
// on a later attempt.
//
// Together they ask whether a commit that failed reached a terminal state. If
// it did, the retry is a second terminal operation and the contract rejects
// it; if it did not, the retry is the same operation continuing. The answer
// here is the second — a failed commit leaves the transaction open — and a
// generator that assumes the first emits a suite failing on the retry it was
// told to expect.
//
// The generator answers it by binding the pairing rather than the mixin: the
// begin of the same tx contract runs first, and the free-standing
// retry claim is withheld because that one would be false here.
package txwithretry

import (
	"context"
	"errors"
)

// ErrClosed reports a terminal operation on a settled transaction.
var ErrClosed = errors.New("txwithretry: transaction already settled")

// ErrTransient reports a commit that may succeed on retry. It is deliberately
// distinct from ErrClosed: the difference between them is the whole question
// this fixture poses.
var ErrTransient = errors.New("txwithretry: transient commit failure")

// TxWithRetry is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=txwithretrytest/
type TxWithRetry interface {
	// Begin is the tx contract's begin role.
	//dokimi:contract tx role=begin
	Begin(ctx context.Context) error

	// Commit is the tx contract's commit role and retries on transient
	// failure. A transient failure leaves the transaction open, so the retry
	// is the same terminal operation rather than a second one.
	//
	// attempts= is declared, and which claim it binds is the pairing's
	// answer. Read alone the mixin says any caller may retry until it
	// converges, and a check driving Commit with nothing open meets a
	// transaction still Idle, reports ErrClosed on every attempt, and reds
	// a subject that honoured both directives. Read beside the begin of the
	// same tx contract, it says the retry converges once the
	// transaction is open — which is what the two of them state and
	// neither states alone.
	//dokimi:contract tx role=commit
	//dokimi:mixin retry-succeeds attempts=2
	Commit(ctx context.Context) error

	// Rollback is the tx contract's rollback role.
	//dokimi:contract tx role=rollback
	Rollback(ctx context.Context) error
}
