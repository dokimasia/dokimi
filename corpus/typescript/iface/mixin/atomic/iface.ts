// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the atomic mixin, which states that a
 * failed write leaves no partial state, so a reader never sees half of it.
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

/**
 * Represents the two-field write that the mixin keeps together.
 *
 * @remarks
 * A failed write applies neither field, and read returns both so a
 * half-applied one is observable.
 */
export interface Entry {
  /** Contains the key of the entry. */
  key: string;
  /** Contains the left half of the write. */
  left: string;
  /** Contains the right half of the write. */
  right: string;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=atomictest/
export interface Mixed {
  /**
   * Applies both of the entry's fields or neither.
   *
   * @remarks
   * One input rather than several, because the law's claim is about one
   * write's halves — the fields travel together so the observation can too.
   * The read= names what a check looks through afterwards: the assertion is
   * about what is NOT there after a failure, and a check that cannot look
   * asserts only that a successful call succeeded.
   */
  //dokimi:mixin atomic read=read
  write(e: Entry): Promise<void>;

  /** Returns the whole entry under the key, so a partial write is observable. */
  read(key: string): Promise<Entry>;
}
