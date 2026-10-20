// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the pure mixin, which states that the
 * method has no side effects and repeated calls agree.
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
//dokimi:out path=puretest/
export interface Mixed {
  /**
   * Returns a result computed from the input alone.
   *
   * @remarks
   * A promise would suggest I/O, which is what purity excludes.
   */
  //dokimi:mixin pure
  derive(input: string): string;
}
