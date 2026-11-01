---
milestone: 0017
title: Every classification of the shape catalog has a claim
status: Planned
depends-on: 0003
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0028
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0017: Every classification of the shape catalog has a claim

## Goal

A Go team whose declaration has a batch reader, a stream reader, an indexed callable, a wrapper of a sibling, or a role `verify`, `release` or `stats` gets a check of the law of that classification. No classification of the shape catalog keeps the stance `gap`.

## Done when

- [ ] RFC-0028 is accepted.
- [ ] The spec has the claims of RFC-0028 with their laws, probes, reasons and defect forms, and its unclaimed table has no entry with the stance `gap`.
- [ ] The package `probe` has the seven probes of RFC-0028, and the suite generator binds the partner of `agrees-with-its-sibling` to a new subject.
- [ ] A fixture of the corpus for each of the eight classifications runs its check, and each check fails on its planted defect with its reason code.
- [ ] `dokimi conform --lang go` passes with the expected outcomes of the new fixtures.

## Why now

The claims are derived checks of a single caller, so milestone 0003 has to exist first. They do not need a new machinery, so they come before the machineries of sequences and concurrency. The generated bindings of milestone 0016 run each of them over generated inputs as well.

## Scope

The claims, laws, reasons, defect forms and probes of RFC-0028, the binding of a partner to a new subject, and the fixtures of the corpus with their expected outcomes.

## Not in this milestone

- The induced failure of `poison-accessor` and a change of a chain: the claims of the mixins `poisonable` and `tamper-evident` check them.
- The claims of `sticky` and `leak-free`: milestone 0009.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-31 | Created at position 10 | Ten classifications of the catalog had the stance `gap`, and eight of them state a law that one caller checks |
