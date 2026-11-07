---
rfc: 0036
title: A simulation of components that run goroutines
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-05
updated: 2026-11-07
discussion: https://github.com/dokimasia/dokimi/issues/58
supersedes: RFC-0025 in part
superseded-by: ADR-0023 in part
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0036: A simulation of components that run goroutines

## Summary

A storage engine, a cache with a background flusher and an RPC service run their work on goroutines of their own. The simulation of `go.dokimi.dev/sim` runs sans-IO machines alone, so it cannot run such a component. This RFC adds a simulated host and runs components on it:

- **The host.** A new module, `go.dokimi.dev/host`, simulates the machine of a component: the incarnations of its process, its disk, a wall clock and a monotonic clock, its network and its dependencies. A world crashes, pauses and restarts an incarnation as a whole, and fences every resource of an incarnation at its end. conform and sim share the host.
- **Two execution modes for components.** Native mode runs the code of a component on goroutines in a `testing/synctest` bubble. It replays the choices of a case, and it does not replay the order of the goroutines. Tasks mode runs code that blocks only through the task scheduler of `stateful`, and replays the whole case. The machines of sim remain a third mode, which replays whole.
- **The world of a subject.** A new family of derived checks runs the subject of a declaration on one host. Concurrent clients make the derived calls of the declaration. Between their calls the world crashes the subject and opens it again, strikes faults of its disk, its dependencies and its clock, and advances time. The claims `linearizable-under-faults` and `completes-after-recovery` judge the history of the calls against the reference of the declaration.
- **Components in a system.** A process of the world of a system can run a Go component on a host of its own. A network of streams and datagrams connects the hosts, with the partitions and the link faults of the world.

The run checks the history of the calls after every step, partition by partition. The waits of tasks and the replays of a failing case are changes of the definition of assert-spec, which the RFC requires.

## Motivation

### A component with goroutines has no simulation

[#58](https://github.com/dokimasia/dokimi/issues/58) asks for a simulation of a Go component that starts goroutines. A storage engine runs its flushes, its compactions and its schedulers on goroutines of its own, and writes through a file system. Each consumer of such a component writes a stateful machine with a double of its file system. That harness repeats the choice of each fault, the crash points, the restart of a crashed process and the proof of each invariant.

The simulation of sans-IO machines runs every machine on the goroutine of the case. It counts a goroutine in a performer as a cause of a replay that differs, so it cannot run such a component.

### The parts exist and no check combines them

| Part | Where | What it supplies |
|---|---|---|
| The door `Durable` | conform | Opens a subject over a `medium.FS`, and opens it again over what a crash of the medium left |
| `Handle.Injector` | conform | Fails, blocks and releases the calls of the dependencies of a subject on demand |
| The door `Reference` | conform | A reference implementation, such as a fake that the generator `fake` writes |
| `probe.LinearizableHistory` | conform | Runs the derived calls of a declaration on concurrent clients, and checks the history against the reference through `history.SpecFrom` |
| `Check.Around` and `prop.Around` | conform and assert-go | Run each run of a check in an environment, such as a `testing/synctest` bubble |
| `medium.Memory` | dokimi-go | Crashes at a persistence point, tears a write, flips a bit, runs out of space, fails a call, and blocks its calls until a release |
| `history`, `stateful` and `prop` | assert-go | The history of calls with unknown outcomes and its incremental check, machines whose steps a case chooses, the task scheduler, and the replay of a case |

For the mixin `durable`, the crash-point probe runs sequences of one to three writes on one goroutine and crashes each of them. The linearizability probe runs concurrent clients, and it does not crash the subject or strike a fault. Both probes run their calls without an advance of the clock, so a flush that a timer starts does not run.

### The history is checked once per case

The run of a simulation is a machine of `stateful` with the state `struct{}` (`sim/run.go:151` of dokimi-go). It checks the history for linearizability once, after the final phase, partition by key (`run.go:540`). So a violation shows only at the end of its case, and each candidate of the shrinker runs every step of the case, also the steps after the violation. `history.Linearizable` of assert-go continues the search of each partition of a growing history through `Resume`, so a check after every step costs the steps since the check before it. A partition in which a call that the check before it read as pending has completed searches again from its first call.

### What Go lets a harness control

A harness controls what a component takes by injection: its file system, its transport and its clock. It does not control the order of goroutines. On Linux the Go runtime seeds its randomness from the auxiliary vector of the kernel (`_AT_RANDOM`), preempts a goroutine after 10 ms of real time (`forcePreemptNS`), and runs goroutines on up to `GOMAXPROCS` threads at once. A probe on Go 1.27.3 ran six goroutines in a bubble, and each of three runs gave a different order.

A bubble also has a limit on what it can wait for:

- **A goroutine that waits for a mutex is not durably blocked.** The documentation of `testing/synctest` lists a locked `sync.Mutex` and `sync.RWMutex` among the waits that are not durable (`testing/synctest/synctest.go:91` of Go 1.27.3). The table of idle waits of the runtime leaves the wait for a mutex out (`runtime/runtime2.go:1386`), and a bubble advances time or returns from `synctest.Wait` only while no goroutine runs (`runtime/synctest.go:133`).
- **A probe on Go 1.27.3 hung in both of its forms.** A goroutine sleeps for one second with a mutex locked while the root goroutine waits for the mutex: the clock of the bubble never advances. A goroutine waits on a channel with a mutex locked while another goroutine waits for the mutex: `synctest.Wait` never returns. Both tests ran into the timeout of 15 s.
- **The limit is known upstream.** golang/go#77687 tracks durable waits for a mutex, with `database/sql` and the Kafka client franz-go as cases. The Go team declined the change in the first design of `synctest`, because the mutex is one of the most performance-critical parts of the standard library.

So a bubble cannot become quiescent while one goroutine waits for a mutex and the goroutine that locked it waits durably: on a sleep, on a channel, on a disk that a fault blocks, or on the reply of a simulated network.

Two more limits shape the design:

- **One reading of time.** `t.Add` adds the same duration to the wall reading and to the monotonic reading of a `time.Time`, and `t.Round(0)` strips the monotonic reading (`time.go:40` to `time.go:46`). A clock that returns the time of the bubble plus an offset makes a jump of the wall clock also a jump of every duration that `Sub` and `Since` measure, which Go code in production never sees.
- **A task that yields remains ready.** `Scheduler.Run` of `stateful` releases tasks while one is ready, and `Scheduler.Yield` makes the running task ready again (`scheduler.go:148` to `scheduler.go:206` of assert-go). A task that waits for a condition by yielding in a loop keeps `Run` from returning. Under the strategy `uniform` each release records a choice, so such a run ends as rejected once it has made the 8,192 choices of a case. Under `pct` a release does not record a choice, so the run never ends.

### Simulators scope every resource to a host

Every simulator of systems that runs real code makes the machine of a component the unit of a crash, and controls the blocking of its code:

| Simulator | Unit of a crash | Time | Blocking that it controls |
|---|---|---|---|
| FoundationDB | A process or a machine, through a kill type: an instant kill, a reboot, a reboot that deletes the data, a failed disk | Virtual. Actors wait on futures and on `delay` | Everything, through its actor runtime |
| madsim | A node. A kill drops every task of the node and loses the data that was not flushed. A restart runs the initial task again | Moves to the next event once no task is ready | Its own runtime of Tokio. Blocking code panics |
| turmoil | A host. A crash loses the writes that a sync did not cover, and a bounce restarts the software | A fixed tick per step | Its own runtime |
| gosim | A machine. A stop ends its goroutines and its timers, and a restart gives it new globals. A crash keeps a random subset of the pending writes | Moves to the next timer | Everything, through a translation of the source |
| TigerBeetle VOPR | A replica. A crash can corrupt the sectors of its writes in flight | A wall clock with linear, periodic and step offsets against a monotonic clock | Its single-threaded design |

Go has no runtime that a library can replace. This RFC keeps the code that production runs. So the host controls everything that a component takes through it, and each mode states the guarantee that its control supports.

## Detailed design

### Components

| Component | Where | Responsibility |
|---|---|---|
| The host | `go.dokimi.dev/host`, a new module | Hosts, incarnations and kills, the two clocks, the network, the dependencies, the fences, and the accounting of goroutines |
| The faults of a medium | `go.dokimi.dev/medium` | A call that fails before, after or in the middle of its effect, a hold of each call, the kill of a process, and a crash that keeps a part of the pending writes |
| The world of a subject | `go.dokimi.dev/conform/probe` | `World`, `Operation` and `Pool`, the steps of a world, its checks, and the probes of the two claims |
| The door of a hosted subject | `go.dokimi.dev/conform` | The option `Hosted`, the modes `Tasks` and `Blocking`, the environment `Bubble` and its `RecordStep`, the methods of the handle that a world calls, and the injector of the host behind the names of conform |
| The claims of a world | The runtime contract | The machinery `quiescence`, the two claims with their laws, probes, reasons and defect forms, and the field `world` of the suite record |
| The mixin `creates` | The shape catalog of eidos | A callable that returns the key of the record that it created |
| The rules of the two claims | dokimi `rules/` | Their licences and their gates |
| The suite generator | dokimi | The operations and the pools of a declaration, and the checks of the two claims |
| Components in a system | `go.dokimi.dev/sim` | Component processes on hosts, the calls of clients on them, and the faults of the network |
| The sim generator | dokimi | The simulation of a system without a machine interface |
| The waits of tasks | assert-spec, then the `stateful` of each language | `wait`, `clock`, `advance` and `max-releases` of the task scheduler |
| The replays of a failing case | assert-spec, then the `prop` of each language | The option `replays` and the count of the replays that reproduce a failure |

### Execution modes

| Mode | Code that runs | Quiescence | Replay | Faults |
|---|---|---|---|---|
| Machines | The sans-IO machines of sim | The run applies one event at a time | The whole case | Every fault of sim |
| Tasks | Code that blocks only through the host under the task scheduler | `Scheduler.Run` and `Scheduler.Advance` return once no task is ready | The whole case | Every fault of the host |
| Native | Any Go code on goroutines in a bubble | `synctest.Wait` returns | The choices of the case alone. A failing case that a replay does not reproduce ends as flaky | The instant faults, and the blocking faults for a subject that states `Blocking` |

A fault is instant when it changes the outcome of a call without making a goroutine wait for the world: a call that fails, a torn or lost write, a full disk, a flipped bit, a crash, a dropped datagram, a reset stream, and a jump or a strobe of the wall clock. A fault is blocking when it makes a goroutine wait for a later step of the world: a hold of the disk or of a dependency, the latency of a link, the bytes of a stream that wait for the repair of a partition, and the pause of a process.

The suite record states the mode of each world check.

### The host

```go
// Package host simulates the machine of a component: the incarnations of its
// process, its disk, its clocks, its network and its dependencies. The world of
// a subject in conform runs one host, and a world of sim runs a host for each
// component process.
package host

// Host is a simulated machine. It keeps its disk, its clocks and its address
// across the incarnations of its process.
type Host struct{ /* unexported */ }

// New returns the host named name on the network n, whose monotonic clock runs
// at a rate error that the case c chooses. With the scheduler s, the host runs
// its goroutines, its sleeps and its timers as tasks of s. With a nil s, it
// runs them on the goroutines of the bubble.
func New(c *prop.Case, name string, n *Network, s *stateful.Scheduler) *Host

// Name returns the name of the host, with which each address of the host
// starts.
func (h *Host) Name() string

// Disk returns the disk of the host, which every incarnation of its process
// opens.
func (h *Host) Disk() *medium.Memory

// Injector returns the injector of the dependencies of the host.
func (h *Host) Injector() *Injector

// Start starts a new incarnation of the process of the host, and returns it
// with the environment that its code receives.
func (h *Host) Start() (*Incarnation, *Env)

// Broadcast wakes each wait of an incarnation of the host on Env.Wait, which
// calls its condition again. A world calls it on each host after each of its
// steps. Under the scheduler it does nothing, because the scheduler checks its
// parked tasks at the start of each Run and each Advance.
func (h *Host) Broadcast()

// Incarnation is one run of the process of a host, from its start to its
// crash. Only the world has it.
type Incarnation struct{ /* unexported */ }

// Go runs f on a goroutine of the incarnation, or as a task under the
// scheduler of the host. The goroutine and every goroutine that it starts
// have the labels of the incarnation.
func (i *Incarnation) Go(f func())

// Crash ends the incarnation by the kind k: it fences every resource of the
// incarnation, ends its lifetime and applies k to the disk.
func (i *Incarnation) Crash(k Kill)

// Pause blocks every call of the environment of the incarnation, and defers
// its timers, until Resume. The clocks of the host keep running.
func (i *Incarnation) Pause()

// Resume ends a pause.
func (i *Incarnation) Resume()

// Leaked returns the stacks of the goroutines of the incarnation that are
// alive after its crash.
func (i *Incarnation) Leaked() []string

// Kill is a kind of crash.
type Kill uint8

const (
	// ProcessCrash ends the process and keeps the page cache, as a kill of
	// the process does: the writes that a sync did not cover remain pending
	// on the disk.
	ProcessCrash Kill = iota + 1
	// PowerLoss ends the process and loses the page cache: the disk keeps a
	// crash state of the writes that a sync did not cover.
	PowerLoss
	// DiskReplacement ends the process, and the next incarnation starts over
	// an empty disk.
	DiskReplacement
)

// Env is the environment of one incarnation. Each method returns the same
// resource for the whole incarnation.
type Env struct{ /* unexported */ }

func (e *Env) Name() string
func (e *Env) Incarnation() int
// Lifetime returns the context of the incarnation, which a crash, a stop of
// the process and the end of the check end.
func (e *Env) Lifetime() context.Context
// Disk returns the view of the incarnation on the disk of the host.
func (e *Env) Disk() medium.FS
func (e *Env) Clock() *Clock
// Rand returns a source whose every value is a choice of the case.
func (e *Env) Rand() rand.Source
// Dependency returns the dependency name of the incarnation, which the
// injector of the host fails, blocks and releases.
func (e *Env) Dependency(name string) *Dependency
// Go runs f on a goroutine, or as a task under the scheduler of the host.
func (e *Env) Go(f func())
// Yield is a scheduling point under the scheduler of the host, and calls
// runtime.Gosched without one.
func (e *Env) Yield()
// Wait blocks until cond returns true. Under the scheduler it parks the task.
// Without one it waits on a condition variable that each Yield of the
// incarnation and each Broadcast of the host signal.
func (e *Env) Wait(cond func() bool)
func (e *Env) Listen(port int) (net.Listener, error)
func (e *Env) Dial(ctx context.Context, addr string) (net.Conn, error)
func (e *Env) ListenPacket(port int) (net.PacketConn, error)

// Park blocks until ctx ends, and returns the error of ctx. For a context that
// derives from the lifetime of an incarnation, under the scheduler of its
// host, it parks the running task. Otherwise it waits for ctx.Done().
func Park(ctx context.Context) error
```

A planted defect has the context of a call and no environment, so it blocks a call through `Park`. The lifetime of each incarnation contains the waiter of its host. The context of each call of the world derives from that lifetime, so `Park` finds the scheduler of the host in it.

A crash of an incarnation runs in this order:

1. It fences the incarnation. Each later call of the incarnation on its view of the disk, on a dependency of `Env.Dependency`, on a connection and on a listener fails with an error that wraps `medium.ErrCrashed`.
2. It ends the lifetime of the incarnation.
3. It applies the kind to the disk:
   - `ProcessCrash` calls `medium.Memory.Kill`, which closes each file of the process and keeps every write.
   - `PowerLoss` calls `medium.Memory.CrashKeeping` with the counts of pending writes that the case chooses, as the section on the medium describes.
   - `DiskReplacement` gives the host a new empty disk.
4. It resets each connection of the incarnation and closes each of its listeners.
5. It ends each hold of its calls: each call that a hold of the disk or of a dependency blocks fails with `medium.ErrCrashed`.

The world then waits for quiescence and reads `Incarnation.Leaked`. The end of a case crashes a running incarnation by `ProcessCrash`, so the world reads `Leaked` for every incarnation of the case. A goroutine inherits the profiler labels of the goroutine that starts it (`runtime/pprof/runtime.go:39`, `runtime/proc.go:5403`). The world opens the subject and makes each call through `Incarnation.Go`, so every goroutine of the subject has the labels `dokimi.host` and `dokimi.incarnation`. `Leaked` reads the goroutine profile, whose text form states the labels of each stack (`runtime/pprof/pprof.go:498`). A goroutine of a crashed incarnation that is alive after quiescence ignored its lifetime, and the world fails the check with its stack.

The host runs a timer for the arrival of a datagram or of the bytes of a stream, for a deadline of a connection, and under the scheduler for each sleep. Such a timer can still run after the crash of the incarnation whose goroutine started it. So each timer takes the profiler labels of the context of its wait, and `Leaked` does not count it. That context carries no labels of an incarnation, because `Incarnation.Go` sets them on its goroutines alone.

`conform.Injector`, `conform.Dependency`, `conform.ErrInjected` and `conform.FaultOf` are aliases of the names of the host. Every consumer of them compiles unchanged. A dependency that `Env.Dependency` returns belongs to its incarnation. The fence of step 1 ends it.

### Clocks

```go
// Clock is the clock of a host, as an incarnation reads it: a wall clock that
// the faults of the clock move, and a monotonic clock that they do not move.
type Clock struct{ /* unexported */ }

// Wall returns the wall time of the host without a monotonic reading, so a
// comparison or a difference of two values is a computation of wall time.
func (c *Clock) Wall() time.Time

// Monotonic returns the time since the start of the host on its monotonic
// clock.
func (c *Clock) Monotonic() time.Duration

// Sleep blocks until the monotonic clock has advanced by d, or until the
// lifetime of the incarnation ends.
func (c *Clock) Sleep(d time.Duration)

// AfterFunc calls f on a goroutine of the incarnation, or as a task under the
// scheduler, once the monotonic clock has advanced by d. stop cancels the
// call, and reports whether it did so before f started.
func (c *Clock) AfterFunc(d time.Duration, f func()) (stop func() bool)

// Bump moves the wall clock of the host by d.
func (h *Host) Bump(d time.Duration)

// Strobe moves the wall clock of the host back and forth by delta every period
// of its monotonic clock, until ResetClock.
func (h *Host) Strobe(delta, period time.Duration)

// ResetClock ends each bump and each strobe of the wall clock of the host.
func (h *Host) ResetClock()

// WallRead reports whether an incarnation of the host has read its wall clock.
func (h *Host) WallRead() bool
```

The wall clock of a host is its start time plus its monotonic time plus the offset of its bumps and its strobe, which are the clock faults of Jepsen. The monotonic clock of a host advances at the rate of the clock of the bubble, or of the clock of the scheduler in Tasks mode, times one plus the rate error of the host. Its rate error is a choice of the case within the range of the drift of sim. A host under the scheduler takes its monotonic time from `Scheduler.Clock`, and its sleeps and timers park their tasks on that clock until `Scheduler.Advance` moves it to their deadlines.

A world strikes the faults of a clock only once the subject has read its wall clock. A world whose subject never reads it records the faults of the clock as vacuous in the suite record. A forbidigo rule on `time.Now` and `time.Since` in the package of a subject finds code that bypasses the clock of the host. The handle of a check outside a world keeps its controlled clock.

### The network

| Part | Model |
|---|---|
| Address | The name of a host and a port, such as `node-2:7000` |
| Stream | Bytes arrive in order, after the latency of the link plus their size over its throughput. The case chooses the latency of each link once per case from the latency range of the network, which is 0 to 100 ms by default as in turmoil. It chooses the throughput from the throughput range, which has no limit by default. Both shrink to a link without delay |
| Datagram | Each datagram arrives after the latency of its link, or not at all under `drop` |
| Partition | The shapes `isolate`, `halves`, `majorities-ring` and `bridge` cut the links between their sides until `repair` |

```go
// Network connects the hosts of a world through streams and datagrams.
type Network struct{ /* unexported */ }

// NewNetwork returns a network of the case c without hosts, which delays the
// bytes of each stream for blocking.
func NewNetwork(c *prop.Case, blocking bool, opts ...NetworkOption) *Network

// NetworkOption is an option of NewNetwork.
type NetworkOption struct{ /* unexported */ }

// Latency states the latency range of a blocking network. The case chooses
// the latency of each link as turmoil chooses the latency of a message: lowest
// plus highest − lowest times a value of the exponential distribution with
// the rate 5, at most highest. The range is 0 to 100 ms without the option.
func Latency(lowest, highest time.Duration) NetworkOption

// Throughput states the lowest throughput of a link of a blocking network, in
// bytes per second. The case chooses the time of a byte on each link in the
// same way, from no time to the time of a byte at lowest. A link has no limit
// without the option.
func Throughput(lowest int) NetworkOption

// Partition cuts every link between two sides until Repair.
func (n *Network) Partition(sides ...[]string)
func (n *Network) Repair()

// Reset fails every call on both ends of each stream between the hosts a and
// b with an error that wraps ErrReset.
func (n *Network) Reset(a, b string)

// Slow raises the latency of the link between the hosts a and b to latency,
// until Restore.
func (n *Network) Slow(a, b string, latency time.Duration)
func (n *Network) Restore(a, b string)

// Drop drops the oldest datagram in flight to the host named to. Delay moves
// that datagram behind the next datagram in flight to the host, and Duplicate
// sends it once more. Each reports whether the datagrams that it moves were in
// flight.
func (n *Network) Drop(to string) bool
func (n *Network) Delay(to string) bool
func (n *Network) Duplicate(to string) bool

// InFlight returns the number of datagrams in flight to the host named to.
func (n *Network) InFlight(to string) int

// ErrReset is the error of a call on a stream that a reset ended.
var ErrReset = errors.New("host: the connection was reset")

// ErrNoNetwork is the error of a call of the network of a host without one,
// such as the one host of the world of a subject.
var ErrNoNetwork = errors.New("host: the host has no network")
```

A network built with `blocking` set gives each link the latency and the throughput of the case. It delays the bytes of each stream across a cut until the repair. A network without `blocking` delivers each byte at once and refuses `Slow`. It resets each stream across a cut and refuses each dial across it. Each read of a connection and each accept waits on a condition variable of the bubble, or parks its task under the scheduler. So a component that waits for its network waits durably. `drop`, `delay` and `duplicate` strike datagrams through `Drop`, `Delay` and `Duplicate`, as they strike the messages of machines.

### Native mode

Native mode keeps `synctest.Wait` as its quiescence. It enforces the limit of a bubble:

- **A contract for blocking faults.** A world strikes only instant faults by default. It strikes the blocking faults for a subject whose application states the option `conform.Blocking[S]()`, or for a system that states `Blocking`. The option states that no goroutine of the subject waits on a channel, on a sleep, on the disk or on the network while it has locked a mutex that another goroutine waits for.
- **A watchdog.** The environment `conform.Bubble`, which a world check states as its `Around`, runs the check in `synctest.Test` and samples the goroutines of the bubble from a goroutine outside the bubble once a second. A stack dump marks each goroutine of a bubble with the bubble and each durable wait with `(durable)` (`runtime/traceback.go:1316` to `runtime/traceback.go:1331`), and the header of each goroutine has its labels (`runtime/traceback.go:1334`). A world calls `conform.RecordStep` after each of its steps, which counts the steps of the worlds of its bubble. A stall is two samples with the same stacks and the same count of world steps in which every goroutine of the bubble waits and one or more of them wait for a mutex or a system call. The equal count keeps a world whose stacks repeat between its steps from a report. The equal stacks keep a step whose goroutines move for more than a second from a report. The watchdog then writes the goroutines that do not wait durably and then the durable ones to standard error, with their incarnation labels and the limit of golang/go#77687, and ends the test binary with the exit status 2. The exit runs the exit hooks of the runtime, so the test binary writes its coverage, which a fatal panic does not (`runtime/proc.go:331` to `runtime/proc.go:352`). The sample period sets how soon a stall shows. A stall is a state of the bubble, and the period limits no run.
- **Replays.** A world check in native mode states `prop.Replays(10)`. A failing case is replayed ten times before the shrink, and the record states how many replays reproduced the failure, so a failure that reproduces 1 time in 10 is not reported as one that reproduces 9 times in 10. A run whose replays reproduce nothing ends as flaky, with the difference of its first replay.

```go
// Bubble runs run in a testing/synctest bubble on t, under the watchdog of a
// stall.
func Bubble(t *testing.T, run func(*testing.T))

// RecordStep records a step of the world that runs in the bubble of the
// calling goroutine, for the watchdog of Bubble. A call outside a bubble of
// Bubble does nothing.
func RecordStep()
```

### Tasks mode

A subject whose application states the option `conform.Tasks[S]()`, and a system that states `Tasks`, runs its hosts under a task scheduler of the case:

| The host | The task scheduler |
|---|---|
| `Env.Go` and `Incarnation.Go` | `Scheduler.Spawn` |
| `Env.Yield` | `Scheduler.Yield` |
| `Env.Wait` | `Scheduler.Wait` |
| `host.Park` of a context of an incarnation | `Scheduler.Wait` until the context ends |
| `Clock.Sleep` | The `Sleep` of `Scheduler.Clock` |
| `Clock.AfterFunc` | A task that sleeps on `Scheduler.Clock` until its deadline, and then calls the function unless `stop` came first |
| `Clock.Monotonic` | The time of `Scheduler.Clock` since the start of the host, times one plus the rate error of the host |
| The step `advance` of the world | `Scheduler.Advance`, which releases the tasks of each deadline at the time of that deadline |

Quiescence is the return of `Scheduler.Run` or `Scheduler.Advance`, which happens once no task is ready. Tasks run one at a time, so code between two waits does not need a mutex. A task that blocks on a channel, a mutex or a system call blocks the scheduler, and the watchdog of the bubble reports it. Under Tasks the world passes the checksum of its choices and of each output to `Case.Observe`, and the check states `prop.Confirm(10)`, so a replay of every tenth valid case finds a decision that the subject takes outside the case, such as an iteration of a map.

`Scheduler.Wait`, `Scheduler.Clock`, `Scheduler.Advance` and the option `MaxReleases` are new members of the definition of assert-spec, which its RFC on tasks that wait and sleep proposes. Each scheduler has a clock of its own, which starts at midnight UTC 2000-01-01 as the clock of a bubble does, so a replay of a case reads the same instants as the case.

### The world of a subject

A world runs as the property of a derived check of a suite, on one host:

```go
// Operation is one callable of a declaration as the world of a subject calls
// it.
type Operation[S any] struct {
	// Input requests the input of a call from the case. It picks its keys
	// from the pools of the case.
	Input func(c *prop.Case) any
	// Args returns the arguments of the call that the history records.
	Args func(input any) []any
	// Keys returns the keys of the call, which partition the check of the
	// history. A call without keys is in every partition.
	Keys func(input any) []any
	// Call calls the callable on s with the input under ctx, and returns the
	// output of the call and the error of the callable, nil for a call that
	// returned no error.
	Call func(ctx context.Context, s S, input any) (conform.Output, error)
	// Produced adds to the pools of the case the keys that a call returned
	// with success, such as the key of a record that it created. It is nil
	// for a callable that returns no key.
	Produced func(out conform.Output)
	// Name is the operation in the history, the name of the callable.
	Name string
}

// World is the world of the subject of a handle: the operations of its
// clients, the reference that judges their history, the error of a corrupted
// record, the number of clients and the bound of a call after the recovery.
type World[S any] struct {
	// Operations returns the operations of the case c. It chooses the
	// starting values of the pools of c once.
	Operations func(c *prop.Case) []Operation[S]
	// Reference applies an operation of the history to a reference.
	Reference Reference[S]
	// Corrupt is the error with which the subject reports a corrupted record,
	// and nil for a declaration that states none. The world flips bits at
	// rest only for a subject with Corrupt.
	Corrupt error
	// Clients is the number of clients, and 1 for 0.
	Clients int
	// Bound is the monotonic time of the host within which a call of the
	// recovery completes, and 1 s for 0.
	Bound time.Duration
}

// WorldOptions returns the options of the property run of a world check of
// h: prop.Confirm(10) in Tasks mode, and prop.Replays(10) in native mode.
func WorldOptions[S any](h *conform.Handle[S]) []prop.Option

// Amnesiac is the door Host of the defect form amnesiac-recovery. From the
// second incarnation of the host of env on, it removes every file and every
// directory of the disk of env. It then opens the subject of the application
// over env through Handle.FreshHost.
func Amnesiac[S any](h *conform.Handle[S], env *host.Env) S

// Pool is a pool of the keys of one type that the operations of a case pick
// from: the keys that the case chose at its start, and the keys that the
// calls of the subject produced. The world calls Pick and Add on the
// goroutine of the case.
type Pool[T any] struct{ /* unexported */ }

// NewPool returns a pool whose starting keys the case c chooses from g under
// label.
func NewPool[T any](c *prop.Case, g prop.Generator[[]T], label string) *Pool[T]

// Pick returns a key of the pool, as a choice of the case under label.
func (p *Pool[T]) Pick(c *prop.Case, label string) T

// Add adds the key v, which a call produced.
func (p *Pool[T]) Add(v T)

// FaultyWorldHistory runs the world w over the subject of h, and checks the
// history of the calls of its clients after every step with
// history.Linearizable and opts. It fails when a partition of the history has
// no order of its calls that respects their real-time order and that the
// reference accepts.
func FaultyWorldHistory[S any](h *conform.Handle[S], c *prop.Case, w World[S], opts ...history.Option)

// WorldRecovery runs the world w over the subject of h, then ends every
// fault, opens a crashed subject again and makes one call of each operation
// on each client. It fails when a call of that recovery has not completed
// within w.Bound.
func WorldRecovery[S any](h *conform.Handle[S], c *prop.Case, w World[S])
```

A world opens its subject through one of two doors:

- **`Durable`**, the door of the mixin `durable`, opens the subject over `Env.Disk()`. A subject that takes nothing else from its machine does not need another door.
- **`Hosted`**, a new option, opens the subject over the whole environment of an incarnation:

```go
// Hosted states the door that opens a subject over the environment of an
// incarnation of a host: its lifetime, its disk, its clock, its network and its
// dependencies.
func Hosted[S any](open func(h *Handle[S], env *host.Env) S) Option[S]

// Host returns a subject over env through the option Hosted, or through the
// door Host of a planted defect. Without the door it ends the check with a
// fault that wraps ErrOption.
func (h *Handle[S]) Host(env *host.Env) S

// FreshHost returns a subject of the application over env through the option
// Hosted, also on a planted defect, as FreshOpen returns a subject over a
// medium. Without the option it opens the subject over env.Disk() through the
// option Durable. Without either option it ends the check with a fault that
// wraps ErrOption.
func (h *Handle[S]) FreshHost(env *host.Env) S

// Hostable reports whether the handle has the door Host, from the option
// Hosted or from the doors of a planted defect.
func (h *Handle[S]) Hostable() bool

// Mode returns the execution mode of the world of a check: spec.ModeTasks
// for an application with the option Tasks, and spec.ModeNative otherwise.
func (h *Handle[S]) Mode() spec.Mode

// Blocking reports whether the application states the option Blocking.
func (h *Handle[S]) Blocking() bool

// RecordWorld adds the record w of the world of one case to the record of the
// check.
func (h *Handle[S]) RecordWorld(w *spec.WorldRecord)
```

`Doors` gains the field `Host` for the door of a planted defect. The probes of a world are in the package `conform/probe`, so they read the door, the mode and the option `Blocking` through the handle, and write the field `world` of the record through `RecordWorld`, as a probe adds the points of a medium through `AddPoints`.

#### The steps

The world is a machine of `stateful` whose steps the case chooses. Each case keeps a subset of the faults of the world, as swarm testing keeps a subset of the features of a test: the case chooses once for each fault whether the world strikes it, and each choice shrinks to a case without the fault. Each case keeps the calls, the cancel, the advance, the crash and the restart. So each case can restart a crashed subject, and each case can make calls.

| Step | Enabled | What the world does |
|---|---|---|
| A call of an operation | A client without an open call | Records the invocation, and makes the call through `Incarnation.Go` |
| `cancel` | A client with an open call | Ends the context of the call, which completes it as unknown |
| `advance` | Every state | Advances the monotonic clock of the host by a duration that the case chooses, from 1 µs to 10 s, which shrinks to 1 µs |
| `crash` | The subject runs | Crashes the incarnation by a kind that the case chooses among `ProcessCrash` and `PowerLoss` |
| `restart` | The subject has crashed | Starts a new incarnation, and opens the subject through its door |
| `pause` | The subject runs, in Tasks mode or under `Blocking` | Pauses the incarnation until the end of the fault or the crash of the incarnation |
| A fault of the disk | The subject runs | Instant: fails the next call before, after or in the middle of its effect, tears the next write, fills the disk until the end of the fault, or flips a bit at rest in a file with data of a subject with `World.Corrupt`. Blocking: blocks each call until a release or a crash |
| A fault of a dependency | The subject has a dependency of the injector | Instant: fails the next call, or each call until the end of the fault. Blocking: blocks each call until a release |
| A fault of the clock | The subject has read the wall clock | Bumps the wall clock, or strobes it |
| The end of a fault | A fault is in force | Releases a hold, lifts the limit of a full disk, ends the failure of a dependency, resumes a pause, or resets the wall clock |

At most one fault is in force at a time. A crash ends a pause and a hold of the disk, and every other fault remains in force across a crash until its end. A blocking fault is enabled only in Tasks mode and under `Blocking`. After each step the world calls `Host.Broadcast` and `conform.RecordStep`, waits for quiescence, and passes the keys of each call that completed with success to `Produced`. A step that the subject cannot take, such as a crash of a subject without a door that opens it, is not enabled.

`DiskReplacement` is a fault of a system, because a single durable subject that starts over an empty disk loses every write by design.

#### The outcomes of a call

| The call | Outcome in the history |
|---|---|
| It returns while its incarnation runs, without an error or with an error of the subject | OK with its output, its error included, as `conform.Complete` records an output |
| It returns with an error that wraps an error that the world injects, `syscall.EIO`, `syscall.ENOSPC`, `host.ErrInjected` or `World.Corrupt`, or an error of a crash or of its context | Unknown with the error, because the fault may have struck before or after the effect of the call |
| Its incarnation crashes while it is open | Unknown with `medium.ErrCrashed` |
| A `cancel` ends its context | Unknown with the error of the context |
| The case ends while it is open | Unknown with `medium.ErrCrashed`, because the end of the case crashes the incarnation |

`history.Linearizable` lets an unknown call take effect or not. A durable subject acknowledges a write once a crash can no longer lose it. The claims of the mixin `durable` state that contract. So the reference never crashes, and a write that completed before a crash must remain visible after it.

#### The keys of a case

`Operations` creates one pool for each type of key of the declaration. Each operation picks its keys from those pools. A pool starts with one to three keys and shrinks to one, so a counterexample shows the fewest keys that it needs. `Pick` follows the hotspot distribution of YCSB with a hot set of one key: the case chooses once per pool whether the pool is hot, and a pick from a hot pool takes the first key in half of its picks. Calls of different clients then collide on one key. Such a collision exposes a race.

A callable with the mixin `creates` returns the key of the record that it created as its first value result. The param `key` of the mixin refers to the field of that value that contains the key, as `//dokimi:mixin creates key=Key` refers to the field `Key` of an entry. The generator writes `Produced` for the callable where a pool has the type of the key, and that pool takes the key. A first value that is a nil pointer does not add a key. A callable without a value result, or with a key of a type that no pool has, does not produce a key, because no operation of a case takes it. In Tasks mode a replay produces the same keys. In native mode a replay can produce other keys, and the difference ends the case as flaky. The keys of a pool never depend on the reference.

#### The checks

`FaultyWorldHistory` checks the history after every step against the spec that `history.SpecFrom` builds over fresh references of the handle. Each check continues the search of each partition of the check before it through `history.Resume`, so the check of the last step covers the whole history. A partition whose search ends undecided at its budget ends the checks of the case and does not engage the check, as in `LinearizableHistory`. The options of the call, such as `history.Budget` and `history.Workers`, set the limits of each check, as the options of `LinearizableHistory` do. The generated checks pass none.

`WorldRecovery` ends each fault and opens a crashed subject again once the steps end. It then makes one call of each operation on each client, advances the monotonic clock of the host by `w.Bound`, and fails for each call of that recovery that is still open.

### The claims of a world

The runtime contract gains one machinery and two claims:

| Machinery | Statement |
|---|---|
| `quiescence` | The run waits before each step until no goroutine or task of the subject can run: every goroutine is durably blocked, or under the task scheduler no task is ready |

| Claim | Law | Probe | Reason | Defect form |
|---|---|---|---|---|
| `linearizable-under-faults` | Through crashes of the subject and faults of its medium, its dependencies and its clock, each partition of the history of concurrent calls has an order that respects their real-time order and that the reference accepts | `faulty-world-history` | `not-linearizable-under-faults` | `amnesiac-recovery`: each restart opens the subject over a new medium |
| `completes-after-recovery` | Once every fault ends and the subject opens again, each call completes within the bound of the world | `world-recovery` | `call-open-after-recovery` | `blocks-after-recovery`: each call after the first restart blocks until its case ends |

Both claims list the machineries `quiescence`, `generated-sequences` and `killable-process-or-breakable-medium`. `linearizable-under-faults` also lists `concurrent-callers`, `call-that-fails-on-demand`, `reference-implementation` and `clock-the-run-advances`. The probes take one input of the kind `world`.

The planted defect of each claim of a world changes the subject only from the second incarnation of its host on. The checks of `synced-writes-survive` and `no-torn-record` also open their subject again after each crash of a medium, so a defect that changed every second open would fail those checks as well. Isolation would then fail it, because a planted defect fails the check that it proves and no other. The defect of `racy-read-write` changes the subject only while calls overlap, for the same reason.

- **`amnesiac-recovery`** is the door Host `probe.Amnesiac`. It erases the disk of each incarnation after the first, and then opens the subject of the application over the environment through `Handle.FreshHost`.
- **`blocks-after-recovery`** is a planted defect whose callables block until their context ends from the second incarnation on, as the generator writes the defect of `racy-read-write`. Its door Host reads the number of the incarnation from the environment. A callable without a context blocks until the lifetime of the incarnation ends. Each callable blocks through `host.Park`, so in Tasks mode its task parks and `Scheduler.Run` returns.

A world opens each incarnation through the door Host of the planted defect, because the defect has the door. A check without a world, such as a durable check, opens the defect through the door Open of the application, so the defect is the subject of the application there.

The rules of dokimi license each claim once for a declaration with a callable with the mixin `durable`. `linearizable-under-faults` also needs a reference, and the gate `returns-a-value` of a callable of the declaration. The world calls each operation from each client, so only a declaration whose each operation has the mixin `concurrent` gets three clients. Any other declaration gets one. An application without a door that opens the subject, or without the door `Reference`, fails each check with a fault that wraps `conform.ErrOption`.

The record of a check in the suite record gains the field `world`:

| Field | Value |
|---|---|
| `mode` | `tasks` or `native` |
| `replays` and `reproduced` | The replays of a failing case, and how many of them reproduced the failure |
| `steps` | The steps of the cases |
| `crashes` | The crashes, by kind |
| `faults` | The faults that the world struck, by fault, and the faults that it recorded as vacuous |
| `undecided` | The partitions whose search ended undecided |
| `leaked` | The goroutines of crashed incarnations that remained alive |

### The generated checks

The suite generator writes the operations of a declaration and one check per licensed claim into its suite:

```go
// storeWorld is the world of the subject of Store.
var storeWorld = probe.World[store.Store]{
	Operations: storeOperations,
	Reference:  storeReference,
	Clients:    3,
}

// storeOperations returns the operations of the world of Store in the case c,
// over the pools of keys that c chooses.
func storeOperations(c *prop.Case) []probe.Operation[store.Store] {
	keys := probe.NewPool(c, prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)),
		"keys.string")
	return []probe.Operation[store.Store]{
		{
			Name: "Put",
			Input: func(c *prop.Case) any {
				return []any{keys.Pick(c, "put.key"), c.Draw(prop.Of[string](), "put.value")}
			},
			Args: storeArgs, Keys: storeFirstKey, Call: storePutCall,
		},
		{
			Name:  "Get",
			Input: func(c *prop.Case) any { return []any{keys.Pick(c, "get.key")} },
			Args:  storeArgs, Keys: storeFirstKey, Call: storeGetCall,
		},
	}
}

{
	ID:          StoreChecks.Get.LinearizableUnderFaults(),
	Statement:   "the history of the calls of Store through crashes and faults has a linearization that the reference accepts",
	Claim:       spec.LinearizableUnderFaults,
	Machineries: []spec.MachineryID{spec.Quiescence, spec.ConcurrentCallers /* and the others */},
	Around:      conform.Bubble,
	Options:     probe.WorldOptions[store.Store],
	Property: func(h *StoreHandle, c *prop.Case) {
		probe.FaultyWorldHistory(h, c, storeWorld, history.Budget(200000))
	},
},
```

The generator writes the functions of each operation from the derived binding of its callable, and the pools from the types of its key inputs, as the actions of `LinearizableHistory` take their keys. A check of a world binds on the first callable of the declaration that passes its gates, as the check of `agrees-with-the-reference` binds. The check of `linearizable-under-faults` binds on the first callable that returns a value, such as `Get`. A check of a world is excused, refused and proven as every derived check is.

The check of `linearizable-under-faults` passes the budget of the history checks of the suite, 200,000 steps, as the check of `linearizable` does. A search of a partition ends undecided once it has spent its budget, and the default budget of `history.Linearizable` is 10,000,000 steps.

### The faults of a medium

```go
// FailBefore makes the next call of the medium or of one of its files fail
// with an error that wraps err, and the call does not change the medium.
func (m *Memory) FailBefore(err error)

// FailAfter makes the next call of the medium or of one of its files take its
// whole effect, and then fail with an error that wraps err.
func (m *Memory) FailAfter(err error)

// FailPart makes the next write of a file of the medium write its first n
// bytes, and then fail with an error that wraps err.
func (m *Memory) FailPart(n int, err error)

// Hold blocks each later call of the medium and of its files until Release
// or Crash. A call that a crash ends fails with ErrCrashed.
func (m *Memory) Hold()

// Release ends the hold of the medium, and each blocked call goes on.
func (m *Memory) Release()

// Kill ends the process that the medium serves, as a kill of the process
// does. Each open file is closed, and its later calls fail with ErrCrashed.
// The data of every write remains, durable or pending.
func (m *Memory) Kill()

// CrashKeeping crashes the medium as Crash does, and keeps the first keep(name,
// pending) pending writes and truncates of each file name, in the order of the
// calls, and drops the others. keep returns 0 to drop every pending call of a
// file, and the count of its pending calls to keep every one. The last kept
// write of a file is torn where Tear armed it. The pending entries of each
// directory are kept or dropped as a whole, by keep of the name of the
// directory with a count of 1.
func (m *Memory) CrashKeeping(keep func(name string, pending int) int)
```

The world strikes `FailBefore`, `FailAfter` and `FailPart` with an error that wraps `syscall.EIO`. A crash by `PowerLoss` calls `CrashKeeping` with counts that the case chooses for each file and each directory, and that shrink to 0, the crash state of `Crash`. So a crash can keep a later write of one file and drop an earlier write of another, which is the reordering that ALICE explores, and gosim keeps a random subset of the pending writes in the same way. A hold blocks each call on a channel that `Hold` makes inside the bubble, so a blocked call waits durably. Under the scheduler, the view of an incarnation on the disk parks the task of each call while `Memory.Held` reports the hold, because the scheduler does not see a wait on a channel.

### Components in a system

A process of the world of a system can run a Go component on a host of its own, beside the machine processes:

```go
// Component is a process of a world that runs a Go value of the type T on a
// host. Each incarnation of the process is the value that open returns over a
// new environment of the host.
type Component[T any] struct{ /* unexported */ }

// SpawnComponent adds a component process named name to w on a new host, and
// starts its first incarnation.
func SpawnComponent[T, E, F any](w *World[E, F], name string, open func(env *host.Env) T) *Component[T]

// Running reports whether an incarnation of the process runs, so that Value
// returns its value. It reports false while the process has crashed and after
// the stop of the world, such as for the check of an invariant that follows
// the stop.
func (c *Component[T]) Running() bool

// Value returns the value of the incarnation that runs. It fails the case of
// the world while the process has crashed.
func (c *Component[T]) Value() T

// ComponentOf returns the component process named name of w, whose values
// have the type T, as SpawnComponent returned it. A function of a
// simulation, such as an action or an invariant, gets the components of its
// world from it. ComponentOf fails the case of w for a name that no process
// has, a process that is no component and a component whose values have
// another type.
func ComponentOf[T, E, F any](w *World[E, F], name string) *Component[T]

// Invoke records the invocation of operation by client in the history of the
// case with args and keys, and runs call on a goroutine of the incarnation of
// the process.
func (w *World[E, F]) Invoke(
	client int, process, operation string, args, keys []any, call func(ctx context.Context) (any, error),
)

// ReplaceDisk crashes the component process named name by a disk
// replacement. Its open calls end as unknown with medium.ErrCrashed, and its
// next incarnation starts over a new empty disk of its host.
func (w *World[E, F]) ReplaceDisk(name string)

// FaultDiskReplacement crashes a component until a restart that starts over a
// new empty disk.
const FaultDiskReplacement Fault = "disk-replacement"
```

`Invoke` completes the call as OK with its output for a nil error. It completes it as unknown for an error that wraps `context.Canceled`, `context.DeadlineExceeded`, `medium.ErrCrashed` or `host.ErrReset`, and as failed for any other error. A crash of the process of the call completes an open call as unknown with `medium.ErrCrashed`, and the end of the run completes it as unknown with `ErrRunEnded`. The history allows a client one open call at a time.

The simulation states the mode and the network of its components:

```go
type Simulation[E, F, S any] struct {
	// The fields of a system of machines, as before.

	// History are the options of each check of the history against the
	// spec, such as history.Budget, as the options of FaultyWorldHistory are.
	History []history.Option
	// Around runs each case inside an environment, such as conform.Bubble,
	// which a world with a component in native mode needs.
	Around func(t *testing.T, run func(*testing.T))
	// Mode is the mode of the hosts of the components: spec.ModeTasks runs
	// them under a task scheduler of the case, and spec.ModeNative on the
	// goroutines of the bubble.
	Mode spec.Mode
	// Network are the options of the network of the hosts of the components.
	Network []host.NetworkOption
	// Blocking lets a world in native mode strike the blocking faults on its
	// components. A world in Tasks mode strikes them without it.
	Blocking bool
}
```

A component is a process of the world, so every fault of a process strikes it through its host:

- a crash by `ProcessCrash` or `PowerLoss`, a choice of the case that shrinks to `ProcessCrash`
- a restart that opens a new incarnation and waits for the open, as the world of a subject waits for it
- a pause of its incarnation
- a jump of its wall clock through `Host.Bump`
- a torn write, a flipped bit and a full disk on the disk of its host
- the sides of a partition on the network

The world reads the goroutines of a crashed incarnation as the world of a subject reads them. It fails the case for each one that is alive after quiescence. A strobe strikes machine processes alone. The strobe of a machine counts the steps of its process, and a component takes no step.

The fault `disk-replacement` crashes a component by `DiskReplacement`. Its restart opens the next incarnation over a new empty disk. Such a crash loses every write of a single component, so a simulation tolerates the fault only for a system that keeps a copy of each write on another component. `Tolerate` bounds the components that a disk replacement keeps down at once, apart from the bound of `crash`. `Faults` lists such a component as `disk-replacement:a` until its restart.

`drop`, `delay` and `duplicate` strike the oldest datagram in flight to a component, and each is enabled while `Network.InFlight` reports enough datagrams. `reset` resets the streams between two components, and `slow-link` raises the latency of the link between two components until `restore-link`, as a blocking fault.

The step `advance` of a world with a component is always enabled. While a machine has a timer, the step advances by the time until the next timer and up to as much again. Otherwise it advances by 1 µs to 10 s, as the world of a subject advances.

A world with a component in native mode runs each case in the environment `conform.Bubble` through `Around`. The seat of `sim.Test` then has a `*testing.T` for such a world. After every step the world calls `Host.Broadcast` on each host and `conform.RecordStep`, and it waits for quiescence. The methods of the world act on its components:

- `World.Advance(d)` advances the monotonic clocks of the hosts by `d` times their rate errors, and the clocks of the machine processes by `d` times their drift.
- `World.Crash` crashes the incarnation of a component process, and `World.ReplaceDisk` crashes it by `DiskReplacement`.
- `World.Restart` starts the next incarnation of a crashed component, and `World.Pause` pauses the incarnation that runs.
- `World.Reset` and `World.SlowLink` strike the network through `Network.Reset` and `Network.Slow`. The fault `slow-link` is a blocking fault.
- `World.Deliver`, `World.Call` and `World.After` fail the case for a component process. A component takes its calls through `Invoke` alone.

A run crashes each component that runs by `ProcessCrash` at the end of each case, also after a failure. No goroutine of a component then runs after the end of its case.

#### A system without a machine interface

```go
// None is the type of the events and the effects of a world without machine
// processes. No type implements it.
type None interface{ none() }
```

A system of components alone is a `World[None, None]`. The sim generator writes the simulation of a package with `//dokimi:system` and without a machine interface. Such a package states `tolerate`, `strategy`, `require` and `isolation` on `//dokimi:system`, and states the modes `tasks` and `blocking` as values of `tolerate`. Its functions under `action`, `spec`, `invariant`, `converges`, `fingerprint` and `scenario` take a `*sim.World[sim.None, sim.None]`, and get each component of the world through `sim.ComponentOf`. The generated simulation runs each case in `conform.Bubble`, and sets `Mode` and `Blocking` from the values of `tolerate`. It has no rule of sim, because each rule of sim judges the steps of machines.

### The history after every step

The world of a subject and the run of a system both check each partition of the history after every step through `history.Resume`, which continues the search of each partition of the check before it. So a violation ends its case at the step where it occurs, and the check after the final phase covers the calls of that phase. A partition in which a pending call has completed since the check before it searches again from its first call.

### Migration

| Step | Repository | Change |
|---|---|---|
| 1 | assert-go | `Resume` for a partitioned check, with its specs |
| 2 | assert-spec | The waits of tasks, and the replays of a failing case, each through an RFC of assert-spec, its names and its vectors |
| 3 | assert-go | The new members of `stateful` and `prop`, with their specs |
| 4 | dokimi-spec | The machinery `quiescence`, the two claims with their laws, probes, reasons and defect forms, the input kind `world`, and the field `world` of the suite record |
| 5 | dokimi-go | `host`: the new module. `medium`: the seven methods. `conform`: `Hosted`, `Tasks`, `Blocking`, `Bubble`, `RecordStep`, the methods of `Handle` and the aliases of the injector. `conform/probe`: `World`, `Operation`, `Pool`, the two probes and their defect forms. `sim`: components on hosts, `ComponentOf`, `Invoke`, `None`, the faults of the network and the check after every step. The copy of the spec follows its main branch |
| 6 | eidos | The mixin `creates` in the shape catalog |
| 7 | dokimi | The rules of the two claims, the operations, pools and checks of the suite generator with `Produced` for a callable with the mixin `creates`, the systems without a machine interface of the sim generator, and their goldens |

Each step needs the push of the step before it. The corpus gains two fixtures:

- `iface/door/hosted`, the fixture of the door Hosted, whose declaration has two subjects:
  - a durable store whose writes a goroutine flushes to its medium on a timer, in native mode under `Blocking`, the component that #58 describes
  - the same store written to the host's task interface, in Tasks mode
- `sim/mirrored`, a system of two components that replicate a register over a stream

### What is in and what is out

| Capability | Verdict | Reason, and what reverses it |
|---|---|---|
| The host, shared by conform and sim | In | Every simulator of real code scopes its resources to a machine. One host keeps the fences, the clocks and the network of both worlds alike |
| The kinds of crash: a process crash, a power loss with crash states, a disk replacement and a pause | In | A system that relies on the page cache survives a process crash and not a power loss, and the world tests what the subject claims |
| A disk replacement as a fault of its own in sim | In | A single component loses every write at a disk replacement, so a crash that can replace the disk fails the simulation of every single component. A system that keeps copies of its writes states the fault in `Tolerate` |
| The accounting of goroutines by incarnation | In | Profiler labels pass to every goroutine that a labelled goroutine starts |
| A wall clock and a monotonic clock | In | `time.Time` couples the two readings, so one clock with an offset fakes durations that production never measures |
| Native mode, with instant faults by default and blocking faults under `Blocking` | In | A bubble cannot become quiescent while a goroutine waits for a mutex that a durably blocked goroutine has locked. Durable mutexes in Go, as golang/go#77687 tracks, reverse the default |
| The watchdog of a bubble | In | A stall that `synctest` cannot detect otherwise hangs until the timeout of `go test`, without its cause |
| Settings of the latency range and the lowest throughput of the links | In | The timeouts of each subject decide which latencies matter, so no fixed range fits every system. The latency range of turmoil is the default, and turmoil does not limit the throughput of a link |
| Tasks mode with waits and sleeps of the scheduler | In | A task that waits by yielding keeps `Run` from returning. A scheduler that parks its waiting tasks returns once no task is ready, and its clock moves from deadline to deadline, as the clocks of madsim and `testing/synctest` move to the next timer |
| Replays of a failing case with a reproduction count | In | One replay cannot tell a failure that reproduces 1 time in 10 from one that reproduces 9 times in 10 |
| Pools of keys per case, with the keys that the subject produced | In | The pools are the state of the machine of the world, which the definition of assert-spec prefers to bundles |
| A value shared by key across a case in the property engine | Out | The pools of a world are values that the operations of a case create, so they need no member of the definition. A second use of shared values outside a machine reverses it |
| Synchronization primitives under a build tag | Out | They make a durable mutex only of the locks of the subject, and the libraries that a subject calls still lock `sync.Mutex`. Tasks mode serves code that adopts the interface of the host |
| A WebAssembly host that serves the clock and the randomness of the runtime | Out | It runs the wasip1 port of a component, not the code that production runs, and it needs a WebAssembly runtime as a dependency |
| A translation of the source onto a scheduler, as gosim does | Out | It runs a rewrite of the component, not the code that production runs |
| Inputs that read the state of the reference | Out | They make each replay of a racy case diverge |
| A crash of a declaration whose acknowledgement is not durable | Out | Its reference needs the semantics of a crash. A specification that models a crash as an operation reverses it |

## Alternatives considered

### A WebAssembly host for the whole program

The test binary, or each component, is built with `GOOS=wasip1` and runs under wazero. The host serves `random_get` from a seed, and serves `clock_time_get` from a clock that advances only when the program sleeps. Polar Signals runs the deterministic simulation of FrostDB this way, with a modified Go runtime. A probe on Go 1.27.3 and wazero v1.12.0 did not need a modified runtime. Seed 1 gave the same order of 1,696 contended operations in two runs, and seed 2 gave another order.

**Why not:** the simulation would run the wasip1 port of a component instead of the code that production runs. The wasip1 port lacks cgo, mmap and the race detector. Each call of its file system stalls every goroutine of the program. The simulation would also depend on a WebAssembly runtime.

### A translation of the source

gosim translates a Go program onto a coroutine scheduler, and detsim rewrites the source of a package with `go/ast` onto a scheduler of its own, through an overlay of the build.

**Why not:** both run a rewrite of the component under test. detsim covers goroutines, channels, `sync`, `time`, `math/rand` and basic network and file calls, and rejects code with unsafe fields of `sync` types.

### One clock with an offset

The clock of a host returns the time of the bubble plus the offset of its faults, through one method `Now`.

**Why not:** `t.Add` moves the monotonic reading with the wall reading. A bump of the wall clock would also move every duration that a subject measures with `Sub` or `Since` between two reads, which Go code in production never sees. A clock that strips the monotonic reading moves those durations as well.

### A round of the scheduler

`Scheduler.Run` releases each task that is ready at the start of a round once, and returns.

**Why not:** a task that waits by yielding then takes one release in each round of each step. The run does not record that the task waits. The return of a round would also not mean that the subject cannot move, so the world would take its next step while a task can still run. A task that parks on `Wait` until its condition is true gives quiescence an exact meaning.

### A value shared by key across a case

The property engine returns, for each key, the value of the first draw of a case to every later draw of the case, as `st.shared` of Hypothesis does.

**Why not:** the definition of assert-spec rejected bundles for machines, because the state of a machine already gives a value to `sampled-from`. The pools of a world are that state. They are values that the operations of a case create, and the keys that the subject produces join them through `Produced`.

### Durable mutexes through a build tag

A package of dokimi-go replaces `sync.Mutex` under a build tag with a mutex over a channel of the bubble, which waits durably.

**Why not:** the consumer changes the imports of its locks, and the libraries that the subject calls, such as `database/sql`, still lock `sync.Mutex`. Tasks mode gives code that adopts the interface of the host a whole replay, and the watchdog states the cause of a stall for the rest.

### A harness of its own for each subsystem

A generator of a subsystem lists its components and its entry points in directives, and the consumer writes workloads and invariants in types of the harness.

**Why not:** the suite of a declaration already derives the calls with their inputs, the reference and the doors of the subject. A second set of types would repeat them. A check of a world binds a claim, has a planted defect and an excusal, and shows in the suite record as every derived check does.

### The world of a subject in sim

`go.dokimi.dev/sim` runs the world of one subject as a system of one component.

**Why not:** the world of a subject judges a claim of its declaration, so it belongs to the suite of the declaration and to the doors of its handle. sim runs systems of processes. Both run on the host, so they share its disk, its clocks, its network and its dependencies.

### Checks only at the end of a case

The world checks the whole history once, after its last step, as the run of a system does today.

**Why not:** a check after every step ends a case at the step of the violation. The shrinker then reruns fewer steps for each candidate.

## Drawbacks

- **The order of the change.** The waits of tasks and the replays of a failing case change the definition of assert-spec in every language. The claims change the runtime contract. The change runs through assert-spec, assert-go, dokimi-spec, dokimi-go and dokimi in that order.
- **New API.**
  - The host has `Host` with nine methods, `Incarnation` with five, `Env` with thirteen, `Clock` with four, `Network` with nine, `Kill` with three values, the options `Latency` and `Throughput`, `ErrReset`, `ErrNoNetwork`, `Park`, and the injector that moves from conform.
  - medium gains seven methods.
  - conform gains `Hosted`, `Tasks`, `Blocking`, `Bubble`, `RecordStep`, the methods `Host`, `FreshHost`, `Hostable`, `Mode`, `Blocking` and `RecordWorld` of `Handle`, and a field of `Doors`.
  - probe gains `World`, `Operation`, `Pool`, `WorldOptions`, `Amnesiac` and two probes.
  - sim gains `Component`, `SpawnComponent`, `ComponentOf`, `Invoke`, `ReplaceDisk`, `None`, the fields `History`, `Around`, `Mode`, `Network` and `Blocking` of `Simulation`, and the faults `reset`, `slow-link` and `disk-replacement`.
- **A contract for the subject.** The subject takes its disk, its network, its clock and its dependencies from its environment, and ends its goroutines once its lifetime ends. A subject that calls the `os`, `net` or `time` packages directly does not receive a fault through those calls. A subject that waits durably with a mutex locked that another goroutine waits for needs Tasks mode or the default of instant faults.
- **Two clocks for the subject.** A subject reads timestamps from `Wall` and durations from `Monotonic`. A subject that measures a timeout with two reads of `Wall` sees the faults of the wall clock in that timeout, as it would on a machine whose clock steps.
- **Flaky runs in native mode.** A race inside a subject ends a run as flaky when no replay reproduces it, and a person reads such a run. The shrinker of such a case can stop at a larger case.
- **A restart in one process.** A crash fences the disk, the dependencies, the network and the lifetime of an incarnation. Its package variables and its `sync.Once` values remain for the next incarnation, which a restart of a real process does not keep. gochecknoglobals finds such variables in the package of a subject.
- **The cost of a check.** A world check runs up to 100 steps per case. Each step ends with a wait for quiescence and an incremental check of the history, and a failing case in native mode costs ten replays. On 4 cores, the two world checks of the corpus fixture `iface/door/hosted` take 0.7 s and 0.9 s, and 6.5 s and 10.5 s under the race detector.

## Open questions

None.

## Unresolved and future work

- A declaration whose writes survive only an explicit sync needs a reference that models a crash. That specification is not proposed here.
- A check that a new version of a durable subject reads the data of an earlier version is not proposed here.
- The runtime of TypeScript is not proposed here.

## References

| What | Where |
|---|---|
| The request of a world that runs goroutines | <https://github.com/dokimasia/dokimi/issues/58> |
| `testing/synctest` of Go | <https://pkg.go.dev/testing/synctest> |
| testing/synctest: support durably blocking on mutex acquisition | <https://github.com/golang/go/issues/77687> |
| The monotonic clock of `time.Time` | <https://pkg.go.dev/time#hdr-Monotonic_Clocks> |
| The preemption, the startup randomness and the stack dumps of the Go runtime | <https://github.com/golang/go/tree/go1.27.3/src/runtime> |
| The partitioned `Resume` of assert-go | <https://github.com/dokimasia/assert-go/blob/main/docs/rfc/0025-resuming-each-partition-of-a-growing-history.md> |
| The task scheduler and the bundles of the definition of assert-spec | <https://github.com/dokimasia/assert-spec/blob/main/docs/rfc/0012-machines-simulation-and-campaigns.md> |
| Tasks that wait and sleep under the task scheduler of assert-spec | <https://github.com/dokimasia/assert-spec/blob/main/docs/rfc/0056-tasks-that-wait-and-sleep.md> |
| The replays of a failing case of assert-spec | <https://github.com/dokimasia/assert-spec/blob/main/docs/rfc/0057-the-replays-of-a-failing-case.md> |
| FoundationDB, simulation and testing | <https://apple.github.io/foundationdb/testing.html> |
| madsim, a deterministic simulator of Rust | <https://github.com/madsim-rs/madsim> |
| turmoil, the simulation of hosts and networks of Tokio | <https://github.com/tokio-rs/turmoil> |
| gosim, a deterministic simulation of Go programs by translation | <https://github.com/jellevandenhooff/gosim> |
| detsim, a deterministic scheduler for Go by source rewriting | <https://github.com/arshnah/detsim> |
| TigerBeetle, the deterministic simulator VOPR | <https://github.com/tigerbeetle/tigerbeetle/blob/main/docs/internals/vopr.md> |
| Jepsen, the clock nemesis | <https://github.com/jepsen-io/jepsen/blob/main/jepsen/src/jepsen/nemesis/time.clj> |
| Pillai et al., "All File Systems Are Not Created Equal: On the Complexity of Crafting Crash-Consistent Applications", OSDI 2014 | <https://www.usenix.org/conference/osdi14/technical-sessions/presentation/pillai> |
| YCSB, the request distributions of the core workload | <https://github.com/brianfrankcooper/YCSB> |
| Hypothesis, shared strategies and bundles | <https://hypothesis.readthedocs.io/en/latest/stateful.html> |
| Polar Signals, (Mostly) Deterministic Simulation Testing in Go | <https://www.polarsignals.com/blog/posts/2024/05/28/mostly-dst-in-go> |
| Elle, inferring isolation anomalies from histories | <https://arxiv.org/abs/2003.10554> |
| assert-go, the search of `Linearizable` over calls of unknown outcome | <https://github.com/dokimasia/assert-go/issues/89> |
