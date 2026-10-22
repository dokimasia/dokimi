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

/** Represents the level of a log record, whose module declares the error of its parse function. */
//dokimi:enum
export enum Level {
  Debug = 1,
  Info,
  Warn,
}

/** Represents the phase of a job, whose module declares its type guard. */
//dokimi:enum
export enum Phase {
  Queued,
  Running,
}

/** Reports whether a value is a phase. */
export function isPhase(value: unknown): value is Phase {
  return value === Phase.Queued || value === Phase.Running;
}

/** Is the error of a text of no level. */
export class UnknownLevelError extends Error {
  /** text is the text of no level. */
  readonly text: string;

  /** Creates the error of the text of no level. */
  constructor(text: string) {
    super(`status: unknown Level value ${JSON.stringify(text)}`);
    this.name = "UnknownLevelError";
    this.text = text;
  }
}
