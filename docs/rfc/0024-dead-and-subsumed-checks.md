---
rfc: 0024
title: Dead and subsumed checks from mutation records
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Draft
created: 2026-10-28
updated: 2026-10-28
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

`dokimi report --mutation <record>` reads the record of a mutation run of mutate-spec beside the suite records of a run, and attributes each killed mutant to the checks that killed it. A check that ran on a subject whose code the mutation run changed, and that killed no mutant that it covered, is dead: an error. A check whose killed mutants are all killed by another check of the same suite is subsumed: a warning. The attribution needs every test that kills a mutant, and a mutation run stops at the first. dokimi requests a mode of mutate-spec that records every killing test.

## Motivation

A planted defect proves that a check catches its defect. It does not show what else the check catches, or whether another check already catches all of it. A suite can contain checks that kill nothing beyond their own defect, and checks that duplicate another check. Both cost run time, and a dead check shows that the law it states is not what the code exercises.

A mutation run measures what the tests of a finished suite detect, and mutate-spec defines its record. Two fields of each mutant matter here (`mutate-spec/spec/record.schema.json:262-267`):

- `tests` lists the tests that failed on the mutant.
- `coveredBy` lists the tests whose run executed the site of the mutant.

conform runs each check as a test named by the identity of the check, so a name in either field is the identity of a check. A run stops at the first failing test of a mutant (`mutate-spec/docs/rfc/0002-a-mutation-run-and-its-record.md:145`), so `tests` lists only the first killer, and subsumption needs every killer.

## Detailed design

### Attribution

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

### The request to mutate-spec

A mode of the run that does not stop at the first failing test of a mutant, and records every failing test in `tests`:

```json
{ "verdict": "killed", "tests": ["TestStore/Store.Get/zero-on-error", "TestStore/Store.Get/read-after-write"], "allKillers": true }
```

The field `allKillers` states that `tests` is complete. Without it, `dokimi report` reports dead checks and no subsumed checks, because a check that is not the first killer would look subsumed.

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

- **A request to mutate-spec.** Subsumption waits for a mode that records every killer, and that mode costs one run of every test per mutant.
- **A mutation run per report.** Attribution needs a mutation run, which costs hours on a large suite.
- **Names as the join.** The join relies on the name of each test being the identity of its check. A runner that renames subtests breaks the join.

## Open questions

- Does the TypeScript runtime name each check as the test of Vitest under the identity of the check, so that Stryker records it?

## References

| What | Where |
|---|---|
| mutate-spec, a mutation run and its record | <https://github.com/dokimasia/mutate-spec/blob/main/docs/rfc/0002-a-mutation-run-and-its-record.md> |
