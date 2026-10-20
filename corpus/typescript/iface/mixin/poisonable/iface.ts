// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the poisonable mixin, which states that
 * a subject driven into a failed state reports it consistently rather than
 * intermittently.
 *
 * @remarks
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=poisonabletest/
export interface Mixed {
  /** Drives the subject into the failed state that probe reports. */
  fail(): Promise<void>;

  /**
   * Returns the failure of the subject, or undefined while it is healthy.
   *
   * @remarks
   * probe takes nothing, which is the shape a stored-error accessor has: the
   * answer is state rather than a lookup.
   */
  //dokimi:mixin poisonable induce=fail
  probe(): Error | undefined;
}
