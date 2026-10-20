// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package dokimi composes the plugins of dokimi into the workspace of the
// binary.
//
// [Base] registers the brand, the frontends, the rules and the targets of Go
// and TypeScript, and the annotators. [Plans] returns the plan of each
// language, and [Compose] returns Base with the plans:
//
//	w, err := dokimi.Compose().Build()
//
// # Dependency position
//
// The package imports the annotators and the generators of the module, the
// packages authored and workspace of eidos-core, the root, annotate,
// backend, frontend and rules packages of the Go satellite of eidos, the
// root, backend, frontend and rules packages of its TypeScript satellite,
// the annotators of the shape catalog of eidos, and the packages output and
// plugin of the SDK facade.
package dokimi
