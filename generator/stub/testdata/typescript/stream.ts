/** Represents a feed that returns its contents as sequences. */
//dokimi:stub
export interface Feed {
  /** Returns a sequence of values. */
  names(): Iterable<string>;
  /** Returns a sequence that arrives asynchronously. */
  rows(): AsyncIterable<number>;
}
