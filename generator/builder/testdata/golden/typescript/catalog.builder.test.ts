import { byIdentity } from "@dokimi/assert";
import { describe, expect, it } from "vitest";
import { Account, Address, type Bin, type Counter, Derived, HTTPClient, Header, Heir, type Lines, type Loose, MD5Hash, type Page, Tier, UUID, UserUrl } from "./catalog";
import { AccountBuilder, AddressBuilder, BinBuilder, CounterBuilder, DerivedBuilder, HTTPClientBuilder, HeaderBuilder, HeirBuilder, LooseBuilder, Md5HashBuilder, PageBuilder, TierBuilder, UUIDBuilder, UserUrlBuilder } from "./catalog.builder";

function addressBuilderCreateCases(): void {
  it("applies the declared default of street", () => {
    let want: string;
    want = "Dam 1";
    const got = AddressBuilder.create().build();
    expect(got.street).isEqualTo(want, "create applies the declared default of street");
  });
}

function addressBuilderFromCases(): void {
  it("keeps each member of its value", () => {
    const seed = new Address();
    seed.street = "street-a";
    seed.city = "city-a";
    const got = AddressBuilder.from(seed).build();
    expect.soft(got.street).isEqualTo(seed.street, "from keeps street");
    expect.soft(got.city).isEqualTo(seed.city, "from keeps city");
  });

  it("keeps the value of its argument", () => {
    const want = AddressBuilder.create().build();
    expect(AddressBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function addressBuilderWithStreetCases(): void {
  it("writes its argument to street", () => {
    let want: string;
    want = "street-a";
    const got = AddressBuilder.create().withStreet(want).build();
    expect(got.street).isEqualTo(want, "withStreet writes its argument to street");
  });

  it("replaces the value of street", () => {
    let first: string;
    let want: string;
    first = "street-a";
    want = "street-b";
    const got = AddressBuilder.create().withStreet(first).withStreet(want).build();
    expect(got.street).isEqualTo(want, "the second call of withStreet replaces the first");
  });
}

function addressBuilderWithCityCases(): void {
  it("writes its argument to city", () => {
    let want: string;
    want = "city-a";
    const got = AddressBuilder.create().withCity(want).build();
    expect(got.city).isEqualTo(want, "withCity writes its argument to city");
  });

  it("replaces the value of city", () => {
    let first: string;
    let want: string;
    first = "city-a";
    want = "city-b";
    const got = AddressBuilder.create().withCity(first).withCity(want).build();
    expect(got.city).isEqualTo(want, "the second call of withCity replaces the first");
  });
}

function addressBuilderWithFloorCases(): void {
  it("writes its argument to floor", () => {
    let want: number;
    want = 1.5;
    const got = AddressBuilder.create().withFloor(want).build();
    expect(got.floor).isEqualTo(want, "withFloor writes its argument to floor");
  });

  it("replaces the value of floor", () => {
    let first: number;
    let want: number;
    first = 1.5;
    want = 2.5;
    const got = AddressBuilder.create().withFloor(first).withFloor(want).build();
    expect(got.floor).isEqualTo(want, "the second call of withFloor replaces the first");
  });
}

function addressBuilderWithLinesCases(): void {
  it("writes its argument to lines", () => {
    let want: string[];
    want = ["lines-a"];
    const got = AddressBuilder.create().withLines(...want).build();
    expect(got.lines).isEqualTo(want, "withLines writes its argument to lines");
  });

  it("replaces the value of lines", () => {
    let first: string[];
    let want: string[];
    first = ["lines-a"];
    want = ["lines-b"];
    const got = AddressBuilder.create().withLines(...first).withLines(...want).build();
    expect(got.lines).isEqualTo(want, "the second call of withLines replaces the first");
  });
}

function addressBuilderWithLabelsCases(): void {
  it("writes its argument to labels", () => {
    let want: Map<string, string>;
    want = new Map<string, string>([["labels-a", "labels-a"]]);
    const got = AddressBuilder.create().withLabels(want).build();
    expect(got.labels).isEqualTo(want, "withLabels writes its argument to labels");
  });

  it("replaces the value of labels", () => {
    let first: Map<string, string>;
    let want: Map<string, string>;
    first = new Map<string, string>([["labels-a", "labels-a"]]);
    want = new Map<string, string>([["labels-b", "labels-a"]]);
    const got = AddressBuilder.create().withLabels(first).withLabels(want).build();
    expect(got.labels).isEqualTo(want, "the second call of withLabels replaces the first");
  });
}

function addressBuilderWithTagsCases(): void {
  it("writes its argument to tags", () => {
    let want: Set<string>;
    want = new Set<string>(["tags-a"]);
    const got = AddressBuilder.create().withTags(want).build();
    expect(got.tags).isEqualTo(want, "withTags writes its argument to tags");
  });

  it("replaces the value of tags", () => {
    let first: Set<string>;
    let want: Set<string>;
    first = new Set<string>(["tags-a"]);
    want = new Set<string>(["tags-b"]);
    const got = AddressBuilder.create().withTags(first).withTags(want).build();
    expect(got.tags).isEqualTo(want, "the second call of withTags replaces the first");
  });
}

function addressBuilderAppendLinesCases(): void {
  it("appends to the elements of lines", () => {
    let first: string[];
    let want: string[];
    first = ["lines-a"];
    want = ["lines-b"];
    const got = AddressBuilder.create().withLines(...first).appendLines(...want).build();
    expect(got.lines).isEqualTo([...first, ...want], "appendLines appends its elements after those of withLines");
  });
}

function addressBuilderWithLabelsEntryCases(): void {
  it("adds one key to labels", () => {
    let want: Map<string, string>;
    want = new Map<string, string>([["labels-a", "labels-a"]]);
    const builder = AddressBuilder.create();
    for (const [key, item] of want) {
      expect(builder.withLabelsEntry(key, item).build().labels.get(key)).isEqualTo(item, "withLabelsEntry sets the key");
    }
  });
}

function addressBuilderWithLabelsEntriesCases(): void {
  it("adds each key to labels", () => {
    let want: Map<string, string>;
    want = new Map<string, string>([["labels-a", "labels-a"]]);
    const got = AddressBuilder.create().withLabelsEntries(want).build();
    for (const [key, item] of want) {
      expect(got.labels.get(key)).isEqualTo(item, "withLabelsEntries copies each entry");
    }
  });
}

function addressBuilderWithTagsEntryCases(): void {
  it("adds one key to tags", () => {
    let want: Set<string>;
    want = new Set<string>(["tags-a"]);
    const builder = AddressBuilder.create();
    for (const key of want) {
      expect(builder.withTagsEntry(key).build().tags.has(key)).isTrue("withTagsEntry adds the key");
    }
  });
}

function addressBuilderWithTagsEntriesCases(): void {
  it("adds each key to tags", () => {
    let want: Set<string>;
    want = new Set<string>(["tags-a"]);
    const got = AddressBuilder.create().withTagsEntries(...want).build();
    for (const key of want) {
      expect(got.tags.has(key)).isTrue("withTagsEntries adds each key");
    }
  });
}

function addressBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: Address[] = [];
    const record = (value: Address): void => {
      seen.push(value);
    };
    AddressBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function addressBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = AddressBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = AddressBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });

  it("copies lines into the clone", () => {
    let want: string[];
    let other: string[];
    want = ["lines-a"];
    other = ["lines-b"];
    const original = AddressBuilder.create().withLines(...want);
    const clone = original.clone();
    const message = "a change of lines in the clone leaves the original unchanged";
    clone.mutate((value) => {
      value.lines.splice(0, other.length, ...other);
    });
    expect(original.build().lines).isNotEqualTo(other, message);
  });

  it("copies labels into the clone", () => {
    let want: Map<string, string>;
    let other: Map<string, string>;
    want = new Map<string, string>([["labels-a", "labels-a"]]);
    other = new Map<string, string>([["labels-b", "labels-a"]]);
    const original = AddressBuilder.create().withLabels(want);
    const clone = original.clone();
    const message = "a change of labels in the clone leaves the original unchanged";
    clone.mutate((value) => {
      for (const [key, item] of other) {
        value.labels.set(key, item);
      }
    });
    expect(original.build().labels).hasLength(1, message);
  });

  it("copies tags into the clone", () => {
    let want: Set<string>;
    let other: Set<string>;
    want = new Set<string>(["tags-a"]);
    other = new Set<string>(["tags-b"]);
    const original = AddressBuilder.create().withTags(want);
    const clone = original.clone();
    const message = "a change of tags in the clone leaves the original unchanged";
    clone.mutate((value) => {
      for (const key of other) {
        value.tags.add(key);
      }
    });
    expect(original.build().tags).hasLength(1, message);
  });

  it("builds the lines of its original", () => {
    let want: string[];
    want = ["lines-a"];
    const original = AddressBuilder.create().withLines(...want);
    expect(original.clone().build().lines).isEqualTo(want, "the clone builds the lines of its original");
  });

  it("builds the labels of its original", () => {
    let want: Map<string, string>;
    want = new Map<string, string>([["labels-a", "labels-a"]]);
    const original = AddressBuilder.create().withLabels(want);
    expect(original.clone().build().labels).isEqualTo(want, "the clone builds the labels of its original");
  });

  it("builds the tags of its original", () => {
    let want: Set<string>;
    want = new Set<string>(["tags-a"]);
    const original = AddressBuilder.create().withTags(want);
    expect(original.clone().build().tags).isEqualTo(want, "the clone builds the tags of its original");
  });
}

function addressBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = AddressBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });

  it("returns the value under construction", () => {
    const seed = new Address();
    seed.street = "street-a";
    seed.city = "city-a";
    const builder = AddressBuilder.create().mutate((value) => {
      Object.assign(value, seed);
    });
    expect(builder.build()).isEqualTo(seed, "build returns the value that mutate sets");
  });
}

describe("AddressBuilder", () => {
  describe("create", addressBuilderCreateCases);
  describe("from", addressBuilderFromCases);
  describe("withStreet", addressBuilderWithStreetCases);
  describe("withCity", addressBuilderWithCityCases);
  describe("withFloor", addressBuilderWithFloorCases);
  describe("withLines", addressBuilderWithLinesCases);
  describe("withLabels", addressBuilderWithLabelsCases);
  describe("withTags", addressBuilderWithTagsCases);
  describe("appendLines", addressBuilderAppendLinesCases);
  describe("withLabelsEntry", addressBuilderWithLabelsEntryCases);
  describe("withLabelsEntries", addressBuilderWithLabelsEntriesCases);
  describe("withTagsEntry", addressBuilderWithTagsEntryCases);
  describe("withTagsEntries", addressBuilderWithTagsEntriesCases);
  describe("mutate", addressBuilderMutateCases);
  describe("clone", addressBuilderCloneCases);
  describe("build", addressBuilderBuildCases);
});

function derivedBuilderCreateCases(): void {
  it("applies the declared default of createdBy", () => {
    let want: string;
    want = "system";
    const got = DerivedBuilder.create().build();
    expect(got.createdBy).isEqualTo(want, "create applies the declared default of createdBy");
  });
}

function derivedBuilderFromCases(): void {
  it("keeps each member of its value", () => {
    const seed = new Derived();
    seed.active = true;
    seed.createdBy = "createdBy-a";
    const got = DerivedBuilder.from(seed).build();
    expect.soft(got.active).isEqualTo(seed.active, "from keeps active");
    expect.soft(got.createdBy).isEqualTo(seed.createdBy, "from keeps createdBy");
  });

  it("keeps the value of its argument", () => {
    const want = DerivedBuilder.create().build();
    expect(DerivedBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function derivedBuilderWithActiveCases(): void {
  it("writes its argument to active", () => {
    let want: boolean;
    want = true;
    const got = DerivedBuilder.create().withActive(want).build();
    expect(got.active).isEqualTo(want, "withActive writes its argument to active");
  });

  it("replaces the value of active", () => {
    let first: boolean;
    let want: boolean;
    first = true;
    want = false;
    const got = DerivedBuilder.create().withActive(first).withActive(want).build();
    expect(got.active).isEqualTo(want, "the second call of withActive replaces the first");
  });
}

function derivedBuilderWithCreatedByCases(): void {
  it("writes its argument to createdBy", () => {
    let want: string;
    want = "createdBy-a";
    const got = DerivedBuilder.create().withCreatedBy(want).build();
    expect(got.createdBy).isEqualTo(want, "withCreatedBy writes its argument to createdBy");
  });

  it("replaces the value of createdBy", () => {
    let first: string;
    let want: string;
    first = "createdBy-a";
    want = "createdBy-b";
    const got = DerivedBuilder.create().withCreatedBy(first).withCreatedBy(want).build();
    expect(got.createdBy).isEqualTo(want, "the second call of withCreatedBy replaces the first");
  });
}

function derivedBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: Derived[] = [];
    const record = (value: Derived): void => {
      seen.push(value);
    };
    DerivedBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function derivedBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = DerivedBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = DerivedBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });
}

function derivedBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = DerivedBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });

  it("returns the value under construction", () => {
    const seed = new Derived();
    seed.active = true;
    seed.createdBy = "createdBy-a";
    const builder = DerivedBuilder.create().mutate((value) => {
      Object.assign(value, seed);
    });
    expect(builder.build()).isEqualTo(seed, "build returns the value that mutate sets");
  });
}

describe("DerivedBuilder", () => {
  describe("create", derivedBuilderCreateCases);
  describe("from", derivedBuilderFromCases);
  describe("withActive", derivedBuilderWithActiveCases);
  describe("withCreatedBy", derivedBuilderWithCreatedByCases);
  describe("mutate", derivedBuilderMutateCases);
  describe("clone", derivedBuilderCloneCases);
  describe("build", derivedBuilderBuildCases);
});

function heirBuilderCreateCases(): void {
  it("applies the declared default of createdBy", () => {
    let want: string;
    want = "system";
    const got = HeirBuilder.create().build();
    expect(got.createdBy).isEqualTo(want, "create applies the declared default of createdBy");
  });
}

function heirBuilderFromCases(): void {
  it("keeps each member of its value", () => {
    const seed = new Heir();
    seed.active = true;
    seed.createdBy = "createdBy-a";
    const got = HeirBuilder.from(seed).build();
    expect.soft(got.active).isEqualTo(seed.active, "from keeps active");
    expect.soft(got.createdBy).isEqualTo(seed.createdBy, "from keeps createdBy");
  });

  it("keeps the value of its argument", () => {
    const want = HeirBuilder.create().build();
    expect(HeirBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function heirBuilderWithActiveCases(): void {
  it("writes its argument to active", () => {
    let want: boolean;
    want = true;
    const got = HeirBuilder.create().withActive(want).build();
    expect(got.active).isEqualTo(want, "withActive writes its argument to active");
  });

  it("replaces the value of active", () => {
    let first: boolean;
    let want: boolean;
    first = true;
    want = false;
    const got = HeirBuilder.create().withActive(first).withActive(want).build();
    expect(got.active).isEqualTo(want, "the second call of withActive replaces the first");
  });
}

function heirBuilderWithCreatedByCases(): void {
  it("writes its argument to createdBy", () => {
    let want: string;
    want = "createdBy-a";
    const got = HeirBuilder.create().withCreatedBy(want).build();
    expect(got.createdBy).isEqualTo(want, "withCreatedBy writes its argument to createdBy");
  });

  it("replaces the value of createdBy", () => {
    let first: string;
    let want: string;
    first = "createdBy-a";
    want = "createdBy-b";
    const got = HeirBuilder.create().withCreatedBy(first).withCreatedBy(want).build();
    expect(got.createdBy).isEqualTo(want, "the second call of withCreatedBy replaces the first");
  });
}

function heirBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: Heir[] = [];
    const record = (value: Heir): void => {
      seen.push(value);
    };
    HeirBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function heirBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = HeirBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = HeirBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });
}

function heirBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = HeirBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });

  it("returns the value under construction", () => {
    const seed = new Heir();
    seed.active = true;
    seed.createdBy = "createdBy-a";
    const builder = HeirBuilder.create().mutate((value) => {
      Object.assign(value, seed);
    });
    expect(builder.build()).isEqualTo(seed, "build returns the value that mutate sets");
  });
}

describe("HeirBuilder", () => {
  describe("create", heirBuilderCreateCases);
  describe("from", heirBuilderFromCases);
  describe("withActive", heirBuilderWithActiveCases);
  describe("withCreatedBy", heirBuilderWithCreatedByCases);
  describe("mutate", heirBuilderMutateCases);
  describe("clone", heirBuilderCloneCases);
  describe("build", heirBuilderBuildCases);
});

function accountBuilderCreateCases(): void {
  it("applies the declared default of owner", () => {
    let want: string;
    want = "owner";
    const got = AccountBuilder.create().build();
    expect(got.owner).isEqualTo(want, "create applies the declared default of owner");
  });
}

function accountBuilderFromCases(): void {
  it("keeps each member of its value", () => {
    const seed = new Account();
    seed.owner = "owner-a";
    seed.plan = "plan-a";
    const got = AccountBuilder.from(seed).build();
    expect.soft(got.owner).isEqualTo(seed.owner, "from keeps owner");
    expect.soft(got.plan).isEqualTo(seed.plan, "from keeps plan");
  });

  it("keeps the value of its argument", () => {
    const want = AccountBuilder.create().build();
    expect(AccountBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function accountBuilderWithOwnerCases(): void {
  it("writes its argument to owner", () => {
    let want: string;
    want = "owner-a";
    const got = AccountBuilder.create().withOwner(want).build();
    expect(got.owner).isEqualTo(want, "withOwner writes its argument to owner");
  });

  it("replaces the value of owner", () => {
    let first: string;
    let want: string;
    first = "owner-a";
    want = "owner-b";
    const got = AccountBuilder.create().withOwner(first).withOwner(want).build();
    expect(got.owner).isEqualTo(want, "the second call of withOwner replaces the first");
  });
}

function accountBuilderWithPlanCases(): void {
  it("writes its argument to plan", () => {
    let want: string;
    want = "plan-a";
    const got = AccountBuilder.create().withPlan(want).build();
    expect(got.plan).isEqualTo(want, "withPlan writes its argument to plan");
  });

  it("replaces the value of plan", () => {
    let first: string;
    let want: string;
    first = "plan-a";
    want = "plan-b";
    const got = AccountBuilder.create().withPlan(first).withPlan(want).build();
    expect(got.plan).isEqualTo(want, "the second call of withPlan replaces the first");
  });
}

function accountBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: Account[] = [];
    const record = (value: Account): void => {
      seen.push(value);
    };
    AccountBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function accountBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = AccountBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = AccountBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });
}

function accountBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = AccountBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });

  it("returns the value under construction", () => {
    const seed = new Account();
    seed.owner = "owner-a";
    seed.plan = "plan-a";
    const builder = AccountBuilder.create().mutate((value) => {
      Object.assign(value, seed);
    });
    expect(builder.build()).isEqualTo(seed, "build returns the value that mutate sets");
  });
}

describe("AccountBuilder", () => {
  describe("create", accountBuilderCreateCases);
  describe("from", accountBuilderFromCases);
  describe("withOwner", accountBuilderWithOwnerCases);
  describe("withPlan", accountBuilderWithPlanCases);
  describe("mutate", accountBuilderMutateCases);
  describe("clone", accountBuilderCloneCases);
  describe("build", accountBuilderBuildCases);
});

function httpClientBuilderFromCases(): void {
  it("keeps each member of its value", () => {
    const seed = new HTTPClient();
    seed.base = "base-a";
    const got = HTTPClientBuilder.from(seed).build();
    expect(got.base).isEqualTo(seed.base, "from keeps base");
  });

  it("keeps the value of its argument", () => {
    const want = HTTPClientBuilder.create().build();
    expect(HTTPClientBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function httpClientBuilderWithBaseCases(): void {
  it("writes its argument to base", () => {
    let want: string;
    want = "base-a";
    const got = HTTPClientBuilder.create().withBase(want).build();
    expect(got.base).isEqualTo(want, "withBase writes its argument to base");
  });

  it("replaces the value of base", () => {
    let first: string;
    let want: string;
    first = "base-a";
    want = "base-b";
    const got = HTTPClientBuilder.create().withBase(first).withBase(want).build();
    expect(got.base).isEqualTo(want, "the second call of withBase replaces the first");
  });
}

function httpClientBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: HTTPClient[] = [];
    const record = (value: HTTPClient): void => {
      seen.push(value);
    };
    HTTPClientBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function httpClientBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = HTTPClientBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = HTTPClientBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });
}

function httpClientBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = HTTPClientBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });

  it("returns the value under construction", () => {
    const seed = new HTTPClient();
    seed.base = "base-a";
    const builder = HTTPClientBuilder.create().mutate((value) => {
      Object.assign(value, seed);
    });
    expect(builder.build()).isEqualTo(seed, "build returns the value that mutate sets");
  });
}

describe("HTTPClientBuilder", () => {
  describe("from", httpClientBuilderFromCases);
  describe("withBase", httpClientBuilderWithBaseCases);
  describe("mutate", httpClientBuilderMutateCases);
  describe("clone", httpClientBuilderCloneCases);
  describe("build", httpClientBuilderBuildCases);
});

function md5HashBuilderFromCases(): void {
  it("keeps each member of its value", () => {
    const seed = new MD5Hash();
    seed.digest = "digest-a";
    const got = Md5HashBuilder.from(seed).build();
    expect(got.digest).isEqualTo(seed.digest, "from keeps digest");
  });

  it("keeps the value of its argument", () => {
    const want = Md5HashBuilder.create().build();
    expect(Md5HashBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function md5HashBuilderWithDigestCases(): void {
  it("writes its argument to digest", () => {
    let want: string;
    want = "digest-a";
    const got = Md5HashBuilder.create().withDigest(want).build();
    expect(got.digest).isEqualTo(want, "withDigest writes its argument to digest");
  });

  it("replaces the value of digest", () => {
    let first: string;
    let want: string;
    first = "digest-a";
    want = "digest-b";
    const got = Md5HashBuilder.create().withDigest(first).withDigest(want).build();
    expect(got.digest).isEqualTo(want, "the second call of withDigest replaces the first");
  });
}

function md5HashBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: MD5Hash[] = [];
    const record = (value: MD5Hash): void => {
      seen.push(value);
    };
    Md5HashBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function md5HashBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = Md5HashBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = Md5HashBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });
}

function md5HashBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = Md5HashBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });

  it("returns the value under construction", () => {
    const seed = new MD5Hash();
    seed.digest = "digest-a";
    const builder = Md5HashBuilder.create().mutate((value) => {
      Object.assign(value, seed);
    });
    expect(builder.build()).isEqualTo(seed, "build returns the value that mutate sets");
  });
}

describe("Md5HashBuilder", () => {
  describe("from", md5HashBuilderFromCases);
  describe("withDigest", md5HashBuilderWithDigestCases);
  describe("mutate", md5HashBuilderMutateCases);
  describe("clone", md5HashBuilderCloneCases);
  describe("build", md5HashBuilderBuildCases);
});

function uuidBuilderFromCases(): void {
  it("keeps each member of its value", () => {
    const seed = new UUID();
    seed.value = "value-a";
    const got = UUIDBuilder.from(seed).build();
    expect(got.value).isEqualTo(seed.value, "from keeps value");
  });

  it("keeps the value of its argument", () => {
    const want = UUIDBuilder.create().build();
    expect(UUIDBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function uuidBuilderWithValueCases(): void {
  it("writes its argument to value", () => {
    let want: string;
    want = "value-a";
    const got = UUIDBuilder.create().withValue(want).build();
    expect(got.value).isEqualTo(want, "withValue writes its argument to value");
  });

  it("replaces the value of value", () => {
    let first: string;
    let want: string;
    first = "value-a";
    want = "value-b";
    const got = UUIDBuilder.create().withValue(first).withValue(want).build();
    expect(got.value).isEqualTo(want, "the second call of withValue replaces the first");
  });
}

function uuidBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: UUID[] = [];
    const record = (value: UUID): void => {
      seen.push(value);
    };
    UUIDBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function uuidBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = UUIDBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = UUIDBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });
}

function uuidBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = UUIDBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });

  it("returns the value under construction", () => {
    const seed = new UUID();
    seed.value = "value-a";
    const builder = UUIDBuilder.create().mutate((value) => {
      Object.assign(value, seed);
    });
    expect(builder.build()).isEqualTo(seed, "build returns the value that mutate sets");
  });
}

describe("UUIDBuilder", () => {
  describe("from", uuidBuilderFromCases);
  describe("withValue", uuidBuilderWithValueCases);
  describe("mutate", uuidBuilderMutateCases);
  describe("clone", uuidBuilderCloneCases);
  describe("build", uuidBuilderBuildCases);
});

function userUrlBuilderFromCases(): void {
  it("keeps each member of its value", () => {
    const seed = new UserUrl();
    seed.href = "href-a";
    seed.userId = "userId-a";
    const got = UserUrlBuilder.from(seed).build();
    expect.soft(got.href).isEqualTo(seed.href, "from keeps href");
    expect.soft(got.userId).isEqualTo(seed.userId, "from keeps userId");
  });

  it("keeps the value of its argument", () => {
    const want = UserUrlBuilder.create().build();
    expect(UserUrlBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function userUrlBuilderWithHrefCases(): void {
  it("writes its argument to href", () => {
    let want: string;
    want = "href-a";
    const got = UserUrlBuilder.create().withHref(want).build();
    expect(got.href).isEqualTo(want, "withHref writes its argument to href");
  });

  it("replaces the value of href", () => {
    let first: string;
    let want: string;
    first = "href-a";
    want = "href-b";
    const got = UserUrlBuilder.create().withHref(first).withHref(want).build();
    expect(got.href).isEqualTo(want, "the second call of withHref replaces the first");
  });
}

function userUrlBuilderWithUserIdCases(): void {
  it("writes its argument to userId", () => {
    let want: string;
    want = "userId-a";
    const got = UserUrlBuilder.create().withUserId(want).build();
    expect(got.userId).isEqualTo(want, "withUserId writes its argument to userId");
  });

  it("replaces the value of userId", () => {
    let first: string;
    let want: string;
    first = "userId-a";
    want = "userId-b";
    const got = UserUrlBuilder.create().withUserId(first).withUserId(want).build();
    expect(got.userId).isEqualTo(want, "the second call of withUserId replaces the first");
  });
}

function userUrlBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: UserUrl[] = [];
    const record = (value: UserUrl): void => {
      seen.push(value);
    };
    UserUrlBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function userUrlBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = UserUrlBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = UserUrlBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });
}

function userUrlBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = UserUrlBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });

  it("returns the value under construction", () => {
    const seed = new UserUrl();
    seed.href = "href-a";
    seed.userId = "userId-a";
    const builder = UserUrlBuilder.create().mutate((value) => {
      Object.assign(value, seed);
    });
    expect(builder.build()).isEqualTo(seed, "build returns the value that mutate sets");
  });
}

describe("UserUrlBuilder", () => {
  describe("from", userUrlBuilderFromCases);
  describe("withHref", userUrlBuilderWithHrefCases);
  describe("withUserId", userUrlBuilderWithUserIdCases);
  describe("mutate", userUrlBuilderMutateCases);
  describe("clone", userUrlBuilderCloneCases);
  describe("build", userUrlBuilderBuildCases);
});

function tierBuilderFromCases(): void {
  it("keeps each member of its value", () => {
    const seed = new Tier();
    seed.level = 1.5;
    const got = TierBuilder.from(seed).build();
    expect(got.level).isEqualTo(seed.level, "from keeps level");
  });

  it("keeps the value of its argument", () => {
    const want = TierBuilder.create().build();
    expect(TierBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function tierBuilderWithLevelCases(): void {
  it("writes its argument to level", () => {
    let want: number;
    want = 1.5;
    const got = TierBuilder.create().withLevel(want).build();
    expect(got.level).isEqualTo(want, "withLevel writes its argument to level");
  });

  it("replaces the value of level", () => {
    let first: number;
    let want: number;
    first = 1.5;
    want = 2.5;
    const got = TierBuilder.create().withLevel(first).withLevel(want).build();
    expect(got.level).isEqualTo(want, "the second call of withLevel replaces the first");
  });
}

function tierBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: Tier[] = [];
    const record = (value: Tier): void => {
      seen.push(value);
    };
    TierBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function tierBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = TierBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = TierBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });
}

function tierBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = TierBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });

  it("returns the value under construction", () => {
    const seed = new Tier();
    seed.level = 1.5;
    const builder = TierBuilder.create().mutate((value) => {
      Object.assign(value, seed);
    });
    expect(builder.build()).isEqualTo(seed, "build returns the value that mutate sets");
  });
}

describe("TierBuilder", () => {
  describe("from", tierBuilderFromCases);
  describe("withLevel", tierBuilderWithLevelCases);
  describe("mutate", tierBuilderMutateCases);
  describe("clone", tierBuilderCloneCases);
  describe("build", tierBuilderBuildCases);
});

function pageBuilderFromCases(): void {
  it("keeps the value of its argument", () => {
    const want = PageBuilder.create<unknown>().build();
    expect(PageBuilder.from<unknown>(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function pageBuilderWithItemsCases(): void {
  it("writes its argument to items", () => {
    let want: unknown[];
    want = ["items-a"];
    const got = PageBuilder.create<unknown>().withItems(...want).build();
    expect(got.items).isEqualTo(want, "withItems writes its argument to items");
  });

  it("replaces the value of items", () => {
    let first: unknown[];
    let want: unknown[];
    first = ["items-a"];
    want = ["items-b"];
    const got = PageBuilder.create<unknown>().withItems(...first).withItems(...want).build();
    expect(got.items).isEqualTo(want, "the second call of withItems replaces the first");
  });
}

function pageBuilderAppendItemsCases(): void {
  it("appends to the elements of items", () => {
    let first: unknown[];
    let want: unknown[];
    first = ["items-a"];
    want = ["items-b"];
    const got = PageBuilder.create<unknown>().withItems(...first).appendItems(...want).build();
    expect(got.items).isEqualTo([...first, ...want], "appendItems appends its elements after those of withItems");
  });
}

function pageBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: Page<unknown>[] = [];
    const record = (value: Page<unknown>): void => {
      seen.push(value);
    };
    PageBuilder.create<unknown>().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function pageBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = PageBuilder.create<unknown>();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = PageBuilder.create<unknown>();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });

  it("copies items into the clone", () => {
    let want: unknown[];
    let other: unknown[];
    want = ["items-a"];
    other = ["items-b"];
    const original = PageBuilder.create<unknown>().withItems(...want);
    const clone = original.clone();
    const message = "a change of items in the clone leaves the original unchanged";
    clone.mutate((value) => {
      value.items.splice(0, other.length, ...other);
    });
    expect(original.build().items).isNotEqualTo(other, message);
  });

  it("builds the items of its original", () => {
    let want: unknown[];
    want = ["items-a"];
    const original = PageBuilder.create<unknown>().withItems(...want);
    expect(original.clone().build().items).isEqualTo(want, "the clone builds the items of its original");
  });
}

function pageBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = PageBuilder.create<unknown>();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });
}

describe("PageBuilder", () => {
  describe("from", pageBuilderFromCases);
  describe("withItems", pageBuilderWithItemsCases);
  describe("appendItems", pageBuilderAppendItemsCases);
  describe("mutate", pageBuilderMutateCases);
  describe("clone", pageBuilderCloneCases);
  describe("build", pageBuilderBuildCases);
});

function counterBuilderFromCases(): void {
  it("keeps the value of its argument", () => {
    const want = CounterBuilder.create().build();
    expect(CounterBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function counterBuilderWithCountCases(): void {
  it("returns its builder", () => {
    const builder = CounterBuilder.create();
    const got = builder.withCount(builder.build().count);
    expect(got).isEqualTo(builder, "withCount returns its builder for the next call", byIdentity());
  });
}

function counterBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: Counter[] = [];
    const record = (value: Counter): void => {
      seen.push(value);
    };
    CounterBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function counterBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = CounterBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = CounterBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });
}

function counterBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = CounterBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });
}

describe("CounterBuilder", () => {
  describe("from", counterBuilderFromCases);
  describe("withCount", counterBuilderWithCountCases);
  describe("mutate", counterBuilderMutateCases);
  describe("clone", counterBuilderCloneCases);
  describe("build", counterBuilderBuildCases);
});

function headerBuilderFromCases(): void {
  it("keeps each member of its value", () => {
    const seed = new Header();
    seed.v2 = "v2-a";
    seed.$ref = "$ref-a";
    const got = HeaderBuilder.from(seed).build();
    expect.soft(got.v2).isEqualTo(seed.v2, "from keeps v2");
    expect.soft(got.$ref).isEqualTo(seed.$ref, "from keeps $ref");
  });

  it("keeps the value of its argument", () => {
    const want = HeaderBuilder.create().build();
    expect(HeaderBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function headerBuilderWithV2Cases(): void {
  it("writes its argument to v2", () => {
    let want: string;
    want = "v2-a";
    const got = HeaderBuilder.create().withV2(want).build();
    expect(got.v2).isEqualTo(want, "withV2 writes its argument to v2");
  });

  it("replaces the value of v2", () => {
    let first: string;
    let want: string;
    first = "v2-a";
    want = "v2-b";
    const got = HeaderBuilder.create().withV2(first).withV2(want).build();
    expect(got.v2).isEqualTo(want, "the second call of withV2 replaces the first");
  });
}

function headerBuilderWith$refCases(): void {
  it("writes its argument to $ref", () => {
    let want: string;
    want = "$ref-a";
    const got = HeaderBuilder.create().with$ref(want).build();
    expect(got.$ref).isEqualTo(want, "with$ref writes its argument to $ref");
  });

  it("replaces the value of $ref", () => {
    let first: string;
    let want: string;
    first = "$ref-a";
    want = "$ref-b";
    const got = HeaderBuilder.create().with$ref(first).with$ref(want).build();
    expect(got.$ref).isEqualTo(want, "the second call of with$ref replaces the first");
  });
}

function headerBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: Header[] = [];
    const record = (value: Header): void => {
      seen.push(value);
    };
    HeaderBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function headerBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = HeaderBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = HeaderBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });
}

function headerBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = HeaderBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });

  it("returns the value under construction", () => {
    const seed = new Header();
    seed.v2 = "v2-a";
    seed.$ref = "$ref-a";
    const builder = HeaderBuilder.create().mutate((value) => {
      Object.assign(value, seed);
    });
    expect(builder.build()).isEqualTo(seed, "build returns the value that mutate sets");
  });
}

describe("HeaderBuilder", () => {
  describe("from", headerBuilderFromCases);
  describe("withV2", headerBuilderWithV2Cases);
  describe("with$ref", headerBuilderWith$refCases);
  describe("mutate", headerBuilderMutateCases);
  describe("clone", headerBuilderCloneCases);
  describe("build", headerBuilderBuildCases);
});

function binBuilderFromCases(): void {
  it("keeps the value of its argument", () => {
    const want = BinBuilder.create().build();
    expect(BinBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function binBuilderWithCodesCases(): void {
  it("writes its argument to codes", () => {
    let want: ReadonlyArray<string>;
    want = ["codes-a"];
    const got = BinBuilder.create().withCodes(...want).build();
    expect(got.codes).isEqualTo(want, "withCodes writes its argument to codes");
  });

  it("replaces the value of codes", () => {
    let first: ReadonlyArray<string>;
    let want: ReadonlyArray<string>;
    first = ["codes-a"];
    want = ["codes-b"];
    const got = BinBuilder.create().withCodes(...first).withCodes(...want).build();
    expect(got.codes).isEqualTo(want, "the second call of withCodes replaces the first");
  });
}

function binBuilderWithMarksCases(): void {
  it("writes its argument to marks", () => {
    let want: Set<string>;
    want = new Set<string>(["marks-a"]);
    const got = BinBuilder.create().withMarks(want).build();
    expect(got.marks).isEqualTo(want, "withMarks writes its argument to marks");
  });

  it("replaces the value of marks", () => {
    let first: Set<string>;
    let want: Set<string>;
    first = new Set<string>(["marks-a"]);
    want = new Set<string>(["marks-b"]);
    const got = BinBuilder.create().withMarks(first).withMarks(want).build();
    expect(got.marks).isEqualTo(want, "the second call of withMarks replaces the first");
  });
}

function binBuilderWithNotesCases(): void {
  it("writes its argument to notes", () => {
    let want: string[];
    want = ["notes-a"];
    const got = BinBuilder.create().withNotes(want).build();
    expect(got.notes).isEqualTo(want, "withNotes writes its argument to notes");
  });

  it("replaces the value of notes", () => {
    let first: string[];
    let want: string[];
    first = ["notes-a"];
    want = ["notes-b"];
    const got = BinBuilder.create().withNotes(first).withNotes(want).build();
    expect(got.notes).isEqualTo(want, "the second call of withNotes replaces the first");
  });
}

function binBuilderWithAliasesCases(): void {
  it("writes its argument to aliases", () => {
    let want: Lines;
    want = ["aliases-a"];
    const got = BinBuilder.create().withAliases(...want).build();
    expect(got.aliases).isEqualTo(want, "withAliases writes its argument to aliases");
  });

  it("replaces the value of aliases", () => {
    let first: Lines;
    let want: Lines;
    first = ["aliases-a"];
    want = ["aliases-b"];
    const got = BinBuilder.create().withAliases(...first).withAliases(...want).build();
    expect(got.aliases).isEqualTo(want, "the second call of withAliases replaces the first");
  });
}

function binBuilderWithRawCases(): void {
  it("writes its argument to raw", () => {
    let want: Uint8Array;
    want = new Uint8Array([1]);
    const got = BinBuilder.create().withRaw(want).build();
    expect(got.raw).isEqualTo(want, "withRaw writes its argument to raw");
  });

  it("replaces the value of raw", () => {
    let first: Uint8Array;
    let want: Uint8Array;
    first = new Uint8Array([1]);
    want = new Uint8Array([2]);
    const got = BinBuilder.create().withRaw(first).withRaw(want).build();
    expect(got.raw).isEqualTo(want, "the second call of withRaw replaces the first");
  });
}

function binBuilderAppendCodesCases(): void {
  it("appends to the elements of codes", () => {
    let first: ReadonlyArray<string>;
    let want: ReadonlyArray<string>;
    first = ["codes-a"];
    want = ["codes-b"];
    const got = BinBuilder.create().withCodes(...first).appendCodes(...want).build();
    expect(got.codes).isEqualTo([...first, ...want], "appendCodes appends its elements after those of withCodes");
  });
}

function binBuilderAppendAliasesCases(): void {
  it("appends to the elements of aliases", () => {
    let first: Lines;
    let want: Lines;
    first = ["aliases-a"];
    want = ["aliases-b"];
    const got = BinBuilder.create().withAliases(...first).appendAliases(...want).build();
    expect(got.aliases).isEqualTo([...first, ...want], "appendAliases appends its elements after those of withAliases");
  });
}

function binBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: Bin[] = [];
    const record = (value: Bin): void => {
      seen.push(value);
    };
    BinBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function binBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = BinBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = BinBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });

  it("copies aliases into the clone", () => {
    let want: Lines;
    let other: Lines;
    want = ["aliases-a"];
    other = ["aliases-b"];
    const original = BinBuilder.create().withAliases(...want);
    const clone = original.clone();
    const message = "a change of aliases in the clone leaves the original unchanged";
    clone.mutate((value) => {
      value.aliases.splice(0, other.length, ...other);
    });
    expect(original.build().aliases).isNotEqualTo(other, message);
  });

  it("copies raw into the clone", () => {
    let want: Uint8Array;
    let other: Uint8Array;
    want = new Uint8Array([1]);
    other = new Uint8Array([2]);
    const original = BinBuilder.create().withRaw(want);
    const clone = original.clone();
    const message = "a change of raw in the clone leaves the original unchanged";
    clone.mutate((value) => {
      value.raw.set(other);
    });
    expect(original.build().raw).isNotEqualTo(other, message);
  });

  it("builds the aliases of its original", () => {
    let want: Lines;
    want = ["aliases-a"];
    const original = BinBuilder.create().withAliases(...want);
    expect(original.clone().build().aliases).isEqualTo(want, "the clone builds the aliases of its original");
  });

  it("builds the raw of its original", () => {
    let want: Uint8Array;
    want = new Uint8Array([1]);
    const original = BinBuilder.create().withRaw(want);
    expect(original.clone().build().raw).isEqualTo(want, "the clone builds the raw of its original");
  });
}

function binBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = BinBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });
}

describe("BinBuilder", () => {
  describe("from", binBuilderFromCases);
  describe("withCodes", binBuilderWithCodesCases);
  describe("withMarks", binBuilderWithMarksCases);
  describe("withNotes", binBuilderWithNotesCases);
  describe("withAliases", binBuilderWithAliasesCases);
  describe("withRaw", binBuilderWithRawCases);
  describe("appendCodes", binBuilderAppendCodesCases);
  describe("appendAliases", binBuilderAppendAliasesCases);
  describe("mutate", binBuilderMutateCases);
  describe("clone", binBuilderCloneCases);
  describe("build", binBuilderBuildCases);
});

function looseBuilderFromCases(): void {
  it("keeps the value of its argument", () => {
    const want = LooseBuilder.create().build();
    expect(LooseBuilder.from(want).build()).isEqualTo(want, "from keeps the value of its argument");
  });
}

function looseBuilderWithItemsCases(): void {
  it("returns its builder", () => {
    const builder = LooseBuilder.create();
    const got = builder.withItems(...builder.build().items);
    expect(got).isEqualTo(builder, "withItems returns its builder for the next call", byIdentity());
  });
}

function looseBuilderAppendItemsCases(): void {
  it("returns its builder", () => {
    const builder = LooseBuilder.create();
    const got = builder.appendItems();
    expect(got).isEqualTo(builder, "appendItems returns its builder for the next call", byIdentity());
  });
}

function looseBuilderWithIndexCases(): void {
  it("returns its builder", () => {
    const builder = LooseBuilder.create();
    const got = builder.withIndex(builder.build().index);
    expect(got).isEqualTo(builder, "withIndex returns its builder for the next call", byIdentity());
  });
}

function looseBuilderWithIndexEntryCases(): void {
  it("returns its builder", () => {
    const builder = LooseBuilder.create();
    const got = builder.withIndexEntry(undefined as never, undefined as never);
    expect(got).isEqualTo(builder, "withIndexEntry returns its builder for the next call", byIdentity());
  });
}

function looseBuilderWithIndexEntriesCases(): void {
  it("returns its builder", () => {
    const builder = LooseBuilder.create();
    const got = builder.withIndexEntries(builder.build().index);
    expect(got).isEqualTo(builder, "withIndexEntries returns its builder for the next call", byIdentity());
  });
}

function looseBuilderWithKeysCases(): void {
  it("returns its builder", () => {
    const builder = LooseBuilder.create();
    const got = builder.withKeys(builder.build().keys);
    expect(got).isEqualTo(builder, "withKeys returns its builder for the next call", byIdentity());
  });
}

function looseBuilderWithKeysEntryCases(): void {
  it("returns its builder", () => {
    const builder = LooseBuilder.create();
    const got = builder.withKeysEntry(undefined as never);
    expect(got).isEqualTo(builder, "withKeysEntry returns its builder for the next call", byIdentity());
  });
}

function looseBuilderWithKeysEntriesCases(): void {
  it("returns its builder", () => {
    const builder = LooseBuilder.create();
    const got = builder.withKeysEntries();
    expect(got).isEqualTo(builder, "withKeysEntries returns its builder for the next call", byIdentity());
  });
}

function looseBuilderWithExtraCases(): void {
  it("returns its builder", () => {
    const builder = LooseBuilder.create();
    const got = builder.withExtra(undefined as never);
    expect(got).isEqualTo(builder, "withExtra returns its builder for the next call", byIdentity());
  });
}

function looseBuilderWithValueCases(): void {
  it("returns its builder", () => {
    const builder = LooseBuilder.create();
    const got = builder.withValue(builder.build().value);
    expect(got).isEqualTo(builder, "withValue returns its builder for the next call", byIdentity());
  });
}

function looseBuilderMutateCases(): void {
  it("passes the value under construction", () => {
    const seen: Loose[] = [];
    const record = (value: Loose): void => {
      seen.push(value);
    };
    LooseBuilder.create().mutate(record).mutate(record);
    expect(seen).hasLength(2, "mutate calls each function once");
    expect(seen[0]).isEqualTo(seen[1], "the two functions get the same value", byIdentity());
  });
}

function looseBuilderCloneCases(): void {
  it("returns a builder of its own", () => {
    const original = LooseBuilder.create();
    expect(original.clone()).isNotEqualTo(original, "clone returns a builder of its own", byIdentity());
  });

  it("builds the value of its original", () => {
    const original = LooseBuilder.create();
    expect(original.clone().build()).isEqualTo(original.build(), "the clone builds the value of its original");
  });
}

function looseBuilderBuildCases(): void {
  it("returns a new object for each call", () => {
    const builder = LooseBuilder.create();
    expect(builder.build()).isNotEqualTo(builder.build(), "build returns a new object", byIdentity());
  });
}

describe("LooseBuilder", () => {
  describe("from", looseBuilderFromCases);
  describe("withItems", looseBuilderWithItemsCases);
  describe("appendItems", looseBuilderAppendItemsCases);
  describe("withIndex", looseBuilderWithIndexCases);
  describe("withIndexEntry", looseBuilderWithIndexEntryCases);
  describe("withIndexEntries", looseBuilderWithIndexEntriesCases);
  describe("withKeys", looseBuilderWithKeysCases);
  describe("withKeysEntry", looseBuilderWithKeysEntryCases);
  describe("withKeysEntries", looseBuilderWithKeysEntriesCases);
  describe("withExtra", looseBuilderWithExtraCases);
  describe("withValue", looseBuilderWithValueCases);
  describe("mutate", looseBuilderMutateCases);
  describe("clone", looseBuilderCloneCases);
  describe("build", looseBuilderBuildCases);
});
