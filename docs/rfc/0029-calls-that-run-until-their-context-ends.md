---
rfc: 0029
title: Calls that run until their context ends
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-01
updated: 2026-11-01
discussion: https://github.com/dokimasia/dokimi/issues/46
supersedes: none
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0029: Calls that run until their context ends

## Summary

The loop of a server or a worker runs until its context ends, as a `Run(ctx context.Context) error` does. A derived check that calls such a callable with the context of the check never returns, and the check `cancel` expects an error that the contract of the loop does not return. The shape catalog of eidos gets the documentary mixin `runs-until-cancelled` for such a callable. The conformance spec gets the gate `returns-before-cancel`, which refuses the claims that the derivation binds through the callable, its context or its shape. The suite generator leaves the callable out of every benchmark and every ceiling that it writes, and it marks the call of the callable with the field `RunsUntilCancelled` of conform. conform then reports a benchmark, a budget or a ceiling of such a call as a problem of the application.

## Motivation

The derivation binds `smoke` to every callable and `cancel` to every callable with a context. The test of a subject declares a benchmark with a level for each callable. A callable that runs until its context ends breaks all three, as dokimi at 6ef4d86 with conform 494250b shows for this subject:

```go
// Runner runs a loop.
//
//dokimi:stub
//dokimi:suite
//dokimi:out path=runnertest/
type Runner interface {
	// Run runs the loop until ctx ends, and returns nil when ctx ends.
	Run(ctx context.Context) error
}

// Run waits for the end of ctx, and returns nil.
func (Loop) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}
```

- `smoke` calls `Run` through `probe.Survives` with the context of the check, which ends only after the check. `go test -timeout 20s` ends with `panic: test timed out after 20s` in `TestLoopConformsToRunner/Run/smoke`.
- `cancel` calls `Run` through `probe.UnderCancelled` with a cancelled context and expects an error. The loop returns nil, as its contract states, and the check fails with `Run returns an error for a cancelled context`.
- `BenchmarkLoopConformsToRunner` repeats the call on one subject with the context of the benchmark. It ends with `fatal error: all goroutines are asleep - deadlock!`, and a subject with a timer waits until the go command kills the test binary.

The storage contract of a ledger declares such a `Run` for its provider, and its suite and the test of each subject hang in `Provider.Run/smoke`.

The shape `lifecycle`, which dokimi gives such a `Run`, also classifies the prompt calls that move a subject to its next stage. The Go corpus has 26 callables of the shape, such as `Close`, `Commit` and `Acquire`, and each of them returns before its context ends, so its checks and its benchmark are valid. A rule over the shape alone would take those checks from all 26. The mixin `waits` states that a call waits until another call changes the state of its subject. A single call of such a callable on a fresh subject returns, so the derivation keeps its checks. A callable that runs until its context ends does not return on a fresh subject either.

Every implementation of such a callable runs until its context ends, so the property belongs to the declaration and not to one subject. The shape catalog states such properties as mixins, and the suite generator already reads the mixins `waits`, `concurrent` and `concurrent-readers`.

## Detailed design

### Invariants

| Invariant | Enforced by |
|---|---|
| No claim that the derivation binds without a mixin or a contract of the declaration applies to a callable with the mixin `runs-until-cancelled` | The gate `returns-before-cancel` of each such claim |
| No benchmark that dokimi writes measures such a callable | The suite generator |
| No ceiling that dokimi writes applies to such a callable | The suite generator, and `BadBudget` at a budget of the callable |
| A benchmark of a call with `RunsUntilCancelled` measures nothing and reports a problem | `Bench` and `BenchParallel` of conform |
| A budget of such a call, and a ceiling that `Ceilings` adds for it, are problems of the application | The problems of an application in conform |
| The claims of the mixins and the contracts that the declaration states on the callable keep their gates | The conformance spec |

### Components

| Component | Change |
|---|---|
| The shape catalog of eidos | The documentary mixin `runs-until-cancelled` |
| The conformance spec | The gate `returns-before-cancel`, in the gates of 12 claims |
| The derivation | The condition of the gate |
| `conform.Call[S]` | The field `RunsUntilCancelled` |
| `conform.Bench`, `conform.BenchParallel` and the problems of an application | The problems of a benchmark, a budget and a ceiling of such a call |
| The suite generator | `RunsUntilCancelled` of each call, the levels of the benchmarks, the ceilings of the test of the suite, and a case of `BadBudget` |
| The template of the suites | The field `RunsUntilCancelled` of each call with the mixin |

### The mixin `runs-until-cancelled`

The mixin is documentary, as `waits` is. It documents the callable and licenses no check. A declaration states it on a method or a function:

```go
// Run runs the loop of the provider until ctx ends, and returns nil when ctx
// ends.
//
//dokimi:mixin runs-until-cancelled
Run(ctx context.Context) error
```

eidos-plugin-shape states this spec:

```yaml
name: runs-until-cancelled
form: mixin
documentary: true
claim: >
  The callable runs until its context ends, such as the loop of a server or
  a worker. A generator of checks does not call it with a context that does
  not end, and does not expect an error when its context ends.
observation: >
  Nothing. The mixin documents the callable and licenses no check.
falsifiability: >
  Nothing turns red, because the mixin licenses no check.
counterexamples:
  refused: >
    a check that calls the callable with the context of the check, which
    ends only after the check
```

### The gate `returns-before-cancel`

The conformance spec gets one gate:

```yaml
returns-before-cancel:
  passes: The callable has no mixin runs-until-cancelled, so a call returns before its context ends.
  why: "{callable} runs until its context ends"
  remedy: "check {callable} with a check of the consumer that ends its context"
  since: 1.0.0
```

The derivation implements its condition once for every language:

```go
spec.ReturnsBeforeCancel: func(f *Facts) bool { return !slices.Contains(f.Class.Mixins, shape.RunsUntilCancelled) },
```

The gate joins the gates of the 12 claims that the derivation binds through the callable, its context or its shape:

| Claim | Licence | Why the claim cannot apply |
|---|---|---|
| `smoke` | Every callable | Its probe calls the callable with the context of the check |
| `cancel` | Every callable with a context | The end of the context completes the call, so a nil error is not an ignored cancellation |
| `failure-is-reported` | Every callable | Its probe calls the callable with the context of the check, and the callable can keep running after its dependency fails |
| `zero-on-error`, `zeroes-on-error`, `zero-on-miss`, `zeroes-on-miss`, `nil-on-miss`, `miss-is-zero`, `deterministic`, `predicate` and `repeat-survives` | A shape | Each probe calls the callable with the context of the check and reads its result |

A claim of every callable refuses without a record, as each gate of such a claim does. `cancel` records its refusal with the reason and the remedy of the gate:

```go
{
	Claim:    spec.Cancel,
	Callable: "Runner.Run",
	Gate:     spec.ReturnsBeforeCancel,
	Why:      "Run runs until its context ends",
	Remedy:   "check Run with a check of the consumer that ends its context",
},
```

The claims of the mixins and the contracts keep their gates. A declaration states such a claim on a callable that runs until its context ends where a call of the claim returns before its context ends, as a call after `Close` returns `ErrClosed` for the mixin `lifecycle-after-close`.

### The field `RunsUntilCancelled`

```go
type Call[S any] struct {
	// The fields Calls, Over, Callable, Budget, Concurrent and Waits as before.

	// RunsUntilCancelled reports that a call runs until its context ends, as
	// the mixin runs-until-cancelled states. A benchmark and the check of a
	// ceiling call the callable with a context that does not end, so a
	// benchmark of the call and a ceiling of the call are problems of the
	// application.
	RunsUntilCancelled bool
}
```

conform reports each of these problems before it runs a check, and it stops the test or the benchmark of the application:

| Entry point | Problem |
|---|---|
| `Bench` and `BenchParallel` | `the call of Runner.Run runs until its context ends, and a benchmark calls it with a context that does not end` |
| `Run`, `RunBudgets`, `Fuzz`, `Bench` and `BenchParallel` | `the budget of Runner.Run repeats a call that runs until its context ends`, for a call with `RunsUntilCancelled` and a budget |
| The same entry points | `the ceiling of Runner.Run repeats a call that runs until its context ends`, for a ceiling that `Ceilings` adds for such a call |

`Run` and `Fuzz` apply the checks of such a call as before, so a check of the consumer that ends the context of its call applies. `Fuzz` skips a callable without a property check.

### The suite generator

The generator sets `RunsUntilCancelled` of a call where the callable has the mixin, as it sets `Waits` where the callable has the mixin `waits`:

```go
{
	Calls:              runnerRunCalls,
	Over:               newRunnerOver,
	Callable:           RunnerChecks.Run.Callable,
	RunsUntilCancelled: true,
},
```

The generated tests change for such a callable:

| Generated test | Change for a callable with `runs-until-cancelled` |
|---|---|
| The test of a subject, such as `TestLoopConformsToRunner` | No check of the 12 claims |
| The benchmark of a subject, such as `BenchmarkLoopConformsToRunner`, and its parallel form | No level of the callable. A test without a level declares no benchmark |
| The fuzz target of the callable for a subject | Skips, because the callable has no derived property check |
| The allocation test of a subject | No ceiling of the callable, because `BadBudget` refuses its budget |
| The benchmark of the test of the suite and its parallel form | No level of the callable |
| The allocation test of the test of the suite | No ceiling of 0 for the callable |
| The benchmark, the parallel benchmark and the allocation test of the suite of a function with the mixin | Not declared |

`BadBudget` reports each budget of a callable with the mixin at its directive `budget`, also the budget of a subject that refers to the callable.

### The contract in each language

| Part | Go | TypeScript |
|---|---|---|
| The mixin | `//dokimi:mixin runs-until-cancelled` | The same directive |
| The gate | `spec.ReturnsBeforeCancel` | The same gate of the spec |
| The field of a call | `RunsUntilCancelled bool` of `Call[S]` | `runsUntilCancelled` of `Call<S>` |

The TypeScript runtime and the template tree of its suites do not exist yet. Their parts take these names when they do.

### Failure handling

| Failure | Result | Recovery |
|---|---|---|
| A callable that runs until its context ends has no mixin | Its check `smoke` waits until the timeout of the test, its check `cancel` fails where it returns nil, and its benchmark reports a deadlock or waits | State the mixin |
| The mixin on a callable that returns before its context ends | The callable has no check of the 12 claims, no benchmark and no budget | Remove the mixin |
| A budget of a callable with the mixin | `BadBudget` at the directive, and no conformance test of a subject with the budget | Remove the budget |
| A hand-written benchmark of a call with `RunsUntilCancelled` | `Bench` reports the problem and stops the benchmark | Remove the benchmark |
| An application adds a ceiling for such a call | The application reports the problem and runs no check | Remove the ceiling |
| A check of the consumer that calls the callable with the context of the check | The check waits until the timeout of the test | End the context of the call in the check |
| An eidos without the mixin | eidos refuses the directive with `EID-0068` | Upgrade eidos |

### Verification

- **The derivation.** The spec of the gates covers `returns-before-cancel` with and without the mixin. The spec of the derivation covers each of the 12 claims on a callable with the mixin, and a claim of a mixin on such a callable.
- **conform.** The specs of `Bench` and `BenchParallel` cover the problem of a call with `RunsUntilCancelled`. The spec of the problems of an application covers a budget of such a call and a ceiling of `Ceilings` for it. The spec of `Run` covers a check of the consumer on such a call.
- **The suite generator.** The specs of the plugin cover `RunsUntilCancelled` of a call, the levels of each benchmark, the ceilings of the test of the suite and the new case of `BadBudget`. The golden files of a fixture with such a callable pin the field, the refusal of `cancel`, and the tests without its levels.
- **The corpus.** The fixture `iface/mixin/runsuntilcancelled` declares a runner whose `Run` waits for the end of its context, with a check of the consumer that cancels the context of its call. `go test` of the fixture ends, and so does `go test -run '^$' -bench . -benchtime 1x`.

### Migration

The change adds a mixin, a gate, a field, the cases of existing problems and a case of an existing code. The gates of 12 claims change, and a callable without the mixin passes the new gate, so a declaration without the mixin generates the same files as before. A generated suite with the field `RunsUntilCancelled` requires a release of conform with the field, and a declaration with the mixin requires a release of eidos-plugin-shape with the mixin.

## Alternatives considered

### The shape `lifecycle`

The generator refuses the checks and leaves out the benchmarks of every callable of the shape `lifecycle`, as the comment of the request proposes.

**Why not:** the 26 callables of the shape in the Go corpus return before their context ends. Each would lose its checks `smoke` and `cancel` and its benchmark, and each of them is valid.

### The mixin `waits`

A declaration states `waits` on the callable, and the claim of `waits` widens to a call that only another call or the end of its context returns, as the request proposes.

**Why not:** `waits` refuses `smoke` through the gate `returns-unaided` and leaves the callable out of the benchmarks and the ceilings. `cancel` and `failure-is-reported` keep their gates, so `cancel` still fails for a loop that returns nil. The gate `returns-unaided` on `cancel` would refuse `cancel` for every callable that waits for another call too. The check of such a callable is valid, because a campaign returns an error once its context ends.

### A check `smoke` with an ended context

`smoke` calls a callable with the mixin with a context that has ended already, as the request proposes.

**Why not:** such a call returns at its first read of the context. The check covers the code before that read alone, and the probe `survives` would need a second form for it.

### A claim that the call returns once its context ends

A claim of the mixin starts the call, ends its context and checks that the call returns with the result that its contract states.

**Why not:** the claim needs a probe that runs the call in a `testing/synctest` bubble and finds a call that has not returned after the end of its context, a law of its own, and a form of planted defect that passes the call a context that never ends. The form `ignores-cancel` returns at once for an ended context, so it does not fit. conform accepts the deadlock of the bubble after the planted defect fails the statement, so the parts exist, and they are a design of their own. The claim is not proposed here, and a check of the consumer covers it, as the corpus fixture of the mixin shows.

### A directive of dokimi

dokimi reads a directive of its own on the method, such as `//dokimi:runs`.

**Why not:** a declaration states what its callables do in the vocabulary of the shape catalog, and the suite generator already reads `waits` from it. A second vocabulary for one property of a callable would split what the consumers of eidos read.

## Drawbacks

- **No generic check of the callable.** A callable with the mixin loses the checks of the 12 claims, its benchmark and its budget. A consumer writes a check of the directive `check` that ends the context of its call.
- **One more gate on 12 claims.** The spec grows from 38 gates to 39, and each of the 12 claims lists one gate more.
- **Releases of conform and eidos first.** A generated suite with the field compiles only against a release of conform with the field, and the directive needs a release of eidos-plugin-shape with the mixin.
- **One line for each callable.** Each such callable needs a directive line, and its call in the generated suite one field.

## Open questions

None.

## Unresolved and future work

- A claim that a callable with the mixin returns once its context ends, with the result that its contract states, is not proposed here.
- The TypeScript runtime and the TypeScript template tree of the suite generator do not exist. The table of the contract fixes their names.

## References

| What | Where |
|---|---|
| The request to leave out the checks and the benchmark of such a callable | <https://github.com/dokimasia/dokimi/issues/46> |
| The request of the mixin `runs-until-cancelled` | <https://github.com/dokimasia/eidos/issues/109> |
| The benchmarks of conform | <https://github.com/dokimasia/dokimi-go/blob/main/conform/bench.go> |
| The probes of a single caller | <https://github.com/dokimasia/dokimi-go/blob/main/conform/probe/single.go> |
