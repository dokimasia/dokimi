---
milestone: 0009
title: A check can make a call fail on demand
status: Planned
depends-on: 0003
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0019
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0009: A check can make a call fail on demand

## Goal

A Go team whose subject depends on another declaration gets checks that the subject returns a failure of the dependency to its caller, leaves no partial write, and succeeds again once the dependency recovers.

## Done when

- [ ] RFC-0019 is accepted.
- [ ] conform has the `Injector`, `FaultOf`, the door `Faults` and the option `Faultable`.
- [ ] A cache over a store in the corpus under `//dokimi:suite machineries=faults` passes `failure-is-reported`, `no-partial-write` and `retry-succeeds-under-faults` with the door, and each claim fails on its planted defect with its reason code. `retry-succeeds-under-faults` and `retry-succeeds` pass the proof of their depth pair.
- [ ] The same fixture without `machineries=faults` lists the classifications of the claims as unclaimed with the stance `needs`.
- [ ] A check of these claims on an application without `Faultable` fails with the name of the option.

## Why now

The claims are derived checks, so milestone 0003 has to exist first. It is the cheapest machinery above a single caller, because the fault strategies and the doubles that arm them already exist. The simulation of milestone 0015 arms the faults of its dependencies through the same `Injector`.

## Scope

The injector and the door in conform, the claims and defect forms of a failing dependency in the spec, and the fixtures of the corpus that use them.

## Not in this milestone

- Faults at the network or in the source of the subject: refused in RFC-0019.
- The same claims in TypeScript: nobody has scheduled them yet.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 9 | The first roadmap of dokimi |
