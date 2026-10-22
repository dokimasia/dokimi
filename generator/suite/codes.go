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

// UnknownCallable reports a function under the directive check whose callable
// no declaration under the directive suite has, or two of them have. The
// finding is an Error at the function, and no suite binds the check.
var UnknownCallable = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number: 4,
	Meaning: "the directive check has a callable that no declaration under the directive suite has, or that " +
		"two of them have",
})

// BadCheck reports a function under the directive check whose doc comment
// states no statement, or whose signature or the signature of its defect is
// not the signature of a check. A check takes the handle of its check alone,
// and its defect takes the handle and returns the planted defect alone.
// BadCheck also reports a check of a function with the identity of a check
// that the test of the suite adds to the forms of its application. The
// finding is an Error at the function, and no suite binds the check.
var BadCheck = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number: 5,
	Meaning: "a function under the directive check has no statement, it or its defect has another signature " +
		"than a check, or it has the identity of a check that the test of its suite adds",
})

// UnknownExcusal reports an excusal under the directive excuse whose identity
// does not match a check of the suite of the declaration of its subject. The
// finding is an Error at the directive, and suite does not record the
// conformance test of the subject.
var UnknownExcusal = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  6,
	Meaning: "the directive excuse names a check that the suite of the declaration of its subject does not have",
})
