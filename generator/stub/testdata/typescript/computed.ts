/** Represents a record whose method has a computed name. */
//dokimi:stub
export interface Computed {
  /** Returns an iterator over the numbers of the record. */
  [Symbol.iterator](): Iterator<number>;
}

/** Represents a record whose property has a computed name. */
//dokimi:stub
export interface Tagged {
  /** Contains the tag of the record. */
  [Symbol.toStringTag]: string;
}
