// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the errors-kind fixture for an error contract reached through
 * inheritance rather than declared outright.
 *
 * @remarks
 * A class's declarations are not its members. `class NotFoundError extends
 * BaseError {}` is the dominant TypeScript idiom for a family of custom
 * errors, and a generator reading only the declarations finds no field and
 * no constructor on it at all — so the module's directive says its errors
 * are a contract and the generated file covers half of them, without
 * reporting that it did.
 *
 * Every half is reached the same way here: the message, the cause and the
 * constructor are all inherited, so a generator reading declarations sees an
 * empty class and emits nothing about it.
 *
 * The subclass adds no field of its own, and that is deliberate rather than
 * minimal. An inherited message cannot mention a field the subclass declared,
 * so a class that both inherits its contract and adds a field is asserting
 * two things at once — and the generated "every field reaches the message"
 * check would fail for a reason that has nothing to do with inheritance.
 *
 * The store fixture is the same contract declared directly, which is what
 * makes the pair able to tell "inheritance is handled" from "the class
 * happened to declare everything".
 *
 * @packageDocumentation
 */
//dokimi:sentinel

/** Reports a missing key, as the cause that the inherited contract wraps. */
export const notFound = new Error("embedded: not found");

/**
 * Represents the whole contract, which subclasses inherit rather than use.
 *
 * @remarks
 * The constructor names each error after the class that constructs it, so a
 * subclass that declares nothing still reports its own name.
 */
export class BaseError extends Error {
  /** Contains the operation that failed. */
  readonly op: string;
  /** Contains the cause of the failure. */
  override readonly cause: unknown;

  constructor(op: string, cause?: unknown) {
    super(`embedded: ${op}: ${String(cause)}`);
    this.name = new.target.name;
    this.op = op;
    this.cause = cause;
  }
}

/**
 * Reports a missing key through the inherited contract.
 *
 * @remarks
 * NotFoundError declares nothing and inherits the entire error contract. It
 * is a distinct class that `instanceof` can single out, which is what the
 * idiom buys, and nothing about it is visible to a generator reading
 * declarations alone.
 */
export class NotFoundError extends BaseError {}
