---
rfc: 0004
title: Doubles, builders, enumerations and errors
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-19
updated: 2026-10-22
discussion: none
supersedes: none
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0004: Doubles, builders, enumerations and errors

## Summary

Four generators of dokimi write the parts of a test that are not suites:

- `stub` writes a typed double of a declaration, which records calls, returns configured results and injects faults.
- `builder` writes a fixture builder of a data declaration, which starts from the declared defaults of the declaration.
- `enum` writes the text surface of an enumeration and the checks of the surface.
- `sentinel` writes the checks of the error values and error types of a package.

The checks of `enum` and `sentinel` cover the checks that the enumeration tests and the error tests of testkit write, and they assert through the assert library of each language. A double and a builder are tools that tests use, so a structural test of the generated code covers them.

A double and a builder work together. Where a test configures nothing, a double returns a new value from the builder of each result whose type has one. A test then configures only the calls that it checks.

## Motivation

A suite builds its planted defects from a double, and consumers test the callers of an interface against a double. A double written by hand for each interface records calls in its own way and cannot inject a failure in the way that a generated check expects. A generated double gives every planted defect and every consumer the same calls and results in each language.

A builder gives the tests of a data declaration a value that starts from its declared defaults. Each setter of the builder changes one field, so a test states only the fields that it depends on.

Enumerations and error values have contracts that recur in every codebase:

- The members of an enumeration have distinct values and distinct names.
- A name parses back to its member, and an unknown name does not parse.
- An error value is set, and its message does not start with the message of another error value.
- Two error values do not match each other, and an error type exposes the cause that it was given.

testkit generates doubles, builders, enumeration tests and error tests for Go. The reference generators of eidos write builders, enumeration methods and error tests for Go alone, and their tests assert through the package `testing`.

## Detailed design

### Directives

| Directive | Generator | Written on | Meaning |
|---|---|---|---|
| `stub` | `stub` | A declaration | Generate the typed double of the declaration |
| `fault` | The annotator `fault`, read by `stub` | A method of a declaration under `stub` | Generate the helpers of the stub of the method that arm the listed errors, a retry schedule and the faults of a partition |
| `builder` | `builder` | A struct, a class or an alias of one | Generate the fixture builder of the declaration. The param `defaults` names a function that seeds the value |
| `enum` | `enum` | An enumeration | Generate the text surface of the enumeration and its checks. The param `methods=off` leaves out the surface, and `text=comment` takes the text of each variant from its trailing comment |
| `text` | The annotator `text`, read by `enum` | A variant of an enumeration under `enum` | Pin the text of the variant, which the surface renders and parses |
| `sentinel` | `sentinel` | A Go package or a TypeScript module | Generate the checks of the errors of the package. The param `prefix` states the opening of each message, and `no-overlap-with` lists other packages |

An interface or a concrete type with `suite` also gets its double, because its planted defects are built from it. The double of a concrete type implements the interface of its method set that the suite generator writes. Each language declares the subject of a generator in a form of its own:

| Generator | Go | TypeScript | Rust | Python | Protobuf |
|---|---|---|---|---|---|
| `stub` | An interface, or a struct with methods | An interface, or a class | A trait, or a struct with methods | A `Protocol`, an abstract class or a class | A service |
| `builder` | A struct | A class | A struct | A dataclass | A message |
| `enum` | A type with constants | An enum | An enum without fields | An `Enum` subclass | An enum |
| `sentinel` | An error value or an error type | An error class | A type that implements `std::error::Error` | A subclass of `Exception` | None |

Each generator is a plugin of eidos with a template tree for each target:

```go
// Package stub generates the typed double of each declaration that the
// directive stub or the directive suite marks.
package stub

// New returns the generator of typed doubles. It provides the capability
// stub, which the suite generator requires.
func New() plugin.Generator
```

The packages `builder`, `enum` and `sentinel` each declare the same constructor.

### The package of the checks

Each generator writes its checks against the exported API of the package of its subject, as a caller in another package uses it. The checks of Go are in the external test package of their directory, such as `catalog_test` beside `catalog`, through `plugin.Output.External` of eidos. A TypeScript test imports what it checks from the module beside it.

The record of each generator declares each function, type and variable that its Go template writes beside the subject, such as `NewAddressBuilder` and `ParseStatus`. The layout of eidos checks their names in the package of the record. It also qualifies each reference from the checks to them, so a check calls `catalog.NewAddressBuilder`. The qualification follows an `out` that moves both files of a subject to another directory, as in `cachetest.NewCacheStub`.

The checks of an unexported declaration cannot name it from another package, so a file with such checks keeps the package of its directory. A generator marks that file through the mechanism that dokimasia/eidos#45 requests.

### The double

The callables of a double are the methods of the member walk of its subject. The walk adds the methods of each interface that an interface embeds or extends, the methods that an embed of a Go struct promotes, and the methods of each base of a TypeScript class. The double leaves out each member at type level, such as a static method, and the constructors of a class. It declares each TypeScript property as a public field, which a test assigns, because a property is not a callable.

For each callable, a double records every call with its arguments and its results. It returns the results of the first of these arms that applies:

1. A fault that the test armed, where the fault fails the call.
2. A function that the test set, which computes the results from the arguments.
3. A result that the test configured.
4. The default results. A result whose type has a builder takes a new value from the builder, as [Doubles and builders](#doubles-and-builders) states. Every other result is its zero value in Go, and `undefined` in TypeScript.

A callable can fail where the error model of its language has a failure channel. A Go method can fail when its last result is `error`. Every TypeScript method can fail, because it can throw or reject. A fault of a TypeScript method throws its error, or rejects with it. A fault that fails a call of a Go method without an error result fails the test, because a double that dropped the fault would let a fault-injection test pass without a fault.

#### The runtime

The runtime of each language implements the parts of a double that do not depend on a signature:

| Part | Go: `go.dokimi.dev/stub` and `go.dokimi.dev/fault` | TypeScript: `@dokimi/runtime` |
|---|---|---|
| The log of calls with the instant of each, its hooks and bench mode | `stub.Recorder[T]` | `Recorder<T>` of `stub/recorder` |
| The stub of one method: faults, latency, strict mode and the expected counts of calls | `stub.MethodStub[C]` | `MethodStub<C>` of `stub/method-stub` |
| The results of one call, from the arms in order | `stub.Answer` with `stub.Arms[C, R]` | `answer`, and `answerAsync` for a method that returns a promise, with `Arms<C, R, I>` of `stub/answer` |
| A setting of the whole double | `stub.Group` over `stub.Configurable` | `Group` over `Configurable` of `stub/group` |
| The order of the calls of two methods | `stub.OrderTracker` | `OrderTracker` of `stub/order-tracker` |
| The contract of each method | `stub.Behaviour` | `behaviour` of `stub/behaviour` |
| The contract of each option, each constructor and the reset of the double | `stub.BenchModeBehaviour`, `StrictBehaviour`, `ResetBehaviour`, `RandSourceBehaviour`, `ClockBehaviour` and `LatencyBehaviour` | `benchModeBehaviour`, `strictBehaviour`, `resetBehaviour`, `randSourceBehaviour`, `clockBehaviour` and `latencyBehaviour` of `stub/double-behaviour` |
| The fault strategies | The package `fault` of the module `go.dokimi.dev/fault` | The module `fault` |

Each module of TypeScript is a subpath export of the package, such as `@dokimi/runtime/stub/answer`.

A double reads time from the clock of the assert library of its language. Without a clock of the test, it reads the clock of the system. With a controlled clock, a latency or the window of a fault costs no real time. A latency applies to a TypeScript method that returns a promise, because a synchronous method cannot wait.

A fault with a probability takes its random numbers from a source that the test passes. A seeded source repeats the faults of a run. Go takes a `rand.Source`. TypeScript takes a function that returns a number in [0, 1), with `Math.random` as its default.

A test arms a fault through the fault strategies:

| Strategy | Go | TypeScript | The calls that the fault fails |
|---|---|---|---|
| Counted | `fault.NewCountedFault` | `counted(error, n)` | Every nth call |
| Retry | `fault.NewRetryFault` | `retry(error, n)` | The calls 1 to n − 1, and no call from call n |
| Probability | `fault.NewProbabilityFault` | `probability(error, p, random)` | Each call with a probability, from the random source of the double |
| Window | `fault.NewWindowedFault` | `windowed(error, deadline)` | Each call before a deadline on the clock of the double |
| Predicate | `fault.NewPredicateFault` | `predicate(error, select)` | Each call for which a function of the call returns true |
| And, Or | `fault.And` and `fault.Or` | `and(...faults)` and `or(...faults)` | Each call that each of the faults fails, or one of them |

A function of the module `fault` of TypeScript returns a fault, which `MethodStub.setFault` arms.

A recorder of Go has a gate, which blocks each recording of a goroutine until the test releases it. The runtime of TypeScript has no gate. JavaScript runs one thread, so no two recordings of a recorder interleave. A test delays the results of an asynchronous method through a function that returns a promise, which the test settles.

#### The generated double

The templates of each language write a double in the conventions of the language. For the method `Get` of the interface `Store`:

| Operation | Go | TypeScript |
|---|---|---|
| The double | `type StoreStub struct`, with a field of each method | `export class StoreStub implements Store`, with a property of each method |
| The stub of a method | The field `OnGet *StoreGetStub`, which embeds `*stub.MethodStub[StoreGetCall]` | The property `onGet: StoreGetStub`, which extends `MethodStub<StoreGetCall>` |
| The call | `type StoreGetCall struct`, with a field of each parameter and of each result | `export interface StoreGetCall`, with a property of each parameter, `result` and `error` |
| The fixed results | `type StoreGetReturn struct`, with a field of each result | The result type of the method |
| Construct | `NewStoreStub(tb, opts...)` | `new StoreStub(seat, ...options)` |
| Construct with a latency on a clock | `NewStoreStubDelayed(tb, d, clk, opts...)` | `StoreStub.newDelayed(seat, ms, clock, ...options)` |
| Construct with each method that can fail failing | `NewStoreStubFailing(tb, err, opts...)` | `StoreStub.newFailing(seat, error, ...options)` |
| The options | `StoreStubStrict()`, `StoreStubDelegateTo(impl)`, `StoreStubWithClock(clk)`, `StoreStubWithRandSource(src)`, `StoreStubBenchMode()`, `StoreStubWithLatency(d)`, `StoreStubFailing(err)`, and `WithStoreGet(fn)` of each method | `StoreStub.strict()`, `StoreStub.delegateTo(impl)`, `StoreStub.withClock(clock)`, `StoreStub.withRandSource(random)`, `StoreStub.benchMode()`, `StoreStub.withLatency(ms)`, `StoreStub.failing(error)`, and `StoreStub.withStoreGet(fn)` of each method |
| Configure the results | `OnGet.Returns(item, nil)` | `onGet.returns(item)` |
| Set a function | `OnGet.Func(fn)` | `onGet.func(fn)` |
| Return a sequence | `OnNames.Yields(items...)`, and `OnRows.YieldsError(items, err)` for a sequence of pairs with an error | `onRows.yields(...items)` for a method that returns an asynchronous iterable |
| Arm a fault | `OnGet.Faults(err, n)`, `FaultsOnce(err)`, `FaultsWhen(pred, err, n)`, `FaultsWithProbability(p, err)`, `FaultsUntil(deadline, err)`, `FaultsFor(d, err)` or `SetFault(f)` | `onGet.faults(error, n)`, `faultsOnce(error)`, `faultsWhen(pred, error, n)`, `faultsWithProbability(p, error)`, `faultsUntil(deadline, error)`, `faultsFor(ms, error)` or `setFault(f)` |
| Strict mode | `OnGet.Strict()` | `onGet.strict()` |
| Expect a count of calls | `OnGet.Times(n)` and `OnGet.TimesAtLeast(n)` | `onGet.times(n)` and `onGet.timesAtLeast(n)` |
| A latency | `OnGet.Latency(d)` | `onGet.latency(ms)` |
| The recorded calls | `OnGet.Calls()`, `OnGet.Timestamped()` and `OnGet.OnRecord(fn)` | `onGet.calls`, `onGet.timestamped` and `onGet.onRecord(fn)` |
| Check the recorded calls | `OnGet.AssertCalledOnce(t, msg)`, `AssertCalledN` and `AssertNotCalled` | `onGet.assertCalledOnce(seat, msg)`, `assertCalledN` and `assertNotCalled` |
| Clear the calls of every method | `ResetCalls()` | `resetCalls()` |
| A property `zone` | None | The field `zone!: string` |

A field of a call has the name of its parameter in the case of the language. A parameter without a name, or with the name `_`, takes `a` followed by its position, as `naming.Params` of eidos names it. A field of the results of Go has the name of its result. A result without a name is `Err` for the error and `Result` for one other result. Where a method has two or more other results, they are `Result0`, `Result1` and onward. A method of a Go double declares the names of its results where each result has a name that no parameter has. A TypeScript call records the result in the field `result` and a thrown error in the field `error`.

The constructor registers the check of the expected counts of each method for the end of the test. Go registers it through the method `Cleanup` of the TB. TypeScript registers it through the method `cleanup` of the `Collector` that the Vitest adapter of the assert library passes to each test. Under the option `Strict`, a double refuses each unconfigured call. The option `DelegateTo` passes each call to an implementation, and the double records the call. The implementation of the double of a Go struct has the type of the interface of the methods of the double. Both a pointer to the struct and a second double implement that interface. The option `Failing` arms each method that can fail with one error.

A method of a Go double with the mixin `deprecated` of the shape catalog logs through the TB that it is deprecated. A seat of the assert library of TypeScript has no log, so a TypeScript double does not log the mixin. A method with the mixin `order-after` refuses a call before its sibling in a strict double, through the order tracker of the double.

A method of the subject that returns an interface of the package of the subject without the directive `stub` returns a handle. The double of the first subject in the package whose method returns the handle also doubles the handle, in identity order. A handle without a member has no double. The layout of eidos routes a declaration without routing of its own with the subject that placed it, so an `out` on the subject also moves the double of the handle.

#### The directive fault

The directive `fault` on a method of a declaration under `stub` lists the errors that the tests of the method arm, a retry schedule, and the parameter that selects a partition:

```go
//dokimi:fault errors=[ErrNotFound] retry=3 partition=key
Get(ctx context.Context, key string) (Value, error)
```

| Param | Value | Helper of the stub of the method |
|---|---|---|
| `errors` | A list of error values of the package | `FaultNotFound()` of `ErrNotFound`, which fails the next call with the error |
| `retry` | The attempt that succeeds, 2 or more | `RetrySchedule(err)`, which fails each call before the attempt `retry` with err |
| `partition` | A parameter of the method | `FaultForPartition(key, err, n)`, which fails every nth call of the partition key, and `FaultForOtherPartitions(key, err, n)`, which fails every nth call of each other partition |

The helpers of TypeScript are `faultNotFound()`, `retrySchedule(error)`, `faultForPartition(key, error, n)` and `faultForOtherPartitions(key, error, n)`. The annotator `fault` of dokimi stamps the params of each instance on its method, and the stub generator reads them. eidos resolves each name of `errors` to a variable of the package and `partition` to a parameter of the method. The annotator stamps the variables of `errors` as one fact, a list of identities in the order of the param. The name of a helper joins `Fault` and the name of the error value without the opening of an error value, such as `FaultNotFound` for `ErrNotFound`. A name that the rules of the language do not read as the name of an error value joins whole, such as `faultNotFound` for `notFound`.

#### The test of a double

The stub generator writes a test of each double with a method, and a test of the stub of each method with a parameter or a result. The test of the double, such as `TestStoreStub`, has a level for each method and each option of the double. It also has the levels `ResetCalls`, `NewStoreStubDelayed` and `NewStoreStubFailing`. The test of the stub of a method, such as `TestStoreGetStub`, has a level for each method of the stub. A TypeScript test nests the same levels in `describe` blocks. The levels run the behaviour suites of the runtime against the double:

| Suite | Level | Checks |
|---|---|---|
| `Behaviour` | Each method | A call is recorded. An unmet count of `Times` and of `TimesAtLeast` fails the test. Strict mode refuses an unconfigured call. An unconfigured call returns the default results. The double calls the function that the test set. An armed fault fails the call, a counted fault fails the nth call, and a call that a fault failed is recorded |
| `ResetBehaviour` | `ResetCalls` | A reset clears the instants of the calls |
| `LatencyBehaviour` | `NewStoreStubDelayed` and `StoreStubWithLatency` | A latency sleeps on the clock of the double. Under the latency of `StoreStub.newDelayed`, a synchronous TypeScript method records its call at once |
| `StrictBehaviour` | `StoreStubStrict` | A strict double refuses an unconfigured call |
| `ClockBehaviour` | `StoreStubWithClock` | A latency sleeps on the clock of the double, and a call is recorded after its latency. A fault with a window stops when the window closes |
| `RandSourceBehaviour` | `StoreStubWithRandSource` | A fault with a probability reads the source of the double |
| `BenchModeBehaviour` | `StoreStubBenchMode` | Bench mode records no call |

Each level also checks what needs a value of the signature of a method:

- **A method of the double.** The log of calls contains the arguments, and OnRecord runs once for each call.
- **`ResetCalls`.** A reset clears the calls of each method, and keeps the results of Returns.
- **`NewStoreStubFailing` and `StoreStubFailing`.** Each method that can fail fails with the error.
- **`StoreStubStrict`.** A strict double refuses a call before its sibling.
- **`StoreStubDelegateTo`.** Each method passes its call to the implementation, and fails with the error of the implementation.
- **The option of each method, such as `WithStoreGet`.** The option sets the function of the method.
- **`Returns` of the stub of a method.** The method returns the samples of Returns, after a check that its samples differ from the default results, because a check with equal values passes for a double that ignores Returns. A method with the results of Returns passes strict mode.
- **`Func` of the stub of a method.** The function receives the arguments of each call.
- **The other methods of the stub of a method.** Further levels cover a sequence, the helper of each error value, the retry schedule and the two partition faults. The helper of an error value fails the next call with the value, and the call after it succeeds.

The arguments take a sample of each parameter, and Returns takes a sample of each result. A variadic parameter takes no sample. In Go the error result of Returns is a new error, which the check matches. The test passes its own context to a parameter that the rules of the language give the role of a context. A helper of each method returns the samples of its arguments, and a second helper calls the method with them, so each check calls the method in one line. Each check of `DelegateTo` builds a double, and a second double as its implementation. The test of a generic double runs at the witnesses of its type parameters. Without witnesses, the test states that the double has no checks.

A double without methods, such as the TypeScript double of an interface with properties alone, does not run a suite, because each suite calls a method of the double. Its test has the level `resetCalls`, with a case for `new`, for each option and for `newDelayed`. Each case builds the double, assigns a sample to each property and then calls `resetCalls`. The case checks that each property keeps its sample.

The tests of the runtime prove each check of its suites with a double that breaks the check. In TypeScript, a suite returns its checks as data. The generated test registers each check with `it` of vitest and runs it on a `Standard` seat of the assert library. A `Standard` seat throws at the first failure. The other cases of the test build each double without a seat. They assert through the assertions of the assert library on the `expect` of vitest. The setup file of the project registers those assertions through `register` of `@dokimi/assert/vitest`.

### The builder

The subject of a builder is a struct or a class under the directive `builder`, or an alias of one without type parameters. The builder of a generic struct builds each instance of the struct, so an alias of an instance has no builder of its own. The plugin derives a model of each builder from the facts of eidos, and records it. The templates of each language write the model:

| Part of the model | Derived from |
|---|---|
| The members | The fields of the member walk that code outside the package can assign on a value: public fields of the instance that are not immutable. A field with a hard private name is private. An embedded field of a public type is a member that the builder sets whole, and a field that an embed promotes is no member. A field that a base class declares is a member. A field with the tag `builder:"-"` is no member. A builder assigns a field after a dot, so a TypeScript property whose name is not an identifier, such as `'content-type'` or the computed name `[Symbol.toStringTag]`, is no member |
| The names of the methods | The part of each language names the constructors and the setters of a member from the name of the builder, the name of the field and its form |
| The form of a member | The canonical shape of the type of the field: a list, a map, bytes, an optional or a scalar. A field that may be absent is an optional of its type. The type follows each transparent alias. A Go defined type, such as `type Codes []string`, is a reference to its declaration, so the builder sets it whole. The shapes have no set, so the part of each language states the types that it writes as sets |
| The seed | The function that the param `defaults` names, else the companion function of the declaration, else the start value of the language |
| The defaults | The default that the annotator `defaults` stamps on each member, a member of a base included. The rules of the language read the text of a default as a value in the scope of the declaration. For a field of text, a text that the rules do not read as a name, a call or a string is the content of the string |
| The copies | Each list, map, set and bytes value of a member, which `Clone` copies |

The part of each language states the companion function, the start value and the sets:

| Language | Companion function of `Address` | Start value | Set |
|---|---|---|---|
| Go | `AddressDefaults()` | The zero value | A map whose values are `struct{}` |
| TypeScript | `addressDefaults()` | `new Address()` | `Set<T>` |

A builder has the same operations in every language, and the templates of each language write them in the conventions of the language:

| Operation | Go | TypeScript |
|---|---|---|
| The builder of `Address` | `type AddressBuilder struct` | `export class AddressBuilder` |
| Start from the seed and the declared defaults | `NewAddressBuilder() *AddressBuilder` | `AddressBuilder.create(): AddressBuilder` |
| Start from a value | `NewAddressBuilderFrom(v Address) *AddressBuilder` | `AddressBuilder.from(value: Address): AddressBuilder` |
| Set a field | `WithStreet(v string)` | `withStreet(value: string): this` |
| Replace a list | `WithLines(v ...string)` | `withLines(...values: string[]): this` |
| Append to a list | `AppendLines(v ...string)` | `appendLines(...values: string[]): this` |
| Set an entry of a map | `WithLabelsEntry(k, v string)` | `withLabelsEntry(key: string, value: string): this` |
| Copy the entries of another map | `WithLabelsEntries(entries map[string]string)` | `withLabelsEntries(entries: Map<string, string>): this` |
| Add keys to a set | `WithTagsEntry(k string)` and `WithTagsEntries(keys ...string)` | `withTagsEntry(key: string)` and `withTagsEntries(...keys: string[])` |
| Set bytes from text | `WithRawString(s string)` | None. A TypeScript caller encodes text with `TextEncoder` |
| Set an optional field | `WithFloor(v int)`, which points the field at `v` | `withFloor(value: number): this` |
| Change the value in a function | `Mutate(fn func(*Address))` | `mutate(fn: (value: Address) => void): this` |
| Copy the builder | `Clone() *AddressBuilder`, which copies each list, map, set and bytes value | `clone(): AddressBuilder`, with the same copies |
| Build | `Build() Address` | `build(): Address`, a shallow copy of the value |

A Go builder named `AddressBuilder` has constructors whose names open with `New` and the name of the builder, so they never clash with a constructor of `Address` in the same package. A TypeScript builder keeps its value in a private property, and it constructs a class with `new` and no arguments. A class whose public constructors all take a required argument refuses under `Unconstructed`.

The builder has three options:

| Option | Default | Meaning |
|---|---|---|
| `suffix` | `builder` | The word that the name of each builder joins to the name of its declaration |
| `companion-word` | `defaults` | The word that the name of a companion function joins to the name of its declaration |
| `no-tests` | false | Leaves out the checks of every builder |

The family of the tests writes the test of each builder, such as `TestAddressBuilder`, which asserts through the assert library of its language. The test has a level for each constructor, each setter, each probe, `Mutate`, `Clone` and `Build`:

| Level | Checks |
|---|---|
| `NewAddressBuilder` | The constructor sets each declared default other than the zero value of its type |
| `NewAddressBuilderFrom` | `From` keeps each scalar member of a seed that sets a sample of each. It keeps the value that it gets |
| Each setter, such as `WithStreet` | The setter writes its member, which then equals the sample, and an optional points at the sample. A setter called twice keeps the second value |
| Each probe, such as `AppendLines` | A probe calls each method that the checks of the setters leave out: the appender of a list, the setters of the entries of a map and of the keys of a set, the setter of bytes from text, and the setter of a member without a sample |
| `Mutate` | `Mutate` passes one value to each function |
| `Clone` | A clone is a builder of its own, and it builds the value of its original. Each list, map, set and bytes member of a clone is not shared with the original. The check writes the alternate sample into the member of the clone in place, and the member of the original keeps the sample |
| `Build` | `Build` returns the value that `Mutate` sets, where the seed sets a sample of a member. A TypeScript build returns a new object for each call |

The level of a constructor exists where the builder has a declared default. In TypeScript the levels are `create`, `from`, `withStreet`, `mutate`, `clone` and `build`. A TypeScript test declares the cases of each level in a function of its own, such as `addressBuilderWithStreetCases`, which `describe("withStreet")` runs.

A check of a setter needs a sample of the argument of the setter alone, so a field of the type `*int` has a check through `int`. A check of a replacement and a check that a clone does not share a member also need an alternate sample. An embedded field and a member without a sample have a probe of their setter instead.

The probe of an appender passes the sample and then the alternate of its member, or the sample twice without an alternate, and expects the elements of both calls. In Go, the probe of the entries of a map with an alternate copies the sample and then the alternate into a nil map, and expects both. Every other Go probe, the probe of an appender of a member without a sample included, passes the zero value of each part as a literal, such as `""` and `nil`, and compares the whole member. The probe of a scalar without a sample checks that the setter returns its builder. In TypeScript, the other probes pass the sample of their member, and the probe of a member without a sample checks that the setter returns its builder. The Go check of `Clone` sets each list, map, set and bytes member before it clones the builder. The TypeScript checks of a clone set the sample of each such member and compare the member of the clone with it.

A generic declaration is checked at the witnesses of its type parameters. Without witnesses, the test file states that the declaration has no checks.

### Doubles and builders

Where a test configured nothing, a double returns a value from a builder. The test then configures only the calls that it checks. Every other call of the code under test receives a value that the declaration of its type accepts, such as the account of a companion function that a store validates.

| Result of a method | Default result |
|---|---|
| A struct, a class or an alias under `builder`, such as `Account` | A new value from its builder: `NewAccountBuilder().Build()` in Go, and `AccountBuilder.create().build()` in TypeScript. A generic struct takes the type arguments of the result |
| A Go pointer to such a struct, or a TypeScript `Account \| undefined` | The same new value. Go returns its address |
| A TypeScript promise of such a type | A promise of the same new value |
| Every other result, such as an error, a flag, a list or a type without a builder | Its zero value in Go, and `undefined` in TypeScript |

Each call builds a new value, so two calls never share a value that a caller changes. A flag keeps its zero value, so a lookup with the results `(Account, bool)` reports a miss until the test configures a hit. Strict mode still refuses an unconfigured call. A function, a configured result and a fault come before the default results.

The builder states its constructor as a fact, and the stub writes the default results from it:

1. The builder stamps two facts on each struct, class and alias under its directive, in an annotate role beside its generate role. `builder.name` is the name of its builder, and `builder.seed` is the name of the constructor from the seed: the function `NewAccountBuilder` in Go, and the static method `create` in TypeScript.
2. The stub reads both facts on the type of each result. Its records refer to the builder through a bare reference, which the layout of eidos qualifies for the file of the double and for the external test package of its test.
3. The templates of the stub write the default arm of each method with such a result, and the default of its subject in the test.

The annotate role reports each declaration that has no builder and stamps no fact on it. The generate role runs on each declaration with the facts and writes its builder under the stamped name, so each fact contains the name of a builder that the plan writes. The builder derives the name of its own constructor, so the stub does not repeat the naming and the option `suffix` of the builder. A composition without the builder stamps no fact, and its doubles return zero values.

The runtime of each language takes the default arm. `stub.Arms[C, R]` has the field `Default func() R`, and `Arms<C, R, I>` has the property `default`. `stub.Answer` and `answer` return the default results after the configured result and before the zero values. The behaviour suite compares an unconfigured call with the default of the subject, `Subject.Default` in Go and `default` in TypeScript, and with the zero values where the subject has no default.

The layout of eidos resolves a bare reference within the source package of its unit alone, and a TypeScript module is a package of its own. So a result takes its default result where the package of the subject declares its type, and a type of another package keeps its zero value. The layout resolves a reference by its origin and its name only where the reference is in another language. dokimasia/eidos#53 requests the same resolution in the language of the plan. With it, the stub refers to the builder of a type of another package by its origin, and that type takes its default result too.

The corpus states the integration in both languages:

- `integration/validated` puts the directive `stub` on `Store`. An unconfigured `Get` of its double returns the account of `AccountDefaults`, which the store accepts.
- `integration/built` declares a method for each row of the table of default results, beside a type without a builder and a type whose builder is in another package. The test of its double checks the default results of each method.

### The enumeration

The enum generator reads the variants of an enumeration and their texts from the rules of eidos. The text of a variant of an enumeration of text is its value, such as `us-east` for `US Region = "us-east"`. The text of any other variant is its name without the name of its type, such as `Draft` for `StatusDraft`.

Under the param `text=comment`, the trailing comment of a variant is its text, without the spaces around it. A variant without a comment keeps the text that the rules derive:

```go
//dokimi:enum text=comment
type Grade int

const (
	Poor Grade = iota // poor
	Fair              // fair enough
	Good
)
```

`Poor` renders as `poor`, `Fair` as `fair enough` and `Good` as `Good`. A TypeScript enum states the text of a member in the same way, as in `Poor, // poor`.

The directive `text` on a variant pins its text, where a protocol spells the variant in another way than its name or its value:

```go
//dokimi:enum
type Phase int

const (
	// Running is the phase of a job that a worker runs.
	//
	//dokimi:text in-progress
	Running Phase = iota
	Done
)
```

`Running` renders as `in-progress`, and `Done` as `Done`. The directive takes precedence over a trailing comment under `text=comment`. The annotator `text` of dokimi stamps the text of each instance on its variant, through the rule `OnEnumVariant` of eidos, and the enum generator reads it. A TypeScript member takes the directive in the same way. The generator writes the text surface of the enumeration and the checks of the surface. The table shows the names under the default options:

| Member | Go | TypeScript |
|---|---|---|
| Render a value as text | `func (v Status) String() string` | `statusToString(value: Status): string` |
| Parse a text | `ParseStatus(s string) (Status, error)` | `parseStatus(text: string): Status` |
| The error of a text of no variant | `var ErrUnknownStatus = errors.New("status: unknown Status value")` | `class UnknownStatusError extends Error` |
| List the variants | `StatusValues() []Status`, a new slice on each call | `statusValues(): Status[]`, a new array on each call |
| Report a variant | `func (v Status) IsValid() bool` | `isStatus(value: unknown): value is Status` |
| Write a value as text | `func (v Status) MarshalText() ([]byte, error)` | None. JSON writes the value of a member |
| Read a value from text | `func (v *Status) UnmarshalText(text []byte) error` | None |

The render of a value outside the variants returns the value itself: `Status(42)` in Go, and the value of the member in TypeScript. The error of `ParseStatus` wraps `ErrUnknownStatus` with the quoted text, as in `status: unknown Status value "Pending"`. The error class of TypeScript has the text in its message and in its field `text`.

The Go render of an enumeration of integers takes the form that the tool stringer writes:

- A compile-time check indexes an array of one element with each variant minus its value, so a changed value fails the build until the generator runs again.
- The render slices the text of a value out of one string constant through an index of offsets.
- A switch selects among the runs of consecutive values. An enumeration with more than ten runs looks the text up in a map from each value.

The frontend of Go stamps the exact value of each variant, and a variant whose value reads a constant of another package has none. The render of such an enumeration, of an enumeration of floats and of an enumeration of truth values switches over the variants. The render of an enumeration of text whose texts are its values is the conversion `string(v)`.

The surface writes each member that the package does not declare. A member is declared where the package has a declaration of its name. In Go that is a function, a variable, a constant or a type of the package, or a method of the enumeration. In TypeScript it is a function, a class or a variable of the module of the enum. The parse function, its error and the encoding pair depend on other members:

- The parse function reads the texts that the render writes. So the surface writes a parse function only beside a render that it writes, because a declared render can write other texts.
- The parse function returns its error. A parse function of the surface wraps the error that the package declares, and the surface writes the error only where the package declares none.
- The decode reads with the parse function, and the encode writes with the render. The surface writes a decode where a parse function exists, and an encode only where a decode exists, so the two stay a pair.

The list of the variants and the validity method depend on no other member. The directive `enum methods=off` leaves out the whole surface and keeps the checks.

A member of the surface and a check call each declared member, so a declared member has the form of its row:

| Member | Go | TypeScript |
|---|---|---|
| The render | A method of the enumeration without a parameter that returns a text | A function of a value of the enum that returns a text |
| The parse function | A function of a text that returns a variant and an error | A function of a text that returns a variant |
| The error of the parse function | A variable | A class |
| The list of the variants | A function without a parameter that returns a list of the enumeration | A function without a parameter that returns a list of the enum |
| The validity method | A method of the enumeration without a parameter that returns a truth value | A function of one value |
| The encode | A method of the enumeration without a parameter that returns bytes and an error | — |
| The decode | A method of the enumeration of bytes that returns an error | — |

A declaration of another form under the name of a member conflicts with the surface. The generator reports `Conflict` at the declaration and writes nothing for the enumeration.

The enum generator has five options. Four of them are name templates of the members whose names the language does not fix. The table shows the default of each and the names that it renders for `Status`:

| Option | Default | Member | Go | TypeScript |
|---|---|---|---|---|
| `parse-function` | `Parse{{.Enum}}` | The parse function | `ParseStatus` | `parseStatus` |
| `values-function` | `{{.Enum}}Values` | The function that lists the variants | `StatusValues` | `statusValues` |
| `parse-error` | `Unknown{{.Enum}}` | The error of the parse function | `ErrUnknownStatus` | `UnknownStatusError` |
| `valid-method` | `IsValid` | The validity method | `IsValid` | None. TypeScript names its type guard `isStatus` |

The option `no-tests` leaves out the checks of every enumeration.

A name template has the syntax of `text/template` of Go, as the name options of mockery and GoReleaser have. Its only action is `{{.Enum}}`, the name of the enumeration. The text around the actions contains letters, digits, spaces, hyphens and underscores, and the text before the first action starts with a letter. The templates of `parse-function`, `values-function` and `parse-error` use `{{.Enum}}`, because two enumerations of a package would otherwise declare one name twice. The validity method is a method of its enumeration, so the template of `valid-method` can leave the action out.

Each language writes a rendered name in its own case convention:

- Go writes a function and a method in PascalCase. It writes the error as a variable with the prefix `Err`. staticcheck requires that prefix of an error variable in its check ST1012.
- TypeScript writes a function in camelCase. It writes the error as a class in PascalCase with the suffix `Error`. The built-in error classes of JavaScript, such as `TypeError`, have that suffix.

So `parse {{.Enum}} text` renders `ParseStatusText` in Go and `parseStatusText` in TypeScript. A codebase that declares the validity method as `Valid` and names its parse functions as enumer does, such as `StatusString`, sets these options in the config file:

```yaml
options:
  enum:
    parse-function: "{{.Enum}}String"
    valid-method: Valid
```

YAML needs quotes around a value that starts with `{`. The declared method `Valid` then is the validity method that the checks call, and the surface does not write a second one.

A template that breaks a rule stops the build of the composition. The error contains the option and the rule. A template of `parse-function`, `values-function` or `parse-error` without `{{.Enum}}` stops the build too.

The checks call each member of the surface, and each member that the author declares on the enumeration or in its package. They assert through the assert library of the language, in the test of the enumeration, such as `TestStatus`. Its first case checks that no two variants share a value, and the test has a level for each member that it calls:

| Level | Checks |
|---|---|
| `StatusValues` | The list contains the variants in the order of their declaration, and each call returns a new slice |
| `String` | A variant renders the text that its value, its trailing comment or its directive `text` states. No two variants share a text. A value outside the variants renders another text than each variant, where the surface writes the render |
| `IsValid` | Each variant is valid, and a value outside the variants is not. The zero value is the variant that declares it, or no variant |
| `ParseStatus` | Each variant parses back from its text. A text of no variant returns an error that matches the error of the parse function, where the surface or the package has one. In Go the error of a parse function of the surface contains the text |
| `MarshalText` | Each variant writes its text |
| `UnmarshalText` | A variant that the surface writes and reads makes the round trip, and the read refuses a text of no variant |

A level checks each variant in a case of its own, the row of a table.

The text of no variant is a run of question marks one longer than the longest text. The rules of each language derive the value outside the variants, such as the value after the largest number or a text that no variant has. The checks of TypeScript leave out the zero value, which TypeScript does not have for a type. They also leave out the encoding pair.

In Go, the surface is in the package of the enumeration, because Go declares a method in the package of its type. Its checks are in the external test package of the directory, such as `status_test`. In TypeScript the surface and its checks are in modules beside the module of the enum, such as `status.enum.ts`.

### The checks of errors

An author writes the directive `sentinel` above the package clause of a Go package or in the header of a TypeScript module. The generator reads the errors of the package from the rules of eidos:

| Errors | Go | TypeScript |
|---|---|---|
| The error values | The variables whose name the rules of Go accept for an error value: `Err` and an exported identifier | None. The rules of TypeScript accept no name for an error value |
| The error types | The public structs without type parameters whose method set has `Error() string`, with the methods that their embeds promote | The public classes without type parameters that extend a global error class directly or through other classes |
| The value of an error type | The zero value, or its address where the method `Error` has a pointer receiver | A construction through the first public constructor whose required parameters have samples, with a sample of each parameter up to the first optional parameter without one |
| The cause of an error type | The field that the hook `Unwrap() error` returns, declared or promoted, unless an embed by pointer promotes the field | The field `cause`, where the constructor has a parameter of that name and no required parameter after it |
| The members of an error type | The public fields that the type declares, each with a sample | The public fields that the class declares, without the fields `name`, `message`, `stack` and `cause` of `Error` |

The directive has two params:

| Param | Meaning |
|---|---|
| `prefix` | The opening of each message without its separator. `prefix=off` leaves out the checks of the opening |
| `no-overlap-with` | The paths of the packages whose error values the error values of the package do not match, such as `no-overlap-with=[example.com/shop/neighbour]` |

Without the param `prefix`, the opening of a Go message is the name of the package and `: `, such as `store: `. The option `prefix` of the generator replaces the name of the package in that opening. TypeScript has no convention for the opening of a message, so its checks of the opening exist only where the directive states `prefix`.

The generated file of a package writes the checks of the error tests of testkit, in a test of each error value and of each error type. They assert through the assert library of the language:

| Test | Level | Checks |
|---|---|---|
| `TestErrNotFound` | `Error` | The error value is set and has a message with the opening. Its message differs from the message of each other error value, and opens with the whole message of none |
| | `errors.Is` | The error value matches no other error value, and no error value of a package of `no-overlap-with` |
| `TestNotFoundError` | `Error` | The value of the type reports a message with the opening. Each member of text is in the message |
| | `Is` | A type that compares itself and has no hook `Unwrap` agrees with `errors.Is` on each error value of the package |
| | `Unwrap` | A type with a cause returns the cause that a check gives it |
| | `errors.Is` | The value of the type matches no other error type of the package |

The test of an error value loops over the other error values and reports each failure through the package `expect` of assert-go with the name of the other error value, so the loop continues after a failure. A check of one value asserts through assert-go. The test of an error class of TypeScript, `describe("NotFoundError")`, checks that the class matches no other error class of its module, and has the levels `message` and `cause`. It matches two errors through `errorIsNot` of the assert library. `errorIsNot` follows the chain of `cause`. It matches an instance of the same name and message. A TypeScript module has no error values, and its classes have no hooks. Its test has no checks of error values and no check of the agreement.

| Part | Go | TypeScript |
|---|---|---|
| The file | `<file>_sentinel_test.go`, in the external test package of the directory | `<module>.sentinel.test.ts`, beside the module |
| The value of an error type | `&NotFoundError{}`, and `&NotFoundError{Key: "Key-a"}` in the check of the members | `new NotFoundError("key-a")` |
| The cause | `got.Cause = cause` | `new WrappedError("op-a", cause)` |
| The other error types | `[]error{&WrappedError{}}` | `[new WrappedError("op-a")]` |

The file of the checks is beside the source file of the first error value, as testkit writes it. A package without an error value has its checks beside the source file of its first error type. The generator writes the file through `Emitter.FileOf` of eidos.

### Invariants

| Invariant | Enforced by |
|---|---|
| A double decides each result in the order fault, function, configured result, default result | The runtime of each language and its behaviour suite |
| An unconfigured call of a double that is not strict returns a new value from the builder of each result whose type has one | The behaviour suite of each runtime, with the default of the subject |
| A fault armed on a callable that cannot fail fails the test | The runtime of each language |
| A handle has one double in its package | The stub generator, which writes the double of a handle at the first subject of the package that returns it |
| A check of `sentinel` covers every error of its kind in its package | The directive on the package, and the generator, which reads the package whole |
| A check of `sentinel` assigns no cause through an embed by pointer, whose zero value is nil | The generator, which gives such a type no cause |
| A clone of a builder shares no list, map, set or bytes value of a member with the builder | The test of each builder |
| `plugin.go` of a generator contains no word and no syntax of a language | The parts of the languages and the templates of each target |
| A generated double, builder or test imports the runtime and the assert library of its language and no package of the binary | The template trees, and a test of each binding |
| A generated check of Go uses the exported API of the package of its subject | The external test family of each generator, and the declarations of each record that the layout qualifies |

### Failure handling

| Failure | Result | Recovery |
|---|---|---|
| A test arms a fault on a callable that cannot fail | The test fails and names the callable | Remove the fault, or give the callable an error result |
| A Go interface under `stub` has type-set elements | The stub generator reports `Constraint` at the interface and writes no double | Put the directive on an interface that values implement |
| The member walk of the subject of a double has a gap, such as an embed that the load did not resolve or an overloaded method | The stub generator reports `IncompleteMembers` at the subject and writes no double | Load the package of the embed, or give the method one signature |
| A member of the subject has no implementation in a double: a method with type parameters, an index signature, a construct signature, an accessor, a member with a computed name such as `[Symbol.iterator]`, or a member of a class that is not public | The stub generator reports `Unimplementable` at the subject and writes no double | Put the directive on an interface without the member |
| The subject of a TypeScript double is not exported | The stub generator reports `Unimplementable` at the subject and writes no double | Export the subject |
| A parameter, a result or a property of the subject has no declared type, such as a TypeScript method whose result type the compiler infers | The stub generator reports `Unimplementable` at the subject and writes no double, because a double states the type of each field of a call, of each result and of each property | Declare the type |
| Two members of a double have one name, such as the field `OnClose` of the method `Close` and a method `OnClose` | The stub generator reports `Collision` at the subject and writes no double | Rename one of the methods |
| The subject of a double has no member that a double can implement | The stub generator reports `Empty` at the subject and writes no double | Remove the directive, or add a member |
| eidos cannot sample the type of a parameter or a result of a callable | The test of a Go double leaves the argument or the result at its zero value. The test of a TypeScript double passes `undefined` as the argument, and checks no result that needs the value, because TypeScript has no zero value of every type | Author a sample with the directive `sample` |
| The directive `fault` is on a Go method without an error result | The annotator `fault` reports the instance, and the double has no helper of it | Remove the directive, or give the method an error result |
| A name of `errors` resolves to no variable of the package, or `partition` to no parameter of the method | eidos reports the instance at the directive | Correct the name |
| The param `retry` of the directive `fault` is below 2 | The annotator `fault` reports the instance, and the double has no retry schedule | State the attempt that succeeds |
| The param `partition` names a variadic parameter, or a parameter whose type mentions a type parameter of the double | The stub generator reports `Incomparable` at the declaration, and the double has no partition helpers, because a partition fault compares each argument with a key | Select the calls with `FaultsWhen`, or `faultsWhen` in TypeScript |
| A generic subject of a double has no witnesses | The double has no test | Author the witnesses with the directive `witness` |
| A sample of Returns equals the default result, such as a sample that equals the value of a companion function | The check of Returns fails and names the method | Author another sample with the directive `sample` |
| An `out` routes the builder of a result type into another Go package than the double, such as `accounttest` beside a double in `account` | The package of the builder imports the package of the double, and the double imports the package of the builder, so the double does not compile | Give the double the same `out` as the builder |
| A default is no value of the type of its field | The builder reports `BadDefault` at the field and writes nothing for the declaration | Correct the default |
| A field has the tag `builder` with another value than `-` | The builder reports `BadTag` at the field, and the field keeps its setter | Correct the tag |
| A declaration under `builder` has no field that a builder can set | The builder reports `Empty` at the declaration and writes nothing for it | Remove the directive, or make a field public |
| An alias under `builder` refers to no struct of the workspace, or the alias or its struct has type parameters | The builder reports `NotStruct` at the alias | Put the directive on the struct |
| Each public constructor of a class takes a required argument | The builder reports `Unconstructed` at the class | Add a constructor without required parameters |
| Two methods of a builder have one name, such as the setters of the TypeScript properties `name` and `Name`, or `WithLabelsEntry` of the Go map `Labels` and of the field `LabelsEntry` | The builder reports `Collision` at the declaration and writes nothing for it | Rename one of the fields |
| eidos cannot sample the type of a field | The seed of the test leaves the field at the start value. The test has no check of a replacement of the field and no check that a clone does not share it. A probe calls each setter of the field | Author a sample with the directive `sample` |
| eidos samples the type of a field and derives no alternate, such as for `struct{}` | The test has no check of a replacement of the field and no check that a clone does not share it | Author an alternate with the directive `sample` |
| A TypeScript property has no type annotation | The setter takes the type of the property, as `Counter["count"]`. The test checks only that the setter returns its builder. A default of the property reports `BadDefault` | Annotate the type of the property |
| An enumeration under `enum` has no variant | The enum generator reports `Empty` at the enumeration and writes nothing for it | Add the variants, or remove the directive |
| Two variants of an enumeration have one value, or one text such as `StatusDraft` and `Draft`, two trailing comments of one text under `text=comment` or two directives `text` of one text | The enum generator reports `Duplicate` at the enumeration and writes nothing for it | Give each variant a value, a name, a comment and a directive `text` of its own |
| An enumeration renders itself and its package has no parse function | The surface leaves out the parse function and the encoding pair, and the checks call the render of the author | Add the parse function, or accept the shorter surface |
| A declaration of the package has the name of a member of the surface and another form, such as a function `ParseStatus` of bytes | The enum generator reports `Conflict` at the declaration and writes nothing for the enumeration | Give the declaration the form of the member, rename the declaration, or set the option that names the member |
| A name template does not parse, has an action other than `{{.Enum}}`, has another character than a letter, a digit or a separator, or starts with another character than a letter | The build of the composition fails and names the option and the rule | Correct the template |
| The template of a function or of the error of the parse function does not use `{{.Enum}}` | The build of the composition fails and names the option | Add `{{.Enum}}` to the template |
| A Go variable has the name of the error of the parse function and another type than `error` | The generated Go code does not compile | Rename the variable, or set `parse-error` |
| A package under `sentinel` has no error value and no error type | The sentinel generator reports `NoErrors` at the package and writes nothing for it | Remove the directive |
| The directive `sentinel` is above a declaration instead of a package clause | eidos reports `EID-0072` at the directive | Move the directive above the package clause |
| A package lists itself under `no-overlap-with` | The sentinel generator reports `OwnNeighbour` at the package and leaves the path out | Remove the path |
| A package under `no-overlap-with` is not in the run | The sentinel generator warns under `Unloaded` at the package, and the check of that package compares with nothing | Bring the package into the scope of the plan |
| A variable has the name of an error value and another type, such as `ErrCount int` | The generated Go test does not compile | Rename the variable |
| A TypeScript error class has no public constructor whose required parameters have samples, such as a constructor that takes a function | The sentinel generator warns under `Unconstructed` at the class and leaves the class out of the checks. A class without a constructor of its own takes the constructor of the global error class that it extends | Author a sample with the directive `sample`, or declare a public constructor |

## Alternatives considered

### The reference generators of eidos

The composition adds the generators `builder`, `enum` and `sentinel` of the reference plugins of eidos.

**Why not:** they write Go alone, and their generated tests assert through the package `testing` of Go and not through assert-go. Their names and options are fixed. They exist to exercise eidos, so they change with the needs of eidos and not with the needs of dokimi.

### Enumeration checks without a surface

The enum generator writes the checks of an enumeration alone, and the author writes `String`, `Parse` and `MarshalText` by hand or with stringer.

**Why not:** a Go enumeration is a defined type with constants. Go does not relate the text of a variant to its value, and a conversion admits a value outside the variants. The surface states the relation once from the declaration, and the checks pin it. The surface leaves out each member that the author declares, and `methods=off` leaves out all members. An author who keeps a member written by hand gets the checks of it.

### A param of the directive enum that lists the text of each variant

The directive `enum` takes the text of each pinned variant in a param of its own, beside the name of the variant.

**Why not:** the text is then away from its variant, so a renamed or a removed variant leaves an entry whose name matches no variant. A directive on the variant moves with the variant and goes with it, as a trailing comment does.

### The trailing comment as the text by default

The trailing comment of a variant is its text without a param, as in `Poor Grade = iota // poor` under `//dokimi:enum`.

**Why not:** authors write trailing comments that explain a variant, such as `Calm Mood = iota // the reader reads on`. As the text, such a comment would change the render and the parse function of an enumeration that adds the directive, and every value that it encodes. The source would state that change nowhere. The flag `-linecomment` of stringer is off by default for the same reason. The param `text=comment` states the choice in the directive that the author writes anyway.

### Words joined to the name of the enumeration

Each option is a word, and the part of each language joins the word to the name of the enumeration, as `parse-word: parse` gives `ParseStatus`.

**Why not:** the key of a word option does not state where the word goes. The parse function puts its word before the name of the enumeration, and the function that lists the variants puts its word after it. A word also cannot give a name whose word is on the other side, such as the parse function `StatusString` of enumer.

### A prefix and a suffix of each name

Each member has an option for the prefix and an option for the suffix of its name, as OpenAPI Generator has `--model-name-prefix` and `--model-name-suffix`.

**Why not:** the name of a method does not contain the name of the enumeration, so the validity method needs an option of another kind. The error of the parse function also has the prefix `Err` in Go and the suffix `Error` in TypeScript. A prefix option of the error is the prefix of neither name. A template states each whole name with one kind of option.

### Every action of text/template

A name template takes every action of `text/template`, such as a pipeline, a function or a condition, as the templates of mockery do.

**Why not:** a pipeline can fail for one enumeration and pass for another. It can also give a character outside a name for one enumeration alone. The build of the composition then could not check a template before the run. The part of each language writes the case of a name, so a template does not need a function of case.

### A name for one enumeration of a package

The template of the parse function leaves out `{{.Enum}}`, so a package with one enumeration declares `Parse`, as in `color.Parse`.

**Why not:** an option applies to every package of the composition, so a package with two enumerations would declare `Parse` twice. A param of the directive `enum` that sets the name of a member of one enumeration would allow such a name. A codebase that needs the name reverses this verdict.

### The validity method found by its signature

The generator takes a declared method without a parameter that returns a truth value as the validity method, whatever its name.

**Why not:** an enumeration can have other such methods, such as `IsZero` and `IsTerminal`. A wrong match would make the checks call a method of another meaning.

### Slots of the builder

The builder hands out slots for setters on the builder, for statements inside `Build` and for checks in the test file. Other generators fill them.

**Why not:** no generator of dokimi adds a setter, a statement or a check to a builder.

### The builder weaves the default results into the double

The stub declares a named slot in each method of a double. The builder runs in a later turn and triggers on each such method through `sdk.OnEmit`. It appends the call of its constructor to the slot through `Emitter.Slot`.

**Why not:** a slot of a body contains statements. The layout of eidos qualifies a bare reference in a type and leaves a call in a statement as written. The test of a double is in the external test package, so the call would not compile there. RFC-0057 of eidos records the limit for a reference from one emitted declaration to another in a file of another package.

### The stub derives the name of the constructor

The stub reads that a type has a builder, and joins the name of the constructor by the naming of the builder.

**Why not:** the stub would repeat the naming of the builder and its option `suffix`, so a change of the naming of the builder would break each double. The facts of the builder name the constructor.

### A builder that fills an interface field with a double

The builder fills each field whose type is an interface under `stub` with a new double.

**Why not:** a double that a builder constructs has no seat of a test, so strict mode and the expected counts of calls do not report a failure. To check the calls of a double, a test constructs it with its seat and passes it to the setter, as in `WithStore(store)`. The verdict reverses for a builder that takes the seat of the test.

### Samples from the builder in the test of a double

The test of a double takes the samples of a type with a builder from the builder.

**Why not:** the check of Returns needs samples that differ from the default results, and the builder returns the default results.

### A mocking library of each language

Consumers and planted defects use a mocking library of each language, such as gomock in Go or the mocks of Vitest in TypeScript.

**Why not:** each library has its own recording, its own order of results and no shared fault strategies. A planted defect in one language would then differ from a planted defect in another. The runtime of dokimi implements one double contract in every language, and the generated double is typed by the declaration.

### Claims of the conformance spec for enumerations and errors

The enum and sentinel generators bind claims of the conformance spec. Each claim has a planted defect and calls a probe of the spec, like a check of a suite.

**Why not:** the enumeration tests and the error tests of testkit are the reference output of both generators, and they assert without probes. A claim would replace each check of testkit with a probe call and a planted defect, so the generated tests would differ from the tests that testkit writes for the same source.

### A check that a wrapped error value matches the value

The sentinel generator also checks that an error that wraps an error value with `%w` matches the value.

**Why not:** `errors.Is` compares the identity of a comparable value before it calls a method of its type. Every error value of a comparable type then matches through `%w` and `errors.Join`, so the check tests the standard library. The error tests of testkit leave it out for the same reason.

### A separate directive for the packages that must not overlap

The author lists each other package in a directive of its own, such as `//dokimi:sentinel-no-overlap-with example.com/shop/neighbour`.

**Why not:** a param of the directive `sentinel` takes a list in the grammar of directives, such as `no-overlap-with=[example.com/shop/neighbour]`. A second directive would need its own schema, and it would mean nothing on a package without `sentinel`.

### A directive on each error

The author writes the directive `sentinel` above each error value and each error type.

**Why not:** a check that no error value matches another then covers the errors with the directive alone, and an error without the directive passes unchecked. The directive on the package covers every error of the package. The rule `sdk.OnPackage` of eidos runs the generator for the directive on a Go package and on a TypeScript module.

### Checks without a directive

The sentinel generator checks every package with an error value, as the reference generator of eidos does. An author excludes a declaration with the kernel directive `skip`.

**Why not:** the generator would write a test into every package of the workspace with an error. Every other generator of dokimi writes output only where an author writes its directive.

### The errors of the directive fault as positional arguments

The directive `fault` lists its errors as positional arguments, as `//testkit:fault ErrNotFound retry=3 partition=Key` of testkit does.

**Why not:** the grammar of directives of eidos has no variadic positional param. A param of the type of a list of references takes the names in brackets, and eidos resolves each name to a variable of the package before a plugin reads it.

### A double without methods declares its properties alone

The double of an interface without methods declares its properties, and leaves out the options, `newDelayed` and `resetCalls`.

**Why not:** every double has the same constructor, the same options and the same reset, so code that builds doubles builds a handle the same way as its subject. The options of such a double do not change its behaviour. Its test checks that each of them applies.

### Planted defects for doubles and builders

The tests of a double and of a builder come with planted defects, like the checks of a suite.

**Why not:** a double and a builder are code that dokimi generates, not code of the consumer. A planted defect would prove that the test of the double can fail, which the behaviour suite of the runtime already proves once for every double. The corpus checks the output of both generators.

## Drawbacks

- **Four more template trees per language.** With eight targets, the four generators need 32 template trees.
- **A double runtime in each runtime.** Each runtime implements the recording, the order of results, six fault strategies, strict mode and expectations. The stub runtime of testkit has 1,485 lines of Go and its fault package 328 lines.
- **Checks without planted defects.** The checks of `enum` and `sentinel` are fixed templates, so their output, the golden trees of the generators and the corpus show that they pass. No planted defect shows that each of them can fail.
- **Name and parse functions per language.** Each binding maps the name and the parse of a member to a construct of its language, and a language without one has no surface of that member.

## Open questions

None.

## References

| What | Where |
|---|---|
| The stub runtime and the fault strategies of testkit | <https://github.com/thesm-os/testkit/tree/main/stub> |
| The reference plugins of eidos | <https://github.com/dokimasia/eidos/tree/main/eidos-reference> |
| The facts of an owner and the two-package weave of eidos | <https://github.com/dokimasia/eidos/blob/main/docs/architecture/06b-authoring.md> |
| The qualification of a reference by the layout of eidos | <https://github.com/dokimasia/eidos/blob/main/eidos-core/layout/reference.go> |
| The limit of a reference between emitted declarations, in RFC-0057 of eidos | <https://github.com/dokimasia/eidos/blob/main/docs/rfc/0057-template-helpers-and-scaffold-forms.md> |
| A reference by origin and name in the language of the plan | <https://github.com/dokimasia/eidos/issues/53> |
| The error tests of testkit | <https://github.com/thesm-os/testkit/tree/main/conformance/corpus/errors> |
| The tool stringer, whose forms the Go render of an enumeration of integers takes | <https://pkg.go.dev/golang.org/x/tools/cmd/stringer> |
| `errors.Is` | <https://pkg.go.dev/errors#Is> |
| The syntax of a name template | <https://pkg.go.dev/text/template> |
| The templated names `structname`, `filename` and `pkgname` of mockery | <https://vektra.github.io/mockery/latest/configuration/> |
| The option `name_template` of GoReleaser | <https://goreleaser.com/customization/templates/> |
| The names that enumer generates, such as `StatusString` and `IsAStatus` | <https://github.com/dmarkham/enumer> |
| The check ST1012 of staticcheck on the names of error variables | <https://staticcheck.dev/docs/checks/#ST1012> |
| The options `--model-name-prefix` and `--model-name-suffix` of OpenAPI Generator | <https://openapi-generator.tech/docs/usage/> |
