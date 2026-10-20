// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the monotonic mixin, which states that
 * successive reads never go backwards.
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
//dokimi:out path=monotonictest/
export interface Mixed {
  /**
   * Returns the version of the subject, which never decreases across calls.
   *
   * @remarks
   * A single read cannot violate that, so the law is over a sequence and
   * advance is what moves it.
   */
  //dokimi:mixin monotonic
  version(): Promise<number>;

  /** Advances the version that the law watches. */
  advance(): Promise<void>;
}
