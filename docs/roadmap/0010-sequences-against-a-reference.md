---
milestone: 0010
title: Laws of a history are judged over generated sequences
status: Planned
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

A Go team whose declaration has a mixin of a history, such as `monotonic-reads` or `accumulates`, gets a check of that law over generated sequences of calls. A team that states `machineries=reference` and passes a reference implementation also gets a check that every sequence returns on its subject what it returns on the reference. A failure shrinks to the smallest failing sequence.

## Done when

- [ ] RFC-0020 is accepted.
- [ ] The suite generator writes the actions of each declaration and a check for each judge whose mixin the declaration has: `ordered`, `invariant`, `folds`, `conserved` and `after-fault`.
- [ ] The suite generator writes `agrees-with-the-reference` for a declaration with `machineries=reference`, and lists it as unclaimed with the stance `needs` without the param.
- [ ] Each defect form of a judge fails its check at the seed of falsification, which comes from the key of the check, with its reason code, and every check of a single caller passes on it.
- [ ] A failing run on a subject prints the smallest sequence, and a rerun with its seed repeats it.
- [ ] A fixture of the corpus for each mixin of a judge runs its check, and a fixture for a store, a log and a queue runs `agrees-with-the-reference` with its reference.

## Why now

The checks are derived checks with reason codes, so milestone 0003 has to exist first. The judges of a single client come before concurrent callers because the concurrent callers of milestone 0011 run the same actions on several clients. The simulation of milestone 0015 runs its world as actions of the same kind, and records the calls of its clients in the same history.

## Scope

The generated actions, which record their calls in the history of the case, the judges of the package `probe`, the claims and defect forms of the judges, the param `normalize`, and the fixtures of the corpus.

## Not in this milestone

- Concurrent callers: go to milestone 0011.
- A model written in the spec: refused in RFC-0020.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 10 | The first roadmap of dokimi |
