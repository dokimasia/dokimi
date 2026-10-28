---
milestone: 0012
title: A downstream team uses a verified fake
status: Planned
depends-on: 0003, 0004
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0022
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0012: A downstream team uses a verified fake

## Goal

A downstream Go team replaces a container or a hand-written fake of a persister, a cache or an appender with a generated fake that passes the same suite as the production implementation.

## Done when

- [ ] RFC-0022 is accepted.
- [ ] The spec has reference forms for the contracts `persister`, `cache` and `appender`.
- [ ] The generator `fake` writes the fake and its test for each fixture of those contracts, and the test passes the suite of the declaration.
- [ ] A declaration with a callable outside every reference form gets no fake and a refusal that names the callable.
- [ ] The corpus applies the suite of each fixture of those contracts to the generated fake beside its hand-written `InMemory` subject and its controls, and the corpus passes.

## Why now

The derived checks of its declaration verify a fake, so milestone 0003 has to exist first. The controls of milestone 0004 check each fake beside the derived checks. A fake needs none of the machineries, so its place after them is arbitrary. It comes after the corpus gate and TypeScript because those serve every consumer, and a fake serves the consumers of three contracts.

## Scope

The reference forms in the spec, the generator `fake` with its Go template tree, and the fixtures of the corpus.

## Not in this milestone

- Fakes in TypeScript: nobody has scheduled them yet.
- Fakes of contracts beyond the three: nobody has scheduled them yet.
- A fake as the subject or a control of a corpus fixture: refused in RFC-0022.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 12 | The first roadmap of dokimi |
| 2026-10-28 | Moved to position 13 | Milestone 0015 took position 11, because the first milestone of sharder waits on it |
