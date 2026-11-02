---
milestone: 0019
title: Transactions are checked for their isolation level
status: Planned
depends-on: 0011
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0033
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0019: Transactions are checked for their isolation level

## Goal

A Go team whose declaration has the contract `tx` and the mixin `serializable` or `snapshot-isolation` gets a check of that isolation level over the histories of its concurrent transactions.

## Done when

- [ ] RFC-0033 is accepted.
- [ ] The actions of a declaration with the contract `tx` run whole transactions and record each as one call of `txn`.
- [ ] The suite generator writes the checks `serializable` and `snapshot-isolation` for a declaration with the mixin of the same name, and refuses each with the gate `transactions` where the scope of a transaction does not fit.
- [ ] Each planted defect fails its check with its reason code at the seed of falsification under the strategy `PCT`, and every check of a single transaction passes on it.
- [ ] The serializable and snapshot fixtures of the corpus run their checks on three clients.

## Why now

The checks run their transactions on the concurrent clients and the schedule of milestone 0011, and they take their licence from a mixin of one callable as `linearizable` does. Nothing waits on this milestone.

## Scope

The two claims, their probes, reasons and defect forms, the gate `transactions`, the action of a transaction, and the transactional fixtures of the corpus.

## Not in this milestone

- A begin that returns a handle whose own methods write and read, until the first open question of RFC-0033 has an answer.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-11-02 | Created at position 15 | The claims `serializable` and `snapshot-isolation` move from milestone 0011, because their checks need the actions of a transaction, which RFC-0033 designs |
