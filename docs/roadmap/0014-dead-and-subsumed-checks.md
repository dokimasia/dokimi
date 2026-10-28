---
milestone: 0014
title: A report names dead and subsumed checks
status: Planned
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

- [ ] RFC-0024 is accepted.
- [ ] `dokimi report run.json --mutation mutation.json` reports `DeadCheck`, `SubsumedCheck` and `UnattributedMutant` for a fixture of the corpus with a planted dead check.
- [ ] mutate-spec has released the mode that records every killing test, and dokimi-mutate-go implements it.
- [ ] A mutation record without that mode gives dead checks and no subsumed checks.

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
