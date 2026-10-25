---
adr: 0007
title: A suite does not ask for the double of its declaration
status: Accepted
date: 2026-10-25
supersedes: RFC-0003 in part
superseded-by: none
rfc: 0003
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0007: A suite does not ask for the double of its declaration

## Status

Accepted

## Context

RFC-0003 gives each interface and each concrete type under the directive `suite` the typed double of the stub generator, because the planted defects of its suite are built from that double. The annotate half of the suite generator stamps the fact `suite.double` on each such declaration, and the stub generator writes a double for the fact.

- The generated suite refers to no double. The planted defect of a ceiling is the type `<Declaration>Over` of the suite itself. The planted defect of a check of the consumer is a function that the consumer writes.
- A value type under `suite` gets a double, and no code calls it. A struct that embeds `id.ID128` of another module gets `WrapperStub` with the methods `ID`, `Compare` and `String` that it promotes, and a spec of each of them.
- In the Go corpus, 134 of the 136 declarations under `suite` also carry the directive `stub`. No hand-written code calls the double of the other two. Each declaration under `suite` in sharder carries `stub` as well.

## Decision

We will write the double of a declaration only where the consumer asks for it with the directive `stub`.

- The suite generator does not stamp the fact of a double. Neither generator declares a capability for the other.
- The stub generator writes the double of each declaration under `stub` and of each handle that such a declaration returns. A declaration under `suite` alone that a method of a declaration under `stub` returns is such a handle.
- A consumer that builds the planted defects of its checks on the double puts the directive `stub` on the declaration.

## Alternatives Considered

### A double for each interface under suite

The suite generator could keep the request for an interface and drop it for a struct. A planted defect of an interface often delegates to the double. The directive `stub` already states that need, and 134 of the 136 declarations of the corpus write it.

## Consequences

**Positive:**

- A value type under `suite` does not get a double or a spec of a double.
- The suite generator and the stub generator have no order between them.

**Negative:**

- A declaration under `suite` alone whose planted defects call its double needs the directive `stub` as well.

## References

| What | Where |
|---|---|
| RFC-0003, which gives each declaration under suite a double | [0003-the-suite-generator.md](../rfc/0003-the-suite-generator.md) |
| Issue 23, the double of a struct that embeds a type of another module | <https://github.com/dokimasia/dokimi/issues/23> |
