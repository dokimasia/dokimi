// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the bounded mixin, which states that
 * the returned collection never exceeds a declared ceiling.
 *
 * @remarks
 * The interface carries whatever methods the law needs to be stateable. A
 * mixin whose law spans two calls cannot be hosted by a single method: there
 * would be nothing to compare against, and the generated subtest would pass
 * by having nothing to check.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=boundedtest/
export interface Mixed {
  /**
   * Adds an item to the collection that the ceiling bounds.
   *
   * @remarks
   * Without a writer the subject stays empty and the bound holds by vacancy —
   * a limit nothing can approach is a claim nothing checks.
   */
  add(item: string): Promise<void>;

  /**
   * Returns the items of the collection.
   *
   * @remarks
   * The directive carries the ceiling as a parameter, because the law cannot
   * be stated without a number to compare against. Five, deliberately: a
   * bound the drawn sequences cross in a handful of writes is one a broken
   * subject fails inside a single iteration.
   */
  //dokimi:mixin bounded limit=5
  list(): Promise<string[]>;
}
