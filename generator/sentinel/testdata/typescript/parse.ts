/**
 * Declares the errors of the parser, whose messages open without a prefix.
 *
 * @packageDocumentation
 */
//dokimi:sentinel

/** Reports a text that does not parse. */
export class ParseError extends Error {
  /** Contains the text that does not parse. */
  readonly text: string;
  /** Contains the line of the fault. */
  readonly line: number;

  constructor(text: string, line: number, ...notes: string[]) {
    super(`cannot parse ${text} at line ${line}${notes.join("")}`);
    this.name = "ParseError";
    this.text = text;
    this.line = line;
  }
}

/** Reports a text that does not split into tokens. */
export class LexError extends Error {
  constructor(text: string) {
    super(`cannot split ${text}`);
    this.name = "LexError";
  }
}
