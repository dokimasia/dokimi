---
adr: 0010
title: A suite refers only to declarations that its file can import
status: Accepted
date: 2026-10-26
supersedes: RFC-0006 in part
superseded-by: none
rfc: 0006
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0010: A suite refers only to declarations that its file can import

## Status

Accepted

## Context

RFC-0006 gives each parameter of a callable a value by its kind, and the rules apply to a type at every depth. ADR-0008 checks each reference that a generator writes against the package of its routed file, and the suite generator checks the signature of each callable. A reference at some depth of an input still failed the template of the whole file under EID-0023, with the suites of every other declaration of the file:

- The generated code of a type can refer to a declaration that the file cannot import, such as the members of an enumeration of an internal package of another tree, or the signatures of the methods of an interface.
- A starting value or a declared error can refer to such a declaration, such as the default `deep.High` of a field or the error `deep.ErrWorn` of the directive `fault`. eidos does not check the references of a value for a plugin, as dokimasia/eidos#85 reports.

## Decision

We will write only references that the file of the suite can import, and take the next source of a value where the file cannot import the first.

- A type whose generated code refers to a declaration that the file cannot import takes the generator of its registration, as an interface of the standard library does under RFC-0006. The generated code of a type is its override, the implementation of an interface or the generated function of a function type.
- A starting value takes the first source in the order of RFC-0006 whose references the file can import: the authored sample, the value of the builder, then the declared defaults. A field default that the file cannot import leaves the field at its zero value.
- A declared error takes the first source whose references the file can import: the errors of the directive `fault`, the error value of a mixin, then `conform.ErrGenerated`.

The first decision applies since dokimi v0.2.0. The other two wait for eidos#85.

## Alternatives Considered

### Unimportable and no suite

The suite could report `Unimportable` at the declaration, as ADR-0008 does for a callable whose signature the file cannot import. A registration of the test process can state the generator of such a type, so the suite works with it. One field default would then cost the suite of its declaration.

### A check of values in dokimi

dokimi could walk each value and check each callee with `Bound.Importable`. The types of a value are `emit.TypeRef`, and `Bound.Importable` takes a `node.TypeRef`. dokimi would also copy the check that the backend runs when it spells a value.

## Consequences

**Positive:**

- A suite whose inputs refer to an internal package of another tree has a file, and so do the other suites of that file.

**Negative:**

- A type that the file cannot import takes its shape or a registration. An enumeration of an internal package of another tree takes any integer.
- A starting value can differ from the declared one where the file cannot import the declared one.

## References

| What | Where |
|---|---|
| RFC-0006, the values of the parameters of a callable | [0006-generated-inputs-fuzz-targets-benchmarks-and-budgets.md](../rfc/0006-generated-inputs-fuzz-targets-benchmarks-and-budgets.md) |
| ADR-0008, the check of the references of a generator | [0008-a-generator-checks-its-references-against-the-package-of-its-routed-file.md](0008-a-generator-checks-its-references-against-the-package-of-its-routed-file.md) |
| Issue 25, a field of an interface type in a request | <https://github.com/dokimasia/dokimi/issues/25> |
| eidos#85, the references of a value | <https://github.com/dokimasia/eidos/issues/85> |
