// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package fake

import "go.dokimi.dev/eidos/sdk/diag"

// CodePrefix is the prefix of the codes that fake reports under.
const CodePrefix diag.Prefix = "DOKIMIFAKE"

// NoSuite reports an interface under the directive fake without the directive
// suite. The test of a fake applies the suite of its declaration, so the
// finding is an Error at the interface, and fake writes nothing for it.
var NoSuite = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  1,
	Meaning: "an interface under the directive fake has no directive suite",
})

// BadAbsent reports a param absent of the directive fake that is no name
// template of a sentinel. The finding is an Error at the interface, and fake
// writes nothing for it.
var BadAbsent = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  2,
	Meaning: "the param absent of the directive fake is no name template of a sentinel",
})

// Refused reports an interface under the directive fake with a callable that
// no reference form covers. The finding states the callable, what of it no
// form covers and the remedy of the gate reference-role. It is a Warning at the
// callable, and fake writes nothing for the interface.
var Refused = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  3,
	Meaning: "an interface under the directive fake has a callable that no reference form covers",
})

// Unimplementable reports an interface under the directive fake that the file
// of a fake cannot implement, or that the test of a fake cannot refer to. The
// finding is an Error at the interface, and fake writes nothing for it.
var Unimplementable = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  4,
	Meaning: "a fake cannot implement an interface under the directive fake, or its test cannot refer to the interface",
})
