/**
 * Declares the errors of a relay to a remote store, which carries the message
 * of the remote store unchanged.
 *
 * @packageDocumentation
 */
//dokimi:sentinel prefix=relay remote=[CallError]

/** Carries the error of the remote store in place of the result of a call. */
export class CallError extends Error {
  /** Contains the class of the remote error. */
  readonly code: number;

  constructor(message: string, code: number) {
    super(message);
    this.name = "CallError";
    this.code = code;
  }
}

/** Reports a call that the remote store did not answer in time. */
export class TimeoutError extends Error {
  /** Contains the operation that timed out. */
  readonly op: string;

  constructor(op: string) {
    super(`relay: ${op} timed out`);
    this.name = "TimeoutError";
    this.op = op;
  }
}
