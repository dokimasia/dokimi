---
milestone: 0001
title: A lost check fails the generation gate
status: Planned
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

- [ ] A run of dokimi writes `dokimi.lock` for the Go plan and for the TypeScript plan, with one line per check, per refusal and per unclaimed classification, and the partners of each check.
- [ ] Removing a directive from a corpus fixture and running `dokimi run --check` fails with `LostCheck` for each check that the directive licensed.
- [ ] Removing the samples that a check needs, so that the check becomes a refusal, fails `dokimi run --check` with `LostCheck`.
- [ ] Committing the regenerated lock makes the same `dokimi run --check` pass.
- [ ] A run without a change of the source writes a lock with the same bytes.
- [x] eidos has the data family and the line-set family that the lock needs, and dokimi requires the commit of eidos that adds them.

## Why now

Nothing has to exist first: the identities of checks exist today. Every later milestone adds or changes checks, and the lock makes each of those changes reviewable. The corpus gate relies on the lock for its phase of generation.

## Scope

The lock as a line-set data family of eidos, with its refusals, unclaimed classifications and partners, and the finding `LostCheck`.

## Not in this milestone

- The comment on a pull request that names a lost claim: the service builds it.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 1 | The first roadmap of dokimi |
