---
rfc: 0026
title: The workloads and checkers of a simulation on a real cluster
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Draft
created: 2026-10-28
updated: 2026-10-28
discussion: none
supersedes: none
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0026: The workloads and checkers of a simulation on a real cluster

## Summary

The simulation of sans-IO machines tests the machines of a system under a runtime that dokimi generates. It cannot test the real binaries of the system over a real network stack, kernel and disk. A new module of dokimi-go, `go.dokimi.dev/cluster`, runs those binaries as the nodes of a cluster of containers or virtual machines on the host of the test. It drives the operations of the workload through the real client of the system, and injects faults through the kernels of the nodes. It records each call as a `history.Interval` on one monotonic clock of the runner. `history.FromIntervals` turns the intervals into the history that the same spec and checkers of the simulation read: linearizability, serializability and snapshot isolation. The operations, the spec, the faults and their bound have the names of the simulation, so one declaration serves both.

## Motivation

Jepsen tests a system that it does not compile: binaries in any language, deployed on nodes, reached over a network. It injects faults with the tools of each node. `jepsen.net` cuts links with iptables and shapes them with `tc` and netem, `jepsen.nemesis.time` installs the C programs `bump-time` and `strobe-time`, which set the clock of a node, and `jepsen.lazyfs` mounts LazyFS under the data of a database to lose its unsynced writes (`jepsen/src/jepsen/net.clj`, `nemesis/time.clj`, `lazyfs.clj`). A run needs a control node, SSH to each node and a Clojure toolchain, and its checkers run in that process.

The simulation of RFC-0025 checks the same histories deterministically, but on the machines of a sans-IO system and the doubles of their dependencies. A defect that only the real system shows escapes it:

- a driver that delivers a result in the same step as its effect
- a client library that retries a write after a timeout and applies it twice
- a store whose sync does not make its writes durable on the disk
- a kernel that reorders the writes of two files

assert-go already checks a history that a real run records. `history.Interval` states a call with a start and an end on one clock and an outcome of `OK`, `Fail`, `Unknown` or pending, and `history.FromIntervals` orders such calls into a `History` (`assert-go/history/interval.go`). The missing part is the runner that deploys the cluster, drives the calls and injects the faults.

## Detailed design

### Components

| Component | Where | Responsibility |
|---|---|---|
| The directive `cluster` | A function of the consumer | The topology: the nodes, their images, commands, ports and data directories |
| The directive `client` | A function of the consumer | One operation of the workload through the real client of the system, with the name of the action of the simulation |
| The generator `sim` | dokimi | Writes `TestCluster` beside `TestSimulation` when the declaration has a function under `cluster` |
| `go.dokimi.dev/cluster` | dokimi-go | The nodes, the faults, the clients, the phases, the intervals and the artifacts of a run |
| `go.dokimi.dev/assert/history` | assert-go | The history from the intervals, and the checkers |

### The topology

```go
// Package cluster runs the real binaries of a system as the nodes of a
// cluster on the host of the test, drives a workload through its client,
// injects faults through the kernels of the nodes, and checks the history
// of the calls.
package cluster

// Kind is the kind of the nodes of a cluster.
type Kind uint8

const (
	// Containers runs each node as a rootless container of podman, with a
	// network namespace of its own. The nodes share the kernel and the
	// realtime clock of the host.
	Containers Kind = iota + 1
	// Machines runs each node as a virtual machine of QEMU with KVM, with a
	// kernel and a realtime clock of its own.
	Machines
)

// Topology is the cluster of one run.
type Topology struct {
	Kind  Kind
	Nodes []Node
}

// Node is one node: its image, its command, its ports and its data.
type Node struct {
	// Name is the name of the node, which faults and logs use.
	Name string
	// Image is an OCI image for Containers, and a disk image for Machines.
	Image string
	// Command is the command that starts the system on the node.
	Command []string
	// Ports are the ports of the node that the clients connect to.
	Ports []int
	// Data is the directory of the node whose writes the disk faults
	// strike. The runner mounts LazyFS there.
	Data string
}
```

The function of `//dokimi:cluster` returns the topology. The runner starts every node, waits for each port, and stops and removes every node and network that it created when the run ends, also after a failure.

### The operations

```go
// Op is one operation of the workload: a call through the real client of
// the system.
type Op struct {
	// Name is the operation in the history. It equals the name of the
	// action of the simulation, so one spec reads both histories.
	Name string
	// Weight is how often a client takes the operation, relative to the
	// others.
	Weight int
	// Final runs the operation once on each key after the faults end.
	Final bool
	// Call chooses its arguments with r and calls the system through the
	// client of the given number. It returns the arguments, the keys, the
	// output and the error. The runner ends ctx at the timeout of the
	// operation.
	Call func(ctx context.Context, client int, r *rand.Rand) (args []any, keys []any, output any, err error)
}
```

Each function of `//dokimi:client` becomes one `Op`. The runner runs each client on a goroutine of its own and records each call as an `Interval`: the start and the end on the monotonic clock of the runner, the arguments, the keys and the outcome. A call that returns an error completes as `Fail` when the error is one that the system states as a refusal, through `//dokimi:client fail=<function>`, and as `Unknown` otherwise, because a timeout or a broken connection leaves the outcome open. A call that has not returned when the run ends is pending. All clients run in the process of the test, so every interval is on one clock, as `FromIntervals` requires.

### The faults

The faults have the names of the faults of the simulation:

| Fault | Containers | Machines |
|---|---|---|
| `partition` in the shapes `isolate`, `halves`, `majorities-ring` and `bridge`, and `repair` | nftables rules in the network namespace of each node | nftables rules in the guest |
| `drop`, `delay` and `duplicate` | `tc` with netem in the network namespace of each node | `tc` with netem in the guest |
| `crash` and `restart` | `podman kill` with `SIGKILL`, and `podman start` | A reset of the virtual machine through the QEMU monitor |
| `pause` | `podman pause`, which freezes the cgroup of the node | `stop` and `cont` through the QEMU monitor |
| `clock-jump` and `clock-strobe` | Refused | `settimeofday` in the guest, as Jepsen's `bump-time` and `strobe-time` set it |
| `lose-unsynced` and `torn-write` | LazyFS under `Data`: the command `lazyfs::clear-cache` through its fault pipe, and an injection of the type `torn-op` in its configuration, which the runner writes before the node starts | The same, in the guest |

The nodes of `Containers` share the realtime clock of the host. A time namespace virtualizes `CLOCK_MONOTONIC` and `CLOCK_BOOTTIME` alone (`time_namespaces(7)`), and a Go binary bypasses libfaketime, because it reads the clock through the vDSO (the README of libfaketime). A topology of `Containers` with a clock fault in its `tolerate` therefore fails before the run with the remedy to use `Machines`.

The faults in force at once never exceed `tolerate`, as in the simulation, and the runner calls a function of `//dokimi:action`, such as the join of a node, as one more kind of fault.

### A run

| Phase | What happens |
|---|---|
| Active | For the duration of the run, every client takes operations, and the nemesis takes faults within `tolerate` at intervals that the seed chooses |
| Recovery | The nemesis ends every fault, restarts every crashed node and waits for its ports |
| Final | Each client runs each final operation |
| Check | `FromIntervals` builds the history, and the runner checks it against the spec with `history.Linearizable` partitioned by key, and with `history.Serializable` or `history.HasSnapshotIsolation` under `isolation` |

```go
// Run runs the workload on the topology for the duration, with the faults
// of opts, and checks the history of the calls. It fails t with the
// record of the checker, and keeps the artifacts of the run.
func Run(t *testing.T, topo Topology, ops []Op, opts ...Option)
```

The seed chooses the operations, their arguments and the schedule of the faults. The real system then decides the order of the events, so a seed does not replay a run, and the runner does not shrink one. The artifacts of each run record what a person needs to read a failure:

- the intervals and the history, in the JSON form of assert-go
- the schedule of the faults, with the time and the target of each
- the log of each node
- the record of the checker, with the anomaly or the call that no linearization explains

The runner can replay the schedule of the faults of a failed run against a fresh cluster, which repeats the faults but not the order of the events.

### The generated test

`TestCluster` runs under the build tag `dokimicluster`, because it needs podman, or QEMU with KVM, on the host, and minutes per run. A plain `go test` leaves it out. A consumer runs it in a job of its own:

```sh
go test -tags dokimicluster -run TestCluster -timeout 30m ./...
```

## Alternatives considered

### Jepsen itself

The consumer writes a Jepsen test in Clojure for the real system.

**Why not:** the workload, the spec and the faults would be written twice, once in Go for the simulation and once in Clojure, and the two would drift. The checkers of Jepsen exist in assert-go, so a runner in Go applies one spec to both histories.

### Clock faults through libfaketime or a time namespace

Containers are the only kind of node, and a clock fault shifts the clock of one container.

**Why not:** a time namespace shifts the monotonic and boot-time clocks alone, so the realtime clock of a container cannot move apart from the host. libfaketime preloads into the C library, and a Go binary reads the clock without it. A node with a clock of its own needs a kernel of its own.

### A cluster on Kubernetes with Chaos Mesh

The runner deploys the system on Kubernetes and injects faults through the resources of Chaos Mesh.

**Why not:** it needs a cluster of Kubernetes in the CI of each consumer, and the faults depend on the operator of Chaos Mesh. The runner on podman and QEMU needs the host of the test alone. A consumer that runs on Kubernetes in production can still run the same workload against a staging cluster, through a topology of its own.

### A proxy for the faults of the network

Toxiproxy runs between the clients and the nodes, and cuts or slows their connections.

**Why not:** a proxy sees the connections of the clients, not the traffic between the nodes, so it cannot partition one node from another. The rules in the network namespace of each node cut every link.

## Drawbacks

- **No replay and no shrinking.** A real run does not repeat. A failure is a history, a schedule of faults and the logs of the nodes, which a person reads.
- **Heavy runs.** A topology of virtual machines needs the memory and the boot time of each guest, and a run lasts its whole duration plus recovery. The tests run in a job of their own.
- **Clock faults need virtual machines.** Containers cannot move a realtime clock apart from the host.
- **Images of the consumer.** The consumer builds an OCI image or a disk image of its system for the runner.
- **A history of wall-clock intervals.** A real run orders calls by their intervals on the clock of the runner. Two calls whose intervals overlap may have taken effect in either order, so the checker of linearizability searches more orders than in the simulation, and needs a budget.
- **LazyFS under the data.** The disk faults need FUSE on the host or in the guest, and LazyFS slows every write of the node.

## Open questions

None.

## References

| What | Where |
|---|---|
| Jepsen, the faults of the network, the clock and LazyFS | <https://github.com/jepsen-io/jepsen/tree/main/jepsen/src/jepsen> |
| Elle, inferring isolation anomalies from histories | <https://arxiv.org/abs/2003.10554> |
| LazyFS, its faults of lost and torn writes | <https://github.com/dsrhaslab/lazyfs> |
| libfaketime, its limits with Go | <https://github.com/wolfcw/libfaketime> |
| time_namespaces(7) | <https://man7.org/linux/man-pages/man7/time_namespaces.7.html> |
| assert-go, the history of intervals | <https://github.com/dokimasia/assert-go/blob/main/history/interval.go> |
| Chaos Mesh | <https://github.com/chaos-mesh/chaos-mesh> |
| Toxiproxy | <https://github.com/Shopify/toxiproxy> |
