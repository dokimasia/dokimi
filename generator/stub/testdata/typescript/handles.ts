/**
 * Represents a pinger with a property and a method without a declared result
 * type.
 *
 * @remarks
 * The pinger has no double. It sorts before every other producer of the
 * module.
 */
//dokimi:stub
export interface Pinger {
  /** Contains the identifier of the pinger. */
  readonly id: string;
  /** Sends a ping. */
  ping();
}

/** Represents a producer that opens cursors, a handle without a directive. */
//dokimi:stub
export interface Producer {
  /** Returns a new cursor. */
  open(): Promise<Cursor>;
}

/** Represents a handle that {@link Producer.open} returns. */
export interface Cursor {
  /** Returns the next value. */
  next(): Promise<string>;
  /** Releases the cursor. */
  close(): void;
}
