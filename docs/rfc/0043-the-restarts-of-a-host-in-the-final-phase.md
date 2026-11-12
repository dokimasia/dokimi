---
rfc: 0043
title: The restarts of a host in the final phase
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-08
updated: 2026-11-08
discussion: https://github.com/dokimasia/dokimi/issues/80
supersedes: RFC-0039 in part
superseded-by: RFC-0051 in part
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0043: The restarts of a host in the final phase

## Summary

We propose that the final phase starts each host that exits again, as a supervisor starts a process that exits. The delays are those of the CrashLoopBackOff of the kubelet:

- The first exit of a host in the final phase restarts the host at once.
- Each later exit waits before the restart. The delay is 10 s for the second restart and doubles for each later one, up to 300 s.
- A host that runs for more than 600 s after its last restart starts again at once.

## Motivation

`Recover` starts each crashed host and each exited host once, when the final phase starts. The final phase then delivers the pending events and advances the time until the bound of each convergence has passed (`sim/run.go:667-693` and `sim/run.go:726-748` of dokimi-go at sim v0.10.1). It does not start a host. A host whose machine stops with an error in the final phase remains down until the case ends.

The simulation of sharder fails a case this way:

- One host starts the namespace, and a second host joins. The action `recreate` deletes the scopes of the namespace in the last step of the active phase.
- In the final phase the term of the root of the first host ends, and the root creates the namespace again with a new incarnation. The agent of the first host stops with the error of the incarnation mismatch, and the host exits at 0 s of the run.
- The log lists `exit:host-1` at each later step, up to the end of the final phase at 9m40s. The convergence fails with "routes: the router has stopped", because the router of the first host stopped with its host.

sharder converges in this case in a deployment, because a supervisor starts the process again and its agents join the new incarnation.

The restarts of a machine that stops with an error after each start need a separate bound. If a supervisor started it at once after each exit, the machine would take a step at each start, and the final phase would end at its limit of steps without an advance of the time.

## Detailed design

### Components

| Component | Where | Change |
|---|---|---|
| The supervisor of the final phase | `go.dokimi.dev/sim`, `run.supervise` | Arms the restart of each host that exits in the final phase |
| The delays of the restarts | `go.dokimi.dev/sim`, `restartBase` and `restartCap` | 10 s and 300 s |

### The restarts

```go
// restartBase is the delay of the second restart of a host in the final
// phase, and restartCap the longest delay of a restart. They are the defaults
// of the CrashLoopBackOff of the kubelet.
const (
	restartBase = 10 * time.Second
	restartCap  = 5 * time.Minute
)

// supervise arms a restart of each host that has exited in the final phase, as
// a supervisor starts a process that exits. The first restart of a host is due
// at once, and each later one after the delay of its backoff on the monotonic
// clock of the host. A host that ran for more than twice restartCap after its
// last restart starts again at once.
func (r *run[E, F, S]) supervise(backoffs map[*Host[E, F]]*backoff)
```

`settleDown` calls `supervise` before each of its steps. A restart is a timer of the host, so `World.Next` counts it, and an advance of the final phase moves the time to it. The timer starts the host as `Recover` starts it: each process that `Spawn` added takes a new machine in a new epoch and its start.

| Exit of the host in the final phase | Delay of the restart |
|---|---|
| The first | None |
| A later one, at most 600 s after the last restart | 10 s for the second restart, then twice the delay of the restart before it, at most 300 s |
| A later one, more than 600 s after the last restart | None, and the next delay is 10 s again |

The delays are on the monotonic clock of the host, from its exit to its restart. A restart is no fault of the run. The final phase strikes no fault, and `Faults` lists the exited host until its restart.

### Failure handling

| Failure | Behaviour |
|---|---|
| A machine that stops with an error at each start | The host restarts at 0 s, 10 s, 30 s, 70 s, 150 s and 310 s after its first exit, and then every 300 s. The final phase ends within its limit of steps |
| A host that exits after the bound of each convergence has passed | The final phase ends, and the host remains down |
| A host whose restart is due after the bound | The final phase ends at the bound without the restart |

### Migration

- A simulation whose host exits in the final phase takes the steps of its restart. Its stored cases can stop replaying the same run.
- The active phase does not change. The action `restart` starts a host there, within the bound of `Tolerate`.
- The fixtures of the corpus pass unchanged. No host of them exits in the final phase.

### Verdicts

| Item | Verdict | Reason | Reverses the verdict |
|---|---|---|---|
| A restart of a host that exits in the final phase | in | A deployment starts the process again, and the convergence of a system includes that start | none known |
| The delays of the kubelet | in | A delay that doubles keeps the restarts of a machine that stops at each start to a count that grows with the logarithm of the bound. The kubelet documents 10 s and 300 s | A consumer whose deployment starts a process on another schedule, and whose convergence depends on that schedule |
| The first restart at once | in | The kubelet starts a container without a backoff at once, as `Recover` starts an exited host | none known |
| The reset after 600 s | in | The kubelet forgets the backoff of a container that ran for more than twice the longest delay | none known |
| A field of `Simulation` for the delays | out | Each measured exit in a final phase restarts once, so no consumer depends on the delays | A consumer whose convergence depends on the delays of its supervisor |
| Restarts in the active phase | out | The action `restart` starts a host there, at a step that the case chooses | none known |

## Alternatives considered

### A restart at once after each exit

The final phase could start an exited host at once after each exit, as `Recover` does. A machine that stops at each start would then take a step at each start without an advance of the time, and the case would fail at the limit of steps.

### An intensity limit

An Erlang/OTP supervisor restarts a child at once and gives up after more than `MaxR` restarts within `MaxT` seconds. The final phase could end the restarts of a host the same way. A host that its supervisor gave up would then remain down, as in sim v0.10.1. The limit would also need values that no deployment of the measured systems documents.

### A call of `Recover` at each advance

The final phase could call `Recover` before each advance. It would start every exited host at the time of the next timer, whatever its count of exits. A machine that stops at each start would then take a step at each advance.

## Drawbacks

- A machine that stops with an error at each start takes the steps of 6 restarts within a bound of 9m40s, where sim v0.10.1 took none.
- The delays are fixed. A system whose convergence within a short bound needs a fast second restart fails after its second exit.
- `settleDown` reads the state of every host before each of its steps.

## Open questions

None.

## References

| What | Where |
|---|---|
| Issue 80, a host that exits during the final phase remains down | <https://github.com/dokimasia/dokimi/issues/80> |
| RFC-0039, hosts that exit, the steps of a case and performers of a pointer | 0039-hosts-that-exit-the-steps-of-a-case-and-performers-of-a-pointer.md |
| The restart policy of a Pod, Kubernetes documentation | <https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#restart-policy> |
| `doBackOff` of the kubelet | <https://github.com/kubernetes/kubernetes/blob/master/pkg/kubelet/kuberuntime/kuberuntime_manager.go> |
| `flowcontrol.Backoff` of client-go | <https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/client-go/util/flowcontrol/backoff.go> |
| Supervisor principles of Erlang/OTP | <https://www.erlang.org/doc/system/sup_princ.html> |
