import { type Set, Shadowed } from "./shadow";

/**
 * ShadowedBuilder constructs Shadowed values. ShadowedBuilder.create and
 * ShadowedBuilder.from return a new builder. Each setter sets one member
 * and returns the builder, so the calls chain. build returns the value.
 */
export class ShadowedBuilder {
  private readonly v: Shadowed;

  private constructor(v: Shadowed) {
    this.v = v;
  }

  /** create returns a builder of a new Shadowed. */
  static create(): ShadowedBuilder {
    return new ShadowedBuilder(new Shadowed());
  }

  /** from returns a builder of a copy of value. */
  static from(value: Shadowed): ShadowedBuilder {
    return new ShadowedBuilder(Object.assign(new Shadowed(), value));
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
  mutate(fn: (value: Shadowed) => void): this {
    fn(this.v);
    return this;
  }

  /**
   * clone returns a builder of a copy of the value under construction. The
   * copy has its own arrays, maps, sets and byte arrays, so a change of one
   * builder leaves the other unchanged. The two builders share the elements of
   * those members and every other member that refers to an object.
   */
  clone(): ShadowedBuilder {
    const out = new ShadowedBuilder(Object.assign(new Shadowed(), this.v));
    return out;
  }

  /**
   * build returns a shallow copy of the value under construction. The builder
   * stays usable, and a later call returns the value with the changes since.
   * The copy shares its arrays, maps and sets with the builder, and clone
   * returns a builder that shares none of them.
   */
  build(): Shadowed {
    return Object.assign(new Shadowed(), this.v);
  }
}
