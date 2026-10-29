---
milestone: 0003
title: A Go declaration gets derived checks
status: Planned
depends-on: 0002
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0014, RFC-0027
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0003: A Go declaration gets derived checks

## Goal

A Go team that writes `//dokimi:suite` on an interface gets a check for each claim that the classifications of its callables license, each proven against its planted defect by its reason code, without writing a check.

## Done when

- [ ] RFC-0014 is accepted.
- [ ] `dokimi run` on a reader with the mixin `read-after-write` and its param `write` writes the check `read-after-write` with its planted defect, and a reader whose param names no callable lists the claim among its refusals with the remedy.
- [ ] The module `go.dokimi.dev/spec` contains the copy of the conformance spec with its constants, and the module `go.dokimi.dev/conform` contains the package `probe`, with a Go probe for every probe of the spec that a single caller runs.
- [ ] Falsification of a derived check fails with both reason codes when the planted defect fails for another reason.
- [ ] Under `conform.Isolate`, falsification of a pair of `entails` fails with both identities when the deeper check passes on the defect of the shallower one, and the corpus passes the option for the subject of each fixture.
- [ ] Each of the seven partner modes binds a partner in a fixture of the corpus, and refuses with `partner-found` and `one-partner` in two more.
- [ ] A declaration without `machineries=reference` lists a claim that needs a reference as unclaimed with the stance `needs`.
- [ ] Each fixture of the Go corpus whose laws the spec states uses derived checks, and its tests pass.
- [ ] An authored check of the directive `check` proves by the contract of its failure record.
- [ ] `conform.Run` falsifies each derived check on a planted defect that wraps a fresh subject before it applies the check, and writes the falsification record and the record of the subject.

## Why now

The claims, probes and reason codes of milestone 0002 have to exist first. Vacuity, the suite record, the corpus gate, TypeScript and every machinery above a single caller build on derived checks.

## Scope

`Derive` of RFC-0003 with the partner modes and the param `machineries` of RFC-0014, the generated checks and the planted defects that wrap a fresh subject, the module `spec` and the package `probe` of conform, the fields `Claim`, `Reason`, `Law`, `Machineries` and `Shallow` of a check, the methods `Expect` and `Fresh` of the handle, the falsification of derived checks inside `conform.Run`, the isolation and the edges under `conform.Isolate`, the generators `enum` and `sentinel` on the claims of the spec through `conform.Falsify` and `conform.Run`, and the migration of the Go corpus.

## Not in this milestone

- Vacuous checks and controls: go to milestone 0004.
- The suite record: goes to milestone 0005.
- Derived checks in TypeScript: go to milestone 0008.

## Risks to the sequence

| Risk | What it delays | What we would do |
|---|---|---|
| The migration of 146 fixtures takes longer than the derivation | Milestone 0004 and every later milestone | Migrate the fixtures of conformance level 1 first, and keep the others authored until the corpus gate needs them |
| A gate needs a fact that the rules of eidos do not assign | The claims of that gate | File the request with eidos, and refuse the claim with the gate until it is released |

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 3 | The first roadmap of dokimi |
| 2026-10-28 | The isolation and the edges run under `conform.Isolate`, a check gains `Machineries` and `Shallow`, and the corpus counts 146 fixtures | A defect of the subject of a consumer fails the isolation of every planted defect that wraps it, so the controls of the corpus isolate the derived defects. The suite record needs the machineries of each check |
| 2026-10-28 | The corpus isolates the derived defects over the subjects of its fixtures | The controls of the corpus belong to milestone 0004, and the subject of each fixture is a correct implementation |
| 2026-10-28 | The scope gains the generators `enum` and `sentinel` on the claims of the spec | Their checks are claims of the spec, and no milestone contained the change of the generators |
| 2026-10-29 | The copy of the spec is in the module `go.dokimi.dev/spec` at the root of dokimi-go | Every module of the Go runtime reads one copy, and `go:embed` reads the files of its own module alone |
| 2026-10-29 | The Go probes follow RFC-0027, and the probes of enumerations gain the inputs `list`, `valid`, `outside` and `name` | The generated checks call each probe with a fixed Go signature, and every runtime passes the functions that a probe of an enumeration calls |
