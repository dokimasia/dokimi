---
rfc: 0021
title: Concurrent callers and the consistency of their histories
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-28
updated: 2026-11-02
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

A declaration with the mixin `concurrent` that states the machinery `reference` gets the claim `linearizable`: every history of concurrent calls on the subject has an order of its calls in which the reference, called one at a time, returns the same results. Five more claims check concurrent readers, a single flight of one computation, causal order, a snapshot under concurrent writes and an exclusive lease. Each generated check runs the actions of RFC-0020 on several clients through the scheduler of assert-go, records the history of calls and returns, and checks it. The planted defect of `linearizable` is a subject whose read and write race: correct when called one at a time, and wrong when a read overlaps a write.

## Motivation

The mixin `concurrent` states that a declaration is safe for concurrent use. The shape catalog classifies it, and the suite generator reads it only to run a benchmark in parallel (`generator/suite/call.go:143-146`). No check tests the claim. A subject that is correct for one caller and loses a write under two passes every check of a single caller and every judge of a single client.

The race detector of Go finds a data race on a path that runs, and it finds no logical race: two operations that each lock correctly and whose combination is wrong. A history checker finds those, from the calls and returns that a run records. assert-go implements the checkers: `history.Linearizable` with its options of budget, memory limit and workers, `SpecFrom`, which turns a sequential model into a specification, `history.Serializable` and `history.HasSnapshotIsolation`. Its package `stateful` runs a concurrent section of a machine through `NewScheduler` with the strategies `Uniform` and `PCT` (`assert-go/history`, `assert-go/stateful`). assert-spec defines `linearizable`, `serializable` and `snapshot-isolation`, and assert-typescript implements all three.

Porcupine checks linearizability for Go programs in the same way, and partitions a history by key to make checking practical. Checking is NP-complete in general, so a check states a budget.

## Detailed design

### The claims

| Claim | Licence | Machineries | Probe | Reason | Defect form |
|---|---|---|---|---|---|
| `linearizable` | The mixin `concurrent` | `concurrent-callers`, `generated-sequences`, `reference-implementation` | `linearizable-history` | `not-linearizable` | `racy-read-write` |
| `reads-agree-under-concurrency` | The mixin `concurrent-readers` | `concurrent-callers`, `generated-sequences` | `concurrent-reads-agree` | `inconsistent-concurrent-read` | `racy-read` |
| `computes-once` | The contract `single-flight`, role `fn` | `concurrent-callers`, `call-that-fails-on-demand` | `shares-one-computation` | `repeated-computation` | `computes-per-call` |
| `causal-history` | The mixin `causal` with its param `version` | `concurrent-callers`, `generated-sequences` | `versions-follow-the-causal-order` | `unobserved-write` | `racy-causal-read` |
| `pins-the-snapshot-under-writes` | The mixin `point-in-time` with its param `write` | `concurrent-callers`, `generated-sequences` | `pinned-under-concurrent-writes` | `unpinned-read` | `racy-pin` |
| `one-holder-at-a-time` | The contract `lease`, role `acquire`, with its param `taken` | `concurrent-callers` | `exclusive-under-concurrency` | `unguarded-write` | `check-then-grant` |

`linearizable` has the scope `declaration`. A declaration with the mixin `concurrent` and without `machineries=[reference]` lists `linearizable` as unclaimed with the stance `needs`.

This RFC does not propose the claims of the mixins `serializable` and `snapshot-isolation`. Their checks drive transactions, and the contract `tx` names only begin, commit and rollback, so the actions of a transaction need a design of their own.

A mixin classifies a callable, so a claim of the declaration takes its licence from a callable of the declaration with the mixin. The claim binds once, on the first such callable that passes the gates of the claim, or on the first such callable where none passes. Its identity starts with the name of the declaration, as in `Store/linearizable`.

The other five claims bind to the callable that their mixin or contract classifies:

- `reads-agree-under-concurrency` runs the generated sequences of the declaration on client 0, so the subject has state, and then a concurrent section of the read alone on three clients. Each read of the section on a client other than 0 must return what the read returns with the same arguments on client 0 after the section. A fresh subject returns one result for every input of a read, so reads without the sequences cannot tell a defect from the subject. `racy-read` keeps the results of a read in a field of its own, yields to the scheduler and returns the field, so a read that overlaps another returns the results of the other.
- `computes-once` calls `Hold` of the injector of RFC-0019 and makes three clients call the callable with one key. It fails when the dependencies record more than one computation. The check runs inside a `testing/synctest` bubble through the field `Around` of the check, so the probe waits until every client blocks before it counts the calls of the dependencies and calls `Release`. A call after the computation finished starts another computation by design, so the dependencies wait until every client has called. The generated call passes each parameter of a function type a computation that calls a dependency of the injector, so `Hold` keeps the computation waiting and `Calls` counts it. A subject that computes through its own dependencies records its computations through the door `Faults`. `computes-per-call` calls a fresh subject for each call, so each call computes.
- `pins-the-snapshot-under-writes` runs the generated sequences of the declaration with a concurrent section of the read and its `write` on three clients. Each other callable runs on client 0 alone. Each read that succeeds must return the output of the first read with the same arguments that succeeded.
- `causal-history`, `pins-the-snapshot-under-writes` and `one-holder-at-a-time` are the deeper claims of the laws of `causal`, `point-in-time` and `lease`, whose depth is `concurrent-callers`. Each entails the claim of a single caller and reports its reason code. Each defect is correct when called one at a time and races through a yield of the scheduler between a check and its effect, as `racy-read-write` does.
- `causal-history` judges the history of RFC-0020 with the judge `ordered` over the stamps that `version` reads. A write that a client issued after another write of the same client, or after a read that observed the other write, has a higher stamp, and each read observes the stamps of the writes before it in that order.
- A read of `causal-history` observed the write with its keys and the stamp that it returns, where that write started before the read returned. A stamp that two writes with the same keys share names neither write. Each later call of a client follows the stamp of each read of that client, so a client does not read a lower stamp of a key than it read before. A read that follows a write of its keys returns the stamp of a write of those keys. The claim takes the gate `write-returns-the-value`, because the probe reads the stamp of each write from the value that the write returns.

### The generated check

```go
{
	ID:        conform.ID[catalog.Store]("Store/linearizable"),
	Statement: "every concurrent history of Store has an order in which the reference returns the same results",
	Claim:     spec.Linearizable,
	Reason:    spec.NotLinearizable,
	Property: func(h *conform.Handle[catalog.Store], c *prop.Case) {
		probe.LinearizableHistory(h, c, storeActions, storeReference, 2,
			history.Budget(200_000))
	},
	Defect: func(h *conform.Handle[catalog.Store]) catalog.Store {
		return &storeRacyReadWrite{Store: h.Fresh()}
	},
},
```

```go
// LinearizableHistory runs generated sequences of the actions on the subject
// of h through stateful.Steps with stateful.Clients(clients) and the schedule
// of h. The actions record each call in the history of c, as RFC-0020 states.
// After the last step it checks the history once with history.Linearizable
// and opts, against the spec that history.SpecFrom builds over fresh
// references of h through reference. The check partitions the history by the
// keys of the calls. The probe engages h only when the checker decides the
// history.
func LinearizableHistory[S any](
	h *conform.Handle[S],
	c *prop.Case,
	actions Actions[S],
	reference Reference[S],
	clients int,
	opts ...history.Option,
)
```

A concurrent section of `stateful.Clients(2)` runs on the clients 0 to 2, so three clients call the subject at once. The actions and `storeReference` are those of RFC-0020, so the agreement with the reference and the check of linearizability read one history and one reference. The machine of the probe has no spec. A spec on the machine would check the history after every step with every call in one partition, and the probe checks it once after the last step, split by key, within the budget of the check. A call whose context ends before the subject returns completes as `Unknown`, and the checker lets it take effect or not. A history that the budget cannot decide does not engage the case. A property check engages when at least one case engages, so a check whose every case remains undecided is vacuous, and the census of RFC-0015 fails it.

### The schedule of the clients

A check of a subject runs the clients of its concurrent section on threads, `stateful.Repeat` times. A subject that never yields runs each call whole as a task, so only threads interleave its calls. A check of a planted defect runs the clients as tasks of a scheduler with the strategy `Uniform`, which chooses the next task at each yield, so the race of the defect shows in a few cases and replays from the seed of its case. A probe of concurrent callers keeps every action in each case, so each case can show a race between two actions. The handle gives the probe the option and the planted defect the yield:

```go
// Schedule returns the option of stateful.Steps that runs the clients of a
// concurrent section of c: on threads for a subject, and as tasks of a
// scheduler of c with the strategy Uniform for a planted defect.
func (h *Handle[S]) Schedule(c *prop.Case) stateful.Option

// Yield lets the scheduler of the case run another client while the clients
// of a planted defect run as tasks, and returns at once otherwise. A planted
// defect calls it between the two halves of an operation that a correct
// subject runs whole.
func (h *Handle[S]) Yield()
```

### The probes of the other claims

Each probe of generated sequences runs its clients through `stateful.Steps` with `stateful.Clients(2)`, the schedule of the handle and `stateful.Swarm(false)`, as `LinearizableHistory` does. `SharesOneComputation` runs its three clients on goroutines of the bubble of its check, because `synctest.Wait` returns once every call waits on a held dependency.

```go
// ConcurrentReadsAgree runs generated sequences of the actions on the subject
// of h on client 0, so the subject has state, and then a concurrent section of
// the action read alone on the clients 0 to clients, which the schedule of h
// runs. Each case keeps every action. After the section it calls again with
// the arguments of each read that a client other than 0 completed, and fails
// when the output differs. Each other action runs on client 0 alone, because
// Steps leaves an action with Enabled out of a concurrent section. A case
// without a completed read on a client other than 0 does not engage h.
func ConcurrentReadsAgree[S any](
	h *conform.Handle[S], c *prop.Case, actions Actions[S], read string, again Apply, clients int,
)

// PinnedUnderConcurrentWrites runs generated sequences of the actions on the
// subject of h, with a concurrent section of the actions read and write on
// the clients 0 to clients, which the schedule of h runs. Each other action
// runs on client 0 alone. After the last step it compares each read that
// succeeds with the first read that succeeded with the same arguments, and
// fails when the outputs differ. A case without a read that succeeds after
// another read with the same arguments does not engage h.
func PinnedUnderConcurrentWrites[S any](
	h *conform.Handle[S], c *prop.Case, actions Actions[S], read, write string, clients int,
)

// VersionsFollowTheCausalOrder runs generated sequences of the actions on the
// subject of h, with a concurrent section of the actions read and write on
// the clients 0 to clients, which the schedule of h runs. Each other action
// runs on client 0 alone. After the last step it orders the calls of read and
// write that succeed. A call follows each call of its client before it. A
// read follows the write with the same keys and the version that the read
// returns, where that write started before the read returned and no other
// write with the same keys has that version. The probe fails when a write
// returns a version that is not above the version of a call with the same
// keys that it follows, and when a read returns a version below it. It also
// fails when a read that follows a write with the same keys returns a version
// that no write with those keys returned. Each call of write returns the value
// that it wrote. A case without a call that follows a call with the same keys
// does not engage h.
func VersionsFollowTheCausalOrder[S, V any, O cmp.Ordered](
	h *conform.Handle[S], c *prop.Case, actions Actions[S], read, write string, version func(v V) O, clients int,
)

// SharesOneComputation holds the dependencies of h, and the clients 0 to
// clients make the call at once, each on a goroutine of its own. Once every
// goroutine of the bubble of the check is blocked, it counts the calls that
// the dependencies recorded since the hold, releases the dependencies and
// waits for each call. It fails when the dependencies recorded more than one
// call, each call one computation. A check whose dependencies do not record a
// call does not engage h. The probe waits through testing/synctest, so its
// check runs inside a bubble that the field Around of the check sets. A call
// that panics fails the check with the value of the first panic once each
// call has returned.
func SharesOneComputation[S, V any](h *conform.Handle[S], call Call[V], clients int)

// ExclusiveUnderConcurrency runs generated sequences of the call on the
// subject of h, with a concurrent section on the clients 0 to clients that
// the schedule of h runs. The sequences run the call alone, so the lease that
// the first call to succeed takes stays taken, and each other call fails. It
// fails when more than one call succeeds, and when a call fails with an error
// that does not match sentinel. A case without a call that succeeds does not
// engage h.
func ExclusiveUnderConcurrency[S any](h *conform.Handle[S], c *prop.Case, call Act, sentinel error, clients int)
```

The machine of `ExclusiveUnderConcurrency` has the one action of its call, so the check of `one-holder-at-a-time` does not need the sequences of the declaration. The generated check is still a property check, because the scheduler that runs the clients of a planted defect belongs to the case:

```go
Property: func(h *LockerHandle, c *prop.Case) {
	probe.ExclusiveUnderConcurrency(h, c, func(ctx context.Context) error {
		return h.Subject().Lock(ctx, string("key-a"))
	}, ErrLocked, 2)
},
```

The generated call of `computes-once` passes a computation to each parameter of a function type, also of a defined function type. The computation records a call of the callable on a new dependency of the injector and returns the zero value of each result:

```go
probe.SharesOneComputation(h, func(ctx context.Context) (string, error) {
	return h.Subject().Do(ctx, string("key-a"), func(context.Context, string) (r0 string, r1 error) {
		_ = h.Injector().Dependency().Call("Do")
		return r0, r1
	})
}, 2)
```

### The planted defect

`racy-read-write` wraps a fresh subject. Its write stores the value in two steps with a call of `Handle.Yield` between them, so a read that the scheduler places between the steps returns a half-written state. Called one at a time it is correct, so the checks of a single caller and the judges of a single client pass on it.

The proof gives the same result on every run, because the yield of the defect goes through the scheduler of the case. Falsification runs the check with the strategy `Uniform` and a seed from the key of the check. `PCT(3)` switches tasks only at its two change points, so it proved the torn store of the specs at 15 of 40 seeds. `Uniform` with every action in each case proved it at all 40, so the proof does not depend on the name of the declaration. The racy lease of the specs of `ExclusiveUnderConcurrency` is proven at all 40 seeds as well.

The check of linearizability compares each call of each callable with the reference, as the check of agreement with a reference does. So the isolation of a planted defect leaves out both checks, for the defect of any claim.

Each form wraps a fresh subject, as every planted defect does, and races only through `Handle.Yield`. Called one at a time, each form returns what the subject returns:

| Form | Changes | What the defect does |
|---|---|---|
| `racy-read-write` | declaration | Every callable but the first callable with a value result marks its call, yields, and then calls the subject. The first callable with a value result returns its second sample while another call is marked |
| `racy-read` | callable | The callable keeps its results in a field of the defect, yields, and returns the field, so a call that overlaps another returns the results of the other |
| `computes-per-call` | callable | The callable calls a fresh subject for each call, so each call computes |
| `racy-causal-read` | callable | The read of `causal` keeps its results in a field, yields, and returns the field, as `racy-read` does |
| `racy-pin` | callable | The pinned read keeps its results in a field, yields, and returns the field, as `racy-read` does |
| `check-then-grant` | callable | Each call of acquire marks itself and yields. A call that finds another call marked returns the zero value of each result and no error without the subject. Any other call calls the subject |

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
