// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the answeringwriter shape: one value
 * in, the same type out — a write that answers the stored state it produced.
 *
 * @remarks
 * The boundary is type identity on a named type. A basic-typed coincidence —
 * a string-keyed read answering a string — stays with the reader, because
 * the answered state this shape exists for is a record carrying identity and
 * stamps.
 *
 * No classifying directive appears here: the classification comes from the
 * signature alone, so a detector that misfires shows up as wrong generated
 * output rather than as a directive that was misread.
 *
 * @packageDocumentation
 */

/**
 * Represents the stored form that the write returns.
 *
 * @remarks
 * The stamped state is what separates this shape from a transformer.
 */
export interface Value {
  /** Contains the key of the value. */
  key: string;
  /** Contains the body of the value. */
  body: string;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=answeringwritertest/
export interface AnsweringWriter {
  /**
   * Stores the value and returns the stored state.
   *
   * @remarks
   * The same named type in and out is the identity the detector draws on.
   */
  put(v: Value): Promise<Value>;
}
