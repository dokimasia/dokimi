---
rfc: 0018
title: The TypeScript binding of the suite generator
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-28
updated: 2026-11-01
discussion: none
supersedes: none
superseded-by: RFC-0050 in part
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0018: The TypeScript binding of the suite generator

## Summary

The suite generator gains a template tree for TypeScript, and the TypeScript plan runs it and writes the lock `dokimi.typescript.lock`. `@dokimi/runtime` gains the modules `spec`, `conform` and `probe`. A TypeScript declaration under `//dokimi:suite` gets a suite of derived checks with planted defects, which Vitest runs and proves, and each application writes its suite records into the meta of a Vitest test. The binding covers conformance level 1 of the corpus, whose fixtures bind the claims `smoke`, `cancel`, `deterministic`, `idempotent`, `repeat-survives` and `zero-on-error`. A TypeScript callable raises its errors, so `zero-on-error` refuses through the gate `returns-beside-error`, and a probe counts a throw of an error of the engine as a panic.

## Motivation

The TypeScript corpus declares `//dokimi:suite` on 132 declarations, with 77 directives `mixin` and 72 directives `contract`, and none of them gets a suite (`corpus/typescript`). The suite generator has a Go template tree alone (`generator/suite/templates/golang`). `@dokimi/runtime` contains the runtime of the doubles alone: the recorder, the stub of a method, `answer`, the order tracker, the group, the behaviour suites and the fault strategies. The stub, builder, enum and sentinel generators already write TypeScript, so a TypeScript consumer gets doubles and builders and no proof of a contract.

A second language also tests the claim of RFC-0001 that one derivation serves every language. The gates of the derivation read the roles and facts that the rules of each language assign. If TypeScript needs a gate that Go does not, or a claim with another meaning, the neutral facts have failed, and the TypeScript binding is where that shows. A difference in the facts is not such a failure. The rules of eidos give every TypeScript callable the error model `raised`, so a failed call returns no value, and the gate `returns-beside-error` refuses the claims of a value beside an error.

assert-typescript implements assert-spec in TypeScript, with the property engine, the machines and the checkers of histories. Its overlay declares seven assertions blocked, the counts of allocations and of live bytes among them (`assert-spec/overlays/typescript.json`).

## Detailed design

### The generated modules

For `iface.ts` with a declaration under `//dokimi:suite`, the generator writes two modules beside the source. For a type under the directive `subject` in `inmemory.ts`, it writes one module beside the subject, as the Go binding writes `inmemory_suite_test.go`:

| Module | Contents |
|---|---|
| `iface.suite.ts` | The names of the callables and of the checks, `readerSuite()`, the planted defects, and `runReaderSuite` and `runReaderSuiteFromHandle` |
| `iface.suite.test.ts` | The Vitest file that falsifies the checks of the consumer through `falsify` |
| `inmemory.suite.test.ts` | The Vitest file that applies the suite to the subject, with the options of the directive `suite` such as `isolate` |

The names follow the contract of a suite in TypeScript. `readerChecks.read` is a `Callable<Reader>` with a function for each check, `readerSuite()` returns a `Suite<Reader>`, and the types `ReaderHandle`, `ReaderCheck` and `ReaderOptions` name the parts of one declaration. A planted defect is a class that implements the declaration over a fresh subject and breaks one law. Each identifier of the spec in a derived check is a string literal, which the union type of its field admits. The derived check of `read` that `cancel` licenses:

```ts
{
  id: readerChecks.read.cancel(),
  statement: "read returns an error for a cancelled context",
  claim: "cancel",
  probe: "under-cancelled",
  reason: "ignored-cancel",
  law: "cancelled-context-is-reported",
  machineries: ["nothing-beyond-a-caller"],
  apply: async (h) => {
    await underCancelled(h, (signal) => h.subject().read("Key-a", signal));
  },
  defect: (h) => new ReaderReadIgnoresCancel(h.fresh()),
},
```

### The package of a declaration

eidos gives each TypeScript module the module of its nearest `package.json`. The package of a declaration is the name of that `package.json`, a slash, and the directory of the module relative to it, such as `@dokimi/corpus-typescript/iface/mixin/ttl`. The lock and the suite record contain this package, as they contain the import path of a Go package. The gate of RFC-0017 removes the name of the `package.json` of the corpus, as it removes the module path of the Go corpus.

### The lock

The TypeScript plan runs the suite generator and writes the lock of its suites into `dokimi.typescript.lock`, which the composition names after the plan as it names `dokimi.go.lock`. Its lines have the five fields of the lock of Go, with the keys in the key form of the spec, so `Reader.read` and `Reader.Read` have one key.

### The runtime

Each file of the runtime is a module of its own, as each module of the doubles is, such as `@dokimi/runtime/stub/recorder`. No module re-exports another. `@dokimi/runtime` gains three directories of modules:

| Modules | Contents |
|---|---|
| `@dokimi/runtime/spec/<section>` | The union type of the identifiers of one section of the spec, such as `ClaimId` in `spec/claims`. `spec/probes` also maps each probe to its reasons and to the condition under which it does not engage |
| `@dokimi/runtime/spec/key` | `key`, which returns the key form of an identity |
| `@dokimi/runtime/spec/version` | The version of the copy of the spec |
| `@dokimi/runtime/conform/<part>` | The types of a suite in `conform/suite`, `Handle` in `conform/handle`, `run` in `conform/run` and `falsify` in `conform/falsify` |
| `@dokimi/runtime/probe/<probe>` | The probe of a claim of level 1, in the module of its identifier, such as `underCancelled` in `probe/under-cancelled`, and the types of a call in `probe/call` |

The repository of the runtime implements the definition `dokimasia/dokimi` with a copy of the spec under `spec/`, which `spec sync` writes and `spec check` verifies. The generator `tools/spec-types.ts` writes the modules of the sections from the copy, as the copy of the Go runtime writes `names.gen.go`. A union type of string literals checks an identifier at compile time as a Go constant does, so a generated check writes the identifier of the spec as it is.

The modules of `conform` mirror the package `conform` of Go for the machinery `nothing-beyond-a-caller`. Each other machinery adds its part through the RFC that designs it, and RFC-0008 and RFC-0009 already state the TypeScript names of their parts:

```ts
/** The identity of a check of a suite over subjects of the type S, such as "Reader.read/cancel". */
export type Id<S> = string & { readonly [subjectOf]?: S };

/** One check of a suite over subjects of the type S. */
export interface Check<S> {
  readonly id: Id<S>;
  readonly statement: string;
  readonly claim?: ClaimId;
  readonly probe?: ProbeId;
  readonly reason?: ReasonId;
  readonly law?: LawId;
  readonly machineries?: readonly MachineryId[];
  readonly shallow?: boolean;
  readonly entails?: readonly Id<S>[];
  readonly apply: (h: Handle<S>) => Promise<void>;
  readonly defect?: (h: Handle<S>) => S;
}

/** A licensed claim whose gate fails for a callable. */
export interface Refusal {
  readonly claim: ClaimId;
  readonly callable: string;
  readonly gate: GateId;
  readonly why: string;
  readonly remedy: string;
}

/** A classification of a callable without a check. */
export interface Unclaimed {
  readonly classification: string;
  readonly callable: string;
  readonly stance: StanceId;
  readonly note?: string;
}

/** The checks of one declaration, with what the generator declined. */
export interface Suite<S> {
  readonly declaration: string;
  readonly package: string;
  readonly generator: { readonly name: string; readonly version: string };
  readonly checks: readonly Check<S>[];
  readonly refusals: readonly Refusal[];
  readonly unclaimed: readonly Unclaimed[];
}

/** The options of an application of a suite. */
export interface Options<S> {
  /** The subject of the records, as the generator writes it from the directive subject. */
  readonly subject?: string;
  readonly checks?: readonly Check<S>[];
  readonly excuse?: Readonly<Record<string, string>>;
  readonly isolate?: boolean;
}

/**
 * Registers a Vitest test for each check of suite and of the option checks.
 * The test of a derived check falsifies the check on its planted defect over
 * a fresh subject of newSubject, and then applies the check to the subject.
 * The last test of the application writes its falsification record and its
 * subject record into its meta.
 */
export function run<S>(suite: Suite<S>, newSubject: (h: Handle<S>) => S, options?: Options<S>): void;

/**
 * Registers a Vitest test for each check of the consumer in suite, applied to
 * its planted defect, and fails each check whose first failure does not have
 * its statement as the contract. A derived check, whose planted defect wraps a
 * fresh subject of an application, has its falsification in run. The last
 * test writes the falsification record into its meta.
 */
export function falsify<S>(suite: Suite<S>): void;
```

`Handle<S>` implements the `Seat` of assert-typescript, with `observeCall`, so every assertion of a probe reports to it and an authored check engages through its statement. It has the methods `subject()`, `new()`, `fresh()`, `clock()`, `advance(ms)`, `statement()`, `expect(reason)` and `engage()` of RFC-0003, RFC-0014 and RFC-0015, and the property `signal`, which is the signal of the test of the check.

Go writes a record on the test of the application. A Vitest report writes the meta of a test and not that of a `describe` block, so `run` registers a block named after the application, such as `runReaderSuite`, with a test for each check and a last test, `suite records`. The block runs its tests in order, and it is neither concurrent nor shuffled. The test of a check has the identity of the check without the declaration, such as `read/cancel`. The last test reads the state of each other test of the block, because the fixture of assert-typescript fails a test for a failure of `soft` after its body ends. It then writes the records into its meta. `falsify` registers the block `readerSuite` in the same way. `writeRecord` of assert-typescript appends each record to the array under the key `dokimi.suite`. The subject of a record is the option `subject`. A TypeScript value has no import path at run time, and `iface/composite/causal-chain` has two subjects named `Log`, so the generated application of a subject passes the subject that the directive `subject` declares, such as `@dokimi/corpus-typescript/iface/mixin/ttl/ttltest.InMemory`. An application without the option has no subject. The runtime names itself `@dokimi/runtime` with the version of its `package.json`, and the spec of a record is the version of its copy of the spec.

### Level 1

The claims of the fixtures of level 1 have these probes and gates in TypeScript:

| Claim | Probe | Gates | The fact in TypeScript |
|---|---|---|---|
| `smoke` | `survives` | `returns-unaided`, `returns-before-cancel` | Every callable without the mixins `waits` and `runs-until-cancelled` |
| `cancel` | `underCancelled` | `cancellation`, `failure`, `returns-before-cancel` | A parameter of the type `AbortSignal` has the role context, and every callable can fail |
| `deterministic` | `repeatAgrees` | `one-value`, `value-samples`, `returns-unaided`, `returns-before-cancel` | As in Go |
| `idempotent` | `repeatIsQuiet` | `cancellation`, `failure` | As `cancel` |
| `repeat-survives` | `survivesRepeat` | `no-value`, `returns-unaided`, `returns-before-cancel` | As in Go |
| `zero-on-error` | — | `failure`, `returns-beside-error`, `one-value`, `value-samples`, `returns-unaided`, `context-or-input`, `returns-before-cancel` | Refuses through `returns-beside-error`, because the error model of every TypeScript callable is `raised` |

A claim whose gate fails refuses with the reason and the remedy of the gate, such as "author a sample of the result type with the directive sample". The overlay of TypeScript in the corpus declares each refusal of a fixture that Go does not refuse, with its gate.

A probe takes the handle and a call. The call takes the signal that the probe passes, and it returns, throws or settles its promise:

```ts
/** A call of the callable whose value a probe reads. */
export type Call<V> = (signal: AbortSignal) => V | Promise<V>;

/** A call of the callable whose value a probe does not read. */
export type Act = (signal: AbortSignal) => unknown;
```

| Probe | Signature | Procedure |
|---|---|---|
| `survives` | `(h, call: Act)` | Calls call, and fails when it throws an error of the engine |
| `underCancelled` | `(h, call: Act)` | Calls call with an aborted signal, and fails when it settles without an error |
| `repeatAgrees` | `<V>(h, call: Call<V>)` | Calls call twice, and fails when the two values differ. A first call that throws ends the probe |
| `survivesRepeat` | `(h, call: Act)` | Calls call twice, and fails when the second call throws an error of the engine |
| `repeatIsQuiet` | `(h, call: Act)` | Calls call twice, and fails when the second call throws after a first call that settled |

A TypeScript callable may throw on any call, so a throw alone is the failure of the callable and not a panic. The reason `panicked` reads "The call panicked or threw where it had to return". In TypeScript a panic is a throw of an error that the engine raises for a fault of the program: a `TypeError`, a `ReferenceError` or a `RangeError`. The planted defect of `smoke` throws a `TypeError`. A probe other than `survives` and `survivesRepeat` throws such an error again, so the check fails outside its assertions, as a panic ends a probe of Go. A consumer whose callable throws a `RangeError` for an input that the derivation samples excuses the check with the reason.

### What TypeScript refuses

| Capability | Behaviour | Reason |
|---|---|---|
| The ceilings `allocs` and `bytes` of a budget | The generator reports each at the directive as a warning and writes no ceiling. The overlay of TypeScript in the corpus declares each allocation ceiling of a fixture as absent, with the stance `blocked` | The overlay of assert-spec declares `max-allocs`, `bench-max-allocs` and `bench-max-bytes` blocked, because V8 counts no allocations of a call |

A ceiling is not a claim, so it has no refusal in the record. The gate of RFC-0017 accepts an entry with the actual outcome `absent` and the stance `blocked` as the declaration of a missing key.

The ceilings `p99` and `mean` of a budget run in the module `<file>.suite.bench.ts` of RFC-0006, and the fuzz targets of a TypeScript suite are those of RFC-0006, through `fuzzReaderSuite` and Jazzer.js.

### The report and the gate

`dokimi report` reads the JSON report of Vitest beside the events of `go test -json`. The records of a Vitest report are the elements of `meta["dokimi.suite"]` of each test in `testResults[].assertionResults`. The command tells the two forms apart by the first JSON value of a file: a Vitest report is one object with the key `testResults`, and `go test -json` writes one event per line.

`dokimi conform --lang typescript` runs `dokimi run --check` on `corpus/typescript`, then `pnpm exec vitest run --reporter=json`. It reads `dokimi.typescript.lock`, and takes the name of the corpus from its `package.json`.

### The corpus of TypeScript

The falsification of a derived check needs a fresh subject of an application, so each fixture of level 1 that has a subject in Go has one in TypeScript at the same path, such as `iface/lang/embedded/embeddedtest/inmemory.ts`. The TypeScript corpus lacks 11 fixtures of level 1 that the Go corpus has. It gains a port of each fixture whose axis TypeScript has, and its overlay blocks each fixture whose axis is a construct of Go:

| Fixture | TypeScript | Reason |
|---|---|---|
| `iface/lang/dynamic` | Port | A parameter of the type `unknown` |
| `iface/lang/nested` | Port | Fields of function types and of interface types |
| `iface/lang/recursive` | Port | A recursive type |
| `iface/lang/sampledparam` | Port | The directive `sample` on a parameter |
| `iface/lang/suitedhandle` | Port | A declaration that a method of a double returns |
| `iface/lang/crosspackage` | Blocked | The internal packages of Go |
| `iface/lang/kinds` | Blocked | The kinds of `reflect` and the channels of Go |
| `iface/lang/namedreturns` | Blocked | The named results of Go |
| `iface/lang/selfnamed` | Blocked | The field of an embedded interface |
| `iface/lang/opaque` | Blocked | A value that starts from the zero value of its type, which TypeScript does not have |
| `iface/lang/resolved` | Blocked | As `opaque` |

## Alternatives considered

### Decorators as the directives of TypeScript

A consumer writes `@suite()` and `@mixin.ttl({...})` on a class or an interface.

**Why not:** a decorator applies to a class and its members, not to an interface or a type alias, and a decorator runs at load time in the code of the consumer. The comment directives that eidos reads in TypeScript work on every declaration and run nothing. The corpus already uses them on 132 declarations.

### A runtime of probes shared with Go

The TypeScript runtime calls the Go probes through WebAssembly.

**Why not:** a probe calls the subject, which is TypeScript code in the process of Vitest. Each call would cross into WebAssembly and back, and a probe would assert through a second assert library. RFC-0017 records the general choice against a shared implementation.

### The subject from the constructor of the value

`run` writes the name of the constructor of the first subject, as conform of Go writes the type of the value.

**Why not:** a constructor has a name and no module at run time, and two subjects of one fixture can share a name, such as `archive.Log` and `replica.Log`. The generator knows the module of the subject from the directive, so the generated application passes it.

### Smoke as a refusal in TypeScript

A gate refuses `smoke` in every language whose callables raise their errors, because a throw is the failure of the callable.

**Why not:** `smoke` is the claim of every callable, and a refusal would leave TypeScript without its most common check. The errors that the engine raises for a fault of the program separate a panic from a failure.

### The falsification in a Vitest project of its own

The falsification of the checks runs in a second project of the Vitest configuration.

**Why not:** Go runs the falsification in the test binary of the package, and a second project would transform every module of the corpus twice.

## Drawbacks

- **Every probe twice.** Each of the five probes of level 1 exists in Go and in TypeScript, and each later probe needs both.
- **A template tree per generator.** The suite generator gains its largest template tree, beside the Go tree of 1,598 lines in three templates.
- **Asynchronous probes.** Every probe and every planted defect in TypeScript awaits its calls, so a defect of a synchronous callable still returns its value through a promise where the declaration does.
- **Refusals that Go does not have.** A TypeScript consumer sees refusals of `zero-on-error` and warnings for allocation budgets that a Go consumer does not see.
- **A panic by the class of the error.** A callable that throws a `TypeError`, a `ReferenceError` or a `RangeError` by design fails `smoke` for an input that the derivation samples, and the consumer excuses the check.
- **The subject in two places.** The directive `subject` and the option `subject` of the generated application both state the subject of an application.

## Open questions

None.

## References

| What | Where |
|---|---|
| Vitest, the meta of a test and its reporters | <https://vitest.dev/api/advanced/metadata> |
| assert-spec, the overlay of TypeScript | <https://github.com/dokimasia/assert-spec/blob/main/overlays/typescript.json> |
| assert-typescript | <https://github.com/dokimasia/assert-typescript> |
