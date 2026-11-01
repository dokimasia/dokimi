---
rfc: 0020
title: Generated sequences and the judges of their histories
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-28
updated: 2026-11-01
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

A law about a history is a statement about a sequence of calls: a counter never decreases, two drains without a write between them agree on their order, a read right after a write returns that write. The suite generator writes, for each declaration, the actions of a machine of assert-go. Each action calls one callable of the declaration with generated inputs, and records the call and its result in the history of the case. The probe of each claim judges that history after each step. The mixins and the contracts of the shape catalog license 18 such claims. A declaration that states the machinery `reference` also gets a check that every sequence returns on the subject what it returns on a reference implementation. The property engine generates, runs and shrinks the sequences. Falsification runs each check on its planted defect at a fixed seed and without shrinking, so a proof gives the same result on every run.

## Motivation

A check of a single caller tests one call, or a fixed short sequence, against a law. A defect that appears only after a particular order of calls passes every check of a single caller. A counter that decreases after its third increment passes a check that increments twice.

The shape catalog states laws of this kind as mixins, such as `monotonic`, `read-your-writes`, `accumulates` and `no-duplicates`, and as contracts, such as `appender`, `cache` and `workflow`. The catalog describes each such check as repeated calls of the subject whose results the check compares with the earlier ones. Such a law does not need a second implementation. It needs a history of calls and a judge over that history.

A law that the catalog does not state, such as that a store returns what a map would return, needs an expected value. An expected value that the check takes from the subject cannot fail, because a defect moves the expectation with the measurement. A reference written as a separate implementation gives expected values that the defect of the subject does not move.

The parts exist:

- assert-go has the package `stateful`, with `Machine`, `Action` and `Steps`, which runs a generated sequence inside a case of the property engine. Its package `history` has `History`, `Call` and `Event`, and `prop.Case.History` returns the history of a case.
- conform has the door `Reference` and property checks that run on a `prop.Case`.
- The suite generator writes a generator of the calls of each callable for the fuzz targets and the property checks.

The model tools of other languages take the same form. quickcheck-state-machine keeps the model in its own record beside the real system. proptest-state-machine splits a reference state machine from the test of the real system. fast-check runs commands against a model and a real object, and shrinks only the commands that ran.

This proposal has no judge of the delivery of the messages of an outbox, a publisher or a watcher, of an iteration that sees the mutations made while it runs, of the conservation of a quantity, or of the absence of any trace of a failed call. Their classifications keep their stances.

## Detailed design

### The claims

| Claim | Licence | Law | Reason | Defect form | Entails |
|---|---|---|---|---|---|
| `never-goes-back` | The mixin `monotonic` | `results-never-decrease` | `regressed` | `goes-back-after-three` | `monotonic` |
| `reads-never-go-back` | The mixin `monotonic-reads` with its param `version` | `reads-never-go-backwards` | `regressed` | `reads-the-past-after-three` | `monotonic-reads` |
| `keeps-its-order` | The mixin `stable-order` | `two-drains-agree-on-order` | `reordered-drain` | `reorders-after-three` | `stable-order` |
| `writes-apply-in-order` | The mixin `monotonic-writes` with its param `version` | `writes-apply-in-issue-order` | `unobserved-write` | `reorders-writes-after-three` | `monotonic-writes` |
| `reads-follow-own-writes` | The mixin `read-your-writes` with its param `version` | `read-observes-its-own-writes` | `unobserved-write` | `reads-a-stale-version-after-three` | `read-your-writes` |
| `writes-follow-their-reads` | The mixin `writes-follow-reads` with its param `version` | `write-follows-the-read-before-it` | `unordered-write` | `stamps-the-past-after-three` | `writes-follow-reads` |
| `appends-keep-ascending` | The contract `appender`, role `fn` | `append-returns-above-the-last` | `not-ascending` | `descends-after-three` | `appender` |
| `replays-in-order-over-sequences` | The contract `chain`, role `append`, with the mixin `causal` | `history-replays-in-its-stamped-order` | `unordered-history` | `reverses-after-three` | `chain-causal` |
| `yields-each-once` | The mixin `no-duplicates` | `no-element-is-drained-twice` | `repeated-element` | `repeats-after-three` | `no-duplicates` |
| `drains-keep-what-went-in` | The mixin `over-match` | `drain-keeps-what-went-in` | `shrunk-drain` | `drops-after-three` | `over-match` |
| `drains-gain-exactly-what-went-in` | The mixin `permutation` | `drain-gains-exactly-what-went-in` | `miscounted-drain` | `drops-after-three` | `permutation` |
| `compounds-over-sequences` | The mixin `accumulates` with its param `observe` | `repeated-calls-compound` | `uncompounded` | `drops-a-repeat-after-three` | `accumulates` |
| `reads-each-last-write` | The mixin `read-after-write` with its param `write` | `read-returns-the-write-before-it` | `unobserved-write` | `reads-a-stale-value-after-three` | `read-after-write` |
| `persists-each-last-write` | The contract `persister`, role `writer` | `read-returns-the-write-before-it` | `unobserved-write` | `drops-a-write-after-three` | `persister` |
| `cache-agrees-over-sequences` | The contract `cache`, role `cache` | `cache-agrees-with-its-backing` | `cache-disagrees` | `serves-a-stale-value-after-three` | `cache` |
| `agrees-with-the-reference` | Every declaration that states the machinery `reference` | `every-sequence-agrees-with-the-reference` | `diverges-from-the-reference` | `diverges-after-three` | — |
| `keeps-within-its-bounds` | The mixin `bounded` with its params `limit` and `observe` | `use-is-within-its-bounds` | `exceeded-bound` | `grows-past-the-limit-after-three` | — |
| `merges-in-any-order` | The mixin `crdt-merge` with its params `write` and `read` | `merge-order-does-not-matter` | `order-dependent-merge` | `ignores-the-other-replica` | — |
| `follows-its-graph` | The contract `workflow`, role `fn`, with its params `observe` and `transitions` | `states-follow-the-graph` | `forbidden-transition` | `admits-any-transition-after-three` | — |

A claim with an entry in the column Entails witnesses in full the law of a claim of a single caller. It entails that claim and reports the same reason code. Each such law has the depth `generated-sequences`. The derivation marks the claim of a single caller as shallow and the claim of sequences as deep. The law `repeated-calls-compound` has the depth `nothing-beyond-a-caller` today and takes the depth `generated-sequences`. The claim `accumulates` then becomes shallow like the others. The proof of each pair checks that the deeper claim fails on the planted defect of the shallower one. It also checks that the shallower claim passes on the planted defect of the deeper one.

The mixins `monotonic` and `monotonic-reads` have two laws. Each of them gets a claim of its own. A claim without an entry in the column Entails has a new law with the depth `generated-sequences`. Its law is a statement about a sequence of calls. No claim of a single caller checks it.

| Law | Statement |
|---|---|
| `every-sequence-agrees-with-the-reference` | Each call of any sequence returns on the subject what it returns on the reference. |
| `use-is-within-its-bounds` | The use that `observe` reads is between `min` and `limit` after each call of any sequence. |
| `merge-order-does-not-matter` | Two replicas that any sequences of writes made diverge read the same state after a merge in either order. |
| `states-follow-the-graph` | Each state that `observe` reads after a call is the state before the call or one of its successors in `transitions`. |

| Reason | Meaning |
|---|---|
| `diverges-from-the-reference` | A call returned on the subject something other than what it returned on the reference. |
| `exceeded-bound` | The use that `observe` read was below `min` or above `limit`. |
| `order-dependent-merge` | Two merges in opposite orders read two states. |
| `forbidden-transition` | A call moved the state to a state that `transitions` does not list as a successor. |

Each claim lists `generated-sequences` first among its machineries. `agrees-with-the-reference` lists `reference-implementation` after it. Each claim with an entry in the column Entails takes the gates of the claim that it entails, and those of its partners that its probe calls. The other claims take their partners with the mode `named`:

- `keeps-within-its-bounds` takes `observe`.
- `merges-in-any-order` takes `write` and `read`.
- `follows-its-graph` takes `observe`.

A declaration that does not state a param of a partner gets the refusal `partner-found`. This applies to `observe` of `bounded` and of `workflow`, and to `write` and `read` of `crdt-merge`. A declaration that does not state a param that a probe reads gets the refusal `classification-params`, as for `limit` of `bounded` and `transitions` of `workflow`. The param `min` of `bounded` is optional. A binding of a param gains a default, and `min` has the default 0.

A probe calls the partner of an input of the kind `operation` by its name, and the partner of an input of the kind `call` with its own samples. The derivation does not bind the arguments of such a partner, so a partner whose parameters select nothing of the callable does not get the refusal `partner-arguments`. `Merge(ctx, peer Replica)` and `Add(ctx, item string)` do not share an argument, and `merges-in-any-order` still binds. The partner of an input of the kind `apply` or `projection` keeps the refusal, because the probe applies it to the arguments of another call.

The claims `writes-follow-reads` and `writes-follow-their-reads` take the gate `write-returns-the-value`, because the probe `writes-advance` reads the version of the second write from its result. The gate passes where each partner of the claim returns first a value of the type that the callable returns. A write that returns no value or a value of another type gets the refusal of the gate.

Three claims need a gate of a type that the suite generator reads. The generator judges each of them, as it judges `importable-samples`:

| Gate | Claim | Passes |
|---|---|---|
| `countable-use` | `keeps-within-its-bounds` | The callable that `observe` names returns an integer, or a list, a map or a text whose length is the use. |
| `takes-a-replica` | `merges-in-any-order` | The callable takes one input that another subject can be, so a merge can pass it another replica: a pointer to the struct of the declaration, or an interface that the subject type implements. |
| `graph-of-states` | `follows-its-graph` | The variable that `transitions` names is a map from the type that `observe` returns to a list of that type. |

### The history of a case

The calls of a sequence go into the history of the case, `prop.Case.History()`. Each action records its call through `Invoke` and completes it through `conform.Complete`.

A call that returns is an observation of the subject. It completes as `OK`. The history drops each call that completes as `Fail` from its checks. `Fail` states a call that took no effect and returned nothing that a check reads. An error such as a missing key is something that the subject observed. It is an output of `OK` as well. The output records whether the call failed. It drops the values beside an error, because a value beside an error has no meaning. A call whose error wraps the error of its ended context completes as `Unknown`. The context ended before the call returned an outcome. `Complete` takes the context of the call, so an action passes `c.Context()` and its spec passes a context that it cancels. A partner or a reference that a probe applies returns the output that `OutputOf` builds by the same rule.

```go
// Output is the result of a call that the history of a case records. A
// failed call has no values in its Output, because a value beside an error
// has no meaning.
type Output struct {
	// Values are the values that the callable returned in order. They are nil
	// for a call that failed.
	Values []any

	// Failed reports a call that returned an error.
	Failed bool
}

// OutputOf returns the Output of a call that returned values and err.
func OutputOf(err error, values ...any) Output

// Complete completes call in the history of a case. It records a call whose
// err wraps the error of the ended context ctx as Unknown. It records every
// other call as OK with the Output that [OutputOf] returns for err and
// values. An error that a call returns is an observation of the subject.
func Complete(ctx context.Context, call history.Call, err error, values ...any)
```

A probe that compares two outputs accepts two failed calls, and two calls that succeeded with values that are equal under the equality of assert-go.

### The probes

The package `probe` of conform gains a probe for each claim. Each probe runs generated sequences of the actions on the subject of the handle through `stateful.Steps`. The machine has no spec. The probe judges the history after each step in the invariant of the machine. It reports each failure through the handle with the reason code of its claim. It engages the handle once it has compared two calls. Every assertion is inside a probe. The generated code only calls the subject and records the calls.

```go
// Actions returns the actions of a machine over the subject that subject
// returns. Each action calls one callable of a declaration with the inputs
// that it requests from the case, and completes the call in the history of
// the case through [conform.Complete]. The actions take their keys from pools
// that Actions draws from c.
type Actions[S any] func(c *prop.Case, subject func() S) []stateful.Action[struct{}]

// Apply calls a partner with the arguments of a call that the history of a
// case recorded, and returns the output of the partner. The generated check
// binds the arguments of the call to the inputs of the partner.
type Apply func(ctx context.Context, args []any) conform.Output

// Reference returns the function that applies an operation of the history of
// a case to the subject r by the name of its callable. The function returns
// the output of the call.
type Reference[S any] func(r S) func(ctx context.Context, operation string, args []any) conform.Output

// ResultsNeverFall runs generated sequences of the actions. It fails when a
// call of read returns a result below the result of the call of read with the
// same arguments before it.
func ResultsNeverFall[S any, O cmp.Ordered](h *conform.Handle[S], c *prop.Case, actions Actions[S], read string)

// SequencesAgreeWithTheReference runs generated sequences of the actions, and
// applies each call to a reference from the door Reference of h through
// reference. It fails when the output of the reference differs from the
// output of the call.
func SequencesAgreeWithTheReference[S any](
	h *conform.Handle[S], c *prop.Case, actions Actions[S], reference Reference[S],
)

// MergesCommute writes to two replicas, the subject of h and a new subject,
// through generated sequences of the action of write. It merges the second
// replica into the first, applies the same writes to two more new replicas,
// and merges the first of those into the second. It fails when the action of
// read returns two outputs on the two merged replicas. apply calls merge on a
// replica with the other replica as its argument, and the guard of the probe
// fails a merge that fails.
func MergesCommute[S any](
	h *conform.Handle[S], c *prop.Case, actions Actions[S], apply Reference[S], merge, write, read string,
)
```

The actions take a function that returns the subject, so a probe can run one action on more than one subject. A probe reads the first value of an output. The probes of one call, `results-never-fall`, `drains-agree-between-writes` and `each-call-moves-the-observation`, compare the calls that have the same arguments. The other probes that compare two calls compare the calls that have the same keys, and a call without keys with the calls that have the same arguments. `positions-rise` and the probes of drains take a call without keys as a call of every key, so a drain without keys gains the writes of every key.

| Probe | Claims | Inputs after the handle, the case and the actions | Judges after each step |
|---|---|---|---|
| `results-never-fall` | `never-goes-back` | `read`: operation | The result of each call of `read` is not below the result of the call of `read` with the same arguments before it |
| `read-versions-never-fall` | `reads-never-go-back` | `read`: operation, `version`: projection | The version of each value that `read` returns is not below the one before it |
| `drains-agree-between-writes` | `keeps-its-order` | `drain`: operation | Two calls of `drain` with the same arguments and no call between them return one list |
| `written-versions-rise` | `writes-apply-in-order` | `write`: operation, `read`: apply, `version`: projection | After each call of `write`, `read` returns a version above the version that it returned after the write of its keys before it |
| `reads-see-the-last-write` | `reads-follow-own-writes` | `read`: operation, `write`: operation, `readback`: apply, `version`: projection | Each call of `read` returns a version not below the highest version of a write of its key |
| `writes-rise-above-reads` | `writes-follow-their-reads` | `write`: operation, `read`: operation, `readback`: apply, `version`: projection | The version of each call of `write` is above the version that the last call of `read` of its key returned |
| `positions-rise` | `appends-keep-ascending` | `append`: operation | The position that each call of `append` returns is above the one before it |
| `replays-follow-the-stamps` | `replays-in-order-over-sequences` | `replay`: operation, `stamp`: projection | Each call of `replay` returns its entries in the order of their stamps |
| `drains-without-repeats` | `yields-each-once` | `drain`: operation | No call of `drain` returns an element twice |
| `drains-contain-the-writes` | `drains-keep-what-went-in` | `drain`: operation, `write`: operation, `written`: projection | Each call of `drain` gains over the drain before it each element written since |
| `drains-gain-the-writes` | `drains-gain-exactly-what-went-in` | `drain`: operation, `write`: operation, `written`: projection | Each call of `drain` gains over the drain before it exactly the elements written since, each as often as a write wrote it |
| `each-call-moves-the-observation` | `compounds-over-sequences` | `call`: operation, `observe`: apply | `observe` reads another value for the arguments of each call of `call` than it read before the call |
| `reads-return-the-last-write` | `reads-each-last-write`, `persists-each-last-write` | `write`: operation, `readback`: apply, `written`: projection | After each call of `write`, `readback` returns the value that the write wrote |
| `reads-agree-with-the-backing` | `cache-agrees-over-sequences` | `read`: operation, `backing`: apply | The output of each call of `read` agrees with the output of `backing` for its arguments |
| `sequences-agree-with-the-reference` | `agrees-with-the-reference` | `reference`: reference | The output of each call agrees with the output of the same call on the reference |
| `use-is-within-the-bounds` | `keeps-within-its-bounds` | `observe`: call, `min`: int, `limit`: int | The use that `observe` reads is between `min` and `limit` |
| `merges-commute` | `merges-in-any-order` | `apply`: reference, `merge`: operation, `write`: operation, `read`: operation | At the end of a sequence of writes, the two merges read one state |
| `states-move-along-the-graph` | `follows-its-graph` | `fn`: operation, `observe`: apply, `transitions`: value | The state that `observe` reads after each call of `fn` is the state before it or one of its successors |

`reads-return-the-last-write` serves two claims. The probe `write-then-read` serves four claims of a single caller in the same way. `each-call-moves-the-observation` and `states-move-along-the-graph` wrap the action of their callable, which applies `observe` to the arguments of each call right before and right after the call. A change that another operation makes then does not count as the effect of the call, and the first call of an instance has a state before it. The wrapper reads the arguments from the input of the action.

`reads-see-the-last-write` and `writes-rise-above-reads` take the version of a write from the value that the write returns, which the claims of a single caller of the same laws read as well. A write that returns no such value has the version that the readback returns for its arguments. The readback calls the callable of the claim, so a planted defect that returns a stale value also changes the readback. `reads-see-the-last-write` compares each read with the highest version of a write of its key. `written-versions-rise` reads the version back after each write, because the mixin `monotonic-writes` states what a read observes.

A drain gains the elements that it returns beyond the elements of the drain before it, each element of the earlier drain cancelling one equal element. A drain that removes what it returns and a drain that returns a snapshot gain the same elements, unless an element written since equals an element of the drain before. The probes of drains judge no such drain, because they cannot tell the two kinds of drain apart.

`merges-commute` builds two replicas: the subject of the handle and a new subject. It runs the action of `write` on them, with an action of its own for each replica, because two actions of one machine cannot have one name. At the end of the sequence the probe merges the second replica into the first. It builds two new replicas and runs the action of `write` on them with the inputs of the sequence. It merges the first new replica into the second. It then runs the action of `read` on both merged replicas with one input, and compares the two outputs.

The spec gains the input kinds `actions`, `operation`, `apply` and `reference`, beside the kinds that it has. An input of the kind `actions` or `reference` takes the source `declaration`. Its value is a function that the generator writes once for the declaration.

### The generated actions and checks

The generator writes one function of actions for each declaration. It generates a pool of one to three keys of each key type for the case, and then returns one action for each callable:

```go
func storeActions(c *prop.Case, subject func() catalog.Store) []stateful.Action[struct{}] {
	keysString := c.Draw(prop.List(prop.Of[string](), prop.Unique(), prop.MinSize(1), prop.MaxSize(3)), "keys.string")
	return []stateful.Action[struct{}]{
		{
			Name: "Put",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "put.key")
				item := c.Draw(prop.Of[catalog.Item](), "put.item")
				return []any{key, item}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Put", args, args[0])
				values, err := storePutCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
		{
			Name: "Version",
			Input: func(c *prop.Case, _ struct{}) any {
				key := c.Draw(prop.SampledFrom(keysString...), "version.key")
				return []any{key}
			},
			Run: func(c *prop.Case, client int, input any) {
				args, _ := input.([]any)
				call := c.History().Invoke(client, "Version", args, args[0])
				values, err := storeVersionCall(c.Context(), subject(), args)
				conform.Complete(c.Context(), call, err, values...)
			},
		},
	}
}

// storeVersionCall calls Version on subject with the arguments args, and
// returns the values of the call beside its error.
func storeVersionCall(ctx context.Context, subject catalog.Store, args []any) ([]any, error) {
	a0, _ := args[0].(string)
	r0, err := subject.Version(ctx, a0)
	return []any{r0}, err
}
```

The input of an action is the list of the arguments of its call, which the history records. For a callable, the generator writes a function that calls it with the arguments of such a list, such as `storeVersionCall`. The actions, the function of the reference and the `Apply` of a partner call the callable through that function, so they share one call site of the callable. A call that returns a stream returns the list of its elements. The action generates each input with the rules of the call generator of the callable, so a sample, a generated function and a stream of an input are the same as in a fuzz target. Two inputs take their values from the pool of their type:

- an input that the binding `key` of the shape of its callable names
- a field of a written value that a pair of a claim of the declaration binds to such an input

In most steps a read then takes a key that a write stored, and the probes of a key compare calls. The history records that input as the key of each call. A check of linearizability then partitions the history by key. Each label has the form of a key, such as `put.key`.

A callable with the mixin `waits` or `runs-until-cancelled` has no action, because a call of it does not return on its own. A declaration with such a callable gets the claims through its other callables.

The check of each claim is a property check:

```go
{
	ID:          StoreChecks.Version.NeverGoesBack(),
	Statement:   "Version never returns a result below one that it returned before",
	Claim:       spec.NeverGoesBack,
	Probe:       spec.ResultsNeverFall,
	Reason:      spec.Regressed,
	Law:         spec.ResultsNeverDecrease,
	Machineries: []spec.MachineryID{spec.GeneratedSequences},
	Entails:     []conform.ID[catalog.Store]{StoreChecks.Version.Monotonic()},
	Property: func(h *StoreHandle, c *prop.Case) {
		probe.ResultsNeverFall[catalog.Store, int](h, c, storeActions, "Version")
	},
	Defect: storeVersionGoesBackAfterThree,
},
```

The param `normalize=<function>` of the directive `suite` names a function `func(operation string, values []any) []any`. It maps the values of a callable to the values that the history records, for values whose equality is not their value, such as an identifier that a server assigns. The function of the call of each callable passes its values through it, so the actions, the function of the reference and each `Apply` of a partner see the same values. A function of another signature gets the error `BadNormalize` of the generator at the directive, and the suite then passes the values as they are.

### The check of a whole declaration

`agrees-with-the-reference` has the scope `declaration`, and its identity is `Store/agrees-with-the-reference`. The spec gains a licence of every declaration for it. The generator writes `storeReference`, which applies an operation of the history to a reference by the name of its callable and returns its output through `OutputOf`. The probe builds the reference of each case from the door `Reference` of the handle.

A declaration without `machineries=[reference]` gets neither the check nor an entry in its suite record, because the unclaimed table lists classifications of callables. A declaration that states the machinery gets the check, and an application without the door `Reference` fails the check with the name of the option, as every door of conform does. The fuzz targets of a suite run the checks of one callable, so they do not run this check.

The licence `declaration` binds the claim on one callable of the declaration: the first callable that passes the gates of the claim, or the first callable where none does. The claim has the gate `returns-a-value`, which the derivation judges. A callable passes it when it returns a value and has neither the mixin `waits` nor the mixin `runs-until-cancelled`, so the actions call it. The planted defect changes the value of that callable. A declaration without such a callable gets the refusal of the gate, which names the declaration.

### The planted defects

| Form | Changes | Does |
|---|---|---|
| `goes-back-after-three` | callable | Returns the first result that it returned, from its fourth call on |
| `reads-the-past-after-three` | callable | Returns the first value that it returned, from its fourth call on |
| `reorders-after-three` | callable | Reverses the elements that it returns, from its fourth call on |
| `reorders-writes-after-three` | partner:write | Stores nothing, from its fourth call on |
| `reads-a-stale-version-after-three` | callable | Returns the first value that it returned, from its fourth call on |
| `stamps-the-past-after-three` | partner:write | Returns each value with the version 0, from its fourth call on |
| `descends-after-three` | callable | Returns the position 0, from its fourth call on |
| `reverses-after-three` | partner:replay | Reverses the entries that it replays, from its fourth call on |
| `repeats-after-three` | callable | Returns its first element twice, from its fourth call on |
| `drops-after-three` | callable | Drops the last element that it returns, from its fourth call on |
| `drops-a-repeat-after-three` | callable | Returns no error and does nothing, from its fourth call on |
| `reads-a-stale-value-after-three` | callable | Returns the first value that it returned, from its fourth call on |
| `drops-a-write-after-three` | callable | Returns no error and does nothing, from its fourth call on |
| `serves-a-stale-value-after-three` | callable | Returns the first value that it returned, from its fourth call on |
| `diverges-after-three` | declaration | Returns the second sample of the value of the callable on which the check binds, from its fourth call on |
| `grows-past-the-limit-after-three` | partner:observe | Returns a use one above `limit`, from its fourth call on |
| `ignores-the-other-replica` | callable | Returns no error and does nothing, on every call |
| `admits-any-transition-after-three` | callable | Makes its transition twice, from its fourth call on |

Each form wraps a fresh subject. A form that changes its callable from its fourth call on has the `min-steps` 4. A check of a single caller of the same law makes fewer than four calls of the callable and passes on the defect. Because a case merges twice, `ignores-the-other-replica` changes every call. Its `min-steps` is 1. No claim of a single caller checks the law of a merge. The value `declaration` of the field `changes` names the callable on which the check of the declaration binds.

### Proofs, shrinking and replay

Falsification runs the cases of each check at a seed of the identity of the check. As for every property check, it runs without a store and without shrinking. `stateful` runs a mean of 30 and at most 100 steps in a case. A case keeps each action with odds of 1 in 2. The cases of a falsification include sequences with four calls of each callable. The corpus pins that each defect form fails its check at its seed.

On a subject the property engine shrinks a failing sequence. It removes steps, shrinks inputs and reports the smallest sequence that fails the probe. The suite record of the check contains the seed. `DOKIMI_SEED` replays it.

### Isolation

Isolation runs each other check of a suite on each planted defect. A planted defect proves its check only when no other check fails on it. A sequence check calls every callable of its declaration through its actions. A planted defect of another callable affects it by design. Isolation leaves out a sequence check for the planted defect of another callable. It leaves out the check of agreement with a reference for every planted defect but its own. That check compares every result of every callable.

A check of a single caller that calls a callable more than three times, such as the check of the burst of a rate limit, can fail on a planted defect of a sequence check of the same callable. The corpus gives such a pair of classifications fixtures of their own.

## Alternatives considered

### A model written in the spec

The spec states a model of each contract, such as a map for a store, and the generator writes the model from it.

**Why not:** a model generated from the same spec as the check agrees with the check by construction, so it is not an independent oracle. A reference written by a person, or a second implementation that already exists, is independent.

### Agreement with a reference alone

The machinery has one claim, agreement with a reference, and the laws of the mixins are the laws that the reference satisfies.

**Why not:** each law of a mixin is a statement about the history of the subject alone. A judge of that law does not need a reference, so a declaration without one still gets the check. A reference that shares a defect with the subject also hides the defect from agreement, and a judge of the law still finds it.

### The reference as the spec of the machine

The machine of a check of agreement takes the spec that `history.SpecFrom` builds over the reference, and checks the history against it after each step.

**Why not:** the machine reports a failed check of its spec to the case, under its own contract that the history is linearizable. Falsification then never sees the reason code of the claim and reports a wrong reason. A probe that applies each call to the reference reports through the handle with the reason code.

### A failed call as `Fail`

Each action completes a call that returned an error as `Fail`.

**Why not:** the history drops a failed call from every check. A store that reports a missing key would never be compared with the reference, and a defect that returns the error of another key would pass.

### Inputs that each action generates on its own

Each action generates each input from the generator of its type.

**Why not:** a read then rarely takes a key that a write stored. The judges of a key compare few calls, and a defect that a read of a written key shows passes most cases.

### One generator of sequences per claim

Each claim has a probe that generates its own sequence of calls.

**Why not:** every claim of a declaration calls the same callables with inputs from the same generators. One set of actions per declaration serves every probe, and the property engine shrinks every failure in the same way.

### Proofs on random sequences

Falsification runs the property with the seed of the run and shrinks, as an application to a subject does.

**Why not:** a random sequence triggers a defect after its third call only by chance, so a proof could pass on one run and fail on the next. A gate that flakes is ignored. A fixed seed makes the proof the same on every run.

## Drawbacks

- **A slower suite.** A property check runs 100 cases of a mean of 30 steps each by default.
- **Probes per runtime.** Each of the 18 probes exists in the probes of each runtime.
- **The equivalence of results.** Results that differ legitimately need the function of `normalize`, and a wrong function hides a divergence.
- **A reference that shares a defect.** A reference with the same defect as the subject hides it from `agrees-with-the-reference`. The planted defect proves the check against a divergence, not the reference against the law.
- **A fixed proof.** The proof at the seed of the identity shows that the check detects the defect form. It does not show that random sequences trigger every defect of the same form.
- **Isolation with exceptions.** A sequence check does not prove its isolation from the planted defects of the other callables of its declaration.
- **Five gates, four input kinds and a source.** The spec gains the gates `countable-use`, `takes-a-replica`, `graph-of-states`, `returns-a-value` and `write-returns-the-value`, the input kinds `actions`, `operation`, `apply` and `reference`, and the source `declaration`.
- **A refusal of a single caller more.** `writes-follow-reads` refuses a write that returns no value of the type of its read, where its generated suite failed before.
- **A shallow claim more.** `accumulates` becomes shallow, because its law takes the depth `generated-sequences`.

## Open questions

None.

## References

| What | Where |
|---|---|
| quickcheck-state-machine | <https://github.com/stevana/quickcheck-state-machine> |
| proptest-state-machine | <https://docs.rs/proptest-state-machine> |
| fast-check, model-based testing | <https://fast-check.dev/docs/advanced/model-based-testing/> |
| assert-go, the package stateful | <https://github.com/dokimasia/assert-go/tree/main/stateful> |
| assert-go, the package history | <https://github.com/dokimasia/assert-go/tree/main/history> |
| The mixins of the shape catalog | <https://github.com/dokimasia/eidos/tree/main/eidos-plugin-shape/spec/mixins> |
| Oracles that cannot fail: anchoring, arXiv 2608.17214 | <https://arxiv.org/abs/2608.17214> |
