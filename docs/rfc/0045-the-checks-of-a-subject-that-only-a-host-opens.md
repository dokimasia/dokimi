---
rfc: 0045
title: The checks of a subject that only a host opens
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-09
updated: 2026-11-09
discussion: https://github.com/dokimasia/dokimi/issues/82
supersedes: RFC-0036 in part
superseded-by: ADR-0025 in part
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0045: The checks of a subject that only a host opens

## Summary

We propose that the option `Hosted` of conform builds the subject of every check of an application, as the option `Durable` builds it over a medium:

- Each check of an application with `Hosted` runs in `Bubble`, unless the check states an `Around`.
- The handle of a check starts a separate host from its case. `Handle.Subject`, `Handle.New` and `Handle.Fresh` open their subjects through the door Host, on the environment of the first incarnation of that host.
- A check without a property runs in one case of the property engine, so its handle has a case that starts the host. The check itself does not change.
- An application with the option `Tasks` keeps the factory of the consumer for its checks, because a subject of `Tasks` waits through the task scheduler that only the world of a check runs.

## Motivation

A subject that runs on a `host.Host` needs the case of a property, because `host.New` takes a `*prop.Case`. conform builds the handle of a check without a property on the seat of its test (`conform/check.go:385-388` of dokimi-go at conform v0.16.3). It passes that handle to the factory of the application. The factory then cannot start a host, so the suite of an interface cannot check an implementation that only a host opens.

The HLC of a simulated pool of sharder runs over the clock of its host, with `hlc.NewOver` of core. Its physical clock returns a `clock.Timer` that `host.Clock.AfterFunc` sends. Core generates the suite of `clock.Timer`, and all seven checks of the suite are checks without a property. Each of the seven checks and the three smoke checks fails in the factory, before a call of the timer. The derived check `Reset/deterministic#generated` runs as a property, and its factory starts the host and passes.

The world of a check opens each incarnation of a subject through the door Host (`conform/suite.go:263-272` at conform v0.16.3). Under the option `Durable`, an application opens each subject over a new `medium.Memory` through the door Open, in place of the factory of the consumer (`conform/suite.go:509-514`).

## Detailed design

### Components

| Component | Where | Change |
|---|---|---|
| The factory of the checks | `go.dokimi.dev/conform`, `config.subjects` | Under `Hosted` without `Tasks`, opens each subject through the door Host on the host of its handle |
| The host of a handle | `go.dokimi.dev/conform`, `Handle.incarnation` | Starts from the case of the handle, and ends with the case |
| The environment of a check | `go.dokimi.dev/conform`, `Check.environment` | `Bubble` for a check without `Around` of an application with `Hosted` and without `Tasks` |
| A check without a property | `go.dokimi.dev/conform`, `Check.exercise` | Runs in one case under `Hosted` without `Tasks`, the replay of the simplest case |

### The subjects

`Run` and `Check.Run` take the factory of their checks from `config.subjects`:

```go
// subjects returns the factory of the subjects of the checks of an
// application: under the option Hosted a factory that opens each subject
// through the door Host on the host of its handle, and the factory of factory
// otherwise.
func (c *config[S]) subjects(newSubject func(h *Handle[S]) S) func(h *Handle[S]) S

// hosts reports whether the handle of each check of an application opens its
// subjects on a host that the case of the handle starts: under the option
// Hosted, and without the option Tasks. A subject of Tasks waits through the
// task scheduler that the world of a check runs, and a check without a world
// does not run one.
func (c *config[S]) hosts() bool
```

A subject of `Tasks` blocks only through its host. It waits through `Env.Wait` until another task changes a condition. The task scheduler checks the condition of each waiting task after each turn and after each move of its clock. A host without a scheduler checks a wait again only at a `Yield`, at the return of a function of `Incarnation.Go` and at `Host.Broadcast`. The subject `Scheduled` of the corpus waits for its flusher that way. Its `Put` never returns on a host without a scheduler. Only the world of a check runs the host under a task scheduler of its case. An application with `Tasks` keeps the factory of the consumer for each check. It runs only the checks that state an `Around` inside an environment.

`Handle.Subject` and `Handle.New` build their subjects through the factory that `config.subjects` returns. `Handle.Fresh` builds its subject through the copy of that factory in the field `fresh` of the configuration. A handle starts its host when it builds its first subject. The handle of a check with `Apply` starts its host before the check runs:

- **The host.** `host.New` of the case of the handle, named `subject`, without a network and without a scheduler, so the goroutines of the host are goroutines of the bubble.
- **The incarnation.** The first incarnation of the host. Every subject of the handle opens on its environment. The lifetime of the incarnation derives from the context of the case, so the end of the case ends it. A wait on the lifetime and a sleep on the clock of the host then return, and the end of the lifetime cancels each pending call of `Clock.AfterFunc`.
- **The medium.** The handle registers the disk of the host as the medium of each subject that it opens. The doors Restart and Crash of `Durable` then open the subject again over that disk.

`Hosted` takes precedence over `Durable` for the factory of the checks, because the environment of an incarnation contains its disk. `Fuzz`, `RunBudgets`, `FalsifyAllocs` and `Bench` build their subjects through the factory of the consumer. A benchmark measures time on the clock of the process and runs outside a bubble.

### The runs of a check

| Check of an application with `Hosted` and without `Tasks` | Run |
|---|---|
| A check with `Around` | Inside its `Around`, as before |
| A check with `Property` and without `Around` | Each case in `Bubble` |
| A check with `Apply` and without `Around` | In `Bubble`, in one case of the property engine: the replay of the simplest case |

A check with `Apply` keeps its semantics in that case:

- Its handle starts its host before the check runs, so the case chooses the drift of the clock of the host. The property engine ends a run in which no case makes a choice as vacuous, and a check that does not open a subject of the door Host, such as a check over a medium of `Durable` or a check of the door Faults, would otherwise fail that way.
- Its handle takes its random values from the seed of the check, as the handle of a check with `Apply` takes them, and not from the case.
- Its record counts no cases.
- A failure of the check fails the case, and falsification, isolation and the edges read the failure of the case, as they read the failure of a property check.

Falsification runs each check of an application with `Hosted` and without `Tasks` in `Bubble` too, and so do the isolation of another planted defect and the edges of a check.

### Failure handling

| Failure | Behaviour |
|---|---|
| A handle without a case builds a subject of the door Host, such as the handle of `Check.Options` | The check ends with a fault that wraps `ErrConfig` |
| A subject that keeps a goroutine running | The end of the case ends the lifetime of the incarnation. A goroutine that waits on the lifetime or sleeps on the clock of the host then returns |
| A goroutine of the subject that waits for a mutex while the bubble cannot advance | The watchdog of `Bubble` reports the stall, as for a world |

### Migration

- An application with `Hosted` and without `Tasks` checks the subject of its door Host with every check. The factory of the consumer builds no subject of a check, as under `Durable`.
- An application with `Hosted` and `Tasks` checks its subjects as before.
- The fixture `iface/door/hosted` of the corpus applies its suite to the hosted form of `Flushed`. `Scheduled`, whose application states `Tasks`, keeps its factory.
- The test output of a failed check without a property states the failure of its one case.

### Verdicts

| Item | Verdict | Reason | Reverses the verdict |
|---|---|---|---|
| The door Host builds the subject of every check of `Hosted` | in | `Durable` replaces the factory the same way, and the checks then test the subject that a deployment runs | A consumer whose checks without a world need the factory beside `Hosted` |
| A case for a check without a property | in | `host.New` needs the case, and one replay of the simplest case runs the check once | A host of dokimi that starts without a case |
| `Bubble` for each check of `Hosted` | in | A host without a scheduler runs its goroutines and its timers in the bubble | none known |
| The door Host for the checks of an application with `Tasks` | out | A subject of `Tasks` waits through the task scheduler, and a check without a world runs no scheduler | A design that runs the body of a check as a task of a scheduler of its case |
| A new option for subjects that only a host opens | out | The option would state a second door that opens a subject over an incarnation | A consumer that needs both forms of one subject in one application |
| The hosts in `Fuzz`, `RunBudgets` and `Bench` | out | A fuzz target applies one `Around` to every input, and a benchmark measures outside a bubble | A consumer that needs a budget of a subject that only a host opens |

## Alternatives considered

### A nil factory

`Run` could open the subjects through the door Host when the factory of the consumer is nil. The generated entry points wrap the factory in a separate function. conform then never receives a nil factory unless the generator changes.

### A new option

conform could take an option that states that only a host opens a subject. An application would then state the door Host twice, through `Hosted` for its worlds and through the new option for its other checks.

### A case of the check in the generated suite

The generator could give each check without a property a function `Property` that takes no input from its case. Every generated suite would then change, and so would the records of its checks.

## Drawbacks

- Each check of an application with `Hosted` starts a host and runs in a bubble. Its checks take longer than over the factory of the consumer.
- The fixture `iface/door/hosted` no longer checks the stores that its factory builds without a host.
- A failed check without a property reports the failure of a case, with the replay token of the simplest case.

## Open questions

None.

## References

| What | Where |
|---|---|
| Issue 82, a check without a property cannot build a subject that only a host of dokimi opens | <https://github.com/dokimasia/dokimi/issues/82> |
| RFC-0036, a simulation of components that run goroutines | 0036-components-that-run-goroutines.md |
