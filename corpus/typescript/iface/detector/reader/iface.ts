// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the reader shape: a single-key fetch
 * that can fail.
 *
 * @remarks
 * The reader is the commonest shape in any store-like interface, and what
 * every miss-sentinel law attaches to.
 *
 * No classifying directive appears here: the classification comes from the
 * signature alone, so a detector that misfires shows up as wrong generated
 * output rather than as a directive that was misread.
 *
 * @packageDocumentation
 */

/** Reports a key that the reads of the fixture miss. */
export const notFound = new Error("reader: not found");

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
//dokimi:out path=readertest/
export interface Reader {
  /**
   * Returns the value under the key.
   *
   * @remarks
   * The fault line carries all three settings at once — a sentinel, a retry
   * attempt, and the partition field — because the settings are independent
   * and a projection that dropped one when another was present would still
   * pass every fixture carrying them apart. The double is routed into
   * readertest, so the sentinel reference this renders has to import from the
   * module of the fixture.
   */
  //dokimi:mixin not-found sentinel=notFound
  //dokimi:fault errors=[notFound] retry=3 partition=key
  get(key: string): Promise<Value>;
}
