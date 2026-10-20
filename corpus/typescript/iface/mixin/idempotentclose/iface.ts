// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the idempotent mixin on a teardown.
 *
 * @remarks
 * The mixin already has a write-shaped fixture; this one carries the same
 * claim on a lifecycle method, where "repeating the call changes nothing" is
 * the double close a cleanup relies on.
 *
 * The interface pairs the teardown with an observation, because the law
 * reads state across the second close and a close-only interface offers no
 * method to read it through.
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
//dokimi:out path=idempotentclosetest/
export interface Closer {
  /** Tears the subject down, and declares the teardown safe to repeat. */
  //dokimi:mixin idempotent
  close(): Promise<void>;

  /**
   * Returns the count of the open resources.
   *
   * @remarks
   * The second close must leave this observation unchanged, which a
   * close-only interface cannot state.
   */
  stats(): Promise<number>;
}
