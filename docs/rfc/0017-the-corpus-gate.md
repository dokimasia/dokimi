---
rfc: 0017
title: The corpus gate of each language
author: Roy Klopper <roy.klopper@stealthscale.io>
status: Accepted
created: 2026-10-28
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

# RFC-0017: The corpus gate of each language

## Summary

`dokimi conform --lang <language>` judges a language binding by the lock and the records that its corpus produces. It checks that the generated files of the corpus of the language are current, runs the tests with the toolchain of the language, reads the suite records from the artifact of the runner, and compares the outcomes of each fixture with the expected outcomes of the fixture. The expected outcomes are written once for every language, in the key form of the spec. A language states each known difference in an overlay, as a refusal by a gate of the spec or as a stance of the spec, with the reason why the language differs. The gate fails on a difference that the overlay does not declare, on an overlay entry that matches no difference, and on an overlay entry whose outcome now meets the expectation. A language declares a conformance level, and the gate judges the fixtures of that level and below.

## Motivation

The Go corpus is a module of 168 fixtures whose generated suites build and run under `go test` (`corpus/go`). Its gate is the test run: a fixture passes when its tests pass. The run shows that the generated code works. It does not compare what each language claims, proves and declines with what the fixture should produce. The TypeScript corpus beside it has 152 fixtures and is judged the same way.

When a second language binds the suite generator, a fixture can pass in both languages while the two disagree. TypeScript can refuse a claim that Go binds, emit a check under another identity, or prove a check with another reason code. Each language passes its own tests, and the disagreement shows only when a consumer compares them. Pact ran one specification in ten languages, and its compliance matrix records the drift: Ruby, Python and PHP remained at version 2 of the specification while JVM, JS, Go and Rust reached version 4.

The protobuf conformance runner shows a gate that catches drift in both directions. Each language keeps a list of the tests that it fails, with the expected failure message. The runner fails when a listed test passes, when an entry matches no test, and when a failure message changes. Bowtie and the OpenTelemetry compliance matrix publish partial compliance, and neither can detect that a known divergence has gone stale.

## Detailed design

### Fixtures, subjects and controls

A fixture is a directory of the corpus of a language, such as `iface/mixin/ttl` under `corpus/go/`. The corpus of each language puts one fixture at one path, so the path identifies the fixture in every language. The directories above the fixtures are its groups: `iface/composite`, `iface/contract`, `iface/detector`, `iface/door`, `iface/lang`, `iface/mixin`, `enum`, `errors`, `integration` and `struct`.

The subjects of a fixture are its types whose doc comment has the directive `subject`, other than its controls. 143 fixtures of the Go corpus have one subject. `iface/composite/causal-chain` and `iface/lang/suitedhandle` have two, and 23 fixtures have none, such as the fixtures of enumerations and errors, whose generated tests do not apply a suite. A control is a correct implementation of the declaration of a subject with another representation of its state. A check that fails on a control rejects behaviour that the declaration permits. The corpus writes a subject and a control as the directory of its file relative to the corpus, a dot and the name of its type, such as `iface/mixin/ttl/ttltest.InMemory`.

The suite record of an application contains its subject as the type of the value that the factory of the application builds: the import path of its package, a dot and its name. An instance of a generic type has the generic type as its subject, so `memory.Store[int,int]` is the subject `iface/lang/generic/memory.Store`. The gate removes the module path of the corpus, so `go.dokimi.dev/dokimi/corpus/go/iface/mixin/ttl/ttltest.InMemory` becomes the subject `iface/mixin/ttl/ttltest.InMemory`. The records of an application to any other type, such as the tests that a generated suite runs on itself, belong to no phase.

### Expected outcomes

Each fixture has one file of expected outcomes, `corpus/expect/<fixture>.yaml`, such as `corpus/expect/iface/mixin/ttl.yaml`. The file uses terms that do not depend on a language:

```yaml
fixture: iface/mixin/ttl
level: 2
subjects:
  - type: iface/mixin/ttl/ttltest.InMemory
controls:
  - iface/mixin/ttl/deadline.Store
declarations:
  mixed:
    checks:
      mixed.put/cancel:
        claim: cancel
        law: cancelled-context-is-reported
        reason: ignored-cancel
      mixed.put/smoke:
        claim: smoke
        law: no-panic-on-any-input
        reason: panicked
      mixed.read/cancel:
        claim: cancel
        law: cancelled-context-is-reported
        reason: ignored-cancel
      mixed.read/expires:
        claim: expires
        partners:
          write: mixed.put
        law: entry-lapses-on-the-clock
        reason: not-expired
      mixed.read/smoke:
        claim: smoke
        law: no-panic-on-any-input
        reason: panicked
      mixed.read/zero-on-error:
        claim: zero-on-error
        law: no-value-beside-an-error
        reason: value-beside-error
    authored: []
    ceilings:
      - mixed.read/allocation-budget
      - mixed.read/latency-budget
      - mixed.read/mean-latency-budget
    refusals: {}
    unclaimed:
      - classification: shape:writer
        callable: mixed.put
        stance: silent
```

| Field | Value |
|---|---|
| `fixture` | The path of the fixture |
| `level` | The conformance level of the fixture |
| `subjects` | The type of each subject of the fixture, and under `excused` the keys of the checks that the subject excuses with a reason |
| `controls` | The controls that the corpus applies the suites of the fixture to |
| `declarations` | The expected outcomes of each declaration under the directive `suite`, by the key form of its name |
| `checks` | Each derived check by its key, with its claim, its partners, its law and the reason code of its probe |
| `authored` | The keys of the checks of the directive `check` |
| `ceilings` | The keys of the ceilings of the budgets of the declaration |
| `refusals` | The gate of each refusal by the key of its check |
| `unclaimed` | Each unclaimed classification, with its callable and its stance |

An excusal belongs to one subject. `replica.Log` of `iface/composite/causal-chain` refuses each append, so it excuses four checks that `archive.Log` of the same fixture runs.

Each key and each callable is in the key form of the spec, so `Mixed.Read` in Go and `Mixed.read` in TypeScript match one key. The records contain the callable of a refusal and of an unclaimed classification by its name, such as `Mixed.Put`, and the gate converts it with `spec.Key`. The lock contains a refusal by its key and no gate, and a record contains it by its callable and its claim with its gate. The gate joins the two by the callable and the claim.

`level` is the conformance level of the fixture. Its group decides it:

| Level | Fixtures |
|---|---|
| 1 | `iface/lang` and `struct`, the forms of the declarations of a language, with the claims that every callable with a failure channel or a context gets: smoke, cancellation, the zero value and absence |
| 2 | `iface/mixin`, the claims of the mixins, and `iface/detector`, which classifies a callable by its shape as a mixin classifies it by its directive |
| 3 | `iface/contract`, the claims of the contracts, and `enum` and `errors`, which have the fixed claims of enumerations and errors |
| 4 | `iface/composite`, the composites of several axes, and `iface/door` and `integration`, which combine several parts |

### Writing the expected outcomes

`dokimi conform --lang <language> --update` rewrites the field `declarations` of each file from the lock of the run, from its records and from the claims of the spec. The law and the reason code of each check come from its claim in the spec, and the gate of each refusal comes from the records. The checks, the authored checks, the ceilings, the refusals, the unclaimed classifications and the partners come from the lock. `--update` keeps the fields `level`, `subjects` and `controls`, which a person writes. The new file keeps the comment lines that open the file that it replaces. A person reviews each change of a file as a change of the contract of its fixture.

The corpus of Go is the reference corpus. Its run writes the first expected outcomes, and the overlay of Go starts empty. The first files take the subjects of each fixture from the directive `subject`, the controls from the shape families of the Go corpus, the excusals from the subject records of a run, and the level from the group of the fixture.

### The overlay of a language

Each language has one overlay, `corpus/<language>/overlay.yaml`:

```yaml
language: typescript
level: 1
diverge:
  - fixture: iface/mixin/ttl
    key: mixed.read/zero-on-error
    expect: check
    actual: refused
    gate: returns-beside-error
    why: A TypeScript call that fails throws and returns no value
```

An entry weakens an expectation through a refusal or through a stance:

- **A refusal.** The language refuses a claim that the expectation binds, through a gate of the spec. The entry has the actual outcome `refused` and the gate, and the record of the refusal must contain the same gate.
- **A stance.** The language lists a classification as unclaimed that the expectation binds or refuses. The entry has the actual outcome `unclaimed` and one of the stances of the spec, and the lock must list a classification of the callable of the key with the same stance. The gate counts the missing check and the unclaimed classification as one difference.

Each entry has the fixture, the key of an outcome, the expected outcome, the actual outcome, the gate or the stance, and the reason why the language differs. A key repeats across fixtures, such as `mixed.read/zero-on-error` in each fixture of a mixin whose declaration is `Mixed`, so each entry has a fixture. The why and the remedy of a refusal come from the gate in the spec, so the entry does not repeat them.

A fixture that a language cannot state, such as a Go channel in TypeScript, has one entry with the fixture, the stance `blocked`, the reason and no key. The gate fails on a fixture that a language lacks without such an entry. The fixtures of a language are the directories of the groups of the expected fixtures in its corpus, so the gate also fails on a fixture without a file of expected outcomes.

### The phases

| Phase | Compares | Fails on |
|---|---|---|
| Generation | The lines of the fixture in the lock of the plan of the language, and the gate of each refusal in the records, with the checks, the authored checks, the ceilings, the partners, the refusals and the unclaimed classifications of the expectations | A line that is missing, added or changed, and a refusal with another gate, unless an entry of the overlay declares the difference |
| Falsification | Each falsification record of the declarations of the fixture with the claim, the law and the reason code of each check | A check that its planted defect does not prove, such as a check that passes on it, fails for another reason or is vacuous on it, and a derived check without a falsification record |
| Subject | The records of the applications to each subject | A subject without a record, a check that fails, is vacuous or is skipped, an excusal of the subject that the record does not show, and a derived or authored check without a record on the subject |
| Controls | The records of the applications to each control | A control without a record, and a check that fails on a control |
| Declaration | The entries of the overlay with the differences of the generation | An entry without a difference, such as an entry whose outcome now meets the expectation |

A falsification, subject or control record whose field `ended` lists problems fails its phase, because its application stopped before its checks. A falsification that excuses or skips a check has no result to judge. A check whose key the generation finds changed has no record in the later phases, so those phases do not report its absence again.

The application to a subject builds a value only for a check that it runs. For a fixture whose declarations have only refusals and unclaimed classifications, such as `iface/detector/streamconsumer`, the application does not build a value, and its records do not have a subject. The subject phase does not need a record of such a subject. A ceiling of latency runs in a benchmark, which the test command does not run, so the subject phase judges a ceiling only where it has a record.

An overlay declares only differences of the generation. The overlay weakens what a language generates, and the later phases judge what the runtime of the language proves and applies.

### The command

```text
$ dokimi conform --lang go
files: 0 create, 0 update, 766 unchanged, 0 stale, 0 drifted, 0 foreign, 0 withheld; plans: 0 committed, 0 failed, 0 cancelled, 2 prepared, 0 skipped
iface/contract/appender
  falsification  FAIL  appender.run/cancel: want proven by ignored-cancel, got wrong-reason by wrong-sentinel
go: level 4, 168 fixtures, 167 conformant, 0 declared differences
```

The command checks with `dokimi run --check` that the generated files of `corpus/<language>` are current, then runs the test command of the language: `go test -count=1 -json ./...` for Go and `pnpm exec vitest run --reporter=json` for TypeScript. In a run of `go test` that the cache can replay, the test binary writes a test log, and the runtime does not count allocations. The command runs the tests of Go without the cache, so each allocation ceiling has a result. The command reads the lock of the plan of the language and the suite records of the artifact. It writes each finding under its fixture, and a summary of the fixtures that it judged, those without a finding, and the differences that the overlay declared. A fixture without a finding has no line. The command exits with a status other than 0 when the generated files are out of date, when a test fails, when a record does not validate, and on any finding. Under `--update` it writes the expected outcomes in place of the judgement.

The gate of dokimi runs the command for every language that has a binding. A spec of the command runs it on the corpus of the repository, so the Go tests of the gate judge the corpus of Go.

The package `internal/corpus` contains the types and the functions of the gate:

```go
// Expectation is the expected outcome of one fixture, as its file under
// corpus/expect states it. A file writes the fixture, its level, its
// subjects, its controls and its declarations in that order.
type Expectation struct {
	Fixture      string                 `yaml:"fixture"`
	Declarations map[string]Declaration `yaml:"declarations"`
	Subjects     []Subject              `yaml:"subjects"`
	Controls     []string               `yaml:"controls"`
	Level        int                    `yaml:"level"`
}

// Subject is a subject of a fixture, with the keys of the checks that the
// subject excuses with a reason.
type Subject struct {
	Type    string   `yaml:"type"`
	Excused []string `yaml:"excused,omitempty"`
}

// Declaration is the expected outcome of one declaration under the
// directive suite.
type Declaration struct {
	Checks    map[string]Check       `yaml:"checks"`
	Authored  []string               `yaml:"authored"`
	Ceilings  []string               `yaml:"ceilings"`
	Refusals  map[string]spec.GateID `yaml:"refusals"`
	Unclaimed []Unclaimed            `yaml:"unclaimed"`
}

// Check is the expected outcome of one derived check. Partners maps the
// name of each partner to the key of its callable.
type Check struct {
	Claim    spec.ClaimID      `yaml:"claim"`
	Partners map[string]string `yaml:"partners,omitempty"`
	Law      spec.LawID        `yaml:"law"`
	Reason   spec.ReasonID     `yaml:"reason"`
}

// Unclaimed is an expected classification of a callable without a check.
type Unclaimed struct {
	Classification string        `yaml:"classification"`
	Callable       string        `yaml:"callable"`
	Stance         spec.StanceID `yaml:"stance"`
}

// Overlay is the overlay of one language: the differences that it declares
// and the highest level that the gate judges.
type Overlay struct {
	Language string  `yaml:"language"`
	Diverge  []Entry `yaml:"diverge"`
	Level    int     `yaml:"level"`
}

// Entry declares one difference of a language from the expected outcomes of
// a fixture. Gate or Stance states how the language differs, and Why states
// the reason. An entry with the stance blocked and no key declares that the
// language cannot state the fixture.
type Entry struct {
	Fixture string        `yaml:"fixture"`
	Key     string        `yaml:"key,omitempty"`
	Expect  string        `yaml:"expect,omitempty"`
	Actual  string        `yaml:"actual,omitempty"`
	Gate    spec.GateID   `yaml:"gate,omitempty"`
	Stance  spec.StanceID `yaml:"stance,omitempty"`
	Why     string        `yaml:"why"`
}

// Phase is a phase of the gate: generation, falsification, subject,
// controls or declaration.
type Phase string

// Outcome is the outcome of one fixture in a run: its lines of the lock and
// the suite records of its declarations.
type Outcome struct {
	Fixture string
	Lock    []string
	Records []report.Record
}

// Finding is one failure of a phase of a fixture, with the expected and the
// actual outcome.
type Finding struct {
	Fixture string
	Phase   Phase
	Key     string
	Want    string
	Got     string
}

// ReadExpectations returns the expectation of each file of fsys whose name
// ends in .yaml, in the lexical order of the files. It returns an error that
// wraps ErrFile for a file that does not decode, for an unknown field, and
// for a fixture that is not the path of its file.
func ReadExpectations(fsys fs.FS) ([]Expectation, error)

// ReadOverlay returns the overlay of the file name of fsys.
func ReadOverlay(fsys fs.FS, name string) (Overlay, error)

// WriteExpectation writes e into its file under dir, and keeps the comment
// lines that open the file that it replaces.
func WriteExpectation(dir string, e *Expectation) error

// Fixtures returns the fixtures of the corpus fsys: each directory in a group
// of the fixtures of expected, as its path, sorted.
func Fixtures(fsys fs.FS, expected []Expectation) ([]string, error)

// Outcomes groups the lines of a lock and the suite records of a run by
// fixture, and removes module, the module path of the corpus, from each
// package and each subject.
func Outcomes(module string, fixtures, lock []string, records []report.Record) []Outcome

// Judge returns the findings of the outcomes of each fixture of expected up
// to the level of overlay, sorted by fixture and then by phase, and the
// number of the entries of overlay that declare a difference. Judge also
// reports each outcome that lacks an expectation, and each expectation of
// the level that lacks both an outcome and a blocked entry.
func Judge(expected []Expectation, overlay Overlay, outcomes []Outcome) (found []Finding, declared int)

// Update returns expected with the declarations of each fixture rewritten
// from its outcome, and keeps the level, the subjects and the controls.
func Update(expected []Expectation, outcomes []Outcome) []Expectation
```

## Alternatives considered

### Expected outcomes per language

Each language keeps its own expected outcomes beside its fixtures.

**Why not:** two sets of expectations drift apart in the same way as two implementations. One set, with declared overlays, keeps every difference visible in one file per language.

### The lock of Go as the expected outcomes

Each language compares its lock with the lock of Go, and declares its differences from Go.

**Why not:** a refusal of Go depends on the signatures of Go, such as a callable that returns two values. The expectations state what each fixture should produce, and each language, Go included, declares its differences in its overlay.

### The subject as the name of its factory

The record contains the name of the factory that the consumer passed to the application, such as `storetest.NewInMemory`.

**Why not:** each generated suite passes an adapter of the factory to the runtime, so the record contains the name of the adapter. An option that passes the name of the factory would add an API to every generated suite. The type of the value that the factory builds identifies the subject without an option, through any adapter.

### Divergence in the harness

A language skips a fixture or a check in the code of its harness, as Bowtie does.

**Why not:** a skip in code has no stance and no remedy, and the gate cannot tell when it goes stale. An overlay entry is data that the gate checks against the outcome.

### A shared implementation of the probes

Every language calls one implementation of the probes through a foreign-function interface, as Pact moved to a core in Rust.

**Why not:** a probe runs inside the test process of the consumer and asserts through the assert library of its language. A foreign library in every test process needs a native build for every platform. The derivation is already shared in the binary, so only the probes and the runtime are written per language.

## Drawbacks

- **One file of expectations per fixture.** 168 files, one for each path of the Go corpus, which the TypeScript corpus shares.
- **The subjects and the controls in two places.** The directive `subject` and `corpus/go/families.go` list the subjects and the controls of the Go corpus, and the expectations list them again.
- **The toolchain of every language in the gate.** The gate of dokimi runs `go` and `pnpm`, and each later language adds its toolchain.
- **A slower gate.** On the Go corpus, `dokimi run --check` takes 1.5 s and the test run 26 s over 188 packages with a warm build cache. With the test binaries in the build cache, the spec of the command judges the whole corpus in 6.5 s. The tests of the gate run the corpus three times: in the test step of its module, and in the spec of the command under the test step and under the race step.
- **Levels need agreement.** The group of a fixture decides its level, and a disputed group changes which fixtures a language must pass.

## Open questions

None.

## References

| What | Where |
|---|---|
| protobuf conformance runner and its failure lists | <https://github.com/protocolbuffers/protobuf/blob/main/conformance/conformance_test.cc> |
| Bowtie, the harness of JSON Schema implementations | <https://docs.bowtie.report/en/stable/implementers/> |
| Pact, the compliance matrix of its specification | <https://docs.pact.io/getting_started/specification> |
| web-platform-tests, expectation metadata | <https://web-platform-tests.org/tools/wptrunner/docs/expectation.html> |
| OpenTelemetry, the compliance matrix | <https://github.com/open-telemetry/opentelemetry-specification/blob/main/spec-compliance-matrix.md> |
