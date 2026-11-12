---
milestone: 0020
title: The TypeScript runtime implements spec 2.1.0
status: In progress
depends-on: 0008
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0050
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0020: The TypeScript runtime implements spec 2.1.0

## Goal

A TypeScript team applies a suite of any machinery of spec 2.1.0 through `@dokimi/runtime`: calls that fail on demand, generated sequences, concurrent callers, a breakable medium, a simulated host, and a simulation of machines and components.

## Done when

- [x] RFC-0050 is accepted.
- [ ] `@dokimi/runtime` exports each module of the tables of RFC-0050.
- [ ] `spec try` of dokimi-spec runs the tests of the TypeScript runtime, and they pass.
- [ ] The tests of each module of the runtime cover 100% of its statements under `make check`.
- [x] The doubles run the 13 checks of `behaviour` and the two checks of `resetBehaviour`, and the TypeScript tree of the stub generator does not write `newDelayed` or `newFailing`.
- [ ] A world of a subject replays a failing case whole in native mode and in Tasks mode.
- [ ] A simulation of two machines and one of two components fail their planted defects with the reason codes of Go.

## Why now

The suites of level 1 run on the base of `conform` that milestone 0008 builds. Each later level of TypeScript calls a machinery of this milestone. A mixin of level 2 needs the medium and the injector, a contract of level 3 the sequences and the concurrent callers, and a door of level 4 the host. The runtime calls the task scheduler, the replays and `keep` of definition 16.2.0, so `@dokimi/assert` has to implement that definition first.

## Scope

The parts of `conform` beyond the machinery `nothing-beyond-a-caller`, the probes of the other machineries, the behaviour suites of the doubles and the asynchronous decision of a fault, the modules `medium`, `host` and `sim`, the TypeScript tree of the stub generator, and the listing of the runtime in dokimi-spec.

## Not in this milestone

- The TypeScript trees of the suite, fake and sim generators beyond level 1: milestones 0021 to 0024.
- Allocation ceilings, parallel benchmarks, the race detector and the watchdog of a stall: out in RFC-0050.

## Risks to the sequence

| Risk | What it delays | What we would do |
|---|---|---|
| `@dokimi/assert` lacks definition 16.2.0 | This milestone and milestones 0021 to 0024 | Request the catch-up from assert-typescript, and build the modules that definition 13.6.0 covers |
| The definition of assert-spec lacks the resumed searches of a growing history | The check of the history after each step of a simulation, and the linearizability probe | Request them from assert-spec, as assert-go declares them beside the definition |

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-11-11 | Created at position 20 | RFC-0050 designs the TypeScript runtime at spec 2.1.0 |
| 2026-11-11 | `conform/output`, `probe/sequence` and `readVersionsNeverFall` move to milestone 0008 | The fixture `partnernaming` of level 1 checks `reads-never-go-back` over generated sequences |
| 2026-11-11 | Status set to In progress | The build of the medium, the host and the doubles has started |
