// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the sideeffect mixin, which states that
 * the method mutates state beyond its return value.
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
//dokimi:out path=sideeffecttest/
export interface Mixed {
  /**
   * Marks the key as used.
   *
   * @remarks
   * touch returns nothing useful, so its entire effect is out of band. That
   * is what the mixin declares, and observed is what makes it visible.
   */
  //dokimi:mixin side-effect observe=observed
  touch(key: string): Promise<void>;

  /** Returns the count of the touches of the key, which exposes the out-of-band effect. */
  observed(key: string): Promise<number>;
}
