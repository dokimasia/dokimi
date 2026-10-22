// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package enum generates the text surface of each enumeration under the
// directive enum, and the checks of the surface.
//
// The plugin reads the variants of an enumeration and their texts from the
// rules of eidos, and records a model of the surface and of its checks. The
// templates of each target write the model in the conventions of their
// language:
//
//	s, err := ParseStatus("Published") // errors.Is(err, ErrUnknownStatus) for any other text
//	s.String()                        // "Published"
//	s.IsValid()                       // true
//	StatusValues()                    // every variant, in order
//
//	const s = parseStatus("Published"); // throws UnknownStatusError for any other text
//	statusToString(s);                  // "Published"
//	isStatus(s);                        // true
//	statusValues();                     // every member, in order
//
// The text of a variant of an enumeration of text is its value, and the text
// of any other variant is its name. Under the directive enum text=comment, the
// trailing comment of a variant is its text:
//
//	Poor Grade = iota // poor
//
// The checks pin the number of the variants, the distinct values and texts
// of the variants, each text that a value or a comment states, the round trip
// of each variant through its text, the refusal of a text that no variant
// has, the validity of a value outside the variants and of the zero value. In
// Go they also pin the refused text in the error of the parse function, and
// the round trip through MarshalText and UnmarshalText.
//
// # Options
//
// The options parse-function, values-function and parse-error are the name
// templates of the parse function, of the function that lists the variants
// and of the error of the parse function. The option valid-method is the
// method template of the validity method of Go. A template has the syntax of
// text/template, and {{.Enum}} is the name of the enumeration. Each language
// writes a rendered name in its own case convention:
//
//	parse-function: Parse{{.Enum}}     # ParseStatus, parseStatus
//	values-function: "{{.Enum}}Values" # StatusValues, statusValues
//	parse-error: Unknown{{.Enum}}      # ErrUnknownStatus, UnknownStatusError
//	valid-method: IsValid              # IsValid
//
// The option no-tests leaves out the checks.
//
// # Languages
//
// The file golang.go contains the part of Go, and typescript.go the part of
// TypeScript. Each part declares the target of its language with its
// template tree and its helpers, the names of the members of the surface,
// the reference to a variant, the literal of a text, and the conversion that
// prints a value outside the variants. enum leaves an enumeration of a
// language without a part alone.
//
// The part of Go also returns the table of the texts of an enumeration of
// integers, and the render slices the text of a value out of it in the form
// of the tool stringer: a compile-time check of the value of each variant,
// the texts in one string with an index of offsets, and a switch over the
// runs of consecutive values, or a map where the runs are more than ten. The
// render of an enumeration of text whose texts are its values converts the
// value to text. Every other render switches over the variants.
//
// # Diagnostics
//
// enum reports [Empty] at an enumeration without a variant, and [Duplicate]
// at an enumeration whose two variants have one value or one text. It reports
// [Conflict] at a declaration of the package that has the name of a member of
// the surface and another form. enum records nothing for an enumeration that
// it reports.
//
// # Dependency position
//
// The package imports the package internal/lang of the module, the root
// packages of the Go and the TypeScript satellites in golang.go and
// typescript.go, the package naming of eidos-lang, the root, diag,
// directive, emit, node, plugin, rules and symbol packages of the SDK
// facade of eidos, and the Go stdlib.
package enum
