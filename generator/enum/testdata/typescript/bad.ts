/** Represents an enumeration without a member. */
//dokimi:enum
export enum Nothing {}

/** Represents a shade, whose two members have one value. */
//dokimi:enum
export enum Shade {
  Dark = 1,
  Black = 1,
}

/** Represents a tint, whose module declares a parse function of a number. */
//dokimi:enum
export enum Tint {
  Pale = 1,
  Deep,
}

/** Returns the tint of a number. */
export function parseTint(value: number): Tint {
  return value === 2 ? Tint.Deep : Tint.Pale;
}

/** Represents a hue, whose module declares a function under the name of the error of its parse function. */
//dokimi:enum
export enum Hue {
  Warm = 1,
  Cold,
}

/** Returns an error of a hue. */
export function UnknownHueError(): Error {
  return new Error("hue");
}
