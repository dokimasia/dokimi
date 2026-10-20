/**
 * Declares the errors of the catalog of the web.
 *
 * @packageDocumentation
 */
//dokimi:sentinel

/** Reports an account that does not exist. */
export class MissingAccountError extends Error {
  /** Contains the owner of the account. */
  readonly owner: string;

  constructor(owner: string) {
    super(`no account of ${owner}`);
    this.name = "MissingAccountError";
    this.owner = owner;
  }
}

/** Reports a failed request with its cause. */
export class RequestError extends Error {
  /** Contains the cause of the failure. */
  override readonly cause: unknown;

  constructor(cause?: unknown) {
    super("request failed");
    this.name = "RequestError";
    this.cause = cause;
  }
}
