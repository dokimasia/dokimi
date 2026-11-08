---
rfc: 0039
title: Hosts that exit, the steps of a case and performers of a pointer
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-08
updated: 2026-11-08
discussion: https://github.com/dokimasia/dokimi/issues/72
supersedes: RFC-0025 in part, RFC-0037 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0039: Hosts that exit, the steps of a case and performers of a pointer

## Summary

We propose three changes to the simulation of sans-IO machines:

- The error of a machine ends its host. The world passes the stop event to each other process of the host, and the host exits. The action `restart` and `Recover` start an exited host again with each process that `Spawn` added, as a supervisor starts a process that has exited. The bounds of `tolerate` and `require` do not count an exit.
- A simulation sets the mean and the most steps of the active phase of a case. `Simulation.MeanSteps` and `Simulation.MaxSteps` pass them to `stateful.Steps`, and the params `mean-steps` and `max-steps` of the directives `sim` and `system` set them.
- A performer can take a pointer to its effect. The generated performer then passes the address of the effect.

## Motivation

- **The error of a machine.** A process of sharder ends every machine that it runs once one of them stops with an error, and a supervisor starts the process again (`controller/controller.go` and `embedded/process.go` of sharder). A world of sim v0.7.0 stops only the process of that machine (`World.take`). The other processes of its host keep running, and neither `Recover` nor the action `restart` starts the stopped process again, because it did not crash. A system reads the name, the epoch and the clock of a process, so it cannot end the host itself. In a probe of the simulation of sharder, a case recreated a namespace. The agents of the old incarnation stopped with the error of the incarnation, and the root of their host kept running. No run of such a case converges, and a violation that a run finds in that state is no defect of sharder.
- **The length of a case.** `Run` calls `stateful.Steps` without options, so the active phase of a case takes a mean of 30 steps and at most 100. A namespace of sharder needs more than 160 deliveries before its first cell machine starts. The actions and the faults of a case then strike a system that is still starting. In 300 cases the action `disrupt` ran in 36.6% of the cases, and no cell took a disruption. `Simulation` has no field for the steps of a case.
- **A performer of a large effect.** The form of a performer takes its effect by value, so gocritic's `hugeParam` reports each performer whose effect has 80 bytes or more. The lint baseline of ergon excludes `hugeParam` for `_test.go` files alone, and a simulation kit is no test file. 9 of the 21 performers of sharder take an effect of 80 to 128 bytes.

## Detailed design

### Components

| Component | Where | Change |
|---|---|---|
| The exit of a host | `go.dokimi.dev/sim` | The error of a machine ends its host, and `Restart`, `Recover` and the action `restart` start an exited host again |
| The steps of a case | `go.dokimi.dev/sim` | `Simulation.MeanSteps` and `Simulation.MaxSteps` |
| The sim generator | dokimi | The params `mean-steps` and `max-steps` of `sim` and `system`, and the performer of a pointer |

### The exit of a host

The world ends the host of a machine that reports an error after a step:

1. It stops the process of the machine as sim v0.7.0 does. The machine takes the stop event, and the world performs the replies of that step.
2. It stops each other process of the host that runs in the same way, in the order of their spawns. Each stop calls the `Cancel` of the waits of its process, and ends the open calls of the process as unknown with `ErrRunEnded`.
3. The host exits. Its pause and its strobe end, as at a crash.

`Restart` and `Recover` start an exited host as they start a crashed one. Each process of the host that the exit stopped and that `Spawn` added takes a new machine in a new epoch, and its start. A process that an effect started remains stopped until an effect starts it again, as after a crash. The action `restart` of a run chooses among the crashed and the exited hosts. A crash of an exited host crashes its media. Its restart then starts the processes that the exit stopped.

`Faults` lists an exited host as `exit:<host>` in the log of a run. `exit` is no kind of `Fault`. The params `tolerate` and `require` cannot name it, and the bound of `crash` counts only the crashes that the case chooses.

```go
// Restart restarts the crashed or exited host named name. Each process of the
// host that the crash crashed or the exit stopped, and that Spawn added,
// takes a new machine in a new epoch, and its start. A component process
// starts a new incarnation and opens its value. A process that an effect
// started remains down until an effect starts it again. Restart fails the
// case for a name that no host has, and for a host that has neither crashed
// nor exited.
func (w *World[E, F]) Restart(name string)
```

### The steps of a case

```go
type Simulation[E, F, S any] struct {
	// The fields of a simulation of sim v0.7.0.

	// MeanSteps is the mean number of steps of the active phase of a case,
	// and the mean of the package stateful, 30, for 0.
	MeanSteps int
	// MaxSteps is the most steps of the active phase of a case, and the most
	// of the package stateful, 100, for 0.
	MaxSteps int
}
```

`Run` passes `stateful.Mean(s.MeanSteps)` for a `MeanSteps` above 0, and `stateful.Max(s.MaxSteps)` for a `MaxSteps` above 0. The final phase keeps its bound of 10,000 steps.

| Directive | Params | Finding |
|---|---|---|
| `//dokimi:sim start=Started stop=Stopped mean-steps=200 max-steps=1000` | `mean-steps` and `max-steps`, both integers | `BadParam` for a value below 1, and for a `max-steps` below `mean-steps`. Without `max-steps`, `mean-steps` compares with the most of 100 |
| `//dokimi:system mean-steps=200 max-steps=1000` | The same params on a system of components | The same findings |

### The performer of a pointer

| Directive | Forms of the function |
|---|---|
| `perform func=` | `func(w *sim.World[E, F], p *sim.Process[E, F], f T) sim.Outcome[E, F]` and `func(w *sim.World[E, F], p *sim.Process[E, F], f *T) sim.Outcome[E, F]` |

The generated performer passes the address of the effect to a function of the second form. The field `Pointer` of the model of a performer marks such a function, and the template writes `&f` for it.

### Failure handling

| Failure | Behaviour |
|---|---|
| A machine stops with an error on a host with one process | The host exits, and `Recover` starts it again |
| A machine stops with an error while another process of its host has an open call | The other process takes its stop event and can reply. A call without a reply ends as unknown with `ErrRunEnded` |
| A crash action while a host has exited | The action chooses among the hosts with a process that runs, so it leaves the exited host |
| A crash point of a scenario after the step that ends a host | The crash crashes the media of the exited host, and the restart starts its processes |
| A machine that stops with an error after each start | Each start ends in an exit again. `Recover` starts the host once, and the final phase ends within its bound of steps |

### Migration

- The exit of a host changes every world whose machine stops with an error: the other processes of its host stop, and `Recover` starts the host again. Its changeset starts with `**Breaking:**`.
- `MeanSteps` and `MaxSteps` are new fields, and 0 keeps the steps of sim v0.7.0. A generated simulation without the params does not change.
- The pointer form is a second form. A performer of a value keeps its form.

### Verdicts

| Item | Verdict | Reason | Reverses the verdict |
|---|---|---|---|
| The error of a machine ends its host | in | A process of a deployment ends at the error of one of its parts, and the world otherwise explores states that no deployment reaches | A consumer whose processes on one host exit apart |
| `Restart`, `Recover` and the action `restart` start an exited host | in | A supervisor starts a process that has exited | none known |
| An exit apart from the bounds of `crash` | in | The case does not choose an exit | none known |
| A strategy of a supervisor for each host | out | sharder needs one strategy, and a host is already the domain of the crashes, the pauses and the clock of its processes | A consumer with processes on one host that exit apart |
| The steps of a case as params of the directive | in | The length of the start of a system is a property of the system | none known |
| A performer of a pointer | in | A pointer avoids the copy that `hugeParam` reports, and `http.Handler` takes its request the same way | none known |

## Alternatives considered

### A strategy of a supervisor for each host

An Erlang supervisor chooses among `one_for_one`, `one_for_all` and `rest_for_one`. A host could take such a strategy. sharder needs `one_for_all` alone, and a host is already the domain of the crashes, the pauses and the clock of its processes. A strategy for each host would have one value in use, the value that sharder needs.

### The system ends the host

The system could watch the errors of its machines and crash their host. It reads the name, the epoch and the clock of a process alone, so it cannot see the error of a machine. Its crash of the host would also lose the writes that a sync did not cover, and the exit of a process keeps those writes.

### A longer default of the steps

The package stateful fixes its defaults of 30 and 100 in its definition. A start of 160 deliveries is a property of sharder, and a longer default lengthens every case of every other simulation.

### A warm-up phase without faults

The world could deliver events without actions and faults until the system serves. The faults that strike a system while it starts are part of what a simulation explores, and the warm-up would leave them out.

### An exclusion of hugeParam

The baseline of ergon could exclude `hugeParam` for the packages of a simulation. The copy of an effect of 128 bytes at each performed effect is the cost that `hugeParam` reports, and a pointer removes it.

## Drawbacks

- A world cannot model two processes of a deployment that share a machine and exit apart. One host couples their exits, and two hosts separate their clocks and their crashes.
- The exit changes the outcome of the existing simulations whose machines stop with an error, so their consumers check their oracles again after the upgrade.
- A case with a mean of 300 steps runs ten times the steps of a case with the default mean of 30.
- The template writes two forms of the call of a performer.

## Open questions

None.

## References

| What | Where |
|---|---|
| Issue 72, a machine that stops with an error leaves the other machines of its host running | <https://github.com/dokimasia/dokimi/issues/72> |
| Issue 73, a case ends its actions and faults before a system of sharder serves | <https://github.com/dokimasia/dokimi/issues/73> |
| Issue 74, a performer of an effect of 80 bytes or more fails hugeParam | <https://github.com/dokimasia/dokimi/issues/74> |
| RFC-0025, a deterministic simulation of sans-IO machines | 0025-a-simulation-of-sans-io-machines.md |
| RFC-0037, hosts, started machines and the inputs of a system | 0037-hosts-started-machines-and-the-inputs-of-a-system.md |
| The supervisor behaviour of Erlang/OTP | <https://www.erlang.org/doc/system/sup_princ.html> |
| The pod lifecycle of Kubernetes | <https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/> |
