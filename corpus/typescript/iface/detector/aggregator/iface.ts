// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the aggregator shape:
 * whole-collection state reduced to one value.
 *
 * @remarks
 * Taking no key is the distinction from a reader: there is nothing to look
 * up, only something to compute.
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
//dokimi:out path=aggregatortest/
export interface Aggregator {
  /** Returns the count of the entries of the collection. */
  count(): Promise<number>;
}
