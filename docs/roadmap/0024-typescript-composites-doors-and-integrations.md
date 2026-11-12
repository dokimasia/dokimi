---
milestone: 0024
title: A TypeScript declaration gets the checks of its composites, doors and integrations
status: Planned
depends-on: 0022, 0023
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

# Milestone 0024: A TypeScript declaration gets the checks of its composites, doors and integrations

## Goal

A TypeScript team whose declaration combines two or more axes, or whose subject opens through a door, gets the derived checks of conformance level 4. The TypeScript corpus passes the gate at every level that the Go corpus passes.

## Done when

- [ ] The suite generator writes the TypeScript checks of the composites, the doors and the integrations, with the world checks of a subject in native mode and in Tasks mode.
- [ ] The TypeScript corpus has the ports of the doors of level 4 of RFC-0050, and each fixture of level 4 that has a subject in Go has a subject at the same path.
- [ ] `dokimi conform --lang typescript` passes at level 4, with an overlay that declares each difference from Go.
- [ ] npm has releases of `@dokimi/runtime` and `@dokimi/assert`, the TypeScript corpus requires them, and the gate of dokimi runs it beside the gate of Go.

## Why now

Milestone 0022 comes first, because the gate judges level 4 with the levels below it. The fixtures of `sim/` are of level 4, and milestone 0023 gives them their simulation.

## Scope

The bindings of the claims of the composites, of the doors and of the integrations in the TypeScript tree of the suite generator, the ports and the subjects of level 4 in the TypeScript corpus, the entries of level 4 of the overlay, and the run of the gate of TypeScript in the gate of dokimi.

## Not in this milestone

- Allocation ceilings, parallel benchmarks, the race detector and the watchdog of a stall: out in RFC-0050.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-11-11 | Created at position 24 | RFC-0050 designs the TypeScript runtime at spec 2.1.0 |
