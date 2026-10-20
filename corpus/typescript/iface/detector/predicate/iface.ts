// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the predicate shape: a question
 * about current state, with no arguments and no failure mode.
 *
 * @remarks
 * The empty parameter list is the signature's distinguishing feature.
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
//dokimi:out path=predicatetest/
export interface Predicate {
  /** Reports whether the subject has no entries. */
  isEmpty(): boolean;
}
