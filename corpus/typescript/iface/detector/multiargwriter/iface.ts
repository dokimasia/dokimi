// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the multiargwriter shape: three or
 * more parameters.
 *
 * @remarks
 * Three parameters are past the composite pair, at the arity boundary the
 * detector documents and draws.
 *
 * No classifying directive appears here: the classification comes from the
 * signature alone, so a detector that misfires shows up as wrong generated
 * output rather than as a directive that was misread.
 *
 * @packageDocumentation
 */

/** Reports a key that {@link MultiArgWriter.set} never filed. */
export const notFound = new Error("multiargwriter: not found");

/** Represents the stored form. */
export interface Value {
  /** Contains the key of the value. */
  key: string;
  /** Contains the body of the value. */
  body: string;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=multiargwritertest/
export interface MultiArgWriter {
  /**
   * Files the body under the key with its media type.
   *
   * @remarks
   * set carries three parameters — past the composite pair, into the boundary
   * the detector draws at three or more.
   */
  set(key: string, body: string, mime: string): Promise<void>;

  /**
   * Returns the body that set filed under the key.
   *
   * @remarks
   * get is here so the model tier has something to compare. A lone writer
   * answers nothing, so both sides return the same for every call a correct
   * subject makes and the differential has nothing to disagree about.
   */
  get(key: string): Promise<string>;
}
