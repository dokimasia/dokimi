---
milestone: 0008
title: A TypeScript declaration gets a proven suite
status: Planned
depends-on: 0007
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0018
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0008: A TypeScript declaration gets a proven suite

## Goal

A TypeScript team that writes `//dokimi:suite` on an interface gets a Vitest suite of derived checks at conformance level 1, each proven against its planted defect by the same reason code as in Go.

## Done when

- [ ] RFC-0018 is accepted.
- [ ] `dokimi run` writes `iface.suite.ts` and `iface.suite.test.ts` for each declaration under `//dokimi:suite` in the TypeScript corpus.
- [ ] `@dokimi/runtime` exports `conform`, `spec` and `probe`, with a probe for each claim of level 1.
- [ ] Each application writes its suite record into the meta of its Vitest test, and `dokimi report` reads it.
- [ ] `dokimi conform --lang typescript` passes at level 1, with an overlay that declares each difference from Go.
- [ ] Removing an entry from the overlay of TypeScript fails the gate with the key of the difference.

## Why now

The corpus gate of milestone 0007 has to exist first, so that TypeScript is judged against the same expectations as Go from its first check. It comes before the machineries above a single caller because each machinery added before it is one more probe to write twice.

## Scope

The template tree of the suite generator for TypeScript, the modules `conform`, `spec` and `probe` of `@dokimi/runtime`, the probes of level 1, the overlay of TypeScript and the controls of its corpus.

## Not in this milestone

- Levels 2 to 4 in TypeScript: nobody has scheduled them yet.
- Allocation budgets in TypeScript: refused in RFC-0018.
- Fuzz targets in TypeScript through Jazzer.js: designed in RFC-0006, and nobody has scheduled them yet.

## Risks to the sequence

| Risk | What it delays | What we would do |
|---|---|---|
| More than a third of the claims of level 1 refuse in TypeScript | Every later language binding | Stop, and revise the gates of the derivation before a third language |
| A claim needs another meaning in TypeScript | Every later language binding | Stop, because the shared derivation has failed, and revise RFC-0014 |

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 8 | The first roadmap of dokimi |
