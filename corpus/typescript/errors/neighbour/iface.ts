// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the target of the non-overlap declaration of another module.
 *
 * @remarks
 * The `no-overlap-with` param of `sentinel` names a second module whose error
 * values must not match this one's. That check needs two modules by
 * construction, so a corpus with a single errors fixture leaves the directive
 * unexercised — it would parse, name nothing reachable, and pass.
 *
 * The rules of TypeScript name no error values, so the check compares nothing
 * here. The fixture pins that the param resolves the module: a path that the
 * run did not load reports Unloaded.
 *
 * The errors here deliberately read like a store's. Distinct errors that
 * describe the same conditions are exactly the case where a copy-paste
 * between modules produces an accidental alias, which is what the
 * cross-module check exists to catch.
 *
 * @packageDocumentation
 */
//dokimi:sentinel

/**
 * Reports a missing key.
 *
 * @remarks
 * Its message differs from the store module's by prefix alone, which is the
 * near-miss the check must accept as distinct.
 */
export class NotFoundError extends Error {
  constructor() {
    super("neighbour: not found");
    this.name = "NotFoundError";
  }
}

/** Reports a losing write. */
export class ConflictError extends Error {
  constructor() {
    super("neighbour: conflict");
    this.name = "ConflictError";
  }
}
