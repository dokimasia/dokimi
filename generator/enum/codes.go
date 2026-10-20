// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package enum

import "go.dokimi.dev/eidos/sdk/diag"

// CodePrefix is the prefix of the codes that enum reports under.
const CodePrefix diag.Prefix = "DOKIMIENUM"

// Empty reports an enumeration under the directive enum without a variant. A
// surface of no variants has nothing to render or parse, so the finding is an
// Error at the enumeration, and enum writes nothing for it.
var Empty = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  1,
	Meaning: "an enumeration under the directive enum has no variant",
})

// Duplicate reports an enumeration under the directive enum whose two
// variants have one value or one text. A render cannot tell two variants of
// one value apart, and a parse cannot reach one of two variants of one text,
// so the finding is an Error at the enumeration, and enum writes nothing for
// it.
var Duplicate = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  2,
	Meaning: "two variants of an enumeration under the directive enum have one value or one text",
})
