import { byIdentity } from "@dokimi/assert";
import { describe, expect, it } from "vitest";
import type { Mirror } from "./aliased";
import { MirrorBuilder } from "./aliased.builder";
import { Account } from "./catalog";

function mirrorBuilderCreateCases(): void {
  it("applies the declared default of owner", () => {
    let want: string;
    want = "owner";
    const got = MirrorBuilder.create().build();
    expect(got.owner).isEqualTo(want, "create applies the declared default of owner");
  });
}

function mirrorBuilderFromCases(): void {
  it("keeps each member of its value", () => {
    const seed = new Account();
    seed.owner = "owner-a";
    seed.plan = "plan-a";
    const got = MirrorBuilder.from(seed).build();
    expect.soft(got.owner).isEqualTo(seed.owner, "from keeps owner");
    expect.soft(got.plan).isEqualTo(seed.plan, "from keeps plan");
  });

  it("keeps the value of its argument", () => {
    const want = MirrorBuilder.create().build();
    expect(MirrorBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function mirrorBuilderWithOwnerCases(): void {
  it("writes its argument to owner", () => {
    let want: string;
    want = "owner-a";
    const got = MirrorBuilder.create().withOwner(want).build();
    expect(got.owner).isEqualTo(want, "withOwner writes its argument to owner");
  });

  it("replaces the value of owner", () => {
    let first: string;
    let want: string;
    first = "owner-a";
    want = "owner-b";
    const got = MirrorBuilder.create().withOwner(first).withOwner(want).build();
    expect(got.owner).isEqualTo(want, "the second call of withOwner replaces the first");
  });
}

function mirrorBuilderWithPlanCases(): void {
  it("writes its argument to plan", () => {
    let want: string;
    want = "plan-a";
    const got = MirrorBuilder.create().withPlan(want).build();
    expect(got.plan).isEqualTo(want, "withPlan writes its argument to plan");
  });

  it("replaces the value of plan", () => {
    let first: string;
    let want: string;
    first = "plan-a";
    want = "plan-b";
    const got = MirrorBuilder.create().withPlan(first).withPlan(want).build();
    expect(got.plan).isEqualTo(want, "the second call of withPlan replaces the first");
  });
}

function mirrorBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: Mirror[] = [];
    const record = (value: Mirror): void => {
      seen.push(value);
    };
    MirrorBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function mirrorBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = MirrorBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = MirrorBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });
}

function mirrorBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = MirrorBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });

  it("returns the value under construction", () => {
    const seed = new Account();
    seed.owner = "owner-a";
    seed.plan = "plan-a";
    const builder = MirrorBuilder.create().mutate((value) => {
      Object.assign(value, seed);
    });
    expect(builder.build()).isEqualTo(seed, "build returns the value that mutate sets");
  });
}

describe("MirrorBuilder", () => {
  describe("create", mirrorBuilderCreateCases);
  describe("from", mirrorBuilderFromCases);
  describe("withOwner", mirrorBuilderWithOwnerCases);
  describe("withPlan", mirrorBuilderWithPlanCases);
  describe("mutate", mirrorBuilderMutateCases);
  describe("clone", mirrorBuilderCloneCases);
  describe("build", mirrorBuilderBuildCases);
});
