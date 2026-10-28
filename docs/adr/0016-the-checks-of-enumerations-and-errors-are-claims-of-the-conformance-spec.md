---
adr: 0016
title: The checks of enumerations and errors are claims of the conformance spec
status: Accepted
date: 2026-10-28
supersedes: RFC-0004 in part
superseded-by: none
rfc: 0004
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0016: The checks of enumerations and errors are claims of the conformance spec

## Status

Accepted

## Context

RFC-0004 writes the checks of an enumeration and of the errors of a package as fixed templates that assert through the assert library, without probes and without planted defects. It rejects claims of the conformance spec for them, so that the generated tests equal the tests that testkit writes for the same source. Its drawbacks state that no planted defect shows that each of these checks can fail.

RFC-0001 requires a planted defect for every check, and RFC-0003 already gives a claim the generator `enum` or `sentinel`. The parity with the output of testkit served the move from testkit, and the corpus of dokimi now replaces testkit as the reference.

## Decision

We will state the checks of the enum generator and the sentinel generator as claims of the conformance spec, with probes, reason codes and planted defects.

- Each level of the tables of RFC-0004 becomes one claim with the generator `enum` or `sentinel`. A claim of either generator has no licences and binds to each subject of its generator, as RFC-0003 states.
- The inputs of these claims are the sources of RFC-0003 for enumerations and errors: the members, the name and parse functions, the error values, the error types and the causes.
- The params of the directive `sentinel`, `prefix` and `no-overlap-with`, are params of the claims that read them. `prefix=off` withdraws the claims of the opening.
- The decisions of ADR-0003, ADR-0004 and ADR-0005 remain in force. An unexported enumeration has no checks, the encode of an enumeration refuses a value outside the variants, and a remote error type has the message of another process.
- The generated test of an enumeration or of the errors of a package calls `conform.Falsify` and the application of its suite, and keeps one test per enumeration and per error, as RFC-0004 names them.

## Alternatives Considered

### Fixed templates with planted defects

The templates of RFC-0004 remain, and each gains a planted defect. Each runtime would then implement a second kind of proof beside the probes and the reason codes, and the corpus gate of each language would compare checks that the spec does not state.

## Consequences

**Positive:**

- Every check of dokimi has a planted defect and a reason code, and the corpus gate compares the checks of enumerations and errors in every language.

**Negative:**

- The generated tests of enumerations and errors differ from the tests that testkit writes for the same source.
- The spec gains the claims, probes, reason codes and defect forms of the 12 levels of the two tables.

## References

| What | Where |
|---|---|
| RFC-0004, the checks of enumerations and errors | [0004-doubles-builders-enumerations-and-errors.md](../rfc/0004-doubles-builders-enumerations-and-errors.md) |
| RFC-0003, the claims of the generators enum and sentinel | [0003-the-suite-generator.md](../rfc/0003-the-suite-generator.md) |
