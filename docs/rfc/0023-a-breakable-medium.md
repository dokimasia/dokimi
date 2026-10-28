---
rfc: 0023
title: A breakable medium
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

# RFC-0023: A breakable medium

## Summary

A durable subject writes to storage, and its laws are about what survives a crash. dokimi-go gains the module `go.dokimi.dev/medium`: a small filesystem interface with `Sync`, an implementation over the operating system, and an in-memory implementation that keeps synced and unsynced data apart. The in-memory medium can drop unsynced data at a crash, tear a write, corrupt a block on read and run out of space. conform uses it for the claims of the mixin `durable` on a declaration that states the machinery `medium`. A check enumerates every crash point of every sequence of up to three calls, crashes the medium there, restarts the subject on what survives, and checks that each acknowledged durable write is readable and that recovery returns no torn record.

## Motivation

conform has the doors `Restart` and `Crash`. The consumer writes both, and the crash runs one call and crashes the medium of the consumer at one point that `h.Rand()` draws (`conform/handle.go:60-66`). The corpus has one fixture of the door, `durable`, with a medium that the fixture writes itself. Each durable subject needs its own model of what a crash keeps, and each consumer writes it again.

A real disk does not produce on demand the faults that break durable code:

- a write torn by a power loss in the middle of a block
- buffered data lost because no sync ran
- a bit flipped at rest and found on read
- a full disk in the middle of a write

Crash-consistency studies show where the defects are. CrashMonkey and ACE reproduced 24 of 26 crash-consistency bugs reported in Linux filesystems over five years and found 10 new ones, with workloads of three or fewer operations, and most bugs followed a call related to fsync. ALICE found 60 crash vulnerabilities in 11 applications by enumerating the crash states that a model of persistence allows. LazyFS models lost unsynced data with a page cache of its own.

Each such tool needs the subject to accept its medium: a mount point, a block device, or an injected interface. For a check inside `go test`, an injected interface is the only one that needs no root and no kernel module.

## Detailed design

### The module medium

```go
// Package medium is the storage of a durable subject: a filesystem with an
// explicit sync, implemented over the operating system and in memory. The
// in-memory medium can crash, losing every write that no sync covered.
package medium

// FS is the storage of a subject.
type FS interface {
	OpenFile(name string, flag int, perm fs.FileMode) (File, error)
	Remove(name string) error
	Rename(oldName, newName string) error
	MkdirAll(name string, perm fs.FileMode) error
	ReadDir(name string) ([]fs.DirEntry, error)
	Stat(name string) (fs.FileInfo, error)
	// SyncDir makes the entries of the directory durable, such as a file
	// created or renamed in it.
	SyncDir(name string) error
}

// File is an open file of an FS.
type File interface {
	io.Reader
	io.ReaderAt
	io.Writer
	io.WriterAt
	io.Seeker
	// Sync makes the written data of the file durable.
	Sync() error
	Truncate(size int64) error
	Stat() (fs.FileInfo, error)
	Close() error
}

// OS returns the FS of the operating system under root.
func OS(root string) FS

// Memory is an FS in memory that keeps the data that a sync made durable
// apart from the data written since.
type Memory struct{ /* unexported */ }

// NewMemory returns an empty Memory.
func NewMemory() *Memory

// Crash drops every write that no sync covered, and every entry of a
// directory that no SyncDir covered. Each open file is closed.
func (m *Memory) Crash()

// Tear makes the next write of a file keep only its first n bytes when a
// crash follows it before a sync.
func (m *Memory) Tear(n int)

// Corrupt flips one bit of the durable data of the file at the offset, so a
// later read returns the flipped byte.
func (m *Memory) Corrupt(name string, offset int64)

// Fill makes every write after the next n bytes fail with an error that
// wraps syscall.ENOSPC.
func (m *Memory) Fill(n int64)

// Points returns the number of persistence points so far: each Sync and
// each SyncDir.
func (m *Memory) Points() int
```

A consumer writes its subject over `medium.FS` and passes `medium.OS(dir)` in production.

### The door

```go
// Durable sets the factory of a subject over a medium, and the factory of a
// subject that recovers from the medium after a crash. It also sets the
// doors Restart and Crash over the medium.
func Durable[S any](open func(h *Handle[S], m medium.FS) S) Option[S]
```

The same function opens a fresh subject on an empty medium and recovers a subject from a medium after a crash. `Durable` implements the doors of `Restartable` and `Crashable` from it:

- `Restart` opens the subject again on the same medium.
- `Crash` runs the call, calls `Memory.Crash` and opens the subject again on what survived.

An application with `Durable` therefore needs neither `Restartable` nor `Crashable`, and an application that passes both kinds of option fails with the names of the options. `Restartable` and `Crashable` remain for a subject that does not write through `medium.FS`, such as a process that the consumer restarts. A check through `Crashable` keeps the crash point that it draws from `h.Rand()`, as RFC-0008 states, because only the medium of the consumer knows its units. The enumeration of every crash point needs `medium.Memory`.

### The claims

| Claim | Licence | Law | Defect form |
|---|---|---|---|
| `synced-writes-survive` | The mixin `durable` | Each write that returned success and that a sync covered is readable after a crash | `acknowledges-before-sync` |
| `no-torn-record` | The mixin `durable` | Recovery after a torn write returns each record whole or not at all | `trusts-a-partial-record` |
| `corruption-is-reported` | The mixin `durable` with its param `corrupt` | A read of a corrupted record fails with the error that `corrupt` names and returns no value | `ignores-checksum` |
| `full-medium-is-reported` | The mixin `durable` | A write that the medium cannot hold fails with an error, and the state is as before | `half-written-on-full` |

Each claim lists `killable-process-or-breakable-medium` among its machineries. A declaration without `machineries=medium` lists the mixin `durable` as unclaimed with the stance `needs`.

The shape catalog of eidos has the mixin `durable` (`eidos-plugin-shape/spec/mixins/durable.yaml`). Its param `corrupt` refers to the error that a read of a corrupted record reports, so a declaration states with it that the subject checks the integrity of a record on read. A declaration without `corrupt` gets no claim `corruption-is-reported`.

### Enumeration

The probe of `synced-writes-survive` generates every sequence of up to three calls of the writers of the declaration, with inputs from their generators. For each sequence it runs the calls on a fresh medium, counts the persistence points, and for each point replays the sequence on a fresh medium, crashes at that point, opens the subject again and reads back each write that the sequence acknowledged before the point. A sequence of three calls with two syncs each has seven crash points, so a callable with ten sampled inputs gives a bounded number of runs that a check states in its budget.

## Alternatives considered

### afero

`spf13/afero` is a widely used filesystem interface for Go with an in-memory backend.

**Why not:** its in-memory backend has no model of durability: a sync does nothing, so a crash cannot drop unsynced data. Its interface has no sync of a directory. A wrapper over afero would contain the model that `medium.Memory` states, plus a dependency.

### A FUSE filesystem such as LazyFS

The check mounts a filesystem that drops unsynced data, and the subject uses the real filesystem of the operating system under the mount.

**Why not:** a mount needs FUSE, privileges and a process outside the test. It cannot run in the default CI of a consumer. LazyFS and CrashMonkey remain useful for the chaos experiments against running systems.

### Random crash points

The probe crashes at a point that the random source draws, as the door `Crash` does today.

**Why not:** a random point finds a defect at one crash point out of many only by chance. ACE showed that bounded exhaustive sequences of three operations find most crash-consistency defects. Enumerating every point of short sequences is cheaper than running many random long ones.

## Drawbacks

- **The consumer writes against medium.FS.** A subject that calls the `os` package directly gets no claims of durability until it takes an `FS`.
- **A production dependency.** The consumer's production code imports `go.dokimi.dev/medium` for the interface and `OS`. The module imports the standard library alone.
- **A model of persistence.** `medium.Memory` models the guarantees that POSIX states for `fsync`. A filesystem that keeps less, such as one that reorders the sync of a directory, can lose data that the model keeps.
- **Runs that grow with the sequence.** Each crash point is one more run of the sequence. Sequences longer than three calls are not enumerated.

## Open questions

- Does the TypeScript runtime get a medium over the promises of `node:fs`, or does TypeScript refuse the claims of durability?

## References

| What | Where |
|---|---|
| Mohan et al., Finding Crash-Consistency Bugs with Bounded Black-Box Crash Testing, OSDI 2018 | <https://www.usenix.org/conference/osdi18/presentation/mohan> |
| Pillai et al., All File Systems Are Not Created Equal, OSDI 2014 | <https://www.usenix.org/conference/osdi14/technical-sessions/presentation/pillai> |
| LazyFS | <https://github.com/dsrhaslab/lazyfs> |
| TigerBeetle, the storage of its simulator | <https://github.com/tigerbeetle/tigerbeetle/blob/main/src/testing/storage.zig> |
| spf13/afero | <https://pkg.go.dev/github.com/spf13/afero> |
