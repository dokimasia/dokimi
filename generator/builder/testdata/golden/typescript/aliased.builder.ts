import type { Mirror } from "./aliased";
import { Account } from "./catalog";

/**
 * MirrorBuilder constructs Mirror values. MirrorBuilder.create and
 * MirrorBuilder.from return a new builder. Each setter sets one member and
 * returns the builder, so the calls chain. build returns the value.
 */
export class MirrorBuilder {
  private readonly v: Mirror;

  private constructor(v: Mirror) {
    this.v = v;
  }

  /**
   * create returns a builder of the declared defaults. A member without a
   * declared default has the value of a new Mirror.
   */
  static create(): MirrorBuilder {
    const v = new Account();
    v.owner = "owner";
    return new MirrorBuilder(v);
  }

  /** from returns a builder of a copy of value. */
  static from(value: Mirror): MirrorBuilder {
    return new MirrorBuilder(Object.assign(new Account(), value));
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
  mutate(fn: (value: Mirror) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): MirrorBuilder {
    const out = new MirrorBuilder(Object.assign(new Account(), this.v));
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): Mirror {
    return Object.assign(new Account(), this.v);
  }
}
