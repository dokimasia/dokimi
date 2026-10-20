// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the concurrentreaders mixin, which
 * states that concurrent reads stay safe while a writer runs.
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
//dokimi:out path=concurrentreaderstest/
export interface Mixed {
  /**
   * Returns the value under the key.
   *
   * @remarks
   * The law drives get concurrently against put. Reads alone cannot violate
   * this, so the interface has to carry the writer that contends with them.
   */
  //dokimi:mixin concurrent-readers
  get(key: string): Promise<string>;

  /** Stores the value under the key, as the contending writer. */
  put(key: string, value: string): Promise<void>;
}
