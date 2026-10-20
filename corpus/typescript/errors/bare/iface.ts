// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the minimal half of the errors-kind corpus: error classes without
 * fields and without a cause, and nothing else.
 *
 * @remarks
 * The generator emits checks that cover any set of error classes: that every
 * message is non-empty and carries the module prefix, and that no class
 * matches a sibling. It then appends further checks for each field and cause
 * it detects. With neither here there is no field to find in a message and no
 * cause to expose — so this fixture is the floor, and the store fixture is
 * the ceiling.
 *
 * The rules of TypeScript name no error values, so the TypeScript floor is a
 * set of classes rather than the module-level sentinels of the Go corpus.
 *
 * A corpus holding only the rich case cannot tell "the optional checks were
 * correctly omitted" from "the optional checks were silently dropped".
 *
 * @packageDocumentation
 */
//dokimi:sentinel

/** Reports an operation on an empty collection. */
export class EmptyError extends Error {
  constructor() {
    super("bare: empty");
    this.name = "EmptyError";
  }
}

/** Reports an operation on a saturated collection. */
export class FullError extends Error {
  constructor() {
    super("bare: full");
    this.name = "FullError";
  }
}

/** Reports input that the module rejects. */
export class InvalidError extends Error {
  constructor() {
    super("bare: invalid");
    this.name = "InvalidError";
  }
}
