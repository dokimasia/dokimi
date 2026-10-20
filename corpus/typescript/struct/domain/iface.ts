// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the struct-kind fixture: the input the builder generator reads.
 *
 * @remarks
 * A builder is derived per field, and the emitted setter depends on the
 * field's type rather than on its name — an array owes an appender as well as
 * a setter, a map owes withEntry and withEntries, a set owes an entry setter
 * without a value, and a private field owes nothing at all. So the fixture's
 * job is to carry one field of every shape that changes the answer, not to
 * look like a realistic domain type.
 *
 * {@link userDefaults} and {@link addressDefaults} are the companions that
 * the generated builders seed from. The module declares several classes on
 * purpose: a companion named defaults would collide on the second, which is
 * why the convention carries the name of the class.
 *
 * The classes are split by what varies rather than collected into one class:
 * {@link Primitives} covers the scalar spread, {@link Containers} covers the
 * shapes needing more than one setter, and {@link User} covers composition —
 * inheritance, nesting, and self-reference.
 *
 * Routing is declared once for the module rather than repeated on each
 * class: every builder here lands in the same companion directory, so a
 * per-class directive is the same statement written N times, and the Nth
 * copy is the one that gets forgotten.
 *
 * @packageDocumentation
 */
//dokimi:out path=domaintest/

/**
 * Represents a day of the week as a number.
 *
 * @remarks
 * A TypeScript alias is transparent, so the setter takes the alias spelling
 * and accepts any number.
 */
export type Weekday = number;

/** Represents a temperature in degrees Celsius. */
export type Celsius = number;

/** Represents an identifier. */
export type Id = string;

/**
 * Represents bytes through an alias.
 *
 * @remarks
 * An alias is its target, so it must take the same setter as Uint8Array and
 * not a distinct one.
 */
export type Bytes = Uint8Array;

/** Represents a source of bytes, which a builder can set and never construct. */
export interface Reader {
  /** Returns the next bytes of the source. */
  read(): Promise<Uint8Array>;
}

/**
 * Represents a record with one field of every primitive type.
 *
 * @remarks
 * The width of a number does not change the setter's shape, so number and
 * bigint are the whole spread, beside the aliases over them.
 */
//dokimi:builder
export class Primitives {
  /** Is true when the flag is set. */
  flag: boolean = false;
  /** Contains a text. */
  text: string = "";
  /** Contains a number. */
  count: number = 0;
  /** Contains an integer beyond the safe range of a number. */
  large: bigint = 0n;
  /** Contains a day through an alias over a number. */
  day: Weekday = 0;
  /** Contains a temperature through an alias over a number. */
  temp: Celsius = 0;
  /** Contains a reference through an alias over a string. */
  ref: Id = "";
}

/** Represents a record with every shape that owes more than a single setter. */
//dokimi:builder
export class Containers {
  /** Lists texts, which owe a setter and an appender. */
  slice: string[] = [];

  /** Lists addresses, whose appender takes an {@link Address}. */
  sliceOfStruct: Address[] = [];

  /** Contains a tuple of a fixed length, which takes no appender. */
  array: [number, number, number] = [0, 0, 0];

  /** Contains bytes. */
  body: Uint8Array = new Uint8Array(0);

  /** Contains bytes through an alias, which takes the same setter. */
  aliased: Bytes = new Uint8Array(0);

  /** Maps texts to texts, which owe withEntry for one pair and withEntries for many. */
  map: Map<string, string> = new Map();

  /** Maps texts to addresses, so the second parameter of withEntry is no scalar. */
  mapOfStruct: Map<string, Address> = new Map();

  /**
   * Contains a set of texts.
   *
   * @remarks
   * Its entry setter takes no value parameter at all: a set has no value to
   * vary.
   */
  set: Set<string> = new Set();

  /**
   * Contains an optional number.
   *
   * @remarks
   * The field distinguishes unset from zero, so its setter takes a value.
   */
  pointer?: number;

  /** Contains a reader, which can only be set, never constructed by the builder. */
  reader?: Reader;

  /**
   * Contains a callback.
   *
   * @remarks
   * A function is settable but not comparable, which breaks any builder that
   * derives equality from the field set.
   */
  callback: () => Promise<void> = () => Promise.resolve();

  /**
   * Contains a stream of events.
   *
   * @remarks
   * A builder must not construct a stream, since its source is a decision the
   * caller owns.
   */
  events?: AsyncIterable<string>;

  /** Contains an arbitrary value, which must not be special-cased into another setter. */
  extra: unknown = undefined;

  /**
   * Contains an error.
   *
   * @remarks
   * The type has a well-known name, present because a builder that
   * special-cases the identifier rather than the type gets it wrong.
   */
  err?: Error;
}

/**
 * Represents an address.
 *
 * @remarks
 * An address is nested inside {@link User} and used as an array and a map
 * element in {@link Containers}, so its own builder has to be reachable from
 * three directions.
 */
//dokimi:builder
export class Address {
  /** Contains the street and the house number. */
  street: string = "";
  /** Contains the city. */
  city: string = "";
  /** Contains the postal code. */
  postal: string = "";
}

/**
 * Represents a role, which {@link User} extends.
 *
 * @remarks
 * The fields of a role are inherited. A builder reading only declared fields
 * misses name and level entirely.
 */
//dokimi:builder
export class Role {
  /** Contains the name of the role. */
  name: string = "";
  /** Contains the level of the role. */
  level: number = 0;
}

/**
 * Represents an audit, which {@link User} nests as an optional field.
 *
 * @remarks
 * The fields of an audit are only reachable once the audit exists, so a
 * builder has to allocate it before it sets a field through it.
 */
//dokimi:builder
export class Audit {
  /** Contains the instant of the creation. */
  createdAt: Date = new Date(0);
  /** Contains the creator. */
  createdBy: string = "";
}

/**
 * Represents a user.
 *
 * @remarks
 * The class covers composition: inheritance of a base, nesting, self
 * reference, and the private field a builder must skip.
 */
//dokimi:builder
export class User extends Role {
  /** Contains the audit of the user, where it has one. */
  audit?: Audit;

  /** Contains the username. */
  //dokimi:default "anonymous"
  username: string = "";

  /** Contains the age of the user. */
  age: number = 0;

  /** Is true when the user is active. */
  active: boolean = false;

  /** Contains the home address, which a builder sets as a whole. */
  home: Address = new Address();

  /**
   * Contains the manager of the user, where it has one.
   *
   * @remarks
   * The field is self-referential. A builder that recurses into nested class
   * types without a cycle check never terminates here.
   */
  manager?: User;

  /** Lists the reports of the user, the same cycle by another route. */
  reports: User[] = [];

  /**
   * Contains the deadline.
   *
   * @remarks
   * A Date is a library class that the builder must treat as opaque rather
   * than walk into.
   */
  deadline: Date = new Date(0);

  /** Contains the timeout of the user, in milliseconds. */
  timeout: number = 0;

  /**
   * Contains a note that is not settable from another module, so a generated
   * builder must skip it.
   */
  #note: string = "";

  /** Returns the private note. */
  note(): string {
    return this.#note;
  }
}

/**
 * Returns the seed of the builder of {@link User}.
 *
 * @remarks
 * The values are deliberately different from the values of a new user across
 * several field shapes. A seed that happened to equal a new user would make
 * the companion path indistinguishable from having no companion at all, which
 * is the branch the plain fixture covers.
 */
export function userDefaults(): User {
  const user = new User();
  user.name = "member";
  user.level = 1;
  user.username = "anonymous";
  user.age = 30;
  user.active = true;
  user.home = new Address();
  user.home.city = "Amsterdam";
  user.deadline = new Date(0);
  user.timeout = 5000;
  return user;
}

/**
 * Returns the seed of the builder of {@link Address}.
 *
 * @remarks
 * Its presence is the point: two companions in one module is what the
 * convention of the name exists for, and a fixture with one cannot
 * demonstrate it.
 */
export function addressDefaults(): Address {
  const address = new Address();
  address.street = "Keizersgracht 1";
  address.city = "Amsterdam";
  address.postal = "1015";
  return address;
}
