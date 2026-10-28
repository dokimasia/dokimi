---
rfc: 0022
title: Verified fakes
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

# RFC-0022: Verified fakes

## Summary

A new generator, `fake`, writes an in-memory implementation of a declaration whose callables all belong to contracts of the shape catalog with a reference form in the conformance spec, such as `persister`, `cache` and `appender`. The generated test of the fake applies the suite of the declaration to the fake, so the fake passes the same checks as every production implementation that the consumer applies the suite to. A downstream team uses the fake in its own tests in place of a container or a hand-written fake. The suite record of the fake states that it passed, at which version of the spec, so a consumer can see which laws the fake passed. A declaration with a callable outside every reference form gets no fake and a refusal that names the callable.

## Motivation

A downstream team tests its code against a dependency in three common ways, and each has a cost:

- A container of the real service costs seconds per test and a daemon in CI.
- A hand-written fake drifts from the real implementation, because no suite checks it against the same contract.
- A double returns what the test configures and states no behaviour.

The practice that checks a fake against its contract is a shared suite. Google requires that "a fake must have its own tests to ensure that it conforms to the API of its corresponding real implementation", and that the owner of the real implementation writes the fake. gocloud runs one conformance suite, `drivertest.RunConformanceTests`, against the in-memory blob driver and against the S3 driver. The standard library checks every `fs.FS`, `MapFS` included, with `fstest.TestFS`. No published study measures how far fakes drift from the implementations that they replace.

The Go corpus contains 141 hand-written `InMemory` subjects, one per fixture. Most of them implement a store, a cache, a log or a queue in a map or a slice. A downstream team writes the same kind of in-memory implementation for its own tests, and no suite checks it.

## Detailed design

### Reference forms

The spec gains `references.yaml`. A reference form states the state of a contract and the effect of each role on it, in terms that each language renders:

```yaml
references:
  persister:
    contract: persister
    state: {kind: map, key: key, value: value}
    roles:
      writer: {writes: [key, value]}
      reader: {reads: key, absent: not-found}
    since: 1.0.0
```

A declaration gets a fake when each of its callables has a role of a contract with a reference form. The generator writes the fake from the forms. A callable without a role refuses with the gate `reference-role`, whose remedy names the contract roles that would cover it.

### The generated files

| File | Contents |
|---|---|
| `store_fake.go` | `type StoreFake struct`, safe for concurrent use, with `NewStoreFake() *StoreFake` |
| `store_fake_test.go` | `TestStoreFake`, which runs `conform.Run` with the suite of the declaration and the factory `NewStoreFake` |

The fake implements the declaration, so a downstream team passes it wherever the code takes the declaration. Its failures are the error values that the declaration names through its mixins, such as the sentinel of `not-found`.

### The suite checks the fake

The test of the fake applies every derived check and every authored check of the suite. A fake that fails a check fails the build of the generating repository, before any downstream team uses it. The suite record of the fake states the version of the spec and the checks that it passed, so a downstream team sees the laws that the fake passed.

The fake and the derived checks come from one spec, so the derived checks alone cannot show that the fake is wrong. The production implementation, the authored checks and the corpus check it from outside the spec:

- The consumer applies the same suite to its production implementation, so a check that the fake passes and the production implementation fails shows a difference between the two.
- The authored checks of the consumer apply to the fake.
- The corpus applies the suite of each fixture with a reference form to the fake, beside the hand-written subject and the controls of the fixture.

A fake is never the subject or a control of a corpus fixture. The subjects and controls of the corpus are written by hand, because they judge the probes and the fakes that come from the spec.

## Alternatives considered

### A fake from the reference of the generated sequences

The consumer's reference implementation of RFC-0020 becomes the fake.

**Why not:** a reference is the consumer's code, kept for tests of the producer. A fake is published for downstream teams and must be safe for concurrent use, documented and stable. The generator produces it from the spec, so its behaviour is the same for every consumer of a contract.

### Verification of a hand-written fake

dokimi verifies a fake that the consumer writes, by applying the suite to it.

**Why not:** that is already possible with `conform.Run` and needs no design. This RFC removes the hand-written code for the contracts that the spec can state.

## Drawbacks

- **Reference forms per contract.** Each contract that gets fakes needs a form in the spec and templates in each language.
- **Only behaviour that a law covers.** A fake agrees with the real implementation on what the suite checks. Behaviour outside the laws, such as the order of a listing that no law fixes, can differ.
- **A fake that agrees with the derived checks by construction.** The fake and the derived checks come from one spec. The production implementation, the authored checks and the subjects and controls of the corpus check the fake against behaviour that the spec did not produce.
- **One more generator.** Its template trees in each language and its corpus fixtures.

## Open questions

- Which contracts of the catalog get a reference form: the three of the summary, or every contract with a state?

## References

| What | Where |
|---|---|
| Software Engineering at Google, chapter 13, Test doubles | <https://abseil.io/resources/swe-book/html/ch13.html> |
| gocloud.dev, the conformance tests of blob drivers | <https://pkg.go.dev/gocloud.dev/blob/drivertest> |
| testing/fstest | <https://pkg.go.dev/testing/fstest> |
| Martin Fowler, ContractTest | <https://martinfowler.com/bliki/ContractTest.html> |
