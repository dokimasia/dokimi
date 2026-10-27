---
rfc: 0011
title: The environment of the runs of a check
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Draft
created: 2026-10-27
updated: 2026-10-27
discussion: https://github.com/dokimasia/dokimi-go/issues/17
supersedes: RFC-0006 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0011: The environment of the runs of a check

## Summary

A check states the environment of its runs through the field `Around` of `conform.Check`, such as a `testing/synctest` bubble through `synctest.Test`. conform runs every run of the check inside it, also the falsification on the planted defect and the isolation of another planted defect. A property check runs each of its cases inside it. The directive `check` sets the field through its param `around`. conform recovers a failure of the environment, such as the panic of a bubble that still has blocked goroutines when its root goroutine exits. An application reports that failure as a failure of the check. Falsification accepts it once the check has failed at its statement. Fuzz targets leave a check with `Around` out.

## Motivation

A check whose body waits with `synctest.Wait` needs a bubble. A conformance kit of the volumes of a storage engine checks three laws. A failed sync ends each write that waits for it. Close completes each admitted write. Close ends a write that waits for room in the queue. Each check starts its writes on goroutines and waits with `synctest.Wait` until they block in the volume. It then fails a sync, ends a context or closes the volume. A bubble also fails at once a volume that leaves a write blocked, where a check outside a bubble waits until the timeout of the test.

Today a check states such an environment through `prop.Around` among its `Options`. An application applies the options, and falsification leaves them out (`conform/check.go:50` at conform v0.7.0). Falsification runs the check on its planted defect on the seat of `assert.Rejects`. The isolation runs every other check of the suite on an `assert.Recorder`. Both seats lack a chain to a `*testing.T`, so `prop.Around` cannot open a bubble there and ends the run with a fault (dokimasia/assert-go#75). The check of the kit fails its falsification:

```text
the planted defect of Valuer.Value/returns-one fails the check before an assertion fails. Falsification needs a failure of the assertion whose message is the statement "Value returns 1". The check failed with this message:
    goroutine is not in a bubble
```

The planted defect of such a check can also leave goroutines blocked by design. The defect of the third law never ends the write that waits for room, so the write is still blocked when the check ends. `synctest.Test` panics when blocked goroutines remain in a bubble after its root goroutine exits (`testing/synctest/synctest.go:69-75` and `runtime/synctest.go:243-247` at Go 1.27.3). prop raises a panic of the wrapper of `Around` again on the goroutine of the run (`prop/option.go:254` at assert v0.9.0). The panic would end the falsification before it reads the failure at the statement.

## Detailed design

### Invariants

| Invariant | Enforced by |
|---|---|
| Each run of a check with `Around` runs inside it: the application, the falsification and the isolation | conform |
| Each case of a property check with `Around` runs inside it | `prop.Around`, which conform passes |
| A failure of the environment in an application fails the test of the check | conform |
| A failure of the environment on a planted defect does not change the verdict of a falsification in which the check failed at its statement | conform |
| A fuzz target leaves out each check with `Around` | `conform.Fuzz` |
| The function of the param `around` has the signature of `synctest.Test` | `BadCheck` of the suite generator |

### Components

| Component | Change |
|---|---|
| `conform.Check[S]` | The field `Around` |
| The application of a check in conform | Runs the check inside `Around` |
| The falsification and the isolation in conform | Run the check inside `Around` on the test of the falsification |
| `conform.Fuzz` | Leaves out a check with `Around` |
| assert-go | A property with `Around` under `Rejects` and on a `Recorder`, which dokimasia/assert-go#75 requests |
| The directive `check` | The param `around` |
| The code `BadCheck` | The signature of the function of `around` |
| The template of the suites | The field `Around` of each check with the param |

### The field `Around`

```go
type Check[S any] struct {
	// The fields Apply, Property, Defect, Doors, Options, ID, Statement and
	// Counted as before.

	// Around runs each run of the check inside an environment, such as a
	// testing/synctest bubble through synctest.Test. It calls run once with
	// the test of the environment. An application runs the check inside
	// Around on the test of the check. Falsification runs the check on its
	// planted defect inside Around on the test of the falsification, and so
	// does the isolation of another planted defect. A property check runs
	// each of its cases inside Around. Fuzz targets leave a check with
	// Around out.
	Around func(t *testing.T, run func(*testing.T))
}
```

The signature is the signature of `synctest.Test` and of the wrapper of `prop.Around`, so a check states `Around: synctest.Test` directly.

### The runs inside `Around`

| Run | A check with `Apply` | A check with `Property` |
|---|---|---|
| An application: `Run`, `RunBudgets` and `Check.Run` | conform calls `Around` with the test of the check. It builds the handle on the test that `Around` hands `run` | conform passes `prop.Around` with `Around` after the options of the run |
| Falsification: `Falsify`, `FalsifyAllocs`, `Check.Falsify` and the falsification of a check that `With` adds | conform calls `Around` with the test of the falsification. It calls `Rejects` on the test that `Around` hands `run` | conform passes `prop.Around` with `Around`, and the run under `Rejects` finds the test of the falsification |
| The isolation of another planted defect | conform calls `Around` of the other check with the test of the falsification. It runs the check on a `Recorder` with the context of the test that `Around` hands `run` | conform passes `prop.Around` with `Around` of the other check, on a `Recorder` that takes the test of the falsification |
| A fuzz target | Left out | Left out |

The property runs under falsification and under the isolation need dokimasia/assert-go#75. That request makes a property with `Around` run inside the wrapper under the seat of `Rejects`, and on a `Recorder` that takes a test.

The test that `Around` receives is the first `*testing.T` of the chain of the seat, as for `prop.Around`. `Check.Run` and `Check.Falsify` on a seat without one end with a fault that wraps `ErrConfig`:

```text
conform: the check Valuer.Value/returns-one runs inside Around, and the seat is no *testing.T and wraps none
```

`Options` runs before the run and outside `Around`, as before. A check that states `prop.Around` among its options and the field `Around` runs inside the field, because conform passes the field after the options.

### Failures of the environment

conform passes `Around` inside a wrapper of its own. The wrapper calls `Around` and recovers a panic of the environment. It keeps the panic for the run. `synctest.Test` panics this way when blocked goroutines remain in its bubble. The wrapper then returns, so prop states the verdict of the case from the records of its body. Each run handles a kept failure of the environment in one way:

| Run | A failure of the environment |
|---|---|
| An application | Fails the test of the check: `the environment of Valuer.Value/returns-one fails: deadlock: main bubble goroutine has exited but blocked goroutines remain` |
| Falsification, after the check failed at its statement | Changes nothing. The planted defect fails the check |
| Falsification, before the check failed at its statement | Fails the test of the falsification as a failure outside an assertion, with the message of the panic |
| The isolation of another planted defect | Fails the other check on that defect: `Valuer.Value/returns-one fails its environment with "deadlock: main bubble goroutine has exited but blocked goroutines remain"` |

A planted defect that leaves goroutines blocked also fails the bubble of another check that drives it into that state. So the isolation counts the failure of the environment of another check as a failure of that check, as it counts a panic.

### The directive `check`

The param `around` contains the name of a function of the package of the check, as `defect` and `options` do:

```go
// FailedSyncEndsTheWrites checks that a failed sync ends each write that
// waits for it.
//
//dokimi:check Volume.Write defect=ignoresFailedSyncs around=bubble
func FailedSyncEndsTheWrites(h *VolumeHandle, c *prop.Case) { … }

// bubble runs each run of a check in a testing/synctest bubble of its own.
func bubble(t *testing.T, run func(*testing.T)) { synctest.Test(t, run) }
```

The annotate step stamps the identity of the function as a fact of the check, under the key `suite.around`. It reports `BadCheck` at the directive for a function that does not take a `*testing.T` and a `func(*testing.T)`, or that returns a result. The generator writes the field into the check of the suite:

```go
{
	ID:        VolumeChecks.Write.FailedSyncEndsTheWrites(),
	Statement: "a failed sync ends each write that waits for it",
	Property:  FailedSyncEndsTheWrites,
	Defect:    ignoresFailedSyncs,
	Around:    bubble,
},
```

### The contract in each language

| Part | Go | TypeScript |
|---|---|---|
| The environment of a check | The field `Around func(t *testing.T, run func(*testing.T))` of `Check[S]` | The field `around` of `Check<S>`, a function that runs each run of the check inside an environment |
| The param of the directive | `around` of `check` | The same param |

The TypeScript runtime and its template tree do not exist yet. Their parts take these names when they do.

### Failure handling

| Failure | Result | Recovery |
|---|---|---|
| A check with `Around` on a seat without a test | A fault that wraps `ErrConfig`, and the check does not run | Run the check on a `*testing.T`, or on a seat that wraps one |
| The environment fails in an application | The test of the check fails with the message of the environment | Make the subject end its goroutines, or correct the check |
| The environment fails on a planted defect before the check fails at its statement | The test of the falsification fails as a failure outside an assertion | Change the defect or the check, so that the statement fails first |
| A function of `around` with another signature | `BadCheck` at the directive. The check is not bound | Correct the function |
| `Around` returns without a call of `run`, or calls it twice | The fault of `prop.Around` for a property check. conform ends an `Apply` check with the same fault | Call `run` once |

### Verification

- **conform.** The spec of a check covers `Around` in an application, in falsification and in the isolation, for `Apply` and for `Property`. It covers a bubble with blocked goroutines in each run, a seat without a test, and a wrapper that does not call `run`. The spec of `Fuzz` covers a target that leaves a check with `Around` out.
- **The suite generator.** The specs of the plugin cover the param `around` and its `BadCheck`. The golden files pin the field `Around` of a generated check.
- **The corpus.** A fixture of a queue whose Close ends each write that waits for room declares its check in a bubble. Its planted defect leaves the waiting write blocked, and the falsification of the check passes.

### Migration

The change adds a field, a param and a case of an existing code, so no existing call breaks. A check without `Around` runs as before. conform implements the property runs under falsification and under the isolation once assert-go releases dokimasia/assert-go#75, and the suite generator follows the release of conform.

## Alternatives considered

### `prop.Around` among the options, also under falsification

Falsification passes the options of a check, and conform keeps `prop.Around` among them.

**Why not:** the options also contain the shares of a run. A planted defect changes those shares, so falsification leaves the options out. An `Apply` check has no property run to pass options to. The environment applies to every run of a check, and the options apply to one kind of run.

### An option of the application

An application passes the environment of each check, such as `conform.Around(id, synctest.Test)`.

**Why not:** the body of the check calls `synctest.Wait`, so the check needs the bubble in every application. Each subject of the suite would repeat the option.

### A bubble inside the body of the check

The check calls `synctest.Test` itself.

**Why not:** `synctest.Test` needs a `*testing.T`, and the handle of a check under falsification has none. The subject and the context of the handle must also belong to the bubble, and the handle builds both before the body runs.

### One bubble around the whole property run under falsification

conform calls `Around` once around the property run on the seat of `Rejects`, without a change of assert-go.

**Why not:** the cases would share one clock and the goroutines of earlier cases. The cases of an application each start a bubble of their own, so a check could pass in an application and fail under falsification for another reason.

### A seat of conform that leads to the test

conform wraps the seat of `Rejects` and a `Recorder` in a seat whose method `Unwrap` returns the test.

**Why not:** assert looks for each method that a seat lacks along the chain. A call of a method such as `Skipf` or `Logf` from inside a falsification would then go to the test. A second recording seat in conform would also duplicate the seats of assert.

### Fuzz targets inside the environment

A fuzz target runs a check with `Around` inside it.

**Why not:** one target draws among the property checks of a callable, and prop applies one wrapper to every case of a target. A check with `Around` would need a target of its own.

## Drawbacks

- **A wait on assert-go.** The property runs under falsification and under the isolation need dokimasia/assert-go#75 before conform can implement them.
- **No fuzz target.** A check with `Around` does not run under the fuzzer.
- **Blocked goroutines for the rest of the process.** conform recovers the panic of a bubble with blocked goroutines. Those goroutines remain in the test process until the process ends.
- **A bubble per case.** Each case of a property check with `Around` starts a bubble of its own. We have not measured what a bubble costs.

## Open questions

None.

## References

| What | Where |
|---|---|
| The request of the environment of a check | <https://github.com/dokimasia/dokimi-go/issues/17> |
| The request of `Around` under `Rejects` and on a `Recorder` | <https://github.com/dokimasia/assert-go/issues/75> |
| `prop.Around` of the assert library | <https://github.com/dokimasia/assert-go/blob/main/prop/option.go> |
| The bubbles of `testing/synctest` | <https://pkg.go.dev/testing/synctest> |
