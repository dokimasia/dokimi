---
rfc: 0006
title: Generated inputs, fuzz targets, benchmarks and budgets
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-22
updated: 2026-10-22
discussion: none
supersedes: none
superseded-by: ADR-0002, ADR-0006, ADR-0009, ADR-0010, ADR-0011, ADR-0012, RFC-0007, RFC-0009, RFC-0010, RFC-0011 and RFC-0012 in part
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0006: Generated inputs, fuzz targets, benchmarks and budgets

## Summary

The suite generator gains two machineries. Under `generated-inputs`, the property engine of the assert library generates the inputs of a check and shrinks each failure to its smallest case. Under `a-benchmark-run`, the benchmark contract of the assert library measures a callable. Both cover every callable that a suite covers: the methods of an interface or of a concrete type, functions, and generic declarations at their witnesses.

The generator writes four kinds of output with them:

- **Generated bindings.** A claim whose check passes sample inputs to a call also gets a check that generates those inputs from the shapes of their types. The claims `commutative` and `associative` bind under `generated-inputs` alone. A check of a consumer that takes a case of the property engine is a property check too.
- **Fuzz targets.** Each callable of each application gets a fuzz target. The target runs the bodies of the property checks of the callable under the fuzzer of the language. The same bodies run under `ForAll` in every test run and under the campaign profile of the property engine.
- **Budgets.** The directive `budget` states the ceilings of a callable for allocations, bytes and latency. Each ceiling is a check with a planted defect. Allocation ceilings run in every test run over generated inputs, and the other ceilings run in benchmarks.
- **Benchmarks.** Each callable of each application gets a benchmark of 100 generated calls under the ceilings of its budget.

Each generated value starts from the starting value of its type: its authored sample, the value of its builder, or its declared default. A case keeps or changes each field of that value, and a failure shrinks back to it. Enumerations, builders, declared errors and the faults of doubles apply at every depth of a value.

## Motivation

A check of a single caller passes the two samples of each type to its call, so it misses a subject that fails at other inputs:

- a store that truncates keys of more than 255 bytes
- a parser that panics on an empty document
- a decoder that loses the last element of a long list

The property engine of the assert library finds such inputs and shrinks each one to the smallest input that fails. A person still writes a property by hand for each callable and each law.

A codebase that decodes outside data keeps a fuzz target for each decoder and a property with the same body that every test run checks. A codebase with allocation budgets keeps a benchmark of each hot path and a test that checks the same ceiling without `-bench`. Both repeat the same code for every callable. A ceiling in either codebase comes without a proof that it can fail.

testkit generates benchmarks for the interfaces of Go. Its design has five limits:

- It covers interfaces alone. A function and a method of a concrete type do not get a benchmark.
- It passes zero values to each method, so it measures one input. A method that panics on a zero value skips its benchmark.
- It does not state a ceiling. A consumer adds a primitive of the shape of each method to an option of the method. The consumer also writes out the type arguments of each primitive, as in `bench.ReaderAllocsWithin[allshapes.Service, string, allshapes.Item]("seed-1", 0)`.
- A ceiling does not come with a planted defect, so nothing proves that it can fail.
- It does not generate fuzz targets. Its generator of codecs is a plan in its documentation.

The property engine derives a generator from the shape of any type, and two types of one shape generate the same values from one seed in every language of assert-spec. dokimi derives the rest from the facts of eidos:

- which laws a check can verify on every input
- how each parameter gets a value
- which subject receives each call
- which planted defect proves each check

## Detailed design

### Concepts

| Term | Meaning |
|---|---|
| Call generator | A generator of the property engine that dokimi writes for each callable. Each of its values is a call of the callable with arguments that the case generated |
| Starting value | The simplest generated value of a type: its authored sample, the value of its builder, its declared default or the simplest value of its shape, in that order |
| Sample binding | The binding of the inputs of a claim to the samples of eidos. Every check of a single caller runs it |
| Generated binding | The binding of the same claim to generated inputs, under the machinery `generated-inputs` |
| Property check | A check whose body runs once for each case of the property engine. A generated binding and a check of a consumer that takes a case are property checks |
| Ceiling | One limit of a budget: the allocations, the bytes, the p99 latency or the mean latency of one call |

### The machineries

The conformance spec gains two machineries:

| Machinery | What the run supplies |
|---|---|
| `generated-inputs` | Inputs of the calls of a check that the property engine generates and shrinks |
| `a-benchmark-run` | A run of the benchmarks of the test runner, such as `go test -bench` or `vitest bench` |

`generated-inputs` is deeper than `nothing-beyond-a-caller` and shallower than `generated-sequences`. Under `generated-sequences` the engine also chooses the calls, and the claims of that machinery have a design of their own. The sample binding of a claim whose law has the depth `generated-inputs` witnesses the law at its shallow end, and the suite record marks that check as shallow.

### The values of parameters

The call generator of a callable gives each parameter a value by its kind:

| Parameter | Go | TypeScript |
|---|---|---|
| In the role `context` | The context of the case | The signal of the case |
| Of a type whose shape the assert library reads | A value of `prop.Of[T]()` | A value of the shape definition that dokimi writes from the facts of eidos |
| A pointer, or a type with `undefined` | A pointer to a value of the type, and also nil under the mixin `nil-safe` | A value of the type, and also `undefined` under the mixin `nil-safe` |
| Variadic | A list of values of the element type. A variadic parameter of functions or of interfaces outside the standard library takes no values | A list of values of the element type |
| Of an enumeration | One of its members | One of its members |
| Of a struct or a class with a builder | The value that the builder builds, with each changed field from a value of the shape of the type | The instance that the builder builds from the starting value of the class, with the changed fields through its setters |
| Of a function type | A function whose nth call returns the nth value of the generated list of each result. After the last values it returns its zero results with the end error of each error result. An error of a list is nil or a declared error | A function that returns or resolves with values that it generates at each call, or throws a declared error |
| Of an interface | An implementation that the suite declares for an interface outside the standard library. Its methods return generated results as a function does. An interface of the standard library takes the generator of its registration | The double of an interface with a double of dokimi, whose methods generate their results as a function does. It takes its random numbers from the case |
| A channel that the callable receives from | A closed channel of a generated list | — |
| Of any other type | The generator that a registration of the test process states for the type | The generator that a registration of the test file states for the type |
| The receiver of a method | A fresh subject of the factory of the application for each case | A fresh subject of the factory of the application for each case |
| A type parameter | The type argument of its witness, as in the suite | The type argument of its witness, as in the suite |

The label of each value is the key form of the name of its parameter, such as `user-id`. A counterexample then shows the same labels in every language.

The shape of a type admits every value of the type. A declaration narrows or widens the domain of a callable:

- The tags of `prop` on the fields of a struct narrow a field, such as `prop:"pattern=[a-z0-9-]{1,64}"` on the field of an ID.
- A registration replaces the generator of a type. In Go `prop.Register` registers it in an `init` function of the package of the suite or in `TestMain`. In TypeScript `register` registers it in the setup file of vitest.
- The mixin `nil-safe` adds nil to the domain of a pointer.

The header of the suite lists each parameter of any other type with the registration that it needs. Without the registration the property engine fails the check before its first case, and its failure contains the name of the type.

### The starting value of a type

The simplest generated value of a type is its starting value. The first source that states one decides it:

1. The sample that the kernel directive `sample` authors for the type.
2. The value of the builder of the type. `NewConfigBuilder()` starts from the companion function `ConfigDefaults()` and the declared defaults of the fields.
3. The default that the annotator `defaults` stamps on the type, and for a struct the default of each field. The rules of the language read the text of a default. A type of text takes a text that the rules do not read as its content.
4. The simplest value of the shape of the type, such as the zero value.

This is the order of the first sample in RFC-0003, with the builder added. A callable with the mixin `sample` takes the value of the builder of the mixin as the starting value of its inputs.

A generated value of a struct keeps or changes each field of its starting value, by a choice of the case. So the cases cluster around a value that the declaration accepts, and a failure shrinks back towards that value. A counterexample then shows only the fields that differ from it. `config.Config` of the corpus declares the defaults `localhost`, `8080` and `0.75`, so a counterexample such as `Port: 0` shows the one field that breaks the call.

The generator of a type draws a value of the shape of the type and assigns each changed field of that value to the starting value. Inside the generator of a type, `prop.Of` of the type takes the shape without the generator, so the tags of `prop` on its fields apply. A field with the tag `prop:"-"` keeps its starting value. A defined type that is no struct, such as `type Tag string`, keeps or changes its starting value whole. A generic type takes the simplest value of its shape. The store of the property engine replays a case from its choices alone, so a case of such a generator replays as every other case does.

### Types inside other types

The rules of the values of parameters apply to a type at every depth. An enumeration inside a struct takes one of its members, and a type with a builder inside a list takes the value of its builder. The shape reader of `prop` reads a defined type over an integer over every integer, so the generator overrides each enumeration. The generator leaves the types of the standard library to the shape reader, because their sets of constants include flags such as `fs.FileMode` and units such as `time.Duration`.

Each property run passes an override of each such type through the option `Using` of `prop`. An override applies to each value of the type inside the inputs of the run, and to no other run. A registration applies to the whole test process. A second registration of one type panics, and two suites in one test binary can need the same type. `Using` is an option of the property forms alone, so this design requests it for `ForAll` and `Fuzz` from assert-spec, in every language.

### Declared errors

An error result of a generated function or of a double is nil, or an error that the case chooses from the first source that states one:

1. The errors that the directive `fault errors=` lists for the method.
2. The error value that a mixin of the method names, such as the sentinel of the mixin `not-found`.
3. `conform.ErrGenerated`.

A caller of a dependency then meets the errors that its code handles, such as `ErrNotFound`. A TypeScript function or double throws the error, or rejects with it.

The nth call of a generated function of Go returns the nth value of each of its lists. After the last values it returns its zero results. Each error result of those zero results is the end error: the first declared error of the method, or `conform.ErrGenerated`. A caller that reads until an error, such as the loop over a feed, then ends.

Each double that a property check builds takes the random source of the case through its option `WithRandSource`. Its probabilistic faults then replay and shrink with the case. A factory reaches the same source through the method `Rand` of the handle.

### Generated bindings

A claim gets a generated binding on a callable when its sample binding passes sample inputs to a call. The sources of such inputs are `FromCallable`, `FromPartner` and `FromRole`. The generated binding generates every input of every call. It keeps each other input of the probe: the params of the classification, the zero values, the samples of results and the clock.

The calls of one case share the value of each parameter with the same key and type. So a write and a read that both take a parameter `key string` write and read the same key. A claim whose calls do not share a parameter keeps its sample binding alone, and the header states the reason. The mixin `read-after-write` with `write=Put` on `Get(ctx, id string)` is such a claim when `Put` takes a `Session`. No parameter relates the session to its ID. A property check of the consumer states such a relation, as [Property checks of a consumer](#property-checks-of-a-consumer) shows.

A generated binding has the law, the probe, the reason code and the planted defect of its claim, with the machinery `generated-inputs`. Its identity adds the case `generated` to the identity of the sample binding: `Store.Get/zero-on-error#generated`. For a claim that binds twice on one callable, the case of each binding gains the suffix `-generated`, as in `Store.Get/expires#0-generated`. In Go `StoreChecks.Get.ZeroOnErrorGenerated()` returns the identity, and in TypeScript `storeChecks.get.zeroOnErrorGenerated()` does.

The check runs the property engine with its statement as the contract. Each case runs on a handle of its own:

- The handle wraps the case. Each assertion of the probe then reports to the case, and the engine shrinks the failure.
- `Subject` returns a fresh subject of the factory for each case. A factory that takes the handle gets the case through `Unwrap`, so it can generate a receiver or prepare a state.
- `Clock` returns a fresh controlled clock for each case, and `Context` returns the context of the case.
- `Rand` returns the random source of the case, whose every value is a choice of the case.

A probe takes the same inputs in both bindings, so one implementation of each probe serves both. The property engine reports the record of the run to the handle of the check. The failure of the smallest case contains the assertion of the probe, and the assertion identifies the reason code. A run in which no case exercises the claim is vacuous.

The Go suite of `Store` declares the call generator of `Get` and its generated binding of `zero-on-error`:

```go
// storeGetCalls generates calls of Get with an ID that the case generates.
var storeGetCalls = prop.Composite(func(c *prop.Case) func(context.Context, session.Store) (session.Session, error) {
	id := c.Draw(prop.Of[string](), "id")
	return func(ctx context.Context, s session.Store) (session.Session, error) {
		return s.Get(ctx, id)
	}
})

// One check of StoreSuite.
{
	ID:        StoreChecks.Get.ZeroOnErrorGenerated(),
	Statement: "Get returns the zero value beside an error",
	Property: func(h *StoreHandle, c *prop.Case) {
		call := c.Draw(storeGetCalls, "call")
		s := h.Subject()
		probe.ZeroBesideError(h, func(ctx context.Context) (session.Session, error) {
			return call(ctx, s)
		}, session.Session{})
	},
	Defect: valueBesideError,
},
```

### Claims of every pair and triple of inputs

The mixins `commutative` and `associative` state a law about every pair or every triple of inputs, so the samples of a type do not witness it. Their claims bind under `generated-inputs` alone, so their identities do not have a case:

| Claim | Licence | Gates | Probe | Reason code | Law | Defect form |
|---|---|---|---|---|---|---|
| `commutative` | Mixin commutative | two-inputs-of-one-type | `either-order-agrees` | `order-dependent` | `combining-ignores-the-order` | `order-dependent` |
| `associative` | Mixin associative | two-inputs-of-one-type, result-of-the-input-type | `either-grouping-agrees` | `grouping-dependent` | `combining-ignores-the-grouping` | `grouping-dependent` |

`either-order-agrees` asserts `commutative` of assert-spec on two generated inputs, and `either-grouping-agrees` asserts `associative` on three. The gate `two-inputs-of-one-type` passes when the callable takes exactly two inputs of one type. The gate `result-of-the-input-type` passes when the callable returns one value of that type. The section `unclaimed` of the spec loses its entries of the two mixins.

### Property checks of a consumer

The directive `check` binds a function of the handle and a case as a property check. The suite runs it under the property engine, and the fuzz target of its callable runs it under the fuzzer. Its decisive assertion passes `h.Statement()` as its message, as the assertion of every check of a consumer does. Falsification passes when the run fails on the planted defect and the failure of its smallest case has the statement as its contract.

The team of the session store relates a session to its ID in a property check:

```go
// KeepsTheLastPut checks that Get returns the session of the last Put of its
// ID.
//
//dokimi:check Store.Get defect=keepsTheFirstPut
func KeepsTheLastPut(h *StoreHandle, c *prop.Case) {
	s := h.Subject()
	first := c.Draw(prop.Of[session.Session](), "first")
	second := c.Draw(prop.Of[session.Session](), "second")
	second.ID = first.ID
	assert.NoError(h, s.Put(h.Context(), first), "Put stores the first session")
	assert.NoError(h, s.Put(h.Context(), second), "Put stores the second session")
	got, err := s.Get(h.Context(), first.ID)
	assert.NoError(h, err, "Get reads the session")
	assert.Equal(h, got, second, h.Statement(), assert.EquateEmpty())
}

// keepsTheFirstPut returns a Store whose Put keeps the first session of each
// ID.
func keepsTheFirstPut(h *StoreHandle) session.Store {
	kept := map[string]session.Session{}
	return NewStoreStub(h,
		WithStorePut(func(_ context.Context, s session.Session) error {
			if _, ok := kept[s.ID]; !ok {
				kept[s.ID] = s
			}
			return nil
		}),
		WithStoreGet(func(_ context.Context, id string) (session.Session, error) {
			s, ok := kept[id]
			if !ok {
				return session.Session{}, session.ErrNotFound
			}
			return s, nil
		}))
}
```

### Runs

A property check runs under the settings of the property engine. The runtime adds three settings:

- The seed comes from `DOKIMI_SEED` when the variable is set. Otherwise the seed depends on the profile of the engine. The profile `default` takes a random seed, and the profile `ci` and a test binary of a mutation run take a seed from the contract. The suite record contains the seed of each run.
- A run tries 100 valid cases. The option `conform.Cases(n)` and the param `cases` of the directive `subject` set another number for one application.
- The engine stores the smallest case of each failure under `testdata/prop/<test name>`, and the next run tries it first. A team commits the case of a real defect as a regression case.

An application with an invalid setting reports it as a problem of its options and does not run a check. The invalid settings are a `DOKIMI_SEED` that is not a decimal number below 2^64, and the option `Cases` with fewer than one case. `Check.Run` ends with a fault that wraps `conform.ErrConfig` for such a setting.

The profile `campaign` turns each run into a search for the seconds of `DOKIMI_ASSERT_PROP_BUDGET`. The labels, the fingerprints and the scores of its cases guide the search:

```sh
DOKIMI_ASSERT_PROP_PROFILE=campaign DOKIMI_ASSERT_PROP_BUDGET=600 \
	go test -run 'TestStoreConformsToStore/Get/zero-on-error#generated' ./session/postgres
```

Falsification runs a property hermetically with a seed from the key of the check. It turns the store and the shrinking off. So the cases of a planted defect never enter `testdata/prop`, and the first failing case ends the run. A campaign does not extend a falsification either.

### Fuzz targets

Each callable of each application has a fuzz target. For each input of the fuzzer, its body generates the identity of one property check of the callable under the label `check`, and runs the body of that check on a handle of the case. It chooses among the checks that the application does not excuse. A ceiling of a budget is not among them, because the instrumentation of a fuzzer changes the allocations and the latency of a call. One target for each callable gives each target one entry point, so `go test -fuzz` selects each one alone.

In Go the conformance test of a subject declares `Fuzz<Subject><Callable>`, and the test of the suite of a function declares `Fuzz<Function>`. A subject of two declarations with one callable name gets `Fuzz<Subject><Declaration><Callable>`. Each target calls the generated function `FuzzStoreSuite`, which calls `prop.Fuzz`:

- `go test` replays the stored cases of the target and the seed corpus in `testdata/fuzz/<target>`.
- `go test -fuzz '^FuzzStoreGet$'` runs the fuzzer of Go. The bridge of the property engine decodes each input of the fuzzer into a valid case. A failing input then shrinks, and the engine stores it under `testdata/prop/FuzzStoreGet`.

In TypeScript the conformance test of a subject awaits `fuzzStoreSuite` for each callable. Under vitest the call replays the stored cases and the seed files. The module `<file>.suite.fuzz.ts` exports the target of each callable for Jazzer.js. The consumer installs Jazzer.js and runs a target:

```sh
npx jazzer session/postgres/store.suite.fuzz --fuzzEntryPoint=fuzzStoreGet -- -max_total_time=600
```

### Budgets

The directive `budget` states the ceilings of a callable:

```go
// Get returns the session with the ID id until its lifetime ends.
//
//dokimi:budget allocs=1 bytes=96 p99=2us mean=1us
Get(ctx context.Context, id string) (Session, error)
```

| Param | Value | Check | Runs in | Assertion | Reason code |
|---|---|---|---|---|---|
| `allocs` | The most allocations of one call | `allocation-budget` | Every test run, over generated inputs, and the benchmark | `max-allocs`, and `bench-max-allocs` in the benchmark | `allocations-over-budget` |
| `bytes` | The most bytes that one call allocates | `byte-budget` | The benchmark | `bench-max-bytes` | `bytes-over-budget` |
| `p99` | The highest p99 latency of a call | `latency-budget` | The benchmark | `bench-max-latency` | `latency-over-budget` |
| `mean` | The highest mean latency of a call | `mean-latency-budget` | The benchmark | `bench-max-mean` | `mean-over-budget` |
| `warmup` | The number of iterations that each benchmark of the callable runs before it measures | — | The benchmark | — | — |

A warm-up is no ceiling and has no check. The benchmark contract runs it through its option `Warmup`, so the first calls build their caches outside the measurement, such as the cache of an interface conversion. Without `warmup` each benchmark runs each of its calls once before it measures. A budget may state a warm-up without a ceiling.

A budget on a callable of a declaration applies to every subject of the declaration. For a function and for a method of a concrete type the declaration is the implementation. A budget on a subject applies to that subject alone. Its first argument is the identity of the callable, as in the directive `excuse`:

```go
// Store is the Store on PostgreSQL.
//
//dokimi:subject session.Store new=NewTestStore
//dokimi:budget Store.Get p99=3ms mean=1ms
type Store struct{ … }
```

eidos refuses a directive that states one ceiling twice, and an `allocs` or a `bytes` that is not a whole number, before the generator reads the directive. The generator reports each of these mistakes at the directive under `BadBudget`:

- a budget without a ceiling and without a warm-up
- an `allocs`, a `bytes` or a `warmup` below 0, and a `p99` or a `mean` that is not a duration above 0
- a `mean` without a `p99`, because a mean alone passes a tail that grows
- a second budget of one callable
- a callable argument in a budget on a method or a function, because only a budget on a subject takes one
- a budget on a method that is not exported, or on a callable of a declaration without the directive `suite`
- a budget on a struct without the directive `subject`
- a budget on a subject without an argument, or with an argument that is not a callable of its declaration
- a budget on a subject for a callable whose declaration states a budget. The subject excuses the check of the declaration instead

The names of the callables of a suite have a method for each ceiling of a budget, such as `StoreChecks.Get.LatencyBudget()`, which returns the identity of its check. An excusal of a ceiling takes that identity. A check of a consumer with the segment of a ceiling, such as the function `LatencyBudget`, would have the identity of that ceiling, so the generator refuses it under `BadCheck`.

TypeScript refuses `allocs` and `bytes` with the reason of its overlay in assert-spec: V8 does not count the allocations of an iteration.

The check `allocation-budget` is a property check. Each case takes a call from the call generator and builds a fresh subject outside the count. It asserts `max-allocs` on the call. The assertion warms the call once and counts the 100 calls after it. The case also records the count as its score with `Target`, so the profile `campaign` searches for the input with the most allocations. A failure shrinks to the smallest input that allocates more than the ceiling. The test `Test<Subject>Allocs` applies these checks without `t.Parallel`, because the count covers the whole process.

The assert library does not check an allocation ceiling in a build that counts allocations differently. It does not check one either in a run whose result go test can cache. Such builds run under the race detector, msan or asan, without optimisation, or as a test binary of a mutation run. There the check reports the result `unchecked` with the reason of the library, and falsification reports the same. `go test -count=1` without the race detector checks every ceiling.

Each ceiling has a planted defect whose callable returns its zero results and exceeds the ceiling:

| Defect form | The callable of the defect |
|---|---|
| `allocates-over-budget` | Allocates one time more per call than the ceiling admits |
| `allocates-bytes-over-budget` | Allocates more bytes per call than the ceiling admits |
| `waits-over-budget` | Returns after twice the ceiling. In Go it sleeps, because the benchmark contract of assert-go times a call on the clock of the system. In TypeScript it advances the controlled clock of its handle. The contract of assert-typescript reads that clock |

A defect only adds allocations, bytes and time to a call, and a count or a clock does not go back. So each defect fails its check in every run. RFC-0001 reverses its verdict on benchmarks when a budget check cannot fail on its planted defect without flaking. These defects fail without flaking.

In Go the call of each callable has a function `Over`. It returns a subject whose callable spends a cost and returns its zero results. The runtime passes the cost of the defect form of the ceiling. Other checks of the suite can fail on a defect that returns zero results, so the falsification of a ceiling does not apply another check to its defect.

### Benchmarks

Each callable of each application has a benchmark:

1. It builds one subject from the factory, outside the measurement.
2. It takes 100 calls from the call generator of the callable, under the seed of the benchmark. The seed comes from `DOKIMI_SEED` when the variable is set, and from the key of the callable otherwise, so every run measures the same inputs.
3. It repeats the calls in their order under a benchmark contract with the ceilings of the budget of the callable. The contract first runs the warm-up of the budget, or each call once, through its option `Warmup`. A parallel benchmark runs the warm-up on its goroutines.
4. It publishes the share of the calls that returned an error as `errors/op`. The share includes the calls of the warm-up. The contract publishes the allocations, the bytes and the latencies.

A benchmark does not assert the results of its calls, because the checks of the suite assert them. The share of errors shows which path the benchmark measured. A benchmark of `Get` on an empty store measures the path of a miss, with `errors/op` at 1. A factory that stores sessions first makes it measure the path of a hit.

In Go the conformance test of a subject declares `Benchmark<Subject>`, with a level for each callable, such as `BenchmarkStore/Get`. A subject of two declarations has a level for each declaration above the levels of its callables. `Benchmark<Subject>Parallel` measures the calls of each callable with the mixin `concurrent` or `concurrent-readers` under `RunParallel`. Each level of both benchmarks passes the options of the application, which the benchmark declares once as `opts`. The test of the suite of a function declares `Benchmark<Function>`, and `Benchmark<Function>Parallel` for a function with the mixin `concurrent` or `concurrent-readers`.

In TypeScript the file `<file>.suite.bench.ts` declares a test for each callable. `vitest bench` runs the file, and `vitest run` leaves it out. The test measures the calls with the `bench` fixture of vitest. A budget adds a `Contract` of `@dokimi/assert` with `maxLatency` and `maxMean` over the same calls.

### The generated files

| Part | Go | TypeScript |
|---|---|---|
| The call generators, the property checks and the budgets of a declaration | `<file>_suite.go`, in the value of `StoreSuite()` | `<file>.suite.ts`, in the value of `storeSuite()` |
| The falsification of the property checks and of the ceilings of benchmarks | The level `StoreSuite` of `TestStoreSuite` | The `describe` of `storeSuite` |
| The falsification of the allocation ceilings | `TestStoreSuiteAllocs`, without `t.Parallel` | — |
| The application to a subject, with its property checks | `Test<Subject>ConformsTo<Declaration>` | The conformance test of the subject |
| The allocation ceilings of a subject | `Test<Subject>Allocs`, without `t.Parallel` | — |
| The benchmarks of a subject | `Benchmark<Subject>` and `Benchmark<Subject>Parallel` | `<file>.suite.bench.ts` |
| The fuzz targets of a subject | `Fuzz<Subject><Callable>` | The conformance test of the subject, and `<file>.suite.fuzz.ts` for Jazzer.js |

A function is the only subject of its suite, so the test of its suite applies the suite to the function. The level `<Function>` of `Test<Function>Suite` contains that application. The test also declares `Test<Function>Allocs`, `Benchmark<Function>` and `Fuzz<Function>`. A concrete type becomes the subject of its own suite through the directive `subject` on the type, as in `//dokimi:subject Store new=NewStore`. Its methods then get the conformance test, the benchmarks and the fuzz targets of a subject.

The Go suite of `Store` declares four entry points beside `RunStoreSuite` and `RunStoreSuiteFromHandle`:

```go
// RunStoreBudgets applies the allocation ceilings of Store to a fresh subject
// of newSubject for each case.
func RunStoreBudgets[T session.Store](t *testing.T, newSubject func(h *StoreHandle) T, opts ...StoreOption)

// BenchStoreSuite measures callable on a subject of newSubject, under the
// ceilings of its budget.
func BenchStoreSuite[T session.Store](
	b conform.B,
	callable conform.Callable[session.Store],
	newSubject func(h *StoreHandle) T,
	opts ...StoreOption,
)

// BenchStoreSuiteParallel measures callable on a subject of newSubject on
// GOMAXPROCS goroutines at once.
func BenchStoreSuiteParallel[T session.Store](
	b *testing.B,
	callable conform.Callable[session.Store],
	newSubject func(h *StoreHandle) T,
	opts ...StoreOption,
)

// FuzzStoreSuite registers the fuzz target of the property checks of
// callable on f.
func FuzzStoreSuite[T session.Store](
	f *testing.F,
	callable conform.Callable[session.Store],
	newSubject func(h *StoreHandle) T,
	opts ...StoreOption,
)
```

Each takes the factory of the handle alone, because a subject that reads the clock needs the handle in a benchmark and in a fuzz target as well. The suffix `FromHandle` separates the two forms of `RunStoreSuite`, so a function of one form does not take it. A benchmark and a fuzz target take one callable. The generated test writes a level of a benchmark and a fuzz target for each callable. The suite of a function takes the function itself, as in `FuzzSlugSuite(f, Slug)`.

### The runtime

| Part | Go: `go.dokimi.dev/conform` | TypeScript: `@dokimi/runtime` |
|---|---|---|
| A property check | The field `Property func(h *Handle[S], c *prop.Case)` of `Check[S]`, in place of `Apply` | The field `property` of `Check<S>` |
| The call generator of a callable | `Call[S]`: the callable, its call generator, its budget, whether it admits concurrent calls, and `Over`, which builds the planted defect of each kind of ceiling | `Call<S>` |
| A budget | `Budget`, with the optional fields `Allocs`, `Bytes`, `Latency`, `Mean` and `Warmup` | `Budget`, with the optional fields `latency` and `mean` |
| The calls of a suite | The field `Calls` of `Suite[S]` | The field `calls` of `Suite<S>` |
| The handle of a case | `Handle[S]`, whose `Unwrap` returns the case | `Handle<S>`, on the case |
| The random source of a check | The method `Rand` of `Handle[S]`: the source of the case in a property check, and a source of the seed of the run in any other check | The method `rand` of `Handle<S>` |
| The overrides of a suite | The field `Using` of `Suite[S]`, which the runtime passes to each property run | The field `using` of `Suite<S>` |
| The number of cases of an application | The option `Cases(n)` | The field `cases` of `Options<S>` |
| The budget of a subject | The option `Ceilings(callable, budget)` | The field `ceilings` of `Options<S>` |
| The application of the allocation ceilings | `RunBudgets(t, suite, newSubject, opts...)` | — |
| The benchmarks | `Bench(b, suite, callable, newSubject, opts...)` on a `conform.B`, which `*testing.B` implements, and `BenchParallel(b, suite, callable, newSubject, opts...)` | `bench(fixture, suite, newSubject, options)` |
| A fuzz target | `Fuzz(f, suite, callable, newSubject, opts...)` | `fuzz(seat, suite, callable, newSubject, options)`, which returns a promise of a `Target` |
| The falsification of the allocation ceilings | `FalsifyAllocs(t, suite)` | — |
| The error of a generated function | `ErrGenerated` | `ErrGenerated` |

The runtime derives the check of each ceiling from the calls of the suite. The check has its identity, such as `Store.Get/allocation-budget`, and the planted defect of its kind of ceiling. The check of an allocation ceiling also has a statement, such as "Get allocates at most 1 time per call", which is the contract of its property run. A ceiling that the option `Ceilings` adds for one subject is a check of that application alone. The application falsifies it before it applies it, as it does with a check that `conform.With` adds.

### Falsification

| Check | Falsification runs | Proven when |
|---|---|---|
| A sample binding | The check on its planted defect, as RFC-0003 states | The probe fails with the reason code of the claim |
| A generated binding | The property engine on the planted defect, hermetically | The run fails at the assertion of the reason code of the claim |
| A property check of a consumer | The property engine on the planted defect, hermetically | The run fails, and the failure of its case has the statement as its contract |
| `allocation-budget` | The check on its planted defect in `TestStoreSuiteAllocs` | `max-allocs` fails, or the run reports `unchecked` |
| `byte-budget`, `latency-budget` and `mean-latency-budget` | 10 calls of the planted defect under the benchmark contract, on a benchmark of the runtime that records each failure | The ceiling of the check fails, or the run reports `unchecked` for bytes |

In Go a planted defect of a latency ceiling of 50 µs sleeps 1 ms in all, for 10 calls of twice the ceiling.

### The suite record, the header and the lock

The suite record gains three fields of a check and one result:

| Field | Values |
|---|---|
| `checks[].cases` | The number of valid cases of a property check |
| `checks[].seed` | The seed of the run of a property check, as a decimal string |
| `checks[].counterexample` | The values of the smallest failing case by label, as the typed literals of assert-spec |
| `checks[].result` | Also `unchecked` for an allocation ceiling that the build or the run does not check, with the reason in `checks[].why` |

The version of the spec changes at its minor number, because every change adds a field, a result or an entry.

The header lists the generated bindings and the claims that keep their sample binding alone, with the reason. It also lists each parameter that needs a registration, and the ceilings of each budget. The lock contains the keys of the generated bindings and of the ceilings.

An excusal of a sample binding also excuses its generated binding, because the reason of an excusal concerns the claim. An excusal of a generated binding excuses it alone, such as for a subject whose calls take a second each.

### A worked example

The session store of RFC-0003 gets nine generated bindings:

| Callable | Generated bindings | Generated values |
|---|---|---|
| `Put` | `smoke#generated` and `cancel#generated` | `s` |
| `Get` | `smoke#generated`, `cancel#generated`, `zero-on-error#generated` and `miss#generated` | `id` |
| `Delete` | `smoke#generated`, `cancel#generated` and `delete-removes#generated` | `id` |

`Get/read-after-write` and `Get/expires` keep their sample bindings alone. `Put` takes a session and `Get` takes an ID, so no parameter relates the two calls. The team adds the property check `KeepsTheLastPut` of [Property checks of a consumer](#property-checks-of-a-consumer).

The team applies the suite to the store on PostgreSQL with 20 cases for each property check and a latency budget of `Get`:

```go
// Store is the Store on PostgreSQL.
//
//dokimi:subject session.Store new=NewTestStore cases=20
//dokimi:excuse Store.Get/expires reason="PostgreSQL computes the expiry with now(), which a test cannot advance"
//dokimi:budget Store.Get p99=3ms mean=1ms
type Store struct{ … }
```

The suite generator writes the test of the subject:

```go
// TestStoreConformsToStore applies the checks of the suite of Store to Store.
func TestStoreConformsToStore(t *testing.T) {
	t.Parallel()

	sessiontest.RunStoreSuite(
		t,
		postgres.NewTestStore,
		conform.Cases[session.Store](20),
		conform.Excuse(
			sessiontest.StoreChecks.Get.Expires(),
			"PostgreSQL computes the expiry with now(), which a test cannot advance",
		),
	)
}

// BenchmarkStore measures each callable of the suite of Store on Store.
func BenchmarkStore(b *testing.B) {
	opts := []sessiontest.StoreOption{
		conform.Ceilings(
			sessiontest.StoreChecks.Get.Callable,
			conform.Budget{
				Latency: 3 * time.Millisecond,
				Mean:    time.Millisecond,
			},
		),
		conform.Excuse(
			sessiontest.StoreChecks.Get.Expires(),
			"PostgreSQL computes the expiry with now(), which a test cannot advance",
		),
	}
	b.Run("Get", func(b *testing.B) {
		sessiontest.BenchStoreSuite(
			b,
			sessiontest.StoreChecks.Get.Callable,
			func(*sessiontest.StoreHandle) *postgres.Store { return postgres.NewTestStore() },
			opts...)
	})
}

// FuzzStoreGet fuzzes the property checks of Store.Get on Store.
func FuzzStoreGet(f *testing.F) {
	sessiontest.FuzzStoreSuite(
		f,
		sessiontest.StoreChecks.Get.Callable,
		func(*sessiontest.StoreHandle) *postgres.Store { return postgres.NewTestStore() },
		conform.Excuse(
			sessiontest.StoreChecks.Get.Expires(),
			"PostgreSQL computes the expiry with now(), which a test cannot advance",
		),
	)
}
```

The levels `Put` and `Delete` of `BenchmarkStore`, `FuzzStorePut` and `FuzzStoreDelete` follow the same form. Each entry point takes the options of the application and checks its excusals against every check of the suite. The budget does not state an allocation ceiling, so the test does not declare `TestStoreAllocs`.

Among the failures of the first run, `Get/miss#generated` shrinks to the ID `"\x00"`. PostgreSQL refuses the byte 0x00 in a value of the type `text`, and `Get` returns the error of the driver in place of `ErrNotFound`. `KeepsTheLastPut` fails at the same byte, because `Put` returns that error for a session with the ID `"\x00"`. The team makes `Get` report `ErrNotFound` for an ID that the store cannot keep. It also states the domain of a session ID on the declaration, so every property check of a session generates IDs that the store keeps:

```go
// Session is the session of one signed-in user.
type Session struct {
	ID    string `prop:"pattern=[a-z0-9-]{1,64}"`
	User  string
	Roles []string
}
```

The engine keeps the failing case of `Get/miss#generated` under `session/postgres/testdata/prop`, and the team commits it as a regression case. The failing case of `KeepsTheLastPut` is outside the stated domain, so the team deletes it.

A function gets the same parts. The function `Slug` with the directives `suite` and `budget allocs=4` gets these declarations beside `TestSlugSuite`:

```go
// TestSlugAllocs applies the allocation ceiling of Slug to Slug.
func TestSlugAllocs(t *testing.T) {
	RunSlugBudgets(t, Slug)
}

// BenchmarkSlug measures Slug.
func BenchmarkSlug(b *testing.B) {
	BenchSlugSuite(b, Slug)
}

// FuzzSlug fuzzes the property checks of Slug.
func FuzzSlug(f *testing.F) {
	FuzzSlugSuite(f, Slug)
}
```

### Requests

| Project | Request | Why |
|---|---|---|
| assert-spec | The values that a generator yields for a seed, outside a property ([dokimasia/assert-spec#16](https://github.com/dokimasia/assert-spec/issues/16)) | A benchmark measures a fixed set of generated calls in every run |
| assert-spec | A query of whether a run checks allocation ceilings, with the reason when it does not ([dokimasia/assert-spec#17](https://github.com/dokimasia/assert-spec/issues/17)) | Without it, an allocation ceiling passes in a build that does not count, and its falsification fails |
| assert-spec | The option `Using` for `ForAll` and `Fuzz` ([dokimasia/assert-spec#18](https://github.com/dokimasia/assert-spec/issues/18)) | A property run overrides the generator of an enumeration or of a type with a builder at every depth, without a registration of the test process |
| assert-spec | The generators of `using` for `values` ([dokimasia/assert-spec#19](https://github.com/dokimasia/assert-spec/issues/19)) | A benchmark takes its calls from the same domain as a property run of the suite, with the overrides of its enumerations and builders |

The benchmarks, the allocation ceilings and the property checks each need one of these requests. The fuzz targets run the property checks and need `Using` as well. #19 requests the overrides of the suite for the calls that a benchmark takes from `values`.

### Invariants

| Invariant | Enforced by |
|---|---|
| Every property check has a planted defect, and falsification fails when the run of the property passes on it | The falsification test of each generated suite |
| A property check takes its inputs from the shapes of their types, the registrations of the test and the case, and never from the subject | The call generators that the suite generator writes |
| The calls of one case share the value of each parameter with the same key and type | The call generators that the suite generator writes |
| A generated value starts from the starting value of its type, and a failure shrinks back to it | The generators that the suite generator writes for each type with a starting value |
| An enumeration generates its members alone, at every depth of a value | The overrides of each property run and of the calls of each benchmark |
| Each double of a property check takes its random numbers from the case | The call generators that the suite generator writes |
| Each case of a property check runs on a fresh subject and a fresh clock | The handle of a case |
| The falsification of a property check does not write into `testdata/prop` | The runtime, which falsifies hermetically without the store |
| Every ceiling of a budget is a check with a planted defect that fails it in every run | The defect forms of budgets, and their falsification |
| An allocation ceiling that a build does not count reports `unchecked` and never `pass` | `RunBudgets` and `FalsifyAllocs` |
| A benchmark measures the same calls in every run of one seed | The seed of the benchmark |
| Every callable of every application has a fuzz target and a benchmark | The suite generator |
| A counterexample names each input by the key of its parameter in every language | The labels of the call generators |

### Failure handling

| Failure | Result | Recovery |
|---|---|---|
| A property check fails on a subject at an input that the callable does not accept | The check fails with the smallest such input | Make the callable handle the input, or state its domain with tags, a registration or the mixin `nil-safe` |
| A parameter has a type that the shape reader refuses, and no registration covers it | The header lists the type. The property engine fails the check before its first case | Register a generator of the type |
| The calls of a claim do not share a parameter | The claim keeps its sample binding alone, and the header states why | A property check of the consumer that relates the calls |
| A generated binding is vacuous on its planted defect | Falsification fails with the identity of the check | Fix the call generator, the probe or the defect form |
| A replay of a case gives another result, because the subject reads a state outside the case | The property engine reports the run as flaky, with the first difference | Build the subject from the handle, and remove the shared state |
| A generated function is called from another goroutine | The property engine reports the run as flaky | A registration of a generator of the parameter whose function generates its results before the call |
| A `budget` directive is malformed | Generation reports `BadBudget` at the directive | Correct the directive |
| A build or a run does not count allocations | Each allocation ceiling reports `unchecked`, with the reason of the assert library | Run `go test -count=1` without the race detector |
| A benchmark measures failing calls alone | `errors/op` is 1 | Prepare the subject in the factory |
| The fuzzer finds a failing input | The engine stores the case under `testdata/prop/<target>`, and go test writes the input under `testdata/fuzz/<target>`. Every later `go test` replays both | Fix the subject, and commit the case as a regression case |
| A TypeScript budget states `allocs` or `bytes` | The ceiling refuses with the reason of the overlay | None. The overlay states the reason |

### Capabilities

| Capability | Verdict | Reason | Reverses when |
|---|---|---|---|
| A generated binding of each claim whose sample binding passes sample inputs to a call | In | Two samples are two points of a domain. The property engine searches the rest and shrinks a failure to its smallest case | The corpus does not show a defect that a generated binding finds and its sample binding misses |
| A sample binding beside each generated binding | In | Each binding has an identity, a proof and an excusal of its own, so a slow subject can excuse the generated binding and keep the samples | — |
| The claims `commutative` and `associative` | In | Their laws concern every pair and every triple of inputs | — |
| Property checks of a consumer | In | A consumer states a relation that the declarations do not show, and gets falsification, a fuzz target and stored cases with it | — |
| Every callable of a suite, and generic declarations at their witnesses | In | The facts of a parameter are the same in every kind of callable | — |
| The starting value of a type: its sample, the value of its builder or its declared default | In | The cases cluster around a value that the declaration accepts, and a counterexample shows only the fields that differ from it | — |
| The setters of the builder for the changed fields of a value | Out | Each changed field takes its value from a fresh value of the shape, which no other value shares, so the copies of the builder add nothing | A builder whose setters validate or normalize a field |
| Overrides of enumerations and builders at every depth | In | The shape reader reads an enumeration of integers over every integer | — |
| Declared errors as the generated errors | In | A caller meets the errors that its dependency declares, such as `ErrNotFound` | — |
| The random source of the case in each double | In | The probabilistic faults of a double replay and shrink with the case | — |
| The builder of the mixin `sample` as the starting value of the inputs of its callable | In | The mixin names the source of the inputs that its callable accepts | — |
| Property checks of the generated builders and doubles | Out | RFC-0004 keeps their tests structural, and samples cover each path of their templates | A defect of a generated builder or double shows at some values alone |
| A receiver that a factory generates from the case | In | The receiver of a value type is most of the input of its methods, and the handle of a case gives the factory the case | — |
| Functions and doubles that generate their results | In | A callable that takes a callback or an interface gets inputs without a registration, and the engine shrinks the results | A callee calls them from other goroutines so often that runs become flaky |
| A fuzz target of each callable of each application | In | Each target has one entry point, and the fuzzer of the language runs the bodies of the property checks | — |
| The campaign profile over the same bodies | In | It searches in every language of assert-spec, with or without a fuzzer | — |
| The directive `budget` with allocations, bytes, p99 and mean | In | It states the costs of a callable as checks with planted defects | — |
| Allocation ceilings over generated inputs in every test run | In | A count of allocations repeats on every machine, and a failure shrinks to the smallest input that allocates more | — |
| A search for the input with the most allocations | In | A search guided by a measured cost finds pathological inputs that random inputs miss, as PerfFuzz and SlowFuzz measured | — |
| Latency ceilings in ordinary test runs | Out | The latency of a call differs per machine and per load | A class of runners repeats each latency within its ceiling |
| A search for the slowest input | Out | A campaign reports its failures, but not the case with the best score of a label | The campaign of assert-spec reports the case with the best score of each label |
| A benchmark of each callable over 100 generated calls | In | A trend needs the same inputs in every run, and one input measures one path | — |
| A warm-up of the benchmarks of a callable in its budget | In | The first calls of a callable build caches that later calls reuse. A ceiling of 0 allocations then needs a warm-up of thousands of calls | — |
| Parallel benchmarks of concurrent callables | In | The cost under contention differs from the serial cost, and the mixins state which callables admit concurrent calls | — |
| A benchmark that prepares the subject through a declared partner | Out | The declarations do not state a relation between the value of a write and the key of a read | The shape catalog states that relation |
| A committed baseline of benchmark results | Out | A measurement of one machine fails on another. The service compares the runs of one class of runners | CI runs the benchmarks on a dedicated runner that repeats within the threshold |
| Classes of complexity over the sizes of inputs | Out | The generated inputs of a run vary in size, and a fit of latencies flakes on shared runners | A consumer needs a ceiling for each element of an input |
| Sequences of calls that the engine chooses | In, as claims of the machinery `generated-sequences` | The judges of a history and the comparison with a reference implementation run over generated sequences | — |
| Comparisons with a reference or a second subject over generated inputs | In, as claims of the machinery `generated-sequences` | A comparison of one call is a sequence of one call | — |

## Alternatives considered

### Separate generators of benchmarks and fuzz targets

A generator of benchmarks and a generator of fuzz targets each take a directive of their own, as testkit separates its benchmarks from its suites.

**Why not:** the declarations, the subjects, the factories and the handles of the suite define where each call goes, and the derivation of the suite defines which laws apply. A second generator would repeat both. A budget is a check that needs a planted defect and the falsification of the suite.

### One binding with the samples as its first cases

The generated binding replaces the sample binding and tries the samples first.

**Why not:** `prop.ForAll` takes one stated case through the option `Draws`, so its body cannot try both samples first. One identity would also hide whether a sample or a generated input failed, and an excusal for a slow subject would excuse the samples as well.

### Zero values as the inputs of benchmarks

Each benchmark passes the zero value of each parameter, as testkit does.

**Why not:** a zero value is one input, and it often takes the path of a miss or an error. A callable that panics on it does not get a benchmark.

### Ceilings through primitives for each shape

A consumer passes a primitive of the shape of a method to an option of the method, such as `ReaderAllocsWithin`, as in testkit.

**Why not:** the consumer writes a call for each method with its type arguments, and nothing proves the ceiling. A directive states the numbers, and the generator writes the calls and their planted defects.

### Fuzz targets of raw bytes

Each fuzz target takes the bytes of the fuzzer and decodes the arguments itself, as fzgen writes its targets for Go.

**Why not:** the bridge of the property engine decodes each input of a fuzzer into a valid case. Typed inputs, shrinking, the store and the replay token then apply to each input that the fuzzer finds, and the same body runs under `ForAll`.

### A fixed seed in every run

The runtime derives the seed of each property check from its key in every run.

**Why not:** the profile `ci` of the property engine fixes the seeds where a run must repeat, and a mutation run fixes them as well. The profile `default` tries new inputs in each local run. One mechanism of seeds then serves the assert library and dokimi.

### Struct values composed field by field

The call generator composes each struct of the workspace from its fields, and takes `prop.Of` for the values of the leaves alone.

**Why not:** the tags of `prop` on the fields of a composed struct would apply only where dokimi reads them again, which repeats the shape reader of the assert library. An override through `Using` keeps that reader for every field.

### A registration of each type for the test process

Each suite package registers its enumerations and its types with a builder through `prop.Register` in an `init` function.

**Why not:** a second registration of a type panics, and a test binary that imports two such suites registers a type twice.

### A budget as a mixin of the shape catalog

The shape catalog of eidos gains a mixin `budget`, and a claim of the spec reads it.

**Why not:** a budget states a cost of an implementation, and the cost differs per subject. The shape catalog classifies the behaviour of a declaration. The verdict reverses when the shape catalog gains classifications of costs.

## Drawbacks

- **Run time.** Each generated binding runs 100 cases on a fresh subject each, about 100 times the calls of its sample binding. A subject that takes a second to build costs 100 seconds for each generated binding, unless its application lowers `cases`.
- **Inputs outside the intended domain.** A generated input that the callable does not mean to accept fails its check until the declaration states the domain. Tags apply to the fields of a struct, so a parameter of a scalar type cannot have one. The callable of such a parameter then accepts every value of its type, or the declaration gives the parameter a type of its own.
- **Ceilings over every input.** An allocation ceiling applies to every generated input. The tags `min_size` and `max_size` of `prop` bound the sizes of the input of a callable whose allocations grow with the input.
- **Requests to assert-spec.** The benchmarks, the allocation ceilings and the property checks wait for them.
- **Unchecked ceilings in CI.** The allocation ceilings need `go test -count=1` without the race detector. A CI without such a run reports every allocation ceiling as unchecked.
- **Sleeps in Go.** A planted defect of a latency ceiling sleeps twice the ceiling for each of its 10 calls.
- **A fuzzer for TypeScript.** A campaign of the fuzzer needs Jazzer.js, which the consumer installs.
- **More generated code.** Each callable gets a call generator, a fuzz target and a level of a benchmark.

## Open questions

None.

## References

| What | Where |
|---|---|
| The purpose and scope of dokimi | [RFC-0001](0001-the-architecture-of-dokimi.md) |
| The conformance spec | [RFC-0002](0002-the-conformance-spec.md) |
| The suite generator | [RFC-0003](0003-the-suite-generator.md) |
| The property engine of assert-go | <https://github.com/dokimasia/assert-go/blob/main/prop/doc.go> |
| The benchmark contract of assert-go | <https://github.com/dokimasia/assert-go/blob/main/bench/doc.go> |
| The fuzz targets of assert-typescript | <https://github.com/dokimasia/assert-typescript/blob/main/src/prop/fuzz.ts> |
| The benchmark contract of assert-typescript | <https://github.com/dokimasia/assert-typescript/blob/main/src/bench.ts> |
| The assertions of assert-spec | <https://github.com/dokimasia/assert-spec/blob/main/spec/assertions.yaml> |
| The mixins of the shape catalog of eidos | <https://github.com/dokimasia/eidos/tree/main/eidos-plugin-shape/spec/mixins> |
| The filenames of the families of a plugin in eidos | <https://github.com/dokimasia/eidos/blob/main/docs/architecture/18-routing-and-layout.md> |
| The generator of benchmarks of testkit | <https://github.com/thesm-os/testkit/blob/v0.10.0/docs/testkit/generators/bench.md> |
| The planned generator of codecs of testkit | <https://github.com/thesm-os/testkit/blob/v0.10.0/docs/testkit/generators/codec.md> |
| The validator of benchmark regressions of testkit | <https://github.com/thesm-os/testkit/blob/v0.10.0/docs/testkit/validators/bench-regression.md> |
| Fuzzing in Go | <https://go.dev/doc/security/fuzz/> |
| Jazzer.js, and its option `fuzzEntryPoint` | <https://github.com/CodeIntelligenceTesting/jazzer.js/blob/main/docs/fuzz-settings.md> |
| Benchmarks in vitest | <https://vitest.dev/guide/benchmarking> |
| Lemieux, Padhye, Sen and Song, PerfFuzz: Automatically Generating Pathological Inputs, ISSTA 2018 | <https://www.carolemieux.com/perffuzz-issta2018.pdf> |
| Petsios, Zhao, Keromytis and Jana, SlowFuzz: Automated Domain-Independent Detection of Algorithmic Complexity Vulnerabilities, CCS 2017 | <https://arxiv.org/abs/1708.08437> |
| The ghostwriter of Hypothesis, which writes property tests from signatures | <https://hypothesis.readthedocs.io/en/latest/reference/integrations.html#ghostwriter> |
| fzgen, which writes fuzz targets for the functions and methods of Go | <https://github.com/thepudds/fzgen> |
