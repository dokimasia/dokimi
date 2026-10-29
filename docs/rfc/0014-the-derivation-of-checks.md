---
rfc: 0014
title: The derivation of checks from the conformance spec
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-28
updated: 2026-10-28
discussion: none
supersedes: RFC-0003 in part
superseded-by: none
produces-adr: ADR-0015, ADR-0017
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
- The isolation of a derived defect and its edges run only under the option `conform.Isolate`, which the corpus passes for each control, as ADR-0017 decides.

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
| The directive `suite` has no parameters | It has the params `machineries`, a list of the machineries that the consumer supplies, and `normalize`, the function that maps a result to the value that a history records |
| `Value`, a parameter or a return with its shape, samples and zero | Gains `Fields`, the fields of a struct type with their names and shapes from `Bound.FieldTypeOf`, which the modes `derived`, `interposed` and `rewriter` read |
| `Claim.Generator`: `suite`, `enum` or `sentinel` | Also `sim` |
| `Claim` binds to a callable | Gains `Scope`: `callable`, the default, or `declaration` |
| `Claim.Machinery`, one machinery | `Claim.Machineries`, a list whose first entry is the deepest |
| `Claim.Entails`, one segment | A list of segments |
| `Licence` with `Callable`, `Context`, `Shape`, `Mixin`, `Contract`, `Role` and `With` | Gains `NeedsParam` and `Conflicts` |
| `Source` with `FromPartner` and `FromRole` | `Claim` gains `Partners`, each with a name and one of seven modes. `Source` gains `Partner`, the name of the partner of a `FromPartner` source. The mode `role` replaces `FromRole` |
| Step 2 of `Derive`: a machinery that the suite generator or the runtime does not supply | Also a machinery of the consumer that the param `machineries` of the declaration does not state |
| The reason of a failure from its assertion | The reason that the probe sets with `Handle.Expect` |
| A planted defect built from the double | A wrapper of a fresh subject of the application |
| Falsification of a derived check in the generated test of the suite | Falsification inside `conform.Run`, before the check runs on the subject |
| The isolation of a planted defect and the edges, in the generated test of the suite | Inside `conform.Run` under the option `conform.Isolate`, which the corpus passes for each control |
| The gates of RFC-0003 | Gain `returns-beside-error`, which passes when the error model of the callable is `last-return`, so a failed call returns its values beside the error. A call whose error model is `thrown`, `raised` or `result-type` returns no value when it fails, whatever the zero of its type. The claims `zero-on-error` and `zeroes-on-error` list the gate |

This RFC leaves the rest of RFC-0003 unchanged: the order of `Derive`, the claims and the stances, the identities, the lock and the checks of the consumer.

### Partner modes

```go
// PartnerMode is how the derivation finds the partner of a claim. The
// package spec declares it as the type of the field Mode of each Partner.
type PartnerMode string

// Partner resolves the partner p of a claim on the callable of f, among the
// callables of decl, by the mode of p, and binds the arguments of the read.
// It returns the gate that refuses the claim instead: the gate partner-found
// when no callable matches, one-partner when two do, partner-arguments when
// a parameter of the read selects nothing of the write, and rewritable when
// a second write could not change what the read returns.
func Partner(p *spec.Partner, f *Facts, decl []Facts) (*Pair, *spec.Gate)

// Pair is the partner of a claim, with the read, the write, the index of the
// written value in the inputs of the write, the arguments of the read, and
// the field that a second write changes.
type Pair struct {
	Partner     *Facts
	Read, Write *Facts
	Changed     string
	Arguments   []Argument
	Written     int
}
```

| Mode | Reads | In RFC-0003 |
|---|---|---|
| `named` | The param of the licence that names the partner, such as `write` of the mixin `read-after-write`. The licensed callable reads | `FromPartner` |
| `reader` | The same, with the halves swapped: the licensed callable writes, and the partner reads | New |
| `derived` | The result field that a param of the licence names, and the one other callable that takes the type of that field as a parameter | New |
| `role` | The callable with the named role of the same contract instance | `FromRole` |
| `interposed` | A `named` partner that the probe calls twice, and a field of the written value that the second call changes | New |
| `rewriter` | A `reader` partner that the probe calls twice, with the same field | New |
| `establishing` | The one other callable with the classification of the entry whose written value the read of the licensed callable selects by the rule of the arguments | New |

A value whose every field the read selects by refuses `interposed` and `rewriter` with the gate `rewritable`, because no second write could change what the read returns.

### The arguments of a partner

The probe calls the licensed callable and its partner on one state, so the read must select what the write stored. The derivation binds the arguments of the two calls by one rule:

1. A parameter that both callables take, with the same name and the same type, takes one sample in both calls.
2. Each other parameter of the read takes the field of the written value with its name and its type, such as the argument `id` of `Get(ctx, id)` from the field `ID` of the session that `Put(ctx, s)` stored.
3. A parameter in the role context takes the context of the check.

A read parameter that neither rule binds refuses the claim with the gate `partner-arguments`, whose remedy is to align the name of the parameter with the write. The read is the licensed callable in the modes `named`, `derived` and `interposed`, and the partner in the modes `reader`, `rewriter` and `establishing`. For a partner of the mode `role`, the read is the role that returns a value and the write the role that returns none. Where both roles return a value or neither does, as the roles `cache` and `backing` of the contract `cache`, the rule binds nothing and the probe passes the result of one call to the other.

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
	Machineries: []spec.MachineryID{spec.NothingBeyondACaller},
	Apply: func(h *StoreHandle) {
		probe.WriteThenRead(h,
			func(ctx context.Context) (catalog.Item, error) { return h.Subject().Read(ctx, "a") },
			func(ctx context.Context) error { return h.Subject().Put(ctx, "a", catalog.Item{Name: "b"}) })
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
2. Under `conform.Isolate`, it runs every other check on the same defect. Each of them must pass, apart from the checks that share the defect and the deeper checks of an edge of the check, as RFC-0003 states.
3. Under `conform.Isolate`, it runs each edge of the check in both directions, as RFC-0003 states.
4. It applies the check to the subject.

A failure of falsification fails the test of the check before the check runs on the subject. `Run` writes the falsification record and the record of the subject. `conform.Falsify` in the generated test of the suite falsifies the checks of the consumer and the ceilings, whose defects need no subject. It also isolates their defects and checks their edges.

```go
// Isolate makes an application run every other derived check on the planted
// defect of each derived check, and check each edge in both directions. The
// corpus passes it for each control.
func Isolate[S any]() Option[S]
```

Isolation proves that the template of a defect form breaks one law alone, and a correct subject is the condition of that proof. The controls of the corpus are correct subjects, so the corpus proves every defect form of every claim that it binds. The subject of a consumer is not known to be correct. A subject that fails one check fails it again under every planted defect that wraps the subject, so isolation over that subject reports N−1 failures of planted defects before the application reports the defect of the subject. Isolation also costs N×N runs of checks for a suite of N derived checks, such as 900 runs for 30 checks.

A defect that overrides a partner, such as `forgets-writes` on `Put`, also fails the other checks whose law needs that partner, such as a check of `Put` itself. The derivation gives each such check the same defect, so the exemption of the checks that share a defect covers them.

### The runtime

The Go runtime gains the module `go.dokimi.dev/spec` and the package `probe` of the module `go.dokimi.dev/conform`, as RFC-0002 and RFC-0005 state:

- `spec` embeds the vendored copy of the conformance spec, declares a constant for each entry with its spelling in the naming table, and returns the reasons of each probe.
- `probe` implements every probe of the spec. A probe asserts through assert-go, passes the statement of the check as the contract of each assertion, and calls `Expect` before each assertion.

`conform.Check[S]` gains seven fields:

```go
// Claim is the claim of the conformance spec that the check binds, and
// empty for a check of the consumer.
Claim spec.ClaimID

// Probe is the probe of the claim, whose reasons Handle.Expect accepts,
// and empty for a check of the consumer.
Probe spec.ProbeID

// Reason is the reason code of the planted defect of the check, and empty
// for a check of the consumer.
Reason spec.ReasonID

// Law is the law that the check witnesses, and empty for a check of the
// consumer.
Law spec.LawID

// Machineries are the machineries of the claim, the deepest first, and nil
// for a check of the consumer.
Machineries []spec.MachineryID

// Shallow reports that the machineries of the claim do not include the
// depth of its law.
Shallow bool

// Entails are the identities of the shallower checks of the suite whose
// laws the law of the check implies, one for each edge of the derivation.
Entails []ID[S]
```

The runtime writes these fields into the suite record, because it does not read a claim at run time. The isolation and the edges of `Isolate` read `Entails`, and `Expect` reads the reasons of `Probe` from the copy of the spec.

### Migration of the corpus

The Go corpus contains 146 fixtures with the directive `suite`. Seven files of them bind 12 checks of the consumer with the directive `check`, and the other fixtures have no checks. Each fixture gains the derived checks of its classifications. The derived check of the same identity replaces a check of the consumer whose law the spec derives. A check whose law the spec does not state remains. The ledger of the corpus records for each fixture which checks are derived.

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

### Isolation in every application

`conform.Run` runs every other derived check on each planted defect and checks each edge on the subject of every consumer.

**Why not:** a defect of the subject fails its check under every planted defect that wraps the subject. The application then blames N−1 planted defects for one defect of the subject. The run also costs N×N checks per application. The template of a defect form is the same for every subject, so the controls of the corpus prove its isolation once. The verdict reverses when a defect form passes isolation on every control and fails it on a subject that passes its own suite.

## Drawbacks

- **A derived check is proven only in an application.** The generated test of a suite has no subject, so it proves the checks of the consumer and the ceilings alone. The corpus and each test of a consumer prove the derived checks.
- **Isolation only over the controls.** The application to the subject of a consumer does not show that a planted defect breaks one law alone. A defect form that breaks a second law only over some subjects passes the corpus.
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
