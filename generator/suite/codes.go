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

// BadFactory reports a subject whose factory takes more than one parameter,
// returns another number of results than one, or is a function that the
// conformance test cannot call. An application calls a factory with nothing
// or with the handle of a check, and takes one subject from it. The
// conformance test is in the external test package of the directory of the
// subject, or of the directory of its directive out. So it calls an exported
// factory or a factory of its own package. The finding is an Error at the
// subject, and suite records no conformance test for it.
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
// not the signature of a check. A check takes the handle of its check, and a
// property check also a *prop.Case after it. Its defect takes the handle and
// returns the planted defect alone. BadCheck also reports a door of the
// defect with another signature than its door. Each door returns its subject
// alone. A reference and a second take the handle, a restart the handle and
// the planted defect, and a crash the handle, the planted defect and a func()
// of the call that it crashes. BadCheck also reports the options of a check
// without a *prop.Case, and options that do not take the handle alone and
// return a []prop.Option.
// BadCheck also reports a check of a function with the identity of a check
// that the test of the suite adds to the forms of its application, and a
// check with the segment of a ceiling of a budget, such as allocation-budget.
// The finding is an Error at the function, and no suite binds the check.
var BadCheck = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number: 5,
	Meaning: "a function under the directive check has no statement, it, its defect, a door of its defect or " +
		"its options have another signature than a check, or it has the identity of a check that the test of its " +
		"suite or a budget adds",
})

// UnknownExcusal reports an excusal under the directive excuse whose identity
// does not match a check of the suite of the declaration of its subject. The
// finding is an Error at the directive, and suite does not record the
// conformance test of the subject.
var UnknownExcusal = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  6,
	Meaning: "the directive excuse names a check that the suite of the declaration of its subject does not have",
})

// BadBudget reports a directive budget that states a ceiling that no check
// can apply: a count of allocations or bytes below 0, a p99 or a mean that is
// no duration above 0, and a mean without a p99. It reports a warm-up below 0,
// and a budget without a ceiling and without a warm-up, which no benchmark can
// apply. It also reports a budget that no suite applies: a second budget of one
// callable, a budget of a callable that no suite has, a budget of a method
// that names a callable, a budget of a subject that names no callable of its
// declaration or a callable whose declaration states a budget, and a budget of
// a struct without the directive subject. eidos refuses a count that is no
// integer, and a param that the directive states twice. The finding is an
// Error at the directive. suite stamps nothing for the budget, and records no
// conformance test of a subject with such a budget.
var BadBudget = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number: 7,
	Meaning: "the directive budget states a ceiling that no check can apply, a warm-up that no benchmark can apply, " +
		"or a budget that no suite applies",
})

// RepeatedSubject reports a second directive subject on a struct that names
// the declaration of an earlier one. A struct has one conformance test for
// each declaration. The finding is an Error at the second directive, and suite
// records the conformance test of the first one alone.
var RepeatedSubject = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  8,
	Meaning: "a struct has two directives subject of one declaration",
})

// Unimportable reports an interface or a struct under the directive suite with
// a callable whose signature refers to a declaration that the file of its
// suite cannot import, such as a type of an internal package of another tree
// that the struct promotes from an embedded type. The suite spells the
// signature of each callable, so the finding is an Error at the declaration,
// and suite records nothing for it.
var Unimportable = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number: 9,
	Meaning: "a callable of a declaration under the directive suite refers to a declaration that the file of its " +
		"suite cannot import",
})

// Ungenerated reports a type in the inputs of a callable that no generator of
// the suite covers. An interface with an unexported method is such a type when
// the standard library declares it, or when its package does not export an
// implementation that the file of the suite can import. A struct with two
// fields of one JSON name is such a type. So is a type of a kind that the
// property engine refuses, when its generated code refers to a declaration
// that the file cannot import. Each run of the
// callable fails before its first case until the test process registers a
// generator of the type. The finding is a Warning at the declaration, and
// suite records the suite of the declaration.
var Ungenerated = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  10,
	Meaning: "a type in the inputs of a callable of a declaration under the directive suite has no generator",
})

// BadOptions reports a subject whose function of the param options the
// conformance test cannot call. The test calls the function without
// arguments and passes its one result, a list, to the application. So the
// function takes no parameter and returns one list, and the package of the
// test can call it. The finding is an Error at the subject, and suite records
// no conformance test for it.
var BadOptions = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number: 11,
	Meaning: "the function of the param options of the directive subject is not a function of the options of " +
		"an application",
})
