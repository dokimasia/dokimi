---
milestone: 0014
title: A report names dead and subsumed checks
status: Done
depends-on: 0005
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0024
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0014: A report names dead and subsumed checks

## Goal

A team that runs a mutation run beside its suite learns which checks killed no mutant of the code they cover and which checks another check already covers.

## Done when

- [x] RFC-0024 is accepted.
- [x] `dokimi report run.json --mutation mutation.json` reports `DeadCheck`, `SubsumedCheck` and `UnattributedMutant` for a fixture of the corpus with a planted dead check.
- [x] mutate-spec has released the modes that record every killing test and the coverage of each subtest, and dokimi-mutate-go implements them.
- [x] A mutation record without those modes gives unattributed mutants but not dead or subsumed checks, and `dokimi report` prints the mode that the record lacks.

## Why now

The attribution joins the mutation record with the suite records of milestone 0005, so that milestone has to exist first. Nothing waits on it.

## Scope

The attribution in `dokimi report`, the findings, and the request to mutate-spec.

## Not in this milestone

- Mutation as the proof of a check: refused in RFC-0001.
- The view of mutation in the service: the service builds it.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 14 | The first roadmap of dokimi |
| 2026-10-28 | Moved to position 15 | Milestone 0015 took position 11, because the first milestone of sharder waits on it |
| 2026-11-03 | In progress | RFC-0024 is accepted |
| 2026-11-03 | The milestone waits for a second mode of mutate-spec, and a record without the modes gives no dead checks | A mutation run records the coverage of each top-level test, and conform runs each check as a subtest. With the first killer alone, a check that runs after another killer of each of its mutants looks dead. [dokimasia/mutate-spec#36](https://github.com/dokimasia/mutate-spec/issues/36) and [dokimasia/mutate-spec#38](https://github.com/dokimasia/mutate-spec/issues/38) request the modes |
| 2026-11-03 | Done | `dokimi report --mutation` reads a record of dokimi-mutate-go v0.3.0 over the corpus fixture `iface/mutation/deadcheck`, and reports its dead check, its four subsumed checks and the two survivors of its method `Len` |
