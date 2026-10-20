import { Account, Address, Bin, Counter, Derived, HTTPClient, Header, Heir, Loose, MD5Hash, type Notifier, Page, Tier, UUID, UserUrl, accountDefaults, httpClientDefaults, md5HashDefaults, userUrlDefaults, uuidDefaults } from './catalog';

/**
 * AddressBuilder constructs Address values.
 *
 * Start with AddressBuilder.create or AddressBuilder.from, set what you
 * care about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to Address breaks every
 * literal at once and a reader cannot tell which members a given call site
 * actually cares about. Here the constructor supplies the rest, and each
 * call states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class AddressBuilder {
  private readonly v: Address;

  private constructor(v: Address) {
    this.v = v;
  }

  /**
   * create returns a builder seeded with the declared defaults.
   *
   * Only the members that declare one. Everything else holds the value that
   * a new Address has, so a caller that sets nothing gets a mixture
   * rather than either extreme — which is what makes a check against a built
   * value worth writing.
   */
  static create(): AddressBuilder {
    const v = new Address();
    v.street = 'Dam 1';
    return new AddressBuilder(v);
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from(value: Address): AddressBuilder {
    return new AddressBuilder(Object.assign(new Address(), value));
  }

  /**
   * withStreet replaces street.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withStreet(value: string): this {
    this.v.street = value;
    return this;
  }

  /**
   * withCity replaces city.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withCity(value: string): this {
    this.v.city = value;
    return this;
  }

  /**
   * withFloor replaces floor.
   *
   * Use mutate to clear the member.
   */
  withFloor(value: number): this {
    this.v.floor = value;
    return this;
  }

  /**
   * withLines replaces lines with the elements given.
   *
   * Variadic, and replacing: whatever the collection held is discarded,
   * including anything a constructor seeded. appendLines is the counterpart
   * that keeps it. Called with no arguments this empties the member rather
   * than leaving it alone, which is how a caller clears a seeded default.
   */
  withLines(...values: string[]): this {
    this.v.lines = values;
    return this;
  }

  /**
   * appendLines adds to lines, keeping what is already there.
   *
   * The counterpart to withLines, which discards it. Both are offered
   * because a caller extending a seeded default and a caller replacing it
   * want opposite things, and one method cannot mean both.
   */
  appendLines(...values: string[]): this {
    this.v.lines = [...this.v.lines, ...values];
    return this;
  }

  /**
   * withLabels replaces labels wholesale.
   *
   * The whole mapping, not one entry: what the member held is discarded
   * rather than merged into, and the map is kept rather than copied — so a
   * caller that writes to it afterwards changes what was already built.
   * withLabelsEntry is what sets a single key.
   */
  withLabels(value: Map<string, string>): this {
    this.v.labels = value;
    return this;
  }

  /**
   * withLabelsEntry sets one key of labels.
   *
   * A key already present takes the new value.
   */
  withLabelsEntry(key: string, value: string): this {
    this.v.labels = new Map([...this.v.labels, [key, value]]);
    return this;
  }

  /**
   * withLabelsEntries merges entries into labels.
   *
   * Merges rather than replaces: a key already present takes the new value,
   * and keys absent from entries are left alone. Use withLabels to discard
   * what is there.
   */
  withLabelsEntries(entries: Map<string, string>): this {
    this.v.labels = new Map([...this.v.labels, ...entries]);
    return this;
  }

  /**
   * withTags replaces tags wholesale.
   *
   * The whole collection, not one key: what the member held is discarded
   * rather than merged into, and the set is kept rather than copied.
   * withTagsEntry is what adds to what is already there.
   */
  withTags(value: Set<string>): this {
    this.v.tags = value;
    return this;
  }

  /**
   * withTagsEntry adds one key to tags.
   *
   * A set entry is a key, so there is nothing else for a caller to supply.
   */
  withTagsEntry(key: string): this {
    this.v.tags = new Set([...this.v.tags, key]);
    return this;
  }

  /**
   * withTagsEntries adds each key to tags, keeping those it holds.
   *
   * Additive like withTagsEntry, so two calls union rather than the second
   * replacing the first. A key already present is left as it stands.
   */
  withTagsEntries(...keys: string[]): this {
    this.v.tags = new Set([...this.v.tags, ...keys]);
    return this;
  }

  /**
   * mutate applies fn to the value under construction.
   *
   * The way out for what the setters do not reach: a member this builder
   * cannot name, a shape it does not model, or a change that touches several
   * members at once and only makes sense as one edit.
   *
   * Prefer a setter where one exists. A setter says which member it writes,
   * and a reader of the call site learns that without following a closure.
   */
  mutate(fn: (value: Address) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder that can be configured without affecting this
   * one.
   *
   * One level deep, and the boundary matters. The array, byte-array and keyed
   * members are copied, so writing through one clone is not visible through
   * another. What is shared is anything held inside those, and every other
   * member that refers to an object. Use mutate to replace such a member
   * outright rather than editing it in place.
   */
  clone(): AddressBuilder {
    const out = new AddressBuilder(Object.assign(new Address(), this.v));
    out.v.lines = [...out.v.lines];
    out.v.labels = new Map(out.v.labels);
    out.v.tags = new Set(out.v.tags);
    return out;
  }

  /**
   * build returns the constructed Address.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): Address {
    return Object.assign(new Address(), this.v);
  }
}

/**
 * DerivedBuilder constructs Derived values.
 *
 * Start with DerivedBuilder.create or DerivedBuilder.from, set what you
 * care about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to Derived breaks every
 * literal at once and a reader cannot tell which members a given call site
 * actually cares about. Here the constructor supplies the rest, and each
 * call states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class DerivedBuilder {
  private readonly v: Derived;

  private constructor(v: Derived) {
    this.v = v;
  }

  /**
   * create returns a builder seeded with the declared defaults.
   *
   * Only the members that declare one. Everything else holds the value that
   * a new Derived has, so a caller that sets nothing gets a mixture
   * rather than either extreme — which is what makes a check against a built
   * value worth writing.
   */
  static create(): DerivedBuilder {
    const v = new Derived();
    v.createdBy = 'system';
    return new DerivedBuilder(v);
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from(value: Derived): DerivedBuilder {
    return new DerivedBuilder(Object.assign(new Derived(), value));
  }

  /**
   * withActive replaces active.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withActive(value: boolean): this {
    this.v.active = value;
    return this;
  }

  /**
   * withCreatedBy replaces createdBy.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withCreatedBy(value: string): this {
    this.v.createdBy = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction.
   *
   * The way out for what the setters do not reach: a member this builder
   * cannot name, a shape it does not model, or a change that touches several
   * members at once and only makes sense as one edit.
   *
   * Prefer a setter where one exists. A setter says which member it writes,
   * and a reader of the call site learns that without following a closure.
   */
  mutate(fn: (value: Derived) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder that can be configured without affecting this
   * one.
   *
   * One level deep, and the boundary matters. The array, byte-array and keyed
   * members are copied, so writing through one clone is not visible through
   * another. What is shared is anything held inside those, and every other
   * member that refers to an object. Use mutate to replace such a member
   * outright rather than editing it in place.
   */
  clone(): DerivedBuilder {
    const out = new DerivedBuilder(Object.assign(new Derived(), this.v));
    return out;
  }

  /**
   * build returns the constructed Derived.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): Derived {
    return Object.assign(new Derived(), this.v);
  }
}

/**
 * HeirBuilder constructs Heir values.
 *
 * Start with HeirBuilder.create or HeirBuilder.from, set what you care
 * about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to Heir breaks every literal
 * at once and a reader cannot tell which members a given call site
 * actually cares about. Here the constructor supplies the rest, and each
 * call states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class HeirBuilder {
  private readonly v: Heir;

  private constructor(v: Heir) {
    this.v = v;
  }

  /**
   * create returns a builder seeded with the declared defaults.
   *
   * Only the members that declare one. Everything else holds the value that
   * a new Heir has, so a caller that sets nothing gets a mixture
   * rather than either extreme — which is what makes a check against a built
   * value worth writing.
   */
  static create(): HeirBuilder {
    const v = new Heir();
    v.createdBy = 'system';
    return new HeirBuilder(v);
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from(value: Heir): HeirBuilder {
    return new HeirBuilder(Object.assign(new Heir(), value));
  }

  /**
   * withActive replaces active.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withActive(value: boolean): this {
    this.v.active = value;
    return this;
  }

  /**
   * withCreatedBy replaces createdBy.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withCreatedBy(value: string): this {
    this.v.createdBy = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction.
   *
   * The way out for what the setters do not reach: a member this builder
   * cannot name, a shape it does not model, or a change that touches several
   * members at once and only makes sense as one edit.
   *
   * Prefer a setter where one exists. A setter says which member it writes,
   * and a reader of the call site learns that without following a closure.
   */
  mutate(fn: (value: Heir) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder that can be configured without affecting this
   * one.
   *
   * One level deep, and the boundary matters. The array, byte-array and keyed
   * members are copied, so writing through one clone is not visible through
   * another. What is shared is anything held inside those, and every other
   * member that refers to an object. Use mutate to replace such a member
   * outright rather than editing it in place.
   */
  clone(): HeirBuilder {
    const out = new HeirBuilder(Object.assign(new Heir(), this.v));
    return out;
  }

  /**
   * build returns the constructed Heir.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): Heir {
    return Object.assign(new Heir(), this.v);
  }
}

/**
 * AccountBuilder constructs Account values.
 *
 * Start with AccountBuilder.create or AccountBuilder.from, set what you
 * care about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to Account breaks every
 * literal at once and a reader cannot tell which members a given call site
 * actually cares about. Here the constructor supplies the rest, and each
 * call states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class AccountBuilder {
  private readonly v: Account;

  private constructor(v: Account) {
    this.v = v;
  }

  /**
   * create returns a builder seeded from the module's companion, with any
   * declared per-member default applied over it.
   *
   * Two layers, in that order: whatever the companion supplies, then each
   * member that declares its own default. A member the companion sets and
   * the declaration does not keeps the companion's value, so the two sources
   * compose rather than one replacing the other.
   */
  static create(): AccountBuilder {
    const v = accountDefaults();
    v.owner = 'owner';
    return new AccountBuilder(v);
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from(value: Account): AccountBuilder {
    return new AccountBuilder(Object.assign(new Account(), value));
  }

  /**
   * withOwner replaces owner.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withOwner(value: string): this {
    this.v.owner = value;
    return this;
  }

  /**
   * withPlan replaces plan.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withPlan(value: string): this {
    this.v.plan = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction.
   *
   * The way out for what the setters do not reach: a member this builder
   * cannot name, a shape it does not model, or a change that touches several
   * members at once and only makes sense as one edit.
   *
   * Prefer a setter where one exists. A setter says which member it writes,
   * and a reader of the call site learns that without following a closure.
   */
  mutate(fn: (value: Account) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder that can be configured without affecting this
   * one.
   *
   * One level deep, and the boundary matters. The array, byte-array and keyed
   * members are copied, so writing through one clone is not visible through
   * another. What is shared is anything held inside those, and every other
   * member that refers to an object. Use mutate to replace such a member
   * outright rather than editing it in place.
   */
  clone(): AccountBuilder {
    const out = new AccountBuilder(Object.assign(new Account(), this.v));
    return out;
  }

  /**
   * build returns the constructed Account.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): Account {
    return Object.assign(new Account(), this.v);
  }
}

/**
 * HTTPClientBuilder constructs HTTPClient values.
 *
 * Start with HTTPClientBuilder.create or HTTPClientBuilder.from, set what
 * you care about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to HTTPClient breaks every
 * literal at once and a reader cannot tell which members a given call site
 * actually cares about. Here the constructor supplies the rest, and each
 * call states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class HTTPClientBuilder {
  private readonly v: HTTPClient;

  private constructor(v: HTTPClient) {
    this.v = v;
  }

  /**
   * create returns a builder seeded from the module's companion, with any
   * declared per-member default applied over it.
   *
   * Two layers, in that order: whatever the companion supplies, then each
   * member that declares its own default. A member the companion sets and
   * the declaration does not keeps the companion's value, so the two sources
   * compose rather than one replacing the other.
   */
  static create(): HTTPClientBuilder {
    const v = httpClientDefaults();
    return new HTTPClientBuilder(v);
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from(value: HTTPClient): HTTPClientBuilder {
    return new HTTPClientBuilder(Object.assign(new HTTPClient(), value));
  }

  /**
   * withBase replaces base.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withBase(value: string): this {
    this.v.base = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction.
   *
   * The way out for what the setters do not reach: a member this builder
   * cannot name, a shape it does not model, or a change that touches several
   * members at once and only makes sense as one edit.
   *
   * Prefer a setter where one exists. A setter says which member it writes,
   * and a reader of the call site learns that without following a closure.
   */
  mutate(fn: (value: HTTPClient) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder that can be configured without affecting this
   * one.
   *
   * One level deep, and the boundary matters. The array, byte-array and keyed
   * members are copied, so writing through one clone is not visible through
   * another. What is shared is anything held inside those, and every other
   * member that refers to an object. Use mutate to replace such a member
   * outright rather than editing it in place.
   */
  clone(): HTTPClientBuilder {
    const out = new HTTPClientBuilder(Object.assign(new HTTPClient(), this.v));
    return out;
  }

  /**
   * build returns the constructed HTTPClient.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): HTTPClient {
    return Object.assign(new HTTPClient(), this.v);
  }
}

/**
 * Md5HashBuilder constructs MD5Hash values.
 *
 * Start with Md5HashBuilder.create or Md5HashBuilder.from, set what you
 * care about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to MD5Hash breaks every
 * literal at once and a reader cannot tell which members a given call site
 * actually cares about. Here the constructor supplies the rest, and each
 * call states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class Md5HashBuilder {
  private readonly v: MD5Hash;

  private constructor(v: MD5Hash) {
    this.v = v;
  }

  /**
   * create returns a builder seeded from the module's companion, with any
   * declared per-member default applied over it.
   *
   * Two layers, in that order: whatever the companion supplies, then each
   * member that declares its own default. A member the companion sets and
   * the declaration does not keeps the companion's value, so the two sources
   * compose rather than one replacing the other.
   */
  static create(): Md5HashBuilder {
    const v = md5HashDefaults();
    return new Md5HashBuilder(v);
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from(value: MD5Hash): Md5HashBuilder {
    return new Md5HashBuilder(Object.assign(new MD5Hash(), value));
  }

  /**
   * withDigest replaces digest.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withDigest(value: string): this {
    this.v.digest = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction.
   *
   * The way out for what the setters do not reach: a member this builder
   * cannot name, a shape it does not model, or a change that touches several
   * members at once and only makes sense as one edit.
   *
   * Prefer a setter where one exists. A setter says which member it writes,
   * and a reader of the call site learns that without following a closure.
   */
  mutate(fn: (value: MD5Hash) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder that can be configured without affecting this
   * one.
   *
   * One level deep, and the boundary matters. The array, byte-array and keyed
   * members are copied, so writing through one clone is not visible through
   * another. What is shared is anything held inside those, and every other
   * member that refers to an object. Use mutate to replace such a member
   * outright rather than editing it in place.
   */
  clone(): Md5HashBuilder {
    const out = new Md5HashBuilder(Object.assign(new MD5Hash(), this.v));
    return out;
  }

  /**
   * build returns the constructed MD5Hash.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): MD5Hash {
    return Object.assign(new MD5Hash(), this.v);
  }
}

/**
 * UUIDBuilder constructs UUID values.
 *
 * Start with UUIDBuilder.create or UUIDBuilder.from, set what you care
 * about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to UUID breaks every literal
 * at once and a reader cannot tell which members a given call site
 * actually cares about. Here the constructor supplies the rest, and each
 * call states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class UUIDBuilder {
  private readonly v: UUID;

  private constructor(v: UUID) {
    this.v = v;
  }

  /**
   * create returns a builder seeded from the module's companion, with any
   * declared per-member default applied over it.
   *
   * Two layers, in that order: whatever the companion supplies, then each
   * member that declares its own default. A member the companion sets and
   * the declaration does not keeps the companion's value, so the two sources
   * compose rather than one replacing the other.
   */
  static create(): UUIDBuilder {
    const v = uuidDefaults();
    return new UUIDBuilder(v);
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from(value: UUID): UUIDBuilder {
    return new UUIDBuilder(Object.assign(new UUID(), value));
  }

  /**
   * withValue replaces value.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withValue(value: string): this {
    this.v.value = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction.
   *
   * The way out for what the setters do not reach: a member this builder
   * cannot name, a shape it does not model, or a change that touches several
   * members at once and only makes sense as one edit.
   *
   * Prefer a setter where one exists. A setter says which member it writes,
   * and a reader of the call site learns that without following a closure.
   */
  mutate(fn: (value: UUID) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder that can be configured without affecting this
   * one.
   *
   * One level deep, and the boundary matters. The array, byte-array and keyed
   * members are copied, so writing through one clone is not visible through
   * another. What is shared is anything held inside those, and every other
   * member that refers to an object. Use mutate to replace such a member
   * outright rather than editing it in place.
   */
  clone(): UUIDBuilder {
    const out = new UUIDBuilder(Object.assign(new UUID(), this.v));
    return out;
  }

  /**
   * build returns the constructed UUID.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): UUID {
    return Object.assign(new UUID(), this.v);
  }
}

/**
 * UserUrlBuilder constructs UserUrl values.
 *
 * Start with UserUrlBuilder.create or UserUrlBuilder.from, set what you
 * care about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to UserUrl breaks every
 * literal at once and a reader cannot tell which members a given call site
 * actually cares about. Here the constructor supplies the rest, and each
 * call states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class UserUrlBuilder {
  private readonly v: UserUrl;

  private constructor(v: UserUrl) {
    this.v = v;
  }

  /**
   * create returns a builder seeded from the module's companion, with any
   * declared per-member default applied over it.
   *
   * Two layers, in that order: whatever the companion supplies, then each
   * member that declares its own default. A member the companion sets and
   * the declaration does not keeps the companion's value, so the two sources
   * compose rather than one replacing the other.
   */
  static create(): UserUrlBuilder {
    const v = userUrlDefaults();
    return new UserUrlBuilder(v);
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from(value: UserUrl): UserUrlBuilder {
    return new UserUrlBuilder(Object.assign(new UserUrl(), value));
  }

  /**
   * withHref replaces href.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withHref(value: string): this {
    this.v.href = value;
    return this;
  }

  /**
   * withUserId replaces userId.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withUserId(value: string): this {
    this.v.userId = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction.
   *
   * The way out for what the setters do not reach: a member this builder
   * cannot name, a shape it does not model, or a change that touches several
   * members at once and only makes sense as one edit.
   *
   * Prefer a setter where one exists. A setter says which member it writes,
   * and a reader of the call site learns that without following a closure.
   */
  mutate(fn: (value: UserUrl) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder that can be configured without affecting this
   * one.
   *
   * One level deep, and the boundary matters. The array, byte-array and keyed
   * members are copied, so writing through one clone is not visible through
   * another. What is shared is anything held inside those, and every other
   * member that refers to an object. Use mutate to replace such a member
   * outright rather than editing it in place.
   */
  clone(): UserUrlBuilder {
    const out = new UserUrlBuilder(Object.assign(new UserUrl(), this.v));
    return out;
  }

  /**
   * build returns the constructed UserUrl.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): UserUrl {
    return Object.assign(new UserUrl(), this.v);
  }
}

/**
 * TierBuilder constructs Tier values.
 *
 * Start with TierBuilder.create or TierBuilder.from, set what you care
 * about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to Tier breaks every literal
 * at once and a reader cannot tell which members a given call site
 * actually cares about. Here the constructor supplies the rest, and each
 * call states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class TierBuilder {
  private readonly v: Tier;

  private constructor(v: Tier) {
    this.v = v;
  }

  /**
   * create returns a builder holding a new Tier.
   *
   * No member declares a default, so there is nothing to seed. A caller that
   * builds immediately gets the same value that a new Tier has,
   * and the setters are what make it differ.
   */
  static create(): TierBuilder {
    return new TierBuilder(new Tier());
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from(value: Tier): TierBuilder {
    return new TierBuilder(Object.assign(new Tier(), value));
  }

  /**
   * withLevel replaces level.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withLevel(value: number): this {
    this.v.level = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction.
   *
   * The way out for what the setters do not reach: a member this builder
   * cannot name, a shape it does not model, or a change that touches several
   * members at once and only makes sense as one edit.
   *
   * Prefer a setter where one exists. A setter says which member it writes,
   * and a reader of the call site learns that without following a closure.
   */
  mutate(fn: (value: Tier) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder that can be configured without affecting this
   * one.
   *
   * One level deep, and the boundary matters. The array, byte-array and keyed
   * members are copied, so writing through one clone is not visible through
   * another. What is shared is anything held inside those, and every other
   * member that refers to an object. Use mutate to replace such a member
   * outright rather than editing it in place.
   */
  clone(): TierBuilder {
    const out = new TierBuilder(Object.assign(new Tier(), this.v));
    return out;
  }

  /**
   * build returns the constructed Tier.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): Tier {
    return Object.assign(new Tier(), this.v);
  }
}

/**
 * PageBuilder constructs Page values.
 *
 * Start with PageBuilder.create or PageBuilder.from, set what you care
 * about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to Page breaks every literal
 * at once and a reader cannot tell which members a given call site
 * actually cares about. Here the constructor supplies the rest, and each
 * call states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class PageBuilder<T> {
  private readonly v: Page<T>;

  private constructor(v: Page<T>) {
    this.v = v;
  }

  /**
   * create returns a builder holding a new Page.
   *
   * No member declares a default, so there is nothing to seed. A caller that
   * builds immediately gets the same value that a new Page has,
   * and the setters are what make it differ.
   */
  static create<T>(): PageBuilder<T> {
    return new PageBuilder<T>(new Page<T>());
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from<T>(value: Page<T>): PageBuilder<T> {
    return new PageBuilder<T>(Object.assign(new Page<T>(), value));
  }

  /**
   * withItems replaces items with the elements given.
   *
   * Variadic, and replacing: whatever the collection held is discarded,
   * including anything a constructor seeded. appendItems is the counterpart
   * that keeps it. Called with no arguments this empties the member rather
   * than leaving it alone, which is how a caller clears a seeded default.
   */
  withItems(...values: T[]): this {
    this.v.items = values;
    return this;
  }

  /**
   * appendItems adds to items, keeping what is already there.
   *
   * The counterpart to withItems, which discards it. Both are offered
   * because a caller extending a seeded default and a caller replacing it
   * want opposite things, and one method cannot mean both.
   */
  appendItems(...values: T[]): this {
    this.v.items = [...this.v.items, ...values];
    return this;
  }

  /**
   * mutate applies fn to the value under construction.
   *
   * The way out for what the setters do not reach: a member this builder
   * cannot name, a shape it does not model, or a change that touches several
   * members at once and only makes sense as one edit.
   *
   * Prefer a setter where one exists. A setter says which member it writes,
   * and a reader of the call site learns that without following a closure.
   */
  mutate(fn: (value: Page<T>) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder that can be configured without affecting this
   * one.
   *
   * One level deep, and the boundary matters. The array, byte-array and keyed
   * members are copied, so writing through one clone is not visible through
   * another. What is shared is anything held inside those, and every other
   * member that refers to an object. Use mutate to replace such a member
   * outright rather than editing it in place.
   */
  clone(): PageBuilder<T> {
    const out = new PageBuilder<T>(Object.assign(new Page<T>(), this.v));
    out.v.items = [...out.v.items];
    return out;
  }

  /**
   * build returns the constructed Page.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): Page<T> {
    return Object.assign(new Page<T>(), this.v);
  }
}

/**
 * CounterBuilder constructs Counter values.
 *
 * Start with CounterBuilder.create or CounterBuilder.from, set what you
 * care about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to Counter breaks every
 * literal at once and a reader cannot tell which members a given call site
 * actually cares about. Here the constructor supplies the rest, and each
 * call states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class CounterBuilder {
  private readonly v: Counter;

  private constructor(v: Counter) {
    this.v = v;
  }

  /**
   * create returns a builder holding a new Counter.
   *
   * No member declares a default, so there is nothing to seed. A caller that
   * builds immediately gets the same value that a new Counter has,
   * and the setters are what make it differ.
   */
  static create(): CounterBuilder {
    return new CounterBuilder(new Counter());
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from(value: Counter): CounterBuilder {
    return new CounterBuilder(Object.assign(new Counter(), value));
  }

  /**
   * withCount replaces count.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withCount(value: Counter["count"]): this {
    this.v.count = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction.
   *
   * The way out for what the setters do not reach: a member this builder
   * cannot name, a shape it does not model, or a change that touches several
   * members at once and only makes sense as one edit.
   *
   * Prefer a setter where one exists. A setter says which member it writes,
   * and a reader of the call site learns that without following a closure.
   */
  mutate(fn: (value: Counter) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder that can be configured without affecting this
   * one.
   *
   * One level deep, and the boundary matters. The array, byte-array and keyed
   * members are copied, so writing through one clone is not visible through
   * another. What is shared is anything held inside those, and every other
   * member that refers to an object. Use mutate to replace such a member
   * outright rather than editing it in place.
   */
  clone(): CounterBuilder {
    const out = new CounterBuilder(Object.assign(new Counter(), this.v));
    return out;
  }

  /**
   * build returns the constructed Counter.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): Counter {
    return Object.assign(new Counter(), this.v);
  }
}

/**
 * HeaderBuilder constructs Header values.
 *
 * Start with HeaderBuilder.create or HeaderBuilder.from, set what you care
 * about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to Header breaks every literal
 * at once and a reader cannot tell which members a given call site
 * actually cares about. Here the constructor supplies the rest, and each
 * call states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class HeaderBuilder {
  private readonly v: Header;

  private constructor(v: Header) {
    this.v = v;
  }

  /**
   * create returns a builder holding a new Header.
   *
   * No member declares a default, so there is nothing to seed. A caller that
   * builds immediately gets the same value that a new Header has,
   * and the setters are what make it differ.
   */
  static create(): HeaderBuilder {
    return new HeaderBuilder(new Header());
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from(value: Header): HeaderBuilder {
    return new HeaderBuilder(Object.assign(new Header(), value));
  }

  /**
   * withV2 replaces v2.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withV2(value: string): this {
    this.v.v2 = value;
    return this;
  }

  /**
   * with$ref replaces $ref.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  with$ref(value: string): this {
    this.v.$ref = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction.
   *
   * The way out for what the setters do not reach: a member this builder
   * cannot name, a shape it does not model, or a change that touches several
   * members at once and only makes sense as one edit.
   *
   * Prefer a setter where one exists. A setter says which member it writes,
   * and a reader of the call site learns that without following a closure.
   */
  mutate(fn: (value: Header) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder that can be configured without affecting this
   * one.
   *
   * One level deep, and the boundary matters. The array, byte-array and keyed
   * members are copied, so writing through one clone is not visible through
   * another. What is shared is anything held inside those, and every other
   * member that refers to an object. Use mutate to replace such a member
   * outright rather than editing it in place.
   */
  clone(): HeaderBuilder {
    const out = new HeaderBuilder(Object.assign(new Header(), this.v));
    return out;
  }

  /**
   * build returns the constructed Header.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): Header {
    return Object.assign(new Header(), this.v);
  }
}

/**
 * BinBuilder constructs Bin values.
 *
 * Start with BinBuilder.create or BinBuilder.from, set what you care
 * about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to Bin breaks every literal at
 * once and a reader cannot tell which members a given call site actually
 * cares about. Here the constructor supplies the rest, and each call
 * states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class BinBuilder {
  private readonly v: Bin;

  private constructor(v: Bin) {
    this.v = v;
  }

  /**
   * create returns a builder holding a new Bin.
   *
   * No member declares a default, so there is nothing to seed. A caller that
   * builds immediately gets the same value that a new Bin has,
   * and the setters are what make it differ.
   */
  static create(): BinBuilder {
    return new BinBuilder(new Bin());
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from(value: Bin): BinBuilder {
    return new BinBuilder(Object.assign(new Bin(), value));
  }

  /**
   * withCodes replaces codes with the elements given.
   *
   * Variadic, and replacing: whatever the collection held is discarded,
   * including anything a constructor seeded. appendCodes is the counterpart
   * that keeps it. Called with no arguments this empties the member rather
   * than leaving it alone, which is how a caller clears a seeded default.
   */
  withCodes(...values: string[]): this {
    this.v.codes = values;
    return this;
  }

  /**
   * appendCodes adds to codes, keeping what is already there.
   *
   * The counterpart to withCodes, which discards it. Both are offered
   * because a caller extending a seeded default and a caller replacing it
   * want opposite things, and one method cannot mean both.
   */
  appendCodes(...values: string[]): this {
    this.v.codes = [...this.v.codes, ...values];
    return this;
  }

  /**
   * withMarks replaces marks.
   *
   * Use mutate to clear the member.
   */
  withMarks(value: Set<string>): this {
    this.v.marks = value;
    return this;
  }

  /**
   * withNotes replaces notes.
   *
   * Use mutate to clear the member.
   */
  withNotes(value: string[]): this {
    this.v.notes = value;
    return this;
  }

  /**
   * withAliases replaces aliases with the elements given.
   *
   * Variadic, and replacing: whatever the collection held is discarded,
   * including anything a constructor seeded. appendAliases is the counterpart
   * that keeps it. Called with no arguments this empties the member rather
   * than leaving it alone, which is how a caller clears a seeded default.
   */
  withAliases(...values: string[]): this {
    this.v.aliases = values;
    return this;
  }

  /**
   * appendAliases adds to aliases, keeping what is already there.
   *
   * The counterpart to withAliases, which discards it. Both are offered
   * because a caller extending a seeded default and a caller replacing it
   * want opposite things, and one method cannot mean both.
   */
  appendAliases(...values: string[]): this {
    this.v.aliases = [...this.v.aliases, ...values];
    return this;
  }

  /**
   * withRaw replaces raw with the bytes given.
   *
   * The bytes are kept, not copied: the builder and every value built from
   * it share the array the caller passed, so a caller that mutates it
   * afterwards changes what was already built.
   */
  withRaw(value: Uint8Array): this {
    this.v.raw = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction.
   *
   * The way out for what the setters do not reach: a member this builder
   * cannot name, a shape it does not model, or a change that touches several
   * members at once and only makes sense as one edit.
   *
   * Prefer a setter where one exists. A setter says which member it writes,
   * and a reader of the call site learns that without following a closure.
   */
  mutate(fn: (value: Bin) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder that can be configured without affecting this
   * one.
   *
   * One level deep, and the boundary matters. The array, byte-array and keyed
   * members are copied, so writing through one clone is not visible through
   * another. What is shared is anything held inside those, and every other
   * member that refers to an object. Use mutate to replace such a member
   * outright rather than editing it in place.
   */
  clone(): BinBuilder {
    const out = new BinBuilder(Object.assign(new Bin(), this.v));
    out.v.aliases = [...out.v.aliases];
    out.v.raw = out.v.raw.slice();
    return out;
  }

  /**
   * build returns the constructed Bin.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): Bin {
    return Object.assign(new Bin(), this.v);
  }
}

/**
 * LooseBuilder constructs Loose values.
 *
 * Start with LooseBuilder.create or LooseBuilder.from, set what you care
 * about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to Loose breaks every literal
 * at once and a reader cannot tell which members a given call site
 * actually cares about. Here the constructor supplies the rest, and each
 * call states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class LooseBuilder {
  private readonly v: Loose;

  private constructor(v: Loose) {
    this.v = v;
  }

  /**
   * create returns a builder holding a new Loose.
   *
   * No member declares a default, so there is nothing to seed. A caller that
   * builds immediately gets the same value that a new Loose has,
   * and the setters are what make it differ.
   */
  static create(): LooseBuilder {
    return new LooseBuilder(new Loose());
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from(value: Loose): LooseBuilder {
    return new LooseBuilder(Object.assign(new Loose(), value));
  }

  /**
   * withItems replaces items with the elements given.
   *
   * Variadic, and replacing: whatever the collection held is discarded,
   * including anything a constructor seeded. appendItems is the counterpart
   * that keeps it. Called with no arguments this empties the member rather
   * than leaving it alone, which is how a caller clears a seeded default.
   */
  withItems(...values: Notifier[]): this {
    this.v.items = values;
    return this;
  }

  /**
   * appendItems adds to items, keeping what is already there.
   *
   * The counterpart to withItems, which discards it. Both are offered
   * because a caller extending a seeded default and a caller replacing it
   * want opposite things, and one method cannot mean both.
   */
  appendItems(...values: Notifier[]): this {
    this.v.items = [...this.v.items, ...values];
    return this;
  }

  /**
   * withIndex replaces index wholesale.
   *
   * The whole mapping, not one entry: what the member held is discarded
   * rather than merged into, and the map is kept rather than copied — so a
   * caller that writes to it afterwards changes what was already built.
   * withIndexEntry is what sets a single key.
   */
  withIndex(value: Map<string, Notifier>): this {
    this.v.index = value;
    return this;
  }

  /**
   * withIndexEntry sets one key of index.
   *
   * A key already present takes the new value.
   */
  withIndexEntry(key: string, value: Notifier): this {
    this.v.index = new Map([...this.v.index, [key, value]]);
    return this;
  }

  /**
   * withIndexEntries merges entries into index.
   *
   * Merges rather than replaces: a key already present takes the new value,
   * and keys absent from entries are left alone. Use withIndex to discard
   * what is there.
   */
  withIndexEntries(entries: Map<string, Notifier>): this {
    this.v.index = new Map([...this.v.index, ...entries]);
    return this;
  }

  /**
   * withKeys replaces keys wholesale.
   *
   * The whole collection, not one key: what the member held is discarded
   * rather than merged into, and the set is kept rather than copied.
   * withKeysEntry is what adds to what is already there.
   */
  withKeys(value: Set<Notifier>): this {
    this.v.keys = value;
    return this;
  }

  /**
   * withKeysEntry adds one key to keys.
   *
   * A set entry is a key, so there is nothing else for a caller to supply.
   */
  withKeysEntry(key: Notifier): this {
    this.v.keys = new Set([...this.v.keys, key]);
    return this;
  }

  /**
   * withKeysEntries adds each key to keys, keeping those it holds.
   *
   * Additive like withKeysEntry, so two calls union rather than the second
   * replacing the first. A key already present is left as it stands.
   */
  withKeysEntries(...keys: Notifier[]): this {
    this.v.keys = new Set([...this.v.keys, ...keys]);
    return this;
  }

  /**
   * withExtra replaces extra.
   *
   * Use mutate to clear the member.
   */
  withExtra(value: Notifier): this {
    this.v.extra = value;
    return this;
  }

  /**
   * withValue replaces value.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withValue(value: Notifier): this {
    this.v.value = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction.
   *
   * The way out for what the setters do not reach: a member this builder
   * cannot name, a shape it does not model, or a change that touches several
   * members at once and only makes sense as one edit.
   *
   * Prefer a setter where one exists. A setter says which member it writes,
   * and a reader of the call site learns that without following a closure.
   */
  mutate(fn: (value: Loose) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder that can be configured without affecting this
   * one.
   *
   * One level deep, and the boundary matters. The array, byte-array and keyed
   * members are copied, so writing through one clone is not visible through
   * another. What is shared is anything held inside those, and every other
   * member that refers to an object. Use mutate to replace such a member
   * outright rather than editing it in place.
   */
  clone(): LooseBuilder {
    const out = new LooseBuilder(Object.assign(new Loose(), this.v));
    out.v.items = [...out.v.items];
    out.v.index = new Map(out.v.index);
    out.v.keys = new Set(out.v.keys);
    return out;
  }

  /**
   * build returns the constructed Loose.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): Loose {
    return Object.assign(new Loose(), this.v);
  }
}
