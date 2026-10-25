---
adr: 0009
title: A suite draws a value of any type and the pointer shapes of the property engine
status: Accepted
date: 2026-10-25
supersedes: RFC-0006 in part
superseded-by: none
rfc: 0006
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0009: A suite draws a value of any type and the pointer shapes of the property engine

## Status

Accepted

## Context

RFC-0006 gives each parameter of a callable a value by its kind. A pointer without the mixin `nil-safe` takes a pointer to a drawn value of its type. An interface of the standard library, and any other type that the shape reader of the property engine does not read, take the generator that a registration of the test process states.

- The shape reader reads `*big.Int` as an integer, `*big.Rat` as a decimal and `*time.Location` as a zone. It does not read a value of `big.Int`, `big.Rat` or `time.Location`, because none of them has an exported field. The suite drew the value and took its address, so each call of a callable that takes such a pointer failed with "time.Location has no exported field to read".
- The shape reader refuses `any`, an interface without registered variants. A suite cannot register `any`: a second registration of one type panics, and two suites of one test binary would both register it. So each call of a callable that takes `any`, such as `Print(msg string, values ...any)`, failed at its first case.
- A registration of the test process can state the generator of any other type that the shape reader refuses. dokimi does not see such a registration.

## Decision

We will draw a value of any type and the pointer shapes of the property engine in the suite, and leave every other refused type to a registration.

- A pointer to `big.Int`, `big.Rat` or `time.Location` takes the pointer that the property engine draws.
- A suite whose calls take a value of `any` at some depth of their inputs passes `prop.Using` of a generator of `any` to each property run and each benchmark. The generator draws a string, an int or a bool.
- The suite does not refuse another type that the shape reader refuses. A registration of the test process states its generator, as RFC-0006 designs.

## Alternatives Considered

### Every pointer drawn as a pointer

The suite could draw `prop.Of[*T]()` for each pointer, with a filter of the nil pointers. Each pointer parameter then rejects the cases with nil, and a callable with three pointer parameters rejects most of its cases.

### A refusal of a type without exported fields

The suite could report an Error at a declaration with a callable that takes a type without exported fields. A registration of the test process can state the generator of such a type, and the refusal would forbid that registration.

### Fixed values of any

`prop.SampledFrom[any]("text", 1, true)` draws three values alone. The generator of the decision draws from the whole domain of each of the three types.

## Consequences

**Positive:**

- A callable that takes a `*time.Location`, a `*big.Int` or a `*big.Rat` gets calls without a registration.
- A callable that takes `any`, a variadic `...any` or a struct with a field of `any` gets calls without a registration.

**Negative:**

- The generator of `any` of a suite takes precedence over a registration of `any` in the test process, because `prop.Using` takes precedence over a registration.
- The draw of `any` needs the property engine to apply `prop.Using` to a type that its shape reader refuses, which assert-go#67 requests.

## References

| What | Where |
|---|---|
| RFC-0006, the values of the parameters of a callable | [0006-generated-inputs-fuzz-targets-benchmarks-and-budgets.md](../rfc/0006-generated-inputs-fuzz-targets-benchmarks-and-budgets.md) |
| assert-go#67, `prop.Of` refuses a type before `prop.Using` applies | <https://github.com/dokimasia/assert-go/issues/67> |
