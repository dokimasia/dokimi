---
rfc: 0035
title: The repeated parts of the generated Go code
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-11-04
updated: 2026-11-04
discussion: none
supersedes: RFC-0003, RFC-0004, RFC-0006 and RFC-0014 in part
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0035: The repeated parts of the generated Go code

## Summary

The generated Go code of the corpus repeats facts that a runtime or another part of the same file already contains. Patterns of that kind account for 82,787 of its 299,302 lines. This RFC removes each of them:

- `conform.ExcuseSuite` excuses every check of a suite with one option.
- conform reads the probe, the reason code and the law of a derived check from its claim, so a generated check does not state them.
- `conform.ID.Generated` returns the identity of a generated binding, so a suite does not declare an identity method for one.
- A generated suite has no header comment.
- `stub.Behaviour` and `stub.ResetBehaviour` run the cases that the generated test of a double repeats for each method.
- A generated double has no constructors `New<Double>Delayed` and `New<Double>Failing`.

The checks, their identities, the lock and the suite records do not change. conform adds to its API. stub breaks its API in a minor release, because the module is below 1.0.0 and has no consumers outside dokimi.

## Motivation

The corpus of Go has 299,302 generated lines in 843 files, beside 24,249 lines that people wrote.

| File | Lines | Share |
|---|---|---|
| `<file>_suite.go` | 97,605 | 32.6% |
| `<file>_stub_test.go` | 78,929 | 26.4% |
| `<file>_suite_test.go` | 69,041 | 23.1% |
| `<file>_stub.go` | 42,378 | 14.2% |
| The other generated files | 11,349 | 3.8% |

These patterns repeat what another part of the code contains:

| Pattern | Lines | The repeated fact |
|---|---|---|
| The excusals in the tests of a suite | 34,736 | Each test of a suite excuses each check of the suite by its identity. The tests make 8,684 calls of `conform.Excuse`, and the tests of `partnernaming` alone make 540 |
| The cases of each method in the test of a double | 30,480 | 314 methods repeat the same cases of their stubs: the arguments of a call, the hook `OnRecord`, the reset, the fixed results and the function of `Func` |
| The header comment of a suite | 6,927 | The statement of each check, its depth, the refusals and the unclaimed classifications, which the value of the suite also contains |
| The fields `Probe`, `Reason` and `Law` of a derived check | 4,614 | The entry of the claim in `go.dokimi.dev/spec` |
| The identity methods of generated bindings | 3,729 | 621 methods that each return the identity of a sample binding with the case `generated` |
| The constructors `Delayed` and `Failing` of a double | 2,301 | Two options of the double, which each constructor passes |

The counts come from the corpus at commit 2f3b9c8a. Each line costs a reader of a generated diff, and each line compiles in a test binary. The 208 test binaries of the corpus take 57.5 s to build and 210 s of CPU time on four cores, after a build of the standard library and the runtime modules.

The tests also do work that no check needs. 166 test files build the planted defect of a ceiling as `<Declaration>Suite().Calls[i].Over(h, func() {})`. So each subject of a budget case or of a fuzz input builds every check and every closure of the suite again.

## Detailed design

### The runtime of suites

conform gains one option and one method:

```go
// ExcuseSuite excuses each check of the field Checks of the suite from an
// application. The application skips each such check with the reason. It
// does not excuse a check that With adds or a ceiling of a budget. An
// application fails when the reason is empty, and when the suite has no
// check.
func ExcuseSuite[S any](reason string) Option[S]

// Generated returns the identity of the generated binding of the check with
// the identity id. It adds the case generated to an identity without a case,
// such as Store.Get/keeps-the-value#generated, and the suffix -generated to
// the case of any other identity, such as Store.Get/expires#0-generated.
func (id ID[S]) Generated() ID[S]
```

`Generated` replaces the unexported method `generated`, which the excusal of a sample binding already reads to excuse its generated binding.

A derived check may leave the fields `Probe`, `Reason` and `Law` of `Check[S]` empty. Each claim of the spec has a probe, the reason code of its defect and a law, and conform reads each empty field from the claim through `spec.ClaimOf`:

- The vacuous condition of a check and `Handle.Expect` take the probe of the claim.
- Falsification compares the reason that the probe set with the reason code of the defect of the claim.
- The suite record takes the law of the claim.

A check that states a field keeps its value. The specs of the probes of conform state the probe and the reason of each check that they apply, because their checks do not bind a claim of the spec. A check of the consumer does not have a claim, so conform does not fill an empty field of such a check.

```go
// Probe is the probe that the check calls, whose reasons Handle.Expect
// accepts. A derived check without a probe takes the probe of its claim. A
// check with a probe engages its claim through Handle.Engage alone.
Probe spec.ProbeID

// Reason is the reason code of the planted defect of the check. A derived
// check without a reason takes the reason code of the defect of its claim.
Reason spec.ReasonID

// Law is the law that the check witnesses. A derived check without a law
// takes the law of its claim.
Law spec.LawID
```

An application reports a derived check without a reason code as a problem where the copy of the spec does not contain its claim. The message of the problem contains the claim and the version of the copy.

### The runtime of doubles

`stub.Subject` gains three fields, and its field `Override` passes the call of the arguments:

```go
// Subject is one method of a generated double, as [Behaviour] and
// [ResetBehaviour] drive it. C is the type of the calls of the method, and R
// the type of its results.
//
// Stub, Call and Reset are required. Arguments is nil for a method without
// parameters, and Fails for a method that cannot fail. Result and Returns are
// nil for a method without results, and Override is nil for a method whose
// function a check cannot set. Behaviour and ResetBehaviour leave out each
// check that needs a nil field.
type Subject[C, R any] struct {
	// Stub is the stub of the method on the double.
	Stub *MethodStub[C]
	// Call calls the method with the samples of its arguments, and drops its
	// results.
	Call func()
	// Arguments checks on tb that the call c has the samples of the arguments
	// of Call.
	Arguments func(tb assert.TB, c C)
	// Reset resets the double, as its method ResetCalls does.
	Reset func()
	// Fails calls the method and returns its error.
	Fails func() error
	// Result calls the method and returns its results.
	Result func() R
	// Default returns the results of an unconfigured call, and is nil where
	// they are the zero results.
	Default func() R
	// Returns sets the fixed results of the method to samples, and returns
	// the samples.
	Returns func() R
	// Override sets a function of the method that passes the call of its
	// arguments to seen and returns the zero results.
	Override func(seen func(c C))
}

// ResetBehaviour checks the reset of a generated double through the method of
// the subject that newSubject returns. Each check is a parallel subtest of t.
// name is the name of the method in each message.
func ResetBehaviour[C, R any](t *testing.T, name string, newSubject func(tb assert.TB) Subject[C, R])
```

`Behaviour` keeps its signature and runs these checks as parallel subtests:

| Check | Fields | Status |
|---|---|---|
| records the call | — | Kept |
| records the arguments of the call | `Arguments` | New |
| runs the OnRecord hook once for each call | — | New |
| reports an unmet call count | — | Kept |
| reports an unmet minimum call count | — | Kept |
| refuses an unconfigured call in strict mode | — | Kept |
| returns the default results of an unconfigured call | `Result` | Kept |
| returns the fixed results of Returns | `Result` and `Returns` | New. It first checks that the samples differ from the default results, because a check of equal values passes for a double that ignores Returns |
| accepts a call with fixed results in strict mode | `Returns` | New |
| passes the arguments of each call to the function of Func | `Override` and `Arguments` | Replaces the check that the double calls the function |
| returns an injected fault | `Fails` | Kept |
| fails the call that the count of Faults names | `Fails` | Kept |
| records the injected fault on the call | `Fails` | Kept |

`ResetBehaviour` takes the subject of a method in place of an `Instance`, and runs two checks:

| Check | Fields |
|---|---|
| clears the calls of the method and their instants | — |
| keeps the fixed results of Returns | `Result` and `Returns` |

`Instance[C]` loses its field `Reset`, because no check of an option reads it.

### The generated suite

The suite of `Store` has no header comment. The value of `StoreSuite()` contains each fact of the header: the statement of each check, the field `Shallow`, each refusal with its reason and its remedy, and each unclaimed classification with its stance and its note. The lock lists the key of each check, so a lost check still shows in review as a deleted line.

`StoreChecks` declares an identity method for each check, except for a generated binding beside a sample binding. The value of the suite names such a binding through `Generated`. A claim that binds under `generated-inputs` alone, such as `commutative`, keeps its identity method. The value of a derived check does not have the fields `Probe`, `Reason` and `Law`:

```go
{
	ID:          ContractChecks.Put.Smoke().Generated(),
	Statement:   "Put returns without a panic for derived inputs",
	Claim:       spec.Smoke,
	Machineries: []spec.MachineryID{spec.GeneratedInputs, spec.NothingBeyondACaller},
	Property: func(h *ContractHandle, c *prop.Case) {
		drawnV := c.Draw(prop.Of[ifabsent.Value](), "v")
		probe.Survives(h, func(ctx context.Context) error {
			return h.Subject().Put(ctx, drawnV)
		})
	},
	Defect: contractPutPanics,
},
```

A test that excuses a generated binding alone, such as the test of a subject under the directive `excuse`, passes `conform.Excuse(storetest.StoreChecks.Get.ZeroOnError().Generated(), reason)`.

### The generated tests of a suite

Each test file of a suite builds the suite once, in a variable of the package. Each test of the file that excuses the checks of the suite passes `ExcuseSuite` once:

```go
// contractSuite is the suite of Contract that the tests of this file apply.
var contractSuite = ifabsenttest.ContractSuite()

// contractIdle returns the planted defect of a ceiling of Contract with a
// spend of nothing, which returns the zero results of each call at once.
func contractIdle(h *ifabsenttest.ContractHandle) ifabsent.Contract {
	return contractSuite.Calls[0].Over(h, func() {})
}

// TestContractSuiteBudgets applies the allocation ceilings of the suite of
// Contract to contractIdle with every check of the suite excused, and a
// ceiling of 0 to each callable without a budget whose calls build no input
// of their own.
func TestContractSuiteBudgets(t *testing.T) {
	ifabsenttest.RunContractBudgets(
		t,
		contractIdle,
		conform.Ceilings(ifabsenttest.ContractChecks.Put.Callable, conform.Budget{Allocs: new(uint64(0))}),
		conform.ExcuseSuite[ifabsent.Contract]("the planted defect of a ceiling returns the zero results of each call"),
	)
}
```

The test of the two forms of the application, the test of the budgets and the fuzz target of each callable make the same excusal with the reasons of today.

### The generated double

The double of `Store` loses `NewStoreStubDelayed(tb, d, clk, opts...)` and `NewStoreStubFailing(tb, err, opts...)`. A test passes the options that each constructor passed:

```go
s := storetest.NewStoreStub(t, storetest.StoreStubWithClock(clk), storetest.StoreStubWithLatency(d))
f := storetest.NewStoreStub(t, storetest.StoreStubFailing(err))
```

No code of the corpus outside the files of the doubles calls either constructor.

### The generated test of a double

The subject of each method sets the new fields:

```go
// contractStubPutSubject returns the subject of [stub.Behaviour] and
// [stub.ResetBehaviour] for Put. The subject calls a new double that
// reports to tb.
func contractStubPutSubject(tb assert.TB) stub.Subject[ifabsenttest.ContractPutCall, ifabsenttest.ContractPutReturn] {
	tb.Helper()
	s := ifabsenttest.NewContractStub(tb)
	return stub.Subject[ifabsenttest.ContractPutCall, ifabsenttest.ContractPutReturn]{
		Stub:  s.OnPut.MethodStub,
		Reset: s.ResetCalls,
		Call: func() {
			_ = contractStubPut(tb, s)
		},
		Arguments: func(tb assert.TB, c ifabsenttest.ContractPutCall) {
			tb.Helper()
			a0, a1 := contractStubPutArgs(tb)
			expect.Equal(tb, c.Ctx, a0, "the call records Ctx", expect.ByIdentity())
			expect.Equal(tb, c.V, a1, "the call records V")
		},
		Result: func() ifabsenttest.ContractPutReturn {
			got0 := contractStubPut(tb, s)
			return ifabsenttest.ContractPutReturn{Err: got0}
		},
		Returns: func() ifabsenttest.ContractPutReturn {
			want0 := contractStubPutPins()
			s.OnPut.Returns(want0)
			return ifabsenttest.ContractPutReturn{Err: want0}
		},
		Override: func(seen func(ifabsenttest.ContractPutCall)) {
			s.OnPut.Func(func(p0 context.Context, p1 ifabsent.Value) error {
				seen(ifabsenttest.ContractPutCall{Ctx: p0, V: p1})
				var z0 error
				return z0
			})
		},
		Fails: func() error {
			r0 := contractStubPut(tb, s)
			return r0
		},
	}
}
```

The generator leaves the field `Arguments` out for a method without parameters. It leaves the field `Returns` out for a method without a sample of a result and without an error result.

The test of the double runs `Behaviour` in the level of each method, and `ResetBehaviour` for each method in the level `ResetCalls`:

```go
t.Run("Put", func(t *testing.T) {
	t.Parallel()
	stub.Behaviour(t, "Put", contractStubPutSubject)
})

t.Run("ResetCalls", func(t *testing.T) {
	t.Parallel()
	stub.ResetBehaviour(t, "Put", contractStubPutSubject)
})
```

These cases of the generated test move into the runtime:

| Case | Level today | Check of the runtime |
|---|---|---|
| records the arguments of the call | Each method | `Behaviour` |
| runs the OnRecord hook once for each call | Each method | `Behaviour` |
| clears the recorded calls of the method | `ResetCalls` | `ResetBehaviour` |
| keeps the results of Returns of the method | `ResetCalls` | `ResetBehaviour` |
| sets the results of each call of the method | `Returns` of the stub of the method | `Behaviour` |
| configures the method for strict mode | `Returns` of the stub of the method | `Behaviour` |
| passes the arguments of each call of the method to the function | `Func` of the stub of the method | `Behaviour` |

The levels `New<Double>Delayed` and `New<Double>Failing` go with their constructors. The levels of the options keep their cases, because each option builds a double of its own: `Strict`, `DelegateTo`, `WithClock`, `WithRandSource`, `BenchMode`, `WithLatency`, `Failing` and the option of each method. The test of the stub of a method keeps the levels of its other methods, such as `Yields` and the helper of each error value. The stub of a method without such a method has no test of its own.

### The size of the change

| Pattern | Lines removed | Lines added |
|---|---|---|
| The excusals in the tests of a suite | 34,736 | One option for each test that excuses the checks of a suite |
| The cases of each method in the test of a double | 30,480 | Up to three fields in each of the 314 subjects of a method, and one call of `ResetBehaviour` for each method |
| The header comment of a suite | 6,927 | None |
| The fields `Probe`, `Reason` and `Law` | 4,614 | None |
| The identity methods of generated bindings | 3,729 | None |
| The constructors `Delayed` and `Failing` | 2,301 | None |
| The builds of a suite for each subject | 166 calls | One variable in each of 166 test files |

The corpus that this design generates has 222,864 generated lines in the same 843 files, 76,438 lines fewer.

### Migration

1. dokimi-go changes conform and stub, with their specs and changesets. Each module keeps its coverage and its mutation score at 100%. The release of each module raises its minor version, for the new API of conform and for the breaking changes of stub.
2. dokimi changes the templates of the suite generator and of the stub generator, their golden files, and the specs that pin the header. It regenerates the corpus with a binary of the working tree. Each module of dokimi requires the new releases through `go get` and `go mod tidy`.

The second step needs the push of the first, because the generated code calls the new API.

## Alternatives considered

### One binding with the samples as its first cases

The generated binding replaces the sample binding and runs the samples as its first cases. The corpus would lose 621 checks with their identity methods, their values in the suite and their excusals.

**Why not:** each binding has a verdict and an excusal of its own. A subject whose calls take a second each can excuse the generated binding and keep the samples. One identity would also hide whether a sample or a generated input failed. The excusal of a sample binding already excuses its generated binding, so a test that excuses a claim names it once.

### Applications through conform alone

A test calls `conform.Run(t, storetest.StoreSuite(), newSubject)`. The suite does not declare entry points.

**Why not:** Go does not accept a type parameter as the constraint of another type parameter. The compiler refuses `func run[S any, T S](newSubject func() T) S` with "cannot use a type parameter as constraint". So `conform.Run` cannot take a factory of `*InMemory` for the suite of `Store`. Each test would write a function literal that converts its subject. The entry points of a suite are generic over a type that implements the declaration, so `RunStoreSuite(t, NewInMemory)` compiles.

### One form of the application

`RunStoreSuite` takes a factory of the handle alone. The suite does not declare `RunStoreSuiteFromHandle`.

**Why not:** most constructors do not read the handle. Each test would then wrap its constructor in a function literal. The two forms take 17 lines in the suite of a declaration.

### A header without the statements

The header lists the identity of each check, without its statement.

**Why not:** the value of the suite contains the identity of each check beside its statement, and the lock contains the key of each check. The header would be a third list of the same checks. godoc renders only the comments directly above declarations, so the header has no reader outside the file.

### No fields for the probe, the reason and the law

`Check[S]` loses the fields `Probe`, `Reason` and `Law`, and conform always reads them from the claim.

**Why not:** the specs of the probes of conform apply each probe through a check that states the probe and its reason. Some of those checks do not bind a claim, such as the check of the vacuous condition of a probe, which runs without a planted defect. Without the fields, each such check would need a claim of the spec and a planted defect that fails it.

### The options of a double in the runtime

`stub.WithClock(clk)` configures any double, so a double does not declare an option for a setting of its group.

**Why not:** Go cannot infer the type of the double from the arguments of such an option. Each call would then name the type, as in `stub.WithClock[*storetest.StoreStub](clk)`. A method of the double that returns its group would need a name that no method of a declaration has.

### A helper of each test file that returns the excusals

The generator writes a function that returns the excusal of each check of the suite for a reason, and each test calls it.

**Why not:** the function still lists each check of the suite in each test file, and each new check of the suite changes it. The option `ExcuseSuite` states the excusal once, in the runtime.

### Identities as fields

`StoreChecks.Get.ZeroOnError` is a field of the type `conform.ID[store.Store]`, which halves the lines of an identity.

**Why not:** a test can assign a field of the exported variable `StoreChecks`, and each later read of the identity then takes the assigned value. A method returns one identity.

## Drawbacks

- **Breaking changes in stub.** `stub.Instance` loses its field `Reset`, `stub.Subject.Override` changes its signature, and `stub.ResetBehaviour` takes the subject of a method.
- **Two sources of the probe of a derived check.** A derived check that states a probe, a reason code or a law keeps it, also where its claim has another one. The generator states none of them.
- **A longer name of a generated binding.** A test names `StoreChecks.Get.ZeroOnError().Generated()` in place of `StoreChecks.Get.ZeroOnErrorGenerated()`.
- **No summary at the top of a suite.** A reader finds the checks of a suite in the value of the suite or in the lock.
- **A type argument at each excusal of a suite.** `conform.ExcuseSuite[store.Store](reason)` names the subject type, because no argument of the option has that type.
- **Failures in the runtime.** A failed case of a double reports a line of `stub` and the name of the method, in place of a line of the generated test.
- **More runs of the reset checks.** `ResetBehaviour` runs once for each of the 314 methods of the corpus. Today the check of the instants runs once for each of the 152 doubles.

## Open questions

None.

## Unresolved and future work

- The TypeScript binding of the suite and of the doubles is not proposed here.
- The race run of the corpus with `-p=1` takes 245 s, and the race runs of the specs of `generator/suite`, `cmd/dokimi` and `generator/sim` take 400 s. Both cost more time in CI than the build of the generated code, and neither is addressed here.

## References

| What | Where |
|---|---|
| The type parameters of a declaration in the Go specification | <https://go.dev/ref/spec#Type_parameter_declarations> |
| The property engine of assert-go | <https://github.com/dokimasia/assert-go/blob/main/prop/doc.go> |
