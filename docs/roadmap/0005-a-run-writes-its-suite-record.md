---
milestone: 0005
title: A run writes its suite record
status: Done
depends-on: 0003, 0004
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0016
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0005: A run writes its suite record

## Goal

A developer runs `go test -json ./... > run.json` and `dokimi report run.json`, and reads what each suite claimed, proved, refused and found vacuous, without the service.

## Done when

- [x] RFC-0016 is accepted.
- [x] Each application of a suite in the Go corpus writes a suite record under `dokimi.suite.<seq>` in the artifact of `go test -json`, split as assert-spec splits a call record, and the record validates against the schema of the spec.
- [x] `dokimi report` lists each check that was vacuous in every application of the run.
- [x] The record contains the seed of the run, and a rerun with `DOKIMI_SEED` set to it repeats the run.
- [x] `dokimi report run.json` prints the claims, proofs, refusals, unclaimed classifications and vacuous checks of every suite of the corpus.
- [x] `dokimi report` exits with a failure for a record whose falsification result is not `proven`.
- [x] assert-spec has accepted the reservation of the prefix `dokimi.` in its overlays, and lists `dokimi.suite` among their kinds.

## Why now

The record contains the reason codes of milestone 0003 and the vacuous result of milestone 0004, so both have to exist first. The corpus gate of milestone 0007 reads the records, and so does the service.

## Scope

The fields `Package`, `Refusals` and `Unclaimed` of `conform.Suite`, the writing of the record in Go, and the command `dokimi report`.

## Not in this milestone

- The record in TypeScript: goes to milestone 0008.
- The upload of a record to the service: the service builds it.
- A converter to CTRF or JUnit XML: nobody has scheduled it yet.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 5 | The first roadmap of dokimi |
| 2026-10-29 | In progress | RFC-0016 is accepted |
| 2026-10-30 | Done | The Go runtime writes the suite records of each application, and `dokimi report` reads the records of a run of the Go corpus |
