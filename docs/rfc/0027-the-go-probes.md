---
rfc: 0027
title: The Go probes
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-29
updated: 2026-10-29
discussion: none
supersedes: none
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0027: The Go probes

## Summary

The package `go.dokimi.dev/conform/probe` implements each probe of the conformance spec as one Go function. A probe takes the handle of its check first and then one argument for each input of its entry in the spec. Each input kind has one Go form. A `call` has one of four call shapes, and the generator adapts each callable to the shape that the probe of its claim takes. The probe sets the reason of its claim through `Handle.Expect` before each decisive assertion, and passes the statement of the check as the contract of each assertion.

## Motivation

A generated derived check calls one probe, and the check is only as good as the probe. The conformance spec states each probe in language-neutral terms: a summary, inputs of eleven kinds, the reasons and the vacuous condition. The kinds state what the generated check passes, such as "a function that calls the subject with bound arguments". They do not state the Go type of that function.

A callable returns a value beside an error or beside a presence flag, two values, a stream, an error alone, or nothing. The 71 probes of the spec take 102 inputs of the kind `call`. A Go signature per probe and per result form would multiply each of those inputs by the forms that its callable can have. testkit took that path in part, with 68 probes for fewer laws, such as `RepeatAgrees`, `RepeatAgreesValue`, `RepeatAgreesSlice` and `RepeatAgreesFlag` for one law.

This RFC fixes the Go form of each input kind once, so the generator and the probes agree on every signature.

## Detailed design

### Conventions of a probe

Every probe follows these rules:

- **The handle first.** The first parameter is the handle `h *conform.Handle[S]` of the check, where `S` is the type of the subjects. The probe reads the context of each call from `h.Context()`, the statement from `h.Statement()` and the clock from `h.Clock()`.
- **The reason before the assertion.** The probe calls `h.Expect` with the reason of a failure before each decisive assertion. A guard assertion, such as the success of an establishing write, does not set a reason. Its failure has no reason, and falsification reports it as a failure without a reason.
- **The statement as the contract.** Each assertion passes `h.Statement()` as its message. A probe over `values` passes the statement followed by the names of the values at fault.
- **Equality by structure.** A probe compares two values through `assert.Equal`, which compares any two values by structure, so the type parameters do not need the constraint `comparable`. A probe that orders values requires `cmp.Ordered`.
- **No vacuity yet.** A probe in this RFC does not report that it did not exercise its claim. A check whose precondition is not reached passes. The handle method that marks a check as engaged belongs to the design of vacuous checks.

```go
// Package probe implements the probes of the conformance spec of dokimi.
// Each function is one probe. It takes the handle of its check first and
// one argument for each input of its entry in the spec, in the order of
// the entry.
package probe
```

### The input kinds in Go

| Kind | Go form | Notes |
|---|---|---|
| `call` | One of the call shapes in the next section | The generator binds the arguments of the callable, as RFC-0014 states. A write or a registration returns the value that it wrote or the delay that it set, which the probe compares |
| `value` | A value of a type parameter, such as `zero V` | A sample, a zero value or a declared edge |
| `error` | `error` | An error value, compared with `errors.Is` |
| `duration` | `time.Duration` | |
| `count` | `int` | |
| `projection` | `func(v V) O` with `O cmp.Ordered` | Returns one part of a value, such as its version or its cursor |
| `scope` | `func(ctx context.Context, body probe.Act) error` | Runs body inside a scope of the subject, such as a transaction |
| `hook` | `func(ctx context.Context, hook func()) error` | Installs hook through the subject |
| `predicate` | `probe.Call[bool]` | Reports whether the subject is in a state, such as leadership |
| `clock` | No parameter | The probe advances the clock of the handle through `h.Advance` |
| `values` | `[]conform.Named[T]` | The values that the check covers as one set, each with its name |

```go
// Named is a value with the name by which the source declares it, such as
// a variant of an enumeration or an error value of a package. A probe over
// values names each value at fault in its contract.
type Named[T any] struct {
	Name  string
	Value T
}
```

`conform.Named` is declared in the package `conform`, because the generated suites of the enumerations and the errors pass it without importing the package `probe`.

### The call shapes

A `call` input has one of four shapes. The generator writes a closure that calls the callable with its bound arguments and returns the results in the shape that the probe takes:

```go
// Call calls a callable and returns its value, and the error of a call
// that fails. A callable without an error result returns a nil error. A
// callable with two values returns them as a Pair. A callable whose value
// is a stream returns the elements that the stream yields, collected into
// a slice, and the error that ends the stream.
type Call[V any] func(ctx context.Context) (V, error)

// Act calls a callable for its effect, and returns the error of a call
// that fails. The probe does not read the values of the callable.
type Act func(ctx context.Context) error

// Lookup calls a callable with a presence flag, and returns its value and
// the flag. A callable with two values beside the flag returns them as a
// Pair.
type Lookup[V any] func(ctx context.Context) (V, bool)

// Pair contains the two values of a callable that returns two values.
type Pair[A, B any] struct {
	First  A
	Second B
}
```

| Callable returns | Shape that a probe of a value takes | Shape that a probe of an effect takes |
|---|---|---|
| `(V, error)` | `Call[V]`: the results as they are | `Act`: the error |
| `V` | `Call[V]`: the value and nil | `Act`: nil |
| `(V, bool)` | `Lookup[V]`; `Call[Pair[V, bool]]` for a probe that compares whole results | `Act`: nil |
| `(A, B, error)` | `Call[Pair[A, B]]` | `Act`: the error |
| `(A, B, bool)` | `Lookup[Pair[A, B]]` | `Act`: nil |
| `error` | — | `Act`: the error |
| nothing | — | `Act`: nil |
| a stream of `E` with an error | `Call[[]E]`: the collected elements | `Act`: the error of the stream |

The gates of each claim select the callables whose results fit the shape of its probe. `zero-on-miss` has the gate `presence-flag`, so its callable has a flag, and its probe takes a `Lookup`.

### The probes

`S` is the type of the subjects, `V` and `E` are the types of a value and of an element, and `O` is an ordered type. Each signature leaves out the first parameter `h *conform.Handle[S]`. Each probe sets the one reason of its entry in the spec.

#### A single call

| Probe | Signature | Procedure |
|---|---|---|
| `Survives` | `(call Act)` | Calls call, and fails when the call panics |
| `UnderCancelled` | `(call Act)` | Calls call with a cancelled context, and fails when it returns nil |
| `ZeroBesideError` | `[V](call Call[V], zero V)` | Calls call. Where it succeeds, calls it again with a cancelled context. Fails when a call that returns an error returns a value other than zero |
| `ZeroPairBesideError` | `[A, B](call Call[Pair[A, B]], zero Pair[A, B])` | As `ZeroBesideError`, for both values |
| `ZeroBesideFalse` | `[V](call Lookup[V], zero V)` | Calls call, and fails when it returns a false flag beside a value other than zero |
| `ZeroPairBesideFalse` | `[A, B](call Lookup[Pair[A, B]], zero Pair[A, B])` | As `ZeroBesideFalse`, for both values |
| `NilBesideAMiss` | `[V](call Call[V])` | Calls call on the fresh subject, and fails when its value is not nil |
| `ZeroBesideAMiss` | `[V](call Call[V], zero V)` | Calls call on the fresh subject, and fails when its value is not zero |
| `RepeatAgrees` | `[V](call Call[V])` | Calls call twice, and fails when the two values differ. A first call that fails ends the probe |
| `SurvivesRepeat` | `(call Act)` | Calls call twice, and fails when the second call panics |
| `ReportsSentinel` | `[V](call Call[V], sentinel error)` | Calls call on the fresh subject, and fails unless its error matches sentinel |
| `RepeatIsQuiet` | `(call Act)` | Calls call twice, and fails when the second call fails after a first that succeeded |
| `NeverDecreases` | `[O cmp.Ordered](call Call[O])` | Calls call twice, and fails when the second value is below the first |
| `RepeatAscends` | `[O cmp.Ordered](call Call[O])` | Calls call twice, and fails unless the second value is above the first |
| `EscapesTheUnsafe` | `(call Call[string], unsafe string)` | Calls call, and fails when its output contains unsafe |
| `AcceptsTheEdge` | `[V](call func(ctx context.Context, v V) error, edge V)` | Calls call with edge, and fails when the call fails |
| `RemainsPoisoned` | `(call, induce Act)` | Calls induce, then call twice, and fails when a call succeeds |
| `SucceedsWithin` | `(call Act, attempts int)` | Calls call up to attempts times, and fails when no call succeeds |
| `RefusesTheSecondWrite` | `(call Act, sentinel error)` | Calls call twice, and fails unless the second error matches sentinel |
| `RefusesPastTheBurst` | `(call Act, burst int, sentinel error)` | Calls call burst times, then once more, and fails unless the last error matches sentinel |
| `RefusesEarly` | `(call Act, sentinel error)` | Calls call on the fresh subject, and fails unless its error matches sentinel |

#### A call and a partner

| Probe | Signature | Procedure |
|---|---|---|
| `WriteThenRead` | `[V](read, write Call[V])` | Calls write, which returns the value that it wrote, then read, and fails when read does not return that value |
| `PinnedAcrossWrite` | `[V](read Call[V], write Act)` | Calls read, write and read on one snapshot, and fails when the two reads differ |
| `SecondWriteReplaces` | `[V](call, read Call[V])` | Calls call twice, each time with another sample and returning the value that it wrote, then read, and fails when read does not return the second value |
| `KeepsPartitionsApart` | `[V](write, read Call[V])` | Calls write twice, in the partition of the read and in another, then read, and fails when read does not return the first value |
| `RoundTrips` | `[V, W](call Call[Pair[V, W]], inverse func(ctx context.Context, w W) (V, error))` | Calls call, which returns its input and its result, and inverse with the result, and fails when inverse does not return the input |
| `WorksAfter` | `(call, ready Act)` | Calls ready, then call, and fails when call fails |
| `ReportsSealed` | `(call, close Act, sentinel error)` | Calls close, then call, and fails unless the error of call matches sentinel |
| `SealsAfterClose` | `[H](call Call[H], close, after func(ctx context.Context, h H) error, sentinel error)` | Opens a handle with call, closes it with close, calls after on it, and fails unless the error matches sentinel |
| `MovesObservation` | `[V](call Act, observe Call[V])` | Observes, calls call, observes again, and fails when the two observations are equal |
| `EachCallCompounds` | `[V](call Act, observe Call[V])` | Observes around each of two calls, and fails when one call leaves the observation unchanged |
| `SettlesThenObserves` | `[V](call, settle Act, observe Call[V])` | Observes, calls call and settle, observes again, and fails when the observations are equal |
| `DetectsTamper` | `(call, verify, tamper Act)` | Calls call, verify, tamper and verify, and fails when the last verify succeeds |
| `ConvergesAfter` | `(call, begin Act, attempts int)` | Calls begin, then call up to attempts times, and fails when no call succeeds |
| `RunsTheHook` | `(call Act, register func(ctx context.Context, hook func()) error)` | Registers a hook, calls call, and fails when the hook did not run |
| `RefusesTheInvalid` | `[V](call, validate func(ctx context.Context, v V) error, invalid V)` | Calls validate with invalid, then call with invalid, and fails when call succeeds. A validator that accepts invalid ends the probe |
| `RefusesWhatThePredicateRejects` | `(call Act, match Call[bool])` | Calls call, which writes a first value. Then calls match, which reports whether the predicate accepts the value of the second call, then call again, and fails when the second call succeeds. The two values differ in the field that the param field of the contract names. A first call that fails ends the probe, and so does a predicate that accepts the value |
| `ResignsLeadership` | `(call, resign Act, leading Call[bool])` | Calls call and resign, and fails when leading reports leadership |
| `RollsBackOnFailure` | `(scope func(ctx context.Context, body Act) error, write Act, read Act, sentinel error)` | Runs scope with a body that writes and then fails, then read, and fails unless the error of read matches sentinel |
| `AgreesWithTheBacking` | `[V](call, backing Call[V], put Act)` | Calls put, then call and backing, and fails when their values differ |
| `TakesBackWhatItLent` | `[V](call Call[V], put func(ctx context.Context, v V) error)` | Gets a value with call and puts it back, and fails when put fails |
| `ContinuesFromTheCursor` | `[V, C](call func(ctx context.Context, cursor C) (V, error), cursor func(v V) C)` | Reads the first page, reads with its cursor, and fails when the second page has the cursor of the first |
| `CompensatesTheStep` | `(call, compensate Act)` | Calls call, then compensate, and fails when compensate fails |
| `RefusedLeavesNothing` | `[V, R](call func(ctx context.Context, r R) error, reader Call[V], refused R, zero V)` | Calls call with refused, then reader, and fails when reader returns a value other than zero. A call that accepts refused ends the probe |
| `RemovesOnDelete` | `(call, put, read Act, sentinel error)` | Calls put, call and read, and fails unless the error of read matches sentinel |

#### Drains, histories and versions

| Probe | Signature | Procedure |
|---|---|---|
| `NoElementTwice` | `[E](drain Call[[]E], write Act)` | Writes twice, drains once, and fails when an element repeats |
| `DrainsInOneOrder` | `[E](drain Call[[]E], write Act)` | Writes twice, drains twice, and fails when the drains differ in order |
| `GainsWhatWentIn` | `[E](drain Call[[]E], write Act)` | Drains, writes twice, drains, and fails unless the drain gains exactly two elements |
| `KeepsWhatWentIn` | `[E](drain Call[[]E], write Act)` | As `GainsWhatWentIn`, and fails when the drain gains fewer than two |
| `ReadsNeverRegress` | `[V, O cmp.Ordered](read Call[V], write Act, version func(v V) O)` | Reads, writes, reads, and fails when the second version is below the first |
| `WritesAdvance` | `[V, O cmp.Ordered](read Call[V], write Call[V], version func(v V) O)` | Reads, then writes, and fails unless the version of the write is above the version of the read |
| `ReplaysInStampOrder` | `[O cmp.Ordered](call Act, replay Call[[]O])` | Appends twice, replays, which returns the stamp of each entry, and fails when a stamp is below the stamp before it |
| `ReplayGrows` | `[E](call Act, replay Call[[]E])` | Replays, appends, replays, and fails unless the replay grows by one entry |

#### The clock

Each probe of the clock runs inside a `testing/synctest` bubble, which the generated check sets through its field `Around`. Before it advances the clock, the probe calls `synctest.Wait`, so every goroutine of the subject that waits on the clock has started to wait.

| Probe | Signature | Procedure |
|---|---|---|
| `ReturnsWithin` | `(call Act, budget time.Duration)` | Calls call on a goroutine, advances the clock past budget, and fails when the call has not returned |
| `CountsWithinWindow` | `[N cmp.Ordered](call Act, count Call[N], window time.Duration)` | Calls call, reads count inside the window and after it, and fails unless the first count includes the call and the second does not |
| `RunsOnSchedule` | `[N cmp.Ordered](call Call[time.Duration], count Call[N])` | Registers the task with call, which returns the delay of its instant, reads count halfway to the instant and after it, and fails unless the task ran only after its instant |
| `MovesWithTheClock` | `[V](call Call[V], put Act)` | Calls put and call, advances the clock in doubling steps, and fails when the value never changes |
| `SteadyWhileTheClockStops` | `[V](call Call[V], put Act)` | Calls put, then call twice without an advance, and fails when the values differ |
| `ExpiresAfter` | `[V](read Call[V], write Act, lifetime time.Duration, sentinel error)` | Writes, advances the clock past lifetime, reads, and fails unless the error matches sentinel |

#### Enumerations and errors

These probes take no subject. Their handle has the subject type `struct{}`. Each input of a function or of an outside value has its entry in the spec, with the sources `list`, `valid`, `outside` and `name` of the claim. The gates `list-function`, `valid-function` and `outside-value` refuse a claim of an enumeration that lacks the function or the value.

| Probe | Signature | Procedure |
|---|---|---|
| `ListsTheMembers` | `[E](members []conform.Named[E], list func() []E)` | Fails unless list returns the members once each, in their order |
| `NamesDiffer` | `[E](members []conform.Named[E], name func(e E) string, outside E)` | Fails when two members, or a member and outside, have one name |
| `JudgesValidity` | `[E](members []conform.Named[E], valid func(e E) bool, outside E)` | Fails unless valid accepts each member and refuses outside |
| `ParsesEachText` | `[E](members []conform.Named[E], parse func(text string) (E, error), name func(e E) string)` | Fails unless the name of each member parses to the member, and a text of no member fails |
| `EncodesEachVariant` | `[E encoding.TextMarshaler](members []conform.Named[E], name func(e E) string)` | Fails unless each member encodes its name |
| `DecodesEachText` | `[E any, P interface{ *E; encoding.TextUnmarshaler }](members []conform.Named[E], name func(e E) string)` | Fails unless the name of each member decodes to the member, and a text of no member fails |
| `MessagesDiffer` | `(values []conform.Named[error], opening string)` | Fails unless each message opens with opening and no message equals or opens with another |
| `ErrorsDiffer` | `(values, neighbours []conform.Named[error])` | Fails when an error value matches another value or a neighbour |
| `TypeMessage` | `(types []conform.Named[error], opening string)` | Fails unless the message of a value of each type opens with opening |
| `ComparesWithIdentity` | `(types, values []conform.Named[error])` | Fails when the method Is of a type disagrees with `errors.Is` on a value |
| `ReturnsTheCause` | `(causes []conform.Named[func(cause error) error])` | Fails unless each function makes an error that unwraps to the cause it receives |
| `TypesDiffer` | `(types []conform.Named[error])` | Fails when a value of one type matches a value of another |

### An example

The check of RFC-0014 calls `WriteThenRead` with the read and a write that returns the value that it wrote:

```go
Apply: func(h *StoreHandle) {
	probe.WriteThenRead(h,
		func(ctx context.Context) (catalog.Item, error) { return h.Subject().Read(ctx, "a") },
		func(ctx context.Context) (catalog.Item, error) {
			item := catalog.Item{Name: "b"}
			return item, h.Subject().Put(ctx, "a", item)
		})
},
```

## Alternatives considered

### One function for each result form

testkit declares a probe for each result form of a law, such as `RepeatAgrees`, `RepeatAgreesValue`, `RepeatAgreesSlice` and `RepeatAgreesFlag`.

**Why not:** the spec has one probe for each law, and its reasons belong to the probe. A function per form multiplies the probes of the runtime past the probes of the spec, and each new form adds functions to every probe that it fits. Four call shapes and a pair cover every form that the gates admit.

### A zero value that the probe makes

testkit makes the zero value inside the probe with `var zero V`, so the generated check passes no zero.

**Why not:** the spec gives the zero an input of the kind `value`, because a language without zero values passes the value that eidos derives. A Go probe that ignores the input differs from the probe of every other runtime.

### A waiter count on the controlled clock

testkit waits until a goroutine of the subject waits on its clock through `AwaitWaitersUntil`. The controlled clock of assert has no such method.

**Why not:** `testing/synctest` waits until every goroutine of a bubble is durably blocked, also on a timer of the subject that does not read the clock of the handle. A waiter count sees only the waiters of one clock, and a bubble does not need a new API of assert.

## Drawbacks

- **One closure per call input.** A generated check writes a closure for each `call` input of its probe, 102 inputs across the 71 probes, each adapted to its shape. A callable with two values returns a `Pair`.
- **A bubble for each clock check.** Six probes run inside `testing/synctest`. A subject that starts a goroutine that never ends fails such a check, because the bubble does not end.
- **No vacuity.** A probe of this RFC passes where its precondition is not reached. Each of the 71 probes changes once more for the design of vacuous checks.
- **Three gates of enumerations.** An enumeration without a list function, a validity function or a value outside its variants refuses the claims that need them, where the generated test of the enumeration checked the rest.

## References

| What | Where |
|---|---|
| testkit, the probes of its generated suites | <https://github.com/thesm-os/testkit/blob/main/engine/probe/probe.go> |
| testing/synctest | <https://pkg.go.dev/testing/synctest> |
