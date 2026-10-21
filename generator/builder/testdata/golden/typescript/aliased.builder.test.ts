import { byIdentity, check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { describe } from 'vitest';
import type { Mirror } from './aliased';
import { MirrorBuilder } from './aliased.builder';
import { Account } from './catalog';

describe("MirrorBuilder construction", () => {
  test("from round-trips what it was seeded with", ({ seat }) => {
    const seed = new Account();
    seed.owner = 'owner-a';
    seed.plan = 'plan-a';
    const got = MirrorBuilder.from(seed).build();
    check.equal(seat, got.owner, seed.owner, "from keeps the member it was given");
    check.equal(seat, got.plan, seed.plan, "from keeps the member it was given");
  });

  test("owner carries its declared default", ({ seat }) => {
    let want: string;
    want = 'owner';
    const got = MirrorBuilder.create().build();
    check.equal(seat, got.owner, want, "the constructor applies the declared default");
  });
});

describe("MirrorBuilder setters", () => {
  test("withOwner writes owner", ({ seat }) => {
    let want: string;
    want = 'owner-a';
    const got = MirrorBuilder.create().withOwner(want).build();
    check.equal(seat, got.owner, want, "the setter writes what it was given");
  });

  test("withOwner replaces what was already there", ({ seat }) => {
    let first: string;
    let want: string;
    first = 'owner-a';
    want = 'owner-b';
    const got = MirrorBuilder.create().withOwner(first).withOwner(want).build();
    check.equal(seat, got.owner, want, "the later call wins");
  });

  test("withPlan writes plan", ({ seat }) => {
    let want: string;
    want = 'plan-a';
    const got = MirrorBuilder.create().withPlan(want).build();
    check.equal(seat, got.plan, want, "the setter writes what it was given");
  });

  test("withPlan replaces what was already there", ({ seat }) => {
    let first: string;
    let want: string;
    first = 'plan-a';
    want = 'plan-b';
    const got = MirrorBuilder.create().withPlan(first).withPlan(want).build();
    check.equal(seat, got.plan, want, "the later call wins");
  });
});

describe("MirrorBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: Mirror[] = [];
    const record = (value: Mirror): void => {
      seen.push(value);
    };
    MirrorBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = MirrorBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, MirrorBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});
