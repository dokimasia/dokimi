<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# Corpus claims ledger

What the suite generator claims for every classification the corpus
states, one disposition per row. This is the review artifact for the
per-axis batches: a stamp with no row here is the failure mode, and a
row's disposition is a decision someone made, not a green run nobody
read.

Every claimed row below has met the adversarial definition of done: a
minimal bug planted in its inmemory subject — one per claimed
classification — red the claim's own check by identity, and was
reverted (last full sweep: 2026-08-26, twelve plants, twelve reds:
zero-on-error, miss, idempotent twice, three dropped writes, a dropped
append, two smoke plants at witness instantiations, an ignored cancel,
and a loud repeat; batch 1 added seven more: zero-on-miss, nil-on-miss,
miss-is-zero, deterministic, repeat-survives, zeroes-on-error and
zeroes-on-miss; batch 2 added fifteen: cacheable, pure, monotonic,
monotonicreads, writesfollowreads, nilsafe, lifecycleafterclose,
orderafter, deleteremoves, sideeffect, poisonable, tamperevident,
retrysucceeds, defaultonerror and hooks). A row whose seat EXCUSES a
conditional claim — an in-memory subject that cannot reach the error
path — says so; those claims' red-capability is held by the generated
CanFail companions, and the plant discipline applies to the claims a
seat engages.

Two of batch 2's plants are worth naming, because they are the ones a
weaker check would have passed. The orderafter plant refused early for
an unrelated reason — a nil map, not the declared sentinel — and the
claim caught it by identity, which is exactly what the mixin's
`unready=` exists to rule out. The poisonable plant reported the
induced state once and cleared it, and the claim caught the flicker
rather than the first honest answer.

Dispositions:

- **claimed** — a rule binds the classification's claim at the
  example-based floor; the corpus fixture is suite-wired and its
  header, refusal paths and planted proof were reviewed.
- **refuses** — the classification's gates are verified to refuse with
  the right wording where a signature cannot carry the claim; noted on
  the claimed row rather than carried separately.
- **model tier** — the claim needs generated sequences, a reference, or
  laws; no honest example-floor form exists. The stamp stays
  header-listed as unclaimed until that tier lands.
- **door** — the claim needs a capability door (clock, crash, induce)
  before any floor form is honest.
- **pending** — not yet reviewed; the batch that owns it is named.
- **no suite** — the fixture cannot carry one, with the reason.

## Language axis (batch 0 — reviewed)

| Fixture | Disposition | Notes |
|---|---|---|
| generic | claimed | The whole value surface renders at the derived witness instantiation `Store[string, int]`; the header states the types every claim is proven at. |
| genericbound | claimed | Authored witnesses (`//testkit:witness K=int V=….Score`); the qualified witness import is registered by the render path. |
| receivercollision | no suite | The parameter collision is the fixture's purpose, and one role cannot hold two types; the decline naming both declarations is the verified behaviour. |
| seededreader | claimed | The fixture key is the `//testkit:default` the Key type declares, spelled as a conversion of the role's own type; the seat seeds the catalog under that key and excuses both zero-on-error claims by identity — the catalog answers the zero for an unseeded key, and its count answers from memory. |
| sampled | claimed | Authored values, load-bearing: Apply validates the tag, a refused first call leaves the idempotent claim vacuous, and the census fails the run — so only the `//testkit:sample` value keeps the suite green (witnessed by a config-override probe). Level carries sample and default at once, pinning the precedence — sample over default over derivation. |
| embedded, embeddedforeign, multireturn, namedreturns, nocontext, partnernaming, roledtypes, variadic | pending (batch 4) | |
| function | no suite | A function type declares no method set to derive from. |

## Detector axis (batch 1 — reviewed)

All 23 fixtures are suite-wired, every generated header carries the
per-method "what the declaration carries" block, and eidos's full
detector vocabulary is covered — closer and deleter gained fixtures in
this batch. A floor claim is stated only where it is universal for the
shape; everything conditional or partner-needing is the model tier's.

| Classification | Disposition | Notes |
|---|---|---|
| reader | claimed | zero-on-error, gated on value+sample+comparability; refusal wordings pinned; planted. |
| aggregator | claimed | zero-on-error under the shared licence; the in-memory seat excuses it — the reductions never fail. |
| answeringwriter | claimed | zero-on-error; excused on the in-memory seat, whose write always succeeds. |
| streamconsumer | claimed | zero-on-error on Ingest; the Source values are authored via `//testkit:sample`; excused on the in-memory seat. |
| multireader | claimed | zeroes-on-error — both slots owe their zeros beside the error; planted. |
| multiaggregator | claimed | zeroes-on-error; excused on the in-memory seat. |
| lookup | claimed | zeroes-on-miss — both slots beside the false flag; planted. |
| readerwithbool | claimed | zero-on-miss — the comma-ok discipline; planted. |
| readernoerror | claimed | miss-is-zero — the value is the only channel a miss can ride; planted. |
| pointerreader | claimed | nil-on-miss; cancel refuses in the header (no error channel); planted. |
| pure | claimed | deterministic — a repeat agrees, the law the shape's own doc carries; planted. |
| voidlifecycle | claimed | repeat-survives — the void teardown's one remaining law; planted. |
| lifecycle | model tier | The closed-sentinel law spans operations after teardown; the fixture's authored idempotent mixin is claimed and engaged. |
| closer | model tier | New fixture closing the vocabulary gap. The name-gated bare teardown; a repeat may legitimately answer differently, so no floor claim is universal. |
| deleter | model tier | New fixture closing the vocabulary gap. Removal-observability needs a read no signature names. |
| writer, compositewriter, multiargwriter, mutator | model tier | Write-observability needs a partner no signature names; multiargwriter's Get is claimed under the reader licence and engaged. |
| batchreader | model tier | Batch semantics need the reference; the slice answer refuses comparability at the floor. |
| predicate | model tier | Determinism-at-rest is not universal for the shape — an expiry probe drifts legitimately. |
| poisonaccessor | model tier | The latch outlives a call; the floor cannot poison. |
| streamreader | model tier | Drain laws; a stream may be consume-once, so two drains agreeing is not universal. |

## Mixin axis (batch 2 — reviewed)

All 60 fixtures are suite-wired and seated. Fifteen classifications
gained a floor claim; the rest are the model tier's or the clock
door's, each with its reason. Three of eidos's mixins owe
documentation rather than a check by their own vocabulary's word —
they are reported apart from the unclaimed gaps, so a reader is not
sent looking for a rule that is deliberately never coming.

| Classification | Disposition | Notes |
|---|---|---|
| notfound | claimed | miss, sentinel bound as data. |
| idempotent | claimed | repeat-is-quiet; the idempotentclose fixture carries the same claim on a teardown. |
| readafterwrite | claimed | write-then-read; partner named by `write=`. |
| causal, readyourwrites, monotonicwrites | claimed | write-then-read at the example floor: one client, one write, the read that follows; the ordering laws beyond that floor are the model tier's. Their answering Store also carries the answeringwriter zero-on-error claim, excused on the in-memory seats. |
| cacheable | claimed | A repeat agrees — the law the mixin's own doc states, over whichever value-answering shape the method has. The subject computes from the key, so the error path is excused by identity. |
| pure | claimed | The same agreement claim under the behavioural stamp, beside the structural detector's; both engage on the fixture. |
| monotonic | claimed | never-decreases: two reads with nothing between, gated on a builtin ordered result — a named type over an int is comparable, and its order is a claim about the declaration rather than the type. |
| monotonicreads | claimed | The per-client floor: one write, two reads, stamps held in order by the member `version=` names. Trace-replaying forms are the model tier's. |
| writesfollowreads | claimed | write, read, write again: the second stamp must exceed what the read observed. Needs a write that ANSWERS the stored state — a silent writer refuses, since the stamp cannot be read off it. |
| nilsafe | claimed | The nilable argument goes nil and the call must survive; gated on a parameter that can be nil at all. |
| lifecycleafterclose | claimed | The closed sentinel by identity, with both halves authored: which method closes, and which error every later call owes. |
| orderafter | claimed | The unready sentinel by identity — an implementation failing early for an unrelated reason must not pass as ordering enforcement, which is the whole reason `unready=` exists. |
| deleteremoves | claimed | The establishing write is DERIVED (the one sibling drawing the deleted roles), so the post-delete read cannot hold by vacancy; the miss is judged against the declared sentinel. |
| sideeffect | claimed | The named observation must move across the call; gated on an observation answering a comparable value. |
| poisonable | claimed | The latch, not the flicker: induced once, then reported on the call AND its repeat. |
| tamperevident | claimed | store, verify clean, tamper, verify again — every step before the last is vacuous on refusal, because each is a question about the fixture. |
| retrysucceeds | claimed | Converges within the declared `attempts=`; an unparsable bound refuses in the header naming the value. |
| defaultonerror | claimed | The declared default beside the error, read from the result type's own `//testkit:default`. REPLACES the reader's zero-on-error licence — the two contradict on one method, and the authored statement wins. |
| hooks | claimed | A hook installed through the named registrar must fire; the adapter synthesizes the callback from the registrar's own parameter types. |
| errors, scope, deprecated | documentary | Their own eidos vocabulary says no check is owed: errors "derives nothing by design", scope "owes documentation, not a check" (partition carries the checkable half), deprecated is an annotation the header surfaces. Reported apart from the unclaimed. |
| ttl, ttlperwrite, timeout, scheduled, windowed, pointintime, timeaware, eventually | door | Time-reading claims; wait on the clock door. |
| accumulates, associative, commutative, conservative | model tier | Algebraic folds: the claim is over regroupings and reorderings of a generated sequence, which the floor cannot draw. |
| atomic | model tier | A failed write leaving no partial state needs a write that CAN fail; `read=` now names the observer the model tier reads through. |
| bounded, boundedcache | model tier | A ceiling needs enough drawn writes to approach it; one call cannot cross a bound. |
| concurrent, concurrentreaders | model tier | Race freedom is the detector's business, not a sequential floor's. |
| crdtmerge | model tier | Convergence compares two replicas merged in opposite orders. |
| indexed | model tier | The bound is a runtime fact about the seeded subject, so the draw needs the sizing method read first. |
| injectionsafe, xsssafe | model tier | Adversarial payload corpora, not one derived draw. |
| integrationonly | model tier | Needs the infrastructure door before any run is honest. |
| leakfree | model tier | A balance over repeated acquire/release cycles. |
| noduplicates, overmatch, permutation, stableorder | model tier | Drain laws over a generated collection. |
| partition | model tier | Isolation needs two scopes written and read across. |
| restrictedpool | model tier | The pool-provenance fixture: its point is what a run may draw, which the floor does not vary. |
| sample | model tier | The builder names the input a sequence draws from. |
| serializable, snapshotisolation | model tier | Anomalies are patterns across transactions, not properties of one call. |
| sticky | model tier | Resolution stability across a sequence of reads and writes. |
| streamreflectsmutations | model tier | A stream opened after a mutation, compared against the mutation set. |
| total | model tier | Answering for every input in a domain needs the domain enumerated; the fixture's own Classify carries the reader licence, excused on the in-memory seat. |
| validates | model tier | Rejection needs invalid inputs the derivation does not produce. |
| wrappedvia | model tier | Delegation is a fact about the implementation's insides; nothing at the interface floor says whether one method ran through another. |

## Contract axis (batch 3)

| Classification | Disposition | Notes |
|---|---|---|
| chain (append role) | claimed | replay-grows, count form; the entry's content is the model tier's question. Other roles pending. |
| appender, batch-writer, cache, cas, circuit-breaker, codec, codec-lossy, cursor, cursor-open, if-absent, if-match, leader-election, lease, outbox, pagination, persister, pool, publisher (and delivery variants), rate-limit, saga, singleflight, transaction, tx, updater, upserter, watcher, workflow | pending (batch 3) | Role-by-role; per-role unclaimed accounting lands with the first second-role rule. |

## Composite (batch 4)

| Fixture | Disposition | Notes |
|---|---|---|
| causal-chain | claimed (chain/append) | The bare causal mixin beside it refuses on the result shape, as reviewed. The replica seat proves the open surface: a subject that cannot append excuses the chain claim by identity (Harness.Excuse), a replica-only run declines it with Without+Because, and the observed census holds excused apart from passed. |
| batched-mixins, leased-idempotent-writer, paginated-reader, tx-with-retry | pending (batch 4) | Interaction review after their axes settle. |
