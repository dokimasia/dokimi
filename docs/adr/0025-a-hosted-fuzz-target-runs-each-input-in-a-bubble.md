---
adr: 0025
title: A fuzz target of a hosted application runs each input in a bubble
status: Accepted
date: 2026-11-13
supersedes: RFC-0045 in part
superseded-by: none
rfc: 0045
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0025: A fuzz target of a hosted application runs each input in a bubble

## Status

Accepted

## Context

Under the option `Hosted` of conform, `Run` opens the subject of each check through the door Host and runs each check in `Bubble`. `Fuzz` builds the subject of each input through the factory of the consumer and runs the input outside a bubble (conform v0.17.1). For each input, a property check that sleeps on the clock of the host of its subject sleeps on the clock of the process. `synctest.Wait` panics in a check outside a bubble.

The membership of the quorum store of sharder runs three processes of etcd-io/raft. Its property check `Membership.Join/views-list-the-live-sessions` advances the clock of its subject and then calls `synctest.Wait`. Outside a bubble the operation `Advance` skips the input, so no input of the fuzz target of `Join` reaches an expiry or a lost lease.

The property engine applies one wrapper of `prop.Around` to every input of a fuzz target. `Fuzz` leaves out each check with an `Around`, and under `Hosted` each other check runs in `Bubble`, so one wrapper serves every check that `Fuzz` runs.

## Decision

We will open the subject of each input of a fuzz target of an application with `Hosted` and without `Tasks` through the door Host in `Bubble`, because `Run` runs the same checks in that environment.

- `Fuzz` takes the factory of `config.subjects`, as `Run` does. The handle of each input starts a host from the case of the input, and the end of the case ends the lifetime of its incarnation.
- `Fuzz` passes `prop.Around(Bubble)` to `prop.Fuzz` under `Hosted` without `Tasks`. `prop.Fuzz` hands the wrapper the test of each input, the inputs of the seed corpus included.
- An application with `Tasks` keeps the factory of the consumer for its fuzz target, as it keeps it for its checks.
- `RunBudgets`, `FalsifyAllocs` and `Bench` keep the factory of the consumer, because a budget and a benchmark measure on the clock of the process.
- The TypeScript runtime does the same: under `hosted` without `tasks`, `fuzz` opens the subject of each input through the door `hosted` and runs the input inside `fakeTimers`, as `run` runs each check.

## Alternatives Considered

### The factory of the consumer

The fuzz target can keep the factory of the consumer outside a bubble, as it does at conform v0.17.1. A fuzz target of a subject that only a host opens then checks another subject than `Run` checks. A property check that waits for the quiescence of a bubble cannot run under `-fuzz` at all.

## Consequences

**Positive:**

- The fuzz target of a hosted application checks the subject and the environment that `Run` checks.
- A property check that sleeps on the clock of its host and waits for the quiescence of the bubble runs under `-fuzz`.

**Negative:**

- Each input starts a host and runs in a bubble, so an input takes longer than over the factory of the consumer.
- The factory of the consumer builds no subject of a fuzz target of `Hosted` without `Tasks`.
- `Fuzz` still leaves out each check with an `Around`.

## References

| What | Where |
|---|---|
| A fuzz target of a Hosted application runs its inputs outside a bubble | <https://github.com/dokimasia/dokimi-go/issues/51> |
| The checks of a subject that only a host opens | [RFC-0045](../rfc/0045-the-checks-of-a-subject-that-only-a-host-opens.md) |
| `prop.Around` and the inputs of `prop.Fuzz`, assert v0.22.0 | <https://github.com/dokimasia/assert-go/blob/v0.22.0/prop/option.go#L225-L267> |
