---
milestone: 0001
title: A lost check fails the generation gate
status: Done
depends-on: none
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0013
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0001: A lost check fails the generation gate

## Goal

A reviewer sees every check that a change removes as one line in one file, and CI fails a change that loses a check without that line.

## Done when

- [x] RFC-0013 is accepted.
- [x] A run of dokimi writes `dokimi.go.lock` for the plan go, with one line per check, per ceiling, per refusal and per unclaimed classification, and the partners of each check.
- [x] The suite record of each check has the same key as its line of the lock.
- [x] Removing a directive from a corpus fixture and running `dokimi run --check` fails with `LostCheck` for each check that the directive licensed.
- [x] Removing the samples that a check needs, so that the check becomes a refusal, fails `dokimi run --check` with `LostCheck`.
- [x] Committing the regenerated lock makes the same `dokimi run --check` pass.
- [x] A run without a change of the source writes a lock with the same bytes.
- [x] eidos has the data family and the line-set family that the lock needs, and dokimi requires the commit of eidos that adds them.

## Why now

Nothing has to exist first: the identities of checks exist today. Every later milestone adds or changes checks, and the lock makes each of those changes reviewable. The corpus gate relies on the lock for its phase of generation.

## Scope

The lock as a line-set data family of eidos, with its ceilings, refusals, unclaimed classifications and partners, and the finding `LostCheck`. The key form in the runtime, so the suite records and the lock name each check alike.

## Not in this milestone

- The comment on a pull request that names a lost claim: the service builds it.
- The lock of the plan typescript: goes to milestone 0008, whose TypeScript suites run the suite generator in that plan.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 1 | The first roadmap of dokimi |
| 2026-10-30 | In progress | RFC-0013 is accepted |
| 2026-10-30 | The lock is `dokimi.go.lock`, the lock of the plan typescript goes to milestone 0008, and the records take the key form of the lock | One root contains the locks of every plan, the plan typescript runs no suite generator, and the corpus gate joins the lock and the records by key |
| 2026-10-30 | Done | The corpus lock has 1,018 lines, and each of its 766 checks, authored checks and ceilings has a suite record with the same key and claim. Without its directive `suite`, the fixture `iface/mixin/xsssafe` fails `dokimi run --check` with `LostCheck` for each of its four checks |
