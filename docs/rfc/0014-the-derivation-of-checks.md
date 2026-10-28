---
rfc: 0014
title: The derivation of checks from the conformance spec
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Draft
created: 2026-10-28
updated: 2026-10-28
discussion: none
supersedes: RFC-0003 in part
superseded-by: none
produces-adr: ADR-0015
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# RFC-0014: The derivation of checks from the conformance spec

## Summary

RFC-0003 designs the derivation of the suite generator: `Derive` over the `Facts` of each callable, the conditions of the gates, the claims table with 73 claims, the stances of the unclaimed classifications, the edges between a deeper and a shallower check, and the falsification of each check by the reason code of its claim. This RFC completes that design where the conformance spec of RFC-0002 grew beyond it, and where its planted defects cannot work:

- Five partner modes join the two sources of a partner that RFC-0003 has.
- A probe sets the reason code of each assertion on the handle, and the spec lists the reasons of each probe.
- The declaration states the machineries that its consumer supplies.
- A claim entails a list of shallower claims.
- A planted defect wraps a fresh subject of the application, and `conform.Run` falsifies each derived check, as ADR-0015 decides.

## Motivation

dokimi exists to derive checks from declarations, and today it derives none. Every check of a suite is a function that the consumer writes and binds with `//dokimi:check Store.Get defect=forgetful` (`generator/suite/doc.go:21-34`). RFC-0003 designs the derivation, and no code implements it.

Four parts of that design need a change before it can be built:

- **Partners.** RFC-0003 finds a partner through a param of the classification (`FromPartner`) or a role of the same contract (`FromRole`). The partners that testkit resolved need five more modes: the halves of a pair the other way round, a partner that a field of a result selects, a partner that a probe calls twice, and the one writer that puts the subject into a state.
- **Reason codes.** RFC-0003 has the handle derive the reason from the assertion of a failure. One probe can fail at two places of one assertion for two reasons, so the probe states the reason.
- **Machineries of the consumer.** RFC-0003 makes a claim unclaimed when the suite generator or the runtime does not supply its machinery. A reference implementation, faults on demand, a second subject and a breakable medium come from the consumer. The derivation cannot see them in the source, so the declaration states them.
- **Planted defects.** RFC-0003 builds each defect from the double of the declaration, and ADR-0007 removed the double from every suite. A double is not a correct implementation, so every check whose law depends on state fails on every defect built from it.

## Detailed design

### The changes to RFC-0003

| Part of RFC-0003 | Change |
|---|---|
| `Claim.Machinery`, one machinery | `Claim.Machineries`, a list whose first entry is the deepest |
| `Claim.Entails`, one segment | A list of segments |
| `Licence` with `Callable`, `Context`, `Shape`, `Mixin`, `Contract`, `Role` and `With` | Gains `NeedsParam` and `Conflicts` |
| `Source` with `FromPartner` and `FromRole` | A partner has one of seven modes, and `FromPartner` reads the mode of the partner |
| Step 2 of `Derive`: a machinery that the suite generator or the runtime does not supply | Also a machinery of the consumer that the param `machineries` of the declaration does not state |
| The reason of a failure from its assertion | The reason that the probe sets with `Handle.Expect` |
| A planted defect built from the double | A wrapper of a fresh subject of the application |
| Falsification of a derived check in the generated test of the suite | Falsification inside `conform.Run`, before the check runs on the subject |
| The gates of :284-296 | Gain `value-beside-error`, which passes when the error model of the callable is `last-return`, so a failed call returns its values beside the error. A call whose error model is `thrown`, `raised` or `result-type` returns no value when it fails, whatever the zero of its type |

This RFC leaves the rest of RFC-0003 unchanged: the facts, the order of `Derive`, the claims and the stances, the identities, the lock, the checks of the consumer, the isolation of a defect and the edges.

### Partner modes

```go
// PartnerMode is how the derivation finds the partner of a claim.
type PartnerMode string

// Partner resolves the partner of a claim on the callable of f, among the
// callables of decl. It returns the facts of the partner, or the gate that
// refuses the claim: the gate partner when no callable matches, and the gate
// one-partner when two do.
func Partner(mode PartnerMode, p spec.Partner, f *Facts, decl []Facts) (*Facts, *spec.Gate)
```

| Mode | Reads | In RFC-0003 |
|---|---|---|
| `named` | The param of the licence that names the partner, such as `write` of the mixin `read-after-write`. The licensed callable reads | `FromPartner` |
| `reader` | The same, with the halves swapped: the licensed callable writes, and the partner reads | New |
| `derived` | The result field that a param of the licence names, and the one other callable that takes the type of that field as a parameter | New |
| `role` | The callable with the named role of the same contract instance | `FromRole` |
| `interposed` | A `named` partner that the probe calls twice, and a field of the written value that the second call changes | New |
| `rewriter` | A `reader` partner that the probe calls twice, with the same field | New |
| `establishing` | The one other callable with the classification of the entry that takes every input of the licensed callable | New |

A value whose every field the read selects by refuses `interposed` and `rewriter` with the gate `rewritable`, because no second write could change what the read returns.

### The reason on the handle

```go
// Expect sets the reason code of the assertions that the probe calls next.
// The handle attaches it to the failure of such an assertion, and
// falsification reads it from there. Expect fails the check when the entry
// of the probe in the spec does not list reason.
func (h *Handle[S]) Expect(reason spec.ReasonID)
```

A probe that compares two results for two reasons calls `Expect` before each comparison. A check of the consumer does not call `Expect`. Its proof is the contract of its first failure, as RFC-0003 states.

### The generated check and its planted defect

```go
{
	ID:        StoreChecks.Read.ReadAfterWrite(),
	Statement: "Read returns what Put stored before it",
	Claim:     spec.ReadAfterWrite,
	Reason:    spec.UnobservedWrite,
	Law:       spec.AReadReturnsTheWriteBeforeIt,
	Apply: func(h *StoreHandle) {
		probe.WriteThenRead(h,
			func(ctx context.Context, k string) (catalog.Item, error) { return h.Subject().Read(ctx, k) },
			func(ctx context.Context, k string, v catalog.Item) error { return h.Subject().Put(ctx, k, v) },
			storeKeys(h), storeItems(h))
	},
	Defect: func(h *StoreHandle) catalog.Store {
		return storeForgetsWrites{Store: h.Fresh()}
	},
},
```

The defect form `forgets-writes` writes a type that embeds the subject type and overrides the callable of the form:

```go
// storeForgetsWrites is the planted defect of Store.Read/read-after-write: a
// Store whose Put returns success and stores nothing.
type storeForgetsWrites struct{ catalog.Store }

// Put returns success and stores nothing.
func (storeForgetsWrites) Put(context.Context, string, catalog.Item) error { return nil }
```

`Fresh` returns a fresh subject from the factory of the application, also under falsification:

```go
// Fresh returns a fresh subject from the factory of the application. Under
// falsification Subject and New return the planted defect, and Fresh still
// returns a subject of the factory, which a planted defect wraps.
func (h *Handle[S]) Fresh() S
```

### Falsification inside the application

`conform.Run` falsifies each derived check before it applies the check to the subject:

1. It runs the check on its planted defect over a fresh subject. A failure with the reason of the claim proves the check. Another reason is `wrong-reason`, and a pass is `passed`.
2. It runs every other check on the same defect. Each of them must pass, apart from the checks that share the defect and the deeper checks of an edge of the check, as RFC-0003 states.
3. It runs each edge of the check in both directions, as RFC-0003 states.
4. It applies the check to the subject.

A failure of falsification fails the test of the check before the check runs on the subject. `Run` writes the falsification record and the record of the subject. `conform.Falsify` in the generated test of the suite falsifies the checks of the consumer and the ceilings, whose defects need no subject.

A defect that overrides a partner, such as `forgets-writes` on `Put`, also fails the other checks whose law needs that partner, such as a check of `Put` itself. The derivation gives each such check the same defect, so the exemption of the checks that share a defect covers them.

### The runtime

The module `go.dokimi.dev/conform` gains the packages `spec` and `probe`, as RFC-0002 and RFC-0005 state:

- `spec` decodes the vendored JSON files of a release of dokimi and declares a constant for each law, claim, probe and reason code.
- `probe` implements every probe of the spec. A probe asserts through assert-go, passes the statement of the check as the contract of each assertion, and calls `Expect` before each assertion.

`conform.Check[S]` gains three fields:

```go
// Claim is the claim of the conformance spec that the check binds, and
// empty for a check of the consumer.
Claim spec.ClaimID

// Reason is the reason code of the planted defect of the check, and empty
// for a check of the consumer.
Reason spec.ReasonID

// Law is the law that the check witnesses, and empty for a check of the
// consumer.
Law spec.LawID
```

### Migration of the corpus

The Go corpus contains 168 fixtures whose checks are written by the consumer. Each fixture whose laws the spec derives replaces those checks with the derived checks of the same identities, and keeps the checks that the spec does not state. The ledger of the corpus records for each fixture which checks are derived.

## Alternatives considered

### Partner resolution in Go functions

Each claim with a partner gets a Go function that searches the declaration, as testkit did with about 40 resolver functions in seven modes.

**Why not:** a function per claim is a second statement of the claim in code. Each mode is one function that every claim of the mode shares, and a claim states its mode as data.

### One query form for every partner

Each partner is a query over the other callables: a classification and the inputs that must have the same type.

**Why not:** a query covers the modes `named` and `establishing`. It cannot state a partner that a field of a result type selects, a partner that the probe calls twice with a changed field, or a partner of the other half of a pair. testkit needed all three.

### The reason from the assertion

The entry of a probe maps each assertion to one reason, and the handle reads the assertion of the failure record, as RFC-0003 designs.

**Why not:** a probe that compares two results through `equal` for two reasons would split into two probes that repeat one procedure.

### Planted defects from the double

ADR-0015 records the alternatives of the double of RFC-0003 and of the double that delegates to the generated fake.

## Drawbacks

- **A derived check is proven only in an application.** The generated test of a suite has no subject, so it proves the checks of the consumer and the ceilings alone. The corpus and each test of a consumer prove the derived checks.
- **A probe exists once per runtime.** Each probe of the spec has a Go implementation in conform and an implementation in each later runtime.
- **The corpus migrates.** Each fixture with derived laws rewrites its checks to derived ones, and its goldens change.
- **A defect form per kind of break.** Forms are shared between claims, but each new kind of break needs a form, a template in each language and a sentence in `defects.yaml`.

## Open questions

None.

## Unresolved and future work

- Claims over the fields of a struct, such as the invariants of a value, are not proposed.

## References

| What | Where |
|---|---|
| The shape catalog of eidos | <https://github.com/dokimasia/eidos/tree/main/eidos-plugin-shape/spec> |
| testkit, the claims table and the resolvers of partners | <https://github.com/thesm-os/testkit/tree/main/docs> |
