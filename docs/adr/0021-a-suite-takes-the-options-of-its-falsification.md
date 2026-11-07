---
adr: 0021
title: A suite takes the options of its falsification
status: Accepted
date: 2026-11-07
supersedes: RFC-0008 in part
superseded-by: none
rfc: 0008
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0021: A suite takes the options of its falsification

## Status

Accepted

## Context

The generated test of a suite falsifies the suite without options, as `conform.Falsify(t, storetest.StoreSuite())`. Each handle on a planted defect takes its doors from the check of the defect, and each door of a check is a function of the package of the check. A kit cannot build every value that its planted defects take:

- The planted defects of a kit of replication calls wrap a reference over a leader. Only the composition root of the program builds a leader. The import rules of the program keep each kit below the composition root, so neither the kit nor a door of its checks can build the reference.
- The external test package of the kit can import the composition root. The generated test of the suite is in that package.
- `conform.Falsify` and `conform.FalsifyAllocs` take the options of an application, as `conform.Run` takes them. A handle on a planted defect takes each door that the check of the defect lacks from these options. `Handle.Reference` and `conform.Open` then return the reference and the values of the options.

## Decision

We will add the param `falsify` to the directive `suite` for the function of the options of its falsification, because only a package above a kit can build some values that the planted defects of the kit take.

- The value of the param is a function without parameters that returns the list of the options, as the value of the param `options` of the directive `subject` is. eidos resolves a bare name in the package of the declaration first. It then resolves the name in the external test package of the directory of the declaration, and in the external test package of the directory of its directive `out`.
- The test of the suite passes the list to `conform.Falsify` and `conform.FalsifyAllocs` after the suite. The other tests of the suite excuse every check of the suite and take no options.
- The suite generator reports `BadFalsify` at a declaration whose function the test cannot call without arguments, or whose result is not one list. It records the suite of such a declaration without its test.

```go
// Store is a store of the replicas of a leader.
//
//dokimi:suite falsify=falsifyOptions
//dokimi:out path=storetest/
type Store interface{ … }
```

```go
package storetest_test

// falsifyOptions returns the options of the falsification of Store: the
// reference over a leader of the composition root.
func falsifyOptions() []storetest.StoreOption {
	return []storetest.StoreOption{conform.Reference(newLeaderStore)}
}
```

The test of the suite then calls `conform.Falsify(t, storetest.StoreSuite(), falsifyOptions()...)`.

## Alternatives Considered

### The doors of the checks alone

Each door of a planted defect remained a function of the package of its check, such as `reference=newReference`.

The package of a check is the kit. The kit is below the composition root, so it cannot import the package that builds the reference, and no door of its checks can build it.

### A param of the directive suite for each door

The directive `suite` took the params `reference`, `second`, `restart` and `crash`, and the test of the suite passed each door as an option.

The doors of an application share state, such as the leader that a reference wraps, and one function builds that state. A function of options also passes `conform.Supplies` and every later option without another param. The directive `subject` takes one function of options for the same reasons.

## Consequences

**Positive:**

- A kit takes the reference and the values of its planted defects from its external test package, which can import a package above the kit.
- The test of a suite without the param does not change.

**Negative:**

- A door of a planted defect has two sources: the check of the defect and the function of the param. A reader of a kit reads its external test package to learn which reference a planted defect wraps.
- The planted defects of such a kit pass their falsification only with the options. A call of `conform.Falsify` without them ends each check that calls the door with a fault that wraps `conform.ErrOption`.
- The generator does not check that the options supply each door that the checks call. A door that neither source supplies fails at run time, in the falsification of its check.
- The directive `suite` grows from two params to three, and the suite generator from 16 codes to 17.

## References

| What | Where |
|---|---|
| RFC-0008, whose falsification of a suite without options this decision opens to a function of options | [0008-the-doors-of-a-subject-and-of-a-planted-defect.md](../rfc/0008-the-doors-of-a-subject-and-of-a-planted-defect.md) |
| Issue 60, the generated spec of a suite cannot pass options to Falsify | <https://github.com/dokimasia/dokimi/issues/60> |
| Issue 31 of dokimi-go, a planted defect cannot wrap a reference that the caller of its kit builds | <https://github.com/dokimasia/dokimi-go/issues/31> |
