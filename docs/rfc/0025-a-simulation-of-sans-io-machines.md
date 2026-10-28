---
rfc: 0025
title: A deterministic simulation of sans-IO machines
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Draft
created: 2026-10-28
updated: 2026-10-28
discussion: https://github.com/dokimasia/dokimi/issues/43
supersedes: none
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0025: A deterministic simulation of sans-IO machines

## Summary

A sans-IO machine applies one event at a time and returns the effects that a runtime performs, such as a commit to a store or a reply to a caller. A system of such machines fails under orders of events and faults that no check of one machine produces: a reply that arrives after a crash, a commit that a fault loses after it applied, a process whose clock runs fast. A new generator, `sim`, writes a deterministic simulation of a system of sans-IO machines from directives on their declarations. The simulation runs every machine of the system on the goroutine of a case of the property engine. Every choice of a run is a choice of the case, so a failing run shrinks and replays from its seed:

- which process receives its next event, under a strategy of scheduling
- how far a clock advances
- which fault strikes, within the faults that the system declares it tolerates
- which call a client makes, with which inputs

The simulation records every call of a client in the history of the case, with an unknown outcome for a call that a crash or a timeout leaves open. It checks that history against a sequential specification for linearizability, and against the anomalies of Elle for serializability and snapshot isolation. It also checks five rules of a runtime, the invariants of the consumer after every step, and convergence once the faults end. It checks that the run itself is deterministic, by a checksum of every step that a repeated run compares. A new module of dokimi-go, `go.dokimi.dev/sim`, is the runtime of the simulation.

The workloads, the faults and the checkers cover what Jepsen tests on a system of sans-IO machines, and every failure shrinks and replays from a seed.

## Motivation

sharder runs each component of a namespace as a sans-IO machine (`internal/machine/machine.go`). `Step(ev Event, out []Effect) []Effect` applies one event and appends its effects, and `Err() error` reports a machine that stopped itself. The package declares 32 event types and 21 effect types, and an unexported marker method closes each set. A `Token` links each result to its effect, and a `ReplyToken` links each `Reply[T]` to its inbound call (`internal/machine/inbound.go:202`). The package comment states the rules of a runtime (`internal/machine/doc.go`).

The specs of the four machines of sharder each contain a runtime for one machine, of 429 to 816 lines (`internal/machine/{root,cell,agent,routes}/helpers_test.go`). The production runtime, `internal/driver`, implements the same rules again in 1,150 lines. None of the five runs the machines of a namespace together, and none injects faults at the network, the store or the processes. A simulator that sharder writes by hand would be a sixth runtime. Its scheduler, virtual clocks, routing of results by token, faults and traces do not depend on sharder.

Jepsen is the established test of such a system. It runs the real system on a cluster, drives it with concurrent clients, injects faults through a nemesis, and records a history whose operations end as `:ok`, `:fail` or `:info`, the last for an outcome that the client cannot know. Knossos checks the history for linearizability, and Elle infers the anomalies of serializability and snapshot isolation from the dependencies that the reads of list-append transactions reveal. A run of Jepsen needs a cluster. Its schedule depends on the real network and kernel, so a failure seldom repeats, and nothing shrinks it.

The parts that a deterministic replacement needs exist in assert-go and dokimi:

- `history.History` records invocations and completions by client, with the outcomes `OK`, `Fail` and `Unknown`, and with the keys that partition a check (`assert-go/history/history.go`, `call.go`).
- `history.Linearizable` checks a history against a `history.Spec`, whose `Next` may return several states and covers calls of an unknown outcome (`spec.go:49-65`). `history.SpecFrom` builds a spec from a reference.
- `history.Serializable` and `history.HasSnapshotIsolation` check list-append and read-write register transactions for the anomalies of Elle, from assert-go v0.11.0: G0, G1c, G-single, G-nonadjacent and G2, a garbage read, a duplicate append, an internal inconsistency, an aborted read and an intermediate read (`isolation.go`, `anomaly.go`).
- `history.FromIntervals` builds a history from timed intervals, for a run outside a simulation (`interval.go`).
- A machine of `stateful` runs generated actions as choices of a case, checks its `Spec` against the history after every step, and runs `Invariant` and `Settle` (`stateful/machine.go`). `prop.Case` records labels through `Classify`, takes a fingerprint of a state through `Observe`, and `prop.Require` demands a share of cases for a label.
- The fault strategies of `fault` and the doubles of `stub` make a dependency fail on demand, through the `Injector` of the calls that fail on demand.

The deterministic simulator of TigerBeetle, turmoil and madsim run every node of a system on one thread under a seeded scheduler, inject faults at the network and the storage, and replay a failure from its seed. gosim does the same for a whole Go program by translating its source onto a coroutine scheduler. Its scheduler picks a runnable goroutine uniformly at random at each step, and it hashes each pick into a checksum that a repeated run compares (`rt/runtime/runtime.go:291`, `checksum.go`). A sans-IO machine performs no I/O, so a simulation of sans-IO machines needs neither translation nor interception.

## Detailed design

### Components

| Component | Where | Responsibility |
|---|---|---|
| The directives of `sim` | The declarations of the consumer | State the machine, its tokens and replies, the performer of each effect, the system, the faults it tolerates, the oracles, the scenarios and the actions |
| The generator `sim` | dokimi | Writes the typed glue of the world and the tests of the simulation |
| `go.dokimi.dev/sim` | dokimi-go | The world, its processes, clocks and epochs, its faults, the scheduling strategies, the run, the log and its checksum |
| `go.dokimi.dev/assert/history` and `stateful` | assert-go | The history of the calls of clients, the machine of the run, and the checkers |
| The claims of the runtime rules | The conformance spec | The five rules, with their probes, reason codes and defect forms |

### The directives

| Directive | On | States |
|---|---|---|
| `//dokimi:sim start=Started stop=Stopped cancel=Cancelled` | The machine interface | The step method and the error method, the closed sets of events and effects of their signatures, and the events of the lifecycle |
| `tolerate=crash:1,partition:1,pause:1` on `sim` | The machine interface | The faults that the system tolerates at once. The run never exceeds them while faults are active |
| `strategy=uniform` or `strategy=pct:<depth>` on `sim` | The machine interface | The strategy that chooses the next process. `pct:3` by default |
| `require=crash:0.2,partition:0.1` on `sim` | The machine interface | The share of cases in which each fault must strike, through `prop.Require` |
| `isolation=serializable` or `isolation=snapshot-isolation` on `sim` | The machine interface | That the transactions of the clients are list-append or register transactions, the operation `txn` of assert-go, and the isolation level that the history must have |
| `//dokimi:token` | A token type | That a field of this type in an effect links the effect to the event of its result |
| `//dokimi:token reply` | A token type | That a field of this type in an event links an inbound call to the effect that replies to it |
| `//dokimi:reply value=Value err=Err to=To` | The reply effect type | The fields of a reply that contain its value, its error and the token of its call |
| `//dokimi:perform call=<callable> result=<event>` | An effect type | The callable of a dependency that performs the effect, such as `store.Writer.Commit`, and the event of its result, such as `Committed` |
| `//dokimi:perform func=<function>` | An effect type | A function of the consumer that performs the effect, for an effect whose result needs the logic of the consumer |
| `//dokimi:system` | A function of the consumer | The processes of the system and their dependencies |
| `//dokimi:action` | A function of the consumer | A change that starts outside the machines, such as the start of a node, with inputs that the case chooses |
| `//dokimi:action call keys=<function>` | A function of the consumer | An inbound call of a client, which the history records, and the function that returns the keys that the call touches |
| `//dokimi:action call final` | A function of the consumer | A call that the run makes after the faults end, such as a read of every key |
| `//dokimi:spec` | A function of the consumer | The sequential specification of the calls of the clients, as a `history.Spec`, or a reference that `history.SpecFrom` wraps |
| `//dokimi:invariant defect=<function>` | A function of the consumer | A safety property that the run checks after every step, and a function that builds the system with one planted defect of the property |
| `//dokimi:converges within=<function>` | A function of the consumer | A liveness property that the run checks once the faults end, and the bound of virtual time within which it must become true |
| `//dokimi:fingerprint` | A function of the consumer | A hash of the state of the system, which the run passes to `Case.Observe` after every step |
| `//dokimi:scenario` | A function of the consumer | A fixed run whose events the crash points cover |

The generator refuses a machine interface whose step method does not have the form `Step(ev E, out []F) []F` with interface types `E` and `F`, or whose error method is not `Err() error`. The refusal names the method and the form. The set of events is the types of the package of `E` that implement `E`, and the set of effects is the types that implement `F`. An effect type without `perform` refuses with the type and the remedy of the directive. A system with an action `call` and without `reply` refuses, because the run cannot complete a call in the history.

### The world

```go
// Package sim is a deterministic simulation of a system of sans-IO machines.
// Every machine of a world runs on the goroutine of one case of the property
// engine, and every choice of a run is a choice of the case.
package sim

// Machine is a sans-IO machine: it applies one event at a time and appends
// the effects of the event to out.
type Machine[E, F any] interface {
	Step(ev E, out []F) []F
	Err() error
}

// World is the system of one run: its processes, the events that wait for
// each process, the effects in flight, the faults, the log of its steps and
// the history of the calls of its clients.
type World[E, F any] struct{ /* unexported */ }

// Spawn adds a process named name to the world. newMachine builds the
// machine of the process from its clock. A process that restarts builds a
// new machine with the same function, in a new epoch.
func (w *World[E, F]) Spawn(name string, newMachine func(clock *Clock) Machine[E, F])

// Deliver queues ev for the process named to, in the current epoch of the
// process. An action of the consumer calls it for a change outside the
// machines or for a forged message.
func (w *World[E, F]) Deliver(to string, ev E)

// Call queues ev, an inbound call with the reply token of the event, for the
// process named to, and records its invocation by client in the history of
// the case with the operation, the arguments and the keys. The reply with
// the same token completes the call.
func (w *World[E, F]) Call(client int, to string, operation string, args []any, keys []any, ev E)

// After queues ev for the process p once the clock of p has advanced by d.
// A performer of a timer calls it.
func (w *World[E, F]) After(p *Process[E, F], d time.Duration, ev E)

// Process returns the process named name, and nil for an unknown name.
func (w *World[E, F]) Process(name string) *Process[E, F]

// Log returns the steps of the machines so far.
func (w *World[E, F]) Log() []Entry[E, F]

// History returns the history of the calls of the clients, which is the
// history of the case.
func (w *World[E, F]) History() *history.History

// Checksum returns the checksum of the log so far.
func (w *World[E, F]) Checksum() uint64

// System builds the processes of a world and their dependencies. The
// function of the directive system has this type, and so has the defect
// function of each invariant.
type System[E, F any] func(w *World[E, F])

// Process is one process of a world: its name, its epoch, its clock and its
// machine.
type Process[E, F any] struct{ /* unexported */ }

// Name returns the name of the process.
func (p *Process[E, F]) Name() string

// Epoch returns the epoch of the process, from 1, which each restart raises.
func (p *Process[E, F]) Epoch() uint64

// Clock returns the clock of the process.
func (p *Process[E, F]) Clock() *Clock

// Entry is one step of a machine: the process and its epoch, the event, the
// effects that the machine returned, the time of its clock and the faults in
// force.
type Entry[E, F any] struct {
	Process string
	Epoch   uint64
	Event   E
	Effects []F
	Time    time.Time
	Faults  []string
}

// Check is the test handle of an invariant or a convergence. It implements
// assert.TB, and Statement returns the statement of its directive.
type Check struct{ /* unexported */ }

// Statement returns the statement of the check: the first sentence of the
// doc comment of its function.
func (c *Check) Statement() string

// Clock is the virtual clock of one process. It advances when the run
// advances time, at the rate of its drift, and it does not advance while its
// process is paused.
type Clock struct{ /* unexported */ }

// Now returns the time of the clock.
func (c *Clock) Now() time.Time
```

The log of the world records the steps of the machines, and the history of the case records the calls of the clients. The two are different records. A step of a machine is no call, and the checkers of assert-go read calls.

The machine of the consumer reads the time through its own clock interface, so the system function adapts `*sim.Clock` to it. The world advances time in jumps to the next timer of any process, by an amount that the case chooses, and each clock advances by that amount times its drift.

### Epochs

Each queued event, each effect in flight and each pending wait has the epoch of the process that it belongs to. A restart raises the epoch of the process. The world drops an event of an earlier epoch before its delivery and records the drop in the log. The world therefore never delivers a result of an effect that the old machine started to the new machine, also when the new machine has issued a token of the same number.

### Performing an effect

The generator writes one switch over the closed set of effects. For an effect with `call=`, the generated performer calls the dependency through the double that `stub` writes for it and builds the event of the result. For an effect with `func=`, it calls the function of the consumer. For the reply effect, the world completes the call of its token in the history: `OK` with the value when the error is nil, and `Fail` with the error otherwise.

```go
// Outcome is the result of performing one effect.
type Outcome[E any] struct {
	// Event is the event of the result, and the zero E for an effect without
	// a result, such as a reply.
	Event E
	// Ready is closed once a wait has a next event, and nil for an effect
	// that does not wait.
	Ready <-chan struct{}
	// Next returns the next event of a wait, and false once the wait ends.
	Next func() (E, bool)
}

// Perform performs the effect f of the process p.
type Perform[E, F any] func(w *World[E, F], p *Process[E, F], f F) Outcome[E]
```

A wait is an effect whose result is a stream of events, such as a watch or a campaign for a term. The world checks the `Ready` channel of each pending wait after every step without blocking, in the order in which the waits started, and queues the next event of each ready wait. A performer must not start a goroutine. The channels that it returns are closed by the doubles of the world, which run on the goroutine of the case.

### Calls of an unknown outcome

A call ends as `Unknown` in the history when its outcome cannot be known:

- The runtime cancels it, because its caller stopped waiting, and the world delivers `Cancelled`.
- The process of the call crashes before it replies.
- The run ends while the call is open.

`history.Linearizable` lets an unknown call take effect or not, and its client continues on a new process, as Jepsen treats an `:info` operation.

### Faults

| Fault | Strikes | Effect on the run |
|---|---|---|
| `drop` | A message to a process | The event never arrives |
| `delay` | A message | The event arrives after a later event |
| `duplicate` | A message | The event arrives twice |
| `partition` | The processes, in one of the shapes `isolate`, `halves`, `majorities-ring` and `bridge` | Every message across the cut is dropped until `repair` |
| `repair` | A partition | The cut ends |
| `lose-before` | An effect on a dependency | The call fails, and the dependency does not apply it |
| `lose-after` | An effect on a dependency | The dependency applies the call, and the call fails |
| `outage` | A dependency | Every call fails until the outage ends |
| `slow` | A dependency | Every call takes longer than its timeout until the fault ends |
| `crash` | A process | The machine and its effects in flight are gone, and its open calls end as unknown. The dependencies keep what they applied |
| `pause` | A process | The process receives no event, and its clock does not advance while the other clocks advance |
| `restart` | A crashed process | A new machine of the process receives `Started` in a new epoch |
| `clock-jump` | The clock of a process | The clock moves forwards or backwards by an amount that the case chooses |
| `clock-strobe` | The clock of a process | The clock alternates between two offsets for a number of steps |

The shapes of a partition are those of the nemesis of Jepsen. `isolate` cuts one process off, `halves` splits the processes in two, `majorities-ring` gives each process a majority that overlaps its neighbours, and `bridge` joins two halves through one process. The faults of a dependency arm the `Injector` of the calls that fail on demand on the double of the dependency, with the errors of its directive `fault`. A dependency that persists through `medium.FS` takes the faults of the breakable medium. A forged message is an action of the consumer that calls `Deliver` with an event that no effect produced.

A system that tolerates one crash fails convergence trivially under two, so the run keeps the faults in force at once within the param `tolerate`. An action of a fault that would exceed it is not enabled in that state. A system without `tolerate` tolerates no fault, and the run checks it without faults.

### The run

The generator writes a machine of `stateful` over the world, in three phases:

| Phase | Actions |
|---|---|
| Active | `deliver`, `advance`, one action per fault within `tolerate`, and one action per `//dokimi:action` that is not `final` |
| Recovery | `Settle` ends every fault and restarts every crashed process |
| Final | Each `final` action, then `deliver` and `advance` until the bound of each `//dokimi:converges` passes |

`deliver` asks the strategy for a process with a pending event, and the process steps on its oldest event of its current epoch. `uniform` picks a process uniformly. `pct:<depth>` draws a priority for each process and `depth-1` change points from the case, and picks the pending process of the highest priority, lowering it at each change point, as the strategy `PCT` of the stateful scheduler of assert-go orders its tasks. A bug that needs a certain order of `depth` events has a probability of discovery that PCT bounds from below, where a uniform pick gives none.

After every step, `Invariant` checks:

- the five rules of a runtime over the log
- each invariant of the consumer
- the history against `Machine.Spec`, which the run sets to the spec of `//dokimi:spec`

In `Settle`, once the final phase ends, the run checks each convergence and the whole history: `history.Linearizable` against the spec, partitioned by the keys of each call, and `history.Serializable` or `history.HasSnapshotIsolation` when the directive states `isolation`.

After every step the run also calls `Case.Classify` with the kind of each fault that struck and of each effect that the world performed, and `Case.Observe` with the fingerprint of `//dokimi:fingerprint`. `require` turns the labels of the faults into `prop.Require`, so a suite whose cases seldom crash a process fails instead of passing on calm runs.

For each `//dokimi:scenario`, the generated test runs the scenario once to count its events. It then runs it again for each event, crashing the process of the event once before the event and once after it, and restarting it. The invariants, the rules of the runtime and the history run after every step of each run.

### Determinism

The world hashes each entry of its log into a checksum: the process, the epoch, the typed literal of the event and of each effect, and the time of the clock. The run repeats every failing case, and every 16th case of a passing run, with the same choices, and compares the two checksums after each step. A difference fails the test with the first step at which the two runs differ, and with the causes that the world cannot see:

- iteration over a Go map in a machine or a performer
- the global source of `math/rand`
- a read of the wall clock
- a goroutine in a performer

A machine whose effects depend on the order of a map therefore fails at once, instead of producing a seed that does not replay.

### The rules of a runtime

The spec states five claims with the generator `sim` and the machinery `a-simulated-system`. A claim of `sim` binds to each machine interface under the directive, as a claim of `enum` binds to each enumeration:

| Claim | Law | Reason | Defect form |
|---|---|---|---|
| `replies-once` | A machine returns at most one reply to each inbound call | `second-reply` | `replies-twice` |
| `replies-before-stop` | On the stop event, a machine replies to each inbound call that it has not replied to | `call-without-reply` | `forgets-a-call-on-stop` |
| `quiet-after-stop` | After the stop event, a machine returns replies alone | `effect-after-stop` | `writes-after-stop` |
| `quiet-after-error` | A machine returns no effect after the step in which its error method returns an error | `effect-after-error` | `steps-after-error` |
| `fresh-tokens` | Each effect that starts a call or a wait has a token that no earlier effect of its machine in the same epoch has | `reused-token` | `reuses-a-token` |

Each probe is a judge over the log of the world. Each defect form wraps a fresh machine and changes the effects of its step, such as appending a second reply. Falsification runs each claim on its defect at a seed from the key of the check, as every proof of a property check runs.

### The invariants of the consumer

An invariant is a function of the consumer that asserts through assert-go with the statement of its directive as the contract:

```go
// No two agents hold a lease on one shard at once.
//
//dokimi:invariant defect=systemWithTwoLeaders
func exclusiveLeases(h *sim.Check, w *sim.World[machine.Event, machine.Effect]) {
	expect.Empty(h, overlappingLeases(w), h.Statement())
}
```

The defect is a function with the signature of the system function, which builds the same system with one planted defect, such as an agent that keeps its lease after a failed renewal. Falsification runs the scenarios and the random runs on the defect system at a seed from the key of the invariant, and proves the invariant by the contract of its first failure, as every check of the consumer is proven.

### The generated tests

| Test | What it runs |
|---|---|
| `TestSimulation` | The property over random runs, the crash points of each scenario, and the determinism check |
| `TestSimulationFalsify` | Each rule of the runtime on its defect form, and each invariant on its defect system |
| `FuzzSimulation` | The choices of the run as the input of a fuzz target, through `prop.Fuzz` of assert-go |

A failing run reports the shrunk log with the step of each process, its epoch, its clock and the faults in force, the history of the calls with the anomaly or the call that no linearization explains, and the seed that replays it through `DOKIMI_SEED`.

### The capabilities of Jepsen

| Jepsen | The simulation |
|---|---|
| A generator of operations from concurrent clients | One action `call` per operation, with inputs from the case. Each open call keeps its client busy, so the clients run concurrently |
| `:ok`, `:fail` and `:info` | `OK`, `Fail` and `Unknown` in `history.History` |
| Knossos, linearizability against a model | `history.Linearizable` against the spec of `//dokimi:spec`, partitioned by key, and `Machine.Spec` after every step |
| Elle on list-append transactions | `history.Serializable` and `history.HasSnapshotIsolation` |
| Elle on read-write registers | The same checks over register transactions. Blind writes reveal no order, so a register history shows fewer anomalies, as in Elle |
| The checkers of a set, a queue, a counter and a bank | A spec of the consumer for linearizability, and an invariant for a quantity that the system conserves, such as the total of a bank |
| A nemesis of partitions, kills, pauses and clocks | The faults `partition` in four shapes, `crash`, `pause`, `restart`, `clock-jump` and `clock-strobe` |
| Faults of the disk, such as LazyFS | The breakable medium, for a dependency that persists through `medium.FS` |
| A final generator after the nemesis ends its faults | The phase `final` |
| Timelines and latencies | The log, with the virtual time of every step and of every call |
| Reproduction of a failure | The seed, on any machine, and a shrunk run |

Jepsen also tests a system that it does not compile: the real binaries of any language, over a real network stack and kernel. The simulation cannot test those, because it runs the machines and the doubles of their dependencies. `history.FromIntervals` builds a history from the timed intervals of a run on a real cluster, so the same specs and checkers apply to such a run. A runner that provisions the cluster, cuts its network and moves its clocks is a design of its own, which this RFC does not propose.

## Alternatives considered

### A simulator that each consumer writes

The consumer writes the world for its own machines, as sharder's specs write a runtime for each machine.

**Why not:** the scheduler, the clocks, the routing by token, the faults and the history do not depend on the machines. sharder already has five runtimes of the same rules, and a simulator would be a sixth that no suite checks.

### A library without a generator

`go.dokimi.dev/sim` alone, with the consumer writing the switch over its effects and the routing of its results.

**Why not:** the switch over 21 effect types and the routing of 32 event types are the code that the five runtimes of sharder repeat. The closed sets, the tokens and the performers are declared in the source, so a generator writes the switch from them, and a new effect type without a performer fails generation.

### Jepsen as the oracle

The consumer runs Jepsen against a cluster of the system, and the simulation tests the machines alone.

**Why not:** a run of Jepsen needs a cluster, its schedule depends on the real network, and a failure neither repeats nor shrinks. Its checkers exist in assert-go as Go functions over a history, so the simulation applies them to every case, in `go test`.

### A translation of the whole program, as gosim does

The simulation translates the production runtime and the machines onto a coroutine scheduler and simulates the network and the filesystem below them.

**Why not:** gosim supports the Go families 1.25 and 1.26 alone (`internal/buildspec/families_gen.go`), and sharder and dokimi build with Go 1.27.3. It calls the private symbols of each Go runtime through `//go:linkname`, runs through its own launcher instead of `go test`, picks the next goroutine uniformly at random, and aborts `AllocsPerRun` and benchmarks (`internal/compat/shims/testing/missing.go`). It tests the production runtime, which this simulation replaces with a generated one, so the two compose: a consumer runs its production runtime under gosim, and its machines under this simulation.

### A hypervisor that runs the whole system

A deterministic hypervisor runs the processes of the system and explores their schedules and faults.

**Why not:** it needs a build of every process and a service outside the test, and it cannot run inside `go test` in the CI of a consumer. A failure in it does not shrink to the steps of the machines.

## Drawbacks

- **Many directives.** A consumer writes the directives on its machine interface, its two token types, its reply type, each effect type, and its functions. sharder writes about 30 of them.
- **Code of the consumer.** The consumer writes the system, the performers of the effects that need its logic, the spec of its calls, the invariants with their planted defects, the scenarios and the actions.
- **The cost of a run.** A case builds a world with every process of the system. A property of 100 cases of 200 steps runs 20,000 steps of the machines, the determinism check repeats one case in 16, and each scenario of n events adds 2n runs for its crash points. A check of linearizability is NP-complete in general, so a history of many concurrent calls needs a budget and a partition by key.
- **A generated runtime.** The simulation checks the machines under a runtime that dokimi generates. A defect of the production runtime alone, such as a driver that delivers a result in the step that returned its effect, needs the production runtime under gosim or the checks of the runtime rules on its own log.
- **One language.** This RFC designs the Go runtime and template tree. The directives and the claims do not depend on a language, and a runtime of another language is not proposed here.

## Open questions

None.

## References

| What | Where |
|---|---|
| The request of the generator | <https://github.com/dokimasia/dokimi/issues/43> |
| sharder, the machine interface and the rules of a runtime | <https://github.com/thesm-os/sharder/tree/main/internal/machine> |
| Jepsen | <https://github.com/jepsen-io/jepsen> |
| Knossos, the linearizability checker of Jepsen | <https://github.com/jepsen-io/knossos> |
| Elle, inferring isolation anomalies from histories | <https://arxiv.org/abs/2003.10554> |
| assert-go, the packages history and stateful | <https://github.com/dokimasia/assert-go/tree/main/history> |
| gosim, a deterministic simulation of Go programs by translation | <https://github.com/neria-cloud/gosim> |
| TigerBeetle, the deterministic simulator VOPR | <https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/internals/vopr.md> |
| turmoil, a simulation of hosts and a network for Rust | <https://github.com/tokio-rs/turmoil> |
| madsim, a deterministic simulator of a Rust runtime | <https://github.com/madsim-rs/madsim> |
| assert-spec, the workload of Elle on read-write registers | <https://github.com/dokimasia/assert-spec/issues/27> |
