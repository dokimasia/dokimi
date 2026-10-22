import { byIdentity } from "@dokimi/assert";
import { describe, expect, it } from "vitest";
import type { Frozen } from "./immutable";
import { FrozenBuilder } from "./immutable.builder";

function frozenBuilderFromCases(): void {
  it("keeps the value of its argument", () => {
    const want = FrozenBuilder.create().build();
    expect(FrozenBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function frozenBuilderWithTagsCases(): void {
  it("returns its builder", () => {
    const builder = FrozenBuilder.create();
    const got = builder.withTags(builder.build().tags);
    expect(got).isEqualTo(builder, "withTags returns its builder for the next call", byIdentity());
  });
}

function frozenBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: Frozen[] = [];
    const record = (value: Frozen): void => {
      seen.push(value);
    };
    FrozenBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function frozenBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = FrozenBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = FrozenBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });
}

function frozenBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = FrozenBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });
}

describe("FrozenBuilder", () => {
  describe("from", frozenBuilderFromCases);
  describe("withTags", frozenBuilderWithTagsCases);
  describe("mutate", frozenBuilderMutateCases);
  describe("clone", frozenBuilderCloneCases);
  describe("build", frozenBuilderBuildCases);
});
