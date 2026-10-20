// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the readerwithbool shape: the
 * map-style fetch.
 *
 * @remarks
 * Presence rides in the result rather than in a rejection, so a miss is an
 * ordinary outcome instead of a failure. TypeScript writes presence as an
 * optional result, so this signature is the pointerreader's.
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
//dokimi:out path=readerwithbooltest/
export interface ReaderWithBool {
  /** Returns the value under the key, or undefined for a miss. */
  load(key: string): Promise<Value | undefined>;
}
