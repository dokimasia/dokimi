---
rfc: 0015
title: Vacuous checks and controls
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-28
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

# RFC-0015: Vacuous checks and controls

## Summary

Falsification proves that a check can fail. It does not prove that the check exercised its claim on the subject, or that the check accepts every correct implementation. This RFC adds both proofs to conform.

- **Vacuity.** A check is vacuous on a subject when its probe never gets the subject into the state of its claim, such as a check of the zero value beside an error on a subject whose calls never fail. An application fails each check that is vacuous on its subject, unless an excusal covers it.
- **Controls.** A control is a correct implementation that differs from the subject. The corpus applies the suite to each control through `conform.Run`, as RFC-0003 states, and fails each check that a control fails, because such a check forbids what the declaration permits. The corpus has a control for each shape family.

## Motivation

A probe of a reader asserts that a failed call returns the zero value. When no call of the subject fails, the assertion never runs, and the check passes. The suite is green and asserts nothing about that law. No runner output reports it. Today conform has no vacuous outcome, no census and no control (`conform/check.go`, `conform/suite.go`).

The literature names both defects. A formula holds vacuously when a subformula does not affect its truth, and IBM measured that about 20% of specifications passed vacuously on their first runs (Kupferman and Vardi). Ball and Kupferman applied the idea to tests: when a mutated system satisfies a specification, part of the system is not covered. That second case is a planted defect that passes, which falsification already catches. The first case, a claim that the subject never exercises, is what this RFC adds. QuickCheck states it as coverage: `cover` with `checkCoverage` fails a property whose cases do not meet a labelled condition often enough. JUnit separates an aborted test, whose assumption failed, from a skipped one.

Over-strong checks are the opposite defect. Test262 refuses "any test that restricts potentially valid extensions" at review. Jahangirova and colleagues call an assertion that fails on a correct state a false positive, and find them by running the assertion against generated correct states. A check that only the subject passes looks like a working check until someone writes a second correct implementation. In the Go corpus every fixture applies its suite to its own `InMemory` subject alone, so no fixture tests that its checks accept another correct implementation.

## Detailed design

### Engagement

A probe calls `Handle.Engage` when it has put the subject into the state of its claim:

```go
// Engage records that the check put the current subject into the state of
// its claim, such as a call that failed for a check of the zero value beside
// an error. A check whose run ends without a call of Engage is vacuous on
// that subject.
func (h *Handle[S]) Engage()
```

The spec states the condition of each probe in its field `vacuous`, and the Go implementation of the probe calls `Engage` where that condition stops being true. An authored check is engaged when the assertion whose contract is its statement runs at least once, whatever its verdict. The handle implements `assert.Observer` of assert-go, from v0.11.0, so it observes every assertion call, a passing one included:

```go
// ObserveCall receives one assertion call: the id of the assertion and the
// contract that the caller passed. The first Observer of the chain of a
// seat receives each call once, before the seat receives its failure.
func (h *Handle[S]) ObserveCall(call assert.Call)
```

The handle marks an authored check engaged when a call has the statement of the check as its contract. assert-typescript gives a `Seat` the optional member `observeCall`, which receives a `Call` with `assertion` and `contract`, so the TypeScript handle observes the calls in the same way.

### The outcome of a check on a subject

| Result | Meaning |
|---|---|
| `pass` | The check passed and engaged |
| `vacuous` | The check passed and never engaged |
| `fail` | The check failed |
| `excused` | An excusal of the consumer covers the check |
| `skipped` | The check skipped through `Handle.Skipf`, because the environment lacks what it needs, and no excusal covers it |

A property check engages when at least one of its cases engages. A case of a probe whose checker can leave a case undecided, such as a history that exceeds the budget of a checker of linearizability, engages only when the checker decides the case. A `skipped` check does not fail the application. The suite record reports it with the reason of the skip, so a run on a system without direct I/O passes and shows what it did not check.

### The census of an application

`conform.Run` counts, for each check, whether it engaged on the subject of the application. When the application ends, it fails the test of each check that passed without engaging and that no excusal covers. An application to a control is an application to a subject, so a check that a control never engages fails it the same way, and the corpus excuses such a check for that control with a reason.

Each application keeps its own census, because `go test` runs the applications of a package as separate tests and does not call a hook of conform after the last of them. `dokimi report` reads the records of every application of a run and lists each check that was vacuous in every one. The runtime spec contains the vacuous condition of each probe. The message of `Run` names the check, the vacuous condition of its probe as the spec states it, and both remedies. The second remedy contains both forms of an excusal:

```text
Store.Read/zero-on-error never engaged: no call fails, with the sample
inputs or with a cancelled context. Pass a subject that engages the check,
or excuse the check with the directive excuse or conform.Excuse and a
reason.
```

### Controls

A control is a correct implementation of the declaration that differs from the subject in its representation, such as a store backed by a slice beside a store backed by a map. The corpus applies the suite to each control through `RunStoreSuite` or `RunStoreSuiteFromHandle`, as it applies the suite to a subject. A check that fails on a control forbids behaviour that the declaration permits, and the test of the control fails with the identity of the check.

### The corpus

A shape family is a group of the shapes of the catalog whose subjects in the corpus share one representation. A control of a family tests the probes and the templates that every fixture of the family shares. A family per shape of the catalog would need about 100 controls, and most of those shapes share a map. Each shape family of the corpus has at least one control that is not the subject:

| Family | Subject | Control |
|---|---|---|
| Stores and readers | A map | A sorted slice with binary search |
| Appenders and logs | A slice | A chunked list that grows by blocks |
| Caches | A map in front of a backing | A cache of one entry that evicts on each put |
| Queues | A buffered channel | A ring buffer |

A test of the corpus declares the family of each fixture in a table. It fails for a fixture outside every family and for a family without a control. The ledger of the corpus lists the control of each fixture. The verdict reverses when a check is over-strong on one shape of a family alone.

## Alternatives considered

### A minimum share of engaged cases

Each check states a share of its cases that must engage, as QuickCheck `cover` states a percentage with a confidence.

**Why not:** a derived check runs a fixed sequence on most subjects, so a share has one case to count. One engagement per check across the subjects of an application is the claim that a single caller can make. A property check can still require shares through `prop.Require` in its options.

### Vacuity as a skip

A vacuous check skips, so the test runner shows it.

**Why not:** a skip reads as an intended omission, and every runner hides skips in its summary. A vacuous check that nobody excused is a defect of the suite or of the subject, so it fails.

### Vacuity reported by the probe

A probe reports vacuity where it finds it, as testkit did with a skip whose message starts with `vacuous:`. A check that reports nothing is engaged.

**Why not:** a probe that forgets the report passes vacuously and nothing shows it. With `Engage`, a probe that forgets the call makes every check of its claim vacuous on a correct subject, and the corpus fails at once.

### One census over every application of a package

Every application of a test binary adds to one census, and the census fails a check that engaged on no subject of any application.

**Why not:** `go test` does not call a hook of a library after the last test of a package. Only a `TestMain` of the package could fail the census, and a generated `TestMain` conflicts with the one that a consumer writes. `dokimi report` reads every record of the run and gives the same view.

### An application to controls

`conform.Controls` applies the suite to a set of controls, and the suite record has the application `controls`, so a reader blames a failure on the check and not on the implementation.

**Why not:** the runtime runs a control as it runs a subject, so the two applications differ only in the record. The controls of dokimi are in its corpus, which fails on every failed check, so the record adds nothing that the corpus does not enforce. RFC-0003 rejects the same application. The verdict reverses when the service needs to tell the failures of a control of a consumer from the failures of a subject.

### Generated controls

The generator writes a second implementation of each declaration to serve as the control.

**Why not:** a generated control comes from the same declaration and the same templates as the derived checks. A check and a control that both come from one source agree by construction, so the control proves nothing about the check. A control from another author does.

## Drawbacks

- **A newer assert library.** The vacuous outcome of an authored check needs assert-go v0.11.0, and the release of assert-typescript with `observeCall`.
- **Every probe must call Engage correctly.** A probe that never calls it makes every check of its claim vacuous. A probe that calls it too early hides vacuity. The corpus tests each probe against a subject that does not engage it.
- **A control per shape family in the corpus.** Each family needs a second correct implementation, written by hand. A fixture whose subject has a representation outside the table adds a family and its control.
- **More failures for consumers.** A suite that passes today can fail with a vacuous check after the upgrade, and the consumer adds a subject that engages it or an excusal.
- **A census per application.** A check that is vacuous on the subject of one application fails that application, also when another application of the same package engages it. The consumer excuses it in the application that cannot engage it.

## References

| What | Where |
|---|---|
| Kupferman and Vardi, Vacuity detection in temporal model checking, STTT 2003 | <https://link.springer.com/article/10.1007/s100090100062> |
| Ball and Kupferman, Vacuity in Testing, TAP 2008 | <https://www.cs.huji.ac.il/~ornak/publications/tap08.pdf> |
| QuickCheck, cover and checkCoverage | <https://hackage.haskell.org/package/QuickCheck/docs/Test-QuickCheck.html> |
| JUnit, assumptions and aborted tests | <https://docs.junit.org/current/writing-tests/assumptions.html> |
| Test262, the rule on valid extensions | <https://github.com/tc39/test262/blob/main/CONTRIBUTING.md> |
| Jahangirova, Clark, Harman and Tonella, Test oracle assessment and improvement, ISSTA 2016 | <https://dl.acm.org/doi/10.1145/2931037.2931062> |
