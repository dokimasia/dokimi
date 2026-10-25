---
adr: 0006
title: A check that counts allocations runs outside the parallel tests
status: Accepted
date: 2026-10-24
supersedes: RFC-0003 and RFC-0006 in part
superseded-by: none
rfc: 0003
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0006: A check that counts allocations runs outside the parallel tests

## Status

Accepted

## Context

RFC-0003 binds a function of a consumer to a callable with the directive `check`. `conform.Run` and `conform.Falsify` run each check of a suite as a parallel test. Issue 20 and issue 21 leave a ceiling of a steady state and a ceiling for each element of an input to such a check, which counts the calls with `assert.MaxAllocs`.

- `assert.MaxAllocs` counts the allocations of the whole process, so its test does not call `t.Parallel`. Each check of a suite builds and calls a subject of its own, so a count in a parallel test includes the allocations of every check that runs concurrently.
- RFC-0006 runs the allocation ceilings of a budget in `Test<Subject>Allocs` and `Test<Declaration>SuiteAllocs`, which do not run in parallel. `RunBudgets` and `FalsifyAllocs` run those tests.
- sharder needs three such checks: `placement.State.Apply` after a sequence of moves, `placement.NewState` for a cell of 1,000 shards, and the view of a cell in its routes machine.

## Decision

We will run a check that counts allocations with the allocation ceilings, in the tests that do not run in parallel.

- The directive `check` takes the boolean param `counted`, as in `//dokimi:check State.Apply defect=leaksAMove counted=true`. The suite marks such a check with the field `Counted` of `conform.Check`.
- `RunBudgets` applies each counted check of the suite and of the option `With` after the allocation ceilings. It falsifies a counted check of `With` before it applies it, as `Run` does.
- `FalsifyAllocs` falsifies each counted check of the suite on its planted defect.
- `Run`, `Falsify` and `Fuzz` leave a counted check out. A check that runs in a parallel test does not apply a counted check to another planted defect.
- A counted check skips in a run that does not check allocation ceilings, as a ceiling does.
- The suite of a declaration with a counted check has `Test<Declaration>SuiteAllocs`, and each of its subjects has `Test<Subject>Allocs`. `Test<Declaration>SuiteBudgets` excuses every check of the consumer, because the planted defect of a ceiling returns the zero results of each call.

## Alternatives Considered

### A check that runs once without t.Parallel inside Run

`Run` could run a counted check as a subtest without `t.Parallel`. The generated test of a subject calls `t.Parallel` itself, so other top-level tests of the package still run while the check counts.

### A test of the consumer outside the suite

The consumer could keep each counted check in a test of its own, as sharder's `TestStateAllocs` does today. Such a test has no planted defect, so falsification does not prove it.

## Consequences

**Positive:**

- A consumer states a ceiling of a steady state or of a large input as a check of the suite, and falsification proves it.

**Negative:**

- A check of a parallel test does not judge whether a planted defect also fails a counted check.
- A counted check does not run under the race detector or with coverage, because those builds do not count allocations.

## References

| What | Where |
|---|---|
| RFC-0003, which runs every check of a suite in a parallel test | [0003-the-suite-generator.md](../rfc/0003-the-suite-generator.md) |
| RFC-0006, whose tests of the allocation ceilings run the counted checks too | [0006-generated-inputs-fuzz-targets-benchmarks-and-budgets.md](../rfc/0006-generated-inputs-fuzz-targets-benchmarks-and-budgets.md) |
| Issue 9 of dokimi-go, a check that counts allocations runs in parallel | <https://github.com/dokimasia/dokimi-go/issues/9> |
| Issue 20, a ceiling of a steady state | <https://github.com/dokimasia/dokimi/issues/20> |
| Issue 21, a ceiling for each element of an input | <https://github.com/dokimasia/dokimi/issues/21> |
