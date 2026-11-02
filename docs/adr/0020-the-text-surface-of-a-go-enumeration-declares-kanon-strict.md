---
adr: 0020
title: The text surface of a Go enumeration declares kanon.Strict
status: Accepted
date: 2026-11-02
supersedes: RFC-0004 in part
superseded-by: none
rfc: 0004
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0020: The text surface of a Go enumeration declares kanon.Strict

## Status

Accepted

## Context

The surface of a Go enumeration has a text codec: `MarshalText` writes the text of a variant, and `UnmarshalText` reads it back through the parse function. The parse function of the surface accepts the texts of the variants alone. No two variants share a value or a text, and the encode refuses each value that the validity method refuses.

kanon encodes a type through its text methods. A canonical decode of such a type encodes each decoded value again and compares the bytes with its input. A type that declares the marker `StrictKanon()` of `kanon.Strict` states that its decode accepts a text only where its encode writes that text for the decoded value. kanon v0.4.0 then leaves out the second encode.

- Issue 50 reports that the comparison never fails for a generated surface. Its statement in the code that kanon generates runs in no spec, so a module that requires full statement coverage of its generated code fails.
- The issue measured the decode of a canonical record of a string, an enumeration and a `uint32`: 22.80 ns and 2 allocations without the marker, and 10.93 ns and 1 allocation with it.

## Decision

We will write the marker `StrictKanon` of `kanon.Strict` on a Go enumeration whose surface writes the encode, the decode, the parse function and the validity method, because the decode of that surface accepts only the texts that its encode writes.

- The marker is `func (Status) StrictKanon() {}`. It has no receiver name and no body, and the generated file does not import kanon.
- A declared member breaks the guarantee where it differs from the member of the surface. A declared parse function can accept a text in another case, and a declared validity method can refuse a variant whose text the decode accepts. The surface then does not write the marker.
- A method `StrictKanon` of the enumeration is a declared member, so the surface does not write a second one. A method of another form under that name conflicts with the surface, and the generator reports `Conflict`.
- The checks do not call the marker. The checks of the encode and the decode check the round trip of each variant and the refusal of a text that no variant has. `kanontest.RunStrict` checks the guarantee in a package that encodes the enumeration with kanon.
- The surface of TypeScript has no marker, because it has no text codec.

## Alternatives Considered

### The marker beside every pair of text methods

Issue 50 asked for the marker beside `MarshalText` and `UnmarshalText` of each enumeration. A surface that decodes through a declared parse function would then state a guarantee that the parse function can break. If that parse function accepts `READY` and `ready` for one variant, the decode accepts a text that the encode does not write.

### kanon.Exact

`kanon.Exact` states the same guarantee of the decode, and also that the append method fails for the zero value alone. The encode of an enumeration fails for every value that the validity method refuses, so an enumeration cannot declare it.

## Consequences

**Positive:**

- A canonical decode of a generated enumeration encodes no decoded value again, and the code that kanon generates has no statement for the comparison.

**Negative:**

- An enumeration with a declared parse function, validity method, encode or decode keeps the second encode, unless its author declares the marker.
- The template writes the name of a method of kanon without an import of kanon. A rename of `kanon.Strict` breaks the guarantee without a compile error, and the template has to change with it.
- Each generated Go enumeration has one more exported method.

## References

| What | Where |
|---|---|
| RFC-0004, whose surface this decision extends with a marker | [0004-doubles-builders-enumerations-and-errors.md](../rfc/0004-doubles-builders-enumerations-and-errors.md) |
| Issue 50, the generator does not write StrictKanon beside the text methods | <https://github.com/dokimasia/dokimi/issues/50> |
| `kanon.Strict` and its guarantee | <https://github.com/thesm-os/kanon/blob/main/message.go> |
