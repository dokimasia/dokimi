---
adr: 0022
title: A child error value matches its parent and opens with its message
status: Accepted
date: 2026-11-07
supersedes: RFC-0004 in part
superseded-by: none
rfc: 0004
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0022: A child error value matches its parent and opens with its message

## Status

Accepted

## Context

The generated test of each error value checks that its message opens with the message of no other error value of the package, and that it matches no other error value under `errors.Is`.

A package can declare a family of error values. A caller handles a class of errors by the parent and one cause of the class by a child, whose initializer wraps the parent:

```go
// ErrPolicy is returned for a method whose options break a rule.
var ErrPolicy = errors.New("policy: invalid policy")

// ErrExposure is returned for a method without an exposure.
var ErrExposure = fmt.Errorf("%w: no exposure", ErrPolicy)
```

- Each child matches its parent under `errors.Is`, and its message opens with the message of the parent. Each child fails both checks, so a package with a family cannot apply the directive `sentinel`.
- The rules of eidos state the error value that the initializer of an error value wraps at the opening of its message. The Go rules return the first operand of a call of `fmt.Errorf` whose format opens with `%w`.
- The generator left out a check that an error which wraps an error value matches the value, because such a check tests the standard library. A child is an error value of the package. Its match with its parent is the contract of the family on which a caller relies.

## Decision

We will replace the two checks of a child error value against its parent with checks of their family, because a caller handles a class of errors by the parent.

- The parent of an error value is the error value of the same package that the rules of eidos state for its initializer. An error value without a parent of its package keeps every check.
- The ancestors of a child are its parent, the parent of the parent, and so on. The checks of a child that refuse a match and the opening of a message leave out each ancestor. A child matches each ancestor through its parent, and its message opens with the message of each.
- A child has the checks `matches its parent` and `opens with the message of its parent`. Its check `returns a message of its own` compares its message with the message of its parent as well.
- Every other check stays, also between two children of one parent. The test of a parent compares the parent with each child, so it checks that the parent matches none of its children.

## Alternatives Considered

### The checks of every error value

Each error value kept its checks against every other error value of the package. A package with a family then cannot apply the directive, because each child fails two checks.

### The checks against the parent alone

The test of a child left out its parent alone and kept its checks against every other error value. A child of a child matches the parent of its parent through both wraps, and its message opens with the message of that error value, so it fails both checks against it.

## Consequences

**Positive:**

- A package with a family of error values applies the directive, and its checks still pin the identity of each error value.

**Negative:**

- The check that a child matches its parent fails only where the parent is nil when the child is initialized, or where the rules and the runtime disagree.
- A family through another call, such as `errors.Join` or a function of the package, still fails the checks of its children.
- The test of a child has two more checks than the test of an error value without a parent.

## References

| What | Where |
|---|---|
| RFC-0004, whose checks of an error value this decision opens to a family | [0004-doubles-builders-enumerations-and-errors.md](../rfc/0004-doubles-builders-enumerations-and-errors.md) |
| ADR-0123 of eidos, the error value that an error value wraps | <https://github.com/dokimasia/eidos/blob/main/docs/adr/0123-the-rules-state-the-error-value-that-an-error-value-wraps.md> |
| Issue 62, the checks refuse an error value that wraps another error value of its package | <https://github.com/dokimasia/dokimi/issues/62> |
| The documentation of `errors.Is` | <https://pkg.go.dev/errors#Is> |
