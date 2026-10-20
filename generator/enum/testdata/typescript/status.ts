/** Represents the state of an order, whose first member is 1. */
//dokimi:enum
export enum Status {
  Draft = 1,
  Published,
  Archived,
}

/** Represents a region of a cloud, whose text is its value. */
//dokimi:enum
export enum Region {
  US = "us-east",
  EU = "eu-west",
  AP = "ap-south",
}

/** Represents a header of a request with a member whose name is no identifier. */
//dokimi:enum
export enum Header {
  "content-type" = "ct",
  Accept = "a",
  Html = "<html>",
  "1st" = "first",
}

/**
 * Represents the grade of a review.
 *
 * @remarks
 * The directive takes the text of a member from the trailing comment of the
 * member. Good has no comment, so its text is its name.
 */
//dokimi:enum text=comment
export enum Grade {
  Poor, // poor
  Fair, //  fair enough
  Good,
}

/** Represents a mode of a lamp with its own render. */
//dokimi:enum
export enum Mode {
  Off,
  On,
}

/** Returns the name of a mode. */
export function modeToString(mode: Mode): string {
  return mode === Mode.On ? "on" : "off";
}
