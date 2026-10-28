---
rfc: 0020
title: Generated sequences and the judges of their histories
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

# RFC-0020: Generated sequences and the judges of their histories

## Summary

A law about a history is a statement about a sequence of calls: a counter never decreases, every call of an accumulating writer has an effect, a move conserves a quantity. The suite generator writes, for each declaration, the actions of a machine of assert-go. Each action calls one callable of the declaration with inputs from the generators of the suite and records the call in the history of the case, the `history.History` of assert-go. A judge is a probe that checks that history after each step. Five judges read the history of the subject alone, and the mixins of the shape catalog license their claims. A sixth judge compares the subject with a reference implementation, and a declaration gets it when it states the machinery `reference`. The property engine generates, runs and shrinks the sequences. Falsification runs each check on its planted defect with a seed from the key of the check and without shrinking, so a proof gives the same result on every run.

## Motivation

A check of a single caller tests one call, or a fixed short sequence, against a law. A defect that appears only after a particular order of calls passes every check of a single caller. A counter that decreases after its third increment passes a check that increments twice.

The shape catalog states laws of this kind as mixins: `monotonic`, `monotonic-reads`, `accumulates`, `conservative`, `no-duplicates`, `stable-order` and `atomic`. Their text in the catalog describes a check that calls the subject repeatedly and compares each result with the earlier ones (`eidos-plugin-shape/spec/mixins/monotonic.yaml`, `accumulates.yaml`). Such a law does not need a second implementation. It needs a history of calls and a judge over that history.

A law that the catalog does not state, such as that a store returns what a map would return, needs an expected value. An expected value that the check takes from the subject cannot fail, because a defect moves the expectation with the measurement. A reference written as a separate implementation gives expected values that the defect of the subject does not move.

The parts exist:

- assert-go has the package `stateful`, with `Machine[S]`, `Action[S]`, its invariant and its settle step, and `Steps`, which runs a generated sequence inside a case of the property engine (`assert-go/stateful/machine.go`, `steps.go`).
- assert-spec defines machines, and both assert-go and assert-typescript implement them.
- conform has the door `Reference` and property checks that run on a `prop.Case` (`conform/check.go`, `conform/suite.go:134`).
- The suite generator writes a generator of the inputs of each callable for the fuzz targets and the property runs (`conform.Call`, `generator/suite/call.go`).

The model tools of other languages take the same shape. quickcheck-state-machine keeps the model in its own record beside the real system. proptest-state-machine splits a reference state machine from the test of the real system. fast-check runs commands against a model and a real object, and shrinks only the commands that ran.

## Detailed design

### The judges

| Judge | Claim | Licence | Law | Reason | Defect form | Entails |
|---|---|---|---|---|---|---|
| `ordered` | `never-goes-back` | The mixins `monotonic` and `monotonic-reads` | Each result of the callable, projected through the param `version` or taken whole, never compares below the result of the call before it | `regressed` | `goes-back-after-three` | `monotonic` or `monotonic-reads` |
| `ordered` | `keeps-its-order` | The mixin `stable-order` | Two drains without a write between them yield their elements in the same order, after any sequence of writes | `reordered-drain` | `reorders-after-three` | `stable-order` |
| `invariant` | `yields-each-once` | The mixin `no-duplicates` | Each drain after any sequence of writes yields each element once | `repeated-element` | `repeats-after-three` | `no-duplicates` |
| `folds` | `compounds-over-sequences` | The mixin `accumulates` with its param `observe` | The reader that `observe` names shows one effect for each call of any sequence | `uncompounded` | `drops-a-repeat-after-three` | `accumulates` |
| `conserved` | `conserves-the-quantity` | The mixin `conservative` with its param `field` | The sum of the field over the subject is the same after each call | `not-conserved` | `leaks-one-unit` | — |
| `after-fault` | `no-trace-of-a-failed-call` | The mixin `atomic`, with the machinery `faults` | A sequence with a failed call leaves the subject in the state of the same sequence without it | `partial-write` | `keeps-half-a-call` | `no-partial-write` |
| `agrees` | `agrees-with-the-reference` | Every declaration that states the machinery `reference` | Each call of the sequence returns on the subject what it returns on the reference | `diverges-from-the-reference` | `diverges-after-three` | — |

RFC-0003 binds `monotonic`, `monotonic-reads`, `stable-order`, `no-duplicates` and `accumulates` as claims of a single caller, and RFC-0019 binds `no-partial-write`. Each claim of this table is the deeper claim of the same law, so it entails the shallower claim and reports the same reason code. The proof of each depth pair checks that the deeper claim fails on the defect of the shallower one, and that the shallower claim passes on the defect of the deeper one, which needs at least four calls.

Each claim lists `generated-sequences` first among its machineries. `agrees-with-the-reference` has the scope `declaration`, and its identity is `Store/agrees-with-the-reference`. The other claims bind to the callable that the mixin classifies, and their partners come from the modes of RFC-0002, such as `named` for the param `observe`.

### The runtime

The calls of a sequence go into the history of the case, `prop.Case.History()`, which is a `history.History` of assert-go. Each action records its call through `Invoke` and completes it through the `Call` that `Invoke` returns: `OK` with the results, `Fail` with the error result, and `Unknown` when the context of the call ends first. A machine of `stateful` records its steps in the same history, and its `Spec` checks that history after every step. The judges and the checkers of linearizability and isolation therefore read one record.

The package `probe` of conform gains the driver and the judges:

```go
// Actions returns the actions of a machine over subject. Each action calls
// one callable and records the call in the history of the case.
type Actions[S any] func(subject S) []stateful.Action[struct{}]

// Judge checks the history of the case after each step of a sequence. It
// sets the reason code of each assertion on h.
type Judge[S any] func(h *conform.Handle[S], calls *history.History)

// Sequences runs generated sequences of actions on a fresh subject of h and
// applies judge to the history of c after each step. It engages h once
// judge has compared two calls of the callable of the claim.
func Sequences[S any](h *conform.Handle[S], c *prop.Case, actions Actions[S], judge Judge[S])

// Ordered returns a judge that checks that project of each completed call
// of operation is not less than project of the call of operation before it.
func Ordered[S, V any](operation string, project func(history.Event) V, less func(a, b V) bool, reason spec.ReasonID) Judge[S]

// AgreesWithReference runs generated sequences of actions on a fresh
// subject of h, with the Spec of the machine set to the spec that
// history.SpecFrom builds over the reference of h. The machine checks the
// history against that spec after every step, and the probe reports a
// failure with the reason diverges-from-the-reference.
func AgreesWithReference[S any](h *conform.Handle[S], c *prop.Case, actions Actions[S], reference func(r S) history.Subject)
```

`Invariant`, `Folds`, `Conserved` and `AfterFault` have the form of `Ordered`. Every assertion is inside a probe, and the generated code only calls the subject and records the calls.

### The generated actions and checks

For a store whose `Version` has the mixin `monotonic-reads` with `version=Revision`, the generator writes the actions once per declaration:

```go
func storeActions(subject catalog.Store) []stateful.Action[struct{}] {
	return []stateful.Action[struct{}]{
		{
			Name:  "Put",
			Input: func(c *prop.Case, _ struct{}) any { return c.Draw(storePutInput, "Put") },
			Run: func(c *prop.Case, client int, in any) {
				put := in.(storePutArgs)
				call := c.History().Invoke(client, "Put", []any{put.Key, put.Item}, put.Key)
				probe.Complete(call, nil, subject.Put(c.Context(), put.Key, put.Item))
			},
		},
		{
			Name:  "Version",
			Input: func(c *prop.Case, _ struct{}) any { return c.Draw(storeVersionKey, "Version") },
			Run: func(c *prop.Case, client int, in any) {
				key := in.(string)
				call := c.History().Invoke(client, "Version", []any{key}, key)
				v, err := subject.Version(c.Context(), key)
				probe.Complete(call, v, err)
			},
		},
	}
}
```

`probe.Complete` completes a call as `OK` with the output when the error is nil, as `Unknown` when the error wraps the error of the context of the case, and as `Fail` otherwise. The key of each call is the input that the binding `key` of the shape names, so a check of linearizability partitions the history by key.

and one check per claim:

```go
{
	ID:        StoreVersionCallable.Check("never-goes-back"),
	Statement: "Version never returns a Revision below one that it returned before",
	Claim:     spec.NeverGoesBack,
	Reason:    spec.Regressed,
	Property: func(h *conform.Handle[catalog.Store], c *prop.Case) {
		probe.Sequences(h, c, storeActions, probe.Ordered[catalog.Store](
			"Version", storeVersionRevision, storeRevisionLess, spec.Regressed))
	},
	Defect: func(h *conform.Handle[catalog.Store]) catalog.Store {
		return &storeGoesBackAfterThree{Store: h.Fresh()}
	},
},
```

A declaration with `//dokimi:suite machineries=reference` also gets:

```go
{
	ID:        conform.ID[catalog.Store]("Store/agrees-with-the-reference"),
	Statement: "every sequence of calls returns on Store what it returns on the reference",
	Claim:     spec.AgreesWithTheReference,
	Reason:    spec.DivergesFromTheReference,
	Property: func(h *conform.Handle[catalog.Store], c *prop.Case) {
		probe.AgreesWithReference(h, c, storeActions, storeSubject)
	},
	Defect: func(h *conform.Handle[catalog.Store]) catalog.Store {
		return &storeDivergesAfterThree{Store: h.Fresh()}
	},
},
```

The generator writes `storeSubject`, which applies an operation of the history to the reference by the name of its callable, as `history.Subject` states. The spec compares the output that the history recorded with the output of the reference through the equality of assert-go. The param `normalize=<function>` of the directive `suite` names a function that maps the result of a callable to the value that the history records, for results whose equality is not their value, such as an identifier that the server assigns. The actions and the reference both pass their results through it.

An action whose callable can block until another call, the mixin `waits`, is left out of the actions. A declaration with such a callable gets the claims through the callables that do not wait.

### The planted defects

Each defect form wraps a fresh subject and changes one callable after its third call, so `min-steps` of each form is 4. `diverges-after-three` changes the result of the first callable with a value result. `goes-back-after-three` returns the first result that it returned. The other checks of the suite make fewer than four calls of the overridden callable, so the isolation of the other checks on the defect shows that no check of a single caller detects it.

### Proofs, shrinking and replay

Falsification runs each check with a seed from the key of the check, over sequences of at least `min-steps` calls, and does not shrink, as RFC-0002 and RFC-0006 state. The corpus pins that each defect form fails its check at that seed.

On a subject the property engine shrinks a failing sequence by removing actions and shrinking inputs, and reports the smallest sequence that fails the judge. The seed of the run, `DOKIMI_SEED`, replays it. The suite record of the check contains the seed.

### Without a reference

A declaration without `machineries=reference` lists `agrees-with-the-reference` as unclaimed with the stance `needs`. A declaration that states the machinery gets the check, and an application without the door `Reference` fails the check with the name of the option, as every door of conform does.

## Alternatives considered

### A model written in the spec

The spec states a model of each contract, such as a map for a store, and the generator writes the model from it.

**Why not:** a model generated from the same spec as the check agrees with the check by construction, so it is not an independent oracle. A reference written by a person, or a second implementation that already exists, is independent.

### Agreement with a reference alone

The machinery has one claim, agreement with a reference, and the laws of the mixins are the laws that the reference satisfies.

**Why not:** each law of a mixin is a statement about the history of the subject alone. A judge of that law does not need a reference, so a declaration without one still gets the check. A reference that shares a defect with the subject also hides the defect from agreement, and a judge of the law still finds it.

### One generator of sequences per claim

Each claim has a probe that generates its own sequence of calls.

**Why not:** every claim of a declaration calls the same callables with inputs from the same generators. One set of actions per declaration serves every judge, and the property engine shrinks every failure in the same way.

### Proofs on random sequences

Falsification runs the property with the seed of the run and shrinks, as an application to a subject does.

**Why not:** a random sequence triggers a defect after its third call only by chance, so a proof could pass on one run and fail on the next. A gate that flakes is ignored. A fixed seed and a minimum length make the proof the same on every run.

## Drawbacks

- **A slower suite.** A property check runs about 100 sequences of up to 50 calls each by default.
- **Judges per runtime.** Each judge exists in the probes of each runtime.
- **The equivalence of results.** Results that differ legitimately need a function `equal`, and a wrong function hides a divergence.
- **A reference that shares a defect.** A reference with the same defect as the subject hides it from `agrees-with-the-reference`. The planted defect proves the check against a divergence, not the reference against the law.
- **A fixed proof.** The proof at the seed of the key shows that the check detects the defect form. It does not show that random sequences trigger every defect of the same form.

## Open questions

None.

## References

| What | Where |
|---|---|
| quickcheck-state-machine | <https://github.com/stevana/quickcheck-state-machine> |
| proptest-state-machine | <https://docs.rs/proptest-state-machine> |
| fast-check, model-based testing | <https://fast-check.dev/docs/advanced/model-based-testing/> |
| assert-go, the package stateful | <https://github.com/dokimasia/assert-go/tree/main/stateful> |
| The mixins of the shape catalog | <https://github.com/dokimasia/eidos/tree/main/eidos-plugin-shape/spec/mixins> |
| Oracles that cannot fail: anchoring, arXiv 2608.17214 | <https://arxiv.org/abs/2608.17214> |
