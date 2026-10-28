---
adr: 0015
title: A derived check is falsified on a wrapper of a fresh subject of its application
status: Accepted
date: 2026-10-28
supersedes: RFC-0003 and ADR-0007 in part
superseded-by: ADR-0017 in part
rfc: 0014
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0015: A derived check is falsified on a wrapper of a fresh subject of its application

## Status

Accepted

## Context

RFC-0003 builds the planted defect of each derived check from the typed double of the declaration, with the callables of the defect form replaced. It rejects a defect that wraps the subject of the consumer, because falsification would then need that subject. ADR-0007 removed the double from every suite, because no suite had a derived check and no generated file referred to the double.

A double is not a correct implementation. It stores nothing and returns what a test configures. RFC-0003 requires each planted defect to pass every other check of its suite, apart from the checks that share it and the deeper checks of its edges. On a defect built from a bare double, every check whose law depends on state fails, such as `read-after-write` on the defect of `zero-on-error`. Only a correct implementation keeps the other checks passing.

The reason codes remove the risk that made RFC-0003 reject a wrapper. A partner in the subject of the consumer that fails a check for a reason outside the check reports another reason code, and falsification reports `wrong-reason` instead of a proof.

## Decision

We will build the planted defect of a derived check as a wrapper of a fresh subject of the application, and falsify each derived check inside the application.

- The template of a defect form writes a type that embeds the subject type and overrides the callables that the form changes.
- The handle gains `Fresh`, which returns a fresh subject from the factory of the application, also under falsification. `Subject` and `New` keep their meaning: under falsification they return the planted defect.
- `conform.Run` falsifies each derived check on its defect over a fresh subject before it applies the check to the subject, as it already falsifies a check that `conform.With` adds. It writes the falsification record and the record of the subject.
- `conform.Falsify` falsifies the checks of the consumer, whose defects the consumer writes, and the ceilings. It does not falsify a derived check, because the generated test of a suite has no subject.
- The isolation of RFC-0003 is unchanged: each other check passes on the defect, apart from the checks that share it and the deeper checks of its edges.
- The suite does not ask for the double of its declaration, as ADR-0007 decided.

## Alternatives Considered

### The double of RFC-0003

The suite asks for the double again, and each defect replaces the callables of its form on the double. Every check that depends on state fails on every defect, so the isolation of RFC-0003 fails for most claims of the spec.

### The double that delegates to the generated fake

Each defect is the double, delegating to the verified fake of the declaration, with the callables of its form replaced. Falsification needs no subject, but only a declaration whose contracts have a reference form has a fake. Every other declaration would refuse its derived checks.

## Consequences

**Positive:**

- A planted defect is correct everywhere apart from the callables of its form, so claims that depend on state pass the isolation of RFC-0003.
- A declaration needs neither a double nor a fake for its derived checks.

**Negative:**

- A derived check is proven only where an application runs, in the tests of the consumer and in the corpus. The generated test of a suite proves the checks of the consumer and the ceilings alone.
- A subject whose callable is wrong can fail the falsification of another check. The reason code then reports `wrong-reason`, and the consumer reads the failure of the subject in the record of the same application.

## References

| What | Where |
|---|---|
| RFC-0003, the planted defects built from the double and the rejected wrapper | [0003-the-suite-generator.md](../rfc/0003-the-suite-generator.md) |
| ADR-0007, a suite does not ask for the double of its declaration | [0007-a-suite-does-not-ask-for-the-double-of-its-declaration.md](0007-a-suite-does-not-ask-for-the-double-of-its-declaration.md) |
| RFC-0014, the derivation of checks | [0014-the-derivation-of-checks.md](../rfc/0014-the-derivation-of-checks.md) |
