---
rfc: 0032
title: Derived claims that a declaration narrows
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Draft
created: 2026-11-02
updated: 2026-11-02
discussion: https://github.com/dokimasia/dokimi/issues/51
supersedes: RFC-0003 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0032: Derived claims that a declaration narrows

## Summary

dokimi derives the claims of a callable from its signature. The signature does not show whether a call reads a clock, changes its subject, or refuses part of the domain of a parameter with a panic. The author of the declaration states each of these facts with a directive:

- The gate `stable-subject` refuses the claims `deterministic` and `predicate` for a callable with the mixin `time-aware` or `side-effect`.
- dokimi reads the mixin `sample` of the shape catalog. Its builder returns the inputs of every derived call of the callable, so a check calls the callable inside the precondition of its parameters.
- The directive `refuse` withdraws a derived claim of a callable for every subject of its suite, with a reason.

## Motivation

dokimi v0.7.0 derives claims that the interfaces of thesm-os/core do not make. A run of `make generate-go` at core de41002, then `go test ./...`, fails 15 packages:

| Callable | Claim | Why the claim is false or unreachable |
|---|---|---|
| `clock.Clock.Now`, `Time` and `Elapsed` | `deterministic` | Each call reads the clock |
| `rand.Rand.Uint64` | `deterministic` | Each call returns the next value of the source |
| `crypto.AEAD.Seal` | `deterministic` | Seal panics for a nonce that is not `NonceSize()` bytes. The AEAD of `aesgcm.NewRandomNonce` also generates a new nonce for each call |
| `clock.Timer.Stop` | `predicate` | Stop stops the timer and reports whether it was active, so the repeat returns false |
| `crypto.AEAD.Open` | `smoke` | Open panics for a nonce that is not `NonceSize()` bytes, as `cipher.AEAD` documents |
| `clock.Clock.Update` | `zero-on-error` | Update fails only for an instant beyond the max offset of the clock, and no derived input is such an instant |

The derivation reads the signature alone:

- The shape `computation` matches each callable that returns one value and takes no context, so `Len() int` and `Now() Instant` get the same claims.
- The shape `predicate` matches each method without parameters that returns a truth value. A command such as `Stop` has the signature of a query such as `Empty`.
- A derived sample of `[]byte` has any length. Seal and Open take the parameters of `cipher.AEAD` of the standard library, so core cannot put a directive on them.
- The gate `context-or-input` passes for each callable with an input, because a fresh store fails a read of any key. A fresh clock does not fail an update for every input.

An excusal beside each subject passes the run. Each implementation of the interface then repeats the same excusals, because the claims are false for the interface, or unreachable for each of its implementations.

## Detailed design

### Invariants

| Invariant | Enforced by |
|---|---|
| A callable with the mixin `time-aware` or `side-effect` has neither the claim `deterministic` nor the claim `predicate` | The gate `stable-subject` |
| Every derived call of a callable with the mixin `sample` takes its inputs from the builder of the mixin | The suite generator |
| A claim that a directive `refuse` names is a refusal for every subject of the suite, with the reason of the directive | The derivation |
| The claim of a directive `refuse` is a claim that the derivation licenses for its callable | The code `UnknownRefusal` of the suite generator |

### Components

| Component | Change |
|---|---|
| The conformance spec | The gate `stable-subject`, which the claims `deterministic` and `predicate` list |
| The derivation | Refuses each claim that a directive `refuse` of the callable names |
| The suite generator | Reads the builder of the mixin `sample`, and the directive `refuse`, with the codes `BadBuilder` and `UnknownRefusal` |
| The template of the suites | Writes the call of the builder in place of the samples and the generated values of the parameters |
| The Go guide | A row for the builder and a row for the directive `refuse` |

### The gate stable-subject

The spec adds the gate, and the claims `deterministic` and `predicate` list it:

```yaml
stable-subject:
  passes: The callable has no mixin time-aware and no mixin side-effect, so a repeat on an unchanged subject reads the state of the first call.
  why: "{callable} reads the clock or changes its subject, so a repeat can return another value"
  remedy: "remove the mixin time-aware or side-effect where a repeat of {callable} returns the same value"
  since: 1.1.0
```

core states the mixins on the callables of the motivation:

```go
type Clock interface {
	// Now returns the current instant of the clock.
	//
	//dokimi:mixin time-aware
	Now() Instant
}

type Timer interface {
	// Stop stops the timer, and reports whether the timer was active.
	//
	//dokimi:mixin side-effect
	Stop() bool
}
```

The claims of both mixins require a context, so `Now`, `Uint64` and `Stop` refuse them with the gate `cancellation`.

### The builder of the inputs

In the shape catalog, the mixin `sample` means that a check builds its inputs through the callable of its param `builder`, and calls the subject with them. dokimi reads it on a callable of a declaration under `suite`:

- The builder is a function that takes the subject or nothing. It returns one value of the type of each parameter of the callable, in their order, and none for a context.
- The sample binding of each derived claim of the callable takes the values of the builder. So does the call generator of the callable, so its benchmarks, its budgets and its generated bindings call it with them too.
- The builder takes precedence over the sample of a parameter.
- A builder of another signature stops the generation with the code `BadBuilder` at the directive.
- The file of the suite calls the builder, so the file must be able to import it. A builder that it cannot import is `Unimportable`, as for any other reference of a suite.

A method that an interface promotes from an embed takes the mixin on a redeclaration with the same signature. Go accepts an explicit method beside an embed that has it, when the two signatures are identical:

```go
type AEAD interface {
	cipher.AEAD

	// Seal encrypts and authenticates plaintext, as [cipher.AEAD] documents.
	//
	//dokimi:mixin sample builder=SealInputs
	Seal(dst, nonce, plaintext, additionalData []byte) []byte

	// Open authenticates and decrypts ciphertext, as [cipher.AEAD] documents.
	//
	//dokimi:mixin sample builder=OpenInputs
	Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error)
}

// SealInputs returns the inputs of a call of Seal on a, with a nonce of
// NonceSize bytes.
func SealInputs(a AEAD) (dst, nonce, plaintext, additionalData []byte) {
	return nil, make([]byte, a.NonceSize()), []byte("plaintext"), nil
}
```

### The directive refuse

The directive `refuse` on a callable takes one derived claim of the callable and the reason why the declaration does not make it:

```go
type Clock interface {
	// Update merges the instant of a peer into the clock. It returns
	// ErrOffset for an instant beyond the max offset of the clock.
	//
	//dokimi:refuse zero-on-error reason="Update fails only for an instant beyond the max offset of the clock"
	Update(observed Instant) (Instant, error)
}
```

| Part | Rule |
|---|---|
| Carrier | A method of a declaration under `suite`, or a function under `suite` |
| Positional param | The name of the claim, such as `zero-on-error` |
| Param `reason` | Required. The reason why the declaration does not make the claim |
| Repeats | Once for each claim |
| Result | The derivation refuses the claim. The suite record lists the refusal with the reason of the directive, as it lists the refusal of a gate |
| Failure | A claim that the derivation does not license for the callable stops the generation with `UnknownRefusal`, and the message lists the claims of the callable |

A refusal of `deterministic` also covers Seal of an AEAD that generates its own nonce. The catalog has no mixin for a call that reads a random source.

### Verification

- **The derivation.** The specs of the gate cover a callable with each mixin and a callable with neither. The specs of a refusal cover a refused claim, and a claim that another claim entails.
- **The suite generator.** The specs cover a builder with the subject, a builder without it, a builder of another signature, an unimportable builder and a redeclared promoted method. They also cover a refusal of a licensed claim and a refusal of an unknown claim. The golden files pin the call of a builder.
- **The corpus.** One fixture of each case of the motivation: a clock with a time-aware read, a source with a side effect, a timer whose stop is a command, an AEAD with the builders of its nonce, and an update with a refused claim. Each fixture passes its suite and the falsification of each check that it keeps.

## Alternatives considered

### Leave deterministic out for every method of an interface that reads state

Issue 51 proposed it first. The derivation reads the signature, and `Len() int` and `Now() Instant` have one form. Leaving the claim out for every method without parameters drops it for each read that repeats alike.

### Leave predicate out for a method that changes its receiver

Bertrand Meyer's command-query separation tells `Stop() bool` from `Empty() bool` by what each call does, and the signature does not show it. Only a declaration states that `Stop` is a command.

### The sample of a parameter

The sample of a parameter fixes one value of it in every generated call. Seal and Open come from `cipher.AEAD`, so core cannot put a directive on their parameters. One nonce also cannot follow `NonceSize`, which is 12 for `aesgcm.New` and 0 for `aesgcm.NewRandomNonce`.

### A gate on a failing sample

Issue 51 proposed to refuse `zero-on-error` unless a sample makes the call fail. dokimi derives the claims while it generates, so it cannot run a call to find out. A gate on a declared failing input would also refuse the claim for each reader, whose read of any key fails on a fresh subject and engages the check today.

### An excusal on the declaration

The directive `excuse` could go on the interface and cover every subject. An excusal keeps the check and records that a subject does not pass it. A refusal fits the callables of the motivation better, because none of them makes its claim for any subject.

## Drawbacks

- **Builders in the API.** A builder is a function that the file of the suite imports, so core exports `SealInputs` and `OpenInputs` from `crypto`.
- **A repeated signature.** A redeclared promoted method repeats the signature of the embed. The interface does not compile when the embed changes the signature, until the redeclaration changes with it.
- **A refusal for every subject.** A refusal hides the claim also for a subject that would pass it.
- **More refusals in the record.** The claims of `time-aware` and `side-effect` require a context, so each mixin on a callable without one adds the refusals of the gate `cancellation`.
- **One more gate in the spec.** Each implementation of the derivation adds the gate `stable-subject`.

## Open questions

1. Can the builder resolve in the package of the suite, such as `cryptotest`, so that it is not part of the API of the declaration? eidos resolves a callable in the scope of the subject of the directive, and the declaration does not import the package of the suite.

## References

| What | Where |
|---|---|
| Issue 51, four derived claims that the interfaces of core do not make | <https://github.com/dokimasia/dokimi/issues/51> |
| The mixins `sample`, `time-aware` and `side-effect` of the shape catalog | <https://github.com/dokimasia/eidos/tree/main/eidos-plugin-shape/spec/mixins> |
| The interface types of Go, an explicit method beside an embed | <https://go.dev/ref/spec#Interface_types> |
| `cipher.AEAD` and the panic for a nonce of another length | <https://pkg.go.dev/crypto/cipher#AEAD> |
| Bertrand Meyer, Object-Oriented Software Construction, command-query separation | <https://en.wikipedia.org/wiki/Command%E2%80%93query_separation> |
