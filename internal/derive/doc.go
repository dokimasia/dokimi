// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package derive derives the checks of a declaration from the conformance
// spec. It reads the facts that eidos derives for each callable, binds each
// claim whose licence the callable has and whose gates pass, and lists every
// other licensed claim as a refusal and every classification without a
// claim with its stance.
//
// The derivation reads only facts that the rules of eidos derive in the same
// form for every language: the roles of the parameters and the returns, the
// error model, the shapes of the types, the samples and the classifications.
// Its functions are pure, so their tests build the facts as values.
//
// # Dependency position
//
// The package imports the package spec of the module, the shape catalog of
// eidos, the packages emit, rules and symbol of the SDK facade of eidos, and
// the Go stdlib. It imports no satellite of a language.
package derive
