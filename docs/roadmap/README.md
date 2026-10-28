<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Roadmap

What we are building, in what order, and what has to be true first. Milestone numbers are permanent. The order is the order of this table, and it changes.

```mermaid
flowchart LR
    M0001[0001 Lost check fails the gate]
    M0002[0002 Every classification has a law]
    M0003[0003 Derived checks in Go]
    M0004[0004 Vacuous and over-strong checks fail]
    M0005[0005 Suite record]
    M0006[0006 Misspelled brand fails]
    M0007[0007 Corpus judges each language]
    M0008[0008 TypeScript suites]
    M0009[0009 Calls fail on demand]
    M0010[0010 Laws of a history]
    M0011[0011 Concurrent callers]
    M0012[0012 Verified fakes]
    M0013[0013 Breakable medium]
    M0014[0014 Dead and subsumed checks]
    M0015[0015 Simulation of sans-IO machines]
    M0002 --> M0003
    M0003 --> M0004
    M0003 --> M0005
    M0004 --> M0005
    M0001 --> M0007
    M0004 --> M0007
    M0005 --> M0007
    M0007 --> M0008
    M0003 --> M0009
    M0003 --> M0010
    M0010 --> M0011
    M0003 --> M0012
    M0004 --> M0012
    M0003 --> M0013
    M0005 --> M0014
    M0009 --> M0015
    M0010 --> M0015
```

| Order | Milestone | Status | Depends on | Released in |
|---|---|---|---|---|
| 1 | [0001](0001-a-lost-check-fails-the-gate.md) A lost check fails the generation gate | Planned | none | unscheduled |
| 2 | [0002](0002-every-classification-has-a-law.md) Every classification has a law or a stance | In progress | none | unscheduled |
| 3 | [0003](0003-go-declarations-get-derived-checks.md) A Go declaration gets derived checks | Planned | 0002 | unscheduled |
| 4 | [0004](0004-vacuous-and-over-strong-checks-fail.md) Vacuous and over-strong checks fail | Planned | 0003 | unscheduled |
| 5 | [0005](0005-a-run-writes-its-suite-record.md) A run writes its suite record | Planned | 0003, 0004 | unscheduled |
| 6 | [0006](0006-a-misspelled-brand-fails-the-run.md) A misspelled brand fails the run | Planned | none | unscheduled |
| 7 | [0007](0007-the-corpus-judges-each-language.md) The corpus judges each language by its records | Planned | 0001, 0004, 0005 | unscheduled |
| 8 | [0008](0008-typescript-declarations-get-proven-suites.md) A TypeScript declaration gets a proven suite | Planned | 0007 | unscheduled |
| 9 | [0009](0009-calls-fail-on-demand.md) A check can make a call fail on demand | Planned | 0003 | unscheduled |
| 10 | [0010](0010-sequences-against-a-reference.md) Laws of a history are judged over generated sequences | Planned | 0003 | unscheduled |
| 11 | [0015](0015-a-simulation-of-sans-io-machines.md) A system of sans-IO machines runs in a deterministic simulation | Planned | 0009, 0010 | unscheduled |
| 12 | [0011](0011-concurrent-callers.md) Concurrent callers are checked for linearizability | Planned | 0010 | unscheduled |
| 13 | [0012](0012-verified-fakes.md) A downstream team uses a verified fake | Planned | 0003, 0004 | unscheduled |
| 14 | [0013](0013-a-breakable-medium.md) A durable subject survives a breakable medium | Planned | 0003 | unscheduled |
| 15 | [0014](0014-dead-and-subsumed-checks.md) A report names dead and subsumed checks | Planned | 0005 | unscheduled |

## Capabilities without a milestone

RFC-0001 gives each of these the verdict In. None has a position in the order yet.

| Capability | What it needs first |
|---|---|
| Clock skew, drift, rollback and bounded uncertainty | Milestone 0003 |
| Cross-version evolution | Milestone 0013 |
| Replay of recorded call sequences | Milestone 0010 |
| A finding for an expected value that flows from the subject | Milestone 0003 |
| Contradictions between a declaration and its implementation | Milestone 0010 |
| Directive proposals that a person accepts | Milestone 0006 |
| Packs of laws from third parties | Milestone 0008 |
| A signature over a bundle and its lock | Milestones 0001 and 0005 |
| A compiler from a counterexample of a model checker to a check | Milestone 0010 |
| A library of laws from incidents | Milestone 0002 |
| The plan beside the declaration in an editor | Milestone 0003 |
| Debugging that steps back and forth through a trace | Milestone 0015 |
| The workloads and checkers of a simulation against a real cluster | Milestone 0015 |
| Language bindings beyond TypeScript | Milestone 0008, and the rules of each language in eidos |
