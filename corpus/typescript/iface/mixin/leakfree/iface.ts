// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the leakfree mixin, which states that a
 * resource acquired is a resource released.
 *
 * @remarks
 * Repeating the cycle leaves nothing outstanding. The interface carries both
 * halves of the cycle and a way to observe the balance, because the claim is
 * about the pair: an acquire alone leaks by definition, and a release alone
 * has nothing to release.
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
//dokimi:out path=leakfreetest/
export interface Mixed {
  /**
   * Acquires the resource.
   *
   * @remarks
   * The mixin names both halves of the cycle whose balance is the claim.
   */
  //dokimi:mixin leak-free open=acquire close=release
  acquire(): Promise<void>;

  /** Releases the resource. */
  release(): Promise<void>;

  /**
   * Returns the count of the resources that are acquired and not released.
   *
   * @remarks
   * The balance is what makes a leak visible rather than merely eventual.
   */
  outstanding(): Promise<number>;
}
