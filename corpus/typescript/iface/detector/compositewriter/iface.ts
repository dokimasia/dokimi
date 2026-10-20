// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the compositewriter shape: a key
 * beside a value and nothing else.
 *
 * @remarks
 * Exactly two parameters is the boundary the detector documents and draws.
 *
 * No classifying directive appears here: the classification comes from the
 * signature alone, so a detector that misfires shows up as wrong generated
 * output rather than as a directive that was misread.
 *
 * @packageDocumentation
 */

/** Represents the stored form. */
export interface Value {
  /** Contains the key of the value. */
  key: string;
  /** Contains the body of the value. */
  body: string;
  /**
   * Contains the revision of the value.
   *
   * @remarks
   * The implementation assigns the revision, which is why the write has to
   * return it.
   */
  rev: string;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=compositewritertest/
export interface CompositeWriter {
  /**
   * Stores the value under the key.
   *
   * @remarks
   * set is the documented shape whole: a key beside a value and nothing else
   * — exactly two parameters, which is the boundary the detector draws.
   */
  set(key: string, v: Value): Promise<void>;
}
