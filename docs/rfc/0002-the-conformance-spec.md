---
rfc: 0002
title: The conformance spec
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Draft
created: 2026-10-19
updated: 2026-10-19
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

assert-spec defines assertions over values and callables, and it defines no shape. Each probe asserts through assertions of assert-spec, and its entry in the conformance spec has the assertion behind each of its reason codes. A run writes the call records of assert-spec and the suite records of the conformance spec into the artifact of its test runner.

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
| `assertions.json` | The assertions of each probe. A test of the package `spec` resolves each one |
| The failure record | The failure of a check in its suite record |
| The call record | The calls of each probe in a recorded run |
| The controlled clock | The machinery `a-clock-the-run-advances` |
| The property engine | The machinery `generated-sequences`, and the fuzz targets of the claims of a suite |
| The machines | The model tier, which runs generated sequences against a reference implementation |
| `linearizable`, `serializable` and `snapshot-isolation` | The machinery `concurrent-callers`, which records a history and checks it |
| The benchmark ceilings and `max-allocs` | The budgets of a callable |
| The overlay of each language | The refusal of a claim in a language whose overlay declares an assertion of its probe absent |

The overlay is the source of truth for what a language lacks. The overlay of TypeScript declares `bench-max-allocs` absent, because V8 does not count allocations per iteration. A claim whose probe asserts `bench-max-allocs` refuses in TypeScript, and the refusal contains the reason and the remedy of the overlay entry.

### Layout

```text
spec/
├── VERSION            the version of the spec, such as 1.0.0
├── laws.yaml          the laws, edited by people
├── claims.yaml        the claims of each generator
├── gates.yaml         the gates
├── probes.yaml        the probes
├── reasons.yaml       the reason codes
├── machineries.yaml   the machineries
├── stances.yaml       the stances, and the stance of each classification without a claim
├── laws.json          rendered from laws.yaml, as each YAML file is rendered
├── …
├── schemas/           a JSON Schema for each rendered file and for the suite record
├── manifest.json      the SHA-256 digest of each file that a runtime vendors
├── spec.go            the package spec, which embeds the JSON files and decodes them
├── names.gen.go       a Go constant for each entry and for each assertion of assert-spec
├── assert-spec/       the vendored copy of assert-spec
└── mutate-spec/       the vendored copy of mutate-spec
```

`go generate ./spec` renders each YAML file into its JSON file and writes `names.gen.go` and `manifest.json`. The gate of dokimi runs the generators and fails when they change a file. The rendered JSON files are committed, because the runtimes read them and do not parse YAML. The manifest has entries for the files of the conformance spec alone, so a runtime does not vendor the copies of assert-spec and mutate-spec.

### Identifiers

Every entry has an identifier of lowercase ASCII words joined by hyphens:

```text
id := word ("-" word)*
word := [a-z][a-z0-9]*
```

An identifier is unique within its file. The identifier of a removed entry is never used again. Each runtime writes an identifier in the casing of its language: `zero-beside-error` is `ZeroBesideError` in Go, `zeroBesideError` in TypeScript and `zero_beside_error` in Rust and Python.

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

`depth` is the machinery at which a check witnesses the law in full. A check whose machinery is shallower witnesses the law at its shallow end, and the runtime marks the check as shallow in its suite record. `distinct-from` contains each law that can be false while this one is true, with the case that separates them. assert-spec calls a body that passes for every generated input a property, so the conformance spec calls its unit of coverage a law.

### Claims

A claim is a statement that a classification licenses, with the law that its check witnesses. Each claim has licences, gates, a law, a machinery, a probe with its inputs and a planted defect:

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
    gates: [failure, one-value, value-samples]
    law: no-value-beside-an-error
    machinery: nothing-beyond-a-caller
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

The identifier of a claim is the segment of each check of the claim. `generator` is the generator that binds the claim: `suite`, `enum` or `sentinel`. The derivation of the binary reads the claims, and no runtime reads them.

A claim of `enum` or `sentinel` has no licences. The subject of a claim of `enum` is an enumeration, and the subject of a claim of `sentinel` is the errors of a package. Such a claim binds to each subject of its generator and takes its inputs from the subject:

```yaml
claims:
  distinct-errors:
    generator: sentinel
    gates: [two-error-values]
    law: error-values-are-distinct
    machinery: nothing-beyond-a-caller
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

### Gates

A gate is a condition on the facts of the subject of a claim: a callable for `suite`, an enumeration for `enum` and the errors of a package for `sentinel`. The derivation implements each gate once for every language. Each gate in the spec has its condition, the reason of a refusal and the remedy:

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
    asserts:
      equal: value-beside-error
    vacuous: No call fails, with the sample inputs or with a cancelled context.
    since: 1.0.0
```

In `asserts`, each assertion that the probe calls has the reason code of its failure. A probe calls each assertion for one reason only, so the failure record of an assertion identifies its reason code. Two reasons that need the same assertion belong to two probes. `vacuous` is the condition under which the probe reports that it did not exercise its claim.

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

A machinery is what a run supplies before a check exercises its claim:

| Machinery | What the run supplies |
|---|---|
| `nothing-beyond-a-caller` | A fixed sequence of calls on one subject, with derived inputs |
| `a-clock-the-run-advances` | The controlled clock of the assert library, which the check advances |
| `a-call-that-fails-on-demand` | A subject whose next call fails when the check asks for it |
| `generated-sequences` | Sequences of calls that the property engine produces and shrinks |
| `a-reference-implementation` | A second implementation that gives the expected results |
| `concurrent-callers` | Calls from several callers at once, recorded as a history |
| `a-second-subject` | Another implementation that receives the same calls |
| `a-killable-process-or-breakable-medium` | A process that the run can kill, and storage that the run can break |

### Stances

A stance is the reason why a classification has no check:

| Stance | Meaning | What closes it |
|---|---|---|
| `silent` | The classification implies nothing that a check could falsify, such as a call that returns nothing | Nothing |
| `needs` | The claim needs a machinery that the run does not supply, such as concurrent callers | A tier or a runtime that supplies the machinery |
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

A test of the package `spec` walks every classification of the shape catalog of eidos. Each classification must be the licence of a claim or have an entry in `unclaimed`.

### Check identity

A check has two forms of identity. The declared form uses the names that the source declares:

```text
identity := declaration ["." callable] "/" segment ["#" case]
```

`declaration` and `callable` are the names of the declaration and the callable as the source declares them, such as `Store.Read` in Go and `Store.read` in TypeScript. A check of an enumeration or of the errors of a package has no callable. Its declaration is the enumeration or the name of the package. The name of a TypeScript module is the name of its file without the extensions. `segment` is the identifier of the claim, and `case` separates two bindings of one claim on one callable.

To join checks across languages, the key form splits each name into words at an underscore, at a hyphen and before an upper-case letter that follows a lower-case letter or a digit. It writes the words in lower case and joins them with hyphens. `Store.ReadAll/zero-on-error` in Go, `Store.readAll/zero-on-error` in TypeScript and `Store.read_all/zero-on-error` in Rust all have the key `store.read-all/zero-on-error`. Two callables of one declaration with the same key fail generation.

The test runner of each language runs each check as a test of its own, named by the identity of the check. The call records of a check belong to that test. The run record of mutate-spec contains the names of the tests that detected each mutant, so a mutant maps to the checks that detected it.

### The suite record

A suite has three applications, and each application writes one suite record:

- The falsification test applies every check to its planted defect.
- `Run` applies every check to a subject of the consumer.
- `Control` applies every check to the control of the consumer.

A runtime writes the suite record into the artifact of the test runner beside the call records of its assert library, under the key `dokimi.suite`. The entry `records` of the overlay of each language in assert-spec describes where that artifact keeps the records of a test. `schemas/suite.schema.json` describes the format:

```json
{
  "spec": "1.0.0",
  "suite": {"language": "go", "package": "example.com/store/storetest", "declaration": "Store"},
  "generator": {"name": "dokimi", "version": "v0.1.0"},
  "runtime": {"name": "go.dokimi.dev/dokimi-go", "version": "v0.1.0"},
  "seed": "0x5f3a1c9e27d4b860",
  "application": "falsification",
  "checks": [
    {
      "identity": "Store.Read/zero-on-error",
      "key": "store.read/zero-on-error",
      "claim": "zero-on-error",
      "law": "no-value-beside-an-error",
      "machinery": "nothing-beyond-a-caller",
      "shallow": false,
      "result": "proven",
      "reason": "value-beside-error"
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
| `application` | `falsification`, `subject` or `control` |
| `subject` | The name of the subject that the consumer passed to `Run`, for the application `subject` |
| `checks[].result` | `pass`, `fail`, `vacuous` or `excused` for a subject and for the control. `proven`, `passed`, `wrong-reason` or `vacuous` for falsification |
| `checks[].reason` | The reason code that the probe reported |
| `checks[].failure` | The failure record of assert-spec for the assertion that failed, on `fail` and on `wrong-reason` |
| `checks[].shallow` | True when the machinery of the claim is shallower than the depth of its law |
| `unclaimed[].classification` | `shape:<name>`, `mixin:<name>` or `contract:<name>/<role>`, with the names of the shape catalog |
| `unclaimed[].stance` | A stance of the spec |

The results of falsification:

| Result | The check on its planted defect |
|---|---|
| `proven` | Failed with the reason code of its claim |
| `passed` | Passed |
| `wrong-reason` | Failed with another reason code |
| `vacuous` | Did not exercise its claim |

On every result other than `proven` the falsification test fails with a message that contains the identity of the check, the reason code of its claim and the reason code that the probe reported.

A reader skips each field that its version of the schema lacks, and it shows an unknown member of an enumeration as unknown instead of mapping it to another member.

### The bundle and the seed

The bundle of a run is the set of artifacts of its test runner, with their call records and suite records, and the trace files of the tiers that record calls. CI uploads the bundle to the service beside the outputs of the other tools of the run, such as coverage and mutation records.

The seed of a run comes from the environment variable `DOKIMI_SEED`, as a decimal or a `0x` hexadecimal 64-bit number. When the variable is unset, the runtime takes a seed from the clock and prints it with the first failure. A probe that runs the property engine passes the seed of the run to the option `seed` of the engine. Every suite record contains the seed, so a run repeats with `DOKIMI_SEED` set to it.

### The contract between a generated file and its runtime

A generated file and its runtime agree on a compatibility constant, the names of the probes, the constants of the spec and the reason codes:

1. **A compatibility constant.** Each package of a runtime that generated files import exports a constant whose name contains the major version of the runtime, such as `stub.CompatV0` in Go. Every generated file refers to the constant of the major version that the generator targets. A major version of a runtime renames the constant. A mismatch fails the build of the consumer in Go, Rust, Java, Kotlin and C#, and the type check or the module load in TypeScript. In Python and PHP it fails the import of the generated file.
2. **Probe names.** A generated check calls each probe by its identifier in the casing of the language. With a runtime that implements an older spec without the probe, the build or the import fails at that call.
3. **Constants.** Each runtime generates the constants of the entries of the spec from its vendored copy. Generated files refer to those constants, so a file that refers to an entry that the runtime lacks fails the build or the import.
4. **Reason codes.** A probe asserts through the assert library of its language. It passes the library a test handle of the runtime instead of the test itself. The handle receives the failure record of each failed assertion and records the reason code of the assertion in `asserts`. Falsification reads the recorded code and never the text of the failure.

The handle implements the interface that its assert library calls on a test, such as `assert.TB` in Go and `Seat` in TypeScript. In falsification the handle wraps the recorder of the assert library. The recorder keeps the failure records of the check without failing the falsification test.

### The API of a runtime

Each runtime implements the same parts in the conventions of its language:

| Part | Go: `go.dokimi.dev/dokimi-go` | TypeScript: `@dokimi/runtime` |
|---|---|---|
| The constants of the spec | The package `spec`, such as `spec.ErrorsDiffer` and `spec.SharedIdentity` | The module `spec`, such as `probes.errorsDiffer` and `reasons.sharedIdentity` |
| A check | `conform.Check[S]`: the identity and the key, the claim with its law and machinery, the reason of its defect, a function that applies the check to a subject of the type `S`, and a function that returns its planted defect as a subject | `Check<S>`, with the same fields |
| The checks of one declaration | `conform.Suite[S]`: the language, the package, the declaration, the checks, the refusals and the unclaimed classifications | `Suite<S>` |
| The test handle | `*conform.Handle`, which implements `assert.TB` and `assert.Reporter` | `Handle`, which implements `Seat` |
| A probe | A function of the package `probe`, such as `probe.ErrorsDiffer(h, values)` | A function of the module `probe`, such as `errorsDiffer(h, values)` |
| A named value | `conform.Named[T]` | `Named<T>` |
| The application to a subject | `conform.Verify(t, suite, subject)`, where `subject` returns a fresh subject for each check | `verify(suite, subject)` |
| The falsification | `conform.Falsify(t, suite)` | `falsify(suite)` |

`Verify` and `Falsify` run each check as a test of its own, named by the identity of the check. Each check gets a handle of its own, and the probe reads the statement of the check from it. `Verify` fails a check whose probe fails, and `Falsify` fails a check whose planted defect does not fail it with the reason of its claim. Each writes the suite record of its application. The subject of a suite is the implementation that the factory of the consumer returns. The subject of the checks of errors is a value of a type that the generated test declares, with the errors of the package.

### The package spec

```go
// Package spec contains the conformance spec of dokimi: the laws, claims,
// gates, probes, reason codes, machineries and stances that the binary
// shares with the runtimes.
package spec

// LawID identifies a law. names.gen.go declares one constant for each.
type LawID string

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

// AssertionID identifies an assertion of the vendored copy of assert-spec.
type AssertionID string

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

	// Asserts contains the reason code of a failure of each assertion that
	// the probe calls.
	Asserts map[AssertionID]ReasonID

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

The package declares the same pair of functions for laws, claims, gates, reason codes, machineries and stances. It decodes the embedded files on the first call. A test of the package validates each embedded file against its schema and resolves these references:

- the law, the gates, the probe, the reason code and the machinery of each claim
- each classification of a claim, against the shape catalog of eidos
- each input of a claim, against the inputs of its probe
- each assertion of a probe, against the vendored copy of assert-spec
- the depth of each law and each law of `distinct-from`
- every classification of the shape catalog, against the claims and the entries of `unclaimed`

### Vendoring

Each runtime commits a copy of the JSON files, `VERSION` and `manifest.json`, with two scripts that it copies unchanged from dokimi:

- `tools/spec-sync.sh <directory>` fetches the files that the manifest of the tag in the environment variable `SPEC_REF` lists, and writes the tag into the file `REF` beside them. dokimi is private, so the script reads the files through the GitHub API with the token in `GITHUB_TOKEN`.
- `tools/spec-check.sh <directory> [--strict]` fails when a vendored file differs from its digest. With `--strict` it also fails when the manifest differs from the manifest of the tag in `REF`.

A scheduled workflow of each runtime compares its copy with the spec of the latest release tag of dokimi. It opens an issue titled "Refresh the vendored spec to <tag>" when the two differ. The workflow reads dokimi with the token of a GitHub App, because the default token of a workflow cannot read another repository.

`manifest.json` has the format of the manifest of assert-spec. `files` contains the digest `sha256:<hex>` of each path, and `digest` is the SHA-256 of the sorted lines `<path> <digest>`:

```json
{
  "digest": "sha256:5b0e…",
  "files": {
    "VERSION": "sha256:2c1a…",
    "probes.json": "sha256:9f4d…"
  },
  "version": "1.0.0"
}
```

dokimi vendors assert-spec and mutate-spec with the scripts `spec-sync.sh` and `spec-check.sh` of each definition, as assert-go vendors assert-spec. A scheduled workflow compares each copy with the main branch of its definition and opens an issue when the two differ. Each runtime requires the release of its assert library that implements the version of assert-spec in the vendored copy.

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
| Every reference of an entry resolves in the spec, in the vendored copy of assert-spec and in the shape catalog of eidos | A test of the package spec |
| Each assertion of a probe has one reason code | The type of `asserts`, a map from an assertion to one code |
| Every classification of the shape catalog has a claim or a stance | A test of the package spec |
| A vendored copy matches the manifest of the tag in `REF` | `spec-check.sh --strict` in the gate of each runtime |
| Falsification matches a failure by its reason code and never by its text | The test handle of each runtime, and the corpus |
| A generated file and its runtime agree on their major version | The compatibility constant |
| Every suite record contains its seed | The runtime of each language |

### Failure handling

| Failure | Result | Recovery |
|---|---|---|
| A YAML file has an entry that its schema rejects | The test of the package spec fails | Correct the entry |
| An entry refers to a law, gate, probe, reason code, machinery or classification that does not exist | The test of the package spec fails and reports the reference | Add the entry, or correct the reference |
| The entry of a probe refers to an assertion that the vendored copy of assert-spec lacks | The test of the package spec fails | Add the assertion through the RFC process of assert-spec, then vendor the new definition |
| A rendered file, `names.gen.go` or the manifest is stale | The generate step fails the gate of dokimi | Run `go generate ./spec` and commit the result |
| A vendored file differs from its digest | `spec-check.sh` fails the gate of the runtime | Run `spec-sync.sh` with the tag in `REF` |
| A newer release of dokimi has another spec | The scheduled workflow opens an issue | Run `spec-sync.sh` with the newer tag and implement the new entries |
| A generated file calls a probe that the runtime lacks | The build or the import of the consumer fails | Upgrade the runtime |
| A generated file and its runtime differ in their major version | The compatibility constant fails the build or the import | Upgrade the older side |
| A probe calls an assertion that is missing from `asserts` in its entry | The handle records the failure without a code, and falsification reports `wrong-reason` | Fix the probe or its entry |
| The overlay of a language declares an assertion of a probe absent | Each claim of the probe refuses in that language with the reason of the overlay entry | None in dokimi. The overlay contains the reason |

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

The verdict reverses when a claim needs a binding that the spec cannot state as data, such as a partner that only a function can resolve.

### A repository of its own for the spec

A repository dokimi-spec contains the conformance spec, and the binary and each runtime vendor it, as each assert library vendors assert-spec.

**Why not:** every assert library is an equal implementation of its standard, so that standard has a repository of its own. The conformance spec has one producer of checks. A change that adds a claim with a new probe also changes the templates of its defect form and the expected outcomes of the corpus. In the dokimi repository that change is one commit, one review and one release. With dokimi-spec it is a commit there and a sync commit in dokimi, and between the two the generator and its spec disagree.

The verdict reverses when the runtimes are public while dokimi is private, because the spec then needs a public home. It also reverses when a second program generates calls to the probes.

### Runtimes vendor the main branch

A runtime vendors the spec of the main branch of dokimi, as the assert libraries vendor the main branch of assert-spec.

**Why not:** a consumer runs tests that a released binary generated. The main branch can contain a change of meaning that no release contains yet, and a runtime that vendors it implements a spec that the binary of the consumer does not use. A release tag fixes the spec that a binary and a runtime share.

### Failures matched by their text

A planted defect passes when the failure of its check contains a stated piece of text, as in testkit.

**Why not:** the text of a failure is written for a person in the language of the runtime. A TypeScript probe that had to print the same sentence as a Go probe would turn the wording of one runtime into a contract for all runtimes. A reason code is the same in every language, and each runtime words the text for its own readers.

### A report file beside the artifact of the test runner

The runtime writes a JSON report of its own into a directory, apart from the artifact of the test runner.

**Why not:** the service reads the call records from the artifact of the test runner. A second file per suite needs a second upload path and a second parser in the service. A suite record in the same artifact joins the call records of a suite with its claims without another upload.

## Drawbacks

- **A sync commit for each new probe.** Each runtime takes a change of the spec through a release of dokimi and a sync commit of its own. With eight runtimes, an added probe means eight sync commits and eight implementations. A new claim that calls an existing probe does not need a change in a runtime.
- **A token in each runtime.** dokimi is private, so the sync script and the scheduled workflow of each runtime need a token that can read it.
- **A constant generator in each runtime.** Each runtime turns the vendored JSON into the constants of its language with a generator of its own.
- **Every probe in every runtime.** The probes are host code, so each probe exists once per runtime, and the corpus judges each copy.
- **One reason per assertion of a probe.** A probe that checks two violations through one assertion, such as two values that it compares with `equal`, splits into two probes.
- **New assertions wait on assert-spec.** A probe that needs an assertion that assert-spec lacks waits for an RFC of assert-spec, an entry in the overlay of each language and an implementation in each assert library.
- **Two copies of each entry in dokimi.** The YAML file and its rendered JSON file are both committed. The generate step keeps them equal.
- **A YAML library in the dokimi module.** The renderer reads YAML, so the module requires a YAML library that the binary does not use at run time.

## Open questions

None.

## References

| What | Where |
|---|---|
| The shape catalog of eidos | <https://github.com/dokimasia/eidos/blob/main/docs/architecture/12-shape-catalog.md> |
| assert-spec: the definition, its rendering and its vendoring | <https://github.com/dokimasia/assert-spec/blob/main/README.md> |
| The vendoring scripts of assert-spec | <https://github.com/dokimasia/assert-spec/tree/main/tools> |
| The failure record of assert-spec | <https://github.com/dokimasia/assert-spec/blob/main/spec/failure.md> |
| The recording of a run in assert-spec | <https://github.com/dokimasia/assert-spec/blob/main/spec/recording.md> |
| The overlays of assert-spec | <https://github.com/dokimasia/assert-spec/blob/main/spec/overlays.md> |
| The drift workflow of assert-go | <https://github.com/dokimasia/assert-go/blob/main/.github/workflows/drift.yml> |
| mutate-spec: a mutation run and its record | <https://github.com/dokimasia/mutate-spec/blob/main/docs/rfc/0002-a-mutation-run-and-its-record.md> |
| JSON Schema 2020-12 | <https://json-schema.org/draft/2020-12/json-schema-core> |
