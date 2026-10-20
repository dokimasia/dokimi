/**
 * Declares the errors whose parameter cause comes before a required
 * parameter.
 *
 * @packageDocumentation
 */
//dokimi:sentinel

/** Reports a failure with a label, which the field cause contains. */
export class LabelError extends Error {
  /** Contains the label of the failure. */
  override readonly cause: string;
  /** Contains the operation that failed. */
  readonly op: string;

  constructor(cause: string, op: string) {
    super(`${op}: ${cause}`);
    this.name = "LabelError";
    this.cause = cause;
    this.op = op;
  }
}
