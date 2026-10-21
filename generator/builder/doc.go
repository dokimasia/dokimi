// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package builder generates a fixture builder of each struct or class under
// the directive builder, and the checks of each builder.
//
// The plugin derives a model of each builder from the facts of eidos and
// records it: the fields that code outside the package can assign, the
// form of each field, the seed of a new builder and the declared defaults.
// The templates of each target write the model in the conventions of their
// language:
//
//	a := NewAddressBuilder().WithStreet("Dam 1").AppendLines("2nd floor").Build()
//	const a = AddressBuilder.create().withStreet("Dam 1").appendLines("2nd floor").build();
//
// A builder has a setter for each member, Mutate for a change that no setter
// makes, Clone for a copy whose collections the original does not share, and
// Build. The checks of a builder pin these properties:
//
//   - The constructor from a value keeps each member, and the constructor
//     from the seed applies each declared default.
//   - Each setter writes its member. A second call of a setter replaces the
//     value of the first.
//   - A clone does not share a collection with the original.
//   - A probe calls each method that these checks leave out, such as an
//     appender, the setters of the entries of a map and the setter of a
//     member without a sample.
//   - Mutate passes one value to each function. Clone and From keep the
//     value under construction.
//
// # Facts
//
// builder is a plugin with an annotate half and a generate half. The annotate
// half stamps the name of the builder of each declaration and the name of its
// constructor from the seed under [KeyName] and [KeySeed], and reports each
// declaration that has no builder. The generate half writes the builder of
// each declaration with the facts. Another generator calls a builder through
// the facts, as the stub generator does for the default results of a double.
//
// # Options
//
// The option suffix sets the word of the name of each builder, and the
// option companion-word sets the word of the name of each companion
// function. The option no-tests leaves out the checks.
//
// # Languages
//
// The file golang.go contains the part of Go, and typescript.go the part of
// TypeScript. Each part declares the target of its language with its
// template tree and its helpers, the types that the language writes as
// sets, the names of the constructors and of the methods of a member, the
// fields that its language writes after a dot, and the name of the
// companion function of a declaration. builder leaves a declaration of a
// language without a part alone.
//
// # Diagnostics
//
// builder reports [BadDefault] at a field whose default is no value of its
// type, and [BadTag] at a field whose tag builder has another value than -.
// It reports [NotStruct] at an alias whose target is not a struct of the
// workspace, and at an alias that has type parameters or refers to a
// generic struct. It reports [Unconstructed] at a class whose public
// constructors each take a parameter that a caller must pass. It reports
// [Empty] at a declaration without a field that a builder can set, and
// [Collision] at a declaration whose builder would have two methods of one
// name. builder records nothing for a declaration that it reports, except
// after [BadTag], which leaves the field in the builder.
//
// # Dependency position
//
// The package imports the annotator defaults and the package internal/lang
// of the module, the root packages of the Go and the TypeScript satellites
// in golang.go and typescript.go, the package naming of eidos-lang, the root,
// diag, directive, emit, meta, node, plugin, rules and symbol packages of the
// SDK facade of eidos, and the Go stdlib.
package builder
