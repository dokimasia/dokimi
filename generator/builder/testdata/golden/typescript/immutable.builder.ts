import type { Set } from 'immutable';
import { Frozen } from './immutable';

/**
 * FrozenBuilder constructs Frozen values.
 *
 * Start with FrozenBuilder.create or FrozenBuilder.from, set what you care
 * about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to Frozen breaks every literal
 * at once and a reader cannot tell which members a given call site
 * actually cares about. Here the constructor supplies the rest, and each
 * call states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class FrozenBuilder {
  private readonly v: Frozen;

  private constructor(v: Frozen) {
    this.v = v;
  }

  /**
   * create returns a builder holding a new Frozen.
   *
   * No member declares a default, so there is nothing to seed. A caller that
   * builds immediately gets the same value that a new Frozen has,
   * and the setters are what make it differ.
   */
  static create(): FrozenBuilder {
    return new FrozenBuilder(new Frozen());
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from(value: Frozen): FrozenBuilder {
    return new FrozenBuilder(Object.assign(new Frozen(), value));
  }

  /**
   * withTags replaces tags.
   *
   * The later call wins: a setter overwrites rather than merges, so a caller
   * building from a constructor that seeded this member states only what it
   * varies.
   */
  withTags(value: Set<string>): this {
    this.v.tags = value;
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
  mutate(fn: (value: Frozen) => void): this {
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
  clone(): FrozenBuilder {
    const out = new FrozenBuilder(Object.assign(new Frozen(), this.v));
    return out;
  }

  /**
   * build returns the constructed Frozen.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): Frozen {
    return Object.assign(new Frozen(), this.v);
  }
}
