---
rfc: 0037
title: Hosts, started machines and the inputs of a system
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-07
updated: 2026-11-07
discussion: https://github.com/dokimasia/dokimi/issues/63
supersedes: RFC-0025 in part
superseded-by: RFC-0038 and RFC-0039 in part
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0037: Hosts, started machines and the inputs of a system

## Summary

We propose five changes to the simulation of a system of sans-IO machines:

- A host runs one or more machine processes. They share its clock, and a crash, a restart, a pause,
  a jump and a strobe strike each of them at once. `tolerate` and `require` count hosts.
- A performer starts a machine on the host of its process, and stops such a machine. The first event
  of a started machine is an inbound call of the performer, whose reply is the result of the effect.
- `World.Now` returns the time of the run, which has no drift and which no fault of a host moves.
- `perform func=` names the dependency that its function calls, so the faults of the dependency
  strike the effect as they strike an effect under `call=`.
- The function under `system` can take its inputs from the case, as the function of an action
  does. A scenario takes the simplest choice of each input.

## Motivation

sharder runs the sans-IO machines of its components together in one process, and a simulation of
sim v0.5.0 cannot run them that way:

- A controller process of sharder runs the root machine of each namespace, the agent machine of each
  root and the cell machine of each cell that the root places on the process. A crash ends these
  machines at once, and they read one clock. `World.Spawn` gives each machine a process and a clock
  of its own, and `Crash`, `Pause` and `Restart` strike one process (`sim/world.go`). A world in which
  each machine of sharder is a process crashes the agent while its cells keep running, a state that
  the design of sharder excludes.
- The agent starts a cell machine at the epoch of an add, and stops it at a drop. A restart of the
  process starts the root and the agents again, and no cell. A world starts a machine only through
  `Spawn`, which queues `Started` with a reply token of the world, and `Recover` restarts every
  crashed process.
- The in-memory store of sharder keeps the sessions of every node on one clock. A dependency of a
  world reads the time only through the clock of one process, whose drift the case chooses and which
  a pause, a jump and a strobe change.
- The world strikes the faults of a dependency only at an effect under `call=`
  (`sim/fault.go`, `faulted`). None of the 21 effects of sharder can bind through `call=`: an effect
  that waits, that sends a message or whose result needs the logic of the consumer has a performer
  under `func=`. So the world strikes no fault of the store, the membership or the application of
  sharder.
- `System` is `func(w *World[E, F])`, so every case of a simulation runs one system. The case cannot
  choose the number of processes, the configuration of a process or the seed of a machine, and a
  failure that the size of the system causes does not shrink to a smaller system.

The simulator of FoundationDB runs one or more processes on each simulated machine and kills a
process or the whole machine. madsim kills a node with every task of the node. turmoil crashes a
host with the software that runs on it. The VOPR of TigerBeetle chooses the number of the replicas
of a cluster from the seed of the run, and Hypothesis runs each `@initialize` rule of a state machine
with arguments that the case draws.

## Detailed design

### Components

| Component | Where | Change |
|---|---|---|
| Hosts | `go.dokimi.dev/sim` | `Host`, `World.AddHost` and `Host.Spawn`, and the faults of a process on the host of the process |
| Started machines | `go.dokimi.dev/sim` | `Outcome[E, F]` with `Launch` and `Stop` |
| The time of the run | `go.dokimi.dev/sim` | `World.Now` |
| The waits of a dependency | `go.dokimi.dev/sim` | An outage of a dependency ends each wait of an effect of the dependency |
| The inputs of a system | `go.dokimi.dev/sim` | `System[E, F]` takes the case |
| The sim generator | dokimi | The param `dependency` of `perform`, the forms of `system` with inputs, and the outcome of two type arguments |
| The corpus | dokimi | The fixture `sim/placed` and the regenerated simulations |

### Hosts

```go
// Host is a machine of a world on which machine processes run. The processes
// of a host share its clock. A crash, a restart, a pause, a jump and a strobe
// strike each process of the host at once, and a fault of a medium strikes the
// media of each of them.
type Host[E, F any] struct{ /* unexported */ }

// AddHost adds a host named name to w and returns it. The case chooses the
// drift of its clock under the label drift.<name>. AddHost fails the case
// where w already has a host named name.
func (w *World[E, F]) AddHost(name string) *Host[E, F]

// Name returns the name of the host.
func (h *Host[E, F]) Name() string

// Spawn adds a process named name to the host and queues its start, as
// World.Spawn does. newMachine builds the machine of the process from the clock
// of the host. A restart of the host starts the process again in a new epoch.
// Spawn fails the case for a host that does not run, and for a name that a
// process of the world has.
func (h *Host[E, F]) Spawn(name string, newMachine func(clock *Clock) Machine[E, F], deps ...any)
```

- **Spawn.** `World.Spawn(name, newMachine, deps...)` adds a host named `name` with one process of
  that name. It draws the drift of the clock under the label `drift.<name>`, as it does now, so a
  system of one process per host keeps its choices.
- **Names.** A process and a host each have a name of their own. A fault of a host takes the name of
  the host, and `Deliver`, `Call`, `Step`, `Message.To` and the faults of a message take the name of
  a process. `World.Spawn` gives a process the name of its host.
- **The faults of a host.** `Crash`, `Restart`, `Pause`, `Resume`, `Jump`, `Strobe`, `Tear`,
  `Corrupt`, `Fill`, `RestoreDisk` and `ReplaceDisk` take the name of a host and strike each process
  of the host. A strobe counts the steps of every process of the host. `Faults` lists a fault of a
  host under the name of the host, such as `crash:controller`.
- **Partitions.** The sides of a partition are host names. Two processes of one host are never cut
  apart, and a message between them arrives under every partition.
- **The run.** The run chooses the target of a fault of a host among the hosts, so `tolerate` bounds
  the hosts that a fault keeps down at once, and `require` counts the cases in which the fault
  strikes a host. `drop`, `delay` and `duplicate` target a process with a message in flight, as now.
- **Components.** `SpawnComponent` adds a host with the one process of the component. A component
  does not share its host with a machine.

### Machines that an effect starts and stops

```go
// Outcome is the result of performing one effect.
type Outcome[E, F any] struct {
	// Event, Next and Message are the fields of an outcome of sim v0.5.0.
	Event   E
	Next    func() (ev E, ok, done bool)
	Message *Message[E]
	// Launch is the machine that the effect starts on the host of its
	// process, and nil for any other effect.
	Launch *Launch[E, F]
	// Stop is the name of a process that an effect started and that this
	// effect stops with the stop event, and empty for any other effect.
	Stop string
}

// Launch is a machine that an effect starts on the host of its process. The
// first event of the machine is an inbound call, whose reply returns to the
// process of the effect as the result of the effect.
type Launch[E, F any] struct {
	// Name is the name of the process of the machine.
	Name string
	// Machine builds the machine from the clock of the host.
	Machine func(clock *Clock) Machine[E, F]
	// Deps are the dependencies of the process.
	Deps []any
	// Call returns the first event of the machine, an inbound call with the
	// reply token reply, in place of the start event of the anatomy.
	Call func(reply uint64) E
	// Result returns the event of the result of the effect from the value and
	// the error of the reply to the call.
	Result func(value any, err error) E
}

// Perform performs the effect f of the process p.
type Perform[E, F any] func(w *World[E, F], p *Process[E, F], f F) Outcome[E, F]
```

- **Launch.** The world adds the process to the host of the process of the effect, gives it a new
  machine in a new epoch, and queues the event of `Call` as its first event. The reply of the machine
  to that call returns to the process of the effect as the event of `Result`, as the reply to a
  message returns. The timeout of the effect bounds the call, and a crash of the host ends the call as
  it ends a message.
- **A started process again.** A start of a process that runs fails the case. A start of a process
  that an effect started, and that has stopped or crashed, starts it again in a new epoch. A start of
  a process of another host, or of a process that `Spawn` added, fails the case.
- **Stop.** The world passes the stop event of the anatomy to the machine, performs the replies that
  it returns, and stops the process before it queues the event of the outcome. A stop of a process
  that no effect started, that runs on another host or that does not run fails the case.
- **Restart.** A crash of a host crashes each of its processes. A restart of the host starts each
  process that `Spawn` added, and a process that an effect started remains crashed until an effect
  starts it again.

### The time of the run

```go
// Now returns the time of the run: Start, plus the time by which Advance has
// advanced the run. It has no drift, and no fault of a host moves it.
func (w *World[E, F]) Now() time.Time
```

A dependency that every process calls, such as the store of sharder, reads the time through `Now`.
The world polls each wait after each step, so an expiry between two steps reaches its machine at the
next step.

### The dependency of a function performer

```go
// Store is the dependency of the cells.
type Store interface {
	// Write writes value for cell.
	//
	//dokimi:fault errors=[ErrUnavailable]
	Write(ctx context.Context, cell, value string) error
}

// Write writes the value of a cell to the store, within Timeout.
//
//dokimi:perform func=write dependency=Store result=Stored timeout=Timeout
type Write struct {
	Value   string
	Token   Token
	Timeout time.Duration
}
```

- **The param.** `dependency` names the interface of the dependency that the function calls. The
  name of the dependency in the faults of a run is the package and the interface, such as
  `placed.Store`, as for `call=`. `result` names the event of the failure of the effect. The field of
  the token type of that event takes the token of the effect, and its field of the type error takes
  the error.
- **The errors.** A loss of the effect takes an error that the case chooses among the errors that the
  directive `fault` declares on the methods of the interface, in the order of the methods, each error
  once. A dependency without such errors fails with `ErrLost`.
- **The faults.** The world strikes the faults of the dependency at the effect as it strikes them at
  an effect under `call=`. A loss before the effect queues the event of its failure and does not call
  the function. A loss after the effect calls the function and queues the event of its failure in
  place of its outcome. An outage loses each effect of the dependency, and a slowness keeps the
  effect in flight until its timeout or the restore.
- **A wait.** The world keeps the effect of each wait. An outage of a dependency ends each wait of an
  effect of the dependency with the event of its failure.

### The inputs of a system

```go
// System builds the processes of a world and their dependencies from inputs
// that it takes from the case: the function takes the inputs, and the function
// that it returns builds the world.
type System[E, F any] func(c *prop.Case) func(w *World[E, F])
```

- **The run.** A case calls `System` with its case before its first step, so each case can run
  another system, and a failing case shrinks its inputs with its other choices. The defect of an
  invariant is a `System` too.
- **Scenarios.** A scenario runs in a case of its own, the replay of the token `prop1:` of no choices.
  The property engine runs that case alone and gives each draw its simplest value, so the system
  takes the simplest choice of each input. Every run of a scenario, the run that counts its steps and
  the run of each crash point, takes the same inputs. `prop.Cases(1)` would not do this: after the
  simplest case it runs the four edge cases of the engine, which take other choices.
- **The directive.** The function under `system` has the signature `func(w *sim.World[E, F])`, or
  `func(c *prop.Case) func(w *sim.World[E, F])` for a system with inputs. The generator wraps a
  system without inputs in a function that ignores the case. The param `defect` of `invariant` takes
  both forms, and a scenario keeps the form `func(w *sim.World[E, F])`.

### The generator

| Directive | Change |
|---|---|
| `//dokimi:perform func=<function> dependency=<interface> result=<event> timeout=<field>` | `dependency` names the dependency of the function, and `result` the event of its failure. The anatomy returns the name of the dependency, the event of the failure and the errors for the effect |
| `//dokimi:system` | The form with inputs, `func(c *prop.Case) func(w *sim.World[E, F])` |
| `//dokimi:invariant defect=<function>` | The defect takes both forms of a system |
| `perform func=` | The function returns `sim.Outcome[E, F]` |

| Directive of an effect | Finding |
|---|---|
| `dependency` without `func` | `BadPerform`: `call=` names its dependency through its method |
| `dependency` without `result` | `BadPerform`: the world builds the event of a failure from `result` |
| `result` with `func` and without `dependency` | `BadPerform`: a function without a dependency has no failure that the world builds |
| `dependency` that names no interface | `BadPerform` |
| An effect without a field of a token, or a `result` without a field of the token type or of the type error | `BadPerform`, as for `call=` |

### Failure handling

| Failure | Behaviour |
|---|---|
| A name of a host that a host has, or a name of a process that a process has | The case fails |
| A spawn on a host that does not run | The case fails |
| A fault of a host by the name of a process of a shared host | The case fails with the name that no host has |
| A start of a process that runs, that `Spawn` added, or that runs on another host | The case fails |
| A stop of a process that no effect started, that runs on another host or that does not run | The case fails |
| A call to a process that an effect started and that the crash of its host ended | The call ends as unknown with `ErrCrashed` |

### Migration

- **Breaking:** `Outcome` has the type arguments `[E, F]`, and `Perform` returns `Outcome[E, F]`. A
  function under `perform func=` changes its result from `sim.Outcome[E]` to `sim.Outcome[E, F]`.
- **Breaking:** `System[E, F]` takes the case. A simulation that a consumer builds by hand wraps its
  system in `func(*prop.Case) func(w *World[E, F])`. A generated simulation wraps it in that function.
- A system of one process per host through `World.Spawn` keeps its choices, its names and its
  faults.
- The generated simulations of the corpus change, and the corpus gains the fixture `sim/placed`: an
  agent that starts and stops the cells of a host, a store under `perform func= dependency=` whose
  writes carry the time of the run, and a system that takes the timeout of a write from the case.

### Verdicts

| Item | Verdict | Reason | Reverses the verdict |
|---|---|---|---|
| Hosts of machine processes | in | sharder runs the machines of a process together, and every simulator in the motivation crashes a machine with every process on it | none known |
| A component on a shared host | out | A component has a host of `go.dokimi.dev/host` with its own network and disk | a consumer that runs a component beside machines in one process |
| A launch and a stop through the outcome | in | The world applies the outcome, so a loss after the effect drops the launch as it drops a message | none known |
| `Outcome[E, F]` | in | A start builds a `Machine[E, F]`, so the outcome needs the type of the effects | none known |
| `World.Now` | in | A dependency of every process needs one time that no fault moves | none known |
| The errors of a dependency from the methods of its interface | in | The directive `fault` declares the errors on the methods, and a function can call each of them | a consumer that needs other errors for one effect |
| A system with inputs from the case | in | A failure that the size of a system causes shrinks to a smaller system | none known |

## Alternatives considered

### A method of the world that starts a machine

A performer would call `World.Launch` during its perform. The launch would apply also when a loss
after the effect discards the outcome, so the process of the effect would take the result of the
start beside the event of the failure.

### A registry of the kinds of machines

The system would register a builder for each kind of machine on a host, and an outcome of one type
argument would contain the kind. The builder of a cell of sharder needs the epoch and the
configuration of the add, which only the performer has.

### The case on the world

`World.Case` would give the system the case. Every action, performer and invariant could then take
choices from it outside the steps of the run. A choice that a goroutine of a component takes from
the case makes the run nondeterministic in native mode.

### The errors on the directive perform

`perform` would take the errors of the dependency as a param. The directive `fault` already declares
them on the methods, so the param would be a second copy that can disagree.

## Drawbacks

- Every function under `perform func=` changes its result type, and every simulation that a consumer
  builds by hand wraps its system. Both changes break the code of a consumer.
- The faults of a process take the name of a host, so a test that crashes a process of a shared host
  names the host.
- The restart of a host does not start a process that an effect started. A call of a client to such
  a process ends as unknown until an effect starts the process again.
- The errors of a dependency under `perform func=` are those of every method of its interface, also
  of a method that the function does not call.

## Open questions

None.

## References

| What | Where |
|---|---|
| Issue 63, the machines of one process do not crash, pause and restart together | <https://github.com/dokimasia/dokimi/issues/63> |
| Issue 64, a dependency that several processes share has no clock of the run | <https://github.com/dokimasia/dokimi/issues/64> |
| Issue 65, a fault of a dependency never strikes an effect that a function performs | <https://github.com/dokimasia/dokimi/issues/65> |
| Issue 66, a system cannot take inputs from its case | <https://github.com/dokimasia/dokimi/issues/66> |
| RFC-0025, a deterministic simulation of sans-IO machines | 0025-a-simulation-of-sans-io-machines.md |
| RFC-0036, a simulation of components that run goroutines | 0036-components-that-run-goroutines.md |
| The simulator of FoundationDB | <https://apple.github.io/foundationdb/testing.html> |
| madsim | <https://github.com/madsim-rs/madsim> |
| turmoil | <https://github.com/tokio-rs/turmoil> |
| TigerBeetle, the VOPR | <https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/internals/vopr.md> |
| Hypothesis, the rules of a state machine | <https://hypothesis.readthedocs.io/en/latest/stateful.html> |
