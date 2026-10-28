---
milestone: 0002
title: Every classification has a law or a stance
status: Planned
depends-on: none
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0002
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0002: Every classification has a law or a stance

## Goal

A reader finds in the directory `spec` of dokimi, for each of the 108 classifications of the shape catalog, the claims that it licenses with their laws, probes and reason codes, or the stance that states why it has none.

## Done when

- [ ] RFC-0002 is accepted.
- [ ] `spec/` contains the YAML files of laws, claims, gates, probes, reason codes, defect forms, machineries and stances, and `go generate ./spec` renders their JSON files, their schemas and the manifest.
- [ ] Every claim with a partner names one of the seven partner modes, and every claim of a mixin that implies a shallower claim states the pair in `entails`.
- [ ] The test of the package `spec` passes, and it fails for a classification of the shape catalog without a claim or an entry in `unclaimed`.
- [ ] Every claim names a probe, a defect form and a reason code that the spec defines, and every probe names assertions that the vendored copy of assert-spec defines.
- [ ] A release of dokimi contains the spec, so a runtime can vendor it from a tag.

## Why now

Nothing has to exist first. The derivation of checks in milestone 0003 reads the claims, and every later milestone names a claim, a probe or a reason code from this spec.

## Scope

- The spec as data, with its schemas, its rendering, its test and the vendoring scripts that the runtimes copy.
- A claim with its probe and defect form for each classification that a single caller can falsify.
- A stance for each classification that a single caller cannot falsify.

## Not in this milestone

- The Go implementation of each probe: goes to milestone 0003.
- The claims of the machineries above a single caller: go to milestones 0009, 0010, 0011 and 0013.
- The reference forms of fakes: go to milestone 0012.

## Risks to the sequence

| Risk | What it delays | What we would do |
|---|---|---|
| A partner of a claim fits none of the seven partner modes | Milestone 0003 and every milestone that derives checks | Give the claim a stance of `blocked` with the reason, and propose an eighth mode |
| A probe needs an assertion that assert-spec lacks | The claims of that probe | File the assertion with assert-spec, and give the claim the stance `needs` until it is released |

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 2 | The first roadmap of dokimi |
