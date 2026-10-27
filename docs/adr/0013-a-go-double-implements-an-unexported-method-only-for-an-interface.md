---
adr: 0013
title: A Go double implements an unexported method only for an interface, and its checks leave it out
status: Accepted
date: 2026-10-27
supersedes: RFC-0004 in part
superseded-by: none
rfc: 0004
---

<!--
  ~ Copyright Dokimasia B.V. 2026
  ~ SPDX-License-Identifier: Apache-2.0
-->

# ADR-0013: A Go double implements an unexported method only for an interface, and its checks leave it out

## Status

Accepted

## Context

RFC-0004 gives a double each method of the member walk of its subject, and the checks of the double call each method. The member walk of a Go declaration contains its unexported methods. Go refuses a call of such a method from another package, and no type of another package implements it:

- A struct under an `out` to another package had no double where an unexported method referred to an unexported type, such as `add(s step)`. The stub generator reported `Unimplementable`. Issue 37 reports it.
- The double of a struct in its own package declared each unexported method, such as `reset(n int)`. The checks of the double are in the external test package, and their call of `s.reset` did not compile.
- The checks of the double of an interface with an unexported method, such as `seal()`, did not compile either. Under an `out` to another package, the double did not implement the interface.
- The option `DelegateTo` of the double of a struct takes the interface of the methods of the double. A pointer to the struct did not implement that interface where the double was in another package and declared an unexported method.

## Decision

We will implement an unexported method only in the double of an interface, and leave each unexported method out of the checks.

- The double of a struct leaves out each unexported method without a check of its signature. A struct without an exported method reports `Empty`, because its double has no method.
- The double of an interface implements each unexported method, because the double is a value of the interface. The checks of the double leave out the unexported methods.
- The stub generator reports `Unimplementable` at an interface with an unexported method where the double is in another package than the method. An `out` to another package is such a case, and so is an interface that embeds an unexported method of another package.

This decision supersedes the sentences of RFC-0004 that give a double and its checks each method of the member walk, for a Go declaration.

## Alternatives Considered

### Unimplementable at each subject with an unexported method

The stub generator refuses a TypeScript class with a private member, because the double of a class is a value of the class. A Go double of a struct is not a value of the struct, so it does not need the unexported methods. The double of an interface can implement them in the package of the interface. A refusal would leave each such declaration without a double.

### The checks in the package of the subject

The checks could call the unexported methods from an internal test file. ADR-0001 refuses a generated white-box test, and ADR-0004 does not write checks of an unexported enumeration for the same reason.

## Consequences

**Positive:**

- A struct with unexported methods has a double in every package, and a pointer to the struct implements the interface that `DelegateTo` takes.
- An interface with an unexported method has a double in its own package. A test passes the double to the code of that package.

**Negative:**

- A test cannot configure or record an unexported method of a struct through its double.
- The generated checks do not cover the unexported methods of the double of an interface.
- An interface with an unexported method has no double in another package.

## References

| What | Where |
|---|---|
| RFC-0004, whose methods of a double and of its checks this decision narrows for Go | [0004-doubles-builders-enumerations-and-errors.md](../rfc/0004-doubles-builders-enumerations-and-errors.md) |
| ADR-0001, every generated test is a black-box test | [0001-every-generated-test-is-a-black-box-test.md](0001-every-generated-test-is-a-black-box-test.md) |
| ADR-0004, an unexported enumeration has unexported members and no checks | [0004-an-unexported-enumeration-has-unexported-members-and-no-checks.md](0004-an-unexported-enumeration-has-unexported-members-and-no-checks.md) |
| Issue 37, a struct with an unexported method gets no double in another package | <https://github.com/dokimasia/dokimi/issues/37> |
