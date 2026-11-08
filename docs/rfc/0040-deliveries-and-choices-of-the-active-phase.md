---
rfc: 0040
title: The deliveries and the choices of the active phase
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-08
updated: 2026-11-08
discussion: https://github.com/dokimasia/dokimi/issues/76
supersedes: RFC-0025 in part, RFC-0039 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0040: The deliveries and the choices of the active phase

## Summary

We propose two changes to the run of a simulation of sans-IO machines:

- Every case keeps the actions `deliver` and `advance`. A second delivery, `deliver fast`, has the weight of the other actions of the active phase together, and the swarm of a case keeps it with odds of 1 in 2. A case that keeps it delivers in more than half of its steps. A case without it delivers at the weight 1, as a case of sim v0.9.0 delivered.
- `Test`, `Falsify` and `Fuzz` give a case 82 choices for each step of `MaxSteps`. The most steps of a case then bound its active phase, and the most choices of the property engine do not.

## Motivation

A probe ran 200 cases of the simulation of sharder at each mean of 30, 300 and 3,000 steps with sim v0.9.0, with the most steps at ten times the mean:

| Mean of the steps | Deliveries of the active phase, p50 / p90 / max | Cases with a serving cell in the active phase |
|---|---|---|
| 30 | 0 / 4 / 14 | 0 of 200 |
| 300 | 1 / 42 / 188 | 0 of 200 |
| 3,000 | 0 / 124 / 352 | 0 of 200 |

A cell of sharder serves once its cell machine has replied to `Started` without an error. No action and no fault of these cases struck a cell that served. The falsification of sharder missed the planted defects of the admission of a disruption and of the generation of a commit, whose invariants need a cell that serves.

The run causes it through its swarm, its weights and its choices:

- **The swarm.** `stateful.Steps` makes one swarm choice per action, and keeps each action with odds of 1 in 2 (`stateful/steps.go:161-183` of assert-go v0.19.0). `Run` passes `deliver` and `advance` to the swarm as it passes every other action (`sim/run.go:551-566` of dokimi-go at sim v0.9.0). Half of the cases deliver no event in their active phase.
- **The weights.** Each action has the weight 1. sharder has 14 actions of its own and tolerates 13 kinds of faults, so its active phase has 34 actions. A case that keeps `deliver` delivers in about one of the 17 steps that the swarm leaves to it.
- **The choices.** The property engine rejects a case that makes more than 8,192 choices (`prop/option.go:120-131` of assert-go v0.19.0). A step of sharder makes about 4 choices, so no valid case takes more than about 2,100 steps, whatever `MaxSteps` allows. At a mean of 3,000 the run of one planted defect rejected 148 cases against 124 valid cases.

A race that needs a call between the send of a message and its delivery needs cases whose deliveries are slow. The fixture `sim/replicated` plants a primary that acknowledges a write before the backup has it, and its spec reads the backup between the acknowledgement and the replication. A fast rate of delivery in every case then has a cost of its own.

## Detailed design

### Components

| Component | Where | Change |
|---|---|---|
| The field `Keep` of an action | `go.dokimi.dev/assert/stateful` | The field `action.keep` of definition 16.2.0 of assert-spec. Every case keeps an action with `Keep` set, without a swarm choice |
| The deliveries of the active phase | `go.dokimi.dev/sim` | `deliver` and `advance` set `Keep`. The new action `deliver fast` delivers as `deliver` does, at the weight of the other actions |
| The choices of a case | `go.dokimi.dev/sim` | `Test`, `Prove` and `Fuzz` pass `prop.MaxChoices` with 82 choices for each step of `MaxSteps`, and at least the default of the engine |

### The deliveries of the active phase

`Run` builds the actions of the active phase of each case as sim v0.9.0 builds them, with `deliver fast` after `deliver` and two fields set:

| Action | `Keep` | `Weight` |
|---|---|---|
| `deliver` | true | 1 |
| `deliver fast` | false | The count of the actions of the active phase other than `deliver` and `deliver fast` |
| `advance` | true | 1 |
| `cancel`, each fault, the end of each fault and each action of the consumer | false | 1 |

`deliver` and `deliver fast` are enabled while a process has a pending event. Each picks the process through the strategy, and each passes its name to the case as a label. A step lists the kept actions that are enabled, and its index chooses among them by weight. With n other actions, a step of a case that keeps `deliver fast` takes a delivery with probability (n + 1) / (2n + 1) where every action is listed. A step whose list contains fewer of the other actions takes a delivery more often. A case without `deliver fast` takes `deliver` at the weight 1.

The two rates are a feature that the swarm toggles, as swarm testing toggles a feature of a test. A system that needs a few hundred deliveries to serve serves in the half of the cases that deliver fast. A race that needs a message in flight keeps the half of the cases that deliver at the weight 1. A shrink moves the index of a step towards the first listed action, `deliver`, so a minimal counterexample takes `deliver fast` only where its failure needs the rate.

The planted defect of `sim/replicated` with a read of the backup, run at 40 seeds with 100 cases each:

| The deliveries | Seeds that fail | Mean cases to the failure |
|---|---|---|
| sim v0.9.0 | 40 of 40 | 65.5 |
| `deliver` at the weight of the other actions in every case | 39 of 40 | 79.9 |
| `deliver` at the weight 1 in every case, and `deliver fast` in half of them | 40 of 40 | 38.5 |

`advance` keeps the weight 1. While a timer is pending, each advance moves the time to the next timer and further by an amount up to that gap, so each advance ends a timer.

### The choices of a case

```go
// maxChoices returns the most choices of a case of s: the choices of the
// property engine, or stepChoices for each step of MaxSteps where those are
// more.
func (s Simulation[E, F, S]) maxChoices() int {
	return max(engineChoices, stepChoices*s.MaxSteps)
}
```

`engineChoices` is 8,192, the default of the property engine. `stepChoices` is 82, which keeps the ratio of those 8,192 choices to the 100 steps of the package `stateful`, rounded up. `Test`, `Prove` and `Fuzz` pass `prop.MaxChoices(s.maxChoices())` beside their other options, and `Falsify` proves each check through `Prove`. A simulation with a `MaxSteps` of 100 or less keeps the 8,192 choices of the engine. A simulation with a `MaxSteps` of 1,000 gives a case 82,000 choices, and a step of sharder then has about 20 times the choices that it makes.

The doc comment of `Run` states the ratio. A consumer that calls `Run` inside a property of its own passes `prop.MaxChoices` itself.

### Failure handling

| Failure | Behaviour |
|---|---|
| A system without a pending event | Neither delivery is enabled, and the step chooses among the other kept actions |
| A case whose swarm disables every fault and every action of the consumer | The case delivers and advances alone, as a run without faults |
| A case that makes more than 82 choices for each step of `MaxSteps` | The engine rejects the case, as it rejects every case beyond its most choices |
| A trace of sim v0.9.0 | It replays. Its step entries name `deliver`, which every case keeps, and the swarm disables `deliver fast`, which no entry names |
| A replay token of sim v0.9.0 | It no longer replays the same case, because the swarm makes no choice for `deliver` and `advance` and the steps list `deliver fast` |

### Migration

- Every case of every simulation takes other choices. Its changesets start with `**Breaking:**`.
- A counterexample and the labels of a case name `deliver fast` for each delivery that the second action takes.
- sim requires the release of assert-go that implements definition 16.2.0.
- The sim generator does not change. A generated simulation calls `Test`, `Falsify` and `Fuzz`, which pass the choices of its `MaxSteps`.

### Verdicts

| Item | Verdict | Reason | Reverses the verdict |
|---|---|---|---|
| Every case keeps `deliver` and `advance` | in | The delivery and the advance of time are the progress of the system, not features of the test that the swarm omits | A system whose bugs need a case without deliveries in its active phase |
| `deliver fast` at the weight of the other actions, in half of the cases | in | Half of the cases then deliver at a rate that does not fall with the number of actions and faults, and the other half keep the rate at which a message waits in flight. The planted defect of `sim/replicated` failed at 40 of 40 seeds | A measured system that needs the fast rate in more than half of its cases |
| `advance` keeps the weight 1 | in | Each advance ends a timer, and a deployment delivers many messages between the ends of two timers | A measured system that needs a timer between most deliveries |
| 82 choices for each step of `MaxSteps` | in | The ratio of the defaults of the engine and of stateful bounds a long case by its steps | A consumer whose steps make more than 82 choices |
| A param for the rates of delivery | out | One rule fits each measured system, so the generator needs no param for it | A consumer that needs another share |

## Alternatives considered

### A swarm inside the run

`Run` could turn the swarm of `stateful` off and choose the faults and the actions of each case itself. The run would copy the swarm and its conditioning. A trace would not keep the actions that its step entries name, because the choices of the run are no swarm choices of the machine.

### One weight of `deliver` in every case

`deliver` could take the weight of the other actions in every case. The planted defect of `sim/replicated` then failed at 39 of 40 seeds, after a mean of 79.9 cases, where sim v0.9.0 needed 65.5. A message waits in flight for fewer steps, so the read of the backup falls between the acknowledgement and the replication in fewer cases.

### A coin of each case under a label

A case could draw under a label whether `deliver` takes the weight of the other actions, at the same odds. The rates would be those of `deliver fast`. A trace then states the draw after the draws of its system, and the property engine refuses each trace of sim v0.9.0, as in `the draw labelled "fast deliveries" takes the step entry of "send"`.

### A delivery after every step

The run could deliver a number of events that the case chooses after each action of the active phase. The index of a step would no longer choose between a delivery and a fault, so a fault could not strike between two deliveries of one step. The phases of RFC-0025 would change for every simulation.

### `deliver` and `advance` as drain actions

The definition calls a pending message and a due timer drain actions. The drain then takes them after the active phase until neither is enabled or `MaxSteps` steps ran, under the faults in force, before the recovery. `advance` is always enabled in a world with components, so each such case would take `MaxSteps` more steps.

### A weight alone

A weight changes the share of `deliver` among the kept actions. Half of the cases would still disable it.

### A field for the most choices

`Simulation` could take the most choices of a case. Each consumer would then derive the count from the choices of its own steps. The ratio of the defaults already gives a step of sharder about 20 times the choices that it makes.

### No limit on the choices

`Test` could pass the largest count to `prop.MaxChoices`. An input of the consumer that never stops drawing would then run until the timeout of the test, where the limit rejects its case.

## Drawbacks

- A case that keeps `deliver fast` delivers in more than half of its steps, so it checks the invariants and the history after more deliveries. A case of sharder at a mean of 300 steps then takes about 150 deliveries or more in its active phase, where a case of sim v0.9.0 took a median of 1. That count is arithmetic from the weights, not a measurement.
- In such a case each fault and each action of the consumer takes a smaller share of the steps. With n other actions, each takes at most one step in 2n + 1 where all are listed.
- A fault of a message behind another, such as `delay`, strikes in fewer cases, because every case delivers and the queues of the processes are short. In the spec of sim that runs every fault, `delay` struck in 4 of 800 cases at the default mean of 30 steps. The cases of sim v0.9.0 that disabled `deliver` had left their messages in their queues. The required share of 1% held at a mean of 100 steps.
- Two actions deliver an event, and a reader of a counterexample meets both names.
- The replay tokens of every simulation stop replaying the same case.
- sim then requires a newer release of assert-go, and each consumer moves to it.

## Open questions

None.

## References

| What | Where |
|---|---|
| Issue 76, the active phase of a case of sharder ends before its cells serve | <https://github.com/dokimasia/dokimi/issues/76> |
| RFC-0059 of assert-spec, an action that every case keeps | <https://github.com/dokimasia/assert-spec/blob/main/docs/rfc/0059-an-action-that-every-case-keeps.md> |
| RFC-0025, a deterministic simulation of sans-IO machines | 0025-a-simulation-of-sans-io-machines.md |
| RFC-0039, hosts that exit, the steps of a case and performers of a pointer | 0039-hosts-that-exit-the-steps-of-a-case-and-performers-of-a-pointer.md |
| Groce, Zhang, Eide, Chen and Regehr, "Swarm Testing", ISSTA 2012 | <https://agroce.github.io/issta12.pdf> |
