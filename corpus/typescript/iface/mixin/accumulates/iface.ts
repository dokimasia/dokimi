// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the accumulates mixin, which states
 * that repeating the call compounds rather than coalescing.
 *
 * @remarks
 * The second position on the effect axis, beside idempotent — not its
 * negation. A mixin appears only because somebody wrote the directive, so
 * absence already means unclaimed and there is nothing to switch off. What
 * the pair distinguishes is a state absence cannot: "nobody considered
 * whether repeating this is safe" against "it was considered, and it
 * compounds".
 *
 * The interface carries a reader for the same reason idempotent's does: a
 * claim about what two calls leave behind needs something to observe it
 * with, and a single method would have nothing to compare.
 *
 * @packageDocumentation
 */

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=accumulatestest/
export interface Mixed {
  /**
   * Adds the amount to the total under the key.
   *
   * @remarks
   * add compounds. Two calls with one argument leave twice what one call
   * left, which is the whole claim and the opposite of what idempotent's put
   * asserts about the same shape.
   */
  //dokimi:mixin accumulates observe=total
  add(key: string, amount: number): Promise<void>;

  /**
   * Returns the total that the additions under the key left.
   *
   * @remarks
   * Any N ≥ 2 proves the claim, so the mixin takes no count — the number is
   * the check's choice, not the author's assertion.
   */
  total(key: string): Promise<number>;
}
