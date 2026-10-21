// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package suite

import "go.dokimi.dev/eidos/sdk/diag"

// CodePrefix is the prefix of the codes that suite reports under.
const CodePrefix diag.Prefix = "DOKIMISUITE"

// Uninstantiable reports a declaration under the directive suite whose type
// parameters its suite cannot instantiate: a generic interface without
// witnesses, a generic struct and a generic function. A suite needs one type
// of its subjects, so the finding is an Error at the declaration, and suite
// records nothing for it.
var Uninstantiable = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  1,
	Meaning: "a declaration under the directive suite has type parameters that its suite cannot instantiate",
})

// NoSuite reports a subject whose declaration has no directive suite. The plan
// then emits no application of the suite of the declaration, so the finding is
// an Error at the subject, and suite records no conformance test for it.
var NoSuite = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  2,
	Meaning: "the declaration of the directive subject has no directive suite",
})

// BadFactory reports a subject whose factory takes more than one parameter, or
// returns another number of results than one. An application calls a factory
// with nothing or with the handle of a check, and takes one subject from it.
// The finding is an Error at the subject, and suite records no conformance
// test for it.
var BadFactory = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  3,
	Meaning: "the function of the param new of the directive subject is not a factory of subjects",
})
