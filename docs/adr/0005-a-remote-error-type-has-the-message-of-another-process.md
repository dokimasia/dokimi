---
adr: 0005
title: A remote error type has the message of another process
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

# ADR-0005: A remote error type has the message of another process

## Status

Accepted

## Context

RFC-0004 checks that each error type of a package under the directive `sentinel` returns a message for its zero value, and that the message opens with the prefix of the package. An error type that carries the error of another process returns the message of that error unchanged, so it fails both checks.

- Issue 19 reports `protocol.CallError` of sharder. A transport sends it in place of the result of a call, with the class and the message of the error of the server.
- A prefix of `protocol:` would put a second prefix in front of the message of the server, as in `protocol: control: the call names a state that has changed`.
- The directive `sentinel` adjusts its checks for a package through its params `prefix` and `no-overlap-with`.

## Decision

We will add the param `remote` to the directive `sentinel`. It lists the remote error types of the package, the types whose message is the message of an error of another process, unchanged.

- The checks of a remote error type leave out the message of the zero value and the prefix of the message.
- They keep the checks that the type puts each field of text into its message and matches no other error type of the package.
- A type under `remote` that is no public error type of the package reports `UnknownRemote` at the package, and the checks of the package keep every check of their error types.
- The param applies to the error types of Go and the error classes of TypeScript.

The name follows the term remote error of RPC frameworks, an error that a server returns and a client receives.

## Alternatives Considered

### A directive on the error type

A directive such as `//dokimi:remote` on the type would keep the mark beside the type. The params of `sentinel` already adjust its checks for the package in one place, and a second directive would need a schema of its own.

### A message of the zero value that the type writes

The type could return a message of its own for an empty message, such as `protocol: remote error`. The message of a remote error is then wrong for a caller that prints it.

## Consequences

**Positive:**

- A package can put its wire carriers of remote errors under the directive `sentinel`.

**Negative:**

- A remote error type has no check that its message is set. The transport of the package tests that.

## References

| What | Where |
|---|---|
| RFC-0004, whose checks of the message of every error type this decision narrows | [0004-doubles-builders-enumerations-and-errors.md](../rfc/0004-doubles-builders-enumerations-and-errors.md) |
| Issue 19, an error type that carries the message of another process fails the checks of its message | <https://github.com/dokimasia/dokimi/issues/19> |
