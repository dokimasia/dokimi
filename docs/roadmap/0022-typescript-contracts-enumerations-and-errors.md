---
milestone: 0022
title: A TypeScript declaration gets the checks of its contracts, enumerations and errors
status: Planned
depends-on: 0021
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0050
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0022: A TypeScript declaration gets the checks of its contracts, enumerations and errors

## Goal

A TypeScript team whose declaration has a contract, an enumeration or errors gets the derived checks of conformance level 3, and a declaration of a persister, a cache or an appender gets a verified fake.

## Done when

- [ ] The suite generator writes the TypeScript checks of the contracts, the enumerations and the errors, with generated sequences, concurrent callers and the isolation levels of transactions.
- [ ] The fake generator writes `<module>.fake.ts` and `<module>.fake.test.ts` for each declaration under `//dokimi:fake` in the TypeScript corpus, and the suite of the declaration passes on the fake.
- [ ] Each fixture of level 3 of the TypeScript corpus that has a subject in Go has a subject at the same path.
- [ ] `dokimi conform --lang typescript` passes at level 3, with an overlay that declares each difference from Go and blocks `enum/split` and `errors/split`.

## Why now

The gate judges each level with the levels below it, so level 2 of milestone 0021 comes first. The checks of the contracts run the sequences, the concurrent callers and the checkers of histories of milestone 0020.

## Scope

The bindings of the claims of the contracts, of the enumerations and of the errors in the TypeScript tree of the suite generator, the TypeScript tree of the fake generator, the subjects and the controls of level 3 in the TypeScript corpus, and the entries of level 3 of the overlay.

## Not in this milestone

- Level 4 in TypeScript: milestone 0024.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-11-11 | Created at position 22 | RFC-0050 designs the TypeScript runtime at spec 2.1.0 |
