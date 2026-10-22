// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package defaults stamps the declared default value of each field and of
// each type.
//
// defaults is an annotator with two rules. The directive rule reads
// //dokimi:default <value> on a field, and on a type alias or a defined type,
// in Go and in TypeScript. The tag rule reads the tag of the key that the
// option tag-key sets, default unless a config sets another key, in each
// language whose rules state tags:
//
//	type Address struct {
//		//dokimi:default "Dam 1"
//		Street string `default:"Main 1"`
//		Country string `default:"NL"`
//	}
//
//	//dokimi:default "test-key"
//	type Key string
//
// Within one plugin, the first claim in rule order takes precedence, and
// the directive rule comes first. So on a field with both forms, the value
// of the directive becomes the fact: Street defaults to Dam 1, and Country
// to NL. A tag with an empty value declares no default. A type has no tag,
// so the directive is the only form of its default. The option no-tags turns
// the tag rule off. The rules of TypeScript state no tags, so a TypeScript
// field declares its default through the directive alone.
//
// # Facts
//
// defaults stamps the text of a default under [KeyValue], as the kernel's
// sample annotator stamps the text of a sample. A reader turns the text into
// a value of the type of the declaration with [Value]. [KeySource] contains
// the form that declares the value, [SourceDirective] or [SourceTag].
// defaults provides the capability [Capability], so a plugin that reads the
// facts requires it and runs after it.
//
// # Dependency position
//
// The package imports the root, directive, emit, meta, node, plugin, rules
// and symbol packages of the SDK facade of eidos, and the Go stdlib.
package defaults
