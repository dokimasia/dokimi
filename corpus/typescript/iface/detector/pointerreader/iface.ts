// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the pointerreader shape: an optional
 * return.
 *
 * @remarks
 * A miss is expressible as undefined rather than as a value. That makes the
 * nil-safety law meaningful here in a way it is not for a value return.
 *
 * No classifying directive appears here: the classification comes from the
 * signature alone, so a detector that misfires shows up as wrong generated
 * output rather than as a directive that was misread.
 *
 * @packageDocumentation
 */

/** Represents the payload that the fixture reads. */
export interface Value {
  /** Contains the key of the payload. */
  key: string;
  /** Contains the body of the payload. */
  body: string;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=pointerreadertest/
export interface PointerReader {
  /** Returns the value under the key, or undefined for a miss. */
  find(key: string): Promise<Value | undefined>;
}
