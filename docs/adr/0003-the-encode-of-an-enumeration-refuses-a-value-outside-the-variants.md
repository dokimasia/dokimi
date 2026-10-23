---
adr: 0003
title: The encode of an enumeration refuses a value outside the variants
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

# ADR-0003: The encode of an enumeration refuses a value outside the variants

## Status

Accepted

## Context

RFC-0004 gives the surface of a Go enumeration an encode that writes the text of the render for every value. A conversion from the underlying type makes a value outside the variants, such as `Phase(0)`. The encode writes `"Phase(0)"`, and the decode of the surface refuses that text with `ErrUnknownPhase`.

- Issue 17 reports eight enumerations of sharder whose kanon suites fail. kanontest checks that a type accepts a value for its text codec exactly when its validation accepts it.
- The hand-written `MarshalText` of four of those enumerations returned an error for such a value before they moved to the generator.
- `time.Time.MarshalText` returns an error for a year that RFC 3339, the format of its `UnmarshalText`, cannot write.

## Decision

We will write an encode that refuses each value that the validity method refuses.

- The encode returns an error that wraps the error of the parse function for such a value, as in `fmt.Errorf("%w %q", ErrUnknownPhase, v.String())`. The error contains the text of the render.
- Where the package declares a parse function and no error of it, the encode returns an error with the message of the error that the surface would write, such as `declared: unknown Lane value "Lane(9)"`.
- The checks of the encode call it with a value outside the variants, or with the zero value where no variant has it, and check the error.

Each text that the encode writes then decodes with the decode of the surface.

## Alternatives Considered

### The text of the render for every value

The encode and the decode would disagree on a value outside the variants. A program that writes such a value as JSON would write a document that it cannot read back.

## Consequences

**Positive:**

- The text codec of a generated surface and its validity method accept the same values, as kanontest checks.

**Negative:**

- A program that wrote a value outside the variants as text gets an error from the encode instead.

## References

| What | Where |
|---|---|
| RFC-0004, whose encode of the text of every value this decision replaces | [0004-doubles-builders-enumerations-and-errors.md](../rfc/0004-doubles-builders-enumerations-and-errors.md) |
| Issue 17, MarshalText encodes a value that IsValid refuses | <https://github.com/dokimasia/dokimi/issues/17> |
