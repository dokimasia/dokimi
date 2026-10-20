// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package embedded is the language-axis fixture for an interface embedding
// another, so the method set is larger than the declaration.
//
// This axis varies the Go type system rather than the classification, so these
// break generators independently of any directive.
//
// Every embed here is declared in this package. [embeddedforeign] holds the same
// shape with the embed coming from outside the run, and the two must agree —
// which is the property having both is for.
package embedded

import (
	"context"
	"errors"
)

// ErrUnreachable is the sentinel [Base] reports, declared here so a fault
// helper generated onto a double of an embedding interface has something to
// qualify against.
var ErrUnreachable = errors.New("embedded: unreachable")

// Base is embedded into [Composed].
//
//dokimi:stub
//dokimi:out path=embeddedtest/
type Base interface {
	// Ping carries a fault directive so an embedding interface's double is
	// checked for helpers generated from a method it never declared. A
	// contributor reading only the declarations would silently skip this.
	Ping(ctx context.Context) error
}

// Closer is embedded into [Composed] alongside [Base].
//
//dokimi:stub
//dokimi:out path=embeddedtest/
type Closer interface {
	Close(ctx context.Context) error
}

// Composed embeds two interfaces and adds one method. A generator that reads
// only the declared methods emits an incomplete stub, which fails the
// compile-time interface check.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=embeddedtest/
type Composed interface {
	Base
	Closer
	Get(ctx context.Context, key string) (string, error)
}
