---
adr: 0014
title: The directive stub takes the name template of the field of each method
status: Accepted
date: 2026-10-28
supersedes: RFC-0004 in part
superseded-by: none
rfc: 0004
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0014: The directive stub takes the name template of the field of each method

## Status

Accepted

## Context

RFC-0004 names the field that configures a method of a double after the method, as `OnGet` for `Get`. A subject with the methods `Node` and `OnNode` gets a double whose field `OnNode` and method `OnNode` have one name. The stub generator reports `Collision`, and RFC-0004 gives the remedy "Rename one of the methods":

- A rename changes the API of the subject to fit the names of its double. Issue 41 reports a released struct whose method `OnNode` has six callers in three packages.
- The enum generator already takes name templates, such as `Parse{{.Enum}}`.

## Decision

We will add the param `field` to the directive `stub`. Its value is the name template of the field that configures each method of the double.

- The only action of the template is `{{.Method}}`, the name of the method in PascalCase, and the template uses it at least once. The text around the action follows the rules of the name templates of the enum generator.
- The default is `On{{.Method}}`, so a double without the param keeps its names.
- Go writes the rendered name in PascalCase, and TypeScript in camelCase. `//dokimi:stub field={{.Method}}Stub` gives `Node` the field `NodeStub` in Go and `nodeStub` in TypeScript.
- The template applies to each double that the directive writes, the doubles of its handles included.
- The stub generator reports `BadField` at a declaration whose param `field` is not such a template, and does not write a double of it.

The enum generator and the stub generator parse their templates through one package under `internal/`.

## Alternatives Considered

### A recorder apart from the double

gomock records expectations on the recorder that `EXPECT()` returns, so a mock has the methods of its interface and one more. The change would move the configuration of every generated double, as `s.OnGet.Returns` in each test of each consumer.

### An option of the plugin

An option in the config file would rename the fields of every double of a module to resolve the collision of one subject.

## Consequences

**Positive:**

- A subject gets its double whatever the names of its methods, without a change of its API.

**Negative:**

- A double with the param has other field names than the doubles of the same module without it.

## References

| What | Where |
|---|---|
| RFC-0004, whose names of the fields of a double this decision opens to a template | [0004-doubles-builders-enumerations-and-errors.md](../rfc/0004-doubles-builders-enumerations-and-errors.md) |
| Issue 41, a struct with the methods Node and OnNode does not get a double | <https://github.com/dokimasia/dokimi/issues/41> |
