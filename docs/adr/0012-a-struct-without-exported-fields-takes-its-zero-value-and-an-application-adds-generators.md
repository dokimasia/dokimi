---
adr: 0012
title: A struct without exported fields takes its zero value, and an application adds generators
status: Accepted
date: 2026-10-26
supersedes: RFC-0006 and ADR-0009 in part
superseded-by: none
rfc: 0006
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0012: A struct without exported fields takes its zero value, and an application adds generators

## Status

Accepted

## Context

RFC-0006 leaves a type that the shape reader of the property engine refuses to a registration of the test process, and ADR-0009 keeps that rule for every type except `any` and the pointer shapes. The shape reader refuses a struct without an exported field. A suite draws the values of its inputs at every depth, so each property run, budget, fuzz target and benchmark of a callable whose input contains such a struct fails before its first case. Issue 26 reports two such types of theke: `schema.Horizons` and `config.Resolution`.

Each remedy costs the author of the type or of the test:

- A registration applies to the whole test binary, and a second registration of one type panics. Each test binary of a consumer then registers the type, and two kits of one binary cannot both register it.
- A sample of the type has a value in the scope of the type, so the package of the type exports a constructor for its tests alone.

Code outside the package of such a struct builds its zero value alone, unless the package exports a constructor.

## Decision

We will give a struct without exported fields its zero value as its starting value, and let an application add generators of its own.

- A struct outside the standard library without an exported field takes its zero value as its starting value, where no authored sample, builder or declared default states one. The suite passes the generator of that value through `prop.Using`, as it passes the generator of every other starting value. ADR-0010 applies: a struct that the file of the suite cannot import keeps its registration.
- conform gains an option of an application that adds generators through `prop.Using` to the property runs, the fuzz targets and the benchmarks of the application, such as `conform.Using`. The application passes them after the overrides of the suite, so a generator of the option replaces the override of the suite for its type.
- The tests that the suite generator writes for a suite take the zero value. Their subjects are planted defects and idle subjects, which do not read the inputs.

## Alternatives Considered

### A registration of the test process

ADR-0009 leaves such a type to a registration. Each test binary then registers the type, and two kits of one binary cannot both register it.

### A sample in the kit

A sample could have a value in the scope of the package of the directive `out`, as issue 26 suggests. eidos would read the value of a directive in the scope of another package, and the author would write a value in a directive where a generator in Go serves.

### A fallback below a registration

The property engine could take the zero value where neither `using` nor a registration states a generator. A registration then still applies. The engine of every language would need a new level of precedence, and an application would still need a way to state generators of its own.

## Consequences

**Positive:**

- Each generated test of a callable whose input contains a struct without exported fields runs.
- A consumer states values of such a type for its application in Go, such as a `config.Resolution` that its registry builds.

**Negative:**

- A registration of such a type no longer applies in the runs of a suite, because `prop.Using` takes precedence over a registration. The option of the application replaces it.
- A struct whose zero value is no valid input fails the checks of its callables until its declaration states a sample or its application states a generator.

## References

| What | Where |
|---|---|
| RFC-0006, the values of the parameters of a callable | [0006-generated-inputs-fuzz-targets-benchmarks-and-budgets.md](../rfc/0006-generated-inputs-fuzz-targets-benchmarks-and-budgets.md) |
| ADR-0009, the types that a registration states | [0009-a-suite-draws-a-value-of-any-type-and-the-pointer-shapes-of-the-property-engine.md](0009-a-suite-draws-a-value-of-any-type-and-the-pointer-shapes-of-the-property-engine.md) |
| ADR-0010, the references that the file of a suite can import | [0010-a-suite-refers-only-to-declarations-that-its-file-can-import.md](0010-a-suite-refers-only-to-declarations-that-its-file-can-import.md) |
| Issue 26, a field whose type has no exported field | <https://github.com/dokimasia/dokimi/issues/26> |
