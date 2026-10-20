import { check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { describe } from 'vitest';
import { ClosedError, MissingError, NotFoundError, WrappedError } from './store';

describe("ClosedError contract", () => {
  test("reports a message", ({ seat }) => {
    const got = new ClosedError();
    check.notEqual(seat, got.message, "", "an error must report a message");
  });

  test("carries the prefix of its module", ({ seat }) => {
    const got = new ClosedError();
    check.isTrue(seat, got.message.startsWith("kv: "), "an error class must carry the prefix of its module");
  });

  test("does not match the other error classes of its module", ({ seat }) => {
    const got = new ClosedError();
    const others = [new MissingError('key-a'), new NotFoundError('key-a'), new WrappedError('op-a', 'cause-a')];
    for (const other of others) {
      check.errorIsNot(seat, got, other, "an error class must not match another declared here");
    }
  });
});

describe("MissingError contract", () => {
  test("reports a message", ({ seat }) => {
    const got = new MissingError('key-a');
    check.notEqual(seat, got.message, "", "an error must report a message");
  });

  test("carries the prefix of its module", ({ seat }) => {
    const got = new MissingError('key-a');
    check.isTrue(seat, got.message.startsWith("kv: "), "an error class must carry the prefix of its module");
  });

  test("does not match the other error classes of its module", ({ seat }) => {
    const got = new MissingError('key-a');
    const others = [new ClosedError(), new NotFoundError('key-a'), new WrappedError('op-a', 'cause-a')];
    for (const other of others) {
      check.errorIsNot(seat, got, other, "an error class must not match another declared here");
    }
  });
});

describe("NotFoundError contract", () => {
  test("reports a message", ({ seat }) => {
    const got = new NotFoundError('key-a');
    check.notEqual(seat, got.message, "", "an error must report a message");
  });

  test("carries the prefix of its module", ({ seat }) => {
    const got = new NotFoundError('key-a');
    check.isTrue(seat, got.message.startsWith("kv: "), "an error class must carry the prefix of its module");
  });

  test("does not match the other error classes of its module", ({ seat }) => {
    const got = new NotFoundError('key-a');
    const others = [new ClosedError(), new MissingError('key-a'), new WrappedError('op-a', 'cause-a')];
    for (const other of others) {
      check.errorIsNot(seat, got, other, "an error class must not match another declared here");
    }
  });

  test("puts every member it carries into its message", ({ seat }) => {
    const got = new NotFoundError('key-a');
    check.isTrue(seat, got.message.includes(got.key), "key must reach the message");
  });
});

describe("WrappedError contract", () => {
  test("reports a message", ({ seat }) => {
    const got = new WrappedError('op-a', 'cause-a');
    check.notEqual(seat, got.message, "", "an error must report a message");
  });

  test("carries the prefix of its module", ({ seat }) => {
    const got = new WrappedError('op-a', 'cause-a');
    check.isTrue(seat, got.message.startsWith("kv: "), "an error class must carry the prefix of its module");
  });

  test("does not match the other error classes of its module", ({ seat }) => {
    const got = new WrappedError('op-a', 'cause-a');
    const others = [new ClosedError(), new MissingError('key-a'), new NotFoundError('key-a')];
    for (const other of others) {
      check.errorIsNot(seat, got, other, "an error class must not match another declared here");
    }
  });

  test("exposes the cause it was given", ({ seat }) => {
    const cause = new Error("cause");
    const got = new WrappedError('op-a', cause);
    check.equal(seat, got.cause, cause, "an error naming a cause must expose it");
  });

  test("puts every member it carries into its message", ({ seat }) => {
    const got = new WrappedError('op-a', 'cause-a');
    check.isTrue(seat, got.message.includes(got.op), "op must reach the message");
  });
});
