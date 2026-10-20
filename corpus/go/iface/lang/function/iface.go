// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package function is the language-axis fixture for free functions.
//
// eidos annotates package-level functions through the same three axes it
// applies to interface methods — its walker has an OnFunction hook alongside
// OnMethod. Every other fixture in this corpus declares an interface, so a
// generator or a gate that reads only methods would look complete against all
// of them and still miss every classification a consumer put on a plain
// function.
//
// The functions here carry a body rather than being declarations, because a
// free function has nowhere else to live: an interface method is a signature
// and a function is not.
//
// Nothing is generated from them, and that is the gap this fixture records
// rather than closes. Every generator in testkit starts from an interface —
// there is no function hook in any of them — so the directives below reach
// the annotation census and stop there. The census is what holds them: a
// classification reachable only outside an interface would read as uncovered
// without these, which is what the package is for. Emitting a double from a
// free function is testkit's to write, and until it is, no //dokimi:out
// appears here to promise otherwise.
package function

import (
	"context"
	"errors"
)

// ErrNotFound is the sentinel the reading function reports.
var ErrNotFound = errors.New("function: not found")

// Value is the payload these functions carry.
type Value struct{ Key, Body string }

// Get is a reader by signature alone, classified with no directive — the same
// inference that applies to a method of the same shape.
func Get(_ context.Context, key string) (Value, error) {
	if key == "" {
		return Value{}, ErrNotFound
	}
	return Value{Key: key}, nil
}

// Put carries a mixin, proving the directive axes reach a function and not
// only a method.
//
//dokimi:mixin idempotent
func Put(_ context.Context, _ Value) error { return nil }

// Acquire is the lease contract's acquire role on a free function. A
// contract binds the functions of a package as it binds the methods of an
// interface.
//
//dokimi:contract lease role=acquire
func Acquire(_ context.Context, _ string) error { return nil }

// Release is the lease contract's release role.
//
//dokimi:contract lease role=release
func Release(_ context.Context, _ string) error { return nil }
