// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package receivercollision is the language-axis fixture for a parameter whose
// identifier is the one a generated method would bind its receiver to.
//
// Go allows a parameter named `s`; a generated method whose receiver is also
// `s` declares the identifier twice in one signature and does not compile. The
// generator cannot pick the receiver before it knows the parameters, which is
// why the collision guard runs after the projection rather than before it —
// and why nothing catches this by inspection: the double is correct for every
// interface whose author happened not to use that letter.
//
// This axis varies the Go type system rather than the classification, so these
// break generators independently of any directive.
//
// No suite directive, and that is the fixture's second finding. The double
// derives a recorded field per parameter name AND type, so PutS and GetS are
// separate values; the suite derives one fixture field per role, and a role is
// the parameter name alone. One name at two types is a role that cannot hold
// one value, so the suite generator refuses the whole interface with a
// positioned diagnostic naming both methods and both types. Keying the role on
// the type as well would put PutS and GetS in the consumer's Config — names the
// author never wrote — which is the trade the role derivation declines to make.
package receivercollision

import (
	"context"
)

// Session is the payload, named so a parameter of it reads naturally as `s`.
type Session struct{ ID string }

// Store declares the collision on every shape that carries a parameter, so a
// receiver bound before the parameters are known fails here rather than in one
// arm of the template a fixture happened to miss.
//
//dokimi:stub
//dokimi:out path=receivercollisiontest/
type Store interface {
	// Put is the writer form: one parameter, named for its type.
	Put(ctx context.Context, s Session) error

	// Get is the reader form, and takes the letter on the key rather than on
	// the value so the collision is not only ever in last position.
	Get(ctx context.Context, s string) (Session, error)

	// Touch returns nothing, which is the arm that renders through invokeVoid
	// rather than invoke.
	Touch(ctx context.Context, s Session)
}
