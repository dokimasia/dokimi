---
rfc: 0028
title: Claims of a single caller for the classifications without one
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Draft
created: 2026-10-31
updated: 2026-10-31
discussion: none
supersedes: none
superseded-by: none
produces-adr: tbd
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0028: Claims of a single caller for the classifications without one

## Summary

Ten classifications of the shape catalog have no claim in the spec and the stance `gap`. Two of them, the mixins `sticky` and `leak-free`, need a dependency that fails on demand, and RFC-0019 binds them. The other eight state a law that one caller can check with the machinery `nothing-beyond-a-caller`: the shapes `batch-reader`, `poison-accessor` and `stream-reader`, the mixins `indexed` and `wrapped-via`, and the roles `chain/verify`, `lease/release` and `pool/stats`. This RFC gives each of them a claim with its law, its probe, its reason code and its planted defect, so that the unclaimed table of the spec has no entry with the stance `gap`.

## Motivation

The derivation lists a classification without a claim as unclaimed, and the stance `gap` states that the classification implies something that no claim witnesses. A consumer that writes `//dokimi:suite` on a declaration with such a classification does not get a check of it, and its suite record states the gap.

Each classification of the catalog states how a check observes it, in the field `observation` of its entry. Eight of the ten observations need one caller and the samples of RFC-0014:

| Classification | Observation of the catalog |
|---|---|
| `shape:batch-reader` | Stores values under several keys, reads them in one call, and compares the returned list with the stored values |
| `shape:poison-accessor` | Induces a failure of the receiver, and expects every later call to report it. A receiver that has not failed reports no error |
| `shape:stream-reader` | Stores values through a sibling writer, drains the stream, and compares the drained values with the stored ones |
| `mixin:indexed` | Reads the size through the callable of the param `by`, and calls the subject at the first and the last position and at the size |
| `mixin:wrapped-via` | Calls the subject and its sibling with one input, and compares their results |
| `contract:chain/verify` | Appends entries and replays them, then changes an entry and expects the verify to report it |
| `contract:lease/release` | Acquires, acquires again to see the taken error, releases, and acquires again to see success |
| `contract:pool/stats` | Gets and puts objects, and compares the accounting that the stats report with the sequence |

The spec states two of the eight wrongly today. The note of `mixin:wrapped-via` names a law about the wrapping of an error, while the catalog states a delegation to a sibling. The note of `mixin:sticky` names a set value that never changes, while the catalog states that one instance serves the inputs of one key.

## Detailed design

### The claims

| Claim | Licence | Partners and params | Gates | Law | Probe | Reason | Defect form |
|---|---|---|---|---|---|---|---|
| `reads-each-key` | The shape `batch-reader` | `write`: the establishing writer | `cancellation`, `list-value` | `batch-read-returns-each-stored-value` | `reads-what-was-stored` | `lost-a-key` | `drops-the-last-key` |
| `fresh-reports-no-poison` | The shape `poison-accessor` | — | `failure` | `fresh-receiver-is-healthy` | `reports-no-error` | `false-poison` | `reports-poison` |
| `streams-what-was-stored` | The shape `stream-reader` | `write`: the establishing writer | the gates of `over-match` | `drain-keeps-what-went-in` | `keeps-what-went-in` | `shrunk-drain` | `drains-nothing` |
| `respects-its-bounds` | The mixin `indexed` with its param `by` | `by`: named by `by` | `failure` | `positions-are-below-the-size` | `bounded-by-the-size` | `out-of-bounds` | `admits-the-size` |
| `agrees-with-its-sibling` | The mixin `wrapped-via` with its param `fn` | `fn`: named by `fn` | `returns-unaided` | `wrapper-returns-what-its-sibling-returns` | `agrees-with-the-sibling` | `diverged-from-the-sibling` | `stops-delegating` |
| `verifies-what-it-appended` | The contract `chain`, role `verify` | `append`: the role `append` | `failure` | `appended-chain-verifies` | `verifies-after-appends` | `false-tamper` | `reports-tamper` |
| `release-frees-the-lease` | The contract `lease`, role `release` | `acquire`: the role `acquire` | `cancellation`, `failure` | `released-lease-is-granted-again` | `grants-after-release` | `kept-after-release` | `release-keeps-the-lease` |
| `counts-each-loan` | The contract `pool`, role `stats` | `get`: the role `get` | `one-value` | `stats-move-with-a-loan` | `moves-with-a-loan` | `uncounted-loan` | `frozen-stats` |

Each claim lists `nothing-beyond-a-caller` among its machineries, and the depth of each new law is `nothing-beyond-a-caller`. `streams-what-was-stored` binds the law, the probe, the reason code and the defect form of `over-match` to the shape `stream-reader`, because the observation of the shape is the observation of that mixin.

The statement of each claim:

| Claim | Statement |
|---|---|
| `reads-each-key` | `{callable} returns the value that {write} stored under each key` |
| `fresh-reports-no-poison` | `{callable} reports no error for a fresh receiver` |
| `streams-what-was-stored` | `{callable} yields each value that {write} stored` |
| `respects-its-bounds` | `{callable} admits each position below the size that {by} reports and refuses the size` |
| `agrees-with-its-sibling` | `{callable} returns what {fn} returns for one input` |
| `verifies-what-it-appended` | `{callable} reports no error for a chain that {append} produced` |
| `release-frees-the-lease` | `{acquire} grants a lease again after {callable} released it` |
| `counts-each-loan` | `{callable} reports another count after {get} lends an object` |

### The laws, the reasons and the defect forms

| Law | Statement |
|---|---|
| `batch-read-returns-each-stored-value` | A batch read returns the stored value of each key, in the order of the keys. |
| `fresh-receiver-is-healthy` | The accessor of a fresh receiver reports no error. |
| `positions-are-below-the-size` | The callable admits each position below the size of the collection and refuses the position equal to the size. |
| `wrapper-returns-what-its-sibling-returns` | The callable returns what its sibling returns for one input. |
| `appended-chain-verifies` | The verify of a chain that appends alone produced reports no error. |
| `released-lease-is-granted-again` | An acquire after the release of a lease succeeds. |
| `stats-move-with-a-loan` | The stats of a pool change when the pool lends an object. |

| Reason | Meaning |
|---|---|
| `lost-a-key` | A batch read returned no stored value for a key. |
| `false-poison` | The accessor of a fresh receiver reported an error. |
| `out-of-bounds` | The callable refused a position below the size, or admitted the position equal to the size. |
| `diverged-from-the-sibling` | The callable returned something other than its sibling for one input. |
| `false-tamper` | The verify reported a change of a chain that appends alone produced. |
| `kept-after-release` | An acquire after a release reported the taken error. |
| `uncounted-loan` | The stats after a loan were the stats before it. |

| Form | Changes | Does | Minimum steps |
|---|---|---|---|
| `drops-the-last-key` | callable | The callable returns the list of the fresh subject without its last value. | 2 |
| `reports-poison` | callable | The callable fails with an error that no directive declares on every call. | 1 |
| `admits-the-size` | callable | The callable returns its second samples and no error for every position. | 2 |
| `stops-delegating` | callable | The callable returns its second samples and no error on every call. | 2 |
| `reports-tamper` | callable | The verify fails with an error that no directive declares on every call. | 2 |
| `release-keeps-the-lease` | callable | The release returns no error and does nothing. | 3 |
| `frozen-stats` | callable | The stats return the zero value on every call. | 2 |

### The probes

The package `probe` of conform gains seven probes, which follow the conventions of RFC-0027. `streams-what-was-stored` uses the probe `KeepsWhatWentIn` that exists.

```go
// ReadsWhatWasStored writes two values through write, reads both keys through
// call in one batch, and fails unless the batch returns the two written values
// in the order of their keys. A write that fails ends the probe.
func ReadsWhatWasStored[S, V any](h *conform.Handle[S], call Call[[]V], write func(ctx context.Context) (V, error))

// ReportsNoError calls call on the fresh subject, and fails when the call
// returns an error.
func ReportsNoError[S any](h *conform.Handle[S], call Act)

// BoundedByTheSize reads the size through size, calls call at the first and the
// last position and at the size, and fails unless the first two calls succeed
// and the call at the size fails. A size of zero calls call at the size alone.
func BoundedByTheSize[S any](h *conform.Handle[S], call func(ctx context.Context, position int) error, size Call[int])

// AgreesWithTheSibling calls call on the subject and sibling on a new subject
// of the factory, with one input, and fails when their results differ.
func AgreesWithTheSibling[S, V any](h *conform.Handle[S], call, sibling Call[V])

// VerifiesAfterAppends appends twice through append, then calls call, and
// fails when the verify returns an error. An append that fails ends the
// probe.
func VerifiesAfterAppends[S any](h *conform.Handle[S], call, append Act)

// GrantsAfterRelease acquires through acquire, releases through call, and
// acquires again, and fails when the second acquire returns an error. A first
// acquire or a release that fails ends the probe.
func GrantsAfterRelease[S any](h *conform.Handle[S], call, acquire Act)

// MovesWithALoan reads the stats through call, lends an object through get, and
// reads the stats again, and fails when the two stats are equal. A get that
// fails ends the probe.
func MovesWithALoan[S, V any](h *conform.Handle[S], call Call[V], get Act)
```

| Probe | Inputs | Reasons | Vacuous condition |
|---|---|---|---|
| `reads-what-was-stored` | `call`: call, `write`: call | `lost-a-key` | A write fails. |
| `reports-no-error` | `call`: call | `false-poison` | Never, because the probe makes its one call. |
| `bounded-by-the-size` | `call`: call, `size`: call | `out-of-bounds` | The size fails. |
| `agrees-with-the-sibling` | `call`: call, `sibling`: call | `diverged-from-the-sibling` | Never, because the probe calls both. |
| `verifies-after-appends` | `call`: call, `append`: call | `false-tamper` | An append fails. |
| `grants-after-release` | `call`: call, `acquire`: call | `kept-after-release` | The first acquire or the release fails. |
| `moves-with-a-loan` | `call`: call, `get`: call | `uncounted-loan` | The get fails. |

`AgreesWithTheSibling` calls the sibling on a new subject, so a write through the subject does not change the state that the sibling reads. The suite generator binds the call of the partner `fn` to `h.New()` for this probe, and to `h.Subject()` for every other probe. `BoundedByTheSize` takes each integer input of the callable as the position, as the catalog states for the mixin, and binds the other inputs to their first samples.

### The unclaimed table

The unclaimed table of the spec loses its ten entries with the stance `gap`. RFC-0019 removes the entries of `mixin:sticky` and `mixin:leak-free`, and this RFC removes the other eight. The test of the package `spec` that fails for a classification without a claim and without an entry then passes for each of them.

## Alternatives considered

### The stance `gap` until a deeper machinery

The eight classifications keep the stance `gap` until a machinery above a single caller checks them in full, such as generated sequences for `batch-reader`.

**Why not:** each observation of the catalog needs one caller, so a check of a single caller witnesses the law in full. A gap would leave eight classifications without a check that the runtime supplies today.

### The induced failure of `poison-accessor`

The claim of `poison-accessor` induces a failure of the receiver and expects the accessor to report it.

**Why not:** the shape has no param for the callable that induces the failure. The mixin `poisonable` has the param `induce`, and its claim `poisonable` checks the induced failure. The shape alone states what the accessor of a receiver that has not failed returns, which its edge in the catalog states.

### A change of a chain in `verifies-what-it-appended`

The claim changes an entry of the chain and expects the verify to report the change.

**Why not:** the contract `chain` has no role that changes an entry. The mixin `tamper-evident` has the params `tamper` and `verify`, and its claim checks the detection. The role `verify` alone states that an untampered chain verifies, which a check of a single caller witnesses.

## Drawbacks

- **Weak laws for two roles.** `counts-each-loan` checks that the stats change with a loan, not that they count it correctly, because the contract does not name a field of the stats. `fresh-reports-no-poison` checks a receiver that has not failed.
- **Eight more claims.** Each claim adds a probe, a reason, a law and a defect form to the spec, and a check to each suite with the classification: seven probes in the package `probe` and seven defect forms in the generator.
- **A binding of a partner to a new subject.** `agrees-with-the-sibling` needs the binding of its partner to `h.New()`, which no other probe takes.

## Open questions

None.

## References

| What | Where |
|---|---|
| The shape catalog of eidos | <https://github.com/dokimasia/eidos/tree/main/eidos-plugin-shape/spec> |
