// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package stub

import "go.dokimi.dev/eidos/sdk/diag"

// CodePrefix is the prefix of the codes that stub reports under.
const CodePrefix diag.Prefix = "DOKIMISTUB"

// Constraint reports an interface with type-set elements, such as the Go
// interface ~int | ~float64. Such an interface constrains a type parameter
// and is never the type of a value, so no double implements it. The finding
// is an Error at the interface, and stub does not write a double of it.
var Constraint = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  1,
	Meaning: "a double cannot implement an interface that only constrains type parameters",
})

// IncompleteMembers reports a declaration whose members the rules of its
// language cannot list in full, such as an interface that embeds an
// interface that the load did not resolve, or an overloaded method. The
// finding is an Error at the declaration, and stub does not write a double
// of it.
var IncompleteMembers = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  2,
	Meaning: "a double cannot implement a declaration whose members are incomplete",
})

// Unimplementable reports a declaration with a member that a double cannot
// implement: a method with type parameters, an index signature, a construct
// signature, an accessor, a member with a computed name such as
// [Symbol.iterator], a member of a class that is not public, an unexported Go
// method of another package than the double, or a member without a declared
// type. It also reports a declaration that the file of its double cannot
// refer to, such as a TypeScript interface that is not exported, and a method
// whose signature refers to a type that the file cannot import. The finding
// is an Error at the declaration, and stub does not write a double of it.
var Unimplementable = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  3,
	Meaning: "a double cannot implement a member of its declaration, or cannot refer to the declaration",
})

// Collision reports a declaration whose double would have two members of one
// name, such as the field OnClose that configures Close and a method OnClose.
// The finding is an Error at the declaration, and stub does not write a
// double of it.
var Collision = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  4,
	Meaning: "a double would have two members of one name",
})

// Empty reports a declaration under the directive without a method or a
// property that a double can implement, such as an interface that declares
// nothing. The finding is an Error at the declaration, and stub does not
// write a double of it.
var Empty = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  5,
	Meaning: "a double of a declaration without a member implements nothing",
})

// Incomparable reports the param partition of the directive fault on a
// parameter whose arguments a partition fault cannot compare with a key: a
// variadic parameter, or a parameter whose type mentions a type parameter of
// the double. The finding is an Error at the declaration under the directive
// stub, and the double of the declaration has no partition helpers.
var Incomparable = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  6,
	Meaning: "a partition fault needs a parameter whose arguments compare with a key",
})

// BadField reports a declaration whose param field of the directive stub is
// not the name template of a field: a template that does not parse, that has
// an action other than {{.Method}} or that does not use {{.Method}}. The
// finding is an Error at the declaration, and stub does not write a double of
// it.
var BadField = diag.MustRegister(CodePrefix, diag.CodeSpec{
	Number:  7,
	Meaning: "the param field of a double needs a name template that uses {{.Method}}",
})
