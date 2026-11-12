---
rfc: 0051
title: The dependency of each host and the restarts of the active phase
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-11
updated: 2026-11-11
discussion: https://github.com/dokimasia/dokimi/issues/92
supersedes: RFC-0043 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0051: The dependency of each host and the restarts of the active phase

## Summary

We propose two changes to the simulation of a system:

- `//dokimi:perform` gains the param `scope`. Under `scope=host` each host of a world has an instance of the dependency, and the world keeps a state of faults for each instance. A fault of the dependency then strikes the effects of one host. A dependency without the param keeps one state for the world, as a database that every host calls does.
- The run restarts a host that exits in the active phase, with the delays that it applies to such a host in the final phase.

## Motivation

A fault of a dependency strikes the effects of every host at once:

- `Anatomy.Dependency` returns the name of the dependency of an effect (`sim/world.go:121` of dokimi-go at sim v0.13.1).
- `World.dependency` keeps one state of faults for each name (`sim/fault.go:621`), and `faulted` applies that state to the effect of any process (`sim/fault.go:635`).
- `LoseBefore`, `LoseAfter`, `Outage`, `Slow` and `Restore` take the name of a dependency (`sim/fault.go:243-319`).

A system whose hosts each write through a separate instance of a dependency cannot fault the instance of one host. Each member of a Raft group of sharder persists its log through `//dokimi:perform func=performPersist dependency=quorum.Storage result=Persisted`, and each host gives its member a separate storage. An outage of `quorum.Storage` fails the writes of every member at once, and every member stops itself. A slowness of `quorum.Storage` keeps the writes of every member in flight. No election then ends, because each member persists its vote before it sends it.

The run starts an exited host again in the final phase alone. `settleDown` calls `supervise` before each of its steps (`sim/run.go:751`). In the active phase an exited host remains down until the case chooses the action `restart`.

A probe ran 200 cases of the simulation of the member of sharder at the seed 7, on a scratch copy of sim v0.13.1 with each change:

| sim | Runs with a leader before the final phase | Runs with an applied commit |
|---|---|---|
| v0.13.1 | 4 | 0 |
| v0.13.1 with the restarts of the active phase | 4 | 0 |
| A state of faults for each host | 16 | 5 |
| A state of faults for each host, with the restarts | 23 | 4 |
| A state of faults for each host, with the restarts, without faults | 88 | 23 |

The restarts lowered the checks at which a member had stopped itself from 29,305 to 11,843, of about 175,000 checks.

## Detailed design

### Components

| Component | Where | Responsibility |
|---|---|---|
| The param `scope` | dokimi, `generator/sim` | Accepts `world` and `host` on `perform`, and writes the scope of each dependency into the anatomy |
| `Scope` and `Anatomy.Scope` | `go.dokimi.dev/sim` and `@dokimi/runtime/sim/anatomy` | The scope of a dependency |
| The instances of a dependency | `World` of both runtimes | A state of faults for each instance |
| The supervisor of the active phase | `run` of both runtimes | Restarts each host that exits in the active phase |

### The directive

```go
// Storage is the storage of one member.
type Storage interface {
	Persist(ctx context.Context, entries []Entry) error
}

//dokimi:perform func=performPersist dependency=Storage result=Persisted scope=host
type Persist struct {
	Entries []Entry
}
```

| Directive | Meaning |
|---|---|
| `//dokimi:perform … scope=world` | The world has one instance of the dependency. It is the default |
| `//dokimi:perform … scope=host` | Each host has an instance of the dependency |
| `scope` on a `perform` without `dependency` and without `call` | `BadPerform`: an effect without a dependency has no instance |

The param applies to both forms of a dependency. The dependency of `call=` is the interface of its method, and the dependency of `func=` is the interface of `dependency=`.

### The anatomy

```go
// Scope is the scope of a dependency. The world has one instance of a
// dependency of the scope ScopeWorld, and each host has an instance of a
// dependency of the scope ScopeHost.
type Scope uint8

const (
	// ScopeWorld gives the world one instance of the dependency, as a
	// database that every host calls.
	ScopeWorld Scope = iota
	// ScopeHost gives each host an instance of the dependency, as the
	// storage of a member.
	ScopeHost
)

type Anatomy[E, F any] struct {
	// …
	// Scope returns the scope of the dependency named dependency. It is nil
	// for machines whose every dependency has the scope ScopeWorld.
	Scope func(dependency string) Scope
}
```

```ts
// @dokimi/runtime/sim/anatomy
/** The scope of a dependency. The world has one instance of a dependency of the scope world, and each host an instance of a dependency of the scope host. */
export type Scope = "world" | "host";

export interface Anatomy<E, F> {
  // …
  /** Returns the scope of the dependency. It is undefined for machines whose every dependency has the scope world. */
  readonly scope?: (dependency: string) => Scope;
}
```

The generated anatomy of a machine with a dependency of the scope `host` sets `Scope` to a function that returns `ScopeHost` for each such dependency and `ScopeWorld` otherwise. The generated anatomy of every other machine leaves it unset.

### The instances

- The name of an instance of the scope host is the name of its host, a slash and the name of the dependency, such as `node-2/quorum.Storage`. An instance of the scope world has the name of the dependency.
- `faulted` takes the state of the instance of the host of the process that performs the effect. A wait keeps the name of its instance, so an outage of `node-2/quorum.Storage` ends the waits of node-2 alone.
- `LoseBefore`, `LoseAfter`, `Outage`, `Slow` and `Restore` take the name of an instance. The fault actions choose their target among the instances that the world has seen, as they choose among the dependencies in sim v0.13.1. `Faults` and the record of a run list the instance, such as `outage:node-2/quorum.Storage`.
- A crash and a restart of a host keep the state of each instance of the host, as they keep the disk of the host. An outage of the storage of a member continues after a restart of the member.

### The restarts of the active phase

- `run` calls `supervise` before each action of the active phase, as `settleDown` calls it before each step of the final phase.
- One map of backoffs serves both phases. The delays of a host that exits in both phases grow across them.
- A restart whose timer finds its host running does nothing. The action `restart` of the active phase can start an exited host before its timer.
- `supervise` restarts an exited host, and leaves a crashed host to the action `restart`, as it does in the final phase. A restart of the supervisor is no fault, and `Tolerate` does not count it.

### Failure handling

| Case | Result |
|---|---|
| `scope` with a value other than `world` and `host` | The directive validation refuses the param |
| `scope` on a `perform` without a dependency | `BadPerform` at the directive |
| A fault on the name of an instance that no effect performs | The fault strikes no effect, and the record of the run lists it as vacuous |
| The timer of a restart after the action `restart` started the host | The timer does nothing |

### Migration

- A simulation without the param `scope` keeps one state of faults for each dependency. Its generated anatomy does not change.
- The restarts of the active phase change the run of each simulation whose host exits in the active phase. Its stored cases can stop replaying the same run.
- The TypeScript simulation takes both changes in the forms of the preceding sections.

### Verdicts

| Capability | Verdict | Reason | What would reverse it |
|---|---|---|---|
| The param `scope` with `world` and `host` | in | The members of sharder write through separate storages, and a state of faults for each host raised the runs with a leader from 4 to 16 of 200 | A system whose instances of one dependency are shared by some hosts and not by others |
| `world` as the default | in | A database that every host calls has one state, and each existing simulation keeps its runs | none known |
| The name of a host and a slash before the name of an instance | in | The fault actions and the records keep a text target, and the name of a dependency is a package and an interface without a slash | A dependency whose name contains a slash |
| The state of an instance across the restarts of its host | in | The storage of a member remains after its process ends | A dependency that a restart of its process replaces, which a component of the host models instead |
| Restarts of an exited host in the active phase | in | A supervisor of a deployment starts an exited process at any time, and the restarts raised the runs with a leader from 16 to 23 of 200 | A system whose active phase needs an exited host to remain down until the case restarts it |
| A host as a second input of a fault action | out | The name of an instance identifies the host, so a fault action needs no second input | A fault that strikes the instances of several hosts at once |

## Alternatives considered

### Every dependency of the scope host

The world could keep a state for each host of every dependency, without a param. A database that every host calls would then fail for one host and serve the others, which no deployment of such a database does.

### A name of the host in the param dependency

The directive could list the instance of each host, such as `dependency=quorum.Storage@node-2`. The hosts of a world come from its inputs at run time, so the directive cannot name them.

### The action restart alone in the active phase

The active phase could leave an exited host to the action `restart`. The probe measured 16 runs with a leader of 200 that way, against 23 with the supervisor.

## Drawbacks

- A fault of a dependency of the scope host strikes one instance. A world of 5 hosts has 5 instances of `quorum.Storage`, so an outage of one instance strikes a fifth of the writes, and a run needs more faults to strike every member.
- The restarts of the active phase change the runs of existing simulations, and their stored cases can stop replaying.
- The anatomy gains a field, and the generated anatomy of a machine with a dependency of the scope host gains a function.

## Open questions

None.

## References

| Reference | Link |
|---|---|
| The issue of the dependency of each host | <https://github.com/dokimasia/dokimi/issues/92> |
| The nemesis of LazyFS of Jepsen: "The :value should be a list of nodes you'd like to lose un-fsynced writes on." | <https://jepsen-io.github.io/jepsen/jepsen.lazyfs.html> |
| The storage of the simulator of TigerBeetle: "Read/write faults are distributed between replicas according to ClusterFaultAtlas, to ensure that at least one replica will have a valid copy to help others repair." | <https://github.com/tigerbeetle/tigerbeetle/blob/main/src/testing/storage.zig> |
| The restart policy of a Pod, Kubernetes documentation | <https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#restart-policy> |
| The dependency of a function performer | [RFC-0037](0037-hosts-started-machines-and-the-inputs-of-a-system.md) |
| The restarts of a host in the final phase | [RFC-0043](0043-the-restarts-of-a-host-in-the-final-phase.md) |
