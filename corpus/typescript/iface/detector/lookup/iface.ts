// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the lookup shape: two values or
 * nothing.
 *
 * @remarks
 * The second value is what separates it from an optional reader: metadata
 * arrives with the value, not from a second call.
 *
 * No classifying directive appears here: the classification comes from the
 * signature alone, so a detector that misfires shows up as wrong generated
 * output rather than as a directive that was misread.
 *
 * @packageDocumentation
 */

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
//dokimi:out path=lookuptest/
export interface Lookup {
  /** Returns the value under the key with its metadata, or undefined for a miss. */
  inspect(key: string): [Value, Meta] | undefined;
}
