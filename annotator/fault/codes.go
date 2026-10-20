// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package fault

import "go.dokimi.dev/eidos/sdk/diag"

// CodePrefix is the prefix of the codes that fault reports under.
const CodePrefix diag.Prefix = "DOKIMIFAULT"

// Infallible reports the directive fault on a method that cannot fail, such
// as a Go method without an error result. A fault of the double of such a
// method has no error to return. The finding is an Error at the directive,
// and fault stamps nothing of the instance.
var Infallible = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  1,
	Meaning: "a fault of a double needs a method that can fail",
})

// BadRetry reports the param retry below 2. A retry schedule fails each call
// before the attempt of retry, so an earlier attempt fails no call. The
// finding is an Error at the directive, and fault stamps no retry of the
// instance.
var BadRetry = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  2,
	Meaning: "a retry schedule lets an attempt from the second on succeed",
})
