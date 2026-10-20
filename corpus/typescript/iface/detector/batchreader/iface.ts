// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the batchreader shape: many keys in
 * one call.
 *
 * @remarks
 * The rest parameter of keys is the marker, and the returned array need not
 * match the request length — absent keys are simply absent.
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
//dokimi:out path=batchreadertest/
export interface BatchReader {
  /** Returns the values under the keys, without an entry for an absent key. */
  getAll(...keys: string[]): Promise<Value[]>;
}
