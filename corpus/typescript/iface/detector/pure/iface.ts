// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the pure shape: a value derived from
 * nothing but the object, with no arguments and no promise.
 *
 * @remarks
 * Repeated calls must agree, which is the law the shape carries.
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
//dokimi:out path=puretest/
export interface Pure {
  /** Returns a description of the subject. */
  describe(): string;
}
