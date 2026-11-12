// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for the directive refuse, which
 * withdraws a derived claim of a callable for every subject of its suite,
 * with a reason.
 *
 * @remarks
 * update fails only for an instant beyond the max offset of the clock, and no
 * sample of an instant is beyond it. A check of the claim zero-on-error on the
 * samples would never engage, so update refuses the claim.
 *
 * @packageDocumentation
 */

/** Reports an instant beyond the max offset of the clock. */
export const offset = new Error("refused: the instant is beyond the max offset");

/**
 * Represents an instant of a hybrid logical clock: a wall time in
 * nanoseconds since the Unix epoch, and a logical count of the events at that
 * wall time.
 */
export interface Instant {
  /** Contains the wall time in nanoseconds since the Unix epoch. */
  wall: bigint;
  /** Contains the count of the events at the wall time. */
  logical: number;
}

/** Represents a hybrid logical clock. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=refusedtest/
export interface Clock {
  /**
   * Merges the instant of a peer into the clock, and returns the instant of
   * the clock after the merge.
   *
   * @throws {@link offset} for an instant beyond the max offset of the clock.
   */
  //dokimi:refuse zero-on-error reason="update fails only for an instant beyond the max offset of the clock"
  update(observed: Instant): Instant;
}
