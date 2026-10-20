// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the nilsafe mixin, which states that an
 * absent argument is handled rather than throwing.
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

/** Represents the payload that store takes as an optional value, so the absent case exists at all. */
export interface Payload {
  /** Contains the key of the payload. */
  key: string;
  /** Contains the body of the payload. */
  body: string;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=nilsafetest/
export interface Mixed {
  /**
   * Stores the payload.
   *
   * @remarks
   * The parameter may be undefined, which is what makes absence expressible.
   * A required parameter would leave the law with no absence to pass.
   */
  //dokimi:mixin nil-safe
  store(v: Payload | undefined): Promise<void>;
}
