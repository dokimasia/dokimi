import { check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { describe } from 'vitest';
import { LabelError } from './label';

describe("LabelError contract", () => {
  test("reports a message", ({ seat }) => {
    const got = new LabelError('cause-a', 'op-a');
    check.notEqual(seat, got.message, "", "an error must report a message");
  });

  test("puts every member it carries into its message", ({ seat }) => {
    const got = new LabelError('cause-a', 'op-a');
    check.isTrue(seat, got.message.includes(got.op), "op must reach the message");
  });
});
