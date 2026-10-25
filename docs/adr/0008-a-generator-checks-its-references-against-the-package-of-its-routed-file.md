---
adr: 0008
title: A generator checks its references against the package of its routed file
status: Accepted
date: 2026-10-25
supersedes: RFC-0003 in part
superseded-by: none
rfc: 0003
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0008: A generator checks its references against the package of its routed file

## Status

Accepted

## Context

RFC-0003 writes the conformance test of a subject in the external test package of its directory, so the factory of the directive `subject` is an exported function of the package of the struct, or a function of its external test package. The suite generator compared the package of the factory with the external test package of the directory of the struct.

- Issue 24 reports a struct under the directive `out`, whose conformance test is in the external test package of the `out` directory. The suite generator refused a factory of that package, which the test can call. It accepted a factory of the external test package of the struct, which the test cannot import.
- A struct that embeds `*assert.Recorder` promotes `Clock() matcher.Clock`, whose result is in an internal package of assert. The suite of such a struct spells that type in its interface of the methods. The layout of eidos refused the interface under EID-0082, and the template still spelled the type from the data of the record. The render then withheld the whole file under EID-0023, with the suites of every other declaration of the file.
- eidos derives the package of each routed file. `Out.Package` previews it for the handle of a generator, and `Bound.Importable` checks a reference against the import rule of the target. Both are in eidos-core v0.8.0.

## Decision

We will check each reference that a generator writes against the package of its routed file, before the generator records it.

- The suite generator checks the factory of a subject against the package of the conformance test. It reports `BadFactory` where the test cannot call the factory, with the reason of the import rule.
- The suite generator checks each parameter and each result of each callable of an interface or a struct against the package of its suite. It reports `Unimportable` at the declaration where the suite cannot refer to one, and does not record the suite. A function is its own callable, and its package can refer to each type of its signature.
- The stub generator checks each method of a double against the package of the file of the double. It reports `Unimplementable` where the double cannot refer to a parameter or a result.
- Where the plan previews no package, a generator records as before, and the layout checks each reference itself.

## Alternatives Considered

### A suite without the callable

The suite of a struct could leave out each callable whose signature its file cannot spell, and test the others. The subject of the declaration is in another package, so the test of a subject cannot preview the package of the suite, and its levels and excusals could refer to a callable that the suite leaves out.

### The rule in the generator

The Go part of a generator could compare import paths itself, as the suite generator compared the external test package. The rule depends on the route of the file, which eidos derives and a plugin cannot.

## Consequences

**Positive:**

- A struct under `out` takes a factory of the external test package of its kit.
- A declaration that its suite cannot spell gets a finding at the declaration, in place of a template error at its file.

**Negative:**

- A struct that promotes a method with a type of an internal package of another tree has no suite.

## References

| What | Where |
|---|---|
| RFC-0003, the factory of a subject and the suite of a struct | [0003-the-suite-generator.md](../rfc/0003-the-suite-generator.md) |
| Issue 24, a struct under out cannot take its factory from the external test package of its kit | <https://github.com/dokimasia/dokimi/issues/24> |
| ADR-0105 of eidos, the import rule of the target | <https://github.com/dokimasia/eidos/blob/main/docs/adr/0105-the-layout-refuses-a-declaration-whose-file-cannot-import-a-reference.md> |
| ADR-0107 of eidos, the package of the routed file of a handle | <https://github.com/dokimasia/eidos/blob/main/docs/adr/0107-a-generator-previews-the-package-of-the-routed-file-of-a-handle.md> |
