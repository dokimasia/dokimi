// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

// Package container is the generic half of the struct-kind corpus.
//
// Type parameters have to survive end to end: the builder type carries them,
// its constructor carries them, and every setter that touches a parameterised
// field takes the parameter rather than its constraint. A generator that
// erases them produces code which compiles only when T happens to be the
// constraint's underlying type, so the failure shows up for one instantiation
// and not another.
//
// Routing is declared once for the package rather than repeated on each
// type: every builder here lands in the same companion package, so a
// per-struct directive is the same statement written N times, and the Nth
// copy is the one that gets forgotten.
//
//dokimi:out path=containertest/
package container

import "time"

// Container is generic over one unconstrained parameter, which is the case
// where nothing about T can be assumed.
//
//dokimi:builder
type Container[T any] struct {
	Value T

	// Label is not parameterised, so its setter must not acquire a type
	// parameter it does not need.
	Label string

	// Items is a slice of the parameter, so its Append takes T rather than any.
	Items []T
}

// Pair is generic over two parameters with different constraints. A generator
// that assumes a single parameter, or that assumes every parameter is `any`,
// gets this wrong.
//
//dokimi:builder
type Pair[K comparable, V any] struct {
	Key K

	Value V

	// Index is keyed by the comparable parameter, so the map setter has to
	// thread both parameters through.
	Index map[K]V
}

// Keyable bounds a type parameter the deriver's palette cannot answer for.
// Deriving a witness only ever tries a builtin against a constraint's
// printed form, and admits the result solely when that form is exactly
// "any", "interface{}" or "comparable" — a named interface's printed form
// is its name, never one of those three, regardless of what it embeds or
// where it is declared. So every type below that bounds a parameter by
// Keyable derives no witness of its own, and needs //dokimi:witness to
// generate anything beyond the builder itself: the builder is written in
// terms of the parameter and needs no concrete type, but a check has to
// construct an actual value, and derivation is the only thing that could
// have said which type to construct one at.
type Keyable interface {
	~string | ~int
}

// Solo carries no //dokimi:witness, so its one parameter never derives
// and the generator withholds every check for it — the negative case
// [Bucket] below answers. The builder still generates in full: nothing
// about a builder's own code needs to name a concrete type.
//
//dokimi:builder
type Solo[K Keyable] struct {
	Key K
}

// PartialPair witnesses one of its two parameters and leaves the other
// to derive — which it cannot, so PartialPair generates no checks either,
// the same as [Solo]. Witnessing is not per-parameter best-effort: an
// instantiation is the whole parameter list at once, and a partial one
// is worth exactly as much as none. [Bucket] is PartialPair with the
// missing entry supplied.
//
//dokimi:witness V=time.Duration
//dokimi:builder
type PartialPair[K Keyable, V any] struct {
	Key   K
	Value V
}

// Bucket witnesses both of its parameters, so it generates real checks
// where [Solo] and [PartialPair] generate none. K's witness is load
// bearing — Keyable derives nothing on its own — while V's overrides a
// witness derivation would have supplied for an ordinary `any` parameter,
// which is worth doing anyway to prove an authored answer always wins.
// The pair also covers both forms a witness takes: K names a builtin
// needing no import, and V names a type through the package that
// declares it.
//
// TTL is not parameterised, and is what makes V's witness resolvable at
// all: resolving a qualified witness value against the file's own
// imports is the route that works for a single-segment import path like
// "time", where a full import-path spelling has no slash to distinguish
// itself by. TTL's own type is what makes the import real rather than
// directive-only.
//
//dokimi:witness K=string V=time.Duration
//dokimi:builder
type Bucket[K Keyable, V any] struct {
	Key   K
	Value V
	TTL   time.Duration
}

// Triple carries the directive through three parameters: A's witness is
// load bearing for the reason [Bucket]'s K is, and B and C are overrides
// for the reason [Bucket]'s V is — proving the directive's per-parameter
// KV parsing is not tuned to exactly two entries.
//
//dokimi:witness A=string B=int C=bool
//dokimi:builder
type Triple[A Keyable, B any, C comparable] struct {
	First  A
	Second B
	Third  C
}

// Quad is [Triple] with a fourth parameter also bound by [Keyable], so
// two of its four witnesses are load bearing rather than one — proving
// the deriver's all-or-nothing rule is checked at every position it
// applies to, not only the first.
//
//dokimi:witness A=string B=int C=bool D=string
//dokimi:builder
type Quad[A Keyable, B any, C comparable, D Keyable] struct {
	First  A
	Second B
	Third  C
	Fourth D
}
