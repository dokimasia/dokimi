// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the monotonicreads mixin, which states
 * that successive reads by one client never move backwards.
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
  /**
   * Contains the ordering stamp that the subject assigns on write.
   *
   * @remarks
   * The per-client guarantee is defined against this version, which is why
   * the mixin names it.
   */
  rev: number;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=monotonicreadstest/
export interface Mixed {
  /** Stores the value, as the write the read is judged against. */
  store(v: Value): Promise<void>;

  /**
   * Returns the value under the key.
   *
   * @remarks
   * Once a client has seen a version, it never sees an older one — ordered by
   * the rev the store stamped.
   */
  //dokimi:mixin monotonic-reads version=rev
  get(key: string): Promise<Value>;
}
