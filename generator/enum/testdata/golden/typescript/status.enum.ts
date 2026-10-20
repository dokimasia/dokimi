import { Grade, Header, Mode, Region, Status } from './status';

/**
 * statusToString returns the textual form of a Status.
 *
 * A value outside the declared set renders as itself, so a corrupt value is
 * visible rather than indistinguishable from a declared one.
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

/** UnknownStatusError reports text that names no declared variant of Status. */
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
 * parseStatus maps text back to a Status.
 *
 * The inverse of statusToString for every declared variant, which is what
 * lets a value survive a round trip through anything that carries it as
 * text — a log line, a JSON field, a database column.
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
 * statusValues returns every declared variant in declaration order.
 *
 * A fresh array per call: a module-level one would let any caller reorder or
 * overwrite the set every other caller reads.
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
 *
 * A cast admits any value, so this is the only thing standing between an
 * out-of-range value and the rest of the program.
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
 * regionToString returns the textual form of a Region.
 *
 * A value outside the declared set renders as itself, so a corrupt value is
 * visible rather than indistinguishable from a declared one.
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

/** UnknownRegionError reports text that names no declared variant of Region. */
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
 * parseRegion maps text back to a Region.
 *
 * The inverse of regionToString for every declared variant, which is what
 * lets a value survive a round trip through anything that carries it as
 * text — a log line, a JSON field, a database column.
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
 * regionValues returns every declared variant in declaration order.
 *
 * A fresh array per call: a module-level one would let any caller reorder or
 * overwrite the set every other caller reads.
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
 *
 * A cast admits any value, so this is the only thing standing between an
 * out-of-range value and the rest of the program.
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
 * headerToString returns the textual form of a Header.
 *
 * A value outside the declared set renders as itself, so a corrupt value is
 * visible rather than indistinguishable from a declared one.
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

/** UnknownHeaderError reports text that names no declared variant of Header. */
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
 * parseHeader maps text back to a Header.
 *
 * The inverse of headerToString for every declared variant, which is what
 * lets a value survive a round trip through anything that carries it as
 * text — a log line, a JSON field, a database column.
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
 * headerValues returns every declared variant in declaration order.
 *
 * A fresh array per call: a module-level one would let any caller reorder or
 * overwrite the set every other caller reads.
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
 *
 * A cast admits any value, so this is the only thing standing between an
 * out-of-range value and the rest of the program.
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
 * gradeToString returns the textual form of a Grade.
 *
 * A value outside the declared set renders as itself, so a corrupt value is
 * visible rather than indistinguishable from a declared one.
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

/** UnknownGradeError reports text that names no declared variant of Grade. */
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
 * parseGrade maps text back to a Grade.
 *
 * The inverse of gradeToString for every declared variant, which is what
 * lets a value survive a round trip through anything that carries it as
 * text — a log line, a JSON field, a database column.
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
 * gradeValues returns every declared variant in declaration order.
 *
 * A fresh array per call: a module-level one would let any caller reorder or
 * overwrite the set every other caller reads.
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
 *
 * A cast admits any value, so this is the only thing standing between an
 * out-of-range value and the rest of the program.
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
 *
 * A cast admits any value, so this is the only thing standing between an
 * out-of-range value and the rest of the program.
 */
export function isMode(value: unknown): value is Mode {
  switch (value) {
    case Mode.Off:
    case Mode.On:
      return true;
  }
  return false;
}
