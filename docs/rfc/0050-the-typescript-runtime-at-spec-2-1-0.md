---
rfc: 0050
title: The TypeScript runtime at spec 2.1.0
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-11
updated: 2026-11-11
discussion: none
supersedes: RFC-0004, RFC-0018, RFC-0019 and RFC-0023 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0050: The TypeScript runtime at spec 2.1.0

## Summary

`@dokimi/runtime` implements the runtime contract of spec 2.1.0 in TypeScript. It has a module for each part of the seven Go modules of the runtime: the runtime model of the spec, the suites and their handle, the probes, the doubles and their faults, the breakable medium, the simulated host and the simulation. The suite, fake, stub and sim generators each write a TypeScript tree, and `dokimi conform --lang typescript` judges the TypeScript corpus up to level 4. Each contract of the Go runtime keeps its meaning in TypeScript. A mechanism changes where TypeScript has another one:

- An `AbortSignal` replaces a context, and a call throws the error that a Go call returns.
- The event loop and the task scheduler of `@dokimi/assert` run the code of a component in place of goroutines. A wait for an idle event loop replaces `synctest.Wait`, and the fake timers of Vitest replace the clock of a bubble.
- A class or an injection token identifies a value at run time where Go reads its type.
- The meta of a Vitest test contains the records that Go writes as attributes of `go test -json`.

The copy of the spec comes from dokimasia/dokimi-spec, which lists the TypeScript runtime as an implementation. Allocation ceilings, parallel benchmarks, the race detector and the watchdog of a bubble have no TypeScript form.

## Motivation

### TypeScript declarations do not get a suite

- The TypeScript corpus has 152 of the 193 fixtures of the expected outcomes, and 132 declarations under `//dokimi:suite`. dokimi does not write a suite for any of them, because the suite, fake and sim generators have a Go tree alone, such as `generator/suite/templates/golang`. The stub, builder, enum and sentinel generators write TypeScript.
- No fixture of the TypeScript corpus has a subject, although 132 of its fixtures have one in Go. The falsification of a derived check needs a fresh subject.
- `@dokimi/runtime` contains the doubles and their fault strategies. Its base of `conform` and its modules of the spec are unreleased, and their copy is spec 1.0.0 of the definition `dokimasia/dokimi`. dokimi has not published a definition since the runtime contract moved to dokimasia/dokimi-spec, so `spec sync` has no source for that copy.
- The doubles of TypeScript differ from the doubles of Go. `behaviour` runs 9 checks for a method where Go runs 13, and a double still has `newDelayed` and `newFailing`.
- `dokimi conform` judges the corpus of Go alone (`cmd/dokimi/conform.go:34`).

dokimi-go implements spec 2.1.0 in seven modules, with 20,776 lines of source and 37,397 lines of specs. Its suites pass levels 1 to 4 of the corpus, the five fixtures of the simulation included. Its modules require assert-go v0.20.0, which implements definition 16.2.0 of assert-spec. `@dokimi/assert` implements definition 13.6.0, and npm has no release of it.

A consumer runs a TypeScript suite under Vitest in Node, so the suite calls a runtime in TypeScript. The derivation of the checks is in dokimi and serves both languages.

### The execution model of TypeScript

- **One thread.** One thread runs the JavaScript of a test file, and another function runs only where the running one awaits. Node drains the next tick queue once the current operation on the stack has completed, and the microtask queue right after it, before the event loop continues. A callback of `setImmediate` runs in the check phase, after the poll phase. So a callback of a real `setImmediate` runs after every promise job that the code before it queued.
- **Fake timers.** `vi.useFakeTimers` of Vitest fakes `setTimeout`, `setInterval`, their clears, `setImmediate`, `clearImmediate` and `Date` by default. It fakes `nextTick` and `queueMicrotask` only when a test lists them in `toFake`. Vitest 4.1.12 bundles @sinonjs/fake-timers 15.0.0, whose install on the global object also replaces the methods of the module `timers` of Node.
- **No hook into an await.** The task scheduler of assert-typescript releases a task until the task awaits `yield` or settles, because a library cannot intercept a native await (`src/stateful/scheduler.ts:43`). A task must not await a promise that only another task settles (`scheduler.ts:48`).
- **No type at run time.** A value has a class or none. An interface and a type alias do not exist at run time.
- **Asynchronous I/O.** The file system of Node returns promises, and a constructor cannot await one.
- **No count of allocations.** V8 counts no allocations of a call. The overlay of TypeScript in assert-spec declares `max-allocs`, `bench-max-allocs` and `bench-max-bytes` blocked.

## Detailed design

### Components

| Component | Where | Responsibility |
|---|---|---|
| The runtime contract | dokimasia/dokimi-spec | Lists the TypeScript runtime as an implementation, so `spec try` runs its tests |
| The copy and the runtime model | `@dokimi/runtime/spec/*` | The entries of the spec and their lookups, the key form, the version and the types of the suite record |
| The suites | `@dokimi/runtime/conform/*` | Suites, checks, the handle, the applications, falsification, benchmarks, fuzz targets and the records |
| The probes | `@dokimi/runtime/probe/*` | A module for each of 121 probes, and the types of their inputs |
| The doubles | `@dokimi/runtime/stub/*` and `@dokimi/runtime/fault` | The doubles and their fault strategies, with the behaviour suites of the doubles of Go |
| The medium | `@dokimi/runtime/medium/*` | A file system with an explicit sync, over Node and in memory |
| The host | `@dokimi/runtime/host/*` | Hosts, incarnations and kills, the clocks, the network, the dependencies and their injector |
| The simulation | `@dokimi/runtime/sim/*` | Worlds of sans-IO machines and of components, their faults, their runs and their proofs |
| The template trees | dokimi, `generator/{suite,fake,stub,sim}/templates/typescript` | The TypeScript code of each generator |
| The report and the gate | dokimi, `dokimi report` and `dokimi conform` | The records of a Vitest report, and the judgement of the TypeScript corpus |
| The TypeScript corpus | dokimi, `corpus/typescript` | The subjects, the controls, the ports and the overlay |
| The assert library | `@dokimi/assert` | `Seat`, the property engine, the machines, the task scheduler and the histories, at definition 16.2.0 |

### The chain of the contract

The runtime contract has one chain in TypeScript too: dokimi-spec, then the TypeScript runtime, then dokimi.

- The `.spec.yaml` of dokimi-spec lists the TypeScript runtime beside dokimi-go, so `spec try` writes a change of the contract into a checkout of each and runs its tests:

  ```yaml
  implementations:
    - {repo: dokimasia/dokimi-go, languages: [go]}
    - {repo: dokimasia/dokimi-typescript, languages: [typescript]}
  ```

- The `.spec.yaml` of the TypeScript runtime vendors the contract into `spec/`, as dokimi-go does:

  ```yaml
  implements:
    - definition: dokimasia/dokimi-spec
      path: spec
      languages: [typescript]
      test: [npm, test]
  ```

- `npm run generate` runs `tools/spec-types.ts`, which renders a module of `src/spec` for each section of the copy. A spec of the tool renders the copy again and compares each module with its file.
- Neither `@dokimi/runtime` nor `@dokimi/assert` has a release while this design is built. dokimi runs the TypeScript corpus against the checkouts of both packages through a scratch workspace, whose `package.json` takes them through `file:` dependencies. It runs the Go corpus against unreleased modules through a scratch `go.work` in the same way. The gate of dokimi runs the TypeScript corpus against releases of both packages from npm.

### The package

`@dokimi/runtime` is one package. Each file of `src/` is a module, which the exports map publishes as a subpath, such as `@dokimi/runtime/conform/handle`. No module re-exports another. The helpers that only the runtime calls are under `src/internal/`, which the exports map leaves out. The package has the peer dependencies `@dokimi/assert` and `vitest`, and it needs Node 26 or later.

| Directory | Go module | Imports |
|---|---|---|
| `spec/` | `go.dokimi.dev/spec` | No other directory |
| `fault` | `go.dokimi.dev/fault` | No other directory |
| `stub/` | `go.dokimi.dev/stub` | `fault` |
| `medium/` | `go.dokimi.dev/medium` | No other directory |
| `host/` | `go.dokimi.dev/host` | `fault` and `medium` |
| `conform/` | `go.dokimi.dev/conform` | `spec`, `fault`, `medium` and `host` |
| `probe/` | `go.dokimi.dev/conform/probe` | `spec`, `conform`, `medium` and `host` |
| `sim/` | `go.dokimi.dev/sim` | `spec`, `conform`, `medium` and `host` |

A rule of the lint configuration of the repository refuses an import against this direction.

### Go and TypeScript

| Go | TypeScript |
|---|---|
| A package, such as `conform` | A directory of modules, such as `conform/`, with one module for each file |
| An exported function or variable | A function or a constant in camel case: `conform.Run` is `run` |
| An exported type | A type or a class of the same name: `conform.Handle` is `Handle` |
| An option function, such as `conform.Durable(open)` | A field of an options object, such as `durable` of `Options<S>` |
| A sentinel error `ErrX` | A class `XError` that extends `Error`, in the module `errors` of its directory |
| `errors.Is(err, target)` | The error, or an error of its chain of `cause`, is the target or an instance of the target class, as the assertions of `@dokimi/assert` match an error |
| A panic | A throw of a `TypeError`, a `ReferenceError` or a `RangeError` |
| `context.Context` and `ctx.Err()` | `AbortSignal` and `signal.reason` |
| The results `(V, error)` | The result `V`. A failed call throws, or its promise rejects |
| The results `(V, bool)` of a lookup | `V \| undefined`, and the tuple `[V, boolean]` where the declaration returns one |
| A function that blocks | A function that returns a promise |
| A goroutine | A task: an async function that a host or the task scheduler starts |
| `synctest.Wait` | The quiescence wait of the event loop |
| The clock of a bubble | The fake timers of Vitest |
| `time.Duration` | A number of milliseconds |
| `time.Time` | `Date` |
| A checksum or a seed of 64 bits | `bigint` |
| `rand.Source` | `RandomSource` of `fault`, a function that returns a number in [0, 1) |
| `reflect.Type` as the key of a door or of a dependency | `Token<T>`: the class of the value, or an `InjectionToken` for a type without a class |
| `iter.Seq`, `iter.Seq2` and `iter.Pull` | `Iterable`, `AsyncIterable` and `AsyncIterator` |
| A channel that a callable receives from | `AsyncIterable` |
| `*testing.T` and `t.Run` | A Vitest test and a `describe` block |
| `*testing.B` | `bench` of Vitest |
| `*testing.F` | A target of Jazzer.js, which `prop.fuzz` returns |
| The attributes of `go test -json` | The meta of a Vitest test |
| `t.Parallel()` | No equivalent. The tests of a block run in order |
| The flags of `os.OpenFile` and an `fs.FileMode` | The constants `O_*` of `node:fs`, and a number of the mode |
| `syscall.ENOSPC` | An error whose `code` is `ENOSPC` |
| `net.Listener`, `net.Conn` and `net.PacketConn` | `Listener`, a `Socket` of the Sockets API of WinterTC, and `DatagramSocket` |
| `runtime.Gosched` | `await env.yield()` |
| The profiler labels of an incarnation | The tasks that the host started for the incarnation |
| A count of allocations | None |
| The race detector | None |

An identifier of the spec has the TypeScript name of `naming.json` in the copy, such as `underCancelled` for the probe `under-cancelled`. Every other member of the runtime keeps the contract of its Go form under the mapping of the table. The sections that follow give each declaration whose form differs from that mapping.

The earlier designs of the TypeScript binding fix these names, and this design keeps them:

- The fields `checks`, `excuse`, `reference`, `second`, `restart`, `crash`, `faultable` and `supplies` of `Options<S>`.
- The fields `doors`, `options` and `around` of `Check<S>`, and the fields `waits` and `runsUntilCancelled` of `Call<S>`.
- `supply`, the field `values` of `Doors<S>`, and `handle.open` for the door of a value.

### The runtime model

`tools/spec-types.ts` renders one module for each section of the copy. The module of the claims:

```ts
// @dokimi/runtime/spec/claims
import type { DefectFormId } from "./defects.js";
import type { LawId } from "./laws.js";
import type { MachineryId } from "./machineries.js";
import type { ProbeId } from "./probes.js";
import type { ReasonId } from "./reasons.js";

/** The identifier of a claim, such as "cancel": one member for each claim of the copy. */
export type ClaimId = "accumulates" | "agrees-with-its-sibling" | "cancel" | "zero-on-miss";

/** The planted defect of a claim: its defect form and the reason code of its probe. */
export interface Defect {
  readonly form: DefectFormId;
  readonly reason: ReasonId;
}

/** A claim of the spec. */
export interface Claim {
  readonly id: ClaimId;
  /** The sentence of the check, with {declaration}, {callable} and each partner in braces. */
  readonly statement: string;
  readonly law: LawId;
  readonly probe: ProbeId;
  readonly defect: Defect;
  readonly since: string;
  /** The machineries of the claim, the deepest first. */
  readonly machineries: readonly MachineryId[];
}

/** Returns the claims of the copy, in the order of the copy. */
export declare function claims(): readonly Claim[];

/** Returns the claim with the identifier id, and undefined for an id of no claim of the copy. */
export declare function claimOf(id: string): Claim | undefined;
```

The other sections follow the same form, with the fields of their Go types:

| Module | Union types | Entry types | Functions |
|---|---|---|---|
| `spec/claims` | `ClaimId` | `Claim` and `Defect` | `claims` and `claimOf` |
| `spec/laws` | `LawId` | `Law` with `statement` | `laws` and `lawOf` |
| `spec/probes` | `ProbeId` and `InputKind` | `Probe` with `summary`, `vacuous`, `inputs` and `reasons`, and `Input` | `probes` and `probeOf` |
| `spec/reasons` | `ReasonId` | `Reason` with `meaning` | `reasons` and `reasonOf` |
| `spec/defects` | `DefectFormId` | `DefectForm` with `does` | `defectForms` and `defectFormOf` |
| `spec/machineries` | `MachineryId` and `Supplier` | `Machinery` with `supplies`, `by` and `param` | `machineries` and `machineryOf` |
| `spec/gates` | `GateId` | `Gate` | `gates` and `gateOf` |
| `spec/stances` | `StanceId` | `Stance` with `meaning` and `closes` | `stances` and `stanceOf` |
| `spec/key` | None | None | `key` |
| `spec/version` | None | None | `version` |
| `spec/records` | `Mode`, `Result`, `Application` and `Seed` | The types of the suite record | `modes` and `results` |

The module of the records declares the record in the JSON names of `schemas/suite.schema.json` of the copy:

```ts
// @dokimi/runtime/spec/records
export type Mode = "tasks" | "native";
export type Application = "falsification" | "subject";
export type Result =
  | "pass" | "fail" | "vacuous" | "excused" | "skipped" | "unchecked"
  | "proven" | "passed" | "wrong-reason" | "not-isolated" | "not-deeper" | "equivalent";
/** The seed of an application or of a property run, as 0x and 16 hexadecimal digits. */
export type Seed = string;

export interface SuiteRecord {
  readonly spec: string;
  readonly suite: { readonly language: string; readonly package: string; readonly declaration: string };
  readonly generator: RecordTool;
  readonly runtime: RecordTool;
  readonly seed: Seed;
  readonly application: Application;
  readonly subject?: string;
  readonly ended?: string;
  readonly checks: readonly CheckRecord[];
  readonly refusals: readonly RefusalRecord[];
  readonly unclaimed: readonly UnclaimedRecord[];
}

export interface WorldRecord {
  readonly mode: Mode;
  readonly crashes?: Readonly<Record<string, number>>;
  readonly faults: { readonly struck?: Readonly<Record<string, number>>; readonly vacuous?: readonly string[] };
  readonly leaked?: readonly string[];
  readonly replays?: number;
  readonly reproduced?: number;
  readonly steps: number;
  readonly undecided?: number;
}
```

`CheckRecord`, `RefusalRecord`, `UnclaimedRecord` and `RecordTool` have the fields of their Go types under the same JSON names. The runtime has no validator of a record, because the Go runtime does not validate a record while it runs either. `dokimi report` validates each record of the TypeScript corpus against the schema of its copy.

### The suites

A check, a call and a suite keep the fields of their Go forms, without the fields of allocation counts:

```ts
// @dokimi/runtime/conform/suite
import type { prop } from "@dokimi/assert";

/** The identity of a check over subjects of the type S, such as "Store.get/keeps-the-value". */
export type Id<S> = string & { readonly [subjectOf]?: S };

/** The name of a callable of a declaration over subjects of the type S, such as "Store.get". */
export type Callable<S> = string & { readonly [subjectOf]?: S };

/** Returns the identity of the generated binding of the check id, as ID.Generated does in Go. */
export declare function generated<S>(id: Id<S>): Id<S>;

/** One check of a suite. A check has apply or property, and not both. */
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
  readonly apply?: (h: Handle<S>) => Promise<void>;
  /** Applies the check to one case of the property engine, with a handle and a subject of the case. */
  readonly property?: (h: Handle<S>, c: prop.Case) => Promise<void>;
  readonly defect?: (h: Handle<S>) => S | Promise<S>;
  readonly doors?: Doors<S>;
  readonly options?: (h: Handle<S>) => readonly prop.Option[];
  readonly around?: (run: () => Promise<void>) => Promise<void>;
}

/** The budget of one call, in milliseconds. A budget of TypeScript has no allocation ceiling. */
export interface Budget {
  readonly latency?: number;
  readonly mean?: number;
  readonly warmup?: number;
}

/** The call generator of one callable, with its budget. */
export interface Call<S> {
  readonly callable: Callable<S>;
  /** Generates the calls of the callable. The promise of a call rejects with the error of the callable. */
  readonly calls: prop.Generator<(signal: AbortSignal, s: S) => Promise<void>>;
  /** Returns the planted defect of a latency ceiling. Its callable awaits spend in each call. */
  readonly over?: (h: Handle<S>, spend: () => Promise<void>) => S;
  readonly budget?: Budget;
  readonly waits?: boolean;
  readonly runsUntilCancelled?: boolean;
}

/** The checks of one declaration, with what the generator declined. */
export interface Suite<S> {
  readonly declaration: string;
  readonly package: string;
  readonly generator: { readonly name: string; readonly version: string };
  readonly checks: readonly Check<S>[];
  readonly refusals: readonly Refusal[];
  readonly unclaimed: readonly Unclaimed[];
  readonly calls: readonly Call<S>[];
}
```

`Refusal` and `Unclaimed` keep their TypeScript forms. A suite of TypeScript has no `using` and no `sinks`. The generated module builds the generator of each input from the shape definitions that dokimi writes into it. A consumer registers the generator of any other shape definition with `prop.register` in the setup file of Vitest.

The doors of a handle may return a promise, because a subject that reads its medium or its host at its start cannot finish in a constructor:

```ts
// @dokimi/runtime/conform/doors
/** The doors that only the code that builds a subject can supply. The handle awaits each door. */
export interface Doors<S> {
  readonly reference?: (h: Handle<S>) => S | Promise<S>;
  readonly second?: (h: Handle<S>) => S | Promise<S>;
  readonly restart?: (h: Handle<S>, s: S) => S | Promise<S>;
  readonly crash?: (h: Handle<S>, s: S, call: () => Promise<void>) => S | Promise<S>;
  readonly faults?: (h: Handle<S>, injector: Injector) => S | Promise<S>;
  /** Returns a subject over the medium m, as the door Open does in Go. */
  readonly openOver?: (h: Handle<S>, m: FS) => S | Promise<S>;
  readonly host?: (h: Handle<S>, env: Env) => S | Promise<S>;
  readonly values?: readonly Supplied<S>[];
}

/** The door of a value of one token, which supply returns. */
export interface Supplied<S> {
  readonly token: Token<unknown>;
  readonly open: (h: Handle<S>, s: S) => unknown;
}

/** Returns the door of a value of the token, which open opens for a subject s. */
export declare function supply<S, V>(token: Token<V>, open: (h: Handle<S>, s: S) => V | Promise<V>): Supplied<S>;
```

```ts
// @dokimi/runtime/conform/token
/**
 * What identifies the values of the type T at run time: the class of the values, or an injection
 * token that a suite declares for a type without a class, such as an interface.
 */
export type Token<T> = (abstract new (...args: never[]) => T) | InjectionToken<T>;

/** A token of the type T, as the injection token of a dependency-injection container. */
export declare class InjectionToken<T> {
  constructor(description: string);
  readonly description: string;
}
```

The door of a value takes its name from an accepted design, so the door of the medium is `openOver` in TypeScript, where Go names it `Open`.

The options of an application are the fields of one object. Each field has the contract of the Go option in its comment:

```ts
// @dokimi/runtime/conform/options
export interface Options<S> {
  /** The subject of the records, as the generator writes it from the directive subject. */
  readonly subject?: string;
  readonly checks?: readonly Check<S>[]; // With
  readonly excuse?: Readonly<Record<string, string>>; // Excuse, by the identity of the check
  readonly excuseSuite?: string; // ExcuseSuite
  readonly isolate?: boolean; // Isolate
  readonly cases?: number; // Cases
  readonly ceilings?: Readonly<Record<string, Budget>>; // Ceilings, by the callable
  readonly reference?: (h: Handle<S>) => S | Promise<S>; // Reference
  readonly second?: (h: Handle<S>) => S | Promise<S>; // Second
  readonly restart?: (h: Handle<S>, s: S) => S | Promise<S>; // Restartable
  readonly crash?: (h: Handle<S>, s: S, call: () => Promise<void>) => S | Promise<S>; // Crashable
  readonly faultable?: (h: Handle<S>, injector: Injector) => S | Promise<S>; // Faultable
  readonly durable?: (h: Handle<S>, m: FS) => S | Promise<S>; // Durable
  readonly hosted?: (h: Handle<S>, env: Env) => S | Promise<S>; // Hosted
  readonly tasks?: boolean; // Tasks
  readonly supplies?: readonly Supplied<S>[]; // Supplies
}
```

The handle builds the subject of a check before the check runs, so `subject()` returns it at once. Every other method that builds a subject returns a promise:

```ts
// @dokimi/runtime/conform/handle
import type { Clock, Seat, prop, stateful } from "@dokimi/assert";

/** The handle of one check, or of one case of a property check. It implements `Seat` of @dokimi/assert. */
export declare class Handle<S> implements Seat {
  /** The signal of the test of the check, or of the case of a property check. */
  readonly signal: AbortSignal;
  subject(): S;
  new(): Promise<S>;
  fresh(): Promise<S>;
  reference(): Promise<S>;
  second(): Promise<S>;
  restart(s: S): Promise<S>;
  crash(s: S, call: () => Promise<void>): Promise<S>;
  openOver(m: FS): Promise<S>;
  freshOpenOver(m: FS): Promise<S>;
  host(env: Env): Promise<S>;
  freshHost(env: Env): Promise<S>;
  hostable(): boolean;
  /** Returns the value of the token for the subject s, from a door of the defect or of the option supplies. */
  open<V>(token: Token<V>, s: S): Promise<V>;
  injector(): Injector;
  mode(): Mode;
  clock(): Clock;
  advance(ms: number): void;
  rand(): RandomSource;
  tempDir(): string;
  cleanup(f: () => void | Promise<void>): void;
  skip(reason: string): never;
  statement(): string;
  expect(reason: ReasonId): void;
  engage(): void;
  addPoints(n: number): void;
  recordWorld(w: WorldRecord): void;
  schedule(c: prop.Case): stateful.Option;
  yield(): Promise<void>;
}
```

The handle has no `blocking`, because a TypeScript world strikes the blocking faults in both modes. It also has no `unwrap` and no `context`. The handle passes each method of `Seat` that it does not declare to the test of the check, and `signal` takes the place of the context.

The entry points of the suites:

```ts
// @dokimi/runtime/conform/run
/**
 * Registers the block of an application. The block has a Vitest test for each check of suite and
 * of options.checks, and a last test that writes the falsification record and the subject record.
 */
export declare function run<S>(
  suite: Suite<S>,
  newSubject: (h: Handle<S>) => S | Promise<S>,
  options?: Options<S>,
): void;

// @dokimi/runtime/conform/falsify
export declare function falsify<S>(suite: Suite<S>, options?: Options<S>): void;

// @dokimi/runtime/conform/bench
/** Registers a Vitest bench of the callable on a subject of newSubject, under the latency ceilings of its budget. */
export declare function bench<S>(
  suite: Suite<S>,
  callable: Callable<S>,
  newSubject: (h: Handle<S>) => S | Promise<S>,
  options?: Options<S>,
): void;

// @dokimi/runtime/conform/fuzz
/** Replays the stored cases and the seed files of the property checks of the callable, and returns the target of Jazzer.js. */
export declare function fuzz<S>(
  seat: Seat,
  suite: Suite<S>,
  callable: Callable<S>,
  newSubject: (h: Handle<S>) => S | Promise<S>,
  options?: Options<S>,
): Promise<(data: Uint8Array) => Promise<void>>;

// @dokimi/runtime/conform/output
export interface Output {
  readonly values: readonly unknown[];
  readonly failed: boolean;
}
export declare function outputOf(error: unknown, ...values: unknown[]): Output;
/** Completes call in the history of a case: as unknown for an error of the aborted signal, and with its output otherwise. */
export declare function complete(signal: AbortSignal, call: history.Call, error: unknown, ...values: unknown[]): void;

// @dokimi/runtime/conform/named
export interface Named<T> {
  readonly value: T;
  readonly name: string;
}

// @dokimi/runtime/conform/quiescence
/**
 * Returns once every promise job that the code of the test queued has run and no host of the test
 * has a ready task. It awaits turns of the event loop through the setImmediate that the global
 * object had when the module loaded.
 */
export declare function quiescence(): Promise<void>;

// @dokimi/runtime/conform/fake-timers
/** Runs run with the timers and the Date of the process faked from the start of a world, and restores them. */
export declare function fakeTimers(run: () => Promise<void>): Promise<void>;

// @dokimi/runtime/conform/errors
export declare class ClockError extends Error {}
export declare class ConfigError extends Error {}
export declare class GeneratedError extends Error {}
export declare class OptionError extends Error {}
export declare class ReasonError extends Error {}
export declare class SkipError extends Error {}
export declare class TempDirError extends Error {}

// @dokimi/runtime/conform/compat
export declare function compatV1(): void;
```

`runBudgets`, `falsifyAllocs` and `benchParallel` have no TypeScript form. A check of TypeScript counts no allocations, so it has no field `counted`.

### The probes

Each probe of `conform/probe` of Go has a module `@dokimi/runtime/probe/<probe>`, which exports the function of its TypeScript name. Its signature maps the Go signature by the table of this design, and the probe returns a promise:

```ts
// @dokimi/runtime/probe/reads-what-was-stored
export declare function readsWhatWasStored<S, V>(h: Handle<S>, call: Call<readonly V[]>, write: Call<V>): Promise<void>;
```

`zero-beside-error`, `zero-pair-beside-error` and `field-beside-error` have no module. A failed TypeScript call returns no value, so the gate `returns-beside-error` refuses their claims `zero-on-error`, `zeroes-on-error` and `returns-on-error` for each callable of TypeScript. `encodes-each-variant`, `decodes-each-text` and `compares-with-identity` have no module either. A TypeScript enumeration has no encode and no decode, so the gate `text-encoding` refuses the claims `encodes-text` and `decodes-text`. A TypeScript error has no `Is`, so the gate `self-comparing-type` refuses the claim `is-agrees`. That leaves 121 modules for the 127 functions of the Go package. The probes of the rules of a simulation, such as `replies-at-most-once`, are in `sim/rules`, as Go has them in `sim`.

The types of the inputs:

```ts
// @dokimi/runtime/probe/call
export type Call<V> = (signal: AbortSignal) => V | Promise<V>;
export type Act = (signal: AbortSignal) => unknown;
/** A call of a callable with a presence flag: its value and the flag. */
export type Lookup<V> = (signal: AbortSignal) => readonly [V, boolean] | Promise<readonly [V, boolean]>;
export type Pair<A, B> = readonly [A, B];
/** Opens an iteration of a stream. The return of the iterator ends the iteration. */
export type Pull<E> = (signal: AbortSignal) => Promise<AsyncIterator<E>>;

// @dokimi/runtime/probe/sequence
export type Actions<S> = (c: prop.Case, subject: () => S) => readonly stateful.Action<undefined>[];
export type Apply = (signal: AbortSignal, args: readonly unknown[]) => Promise<Output>;
export type Reference<S> = (r: S) => (signal: AbortSignal, operation: string, args: readonly unknown[]) => Promise<Output>;

// @dokimi/runtime/probe/medium
export interface Durable<S, V> {
  readonly write: (signal: AbortSignal, s: S, i: number) => Promise<void>;
  readonly read: (signal: AbortSignal, s: S, i: number) => Promise<V>;
  readonly want: (i: number) => V;
  readonly key: (i: number) => unknown;
}
export declare function unsynced(m: FS): FS;
export declare function unfilled(m: FS): FS;

// @dokimi/runtime/probe/world
export interface Operation<S> {
  readonly name: string;
  readonly input: (c: prop.Case) => unknown;
  readonly args: (input: unknown) => readonly unknown[];
  readonly keys: (input: unknown) => readonly unknown[];
  /** Calls the callable on s with the input. It resolves with the output, and rejects with the error of the callable. */
  readonly call: (signal: AbortSignal, s: S, input: unknown) => Promise<Output>;
  readonly produced?: (out: Output) => void;
}
export declare class Pool<T> {
  constructor(c: prop.Case, g: prop.Generator<readonly T[]>, label: string);
  add(v: T): void;
  pick(c: prop.Case, label: string): T;
}
export interface World<S> {
  readonly operations: (c: prop.Case) => readonly Operation<S>[];
  readonly reference: Reference<S>;
  readonly corrupt?: Token<unknown> | Error;
  readonly clients?: number;
  /** The monotonic time of the host, in milliseconds, within which a call of the recovery completes. */
  readonly bound?: number;
}
export declare function worldOptions<S>(h: Handle<S>): readonly prop.Option[];
export declare function amnesiac<S>(h: Handle<S>, env: Env): Promise<S>;
```

| Input kind | Go | TypeScript |
|---|---|---|
| `call` | `Call[V]`, `Act` or `Lookup[V]` | `Call<V>`, `Act` or `Lookup<V>` |
| `value` | A value of the type, such as its zero value | The value that eidos derives, because TypeScript has no zero value |
| `error` | A sentinel, or a `conform.Named[error]` | An error value or an error class, matched through its chain of `cause` |
| `projection` | `func(V) F` | `(v: V) => F` |
| `duration` | `time.Duration` | A number of milliseconds |
| `clock` | The clock of the handle | The clock of the handle |
| `count` | `int` | `number` |
| `hook` | `func(ctx, hook func()) error` | `(signal: AbortSignal, hook: () => void) => unknown` |
| `predicate` | `Call[bool]` | `Call<boolean>` |
| `scope` | `func(ctx, body Act) error` | `(signal: AbortSignal, body: Act) => Promise<void>` |
| `values` | `[]conform.Named[E]` | `readonly Named<E>[]` |
| `actions`, `apply` and `reference` | `Actions[S]`, `Apply` and `Reference[S]` | `Actions<S>`, `Apply` and `Reference<S>` |
| `operation` | The name of an action | The name of an action |
| `pull` | `Pull[E]` | `Pull<E>` |
| `channel` | `<-chan V` | `AsyncIterable<V>` |
| `handle-method` | A method expression, such as `(*Sub).Next` | A function of the handle of a call, such as `(sub) => sub.next(ms)` |
| `durable` and `world` | `Durable[S, V]` and `World[S]` | `Durable<S, V>` and `World<S>` |

`worldOptions` returns `prop.confirm(10)` in both modes, where Go returns `prop.Replays(10)` in native mode. A replay of a case in native mode repeats the whole case in TypeScript, as the section on the host states.

### The doubles

The doubles take the behaviour suites of Go. `Subject` of `stub/behaviour` gains the fields of Go, and its field `override` passes the call of the arguments:

```ts
// @dokimi/runtime/stub/behaviour
/** One method of a generated double, as behaviour and resetBehaviour drive it. */
export interface Subject<C, R> {
  readonly stub: MethodStub<C>;
  readonly call: () => unknown;
  readonly arguments?: (seat: Seat, call: C) => void;
  readonly reset: () => void;
  /** Calls the method, and resolves with the error that it throws or rejects with. */
  readonly fails?: () => Promise<unknown>;
  readonly result?: () => Promise<R>;
  readonly default?: () => R;
  readonly returns?: () => R;
  readonly override?: (seen: (call: C) => void) => void;
}
export declare function behaviour<C, R>(name: string, newSubject: (seat: Seat) => Subject<C, R>): Check[];
export declare function resetBehaviour<C, R>(name: string, newSubject: (seat: Seat) => Subject<C, R>): Check[];
```

- `behaviour` runs the 13 checks of the Go suite, and `resetBehaviour` its two checks. `Instance` of `stub/double-behaviour` loses its field `reset`.
- A double has no `newDelayed` and no `newFailing`. A test passes the options that each passed: `new StoreStub(seat, StoreStub.withClock(clock), StoreStub.withLatency(ms))`, and `new StoreStub(seat, StoreStub.failing(error))`.
- The decision of a fault may be a promise. `answerAsync` awaits it, so a hold of the injector blocks the call of an asynchronous method. A synchronous method cannot wait, so its double takes its fault from `faultOfSync`, whose decision never waits.

### The medium

```ts
// @dokimi/runtime/medium/fs
/** The storage of a subject. Each name is slash-separated and relative to the root. */
export interface FS {
  /** Opens the file name with the constants O_* of node:fs, and creates it with mode under O_CREAT. */
  openFile(name: string, flags: number, mode?: number): Promise<File>;
  remove(name: string): Promise<void>;
  rename(oldName: string, newName: string): Promise<void>;
  mkdirAll(name: string, mode?: number): Promise<void>;
  /** Returns the entries of the directory name, sorted by name. */
  readDir(name: string): Promise<readonly DirEntry[]>;
  stat(name: string): Promise<FileInfo>;
  /** Makes the entries of the directory durable, such as a file that a call created or renamed in it. */
  syncDir(name: string): Promise<void>;
}

/** An open file of an FS. A call after close rejects with an error whose code is EBADF. */
export interface File {
  read(buffer: Uint8Array): Promise<number>;
  readAt(buffer: Uint8Array, offset: number): Promise<number>;
  write(data: Uint8Array): Promise<number>;
  writeAt(data: Uint8Array, offset: number): Promise<number>;
  seek(offset: number, whence: "start" | "current" | "end"): Promise<number>;
  sync(): Promise<void>;
  truncate(size: number): Promise<void>;
  stat(): Promise<FileInfo>;
  close(): Promise<void>;
}

export interface FileInfo {
  readonly name: string;
  readonly size: number;
  readonly mode: number;
  readonly modified: Date;
  readonly directory: boolean;
}

export interface DirEntry {
  readonly name: string;
  readonly directory: boolean;
}

// @dokimi/runtime/medium/os
/** Returns the FS of the operating system under the directory root, over node:fs/promises. */
export declare function os(root: string): FS;

// @dokimi/runtime/medium/memory
/** An FS in memory with the model of persistence of the Go medium. */
export declare class Memory implements FS {
  crash(): void;
  crashAt(point: number): void;
  crashKeeping(keep: (name: string, pending: number) => number): void;
  kill(): void;
  tear(n: number): void;
  corrupt(name: string, offset: number): void;
  fill(n: number): void;
  failBefore(error: unknown): void;
  failAfter(error: unknown): void;
  failPart(n: number, error: unknown): void;
  hold(): void;
  release(): void;
  held(): boolean;
  points(): number;
  written(): number;
}

// @dokimi/runtime/medium/errors
export declare class CrashedError extends Error {}
```

A method rejects with a `TypeError` for a name that is no valid path, and a write that the fill of `fill` refuses rejects with an error whose `code` is `ENOSPC`. A call that `hold` blocks waits on a promise, which `release` or a crash settles.

### The host

```ts
// @dokimi/runtime/host/host
export type Kill = "process-crash" | "power-loss" | "disk-replacement";

/** A simulated machine. It keeps its disk, its clocks and its address across the incarnations of its process. */
export declare class Host {
  /** With the scheduler, the host runs its tasks, its sleeps and its timers under it. Without one, it runs them on the event loop. */
  constructor(c: prop.Case, name: string, network?: Network, scheduler?: stateful.Scheduler);
  name(): string;
  disk(): Memory;
  injector(): Injector;
  start(): { readonly incarnation: Incarnation; readonly env: Env };
  broadcast(): void;
  bump(ms: number): void;
  strobe(delta: number, period: number): void;
  resetClock(): void;
  wallRead(): boolean;
}

// @dokimi/runtime/host/incarnation
export declare class Incarnation {
  /** Runs task as a task of the incarnation, on the event loop or under the scheduler of the host. */
  spawn(task: () => Promise<void>): void;
  crash(kill: Kill): void;
  pause(): void;
  resume(): void;
  /** Returns the stack at the spawn of each task of the incarnation that has not settled after its crash and quiescence. */
  leaked(): readonly string[];
}

// @dokimi/runtime/host/env
/** The environment of one incarnation. Each method returns the same resource for the whole incarnation. */
export declare class Env {
  name(): string;
  incarnation(): number;
  /** The signal of the incarnation, which a crash and the end of the case abort. */
  lifetime(): AbortSignal;
  disk(): FS;
  clock(): Clock;
  rand(): RandomSource;
  dependency(name: string): Dependency;
  spawn(task: () => Promise<void>): void;
  /** A scheduling point under the scheduler, and a turn of the event loop without one. It wakes each wait. */
  yield(): Promise<void>;
  /** Resolves once cond returns true. The host calls cond again after each yield, each settled task and each broadcast. */
  wait(cond: () => boolean): Promise<void>;
  listen(port: number): Listener;
  /** Connects to the listener at address, such as "node-2:7000", as connect of the Sockets API does. */
  connect(address: string): Socket;
  listenDatagram(port: number): DatagramSocket;
}

// @dokimi/runtime/host/clock
export declare class Clock {
  wall(): Date;
  monotonic(): number;
  sleep(ms: number): Promise<void>;
  /** Calls f once the monotonic clock has advanced by ms. stop cancels the call, and reports whether it did so before f started. */
  afterFunc(ms: number, f: () => void | Promise<void>): () => boolean;
}

// @dokimi/runtime/host/socket
/** A stream between two hosts, with the members of the Socket of the Sockets API of WinterTC. */
export interface Socket {
  readonly readable: ReadableStream<Uint8Array>;
  readonly writable: WritableStream<Uint8Array>;
  readonly opened: Promise<{ readonly remoteAddress: string; readonly localAddress: string }>;
  readonly closed: Promise<void>;
  close(reason?: unknown): Promise<void>;
}

export interface Listener extends AsyncIterable<Socket> {
  readonly address: string;
  accept(): Promise<Socket>;
  close(): Promise<void>;
}

export interface DatagramSocket {
  readonly address: string;
  send(data: Uint8Array, address: string): Promise<void>;
  receive(): Promise<{ readonly data: Uint8Array; readonly address: string }>;
  close(): Promise<void>;
}

// @dokimi/runtime/host/network
export interface NetworkOptions {
  /** The lowest and the highest latency of a link, in milliseconds: 0 and 100 by default. */
  readonly latency?: readonly [number, number];
  /** The lowest throughput of a link, in bytes per second. A link has no limit by default. */
  readonly throughput?: number;
}

export declare class Network {
  constructor(c: prop.Case, options?: NetworkOptions);
  partition(...sides: readonly (readonly string[])[]): void;
  repair(): void;
  reset(a: string, b: string): void;
  slow(a: string, b: string, latency: number): void;
  restore(a: string, b: string): void;
  drop(to: string): boolean;
  delay(to: string): boolean;
  duplicate(to: string): boolean;
  inFlight(to: string): number;
}

// @dokimi/runtime/host/injector
export declare class Injector {
  failNext(error: unknown): void;
  failAt(n: number, error: unknown): void;
  failUntilHealed(error: unknown): void;
  heal(): void;
  hold(): void;
  release(): void;
  failures(): number;
  calls(): readonly DependencyCall[];
  dependency(): Dependency;
}

export declare class Dependency {
  /** Records a call of the callable, waits while a hold is in force, and rejects with the error of a failed call. */
  call(callable: string): Promise<void>;
  /** Records a call of the callable, and throws the error of a failed call. A hold does not block it. */
  callSync(callable: string): void;
}

export interface DependencyCall {
  readonly callable: string;
  readonly dependency: number;
  readonly failed: boolean;
}

/** Returns the fault strategy of an asynchronous method of a double of the dependency. */
export declare function faultOf<C>(dependency: Dependency, callable: string): Fault<C>;
/** Returns the fault strategy of a synchronous method of a double of the dependency. */
export declare function faultOfSync<C>(dependency: Dependency, callable: string): Fault<C>;

// @dokimi/runtime/host/park
/** Waits until signal aborts, and rejects with its reason. Under the scheduler of the host it parks the running task. */
export declare function park(signal: AbortSignal): Promise<never>;

// @dokimi/runtime/host/errors
export declare class InjectedError extends Error {}
export declare class NoNetworkError extends Error {}
export declare class ResetError extends Error {}
```

The network has no flag `blocking`. Each stream delays its bytes by the latency and the throughput of its link, a stream across a cut waits for the repair, and `slow` works on every network.

#### Execution modes

| Mode | Code that runs | Quiescence | Replay | Faults |
|---|---|---|---|---|
| Machines | The sans-IO machines of sim | The run applies one event at a time | The whole case | Every fault of sim |
| Tasks | Code that waits only through its host, as tasks of the task scheduler | `run` and `advance` of the scheduler return once no task is ready | The whole case | Every fault of the host |
| Native | Code on the event loop that waits through its host or through promises that it creates | A turn of the event loop completes no operation of the host | The whole case, for code that reads no real time, makes no real I/O and takes no global source of randomness | Every fault of the host |

A call of the disk that a hold blocks is a pending promise in TypeScript, and the quiescence wait treats it as any other wait. So native mode strikes the blocking faults without an option. The event loop runs the promise jobs in an order that the code and the world decide, so a replay of a case repeats the whole case. Native mode explores no interleaving beyond the choices of the world. Tasks mode explores the interleavings at each wait of the host, and it needs the waits of tasks of definition 16.0.0.

The quiescence wait of native mode has three steps:

1. The host counts each operation that completes: a wait, a sleep, a read, an accept, a callback of a timer and a task that settles.
2. The world awaits a turn of the event loop through the `setImmediate` that the global object had when the runtime loaded. Each promise job that the code queued runs before the callback of that turn.
3. The world repeats the turn while a turn completes an operation of the host.

A world in native mode installs the fake timers of Vitest for each case, at the start of the world, without `setImmediate`, `clearImmediate`, `nextTick` and `queueMicrotask`. `advance` of the world moves them by the same time as the run. A component that calls `setTimeout` or reads `Date` directly then sees the time of the run, as a goroutine in a bubble sees the clock of the bubble. It does not receive a fault of a clock. The world restores the real timers at the end of the case. `fakeTimers` of conform does the same around a check of an application with `hosted`.

The mode has two limits:

- **A stall.** A component that runs without an await, or that queues promise jobs without an end, keeps the turn of the quiescence wait from running. The callback of a watchdog would wait as well, so the runtime has no watchdog of a stall.
- **Leaked tasks.** `leaked` covers the tasks that the host started. The host does not see a promise that the code of a component starts without `spawn`.

### The simulation

```ts
// @dokimi/runtime/sim/machine
/** A sans-IO machine: it applies one event at a time and returns the effects of the event. */
export interface Machine<E, F> {
  step(event: E): readonly F[];
  /** Returns the error with which the machine has stopped itself, and undefined while it runs. */
  err(): unknown;
}
export declare function plant<E, F>(machine: Machine<E, F>, anatomy: Anatomy<E, F>, form: DefectFormId): Machine<E, F>;

// @dokimi/runtime/sim/anatomy
export interface Reply {
  readonly to: number;
  readonly value: unknown;
  readonly error: unknown;
}

export interface Anatomy<E, F> {
  readonly started: (reply: number) => E;
  readonly stopped: () => E;
  readonly cancelled?: (reply: number) => E;
  readonly call: (event: E) => number | undefined;
  readonly reply: (effect: F) => Reply | undefined;
  readonly unknown?: (error: unknown) => boolean;
  readonly token: (effect: F) => number | undefined;
  readonly ends?: (effect: F) => number | undefined;
  readonly dependency: (effect: F) => string | undefined;
  readonly fail: (effect: F, error: unknown) => E | undefined;
  readonly end?: (effect: F, error: unknown) => E | undefined;
  readonly errors: (effect: F) => readonly unknown[];
  /** The timeout of the effect in milliseconds, and 0 for an effect whose call the run alone bounds. */
  readonly timeout: (effect: F) => number;
  readonly perform: Perform<E, F>;
}

/** Performs the effect of the process p. The world awaits the outcome before its next step. */
export type Perform<E, F> = (w: World<E, F>, p: Process<E, F>, effect: F) => Outcome<E, F> | Promise<Outcome<E, F>>;

export interface Outcome<E, F> {
  readonly event?: E;
  readonly next?: () => { readonly event?: E; readonly done: boolean } | Promise<{ readonly event?: E; readonly done: boolean }>;
  readonly message?: { readonly call: (reply: number) => E; readonly result: (value: unknown, error: unknown) => E; readonly to: string };
  readonly launch?: Launch<E, F>;
  readonly cancel?: () => void | Promise<void>;
  readonly stop?: string;
}

export interface Launch<E, F> {
  readonly name: string;
  readonly machine: (clock: Clock) => Machine<E, F>;
  readonly call: (reply: number) => E;
  readonly result: (value: unknown, error: unknown) => E;
  readonly deps?: readonly unknown[];
}

// @dokimi/runtime/sim/world
/** The time at which every clock of a world starts: 2000-01-01T00:00:00Z, in milliseconds since the Unix epoch. */
export declare const start: number;

export declare class World<E, F> {
  constructor(c: prop.Case, anatomy: Anatomy<E, F>);
  spawn(name: string, newMachine: (clock: Clock) => Machine<E, F>, ...deps: unknown[]): void;
  addHost(name: string): Host<E, F>;
  deliver(to: string, event: E): void;
  call(client: number, to: string, operation: string, args: readonly unknown[], keys: readonly unknown[], call: (reply: number) => E): void;
  invoke(client: number, process: string, operation: string, args: readonly unknown[], keys: readonly unknown[], call: (signal: AbortSignal) => Promise<unknown>): void;
  cancel(client: number): boolean;
  after(p: Process<E, F>, ms: number, event: E): void;
  afterFunc(p: Process<E, F>, ms: number, f: () => void): void;
  step(name: string): Promise<boolean>;
  advance(ms: number): Promise<void>;
  next(): number | undefined;
  stop(): Promise<void>;
  now(): Date;
  checksum(): bigint;
  log(): readonly Entry<E, F>[];
  history(): history.History;
  pending(): readonly string[];
  process(name: string): Process<E, F> | undefined;
  faults(): readonly string[];
  signal(): AbortSignal;
  crash(name: string): void;
  restart(name: string): void;
  pause(name: string): void;
  resume(name: string): void;
  partition(...sides: readonly (readonly string[])[]): void;
  repair(): void;
  sides(): readonly (readonly string[])[] | undefined;
  drop(to: string): boolean;
  delay(to: string): boolean;
  duplicate(to: string): boolean;
  loseBefore(dependency: string): void;
  loseAfter(dependency: string): void;
  outage(dependency: string): void;
  slow(dependency: string): void;
  restore(dependency: string): void;
  jump(name: string, ms: number): void;
  strobe(name: string, a: number, b: number, steps: number): void;
  tear(name: string, n: number): void;
  corrupt(name: string, file: string, offset: number): void;
  fill(name: string, n: number): void;
  restoreDisk(name: string): void;
  reset(a: string, b: string): void;
  slowLink(a: string, b: string, latency: number): void;
  restoreLink(a: string, b: string): void;
  replaceDisk(name: string): void;
  recover(): void;
}

// @dokimi/runtime/sim/simulation
export type Fault =
  | "drop" | "delay" | "duplicate" | "partition" | "repair" | "lose-before" | "lose-after"
  | "outage" | "slow" | "crash" | "pause" | "restart" | "clock-jump" | "clock-strobe"
  | "torn-write" | "corruption" | "full-disk" | "reset" | "slow-link" | "restore-link" | "disk-replacement";
export type Isolation = "serializable" | "snapshot-isolation";
export type System<E, F> = (c: prop.Case) => (w: World<E, F>) => void | Promise<void>;

export interface Action<E, F> {
  readonly name: string;
  readonly enabled?: (w: World<E, F>) => boolean;
  readonly input?: (c: prop.Case) => unknown;
  readonly run: (w: World<E, F>, client: number, input: unknown) => void | Promise<void>;
  readonly call?: boolean;
  readonly drain?: boolean;
}

export interface Invariant<E, F> {
  readonly statement: string;
  /** Checks the world between two steps. It must not wait. */
  readonly check: (h: Check, w: World<E, F>) => void;
  readonly defect?: System<E, F>;
}

export interface Convergence<E, F> {
  readonly statement: string;
  readonly check: (h: Check, w: World<E, F>) => void;
  readonly defect: System<E, F>;
  /** The time in milliseconds within which the convergence becomes true once the faults end. */
  readonly within: number;
}

export interface Simulation<E, F, S> {
  readonly anatomy: Anatomy<E, F>;
  readonly system: System<E, F>;
  readonly spec?: history.Spec<S>;
  readonly fingerprint?: (w: World<E, F>) => bigint;
  readonly tolerate?: Readonly<Partial<Record<Fault, number>>>;
  readonly require?: Readonly<Partial<Record<Fault, number>>>;
  readonly actions?: readonly Action<E, F>[];
  readonly rules?: readonly Rule[];
  readonly invariants?: readonly Invariant<E, F>[];
  readonly converges?: readonly Convergence<E, F>[];
  readonly scenarios?: readonly { readonly name: string; readonly run: (w: World<E, F>) => void | Promise<void> }[];
  readonly history?: readonly history.Option[];
  readonly around?: (run: () => Promise<void>) => Promise<void>;
  readonly mode?: Mode;
  readonly network?: NetworkOptions;
  readonly clients?: number;
  readonly meanSteps?: number;
  readonly maxSteps?: number;
  readonly strategy?: Strategy;
  readonly isolation?: Isolation;
}

// @dokimi/runtime/sim/run
export declare function run<E, F, S>(c: prop.Case, s: Simulation<E, F, S>): Promise<void>;
export declare function test<E, F, S>(seat: Seat, s: Simulation<E, F, S>): Promise<void>;
export declare function prove<E, F, S>(seat: Seat, s: Simulation<E, F, S>, check: string): Promise<void>;
/** Registers a Vitest test for each rule, invariant and convergence of s, which proves the check. The tests run in order. */
export declare function falsify<E, F, S>(s: Simulation<E, F, S>): void;
export declare function fuzz<E, F, S>(seat: Seat, s: Simulation<E, F, S>): Promise<(data: Uint8Array) => Promise<void>>;

// @dokimi/runtime/sim/component
export declare class Component<T> {
  running(): boolean;
  value(): T;
}
/** Adds a component process named name on a new host of that name, and resolves once its first incarnation has opened. */
export declare function spawnComponent<T, E, F>(w: World<E, F>, name: string, open: (env: Env) => T | Promise<T>): Promise<Component<T>>;
export declare function componentOf<T, E, F>(w: World<E, F>, name: string): Component<T>;

// @dokimi/runtime/sim/dependency
/** Returns the first dependency of the process p of the token, as Dependency of Go returns it by its type. */
export declare function dependency<T, E, F>(w: World<E, F>, p: Process<E, F>, token: Token<T>): T;

// @dokimi/runtime/sim/errors
export declare class CancelledError extends Error {}
export declare class CrashedError extends Error {}
export declare class LostError extends Error {}
export declare class RunEndedError extends Error {}
```

`Process`, `Host`, `Clock`, `Entry`, `Check`, `Rule` and `Strategy` with `pct` and `uniform` have the members of their Go forms. The events and the effects of a world without machine processes are `never`, where Go has `sim.None`. A simulation has neither the field `blocking` nor the field `serial`. Both modes strike the blocking faults, and the proofs of `falsify` run in order. The world awaits the promise of a performer, of a `next` of a wait and of a `cancel` within its step. When only a later step of the world settles such a promise, the world fails the case with the effect, because the step cannot end.

### The generated code

| Generator | Module | Contents |
|---|---|---|
| suite | `<module>.suite.ts` | The names of the callables and of the checks, the suite, the planted defects, the actions, the operations, the pools, the reference and the applications |
| suite | `<module>.suite.test.ts` | The falsification of the checks of the consumer |
| suite | `<subject module>.suite.test.ts` | The application to the subject, with the options of the directive `suite` |
| suite | `<module>.suite.bench.ts` | A bench of each callable under the latency ceilings of its budget |
| suite | `<module>.suite.fuzz.ts` | The target of Jazzer.js of each callable |
| fake | `<module>.fake.ts` | A class over a `Map` or an array with the effects of the reference forms of the declaration, and its absent error class |
| fake | `<module>.fake.test.ts` | The application of the suite of the declaration to the fake |
| stub | `<module>.stub.ts` and `<module>.stub.test.ts` | The double without `newDelayed` and `newFailing`, and the subjects of `behaviour` and `resetBehaviour` |
| sim | `<module>.sim.ts` | The anatomy, the performer and the simulation, such as `machineAnatomy`, `performMachine` and `machineSimulation` |
| sim | `<module>.sim.test.ts` | `test` and `falsify` of the simulation |
| sim | `<module>.sim.fuzz.ts` | The target of Jazzer.js of the simulation |

The suite generator gains `typescript.go` beside `golang.go`, with the TypeScript forms of the closures, the inputs, the planted defects and the operations. The templates have the reduced form of the Go tree:

- A derived check leaves out its probe, its reason and its law, which the runtime reads from its claim.
- The tests of a file excuse the checks of a suite through `excuseSuite`.
- A generated binding takes its identity from `generated`.
- A generated module has no header comment.

A derived check of generated inputs:

```ts
{
  id: generated(storeChecks.put.smoke()),
  statement: "put returns without a panic for derived inputs",
  claim: "smoke",
  machineries: ["generated-inputs", "nothing-beyond-a-caller"],
  property: async (h, c) => {
    const drawnValue = c.draw(valueShape, "value");
    await survives(h, (signal) => h.subject().put(drawnValue, signal));
  },
  defect: storePutPanics,
},
```

The directives of the sim generator take these TypeScript forms:

| Directive | Go | TypeScript |
|---|---|---|
| `sim` on the machine interface | `Step(ev E, out []F) []F` and `Err() error` | `step(event: E): readonly F[]` and `err(): unknown`. The events are the members of the union type `E`, and the effects the members of `F` |
| `token` and `token reply` | A defined type over `uint64` | A type alias over `number` |
| `reply value= err= to=` | Fields of the reply type | Properties of the reply type |
| `system`, `scenario` and `defect=` of `invariant` | `func(w *sim.World[E, F])` | `(w: World<E, F>) => void \| Promise<void>` |
| `action` | `func(c *prop.Case) func(w *sim.World[E, F])` | `(c: prop.Case) => (w: World<E, F>) => void \| Promise<void>` |
| `action call` | `func(c *prop.Case) func(w *sim.World[E, F], client int)` | `(c: prop.Case) => (w: World<E, F>, client: number) => void \| Promise<void>` |
| `spec` | `func() history.Spec[S]`, or `func() history.Subject` | `() => history.Spec<S>`, or a function whose subject `history.specFrom` wraps |
| `invariant` and `converges` | `func(h *sim.Check, w *sim.World[E, F])` | `(h: Check, w: World<E, F>) => void` |
| `within=` of `converges` | `func() time.Duration` | `() => number`, in milliseconds |
| `fingerprint` | `func(w *sim.World[E, F]) uint64` | `(w: World<E, F>) => bigint` |
| `around` | `func(t *testing.T, run func(*testing.T))` | `(run: () => Promise<void>) => Promise<void>` |
| `perform func=` | `func(w, p, f T) sim.Outcome[E, F]`, or with `f *T` | `(w: World<E, F>, p: Process<E, F>, effect: T) => Outcome<E, F> \| Promise<Outcome<E, F>>` |

A TypeScript object is a reference, so a performer takes the effect as it is. The generator does not write the address of an effect.

### The report and the gate

- `dokimi report` reads the JSON report of Vitest beside the events of `go test -json`, from the elements of `meta["dokimi.suite"]` of each test.
- `dokimi conform` judges a language from a table of its bindings: Go and TypeScript, each with its test command, its lock and the module or the package of its corpus. The gate of dokimi runs the command for both.
- `dokimi conform --lang typescript` runs `dokimi run --check` on `corpus/typescript`, then `pnpm exec vitest run --reporter=json`, and judges the fixtures up to the level of the overlay.

### The TypeScript corpus

Every fixture of the TypeScript corpus whose Go form has a subject gains a subject at the same path: 132 fixtures, and the subjects of the ports. Each of the 11 controls of the expected outcomes gains a TypeScript control at the path of its Go form. The corpus gains a port of each fixture whose axis TypeScript has, and its overlay blocks each fixture whose axis is a construct of Go:

| Fixture | Level | TypeScript | Reason |
|---|---|---|---|
| `iface/lang/arraylength` | 1 | Blocked | An array whose type fixes its length through a constant, which TypeScript does not have |
| `iface/lang/crosspackage` | 1 | Blocked | The internal packages of Go |
| `iface/lang/dynamic` | 1 | Port | A parameter of the type `unknown` |
| `iface/lang/kinds` | 1 | Blocked | The kinds of `reflect` and the channels of Go |
| `iface/lang/namedreturns` | 1 | Blocked | The named results of Go |
| `iface/lang/nested` | 1 | Port | Fields of function types and of interface types |
| `iface/lang/opaque` | 1 | Blocked | A value that starts from the zero value of its type, which TypeScript does not have |
| `iface/lang/packagescope` | 1 | Port | A module that declares the names that its suite imports, so the suite imports them under aliases |
| `iface/lang/recursive` | 1 | Port | A recursive type |
| `iface/lang/refused` | 1 | Port | The directive `refuse` |
| `iface/lang/resolved` | 1 | Blocked | As `opaque` |
| `iface/lang/sampledparam` | 1 | Port | The directive `sample` on a parameter |
| `iface/lang/selfnamed` | 1 | Blocked | The field of an embedded interface |
| `iface/lang/suitedhandle` | 1 | Port | A declaration that a method of a double returns |
| `iface/mixin/aead` | 2 | Port | An interface that extends another, with the builders of the mixin `sample` on the methods that it inherits |
| `iface/mixin/clock` | 2 | Port | A time-aware read without a signal |
| `iface/mixin/durable` | 2 | Port | The mixin `durable` over the medium |
| `iface/mixin/returnsonerror` | 2 | Blocked | A failed call that returns a value beside its error |
| `iface/mixin/runsuntilcancelled` | 2 | Port | A call that runs until its signal aborts |
| `iface/mixin/source` | 2 | Port | A call with a side effect and without a signal |
| `iface/mixin/timer` | 2 | Port | A command with the signature of a query |
| `iface/mixin/waits` | 2 | Port | A call that waits for another call |
| `iface/mutation/deadcheck` | 2 | Port | Its suite. Its attribution of mutants needs a mutation record of TypeScript |
| `enum/split` | 3 | Blocked | The members of an enumeration in another file than its type. A TypeScript enum declares its members in one declaration |
| `errors/split` | 3 | Blocked | The errors of one package in three files. A TypeScript module is one file |
| `iface/door/balancer`, `breaker`, `copier`, `ledger` and `readthrough` | 4 | Port | Calls that fail on demand |
| `iface/door/durable` | 4 | Port | The doors `restart` and `crash` over a medium |
| `iface/door/hosted` | 4 | Port | The door `hosted`, with one subject in native mode and one in Tasks mode |
| `iface/door/queue` | 4 | Port | A check that waits for quiescence |
| `iface/door/reference` | 4 | Port | The door `reference` with the param `normalize` |
| `iface/door/rooted` | 4 | Port | A subject in a composition root |
| `iface/door/signer` | 4 | Port | The door of a value |
| `sim/ledger`, `sim/placed`, `sim/replicated` and `sim/waits` | 4 | Port | Systems of sans-IO machines |
| `sim/mirrored` | 4 | Port | Components on hosts, with a stream between them |

That is 31 ports and 10 blocked fixtures. The overlay judges up to level 4:

```yaml
language: typescript
level: 4
diverge:
  - fixture: iface/lang/arraylength
    stance: blocked
    why: A TypeScript array has no length in its type
  - fixture: iface/mixin/ttl
    key: mixed.read/zero-on-error
    expect: check
    actual: refused
    gate: returns-beside-error
    why: A TypeScript call that fails throws and returns no value
```

The overlay declares the 10 blocked fixtures, each refusal through `returns-beside-error`, and each other difference that the first run of the gate finds, with its gate or its stance and its reason.

### The parts of `@dokimi/assert` that the runtime calls

| What the runtime calls | Definition | In `@dokimi/assert` 0.2.0 |
|---|---|---|
| `Seat` with `observeCall`, `writeRecord`, `runSeed`, `cleanup`, the controlled clock and `doesReject` | 13.6.0 | Yes |
| `prop.forAll`, `prop.fuzz`, `prop.around`, `prop.onResult` and the registry of shapes | 13.6.0 | Yes |
| `stateful.steps`, `clients`, `tasks`, `pct`, `uniform`, `history.isLinearizable`, `isSerializable` and `hasSnapshotIsolation` | 13.6.0 | Yes |
| `prop.confirm` | 14.0.0 | No |
| `wait`, `clock`, `advance` and `maxReleases` of the scheduler, `prop.replays`, and the replays of a run result | 16.0.0 | No |
| `keep` of an action | 16.2.0 | No |
| The resumed searches of a growing history | None. assert-go declares `Resume` and `Checkpoint` beside the definition | No |

This design requires two changes outside dokimi:

- `@dokimi/assert` implements definition 16.2.0. assert-go v0.20.0, which dokimi-go requires, implements the same definition.
- The definition of assert-spec gains the resumed searches of a growing history. The run of sim and `probe.LinearizableHistory` of Go resume them through `Resume` and `Checkpoint`, which assert-go declares beside the definition.

### The verdict of each capability

| Capability | TypeScript | Form or reason |
|---|---|---|
| The checks of level 1 | In | The probes, the gates and the refusals of level 1 of the TypeScript binding, over a copy of spec 2.1.0 |
| Generated inputs, fuzz targets, benchmarks and budgets | In | Shape definitions, `prop.fuzz` with Jazzer.js, and `bench` of Vitest with the latency contract of `@dokimi/assert` |
| The value of every input type | In | The call generator of TypeScript gives each parameter a value by its kind |
| The doors of a subject and of a planted defect | In | `Doors<S>`, whose doors may return a promise |
| The options of the property run of a check | In | `options` of `Check<S>` |
| The benchmarks of a callable that waits | In | `waits` of `Call<S>` |
| The environment of a check | In | `around` of `Check<S>` |
| The sample of a parameter | In | The directive `sample` on a parameter |
| Vacuous checks and controls | In | `engage` and `observeCall` of the handle |
| The suite record in the artifact of a run | In | The meta of the last test of a block |
| The corpus gate of each language | In | `dokimi conform --lang typescript` up to level 4 |
| Calls that fail on demand | In | The injector of `host/injector`. A hold blocks an asynchronous call alone |
| Generated sequences and the judges of their histories | In | `stateful.steps` and the checkers of `history` |
| Concurrent callers and the consistency of their histories | In | The clients run as tasks of a scheduler under `uniform`, on a subject and on a planted defect, because `@dokimi/assert` runs a concurrent section only as tasks of its task scheduler. They interleave at the yields of the scheduler, and `h.yield()` yields to it |
| Verified fakes | In | A class over a `Map` or an array, which the fake generator writes |
| A breakable medium | In | `medium/*`, and the mixin `durable` binds its claims |
| Dead and subsumed checks from mutation records | In for the tests | Each Vitest test of a check has the identity of the check. mutate-spec defines the mutants of Go alone, so the report needs an overlay of TypeScript in mutate-spec and an engine that writes its records |
| A deterministic simulation of sans-IO machines | In | `sim/*` and the TypeScript tree of the sim generator |
| The workloads of a simulation on a real cluster | None | Its design is a Draft |
| The probes of each claim | In, except six | 121 modules. `zero-beside-error`, `zero-pair-beside-error` and `field-beside-error` have none, because the gate `returns-beside-error` refuses their claims for every callable. `encodes-each-variant`, `decodes-each-text` and `compares-with-identity` have none, because the gates `text-encoding` and `self-comparing-type` refuse their claims for every TypeScript declaration |
| Claims of a single caller for the classifications without one | In | The drain of a stream reads an iterable or an async iterable |
| Calls that run until their context ends | In | `runsUntilCancelled` of `Call<S>` |
| The doors of a value | In | `supply(token, open)`, `supplies`, `values` and `h.open(token, s)` |
| Judges of deliveries, iterations, totals and failed calls | In | `Pull<E>`, async iterables, and `quiescence` in place of `synctest.Wait` |
| Derived claims that a declaration narrows | In, except `returns-on-error` | The gate `returns-beside-error` refuses `returns-on-error` |
| The isolation levels of transactions | In | `history.isSerializable` and `history.hasSnapshotIsolation` |
| The runtime contract and the rules of the generator | In | The copy from dokimasia/dokimi-spec and the runtime model. The runtime has no validator of records |
| The repeated parts of the generated code | In | The reduced form of the TypeScript trees, and the behaviour suites of the doubles of Go |
| A simulation of components | In | Native mode on the event loop and Tasks mode on the scheduler, without `blocking` |
| Hosts, started machines and the inputs of a system | In | `addHost`, `Host.spawn` and `Launch` |
| The end of a wait | In | `cancel` of an outcome, and `ends` and `end` of the anatomy |
| Hosts that exit and the steps of a case | In | `meanSteps` and `maxSteps` |
| Performers of a pointer | None | A TypeScript object is a reference |
| The deliveries and the choices of the active phase | In | `keep` of the actions of `stateful`, from definition 16.2.0 |
| The final phase of a long convergence and the timers of a performer | In | `afterFunc` of the world |
| The outcome of an effect at a loss after the effect | In | `loseAfter` of the world |
| The restarts of a host in the final phase | In | The delays of the restarts of the Go run |
| The errors of a component that leave the outcome of a call unknown | In | `unknown` of the anatomy, and the reason of an aborted signal |
| The checks of a subject that only a host opens | In | Each check of `hosted` opens its subject through the door `host`, inside `fakeTimers` |
| A change of the consumer in the final phase | In | `drain` of an action |
| The planted defect of a convergence | In | `defect` of a convergence |
| The steps of a scenario of a world of components | In | A call of `invoke` in a scenario is a step with two crash points |
| An environment of the consumer around each case of a simulation | In, without `serial` | `around` of a simulation |
| Allocation ceilings and the checks that count allocations | Out | V8 counts no allocations of a call |
| Parallel benchmarks | Out | One thread runs the calls of a benchmark |
| The race detector | Out | One thread runs the code of a test, so no two calls access memory at once |
| The watchdog of a bubble | Out | A component that starves the event loop starves each callback of a watchdog too |

The reason of each verdict Out is a mechanism that the platform lacks, and the verdict reverses when the platform gains it.

## Alternatives considered

### A package for each Go module

Each part of the runtime is a package, such as `@dokimi/conform`, `@dokimi/host` and `@dokimi/sim`. Each package has its releases, as each Go module has.

**Why not:** each TypeScript module of the runtime has the peer dependencies `@dokimi/assert` and Vitest, so a split does not isolate a dependency. A generated suite imports conform, probe, spec, stub and fault at once, so a consumer would align five versions for each release. One package with subpath exports keeps each module behind a separate import path.

### Tasks mode alone

The runtime runs every component as tasks of the scheduler, which explores the interleavings at each wait of the host.

**Why not:** a task must not await a promise that only another task settles. A component whose two tasks pass work through a promise that the component creates, such as a queue between a reader and a writer, cannot run under the scheduler. The event loop runs such a component as production runs it, and it replays its case whole.

### Real sockets and files of Node in the host

A component takes `node:net` and `node:fs` through its environment, and the world strikes faults through a proxy.

**Why not:** real I/O completes in the poll phase at a time that the operating system decides. The host counts no operation for a pending real read, so the quiescence wait returns while the read is pending. A replay then differs from its run. A partition, a latency and a torn write need the simulated network and medium anyway.

### Hooks of every promise of an incarnation

`createHook` of `node:async_hooks` records each promise that the code of an incarnation creates, for the quiescence wait and for the leaked tasks.

**Why not:** Node marks the API experimental and does not recommend `createHook`, because of its usability issues, safety risks and performance implications. The tasks that a host starts cover the leaked tasks, and the quiescence wait counts the operations of the host instead of single promises.

### Real time for the code of a component

The world does not fake a timer. A component reads time only through its host.

**Why not:** a timer of a component on the real clock runs its callback after the quiescence wait has returned, at a time that the world does not control. A run then differs from its replay. A goroutine in a bubble sees the clock of the bubble, and the fake timers give a TypeScript component the same view of time.

### Synchronous doors of the handle

`subject()`, `fresh()`, `restart(s)` and the other doors return the subject at once. The TypeScript binding of level 1 has them in this form.

**Why not:** a subject that recovers from its medium reads files. A constructor cannot await a read, so such a subject would defer its recovery to its first call. `subject()` remains synchronous, because the handle builds the subject before the check runs.

### A validator of suite records in the runtime

`validateSuiteRecord` checks each record against the schema of the copy, through a validator of JSON Schema such as Ajv.

**Why not:** the runtime does not validate a record while it runs, and Go validates in `dokimi report` alone. The gate validates each record of the TypeScript corpus with that validator. A validator in the runtime would add a third-party package and no check.

### The injector in conform

The injector keeps the module `conform/injector` of its TypeScript names.

**Why not:** the host fails, blocks and releases the calls of the dependencies of an incarnation through the injector, and the host must not import conform, which imports the host. Go keeps its injector in its host module, and conform declares aliases of it.

### The report of StrykerJS in `dokimi report`

`dokimi report` reads the JSON report of StrykerJS, the mutation engine of JavaScript, which lists the tests that kill each mutant.

**Why not:** `dokimi report` attributes the mutants of a record of mutate-spec, with the fields `allKillers`, `killers: all` and `coverage: leaves`. A second format needs a second attribution of the same findings, and its mutants differ from the mutants that mutate-spec defines for every engine.

## Drawbacks

- **A second runtime.** The Go runtime has 20,776 lines of source in 100 files and 37,397 lines of specs. The TypeScript runtime has 3,736 lines of source today, and it mirrors each module of Go with 121 modules of probes.
- **Template trees.** The Go trees have 2,608 lines for the suite generator in three templates, 96 for the fake generator and 343 for the sim generator. The TypeScript stub tree of 940 lines changes for the behaviour suites.
- **Every new machinery twice.** Each machinery that a design adds needs both runtimes and both template trees.
- **The corpus.** The TypeScript corpus needs 132 subjects, 11 controls and 31 ports before the gate judges it at level 4.
- **Fewer checks than in Go.** TypeScript lacks allocation ceilings, parallel benchmarks, the race detector and the watchdog. It refuses three claims of a value beside an error and three claims of the text and the identity of a value, and `leaked` covers the tasks that the host started.
- **Global fake timers.** The fake timers of a world apply to the whole worker of Vitest. A test of the consumer that runs concurrently in the same file sees them.
- **Asynchronous doors.** The generated code awaits each door of the handle except `subject()`, where Go calls it.
- **Explicit yields in Tasks mode.** A component in Tasks mode marks each point where another task may run.
- **A coverage gate alone.** The specs of the TypeScript runtime run under a coverage threshold of 100%. mutate-spec defines the mutants of Go alone, so the specs have no mutation score.
- **Releases before the gate.** The gate of dokimi runs the TypeScript corpus only against releases of both packages on npm, and it runs Node and pnpm beside Go.

## Open questions

None.

## Unresolved and future work

- The workloads of a simulation on a real cluster in TypeScript are not proposed here.
- A runtime for Deno, Bun or a browser is not proposed here. The runtime needs Node 26 or later.

## References

| What | Where |
|---|---|
| Vitest, `vi.useFakeTimers` and the advance of fake timers | <https://vitest.dev/api/vi.html> |
| Vitest, the meta of a test and its reporters | <https://vitest.dev/api/advanced/metadata> |
| Node.js, the event loop, timers and `process.nextTick` | <https://nodejs.org/en/learn/asynchronous-work/event-loop-timers-and-nexttick> |
| Node.js, `queueMicrotask` and `process.nextTick` | <https://nodejs.org/api/process.html#when-to-use-queuemicrotask-vs-processnexttick> |
| Node.js, `async_hooks` and its stability | <https://nodejs.org/api/async_hooks.html> |
| WinterTC, the Sockets API, Draft of 31 July 2026 | <https://sockets-api.proposal.wintertc.org/> |
| @sinonjs/fake-timers | <https://github.com/sinonjs/fake-timers> |
| Angular, `InjectionToken` | <https://angular.dev/api/core/InjectionToken> |
| Jazzer.js | <https://github.com/CodeIntelligenceTesting/jazzer.js> |
| assert-typescript | <https://github.com/dokimasia/assert-typescript> |
| assert-spec, the naming table | <https://github.com/dokimasia/assert-spec/blob/main/spec/naming.yaml> |
| dokimi-spec | <https://github.com/dokimasia/dokimi-spec> |
| dokimi-go | <https://github.com/dokimasia/dokimi-go> |
| mutate-spec | <https://github.com/dokimasia/mutate-spec> |
