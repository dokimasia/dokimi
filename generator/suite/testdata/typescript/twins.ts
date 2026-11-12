/** Represents a store of keys on the left. */
//dokimi:suite
export interface Left {
  /** Writes a key. */
  put(key: string): void;
}

/** Represents a store of keys on the right. */
//dokimi:suite
export interface Right {
  /** Writes a key. */
  put(key: string): void;
}
