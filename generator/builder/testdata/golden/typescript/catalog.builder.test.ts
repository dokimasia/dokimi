import { byIdentity, check } from '@dokimi/assert';
import { test } from '@dokimi/assert/vitest';
import { describe } from 'vitest';
import { Account, Address, type Bin, type Counter, Derived, HTTPClient, Header, Heir, type Lines, type Loose, MD5Hash, type Page, Tier, UUID, UserUrl } from './catalog';
import { AccountBuilder, AddressBuilder, BinBuilder, CounterBuilder, DerivedBuilder, HTTPClientBuilder, HeaderBuilder, HeirBuilder, LooseBuilder, Md5HashBuilder, PageBuilder, TierBuilder, UUIDBuilder, UserUrlBuilder } from './catalog.builder';

describe("AddressBuilder construction", () => {
  test("from round-trips what it was seeded with", ({ seat }) => {
    const seed = new Address();
    seed.street = 'street-a';
    seed.city = 'city-a';
    const got = AddressBuilder.from(seed).build();
    check.equal(seat, got.street, seed.street, "from keeps the member it was given");
    check.equal(seat, got.city, seed.city, "from keeps the member it was given");
  });

  test("street carries its declared default", ({ seat }) => {
    let want: string;
    want = 'Dam 1';
    const got = AddressBuilder.create().build();
    check.equal(seat, got.street, want, "the constructor applies the declared default");
  });
});

describe("AddressBuilder setters", () => {
  test("withStreet writes street", ({ seat }) => {
    let want: string;
    want = 'street-a';
    const got = AddressBuilder.create().withStreet(want).build();
    check.equal(seat, got.street, want, "the setter writes what it was given");
  });

  test("withStreet replaces what was already there", ({ seat }) => {
    let first: string;
    let want: string;
    first = 'street-a';
    want = 'street-b';
    const got = AddressBuilder.create().withStreet(first).withStreet(want).build();
    check.equal(seat, got.street, want, "the later call wins");
  });

  test("withCity writes city", ({ seat }) => {
    let want: string;
    want = 'city-a';
    const got = AddressBuilder.create().withCity(want).build();
    check.equal(seat, got.city, want, "the setter writes what it was given");
  });

  test("withCity replaces what was already there", ({ seat }) => {
    let first: string;
    let want: string;
    first = 'city-a';
    want = 'city-b';
    const got = AddressBuilder.create().withCity(first).withCity(want).build();
    check.equal(seat, got.city, want, "the later call wins");
  });

  test("withFloor writes floor", ({ seat }) => {
    let want: number;
    want = 1.5;
    const got = AddressBuilder.create().withFloor(want).build();
    check.isNotNil(seat, got.floor, "the setter writes the value it was given");
  });

  test("withLines writes lines", ({ seat }) => {
    let want: string[];
    want = ['lines-a'];
    const got = AddressBuilder.create().withLines(...want).build();
    check.length(seat, got.lines, 1, "the setter replaces the sequence");
  });

  test("withLabels writes labels", ({ seat }) => {
    let want: Map<string, string>;
    want = new Map<string, string>([['labels-a', 'labels-a']]);
    const got = AddressBuilder.create().withLabels(want).build();
    check.length(seat, got.labels, 1, "the setter replaces the collection");
  });

  test("withTags writes tags", ({ seat }) => {
    let want: Set<string>;
    want = new Set<string>(['tags-a']);
    const got = AddressBuilder.create().withTags(want).build();
    check.length(seat, got.tags, 1, "the setter replaces the collection");
  });

  test("appendLines keeps what lines holds", ({ seat }) => {
    let want: string[];
    want = ['lines-a'];
    const got = AddressBuilder.create().withLines(...want).appendLines(...want).build();
    check.length(seat, got.lines, 2, "the appender adds to the sequence");
  });

  test("withLabelsEntry adds one key to labels", ({ seat }) => {
    let want: Map<string, string>;
    want = new Map<string, string>([['labels-a', 'labels-a']]);
    const builder = AddressBuilder.create();
    for (const [key, item] of want) {
      check.equal(seat, builder.withLabelsEntry(key, item).build().labels.get(key), item, "the setter sets the key");
    }
  });

  test("withLabelsEntries adds each key to labels", ({ seat }) => {
    let want: Map<string, string>;
    want = new Map<string, string>([['labels-a', 'labels-a']]);
    const got = AddressBuilder.create().withLabelsEntries(want).build();
    for (const [key, item] of want) {
      check.equal(seat, got.labels.get(key), item, "the setter merges each entry");
    }
  });

  test("withTagsEntry adds one key to tags", ({ seat }) => {
    let want: Set<string>;
    want = new Set<string>(['tags-a']);
    const builder = AddressBuilder.create();
    for (const key of want) {
      check.isTrue(seat, builder.withTagsEntry(key).build().tags.has(key), "the setter adds the key");
    }
  });

  test("withTagsEntries adds each key to tags", ({ seat }) => {
    let want: Set<string>;
    want = new Set<string>(['tags-a']);
    const got = AddressBuilder.create().withTagsEntries(...want).build();
    for (const key of want) {
      check.isTrue(seat, got.tags.has(key), "the setter adds each key");
    }
  });
});

describe("AddressBuilder clone", () => {
  test("lines is not shared with a clone", ({ seat }) => {
    let want: string[];
    let other: string[];
    want = ['lines-a'];
    other = ['lines-b'];
    const original = AddressBuilder.create().withLines(...want);
    const clone = original.clone();
    clone.mutate((value) => {
      value.lines.splice(0, other.length, ...other);
    });
    check.notEqual(seat, original.build().lines, other, "writing through a clone must not reach the original");
  });

  test("labels is not shared with a clone", ({ seat }) => {
    let want: Map<string, string>;
    let other: Map<string, string>;
    want = new Map<string, string>([['labels-a', 'labels-a']]);
    other = new Map<string, string>([['labels-b', 'labels-a']]);
    const original = AddressBuilder.create().withLabels(want);
    const clone = original.clone();
    clone.mutate((value) => {
      for (const [key, item] of other) {
        value.labels.set(key, item);
      }
    });
    check.length(seat, original.build().labels, 1, "writing through a clone must not reach the original");
  });

  test("tags is not shared with a clone", ({ seat }) => {
    let want: Set<string>;
    let other: Set<string>;
    want = new Set<string>(['tags-a']);
    other = new Set<string>(['tags-b']);
    const original = AddressBuilder.create().withTags(want);
    const clone = original.clone();
    clone.mutate((value) => {
      for (const key of other) {
        value.tags.add(key);
      }
    });
    check.length(seat, original.build().tags, 1, "writing through a clone must not reach the original");
  });

  test("clone copies lines", ({ seat }) => {
    let want: string[];
    want = ['lines-a'];
    const original = AddressBuilder.create().withLines(...want);
    check.equal(seat, original.clone().build().lines, want, "a clone builds the same member as its original");
  });

  test("clone copies labels", ({ seat }) => {
    let want: Map<string, string>;
    want = new Map<string, string>([['labels-a', 'labels-a']]);
    const original = AddressBuilder.create().withLabels(want);
    check.equal(seat, original.clone().build().labels, want, "a clone builds the same member as its original");
  });

  test("clone copies tags", ({ seat }) => {
    let want: Set<string>;
    want = new Set<string>(['tags-a']);
    const original = AddressBuilder.create().withTags(want);
    check.equal(seat, original.clone().build().tags, want, "a clone builds the same member as its original");
  });
});

describe("AddressBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: Address[] = [];
    const record = (value: Address): void => {
      seen.push(value);
    };
    AddressBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = AddressBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, AddressBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});

describe("DerivedBuilder construction", () => {
  test("from round-trips what it was seeded with", ({ seat }) => {
    const seed = new Derived();
    seed.active = true;
    seed.createdBy = 'createdBy-a';
    const got = DerivedBuilder.from(seed).build();
    check.equal(seat, got.active, seed.active, "from keeps the member it was given");
    check.equal(seat, got.createdBy, seed.createdBy, "from keeps the member it was given");
  });

  test("createdBy carries its declared default", ({ seat }) => {
    let want: string;
    want = 'system';
    const got = DerivedBuilder.create().build();
    check.equal(seat, got.createdBy, want, "the constructor applies the declared default");
  });
});

describe("DerivedBuilder setters", () => {
  test("withActive writes active", ({ seat }) => {
    let want: boolean;
    want = true;
    const got = DerivedBuilder.create().withActive(want).build();
    check.equal(seat, got.active, want, "the setter writes what it was given");
  });

  test("withActive replaces what was already there", ({ seat }) => {
    let first: boolean;
    let want: boolean;
    first = true;
    want = false;
    const got = DerivedBuilder.create().withActive(first).withActive(want).build();
    check.equal(seat, got.active, want, "the later call wins");
  });

  test("withCreatedBy writes createdBy", ({ seat }) => {
    let want: string;
    want = 'createdBy-a';
    const got = DerivedBuilder.create().withCreatedBy(want).build();
    check.equal(seat, got.createdBy, want, "the setter writes what it was given");
  });

  test("withCreatedBy replaces what was already there", ({ seat }) => {
    let first: string;
    let want: string;
    first = 'createdBy-a';
    want = 'createdBy-b';
    const got = DerivedBuilder.create().withCreatedBy(first).withCreatedBy(want).build();
    check.equal(seat, got.createdBy, want, "the later call wins");
  });
});

describe("DerivedBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: Derived[] = [];
    const record = (value: Derived): void => {
      seen.push(value);
    };
    DerivedBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = DerivedBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, DerivedBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});

describe("HeirBuilder construction", () => {
  test("from round-trips what it was seeded with", ({ seat }) => {
    const seed = new Heir();
    seed.active = true;
    seed.createdBy = 'createdBy-a';
    const got = HeirBuilder.from(seed).build();
    check.equal(seat, got.active, seed.active, "from keeps the member it was given");
    check.equal(seat, got.createdBy, seed.createdBy, "from keeps the member it was given");
  });

  test("createdBy carries its declared default", ({ seat }) => {
    let want: string;
    want = 'system';
    const got = HeirBuilder.create().build();
    check.equal(seat, got.createdBy, want, "the constructor applies the declared default");
  });
});

describe("HeirBuilder setters", () => {
  test("withActive writes active", ({ seat }) => {
    let want: boolean;
    want = true;
    const got = HeirBuilder.create().withActive(want).build();
    check.equal(seat, got.active, want, "the setter writes what it was given");
  });

  test("withActive replaces what was already there", ({ seat }) => {
    let first: boolean;
    let want: boolean;
    first = true;
    want = false;
    const got = HeirBuilder.create().withActive(first).withActive(want).build();
    check.equal(seat, got.active, want, "the later call wins");
  });

  test("withCreatedBy writes createdBy", ({ seat }) => {
    let want: string;
    want = 'createdBy-a';
    const got = HeirBuilder.create().withCreatedBy(want).build();
    check.equal(seat, got.createdBy, want, "the setter writes what it was given");
  });

  test("withCreatedBy replaces what was already there", ({ seat }) => {
    let first: string;
    let want: string;
    first = 'createdBy-a';
    want = 'createdBy-b';
    const got = HeirBuilder.create().withCreatedBy(first).withCreatedBy(want).build();
    check.equal(seat, got.createdBy, want, "the later call wins");
  });
});

describe("HeirBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: Heir[] = [];
    const record = (value: Heir): void => {
      seen.push(value);
    };
    HeirBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = HeirBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, HeirBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});

describe("AccountBuilder construction", () => {
  test("from round-trips what it was seeded with", ({ seat }) => {
    const seed = new Account();
    seed.owner = 'owner-a';
    seed.plan = 'plan-a';
    const got = AccountBuilder.from(seed).build();
    check.equal(seat, got.owner, seed.owner, "from keeps the member it was given");
    check.equal(seat, got.plan, seed.plan, "from keeps the member it was given");
  });

  test("owner carries its declared default", ({ seat }) => {
    let want: string;
    want = 'owner';
    const got = AccountBuilder.create().build();
    check.equal(seat, got.owner, want, "the constructor applies the declared default");
  });
});

describe("AccountBuilder setters", () => {
  test("withOwner writes owner", ({ seat }) => {
    let want: string;
    want = 'owner-a';
    const got = AccountBuilder.create().withOwner(want).build();
    check.equal(seat, got.owner, want, "the setter writes what it was given");
  });

  test("withOwner replaces what was already there", ({ seat }) => {
    let first: string;
    let want: string;
    first = 'owner-a';
    want = 'owner-b';
    const got = AccountBuilder.create().withOwner(first).withOwner(want).build();
    check.equal(seat, got.owner, want, "the later call wins");
  });

  test("withPlan writes plan", ({ seat }) => {
    let want: string;
    want = 'plan-a';
    const got = AccountBuilder.create().withPlan(want).build();
    check.equal(seat, got.plan, want, "the setter writes what it was given");
  });

  test("withPlan replaces what was already there", ({ seat }) => {
    let first: string;
    let want: string;
    first = 'plan-a';
    want = 'plan-b';
    const got = AccountBuilder.create().withPlan(first).withPlan(want).build();
    check.equal(seat, got.plan, want, "the later call wins");
  });
});

describe("AccountBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: Account[] = [];
    const record = (value: Account): void => {
      seen.push(value);
    };
    AccountBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = AccountBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, AccountBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});

describe("HTTPClientBuilder construction", () => {
  test("from round-trips what it was seeded with", ({ seat }) => {
    const seed = new HTTPClient();
    seed.base = 'base-a';
    const got = HTTPClientBuilder.from(seed).build();
    check.equal(seat, got.base, seed.base, "from keeps the member it was given");
  });
});

describe("HTTPClientBuilder setters", () => {
  test("withBase writes base", ({ seat }) => {
    let want: string;
    want = 'base-a';
    const got = HTTPClientBuilder.create().withBase(want).build();
    check.equal(seat, got.base, want, "the setter writes what it was given");
  });

  test("withBase replaces what was already there", ({ seat }) => {
    let first: string;
    let want: string;
    first = 'base-a';
    want = 'base-b';
    const got = HTTPClientBuilder.create().withBase(first).withBase(want).build();
    check.equal(seat, got.base, want, "the later call wins");
  });
});

describe("HTTPClientBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: HTTPClient[] = [];
    const record = (value: HTTPClient): void => {
      seen.push(value);
    };
    HTTPClientBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = HTTPClientBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, HTTPClientBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});

describe("Md5HashBuilder construction", () => {
  test("from round-trips what it was seeded with", ({ seat }) => {
    const seed = new MD5Hash();
    seed.digest = 'digest-a';
    const got = Md5HashBuilder.from(seed).build();
    check.equal(seat, got.digest, seed.digest, "from keeps the member it was given");
  });
});

describe("Md5HashBuilder setters", () => {
  test("withDigest writes digest", ({ seat }) => {
    let want: string;
    want = 'digest-a';
    const got = Md5HashBuilder.create().withDigest(want).build();
    check.equal(seat, got.digest, want, "the setter writes what it was given");
  });

  test("withDigest replaces what was already there", ({ seat }) => {
    let first: string;
    let want: string;
    first = 'digest-a';
    want = 'digest-b';
    const got = Md5HashBuilder.create().withDigest(first).withDigest(want).build();
    check.equal(seat, got.digest, want, "the later call wins");
  });
});

describe("Md5HashBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: MD5Hash[] = [];
    const record = (value: MD5Hash): void => {
      seen.push(value);
    };
    Md5HashBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = Md5HashBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, Md5HashBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});

describe("UUIDBuilder construction", () => {
  test("from round-trips what it was seeded with", ({ seat }) => {
    const seed = new UUID();
    seed.value = 'value-a';
    const got = UUIDBuilder.from(seed).build();
    check.equal(seat, got.value, seed.value, "from keeps the member it was given");
  });
});

describe("UUIDBuilder setters", () => {
  test("withValue writes value", ({ seat }) => {
    let want: string;
    want = 'value-a';
    const got = UUIDBuilder.create().withValue(want).build();
    check.equal(seat, got.value, want, "the setter writes what it was given");
  });

  test("withValue replaces what was already there", ({ seat }) => {
    let first: string;
    let want: string;
    first = 'value-a';
    want = 'value-b';
    const got = UUIDBuilder.create().withValue(first).withValue(want).build();
    check.equal(seat, got.value, want, "the later call wins");
  });
});

describe("UUIDBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: UUID[] = [];
    const record = (value: UUID): void => {
      seen.push(value);
    };
    UUIDBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = UUIDBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, UUIDBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});

describe("UserUrlBuilder construction", () => {
  test("from round-trips what it was seeded with", ({ seat }) => {
    const seed = new UserUrl();
    seed.href = 'href-a';
    seed.userId = 'userId-a';
    const got = UserUrlBuilder.from(seed).build();
    check.equal(seat, got.href, seed.href, "from keeps the member it was given");
    check.equal(seat, got.userId, seed.userId, "from keeps the member it was given");
  });
});

describe("UserUrlBuilder setters", () => {
  test("withHref writes href", ({ seat }) => {
    let want: string;
    want = 'href-a';
    const got = UserUrlBuilder.create().withHref(want).build();
    check.equal(seat, got.href, want, "the setter writes what it was given");
  });

  test("withHref replaces what was already there", ({ seat }) => {
    let first: string;
    let want: string;
    first = 'href-a';
    want = 'href-b';
    const got = UserUrlBuilder.create().withHref(first).withHref(want).build();
    check.equal(seat, got.href, want, "the later call wins");
  });

  test("withUserId writes userId", ({ seat }) => {
    let want: string;
    want = 'userId-a';
    const got = UserUrlBuilder.create().withUserId(want).build();
    check.equal(seat, got.userId, want, "the setter writes what it was given");
  });

  test("withUserId replaces what was already there", ({ seat }) => {
    let first: string;
    let want: string;
    first = 'userId-a';
    want = 'userId-b';
    const got = UserUrlBuilder.create().withUserId(first).withUserId(want).build();
    check.equal(seat, got.userId, want, "the later call wins");
  });
});

describe("UserUrlBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: UserUrl[] = [];
    const record = (value: UserUrl): void => {
      seen.push(value);
    };
    UserUrlBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = UserUrlBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, UserUrlBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});

describe("TierBuilder construction", () => {
  test("from round-trips what it was seeded with", ({ seat }) => {
    const seed = new Tier();
    seed.level = 1.5;
    const got = TierBuilder.from(seed).build();
    check.equal(seat, got.level, seed.level, "from keeps the member it was given");
  });
});

describe("TierBuilder setters", () => {
  test("withLevel writes level", ({ seat }) => {
    let want: number;
    want = 1.5;
    const got = TierBuilder.create().withLevel(want).build();
    check.equal(seat, got.level, want, "the setter writes what it was given");
  });

  test("withLevel replaces what was already there", ({ seat }) => {
    let first: number;
    let want: number;
    first = 1.5;
    want = 2.5;
    const got = TierBuilder.create().withLevel(first).withLevel(want).build();
    check.equal(seat, got.level, want, "the later call wins");
  });
});

describe("TierBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: Tier[] = [];
    const record = (value: Tier): void => {
      seen.push(value);
    };
    TierBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = TierBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, TierBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});

describe("PageBuilder setters", () => {
  test("withItems writes items", ({ seat }) => {
    let want: unknown[];
    want = ['items-a'];
    const got = PageBuilder.create<unknown>().withItems(...want).build();
    check.length(seat, got.items, 1, "the setter replaces the sequence");
  });

  test("appendItems keeps what items holds", ({ seat }) => {
    let want: unknown[];
    want = ['items-a'];
    const got = PageBuilder.create<unknown>().withItems(...want).appendItems(...want).build();
    check.length(seat, got.items, 2, "the appender adds to the sequence");
  });
});

describe("PageBuilder clone", () => {
  test("items is not shared with a clone", ({ seat }) => {
    let want: unknown[];
    let other: unknown[];
    want = ['items-a'];
    other = ['items-b'];
    const original = PageBuilder.create<unknown>().withItems(...want);
    const clone = original.clone();
    clone.mutate((value) => {
      value.items.splice(0, other.length, ...other);
    });
    check.notEqual(seat, original.build().items, other, "writing through a clone must not reach the original");
  });

  test("clone copies items", ({ seat }) => {
    let want: unknown[];
    want = ['items-a'];
    const original = PageBuilder.create<unknown>().withItems(...want);
    check.equal(seat, original.clone().build().items, want, "a clone builds the same member as its original");
  });
});

describe("PageBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: Page<unknown>[] = [];
    const record = (value: Page<unknown>): void => {
      seen.push(value);
    };
    PageBuilder.create<unknown>().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = PageBuilder.create<unknown>();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, PageBuilder.from<unknown>(want).build(), want, "from keeps the value it was given");
  });
});

describe("CounterBuilder setters", () => {
  test("withCount returns its builder", ({ seat }) => {
    const builder = CounterBuilder.create();
    check.equal(
      seat,
      builder.withCount(builder.build().count),
      builder,
      "the setter returns its builder for the next call",
      byIdentity(),
    );
  });
});

describe("CounterBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: Counter[] = [];
    const record = (value: Counter): void => {
      seen.push(value);
    };
    CounterBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = CounterBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, CounterBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});

describe("HeaderBuilder construction", () => {
  test("from round-trips what it was seeded with", ({ seat }) => {
    const seed = new Header();
    seed.v2 = 'v2-a';
    seed.$ref = '$ref-a';
    const got = HeaderBuilder.from(seed).build();
    check.equal(seat, got.v2, seed.v2, "from keeps the member it was given");
    check.equal(seat, got.$ref, seed.$ref, "from keeps the member it was given");
  });
});

describe("HeaderBuilder setters", () => {
  test("withV2 writes v2", ({ seat }) => {
    let want: string;
    want = 'v2-a';
    const got = HeaderBuilder.create().withV2(want).build();
    check.equal(seat, got.v2, want, "the setter writes what it was given");
  });

  test("withV2 replaces what was already there", ({ seat }) => {
    let first: string;
    let want: string;
    first = 'v2-a';
    want = 'v2-b';
    const got = HeaderBuilder.create().withV2(first).withV2(want).build();
    check.equal(seat, got.v2, want, "the later call wins");
  });

  test("with$ref writes $ref", ({ seat }) => {
    let want: string;
    want = '$ref-a';
    const got = HeaderBuilder.create().with$ref(want).build();
    check.equal(seat, got.$ref, want, "the setter writes what it was given");
  });

  test("with$ref replaces what was already there", ({ seat }) => {
    let first: string;
    let want: string;
    first = '$ref-a';
    want = '$ref-b';
    const got = HeaderBuilder.create().with$ref(first).with$ref(want).build();
    check.equal(seat, got.$ref, want, "the later call wins");
  });
});

describe("HeaderBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: Header[] = [];
    const record = (value: Header): void => {
      seen.push(value);
    };
    HeaderBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = HeaderBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, HeaderBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});

describe("BinBuilder setters", () => {
  test("withCodes writes codes", ({ seat }) => {
    let want: ReadonlyArray<string>;
    want = ['codes-a'];
    const got = BinBuilder.create().withCodes(...want).build();
    check.length(seat, got.codes, 1, "the setter replaces the sequence");
  });

  test("withMarks writes marks", ({ seat }) => {
    let want: Set<string>;
    want = new Set<string>(['marks-a']);
    const got = BinBuilder.create().withMarks(want).build();
    check.isNotNil(seat, got.marks, "the setter writes the value it was given");
  });

  test("withNotes writes notes", ({ seat }) => {
    let want: string[];
    want = ['notes-a'];
    const got = BinBuilder.create().withNotes(want).build();
    check.isNotNil(seat, got.notes, "the setter writes the value it was given");
  });

  test("withAliases writes aliases", ({ seat }) => {
    let want: Lines;
    want = ['aliases-a'];
    const got = BinBuilder.create().withAliases(...want).build();
    check.length(seat, got.aliases, 1, "the setter replaces the sequence");
  });

  test("withRaw writes raw", ({ seat }) => {
    let want: Uint8Array;
    want = new Uint8Array([1]);
    const got = BinBuilder.create().withRaw(want).build();
    check.length(seat, got.raw, 1, "the setter replaces the sequence");
  });

  test("appendCodes keeps what codes holds", ({ seat }) => {
    let want: ReadonlyArray<string>;
    want = ['codes-a'];
    const got = BinBuilder.create().withCodes(...want).appendCodes(...want).build();
    check.length(seat, got.codes, 2, "the appender adds to the sequence");
  });

  test("appendAliases keeps what aliases holds", ({ seat }) => {
    let want: Lines;
    want = ['aliases-a'];
    const got = BinBuilder.create().withAliases(...want).appendAliases(...want).build();
    check.length(seat, got.aliases, 2, "the appender adds to the sequence");
  });
});

describe("BinBuilder clone", () => {
  test("aliases is not shared with a clone", ({ seat }) => {
    let want: Lines;
    let other: Lines;
    want = ['aliases-a'];
    other = ['aliases-b'];
    const original = BinBuilder.create().withAliases(...want);
    const clone = original.clone();
    clone.mutate((value) => {
      value.aliases.splice(0, other.length, ...other);
    });
    check.notEqual(seat, original.build().aliases, other, "writing through a clone must not reach the original");
  });

  test("raw is not shared with a clone", ({ seat }) => {
    let want: Uint8Array;
    let other: Uint8Array;
    want = new Uint8Array([1]);
    other = new Uint8Array([2]);
    const original = BinBuilder.create().withRaw(want);
    const clone = original.clone();
    clone.mutate((value) => {
      value.raw.set(other);
    });
    check.notEqual(seat, original.build().raw, other, "writing through a clone must not reach the original");
  });

  test("clone copies aliases", ({ seat }) => {
    let want: Lines;
    want = ['aliases-a'];
    const original = BinBuilder.create().withAliases(...want);
    check.equal(seat, original.clone().build().aliases, want, "a clone builds the same member as its original");
  });

  test("clone copies raw", ({ seat }) => {
    let want: Uint8Array;
    want = new Uint8Array([1]);
    const original = BinBuilder.create().withRaw(want);
    check.equal(seat, original.clone().build().raw, want, "a clone builds the same member as its original");
  });
});

describe("BinBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: Bin[] = [];
    const record = (value: Bin): void => {
      seen.push(value);
    };
    BinBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = BinBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, BinBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});

describe("LooseBuilder setters", () => {
  test("withItems returns its builder", ({ seat }) => {
    const builder = LooseBuilder.create();
    check.equal(
      seat,
      builder.withItems(...builder.build().items),
      builder,
      "the setter returns its builder for the next call",
      byIdentity(),
    );
  });

  test("appendItems returns its builder", ({ seat }) => {
    const builder = LooseBuilder.create();
    check.equal(
      seat,
      builder.appendItems(),
      builder,
      "the setter returns its builder for the next call",
      byIdentity(),
    );
  });

  test("withIndex returns its builder", ({ seat }) => {
    const builder = LooseBuilder.create();
    check.equal(
      seat,
      builder.withIndex(builder.build().index),
      builder,
      "the setter returns its builder for the next call",
      byIdentity(),
    );
  });

  test("withIndexEntry returns its builder", ({ seat }) => {
    const builder = LooseBuilder.create();
    check.equal(
      seat,
      builder.withIndexEntry(undefined as never, undefined as never),
      builder,
      "the setter returns its builder for the next call",
      byIdentity(),
    );
  });

  test("withIndexEntries returns its builder", ({ seat }) => {
    const builder = LooseBuilder.create();
    check.equal(
      seat,
      builder.withIndexEntries(builder.build().index),
      builder,
      "the setter returns its builder for the next call",
      byIdentity(),
    );
  });

  test("withKeys returns its builder", ({ seat }) => {
    const builder = LooseBuilder.create();
    check.equal(
      seat,
      builder.withKeys(builder.build().keys),
      builder,
      "the setter returns its builder for the next call",
      byIdentity(),
    );
  });

  test("withKeysEntry returns its builder", ({ seat }) => {
    const builder = LooseBuilder.create();
    check.equal(
      seat,
      builder.withKeysEntry(undefined as never),
      builder,
      "the setter returns its builder for the next call",
      byIdentity(),
    );
  });

  test("withKeysEntries returns its builder", ({ seat }) => {
    const builder = LooseBuilder.create();
    check.equal(
      seat,
      builder.withKeysEntries(),
      builder,
      "the setter returns its builder for the next call",
      byIdentity(),
    );
  });

  test("withExtra returns its builder", ({ seat }) => {
    const builder = LooseBuilder.create();
    check.equal(
      seat,
      builder.withExtra(undefined as never),
      builder,
      "the setter returns its builder for the next call",
      byIdentity(),
    );
  });

  test("withValue returns its builder", ({ seat }) => {
    const builder = LooseBuilder.create();
    check.equal(
      seat,
      builder.withValue(builder.build().value),
      builder,
      "the setter returns its builder for the next call",
      byIdentity(),
    );
  });
});

describe("LooseBuilder methods", () => {
  test("mutate passes the value under construction", ({ seat }) => {
    const seen: Loose[] = [];
    const record = (value: Loose): void => {
      seen.push(value);
    };
    LooseBuilder.create().mutate(record).mutate(record);
    check.length(seat, seen, 2, "mutate calls each function once");
    check.equal(seat, seen[0], seen[1], "the two functions get the same value", byIdentity());
  });

  test("clone and from keep the value under construction", ({ seat }) => {
    const original = LooseBuilder.create();
    const want = original.build();
    const clone = original.clone();
    check.notEqual(seat, clone, original, "clone returns a builder of its own", byIdentity());
    check.equal(seat, clone.build(), want, "a clone builds the value of its original");
    check.equal(seat, LooseBuilder.from(want).build(), want, "from keeps the value it was given");
  });
});
