---
rfc: 0038
title: The end of a wait
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-08
updated: 2026-11-08
discussion: https://github.com/dokimasia/dokimi/issues/69
supersedes: RFC-0037 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0038: The end of a wait

## Summary

We propose three changes to the waits of a simulation of sans-IO machines:

- An effect can end the wait of an earlier effect of its machine, as an unwatch ends a watch. `perform ends=` marks such an effect. The anatomy returns no token for it, so the rule `fresh-tokens` does not judge it, and the world cancels the wait of its token.
- `perform end=` declares the event that ends the wait of an effect, as `TermEnded` ends a campaign with a granted term. An outage of the dependency ends a wait that has delivered an event with that event, and any other wait with the event of its failure.
- `Outcome.Cancel` ends what an effect left in its dependency, as the cancel of the context of a call ends it in a runtime. The world calls it where it drops or ends a wait: at a loss after the effect, at an outage, at an effect that ends the wait, and at the stop of the process. A crash does not call it.

## Motivation

sharder starts three kinds of waits, and a world of sim v0.6.0 ends none of them as the runtime of sharder ends them:

- **An effect that ends a wait.** `Unwatch` ends a `Watch` or a `WatchViews`, `Resign` ends a `Campaign`, and `Leave` ends a `Join`. Each of these effects has the token of the effect that started the wait, and does not start a call or a wait of its own (`internal/machine/store.go`, `election.go` and `membership.go` of sharder). The anatomy returns the token of an effect that starts a call or a wait, and false for any other effect, such as an unwatch. No directive tells the generator that an effect ends a wait. So the generator returns the token of every effect with one field of the token type (`generator/sim/plugin.go`, `anatomy`), and the rule `fresh-tokens` reports each unwatch, resign and leave as a reuse of the token of its wait. The generated `TestSimulation` of a leader that resigns fails at its first resign.
- **The end of a granted term.** An outage of a dependency ends each wait of an effect of the dependency with the event of the failure of the effect (`sim/fault.go`, `Outage`). A campaign of sharder delivers `Campaigned` with the term, and its runtime passes `TermEnded` once the term ends. The root of sharder takes a `Campaigned` with an error as the end of a campaign without a term, and passes the event `lost` to its lifecycle. The lifecycle refuses that event in the state `leading`, so the root stops itself. A join delivers `Joined`, and its runtime passes `SessionEnded` once the session ends. The agent takes a `Joined` with an error as a failed join and joins again, so its node has two sessions.
- **What a lost wait leaves.** A loss after the effect calls the function of the effect and drops its outcome (`sim/fault.go`, `faulted`). The candidacy or the session that the function registered remains in the store, and nothing ends it. The runtime of sharder cancels the context of a call whose result it lost, and the store then ends the candidacy or the session. In a world, the election can then grant a term to a machine that has taken its campaign as failed.

The client of etcd separates these ends. A watch ends when its caller cancels the context of the watch: "If the context "ctx" is canceled or timed out, returned "WatchChan" is closed", and the context "MUST be canceled, as soon as watcher is no longer being used, to release the associated resources". The end of a term is a signal of its own, apart from the result of `Campaign`: `Session.Done` "returns a channel that closes when the lease is orphaned, expires, or is otherwise no longer being refreshed".

## Detailed design

### Components

| Component | Where | Change |
|---|---|---|
| An effect that ends a wait | `go.dokimi.dev/sim` | `Anatomy.Ends`, and the cancel of the wait of its token |
| The end event of a wait | `go.dokimi.dev/sim` | `Anatomy.End`, which an outage takes for a wait that has delivered an event |
| The cancel of an outcome | `go.dokimi.dev/sim` | `Outcome.Cancel`, which the world calls where it drops or ends a wait |
| The sim generator | dokimi | The params `ends` and `end` of `perform`, and the fields `Ends` and `End` of the anatomy |
| The corpus | dokimi | The fixture `sim/waits` |

### The outcome and the anatomy

```go
// Outcome is the result of performing one effect.
type Outcome[E, F any] struct {
	// Event, Next, Message, Launch and Stop are the fields of an outcome of
	// sim v0.6.0.
	Event   E
	Next    func() (ev E, ok, done bool)
	Message *Message[E]
	Launch  *Launch[E, F]
	Stop    string
	// Cancel ends what the effect left in its dependency, such as the
	// candidacy of a campaign or the session of a join, as a runtime ends it
	// when it cancels the context of the call of the effect. It is nil for an
	// effect without such an end. The world calls it at most once.
	Cancel func()
}

// Anatomy is how a world reads the events and the effects of its machines.
type Anatomy[E, F any] struct {
	// Token returns the token of an effect that starts a call or a wait, and
	// false for any other effect, such as an unwatch or a reply.
	Token func(f F) (token uint64, ok bool)
	// Ends returns the token of the wait that the effect f ends, such as the
	// token of the watch that an unwatch ends, and false for any other effect.
	// It is nil for machines without such an effect.
	Ends func(f F) (token uint64, ok bool)
	// Fail returns the event of the result of the effect f that failed with
	// err, and false for an effect without a result.
	Fail func(f F, err error) (ev E, ok bool)
	// End returns the event that ends the wait of the effect f after the wait
	// has delivered an event, such as the end of a term, with err in its field
	// of the type error where the event has one. It reports false for an
	// effect without such an event, and it is nil for machines without one.
	End func(f F, err error) (ev E, ok bool)
	// The other fields are the fields of an anatomy of sim v0.6.0.
}
```

### An effect that ends a wait

```go
// Resign ends the campaign Token, and resigns its term in the store.
//
//dokimi:perform ends=[Campaign] func=resign
type Resign struct{ Token Token }

// Unwatch ends the watch Token.
//
//dokimi:perform ends=[Watch, WatchViews]
type Unwatch struct{ Token Token }
```

- **The anatomy.** `Ends` returns the token of an effect under `ends`, and `Token` returns false for it. The rule `fresh-tokens` and the planted defect `reuses-a-token` read `Token`, so they judge the effects that start a call or a wait alone.
- **The world.** The world performs the effect as it performs any other effect, through its function where it has one. It then cancels the wait whose effect has the token, among the waits of the process in its epoch: it calls the `Cancel` of the outcome of the wait. The wait remains until its `Next` reports that it has ended, so the double of the dependency returns the last events of the wait. A watch ends without an event, as a watch of etcd closes its channel. A campaign with a granted term ends with `TermEnded`, and a campaign without one with the event of its failure.
- **An unknown token.** An effect does not end a wait when no wait of its process has its token, such as the unwatch of a watch that has ended.
- **The function.** `func` is optional under `ends`. The function of a resign of sharder calls `store.Leadership.Resign` before the cancel. The cancel alone ends a watch, so an unwatch has no function.

### The end event of a wait

```go
// Campaign campaigns for the term of the scope through the store.
//
//dokimi:perform func=campaign dependency=Store result=Campaigned end=TermEnded
type Campaign struct {
	Scope string
	Token Token
}

// TermEnded is the end of the term that the election granted to the campaign
// Token.
type TermEnded struct{ Token Token }
```

- **The param.** `end` names the event that ends the wait of the effect after the wait has delivered an event. Its field of the token type takes the token of the effect. Its field of the type error takes the error of the end where the event has one.
- **The outage.** The world records whether a wait has delivered an event: the event of its outcome, or an event of its `Next`. An outage of the dependency ends a wait that has delivered an event with the event of `end`. It ends any other wait of the dependency with the event of its failure, as sim v0.6.0 does.
- **A watch.** A watch has no end event, so an outage ends it with the event of its failure. The error of that event ends a watch of sharder at any point.

### The cancel of an outcome

The world calls the `Cancel` of an outcome where a runtime cancels the context of the call of the effect:

| Where | Why |
|---|---|
| A loss after the effect | The world drops the outcome, and a runtime cancels the context of a call whose result it lost |
| An outage of the dependency of a wait | The world ends the wait |
| An effect that ends the wait | The machine ends the wait |
| The stop of the process of a wait | A runtime cancels the calls of a process that stops |

- **A crash.** A crash of the host does not call `Cancel`. The dependencies keep what they applied, such as the session of a crashed node until its lease expires.
- **Once.** The world calls the `Cancel` of an outcome at most once, also where an outage ends a wait that an effect of its machine has cancelled. A wait that has ended takes no cancel.
- **The end of a run.** `World.Stop` stops each process, so it cancels each wait. The run checks the invariants once more after the stop, and an invariant then sees the dependencies without what the cancels ended.

### The generator

| Directive | Change |
|---|---|
| `//dokimi:perform ends=[<effect>, ...] func=<function>` | `ends` names the effects whose waits the effect ends. `func` is optional, and the performer returns an outcome without a result for an effect without a function |
| `//dokimi:perform func=<function> dependency=<interface> result=<event> end=<event>` | `end` names the event that ends the wait after the wait has delivered an event |

| Directive of an effect | Finding |
|---|---|
| `ends` with `call`, `dependency`, `result`, `timeout` or `end` | `BadPerform`: an effect that ends a wait starts no call and has no result |
| `ends` on an effect without one field of a token | `BadPerform` |
| `ends` that names a type other than an effect of the machine with one field of the token type of the effect | `BadPerform` |
| `end` without `dependency` | `BadPerform`: only an outage of the dependency ends a wait with its end event |
| `end` that names an event without a field of the token type | `BadPerform`, as for `result` |
| `perform` without `func`, `call` and `ends` | `BadPerform`, as in sim v0.6.0 |

The template writes `Ends` and `End` only for a machine with such an effect, so the generated anatomy of every other machine does not change.

### Failure handling

| Failure | Behaviour |
|---|---|
| An effect that ends a wait, whose token names no wait of its process | The effect ends nothing |
| An outage of a wait that has delivered no event | The event of the failure of the effect, as in sim v0.6.0 |
| A wait whose double does not end it after its cancel | The world polls it until the stop of its process |
| A cancel of a wait that the world has cancelled | The world does not call `Cancel` again |

### Migration

- `Outcome.Cancel`, `Anatomy.Ends` and `Anatomy.End` are new fields. A nil field does not change the behaviour of the world, so an outcome and an anatomy of sim v0.6.0 keep their behaviour, and sim v0.7.0 is a minor release.
- A consumer marks each effect that ends a wait with `ends`, and names the end event of a wait with `end`. The function of an effect that leaves a registration in its dependency returns a `Cancel`.
- The corpus gains the fixture `sim/waits`. Its leader campaigns, watches and resigns. The double of its store ends a campaign and a watch at their cancel, and its system tolerates an outage and a loss after the effect.

### Verdicts

| Item | Verdict | Reason | Reverses the verdict |
|---|---|---|---|
| A marker of an effect that ends a wait | in | `fresh-tokens` cannot tell an unwatch from a watch without it | none known |
| The cancel of a wait at an effect that ends it | in | A runtime cancels the context of the wait, and every performer of sharder would otherwise keep its waits in the state of its double | none known |
| The double decides the last events of a cancelled wait | in | The last event depends on the state of the dependency, such as a granted term | a consumer whose double cannot tell how a wait ends |
| An end event for an outage after the first event | in | A machine of sharder takes the event of a failure after the first event as a failed start | none known |
| `Cancel` at a crash | out | The dependencies keep what a crashed process applied | none known |
| `Cancel` at the stop of a process | in | A runtime cancels the calls of a process that stops, and a stopped machine would otherwise keep its candidacy | a runtime that keeps the registrations of a stopped process |

## Alternatives considered

### The double ends the wait through its own state

The function of an unwatch would mark the watch in the double, and the `Next` of the watch would report its end. sharder can write this form on sim v0.6.0. Every performer of sharder would then keep the waits of its effects, and `fresh-tokens` would still report each unwatch.

### The world drops the wait and builds its last event

The world would end the wait at the unwatch, and queue the end event of `end` or the event of the failure. The last event of a resign of sharder depends on the term: `TermEnded` for a granted term, and a `Campaigned` with an error for a campaign without one. A watch ends without an event. Only the double has the state of the term, and the last event depends on that state.

### A directive of its own for an effect that ends a wait

`//dokimi:ends` would mark the effect beside `perform`. Every effect of a machine has a directive `perform`, which states how the world performs it. The world performs an effect that ends a wait in a way of its own, so the marker is a param of `perform`.

### A method of the double that cancels

The world would call a method `Cancel(token)` of the dependency. The function under `perform func=` builds its outcome from any state of its double, so only the function knows what it registered. A field of the outcome keeps the cancel beside the `Next` that the same function returns.

## Drawbacks

- The anatomy has two more fields, and `perform` has two more params. sharder marks 3 effects with `ends` and 2 effects with `end`.
- A double that ignores the cancel of a wait keeps the wait open until the stop of its process, and the world polls the wait after each step.
- The world cancels each wait at the end of each run, so a double whose cancel changes its state shows that change to the invariants that the run checks after the stop.

## Open questions

None.

## References

| What | Where |
|---|---|
| Issue 69, an outage and a loss after the effect end a wait without the end that a runtime gives it | <https://github.com/dokimasia/dokimi/issues/69> |
| Issue 70, fresh-tokens reports each unwatch, resign and leave as a reuse of the token of its wait | <https://github.com/dokimasia/dokimi/issues/70> |
| RFC-0025, a deterministic simulation of sans-IO machines | 0025-a-simulation-of-sans-io-machines.md |
| RFC-0037, hosts, started machines and the inputs of a system | 0037-hosts-started-machines-and-the-inputs-of-a-system.md |
| The watcher of the client of etcd | <https://github.com/etcd-io/etcd/blob/main/client/v3/watch.go> |
| The concurrency package of the client of etcd | <https://pkg.go.dev/go.etcd.io/etcd/client/v3/concurrency> |
