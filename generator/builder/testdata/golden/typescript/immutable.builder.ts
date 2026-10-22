import type { Set } from "immutable";
import { Frozen } from "./immutable";

/**
 * FrozenBuilder constructs Frozen values. FrozenBuilder.create and
 * FrozenBuilder.from return a new builder. Each setter sets one member and
 * returns the builder, so the calls chain. build returns the value.
 */
export class FrozenBuilder {
  private readonly v: Frozen;

  private constructor(v: Frozen) {
    this.v = v;
  }

  /** create returns a builder of a new Frozen. */
  static create(): FrozenBuilder {
    return new FrozenBuilder(new Frozen());
  }

  /** from returns a builder of a copy of value. */
  static from(value: Frozen): FrozenBuilder {
    return new FrozenBuilder(Object.assign(new Frozen(), value));
  }

  /** withTags replaces tags with value. */
  withTags(value: Set<string>): this {
    this.v.tags = value;
    return this;
  }

  /**
   * mutate applies fn to the value under construction. It changes what no
   * setter changes, such as several members in one step.
   */
  mutate(fn: (value: Frozen) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): FrozenBuilder {
    const out = new FrozenBuilder(Object.assign(new Frozen(), this.v));
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): Frozen {
    return Object.assign(new Frozen(), this.v);
  }
}
