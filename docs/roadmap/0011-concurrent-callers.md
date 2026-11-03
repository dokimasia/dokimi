---
milestone: 0011
title: Concurrent callers are checked for linearizability
status: In progress
depends-on: 0009, 0010
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

A Go team whose declaration has the mixin `concurrent` and a reference gets a check that every history of concurrent calls has an order in which its reference returns the same results.

## Done when

- [x] RFC-0021 is accepted.
- [x] The suite generator writes the check `linearizable` for a declaration with the mixin `concurrent` and `machineries=[reference]`.
- [x] The suite generator writes `reads-agree-under-concurrency`, `computes-once`, `causal-history`, `pins-the-snapshot-under-writes` and `one-holder-at-a-time` for each callable that their mixin or contract classifies.
- [x] `computes-once` runs inside a `testing/synctest` bubble and keeps the computation waiting through the injector of milestone 0009.
- [x] `causal-history`, `pins-the-snapshot-under-writes` and `one-holder-at-a-time` entail the claims of a single caller of their laws, and each pair passes the proof of its depth.
- [ ] Each planted defect fails its check with its reason code at the seed of falsification under the strategy `Uniform`, and the checks of milestone 0010 pass on it.
- [x] A check whose every case the budget cannot decide is vacuous in the census.
- [x] The concurrent fixtures of the corpus run `linearizable` with three clients and a partition by key.

## Why now

The checks run the actions of milestone 0010 on several clients, so that milestone has to exist first. `computes-once` keeps a computation waiting through the injector of milestone 0009. Nothing waits on it.

## Scope

The probe `linearizable-history` over the checker of assert-go, the probes of the other five claims of RFC-0021, the claims and their defect forms, the partition by key, and the concurrent fixtures of the corpus.

## Not in this milestone

- Controlled interleavings: out, in RFC-0001.
- The isolation levels `serializable` and `snapshot-isolation`: milestone 0019.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 11 | The first roadmap of dokimi |
| 2026-10-28 | Moved to position 12 | Milestone 0015 took position 11, because the first milestone of sharder waits on it |
| 2026-10-31 | The scope gains the claims of `concurrent-readers`, `single-flight`, `causal`, `point-in-time` and `lease`, and the milestone depends on 0009 | RFC-0021 binds every classification whose stance or law names concurrent callers. `computes-once` keeps its computation waiting through the injector of RFC-0019 |
| 2026-11-02 | In progress | RFC-0021 is accepted |
| 2026-11-02 | The claims `serializable` and `snapshot-isolation` move to milestone 0019 | Their checks drive transactions, and RFC-0033 designs the actions of a transaction |
| 2026-11-02 | A planted defect runs under the strategy `Uniform` in place of `PCT` | `PCT(3)` proved the planted races at 15 and 8 of 40 seeds, and `Uniform` with every action at 40 of 40 |
| 2026-11-02 | The claims `reads-agree-under-concurrency`, `causal-history` and `pins-the-snapshot-under-writes` take the machinery `generated-sequences` | A racy read shows only when two overlapping reads return distinct values, and the key pools of the generated sequences give reads of several keys |
| 2026-11-03 | The last box waits for an option of the steps of assert-spec that sets the mean number of steps of a concurrent section | A concurrent section averages five steps. The planted defects `racy-pin` and `racy-causal-read` race in few cases, and their proofs held at 33 and 30 of 40 seeds of the specs. [dokimasia/assert-spec#31](https://github.com/dokimasia/assert-spec/issues/31) requests the option |
