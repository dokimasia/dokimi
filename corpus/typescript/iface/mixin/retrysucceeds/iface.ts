// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the retrysucceeds mixin, which states
 * that a transient failure succeeds on a later attempt.
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
//dokimi:out path=retrysucceedstest/
export interface Mixed {
  /**
   * Calls the operation under the key, which fails transiently before it
   * succeeds.
   *
   * @remarks
   * The law drives it repeatedly, so a method that either always works or
   * never does cannot host it.
   *
   * The fault directive carries only a key — no sentinel — which is the shape
   * a retry-loop test wants: retrySchedule takes whatever error the test
   * supplies, and forcing a named sentinel here would invent a helper nobody
   * asked for.
   */
  //dokimi:mixin retry-succeeds attempts=3
  //dokimi:fault retry=3
  call(key: string): Promise<void>;

  /**
   * Returns the count of the attempts that the subject has seen.
   *
   * @remarks
   * The count is how the law distinguishes a retry from a first success.
   */
  attempts(): Promise<number>;
}
