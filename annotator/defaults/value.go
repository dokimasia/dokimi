// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

package defaults

import (
	"go.dokimi.dev/eidos/sdk/emit"
	"go.dokimi.dev/eidos/sdk/node"
	"go.dokimi.dev/eidos/sdk/rules"
	"go.dokimi.dev/eidos/sdk/symbol"
)

// aliasDepth is the most declared types that [Value] follows to the target of
// a type, so a cycle of a broken source ends.
const aliasDepth = 8

// Value returns the value of the text of a declared default of the type ref,
// in the scope of the declaration that declares the default. The rules b of
// the language read the text as a value of the language, such as a literal,
// the name of a declaration or a call. For a type of text, a text that the
// rules do not read is the content of the string, as the text of an authored
// sample is. Value reports false for a text that is no value of the type.
func Value(b *rules.Bound, scope *rules.Scope, ref *node.TypeRef, text string) (emit.Value, bool) {
	value, read := b.LiteralFor(*scope, ref, text)
	if read || !textual(b, ref) {
		return value, read
	}
	return emit.Literal(emit.LiteralString, text), true
}

// textual reports whether the type ref is text. A type that refers to a
// declared type, such as type Code string of Go, is text where the target of
// that declaration is text.
func textual(b *rules.Bound, ref *node.TypeRef) bool {
	shape := b.TypeOf(ref)
	for range aliasDepth {
		if shape.Form != symbol.FormReference {
			break
		}
		decl, _ := b.View().Lookup(shape.Ref)
		alias, declared := decl.(*node.Alias)
		if !declared {
			return false
		}
		shape = b.TypeOf(alias.Target)
	}
	return shape.Form == symbol.FormText
}
