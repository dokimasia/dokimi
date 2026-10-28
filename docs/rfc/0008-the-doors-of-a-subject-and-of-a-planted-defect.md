---
rfc: 0008
title: The doors of a subject and of a planted defect
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-27
updated: 2026-10-27
discussion: https://github.com/dokimasia/dokimi/issues/31
supersedes: RFC-0003 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0008: The doors of a subject and of a planted defect

## Summary

A check of a store can restart its subject. It cannot crash the medium of the subject during a call, and the generated falsification of a suite does not give a planted defect any door. This RFC adds a fourth door, the crash, with the option `conform.Crashable` and the method `Handle.Crash`. It gives each check the doors of its planted defect, in the field `Doors` of `conform.Check` and in the params `reference`, `second`, `restart` and `crash` of the directive `check`. Falsification builds each handle on a planted defect with the doors of that defect. It also gives the directive `subject` the param `options`, a function that returns the options of the application, so the test that dokimi writes for a subject can pass every door.

## Motivation

A door is a method of the handle that only the code that builds a subject can supply: `Reference`, `Second` and `Restart`. An application passes each door as an option, such as `conform.Restartable`. These defects keep the suites of a storage engine out of dokimi:

- **A crash has no door.** `Handle.Restart` returns a subject after a restart, once the calls of the check have returned. A check cannot stop the medium in the middle of a write, so no check proves that a crash leaves the whole object or nothing. The test kits of `blob.Store` and `cas.Store` in core state that law through `WithCrash`, a function that runs a write and crashes the storage at a point of its choice.
- **The planted defect of a suite check has no door.** The generated spec of a suite calls `conform.Falsify(t, StoreSuite())` without options. A check of the directive `check` that calls `h.Restart` ends its falsification with a fault that wraps `conform.ErrOption`, and never with a verdict. conform's own spec pins that behaviour as "ends with the fault of a door that the check calls without its option". Every suite check that restarts, crashes or compares with a reference fails the generated spec of its suite.
- **The directive `subject` cannot pass a door.** It takes `new` and `cases`, and the test that it writes passes only the cases, the budgets and the excusals of the directives. A subject whose checks need a door writes its conformance test, its allocation test, a benchmark and a fuzz target of each callable by hand.

core moves the suites of `blob.Store` and `cas.Store` to dokimi. Their checks of a reopened store need `Restartable`, their checks of a crash need a door that does not exist, and their hash checks need `Reference`. A probe generated both suites with dokimi `a92d631` and conform v0.6.0. Its checks of a crash had no door to call.

conform returns the door of each handle, and the suite generator writes the checks and the tests that supply the doors. No layer above them sees a handle, so the change belongs in these two.

## Detailed design

### Invariants

| Invariant | Enforced by |
|---|---|
| A handle on a planted defect returns the doors of the check of the defect, and each door that the check lacks from the options of the falsification | `Check.falsify` of conform |
| A handle on a subject of an application returns the doors of the options of the application alone | `newHandle` of conform |
| A check of another check on a planted defect reads the doors of the defect, not its own | `Check.isolate` of conform |
| A door that no option and no check supplies ends the check with a fault that wraps `ErrOption` and names the remedy | The doors of `Handle` |
| A door function of a check has the signature of its door | The annotate step of the suite generator, under `BadCheck` |
| The options function of a subject takes no parameter and returns one list | The generate step of the suite generator, under `BadOptions` |
| Each generated test of a subject passes the options of its function after the cases, the budgets and the excusals of the directives | The template of the conformance tests |

### Components

| Component | Change |
|---|---|
| `conform.Doors[S]` | A new type with the four doors of a handle |
| `conform.Check[S]` | The field `Doors`, the doors of the planted defect |
| `conform.Crashable` and `Handle.Crash` | The option and the method of the crash door |
| The falsification of conform | Each handle on a planted defect takes the doors of its check first |
| The directive `check` | The params `reference`, `second`, `restart` and `crash` |
| The directive `subject` | The param `options` |
| The codes of the suite generator | `BadOptions`, and doors in the meaning of `BadCheck` |
| The template of the suites | The field `Doors` of each check with a door |
| The template of the conformance tests | The options of the function in the test, the allocation test, the benchmarks and the fuzz targets of a subject |

### The doors

`Doors` declares the four doors of a handle in one type. An application sets them through its options, and a check sets the doors of its planted defect:

```go
// Doors are the doors of a handle: the subjects, the restart and the crash
// that only the code that builds a subject can supply. A nil door is not
// supplied.
type Doors[S any] struct {
	// Reference returns a fresh reference implementation, which
	// [Handle.Reference] returns.
	Reference func(h *Handle[S]) S

	// Second returns a fresh second subject, which [Handle.Second] returns.
	Second func(h *Handle[S]) S

	// Restart returns the subject s after a restart of its process, which
	// [Handle.Restart] returns.
	Restart func(h *Handle[S], s S) S

	// Crash runs call against s, crashes the medium of s while call runs, and
	// returns the subject over the medium as the crash left it. It draws the
	// point of the crash from h.Rand(). [Handle.Crash] returns its subject.
	Crash func(h *Handle[S], s S, call func()) S
}
```

The options `Reference`, `Second` and `Restartable` set their field of the doors of the application. The option `Crashable` sets the fourth:

```go
// Crashable supplies the crash of the subjects of an application.
// [Handle.Crash] returns the subject that crash returns.
func Crashable[S any](crash func(h *Handle[S], s S, call func()) S) Option[S]
```

`Handle.Crash` is the fourth door of the handle:

```go
// Crash runs call against s and crashes the medium of s while call runs,
// through the option [Crashable] or the door Crash of a planted defect. It
// returns the subject over the medium as the crash left it. Without the door
// it ends the check with a fault that wraps [ErrOption]. A seat that takes the
// fault and goes on gets s after call runs without a crash.
func (h *Handle[S]) Crash(s S, call func()) S
```

The crash function draws the point of the crash from `h.Rand()`. In a property check that source belongs to the case, so a failure replays the same crash point from its seed and shrinks it. The function chooses the unit of the point, such as a write, a byte or a sync of its medium, because only the medium knows its units. `Crash` builds the new subject as `Restart` does, so a read of the clock in the crash function counts as a read by a subject.

The conformance spec has the machinery `killable-process-or-breakable-medium` for `Restart`. `Crash` supplies the same machinery:

| Method | Returns | Machinery |
|---|---|---|
| `Restart(s)` | The subject `s` after a restart of its process | `killable-process-or-breakable-medium` |
| `Crash(s, call)` | The subject over the medium of `s` after a crash during `call` | `killable-process-or-breakable-medium` |

### The doors of a planted defect

`Check` gains the doors of its planted defect:

```go
type Check[S any] struct {
	// The fields Apply, Property, Defect, ID, Statement and Counted as before.

	// Doors are the doors of the planted defect of the check. Falsification
	// builds each handle on the defect with them: in the falsification of the
	// check, and where it applies another check to the defect. A door that
	// Doors lacks comes from the options of the falsification.
	Doors Doors[S]
}
```

The planted defect of the check that a reopened store returns each stored value is a store that forgets on a restart. Its restart returns an empty store, and only that check fails on it. So the restart is part of the defect. Every other check of the suite that restarts the defect sees the same restart, and the isolation of the defect proves that no other check fails on it.

Falsification resolves each door of a handle on a planted defect in this order:

1. The door of the check whose planted defect the handle builds on.
2. The door of the options of the falsification: the options of `Check.Falsify`, or the options of the application that falsifies a check that `With` adds.

The suite of a declaration falsifies its checks without options, so a door of a suite check comes from the check alone. A check that `With` adds keeps the doors of its application where it has none of its own, as before. The isolation of a planted defect applies every other check to the defect with the doors of the defect. A check that calls a door that the defect lacks ends with a fault and has no verdict on the defect, as a check that calls a door without its option has today. For a property check this depends on the property engine: assert-go v0.9.0 fails a run whose case ends with a fault, with a counterexample in place of the fault. Until the engine ends such a run with the fault, a planted defect supplies each door that a property check of its suite calls.

The fault of a handle on a planted defect contains the door of the check, because no application exists in the spec of a suite:

```text
conform: the handle lacks the door: Store.Get/keeps-each-put-across-a-restart calls Restart on its
planted defect, and the check has no door Restart. Give the check the door:
	conform.Doors[durable.Store]{Restart: func(h *conform.Handle[durable.Store], s durable.Store) durable.Store { … }}
```

A handle on a subject of an application keeps the message of the option, such as "calls Restart without the option conform.Restartable. Pass the option to the application". The text of `ErrOption` becomes "conform: the handle lacks the door", because a planted defect lacks a door without an application.

### The directive `check`

The directive `check` takes a param for each door of the planted defect. The value of each is a function of the package of the check, as the value of `defect` is:

| Param | Signature | The function returns |
|---|---|---|
| `reference` | `func(h *StoreHandle) Store` | A reference implementation for the checks on the planted defect |
| `second` | `func(h *StoreHandle) Store` | A second subject beside the planted defect |
| `restart` | `func(h *StoreHandle, s Store) Store` | The planted defect `s` after a restart |
| `crash` | `func(h *StoreHandle, s Store, call func()) Store` | The planted defect over its medium after a crash during `call` |

```go
// KeepsEachPutAcrossARestart checks that Get returns each value that Put
// stored before a restart.
//
//dokimi:check Store.Get defect=forgetsOnRestart restart=forget
func KeepsEachPutAcrossARestart(h *StoreHandle, c *prop.Case) { … }

// forgetsOnRestart returns a journal over a medium in memory, which forget
// replaces with an empty journal.
func forgetsOnRestart(h *StoreHandle) durable.Store { return NewJournal(h) }

// forget returns an empty journal in place of the journal s.
func forget(h *StoreHandle, _ durable.Store) durable.Store { return NewJournal(h) }
```

The annotate step stamps the identity of each door function as a fact of the check, under the keys `suite.reference`, `suite.second`, `suite.restart` and `suite.crash`. It reports `BadCheck` at the directive where a door function has another number of parameters or results than its signature, and where the third parameter of a crash is no function. The generator writes the doors into the check of the suite:

```go
{
	ID:        StoreChecks.Get.KeepsEachPutAcrossARestart(),
	Statement: "Get returns each value that Put stored before a restart",
	Property:  KeepsEachPutAcrossARestart,
	Defect:    forgetsOnRestart,
	Doors:     conform.Doors[durable.Store]{Restart: forget},
},
```

A check without a door param has no field `Doors`, so the generated suites of every other check do not change.

### The directive `subject`

The value of the param `options` of the directive `subject` is a function without parameters that returns the options of the application. eidos resolves it as it resolves `new`, to an exported function of the package of the struct or a function of the package of its test.

```go
// Journal is the Store over a medium in memory.
//
//dokimi:subject durable.Store new=NewJournal options=JournalOptions
type Journal struct{ … }

// JournalOptions returns the doors of a journal: a restart that reopens its
// medium, and a crash that cuts a write of its medium.
func JournalOptions() []StoreOption {
	return []StoreOption{conform.Restartable(reopen), conform.Crashable(crash)}
}
```

The generator reports `BadOptions` at the subject in each of these cases:

- The function takes a parameter.
- The function returns another number of results than one, or a result that is not a list.
- The package of the test cannot call the function.

A subject with the finding has no conformance test, as a subject with `BadFactory` has none.

The test, the allocation test, the benchmarks and the fuzz targets of the subject pass the options after the cases, the budgets and the excusals of the directives. An option of the function then replaces an earlier option of its kind. Go does not accept single arguments and a spread list in one variadic call. So each function with an option of a directive builds the list first:

```go
// TestJournalConformsToStore applies the checks of the suite of Store to Journal.
func TestJournalConformsToStore(t *testing.T) {
	t.Parallel()

	opts := append([]durabletest.StoreOption{
		conform.Cases[durable.Store](20),
	}, durabletest.JournalOptions()...)
	durabletest.RunStoreSuiteFromHandle(t, durabletest.NewJournal, opts...)
}
```

A function without an option of a directive passes the list alone, as `durabletest.JournalOptions()...`. A subject without the param `options` gets the same tests as before.

### The contract in each language

| Part | Go | TypeScript |
|---|---|---|
| The doors | `conform.Doors[S]` with `Reference`, `Second`, `Restart` and `Crash` | `Doors<S>` with `reference`, `second`, `restart` and `crash` |
| The doors of a planted defect | The field `Doors` of `Check[S]` | The field `doors` of `Check<S>` |
| The crash of an application | The option `Crashable(crash)` | The field `crash` of `Options<S>` |
| The crash door | `Handle.Crash(s, call)` | `handle.crash(s, call)`, which awaits `call` |
| The params of the directives | `reference`, `second`, `restart` and `crash` of `check`, and `options` of `subject` | The same params |

The TypeScript runtime and its template tree do not exist yet. Their parts take these names when they do.

### Failure handling

| Failure | Result | Recovery |
|---|---|---|
| A check calls a door that its planted defect lacks, in the spec of its suite | The falsification of the check ends with a fault that wraps `ErrOption` and names the door of the check | Give the check the door, such as `restart=forget` |
| A check calls a door that a subject lacks | The check ends with a fault that wraps `ErrOption` and names the option | Pass the option, or excuse the check |
| Another apply check calls a door that a planted defect lacks | The check has no verdict on the defect | None. The isolation covers the other checks |
| Another property check calls a door that a planted defect lacks | The run of the check fails with a counterexample, so the falsification of the defect fails | Give the defect the door |
| A door function of a check has another signature | `BadCheck` at the directive. The check is not bound | Correct the signature |
| The options function of a subject has another signature, or the test cannot call it | `BadOptions` at the subject. The subject has no conformance test | Correct the function |
| A crash function does not crash during `call` | The check sees the medium after a whole call | Draw the point of the crash among the units of the call |

### Verification

- **conform.** The spec of the handle covers `Crash` with and without its door, each door of a planted defect, each door of the options on a defect whose check lacks it, and the message of each door on a planted defect. The spec of the options covers `Crashable` in an application. The spec of a check covers a door of the check over a door of the options, the doors of a defect in the isolation of another check, and an application that leaves out the doors of a defect. The spec of `Falsify` covers a suite check that restarts its planted defect.
- **The suite generator.** The specs of the plugin cover each door param, each `BadCheck` of a door signature, the param `options` in each of the four tests of a subject, and each `BadOptions`. The golden files of the fixture tree pin the generated checks and tests.
- **The corpus.** The fixture `door/durable` applies a suite with a restart check and a property check of a crash to a journal over a medium in memory. The journal passes its doors through `options`. The generated spec of the suite falsifies both checks on planted defects with a restart and a crash of their own.

### Migration

The change adds a type, a field, an option, a method, five params and a code, so no existing call breaks. The text of `ErrOption` changes, and `errors.Is` matches the fault as before.

## Alternatives considered

### The doors of the suite

The directive `suite` takes the doors of every planted defect, as `//dokimi:suite restart=forget`, and the generated spec passes them to `Falsify`.

**Why not:** a planted defect whose fault is in its restart needs a restart of its own. A restart of the suite applies to every defect, so a restart that forgets fails every check that restarts a correct defect. The doors of a check apply to its defect alone.

### The doors of the application on a planted defect

Falsification passes the options of the subject to every planted defect, as it does for a check that `With` adds today.

**Why not:** the restart of a store reopens the medium of that store, and a planted defect is a double over no such medium. The spec of a suite also runs in the package of the suite without a subject, so no application exists to supply the options. The options remain the second source of a door for a check that `With` adds.

### A planted defect that restarts itself

The runtime asserts an optional method on the defect, such as `Restart(h *Handle[S]) S`, and calls it in place of the door.

**Why not:** the runtime would assert an undeclared interface on a value of the consumer's type. A misspelled method would fall back to the fault without a word about the spelling, and a defect built on the generated double would need a type of its own for the method.

### A param of the directive `subject` for each door

The directive `subject` takes `reference`, `second`, `restart` and `crash`, as the directive `check` does.

**Why not:** the doors of an application share state, such as the medium that a restart reopens and a crash cuts. One function builds that state. A function of options also passes `Using`, `With`, `Cases` and every later option without another param.

### A crash point that the runtime draws

The runtime draws an integer and passes it to the crash function as the point of the crash.

**Why not:** the unit of the point depends on the medium: a write of a journal, a byte of a file or a sync of a database. The crash function draws from `h.Rand()`, which is the source of the case in a property check, so the point replays and shrinks either way.

### A crash through a fault of the double

The check arms a fault of the generated double, such as an error of the next write, in place of a crash.

**Why not:** a fault of a double fails one call of the double. A crash stops the medium of a real store in the middle of a write. The double has no medium, so its fault cannot stop a write halfway.

## Drawbacks

- **Four doors in two places.** An application sets a door through an option, and a check sets the door of its planted defect through `Doors`. A reader of a suite looks in both places to learn which door a handle returns.
- **A model of the medium in each door check.** The planted defect of a crash check needs a medium that the crash can cut, such as a map of records with a count of writes. The fixture `durable` spends 53 lines on its medium, 33 on the doors of its control and 28 on the doors of its two defects.
- **Thinner isolation of a door defect.** Every check that calls a door that a planted defect lacks has no verdict on that defect, so the isolation of a door defect skips those checks.
- **Every door on every defect, for now.** A property check fails the isolation of a planted defect that lacks one of its doors, until the property engine ends a run with the fault of a case. The fixture `durable` gives both of its defects a restart and a crash for that reason.
- **Five params and one code.** The directive `check` grows from two params to six, the directive `subject` from two to three, and the suite generator from ten codes to eleven.

## Open questions

None.

## Unresolved and future work

- No derived claim of the conformance spec needs a crash. A claim of the spec that does, such as the atomicity of a batch writer under a crash, is not proposed here.
- The TypeScript runtime and the TypeScript template tree of the suite generator do not exist. The table of the contract fixes their names.

## References

| What | Where |
|---|---|
| The request of the options of a subject | <https://github.com/dokimasia/dokimi/issues/31> |
| The request of the crash door | <https://github.com/dokimasia/dokimi-go/issues/15> |
| The fault of a property case in assert-go | <https://github.com/dokimasia/assert-go/issues/73> |
| The crash of the test kit of `blob.Store` in core | <https://github.com/thesm-os/core/blob/93715491092dae3b1bf2db2580301cfbfb015046/coretest/blobtest/blobtest.go#L79-L91> |
| The crash of the test kit of `cas.Store` in core | <https://github.com/thesm-os/core/blob/93715491092dae3b1bf2db2580301cfbfb015046/coretest/castest/castest.go#L60-L71> |
| The restart door of conform | <https://github.com/dokimasia/dokimi-go/blob/main/conform/handle.go> |
