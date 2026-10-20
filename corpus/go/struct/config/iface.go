// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package config is the field-directive half of the builder's defaults
// handling.
//
// The generator seeds New<Type>() from one of three sources: a Defaults()
// companion in the test package, a per-field default (a //dokimi:default
// directive or a `default:"…"` struct tag — see [Config] and [TagConfig]),
// or the zero value. Each produces a different constructor, so they are
// three fixtures rather than one — [go.dokimi.dev/dokimi/corpus/go/struct/domain]
// carries the companion, this carries the per-field forms, and
// [go.dokimi.dev/dokimi/corpus/go/struct/plain] carries neither.
//
// There is deliberately no Defaults function here. Adding one would make the
// fixture test precedence between the two sources rather than the directive
// path itself, and precedence is a question for a composite.
//
// Routing is declared once for the package rather than repeated on each
// type: every builder here lands in the same companion package, so a
// per-struct directive is the same statement written N times, and the Nth
// copy is the one that gets forgotten.
//
//dokimi:out path=configtest/
package config

// Config carries a default on every literal kind the directive accepts, since
// each is parsed differently: quoting for strings, base detection for ints,
// keyword for bools, and the untyped nil that has no literal form of its own.
//
//dokimi:builder
type Config struct {
	Host string //dokimi:default "localhost"

	Port int //dokimi:default 8080

	Verbose bool //dokimi:default true

	Ratio float64 //dokimi:default 0.75

	// Retries has a zero default stated explicitly. A generator that treats
	// "zero" as "no directive" emits the same constructor either way, and this
	// is the field that catches it.
	Retries int //dokimi:default 0

	// Fallback is a pointer defaulted to nil. Its literal is the same as the
	// zero value, so it tests that the directive is read rather than inferred.
	Fallback *string //dokimi:default nil

	// Undirected carries no directive and must stay at its zero value even
	// though its neighbours do not.
	Undirected string
}

// TagConfig carries the same defaults as [Config], one literal kind per
// field, declared as `default:"…"` struct tags instead of directive
// comments — the other of the two forms the defaults annotator reads. A
// generator that only wired the directive half of the annotator's stamp
// would silently emit a builder that never seeds a tagged field, which is
// exactly the failure this fixture, driven for real, cannot pass under.
//
// The tag's value is stamped exactly as [reflect.StructTag.Get] returns
// it, on the same terms as a directive's positional argument: unparsed
// source text, ready to render as the literal it names. Go's tag syntax
// consumes one layer of quoting to delimit the value, so `default:"localhost"`
// reads back as the six bare characters `localhost` rather than a quoted
// string. The generator infers the quoting from the field's type, so the
// plain form below stamps a string literal without escaping. Host is the
// fixture for it.
//
//dokimi:builder
type TagConfig struct {
	// Host's tag needs no escaping; see the type's docblock.
	Host string `default:"localhost"`

	Port int `default:"8080"`

	Verbose bool `default:"true"`

	Ratio float64 `default:"0.75"`

	// Retries states a zero default through the tag, for the reason
	// [Config.Retries] states one through the directive.
	Retries int `default:"0"`

	// Fallback is a pointer defaulted to nil, through the tag rather than
	// the directive.
	Fallback *string `default:"nil"`

	// Undirected carries no tag and must stay at its zero value even though
	// its neighbours do not.
	Undirected string
}

// MixedConfig declares one field's default both ways at once, with
// different values, so the generated constructor proves which one the
// annotator prefers.
//
// The directive must win. It is the more specific statement — a tag is
// often copied along with the field it sits on, while a directive is
// written deliberately at the line it governs — and a rule that let the
// tag win would leave no way to correct one.
//
//dokimi:builder
type MixedConfig struct {
	// Label names "tag-wins" through the tag and "directive-wins" through
	// the directive. The generated constructor must seed "directive-wins";
	// seeding "tag-wins" is the precedence rule inverted.
	Label string `default:"tag-wins"` //dokimi:default "directive-wins"
}
