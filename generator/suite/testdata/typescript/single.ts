/** Represents a probe of one callable, whose planted defect that panics reads no fresh subject. */
//dokimi:suite
export interface Single {
  /** Applies the probe. */
  apply(signal: AbortSignal): Promise<void>;
}

/** Represents a state of one callable, whose planted defect that reports an error reads no fresh subject. */
//dokimi:suite
export interface Resettable {
  /** Clears the state. */
  reset(): void;
}
