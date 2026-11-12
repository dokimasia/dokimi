/** Represents a record with the revision that its last write stamped. */
export interface Rec {
  /** Contains the key of the record. */
  key: string;
  /** Contains the revision that the last write of the record stamped. */
  rev: number;
}

/** Represents a store whose reads never go back to an older revision. */
//dokimi:suite
export interface Versions {
  /** Returns the record under the key. */
  //dokimi:mixin monotonic-reads version=rev
  version(key: string, signal: AbortSignal): Promise<Rec>;
  /** Stores a record and stamps its revision. */
  put(r: Rec, signal: AbortSignal): Promise<void>;
  /** Stores a record in place of the record under its key and stamps its revision. */
  replace(r: Rec, signal: AbortSignal): Promise<void>;
}
