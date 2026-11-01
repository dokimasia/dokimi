---
adr: 0019
title: A generated binding draws the inputs of its calls in its check
status: Accepted
date: 2026-10-31
supersedes: RFC-0006 and RFC-0012 in part
superseded-by: none
rfc: 0006
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0019: A generated binding draws the inputs of its calls in its check

## Status

Accepted

## Context

A generated binding is a derived check whose calls take inputs that the property engine generates and shrinks. It generates every input of every call of its probe. The calls of one case share the value of each parameter with one key and one type, so a write and its read use one key.

The call generator of each callable draws the inputs of one call and returns a function of the form `func(context.Context, S) error`. The fuzz targets, the benchmarks and the allocation budgets draw their calls from it. A parameter with an authored sample takes the sample in it.

A call of that form returns no value. A probe that reads the result of a call, such as the probe of `zero-on-error` or of `read-after-write`, cannot read it through a drawn call. Two calls from two call generators also draw two keys, so a write and its read do not share one.

## Decision

We will draw the inputs of the calls of a generated binding in the check of the binding, once for each case and for each pair of a key and a type, because a drawn call of the call generator returns no value and draws a key of its own.

- The check draws each input by the rules of the call generator: the sample of a parameter, the generated functions and implementations, and the overrides of the suite.
- Its closures are the closures of the sample binding, with the drawn values in place of the samples.
- The call generators keep their form, and the fuzz targets, the benchmarks and the budgets keep drawing their calls from them.

## Alternatives Considered

### A second call generator that returns the results

Each callable gets a second call generator, whose calls return the results of the callable. A claim of one call draws one call from it and reads the results.

It lost because two calls from two generators still draw two keys. A claim with a partner, such as `read-after-write`, would keep its sample binding alone, and each callable would get a second generator.

## Consequences

**Positive:**

- A probe gets the results of its calls, so `zero-on-error` and the other claims that read a result get a generated binding.
- A write and its read share their key, so a claim with a partner gets a generated binding too.
- The fuzz targets, the benchmarks and the budgets keep one call generator for each callable.

**Negative:**

- The template states the rules that draw a parameter twice: in the call generator and in the check of each generated binding. A change of the rules changes both.
- Each generated binding of a callable repeats the draws of the callable, so the generated code grows with the number of claims of the callable.

**Neutral:**

- A parameter with an authored sample takes the sample in a generated binding, as it does in its call generator.

## References

| What | Where |
|---|---|
| RFC-0006, the generated bindings and the call generators | [0006-generated-inputs-fuzz-targets-benchmarks-and-budgets.md](../rfc/0006-generated-inputs-fuzz-targets-benchmarks-and-budgets.md) |
| RFC-0012, the sample of a parameter in the call generator | [0012-the-sample-of-a-parameter.md](../rfc/0012-the-sample-of-a-parameter.md) |
| The property engine of assert-go | <https://github.com/dokimasia/assert-go/blob/main/prop/doc.go> |
