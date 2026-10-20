// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package gentest loads a tree of sources into a fixture of the SDK of
// eidos, for the specs of the plugins of dokimi, and reads what a plugin
// declares and records.
//
// [New] loads the packages of a tree under [Brand] with the frontends and
// the rules of Go, TypeScript and the scripted language of eidos. It
// registers the keys and validates the directives of the plugins of a
// [Spec], and the annotators of the spec stamp the fixture, as a run of the
// composition does before the plugin under test runs. [Fixture.Find]
// returns a declaration of the fixture by its kind, its owner and its name.
//
// A spec of a generator reads the result of a phase call through these
// helpers:
//
//   - [Records] returns the records of one type of data by their keys, and
//     checks that each family claims its files through its template.
//   - [Finding] returns the one finding at a declaration, and [Unrecorded]
//     checks that the generator records nothing for a subject.
//   - [Assigned] returns the names that the statements of a record assign.
//   - [Rendered] renders what a generator records through a backend, so a
//     spec matches the files with golden.MatchTree.
//   - [Project] returns the rendered files as the output of a run over a
//     tree, so a spec checks them with the toolchain of their language.
//
// [Templates] checks the template tree that a plugin declares for a target,
// and [Codes] checks the spelling and the meaning of the codes of a plugin.
//
// # Dependency position
//
// The package imports go.dokimi.dev/assert and its package expect, the
// frontend and rules packages of the Go and the TypeScript satellites of
// eidos, the backend, backendtest, diag, emit, frontendtest, node, output,
// plugin, plugintest, rules, rulestest, symbol and toolchain packages of the
// SDK facade, and the Go stdlib. Only test files import it.
package gentest
