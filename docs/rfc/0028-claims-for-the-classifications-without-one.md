---
rfc: 0028
title: Claims of a single caller for the classifications without one
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-31
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

# RFC-0028: Claims of a single caller for the classifications without one

## Summary

Eight classifications of the shape catalog have no claim in the spec and the stance `gap`. Each of them states a law that one caller can check with the machinery `nothing-beyond-a-caller`: the shapes `batch-reader`, `poison-accessor` and `stream-reader`, the mixins `indexed` and `wrapped-via`, and the roles `chain/verify`, `lease/release` and `pool/stats`. This RFC gives each of them a claim with its law, its probe, its reason code and its planted defect, so that the unclaimed table of the spec has no entry with the stance `gap`.

## Motivation

The derivation lists a classification without a claim as unclaimed, and the stance `gap` states that the classification implies something that no claim witnesses. A consumer that writes `//dokimi:suite` on a declaration with such a classification does not get a check of it, and its suite record states the gap.

Each classification of the catalog states how a check observes it, in the field `observation` of its entry. Each of the eight observations needs one caller and the samples of RFC-0014:

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

The spec states one of the eight wrongly today. The note of `mixin:wrapped-via` names a law about the wrapping of an error, while the catalog states a delegation to a sibling. The fixture of the mixin in the corpus follows the note, so it wraps an error and delegates nothing.

## Detailed design

### The claims

| Claim | Licence | Partners and params | Gates | Law | Probe | Reason | Defect form |
|---|---|---|---|---|---|---|---|
| `reads-each-key` | The shape `batch-reader` | `write`: the establishing writer | `cancellation`, `list-value` | `batch-read-returns-each-stored-value` | `reads-what-was-stored` | `lost-a-key` | `drops-the-last-key` |
| `fresh-reports-no-poison` | The shape `poison-accessor` | — | `failure` | `fresh-receiver-is-healthy` | `reports-no-error` | `false-poison` | `reports-poison` |
| `streams-what-was-stored` | The shape `stream-reader` | `write`: the establishing writer | the gates of `over-match` | `drain-keeps-what-went-in` | `keeps-what-went-in` | `shrunk-drain` | `drains-nothing` |
| `respects-its-bounds` | The mixin `indexed` with its param `by` | `by`: named by `by` | `failure` | `positions-are-below-the-size` | `bounded-by-the-size` | `out-of-bounds` | `admits-the-size` |
| `agrees-with-its-sibling` | The mixin `wrapped-via` with its param `fn` | `fn`: named by `fn` | `returns-unaided`, `one-value` | `wrapper-returns-what-its-sibling-returns` | `agrees-with-the-sibling` | `diverged-from-the-sibling` | `stops-delegating` |
| `verifies-what-it-appended` | The contract `chain`, role `verify` | `append`: the role `append` | `failure` | `appended-chain-verifies` | `verifies-after-appends` | `false-tamper` | `reports-tamper` |
| `release-frees-the-lease` | The contract `lease`, role `acquire`, which conflicts with the mixin `idempotent` | `release`: the role `release` | `cancellation`, `failure` | `released-lease-is-granted-again` | `grants-after-release` | `kept-after-release` | `release-keeps-the-lease` |
| `counts-each-loan` | The contract `pool`, role `stats` | `get`: the role `get` | `one-value` | `stats-move-with-a-loan` | `moves-with-a-loan` | `uncounted-loan` | `frozen-stats` |

Each claim lists `nothing-beyond-a-caller` among its machineries, and the depth of each new law is `nothing-beyond-a-caller`. `streams-what-was-stored` binds the law, the probe, the reason code and the defect form of `over-match` to the shape `stream-reader`, because the observation of the shape is the observation of that mixin.

`release-frees-the-lease` binds to the acquire of the lease and takes the release as its partner, as the claim `lease` binds to the acquire. An acquire with the mixin `idempotent` grants a held lease again, so a release that does nothing could not fail the check, and the licence conflicts with that mixin as the licence of `lease` does. The role `release` is the partner of the claim, so it does not need an entry in the unclaimed table.

The statement of each claim:

| Claim | Statement |
|---|---|
| `reads-each-key` | `{callable} returns the value that {write} stored under each key` |
| `fresh-reports-no-poison` | `{callable} reports no error for a fresh receiver` |
| `streams-what-was-stored` | `{callable} yields each value that {write} stored` |
| `respects-its-bounds` | `{callable} admits each position below the size that {by} reports and refuses the size` |
| `agrees-with-its-sibling` | `{callable} returns what {fn} returns for one input` |
| `verifies-what-it-appended` | `{callable} reports no error for a chain that {append} produced` |
| `release-frees-the-lease` | `{callable} grants a lease again after {release} released it` |
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
| `release-keeps-the-lease` | partner `release` | The release returns no error and does nothing. | 3 |
| `frozen-stats` | callable | The stats return the zero value on every call. | 2 |

### The probes

The package `probe` of conform gains seven probes, which follow the conventions of RFC-0027. `streams-what-was-stored` uses the probe `KeepsWhatWentIn` that exists.

```go
// ReadsWhatWasStored writes two values through write, reads both keys through
// call in one batch, and fails unless the batch returns the two written values
// in the order of their keys. A write that fails ends the probe.
func ReadsWhatWasStored[S, V any](h *conform.Handle[S], call Call[[]V], write Call[V])

// ReportsNoError calls call on the fresh subject, and fails when the call
// returns an error.
func ReportsNoError[S any](h *conform.Handle[S], call Act)

// BoundedByTheSize reads the size through size, calls call at the first and the
// last position and at the size, and fails unless the first two calls succeed
// and the call at the size fails. A size of zero calls call at the size alone.
func BoundedByTheSize[S any](h *conform.Handle[S], call func(ctx context.Context, position int) error, size Call[int])

// AgreesWithTheSibling calls call on the subject and sibling on a new subject
// of the factory, with one input. It fails when one of the two calls fails
// and the other does not, and when both succeed with different values. A
// value beside an error has no meaning, so two failed calls agree.
func AgreesWithTheSibling[S, V any](h *conform.Handle[S], call, sibling Call[V])

// VerifiesAfterAppends appends twice through appends, then calls call, and
// fails when the verify returns an error. An append that fails ends the
// probe.
func VerifiesAfterAppends[S any](h *conform.Handle[S], call, appends Act)

// GrantsAfterRelease acquires through call, releases through release, and
// acquires through call again, and fails when the second acquire returns an
// error. A first acquire or a release that fails ends the probe.
func GrantsAfterRelease[S any](h *conform.Handle[S], call, release Act)

// MovesWithALoan reads the stats through call, lends an object through get,
// and reads the stats again, and fails when the two stats are equal or a read
// of the stats fails. A get that fails ends the probe.
func MovesWithALoan[S, V any](h *conform.Handle[S], call Call[V], get Act)
```

| Probe | Inputs | Reasons | Vacuous condition |
|---|---|---|---|
| `reads-what-was-stored` | `call`: call, `write`: call | `lost-a-key` | A write fails. |
| `reports-no-error` | `call`: call | `false-poison` | Never, because the probe makes its one call. |
| `bounded-by-the-size` | `call`: call, `size`: call | `out-of-bounds` | The size fails. |
| `agrees-with-the-sibling` | `call`: call, `sibling`: call | `diverged-from-the-sibling` | Never, because the probe calls both. |
| `verifies-after-appends` | `call`: call, `append`: call | `false-tamper` | An append fails. |
| `grants-after-release` | `call`: call, `release`: call | `kept-after-release` | The first acquire or the release fails. |
| `moves-with-a-loan` | `call`: call, `get`: call | `uncounted-loan` | The get fails. |

### The bindings

The suite generator binds the inputs of the seven probes with the forms of input that it has, and with four that it gains:

| Claim | Input | Binding |
|---|---|---|
| `reads-each-key` | `write` | The write passes the first sample on its first call and the second sample on every later call, and returns what it wrote |
| `reads-each-key` | `call` | The batch read passes the key of the value of each of the two writes, in the order of the writes |
| `streams-what-was-stored` | `drain` | The call drains the stream into a list |
| `respects-its-bounds` | `call` | The call takes the position and passes it to each input of an integer type, and the first sample to each other input |
| `agrees-with-its-sibling` | `sibling` | The call of `fn` runs on `h.New()` in place of `h.Subject()` |

`AgreesWithTheSibling` calls the sibling on a new subject, so a write through the subject does not change the state that the sibling reads. The position of `BoundedByTheSize` is an int, and the call converts it to the type of each integer input, because the catalog states that the mixin applies to every integer input.

A batch read pairs with its writer through the element of its list. Each writer persists one value. The written value is the input of the write whose type is the element of the list that the batch read returns. Each key of the variadic input takes the field of the written value with the name of the input without its last letter `s`. The comparison of the two names ignores case, and the input `keys` takes the field `Key`. The second write sets that field to its second value, and the two writes store two keys. A write stores the keys themselves where its written value has the name and the type of the keys, as `Add(ctx, tag string)` stores the `tags` of `Has(ctx, tags ...string)`. Each key then takes the value of one write. The gate `partner-arguments` refuses a batch read whose keys select neither a field of the written value nor the written value itself.

### The drain of a stream

A Go callable returns a stream as an `iter.Seq` or an `iter.Seq2`. The call of the input `drain` collects the stream into a list:

- each element of an `iter.Seq`
- the first value of each pair of an `iter.Seq2` whose second value is an `error`, until the first error, which ends the drain
- the `probe.Pair` of each pair of any other `iter.Seq2`

The planted defect `drains-nothing` of a callable that returns a stream returns an empty stream, because a range over a nil function panics.

### The planted defects

Each form other than `drops-the-last-key` takes a body of a planted defect that the generator has. `reports-poison` and `reports-tamper` fail with the error that the runtime generates, `admits-the-size` and `stops-delegating` return their second samples beside no error, and `release-keeps-the-lease` and `frozen-stats` return their zero results. `drops-the-last-key` calls the callable of the fresh subject and removes the last element of the list that it returns.

### The fixtures of the corpus

The corpus has a fixture for each of the eight classifications. The fixtures of the shape `batch-reader` and of the mixin `wrapped-via` change:

- The batch reader of `iface/detector/batchreader` gains a writer `Put`, the partner of `reads-each-key`.
- The fixture `iface/mixin/wrappedvia` delegates `TracedGet` to `Get` under the mixin `wrapped-via`, as the catalog states the mixin.

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
- **Eight more claims.** Each claim adds a probe, a reason, a law and a defect form to the spec, and a check to each suite with the classification: seven probes in the package `probe`, seven defect forms and one body of a planted defect in the generator.
- **Four forms of binding.** The generator gains the batch read of two keys, the drain of a stream, the position of an indexed call and the call of a sibling on a new subject, and the derivation a rule that pairs a batch read with its writer.
- **A rule of names.** A batch read finds the key of the written value by the name of its variadic input without its last letter `s`. A batch read whose input has another name, such as `entries` for the field `Entry`, is refused by the gate `partner-arguments`.

## Open questions

None.

## References

| What | Where |
|---|---|
| The shape catalog of eidos | <https://github.com/dokimasia/eidos/tree/main/eidos-plugin-shape/spec> |
