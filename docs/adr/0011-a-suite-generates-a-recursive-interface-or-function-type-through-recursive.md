---
adr: 0011
title: A suite generates a recursive interface or function type through Recursive
status: Accepted
date: 2026-10-26
supersedes: RFC-0006 in part
superseded-by: none
rfc: 0006
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0011: A suite generates a recursive interface or function type through Recursive

## Status

Accepted

## Context

A suite gives each function type and each interface outside the standard library at some depth of an input a generated function or an implementation of the suite, through `prop.Using`. The case draws the results of each generated function. An interface whose results contain the interface, as `Children() []Node` does, then has no value:

- Inside the stated generator of a type, a draw of `Of` of that type decodes without the stated generator, as RFC-0037 of assert-spec states.
- The shape reader refuses an interface, so every case fails.
- A function type that returns itself, such as `type StateFn func(*Lexer) StateFn`, fails in the same way, and so do two interfaces whose methods return each other.

RFC-0011 of assert-spec bounds a recursive shape by 100 nodes. dokimasia/assert-spec#21 requests a bound for a recursion through `using`: inside the extension of a `recursive` generator that `using` states for a type, a draw of `of` of a type that contains the type takes the generator of the extension.

## Decision

We will generate an interface or a function type whose generated results contain the type through `prop.Recursive`, once the property engine resolves the type inside the extension.

- A type is recursive where the results of its generated functions contain the type, directly or through another generated type.
- The base is the implementation or the function that returns its zero results at once, with the end error of each error result.
- The extension is the generator that the suite writes for the type, and each draw inside it keeps its form.
- The suite passes the recursive generator through `prop.Using`, as it passes the generator of every other type at depth.

## Alternatives Considered

### A recursion composed in dokimi

The extension could put its own generator at each position of the type, as `prop.List(self)` for a list and `prop.Dict` for a map. It works with the property engine as it is. A struct or a pointer that contains the type then needs the suite to compose the struct field by field. RFC-0006 rejects that composition, because it repeats the shape reader.

### A lazy reference

A generator could refer to itself through a lazy reference, as `deferred` of Hypothesis and `Deferred` of Rapid do. A lazy reference does not bound the size of a value, and Rapid documents that its `Make` may not terminate on a recursive type.

## Consequences

**Positive:**

- A tree, an iterator or a state function at some depth of an input gets values.
- A failure shrinks towards the base, so a counterexample shows the shallowest value that fails.

**Negative:**

- The suite waits for assert-spec#21 and its implementation in assert-go.
- At the leaves, a method returns its zero results, such as a nil `Iterator`. A subject that calls the nil value panics.
- Each leaf is a whole implementation with lists of its own, so the default budget of 100 leaves can make one case large. The suite takes a smaller budget, set by a measurement.

## References

| What | Where |
|---|---|
| RFC-0006, the values of the parameters of a callable | [0006-generated-inputs-fuzz-targets-benchmarks-and-budgets.md](../rfc/0006-generated-inputs-fuzz-targets-benchmarks-and-budgets.md) |
| assert-spec#21, a recursion through `using` | <https://github.com/dokimasia/assert-spec/issues/21> |
| RFC-0011 of assert-spec, a recursive shape | <https://github.com/dokimasia/assert-spec/blob/main/docs/rfc/0011-property-forms-and-derived-inputs.md> |
| RFC-0037 of assert-spec, a stated generator inside itself | <https://github.com/dokimasia/assert-spec/blob/main/docs/rfc/0037-the-values-of-a-generator-under-using.md> |
| The strategies `recursive` and `deferred` of Hypothesis | <https://hypothesis.readthedocs.io/en/latest/reference/strategies.html> |
| `Deferred` and `Make` of Rapid | <https://pkg.go.dev/pgregory.net/rapid#Deferred> |
