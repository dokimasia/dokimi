/** Represents an address with a field of each form of a setter. */
//dokimi:builder
export class Address {
  /** Contains the origin of every address, which is static and has no setter. */
  static origin: string = "catalog";
  /** Contains the street, which has a default through the directive. */
  //dokimi:default "Dam 1"
  street: string = "";
  /** Contains the city, which has no default. */
  city: string = "";
  /** Contains the floor, which is optional. */
  floor?: number;
  /** Lists the lines of the address, which form a sequence. */
  lines: string[] = [];
  /** Contains the labels of the address, which form a mapping. */
  labels: Map<string, string> = new Map();
  /** Contains the tags of the address, which form a set. */
  tags: Set<string> = new Set();
  /** Contains the identifier, which is readonly and has no setter. */
  readonly id: string = "";
  /** Contains a private number, which has no setter. */
  private hidden: number = 0;
  /** Contains a number with a hard private name, which has no setter. */
  #secret: number = 0;

  /** Returns the sum of the private fields. */
  reveal(): number {
    return this.hidden + this.#secret;
  }
}

/** Represents the base of Derived, which declares a field that Derived inherits. */
export class Base {
  /** Contains the creator, which Derived inherits with its default. */
  //dokimi:default "system"
  createdBy: string = "";
}

/** Represents a class that extends {@link Base}. */
//dokimi:builder
export class Derived extends Base {
  /** Is true when the record is active, and Derived declares it. */
  active: boolean = false;
}

/**
 * Represents a class without fields of its own.
 *
 * @remarks
 * The builder sets the fields of {@link Derived} and of {@link Base}.
 */
//dokimi:builder
export class Heir extends Derived {}

/** Represents an account with the companion {@link accountDefaults}. */
//dokimi:builder
export class Account {
  /** Contains the owner, which has a default through the directive. */
  //dokimi:default "owner"
  owner: string = "";
  /** Contains the plan, which the companion sets. */
  plan: string = "";
}

/** Returns the seed of the builder of {@link Account}. */
export function accountDefaults(): Account {
  const account = new Account();
  account.owner = "seed";
  account.plan = "basic";
  return account;
}

/**
 * Represents a client whose name opens with an acronym.
 *
 * @remarks
 * Its companion is {@link httpClientDefaults}.
 */
//dokimi:builder
export class HTTPClient {
  /** Contains the base address, which the companion sets. */
  base: string = "";
}

/** Returns the seed of the builder of {@link HTTPClient}. */
export function httpClientDefaults(): HTTPClient {
  const client = new HTTPClient();
  client.base = "https://example.com";
  return client;
}

/**
 * Represents a hash whose name opens with an acronym and a digit.
 *
 * @remarks
 * Its companion is {@link md5HashDefaults}.
 */
//dokimi:builder
export class MD5Hash {
  /** Contains the digest, which the companion sets. */
  digest: string = "";
}

/** Returns the seed of the builder of {@link MD5Hash}. */
export function md5HashDefaults(): MD5Hash {
  const hash = new MD5Hash();
  hash.digest = "d41d8cd98f00b204e9800998ecf8427e";
  return hash;
}

/**
 * Represents an identifier whose name has capitals alone.
 *
 * @remarks
 * Its companion is {@link uuidDefaults}.
 */
//dokimi:builder
export class UUID {
  /** Contains the identifier, which the companion sets. */
  value: string = "";
}

/** Returns the seed of the builder of {@link UUID}. */
export function uuidDefaults(): UUID {
  const id = new UUID();
  id.value = "00000000-0000-0000-0000-000000000000";
  return id;
}

/**
 * Represents the address of a user, whose name has an initialism of Go as
 * its second word.
 *
 * @remarks
 * Its companion is {@link userUrlDefaults}.
 */
//dokimi:builder
export class UserUrl {
  /** Contains the address, which the companion sets. */
  href: string = "";
  /** Contains the user, whose name has an initialism of Go as its second word. */
  userId: string = "";
}

/** Returns the seed of the builder of {@link UserUrl}. */
export function userUrlDefaults(): UserUrl {
  const url = new UserUrl();
  url.href = "https://example.com/user";
  return url;
}

/**
 * Represents a tier with the function {@link tierSeed}.
 *
 * @remarks
 * tierSeed seeds the builder where the option companion-word is seed.
 */
//dokimi:builder
export class Tier {
  /** Contains the level, which {@link tierSeed} sets. */
  level: number = 0;
}

/** Returns the seed of the builder of {@link Tier} where the option companion-word is seed. */
export function tierSeed(): Tier {
  const tier = new Tier();
  tier.level = 2;
  return tier;
}

/** Represents a page with a type parameter. */
//dokimi:builder
export class Page<T> {
  /** Lists the items, whose type is the type parameter. */
  items: T[] = [];
}

/**
 * Represents a counter with a property without a type annotation.
 *
 * @remarks
 * The function {@link counterDefaults} has the name of its companion and no
 * return type.
 */
//dokimi:builder
export class Counter {
  /** Contains the count, whose type its initializer gives it. */
  count = 0;
}

/**
 * Returns a new {@link Counter}.
 *
 * @remarks
 * The function has no return type, so it is no companion.
 */
export function counterDefaults() {
  return new Counter();
}

/**
 * Represents a header with properties whose names a property access can write
 * after a dot, and properties whose names it cannot.
 */
//dokimi:builder
export class Header {
  /** Contains the content type, whose name is no identifier and has no setter. */
  "content-type": string = "";
  /** Contains a property whose name is no identifier and has no setter. */
  0: string = "";
  /** Contains a property whose name is an identifier with a digit. */
  v2: string = "";
  /** Contains a property whose name opens with a dollar sign. */
  $ref: string = "";
  /** Contains the tag of the header, whose name is computed and has no setter. */
  [Symbol.toStringTag]: string = "Header";
}

/** Represents lines as a transparent alias of a sequence. */
export type Lines = string[];

/** Represents a record with the forms that {@link Address} does not have. */
//dokimi:builder
export class Bin {
  /** Lists codes in a sequence that no holder can change. */
  codes: ReadonlyArray<string> = [];
  /** Contains an optional set of marks. */
  marks?: Set<string>;
  /** Lists optional notes. */
  notes?: string[];
  /** Lists aliases in a sequence through an alias. */
  aliases: Lines = [];
  /** Contains bytes, and TypeScript has no setter of bytes from text. */
  raw: Uint8Array = new Uint8Array(0);
}

/** Represents an interface with a method, of which eidos derives no sample. */
export interface Notifier {
  /** Sends a notification. */
  notify(): void;
}

/** Represents a record with a property of each form whose type has no sample. */
//dokimi:builder
export class Loose {
  /** Lists notifiers. */
  items: Notifier[] = [];
  /** Maps names to notifiers. */
  index: Map<string, Notifier> = new Map();
  /** Contains a set of notifiers. */
  keys: Set<Notifier> = new Set();
  /** Contains an optional notifier. */
  extra?: Notifier;
  /** Contains a notifier. */
  value: Notifier = { notify: () => undefined };
}
