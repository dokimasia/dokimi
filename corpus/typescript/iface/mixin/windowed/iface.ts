// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the windowed mixin, which states that a
 * count covers a bounded interval rather than all of history.
 *
 * @remarks
 * The directive names the pair it governs and the interval itself, because a
 * window nobody stated is a number the generator would have to invent.
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
//dokimi:out path=windowedtest/
export interface Mixed {
  /** Adds one occurrence of the key at the current reading of the clock. */
  //dokimi:mixin windowed incr=record count=countIn window=1m
  record(key: string): Promise<void>;

  /**
   * Returns the count of the occurrences of the key inside the declared
   * window.
   *
   * @remarks
   * The window is what makes the count bounded rather than cumulative.
   */
  countIn(key: string): Promise<number>;
}
