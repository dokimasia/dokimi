---
adr: 0023
title: The requests of a world depend only on its choices
status: Accepted
date: 2026-11-07
supersedes: RFC-0036 in part
superseded-by: none
rfc: 0036
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0023: The requests of a world depend only on its choices

## Status

Accepted

## Context

Native mode runs the goroutines of a component in a `testing/synctest` bubble. It replays the choices
of a case, and it does not replay the order of the goroutines that one event wakes. A failing case
that no replay reproduces ends as flaky.

After each random case, the property engine runs a prefix case that replays the first choices of
the random case, also when the random case passed. A case that requests other bounds after the same
choices ends the run as flaky. assert-spec names this outcome a flaky outcome that a divergence of
the body ended.

The world of a subject and the run of a system took many bounds from the progress of the
goroutines of the subject:

- A call took a client without an open call, and a cancel took a client with one.
- A fault of a dependency needed a call of a dependency, and a fault of the clock needed a read of
  the wall clock.
- A flip took a file with data and an offset within its size.
- A power loss requested a count for each file and each directory with pending calls.
- A pick of a pool took an index among the keys that the calls produced.
- In a run of sim with a component, a drop and a duplicate took a process with a datagram in flight,
  a delay took one with two, and a corruption took a file of the disks of a process.

The world check `completes-after-recovery` of the corpus fixture `iface/door/hosted` ended flaky in
about 1 run of 20 at GOMAXPROCS=8 and never at GOMAXPROCS=1, and no case of those runs failed. A
traced run showed a random case and its prefix case that took the same 51 steps. After each step
both had the same open calls, fault, dependency calls and persistence points. At the power loss of
step 51, the random case requested a keep count for a pending entry of the disk that the prefix
case did not have. The probable cause is the end of a pause at step 44, which woke the flusher and
a Put at one instant. A cancel had ended that Put during the pause.

The same dependency weakens the replays of a failing case. A replay reproduces a failure only when
it requests the same bounds, so a replay of a racy failure diverges at the first request that reads
the racy state.

## Decision

We will take each bound of a request of a world from the choices of its case and from the state
that the steps of the world set, because the property engine ends a run as flaky when a case
requests other bounds after the same choices, and native mode does not replay the order of the
goroutines.

- The state that the steps set is whether an incarnation or a process runs, the fault in force, the
  faults that the case keeps, and the clients and components that the world spawned. No request
  reads an open call, a dependency call, a read of the wall clock, a datagram in flight, a pending
  write, a file on a disk or a key that a call produced.
- A call of an operation is enabled while the subject runs. The case chooses a client among all
  clients, and the call goes to the first client without an open call from it on, in the order of
  the clients and around. The step has no effect while every client has an open call.
- A cancel is enabled while the subject runs. It ends the open call of the first client with one
  from the chosen client on, and it has no effect without an open call.
- A fault of a dependency and a fault of the clock are enabled for a case that keeps them. At the end
  of a case the world records the faults of the dependencies as vacuous for a subject that made no
  call of a dependency, and the faults of the clock for a subject that never read its wall clock. It
  does not count their strikes.
- A flip is enabled for a case that keeps it. The case chooses a share of the files with data and a
  share of the size of a file, and the step maps both onto the disk when it runs. The step has no
  effect without a file with data.
- A pick of a pool that is not hot takes the key at a share of the pool that the case chooses, which
  shrinks to the first key.
- A power loss requests one list of shares for the crash, under the label `keep`. The crash takes
  the shares in the order in which it asks for the count of each name, the directories and then the
  files, each in the order of their names. A name keeps its share of its pending calls, so the crash
  keeps an independent count of each name, as ALICE explores and as gosim keeps a random subset. A
  name past the end of the list keeps none of its pending calls. The list shrinks to the empty list,
  the crash state of a crash without a torn write.
- In a run of sim with a component, a call and a cancel take their client as the world of a subject
  does. A drop, a delay and a duplicate choose a process among all processes, and have no effect
  without enough datagrams in flight to it. A corruption chooses a process, a share of its files and
  a share of the size of a file. An action of the consumer in such a run takes its bounds from the
  case alone.
- Both modes take their bounds this way, so a check makes the same requests in Tasks mode and in
  native mode. A run of sim without a component keeps its forms, because its machines replay whole.

## Alternatives Considered

### A fixed order of the goroutines that the host wakes

The host would wake the goroutines of a pause, a hold and a broadcast one at a time, each until
quiescence. It lost because it removes only the races that the host starts. Goroutines of one subject
that a timer or a channel wakes at one instant still run in an order that the case does not
choose.

### An option of the property engine for a body whose requests vary

A run would not end as flaky when a passing case diverges. It lost because the option changes the
definition of assert-spec in every language, and a replay of a racy failure would still diverge at
its first racy request.

### The subject Flushed without the option Blocking

The world would strike no pause and no hold on Flushed, so no event would wake two of its goroutines
at one instant. It lost because the corpus needs a subject in native mode under Blocking, the
component that issue 58 describes.

## Consequences

**Positive:**

- A prefix case and a replay request the bounds of the case that they replay. A world check in
  native mode then ends as flaky only for a failing case that its replays do not reproduce.
- A replay of a racy failure differs from the failing case only in the order of the goroutines, so
  the replays reproduce such a failure more often.

**Negative:**

- A step without a target has no effect. A case can spend steps on a call while every client has an
  open call, on a cancel without an open call, and on a fault of a resource that the subject never
  uses. A fault of a dependency or of the clock in force keeps the other faults of the world from
  striking until its end, also for a subject that never uses the resource.
- A name past the end of the list of a power loss keeps none of its pending calls, so a crash keeps
  the pending calls of the first names in the order of the crash more often than those of later
  names.
- The stored cases and the golden files of the world checks change once.
- The world cannot check that an action of the consumer keeps the rule. A consumer action that reads
  the progress of a component makes its run flaky in native mode.

## References

| What | Where |
|---|---|
| The modes, the world of a subject and the components of a system | [RFC-0036](../rfc/0036-components-that-run-goroutines.md) |
| The replays of a failing case and the flaky outcome of a divergence | <https://github.com/dokimasia/assert-spec/blob/main/docs/rfc/0057-the-replays-of-a-failing-case.md> |
| The component that a goroutine flushes on a timer | <https://github.com/dokimasia/dokimi/issues/58> |
