---
adr: 0004
title: An unexported enumeration has unexported members and no checks
status: Accepted
date: 2026-10-23
supersedes: RFC-0004 in part
superseded-by: none
rfc: 0004
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0004: An unexported enumeration has unexported members and no checks

## Status

Accepted

## Context

RFC-0004 writes the surface of each enumeration under the directive `enum` and its checks in the external test package. The name templates rendered with the name of the enumeration as written, so `Parse{{.Enum}}` gave `Parsestate` for the unexported `state`. The settle of eidos respelled that unexported declaration, the checks referred to the old name, and the run failed with EID-0081.

- Issue 16 reports 16 unexported enumerations in 8 packages of sharder, such as the states and the events of its state machines.
- A test in the package `x_test` cannot refer to an unexported type of `x`, and ADR-0001 generates black-box tests alone.

## Decision

We will write the surface of an unexported enumeration with unexported names, and write no checks for it.

- The name templates render with the name of the enumeration in PascalCase, as `State` for `state`. The Go part then writes the functions and the error in camelCase, as `parseState`, `stateValues` and `errUnknownState`.
- The methods keep their names, as `String`, `IsValid`, `MarshalText` and `UnmarshalText`, because the interfaces of the standard library call them.
- The enum generator records no checks of an enumeration that its package does not export. The package tests the enumeration through its own API.

## Alternatives Considered

### The checks in the package of the enumeration

The checks could refer to the unexported type from an internal test file. ADR-0001 refuses a generated white-box test.

### A refusal of an unexported enumeration

A finding would keep such an enumeration on hand-written methods or stringer, as sharder keeps 16 of them today.

## Consequences

**Positive:**

- A package can put its internal state machines under the directive `enum`.

**Negative:**

- The surface of an unexported enumeration has no generated spec. The tests of the package cover it.

## References

| What | Where |
|---|---|
| RFC-0004, whose checks of every enumeration this decision narrows | [0004-doubles-builders-enumerations-and-errors.md](../rfc/0004-doubles-builders-enumerations-and-errors.md) |
| ADR-0001, every generated test is a black-box test | [0001-every-generated-test-is-a-black-box-test.md](0001-every-generated-test-is-a-black-box-test.md) |
| Issue 16, a run refuses an unexported enumeration with EID-0081 | <https://github.com/dokimasia/dokimi/issues/16> |
