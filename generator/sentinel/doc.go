// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package sentinel generates the checks of the error values and of the error
// types of each package under the directive sentinel.
//
// The plugin reads the error values of a package and its error types from
// the rules of eidos, and records a model of their checks. The template of
// each target writes the checks in the conventions of its language:
//
//	//dokimi:sentinel prefix=kv no-overlap-with=[example.com/shop/neighbour]
//	package store
//
// The checks pin that each error value is set and has a message that opens
// with the prefix, that no two error values share a message or match each
// other, and that no message opens with another. They pin that each error
// type reports a message, matches no other error type, agrees with the
// standard comparison, exposes its cause and puts each field of text into
// its message.
//
// A package can declare a family of error values, in which a child wraps its
// parent at the opening of its message:
//
//	var ErrPolicy = errors.New("policy: invalid policy")
//
//	var ErrExposure = fmt.Errorf("%w: no exposure", ErrPolicy)
//
// The rules of eidos state the wrap. The checks of a child pin that it
// matches its parent and that its message opens with the message of its
// parent. The checks against another error value leave out the ancestors of
// the child: its parent, the parent of its parent, and so on.
//
// # Options
//
// The option prefix sets the opening of each message of a package whose
// directive states no prefix, in place of the name of the package.
//
// # Languages
//
// The file golang.go contains the part of Go, and typescript.go the part of
// TypeScript. Each part declares the target of its language with its
// template tree, the opening of a message by default and the fields that
// every error of its language has. A check of Go writes into the zero value
// of an error type. A check of TypeScript constructs an error class through
// its first public constructor whose parameters have samples. sentinel
// leaves a package of a language without a part alone.
//
// # Diagnostics
//
// sentinel reports [NoErrors] at a package without an error value and an
// error type, [OwnNeighbour] at a package that lists itself under the param
// no-overlap-with, and [Unloaded] at a package that lists a package that the
// run did not load. It reports [Unconstructed] at an error type without a
// constructor that a check can call.
//
// # Dependency position
//
// The package imports the package internal/lang of the module, the root
// packages of the Go and the TypeScript satellites in golang.go and
// typescript.go, the package naming of eidos-lang, the root, diag,
// directive, emit, node, plugin, rules and symbol packages of the SDK facade
// of eidos, and the Go stdlib.
package sentinel
