---
milestone: 0013
title: A durable subject survives a breakable medium
status: Planned
depends-on: 0003
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0023
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0013: A durable subject survives a breakable medium

## Goal

A Go team whose subject writes to storage through `medium.FS` gets checks that every synced write survives a crash at every persistence point of short sequences, and that recovery returns no torn record.

## Done when

- [ ] RFC-0023 is accepted.
- [ ] The module `go.dokimi.dev/medium` is released with `FS`, `OS` and `Memory`.
- [x] eidos has the mixin `durable` with its param `corrupt`, and dokimi requires the commit of eidos that adds it.
- [ ] A durable log of the corpus under `//dokimi:suite machineries=medium` passes `synced-writes-survive` and `no-torn-record` through the option `Durable`, and each claim fails on its planted defect with its reason code.
- [ ] An application with `Durable` runs the checks of the doors `Restart` and `Crash` without `Restartable` or `Crashable`.
- [ ] The probe enumerates every crash point of every sequence of up to three calls, and the suite record states the count.
- [ ] The fixture `door/durable` of the corpus uses `medium.Memory` in place of its own medium.

## Why now

The claims are derived checks, so milestone 0003 has to exist first. It is the costliest machinery for a consumer, because the subject must take a `medium.FS`. Nothing waits on it.

## Scope

The module `medium`, the option `Durable` with the doors that it implements, the claims and defect forms of the mixin, and the durable fixtures of the corpus.

## Not in this milestone

- A medium for TypeScript: nobody has scheduled it yet.
- Faults of a FUSE filesystem or a block device: refused in RFC-0023.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 13 | The first roadmap of dokimi |
| 2026-10-28 | Moved to position 14 | Milestone 0015 took position 11, because the first milestone of sharder waits on it |
