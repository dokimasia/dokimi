---
milestone: 0006
title: A misspelled brand fails the run
status: Planned
depends-on: none
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: none
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0006: A misspelled brand fails the run

## Goal

A developer who writes `//dokimy:suite` gets a finding that names the brand `dokimi`, instead of a run that generates nothing for the declaration and reports success.

## Done when

- [x] eidos reports a carrier whose brand is within two edits of a registered brand, with the registered brand in the message: the warning EID-0083.
- [x] dokimi requires the release of eidos with that finding: the commit 482b0bb7 of eidos, which reports the warning and fails `dokimi run --strict`.
- [ ] A fixture of dokimi with `//dokimy:suite` fails `dokimi run` with the finding, and the fixture passes once the brand is fixed.

## Why now

Nothing has to exist first, and nothing waits on it. The order is arbitrary. It is at this position because it is the one silent hole of a declaration that no other milestone closes, and its work is in eidos.

## Scope

The request to eidos and the fixture of dokimi. An unknown name under the correct brand already gets a finding with the nearest name.

## Not in this milestone

- A finding for an expected value that flows from the subject: nobody has scheduled it yet.

## Risks to the sequence

| Risk | What it delays | What we would do |
|---|---|---|
| A near brand matches a comment that is no directive, such as a URL | Nothing else | Limit the finding to a line that opens with `//` and a word followed by `:` and a known name of a directive |

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 6 | The first roadmap of dokimi |

## References

| What | Where |
|---|---|
| eidos, the carrier of a directive | <https://github.com/dokimasia/eidos/blob/main/eidos-core/plugin/unit.go> |
