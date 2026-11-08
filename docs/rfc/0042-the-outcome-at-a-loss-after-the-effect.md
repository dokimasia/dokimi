---
rfc: 0042
title: The outcome of an effect at a loss after the effect
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-08
updated: 2026-11-08
discussion: https://github.com/dokimasia/dokimi/issues/78
supersedes: RFC-0037 in part, RFC-0038 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0042: The outcome of an effect at a loss after the effect

## Summary

We propose that a loss after an effect loses the result of the effect alone. The world applies the rest of the outcome as it applies it without the fault:

- It stops the process of `Stop` before it queues the event of the failure.
- It starts the machine of `Launch` with its call, and it passes the call of `Message` to its receiver. The reply to each call returns to no process.
- It calls `Cancel` in place of the wait of `Next`, and it queues the event of the failure in place of `Event`.

## Motivation

The fault `lose-after` models a call that the dependency applies and whose result the caller loses. The world strikes it at an effect of the dependency: it performs the effect, calls the `Cancel` of the outcome and queues the event of the failure (`sim/fault.go:639-648` of dokimi-go at sim v0.10.0). It drops every other part of the outcome.

- **A launch and a stop.** A machine that the dependency started does not run, and a machine that the dependency stopped keeps running.
- **A message.** The receiver never takes the call, so a loss after the effect of a message strikes as a loss before it.

The simulation of sharder fails a case this way:

- The agent of the root of sharder runs on the host of a controller process. Its `Handle` effect starts a cell machine through `Outcome.Launch` for an add and stops it through `Outcome.Stop` for a drop. The performer records each start and each stop in the state of the host.
- A loss after a drop made the performer record the stop and return the stop. The world dropped the stop, and the cell machine kept running.
- The next add of the cell, at a newer term, returned a launch. The world failed the case with `sim: the process "host-3/payments/cell/00VHND0GSPPJ00YR0EG7WH5WS2" runs, so no effect starts it`.

A deployment of sharder applies an add or a drop whose reply is lost. The host starts or stops the controller of the cell, and the agent of the root gets the error of the deadline of its call.

## Detailed design

### Components

| Component | Where | Change |
|---|---|---|
| The loss after an effect | `go.dokimi.dev/sim`, `World.faulted` | Applies the outcome of the effect with its result lost |
| The apply of an outcome | `go.dokimi.dev/sim`, `World.apply` | Takes whether a loss has taken the result of the effect |
| The doc comments of `Outcome`, `Launch`, `Message`, `FaultLoseAfter` and `World.LoseAfter` | `go.dokimi.dev/sim` | State the outcome at a loss after the effect |

### The outcome at a loss after the effect

| Field of the outcome | Without a fault | At a loss after the effect |
|---|---|---|
| `Stop` | The world stops the process before it queues `Event` | The world stops the process before it queues the event of the failure |
| `Event` | The process of the effect takes it | The process takes the event of the failure in its place |
| `Next` | The world starts the wait | The world starts no wait |
| `Cancel` | The world calls it where it ends the wait | The world calls it once, at the loss |
| `Message` | The receiver takes the call, and its reply is the result of the effect | The receiver takes the call, and its reply returns to no process |
| `Launch` | The machine starts with its call, and its reply is the result of the effect | The machine starts with its call, and its reply returns to no process |

```go
// apply performs the effect f of the process p: it stops the process of the
// stop of the outcome, queues the event of the result for a later step of p,
// starts the wait, sends the message and starts the machine of the launch. A
// wait whose outcome has an event has delivered that event. Where lost is set,
// a loss after the effect has taken the result: apply calls the cancel of the
// outcome in place of the event and the wait, and the replies to the message
// and to the launch return to no process.
func (w *World[E, F]) apply(p *Process[E, F], f F, lost bool)
```

`step` calls `apply` with `lost` unset for each effect that `faulted` leaves to it. `restore` does the same for each effect that a slowness kept in flight. `faulted` calls it with `lost` set at a loss after the effect, and then queues the event of the failure. A reply that returns to no process completes a call whose await already has its result:

- `resolve` ignores the reply.
- The timeout of the effect arms no timer for the call.
- A crash or a stop of the process of the effect does not cancel the call.

### Failure handling

| Failure | Behaviour |
|---|---|
| A loss after an effect whose launch names a process that runs | The world fails the case, as it fails it without the fault |
| A loss after an effect whose stop names a process that does not run | The world fails the case, as it fails it without the fault |
| A lost message across a partition | The call does not arrive, as without the fault |
| A lost message to a crashed or stopped process | The call fails, and the failure returns to no process |
| A crash of the process of the effect before the reply to a lost launch or message | The world does not cancel the call |

### Migration

- A simulation that strikes a loss after an effect with a stop, a launch or a message takes other steps after the loss. Its stored cases can stop replaying the same run.
- The corpus does not change. Its fixtures strike `lose-after` at the writes of `sim/placed` and at the waits of `sim/waits`, whose outcomes have a result and a wait alone.
- A machine that retries an effect after a lost launch can find the started machine running. Its performer then returns an outcome without a second launch.

### Verdicts

| Item | Verdict | Reason | Reverses the verdict |
|---|---|---|---|
| The stop and the launch of a lost effect | in | The dependency applied the effect, so the machine that it started runs and the machine that it stopped does not | none known |
| The message of a lost effect | in | The receiver applies the call of a message as a dependency applies a call. `lose-before` and `drop` model the loss of the call itself | A consumer whose dependency of a message fails before the receiver takes the call |
| A reply that returns to no process | in | The process of the effect has taken the event of the failure, and a second result would contradict it | none known |
| A cancel of the call at the loss | out | The dependency has applied the call. In a deployment of sharder the controller of a cell runs after the agent loses the reply of the add | A consumer whose started machines end at the cancel of the call that started them |
| The cancel of the outcome | in | A runtime cancels the context of a call whose result it lost, and the dependency ends what the effect left in it | none known |

## Alternatives considered

### The drop of the outcome

This is the form of sim v0.10.0. A loss after a drop of sharder leaves the machine of its cell running, and the next add of the cell fails the case at its launch.

### A stop and a launch without the message

The world could apply the stop and the launch alone, as issue 78 asks. A loss after the effect of a message would then still strike as a loss before it.

### The cancel of the call of a lost launch

A loss could also cancel the call of the machine of a lost launch, as the world cancels the open calls of a process that crashes. The machine would then take the end of its first call before its first step. A cell of sharder runs after a lost add in a deployment.

### A field of the outcome for the loss

`Outcome` could state what the world applies at a loss after the effect. Every performer would then state it for each effect, and the dependency of every effect applies its call in the same way.

## Drawbacks

- The steps after a loss change for every simulation that strikes `lose-after` at an effect with a stop, a launch or a message. Its stored cases can stop replaying the same run.
- A performer returns a launch for a machine that runs when its machine retries an effect after a lost launch, and the world then fails the case. sharder records each start in the state of the host, so its performer can return an outcome without the launch.
- `apply` takes a flag, and each of its three callers passes it.

## Open questions

None.

## References

| What | Where |
|---|---|
| Issue 78, a loss after an effect drops the launch and the stop of its outcome | <https://github.com/dokimasia/dokimi/issues/78> |
| RFC-0037, hosts, started machines and the inputs of a system | 0037-hosts-started-machines-and-the-inputs-of-a-system.md |
| RFC-0038, the end of a wait | 0038-the-end-of-a-wait.md |
