import { check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { describe } from 'vitest';
import { ChainError, KeptError, LostError } from './unconstructed';

describe("ChainError contract", () => {
  test("reports a message", ({ seat }) => {
    const got = new ChainError('cause-a');
    check.notEqual(seat, got.message, "", "an error must report a message");
  });

  test("does not match the other error classes of its module", ({ seat }) => {
    const got = new ChainError('cause-a');
    const others = [new KeptError(), new LostError('message-a')];
    for (const other of others) {
      check.errorIsNot(seat, got, other, "an error class must not match another declared here");
    }
  });

  test("exposes the cause it was given", ({ seat }) => {
    const cause = new Error("cause");
    const got = new ChainError(cause);
    check.equal(seat, got.cause, cause, "an error naming a cause must expose it");
  });
});

describe("KeptError contract", () => {
  test("reports a message", ({ seat }) => {
    const got = new KeptError();
    check.notEqual(seat, got.message, "", "an error must report a message");
  });

  test("does not match the other error classes of its module", ({ seat }) => {
    const got = new KeptError();
    const others = [new ChainError('cause-a'), new LostError('message-a')];
    for (const other of others) {
      check.errorIsNot(seat, got, other, "an error class must not match another declared here");
    }
  });
});

describe("LostError contract", () => {
  test("reports a message", ({ seat }) => {
    const got = new LostError('message-a');
    check.notEqual(seat, got.message, "", "an error must report a message");
  });

  test("does not match the other error classes of its module", ({ seat }) => {
    const got = new LostError('message-a');
    const others = [new ChainError('cause-a'), new KeptError()];
    for (const other of others) {
      check.errorIsNot(seat, got, other, "an error class must not match another declared here");
    }
  });
});
