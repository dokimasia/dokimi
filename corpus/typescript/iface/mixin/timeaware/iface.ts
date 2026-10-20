// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the timeaware mixin, which marks a
 * callable whose behaviour depends on a clock.
 *
 * @remarks
 * The mixin does not say which quantity, because the quantities belong to the
 * claims layered on top (`ttl`, `timeout`, `windowed`).
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
//dokimi:out path=timeawaretest/
export interface Mixed {
  /** Records that the key was seen, at whatever the clock reads. */
  touch(key: string): Promise<void>;

  /**
   * Returns how long ago the key was seen, in milliseconds.
   *
   * @remarks
   * The age is the dependency the mixin marks: the answer moves when the
   * clock does and not otherwise.
   */
  //dokimi:mixin time-aware
  ageOf(key: string): Promise<number>;
}
