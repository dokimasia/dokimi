import { Account, Address, Bin, Counter, Derived, HTTPClient, Header, Heir, Loose, MD5Hash, type Notifier, Page, Tier, UUID, UserUrl, accountDefaults, httpClientDefaults, md5HashDefaults, userUrlDefaults, uuidDefaults } from "./catalog";

/**
 * AddressBuilder constructs Address values. AddressBuilder.create and
 * AddressBuilder.from return a new builder. Each setter sets one member
 * and returns the builder, so the calls chain. build returns the value.
 */
export class AddressBuilder {
  private readonly v: Address;

  private constructor(v: Address) {
    this.v = v;
  }

  /**
   * create returns a builder of the declared defaults. A member without a
   * declared default has the value of a new Address.
   */
  static create(): AddressBuilder {
    const v = new Address();
    v.street = "Dam 1";
    return new AddressBuilder(v);
  }

  /** from returns a builder of a copy of value. */
  static from(value: Address): AddressBuilder {
    return new AddressBuilder(Object.assign(new Address(), value));
  }

  /** withStreet replaces street with value. */
  withStreet(value: string): this {
    this.v.street = value;
    return this;
  }

  /** withCity replaces city with value. */
  withCity(value: string): this {
    this.v.city = value;
    return this;
  }

  /** withFloor replaces floor with value. mutate can set floor to undefined. */
  withFloor(value: number): this {
    this.v.floor = value;
    return this;
  }

  /**
   * withLines replaces the elements of lines with values. A call without
   * arguments empties lines.
   */
  withLines(...values: string[]): this {
    this.v.lines = values;
    return this;
  }

  /** appendLines appends values to the elements of lines. */
  appendLines(...values: string[]): this {
    this.v.lines = [...this.v.lines, ...values];
    return this;
  }

  /**
   * withLabels replaces labels with value. The builder keeps value without a
   * copy.
   */
  withLabels(value: Map<string, string>): this {
    this.v.labels = value;
    return this;
  }

  /** withLabelsEntry sets the entry of key in labels to value. */
  withLabelsEntry(key: string, value: string): this {
    this.v.labels = new Map([...this.v.labels, [key, value]]);
    return this;
  }

  /**
   * withLabelsEntries copies each entry of entries into labels. A key of entries
   * replaces the value of that key in labels.
   */
  withLabelsEntries(entries: Map<string, string>): this {
    this.v.labels = new Map([...this.v.labels, ...entries]);
    return this;
  }

  /**
   * withTags replaces tags with value. The builder keeps value without a
   * copy.
   */
  withTags(value: Set<string>): this {
    this.v.tags = value;
    return this;
  }

  /** withTagsEntry adds key to tags. */
  withTagsEntry(key: string): this {
    this.v.tags = new Set([...this.v.tags, key]);
    return this;
  }

  /** withTagsEntries adds each key of keys to tags. */
  withTagsEntries(...keys: string[]): this {
    this.v.tags = new Set([...this.v.tags, ...keys]);
    return this;
  }

  /**
   * mutate applies fn to the value under construction. It changes what no
   * setter changes, such as several members in one step.
   */
  mutate(fn: (value: Address) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): AddressBuilder {
    const out = new AddressBuilder(Object.assign(new Address(), this.v));
    out.v.lines = [...out.v.lines];
    out.v.labels = new Map(out.v.labels);
    out.v.tags = new Set(out.v.tags);
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): Address {
    return Object.assign(new Address(), this.v);
  }
}

/**
 * DerivedBuilder constructs Derived values. DerivedBuilder.create and
 * DerivedBuilder.from return a new builder. Each setter sets one member
 * and returns the builder, so the calls chain. build returns the value.
 */
export class DerivedBuilder {
  private readonly v: Derived;

  private constructor(v: Derived) {
    this.v = v;
  }

  /**
   * create returns a builder of the declared defaults. A member without a
   * declared default has the value of a new Derived.
   */
  static create(): DerivedBuilder {
    const v = new Derived();
    v.createdBy = "system";
    return new DerivedBuilder(v);
  }

  /** from returns a builder of a copy of value. */
  static from(value: Derived): DerivedBuilder {
    return new DerivedBuilder(Object.assign(new Derived(), value));
  }

  /** withActive replaces active with value. */
  withActive(value: boolean): this {
    this.v.active = value;
    return this;
  }

  /** withCreatedBy replaces createdBy with value. */
  withCreatedBy(value: string): this {
    this.v.createdBy = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction. It changes what no
   * setter changes, such as several members in one step.
   */
  mutate(fn: (value: Derived) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): DerivedBuilder {
    const out = new DerivedBuilder(Object.assign(new Derived(), this.v));
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): Derived {
    return Object.assign(new Derived(), this.v);
  }
}

/**
 * HeirBuilder constructs Heir values. HeirBuilder.create and
 * HeirBuilder.from return a new builder. Each setter sets one member and
 * returns the builder, so the calls chain. build returns the value.
 */
export class HeirBuilder {
  private readonly v: Heir;

  private constructor(v: Heir) {
    this.v = v;
  }

  /**
   * create returns a builder of the declared defaults. A member without a
   * declared default has the value of a new Heir.
   */
  static create(): HeirBuilder {
    const v = new Heir();
    v.createdBy = "system";
    return new HeirBuilder(v);
  }

  /** from returns a builder of a copy of value. */
  static from(value: Heir): HeirBuilder {
    return new HeirBuilder(Object.assign(new Heir(), value));
  }

  /** withActive replaces active with value. */
  withActive(value: boolean): this {
    this.v.active = value;
    return this;
  }

  /** withCreatedBy replaces createdBy with value. */
  withCreatedBy(value: string): this {
    this.v.createdBy = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction. It changes what no
   * setter changes, such as several members in one step.
   */
  mutate(fn: (value: Heir) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): HeirBuilder {
    const out = new HeirBuilder(Object.assign(new Heir(), this.v));
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): Heir {
    return Object.assign(new Heir(), this.v);
  }
}

/**
 * AccountBuilder constructs Account values. AccountBuilder.create and
 * AccountBuilder.from return a new builder. Each setter sets one member
 * and returns the builder, so the calls chain. build returns the value.
 */
export class AccountBuilder {
  private readonly v: Account;

  private constructor(v: Account) {
    this.v = v;
  }

  /**
   * create returns a builder of the value that the companion function
   * returns. The declared default of each member applies over that value.
   */
  static create(): AccountBuilder {
    const v = accountDefaults();
    v.owner = "owner";
    return new AccountBuilder(v);
  }

  /** from returns a builder of a copy of value. */
  static from(value: Account): AccountBuilder {
    return new AccountBuilder(Object.assign(new Account(), value));
  }

  /** withOwner replaces owner with value. */
  withOwner(value: string): this {
    this.v.owner = value;
    return this;
  }

  /** withPlan replaces plan with value. */
  withPlan(value: string): this {
    this.v.plan = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction. It changes what no
   * setter changes, such as several members in one step.
   */
  mutate(fn: (value: Account) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): AccountBuilder {
    const out = new AccountBuilder(Object.assign(new Account(), this.v));
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): Account {
    return Object.assign(new Account(), this.v);
  }
}

/**
 * HTTPClientBuilder constructs HTTPClient values. HTTPClientBuilder.create
 * and HTTPClientBuilder.from return a new builder. Each setter sets one
 * member and returns the builder, so the calls chain. build returns the
 * value.
 */
export class HTTPClientBuilder {
  private readonly v: HTTPClient;

  private constructor(v: HTTPClient) {
    this.v = v;
  }

  /**
   * create returns a builder of the value that the companion function
   * returns. The declared default of each member applies over that value.
   */
  static create(): HTTPClientBuilder {
    const v = httpClientDefaults();
    return new HTTPClientBuilder(v);
  }

  /** from returns a builder of a copy of value. */
  static from(value: HTTPClient): HTTPClientBuilder {
    return new HTTPClientBuilder(Object.assign(new HTTPClient(), value));
  }

  /** withBase replaces base with value. */
  withBase(value: string): this {
    this.v.base = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction. It changes what no
   * setter changes, such as several members in one step.
   */
  mutate(fn: (value: HTTPClient) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): HTTPClientBuilder {
    const out = new HTTPClientBuilder(Object.assign(new HTTPClient(), this.v));
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): HTTPClient {
    return Object.assign(new HTTPClient(), this.v);
  }
}

/**
 * Md5HashBuilder constructs MD5Hash values. Md5HashBuilder.create and
 * Md5HashBuilder.from return a new builder. Each setter sets one member
 * and returns the builder, so the calls chain. build returns the value.
 */
export class Md5HashBuilder {
  private readonly v: MD5Hash;

  private constructor(v: MD5Hash) {
    this.v = v;
  }

  /**
   * create returns a builder of the value that the companion function
   * returns. The declared default of each member applies over that value.
   */
  static create(): Md5HashBuilder {
    const v = md5HashDefaults();
    return new Md5HashBuilder(v);
  }

  /** from returns a builder of a copy of value. */
  static from(value: MD5Hash): Md5HashBuilder {
    return new Md5HashBuilder(Object.assign(new MD5Hash(), value));
  }

  /** withDigest replaces digest with value. */
  withDigest(value: string): this {
    this.v.digest = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction. It changes what no
   * setter changes, such as several members in one step.
   */
  mutate(fn: (value: MD5Hash) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): Md5HashBuilder {
    const out = new Md5HashBuilder(Object.assign(new MD5Hash(), this.v));
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): MD5Hash {
    return Object.assign(new MD5Hash(), this.v);
  }
}

/**
 * UUIDBuilder constructs UUID values. UUIDBuilder.create and
 * UUIDBuilder.from return a new builder. Each setter sets one member and
 * returns the builder, so the calls chain. build returns the value.
 */
export class UUIDBuilder {
  private readonly v: UUID;

  private constructor(v: UUID) {
    this.v = v;
  }

  /**
   * create returns a builder of the value that the companion function
   * returns. The declared default of each member applies over that value.
   */
  static create(): UUIDBuilder {
    const v = uuidDefaults();
    return new UUIDBuilder(v);
  }

  /** from returns a builder of a copy of value. */
  static from(value: UUID): UUIDBuilder {
    return new UUIDBuilder(Object.assign(new UUID(), value));
  }

  /** withValue replaces value with value. */
  withValue(value: string): this {
    this.v.value = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction. It changes what no
   * setter changes, such as several members in one step.
   */
  mutate(fn: (value: UUID) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): UUIDBuilder {
    const out = new UUIDBuilder(Object.assign(new UUID(), this.v));
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): UUID {
    return Object.assign(new UUID(), this.v);
  }
}

/**
 * UserUrlBuilder constructs UserUrl values. UserUrlBuilder.create and
 * UserUrlBuilder.from return a new builder. Each setter sets one member
 * and returns the builder, so the calls chain. build returns the value.
 */
export class UserUrlBuilder {
  private readonly v: UserUrl;

  private constructor(v: UserUrl) {
    this.v = v;
  }

  /**
   * create returns a builder of the value that the companion function
   * returns. The declared default of each member applies over that value.
   */
  static create(): UserUrlBuilder {
    const v = userUrlDefaults();
    return new UserUrlBuilder(v);
  }

  /** from returns a builder of a copy of value. */
  static from(value: UserUrl): UserUrlBuilder {
    return new UserUrlBuilder(Object.assign(new UserUrl(), value));
  }

  /** withHref replaces href with value. */
  withHref(value: string): this {
    this.v.href = value;
    return this;
  }

  /** withUserId replaces userId with value. */
  withUserId(value: string): this {
    this.v.userId = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction. It changes what no
   * setter changes, such as several members in one step.
   */
  mutate(fn: (value: UserUrl) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): UserUrlBuilder {
    const out = new UserUrlBuilder(Object.assign(new UserUrl(), this.v));
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): UserUrl {
    return Object.assign(new UserUrl(), this.v);
  }
}

/**
 * TierBuilder constructs Tier values. TierBuilder.create and
 * TierBuilder.from return a new builder. Each setter sets one member and
 * returns the builder, so the calls chain. build returns the value.
 */
export class TierBuilder {
  private readonly v: Tier;

  private constructor(v: Tier) {
    this.v = v;
  }

  /** create returns a builder of a new Tier. */
  static create(): TierBuilder {
    return new TierBuilder(new Tier());
  }

  /** from returns a builder of a copy of value. */
  static from(value: Tier): TierBuilder {
    return new TierBuilder(Object.assign(new Tier(), value));
  }

  /** withLevel replaces level with value. */
  withLevel(value: number): this {
    this.v.level = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction. It changes what no
   * setter changes, such as several members in one step.
   */
  mutate(fn: (value: Tier) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): TierBuilder {
    const out = new TierBuilder(Object.assign(new Tier(), this.v));
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): Tier {
    return Object.assign(new Tier(), this.v);
  }
}

/**
 * PageBuilder constructs Page values. PageBuilder.create and
 * PageBuilder.from return a new builder. Each setter sets one member and
 * returns the builder, so the calls chain. build returns the value.
 */
export class PageBuilder<T> {
  private readonly v: Page<T>;

  private constructor(v: Page<T>) {
    this.v = v;
  }

  /** create returns a builder of a new Page. */
  static create<T>(): PageBuilder<T> {
    return new PageBuilder<T>(new Page<T>());
  }

  /** from returns a builder of a copy of value. */
  static from<T>(value: Page<T>): PageBuilder<T> {
    return new PageBuilder<T>(Object.assign(new Page<T>(), value));
  }

  /**
   * withItems replaces the elements of items with values. A call without
   * arguments empties items.
   */
  withItems(...values: T[]): this {
    this.v.items = values;
    return this;
  }

  /** appendItems appends values to the elements of items. */
  appendItems(...values: T[]): this {
    this.v.items = [...this.v.items, ...values];
    return this;
  }

  /**
   * mutate applies fn to the value under construction. It changes what no
   * setter changes, such as several members in one step.
   */
  mutate(fn: (value: Page<T>) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): PageBuilder<T> {
    const out = new PageBuilder<T>(Object.assign(new Page<T>(), this.v));
    out.v.items = [...out.v.items];
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): Page<T> {
    return Object.assign(new Page<T>(), this.v);
  }
}

/**
 * CounterBuilder constructs Counter values. CounterBuilder.create and
 * CounterBuilder.from return a new builder. Each setter sets one member
 * and returns the builder, so the calls chain. build returns the value.
 */
export class CounterBuilder {
  private readonly v: Counter;

  private constructor(v: Counter) {
    this.v = v;
  }

  /** create returns a builder of a new Counter. */
  static create(): CounterBuilder {
    return new CounterBuilder(new Counter());
  }

  /** from returns a builder of a copy of value. */
  static from(value: Counter): CounterBuilder {
    return new CounterBuilder(Object.assign(new Counter(), value));
  }

  /** withCount replaces count with value. */
  withCount(value: Counter["count"]): this {
    this.v.count = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction. It changes what no
   * setter changes, such as several members in one step.
   */
  mutate(fn: (value: Counter) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): CounterBuilder {
    const out = new CounterBuilder(Object.assign(new Counter(), this.v));
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): Counter {
    return Object.assign(new Counter(), this.v);
  }
}

/**
 * HeaderBuilder constructs Header values. HeaderBuilder.create and
 * HeaderBuilder.from return a new builder. Each setter sets one member and
 * returns the builder, so the calls chain. build returns the value.
 */
export class HeaderBuilder {
  private readonly v: Header;

  private constructor(v: Header) {
    this.v = v;
  }

  /** create returns a builder of a new Header. */
  static create(): HeaderBuilder {
    return new HeaderBuilder(new Header());
  }

  /** from returns a builder of a copy of value. */
  static from(value: Header): HeaderBuilder {
    return new HeaderBuilder(Object.assign(new Header(), value));
  }

  /** withV2 replaces v2 with value. */
  withV2(value: string): this {
    this.v.v2 = value;
    return this;
  }

  /** with$ref replaces $ref with value. */
  with$ref(value: string): this {
    this.v.$ref = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction. It changes what no
   * setter changes, such as several members in one step.
   */
  mutate(fn: (value: Header) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): HeaderBuilder {
    const out = new HeaderBuilder(Object.assign(new Header(), this.v));
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): Header {
    return Object.assign(new Header(), this.v);
  }
}

/**
 * BinBuilder constructs Bin values. BinBuilder.create and BinBuilder.from
 * return a new builder. Each setter sets one member and returns the
 * builder, so the calls chain. build returns the value.
 */
export class BinBuilder {
  private readonly v: Bin;

  private constructor(v: Bin) {
    this.v = v;
  }

  /** create returns a builder of a new Bin. */
  static create(): BinBuilder {
    return new BinBuilder(new Bin());
  }

  /** from returns a builder of a copy of value. */
  static from(value: Bin): BinBuilder {
    return new BinBuilder(Object.assign(new Bin(), value));
  }

  /**
   * withCodes replaces the elements of codes with values. A call without
   * arguments empties codes.
   */
  withCodes(...values: string[]): this {
    this.v.codes = values;
    return this;
  }

  /** appendCodes appends values to the elements of codes. */
  appendCodes(...values: string[]): this {
    this.v.codes = [...this.v.codes, ...values];
    return this;
  }

  /** withMarks replaces marks with value. mutate can set marks to undefined. */
  withMarks(value: Set<string>): this {
    this.v.marks = value;
    return this;
  }

  /** withNotes replaces notes with value. mutate can set notes to undefined. */
  withNotes(value: string[]): this {
    this.v.notes = value;
    return this;
  }

  /**
   * withAliases replaces the elements of aliases with values. A call without
   * arguments empties aliases.
   */
  withAliases(...values: string[]): this {
    this.v.aliases = values;
    return this;
  }

  /** appendAliases appends values to the elements of aliases. */
  appendAliases(...values: string[]): this {
    this.v.aliases = [...this.v.aliases, ...values];
    return this;
  }

  /**
   * withRaw replaces raw with value. The builder keeps value without a
   * copy.
   */
  withRaw(value: Uint8Array): this {
    this.v.raw = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction. It changes what no
   * setter changes, such as several members in one step.
   */
  mutate(fn: (value: Bin) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): BinBuilder {
    const out = new BinBuilder(Object.assign(new Bin(), this.v));
    out.v.aliases = [...out.v.aliases];
    out.v.raw = out.v.raw.slice();
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): Bin {
    return Object.assign(new Bin(), this.v);
  }
}

/**
 * LooseBuilder constructs Loose values. LooseBuilder.create and
 * LooseBuilder.from return a new builder. Each setter sets one member and
 * returns the builder, so the calls chain. build returns the value.
 */
export class LooseBuilder {
  private readonly v: Loose;

  private constructor(v: Loose) {
    this.v = v;
  }

  /** create returns a builder of a new Loose. */
  static create(): LooseBuilder {
    return new LooseBuilder(new Loose());
  }

  /** from returns a builder of a copy of value. */
  static from(value: Loose): LooseBuilder {
    return new LooseBuilder(Object.assign(new Loose(), value));
  }

  /**
   * withItems replaces the elements of items with values. A call without
   * arguments empties items.
   */
  withItems(...values: Notifier[]): this {
    this.v.items = values;
    return this;
  }

  /** appendItems appends values to the elements of items. */
  appendItems(...values: Notifier[]): this {
    this.v.items = [...this.v.items, ...values];
    return this;
  }

  /**
   * withIndex replaces index with value. The builder keeps value without a
   * copy.
   */
  withIndex(value: Map<string, Notifier>): this {
    this.v.index = value;
    return this;
  }

  /** withIndexEntry sets the entry of key in index to value. */
  withIndexEntry(key: string, value: Notifier): this {
    this.v.index = new Map([...this.v.index, [key, value]]);
    return this;
  }

  /**
   * withIndexEntries copies each entry of entries into index. A key of entries
   * replaces the value of that key in index.
   */
  withIndexEntries(entries: Map<string, Notifier>): this {
    this.v.index = new Map([...this.v.index, ...entries]);
    return this;
  }

  /**
   * withKeys replaces keys with value. The builder keeps value without a
   * copy.
   */
  withKeys(value: Set<Notifier>): this {
    this.v.keys = value;
    return this;
  }

  /** withKeysEntry adds key to keys. */
  withKeysEntry(key: Notifier): this {
    this.v.keys = new Set([...this.v.keys, key]);
    return this;
  }

  /** withKeysEntries adds each key of keys to keys. */
  withKeysEntries(...keys: Notifier[]): this {
    this.v.keys = new Set([...this.v.keys, ...keys]);
    return this;
  }

  /** withExtra replaces extra with value. mutate can set extra to undefined. */
  withExtra(value: Notifier): this {
    this.v.extra = value;
    return this;
  }

  /** withValue replaces value with value. */
  withValue(value: Notifier): this {
    this.v.value = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction. It changes what no
   * setter changes, such as several members in one step.
   */
  mutate(fn: (value: Loose) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): LooseBuilder {
    const out = new LooseBuilder(Object.assign(new Loose(), this.v));
    out.v.items = [...out.v.items];
    out.v.index = new Map(out.v.index);
    out.v.keys = new Set(out.v.keys);
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): Loose {
    return Object.assign(new Loose(), this.v);
  }
}
