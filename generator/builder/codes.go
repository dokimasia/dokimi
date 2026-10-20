// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package builder

import "go.dokimi.dev/eidos/sdk/diag"

// CodePrefix is the prefix of the codes that builder reports under.
const CodePrefix diag.Prefix = "DOKIMIBUILDER"

// BadDefault reports a field whose declared default is no value of the type
// of the field. Such a default is neither a literal nor a name or a call
// that resolves in the scope of the struct. builder also reports the
// default of a field without a declared type, such as a TypeScript property
// without a type annotation, because it cannot check that default. The
// constructor of the builder cannot assign such a default, so the finding is
// an Error at the field, and builder writes nothing for the struct.
var BadDefault = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  1,
	Meaning: "the declared default of a field is no value of the type of the field",
})

// NotStruct reports an alias under the directive builder whose target is no
// struct that the workspace declares, and an alias that is generic or refers
// to a generic struct. A builder sets the fields of a struct without type
// parameters, so the finding is an Error at the alias, and builder writes
// nothing for it. The builder of a generic struct builds each instance of the
// struct.
var NotStruct = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number: 2,
	Meaning: "an alias under the directive builder refers to no struct of the workspace, " +
		"or the alias or its struct is generic",
})

// Unconstructed reports a class whose constructor takes a parameter that a
// caller must pass. A builder constructs the value without arguments, so the
// finding is an Error at the class, and builder writes nothing for it.
var Unconstructed = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  3,
	Meaning: "the constructor of a class under the directive builder takes a parameter that a caller must pass",
})

// Collision reports a declaration whose builder would have two methods of
// one name, such as the setters of the TypeScript properties name and Name.
// A type cannot declare both methods, so the finding is an Error at the
// declaration, and builder writes nothing for it.
var Collision = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  4,
	Meaning: "two methods of the builder of a declaration would have one name",
})

// BadTag reports a field whose tag builder has another value than -, the
// one value that leaves the field out of the builder. A typo in the tag would
// otherwise keep the setter without a word, so the finding is an Error at the
// field, and the field stays in the builder.
var BadTag = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  5,
	Meaning: "the tag builder of a field has another value than -",
})

// Empty reports a declaration under the directive builder that has no field
// that a builder can set. A builder without a setter cannot do what the
// directive asks, so the finding is an Error at the declaration, and builder
// writes nothing for it.
var Empty = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  6,
	Meaning: "a declaration under the directive builder has no field that a builder can set",
})
