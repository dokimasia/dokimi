// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Command dokimi generates the tests and the test helpers of the
// declarations of a workspace.
//
// The binary runs the two plans of the composition of dokimi. The plan go
// reads the Go packages of the workspace, and the plan typescript reads its
// TypeScript sources. A source file sets a directive as //dokimi:<name>,
// such as //dokimi:builder on a struct.
//
// # The config file
//
// The config file is .dokimi.yaml at the workspace root, and the state
// directory is .dokimi/. The key sources of a plan replaces the packages that
// the plan reads, and enabled: false leaves the plan out:
//
//	version: 1
//	plans:
//	    go:
//	        sources: {lang: golang, packages: [./internal/...]}
//	    typescript:
//	        enabled: false
//
// # Stores
//
// A run locates the stores of a frontend only when the tree has an input
// file of the frontend. A tree with a Go file needs GOROOT or a go binary on
// PATH.
//
// # Dependency position
//
// The command imports the root package of the module and the package cli of
// eidos.
package main
