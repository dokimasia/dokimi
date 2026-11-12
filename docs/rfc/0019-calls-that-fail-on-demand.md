---
rfc: 0019
title: Calls that fail on demand
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-28
updated: 2026-10-31
discussion: none
supersedes: none
superseded-by: RFC-0050 in part
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0019: Calls that fail on demand

## Summary

A subject often depends on another declaration: a cache on its backing store, a breaker on the operation that it guards, a balancer on its instances. Laws about how the subject handles a failure of a dependency need a call of the dependency that fails when the check asks. conform gains the `Injector`, which fails and records the calls of the dependencies of the subjects of one handle, and the door `Faults` with its option `conform.Faultable`. A declaration states the machinery with `//dokimi:suite machineries=[faults]`. The consumer builds a fresh subject over dependencies that ask the injector. A handle builds the subjects of every check of the machinery through the door, so a derived check and its planted defect do not need another path. The machinery licenses the claims of a reported failure, a failed call without effect, a retry and a transaction under faults, a breaker that opens, a key that keeps its instance, and a release on every failure path.

## Motivation

The fault strategies exist. The package `fault` of dokimi-go has counted, retry, probabilistic, windowed and predicate faults, composable with `And` and `Or`. The stub generator arms them on each method of a double through `Faults`, `FaultsOnce` and `SetFault`. A fault armed on a method without an error result fails the test, so an injection that injects nothing cannot pass.

The doors of conform are a reference, a second subject, a restart and a crash, and none of them makes a call fail. No fixture of the corpus has a door that fails a call. A cache that serves a stale value after its backing fails passes the checks of a single caller when the backing never fails, and so does a writer that applies half a batch before an error.

The shape catalog states the laws that need a failing dependency. The mixin `atomic` states that a callable completes fully or has no observable side effect. The contract `circuit-breaker` states that a callable fails fast after repeated failures of its operation. The mixin `sticky` states that the inputs with one key go to one instance while the instance is healthy. The mixin `leak-free` states that a callable releases each resource that it acquires, on the paths of a failure as well. The spec gives `circuit-breaker` the stance `needs` today, and it has no claim for the three mixins.

RFC-0002 names this machinery `call-that-fails-on-demand`, with the param `faults`. It arrives as the clock and the doors arrived. The subject declares the capability through an option, the check demands it, and a missing option fails the check by name.

## Detailed design

### The injector

```go
// ErrInjected is the error of a call of a dependency that an Injector failed.
// The probes of the machinery call-that-fails-on-demand inject it.
var ErrInjected = errors.New("conform: an injected failure of a dependency")

// Injector fails the calls of the dependencies of the subjects of one handle
// on the request of a probe, and records each call of a dependency.
// The consumer builds each subject over dependencies of the injector, through
// FaultOf for a double and through Dependency.Call for a dependency that the
// consumer writes. Its methods are safe for concurrent use.
type Injector struct{ /* unexported */ }

// FailNext makes the next call of a dependency fail with err.
func (i *Injector) FailNext(err error)

// FailAt makes the n-th call of a dependency from now fail with err. FailAt(1,
// err) arms the same failure as FailNext(err).
func (i *Injector) FailAt(n int, err error)

// FailUntilHealed makes every call of a dependency fail with err until Heal.
// It takes precedence over FailNext and FailAt.
func (i *Injector) FailUntilHealed(err error)

// Heal clears every failure that FailNext, FailAt and FailUntilHealed armed.
func (i *Injector) Heal()

// Hold makes each call of a dependency wait until Release.
func (i *Injector) Hold()

// Release lets each call that Hold keeps waiting go on, and every later call
// pass without waiting.
func (i *Injector) Release()

// Failures returns the number of calls of a dependency that the injector
// failed.
func (i *Injector) Failures() int

// Calls returns each call of a dependency, in the order of the calls.
func (i *Injector) Calls() []DependencyCall

// Dependency returns a new dependency of the injector, such as one instance
// behind a balancer. Each dependency has a number of its own, from 1.
func (i *Injector) Dependency() *Dependency

// DependencyCall is one call of a dependency that an injector recorded.
type DependencyCall struct {
	// Dependency is the number of the dependency.
	Dependency int

	// Callable is the callable of the dependency, such as Get.
	Callable string

	// Failed reports that the injector failed the call.
	Failed bool
}

// Dependency is one dependency of the subjects of an injector.
type Dependency struct{ /* unexported */ }

// Call records a call of the callable of the dependency, and returns the
// error that the call fails with, or nil for a call that succeeds. A
// dependency that the consumer writes calls it at the start of each call.
func (d *Dependency) Call(callable string) error

// FaultOf returns the fault strategy of the method callable of a double of
// the dependency d. The consumer arms it on the method through SetFault, so
// each call of the method calls d.Call. Reset of the fault does nothing,
// because the injector keeps the state of every fault.
func FaultOf[C any](d *Dependency, callable string) fault.Fault[C]
```

`Hold` and `Release` keep a computation in flight while other callers arrive, which the claim of the contract `single-flight` of RFC-0021 needs. A call of a dependency cannot see the context of its caller through a fault strategy, so a call that `Hold` keeps waiting goes on at `Release` alone.

`FaultOf` returns a `fault.Fault[C]`, so the module `go.dokimi.dev/conform` requires `go.dokimi.dev/fault`. The package documentation of conform states the new dependency position.

### The door and the subjects of a check

```go
type Doors[S any] struct {
	// The doors Reference, Second, Restart and Crash as before.

	// Faults returns a fresh subject whose dependencies ask inj, which
	// Handle.Subject, Handle.New and Handle.Fresh return for a check of the
	// machinery call-that-fails-on-demand.
	Faults func(h *Handle[S], inj *Injector) S
}

// Faultable supplies the subjects of an application whose dependencies fail
// on demand.
func Faultable[S any](newFaulty func(h *Handle[S], inj *Injector) S) Option[S]

// Injector returns the injector of the handle, which every subject that the
// handle builds through the door Faults asks.
func (h *Handle[S]) Injector() *Injector
```

A handle builds the subjects of a check through the door `Faults` when the machineries of the check include `call-that-fails-on-demand`:

- `Subject` and `New` build through the door where the application builds through its factory.
- `Fresh` always builds through the door, so a planted defect over a fresh subject wraps a subject whose dependencies fail on the request of the probe.
- Under falsification the planted defect remains the source of `Subject` and `New`.
- Each handle has one injector, which every subject that it builds shares. Each case of a property check has a handle and an injector of its own.

The door resolves as every door of RFC-0008 resolves: the door of the check of a planted defect first, then the option of the application. Without the door, the call ends the check with a fault that wraps `ErrOption` and names `conform.Faultable`:

```text
conform: the handle lacks the door: Cache.Get/failure-is-reported builds a subject without the option
conform.Faultable. Pass the option to the application:
	conform.Faultable(func(h *conform.Handle[cache.Cache], inj *conform.Injector) cache.Cache { … })
```

A derived check binds its calls to `h.Subject()` and its planted defect wraps `h.Fresh()`, so neither the template of a check nor the template of a defect changes. A check of the consumer gets the same subjects when it lists the machinery in its field `Machineries`.

A consumer of a cache in front of a store passes the option through the options function of its directive `subject`, as RFC-0008 states:

```go
// CacheOptions returns the door of a cache over a store whose calls fail on
// demand.
func CacheOptions() []cachetest.CacheOption {
	return []cachetest.CacheOption{
		conform.Faultable(func(h *conform.Handle[cache.Cache], inj *conform.Injector) cache.Cache {
			store := inj.Dependency()
			backing := storetest.NewStoreStub(h, storetest.StoreStubDelegateTo(memory.New()))
			backing.OnGet.SetFault(conform.FaultOf[storetest.StoreGetCall](store, "Get"))
			backing.OnPut.SetFault(conform.FaultOf[storetest.StorePutCall](store, "Put"))
			return cache.New(backing)
		}),
	}
}
```

### The machinery of a declaration

The directive `suite` gains the param `machineries`, a list of the params of the machineries that the consumer supplies: `faults`, `reference`, `second` and `medium`. The suite generator passes each machinery of the list to the derivation. The derivation lists the classifications of a claim that needs another machinery as unclaimed, with the stance `needs` and a note with the param, such as `machineries=[faults]`:

```go
//dokimi:suite machineries=[faults]
type Cache interface { … }
```

A value that is the param of no machinery is an Error at the directive under the new code `UnknownMachinery`, number 13, whose message lists the four params. The suite generator derives the suite without the value.

### The door of a check of the consumer

The directive `check` gains the param `faults`, the door of the planted defect for this machinery, as RFC-0008 gives each door a param. Its function has the signature `func(h *CacheHandle, inj *conform.Injector) Cache`. The annotate step stamps it under the key `suite.faults`, and it reports `BadCheck` for another number of parameters or results. The generator writes it into the field `Doors.Faults` of the check. The check also lists the machinery `call-that-fails-on-demand` in its field `Machineries`, so a handle builds its subjects through the door.

### The claims

| Claim | Licence | Partners and params | Gates | Law | Probe | Reason | Defect form | Entails |
|---|---|---|---|---|---|---|---|---|
| `failure-is-reported` | Every callable | — | `failure` | `failed-dependency-is-reported` | `reports-a-failed-dependency` | `swallowed-failure` | `swallows-failure` | — |
| `no-partial-write` | The mixin `atomic` with its param `read` | `read`: named by `read` | `failure` | `failed-call-leaves-no-effect` | `fails-without-effect` | `partial-write` | `half-applied` | — |
| `retry-succeeds-under-faults` | The mixin `retry-succeeds` with its param `attempts`, as `retry-succeeds` | `attempts` | `cancellation`, `failure` | `call-converges-within-the-declared-attempts` | `succeeds-after-a-failure` | `never-converged` | `gives-up-early` | `retry-succeeds` |
| `transaction-under-faults` | The contract `transaction`, role `fn` | `write` and `read`: named by their params, and `not-found` | `cancellation`, `failure` | `failed-commit-leaves-nothing` | `rolls-back-under-faults` | `not-rolled-back` | `commits-on-failure` | — |
| `fails-fast` | The contract `circuit-breaker`, role `fn` | — | `failure` | `open-breaker-skips-the-operation` | `refuses-without-the-dependency` | `never-opened` | `never-opens` | — |
| `keeps-its-instance` | The mixin `sticky` with its param `key` | — | `returns-unaided` | `one-key-one-instance` | `uses-one-instance` | `moved-instance` | `balances-each-call` | — |
| `releases-on-every-path` | The mixin `leak-free` with its params `open` and `close` | `open` and `close` | `failure` | `failure-path-releases` | `releases-what-it-acquires` | `leaked-resource` | `opens-without-closing` | — |

The statement of each claim:

| Claim | Statement |
|---|---|
| `failure-is-reported` | `{callable} returns an error when a call of its dependency fails` |
| `no-partial-write` | `{read} finds no effect of a {callable} that failed with its dependency` |
| `retry-succeeds-under-faults` | `{callable} succeeds within the declared attempts after a call of its dependency fails` |
| `transaction-under-faults` | `{read} finds nothing that {callable} wrote through {write} before its commit failed` |
| `fails-fast` | `{callable} refuses a call without its operation after the operation failed repeatedly` |
| `keeps-its-instance` | `{callable} sends each call with one key to one instance` |
| `releases-on-every-path` | `{callable} releases each resource that it acquired when a call of its dependency fails` |

Each claim lists `call-that-fails-on-demand` first among its machineries. A declaration without `machineries=[faults]` lists the classifications of these claims as unclaimed with the stance `needs`, and the note gives the param. `failure-is-reported` licenses every callable, as `smoke` does, so it has no classification to list. The derivation does not need a fact about the dependencies of a subject, because the declaration states that its consumer builds the subject over dependencies that fail.

A callable that never calls a dependency, such as the length of a cache, gets a check of `failure-is-reported` that never engages, because the claim licenses every callable that can fail. The consumer excuses such a check.

The law of `retry-succeeds`, `call-converges-within-the-declared-attempts`, already has the depth `call-that-fails-on-demand`, so its checks of a single caller are shallow. `retry-succeeds-under-faults` witnesses the law in full and entails `retry-succeeds`. `fails-always` fails the deeper probe. `gives-up-early` passes the shallower probe, whose dependency never fails, so the proof of the pair succeeds in both directions.

`transaction-under-faults` does not entail a claim. A failed commit and a failed body are distinct causes of a failed scope. A scope that rolls back a failed commit and keeps what a failed body wrote passes this claim and fails `transaction`, so neither law implies the other. `transaction` keeps the full depth of its law `failed-scope-leaves-nothing`.

### The laws, the reasons and the defect forms

The spec gains six laws, each with the depth `call-that-fails-on-demand`:

| Law | Statement |
|---|---|
| `failed-dependency-is-reported` | A call returns an error when a call of its dependency fails. |
| `failed-call-leaves-no-effect` | A reader finds no effect of a call that failed. |
| `failed-commit-leaves-nothing` | A read finds nothing that a scope wrote before its commit failed. |
| `open-breaker-skips-the-operation` | After repeated failures of its operation, the callable refuses a call without calling the operation. |
| `one-key-one-instance` | Calls with one value of the key go to one instance while it is healthy. |
| `failure-path-releases` | The callable releases each resource that it acquired, on every path of a failure. |

`failed-commit-leaves-nothing` is distinct from `failed-scope-leaves-nothing`, with the case of a scope that rolls back a failed commit and keeps what a failed body wrote.

The spec gains four reasons. The meaning of `partial-write` widens to every failed write, because the refused write of its first claim is a failed write too.

| Reason | Meaning |
|---|---|
| `swallowed-failure` | The call returned no error while a call of its dependency failed. |
| `partial-write` | The reader found a part of a write that failed. |
| `never-opened` | The callable called its failing operation on every call. |
| `moved-instance` | Two calls with one key went to two instances. |
| `leaked-resource` | A call acquired more resources than it released. |

The spec gains seven defect forms:

| Form | Changes | Does | Minimum steps |
|---|---|---|---|
| `swallows-failure` | callable | The callable returns its zero results and no error where a call of a dependency fails during the call. | 1 |
| `half-applied` | `partner:read` | The read returns the second sample and no error once a call of a dependency failed. | 2 |
| `gives-up-early` | callable | The callable fails on every call once a call of a dependency failed. | 2 |
| `commits-on-failure` | `partner:read` | The read returns the second sample and no error once a call of a dependency failed. | 2 |
| `never-opens` | callable | Each call goes to a fresh subject, whose breaker is closed. | 2 |
| `balances-each-call` | callable | Each call goes to a fresh subject over new dependencies. | 2 |
| `opens-without-closing` | callable | A call that fails records one more acquisition of its own, which no release follows. | 1 |

The unclaimed table of the spec loses the entries of `mixin:atomic`, `mixin:sticky` and `mixin:leak-free`, and the stance of `contract:circuit-breaker/fn`. The claims of this RFC license them, and the derivation gives a declaration without the machinery the stance `needs`.

### The probes

The package `probe` of conform gains seven probes, which follow the conventions of RFC-0027. Each engages only once the subject calls a dependency, and its vacuous condition states the run in which the subject calls none:

```go
// ReportsAFailedDependency makes every call of a dependency fail, calls call,
// and fails when the call returns no error.
func ReportsAFailedDependency[S any](h *conform.Handle[S], call Act)

// FailsWithoutEffect reads through read and makes every call of a dependency
// fail. It then calls call, calls Heal of the injector and reads again. It
// fails when the second read returns a value other than the first. A call that
// succeeds ends the probe.
func FailsWithoutEffect[S, V any](h *conform.Handle[S], call Act, read Call[V])

// SucceedsAfterAFailure fails the next call of a dependency, then calls call
// up to attempts times, and fails when no call succeeds.
func SucceedsAfterAFailure[S any](h *conform.Handle[S], call Act, attempts int)

// RollsBackUnderFaults runs scope with a body that writes and then fails the
// next call of a dependency, so the commit of the scope fails. It calls Heal
// of the injector, reads, and fails unless the error of read matches sentinel.
// A write inside the scope that fails, and a scope that succeeds, end the
// probe.
func RollsBackUnderFaults[S any](
	h *conform.Handle[S], scope func(ctx context.Context, body Act) error, write, read Act, sentinel error,
)

// RefusesWithoutTheDependency makes every call of a dependency fail and calls
// call up to 100 times. It fails when the subject calls a dependency during
// each of the calls.
func RefusesWithoutTheDependency[S any](h *conform.Handle[S], call Act)

// UsesOneInstance calls call twice, and fails when the subject calls two
// dependencies for the two calls.
func UsesOneInstance[S any](h *conform.Handle[S], call Act)

// ReleasesWhatItAcquires calls call once for each step n from 1, with the
// n-th call of a dependency failing. acquire and release are the callables of
// the dependencies that the params open and close of the mixin leak-free
// name. It fails when the calls of release differ in number from the calls of
// acquire that succeeded, after a step. A call of acquire that failed acquired
// nothing. It stops after the first step whose call makes fewer than n calls
// of a dependency, and after 64 steps.
func ReleasesWhatItAcquires[S any](h *conform.Handle[S], call Act, acquire, release string)
```

| Probe | Inputs | Reasons | Vacuous condition |
|---|---|---|---|
| `reports-a-failed-dependency` | `call`: call | `swallowed-failure` | The subject calls no dependency during the call. |
| `fails-without-effect` | `call`: call, `read`: call | `partial-write` | The call succeeds, or the subject calls no dependency during it. |
| `succeeds-after-a-failure` | `call`: call, `attempts`: count | `never-converged` | The subject calls no dependency during the attempts. |
| `rolls-back-under-faults` | `call`: scope, `write`: call, `read`: call, `sentinel`: error | `not-rolled-back` | The write inside the scope fails, the scope succeeds, or its commit calls no dependency. |
| `refuses-without-the-dependency` | `call`: call | `never-opened` | The subject calls no dependency during the first call. |
| `uses-one-instance` | `call`: call | `moved-instance` | The subject calls no dependency during a call. |
| `releases-what-it-acquires` | `call`: call, `open`: value, `close`: value | `leaked-resource` | The subject calls no dependency during the call. |

The generated call of `uses-one-instance` binds the same inputs on both calls, so both calls pass one value of the key. The values `open` and `close` are the names of the callables that the params of `leak-free` refer to.

### The planted defects

Each defect form wraps a fresh subject of the application, as RFC-0014 states. `swallows-failure`, `gives-up-early`, `half-applied` and `commits-on-failure` read the injector of the handle. `never-opens` and `balances-each-call` build a subject of the handle for each call, and `opens-without-closing` acquires on a dependency of its own. The type of their defect gains a field for the handle or the dependency. The constructor of the defect sets the field from the handle that builds the defect, so the probe and the defect share one injector:

| Form | Body of the override |
|---|---|
| `swallows-failure` | Calls the fresh subject. Where the call fails and the injector failed a call of a dependency during it, it returns its zero results and no error |
| `gives-up-early` | Calls the fresh subject, and fails once `Failures()` of the injector is above zero |
| `half-applied`, `commits-on-failure` | The read calls the fresh subject. Where `Failures()` of the injector was above zero before the read, it returns its second sample and no error |
| `never-opens`, `balances-each-call` | Calls a subject of `h.Fresh()` for each call, so each call goes to a closed breaker and new dependencies |
| `opens-without-closing` | Calls the fresh subject. Where the call fails, it records a call of `open` on a dependency of its own, which no call of `close` follows |

The corpus runs every other check on the planted defect of a check, and the defect has to pass each of them. `swallows-failure` returns the error of a call whose dependency did not fail, such as the error of a missing key. The read of `half-applied` and `commits-on-failure` samples the injector before its call, so a read whose own dependency fails reports the failure. The form `keeps-failed-writes` of the claim `transaction` returns its zero results and no error until a call of a dependency fails, and calls the fresh subject after that. Its defect passes `transaction-under-faults`, whose commit fails before the read.

### In TypeScript

The TypeScript runtime declares the parts of this machinery in `@dokimi/runtime/conform/injector`: the class `Injector` with `failNext`, `failAt`, `failUntilHealed`, `heal`, `hold`, `release`, `failures`, `calls` and `dependency`, the class `Dependency` with `call(callable)`, the error class `InjectedError`, and `faultOf(dependency, callable)`, which returns a fault of `@dokimi/runtime/fault`. The options of `run` gain `faultable`, a function of the handle and the injector that returns a subject, and `h.injector()` returns the injector of a handle. The probes are the modules `@dokimi/runtime/probe/<probe>` of their identifiers.

## Alternatives considered

### Failpoints in the source of the subject

The subject marks the places that can fail, as `pingcap/failpoint` and FreeBSD `fail(9)` mark them, and a check enables a marker by name.

**Why not:** a failpoint is code in the subject that a build step rewrites. A generated check would demand markers that the consumer writes by hand at the right places. A failure of a dependency needs no change of the subject, because the consumer already builds the subject over its dependencies.

### Faults at the network

A proxy such as Toxiproxy fails the connections of the subject.

**Why not:** it needs a subject that talks over a network and a proxy process in the test. Most subjects of a suite are in-process types. A proxy belongs to the chaos experiments against running systems, outside the suite.

### A method of the handle for a faulty subject

A check calls `h.Faulty()`, which returns a subject over failing dependencies and its injector.

**Why not:** a derived check binds its calls to `h.Subject()`, and its planted defect wraps `h.Fresh()`. Each template of a call and each defect form would need a second path to a faulty subject, and the planted defect would need a method that bypasses its own door. The machineries of a check already state that it needs the door.

### An injector per dependency

Each double of a dependency gets an injector of its own.

**Why not:** a probe fails whichever dependency the subject calls first, and `keeps-its-instance` compares the instances that the subject calls for two calls. Both need one record of the calls of every dependency of the subject.

## Drawbacks

- **The consumer wires each dependency.** A subject gets the door only when the consumer builds it over dependencies of the injector: one `Dependency` for each instance and one `FaultOf` for each method of a double, about two lines for each method that can fail.
- **Vacuous checks.** `failure-is-reported` licenses every callable that can fail, so each callable that never calls a dependency has a check that never engages and needs an excusal.
- **Defects through the injector.** Each of the seven defect forms reads the injector, calls a dependency of the injector or rebuilds the subject, so it breaks the law of its claim from outside the state of the subject. A defect proves that its check fails for the reason of its claim. It does not show how a real subject breaks the law.
- **A module dependency.** The module `go.dokimi.dev/conform` requires `go.dokimi.dev/fault`.
- **A machinery that the declaration states.** A declaration whose subject has dependencies gets these claims only after its author adds `machineries=[faults]`. The unclaimed classifications of the suite name the param until then.

## Open questions

None.

## References

| What | Where |
|---|---|
| pingcap/failpoint | <https://github.com/pingcap/failpoint> |
| FreeBSD fail(9) | <https://man.freebsd.org/cgi/man.cgi?query=fail&sektion=9> |
| Toxiproxy | <https://github.com/Shopify/toxiproxy> |
