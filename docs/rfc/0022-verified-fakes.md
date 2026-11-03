---
rfc: 0022
title: Verified fakes
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Draft
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

# RFC-0022: Verified fakes

## Summary

A new generator, `fake`, writes an in-memory implementation of a declaration under the directive `fake`. Each callable of the declaration belongs to a contract of the shape catalog with a reference form in the conformance spec: `persister`, `cache` or `appender`. The generated test of the fake applies the suite of the declaration to the fake, so the fake passes the same checks as every production implementation that the consumer applies the suite to. A downstream team uses the fake in its own tests in place of a container or a hand-written fake. The suite record of the fake states that it passed, at which version of the spec, so a consumer can see which laws the fake passed. A declaration with a callable outside every reference form does not get a fake. It gets a refusal with the name of that callable.

## Motivation

A downstream team tests its code against a dependency in three common ways, and each has a cost:

- A container of the real service costs seconds per test and a daemon in CI.
- A hand-written fake drifts from the real implementation, because no suite checks it against the same contract.
- A double returns what the test configures and states no behaviour.

The practice that checks a fake against its contract is a shared suite. Google requires that "a fake must have its own tests to ensure that it conforms to the API of its corresponding real implementation", and that the owner of the real implementation writes the fake. gocloud runs one conformance suite, `drivertest.RunConformanceTests`, against the in-memory blob driver and against the S3 driver. The standard library checks every `fs.FS`, `MapFS` included, with `fstest.TestFS`. No published study measures how far fakes drift from the implementations that they replace.

The Go corpus contains 141 hand-written `InMemory` subjects, one per fixture. Most of them implement a store, a cache, a log or a queue in a map or a slice. A downstream team writes the same kind of in-memory implementation for its own tests, and no suite checks it.

## Detailed design

### Reference forms

The spec gains `references.yaml`. A reference form contains the state of a contract and the effect of each role on that state. Where the claims of the contract establish their state through a writer without a role, the form also contains the effect of that writer. The claim `cache` establishes its state through its partner `put`.

```yaml
references:
  persister:
    contract: persister
    state: map
    roles:
      writer: store
      reader: load
    since: 1.1.0
  cache:
    contract: cache
    state: map
    roles:
      cache: load
      backing: load
    writers: store
    since: 1.1.0
  appender:
    contract: appender
    state: list
    roles:
      fn: append
    since: 1.1.0
```

A state is a `map` from keys to values or a `list` of values. Each effect acts on that state:

| Effect | State | What the callable does |
|---|---|---|
| `store` | map | Stores its written value under its key. Where its first result has the type of the written value, it returns that value |
| `load` | map | Returns the value under its key, or the zero value and the absent error where the map has none |
| `append` | list | Appends its input to the list, and returns the length of the list as the position of the new element |

The derivation finds the key and the written value of each callable by the rule that pairs a read with a write. A parameter of the same name and type in a load and a store is the key. Any other parameter of the load is the field of the written value with its name and type. The persister fixture of the corpus keys its map by the field `Key` of the value of `Put`, because `Get(ctx, key string)` selects that field. The cache fixture keys its map by the parameter `key`, which `Store`, `Lookup` and `Fetch` share.

A writer of a form is a callable without a role of the contract that the catalog classifies with the shape `writer` or `composite-writer`. The binding `key` of `composite-writer` is its key, and the binding `value` is its written value.

The spec package decodes the forms as it decodes the other files:

```go
// ReferenceForm is the reference form of a verified fake: the state of a
// contract, and the effect of each role and of each writer on that state.
type ReferenceForm struct {
	ID       ReferenceFormID
	Contract string
	State    State

	// Roles contains the effect of each role of the contract.
	Roles map[string]Effect

	// Writers is the effect of a writer without a role of the contract, and
	// empty for a form that takes no such writer.
	Writers Effect

	Since string
}

// State is the state of a reference form: a map or a list.
type State string

// Effect is what a callable of a reference form does to the state of the
// form: store, load or append.
type Effect string

// References returns every reference form in the order of references.json.
func References() []ReferenceForm

// ReferenceOf returns the reference form with the identifier id, and false
// when the spec has no such form.
func ReferenceOf(id ReferenceFormID) (ReferenceForm, bool)
```

A test of the package resolves the contract and each role of a form against the shape catalog of eidos. It checks that each role of the contract has an effect. A form of the state `map` has a `store` and a `load`, and a form of the state `list` has an `append`.

### The declarations that get a fake

A declaration opts in with the directive `fake`, as it opts in to the other generators through `stub` and `suite`. It also needs the directive `suite`, because the test of the fake applies that suite. Each callable of the declaration has a role of a contract with a reference form or is a writer of such a form, and the catalog gives it no other classification. The fake satisfies the laws of its forms and the claims that every callable gets: `smoke`, `cancel` and `zero-on-error`. A mixin such as `ttl` or `monotonic-reads` states a law that the forms do not cover, so its callable refuses.

Any other declaration does not get a fake. It gets the refusal of the gate `reference-role`, which the run reports as a finding of the generator. The finding contains the declaration, the callable, the classification outside every form, and the remedy of the gate:

```yaml
reference-role:
  passes: Each callable has a role of a contract with a reference form, or is a writer of such a form, and has no other classification.
  why: "{callable} has no role of a contract with a reference form"
  remedy: "give {callable} a role of a contract with a reference form, or leave {declaration} without the directive fake"
  since: 1.1.0
```

### The behaviour of a fake

A fake is safe for concurrent use. One mutex guards its state, and each call locks the mutex once. Each callable that takes a context returns the error of a done context first, with the zero value of each result. A load of an absent key returns the zero value of each result beside the absent error.

The absent error is a sentinel that the file of the fake declares, `ErrContractFakeNotFound` for the fake `ContractFake`. The directive `fake` takes the param `absent`, a name template of that sentinel, as the directive `stub` takes the name templates of its fields. Its default is `Err{{.Name}}FakeNotFound`.

### The generated files

The generator writes two files beside the suite of the declaration, in the directory of the directive `out`:

| File | Contents |
|---|---|
| `iface_fake.go` | `type ContractFake struct`, `func NewContractFake() *ContractFake`, and the sentinel `ErrContractFakeNotFound` |
| `iface_fake_test.go` | `TestContractFake`, a black-box test that calls `RunContractSuite(t, NewContractFake)` |

The names come from the source file `iface.go` and the declaration `Contract`, as the names of the files of the suite and the stub do. The fake implements the declaration, so a downstream team passes it wherever the code takes the declaration.

### The suite checks the fake

The test of the fake applies every derived check and every authored check of the suite. A fake that fails a check fails the build of the generating repository, before any downstream team uses it. The suite record of the fake states the version of the spec and the checks that it passed, so a downstream team sees the laws that the fake passed.

The fake and the derived checks come from one spec, so the derived checks alone cannot show that the fake is wrong. The production implementation, the authored checks and the corpus check it from outside the spec:

- The consumer applies the same suite to its production implementation, so a check that the fake passes and the production implementation fails shows a difference between the two.
- The authored checks of the consumer apply to the fake.
- The corpus applies the suite of each fixture with a reference form to the fake, beside the hand-written subject and the controls of the fixture.

### The corpus

The fixtures of `persister`, `cache` and `appender` state the directive `fake`. The expectation file of a fixture lists its fakes beside its subjects and its controls:

```yaml
fixture: iface/contract/persister
level: 2
subjects:
  - type: iface/contract/persister/persistertest.InMemory
controls: []
fakes:
  - iface/contract/persister/persistertest.ContractFake
```

The judge of the corpus reads the subject record of each fake as it reads the record of a control. It requires a pass of each check of the declaration, and it fails a fixture whose expectation lists a fake without a record.

A fake is never the subject or a control of a corpus fixture. The subjects and controls of the corpus are written by hand, because they judge the probes and the fakes that come from the spec.

### Other languages

The forms are data of the spec, so each language renders the same effects through a template tree of its own. A fake in TypeScript is a class over a `Map` or an array, with the same effects and the same absent error. This RFC proposes the Go template tree alone.

## Alternatives considered

### A fake from the reference of the generated sequences

The consumer's reference implementation of RFC-0020 becomes the fake.

**Why not:** a reference is the consumer's code, kept for tests of the producer. A fake is published for downstream teams and must be safe for concurrent use, documented and stable. The generator produces it from the spec, so its behaviour is the same for every consumer of a contract.

### Verification of a hand-written fake

dokimi verifies a fake that the consumer writes, by applying the suite to it.

**Why not:** that is already possible with `conform.Run` and needs no design. This RFC removes the hand-written code for the contracts that the spec can state.

### A fake for each declaration that the forms cover

The generator writes a fake for every declaration whose callables the forms cover, without a directive.

**Why not:** each generator of dokimi writes files only for a declaration under its directive, so the source of a declaration shows which files it gets. A fake also adds an exported type and a sentinel to the package of the kit.

### The sentinel of an absent key from the declaration

The fake returns the sentinel that the declaration names for an absent key.

**Why not:** the catalog has no param of such a sentinel for `persister` and `cache`, and the fixtures of the corpus declare their sentinels beside their subjects. A mixin with such a param, such as `delete-removes`, states a law that the forms do not cover.

## Drawbacks

- **Reference forms per contract.** Each contract that gets fakes needs a form in the spec and templates in each language.
- **Only behaviour that a law covers.** A fake agrees with the real implementation on what the suite checks. Behaviour outside the laws, such as the order of a listing that no law fixes, can differ.
- **A fake that agrees with the derived checks by construction.** The fake and the derived checks come from one spec. The production implementation, the authored checks and the subjects and controls of the corpus check the fake against behaviour that the spec did not produce.
- **Fakes for plain declarations alone.** A declaration of the three contracts with another classification does not get a fake, such as a declaration with a reader under the mixin `ttl`.
- **One more generator.** The generator needs a template tree in each language and a directive with one param. The corpus needs fixtures for it.

## Open questions

None.

## References

| What | Where |
|---|---|
| Software Engineering at Google, chapter 13, Test doubles | <https://abseil.io/resources/swe-book/html/ch13.html> |
| gocloud.dev, the conformance tests of blob drivers | <https://pkg.go.dev/gocloud.dev/blob/drivertest> |
| testing/fstest | <https://pkg.go.dev/testing/fstest> |
| Martin Fowler, ContractTest | <https://martinfowler.com/bliki/ContractTest.html> |
