// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the writer shape: a value in,
 * nothing out.
 *
 * @remarks
 * The absence of a returned value separates it from an answering writer.
 *
 * No classifying directive appears here: the classification comes from the
 * signature alone, so a detector that misfires shows up as wrong generated
 * output rather than as a directive that was misread.
 *
 * @packageDocumentation
 */

/** Represents the payload that the fixture stores. */
export interface Value {
  /** Contains the key of the payload. */
  key: string;
  /** Contains the body of the payload. */
  body: string;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=writertest/
export interface Writer {
  /** Stores the value. */
  put(v: Value): Promise<void>;
}
