---
adr: 0018
title: The tests of enumerations and errors are fixed templates
status: Accepted
date: 2026-10-29
supersedes: ADR-0016 in part
superseded-by: none
rfc: 0004
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0018: The tests of enumerations and errors are fixed templates

## Status

Accepted

## Context

ADR-0016 states the checks of the enum generator and the sentinel generator as claims of the conformance spec. The generated test of an enumeration or of the errors of a package applies a suite of those claims through `conform.Run`, which falsifies each check on a planted defect that wraps a fresh subject.

The probes of these claims are generic. Their inputs are the variants, the render and the parse function of an enumeration, and the error values, the error types and the causes of a package. The specs of the probes in the runtime already prove that each probe fails on a planted defect. A falsification for each enumeration and each package of each consumer proves the same again.

The enum generator writes most of the functions that these tests check, and its golden files pin them. A falsification for each subject needs these additions:

- a struct of the parts of each enumeration and of each package of errors
- probes that take the encode and the decode as functions
- the facts of enumerations and errors in the derivation

The templates have eight checks without a claim. They are plain assertions in either design.

## Decision

We will keep the generated tests of enumerations and errors as the fixed templates of RFC-0004, without the claims of the generators `enum` and `sentinel`.

- The conformance spec and the runtime keep the claims, probes, reason codes and defect forms of the two generators.
- A generated test of an enumeration or of the errors of a package does not call `conform.Falsify` or `conform.Run`.
- The decisions of ADR-0003, ADR-0004 and ADR-0005 remain in force.

## Alternatives Considered

### The claims through conform.Run

The generated test applies a suite whose subject is a struct of the parts of the enumeration or of the errors, as ADR-0016 decided. Each check has a planted defect for each subject. Its proof adds nothing to the specs of the probes, because each probe takes the same kinds of inputs for every subject. Its cost is a struct of parts for each subject, two changed probes and the derivation of enumerations and errors.

## Consequences

**Positive:**

- The enum and sentinel generators and their golden files do not change.
- The derivation has no facts of enumerations and errors, and the probes of the text encoding keep their signatures.

**Negative:**

- A generated check of an enumeration or of the errors of a package has no planted defect for its subject. The verdict reverses when a defect appears that only a falsification for each subject catches, such as a test that binds the wrong function into a check.
- The corpus gate compares the checks of enumerations and errors by the identities of the claims alone.

## References

| What | Where |
|---|---|
| ADR-0016, the checks of enumerations and errors are claims of the conformance spec | [0016-the-checks-of-enumerations-and-errors-are-claims-of-the-conformance-spec.md](0016-the-checks-of-enumerations-and-errors-are-claims-of-the-conformance-spec.md) |
| RFC-0004, the checks of enumerations and errors | [0004-doubles-builders-enumerations-and-errors.md](../rfc/0004-doubles-builders-enumerations-and-errors.md) |
