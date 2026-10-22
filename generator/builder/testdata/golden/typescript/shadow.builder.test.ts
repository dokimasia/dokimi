import { byIdentity } from "@dokimi/assert";
import { describe, expect, it } from "vitest";
import { Set, Shadowed } from "./shadow";
import { ShadowedBuilder } from "./shadow.builder";

function shadowedBuilderFromCases(): void {
  it("keeps each member of its value", () => {
    const seed = new Shadowed();
    seed.tags = Object.assign(new Set<string>(), { items: ["items-a"] });
    const got = ShadowedBuilder.from(seed).build();
    expect(got.tags).isEqualTo(seed.tags, "from keeps tags");
  });

  it("keeps the value of its argument", () => {
    const want = ShadowedBuilder.create().build();
    expect(ShadowedBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function shadowedBuilderWithTagsCases(): void {
  it("writes its argument to tags", () => {
    let want: Set<string>;
    want = Object.assign(new Set<string>(), { items: ["items-a"] });
    const got = ShadowedBuilder.create().withTags(want).build();
    expect(got.tags).isEqualTo(want, "withTags writes its argument to tags");
  });

  it("replaces the value of tags", () => {
    let first: Set<string>;
    let want: Set<string>;
    first = Object.assign(new Set<string>(), { items: ["items-a"] });
    want = Object.assign(new Set<string>(), { items: ["items-b"] });
    const got = ShadowedBuilder.create().withTags(first).withTags(want).build();
    expect(got.tags).isEqualTo(want, "the second call of withTags replaces the first");
  });
}

function shadowedBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: Shadowed[] = [];
    const record = (value: Shadowed): void => {
      seen.push(value);
    };
    ShadowedBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function shadowedBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = ShadowedBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = ShadowedBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });
}

function shadowedBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = ShadowedBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });

  it("returns the value under construction", () => {
    const seed = new Shadowed();
    seed.tags = Object.assign(new Set<string>(), { items: ["items-a"] });
    const builder = ShadowedBuilder.create().mutate((value) => {
      Object.assign(value, seed);
    });
    expect(builder.build()).isEqualTo(seed, "build returns the value that mutate sets");
  });
}

describe("ShadowedBuilder", () => {
  describe("from", shadowedBuilderFromCases);
  describe("withTags", shadowedBuilderWithTagsCases);
  describe("mutate", shadowedBuilderMutateCases);
  describe("clone", shadowedBuilderCloneCases);
  describe("build", shadowedBuilderBuildCases);
});
