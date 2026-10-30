---
milestone: 0004
title: Vacuous and over-strong checks fail
status: Done
depends-on: 0003
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0015
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0004: Vacuous and over-strong checks fail

## Goal

A team learns, from a failing test, that a check never exercised its claim on any subject, or that a correct implementation other than its subject fails a check.

## Done when

- [x] RFC-0015 is accepted.
- [x] An application whose subject never fails a call fails `Store.Read/zero-on-error` with a message that names the vacuous condition and both remedies.
- [x] An excusal of the same check makes the application pass, and the excusal is counted apart from the passes.
- [x] A check that skips because the environment lacks direct I/O has the result `skipped`, and the application passes.
- [x] The application of a suite to a control of the corpus fails a check that the control fails, with the identity of the check.
- [x] Each shape family of the Go corpus has a control that is not its subject, and a test of the corpus fails for a family without one.
- [x] assert-go has released the `Observer` hook, in v0.11.0.
- [x] An authored check that never calls its decisive assertion is vacuous.

## Why now

Engagement is called by the probes of milestone 0003, so the probes have to exist first. The suite record of milestone 0005 writes the vacuous result, and the corpus gate of milestone 0007 judges the controls.

## Scope

`Handle.Engage`, the vacuous and skipped outcomes and the census of an application, the handle as an `Observer` of assert-go, and a control for each shape family of the Go corpus.

## Not in this milestone

- Controls of the TypeScript corpus: go to milestone 0008.
- Generated controls: refused in RFC-0015.

## Risks to the sequence

| Risk | What it delays | What we would do |
|---|---|---|
| Writing a control for each family takes longer than the code | Milestone 0007 | Start the corpus gate with the families that have controls, and fail the others in its phase of controls |

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 4 | The first roadmap of dokimi |
| 2026-10-29 | Done | The Go runtime fails vacuous checks, and the Go corpus has a control for each of its 11 shape families |
