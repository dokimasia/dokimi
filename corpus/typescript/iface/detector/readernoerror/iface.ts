// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the readernoerror shape: an
 * infallible fetch, whose miss the value has to carry.
 *
 * @remarks
 * TypeScript declares no error result, so this signature is the reader's.
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
//dokimi:out path=readernoerrortest/
export interface ReaderNoError {
  /** Returns the value under the key, or an empty value for a miss. */
  lookup(key: string): Promise<Value>;
}
