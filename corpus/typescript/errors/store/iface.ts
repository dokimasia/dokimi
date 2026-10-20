// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the errors-kind fixture: the input the sentinel generator reads.
 *
 * @remarks
 * The directive sits in the header of the module rather than on the classes.
 * sentinel addresses a package, which a TypeScript module is, so what a
 * generator *reads* and where it is opted in are different things — a fixture
 * carrying a class-level directive would never be picked up.
 *
 * The prefix deliberately differs from the module name. An override set to
 * the value the derivation would have produced anyway cannot tell "the
 * override was honoured" from "it was ignored", which is the standard the
 * noprefix fixture beside this one states for the suppression branch. A
 * generator that read prefix= and dropped it derives `iface` here, and every
 * message below then fails the assertion it emits.
 *
 * @packageDocumentation
 */
//dokimi:sentinel prefix=kv no-overlap-with=[errors/neighbour/iface]

/**
 * Reports a key that was never written.
 *
 * @remarks
 * The generated check constructs the error and finds the key in its message,
 * so the key has to survive the message.
 */
export class NotFoundError extends Error {
  /** Contains the key that was never written. */
  readonly key: string;

  constructor(key: string) {
    super(`kv: key ${JSON.stringify(key)} not found`);
    this.name = "NotFoundError";
    this.key = key;
  }
}

/** Reports a write that lost a race. */
export class ConflictError extends Error {
  constructor() {
    super("kv: conflict");
    this.name = "ConflictError";
  }
}

/** Reports a use after teardown. */
export class ClosedError extends Error {
  constructor() {
    super("kv: closed");
    this.name = "ClosedError";
  }
}

/**
 * Reports a failed operation with its cause.
 *
 * @remarks
 * The cause gives the generated cause check something to walk.
 */
export class WrappedError extends Error {
  /** Contains the operation that failed. */
  readonly op: string;
  /** Contains the cause of the failure. */
  override readonly cause: unknown;

  constructor(op: string, cause?: unknown) {
    super(`kv: ${op}: ${String(cause)}`);
    this.name = "WrappedError";
    this.op = op;
    this.cause = cause;
  }
}
