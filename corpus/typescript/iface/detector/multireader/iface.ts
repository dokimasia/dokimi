// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the multireader shape: one key in,
 * two values out or a rejection.
 *
 * @remarks
 * The shape is the failable counterpart of a lookup.
 *
 * No classifying directive appears here: the classification comes from the
 * signature alone, so a detector that misfires shows up as wrong generated
 * output rather than as a directive that was misread.
 *
 * @packageDocumentation
 */

/** Reports a key that the reads of the fixture miss. */
export const notFound = new Error("multireader: not found");

/** Represents the primary return. */
export interface Value {
  /** Contains the key of the value. */
  key: string;
  /** Contains the body of the value. */
  body: string;
}

/** Represents the secondary return. */
export interface Meta {
  /** Contains the revision of the value. */
  revision: number;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=multireadertest/
export interface MultiReader {
  /** Returns the value under the key with its metadata. */
  getWithMeta(key: string): Promise<[Value, Meta]>;
}
