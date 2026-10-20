/**
 * Declares error classes that a check cannot construct, beside three that it
 * can.
 *
 * @packageDocumentation
 */
//dokimi:sentinel

/**
 * Reports a loss.
 *
 * @remarks
 * The class has the constructor of Error, which a check calls.
 */
export class LostError extends Error {}

/**
 * Reports a sealed resource.
 *
 * @remarks
 * The constructor is private.
 */
export class SealedError extends Error {
  private constructor() {
    super("sealed");
  }

  /** Returns a new {@link SealedError}. */
  static make(): SealedError {
    return new SealedError();
  }
}

/**
 * Reports the text of a handler.
 *
 * @remarks
 * The constructor takes a handler, which has no sample.
 */
export class HandlerError extends Error {
  constructor(handler: () => string) {
    super(handler());
  }
}

/**
 * Reports an operation with a code.
 *
 * @remarks
 * The constructor has a parameter with a default and without a sample before
 * a required parameter.
 */
export class DefaultedError extends Error {
  constructor(code: () => number = () => 0, op: string) {
    super(`${op} ${String(code())}`);
  }
}

/**
 * Reports a kept resource.
 *
 * @remarks
 * The class has a public constructor without parameters.
 */
export class KeptError extends Error {
  constructor() {
    super("kept");
  }
}

/**
 * Reports a failure of a chain.
 *
 * @remarks
 * The constructor takes the cause in its first parameter.
 */
export class ChainError extends Error {
  /** Contains the cause of the failure. */
  override readonly cause: unknown;

  constructor(cause?: unknown) {
    super("chain");
    this.cause = cause;
  }
}
