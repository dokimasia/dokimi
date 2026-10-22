import { describe, expect, it } from "vitest";
import { LexError, ParseError } from "./parse";

function lexErrorMessageCases(): void {
  it("is not empty", () => {
    const got = new LexError("text-a");
    expect(got.message).isNotEqualTo("", "the error has a message");
  });
}

describe("LexError", () => {
  it("matches no other error class of its module", () => {
    const got = new LexError("text-a");
    const others = [new ParseError("text-a", 1.5)];
    for (const other of others) {
      expect.soft(got).errorIsNot(other, "the error does not match another error class of the module");
    }
  });

  describe("message", lexErrorMessageCases);
});

function parseErrorMessageCases(): void {
  it("is not empty", () => {
    const got = new ParseError("text-a", 1.5);
    expect(got.message).isNotEqualTo("", "the error has a message");
  });

  it("contains each text member", () => {
    const got = new ParseError("text-a", 1.5);
    expect(got.message).doesContain(got.text, "the message contains text");
  });
}

describe("ParseError", () => {
  it("matches no other error class of its module", () => {
    const got = new ParseError("text-a", 1.5);
    const others = [new LexError("text-a")];
    for (const other of others) {
      expect.soft(got).errorIsNot(other, "the error does not match another error class of the module");
    }
  });

  describe("message", parseErrorMessageCases);
});
