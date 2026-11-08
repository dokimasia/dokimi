---
rfc: 0041
title: The final phase of a long convergence and the timers of a performer
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-08
updated: 2026-11-08
discussion: https://github.com/dokimasia/dokimi/issues/77
supersedes: RFC-0039 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0041: The final phase of a long convergence and the timers of a performer

## Summary

We propose two changes to the run of a simulation of sans-IO machines:

- The final phase takes up to 1,000 steps for each second of the longest `Within` of the convergences of a simulation, and at least the 10,000 steps of sim v0.9.0.
- `World.AfterFunc` calls a function once the clock of a process has advanced by a duration, and then polls each wait. A performer can make the world poll a wait at a time without an event to the machine.

## Motivation

The final phase of a case of the simulation of sharder ends at its limit of 10,000 steps before the bound of the convergence of sharder passes. The case then fails with `sim: the system takes more than 10000 steps to settle`, and the run fails before it checks the convergence.

- `settleDown` delivers each pending event and advances the time until the bound of each convergence has passed, for at most `settleLimit` steps (`sim/run.go:632-656` of dokimi-go at sim v0.9.0). The limit is 10,000, and each delivery and each advance counts as one step.
- The convergence of sharder derives its bound from the specs of its presets: 9m40s of the run. A case with 5 hosts and 2 namespaces has 30 sessions, and each renews every 3.33 s on the clock of its host.
- The world polls the waits after each step of a machine alone. A performer of sharder arms a world timer with an event at each renewal, so the world polls the wait of the session at that time.
- In the edge case whose hosts all run at twice the rate of the run, the final phase took its 10,000 steps after about 7m50s of the run. 8,506 of the 9,557 deliveries of the case were the events of those timers, which every machine ignores.

## Detailed design

### Components

| Component | Where | Change |
|---|---|---|
| The limit of the final phase | `go.dokimi.dev/sim` | `settleLimit`, or `settleRate` steps for each second of the longest `Within` where those are more |
| `World.AfterFunc` | `go.dokimi.dev/sim` | A timer of a process that calls a function and then polls each wait |

### The limit of the final phase

```go
// settleLimit is the fewest steps of the limit of the final phase, and
// settleRate the steps of the limit for each second of the longest Within
// of the convergences of a simulation.
const (
	settleLimit = 10_000
	settleRate  = 1_000
)
```

`settleDown` takes the larger of `settleLimit` and `settleRate` times the whole seconds of the longest `Within`. The final phase of sharder has a bound of 9m40s and takes up to 580,000 steps. The measured case took about 20 steps for each second of its run, and the limit leaves it 50 times that rate. A bound of 10 s or less keeps the limit of sim v0.9.0. A system whose events at one time do not end fails within that limit.

### `World.AfterFunc`

```go
// AfterFunc calls f once the clock of the process p has advanced by d, and then
// polls each wait, as the world polls the waits after a step. A jump or a
// strobe of the clock does not move the timer. The world calls f only while p
// runs in the epoch in which AfterFunc armed the timer. AfterFunc fails the
// case for a component process, which takes no event.
func (w *World[E, F]) AfterFunc(p *Process[E, F], d time.Duration, f func())
```

The timer is a timer of the world, as the timer of `World.After` is. `World.Next` counts it. An advance of the final phase moves the time to it and runs it. A performer of a session arms the timer of its next renewal with a function that renews the session in its dependency. The wait of the session then takes the renewal at its poll, and the renewal costs one advance of the time and no delivery.

`host.Clock.AfterFunc` calls a function once the monotonic clock of a host has advanced by a duration, and a component waits on the clock of its host this way. `World.AfterFunc` gives a performer the same timer on the clock of a process.

### Failure handling

| Failure | Behaviour |
|---|---|
| A system whose events at one time do not end | The final phase fails the case at its limit, within 10,000 steps for a bound of 10 s or less |
| A timer of a process that crashed or restarted before the timer is due | The world calls no function, and polls no wait for the timer |
| `AfterFunc` for a component process | The case fails, as it fails for `After` |
| A function of a timer that panics | The advance that runs the timer panics, and the case fails with the panic |

### Migration

- The limit of the final phase grows for a `Within` of more than 10 s. A simulation that settles within 10,000 steps settles as before.
- `AfterFunc` is new. A performer that arms `After` with an event that every machine ignores can arm `AfterFunc` in its place.

### Verdicts

| Item | Verdict | Reason | Reverses the verdict |
|---|---|---|---|
| A limit of the final phase that follows the longest `Within` | in | The bound of a convergence is the time that the final phase needs, and a fixed count of steps ends a long bound early | none known |
| 1,000 steps for each second of the bound | in | The 10,000 steps of sim v0.9.0 cover 10 s at that rate, and the measured system takes about 20 steps for each second | A measured system that takes more than 1,000 steps for each second of its bound |
| `AfterFunc` with a function | in | A performer can change its dependency at the time of the timer, and `host.Clock.AfterFunc` has the same form | none known |
| A poll of every wait after every advance | out | It moves the events of every wait of every simulation, where `AfterFunc` changes the worlds that call it alone | A consumer whose waits all depend on the time |

## Alternatives considered

### A world timer with an event that every machine ignores

A performer of sim v0.9.0 arms such a timer for each renewal of a session. Each renewal then costs an advance and a delivery. 8,506 of the 9,557 deliveries of the measured case were such events.

### A poll of every wait after every advance

Every simulation with a wait would take the events of its waits at other steps if the world polled each wait after each advance. Its stored cases would then stop replaying.

### A limit of the deliveries between two advances

The final phase could fail a system that takes more than 10,000 deliveries without an advance of the time. It finds a system whose events at one time do not end. A system that arms a timer of 1 ns at each event would then advance without an end, and no limit would end its case.

### A field for the limit

`Simulation` could take the limit of the final phase as a field. Each consumer would then count the steps of its own convergence. The rate of 1,000 steps for each second already leaves the measured system 50 times its rate.

## Drawbacks

- A system that does not settle takes more steps before its case fails when its bound is long: 580,000 steps for the bound of sharder.
- `AfterFunc` runs its function inside an advance, so the function sees the clocks of the hosts after the advance and before the next step.
- A function of a timer that changes the world, such as one that delivers an event, takes effect at the advance and not at a step of the active phase.

## Open questions

None.

## References

| What | Where |
|---|---|
| Issue 77, the final phase ends at 10,000 steps before the bound of the convergence of sharder passes | <https://github.com/dokimasia/dokimi/issues/77> |
| RFC-0039, hosts that exit, the steps of a case and performers of a pointer | 0039-hosts-that-exit-the-steps-of-a-case-and-performers-of-a-pointer.md |
| `host.Clock.AfterFunc` of dokimi-go | <https://github.com/dokimasia/dokimi-go/blob/main/host/clock.go> |
| The function `time.AfterFunc` of Go | <https://pkg.go.dev/time#AfterFunc> |
