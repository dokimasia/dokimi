// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the writesfollowreads mixin, which
 * states that a write made after a read is ordered after the write that read
 * observed.
 *
 * @remarks
 * The interface carries a write and a read because the claim relates them:
 * a read alone has nothing to be consistent with, and the generated subtest
 * would pass by having nothing to check.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/** Represents the payload that the store keeps. */
export interface Value {
  /** Contains the key of the payload. */
  key: string;
  /** Contains the body of the payload. */
  body: string;
  /** Contains the ordering stamp that the subject assigns on write, which the mixin names. */
  rev: number;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=writesfollowreadstest/
export interface Mixed {
  /**
   * Stores the value, and returns the stored state.
   *
   * @remarks
   * store is the write the read is judged against. The rev it assigned rides
   * back to the caller, which is what lets a trace order this client's
   * writes.
   */
  store(v: Value): Promise<Value>;

  /** Returns the value under the key, which a later write is ordered against. */
  //dokimi:mixin writes-follow-reads version=rev
  get(key: string): Promise<Value>;
}
