import { check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { describe } from 'vitest';
import { LexError, ParseError } from './parse';

describe("LexError contract", () => {
  test("reports a message", ({ seat }) => {
    const got = new LexError('text-a');
    check.notEqual(seat, got.message, "", "an error must report a message");
  });

  test("does not match the other error classes of its module", ({ seat }) => {
    const got = new LexError('text-a');
    const others = [new ParseError('text-a', 1.5)];
    for (const other of others) {
      check.errorIsNot(seat, got, other, "an error class must not match another declared here");
    }
  });
});

describe("ParseError contract", () => {
  test("reports a message", ({ seat }) => {
    const got = new ParseError('text-a', 1.5);
    check.notEqual(seat, got.message, "", "an error must report a message");
  });

  test("does not match the other error classes of its module", ({ seat }) => {
    const got = new ParseError('text-a', 1.5);
    const others = [new LexError('text-a')];
    for (const other of others) {
      check.errorIsNot(seat, got, other, "an error class must not match another declared here");
    }
  });

  test("puts every member it carries into its message", ({ seat }) => {
    const got = new ParseError('text-a', 1.5);
    check.isTrue(seat, got.message.includes(got.text), "text must reach the message");
  });
});
