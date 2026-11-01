---
milestone: 0010
title: Laws of a history are judged over generated sequences
status: Done
depends-on: 0003
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0020
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0010: Laws of a history are judged over generated sequences

## Goal

A Go team whose declaration has a mixin of a history, such as `monotonic-reads` or `accumulates`, gets a check of that law over generated sequences of calls. A team that states `machineries=[reference]` and passes a reference implementation also gets a check that every sequence returns on its subject what it returns on the reference. A failure shrinks to the smallest failing sequence.

## Done when

- [x] RFC-0020 is accepted.
- [x] The suite generator writes the actions of each declaration, and a check for each claim of RFC-0020 whose mixin or contract the declaration has.
- [x] Each claim with a claim of a single caller of the same law entails it, and the pair passes the proof of its depth.
- [x] The mixin `bounded` of the shape catalog has the param `observe`, and `keeps-within-its-bounds` runs on a fixture with it.
- [x] The suite generator writes `agrees-with-the-reference` for a declaration with `machineries=[reference]`, and omits the check and its entry for a declaration without the param.
- [x] Each defect form of a judge fails its check at the seed of falsification, which comes from the identity of the check, with its reason code, and every check of a single caller passes on it.
- [x] A failing run on a subject prints the smallest sequence, and a rerun with its seed repeats it.
- [x] A fixture of the corpus for each mixin and each contract of a judge runs its check, and a fixture for a store, a log and a queue runs `agrees-with-the-reference` with its reference.

## Why now

The checks are derived checks with reason codes, so milestone 0003 has to exist first. The judges of a single client come before concurrent callers because the concurrent callers of milestone 0011 run the same actions on several clients. The simulation of milestone 0015 runs its world as actions of the same kind, and records the calls of its clients in the same history.

## Scope

The generated actions, which record their calls in the history of the case, the judges of the package `probe`, the claims and defect forms of the judges, the param `normalize`, and the fixtures of the corpus.

## Not in this milestone

- Concurrent callers: go to milestone 0011.
- The judges of deliveries, iterations, totals and failed calls: go to milestone 0018.
- A model written in the spec: refused in RFC-0020.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 10 | The first roadmap of dokimi |
| 2026-10-31 | The scope gains the claims of 21 classifications, six judges and the param `observe` of the mixin `bounded` | RFC-0020 binds every classification whose stance or law names generated sequences. The mixin `bounded` names no reader of its resource in the shape catalog |
| 2026-11-01 | In progress. The judges `conserved`, `after-fault`, `delivered` and `reflects` and their seven claims go to milestone 0018 | The review of RFC-0020 against assert-go and conform found contradictions and four judges without a mechanism. The revised RFC-0020 is accepted. Milestone 0017 waits for the push of the probes of dokimi-go |
| 2026-11-01 | Done | The 19 claims of RFC-0020 run on the fixtures of the Go corpus, and the fixture `iface/door/reference` runs a store, a journal and a queue against their references. The corpus judges 177 of 177 fixtures conformant |
