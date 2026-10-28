---
adr: 0017
title: The corpus isolates the planted defects of derived checks
status: Accepted
date: 2026-10-28
supersedes: ADR-0015 in part
superseded-by: none
rfc: 0014
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0017: The corpus isolates the planted defects of derived checks

## Status

Accepted

## Context

The planted defect of a derived check wraps a fresh subject of its application. `conform.Run` falsifies each derived check on its defect before it applies the check to the subject. Under RFC-0003 the falsification also isolates the defect: every other check must pass on it, apart from the checks that share it and the deeper checks of its edges. It also checks each edge in both directions.

Isolation proves that the template of a defect form breaks one law alone. The proof is valid only over a correct subject. The subject of a consumer is not known to be correct. A subject that fails one check fails it again under every planted defect that wraps the subject. Isolation over that subject reports N−1 failures of planted defects before the application reports the defect of the subject. Isolation also runs N×N checks for a suite of N derived checks, such as 900 runs for 30 checks.

The template of a defect form is the same for every subject. The controls of the corpus are correct subjects that differ from the fixtures.

## Decision

We will isolate the planted defects of derived checks and check their edges only under the option `conform.Isolate`, which the corpus passes for each control.

- `conform.Run` without `Isolate` runs each derived check on its planted defect and then on the subject.
- `conform.Run` with `Isolate` also runs every other derived check on each planted defect, and checks each edge in both directions.
- The falsification record reports a failed isolation as `not-isolated`, and a failed edge as `not-deeper` or `equivalent`, with the identity of the other check.
- `conform.Falsify` still isolates the defects of the checks of the consumer and checks their edges.

## Alternatives Considered

### Isolation in every application

Every `conform.Run` isolates every derived defect, as ADR-0015 decided. A defect of the subject then appears as a failure of each other planted defect, and each application costs N×N runs of checks.

### Isolation only over the checks that pass on the subject

`conform.Run` applies every check to the subject first, then isolates each defect over the checks that passed. The blame is right, but each application still costs N×N runs, and it proves again what the corpus proves for the same template.

## Consequences

**Positive:**

- A failure of the application of a consumer names the check that the subject fails, and no planted defect.
- The falsification of a derived check costs one run of the check per application.

**Negative:**

- A defect form that breaks a second law only over some subjects passes the corpus. The verdict reverses when a defect form passes isolation on every control and fails it on a subject that passes its own suite.
- Every claim needs a fixture with a control in the corpus before its isolation is proven.

## References

| What | Where |
|---|---|
| ADR-0015, a derived check is falsified on a wrapper of a fresh subject | [0015-a-derived-check-is-falsified-on-a-wrapper-of-a-fresh-subject.md](0015-a-derived-check-is-falsified-on-a-wrapper-of-a-fresh-subject.md) |
| RFC-0014, the derivation of checks | [0014-the-derivation-of-checks.md](../rfc/0014-the-derivation-of-checks.md) |
| RFC-0015, vacuous checks and controls | [0015-vacuous-checks-and-controls.md](../rfc/0015-vacuous-checks-and-controls.md) |
