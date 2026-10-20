// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the readafterwrite mixin, which states
 * that a value written is immediately readable.
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
//dokimi:out path=readafterwritetest/
export interface Mixed {
  /**
   * Returns the value that write stored under the key.
   *
   * @remarks
   * The write parameter names the writer, so the pair has to be in one
   * interface — this is the fixture the single-method template could not
   * express.
   */
  //dokimi:mixin read-after-write write=write
  read(key: string): Promise<string>;

  /** Stores the value under the key, as the writer that the mixin names. */
  write(key: string, value: string): Promise<void>;
}
