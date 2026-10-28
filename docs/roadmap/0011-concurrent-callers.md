---
milestone: 0011
title: Concurrent callers are checked for linearizability
status: Planned
depends-on: 0010
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0021
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0011: Concurrent callers are checked for linearizability

## Goal

A Go team whose declaration has the mixin `concurrent` and a reference gets a check that every history of concurrent calls has an order in which its reference returns the same results. A team whose declaration has the mixin `serializable` or `snapshot-isolation` gets a check of that isolation level over the histories of its transactions.

## Done when

- [ ] RFC-0021 is accepted.
- [ ] The suite generator writes the check `linearizable` for a declaration with the mixin `concurrent` and `machineries=reference`.
- [ ] The suite generator writes the checks `serializable` and `snapshot-isolation` for a declaration with the mixin of the same name.
- [ ] Each planted defect fails its check with its reason code at the seed of falsification under the strategy `PCT`, and the checks of milestone 0010 pass on it.
- [ ] A check whose every case the budget cannot decide is vacuous in the census.
- [ ] The concurrent fixtures of the corpus run `linearizable` with three clients and a partition by key.

## Why now

The checks run the actions of milestone 0010 on several clients, so that milestone has to exist first. Nothing waits on it.

## Scope

The probes `linearizable-history`, `serializable-history` and `snapshot-history` over the checkers of assert-go, the claims and their defect forms, the partition by key, and the concurrent and transactional fixtures of the corpus.

## Not in this milestone

- Controlled interleavings: out, in RFC-0001.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 11 | The first roadmap of dokimi |
| 2026-10-28 | Moved to position 12 | Milestone 0015 took position 11, because the first milestone of sharder waits on it |
