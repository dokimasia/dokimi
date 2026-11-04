// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package rules contains the rules of the generator: the classifications of
// the shape catalog that license each claim of the spec, the gates that bind
// it, the inputs and the partners of its check, the stance of each
// classification without a claim, and the reference forms of the verified
// fakes.
//
// People edit the YAML files of the directory, and go generate renders each
// of them into the JSON file of the same name and its JSON Schema under
// schemas. The package embeds the JSON files. On the first call of a function
// that reads a file, the package joins each rule with the runtime entry of the
// same identifier from go.dokimi.dev/spec.
//
// Each kind of entry has a function that returns every entry, such as
// [Claims], and a function that returns one entry by its identifier, such as
// [ClaimOf].
package rules

//go:generate go run ./internal/render/cmd
