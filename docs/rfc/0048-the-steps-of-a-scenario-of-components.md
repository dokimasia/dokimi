---
rfc: 0048
title: The steps of a scenario of a world of components
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-10
updated: 2026-11-10
discussion: https://github.com/dokimasia/dokimi/issues/88
supersedes: RFC-0025 and RFC-0036 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0048: The steps of a scenario of a world of components

## Summary

We propose that each call of a component in a run of a scenario is a step of the world:

- In a run of a scenario, a call of `World.Invoke` on a component that runs counts as a step.
- The world settles after the call, and the run checks the invariants, as after a step of a machine.
- Each such step has two crash points. The first crashes the host of the component before the call, so the call ends as unknown with `medium.ErrCrashed`. The second crashes the host after the call. A restart of the host follows each crash.
- In the active phase, a call of `Invoke` stays part of the action that makes it. The run settles and checks once after the action.

## Motivation

A scenario of a world of components checks its invariants once, after `World.Stop`, and has no crash points:

- `World.Step` counts the steps of machines alone (`sim/world.go:426-459` of dokimi-go at sim v0.12.0). A scenario of components takes no such step, so its count of steps is 0, and `Test` does not run a crash point.
- The run of a scenario runs the scenario, stops the world, and then checks the invariants and the history (`sim/run.go:386-406`). An invariant that reads a component skips it after the stop, because the component no longer runs.
- `Invoke` starts the call on a goroutine of the incarnation and returns. The world does not wait for quiescence after it, so the call completes only at a later wait, such as the wait of `World.Advance`.

The scenarios of machines check their invariants after every step and have two crash points for each step. The simulator of a pool of theke runs as a world of components. Its invariants have branches for inputs that only some seeds produce, such as a write of deletes alone on a full disk. A scenario is the fixed run that checks such a branch in every run of `TestSimulation`.

## Detailed design

### Components

| Component | Where | Change |
|---|---|---|
| The call of a component | `go.dokimi.dev/sim`, `World.Invoke` | Is a step in a run of a scenario, with two crash points |
| The run of a scenario | `go.dokimi.dev/sim`, `Test` and `Simulation.Prove` | Runs each scenario as a run in which each call of a component is a step |

### The step of a call

In a run of a scenario, `Invoke` on a component that runs takes these steps in order:

1. It counts a step of the world.
2. At the crash point before the step, the host of the component crashes. The call ends as unknown with `medium.ErrCrashed`, and the host restarts. A crash of a machine before its step loses the event of the step in the same way.
3. Otherwise the call starts on a goroutine of the incarnation, as before.
4. The world settles. It wakes each wait of the components on `Env.Wait` and records the step for the watchdog of `conform.Bubble`. It then waits for quiescence, so a call whose goroutine does not block completes within its step.
5. At the crash point after the step, the host crashes and restarts.
6. The run checks the invariants, as after a step of a machine.

A call of a component that does not run is no step, as `World.Step` takes no step of a process that does not run. The history records such a call as unknown, as before.

`Test` runs each scenario once to count its steps. It then runs the scenario twice for each step, with a crash before the step and with a crash after it. The contract of such a run contains the number of the step, as `the scenario set finds no failure with a crash before step 1`. `Prove` runs the scenarios of an invariant the same way, without the crash points.

### The active phase

An action of the active phase settles the world and checks the invariants once after its run. A call of `Invoke` in an action is part of the step of the action, and the run checks the invariants after that step. A check after each call would run the invariants twice for each action that makes a call.

### Failure handling

| Failure | Behaviour |
|---|---|
| An invariant that a call of a scenario breaks | `Test` fails at the run without a crash, with the statement of the invariant |
| An invariant that a crash before or after a call breaks | `Test` fails at the crash point, as `the scenario <name> finds no failure with a crash after step <n>` |
| A scenario that calls a crashed component | The call ends as unknown with `medium.ErrCrashed`, and counts no step |

### Migration

- The consumer does not change a scenario of a world of components. Its calls gain the steps, the checks and the crash points.
- A scenario whose invariant fails after a call, or at one of the crash points, fails `TestSimulation`.
- Each scenario of n calls adds 2n runs to `TestSimulation`.

### Verdicts

| Item | Verdict | Reason | Reverses the verdict |
|---|---|---|---|
| A call of `Invoke` as the step of a scenario | in | A scenario of components drives its components through their calls, as a scenario of machines drives its machines through their events | A scenario that drives a component through another method of the world |
| A crash before the call that loses the call | in | A crash of a machine before its step loses the event of the step. The crash point after the step before it already restarts the host before the call | none known |
| A step for each call in the active phase | out | The action of the call already settles and checks, so each action of a call would check the invariants twice | none known |
| A step for each `World.Advance` of a scenario | out | An advance moves the clocks and takes no input, so its crash point repeats the crash point after the call before it | A scenario whose components act on a timer alone |

## Alternatives considered

### A method of the world for each step

A scenario could mark its steps through a method of the world for components. Each scenario of components would then state its steps by hand, and a scenario that misses one would lose its crash points.

### A check after each call in every run

`Invoke` could settle and check after every call, also in the active phase. The action that makes the call would then settle and check a second time.

## Drawbacks

- Each call of a scenario adds a wait for quiescence and two runs of the scenario for its crash points.
- A call that completed only at a later advance completes within its step. An invariant can observe the effect of the call before the advance.
- The world has one more field, which marks a run of a scenario.

## Open questions

None.

## References

| What | Where |
|---|---|
| Issue 88, a scenario of a world of components checks its invariants only after Stop | <https://github.com/dokimasia/dokimi/issues/88> |
| RFC-0025, a deterministic simulation of sans-IO machines | 0025-a-simulation-of-sans-io-machines.md |
| RFC-0036, components that run goroutines | 0036-components-that-run-goroutines.md |
