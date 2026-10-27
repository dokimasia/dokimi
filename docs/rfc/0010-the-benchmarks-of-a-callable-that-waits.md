---
rfc: 0010
title: The benchmarks of a callable that waits for another call
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-27
updated: 2026-10-27
discussion: https://github.com/dokimasia/dokimi/issues/33
supersedes: RFC-0006 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0010: The benchmarks of a callable that waits for another call

## Summary

A callable with the mixin `waits` of the shape catalog can wait until another call changes the state of its subject. A campaign that waits until the leader resigns is such a callable. A benchmark and the check of an allocation ceiling repeat the calls of a callable on one subject, so neither ends for such a callable. The suite generator leaves the callable out of every benchmark and every ceiling that it writes. It also refuses a budget of the callable. The generated suite marks the call with the field `Waits` of conform. conform then reports a problem of the application when a benchmark or a ceiling would repeat such a call. The checks and the fuzz target of the callable remain, because each of their cases builds a fresh subject.

## Motivation

`conform.Bench` builds one subject and takes 100 calls from the call generator of the callable. It repeats the calls in their order on that subject, with the context of the benchmark (`conform/bench.go:64-73` and `:118-133` at conform v0.7.0). `BenchParallel` repeats them on one subject from GOMAXPROCS goroutines. The check of an allocation ceiling builds a subject for each case, and `assert.MaxAllocs` then calls the callable about 100 times on that subject (`conform/budget.go:175`). A callable that waits until another call releases a resource waits from its second call on, because neither the benchmark nor the check makes that other call.

The test of a subject declares `Benchmark<Subject>ConformsTo<X>` with a level for each callable. The `Campaign` of a store of sharder waits until the caller leads the election of its scope. A term lasts until `Resign`, until the context of the campaign ends, or until the store ends the term. A probe of the benchmark of the memory store of sharder waited in its warm-up for more than 3 minutes. The generated benchmark of a lock whose `Acquire` waits, written by dokimi v0.3.0, waits in the second call of its level `Acquire`:

```text
goroutine 12 [select]:
example.com/locks/memory.(*Lock).Acquire(0x336afb40030, {0xa1fd78, 0x336afbf20a0})
	memory/memory.go:30 +0x7a
go.dokimi.dev/conform.repeat[...]({0xa1fd78?, 0x336afbf20a0}, {0xa1e8c8, 0x336afb40030}, {0x336afbe5c08, 0x64, 0x2}, 0x336afc13d30?)
	conform@v0.7.0/bench.go:123 +0xa7
go.dokimi.dev/conform.Bench[...]({0xa201b8, 0x336afcf0608}, {{0x70f2a5, 0x4}, {0x0, 0x0, 0x0}, {0x336afbe8180, 0x2, 0x2}, ...}, ...)
	conform@v0.7.0/bench.go:71 +0x346
```

A waiting benchmark ends late, or without a message that helps:

- The testing package of Go 1.27.3 stops the alarm of `-timeout` before it runs the benchmarks (`testing/testing.go:2605` and `:2618`).
- The go command kills the test binary one minute after the timeout. At the default timeout of 10 minutes, a waiting benchmark runs for 11 minutes (`cmd/go/internal/test/test.go:851` and `testflag.go:69`).
- A subject without a timer or another goroutine fails at once with `fatal error: all goroutines are asleep - deadlock!`. The message does not mention the callable.

So a consumer applies the suite of such a declaration by hand and leaves the benchmark out. The memory store of sharder calls `RunStoreSuiteFromHandle` in `TestStoreConformance`, so it has no generated test of the subject and no fuzz target of the subject.

Every implementation of such a callable waits, so the wait is a property of the declaration and not of one subject. The shape catalog states such properties as mixins. The suite generator already reads the mixins `concurrent` and `concurrent-readers` to decide which callables a parallel benchmark measures, and `BenchParallel` refuses a call without `Concurrent`.

## Detailed design

### Invariants

| Invariant | Enforced by |
|---|---|
| No benchmark that dokimi writes measures a callable with the mixin `waits` | The suite generator |
| No ceiling that dokimi writes applies to a callable with the mixin `waits` | The suite generator, and `BadBudget` at a budget of the callable |
| A benchmark of a call with `Waits` measures nothing and reports a problem | `Bench` and `BenchParallel` of conform |
| A budget of a call with `Waits`, and a ceiling that `Ceilings` adds for the call, are problems of the application | The problems of an application in conform |
| The checks and the fuzz target of a callable with the mixin remain | The test of a subject, and the test of the suite |

### Components

| Component | Change |
|---|---|
| The shape catalog of eidos | The documentary mixin `waits` of eidos-plugin-shape v0.1.0 |
| `conform.Call[S]` | The field `Waits` |
| `conform.Bench` and `conform.BenchParallel` | A problem for a call with `Waits` |
| The problems of an application | A budget of a call with `Waits`, and a ceiling of `Ceilings` for such a call |
| The suite generator | `Waits` of each call, the levels of the benchmarks, the ceilings of the test of the suite, and a case of `BadBudget` |
| The template of the suites | The field `Waits` of each call with the mixin |
| The Go guide | A row for the mixin `waits` |

### The mixin `waits`

The mixin is documentary, as `deprecated` is. It documents the callable and licenses no check. A declaration states it on a method or a function:

```go
// Campaign waits until the caller leads the election of s, or until ctx ends.
//
//dokimi:mixin waits
Campaign(ctx context.Context, s shard.Scope) (Leadership, error)
```

eidos-plugin-shape v0.1.0 states this spec:

```yaml
name: waits
form: mixin
documentary: true
claim: >
  The callable can wait until another call changes the state of its
  subject, such as a campaign that waits until the leader resigns. A
  generator of checks does not repeat the callable on one subject without
  that other call.
observation: >
  Nothing. The mixin documents the callable and licenses no check.
falsifiability: >
  Nothing turns red, because the mixin licenses no check.
counterexamples:
  refused: >
    a benchmark that repeats the callable on one subject, whose second call
    waits for a call that the benchmark never makes
```

The mixin `timeout` states that a call returns once its context ends. A declaration states both mixins where both claims are true.

### The field `Waits`

```go
type Call[S any] struct {
	// The fields Calls, Over, Callable, Budget and Concurrent as before.

	// Waits reports that a call can wait until another call changes the
	// state of the subject, as the mixin waits states. A benchmark and the
	// check of a ceiling repeat the calls of a callable on one subject, so
	// a benchmark of the call and a ceiling of the call are problems of the
	// application.
	Waits bool
}
```

conform reports these problems. It stops the test or the benchmark of the application before it runs a check:

| Entry point | Problem |
|---|---|
| `Bench` and `BenchParallel` | `the call of Lock.Acquire waits for another call, and a benchmark repeats it on one subject` |
| `Run`, `RunBudgets`, `Fuzz`, `Bench` and `BenchParallel` | `the budget of Lock.Acquire repeats a call that waits for another call`, for a call with `Waits` and a budget |
| The same entry points | `the ceiling of Lock.Acquire repeats a call that waits for another call`, for a ceiling that `Ceilings` adds for such a call |

`Run` and `Fuzz` apply the checks of a call with `Waits` as before, because each check and each case builds a fresh subject.

### The suite generator

The generator sets `Waits` of a call where the callable has the mixin `waits`, as it sets `Concurrent` where the callable has `concurrent` or `concurrent-readers`. The suite writes the field into the call of the callable:

```go
{
	Calls:    lockAcquireCalls,
	Over:     newLockOver,
	Callable: LockChecks.Acquire.Callable,
	Waits:    true,
},
```

The generated tests change for such a callable:

| Generated test | Change for a callable with `waits` |
|---|---|
| The test of a subject, such as `TestLockConformsToLock` | None. Each check builds a fresh subject |
| The benchmark of a subject, such as `BenchmarkLockConformsToLock`, and its parallel form | No level of the callable. A test without a level declares no benchmark |
| The fuzz target of the callable for a subject, such as `FuzzLockConformsToLockAcquire` | None. Each case builds a fresh subject |
| The allocation test of a subject | No ceiling of the callable, because `BadBudget` refuses its budget |
| The benchmark of the test of the suite, such as `BenchmarkLockSuite`, and its parallel form | No level of the callable |
| The allocation test of the test of the suite, such as `TestLockSuiteBudgets` | No ceiling of 0 for the callable |
| The benchmark, the parallel benchmark and the allocation test of the suite of a function with the mixin | Not declared |

`BadBudget` reports each budget of a callable with the mixin `waits` at its directive `budget`, also the budget of a subject that refers to the callable. No suite applies such a budget, so the meaning of the code does not change.

### The contract in each language

| Part | Go | TypeScript |
|---|---|---|
| The mixin | `//dokimi:mixin waits` | The same directive |
| The field of a call | `Waits bool` of `Call[S]` | `waits` of `Call<S>` |

The TypeScript runtime and the template tree of its suites do not exist yet. Their parts take these names when they do.

### Failure handling

| Failure | Result | Recovery |
|---|---|---|
| A callable that waits has no mixin `waits` | Its benchmark waits until the go command kills the test binary, or the runtime reports a deadlock | State the mixin |
| A budget of a callable with the mixin | `BadBudget` at the directive, and no conformance test of a subject with the budget | Remove the budget |
| A hand-written benchmark of a call with `Waits` | `Bench` reports the problem and stops the benchmark | Remove the benchmark |
| An application adds a ceiling for a call with `Waits` | The application reports the problem and runs no check | Remove the ceiling |
| The mixin on a callable that does not wait | The callable has no benchmark and no budget | Remove the mixin |
| An eidos without the mixin | eidos refuses the directive with `EID-0068` | Upgrade eidos |

### Verification

- **conform.** The specs of `Bench` and `BenchParallel` cover the problem of a call with `Waits`. The spec of the problems of an application covers a budget of such a call and a ceiling of `Ceilings` for it. The specs of `Run` and `Fuzz` cover a call with `Waits` whose checks run.
- **The suite generator.** The specs of the plugin cover `Waits` of a call, the levels of each benchmark, the ceilings of the test of the suite and the new case of `BadBudget`. The golden files of a fixture with a waiting callable pin the field `Waits` and the tests without its levels.
- **The corpus.** The fixture `iface/mixin/waits` declares a lock whose `Acquire` waits, with a control that releases the lock when a lease ends. `go test -run '^$' -bench . -benchtime 1x` of the fixture ends.

### Migration

The change adds a mixin, a field, the cases of existing problems and a case of an existing code. A declaration without the mixin generates the same files as before. eidos-plugin-shape v0.1.0 contains the mixin. A generated suite with the field `Waits` requires a release of conform with the field.

## Alternatives considered

### A param of the directive `subject`

The directive `subject` takes a param that leaves out the benchmarks, such as `bench=false`.

**Why not:** every implementation of a waiting callable waits, so each subject of the declaration would repeat the param. A subject without the param waits in its benchmark.

### A directive of dokimi

dokimi reads a directive of its own on the method, such as `//dokimi:bench off`.

**Why not:** a declaration states what its callables do in the vocabulary of the shape catalog, and the suite generator already reads `concurrent` from it. A second vocabulary for one property of a callable would split what the consumers of eidos read.

### A fresh subject for each call of a benchmark

`Bench` builds a new subject before each call.

**Why not:** every benchmark would then measure the factory with each call. A callable whose cost depends on earlier calls, such as the read of a cache, would be measured in its first state alone.

### A deadline on each call of a benchmark

`Bench` passes a context that ends after a short time.

**Why not:** the benchmark of a waiting callable would then measure the deadline and count each call as an error.

### The role `campaign` of the contract `leader-election`

The generator leaves out each callable with the role `campaign`.

**Why not:** the role covers one protocol. The acquire of a lock and the dequeue of a blocking queue have no such contract. The catalog stamps a mixin or a contract by directive alone, so the role does not save a directive.

### The suite generator alone

The generator leaves the callable out, and conform has no field `Waits`.

**Why not:** a hand-written benchmark would still wait without a message, and so would a ceiling that an application adds through `Ceilings`. `BenchParallel` already refuses a call without `Concurrent` in the same way.

## Drawbacks

- **No measurement of the callable.** A callable with `waits` has no benchmark and no budget. No generated test measures its cost per call or its allocations on a real subject.
- **A callable that waits in some states.** The dequeue of a queue waits only while the queue is empty, and it still loses its benchmark in every state.
- **A release of conform first.** A generated suite with the field `Waits` compiles only against a release of conform with the field.
- **One line for each callable.** Each waiting callable needs a directive line, and its call in the generated suite one field.

## Open questions

None.

## References

| What | Where |
|---|---|
| The request to leave out the benchmark of a waiting callable | <https://github.com/dokimasia/dokimi/issues/33> |
| The request of the mixin `waits` | <https://github.com/dokimasia/eidos/issues/94> |
| The benchmarks of conform | <https://github.com/dokimasia/dokimi-go/blob/main/conform/bench.go> |
| The alarm of the testing package | <https://github.com/golang/go/blob/go1.27.3/src/testing/testing.go#L2596-L2618> |
| The kill timer of the go command | <https://github.com/golang/go/blob/go1.27.3/src/cmd/go/internal/test/test.go#L842-L854> |
