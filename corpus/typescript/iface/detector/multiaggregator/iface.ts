// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the multiaggregator shape: several
 * reductions at once.
 *
 * @remarks
 * Returning them together rather than through repeated calls is the point:
 * the numbers only mean anything if they describe one instant.
 *
 * No classifying directive appears here: the classification comes from the
 * signature alone, so a detector that misfires shows up as wrong generated
 * output rather than as a directive that was misread.
 *
 * @packageDocumentation
 */

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=multiaggregatortest/
export interface MultiAggregator {
  /** Returns two reductions of the collection at one instant. */
  stats(): Promise<[number, number]>;
}
