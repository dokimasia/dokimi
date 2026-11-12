---
rfc: 0023
title: A breakable medium
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-28
updated: 2026-11-03
discussion: none
supersedes: none
superseded-by: RFC-0050 in part
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0023: A breakable medium

## Summary

A durable subject writes to storage, and its laws are about what survives a crash. dokimi-go gains the module `go.dokimi.dev/medium`: a small filesystem interface with `Sync`, an implementation over the operating system, and an in-memory implementation that keeps synced and unsynced data apart. The in-memory medium can drop unsynced data at a crash, crash at any persistence point, tear a write, corrupt a bit on read and run out of space. conform opens a subject over such a medium through the option `Durable`, and checks the claims of the mixin `durable` on a declaration that states the machinery `medium`. A check enumerates every crash point of every sequence of up to three calls, crashes the medium there, opens the subject again on what survives, and checks that each acknowledged write is readable and that recovery returns no torn record. The simulation of RFC-0025 strikes the same faults at the media of its processes.

## Motivation

conform has the doors `Restart` and `Crash`. The consumer writes both, and the crash runs one call and crashes the medium of the consumer at one point that `h.Rand()` chooses (`conform/handle.go:72-76`). The corpus has one fixture of the door, `iface/door/durable`, with a medium that the fixture writes itself. Each durable subject needs its own model of what a crash keeps, and each consumer writes it again.

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

// ErrCrashed is the error of each call of a crashed medium, and of each call
// of a file that the medium opened before a crash.
var ErrCrashed = errors.New("medium: the medium crashed")

// FS is the storage of a subject. A name is a slash-separated path relative
// to the root of the FS, as fs.ValidPath states it.
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
// apart from the data written since. It is safe for concurrent use.
type Memory struct{ /* unexported */ }

// NewMemory returns an empty Memory.
func NewMemory() *Memory

// Crash drops every write that a sync did not cover, and every entry of a
// directory that a SyncDir did not cover. Each open file is closed, and its
// later calls fail with ErrCrashed. The medium then serves the calls of a new
// process. A crash after CrashAt does not drop anything more.
func (m *Memory) Crash()

// CrashAt makes the medium crash in place of its persistence point
// point+1: that Sync or SyncDir does not persist anything and fails with
// ErrCrashed, and every later call of the medium and of its files fails with
// ErrCrashed until Crash.
func (m *Memory) CrashAt(point int)

// Tear makes the next write of a file keep only its first n bytes when a
// crash follows it before a sync.
func (m *Memory) Tear(n int)

// Corrupt flips the lowest bit of the durable byte of the file at the
// offset, so a later read returns the flipped byte. An offset outside the
// durable data of the file corrupts nothing.
func (m *Memory) Corrupt(name string, offset int64)

// Fill makes every write after the next n bytes fail with an error that
// wraps syscall.ENOSPC. A write that crosses the limit writes the bytes that
// fit and returns their count beside the error. A negative n lifts the limit.
func (m *Memory) Fill(n int64)

// Points returns the number of persistence points so far: each Sync and
// each SyncDir that returned nil.
func (m *Memory) Points() int

// Written returns the number of bytes that the writes of the medium wrote
// so far.
func (m *Memory) Written() int64
```

A consumer writes its subject over `medium.FS` and passes `medium.OS(dir)` in production.

`Memory` models the guarantees that POSIX states for `fsync`:

- A write, a truncate and a create change the pending state of a file. `Sync` makes the pending state of the file durable.
- A create, a remove and a rename change the pending entries of their directories. `SyncDir` makes the pending entries of the directory durable. A rename is durable once both of its directories are synced.
- A crash drops the pending state of each file and the pending entries of each directory. A torn write keeps the pending writes of its file before it and its own first bytes.

`OS` makes `SyncDir` open the directory and sync it. Windows cannot sync a directory, so `SyncDir` of `OS` does nothing on Windows, as the filesystems of Pebble and RocksDB do.

### The door

```go
// Durable sets the door Open of the application to open, which opens a
// subject over a medium and recovers a subject from a medium after a crash.
// The subjects of the checks of the application are opened over new media in
// memory, and the doors Restart and Crash work over the medium of a subject.
func Durable[S any](open func(h *Handle[S], m medium.FS) S) Option[S]

// Open returns a subject over the medium m through the door Open of the
// application, or through the door Open of a planted defect.
func (h *Handle[S]) Open(m medium.FS) S

// FreshOpen returns a subject over the medium m through the door Open of the
// application, also under falsification, as Fresh returns a fresh subject.
func (h *Handle[S]) FreshOpen(m medium.FS) S

// AddPoints adds n to the points of the medium that the check enumerated,
// which the record of the check states.
func (h *Handle[S]) AddPoints(n int)
```

`Doors` gains the door `Open func(h *Handle[S], m medium.FS) S`, which a planted defect sets through `Check.Doors`. With the door Open:

- `Subject`, `New` and `Fresh` open each subject over a new `medium.Memory`, and the handle keeps the medium of each subject that it opened.
- `Restart(s)` opens the subject again on the medium of `s`.
- `Crash(s, call)` runs the call, calls `Memory.Crash` on the medium of `s` and opens the subject again on what survived.

An application with `Durable` does not need `Restartable` or `Crashable`, and `Run` fails with an error that wraps `ErrOption` and names `conform.Durable` and the other option for an application that passes both kinds. `Restartable` and `Crashable` remain for a subject that does not write through `medium.FS`, such as a process that the consumer restarts. A check through `Crashable` keeps the crash point that it takes from `h.Rand()`, as RFC-0008 states, because only the medium of the consumer knows its units. The enumeration of every crash point needs `medium.Memory`.

### The claims

The mixin `durable` is on the writer of a declaration. Each claim calls the writer and its partner `read`, the one other callable of the declaration whose read pairs with the writer by the rule that pairs a read with a write. The spec gains the partner mode `paired` for it:

```go
// ModePaired is the one other callable of the declaration whose read pairs
// with the licensed callable: the licensed callable writes, and the partner
// reads.
ModePaired PartnerMode = "paired"
```

| Claim | Licence | Law | Reason | Probe | Defect form |
|---|---|---|---|---|---|
| `synced-writes-survive` | The mixin `durable` | Each write that returned success is readable after a crash | `lost-acknowledged-write` | `survives-every-crash-point` | `acknowledges-before-sync` |
| `no-torn-record` | The mixin `durable` | Recovery after a torn write returns each record whole or not at all | `torn-record` | `whole-at-every-crash-point` | `trusts-a-partial-record` |
| `corruption-is-reported` | The mixin `durable` with its param `corrupt` | A read of a corrupted record fails with the error that `corrupt` names and returns no value | `unreported-corruption` | `reports-each-flipped-bit` | `ignores-checksum` |
| `full-medium-is-reported` | The mixin `durable` | A write that the medium cannot hold fails with an error, and the state is as before | `half-written-on-full` | `refuses-on-a-full-medium` | `hides-a-full-medium` |

Each claim lists `killable-process-or-breakable-medium` as its machinery. A declaration without `machineries=[medium]` lists the mixin `durable` as unclaimed with the stance `needs`, which the derivation computes from the machinery that the claims lack.

The shape catalog of eidos has the mixin `durable` (`eidos-plugin-shape/spec/mixins/durable.yaml`). Its param `corrupt` refers to the error that a read of a corrupted record reports, so a declaration states with it that the subject checks the integrity of a record on read. A declaration without `corrupt` gets no claim `corruption-is-reported`.

### The probes

```go
// Durable are the calls of a durable writer and of its reader. Write writes
// the sample i of the writer to the subject s, and Read reads the key of the
// sample i from s. Want returns the value that the sample i writes. Key
// returns the key of the sample i. Two samples with one key replace each
// other.
type Durable[S, V any] struct {
	Write func(ctx context.Context, s S, i int) error
	Read  func(ctx context.Context, s S, i int) (V, error)
	Want  func(i int) V
	Key   func(i int) any
}

func SurvivesEveryCrashPoint[S, V any](h *conform.Handle[S], d Durable[S, V])
func WholeAtEveryCrashPoint[S, V any](h *conform.Handle[S], d Durable[S, V])
func ReportsEachFlippedBit[S, V any](h *conform.Handle[S], d Durable[S, V], corrupt error)
func RefusesOnAFullMedium[S, V any](h *conform.Handle[S], d Durable[S, V])

// Unsynced returns a medium over m whose syncs do not persist anything.
// SyncDir and the Sync of each file that it opens return nil. Every other call
// goes to m.
func Unsynced(m medium.FS) medium.FS

// Unfilled returns a medium over m that hides a full medium. When m refuses a
// write of a file that it opens with an error that wraps syscall.ENOSPC, the
// write returns the length of its data and no error. Every other call goes to
// m.
func Unfilled(m medium.FS) medium.FS
```

A derived check has two samples of each input, so the probes enumerate every sequence of one to three writes over the samples: fourteen sequences. Each probe opens each subject over a new `medium.Memory` through `h.Open`.

- `SurvivesEveryCrashPoint` runs each sequence once and counts its persistence points `n`. For each point `p` from 0 to `n` it runs the sequence again on a new medium with `CrashAt(p)`, crashes the medium, opens the subject again and reads each key that a write acknowledged. Each read returns the value of the last acknowledged write of its key, or the value of a later write of the key in which the crash fell.
- `WholeAtEveryCrashPoint` runs the same crash points, and tears the first write of the call in which the crash falls to one byte. It reads each key that the sequence writes. Each read returns a value that a write of the sequence wrote to its key, or an error.
- `ReportsEachFlippedBit` runs each sequence and lists the bytes of each durable file. For each byte it runs the sequence again on a new medium, flips the lowest bit of the byte with `Corrupt` and opens the subject again. Each read returns the value of the last acknowledged write of its key, or an error that wraps `corrupt` beside the zero value.
- `RefusesOnAFullMedium` runs each sequence and measures with `Written` the bytes of its last write. For each limit below those bytes it runs the sequence again with `Fill` before the last write. The last write fails with an error, and a read of the key of each sample returns the value before the last write.

Each probe passes one point to `h.AddPoints` for each crash point, each flipped byte and each limit of a fill that it runs. The record of the check has the sum of the points in its field `points`. The schema of the suite record gains the field.

### The planted defects

The generator writes each planted defect over the subject that the door Open of the application opens:

- `acknowledges-before-sync` opens the subject over `probe.Unsynced(m)`. The `Sync` and `SyncDir` of that medium return nil and do not persist anything.
- `trusts-a-partial-record` writes the zero value under the key before the first value that the defect writes, and writes each later value once. A crash between the two writes leaves the zero value under a key without an acknowledged write.
- `ignores-checksum` returns the zero value with a nil error for a read that fails with the error of `corrupt`.
- `hides-a-full-medium` opens the subject over `probe.Unfilled(m)`. A write of that medium that m refuses for lack of space returns the length of its data and no error. A write that does not fit on the medium then succeeds.

The corpus applies every check of a declaration to each planted defect, and a defect proves its check only when it passes the other checks. So each planted defect fails the check of its own claim alone. A defect that writes the zero value before every value fails three checks. Its two writes let a crash replace an acknowledged value, and a full medium between them leaves the zero value under a key that had a value. `trusts-a-partial-record` writes the zero value only before the first value. Only a key without an acknowledged write then has the zero value, and only `no-torn-record` reads such a key. `refuses-on-a-full-medium` compares values, and the zero value under a new key reads as the zero value of an absent key.

### The simulation

A process of the simulation of RFC-0025 persists through a `*medium.Memory` that the consumer passes among its dependencies. The world strikes the faults of the medium there:

- A crash of a process calls `Crash` on each medium among its dependencies, so a process loses each write that a sync did not cover.
- The fault `torn-write` calls `Tear` on a medium of a running process, with a count of bytes that the case chooses.
- The fault `corruption` calls `Corrupt` on a durable file of a medium, at an offset that the case chooses.
- The fault `full-disk` calls `Fill` on a medium with a limit that the case chooses, and the act `restore-disk` lifts it.

Each fault of the medium counts against the param `tolerate` of the directive `sim`, as the other faults do.

### Other languages

This RFC proposes the medium of Go alone. A TypeScript declaration with the mixin `durable` lists it as unclaimed with the stance `needs`, because no TypeScript application supplies the machinery.

## Alternatives considered

### afero

`spf13/afero` is a widely used filesystem interface for Go with an in-memory backend.

**Why not:** its in-memory backend has no model of durability: a sync does nothing, so a crash cannot drop unsynced data. Its interface has no sync of a directory. A wrapper over afero would contain the model that `medium.Memory` states, plus a dependency.

### A FUSE filesystem such as LazyFS

The check mounts a filesystem that drops unsynced data, and the subject uses the real filesystem of the operating system under the mount.

**Why not:** a mount needs FUSE, privileges and a process outside the test. It cannot run in the default CI of a consumer. LazyFS and CrashMonkey remain useful for the chaos experiments against running systems.

### Random crash points

The probe crashes at a point that the random source chooses, as the door `Crash` does.

**Why not:** a random point finds a defect at one crash point out of many only by chance. ACE showed that bounded exhaustive sequences of three operations find most crash-consistency defects. Enumerating every point of short sequences is cheaper than running many random long ones.

## Drawbacks

- **The consumer writes against medium.FS.** A subject that calls the `os` package directly gets no claims of durability until it takes an `FS`.
- **A production dependency.** The consumer's production code imports `go.dokimi.dev/medium` for the interface and `OS`. The module imports the standard library alone.
- **A model of persistence.** `medium.Memory` models the guarantees that POSIX states for `fsync`. A filesystem that keeps less, such as one that reorders the sync of a directory, can lose data that the model keeps.
- **Runs that grow with the sequence.** Each crash point is one more run of the sequence. Fourteen sequences of up to three writes with two syncs each give 82 crash points for a claim. Sequences longer than three calls are not enumerated.
- **Two samples of each input.** The probes write the two samples of the derivation, so a defect that shows only for a third key needs a check of the consumer.

## Open questions

None.

## References

| What | Where |
|---|---|
| Mohan et al., Finding Crash-Consistency Bugs with Bounded Black-Box Crash Testing, OSDI 2018 | <https://www.usenix.org/conference/osdi18/presentation/mohan> |
| Pillai et al., All File Systems Are Not Created Equal, OSDI 2014 | <https://www.usenix.org/conference/osdi14/technical-sessions/presentation/pillai> |
| LazyFS | <https://github.com/dsrhaslab/lazyfs> |
| TigerBeetle, the storage of its simulator | <https://github.com/tigerbeetle/tigerbeetle/blob/main/src/testing/storage.zig> |
| Pebble, the filesystem of its storage | <https://github.com/cockroachdb/pebble/tree/master/vfs> |
| spf13/afero | <https://pkg.go.dev/github.com/spf13/afero> |
