// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the deprecated mixin, which states that
 * the method is scheduled for removal and callers should migrate.
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
//dokimi:out path=deprecatedtest/
export interface Mixed {
  /**
   * Returns the value under the key.
   *
   * @remarks
   * old carries no runtime law: the mixin is an annotation the generator
   * surfaces, so the signature is deliberately unremarkable.
   *
   * @deprecated Use {@link Mixed.new} instead.
   */
  //dokimi:mixin deprecated
  old(key: string): Promise<string>;

  /**
   * Returns the value under the key, as the replacement of old.
   *
   * @remarks
   * The name new is in quotes because `new(` opens a construct signature in
   * an interface.
   */
  "new"(key: string): Promise<string>;
}
