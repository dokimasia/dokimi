---
milestone: 0023
title: A TypeScript system runs in a deterministic simulation
status: Planned
depends-on: 0020
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

# Milestone 0023: A TypeScript system runs in a deterministic simulation

## Goal

A TypeScript team gets a deterministic simulation of its sans-IO machines, or of its components on simulated hosts, from the same directives as in Go.

## Done when

- [ ] The sim generator writes `<module>.sim.ts`, `<module>.sim.test.ts` and `<module>.sim.fuzz.ts` for each system of the TypeScript corpus.
- [ ] The TypeScript corpus has the five fixtures of `sim/`, and their tests pass under Vitest.
- [ ] Each rule, each invariant and each convergence of the five fixtures fails on its planted defect in `falsify`.
- [ ] A failing case of `sim/mirrored` replays whole in native mode and in Tasks mode.

## Why now

The generated simulation runs on the module `sim` of milestone 0020, and its components run on the host and the medium of the same milestone.

## Scope

The TypeScript tree of the sim generator with the TypeScript forms of its directives, and the five fixtures of `sim/` in the TypeScript corpus.

## Not in this milestone

- The judgement of the fixtures of `sim/` by the gate, which judges them at level 4: milestone 0024.
- The workloads of a simulation on a real cluster: their design is a Draft.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-11-11 | Created at position 23 | RFC-0050 designs the TypeScript runtime at spec 2.1.0 |
