// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package spec contains the conformance spec of dokimi: the laws, claims,
// gates, probes, reason codes, defect forms, machineries and stances that
// the binary shares with the runtimes.
//
// People edit the YAML files of the directory, and go generate renders each
// of them into the JSON file of the same name, its JSON Schema under
// schemas, and the constants of names.gen.go. The package embeds the JSON
// files and decodes them on the first call of a function that reads them.
//
// Each kind of entry has a function that returns every entry in the order
// of its file, such as [Probes], and a function that returns one entry by
// its identifier, such as [ProbeOf].
package spec

//go:generate go run ./internal/render/cmd
