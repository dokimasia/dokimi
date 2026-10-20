// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the cacheable mixin, which states that
 * repeated reads may be served from a cache without changing the answer.
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
//dokimi:out path=cacheabletest/
export interface Mixed {
  /**
   * Returns the value under the key, which a cache may serve.
   *
   * @remarks
   * The law is that caching is unobservable: two reads agree whether or not
   * the second was served from memory.
   */
  //dokimi:mixin cacheable
  get(key: string): Promise<string>;
}
