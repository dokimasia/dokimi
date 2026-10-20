// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the poisonaccessor shape: a latched
 * failure state.
 *
 * @remarks
 * Once poisoned every later call reports the same error, which is what makes
 * the accessor worth having: the failure outlives the call that caused it.
 *
 * No classifying directive appears here: the classification comes from the
 * signature alone, so a detector that misfires shows up as wrong generated
 * output rather than as a directive that was misread.
 *
 * @packageDocumentation
 */

/** Reports the latched failure of a tripped subject. */
export const poisoned = new Error("poisonaccessor: poisoned");

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=poisonaccessortest/
export interface PoisonAccessor {
  /** Returns the latched failure, or undefined before the subject trips. */
  err(): Error | undefined;
}
