---
rfc: 0019
title: Calls that fail on demand
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

# RFC-0019: Calls that fail on demand

## Summary

A subject often depends on another declaration: a cache on its backing store, a service on its client. Laws about how the subject handles a failure of that dependency need a call of the dependency that fails when the check asks. conform gains the door `Faults` and its option `conform.Faultable`. A declaration states the machinery with `//dokimi:suite machineries=faults`, and the consumer builds a fresh subject over its dependencies and passes each dependency the fault that an `Injector` arms. A probe arms the injector, calls the subject and checks one law, such as that a failed write leaves no partial state. On an application without the door, a check that needs it fails with the name of the option, like every door of conform today.

## Motivation

The fault strategies exist. The package `fault` of dokimi-go has counted, retry-shaped, probabilistic, windowed and predicate faults, composable with `And` and `Or`. The stub generator arms them on each method of a double through `Faults`, `FaultsOnce`, `SetFault` and the helpers of the directive `fault`. A fault armed on a method without an error result fails the test, so an injection that injects nothing cannot pass.

The doors of conform are a reference, a second subject, a restart and a crash (`conform/handle.go:50-67`), and none of them makes a call fail. The ledger of the corpus lists the door that induces a failure as pending (`corpus/go/README.md`). A cache that serves a stale value after its backing fails passes the checks of a single caller when the backing never fails, and so does a writer that applies half a batch before an error.

RFC-0002 names this machinery `a-call-that-fails-on-demand`. It arrives as the clock and the doors arrived. The subject declares the capability through an option, the check demands it, and a missing option fails the check by name.

## Detailed design

### The injector

```go
// Injector arms failures on the dependencies of one subject. A probe arms
// it, and each dependency that the consumer built over it fails while the
// injector is armed. An Injector belongs to one handle and one subject.
type Injector struct{ /* unexported */ }

// FailNext makes the next call of any armed dependency fail with err. The
// calls after it succeed.
func (i *Injector) FailNext(err error)

// FailUntilHealed makes every call of every armed dependency fail with err
// until Heal.
func (i *Injector) FailUntilHealed(err error)

// Heal clears every armed failure.
func (i *Injector) Heal()

// Failures returns the number of calls that failed because of the injector.
func (i *Injector) Failures() int

// FaultOf returns a fault strategy that fails the calls of one dependency
// while the injector is armed. The consumer passes it to the double of the
// dependency through SetFault, or calls ShouldTrigger from a hand-written
// dependency.
func FaultOf[C any](i *Injector) fault.Fault[C]
```

### The door

```go
type Doors[S any] struct {
	// ... existing doors

	// Faults returns a fresh subject whose dependencies fail as inj arms
	// them.
	Faults func(h *Handle[S], inj *Injector) S
}

// Faultable sets the door Faults of an application.
func Faultable[S any](newFaulty func(h *Handle[S], inj *Injector) S) Option[S]

// Faulty returns a fresh subject from the door Faults with its injector. It
// fails the check with the name of the option Faultable when the
// application has no such door.
func (h *Handle[S]) Faulty() (S, *Injector)
```

A consumer of a cache in front of a store:

```go
conform.Run(t, cachetest.CacheSuite, newCache,
	conform.Faultable(func(h *conform.Handle[cache.Cache], inj *conform.Injector) cache.Cache {
		backing := storetest.NewStoreStub(h, storetest.StoreStubDelegateTo(store.NewMemory()))
		backing.OnGet.SetFault(conform.FaultOf[storetest.StoreGetCall](inj))
		return cache.New(backing)
	}))
```

### The claims

The spec gains claims whose machinery is `a-call-that-fails-on-demand`:

| Claim | Licence | Law | Reason | Defect form | Entails |
|---|---|---|---|---|---|
| `no-partial-write` | The shape `writer` with the mixin `atomic` | A write that fails leaves the state of the subject as before | `partial-write` | `half-applied` | — |
| `failure-is-reported` | Every shape with a failure channel | The subject returns a failure of the dependency to the caller as an error | `swallowed-failure` | `swallows-failure` | — |
| `retry-succeeds-under-faults` | The mixin `retry-succeeds` with its param `attempts` | A call succeeds once the failure of its dependency clears, within its declared attempts | `never-converged` | `gives-up-early` | `retry-succeeds` |
| `transaction-under-faults` | The contract `transaction`, role `fn` | A scope whose dependency fails leaves no effect of its body | `not-rolled-back` | `commits-on-failure` | `transaction` |

RFC-0003 binds `retry-succeeds` and `transaction` as claims of a single caller, whose subject fails through its own inputs. The two claims of this table witness the same laws in full, with the failure of a dependency, so each entails the claim of RFC-0003 and reports its reason code.

Each claim lists `a-call-that-fails-on-demand` among its machineries. A declaration without `machineries=faults` lists the classifications of these claims as unclaimed with the stance `needs`, and a note that names the param. The derivation therefore does not need a fact about the dependencies of a subject: the declaration states that its consumer builds the subject over dependencies that fail.

## Alternatives considered

### Failpoints in the source of the subject

The subject marks the places that can fail, as `pingcap/failpoint` and FreeBSD `fail(9)` mark them, and a check enables a marker by name.

**Why not:** a failpoint is code in the subject that a build step rewrites. A generated check would demand markers that the consumer writes by hand at the right places. A failure of a dependency needs no change of the subject, because the consumer already builds the subject over its dependencies.

### Faults at the network

A proxy such as Toxiproxy fails the connections of the subject.

**Why not:** it needs a subject that talks over a network and a proxy process in the test. Most subjects of a suite are in-process types. A proxy belongs to the chaos experiments against running systems, outside the suite.

## Drawbacks

- **The consumer wires each dependency.** A subject gets the door only when the consumer builds it over doubles or over dependencies that call `FaultOf`.
- **More doors in an application.** Each door is one more option that a consumer may need to pass, and a missing door fails the checks that need it.
- **A machinery that the declaration states.** A declaration whose subject has dependencies gets these claims only after its author adds `machineries=faults`. The unclaimed classifications of the suite name the param until then.

## Open questions

None.

## References

| What | Where |
|---|---|
| pingcap/failpoint | <https://github.com/pingcap/failpoint> |
| FreeBSD fail(9) | <https://man.freebsd.org/cgi/man.cgi?query=fail&sektion=9> |
| Toxiproxy | <https://github.com/Shopify/toxiproxy> |
