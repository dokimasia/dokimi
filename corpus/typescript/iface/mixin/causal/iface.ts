// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the causal mixin, which states that
 * reads respect the happens-before order of the writes that produced them.
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
   * Contains the ordering stamp that the store assigns.
   *
   * @remarks
   * The causal claim is judged by this stamp, which the directive's version=
   * names.
   */
  rev: number;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=causaltest/
export interface Mixed {
  /**
   * Stores the value, and returns the stored value.
   *
   * @remarks
   * store is the write the read is judged against. It returns the stored
   * value, because the store assigns rev and a caller who never sees the
   * stamp has no order to be consistent with.
   */
  store(v: Value): Promise<Value>;

  /**
   * Returns the value under the key.
   *
   * @remarks
   * The mixin says a read never observes an effect without having observed
   * its cause, and version= names the stamp that orders the causes.
   */
  //dokimi:mixin causal version=rev
  get(key: string): Promise<Value>;
}
