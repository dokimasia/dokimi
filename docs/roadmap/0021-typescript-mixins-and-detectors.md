---
milestone: 0021
title: A TypeScript declaration gets the checks of its mixins and detected shapes
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

# Milestone 0021: A TypeScript declaration gets the checks of its mixins and detected shapes

## Goal

A TypeScript team whose declaration has a mixin or a detected shape gets the derived checks of conformance level 2, each proven against its planted defect by the same reason code as in Go.

## Done when

- [ ] The suite generator writes the TypeScript checks of each claim of the fixtures of level 2, with generated inputs, the bench module and the fuzz module of each declaration.
- [ ] Each fixture of level 2 of the TypeScript corpus that has a subject in Go has a subject at the same path, and the corpus has the ports of level 2 of RFC-0050.
- [ ] `dokimi conform --lang typescript` passes at level 2, with an overlay that declares each difference from Go.

## Why now

The checks of level 2 call the probes of their mixins, the medium of the mixin `durable` and the injector of milestone 0020, through the template tree that milestone 0008 starts.

## Scope

The bindings of the claims of the mixins and of the detected shapes in the TypeScript tree of the suite generator, the modules `<module>.suite.bench.ts` and `<module>.suite.fuzz.ts`, the subjects, the controls and the ports of level 2 in the TypeScript corpus, and the entries of level 2 of the overlay.

## Not in this milestone

- Levels 3 and 4 in TypeScript: milestones 0022 and 0024.
- The attribution of mutants of `iface/mutation/deadcheck` in TypeScript: it needs a TypeScript overlay of mutate-spec and an engine that writes its records.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-11-11 | Created at position 21 | RFC-0050 designs the TypeScript runtime at spec 2.1.0 |
