---
rfc: 0016
title: The suite record in the artifact of a run
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

# RFC-0016: The suite record in the artifact of a run

## Summary

Each application of a suite writes its suite record into the artifact that the test runner already writes, beside the call records of the assert library. In Go conform writes the record as attribute events of `go test -json` under the key `dokimi.suite.<seq>`. In TypeScript the runtime writes it into the meta object of the test, which the JSON report of Vitest contains. `conform.Suite` gains the refusals and the unclaimed classifications that the generator stated, so the record contains what the suite declined as well as what it proved. The command `dokimi report` reads the artifact of one run or of every package of a repository and prints the claims, the proofs, the refusals and the vacuous checks. assert-spec reserves the key for dokimi through its overlays.

## Motivation

conform writes no record and no file today. Its only uses of the environment are the seed and a temporary directory (`conform/handle.go`). A run of a suite tells a person, through the text of failures, what went wrong. It tells nobody what the suite claimed, what it proved and what it declined. The service cannot read a result that never leaves the process, and neither can a gate or a developer without the service.

RFC-0002 defines the suite record: its identity, the version of the spec, the seed, the result and the reason code of each check, the refusals and the unclaimed classifications. It states that a runtime writes the record into the artifact of the test runner, under the key `dokimi.suite.<seq>`, beside the call records of assert-spec. It does not state where each runner stores it.

The runners can store it without a file of dokimi. Go 1.25 added `T.Attr`, which `go test -json` reports as an event with the action `attr` and the fields `Key` and `Value`. A value cannot contain a newline. assert-spec already uses this channel for its call records: one event per call, with the key `dokimi.assert.<seq>`, split across events with the same key when a record exceeds the line buffer of 4,096 bytes (`assert-spec/overlays/go.json`, entry `records`). Vitest stores a JSON object in the meta of each test, and its JSON reporter writes the meta of every test (`vitest.dev/api/advanced/metadata`). assert-typescript puts its call records under the meta key `dokimi.assert`.

assert-spec reserves the prefix `dokimi.` of the keys of a record for the tools of dokimasia, and the entry `records` of each overlay lists the kinds of record that the artifact can contain, `dokimi.assert` and `dokimi.suite` (`assert-spec/spec/overlays.md`, section Records). Its table of kinds names dokimi as the definition of `dokimi.suite`. CTRF solves the same problem with a rule of its own, which puts every property that its specification does not define inside an object named `extra`.

## Detailed design

### The fields of the suite

The generated suite states its refusals and its unclaimed classifications, so the runtime can write them:

```go
// Refusal is a licensed claim whose gate fails for a callable.
type Refusal struct {
	Claim    spec.ClaimID
	Callable string
	Why      string
	Remedy   string
}

// Unclaimed is a classification of a callable without a check, with its
// stance.
type Unclaimed struct {
	Classification string
	Callable       string
	Stance         spec.StanceID
	Note           string
}

type Suite[S any] struct {
	// ... existing fields

	// Package is the import path of the package of the declaration.
	Package string
	// Refusals are the licensed claims that the generator refused.
	Refusals []Refusal
	// Unclaimed are the classifications without a check.
	Unclaimed []Unclaimed
}
```

### Writing the record in Go

`conform.Falsify` writes one record when its application ends. `conform.Run` writes two: the falsification of the derived checks and the application to the subject.

```go
// record writes the suite record of an application through t.Attr under the
// key dokimi.suite.<seq>, where seq numbers the suite records of t from 1.
// A record whose event does not fit the line buffer of 4,096 bytes of
// test2json is split over consecutive attribute events with the same key,
// whose values join into the record, as assert-spec splits a call record.
func record[S any](t T, seq int, r suiteRecord)
```

A record is written on the test of the application, not on the subtest of a check, so one test of the artifact contains the record of the whole suite. The seed is the seed of the run, from `DOKIMI_SEED` or from the clock. conform writes the record whatever the outcome, and also after a fault that ends the application, with the field `ended` set to the fault.

### Writing the record in TypeScript

`run` and `falsify` of `@dokimi/runtime` write the records into the meta of the test of the application, as an element of the array under the key `dokimi.suite`. assert-typescript keeps its call records in the same form under `dokimi.assert`. Vitest writes the meta into its JSON report.

### The kind of record in assert-spec

The overlays of assert-spec list `dokimi.suite` in the key `kinds` of their entry `records`, and their `location` states the rule of every kind: the kind, then the number of the record in its test, and a split over consecutive events of the same key. A reader of the artifact therefore finds the suite records through the overlay of the language, as it finds the call records.

### dokimi report

The binary reads the artifacts of a run and prints what the suites claimed. It needs no service:

```text
$ go test -json ./... > run.json
$ dokimi report run.json
example.com/shop/catalog  Store
  claims     12   proven 11   wrong-reason 0   passed 1
  subject    12   pass 10     vacuous 1        fail 0     excused 1   skipped 0
  refused     2   Store.Scan/zero-on-error: eidos derives no sample of the parameter visit of Scan
  unclaimed   1   mixin:concurrent on Store.Put: needs concurrent callers
  vacuous         Store.Read/zero-on-error: no call of Read failed
seed 0x5f3a1c9e27d4b860   spec 1.0.0   conform v0.9.0
```

`dokimi report --json` writes one aggregated JSON document for the repository. `dokimi report` exits with a status other than 0 when a record has a falsification result other than `proven`, so CI can run it as a gate after the tests.

## Alternatives considered

### A report file of conform

conform writes a JSON file per suite into a directory that an environment variable names.

**Why not:** the service and `dokimi report` would need a second upload path and a second parser beside the artifact of the runner, which already contains the call records. RFC-0002 records this choice.

### A JUnit XML or CTRF report

conform writes its result as JUnit properties or as CTRF, which many services ingest.

**Why not:** both describe tests, not claims. A suite record would be a member of `extra` or of a property bag, and no tool would read it there. A converter from the artifact of the runner to CTRF can be written once the record exists.

### The record on each check

Each check writes its own part of the record on its subtest.

**Why not:** the refusals and the unclaimed classifications belong to the suite, not to a check. A reader would rebuild the suite from the parts of many subtests and miss a suite whose checks all failed to start.

## Drawbacks

- **Go 1.25 or later.** `T.Attr` exists from Go 1.25. conform already requires a later Go.
- **A larger artifact.** A suite of 200 checks writes a record of about 40 KB in about 11 attribute events.
- **The split of long records.** Each reader joins the consecutive events of one key into one record, as it joins a call record of assert-spec.

## Open questions

- Does `dokimi report` also read the call records of assert-spec, to show the failing assertion of each failed check?

## References

| What | Where |
|---|---|
| Go 1.25, T.Attr | <https://go.dev/doc/go1.25> |
| The source of test2json, the attr action | <https://github.com/golang/go/blob/master/src/cmd/internal/test2json/test2json.go> |
| Vitest, the meta of a test | <https://vitest.dev/api/advanced/metadata> |
| assert-spec, the recording of a run | <https://github.com/dokimasia/assert-spec/blob/main/spec/recording.md> |
| assert-spec, the overlays | <https://github.com/dokimasia/assert-spec/blob/main/spec/overlays.md> |
| CTRF, the rule on extra properties | <https://github.com/ctrf-io/ctrf/blob/main/spec/ctrf.md> |
