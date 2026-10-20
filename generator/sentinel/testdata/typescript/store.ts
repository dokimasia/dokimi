/**
 * Declares the errors of the store, whose messages open with kv.
 *
 * @packageDocumentation
 */
//dokimi:sentinel prefix=kv

/** Reports a key that was never written. */
export class NotFoundError extends Error {
  /** Contains the missing key. */
  readonly key: string;

  constructor(key: string) {
    super(`kv: key ${key} not found`);
    this.name = "NotFoundError";
    this.key = key;
  }
}

/**
 * Reports a key that a delete removed.
 *
 * @remarks
 * The class constructs through the constructor of {@link NotFoundError}.
 */
export class MissingError extends NotFoundError {
  override name = "MissingError";
}

/** Reports a failed operation with its cause. */
export class WrappedError extends Error {
  /** Contains the operation that failed. */
  readonly op: string;
  /** Contains the cause of the failure. */
  override readonly cause: unknown;

  constructor(op: string, cause?: unknown) {
    super(`kv: ${op} failed`);
    this.name = "WrappedError";
    this.op = op;
    this.cause = cause;
  }
}

/** Reports a use after the close. */
export class ClosedError extends Error {
  constructor() {
    super("kv: closed");
    this.name = "ClosedError";
  }
}
