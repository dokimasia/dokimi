import { Grade, Header, Mode, Region, Status } from "./status";

/**
 * statusToString returns the text of value. A value outside the declared set
 * renders as its own value.
 */
export function statusToString(value: Status): string {
  switch (value) {
    case Status.Draft:
      return "Draft";
    case Status.Published:
      return "Published";
    case Status.Archived:
      return "Archived";
  }
  return `Status(${String(value)})`;
}

/** UnknownStatusError is the error of text that names no declared variant of Status. */
export class UnknownStatusError extends Error {
  /** text is the text that names no declared variant. */
  readonly text: string;

  /** Creates the error of the text that names no declared variant. */
  constructor(text: string) {
    super(`status: unknown Status value ${JSON.stringify(text)}`);
    this.name = "UnknownStatusError";
    this.text = text;
  }
}

/**
 * parseStatus returns the variant of Status whose text is text. It is the
 * inverse of statusToString on the declared variants.
 *
 * @throws UnknownStatusError for text that names no declared variant.
 */
export function parseStatus(text: string): Status {
  switch (text) {
    case "Draft":
      return Status.Draft;
    case "Published":
      return Status.Published;
    case "Archived":
      return Status.Archived;
  }
  throw new UnknownStatusError(text);
}

/**
 * statusValues returns the declared variants in the order of their
 * declaration. Each call returns a new array.
 */
export function statusValues(): Status[] {
  return [
    Status.Draft,
    Status.Published,
    Status.Archived,
  ];
}

/**
 * isStatus reports whether value is a declared variant of Status.
 * A cast can make any other value.
 */
export function isStatus(value: unknown): value is Status {
  switch (value) {
    case Status.Draft:
    case Status.Published:
    case Status.Archived:
      return true;
  }
  return false;
}

/**
 * regionToString returns the text of value. A value outside the declared set
 * renders as its own value.
 */
export function regionToString(value: Region): string {
  switch (value) {
    case Region.US:
      return "us-east";
    case Region.EU:
      return "eu-west";
    case Region.AP:
      return "ap-south";
  }
  return String(value);
}

/** UnknownRegionError is the error of text that names no declared variant of Region. */
export class UnknownRegionError extends Error {
  /** text is the text that names no declared variant. */
  readonly text: string;

  /** Creates the error of the text that names no declared variant. */
  constructor(text: string) {
    super(`status: unknown Region value ${JSON.stringify(text)}`);
    this.name = "UnknownRegionError";
    this.text = text;
  }
}

/**
 * parseRegion returns the variant of Region whose text is text. It is the
 * inverse of regionToString on the declared variants.
 *
 * @throws UnknownRegionError for text that names no declared variant.
 */
export function parseRegion(text: string): Region {
  switch (text) {
    case "us-east":
      return Region.US;
    case "eu-west":
      return Region.EU;
    case "ap-south":
      return Region.AP;
  }
  throw new UnknownRegionError(text);
}

/**
 * regionValues returns the declared variants in the order of their
 * declaration. Each call returns a new array.
 */
export function regionValues(): Region[] {
  return [
    Region.US,
    Region.EU,
    Region.AP,
  ];
}

/**
 * isRegion reports whether value is a declared variant of Region.
 * A cast can make any other value.
 */
export function isRegion(value: unknown): value is Region {
  switch (value) {
    case Region.US:
    case Region.EU:
    case Region.AP:
      return true;
  }
  return false;
}

/**
 * headerToString returns the text of value. A value outside the declared set
 * renders as its own value.
 */
export function headerToString(value: Header): string {
  switch (value) {
    case Header["content-type"]:
      return "ct";
    case Header.Accept:
      return "a";
    case Header.Html:
      return "<html>";
    case Header["1st"]:
      return "first";
  }
  return String(value);
}

/** UnknownHeaderError is the error of text that names no declared variant of Header. */
export class UnknownHeaderError extends Error {
  /** text is the text that names no declared variant. */
  readonly text: string;

  /** Creates the error of the text that names no declared variant. */
  constructor(text: string) {
    super(`status: unknown Header value ${JSON.stringify(text)}`);
    this.name = "UnknownHeaderError";
    this.text = text;
  }
}

/**
 * parseHeader returns the variant of Header whose text is text. It is the
 * inverse of headerToString on the declared variants.
 *
 * @throws UnknownHeaderError for text that names no declared variant.
 */
export function parseHeader(text: string): Header {
  switch (text) {
    case "ct":
      return Header["content-type"];
    case "a":
      return Header.Accept;
    case "<html>":
      return Header.Html;
    case "first":
      return Header["1st"];
  }
  throw new UnknownHeaderError(text);
}

/**
 * headerValues returns the declared variants in the order of their
 * declaration. Each call returns a new array.
 */
export function headerValues(): Header[] {
  return [
    Header["content-type"],
    Header.Accept,
    Header.Html,
    Header["1st"],
  ];
}

/**
 * isHeader reports whether value is a declared variant of Header.
 * A cast can make any other value.
 */
export function isHeader(value: unknown): value is Header {
  switch (value) {
    case Header["content-type"]:
    case Header.Accept:
    case Header.Html:
    case Header["1st"]:
      return true;
  }
  return false;
}

/**
 * gradeToString returns the text of value. A value outside the declared set
 * renders as its own value.
 */
export function gradeToString(value: Grade): string {
  switch (value) {
    case Grade.Poor:
      return "poor";
    case Grade.Fair:
      return "fair enough";
    case Grade.Good:
      return "Good";
  }
  return `Grade(${String(value)})`;
}

/** UnknownGradeError is the error of text that names no declared variant of Grade. */
export class UnknownGradeError extends Error {
  /** text is the text that names no declared variant. */
  readonly text: string;

  /** Creates the error of the text that names no declared variant. */
  constructor(text: string) {
    super(`status: unknown Grade value ${JSON.stringify(text)}`);
    this.name = "UnknownGradeError";
    this.text = text;
  }
}

/**
 * parseGrade returns the variant of Grade whose text is text. It is the
 * inverse of gradeToString on the declared variants.
 *
 * @throws UnknownGradeError for text that names no declared variant.
 */
export function parseGrade(text: string): Grade {
  switch (text) {
    case "poor":
      return Grade.Poor;
    case "fair enough":
      return Grade.Fair;
    case "Good":
      return Grade.Good;
  }
  throw new UnknownGradeError(text);
}

/**
 * gradeValues returns the declared variants in the order of their
 * declaration. Each call returns a new array.
 */
export function gradeValues(): Grade[] {
  return [
    Grade.Poor,
    Grade.Fair,
    Grade.Good,
  ];
}

/**
 * isGrade reports whether value is a declared variant of Grade.
 * A cast can make any other value.
 */
export function isGrade(value: unknown): value is Grade {
  switch (value) {
    case Grade.Poor:
    case Grade.Fair:
    case Grade.Good:
      return true;
  }
  return false;
}

/**
 * isMode reports whether value is a declared variant of Mode.
 * A cast can make any other value.
 */
export function isMode(value: unknown): value is Mode {
  switch (value) {
    case Mode.Off:
    case Mode.On:
      return true;
  }
  return false;
}
