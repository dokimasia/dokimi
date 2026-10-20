// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the generic half of the struct-kind corpus.
 *
 * @remarks
 * Type parameters have to survive end to end: the builder type carries them,
 * its constructor carries them, and every setter that touches a parameterised
 * field takes the parameter rather than its constraint. A generator that
 * erases them produces code which compiles only when T happens to be the
 * constraint itself, so the failure shows up for one instantiation and not
 * another.
 *
 * Routing is declared once for the module rather than repeated on each
 * class: every builder here lands in the same companion directory, so a
 * per-class directive is the same statement written N times, and the Nth
 * copy is the one that gets forgotten.
 *
 * @packageDocumentation
 */
//dokimi:out path=containertest/

/** Represents a duration in seconds, which a witness names as a type of this module. */
export type Seconds = number;

/**
 * Represents a container that is generic over one unconstrained parameter.
 *
 * @remarks
 * An unconstrained parameter is the case where nothing about T can be
 * assumed.
 */
//dokimi:builder
export class Container<T> {
  /** Contains the value of the container. */
  value!: T;

  /** Lists the items of the parameter, so its appender takes T rather than unknown. */
  items: T[] = [];

  /**
   * Contains the label of the container.
   *
   * @remarks
   * The label is not parameterised, so its setter must not acquire a type
   * parameter it does not need.
   */
  label: string = "";
}

/**
 * Represents a pair that is generic over two parameters.
 *
 * @remarks
 * A generator that assumes a single parameter gets this wrong.
 */
//dokimi:builder
export class Pair<K, V> {
  /** Contains the key of the pair. */
  key!: K;

  /** Contains the value of the pair. */
  value!: V;

  /**
   * Maps keys to values.
   *
   * @remarks
   * The map is keyed by the first parameter, so its setters have to thread
   * both parameters through.
   */
  index: Map<K, V> = new Map();
}

/**
 * Represents a constraint that the derivation of a witness cannot answer for.
 *
 * @remarks
 * Deriving a witness only ever tries a builtin against an unconstrained
 * parameter. So every class below that bounds a parameter by Keyable derives
 * no witness of its own, and needs //dokimi:witness to generate anything
 * beyond the builder itself: the builder is written in terms of the parameter
 * and needs no concrete type, but a check has to construct an actual value,
 * and derivation is the only thing that could have said which type to
 * construct one at.
 */
export type Keyable = string | number;

/**
 * Represents a class without a witness, so its one parameter never derives.
 *
 * @remarks
 * The generator withholds every check for it — the negative case that
 * {@link Bucket} answers. The builder still generates in full: nothing about a
 * builder's own code needs to name a concrete type.
 */
//dokimi:builder
export class Solo<K extends Keyable> {
  /** Contains the key. */
  key!: K;
}

/**
 * Represents a pair that witnesses one of its two parameters and leaves the
 * other to derive — which it cannot.
 *
 * @remarks
 * PartialPair generates no checks either, the same as {@link Solo}.
 * Witnessing is not per-parameter best-effort: an instantiation is the whole
 * parameter list at once, and a partial one is worth exactly as much as none.
 * {@link Bucket} is PartialPair with the missing entry supplied.
 */
//dokimi:witness V=Seconds
//dokimi:builder
export class PartialPair<K extends Keyable, V> {
  /** Contains the key. */
  key!: K;
  /** Contains the value. */
  value!: V;
}

/**
 * Represents a bucket that witnesses both of its parameters, so it generates
 * real checks where {@link Solo} and {@link PartialPair} generate none.
 *
 * @remarks
 * K's witness is load bearing — Keyable derives nothing on its own — while
 * V's overrides a witness derivation would have supplied for an ordinary
 * unconstrained parameter, which is worth doing anyway to prove an authored
 * answer always wins. The pair also covers both forms a witness takes: K
 * names a builtin, and V names a type of this module.
 */
//dokimi:witness K=string V=Seconds
//dokimi:builder
export class Bucket<K extends Keyable, V> {
  /** Contains the key. */
  key!: K;
  /** Contains the value. */
  value!: V;
  /** Contains the lifetime of the bucket, which is not parameterised. */
  ttl: Seconds = 0;
}

/**
 * Represents a triple that carries the directive through three parameters.
 *
 * @remarks
 * A's witness is load bearing for the reason {@link Bucket}'s K is, and B and
 * C are overrides for the reason Bucket's V is — proving the directive's
 * per-parameter parsing is not tuned to exactly two entries.
 */
//dokimi:witness A=string B=number C=boolean
//dokimi:builder
export class Triple<A extends Keyable, B, C> {
  /** Contains the first element. */
  first!: A;
  /** Contains the second element. */
  second!: B;
  /** Contains the third element. */
  third!: C;
}

/**
 * Represents {@link Triple} with a fourth parameter that {@link Keyable} also
 * bounds.
 *
 * @remarks
 * Two of its four witnesses are load bearing rather than one — proving the
 * deriver's all-or-nothing rule is checked at every position it applies to,
 * not only the first.
 */
//dokimi:witness A=string B=number C=boolean D=string
//dokimi:builder
export class Quad<A extends Keyable, B, C, D extends Keyable> {
  /** Contains the first element. */
  first!: A;
  /** Contains the second element. */
  second!: B;
  /** Contains the third element. */
  third!: C;
  /** Contains the fourth element. */
  fourth!: D;
}
