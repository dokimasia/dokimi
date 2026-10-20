/**
 * Represents a sender with a parameter named result.
 *
 * @remarks
 * result is also the field of the result of the call.
 */
//dokimi:stub
export interface Clashed {
  /** Sends the result, and returns the count of the receivers. */
  send(result: string): number;
}

/**
 * Represents a failure with a parameter named error.
 *
 * @remarks
 * error is also the field of the error of the call.
 */
//dokimi:stub
export interface Failed {
  /** Fails with the error. */
  fail(error: Error): void;
}
