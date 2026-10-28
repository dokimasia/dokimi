---
rfc: 0017
title: The corpus gate of each language
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

# RFC-0017: The corpus gate of each language

## Summary

`dokimi conform --lang <language>` judges a language binding by the records that its generated suites write. For each fixture of the corpus it generates the suite with the binary under test, runs the tests with the toolchain of the language, reads the suite records from the artifact of the runner, and compares them with the expected outcomes of the fixture. The expected outcomes are written once for every language. A language states each known difference in an overlay, as a refusal by a gate of the spec or as a stance of the spec, with the reason why the language differs. The gate fails on a difference that the overlay does not declare, on an overlay entry that matches no outcome, and on an overlay entry whose outcome now meets the expectation. A language declares a conformance level, and the gate judges the fixtures of that level and below.

## Motivation

The Go corpus is a module of 168 fixtures whose generated suites build and run under `go test` (`corpus/go`). Its gate is the test run: a fixture passes when its tests pass. The run shows that the generated code works. It does not compare what each language claims, proves and declines with what the fixture should produce, and the TypeScript corpus beside it is judged the same way.

When a second language binds the suite generator, a fixture can pass in both languages while the two disagree. TypeScript can refuse a claim that Go binds, emit a check under another identity, or prove a check with another reason code. Each language passes its own tests, and the disagreement shows only when a consumer compares them. Pact ran one specification in ten languages, and its compliance matrix records the drift: Ruby, Python and PHP remained at version 2 of the specification while JVM, JS, Go and Rust reached version 4.

The protobuf conformance runner shows a gate that catches drift in both directions. Each language keeps a list of the tests that it fails, with the expected failure message. The runner fails when a listed test passes, when an entry matches no test, and when a failure message changes. Bowtie and the OpenTelemetry compliance matrix publish partial compliance, and neither can detect that a known divergence has gone stale.

## Detailed design

### Expected outcomes

Each fixture has one file of expected outcomes, `corpus/expect/<family>/<fixture>.yaml`, in terms that do not depend on a language:

```yaml
fixture: mixin/ttl
level: 2
declarations:
  Mixed:
    checks:
      mixed.read/expires: {claim: expires, law: an-entry-lapses-on-the-clock, reason: unexpired}
      mixed.read/zero-on-error: {claim: zero-on-error, law: no-value-beside-an-error, reason: value-beside-error}
    refusals: []
    unclaimed:
      - {classification: "mixin:ttl", callable: mixed.put, stance: silent}
    subject: {pass: all, vacuous: none}
    controls: {fail: none}
```

The keys of checks and callables are the key form of RFC-0002, so `Mixed.Read` in Go and `Mixed.read` in TypeScript match one key. `level` is the conformance level of the fixture:

| Level | Fixtures |
|---|---|
| 1 | The claims that every callable with a failure channel or a context gets: smoke, cancellation, the zero value and absence |
| 2 | The claims of the mixins |
| 3 | The claims of the contracts |
| 4 | The composites of several axes |

### The overlay of a language

Each language has one overlay, `corpus/<language>/overlay.yaml`:

```yaml
language: typescript
level: 1
diverge:
  - key: mixed.read/zero-on-error
    expect: check
    actual: refused
    gate: value-beside-error
    why: A TypeScript call that fails throws and returns no value
```

An entry weakens an expectation through a refusal or through a stance:

- **A refusal.** The language refuses a claim that the expectation binds, through a gate of the spec or through an assertion that the overlay of the language in assert-spec declares absent. The entry names the gate or the assertion, and the record must name the same one.
- **A stance.** The language lists a classification as unclaimed that the expectation binds or refuses. The entry names one of the stances of the spec, and the stance of the record must be the same.

Each entry has the key of an outcome, the expected outcome, the actual outcome, the gate or the stance, and the reason why the language differs. The why and the remedy of a refusal come from the gate in the spec, so the entry does not repeat them.

A fixture that a language cannot state, such as a Go channel in TypeScript, has one entry with the path of the fixture as its key, the stance `blocked` and the reason. The gate fails on a fixture that a language lacks without such an entry.

### The phases

| Phase | Compares | Fails on |
|---|---|---|
| Generation | The lock of the run with the checks, the refusals and the unclaimed classifications of the expectations | A check that is missing, added or renamed |
| Falsification | The falsification record with the reason of each check | A check that passes on its planted defect, fails for another reason, or is vacuous on it |
| Subject | The record of the application to the subject of the fixture | A check that fails, is vacuous or is skipped on the subject |
| Controls | The records of the applications to the controls of the fixture | A check that fails on a control |
| Declaration | Every difference of the phases above with the overlay | A difference without an entry, an entry without a difference, and an entry whose outcome meets the expectation |

### The command

```text
$ dokimi conform --lang typescript
mixin/ttl
  generation     ok    2 checks, 0 refused, 1 unclaimed, 1 declared divergence
  falsification  ok    1/1 proven with the reason of its claim
  subject        ok    1 pass, 0 vacuous
  controls       ok    0 failures on 1 control
  declaration    ok    1 divergence declared
contract/appender
  falsification  FAIL  appender.run/cancel failed with wrong-sentinel, expected ignored-cancel
typescript: level 1, 41 fixtures, 40 conformant, 3 declared divergences
```

The command runs `dokimi run` on `corpus/<language>`, then the test command of the language: `go test -json ./...` for Go and `pnpm exec vitest run --reporter=json` for TypeScript. It reads the records as RFC-0016 describes them. It exits with a status other than 0 on any failure, and the gate of dokimi runs it for every language.

## Alternatives considered

### Expected outcomes per language

Each language keeps its own expected outcomes beside its fixtures.

**Why not:** two sets of expectations drift apart in the same way as two implementations. One set, with declared overlays, keeps every difference visible in one file per language.

### Divergence in the harness

A language skips a fixture or a check in the code of its harness, as Bowtie does.

**Why not:** a skip in code has no stance and no remedy, and the gate cannot tell when it goes stale. An overlay entry is data that the gate checks against the outcome.

### A shared implementation of the probes

Every language calls one implementation of the probes through a foreign-function interface, as Pact moved to a core in Rust.

**Why not:** a probe runs inside the test process of the consumer and asserts through the assert library of its language. A foreign library in every test process needs a native build for every platform. The derivation is already shared in the binary, so only the probes and the runtime are written per language.

## Drawbacks

- **One file of expectations per fixture.** About 170 files, written once, and changed with each change of a claim.
- **The toolchain of every language in the gate.** The gate of dokimi runs `go` and `pnpm`, and each later language adds its toolchain.
- **A slower gate.** A full run generates and tests every fixture in every language.
- **Levels need agreement.** The level of a fixture is a judgement, and a disputed level changes which fixtures a language must pass.

## Open questions

- Is the gate fast enough to run on every pull request, or does a pull request run the fixtures that its change touches?

## References

| What | Where |
|---|---|
| protobuf conformance runner and its failure lists | <https://github.com/protocolbuffers/protobuf/blob/main/conformance/conformance_test.cc> |
| Bowtie, the harness of JSON Schema implementations | <https://docs.bowtie.report/en/stable/implementers/> |
| Pact, the compliance matrix of its specification | <https://docs.pact.io/getting_started/specification> |
| web-platform-tests, expectation metadata | <https://web-platform-tests.org/tools/wptrunner/docs/expectation.html> |
| OpenTelemetry, the compliance matrix | <https://github.com/open-telemetry/opentelemetry-specification/blob/main/spec-compliance-matrix.md> |
