import type { Mirror } from './aliased';
import { Account } from './catalog';

/**
 * MirrorBuilder constructs Mirror values.
 *
 * Start with MirrorBuilder.create or MirrorBuilder.from, set what you care
 * about, and call build.
 *
 * The point is what you do not have to write. An object literal restates
 * every member each time, so a member added to Mirror breaks every literal
 * at once and a reader cannot tell which members a given call site
 * actually cares about. Here the constructor supplies the rest, and each
 * call states only what it varies.
 *
 * Every setter returns the builder, so calls chain. clone is how one
 * builder becomes two.
 */
export class MirrorBuilder {
  private readonly v: Mirror;

  private constructor(v: Mirror) {
    this.v = v;
  }

  /**
   * create returns a builder seeded with the declared defaults.
   *
   * Only the members that declare one. Everything else holds the value that
   * a new Mirror has, so a caller that sets nothing gets a mixture
   * rather than either extreme — which is what makes a check against a built
   * value worth writing.
   */
  static create(): MirrorBuilder {
    const v = new Account();
    v.owner = 'owner';
    return new MirrorBuilder(v);
  }

  /**
   * from returns a builder seeded with a copy of value.
   *
   * For the case where a caller already holds a value and wants one member
   * different — a round trip, a near-copy, a regression fixture built from
   * something a previous step produced.
   */
  static from(value: Mirror): MirrorBuilder {
    return new MirrorBuilder(Object.assign(new Account(), value));
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
  mutate(fn: (value: Mirror) => void): this {
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
  clone(): MirrorBuilder {
    const out = new MirrorBuilder(Object.assign(new Account(), this.v));
    return out;
  }

  /**
   * build returns the constructed Mirror.
   *
   * The builder stays usable afterwards, and a second call returns the value
   * as it stands then — so a caller may build, set one more member, and
   * build again. What it returns is a shallow copy of the value under
   * construction: a member holding an array, a map or a set is shared with
   * the builder and with every value built from it, and clone is what
   * separates them.
   */
  build(): Mirror {
    return Object.assign(new Account(), this.v);
  }
}
