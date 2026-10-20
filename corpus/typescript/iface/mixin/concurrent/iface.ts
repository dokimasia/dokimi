// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the concurrent mixin, which states that
 * the method is safe to call while other calls of it are in flight.
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
//dokimi:out path=concurrenttest/
export interface Mixed {
  /**
   * Increments the counter under the key.
   *
   * @remarks
   * The generated stress drives bump concurrently. The law is the absence of
   * a race, so the method must mutate something.
   */
  //dokimi:mixin concurrent
  bump(key: string): Promise<void>;

  /**
   * Returns the count of the mutations.
   *
   * @remarks
   * A race with no observable effect cannot be asserted on.
   */
  count(): Promise<number>;
}
