---
rfc: 0002
title: The conformance spec
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-19
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

# RFC-0002: The conformance spec

## Summary

The conformance spec defines what an implementation must do to conform to a shape, a mixin or a contract of the shape catalog of eidos. It contains the laws, the claims that each classification licenses, the gates of each claim, the probes with their inputs and reason codes, the machineries and the stances. It also fixes the grammar of a check identity, the suite record of a run and the contract between a generated file and its runtime.

assert-spec defines assertions over values and callables, and it defines no shape. Each probe asserts through the assert library of its language, which implements assert-spec. A run writes the call records of assert-spec and the suite records of the conformance spec into the artifact of its test runner.

People edit the YAML files of the directory `spec/` in the dokimi repository, and `go generate` renders them into JSON files, Go constants and a manifest of SHA-256 digests. The binary embeds the JSON files. Each runtime vendors them from a release tag of dokimi and checks its copy against the manifest.

## Motivation

The conformance of a shape is a decision of dokimi:

- A reader conforms when it returns the zero value beside an error.
- A cache conforms when it returns what its backing returns.
- A callable with the mixin `ttl` conforms when an entry stops being readable once its lifetime passes.

eidos detects the shape of a callable from its signature and its directives. It does not define how an implementation of the shape behaves. assert-spec defines assertions such as `equal` and `idempotent`, and every assert library implements them. The shape catalog is not part of assert-spec, so assert-spec cannot define which assertion a reader must pass, with which inputs and against which planted defect. One conformance spec defines all three for every language.

A generated test calls a probe by name, and falsification matches a failure by its reason code. Each name must mean the same thing in every runtime. When each runtime keeps its own list, a generated check calls a probe that one runtime spells differently, or a probe reports a code that another runtime lacks. The disagreement shows only when a consumer runs the generated tests. The service and the gates read the suite records of every runtime with one parser, so every runtime writes the same format.

People edit the definition of assert-spec as YAML. A script renders the JSON that the assert libraries read, and each library vendors the JSON with a manifest of digests. The conformance spec follows that design with two differences:

- JSON Schemas describe the rendered files, so a test of the binary validates each file against its schema.
- A runtime vendors the spec of a release tag of dokimi, because the tests that it runs come from a released binary.

## Detailed design

### The conformance spec, assert-spec and mutate-spec

| Definition | Defines | dokimi uses it for |
|---|---|---|
| The conformance spec | Laws, claims, gates, probes, reason codes, machineries, stances, the check identity and the suite record | The derivation, the probes of each runtime, falsification and the suite records of a run |
| assert-spec | Assertions and their names in each language, the failure record, the call record, the overlays, the property engine, the machines, the checkers of histories, the benchmark ceilings and the controlled clock | The assertions of each probe, and the call records of a run |
| mutate-spec | A mutation run and its record | The reports of dead and subsumed checks |

The conformance spec uses these parts of assert-spec:

| Part of assert-spec | Use in the conformance spec |
|---|---|
| The assertions | The assertions that the probes of each runtime call through its assert library |
| The failure record | The failure of a check in its suite record |
| The call record | The calls of each probe in a recorded run |
| The controlled clock | The machinery `clock-the-run-advances` |
| The property engine | The machineries `generated-inputs` and `generated-sequences`, and the fuzz targets of the claims of a suite |
| The machines | The machinery `generated-sequences`, whose actions the judges of a history and the comparison with a reference implementation run |
| `linearizable`, `serializable` and `snapshot-isolation` | The machinery `concurrent-callers`, which records a history and checks it |
| The benchmark ceilings and `max-allocs` | The budgets of a callable |

The conformance spec does not vendor assert-spec. A probe is host code in each runtime, so the build of the runtime fails on a call of an assertion that its assert library lacks. The spec does not state which probes a language implements.

### Layout

```text
spec/
├── VERSION            the version of the spec, such as 1.0.0
├── laws.yaml          the laws, edited by people
├── claims.yaml        the claims of each generator
├── gates.yaml         the gates
├── probes.yaml        the probes
├── reasons.yaml       the reason codes
├── defects.yaml       the forms of the planted defects
├── references.yaml    the reference forms of the verified fakes
├── machineries.yaml   the machineries
├── stances.yaml       the stances, and the stance of each classification without a claim
├── laws.json          rendered from laws.yaml, as each YAML file is rendered
├── …
├── schemas/           a JSON Schema for each rendered file and for the suite record
├── manifest.json      the SHA-256 digest of each file that a runtime vendors, which spec manifest writes
├── spec.go            the package spec, which embeds the JSON files and decodes them
├── names.gen.go       a Go constant for each entry
└── mutate-spec/       the vendored copy of mutate-spec
```

`go generate ./spec` renders each YAML file into its JSON file and writes `names.gen.go`, and `spec manifest` writes `manifest.json`. The gate of dokimi runs both and fails when they change a file. The rendered JSON files are committed, because the runtimes read them and do not parse YAML. The manifest has entries for the files of the conformance spec alone, so a runtime does not vendor the copy of mutate-spec.

### Identifiers

Every entry has an identifier of lowercase ASCII words joined by hyphens:

```text
id := word ("-" word)*
word := [a-z][a-z0-9]*
```

An identifier is unique across the files of the spec, because each runtime declares the constants of every kind of entry in one namespace. An identifier does not start with the article `a`, `an` or `the`, so that `read-returns-the-write-before-it` is `ReadReturnsTheWriteBeforeIt` in Go. The identifier of a removed entry is never used again. Each runtime writes an identifier in the casing of its language: `zero-beside-error` is `ZeroBesideError` in Go, `zeroBesideError` in TypeScript and `zero_beside_error` in Rust and Python.

### Laws

A law is a statement about a subject that is true or false. The service counts coverage per law:

```yaml
laws:
  no-value-beside-an-error:
    statement: A call that fails returns no value.
    depth: nothing-beyond-a-caller
    distinct-from: []
    since: 1.0.0
```

`depth` is the machinery at which a check witnesses the law in full. A check whose claim does not list that machinery witnesses the law at its shallow end, and the runtime marks the check as shallow in its suite record. `distinct-from` contains each law that can be false while this one is true, with the case that separates them. assert-spec calls a body that passes for every generated input a property, so the conformance spec calls its unit of coverage a law.

### Claims

A claim is a statement that a classification licenses, with the law that its check witnesses. Each claim has licences, gates, a law, its machineries, a probe with its inputs and a planted defect:

```yaml
claims:
  zero-on-error:
    generator: suite
    licences:
      - shape: reader
      - shape: aggregator
      - shape: answering-writer
      - shape: stream-consumer
    unless: default-on-error
    gates: [failure, returns-beside-error, one-value, value-samples]
    law: no-value-beside-an-error
    machineries: [nothing-beyond-a-caller]
    statement: "{callable} returns the zero value beside an error"
    probe: zero-beside-error
    inputs:
      call: {from: callable}
      zero: {from: zero}
    defect:
      form: value-beside
      clause: "a {declaration} whose {callable} returns a value beside an error"
      reason: value-beside-error
    since: 1.0.0
```

The identifier of a claim is the segment of each check of the claim. `generator` is the generator that binds the claim: `suite`, `enum`, `sentinel` or `sim`. The derivation of the binary reads the claims. A runtime does not read a claim at run time. The constant generator of each runtime reads only the identifiers of the claims.

A claim of `enum`, `sentinel` or `sim` has no licences. The subject of a claim of `enum` is an enumeration, the subject of a claim of `sentinel` is the errors of a package, and the subject of a claim of `sim` is a machine interface under the directive `sim`. Such a claim binds to each subject of its generator and takes its inputs from the subject:

```yaml
claims:
  distinct-errors:
    generator: sentinel
    gates: [two-error-values]
    law: error-values-are-distinct
    machineries: [nothing-beyond-a-caller]
    statement: "no error value of {declaration} matches another"
    probe: errors-differ
    inputs:
      values: {from: error-values}
    defect:
      form: repeated-error
      clause: "a {declaration} whose second error value is its first"
      reason: shared-identity
    since: 1.0.0
```

A claim of `suite` binds to a callable, or to the whole declaration when it sets `scope: declaration`. A claim on the whole declaration, such as agreement with a reference over generated sequences, has no callable in its identity.

| Field | Meaning |
|---|---|
| `licences[]` | One classification: `shape`, `mixin` or `contract` with its `role`, or `callable` for every callable and `context` for every callable with a parameter in the role context |
| `licences[].needs-param` | A param of the classification that the directive must set, such as `attempts` of the mixin `retry-succeeds` |
| `licences[].with` | A second classification that the same callable must have. A licence with `with` licenses a composite |
| `licences[].conflicts` | Classifications whose presence on the callable withdraws the licence |
| `unless` | A mixin whose presence on the callable withdraws every licence of the claim. For a claim of `enum` or `sentinel`, a param of the directive and the value that withdraws the claim, such as `{param: prefix, value: "off"}` of the directive `sentinel` |
| `scope` | `callable`, the default, or `declaration` |
| `machineries` | The machineries that the check needs. The first is the deepest |
| `partners` | The other callables that the probe calls, each found by a mode |
| `inputs.<name>.from` | `callable`, `partner`, `sample`, `zero`, `param` or `clock` for a claim of `suite`. `members`, `name`, `parse` or `param` for a claim of `enum`. `error-values`, `error-types`, `causes` or `param` for a claim of `sentinel`, where `param` reads a param of the directive, such as `no-overlap-with` |
| `defect.form` | The identifier of a form in `defects.yaml` |
| `entails` | The claims on the same callable whose laws the law of this claim implies |

### Partners

A claim whose probe calls a second callable names it as a partner. Each partner has a mode, and the derivation implements each mode once:

```yaml
claims:
  read-after-write:
    generator: suite
    licences:
      - mixin: read-after-write
    gates: [cancellation, one-value]
    partners:
      write: {mode: named, param: write}
    law: read-returns-the-write-before-it
    machineries: [nothing-beyond-a-caller]
    statement: "{callable} returns what {write} stored before it"
    probe: write-then-read
    inputs:
      read: {from: callable}
      write: {from: partner, partner: write}
    defect:
      form: forgets-writes
      clause: "a {declaration} whose {callable} does not return what {write} stored"
      reason: unobserved-write
    since: 1.0.0
```

| Mode | The partner |
|---|---|
| `named` | The callable that a param of the licence names, such as `write` of the mixin `read-after-write`. The licensed callable reads, and the partner writes |
| `reader` | The callable that a param of the licence names, with the halves swapped: the licensed callable writes, and the partner reads, such as the role `reader` of the contract `persister` |
| `derived` | The one other callable that takes, as a parameter, the type of the result field that a param of the licence names, such as `version` of the mixin `monotonic-reads` |
| `role` | The callable with the named role of the same contract, such as `backing` of the contract `cache` |
| `interposed` | A `named` partner that the probe calls twice, with reads before and after the second call. The derivation also binds a field of the written value that the second call changes |
| `rewriter` | A `reader` partner of a licensed callable that the probe calls twice, with the same field as `interposed` |
| `establishing` | The one other callable of the declaration that has the classification of the entry, such as `{mode: establishing, shape: writer}`, and whose written value the licensed callable reads by its arguments. An example is the writer that puts the subject of a reader into a state |

A mode binds exactly one callable or the claim refuses. No match refuses with the gate `partner-found`, and two matches refuse with the gate `one-partner`. Both remedies name the param that selects the partner. A param of the directive on the licensed callable that names a callable takes precedence over every mode.

### Defect forms

`defects.yaml` lists the forms of the planted defects. A form names the callable that it changes, what the change does, and the shortest sequence of calls that shows it:

```yaml
forms:
  forgets-writes:
    changes: partner:write
    does: The partner returns success and stores nothing.
    min-steps: 2
    since: 1.0.0
  value-beside:
    changes: callable
    does: A failed call returns the second sample of the value instead of the zero value.
    min-steps: 1
    since: 1.0.0
```

The template tree of each language writes a form as a type that wraps a fresh subject of the application and overrides the callable that the form changes, as ADR-0015 decides. Falsification of a property check generates sequences of at least `min-steps` calls.

### Depth pairs

`entails` states that the law of a claim implies the law of another claim on the same callable. Falsification checks each pair in both directions:

- **Down.** The deeper check fails on the planted defect of the shallower check. Otherwise the deeper law does not contain the shallower one.
- **Up.** The shallower check passes on the planted defect of the deeper check. Otherwise the two checks are equivalent, and the shallower one is redundant.

A failure of either direction fails the falsification of the deeper check with the identities of both checks. The falsification test of a suite checks the pairs of the checks of the consumer. An application checks the pairs of the derived checks under the option `isolate`, which the corpus passes for each control. A test of the package `spec` rejects an `entails` that names a claim without a licence in common, and a cycle of `entails`.

### Gates

A gate is a condition on the facts of the subject of a claim: a callable for `suite`, an enumeration for `enum`, the errors of a package for `sentinel` and a machine interface for `sim`. The derivation implements each gate once for every language. Each gate in the spec has its condition, the reason of a refusal and the remedy:

```yaml
gates:
  value-samples:
    passes: eidos derives two distinct samples of each value result.
    why: "eidos derives no sample of a result of {callable}"
    remedy: "author a sample of the result type with the directive sample"
    since: 1.0.0
```

A licensed claim whose gate fails is a refusal. The refusal contains the `why` and the `remedy` of the gate. A claim without licences whose gate fails has no check and no refusal, because no classification licenses it.

### Probes

A probe is the procedure of a check. It calls the subject through its inputs and asserts through assertions of assert-spec. Each runtime implements every probe in its own language:

```yaml
probes:
  zero-beside-error:
    summary: A call that fails returns the zero value of its result.
    inputs:
      - name: call
        kind: call
      - name: zero
        kind: value
    reasons: [value-beside-error]
    vacuous: No call fails, with the sample inputs or with a cancelled context.
    since: 1.0.0
```

`reasons` lists the reason codes that the probe reports. Before each assertion the probe sets the reason code of a failure of that assertion on its handle, so one assertion can report two reasons at two places of a probe. The handle fails a check whose probe sets a reason that the entry of the probe does not list. `vacuous` is the condition under which the probe reports that it did not exercise its claim.

Each input has a kind from a closed list. The template tree of each language passes each kind in one form:

| Input kind | What the generated check passes |
|---|---|
| `call` | A function that calls the subject with bound arguments and returns the results of the subject. It takes the context of the check where the language has one |
| `value` | A sample value of a type, such as an input that the claim names |
| `error` | An error value or an error type that the probe matches by identity |
| `duration` | A length of time, such as the lifetime of an entry |
| `count` | A whole number, such as the attempts of a retry |
| `projection` | A function that returns one part of a value, such as its version |
| `scope` | A function that runs a body inside a scope of the subject, such as a transaction |
| `hook` | A function that installs a hook through the subject |
| `predicate` | A function that returns whether the subject is in a state, such as leadership |
| `clock` | The controlled clock of the run, which the probe advances and waits on |
| `values` | Named values that the check covers as one set, such as the members of an enumeration, the error values of a package, or a function of each error type that makes an error with a given cause |

Every probe also receives the statement of its claim. It passes the statement as the contract of each assertion, so every failure record and every call record of the probe contains the statement. The contract of a probe over `values` is the statement followed by the names of the values at fault.

### Reason codes

A reason code is the identifier of a violation that a probe reports:

```yaml
reasons:
  value-beside-error:
    meaning: The call failed and returned a value other than the zero value.
    since: 1.0.0
```

A planted defect fails its check with the reason code of its claim. Falsification compares that code with the code that the probe reports.

### Machineries

A machinery is what a run supplies before a check exercises its claim. The runtime supplies some machineries to every check. The consumer supplies the others through a door of the application, and the declaration states each of those with the param `machineries` of the directive `suite`:

```go
// Store maps keys to items.
//
//dokimi:suite machineries=reference,faults
type Store interface { /* ... */ }
```

| Machinery | What the run supplies | Supplied by | Param |
|---|---|---|---|
| `nothing-beyond-a-caller` | A fixed sequence of calls on one subject, with derived inputs | The runtime | — |
| `clock-the-run-advances` | The controlled clock of the assert library, which the check advances | The runtime | — |
| `generated-inputs` | Inputs of the calls of a check that the property engine generates and shrinks | The runtime | — |
| `benchmark-run` | A run of the benchmarks of the test runner, such as `go test -bench` or `vitest bench` | The runtime | — |
| `generated-sequences` | Sequences of calls that the property engine produces and shrinks | The runtime | — |
| `concurrent-callers` | Calls from several callers at once, recorded as a history | The runtime | — |
| `call-that-fails-on-demand` | A subject whose dependencies fail when the check asks for it | The consumer | `faults` |
| `reference-implementation` | A second implementation that gives the expected results | The consumer | `reference` |
| `second-subject` | Another implementation that receives the same calls | The consumer | `second` |
| `killable-process-or-breakable-medium` | A process that the run can kill, and storage that the run can break | The consumer | `medium` |
| `simulated-system` | Every machine of a system on one goroutine, with virtual clocks, faults at the network, the dependencies and the processes, and a trace of every step | The generator `sim` | — |

A licensed claim that needs a machinery of the consumer that the declaration does not state has no check and no refusal. The generator lists its classification under the unclaimed classifications with the stance `needs` and a note that names the param. A declaration that states a machinery gets the check, and an application without the door of the machinery fails the check with the name of the option, as every door of conform does.

### Stances

A stance is the reason why a classification has no check:

| Stance | Meaning | What closes it |
|---|---|---|
| `silent` | The classification implies nothing that a check could falsify, such as a call that returns nothing | Nothing |
| `needs` | The claim needs a machinery that the run does not supply, such as a reference implementation that the declaration does not state | The param `machineries` on the declaration, or a runtime that supplies the machinery |
| `declines` | dokimi could build the check and rules it out because of its cost | A cheaper check |
| `blocked` | The check would read a part of the declaration that no generator can read, such as a state graph that the catalog keeps as text | A change of the vocabulary of the shape catalog |
| `gap` | The classification implies something that no claim witnesses | A new claim in the spec |

A refusal is not a stance. It is a licensed claim whose gate is false for a callable, with the reason and the remedy of the gate.

`stances.yaml` also contains the stance of each classification without a claim:

```yaml
unclaimed:
  mixin:concurrent:
    stance: needs
    note: The claims of the mixin need concurrent callers.
```

A test of the package `spec` walks every classification of the shape catalog of eidos. Each classification must be the licence of a claim, a role of a contract that a claim of the same contract binds as a partner, or an entry in `unclaimed`. The role `backing` of the contract `cache` is such a partner: the claim `cache` calls it to compare the cache with its backing.

### Check identity

A check has two forms of identity. The declared form uses the names that the source declares:

```text
identity := declaration ["." callable] "/" segment ["#" case]
```

`declaration` and `callable` are the names of the declaration and the callable as the source declares them, such as `Store.Read` in Go and `Store.read` in TypeScript. A check of a claim with the scope `declaration`, such as `Store/agrees-with-the-reference`, has no callable. A check of an enumeration or of the errors of a package has no callable either. Its declaration is the enumeration or the name of the package. The name of a TypeScript module is the name of its file without the extensions. `segment` is the identifier of the claim, and `case` separates two bindings of one claim on one callable.

To join checks across languages, the key form splits each name into words with `naming.Words` of eidos. It writes the words in lower case and joins them with hyphens. `naming.Words` splits at a separator, between a lower-case and an upper-case letter, before the last upper-case letter of a run that a lower-case letter follows, and before an upper-case letter that follows a digit. `Store.ReadAll/zero-on-error` in Go, `Store.readAll/zero-on-error` in TypeScript and `Store.read_all/zero-on-error` in Rust all have the key `store.read-all/zero-on-error`. `HTTPServer` in Go, `httpServer` in TypeScript and `http_server` in Rust all have the key `http-server`. Two callables of one declaration with the same key fail generation.

The test runner of each language runs each check as a test of its own, named by the identity of the check. The call records of a check belong to that test. The field `tests` of a killed mutant in the run record of mutate-spec names the tests that failed on it, so a mutant maps to the checks that detected it. A run that stops at the first failure names one of those checks.

### The suite record

A suite has two applications, and each writes a suite record:

- The falsification test applies every check of the consumer and every ceiling to its planted defect.
- The application to a subject first falsifies each derived check on its planted defect over a fresh subject, and then applies every check to the subject. It writes a falsification record and a record of the subject.

The corpus applies each suite to its controls through the application to a subject, so a control is a subject whose record the corpus judges. With the option `isolate`, the application also runs every other derived check on each planted defect and checks each depth pair. The application to the subject of a consumer does not isolate its defects. A defect of that subject would fail the other checks on every planted defect that wraps it.

A runtime writes the suite record into the artifact of the test runner beside the call records of its assert library, in the form that assert-spec gives the call records of the language. In Go the key is `dokimi.suite.<seq>`, where `seq` numbers the suite records of one test from 1. A record that does not fit one event is split as assert-spec splits a call record: over consecutive events with the same key, whose values join into the record. In TypeScript the record is an element of the array under the key `dokimi.suite` in the meta of the test. The entry `records` of the overlay of each language in assert-spec describes where that artifact keeps the records of a test. `schemas/suite.schema.json` describes the format:

```json
{
  "spec": "1.0.0",
  "suite": {"language": "go", "package": "example.com/store/storetest", "declaration": "Store"},
  "generator": {"name": "dokimi", "version": "v0.1.0"},
  "runtime": {"name": "go.dokimi.dev/conform", "version": "v0.1.0"},
  "seed": "0x5f3a1c9e27d4b860",
  "application": "subject",
  "subject": "InMemory",
  "checks": [
    {
      "identity": "Store.Read/zero-on-error",
      "key": "store.read/zero-on-error",
      "claim": "zero-on-error",
      "law": "no-value-beside-an-error",
      "machineries": ["nothing-beyond-a-caller"],
      "shallow": false,
      "result": "pass"
    },
    {
      "identity": "Store.Read/keeps-the-item",
      "key": "store.read/keeps-the-item",
      "authored": true,
      "result": "pass"
    }
  ],
  "refusals": [
    {
      "claim": "zero-on-error",
      "callable": "Store.Scan",
      "why": "eidos derives no sample of the parameter visit of Scan",
      "remedy": "author a sample of the parameter visit with the directive sample"
    }
  ],
  "unclaimed": [
    {
      "classification": "mixin:time-aware",
      "callable": "Store.Stale",
      "stance": "needs",
      "note": "the claim needs a clock that the run advances"
    }
  ]
}
```

| Field | Values |
|---|---|
| `spec` | The version of the conformance spec that the runtime implements |
| `application` | `falsification` or `subject` |
| `subject` | The name of the subject that the consumer passed to the application, for the application `subject` |
| `seed` | The seed of the run, as `0x` and 16 hexadecimal digits |
| `ended` | The fault that ended the application before its last check, absent when every check ran |
| `checks[].authored` | True for a check of the directive `check`. Such a check has no `claim`, `law`, `machineries` or `reason` |
| `checks[].result` | `pass`, `fail`, `vacuous`, `excused`, `skipped` or `unchecked` for a subject. `proven`, `passed`, `wrong-reason`, `vacuous`, `unchecked`, `not-isolated`, `not-deeper` or `equivalent` for falsification |
| `checks[].other` | The identity of the other check of a result `not-isolated`, `not-deeper` or `equivalent` |
| `checks[].why` | The reason of `unchecked`, `skipped` or `excused` |
| `checks[].cases` | The number of valid cases of a property check |
| `checks[].seed` | The seed of the run of a property check, as `0x` and 16 hexadecimal digits |
| `checks[].counterexample` | The values of the smallest failing case by label, as the typed literals of assert-spec |
| `checks[].reason` | The reason code that the probe reported |
| `checks[].failure` | The failure record of assert-spec for the assertion that failed, on `fail` and on `wrong-reason` |
| `checks[].shallow` | True when the machineries of the claim do not include the depth of its law |
| `unclaimed[].classification` | `shape:<name>`, `mixin:<name>` or `contract:<name>/<role>`, with the names of the shape catalog |
| `unclaimed[].stance` | A stance of the spec |

The results of falsification:

| Result | The check on its planted defect |
|---|---|
| `proven` | Failed with the reason code of its claim, or, for an authored check, at the assertion whose contract is its statement |
| `passed` | Passed |
| `wrong-reason` | Failed with another reason code, or at an assertion with another contract |
| `vacuous` | Did not exercise its claim |
| `not-isolated` | Failed with the reason code of its claim, and another check that does not share the defect also failed on it |
| `not-deeper` | Failed with the reason code of its claim. The check is the deeper check of an edge, and it passed on the planted defect of the shallower check |
| `equivalent` | Failed with the reason code of its claim. The check is the deeper check of an edge, and the shallower check also failed on its planted defect |

On every result other than `proven` the falsification test fails with a message that contains the identity of the check, the reason code of its claim and the reason code that the probe reported. A result `not-isolated`, `not-deeper` or `equivalent` also names the other check.

An allocation ceiling that the build or the run does not check reports `unchecked`, as RFC-0006 states, on a subject and in falsification.

The results of a subject:

| Result | The check on the subject |
|---|---|
| `pass` | Passed and engaged its claim |
| `fail` | Failed |
| `vacuous` | Passed and never engaged its claim. The application fails the check |
| `excused` | An excusal of the consumer covers it |
| `skipped` | Skipped because the environment lacks what the check needs, such as direct I/O, and no excusal covers it |
| `unchecked` | An allocation ceiling that the build or the run does not check |

A reader skips each field that its version of the schema lacks, and it shows an unknown member of an enumeration as unknown instead of mapping it to another member.

### The bundle and the seed

The bundle of a run is the set of artifacts of its test runner, with their call records and suite records, and the trace files of the simulations. CI uploads the bundle to the service beside the outputs of the other tools of the run, such as coverage and mutation records.

The seed of a run comes from the environment variable `DOKIMI_SEED`, as a decimal or a `0x` hexadecimal 64-bit number. When the variable is unset, the runtime takes a seed from the clock and prints it with the first failure. A probe that runs the property engine passes the seed of the run to the option `seed` of the engine. Every suite record contains the seed, so a run repeats with `DOKIMI_SEED` set to it.

Falsification does not take the seed of the run for a property check. It runs the property hermetically with a seed from the key of the check, as RFC-0006 states, over sequences of at least the `min-steps` of its defect form. It turns the store and the shrinking of the engine off. A proof therefore gives the same result on every run, and the corpus pins that the proof of each defect form fails at its seed.

### The contract between a generated file and its runtime

A generated file and its runtime agree on a compatibility function, the names of the probes, the constants of the spec and the reason codes:

1. **A compatibility function.** Each package of a runtime that generated files import exports a function whose name contains the major version of its module, such as `conform.CompatV1` and `stub.CompatV1` in Go. Every generated file refers to the function of the major version that the generator targets. A major version of a runtime renames the function. A mismatch fails the build of the consumer in Go, Rust, Java, Kotlin and C#, and the type check or the module load in TypeScript. In Python and PHP it fails the import of the generated file.
2. **Probe names.** A generated check calls each probe by its identifier in the casing of the language. With a runtime that implements an older spec without the probe, the build or the import fails at that call.
3. **Constants.** Each runtime generates the constants of the entries of the spec from its vendored copy. Generated files refer to those constants, so a file that refers to an entry that the runtime lacks fails the build or the import.
4. **Reason codes.** A probe asserts through the assert library of its language. It passes the library a test handle of the runtime instead of the test itself, and it sets the reason code of each assertion on the handle before the assertion. The handle receives the failure record of each failed assertion and records the reason code that was set. Falsification reads the recorded code and never the text of the failure.

The handle implements the interface that its assert library calls on a test, such as `assert.TB` in Go and `Seat` in TypeScript. In falsification the handle wraps the recorder of the assert library. The recorder keeps the failure records of the check without failing the falsification test.

### The API of a runtime

Each runtime implements the same parts in the conventions of its language:

| Part | Go: `go.dokimi.dev/conform` | TypeScript: `@dokimi/runtime` |
|---|---|---|
| The constants of the spec | The package `spec`, such as `spec.ErrorsDiffer` and `spec.SharedIdentity` | The module `spec`, such as `probes.errorsDiffer` and `reasons.sharedIdentity` |
| A check | `conform.Check[S]`: the identity and the key, the statement, the claim with its law and machinery, the reason of its defect, a function that applies the check through its handle, and a function that returns its planted defect as a subject of the type `S`. Its methods `Run` and `Falsify` apply the check on any seat to a subject and to its planted defect | `Check<S>`, with the same fields and functions |
| The checks of one declaration | `conform.Suite[S]`: the language, the package, the declaration, the checks, the refusals and the unclaimed classifications | `Suite<S>` |
| The test handle | `*conform.Handle[S]`, which implements `assert.TB`, `assert.Reporter`, `assert.Clocked` and `assert.Observer`, returns the seat of its check from `Unwrap`, and builds the subjects of its check. `Fresh` returns a fresh subject of the factory also under falsification, for a planted defect to wrap. `Expect` sets the reason code of the assertions that a probe calls next | `Handle<S>`, which implements `Seat` with `observeCall`, with `fresh` and `expect` |
| The name of a callable | `conform.Callable[S]`, whose method `Check` returns the identity `conform.ID[S]` of a check of the callable | `Callable<S>`, with the function `check` |
| A probe | A function of the package `probe`, such as `probe.ErrorsDiffer(h, values)` | A function of the module `probe`, such as `errorsDiffer(h, values)` |
| A named value | `conform.Named[T]` | `Named<T>` |
| The application to a subject | `conform.Run(t, suite, newSubject, opts...)`, where `newSubject` builds a fresh subject from the handle of each check. The generated `RunStoreSuite` and `RunStoreSuiteFromHandle` of RFC-0003 call it | `run(suite, newSubject, options)`, which the generated `runStoreSuite` and `runStoreSuiteFromHandle` call |
| The options of an application | `conform.With`, `conform.Excuse`, `conform.Reference`, `conform.Second`, `conform.Restartable`, `conform.Crashable`, `conform.Cases`, `conform.Using`, `conform.Ceilings` and `conform.Isolate`. The RFC of each machinery of the consumer adds the option of its door | `Options<S>`, with `isolate` |
| The falsification | `conform.Falsify(t, suite)` | `falsify(suite)` |

The module `go.dokimi.dev/conform` contains the packages `conform`, `spec` and `probe` of the Go runtime, as RFC-0005 states. `Run` and `Falsify` run each check as a test of its own, named by the identity of the check. Each check gets a handle of its own, and the probe reads the statement of the check from it. `Run` fails a derived check whose planted defect does not fail it with the reason of its claim, and a check that fails on its subject. Under `Isolate`, `Run` also fails a derived check whose planted defect fails another check, and an edge that fails in either direction. `Falsify` fails a check of the consumer whose planted defect does not fail it at its statement, and a ceiling that its defect does not break. Each writes the suite record of its application. The subject of a suite is the implementation that the factory of the consumer returns. The subject of the checks of errors is a value of a type that the generated test declares, with the errors of the package.

### The package spec

The package `spec` of the module `go.dokimi.dev/dokimi` decodes every file of the spec for the derivation. The package `spec` of each runtime, such as `go.dokimi.dev/conform/spec`, declares the identifier types and the constants alone, from its vendored copy.

```go
// Package spec contains the conformance spec of dokimi: the laws, claims,
// gates, probes, reason codes, machineries and stances that the binary
// shares with the runtimes.
package spec

// LawID identifies a law. names.gen.go declares one constant for each.
type LawID string

// ClaimID identifies a claim. It is the segment of each check of the claim.
type ClaimID string

// GateID identifies a gate.
type GateID string

// ProbeID identifies a probe.
type ProbeID string

// ReasonID identifies a reason code.
type ReasonID string

// MachineryID identifies a machinery.
type MachineryID string

// StanceID identifies a stance.
type StanceID string

// DefectFormID identifies a form of a planted defect.
type DefectFormID string

// ReferenceFormID identifies a reference form of a verified fake.
type ReferenceFormID string

// InputKind is the kind of an input of a probe.
type InputKind string

// Input is one input of a probe.
type Input struct {
	Name string
	Kind InputKind
}

// Probe is one probe of the spec.
type Probe struct {
	ID      ProbeID
	Summary string
	Inputs  []Input

	// Reasons lists the reason codes that the probe reports.
	Reasons []ReasonID

	// Vacuous is the condition under which the probe does not exercise its
	// claim.
	Vacuous string

	Since string
}

// Probes returns every probe in the order of probes.json.
func Probes() []Probe

// ProbeOf returns the probe with the identifier id, and false when the spec
// has no such probe.
func ProbeOf(id ProbeID) (Probe, bool)

// Version returns the version of the spec from the file VERSION.
func Version() string
```

The package declares the same pair of functions for laws, claims, gates, reason codes, defect forms, reference forms, machineries and stances. It decodes the embedded files on the first call. A test of the package validates each embedded file against its schema and resolves these references:

- the law, the gates, the probe, the reason code, the defect form and the machineries of each claim
- the reason code of each claim, against the reasons of its probe
- each classification of a claim and each param that a licence or a partner names, against the shape catalog of eidos
- each input of a claim, against the inputs of its probe
- each claim of `entails`, against the licences of the claim, with no cycle
- the depth of each law and each law of `distinct-from`
- every classification of the shape catalog, against the licences and the partners of the claims and the entries of `unclaimed`

### Vendoring

The program `spec` of stealth-scale/spec vendors the conformance spec, as it vendors assert-spec into each assert library. ergon installs it at the release that the section `spec` of `.ergon.yaml` pins, and renders its targets and jobs from the roles of the repository.

dokimi is the repository of a definition. Its `.spec.yaml` declares the files that a runtime vendors:

```yaml
definition:
  repo: dokimasia/dokimi
  root: spec
  manifest: manifest.json
  languages: [go, typescript]
  files: [VERSION, "*.json", "schemas/*.json"]
  implementations:
    - {repo: dokimasia/dokimi-go, languages: [go]}
```

The section lists one entry for each runtime, with the languages that it implements.

`spec manifest` writes `spec/manifest.json`, with the digest of each file, the version and the digest of the whole set. The gate of dokimi runs `spec manifest --check`, which fails when the manifest differs from the files.

Each runtime is an implementation. Its `.spec.yaml` follows the releases of dokimi, because the tests that it runs come from a released binary:

```yaml
implements:
  - definition: dokimasia/dokimi
    track: release
    path: conformance/spec
    languages: [go]
    test: [make, test-go]
```

| Command | Where | What it does |
|---|---|---|
| `spec sync` | A runtime | Writes the copy from the newest release of dokimi, with a `manifest.json` that states the commit of the copy |
| `spec check` | The gate of a runtime | Fails when a file of the copy differs from its manifest |
| `spec check --remote` | The job `spec` of the CI of a runtime | Fails when the commit of the copy is not a release of dokimi, and reports a newer release |
| `spec propose` | The job `propose-spec` of the nightly workflow of a runtime | Opens a pull request that updates the copy to the newest release |
| `spec try` | dokimi | Writes the working tree of the spec into the checkout of each runtime beside dokimi and runs its tests |

mutate-spec is not a definition of the program `spec`, so dokimi vendors it with the scripts `tools/spec-sync.sh` and `tools/spec-check.sh` of mutate-spec, as mutate-go does. Each runtime requires a release of its assert library that has every assertion that its probes call.

### Versioning

`VERSION` changes independently of the version of the Go module. A release of dokimi that leaves the spec unchanged leaves `VERSION` unchanged. Each runtime declares the version of the spec that it implements.

| Change | Version |
|---|---|
| A law, claim, gate, probe, reason code, machinery or stance is added | minor |
| The statement, summary or meaning of an entry is reworded | patch |
| The inputs, the assertions or the meaning of a probe change | major |
| The licences, the gates or the law of a claim change | major |
| The meaning of a law, a machinery or a stance changes | major |
| A field of the suite record is removed or changes meaning | major |
| A reason code is removed or receives another meaning | never |

### Invariants

| Invariant | Enforced by |
|---|---|
| Every rendered file matches its YAML file, its schema and the digest of the manifest | The generate step of the gate of dokimi, and a test of the package spec |
| Every reference of an entry resolves in the spec and in the shape catalog of eidos | A test of the package spec |
| A probe reports only the reason codes that its entry lists | The test handle of each runtime |
| Every classification of the shape catalog has a claim or a stance | A test of the package spec |
| A vendored copy matches its manifest, and its commit is a release of dokimi | `spec check` in the gate of each runtime, and `spec check --remote` in its CI |
| Falsification proves a derived check by its reason code and an authored check by the contract of its failure record, and never by the text of a failure | The test handle of each runtime, and the corpus |
| A proof of a property check gives the same result on every run | The seed from the key of the check, and the corpus |
| A generated file and its runtime agree on their major version | The compatibility function |
| Every suite record contains its seed | The runtime of each language |

### Failure handling

| Failure | Result | Recovery |
|---|---|---|
| A YAML file has an entry that its schema rejects | The test of the package spec fails | Correct the entry |
| An entry refers to a law, gate, probe, reason code, machinery or classification that does not exist | The test of the package spec fails and reports the reference | Add the entry, or correct the reference |
| A probe calls an assertion that the assert library of its runtime lacks | The build of the runtime fails | Add the assertion through the RFC process of assert-spec, and require the release of the assert library that implements it |
| A rendered file, `names.gen.go` or the manifest is stale | The generate step fails the gate of dokimi | Run `go generate ./spec` and commit the result |
| A vendored file differs from its digest | `spec check` fails the gate of the runtime | Run `spec sync --locked` |
| A newer release of dokimi has another spec | `spec propose` opens a pull request in the runtime | Implement the new entries on the branch of the pull request |
| A generated file calls a probe that the runtime lacks | The build or the import of the consumer fails | Upgrade the runtime |
| A generated file and its runtime differ in their major version | The compatibility function fails the build or the import | Upgrade the older side |
| A probe asserts without setting a reason code first | The handle records the failure without a code, and falsification reports `wrong-reason` | Fix the probe |
| A probe sets a reason code that its entry does not list | The handle fails the check with the probe and the code | Fix the probe or its entry |

## Alternatives considered

### Conformance in assert-spec

assert-spec gains the probes, the reason codes and the suite record, and each assert library implements every probe.

**Why not:** assert-spec defines assertions over values and callables. A shape is a classification of eidos, and the shape catalog is not part of assert-spec. The laws of a shape change with the shape catalog and with the claims of dokimi, and neither belongs to assert-spec. Each assert library would implement probes that only the generated tests of dokimi call. A probe composes assertions of assert-spec, so assert-spec does not change when the conformance spec adds a probe.

### Falsification through rejects

The falsification test runs each check through the assertion `rejects` of assert-spec. It accepts the failure when the first failure record contains the assertion of the claim.

**Why not:** many probes fail through `equal` or `err-is`, so the id of an assertion does not separate one violation from another. A value beside an error and a value for a miss both fail `equal`. Each violation of a law has a reason code of its own, and every runtime reports the same code.

### A Go table of claims

The claims are a Go table in the package of the derivation, and the spec contains only the entries that the runtimes read.

**Why not:** the claims define the conformance of each shape. The corpus, the service and the packs of third parties read that definition, and a Go table would make the source of the binary its only statement. A Go table rejects a misspelled name at compile time. The test of the package `spec` rejects the same names before each release.

The partners of a claim are data too. Each partner names one of seven modes, and the derivation implements each mode once. testkit sorted its partners into the same modes and resolved them with about 40 functions. The verdict reverses when a claim needs a partner that no mode can resolve.

### A repository of its own for the spec

A repository dokimi-spec contains the conformance spec, and the binary and each runtime vendor it, as each assert library vendors assert-spec.

**Why not:** every assert library is an equal implementation of its standard, so that standard has a repository of its own. The conformance spec has one producer of checks. A change that adds a claim with a new probe also changes the templates of its defect form and the expected outcomes of the corpus. In the dokimi repository that change is one commit, one review and one release. With dokimi-spec it is a commit there and a sync commit in dokimi, and between the two the generator and its spec disagree.

The verdict reverses when a second program generates calls to the probes.

### Runtimes vendor the main branch

A runtime vendors the spec of the main branch of dokimi, as the assert libraries vendor the main branch of assert-spec.

**Why not:** a consumer runs tests that a released binary generated. The main branch can contain a change of meaning that no release contains yet, and a runtime that vendors it implements a spec that the binary of the consumer does not use. A release tag fixes the spec that a binary and a runtime share.

### Failures matched by their text

A planted defect passes when the failure of its check contains a stated piece of text, as in testkit.

**Why not:** the text of a failure is written for a person in the language of the runtime. A TypeScript probe that had to print the same sentence as a Go probe would turn the wording of one runtime into a contract for all runtimes. A reason code is the same in every language, and each runtime words the text for its own readers. An authored check exists in one language, and its statement is the contract of its decisive assertion. The contract is a field of the failure record, not the text of the failure.

### One reason code per assertion of a probe

The entry of a probe maps each assertion that it calls to one reason code, and the handle reads the code from the assertion of the failure record.

**Why not:** a probe can compare several results through `equal`, each for its own reason. One reason per assertion would split each such probe into several probes that repeat one procedure. A reason set before each assertion keeps one probe per procedure, and the list in the entry still bounds the codes that it reports.

### A report file beside the artifact of the test runner

The runtime writes a JSON report of its own into a directory, apart from the artifact of the test runner.

**Why not:** the service reads the call records from the artifact of the test runner. A second file per suite needs a second upload path and a second parser in the service. A suite record in the same artifact joins the call records of a suite with its claims without another upload.

## Drawbacks

- **A sync commit for each new probe.** Each runtime takes a change of the spec through a release of dokimi and a sync commit of its own. With eight runtimes, an added probe means eight sync commits and eight implementations. A new claim that calls an existing probe does not need a change in a runtime.
- **A constant generator in each runtime.** Each runtime turns the vendored JSON into the constants of its language with a generator of its own.
- **Every probe in every runtime.** The probes are host code, so each probe exists once per runtime, and the corpus judges each copy.
- **Reason codes that a run checks.** A test of the package `spec` cannot tell which assertion of a probe reports which reason. The handle checks each reported code against the list of the probe when a check runs, and the corpus runs every probe.
- **A declaration of each machinery of the consumer.** A declaration that wants the checks of a reference, of faults, of a second subject or of a medium states the machinery with the param `machineries`. Without it the claims are unclaimed with the stance `needs`.
- **New assertions wait on assert-spec.** A probe that needs an assertion that assert-spec lacks waits for an RFC of assert-spec, an entry in the overlay of each language and an implementation in each assert library.
- **Two copies of each entry in dokimi.** The YAML file and its rendered JSON file are both committed. The generate step keeps them equal.
- **A YAML library and a JSON Schema validator in the dokimi module.** The renderer reads YAML, and the test of the package `spec` validates each rendered file. The module requires `go.yaml.in/yaml/v3` and `github.com/google/jsonschema-go`, which the binary does not use at run time. The renderer infers each schema from the Go types of the package `spec` with `jsonschema.For`, and the test validates each rendered file with `Resolved.Validate`.

## Open questions

None.

## References

| What | Where |
|---|---|
| The shape catalog of eidos | <https://github.com/dokimasia/eidos/blob/main/docs/architecture/12-shape-catalog.md> |
| assert-spec: the definition, its rendering and its vendoring | <https://github.com/dokimasia/assert-spec/blob/main/README.md> |
| spec, the program that vendors a conformance definition | <https://github.com/stealth-scale/spec> |
| The vendoring scripts of mutate-spec | <https://github.com/dokimasia/mutate-spec/tree/main/tools> |
| The failure record of assert-spec | <https://github.com/dokimasia/assert-spec/blob/main/spec/failure.md> |
| The recording of a run in assert-spec | <https://github.com/dokimasia/assert-spec/blob/main/spec/recording.md> |
| mutate-spec: a mutation run and its record | <https://github.com/dokimasia/mutate-spec/blob/main/docs/rfc/0002-a-mutation-run-and-its-record.md> |
| The request to name every test that kills a mutant | <https://github.com/dokimasia/mutate-spec/issues/36> |
| `naming.Words` of eidos | <https://github.com/dokimasia/eidos/blob/main/eidos-lang/naming/words.go> |
| JSON Schema 2020-12 | <https://json-schema.org/draft/2020-12/json-schema-core> |
