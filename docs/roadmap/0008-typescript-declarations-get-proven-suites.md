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

- [x] RFC-0018 is accepted.
- [ ] `dokimi run` writes `iface.suite.ts` and `iface.suite.test.ts` for each declaration under `//dokimi:suite` in the TypeScript corpus.
- [ ] `@dokimi/runtime` exports `conform`, `spec` and `probe`, with a probe for each claim of level 1.
- [ ] Each fixture of level 1 of the TypeScript corpus that has a subject in Go has a subject at the same path, and the corpus has ports of `dynamic`, `nested`, `recursive`, `sampledparam` and `suitedhandle`.
- [ ] Each application writes its suite record into the meta of its Vitest test, and `dokimi report` reads it.
- [ ] A run of dokimi writes `dokimi.typescript.lock` for the plan typescript, and removing a directive from a fixture of the TypeScript corpus fails `dokimi run --check` with `LostCheck`.
- [ ] `dokimi conform --lang typescript` passes at level 1, with an overlay that declares each difference from Go.
- [ ] Removing an entry from the overlay of TypeScript fails the gate with the key of the difference.

## Why now

The corpus gate of milestone 0007 has to exist first, so that TypeScript is judged against the same expectations as Go from its first check. It comes after every milestone of Go, because Go is the primary language of dokimi. Each machinery of Go that comes before it adds probes that the binding writes in TypeScript as well.

## Scope

The template tree of the suite generator for TypeScript with the template of its lock, the modules `conform`, `spec` and `probe` of `@dokimi/runtime` with its copy of the spec, the probes of level 1, the reader of the Vitest report in `dokimi report`, the overlay of TypeScript, and the subjects and the ports of level 1 in its corpus.

## Not in this milestone

- Levels 2 to 4 in TypeScript: nobody has scheduled them yet.
- Allocation budgets in TypeScript: refused in RFC-0018.
- Benchmarks and latency budgets in TypeScript: designed in RFC-0006, and the first fixture with a latency budget is at level 2. The generator reports a latency budget of TypeScript as a warning until then.
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
| 2026-10-30 | The lock of the plan typescript moves here from milestone 0001 | A plan has a lock when it runs the suite generator, and the plan typescript runs it from this milestone on |
| 2026-10-31 | In progress | RFC-0018 is accepted. The falsification of a derived check needs a subject, so the scope gains the subjects of level 1 and five ports of fixtures of Go |
| 2026-10-31 | Planned, and moved from position 8 to the end of the order | Go is the primary language of dokimi, so the milestones of Go come first |
