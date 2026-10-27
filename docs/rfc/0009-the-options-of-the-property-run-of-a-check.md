---
rfc: 0009
title: The options of the property run of a check
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-27
updated: 2026-10-27
discussion: https://github.com/dokimasia/dokimi/issues/32
supersedes: RFC-0006 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0009: The options of the property run of a check

## Summary

A property check states the options of its own run, such as the shares of `prop.Require`, through the field `Options` of `conform.Check` and the param `options` of the directive `check`. The options are a function of a handle, so they can depend on the subject. An application passes them after its own options. Falsification and fuzz targets leave them out, because a planted defect changes the share of each outcome and a fuzzer does not apply a share.

## Motivation

conform runs the property of a check under `prop.ForAll` with the options of the application alone. They are the overrides of the suite, the generators of `conform.Using`, the cases of `conform.Cases` and the seed of `DOKIMI_SEED`. The check itself cannot add one.

A stateful machine needs the share of each outcome of its cases. sharder's machine of commits checks that the commits of a store are linearizable against a model, and it classifies the outcome of each commit. Its run required a share of each outcome, such as 0.1 for `ErrFenced`, so a change of the generators that stops producing an outcome failed the run. As a check of a suite, the machine keeps its calls of `Case.Classify`, and its run checks no share. A change of the generators that no longer produces a fenced commit then passes the check. One share applies only to a store that loses commits, so the options depend on the subject.

An application cannot know which outcomes a check classifies, and every subject of the suite needs the same shares. So the check states its options beside its statement. A share that only some subjects meet, such as the share of a lost commit, comes from the handle that the options read.

## Detailed design

### Invariants

| Invariant | Enforced by |
|---|---|
| An application of a property check passes the options of the check after its own options | `Check.run` of conform |
| Falsification, the isolation of a planted defect and a fuzz target leave the options of a check out | `Check.falsify`, `Check.against` and `Fuzz` of conform |
| A check with options has a property | The problems of an application, and `BadCheck` of the suite generator |
| The function of the options takes the handle of its check and returns a list of `prop.Option` | `BadCheck` of the suite generator |

### Components

| Component | Change |
|---|---|
| `conform.Check[S]` | The field `Options` |
| The application of a check in conform | A handle of its own calls `Options` before the run |
| The problems of an application | A check with options and without a property |
| The directive `check` | The param `options` |
| The code `BadCheck` | The signature of the options, and options on a check without a property |
| The template of the suites | The field `Options` of each check with options |

### The field `Options`

```go
type Check[S any] struct {
	// The fields Apply, Property, Defect, Doors, ID, Statement and Counted
	// as before.

	// Options returns the options of the property run of the check in an
	// application, such as the shares of prop.Require. The application
	// calls it before the run, with a handle of its own whose subject is a
	// fresh subject of the factory, and passes the options after its own.
	// Falsification and fuzz targets leave the options out. A check with
	// Options has a function Property.
	Options func(h *Handle[S]) []prop.Option
}
```

The options come after the options of the application. A later option of one kind replaces the earlier one, so a check can set its own number of cases through `prop.Cases`.

The handle of the options is a handle of the check on the test of the check. The function builds the subject of the handle when it calls `h.Subject()`. A function that does not read the subject does not build one. The handle ends when the function returns. On a seat without `Cleanup` its cleanups run at that point, and on a test they run when the test of the check ends.

### The entry points that apply the options

| Entry point | The options of a check |
|---|---|
| `Run`, `RunBudgets` and `Check.Run` | Applied after the options of the application |
| `Falsify`, `FalsifyAllocs` and `Check.Falsify` | Left out |
| The isolation of a planted defect | Left out |
| `Fuzz` | Left out |

Each entry point that leaves the options out has its own reason:

- The falsification of a check needs a failure at its statement on its planted defect. Its run ends at the first failing case and does not count a share.
- The isolation applies each other check to a planted defect, and a planted defect changes the share of the outcomes. A defect that never returns `ErrFenced` would fail a share of `ErrFenced` of another check. The falsification of the defect would then fail, although no statement fails.
- `prop.Fuzz` does not apply `Cases` or `Require`. One fuzz target also runs every property check of its callable.

### The directive `check`

The param `options` contains the name of a function of the package of the check, as `defect` does:

```go
// HistoriesAreLinearizable checks that the commits of a store are
// linearizable against the model of a scope.
//
//dokimi:check Store.Writer defect=admitsFencedCommits options=machineOptions
func HistoriesAreLinearizable(h *StoreHandle, c *prop.Case) { … }

// machineOptions returns the shares of the outcomes of the machine. A store
// that loses commits also returns ErrOutcomeUnknown.
func machineOptions(h *StoreHandle) []prop.Option {
	opts := []prop.Option{prop.Require(labelFenced, 0.1), prop.Require(labelMismatch, 0.1)}
	if _, faulty := h.Subject().(Faulty); faulty {
		opts = append(opts, prop.Require(labelUnknown, 0.1))
	}
	return opts
}
```

The annotate step stamps the identity of the function as a fact of the check, under the key `suite.options`. It reports `BadCheck` at the directive in each of these cases:

- The check takes no `*prop.Case`, so it has no property run.
- The function takes another number of parameters than one, the handle.
- The function returns another number of results than one, or a result that is not a list of `prop.Option`.

The generator writes the field into the check of the suite:

```go
{
	ID:        StoreChecks.Writer.HistoriesAreLinearizable(),
	Statement: "the commits of a store are linearizable against the model of a scope",
	Property:  HistoriesAreLinearizable,
	Defect:    admitsFencedCommits,
	Options:   machineOptions,
},
```

### The contract in each language

| Part | Go | TypeScript |
|---|---|---|
| The options of a check | The field `Options func(h *Handle[S]) []prop.Option` of `Check[S]` | The field `options` of `Check<S>`, a function of the handle that returns the options of a property run |
| The param of the directive | `options` of `check` | The same param |

The TypeScript runtime and its template tree do not exist yet. Their parts take these names when they do.

### Failure handling

| Failure | Result | Recovery |
|---|---|---|
| A run of a check with a share refutes the share or leaves it unmet | The test of the check fails as `CoverageUnmet` of the property engine | Change the generators, or the share |
| A check with options has no property | The application reports a problem and runs no check | Remove the options, or make the check a property check |
| The function of the options ends with a fault, such as a door that the handle lacks | The test of the check ends with the fault | Pass the option of the door |
| The directive `check` states options with another signature, or on a check without a case | `BadCheck` at the directive. The check is not bound | Correct the function or the check |

### Verification

- **conform.** The spec of a check covers the options in `Check.Run`, their order after the options of the application, a share that the run leaves unmet, the subject of the handle of the options, and a function that builds no subject. It covers the end of the handle before the first case, the options of a check without a property, and the isolation of a planted defect without the options. The spec of `Falsify` covers a suite without the options, and the spec of `Run` the problem of options on a check without a property. The spec of `Fuzz` covers a target without the options.
- **The suite generator.** The specs of the plugin cover the param `options` and each `BadCheck` of the options. The golden files pin the field `Options` of a generated check.
- **The corpus.** The crash check of the fixture `door/durable` classifies whether a crash kept the old or the new value, and its options require a share of each.

### Migration

The change adds a field, a param and the cases of an existing code, so no existing call breaks. A check without `Options` runs as before.

## Alternatives considered

### An option of the application

The application passes the shares as an option, such as `conform.Require(check, label, share)`.

**Why not:** every subject of the suite would repeat the shares of each check. A subject cannot know which outcomes a check classifies, and the shares belong with the statement of the check.

### The options in falsification and in fuzz targets

`Falsify` and `Fuzz` pass the options of a check after their own.

**Why not:** the isolation of a planted defect would fail a check whose share the defect changes, with no failure of a statement. The falsification of the check itself ends at its first failing case. `prop.Fuzz` does not apply a share.

### A list in place of a function

The directive states the options as values, such as `options=[prop.Require(labelFenced, 0.1)]`.

**Why not:** a share can depend on the subject, as the share of `ErrOutcomeUnknown` does. A function of the handle reads the subject. A directive also cannot spell a call of the property engine.

### A fuzz target for each check

Each property check gets a fuzz target of its own, so a target can apply the options of its check.

**Why not:** the fuzzer does not apply `Cases` or `Require`, so the options would change nothing in the target. One target for each callable keeps one entry point for each callable.

## Drawbacks

- **One more subject.** A check whose options read the subject builds one subject more for each application of the check.
- **A share is not falsified.** Falsification proves the statement of a check and leaves its shares out. So a share that a subject cannot meet fails only in an application.
- **Two params named options.** The param `options` of the directive `subject` returns options of an application, and the param `options` of the directive `check` returns options of a property run. The types differ, and a function of the wrong kind fails the compile of the generated file.

## Open questions

None.

## References

| What | Where |
|---|---|
| The request of the options of a check | <https://github.com/dokimasia/dokimi/issues/32> |
| `prop.Require` and `prop.Fuzz` of the assert library | <https://github.com/dokimasia/assert-go/blob/main/prop/option.go> |
