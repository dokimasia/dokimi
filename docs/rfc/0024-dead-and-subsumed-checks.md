---
rfc: 0024
title: Dead and subsumed checks from mutation records
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-28
updated: 2026-11-03
discussion: none
supersedes: none
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0024: Dead and subsumed checks from mutation records

## Summary

`dokimi report --mutation <record>` reads the record of a mutation run of mutate-spec beside the suite records of a run, and attributes each killed mutant to the checks that killed it. A check that ran on a subject whose code the mutation run changed, and that killed no mutant that it covered, is dead: an error. A check whose killed mutants are all killed by another check of the same suite is subsumed: a warning. The attribution needs every test that kills a mutant and the mutants that each check covers. A mutation run stops at the first killer of a mutant, and it records the coverage of each top-level test, while conform runs each check as a subtest. dokimi requests two modes of mutate-spec: one that records every killing test, and one that records the coverage of each subtest.

## Motivation

A planted defect proves that a check catches its defect. It does not show what else the check catches, or whether another check already catches all of it. A suite can contain checks that kill nothing beyond their own defect, and checks that duplicate another check. Both cost run time, and a dead check shows that the law it states is not what the code exercises.

A mutation run measures what the tests of a finished suite detect, and mutate-spec defines its record. Two fields of each mutant matter here (`mutate-spec/spec/record.schema.json:289-294`):

- `tests` lists the tests that failed on the mutant.
- `coveredBy` lists the tests whose run alone executed the site of the mutant.

conform runs each check as a subtest of the test of its application, named by the identity of the check without its declaration, such as `TestInMemoryStore/Get/zero-on-error`. A run stops at the first failing test of a mutant (`mutate-spec/docs/rfc/0002-a-mutation-run-and-its-record.md:142-145`), so `tests` lists only the first killer, and subsumption needs every killer. The run records the coverage of each test that the runner can run on its own, which is a top-level test function in Go (`mutate-spec/docs/rfc/0002-a-mutation-run-and-its-record.md:129-133`). For a check, `coveredBy` names the test of its application, which runs every check of the application.

## Detailed design

### Attribution

`dokimi report` joins a name in `tests` or `coveredBy` with a check through the suite records. The name is the name of the test that wrote a record, a slash and the identity of a check of the record without its declaration. `TestInMemoryStore/Get/zero-on-error` is the check `Store.Get/zero-on-error` of the records of `TestInMemoryStore`. The TypeScript binding registers the Vitest test of each check under the identity of the check without its declaration, inside the block of its application, so a record of an engine for TypeScript joins by the name of the block in the same way.

For each mutant of the record whose `coveredBy` names a check of a suite record:

| Verdict of the mutant | Attribution |
|---|---|
| `killed` | Each check in `tests` killed it |
| `survived` | Each check in `coveredBy` covered it and killed nothing |
| Other verdicts | No attribution |

### The findings

| Finding | Severity | Condition |
|---|---|---|
| `DeadCheck` | Error | The check covered at least one mutant and killed none |
| `SubsumedCheck` | Warning | Every mutant that the check killed is killed by one other check of the suite, which kills at least one mutant more |
| `UnattributedMutant` | Notice | A mutant survived and no check covered it |

A dead check is an error, because a law that no change of the code it covers can break is either stated wrongly or covers code that does not implement it. A subsumed check is a warning, because a check can state a law that the deeper check witnesses only indirectly.

`DeadCheck` and `SubsumedCheck` need both modes of the requests. With the first killer alone, a check that runs after another killer of each of its mutants looks dead and subsumed. With the coverage of the test of an application, each check of the application looks as if it covered every mutant that the application covered. A record without both modes does not give dead or subsumed checks, and `dokimi report` prints the mode that the record lacks. `UnattributedMutant` works without the modes, because a survivor that the test of no application covered is covered by no check.

### The requests to mutate-spec

A mode of the run that does not stop at the first failing test of a mutant, and records every failing test in `tests`:

```json
{ "verdict": "killed", "tests": ["TestInMemoryStore/Get/zero-on-error", "TestInMemoryStore/Get/read-after-write"], "allKillers": true }
```

The field `allKillers` states that `tests` is complete.

A mode of the run that runs each subtest alone, as the run runs each top-level test alone, and names the subtests in `coveredBy`:

```json
{ "verdict": "survived", "coveredBy": ["TestInMemoryStore/Get/zero-on-error", "TestInMemoryStore/Get/cancel"] }
```

The record states the mode, so that a reader can tell the test of a check from the test of an application. A subtest that runs alone also runs the code of its parent test before the subtest starts, so the coverage of a check contains the code of its application that runs before the check.

### The command

```text
$ dokimi report run.json --mutation mutation.json
example.com/shop/catalog  Store
  mutants     41 covered by checks   38 killed   3 survived
  dead        Store.Close/smoke: covered 4 mutants, killed none
  subsumed    Store.Get/smoke: every kill also by Store.Get/read-after-write
  unattributed 3 survivors in store.go:88-94, covered by no check
```

## Alternatives considered

### Mutants as the proof of each check

dokimi drops the planted defects and proves each check by the mutants that it kills.

**Why not:** a mutant changes an arbitrary expression and does not target a law. A killed mutant shows that some check failed, not that the check of a given law can fail for its reason. RFC-0001 records this choice. Mutation measures the suite beside the proofs.

### Attribution by coverage alone

A check that covers code and passes on every mutant of it is dead, computed from `coveredBy` alone.

**Why not:** that is the definition above for the dead check. Subsumption still needs the set of killers, which `coveredBy` does not give.

## Drawbacks

- **Two requests to mutate-spec.** Dead and subsumed checks wait for a mode that records every killer and a mode that records the coverage of each subtest. The first costs one run of every covering test per killed mutant. The second costs one run of each check alone per test binary, such as 180 runs for three applications of a suite of 60 checks.
- **A mutation run per report.** Attribution needs a mutation run, which costs hours on a large suite.
- **Names as the join.** The join relies on the name of each test being the identity of its check. A runner that renames subtests breaks the join.

## Open questions

None.

## References

| What | Where |
|---|---|
| mutate-spec, a mutation run and its record | <https://github.com/dokimasia/mutate-spec/blob/main/docs/rfc/0002-a-mutation-run-and-its-record.md> |
