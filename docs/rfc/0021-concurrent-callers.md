---
rfc: 0021
title: Concurrent callers and the consistency of their histories
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Draft
created: 2026-10-28
updated: 2026-10-28
discussion: none
supersedes: none
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0021: Concurrent callers and the consistency of their histories

## Summary

A declaration with the mixin `concurrent` that states the machinery `reference` gets the claim `linearizable`: every history of concurrent calls on the subject has an order of its calls in which the reference, called one at a time, returns the same results. A declaration with the mixin `serializable` or `snapshot-isolation` gets the claim of that isolation level over the histories of its transactions. Each generated check runs the actions of RFC-0020 on several clients through the scheduler of assert-go, records the history of calls and returns, and checks it with the checker of assert-go for its level. The planted defect of `linearizable` is a subject whose read and write race: correct when called one at a time, and wrong when a read overlaps a write.

## Motivation

The mixin `concurrent` states that a declaration is safe for concurrent use. The shape catalog classifies it, and the suite generator reads it only to run a benchmark in parallel (`generator/suite/call.go:143-146`). No check tests the claim. A subject that is correct for one caller and loses a write under two passes every check of a single caller and every judge of a single client.

The race detector of Go finds a data race on a path that runs, and it finds no logical race: two operations that each lock correctly and whose combination is wrong. A history checker finds those, from the calls and returns that a run records. assert-go implements the checkers: `history.Linearizable` with its options of budget, memory limit and workers, `SpecFrom`, which turns a sequential model into a specification, `history.Serializable` and `history.HasSnapshotIsolation`. Its package `stateful` runs a concurrent section of a machine through `NewScheduler` with the strategies `Uniform` and `PCT` (`assert-go/history`, `assert-go/stateful`). assert-spec defines `linearizable`, `serializable` and `snapshot-isolation`, and assert-typescript implements all three.

Porcupine checks linearizability for Go programs in the same way, and partitions a history by key to make checking practical. Checking is NP-complete in general, so a check states a budget.

## Detailed design

### The claims

| Claim | Licence | Machineries | Probe | Reason | Defect form |
|---|---|---|---|---|---|
| `linearizable` | The mixin `concurrent` | `concurrent-callers`, `a-reference-implementation` | `linearizable-history` | `not-linearizable` | `racy-read-write` |
| `serializable` | The mixin `serializable` | `concurrent-callers` | `serializable-history` | `not-serializable` | `reads-uncommitted` |
| `snapshot-isolation` | The mixin `snapshot-isolation` | `concurrent-callers` | `snapshot-history` | `lost-update` | `writes-without-a-snapshot` |

Each claim has the scope `declaration`. A declaration with the mixin `concurrent` and without `machineries=reference` lists `linearizable` as unclaimed with the stance `needs`. The isolation levels do not need a reference, because their checkers judge the reads and writes of the history alone.

### The generated check

```go
{
	ID:        conform.ID[catalog.Store]("Store/linearizable"),
	Statement: "every concurrent history of Store has an order in which the reference returns the same results",
	Claim:     spec.Linearizable,
	Reason:    spec.NotLinearizable,
	Property: func(h *conform.Handle[catalog.Store], c *prop.Case) {
		probe.LinearizableHistory(h, c, storeActions, storeSubject, 2,
			history.Budget(200_000))
	},
	Defect: func(h *conform.Handle[catalog.Store]) catalog.Store {
		return &storeRacyReadWrite{Store: h.Fresh()}
	},
},
```

```go
// LinearizableHistory runs the actions on a fresh subject of h through
// stateful.Steps with stateful.Clients(clients). The actions record each
// call in the history of c, as RFC-0020 states. It checks the history with
// history.Linearizable and opts against the spec that history.SpecFrom
// builds from subject over the reference of h, partitioned by the keys of
// the calls. It engages h only when the checker decides the history.
func LinearizableHistory[S any](
	h *conform.Handle[S],
	c *prop.Case,
	actions Actions[S],
	subject func(reference S) history.Subject,
	clients int,
	opts ...history.Option,
)
```

A concurrent section of `stateful.Clients(2)` runs on the clients 0 to 2, so three clients call the subject at once. The actions and `storeSubject` are those of RFC-0020, so the agreement with the reference and the check of linearizability read one history and one spec. A call whose context ends before the subject returns completes as `Unknown`, and the checker lets it take effect or not. A history that the budget cannot decide does not engage the case. A property check engages when at least one case engages, so a check whose every case remains undecided is vacuous, and the census of RFC-0015 fails it.

### The planted defect

`racy-read-write` wraps a fresh subject. Its write stores the value in two steps with a yield of the scheduler between them, so a read that the scheduler places between the steps returns a half-written state. Called one at a time it is correct, so the checks of a single caller and the judges of a single client pass on it.

The proof gives the same result on every run, because the yield of the defect goes through the scheduler of the case. Falsification runs the check with the strategy `PCT` and a seed from the key of the check, and that seed places the read between the two steps. The corpus pins that the proof fails at that seed.

### Partitioning

A declaration whose callables take a key states the key through the binding `key` of its shape. The probe partitions the history by key, as porcupine does, so each partition is checked alone.

## Alternatives considered

### Controlled interleavings

A test chooses the schedule of the subject, as loom and Shuttle do for Rust.

**Why not:** Go offers no hook into its scheduler, and both tools need the subject to use their own synchronization types. A check of a history needs no change of the subject. The scheduler of assert-go orders the calls of the clients, not the goroutines inside the subject.

### The race detector alone

The generated suite runs under `-race`.

**Why not:** the race detector finds unsynchronized memory access, not a combination of correctly synchronized operations that breaks a law. A consumer runs the suite with `-race`, so its gate runs the race detector and the history check together.

## Drawbacks

- **Budgets.** A history check can exceed its budget and leave the case undecided. A check whose every case is undecided is vacuous, and the consumer raises the budget or partitions the history.
- **Slow cases.** A case runs three clients and checks a history of up to 150 operations.
- **A reference that is not safe for concurrent use.** The checker calls the reference one operation at a time, so the reference need not be safe, and the generator writes the sequential wrapper.
- **Transactions as actions.** The isolation checkers of assert-go read the operation `txn`, whose micro-operations append to lists and read them, or write registers and read them. The actions of a declaration with the contract `tx` record each scope as one such transaction, with a unique value for each write. A declaration whose scope cannot be driven as either form refuses the claims with the gate `transactions`. A history of registers shows fewer anomalies, because a blind write reveals no order.

## Open questions

None.

## References

| What | Where |
|---|---|
| porcupine, a linearizability checker for Go | <https://github.com/anishathalye/porcupine> |
| Elle, inferring isolation anomalies from histories | <https://arxiv.org/abs/2003.10554> |
| loom | <https://docs.rs/loom> |
| Go race detector | <https://go.dev/doc/articles/race_detector> |
| assert-go, the package history | <https://github.com/dokimasia/assert-go/tree/main/history> |
