// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package store is the errors-kind fixture: the input the sentinel generator
// reads.
//
// The directive sits at package scope rather than on the variables. sentinel
// declares On(node.KindPackage), so what a generator *reads* and where it is
// opted in are different things — a fixture carrying a var-level directive
// would never be picked up.
//
// The prefix deliberately differs from the package name. An override set to
// the value the derivation would have produced anyway cannot tell "the
// override was honoured" from "it was ignored", which is the standard the
// noprefix fixture beside this one states for the suppression branch. A
// generator that read prefix= and dropped it derives `store` here, and every
// message below then fails the assertion it emits.
//
//dokimi:sentinel prefix=kv no-overlap-with=[go.dokimi.dev/dokimi/corpus/go/errors/neighbour]
package store

import (
	"errors"
	"fmt"
)

// The sentinel set. Every message carries the package prefix, none is a
// prefix of another, and none satisfies errors.Is against a sibling — which
// are the three properties the generated checks assert.
var (
	// ErrNotFound reports a key that was never written.
	ErrNotFound = errors.New("kv: not found")

	// ErrConflict reports a write losing a race.
	ErrConflict = errors.New("kv: conflict")

	// ErrClosed reports use after teardown.
	ErrClosed = errors.New("kv: closed")
)

// NotFoundError is a custom error carrying a field. The generated round-trip
// check formats it, parses the field back out, and compares — so the field has
// to survive the format string.
type NotFoundError struct {
	Key string
}

// Error implements error.
func (e *NotFoundError) Error() string {
	return fmt.Sprintf("kv: key %q not found", e.Key)
}

// Is reports NotFoundError as a kind of [ErrNotFound], which is the optional
// method the generator detects and exercises.
func (*NotFoundError) Is(target error) bool { return target == ErrNotFound }

// WrappedError wraps a cause, so the generated unwrap-chain check has
// something to walk.
type WrappedError struct {
	Op    string
	Cause error
}

// Error implements error.
func (e *WrappedError) Error() string {
	return fmt.Sprintf("kv: %s: %v", e.Op, e.Cause)
}

// Unwrap exposes the cause.
func (e *WrappedError) Unwrap() error { return e.Cause }
