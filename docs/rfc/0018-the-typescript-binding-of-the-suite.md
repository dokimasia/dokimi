---
rfc: 0018
title: The TypeScript binding of the suite generator
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

# RFC-0018: The TypeScript binding of the suite generator

## Summary

The suite generator gains a template tree for TypeScript, and `@dokimi/runtime` gains the modules `conform`, `spec` and `probe`. A TypeScript declaration under `//dokimi:suite` gets a suite of derived checks with planted defects, which Vitest runs and proves. The binding implements conformance level 1 of the corpus first: the claims of smoke, cancellation through an `AbortSignal` and the declared error of a miss. Each claim of level 1 that TypeScript cannot check refuses through a gate of the spec, such as the zero value beside an error, because a TypeScript call that fails throws and returns no value.

## Motivation

The TypeScript corpus declares `//dokimi:suite` on 132 declarations, with 81 mixins and 72 contract roles, and none of them gets a suite (`corpus/typescript`). The suite generator has a Go template tree alone (`generator/suite/templates/golang`), and `@dokimi/runtime` contains the runtime of the doubles alone: the recorder, the stub of a method, `answer`, the order tracker, the group and the fault strategies. The stub, builder and enum generators already write TypeScript, so a TypeScript consumer gets doubles and builders and no proof of a contract.

A second language is also the test of the claim of RFC-0001 that one derivation serves every language. The gates of the derivation read the roles and facts that the rules of each language assign. If TypeScript needs a gate that Go does not, or a claim with another meaning, the neutral facts have failed, and the TypeScript binding is where that shows. A difference in the facts is not such a failure: the error model of a TypeScript callable is `thrown`, so a failed call returns no value, and the gate `returns-beside-error` refuses the claims of a value beside an error.

assert-typescript implements assert-spec in TypeScript, with the property engine, the machines and the checkers of histories. Its overlay declares six assertions absent, the counts of allocations among them (`assert-spec/overlays/typescript.json`).

## Detailed design

### The generated modules

For `iface.ts` with a declaration under `//dokimi:suite`, the generator writes two modules beside the source:

| Module | Contents |
|---|---|
| `iface.suite.ts` | `export const readerSuite: Suite<Reader>` with the checks, the refusals and the unclaimed classifications, and the planted defects as classes |
| `iface.suite.test.ts` | The Vitest file that falsifies the checks of the consumer, and applies the suite to the subject of the directive `subject` where the declaration has one |

A check in TypeScript:

```ts
{
  id: readerChecks.read.cancel(),
  statement: "read reports an aborted signal",
  claim: claims.cancel,
  reason: reasons.ignoredCancel,
  law: laws.aCancelledContextIsReported,
  apply: async (h) => {
    await probe.underCancelled(h, (signal) => h.subject().read(readerKeys(h)[0], signal));
  },
  defect: (h) => new ReaderIgnoresCancel(h.fresh()),
},
```

### The runtime

`@dokimi/runtime` gains three modules. Each mirrors its Go counterpart in the conventions of TypeScript:

```ts
/** One check of a suite over subjects of the type S. */
export interface Check<S> {
  readonly id: CheckId<S>;
  readonly statement: string;
  readonly claim?: ClaimId;
  readonly reason?: ReasonId;
  readonly law?: LawId;
  readonly apply: (h: Handle<S>) => Promise<void>;
  readonly defect: (h: Handle<S>) => S;
}

/** The checks of one declaration, with what the generator declined. */
export interface Suite<S> {
  readonly language: "typescript";
  readonly module: string;
  readonly declaration: string;
  readonly checks: readonly Check<S>[];
  readonly refusals: readonly Refusal[];
  readonly unclaimed: readonly Unclaimed[];
}

/**
 * Registers a Vitest test for each check of suite. The test falsifies a
 * derived check on its planted defect over a fresh subject of newSubject,
 * then applies the check to a fresh subject. It writes the falsification
 * record and the record of the subject into the meta of the test of the
 * application.
 */
export function run<S>(suite: Suite<S>, newSubject: (h: Handle<S>) => S, options?: Options<S>): void;

/**
 * Registers a Vitest test for each check of the consumer in suite, applied
 * to its planted defect, and fails each check whose first failure does not
 * have its statement as the contract.
 */
export function falsify<S>(suite: Suite<S>): void;
```

`Handle<S>` implements the `Seat` of assert-typescript, so every assertion of a probe reports to it. It has the controlled clock, the seed, an `AbortSignal` that a probe aborts for the claim of cancellation, and the methods `fresh()`, `expect(reason)` and `engage()` of RFC-0014 and RFC-0015.

### Level 1

The claims and gates are those of RFC-0003, which also lists the construct of each language that passes each gate. The table names the fact of each gate in TypeScript:

| Claim | Gates | The fact in TypeScript |
|---|---|---|
| `smoke` | None | — |
| `cancel` | `cancellation`, `failure` | A parameter of the type `AbortSignal` has the role context, and every callable can fail |
| `zero-on-error` | `failure`, `value-beside-error`, `one-value`, `value-samples` | Refuses through `value-beside-error`, because the error model of every TypeScript callable is `thrown`. The overlay of TypeScript in the corpus lists each refusal |
| `miss` | `failure`, `one-value` | A declared error value of the module that the mixin `not-found` names in its param `sentinel` |

A claim whose gate fails refuses with the reason and the remedy of the gate, such as "declare a default of Value with the directive default". The overlay of TypeScript in the corpus lists each such refusal of a fixture with its gate.

### What TypeScript refuses

| Capability | Refusal | Reason |
|---|---|---|
| Allocation budgets | Each claim of a budget of allocations or bytes | The overlay of assert-spec declares `bench-max-allocs` and `bench-max-bytes` absent |

The fuzz targets of a TypeScript suite are those of RFC-0006, through `fuzzStoreSuite` and Jazzer.js.

## Alternatives considered

### Decorators as the directives of TypeScript

A consumer writes `@suite()` and `@mixin.ttl({...})` on a class or an interface.

**Why not:** a decorator applies to a class and its members, not to an interface or a type alias, and a decorator runs at load time in the code of the consumer. The comment directives that eidos reads in TypeScript work on every declaration and run nothing. The corpus already uses them on 132 declarations.

### A runtime of probes shared with Go

The TypeScript runtime calls the Go probes through WebAssembly.

**Why not:** a probe calls the subject, which is TypeScript code in the process of Vitest. Each call would cross into WebAssembly and back, and a probe would assert through a second assert library. RFC-0017 records the general choice against a shared implementation.

## Drawbacks

- **Every probe twice.** Each probe of level 1 exists in Go and in TypeScript, and each later probe needs both.
- **A template tree per generator.** The suite generator gains its largest template tree, beside the Go tree of about 700 lines.
- **Asynchronous probes.** Every probe and every planted defect in TypeScript is asynchronous, so a defect that breaks a law in a synchronous callable still returns a promise where the declaration does.
- **Refusals that Go does not have.** A TypeScript consumer sees refusals for zero values and allocation budgets that a Go consumer does not see.

## Open questions

- Does the falsification of a TypeScript suite run in the same Vitest project as the application, or in a project of its own?

## References

| What | Where |
|---|---|
| Vitest, the meta of a test and its reporters | <https://vitest.dev/api/advanced/metadata> |
| assert-spec, the overlay of TypeScript | <https://github.com/dokimasia/assert-spec/blob/main/overlays/typescript.json> |
| assert-typescript | <https://github.com/dokimasia/assert-typescript> |
