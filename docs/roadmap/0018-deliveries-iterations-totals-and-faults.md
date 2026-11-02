---
milestone: 0018
title: Deliveries, iterations, totals and failed calls are judged over generated sequences
status: Done
depends-on: 0010
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0031
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0018: Deliveries, iterations, totals and failed calls are judged over generated sequences

## Goal

A Go team whose declaration has an outbox, a publisher, a watcher, a stream that reflects its mutations, a conservative move or an atomic write gets a check of that law over generated sequences of calls. A message that a sequence appended or published arrives as the contract states, an iteration sees the writes and the deletes made while it ran, and a failed call leaves no trace.

## Done when

- [x] RFC-0031 is accepted.
- [x] The mixin `conservative` of the shape catalog has the param `observe`.
- [x] The package `probe` has the seven probes of RFC-0031, and the probes of deliveries run inside a bubble of `testing/synctest`.
- [x] The suite generator writes a check for each claim of RFC-0031 whose mixin or contract the declaration has.
- [x] `no-trace-of-a-failed-call` entails `no-partial-write`, and the pair passes the proof of its depth.
- [x] Each defect form fails its check at the seed of falsification with its reason code, and every check of a single caller passes on it.
- [x] A fixture of the corpus for each mixin and each contract of RFC-0031 runs its check.

## Why now

The probes run in the generated sequences of a declaration and record their calls in the history of the case, so milestone 0010 has to exist first. They come before concurrent callers because a delivery and an iteration of one client are the base cases of their concurrent forms.

## Scope

The seven claims, laws, reasons, probes and defect forms of RFC-0031, the input kinds `pull`, `channel` and `handle-method`, and the fixtures of the corpus.

## Not in this milestone

- Deliveries to several subscribers at once: go to milestone 0011.
- A delivery through a network or a process outside the bubble.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-11-01 | Created at position 12 | The review of RFC-0020 found that the judges `conserved`, `after-fault`, `delivered` and `reflects` have no mechanism, and RFC-0031 designs them |
| 2026-11-01 | The mixin `conservative` has the param `observe` | eidos added it at be6064ca |
| 2026-11-01 | In progress | RFC-0031 is accepted. Milestones 0010 and 0017 are done |
| 2026-11-02 | `redelivers-on-demand` applies to a publisher at least once alone | Under the other modes a second delivery breaks the mode, so the check of the mode judges each redelivery |
| 2026-11-02 | Done | The corpus is conformant at level 4 with the checks of the seven claims. The read-through cache runs the check of a failed call, because it is the fixture with the mixin `atomic` and the machinery `faults` |
