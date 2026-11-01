---
milestone: 0009
title: A check can make a call fail on demand
status: In progress
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

A Go team whose subject depends on another declaration gets checks that the subject returns a failure of the dependency to its caller, does not leave a partial write, and succeeds again once the dependency recovers. A breaker, a balancer and a callable that acquires resources get the checks of their laws under a failing dependency.

## Done when

- [x] RFC-0019 is accepted.
- [x] conform has `ErrInjected`, the `Injector` with its `Dependency` and `FaultOf`, the door `Faults` with the option `Faultable`, and `Handle.Injector`.
- [x] A handle builds the subjects of a check of the machinery `call-that-fails-on-demand` through the door `Faults`, and such a check on an application without `Faultable` fails with the name of the option.
- [x] `//dokimi:suite machineries=[faults]` gives a declaration the claims of RFC-0019. An unknown value is an Error under `UnknownMachinery`, and a declaration without the param lists the classifications of the claims as unclaimed with the stance `needs`.
- [x] The directive `check` takes the param `faults`, and the generator writes it into the field `Doors.Faults` of the check.
- [x] Fixtures of the corpus under `iface/door` have a cache over a store, a transaction over a store, a circuit breaker over an operation, a balancer with the mixin `sticky` over two instances, and a callable with the mixin `leak-free`. Their subjects pass `failure-is-reported`, `no-partial-write`, `retry-succeeds-under-faults`, `transaction-under-faults`, `fails-fast`, `keeps-its-instance` and `releases-on-every-path`, and each claim fails on its planted defect with its reason code.
- [x] `retry-succeeds-under-faults` and `retry-succeeds` pass the proof of their depth pair.
- [ ] `dokimi conform --lang go` passes with the expected outcomes of the new fixtures.

## Why now

The claims are derived checks, so milestone 0003 has to exist first. It is the cheapest machinery above a single caller, because the fault strategies and the doubles that arm them already exist. The simulation of milestone 0015 arms the faults of its dependencies through the same `Injector`, and the claim of `single-flight` of milestone 0011 keeps a computation waiting through it.

## Scope

The injector, the door and its option in conform, the seven probes of the package `probe`, the claims, laws, reasons and defect forms of RFC-0019 in the spec, the param `machineries` of the directive `suite` and the param `faults` of the directive `check`, the defect forms of the suite generator, and the fixtures of the corpus with their expected outcomes.

## Not in this milestone

- Faults at the network or in the source of the subject: refused in RFC-0019.
- The same claims in TypeScript: milestone 0008 binds TypeScript, and each machinery after it adds its TypeScript part.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 9 | The first roadmap of dokimi |
| 2026-10-31 | In progress at position 8. The scope gains the claims of `circuit-breaker`, `sticky` and `leak-free`, and `transaction-under-faults` | RFC-0019 is accepted. It binds every classification whose stance names its machinery, and the two mixins whose laws need a dependency that fails |
| 2026-10-31 | Seven of the eight boxes are ticked | `dokimi conform --lang go` passes with 173 conformant fixtures against the conform and the spec of dokimi-go before their release. The last box waits for that release and the corpus that requires it |
