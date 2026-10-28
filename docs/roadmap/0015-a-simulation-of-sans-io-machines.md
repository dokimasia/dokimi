---
milestone: 0015
title: A system of sans-IO machines runs in a deterministic simulation
status: Planned
depends-on: 0009, 0010
released-in: unscheduled
deadline: none
deadline-source: none
prd: none
rfc: RFC-0025
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Milestone 0015: A system of sans-IO machines runs in a deterministic simulation

## Goal

A Go team whose components are sans-IO machines writes directives on their declarations and gets a simulation of the whole system under faults, with the workloads and checkers that Jepsen runs. A failing run shrinks to a short log and history, and replays from its seed.

## Done when

- [ ] RFC-0025 is accepted.
- [ ] dokimi-go has released the module `go.dokimi.dev/sim` with the world, the clocks, the epochs, the faults, the strategies `uniform` and `pct`, the log and its checksum.
- [ ] `dokimi run` writes the simulation of a fixture of the corpus with two machines, from the directives `sim`, `token`, `reply`, `perform`, `system`, `action`, `spec`, `invariant`, `converges`, `fingerprint` and `scenario`.
- [ ] An effect type without `perform` fails generation with the type and the remedy.
- [ ] Each of the five rules of a runtime fails on its planted defect with its reason code at the seed of falsification, and each invariant of the fixture fails on its defect system.
- [ ] A fixture whose machine loses an acknowledged write under a partition fails `history.Linearizable` against its spec, and a fixture of list-append transactions with a write skew fails `history.Serializable`.
- [ ] A call whose process crashes before it replies completes as `Unknown` in the history, and a result of the old epoch never reaches the restarted machine.
- [ ] A fixture whose machine iterates a Go map into its effects fails the determinism check with the first step at which two runs differ.
- [ ] A failing run of the fixture prints the shrunk log with the step, the epoch, the clock and the faults of each process, and the history of the calls, and a rerun with its seed gives the same log at `GOMAXPROCS` 1 and 4.
- [ ] The run never exceeds the faults of `tolerate`, and `require` fails a suite whose cases seldom crash a process.
- [ ] The crash points of a scenario of the fixture run the invariants after a crash before and after each event.
- [ ] sharder runs the simulation of the machines of one namespace with the generated code.

## Why now

The world runs as actions of a machine of `stateful` and records the calls of its clients in the history of milestone 0010. It arms the faults of its dependencies through the `Injector` of milestone 0009. Both have to exist first. It comes before the milestones that follow them in the order because the first milestone of sharder waits on it, and it needs nothing that those milestones deliver.

## Scope

The directives and the generator `sim` with its Go template tree, the module `go.dokimi.dev/sim`, the machinery `a-simulated-system` and the five claims of the rules of a runtime in the spec, and a fixture of the corpus.

## Not in this milestone

- A simulation in TypeScript: not proposed in RFC-0025.
- Byzantine faults: out, in RFC-0001.
- Debugging that steps through a trace in the service: the service builds it.
- A runner that applies the same checkers to a run on a real cluster: designed in RFC-0026, and nobody has scheduled it yet.

## Risks to the sequence

| Risk | What it delays | What we would do |
|---|---|---|
| The waits of sharder need a performer that the rules of a performer forbid, such as one that starts a goroutine | The simulation of sharder | Give the world a wait that its own doubles drive, and request the change from the dependency |

## Changes

| Date | What changed | Why |
|---|---|---|
| 2026-10-28 | Created at position 11 | dokimasia/dokimi#43, the simulation that the first milestone of sharder needs |
