---
milestone: 0007
title: The corpus judges each language by its records
status: Planned
depends-on: 0001, 0004, 0005
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0017
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0007: The corpus judges each language by its records

## Goal

A maintainer runs `dokimi conform --lang go` and learns, for each fixture, whether Go generates, proves and applies the checks that the fixture expects, and every difference that Go has not declared fails the gate.

## Done when

- [ ] RFC-0017 is accepted.
- [ ] Each fixture of conformance levels 1 to 4 has a file of expected outcomes under `corpus/expect`.
- [ ] `dokimi conform --lang go` passes on the corpus with an overlay of Go.
- [ ] A planted difference in the Go binding, such as a probe that reports another reason code, fails the gate with the fixture and the two codes.
- [ ] An overlay entry that matches no difference fails the gate, and so does an entry whose outcome now meets the expectation.
- [ ] The gate of dokimi runs `dokimi conform` for every language that has a binding.

## Why now

The gate compares the lock of milestone 0001, the controls of milestone 0004 and the suite records of milestone 0005 with the expectations, so all three have to exist first. The TypeScript binding of milestone 0008 is judged by it from its first check.

## Scope

The expected outcomes of the corpus, the overlay of Go, the levels, and the command `dokimi conform`.

## Not in this milestone

- The overlay of TypeScript: goes to milestone 0008.
- A gate that runs only the fixtures that a change touches: nobody has scheduled it yet.

## Risks to the sequence

| Risk | What it delays | What we would do |
|---|---|---|
| A full run of the gate is too slow for each pull request | Milestone 0008, whose changes run the gate often | Run the fixtures of the changed level on each pull request and the whole corpus on the main branch |

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 7 | The first roadmap of dokimi |
