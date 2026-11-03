---
milestone: 0013
title: A durable subject survives a breakable medium
status: Done
depends-on: 0003, 0015
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

- [x] RFC-0023 is accepted.
- [x] The module `go.dokimi.dev/medium` is released with `FS`, `OS` and `Memory`.
- [x] eidos has the mixin `durable` with its param `corrupt`, and dokimi requires the commit of eidos that adds it.
- [x] A durable log of the corpus under `//dokimi:suite machineries=[medium]` passes `synced-writes-survive` and `no-torn-record` through the option `Durable`, and each claim fails on its planted defect with its reason code.
- [x] An application with `Durable` runs the checks of the doors `Restart` and `Crash` without `Restartable` or `Crashable`.
- [x] The probe enumerates every crash point of every sequence of up to three calls, and the suite record states the count.
- [x] The fixture `door/durable` of the corpus uses `medium.Memory` in place of its own medium.
- [x] The world of the simulation strikes the faults of the breakable medium at a dependency that persists through `medium.FS`.

## Why now

The claims are derived checks, so milestone 0003 has to exist first. It is the costliest machinery for a consumer, because the subject must take a `medium.FS`. The world of the simulation has to exist before it takes the faults of the medium, so milestone 0015 comes first. Nothing waits on it.

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
| 2026-11-02 | Scope: the faults of the medium in the simulation | RFC-0025 strikes the faults of the breakable medium in the simulation, and this milestone builds the medium after milestone 0015 builds the world |
| 2026-11-03 | In progress | RFC-0023 is accepted |
| 2026-11-03 | Done | dokimi-go released `medium` 0.1.0, and the Go corpus requires it. `dokimi conform --lang go` judges 180 of 180 fixtures conformant, with the four claims of the mixin `durable` on the fixture `iface/mixin/durable` |
