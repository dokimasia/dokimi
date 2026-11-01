---
milestone: 0016
title: A derived check runs over generated inputs
status: Done
depends-on: 0003, 0005
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0006
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0016: A derived check runs over generated inputs

## Goal

A Go team whose declaration gets a derived check of a single caller also gets the same check over inputs that the property engine generates and shrinks. A declaration with the mixin `commutative` or `associative` gets the check of its law over generated pairs and triples of inputs.

## Done when

- [x] RFC-0006 is accepted.
- [x] The suite generator writes a generated binding with the case `generated` beside each derived check whose sample binding passes sample inputs to a call, under the machinery `generated-inputs`.
- [x] The claims `commutative` and `associative` bind under `generated-inputs`, and the unclaimed table of the spec loses their entries.
- [x] Each generated binding fails on the planted defect of its claim with its reason code at the seed of its key.
- [x] A failing run on a subject writes the counterexample of its smallest case into the suite record, and a rerun with its seed repeats it.
- [x] Each claim whose law has the depth `generated-inputs`, such as `deterministic` and `pure`, has a generated binding that is not shallow wherever its calls take sample inputs.
- [x] `dokimi conform --lang go` passes with the generated bindings in the expected outcomes of the corpus.

## Why now

The derived checks of milestone 0003 and the suite records of milestone 0005 have to exist first. It comes right after milestone 0009 and before the machineries of sequences and concurrency. Every claim of a single caller becomes a property test with the call generators and the property checks that the suite generator already writes for the fuzz targets. Fourteen claims of the spec witness their law in full only over generated inputs.

## Scope

The generated bindings of the derived checks, the claims `commutative` and `associative` with their probes, gates, reasons and defect forms, and the expected outcomes of the corpus.

## Not in this milestone

- Fuzz targets, benchmarks and budgets: the suite generator writes them already.
- Generated bindings in TypeScript: milestone 0008 binds TypeScript, and each capability of Go after it adds its TypeScript part.

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-31 | Created at position 9 | The generated bindings of RFC-0006 had no milestone, and ten claims witness their law in full only over generated inputs |
| 2026-10-31 | In progress | Milestone 0009 waits for the release of conform. ADR-0019 states how a generated binding draws its inputs |
| 2026-10-31 | The scope gains the gate `exact-fidelity` of the claim `codec`, two probes that write before their first read, and a probe of a schedule that skips a delay that is not positive | The generated bindings of the corpus failed on a lossy codec, on a drawn key without a value, and on a delay of 0 |
| 2026-10-31 | Six of the seven boxes are ticked | `dokimi conform --lang go` passes with 173 conformant fixtures against the conform and the spec of dokimi-go before their release. The last box waits for that release and the corpus that requires it |
| 2026-11-01 | Done | The corpus requires the conform and the spec of dokimi-go at 1e87d26, and `dokimi conform --lang go` passes with 176 conformant fixtures through the workspace of the repository |
