---
rfc: 0031
title: Judges of deliveries, iterations, totals and failed calls
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-01
updated: 2026-11-01
discussion: none
supersedes: none
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0031: Judges of deliveries, iterations, totals and failed calls

## Summary

Seven classifications of the shape catalog state laws of a history that a probe cannot check from the calls of a case and their results alone. An outbox, a publisher and a watcher deliver their messages later, on a channel or a subscription. An iteration runs across the steps that mutate its source. A conserved quantity is a total over the subject, and a failed call needs a dependency that fails on demand. This RFC gives each law a probe that runs in the generated sequences of its declaration, with its claim, its law, its reason code and its planted defect.

## Motivation

The probes of a history read the calls and the results that the actions of a case record. Four kinds of law need more than that record:

- **Deliveries.** `outbox` and `publisher` deliver each message on the channel that `subscribe` returns, and `watcher` delivers each update through the method `next` of a subscription. A delivery arrives after the call that caused it, on another goroutine.
- **Iterations.** `stream-reflects-mutations` states that an iteration that has started sees the writes and the deletes that its source takes while it runs. The iteration spans several steps.
- **Totals.** `conservative` states that its callable moves a quantity and neither creates nor destroys any of it. The check reads the total of that quantity before and after each call.
- **Failed calls.** `atomic` states that a call that fails leaves no trace. The check needs a dependency that fails during one call of a sequence.

Without these probes the classifications keep their stances, and a declaration with one of them gets no check of its law over sequences.

## Detailed design

### The claims

| Claim | Licence | Law | Reason | Defect form | Entails |
|---|---|---|---|---|---|
| `conserves-the-quantity` | The mixin `conservative` with its params `field` and `observe` | `moves-conserve-the-total` | `not-conserved` | `leaks-one-unit` | — |
| `no-trace-of-a-failed-call` | The mixin `atomic` with its param `read`, and the machinery `faults` | `failed-call-leaves-no-effect` | `partial-write` | `half-applied-after-three` | `no-partial-write` |
| `iteration-sees-mutations` | The mixin `stream-reflects-mutations` with its params `mutate` and `delete` | `iteration-reflects-its-source` | `stale-iteration` | `iterates-a-copy` | — |
| `delivers-each-append` | The contract `outbox`, role `append`, with its role `subscribe` | `each-append-is-delivered` | `lost-message` | `drops-appends-after-three` | — |
| `delivers-by-its-mode` | The contract `publisher`, role `publish`, with its param `mode` and its role `subscribe` | `deliveries-follow-the-mode` | `wrong-delivery-count` | `breaks-the-mode-after-three` | — |
| `redelivers-on-demand` | The contract `publisher`, role `redeliver`, with its roles `publish` and `subscribe`, where `publish` states the mode `at-least-once` | `redelivery-delivers-again` | `not-redelivered` | `skips-a-redelivery-after-three` | — |
| `watch-sees-each-trigger` | The contract `watcher`, role `watch`, with its params `next` and `stop` and its role `trigger` | `watch-returns-each-trigger` | `missed-update` | `misses-updates-after-three` | — |

Each claim lists `generated-sequences` first among its machineries. `no-trace-of-a-failed-call` lists `call-that-fails-on-demand` after it. It entails `no-partial-write`, and its law `failed-call-leaves-no-effect` takes the depth `generated-sequences`, as each law of a history does. The other six laws are new, with the depth `generated-sequences`:

| Law | Statement |
|---|---|
| `moves-conserve-the-total` | Each call of the callable leaves the total that `observe` reads unchanged. |
| `iteration-reflects-its-source` | The rest of an iteration yields each value that a mutation wrote while it ran and no value that a delete removed while it ran. |
| `each-append-is-delivered` | Each message that an append recorded arrives on the channel of a subscription at least once. |
| `deliveries-follow-the-mode` | Each published event arrives at least once, at most once or exactly once, as `mode` states. |
| `redelivery-delivers-again` | Each event that a redelivery passed arrives once more than its publishes alone deliver. |
| `watch-returns-each-trigger` | A subscription returns each update that a trigger of its key caused after the watch started, in order, and no update after its stop. |

| Reason | Meaning |
|---|---|
| `not-conserved` | A call of the callable changed the total that `observe` reads. |
| `stale-iteration` | The rest of an iteration missed a value that a mutation wrote, or yielded a value that a delete removed. |
| `lost-message` | An appended message did not arrive on the channel of a subscription. |
| `wrong-delivery-count` | A published event arrived another number of times than `mode` allows. |
| `not-redelivered` | A redelivered event did not arrive once more. |
| `missed-update` | A subscription missed an update, returned updates out of order, or returned an update after its stop. |

The mixin `conservative` has the param `observe`, the callable that returns the total, as the mixin `bounded` has one for its use. The partners of the roles `subscribe` and `trigger` come from the mode `role`, and the param `observe` of `conservative` from the mode `named`. The partner resolution refuses `conserves-the-quantity` with the gate `partner-found` where a declaration does not state `observe`. A partner that a claim names counts as bound in the unclaimed table.

The suite generator judges five gates when it writes the code of a check:

| Gate | Claims | Passes when |
|---|---|---|
| `numeric-total` | `conserves-the-quantity` | The callable that `observe` names returns an integer or a float, so a planted defect can add a unit to the total. |
| `delivers-on-a-channel` | `delivers-each-append`, `delivers-by-its-mode`, `redelivers-on-demand` | The callable of the role `subscribe` returns a channel that a caller receives from, and the callable passes a value of the type of its elements. |
| `at-least-once-mode` | `redelivers-on-demand` | The callable of the role `publish` states the mode `at-least-once`, so a redelivery delivers an event once more. |
| `stream-mutations` | `iteration-sees-mutations` | The callable returns a stream, and the callables that `mutate` and `delete` name each take a value of the type of its elements. |
| `subscription-methods` | `watch-sees-each-trigger` | The method that `next` names takes a wait and returns an update beside a presence flag, the method that `stop` names takes and returns nothing, and the trigger passes a value of the type of the update. |

The reason of each gate contains the name of the partner that the claim binds, so a refusal shows which callable does not fit.

A redelivery delivers an event again only under `at-least-once`. Under `exactly-once` and `at-most-once` a second delivery of an event breaks the mode, so a redelivery there offers again only what did not arrive. The actions of a sequence call the redelivery too, so `delivers-by-its-mode` judges each redelivery by the mode that the publisher states.

### The probes

The probes run in the generated sequences of their declaration and judge the history after each step, as the probes of a history do. Each one adds what its law needs to the machine of its case:

- `moves-keep-the-total` wraps the action of the callable, and calls `observe` right before and right after each call. It fails when a call changed the total, and when `observe` fails.
- `failed-calls-leave-no-trace` adds an action that runs the action of the callable with every call of a dependency failing. Before and after that call it reads back through `read` with the arguments of the call. It fails when a call that returned the error of a failed dependency changed what `read` returns.
- `iterations-see-the-mutations` starts an iteration after the first step at which the stream yields a value, because the stream of a fresh subject yields none, and pulls its first value. After the last step it pulls the rest. It fails when the rest misses a value that a call of `mutate` wrote during the iteration, unless the first value is that value or a later delete removed it. It also fails when the rest yields a value that a call of `delete` removed and no later write wrote again.
- `appends-arrive` subscribes before the first step. After the last step it waits until each goroutine of the case is blocked, takes each message that the channel holds, and fails when a message that an append recorded is missing.
- `deliveries-match-the-mode` subscribes before the first step and takes the deliveries after the last step. It compares the published events and the delivered events as multisets. It fails when an event arrived fewer times than it was published under `at-least-once` and `exactly-once`, or more times under `at-most-once` and `exactly-once`.
- `redeliveries-arrive` subscribes before the first step. It fails when an event arrives fewer times than its publishes and its redeliveries together.
- `watches-see-the-triggers` takes the subscription of each call of `watch` from its output. After the last step it reads `next` until it reports no update, and compares the updates with the triggers of the key that completed after the watch. It then stops each subscription, triggers its key once more, and fails when `next` returns that update.

The probes of deliveries run inside a bubble of `testing/synctest`, as the probes of a clock do. `synctest.Wait` returns once each goroutine of the bubble is blocked, so the probe waits exactly until the subject has delivered what it will deliver. The timeout that `next` takes elapses on the clock of the bubble without a real wait.

```go
// Pull returns a pull iteration of the stream of a callable: next returns
// the next value of the stream and whether it had one, and stop ends the
// iteration. The generator writes it with iter.Pull or iter.Pull2, by the
// drain of the stream.
type Pull[E any] func(ctx context.Context) (next func() (E, bool), stop func(), err error)

// MovesKeepTheTotal fails when a call of call changed the total that observe
// reads, and when observe fails.
func MovesKeepTheTotal[S, T any](
	h *conform.Handle[S], c *prop.Case, actions Actions[S], call string, observe Call[T],
)

// FailedCallsLeaveNoTrace runs call with every call of a dependency failing,
// and fails when a failed call changed what readback returns for its
// arguments.
func FailedCallsLeaveNoTrace[S any](h *conform.Handle[S], c *prop.Case, actions Actions[S], call string, readback Apply)

// IterationsSeeTheMutations fails when the rest of an iteration of stream
// misses a value that mutate wrote, or yields a value that remove removed,
// during the iteration. written and removed return the value that a call of
// mutate writes and a call of remove removes.
func IterationsSeeTheMutations[S, E any](
	h *conform.Handle[S], c *prop.Case, actions Actions[S], stream Pull[E],
	mutate, remove string, written, removed func(args []any) E,
)

// AppendsArrive fails when a message that a call of append recorded does not
// arrive on the channel that subscribe returns.
func AppendsArrive[S, V any](
	h *conform.Handle[S], c *prop.Case, actions Actions[S], append string,
	subscribe Call[<-chan V], written func(args []any) V,
)

// DeliveriesMatchTheMode fails when a published event arrives another number
// of times than mode allows.
func DeliveriesMatchTheMode[S, V any](
	h *conform.Handle[S], c *prop.Case, actions Actions[S], publish string,
	subscribe Call[<-chan V], written func(args []any) V, mode string,
)

// RedeliveriesArrive fails when an event arrives fewer times than its
// publishes and its redeliveries together. written and redelivered return the
// event that a call of publish and of redeliver passes.
func RedeliveriesArrive[S, V any](
	h *conform.Handle[S], c *prop.Case, actions Actions[S], publish, redeliver string,
	subscribe Call[<-chan V], written, redelivered func(args []any) V,
)

// WatchesSeeTheTriggers fails when a subscription of watch misses an update
// that a trigger of its key caused after the watch, returns the updates out
// of order, or returns an update after its stop.
func WatchesSeeTheTriggers[S, H, V any](
	h *conform.Handle[S], c *prop.Case, actions Actions[S], watch, trigger string,
	next func(sub H, timeout time.Duration) (V, bool), stop func(sub H),
	retrigger Apply, watchKey, triggerKey func(args []any) any, triggered func(args []any) V,
)
```

`watchKey` and `triggerKey` project the key from the arguments of a call of `watch` and of `trigger`: the list of the arguments of the inputs that the two callables share by name, in the order of their names. `retrigger` calls `trigger` with the arguments of an earlier call of it. The generator writes `next` and `stop` as method expressions of the type of the subscription, such as `(*Ticker).Next`. The spec gains the input kinds `pull`, `channel` and `handle-method` for these inputs.

The iteration of a stream of pairs yields the first value of each pair, the key that a delete takes. The iteration of a stream of values and errors ends on the first error of the stream.

### The planted defects

| Form | Changes | Does | Minimum steps |
|---|---|---|---|
| `leaks-one-unit` | partner:observe | Returns one unit more than on the call before it, on every call | 2 |
| `half-applied-after-three` | partner:read | Returns its second sample once a dependency failed, from its fourth call on | 4 |
| `iterates-a-copy` | callable | Returns a stream over a copy of its source, on every call | 1 |
| `drops-appends-after-three` | callable | Returns no error and records nothing from its fourth call on | 4 |
| `breaks-the-mode-after-three` | callable | From its fourth call on, publishes nothing under `at-least-once` and `exactly-once`, and publishes each event twice under `at-most-once` | 4 |
| `skips-a-redelivery-after-three` | callable | Returns no error and does nothing, from its fourth call on | 4 |
| `misses-updates-after-three` | partner:trigger | Returns no error and does nothing from its fourth call on | 4 |

A form whose claim entails a claim of a single caller changes from its fourth call on, so the shallower claim passes on it. The forms of deliveries change from the fourth call on too, so the checks of a single caller of the callable pass on them. `iterates-a-copy` changes every call, because an iteration over a copy is a correct drain, and the claims of a single caller of the stream pass on it. `leaks-one-unit` changes the reader of the total. Each declaration moves its quantity in a way of its own, so no generic defect breaks the move itself. A reader that reports another total on each call breaks the law at the first move.

## Alternatives considered

### A reader of the total that its type selects

The derivation takes as the reader of the total the one callable without inputs whose result has the type of the field.

**Why not:** a declaration can have several such callables, such as a total and a count. A check of the wrong one passes or fails for a reason that the law does not state.

### A wait of real time for deliveries

The probe waits a fixed time after the last step and then takes the deliveries.

**Why not:** a short wait misses deliveries on a loaded machine, and a long wait slows every case. `synctest.Wait` ends the wait when the subject has nothing left to deliver.

### A probe of delivery that reads the history alone

Each delivery is an action of its own that reads the channel once.

**Why not:** a read of an empty channel blocks the step, and a read that does not block misses a delivery that is still under way. The law is a statement about what arrives at all, so the probe reads after the subject has finished.

## Drawbacks

- **A param more on a declaration.** A declaration with `conservative` states `observe` before it gets the check.
- **A bubble for deliveries.** A subject that delivers through a goroutine outside the bubble, or through a network, does not block inside the bubble, and its probe reports what arrived before the wait.
- **Subscriptions of the actions.** The actions of a case call each callable of the declaration, the role `subscribe` among them, so a sequence opens subscriptions that no step reads. A subject that blocks a delivery on a full channel of such a subscription blocks the sequence, and the bubble fails the case with a deadlock.
- **A defect of the reader.** `leaks-one-unit` proves the check against a total that changes, not against a move that loses a unit.
- **Seven probes per runtime.** Each probe exists in the probes of each runtime, with the input kinds `pull`, `channel` and `handle-method`.

## Open questions

None.

## References

| What | Where |
|---|---|
| The package testing/synctest | <https://pkg.go.dev/testing/synctest> |
| The function iter.Pull | <https://pkg.go.dev/iter#Pull> |
| The mixins and contracts of the shape catalog | <https://github.com/dokimasia/eidos/tree/main/eidos-plugin-shape/spec> |
| RFC-0020, the generated sequences of a declaration | [0020-sequences-against-a-reference.md](0020-sequences-against-a-reference.md) |
