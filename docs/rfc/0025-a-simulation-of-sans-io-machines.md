---
rfc: 0025
title: A deterministic simulation of sans-IO machines
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-28
updated: 2026-11-02
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
| `tolerate=[crash:1, partition:1, pause:1]` on `sim` | The machine interface | The faults that the system tolerates at once. The run never exceeds them while faults are active |
| `strategy=uniform` or `strategy=pct:<depth>` on `sim` | The machine interface | The strategy that chooses the next process. `pct:3` by default |
| `require=[crash:0.2, partition:0.1]` on `sim` | The machine interface | The share of cases in which each fault must strike, through `prop.Require` |
| `isolation=serializable` or `isolation=snapshot-isolation` on `sim` | The machine interface | That the transactions of the clients are list-append or register transactions, the operation `txn` of assert-go, and the isolation level that the history must have |
| `//dokimi:token` | A token type | That a field of this type in an effect links the effect to the event of its result |
| `//dokimi:token reply` | A token type | That a field of this type in an event links an inbound call to the effect that replies to it |
| `//dokimi:reply value=Value err=Err to=To unknown=[ErrStopped]` | The reply effect type | The fields of a reply that contain its value, its error and the token of its call, and the errors of a reply that leave the outcome of its call unknown, such as the error of a stop that can follow the effect of the call. `unknown` is optional |
| `//dokimi:perform call=<callable> result=<event> timeout=<field>` | An effect type | The callable of a dependency that performs the effect, such as `store.Writer.Commit`, the event of its result, such as `Committed`, and the field of the effect that bounds the call, such as `Timeout`. `timeout` is optional |
| `//dokimi:perform func=<function>` | An effect type | A function of the package of the simulation that performs the effect, for an effect whose result needs the logic of the consumer |
| `//dokimi:system` | A function of the consumer | The processes of the system and their dependencies |
| `//dokimi:action` | A function of the consumer | A change that starts outside the machines, such as the start of a node, with inputs that the case chooses |
| `//dokimi:action call` | A function of the consumer | An inbound call of a client, which the history records with the keys that its run passes to `World.Call` |
| `//dokimi:action call final` | A function of the consumer | A call that the run makes after the faults end, such as a read of every key |
| `//dokimi:spec` | A function of the consumer | The sequential specification of the calls of the clients, as a `history.Spec`, or a reference that `history.SpecFrom` wraps |
| `//dokimi:invariant defect=<function>` | A function of the consumer | A safety property that the run checks after every step, and a function that builds the system with one planted defect of the property |
| `//dokimi:converges within=<function>` | A function of the consumer | A liveness property that the run checks once the faults end, and the bound of virtual time within which it must become true |
| `//dokimi:fingerprint` | A function of the consumer | A hash of the state of the system, which the run passes to `Case.Observe` after every step |
| `//dokimi:scenario` | A function of the consumer | A fixed run whose events the crash points cover |

The generator refuses a machine interface whose step method does not have the form `Step(ev E, out []F) []F` with interface types `E` and `F`, or whose error method is not `Err() error`. The refusal names the method and the form. The set of events is the types of the package of `E` that implement `E`, and the set of effects is the types that implement `F`. An effect type without `perform` refuses with the type and the remedy of the directive. A system with an action `call` and without `reply` refuses, because the run cannot complete a call in the history.

### The functions of the consumer

The generator writes the simulation of a machine interface into the package that the directive `out` of the interface names, or into the package of the interface without one. The functions of the consumer under the directives of `sim` are in that package, so the simulation can call a function that the package does not export. A function under a directive in another package belongs to the simulation of another machine interface. The param `func` of `perform` names a function of the package of the simulation by its name. The package of the effect cannot import the package of the simulation, so the param is no reference in the scope of the effect.

| Directive | Signature |
|---|---|
| `system`, `scenario` and the param `defect` of `invariant` | `func(w *sim.World[E, F])` |
| `action` | `func(c *prop.Case) func(w *sim.World[E, F])` |
| `action call` | `func(c *prop.Case) func(w *sim.World[E, F], client int)` |
| `spec` | `func() history.Spec[S]`, or `func() history.Subject` for `history.SpecFrom` |
| `invariant` and `converges` | `func(h *sim.Check, w *sim.World[E, F])` |
| The param `within` of `converges` | `func() time.Duration` |
| `fingerprint` | `func(w *sim.World[E, F]) uint64` |
| `perform func=` | `func(w *sim.World[E, F], p *sim.Process[E, F], f T) sim.Outcome[E]` |

`E` and `F` are the event and the effect interfaces of the machine, and `T` is the type of the effect. The outer function of an action takes the inputs of a step from the case, and the inner function runs the step on the world. The name of an action and of a scenario is the name of its function. The statement of an invariant and of a convergence is the first sentence of the doc comment of its function.

A function with another signature refuses under `BadFunction`, with its directive and the signature that the directive needs. A second system, spec or fingerprint refuses too, and so does a package of a simulation without a system. A generic reply type with one type parameter whose constraint is a union of types, such as `struct{} | Ack`, gets a case of the anatomy for each type of the union, such as `Reply[struct{}]` and `Reply[Ack]`. The union is the constraint of the type parameter, or the one constraint element of the interface that the type parameter names. Any other generic reply refuses under `BadReply`, because its instantiations are not a list: a constraint without a union, such as `any`, a union with an approximation, such as `~int | string`, and two type parameters.

The generator derives the names of the declarations of a simulation from the name of its machine interface, such as `Machine`. The anatomy is the variable `machineAnatomy`, the performer of the effects is the function `performMachine`, and the simulation is the function `MachineSimulation`. The tests call `MachineSimulation` from the external test package.

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
// new machine with the same function, in a new epoch. deps are the
// dependencies of the process, such as the doubles of its stores, which a
// performer reads through Dependency.
func (w *World[E, F]) Spawn(name string, newMachine func(clock *Clock) Machine[E, F], deps ...any)

// Dependency returns the first dependency of the process p of the type T.
// A restart keeps the dependencies of a process. Dependency fails the case
// of w for a process without a dependency of the type T.
func Dependency[T, E, F any](w *World[E, F], p *Process[E, F]) T

// Context returns the context of the case of the world, which a performer
// passes to the call of a dependency.
func (w *World[E, F]) Context() context.Context

// Deliver queues ev for the process named to, in the current epoch of the
// process. An action of the consumer calls it for a change outside the
// machines or for a forged message.
func (w *World[E, F]) Deliver(to string, ev E)

// Call queues the inbound call that call builds from a new reply token for
// the process named to, and records its invocation by client in the history
// of the case with the operation, the arguments and the keys. The world
// numbers the reply tokens of a run from 1. The reply with the same token
// completes the call. A call to a crashed process ends as unknown with
// ErrCrashed, and a call to a stopped process as unknown with ErrRunEnded.
func (w *World[E, F]) Call(
	client int, to string, operation string, args []any, keys []any, call func(reply uint64) E,
)

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

// New returns an empty world whose machines have the anatomy a, and which
// takes its choices from the case c. Run builds a new world for each case,
// and a test of the world builds its own.
func New[E, F any](c *prop.Case, a Anatomy[E, F]) *World[E, F]

// Pending returns the names of the processes with an event of their current
// epoch to step on, in the order of their spawns.
func (w *World[E, F]) Pending() []string

// Step makes the process named name step on its oldest event, and drops each
// event of an earlier epoch before it. It returns false for a process without
// such an event.
func (w *World[E, F]) Step(name string) bool

// Advance advances the time of the run by d, and queues each timer that is
// due by then.
func (w *World[E, F]) Advance(d time.Duration)

// Next returns the time until the next timer of any process, and false when
// no timer is pending.
func (w *World[E, F]) Next() (time.Duration, bool)

// Stop passes the stop event to each process that runs, performs the replies
// that it returns, and ends each open call of the history as unknown.
func (w *World[E, F]) Stop()
```

The faults of the next section are methods of the world as well, such as `Crash(name)`, `Restart(name)` and `Partition(sides ...[]string)`, so a scenario and a test call them on a world without the run.

```go
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
// effects that the machine returned, the error of the machine after the
// step, the time of its clock and the faults in force. Dropped reports an
// event of an earlier epoch that the world dropped instead of a step, which
// has no effects.
type Entry[E, F any] struct {
	Process string
	Epoch   uint64
	Event   E
	Effects []F
	Err     error
	Time    time.Time
	Faults  []string
	Dropped bool
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

The machine of the consumer reads the time through its own clock interface, so the system function adapts `*sim.Clock` to it. The world advances time in jumps to the next timer of any process, by an amount that the case chooses, and each clock advances by that amount times its drift. When the world spawns a process, the case chooses the drift of its clock as an offset from −500 to +1,000 per mille, which shrinks towards 0. So a clock runs between half and twice as fast as the run, and a shrunk run has no drift.

A jump and a strobe move the time that `Now` returns, and they do not move the timers of the process. A timer waits for the clock of its process to advance, as a timer on the monotonic clock of an operating system does. The timeout of an effect is such a timer, so a jump neither delays nor hastens it. A convergence keeps its bound under the faults of a clock too. The `TimeSim` of TigerBeetle also applies its offsets to its real time alone, and its monotonic time advances with its ticks.

`Clock` is a type of this package and not a clock of the assert library. `assert.Clock` reports the time and nothing else. `assert.Controlled` moves only forwards: its `Advance` treats a negative duration as 0, and its `Sleep` waits for another goroutine to advance the clock. A clock of a process drifts, jumps backwards and strobes between two offsets. Every machine of a run steps on one goroutine. On that goroutine a `Sleep` of `assert.Controlled` would wait forever. The assert standard supplies no simulation runtime. It has no network model and no fault library, so the faults of a clock belong to this runtime. Other simulators keep them in their own runtime too. The clock nemesis of Jepsen bumps, strobes and resets a clock, and the `TimeSim` of TigerBeetle has linear, periodic, step and non-ideal offsets. A second consumer of such a clock would justify a request to the assert library. One example is a check of conform for the behaviour of a component after a backward jump of its clock.

An invariant or a convergence must not wait. The time of a run advances only through the world. An assertion that waits on a `Check` waits on the clock of the case, and that clock is `assert.System` unless the test has a clock of its own.

### Epochs

Each queued event, each effect in flight and each pending wait has the epoch of the process that it belongs to. A restart raises the epoch of the process. The world drops an event of an earlier epoch before its delivery and records the drop in the log. The world therefore never delivers a result of an effect that the old machine started to the new machine, also when the new machine has issued a token of the same number.

### Performing an effect

The generator writes one switch over the closed set of effects. For an effect with `call=`, the generated performer calls the dependency through the double that `stub` writes for it and builds the event of the result. The system function spawns each process with its doubles, and the performer reads the double through `Dependency`. The performer binds the call by name and by role:

- Each parameter of the callable takes the field of the effect with its name in PascalCase, such as the field `Key` for the parameter `key`. A parameter in the role context takes the context of the world.
- The event of `result` takes the token of the effect in its field of the token type, the error of the call in its field of the type `error`, and the value of the call in its field of the type of that value.
- A directive that names both `func` and `call`, or neither, or `call` without `result`, refuses the effect with `BadPerform`. So does a parameter without its field, and a result without its field in the event.

For an effect with `func=`, it calls the function of the consumer. For the reply effect, the world completes the call of its token in the history: `OK` with the value when the error is nil, `Unknown` with an error that `Anatomy.Unknown` reports, and `Fail` with any other error.

```go
// Outcome is the result of performing one effect.
type Outcome[E any] struct {
	// Event is the event of the result, and the zero E for an effect without
	// a result, such as a reply.
	Event E
	// Next returns the next event of a wait without blocking. ok reports
	// whether the wait has a next event, and done whether the wait has ended.
	// Next is nil for an effect that does not wait.
	Next func() (ev E, ok, done bool)
	// Message is the inbound call that the effect sends to another process,
	// whose reply is the result of the effect, and nil for any other effect.
	Message *Message[E]
}

// Message is an inbound call that one process sends to another.
type Message[E any] struct {
	// To is the name of the process that takes the call.
	To string
	// Call returns the event of the call with the reply token reply.
	Call func(reply uint64) E
	// Result returns the event of the result of the effect from the value
	// and the error of the reply.
	Result func(value any, err error) E
}

// Perform performs the effect f of the process p.
type Perform[E, F any] func(w *World[E, F], p *Process[E, F], f F) Outcome[E]
```

A wait is an effect whose result is a stream of events, such as a watch or a campaign for a term. After every step, the world calls `Next` of each pending wait in the order in which the waits started, and queues the event of each wait that has one. A performer must not start a goroutine, so `Next` must not block. The doubles of the world run on the goroutine of the case, and a double that returns a channel is read through a `select` with a `default` case. A crash or a stop of the process ends its waits.

### Messages between processes

A process calls another process through an effect whose performer returns a `Message`, such as the callback of a controller to an agent. The world queues the call for the receiver with a new reply token, in the current epoch of the receiver. The reply of the receiver with that token returns to the sender as the event that `Result` builds, in the epoch of the sender at the send. The history does not record a message, because a message is not a call of a client.

The world applies the rules of a runtime to each message:

- It bounds the call with the `Timeout` of the effect, on the clock of the sender. Once the timeout has passed without a reply, it queues `Result(nil, context.DeadlineExceeded)` for the sender and `Cancelled` for the receiver, and it ignores a later reply.
- A crash of the receiver queues `Result(nil, ErrCrashed)` for the sender of each call that the receiver has not replied to.
- A crash of the sender queues `Cancelled` for the receiver of each call that the sender is waiting for.

A message is in flight from its send until the receiver steps on it. The reply of the receiver returns to the sender at once, unless a partition cuts the two processes apart.

### The anatomy of a system

The world reads the tokens, the replies and the lifecycle of the machines through an anatomy, which the generator writes from the directives `sim`, `token`, `reply` and `perform`:

```go
// Anatomy is how a world reads the events and the effects of its machines.
type Anatomy[E, F any] struct {
	// Started returns the first event of a run, the inbound call of the start
	// of a process with the reply token reply.
	Started func(reply uint64) E
	// Stopped returns the last event of a run.
	Stopped func() E
	// Cancelled returns the end of the wait of the caller of the inbound call
	// with the reply token reply.
	Cancelled func(reply uint64) E
	// Call returns the reply token of an inbound call, and false for any
	// other event.
	Call func(ev E) (reply uint64, ok bool)
	// Reply returns the reply token, the value and the error of a reply, and
	// false for any other effect.
	Reply func(f F) (to uint64, value any, err error, ok bool)
	// Unknown reports whether the error err of a reply leaves the outcome of
	// its call unknown, such as the error of a stop that can follow the effect
	// of the call. It is nil for machines whose every error of a reply fails
	// its call.
	Unknown func(err error) bool
	// Token returns the token of an effect that starts a call or a wait, and
	// false for any other effect, such as an unwatch or a reply.
	Token func(f F) (token uint64, ok bool)
	// Dependency returns the name of the dependency that performs the effect
	// f, such as store.Writer, and false for an effect that a function of the
	// consumer performs.
	Dependency func(f F) (name string, ok bool)
	// Fail returns the event of the result of the effect f that failed with
	// err, and false for an effect without a result.
	Fail func(f F, err error) (ev E, ok bool)
	// Errors returns the errors that the directive fault of the dependency
	// declares for the call that performs f.
	Errors func(f F) []error
	// Timeout returns the timeout of the effect f, and 0 for an effect whose
	// call the run alone bounds.
	Timeout func(f F) time.Duration
	// Perform performs each effect other than a reply.
	Perform Perform[E, F]
}
```

The world reads a reply through `Reply` to complete its call in the history, and the error of a reply through `Unknown`. It ends the calls of a crashed process through `Call` and `Cancelled`, and it fails the call of a dependency through `Fail`. The judges of the rules of a runtime read the same functions. The generator writes `Fail` from the event of `result` of the directive `perform`, whose field of the type error takes the error, and `Timeout` from the field of the effect that the param `timeout` of `perform` names.

### Calls of an unknown outcome

A call ends as `Unknown` in the history when its outcome cannot be known:

- The runtime cancels it, because its caller stopped waiting, and the world delivers `Cancelled`.
- The process of the call crashes before it replies.
- The run ends while the call is open.
- The machine replies with an error that the param `unknown` of the directive `reply` names. The generator writes `Anatomy.Unknown` from these errors through `errors.Is`. Any other error of a reply fails its call, because a failed call has no effect in the check of the history.

```go
// Cancel ends the wait of client for its open call: the history records the
// call as unknown with ErrCancelled, and the process of the call takes
// Cancelled. It reports whether the client has an open call. The history
// allows a client one open call at a time.
func (w *World[E, F]) Cancel(client int) bool
```

The run cancels a call through the action `cancel`, which is enabled while a client has an open call and the anatomy has `Cancelled`.

```go
// ErrRunEnded is the error of a call that the stop of its process leaves
// without a reply.
var ErrRunEnded = errors.New("sim: the run ended before the call returned")

// ErrCrashed is the error of a call whose process crashed before it replied.
var ErrCrashed = errors.New("sim: the process of the call crashed before it replied")

// ErrCancelled is the error of a call whose client stopped waiting for it.
var ErrCancelled = errors.New("sim: the client stopped waiting for the call")

// ErrLost is the error of an effect that a fault of its dependency loses,
// for a dependency whose directive fault does not declare an error.
var ErrLost = errors.New("sim: a fault of the dependency lost the call")
```

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

The shapes of a partition are those of the nemesis of Jepsen. `isolate` cuts one process off, `halves` splits the ring of the processes in two at a point that the case chooses, `majorities-ring` gives each process a majority that overlaps its neighbours, and `bridge` joins two halves through one process. The faults of the breakable medium strike a dependency that persists through `medium.FS`. A forged message is an action of the consumer that calls `Deliver` with an event that no effect produced.

Each fault is a method of the world:

```go
// Partition cuts the processes into sides. Two processes exchange messages
// while one side contains both, so the sides of majorities-ring and bridge
// overlap. A message across the cut does not arrive.
func (w *World[E, F]) Partition(sides ...[]string)

// Repair ends the partition.
func (w *World[E, F]) Repair()

// Sides returns the sides of the partition in force, and nil without one.
func (w *World[E, F]) Sides() [][]string

// Drop drops the oldest message in flight to the process named to, and
// reports whether a message was in flight.
func (w *World[E, F]) Drop(to string) bool

// Delay moves the oldest message in flight to the process named to behind
// the next event of the process, and reports whether such a message was in
// flight.
func (w *World[E, F]) Delay(to string) bool

// Duplicate queues the oldest message in flight to the process named to once
// more, as a call with a new reply token. The sender takes the first reply of
// the two calls. Duplicate reports whether a message was in flight.
func (w *World[E, F]) Duplicate(to string) bool

// LoseBefore arms a loss of the next effect of the dependency named
// dependency, before the dependency applies it.
func (w *World[E, F]) LoseBefore(dependency string)

// LoseAfter arms a loss of the next effect of the dependency named
// dependency, after the dependency applies it.
func (w *World[E, F]) LoseAfter(dependency string)

// Outage loses each effect of the dependency named dependency before the
// dependency applies it, until Restore.
func (w *World[E, F]) Outage(dependency string)

// Slow keeps each effect of the dependency named dependency in flight until
// Restore.
func (w *World[E, F]) Slow(dependency string)

// Restore ends the outage and the slowness of the dependency named
// dependency, and performs each effect that its slowness keeps in flight.
func (w *World[E, F]) Restore(dependency string)

// Jump moves the clock of the process named name by d, backwards for a
// negative d.
func (w *World[E, F]) Jump(name string, d time.Duration)

// Strobe moves the clock of the process named name between the offsets a
// and b at each of its next steps steps.
func (w *World[E, F]) Strobe(name string, a, b time.Duration, steps int)

// Recover ends each fault in force and restarts each crashed process.
func (w *World[E, F]) Recover()

// Faults returns the faults in force, each as its kind and its target,
// such as pause:store, or as its kind alone for a partition.
func (w *World[E, F]) Faults() []string
```

A fault method fails the case for a target that the fault cannot strike, such as a crash of a crashed process, a pause of a paused one, a second partition, a repair without a partition, a second outage of a dependency and a restore of a dependency without a fault. `Drop`, `Delay` and `Duplicate` report false for a process without a message in flight instead.

The world strikes the faults of a dependency at the effects that the dependency performs, by the anatomy, with an error that the case chooses among the errors of the directive `fault` of the dependency. A dependency without such errors fails with `ErrLost`:

- `lose-before` does not perform the effect, and queues the event that `Fail` returns.
- `lose-after` performs the effect, so the dependency applies the call, and queues the event of `Fail` in place of the result.
- `outage` strikes `lose-before` at every effect of the dependency until the outage ends.
- `slow` keeps each effect of the dependency in flight, and queues the event of `Fail` with `context.DeadlineExceeded` once its timeout has passed on the clock of its process. An effect without a timeout remains in flight until the fault ends.

The `Injector` of the calls that fail on demand fails a call before its double applies it, and its `Hold` blocks the goroutine of the caller. The world needs a call that fails after the dependency applies it, and it runs every machine on the goroutine of the case, so it fails the effects itself instead of through the injector.

A fault is in force from its strike to its end:

| Fault | In force until |
|---|---|
| `crash` | The restart of the process |
| `pause` | The resume of the process |
| `partition` | The repair |
| `outage` and `slow` | The restore of the dependency |
| `lose-before` and `lose-after` | The effect that the armed loss strikes |
| `clock-strobe` | The last of its steps |

`drop`, `delay`, `duplicate`, `clock-jump`, `restart` and `repair` strike once and are never in force. Each entry of the log records the faults in force at its step, and the world passes the kind of each fault that strikes to `Case.Classify`.

A system that tolerates one crash fails convergence trivially under two, so the run keeps the faults in force at once within the param `tolerate`. An action of a fault that would exceed it is not enabled in that state. An action of a fault that is never in force is enabled while `tolerate` names its kind. A system without `tolerate` does not tolerate a fault, and the run checks it without faults.

The case chooses the target of each fault:

- a running process for `crash` and `pause`
- a process with a message in flight for `drop`, `delay` and `duplicate`
- a dependency that the world has seen for the faults of a dependency

A clock jumps by an amount from −1 hour to +1 hour. A strobe has two such offsets and lasts 1 to 16 steps. Each of these choices shrinks towards 0.

### The run

The module runs a simulation as a machine of `stateful` over a new world in each case. The generator writes the simulation of a declaration as a value:

```go
// Simulation is a system of sans-IO machines with its oracles. S is the
// state of the spec of the calls of its clients.
type Simulation[E, F, S any] struct {
	// Anatomy is how the world reads the events and the effects.
	Anatomy Anatomy[E, F]
	// System builds the processes of each world and their dependencies.
	System System[E, F]
	// Spec is the sequential specification of the calls of the clients.
	Spec history.Spec[S]
	// Actions are the actions of the consumer. A final action sets Drain.
	Actions []Action[E, F]
	// Rules are the checks of the claims of sim that bind to the machines.
	Rules []Rule
	// Clients is the number of clients that make calls, and 3 for 0.
	Clients int
	// Invariants, Converges and Scenarios are the oracles of the consumer.
	Invariants []Invariant[E, F]
	Converges  []Convergence[E, F]
	Scenarios  []Scenario[E, F]
	// Fingerprint returns a hash of the state of the system. A nil
	// Fingerprint observes the checksum of the log alone.
	Fingerprint func(w *World[E, F]) uint64
	// Tolerate bounds the faults in force at once, by their kind.
	Tolerate map[Fault]int
	// Require is the share of the cases in which each kind of fault strikes.
	Require map[Fault]float64
	// Strategy chooses the next process. The zero Strategy is PCT(3).
	Strategy Strategy
	// Isolation is the isolation level of the transactions of the clients,
	// and the zero Isolation checks none.
	Isolation Isolation
}

// Run runs one case of s on a new world.
func Run[E, F, S any](c *prop.Case, s Simulation[E, F, S])

// Test runs the property of s on tb over random runs, with the shares of the
// faults that s requires, and then the crash points of each scenario.
func Test[E, F, S any](tb assert.TB, s Simulation[E, F, S])

// Falsify proves each rule and each invariant of s on its planted defect, in
// a parallel test of t named after the rule or the invariant, as Prove proves
// it.
func Falsify[E, F, S any](t *testing.T, s Simulation[E, F, S])

// Prove proves the check named check on tb with its planted defect: the rule
// whose identity is check on its defect form, or the invariant whose
// statement is check on its defect system.
func (s Simulation[E, F, S]) Prove(tb assert.TB, check string)

// Fuzz registers the choices of a run of s as the input of the fuzz target f.
func Fuzz[E, F, S any](f *testing.F, s Simulation[E, F, S])

// Invariant is a safety property that the run checks after every step, with
// the system of its planted defect.
type Invariant[E, F any] struct {
	Statement string
	Check     func(h *Check, w *World[E, F])
	Defect    System[E, F]
}

// Convergence is a liveness property that the run checks once the faults
// have ended and Within has passed.
type Convergence[E, F any] struct {
	Statement string
	Within    time.Duration
	Check     func(h *Check, w *World[E, F])
}

// Scenario is a fixed run, whose events the crash points cover.
type Scenario[E, F any] struct {
	Name string
	Run  func(w *World[E, F])
}

// Action is an action of the consumer: a change that starts outside the
// machines, or a call of a client.
type Action[E, F any] struct {
	// Name names the action in a counterexample.
	Name string
	// Enabled reports whether a step may take the action, and is nil for an
	// action that every step may take.
	Enabled func(w *World[E, F]) bool
	// Input requests the input of a step from the case, and is nil for an
	// action without an input.
	Input func(c *prop.Case) any
	// Run runs the action on the world. client is a client without an open
	// call for an action with Call set, and 0 otherwise.
	Run func(w *World[E, F], client int, input any)
	// Call marks a call of a client. The run enables the action while a
	// client has no open call.
	Call bool
	// Drain marks an action of the final phase, which the run takes once.
	Drain bool
}

// Rule is the check of a claim of sim that binds to the machines of a
// simulation, such as replies-once. The generator writes it from the claim.
type Rule struct {
	// ID is the identity of the check, such as Machine/replies-once.
	ID string
	// Statement is the sentence that the check proves of the machines.
	Statement string
	// Claim is the claim, Probe its probe, Reason the reason of its defect
	// and Form the form of its defect.
	Claim  spec.ClaimID
	Probe  spec.ProbeID
	Reason spec.ReasonID
	Form   spec.DefectFormID
}

// Fault is the kind of a fault, such as crash or partition, as the faults
// below name it.
type Fault string

// Strategy chooses the process that steps next.
type Strategy struct{ /* unexported */ }

// Uniform returns the strategy that picks a process with a pending event
// uniformly.
func Uniform() Strategy

// PCT returns the strategy of probabilistic concurrency testing of depth
// depth. A run with a depth below 1 fails its case.
func PCT(depth int) Strategy

// Isolation is the isolation level that the transactions of the clients
// have.
type Isolation uint8

const (
	// Serializable checks the history with history.Serializable.
	Serializable Isolation = iota + 1
	// SnapshotIsolation checks the history with history.HasSnapshotIsolation.
	SnapshotIsolation
)
```

`PCT` takes the priority of each process and the change points from the case with the encoding of the task scheduler of assert-go, so a shrunk run steps its processes as `Uniform` does.

The run has three phases, which take the parts of `Steps`:

| Phase | Part of `Steps` | Actions |
|---|---|---|
| Active | The sequential steps | `deliver`, `advance`, `cancel`, one action per fault within `tolerate`, the end of each fault in force, and each action of the consumer without `Drain` |
| Recovery | `Settle` | `Recover`, which ends every fault and restarts every crashed process |
| Final | `Settle` | Each `final` action once, then `deliver` and `advance` until the bound of each convergence has passed |

The recovery and the final phase run in `Settle` and not in the drain, because `Steps` offers a drain action to the sequential steps as well, whenever the swarm of the case keeps it. A recovery in the active phase would end the faults early, and a final action would run before the faults end.

The ends of the faults are `restart`, `resume`, `repair`, `restore-outage` and `restore-slow`. Each end is enabled while its fault is in force, and `restart` needs `tolerate` to name it, as each fault does. The run passes the name of each action that it takes to `Case.Classify`, so a test can require a share of the cases for an action. An action of a fault passes no name, because the world passes the kind of the fault when it strikes. An armed loss that never strikes gives its case no label, so `require` counts only the cases in which the loss strikes. `advance` advances the time to the next timer, and further by an amount up to that gap that the case chooses.

`deliver` asks the strategy for a process with a pending event, and the process steps on its oldest event of its current epoch. `uniform` picks a process uniformly. `pct:<depth>` draws a priority for each process and `depth-1` change points from the case, and picks the pending process of the highest priority, lowering it at each change point, as the strategy `PCT` of the stateful scheduler of assert-go orders its tasks. A bug that needs a certain order of `depth` events has a probability of discovery that PCT bounds from below, where a uniform pick gives none.

After every step, `Invariant` checks:

- the five rules of a runtime over the log
- each invariant of the consumer

In `Settle`, once the final phase ends, the run checks each convergence and the whole history: `history.Linearizable` against the spec, partitioned by the keys of each call, and `history.Serializable` or `history.HasSnapshotIsolation` when the directive states `isolation`. The check of the whole history covers the calls that complete in the final phase, after the last step of `Steps`.

The machine of `stateful` has no `Spec`, so an action of the consumer takes no state of the spec. `stateful` checks its `Spec` after every step with every call in one partition. A crash or a cancel ends a call as unknown, and an unknown call is concurrent with every later call. In a fixture with a primary and a backup, 12 crashes and 12 cancels in one case left 23 unknown calls in one partition, and that check stopped undecided at its budget of 10,000,000 steps. The check in `Settle` searches each key apart.

After every step the run also calls `Case.Classify` with the kind of each fault that struck and of each effect that the world performed, and `Case.Observe` with the fingerprint of `//dokimi:fingerprint`. `require` turns the labels of the faults into `prop.Require`, so a suite whose cases seldom crash a process fails instead of passing on calm runs.

For each `//dokimi:scenario`, the generated test runs the scenario once to count its events. It then runs it again for each event, crashing the process of the event once before the event and once after it, and restarting it. The invariants and the rules of the runtime run after every step of each run, and the check of the history runs at its end.

### Determinism

The world hashes each entry of its log into a checksum: the process, the epoch, the typed literal of the event and of each effect, and the time of the clock. After each step the run passes the checksum to `Case.Observe`. The engine replays every failing case with the same choices and compares its fingerprints, and `Test` passes the option `prop.Confirm(16)`, which replays every 16th case of a passing run the same way. A difference fails the test with the first step at which the two runs differ, and with the causes that the world cannot see:

- iteration over a Go map in a machine or a performer
- the global source of `math/rand`
- a read of the wall clock
- a goroutine in a performer

A machine whose effects depend on the order of a map therefore fails at once, instead of producing a seed that does not replay.

### The rules of a runtime

The spec states five claims with the generator `sim` and the machinery `simulated-system`. A claim of `sim` binds to each machine interface under the directive, as a claim of `enum` binds to each enumeration:

| Claim | Law | Reason | Defect form |
|---|---|---|---|
| `replies-once` | A machine returns at most one reply to each inbound call | `second-reply` | `replies-twice` |
| `replies-before-stop` | On the stop event, a machine replies to each inbound call that it has not replied to | `call-without-reply` | `forgets-a-call-on-stop` |
| `quiet-after-stop` | After the stop event, a machine returns replies alone | `effect-after-stop` | `writes-after-stop` |
| `quiet-after-error` | A machine returns no effect after the step in which its error method returns an error | `effect-after-error` | `steps-after-error` |
| `fresh-tokens` | Each effect that starts a call or a wait has a token that no earlier effect of its machine in the same epoch has | `reused-token` | `reuses-a-token` |

Each probe is a judge over the log of the world, which reads the replies, the calls and the tokens through the anatomy and the error of each step from its entry. Each defect form wraps a fresh machine and changes the effects of its step, through the anatomy alone, so the module writes the five forms for every system. `Plant` wraps a machine with a form:

```go
// Plant returns the machine m with the planted defect of the form, which
// reads the effects of m through the anatomy a.
func Plant[E, F any](m Machine[E, F], a Anatomy[E, F], form spec.DefectFormID) Machine[E, F]
```

The forms are:

- `replies-twice` returns each reply of a step twice.
- `forgets-a-call-on-stop` drops the replies of the stop event.
- `writes-after-stop` returns, on the stop event, the last effect with a token that the machine returned before it.
- `steps-after-error` makes the error method return an error from the second step on, and keeps the effects of every step.
- `reuses-a-token` returns, in each step that starts a call or a wait, the first such effect of the machine once more.

Falsification runs each claim on its defect at a seed from the key of the check, as every proof of a property check runs. The run proves a rule when its first failure is the failure of the judge of the probe of the rule, whose message is the statement of the rule.

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
| `TestSimulation` | The property over random runs, the crash points of each scenario, and the determinism check, through `sim.Test` |
| `TestSimulationFalsify` | Each rule of the runtime on its defect form, and each invariant on its defect system, through `sim.Falsify`, which proves each check through `Simulation.Prove` |
| `FuzzSimulation` | The choices of the run as the input of a fuzz target, through `sim.Fuzz` and `prop.Fuzz` of assert-go |

Each test is one call of the module with the simulation of the declaration, so the generated code contains the anatomy, the performer and the simulation alone.

A failing run reports the shrunk log with the step of each process, its epoch, its clock and the faults in force, the history of the calls with the anomaly or the call that no linearization explains, and the replay token of the property engine. Each case passes its log and its history to `Case.Logf` when it ends, and the engine prints the notes of the failing case after its counterexample. The report writes each event and each effect with its type and the verb `%+v` of fmt. `DOKIMI_ASSERT_PROP_REPLAY` replays the token, and `DOKIMI_ASSERT_PROP_SEED` reruns the seed of a run.

### The capabilities of Jepsen

| Jepsen | The simulation |
|---|---|
| A generator of operations from concurrent clients | One action `call` per operation, with inputs from the case. Each open call keeps its client busy, so the clients run concurrently |
| `:ok`, `:fail` and `:info` | `OK`, `Fail` and `Unknown` in `history.History` |
| Knossos, linearizability against a model | `history.Linearizable` against the spec of `//dokimi:spec`, partitioned by key |
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
| The clock nemesis of Jepsen | <https://github.com/jepsen-io/jepsen/blob/main/jepsen/src/jepsen/nemesis/time.clj> |
| Knossos, the linearizability checker of Jepsen | <https://github.com/jepsen-io/knossos> |
| Elle, inferring isolation anomalies from histories | <https://arxiv.org/abs/2003.10554> |
| assert-go, the packages history and stateful | <https://github.com/dokimasia/assert-go/tree/main/history> |
| gosim, a deterministic simulation of Go programs by translation | <https://github.com/neria-cloud/gosim> |
| TigerBeetle, the deterministic simulator VOPR | <https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/internals/vopr.md> |
| TigerBeetle, the simulated time `TimeSim` | <https://github.com/tigerbeetle/tigerbeetle/blob/main/src/stdx/time.zig> |
| assert-spec, the clock | <https://github.com/dokimasia/assert-spec/blob/main/docs/rfc/0006-the-clock.md> |
| assert-spec, machines, simulation and campaigns | <https://github.com/dokimasia/assert-spec/blob/main/docs/rfc/0012-machines-simulation-and-campaigns.md> |
| eidos, the request for the terms of a union of a constraint | <https://github.com/dokimasia/eidos/issues/112> |
| eidos, the request for a method of a type as a callable in scope | <https://github.com/dokimasia/eidos/issues/114> |
| assert-go, the search of `Linearizable` over calls of unknown outcome | <https://github.com/dokimasia/assert-go/issues/89> |
| Gavin Lowe, testing for linearizability with just-in-time linearization | <http://www.cs.ox.ac.uk/people/gavin.lowe/LinearizabiltyTesting/paper.pdf> |
| turmoil, a simulation of hosts and a network for Rust | <https://github.com/tokio-rs/turmoil> |
| madsim, a deterministic simulator of a Rust runtime | <https://github.com/madsim-rs/madsim> |
| assert-spec, the workload of Elle on read-write registers | <https://github.com/dokimasia/assert-spec/issues/27> |
