import { byIdentity, check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { describe } from 'vitest';
import type { Frozen } from './immutable';
import { FrozenBuilder } from './immutable.builder';

describe("FrozenBuilder setters", () => {
  test("withTags returns its builder", ({ seat }) => {
    const builder = FrozenBuilder.create();
    check.equal(
      seat,
      builder.withTags(builder.build().tags),
      builder,
      "the setter returns its builder for the next call",
      byIdentity(),
    );
  });
});

describe("FrozenBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: Frozen[] = [];
    const record = (value: Frozen): void => {
      seen.push(value);
    };
    FrozenBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = FrozenBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, FrozenBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});
