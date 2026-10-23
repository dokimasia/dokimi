---
adr: 0001
title: Every generated test is a black-box test
status: Accepted
date: 2026-10-23
supersedes: RFC-0003 in part
superseded-by: none
rfc: 0003
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0001: Every generated test is a black-box test

## Status

Accepted

## Context

The family of the tests of the suite generator writes its Go files into the external test package of their directory. RFC-0003 makes one exception in its section on the generated suite. Where the factory of the directive `subject` is not exported, the conformance test of the subject is in `<file>_suite_internal_test.go` inside the package of the subject.

- The testing standard of assert-go requires that every test file declares the package `x_test` for its package `x`. A test inside the package reaches its unexported state, so a refactor of the package can break the test without a change of its API.
- A conformance kit often imports an implementation of its own interface, such as the fake behind its planted defects. A test of that implementation inside its package imports the kit, the kit imports the implementation, and Go refuses the import cycle. Issue 9 reports three such kits.
- dokimi requires eidos at e4ae6a9a since 53d5881. Its Go rules resolve a bare name of the param `new` in the external test package of the directory. A consumer can then declare the factory there by its name alone, as `new=newSubject`, without a constructor for tests in its API.

## Decision

We will generate black-box tests alone.

- The factory of the directive `subject` is an exported function of the package of the struct, or a function of the external test package of its directory.
- A factory that the package of the struct does not export reports `BadFactory` at the subject: "the factory newHidden is not exported, so Hidden has no conformance test. Export it, or declare it in the external test package". suite does not record a conformance test for the subject.
- No generated test declares the package under test. suite does not mark a declaration of its family of tests as internal.

This decision supersedes the sentence of RFC-0003 on `<file>_suite_internal_test.go`.

## Alternatives Considered

### The internal test file for an unexported factory

RFC-0003 chose it because only the package of the subject could call such a factory. Since 53d5881 the external test package can declare the factory by its name alone. The internal file also keeps the import cycle of issue 9 and tests the package from inside.

### A generated exported wrapper of the factory

The generator would add an exported function to the package of the consumer, such as `NewHiddenForTest`. That function is a constructor for tests in the API of the package, which issue 9 asks to avoid.

## Consequences

**Positive:**

- Every generated test of the suite generator is in the external test package of its directory, as the testing standard requires.
- The conformance test of a subject cannot form an import cycle with the kit of its suite.

**Negative:**

- A subject of v0.1.0 whose factory is not exported reports `BadFactory`. Its consumer moves the factory into the external test package or exports it.

## References

| What | Where |
|---|---|
| RFC-0003, the suite generator, whose sentence on the internal test file this decision replaces | [0003-the-suite-generator.md](../rfc/0003-the-suite-generator.md) |
| The testing standard of assert-go | <https://github.com/dokimasia/assert-go/blob/main/docs/standards/testing.md> |
| Issue 9, a subject whose package its kit imports cannot get the generated test of a subject | <https://github.com/dokimasia/dokimi/issues/9> |
