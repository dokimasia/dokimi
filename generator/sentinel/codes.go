// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package sentinel

import "go.dokimi.dev/eidos/sdk/diag"

// CodePrefix is the prefix of the codes that sentinel reports under.
const CodePrefix diag.Prefix = "DOKIMISENTINEL"

// NoErrors reports a package under the directive sentinel without an error
// value and without an error type. A file of checks of nothing would read as
// checks of the errors of the package, so the finding is an Error at the
// package, and sentinel writes nothing for it.
var NoErrors = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  1,
	Meaning: "a package under the directive sentinel has no error value and no error type",
})

// OwnNeighbour reports a package that lists itself under the param
// no-overlap-with. Each error value matches itself, so the check could only
// fail. The finding is an Error at the package, and sentinel leaves the
// package out of the list.
var OwnNeighbour = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  2,
	Meaning: "a package under the directive sentinel lists itself under the param no-overlap-with",
})

// Unloaded reports a package under the param no-overlap-with that the run
// did not load. The run has no error values of that package, so its check
// compares with nothing. The finding is a Warning at the package, and the
// checks keep the empty check, so that the file shows the package.
var Unloaded = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  3,
	Meaning: "a package under the param no-overlap-with of the directive sentinel is not in the run",
})

// Unconstructed reports an error type of a language with constructors that
// has no public constructor whose required parameters have samples. A check
// cannot make a value of the type, so the finding is a Warning at the type,
// and the checks leave the type out.
var Unconstructed = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  4,
	Meaning: "an error type under the directive sentinel has no public constructor that a check can call",
})

// UnknownRemote reports a type under the param remote that is no public error
// type of the package of the directive. The param leaves out the checks of
// the message of the zero value and of the prefix of a remote error type, so
// a type of another package or a type that is no error leaves out nothing.
// The finding is an Error at the package, and the checks of the package keep
// every check of their error types.
var UnknownRemote = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  5,
	Meaning: "a type under the param remote of the directive sentinel is no error type of its package",
})
