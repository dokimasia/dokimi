// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package codeclossy is the contract-axis fixture for the codec contract's
// lossy form: a forward transform that discards information, paired with an
// inverse that recovers only what survived.
//
// `fidelity=lossy` selects the weaker roundtrip law: not the identity, but
// agreement on the second pass — encode, decode, encode again, and the two
// encodings match because everything the transform was going to lose is
// already gone. The exact form lives in the sibling `codec` fixture; both
// exist because the parameter selects between two laws, and a corpus stating
// only one leaves the other's selection untested.
//
// The model tier rides the twin floor and that is not a floor here, for
// the reason the lossless codec beside it gives: a codec holds no state
// for an oracle to model, and the claim is about the round trip.
package codeclossy

import (
	"context"
)

// Contract is the fixture interface.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=codeclossytest/
type Contract interface {
	// Encode is the codec contract's forward role. Its directive sets the
	// fidelity of the round trip.
	//dokimi:contract codec role=forward fidelity=lossy
	Encode(ctx context.Context, in string) (string, error)

	// Decode is the codec contract's inverse role.
	//dokimi:contract codec role=inverse
	Decode(ctx context.Context, in string) (string, error)
}
