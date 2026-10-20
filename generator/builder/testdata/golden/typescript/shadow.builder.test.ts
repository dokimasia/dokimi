import { byIdentity, check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { describe } from 'vitest';
import type { Shadowed } from './shadow';
import { ShadowedBuilder } from './shadow.builder';

describe("ShadowedBuilder setters", () => {
  test("withTags returns its builder", ({ seat }) => {
    const builder = ShadowedBuilder.create();
    check.equal(
      seat,
      builder.withTags(builder.build().tags),
      builder,
      "the setter returns its builder for the next call",
      byIdentity(),
    );
  });
});

describe("ShadowedBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: Shadowed[] = [];
    const record = (value: Shadowed): void => {
      seen.push(value);
    };
    ShadowedBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = ShadowedBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, ShadowedBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});
