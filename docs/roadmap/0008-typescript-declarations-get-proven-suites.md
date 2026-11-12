---
milestone: 0008
title: A TypeScript declaration gets a proven suite
status: In progress
depends-on: 0007
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0018, RFC-0050
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
- [x] RFC-0050 is accepted.
- [ ] dokimi-spec lists the TypeScript runtime as an implementation, and the runtime vendors spec 2.1.0 from it.
- [ ] `dokimi run` writes `iface.suite.ts` and `iface.suite.test.ts` for each declaration under `//dokimi:suite` in the TypeScript corpus.
- [ ] `@dokimi/runtime` exports `conform`, `spec` and `probe`, with the runtime model of spec 2.1.0 and a probe for each claim of level 1.
- [ ] Each fixture of level 1 of the TypeScript corpus that has a subject in Go has a subject at the same path, and the corpus has ports of `dynamic`, `nested`, `packagescope`, `recursive`, `refused`, `sampledparam` and `suitedhandle`.
- [ ] Each application writes its suite record into the meta of its Vitest test, and `dokimi report` reads it.
- [ ] A run of dokimi writes `dokimi.typescript.lock` for the plan typescript, and removing a directive from a fixture of the TypeScript corpus fails `dokimi run --check` with `LostCheck`.
- [ ] `dokimi conform --lang typescript` passes at level 1, with an overlay that declares each difference from Go.
- [ ] Removing an entry from the overlay of TypeScript fails the gate with the key of the difference.

## Why now

The corpus gate of milestone 0007 has to exist first, so that TypeScript is judged against the same expectations as Go from its first check. It comes after every milestone of Go, because Go is the primary language of dokimi. Each machinery of Go that comes before it adds probes that the binding writes in TypeScript as well.

## Scope

The template tree of the suite generator for TypeScript with the template of its lock, the modules `conform`, `spec` and `probe` of `@dokimi/runtime` with its copy of spec 2.1.0 and its runtime model, the probes of level 1, the reader of the Vitest report in `dokimi report`, the overlay of TypeScript, and the subjects and the ports of level 1 in its corpus. The fixture `partnernaming` of level 1 has a check of generated sequences, so the scope also has the actions of TypeScript, `conform/output`, `probe/sequence` and the probe `readVersionsNeverFall`.

## Not in this milestone

- The other machineries of the TypeScript runtime: milestone 0020.
- Levels 2 to 4 in TypeScript: milestones 0021, 0022 and 0024.
- Allocation budgets in TypeScript: refused in RFC-0018.
- Benchmarks, latency budgets and fuzz targets in TypeScript: milestone 0021, which has the first fixture with a latency budget.

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
| 2026-11-11 | The copy of the spec comes from dokimi-spec at 2.1.0, the scope gains the runtime model, and levels 2 to 4 go to milestones 0021, 0022 and 0024 | RFC-0050 designs the TypeScript runtime at spec 2.1.0, and TypeScript comes before the bindings of the other languages |
| 2026-11-11 | In progress. The corpus runs against the checkouts of `@dokimi/runtime` and `@dokimi/assert` | RFC-0050 is accepted, and neither package has a release before the TypeScript milestones are built |
| 2026-11-11 | The ports of level 1 gain `packagescope` and `refused` | RFC-0050 ports both fixtures at level 1, beside the five ports of RFC-0018 |
| 2026-11-11 | The scope gains the actions of TypeScript, `conform/output`, `probe/sequence` and `readVersionsNeverFall` from milestone 0020 | The fixture `partnernaming` of level 1 checks `reads-never-go-back` over generated sequences |
