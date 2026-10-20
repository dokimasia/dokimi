/** Represents a point with properties and no methods, so its double declares fields alone. */
//dokimi:stub
export interface Point {
  /** Contains the distance from the vertical axis. */
  readonly x: number;
  /** Contains the distance from the horizontal axis. */
  readonly y: number;
}

/**
 * Represents a box with a property of its type parameter alone.
 *
 * @remarks
 * No property of its double has a sample.
 */
//dokimi:stub
export interface Box<T> {
  /** Contains the value in the box. */
  value: T;
}
