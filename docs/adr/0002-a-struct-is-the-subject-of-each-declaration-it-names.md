---
adr: 0002
title: A struct is the subject of each declaration that it names
status: Accepted
date: 2026-10-23
supersedes: RFC-0003 in part, RFC-0006 in part
superseded-by: none
rfc: 0003
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0002: A struct is the subject of each declaration that it names

## Status

Accepted

## Context

RFC-0003 gives a struct one directive `subject`, and eidos refuses a second one with EID-0014. A struct that implements two or more interfaces under the directive `suite` gets the generated conformance test of one of them, and its spec applies the other suites by hand. Issue 12 reports such a struct in theke: `crashfs.FS` implements `vfs.FS`, `vfs.Linker` and `vfs.Direct`.

- RFC-0006 names the tests of a subject after the struct alone: `Test<Subject>Allocs`, `Benchmark<Subject>`, `Benchmark<Subject>Parallel` and `Fuzz<Subject><Callable>`. The subjects `vfs.FS`, `vfs.Linker` and `vfs.Direct` of `crashfs.FS` would each declare `BenchmarkFS`. The conformance test itself is `Test<Subject>ConformsTo<Declaration>` already.
- The directives `excuse` and `budget` on a struct name a callable with its declaration, as `Writer.Write/allocation-budget` and `Reader.Read`.

## Decision

We will make the directive `subject` repeatable, with one conformance test for each declaration.

- A struct takes one directive `subject` for each declaration under the directive `suite` whose checks it passes. A second directive of one declaration reports `RepeatedSubject` at that directive, and the first directive keeps its test.
- Each name of the tests of a subject joins the struct and the declaration: `Test<Subject>ConformsTo<Declaration>Allocs`, `Benchmark<Subject>ConformsTo<Declaration>`, `Benchmark<Subject>ConformsTo<Declaration>Parallel` and `Fuzz<Subject>ConformsTo<Declaration><Callable>`. A struct with one directive `subject` takes these names too.
- The test of each directive applies the excusals and the budgets whose callable belongs to its declaration. The first directive reads each excusal and each budget whose callable belongs to no declaration of the struct, and reports it. Every test of the struct is then left out.

This decision supersedes the single directive `subject` of RFC-0003 and the names of the tests of a subject of RFC-0006.

## Alternatives Considered

### The names of the struct alone for a struct with one subject

The names would change when a consumer adds a second directive. A name that depends on the other directives of the struct is harder to find from the directive that produced it.

### One directive with a list of declarations

`//dokimi:subject vfs.FS vfs.Linker new=newFS` would share one factory and one number of cases between all declarations. A subject of two declarations can need another factory for each, and the directives `excuse` and `budget` already name their declaration.

## Consequences

**Positive:**

- A struct gets the conformance test, the allocation test, the benchmarks and the fuzz targets of each suite that it passes.
- No two tests of one struct have one name.

**Negative:**

- Every test of a subject of v0.1.0 has a new name, so a consumer that runs one test by its name changes the pattern.
- The tests of two subjects of one struct share the origin of the struct. dokimi requires eidos at 01e3506c, which gives each record the data of its own call (eidos issue 76).

## References

| What | Where |
|---|---|
| RFC-0003, the suite generator, whose single directive `subject` this decision replaces | [0003-the-suite-generator.md](../rfc/0003-the-suite-generator.md) |
| RFC-0006, whose names of the tests of a subject this decision replaces | [0006-generated-inputs-fuzz-targets-benchmarks-and-budgets.md](../rfc/0006-generated-inputs-fuzz-targets-benchmarks-and-budgets.md) |
| Issue 12, a struct takes one directive subject | <https://github.com/dokimasia/dokimi/issues/12> |
| eidos issue 76, two records of one origin in a unit read the data of the first | <https://github.com/dokimasia/eidos/issues/76> |
| eidos ADR-0103, a record keeps the data of the call that places it | <https://github.com/dokimasia/eidos/blob/main/docs/adr/0103-a-record-keeps-the-data-of-the-call-that-places-it.md> |
