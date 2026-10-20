// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for a bound on a KEYED store, beside
 * `bounded`, which bounds a collection.
 *
 * @remarks
 * The two shapes take the same stamp and want different oracles. A bounded
 * collection answers a clamped list, and every store model answers
 * everything it was fed, so the first clamped read reads as disagreement —
 * that fixture rides the twin floor and says so. A bounded cache answers one
 * key at a time, and a key it evicted is a MISS. Absence is legal there,
 * which is what makes an unbounded reference usable: the comparison is
 * one-sided, a subject hit must agree and a subject miss never disagrees,
 * and eviction stops looking like a wrong answer.
 *
 * The bound is small on purpose. Drawn sequences have to write more distinct
 * keys than the capacity or the eviction path never runs, and a law that
 * never reaches its own subject reads as coverage.
 *
 * @packageDocumentation
 */

/**
 * Identifies one entry.
 *
 * @remarks
 * The alias declares the default of every key, because no field declares
 * one.
 */
//dokimi:default "test-key"
export type Key = string;

/** Represents what an entry contains. */
export interface Value {
  /** Contains the body of the entry. */
  //dokimi:default "test-value"
  body: string;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=boundedcachetest/
export interface Mixed {
  /**
   * Stores the value under the key.
   *
   * @remarks
   * At capacity the policy picks a victim to make room; overwriting a held
   * key evicts nothing. WHICH entry dies is deliberately unstated — a victim
   * policy is the consumer's claim, not this interface's.
   */
  put(key: Key, value: Value): Promise<void>;

  /**
   * Returns the value under the key, or undefined where it is not there.
   *
   * @remarks
   * A miss is an answer rather than a rejection, which is the half of the
   * declaration that makes the bound checkable: with a rejection a caller
   * cannot tell an evicted key from a broken read.
   */
  //dokimi:mixin cacheable
  get(key: Key): Promise<Value | undefined>;

  /** Returns the count of the entries, which never exceeds the bound. */
  //dokimi:mixin bounded limit=2
  len(): Promise<number>;
}
