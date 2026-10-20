// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the streamreader shape: values
 * yielded lazily as an async iterable, which rejects for a failed element.
 *
 * @remarks
 * Consumers may stop early, so the implementation must not assume the
 * sequence is drained.
 *
 * No classifying directive appears here: the classification comes from the
 * signature alone, so a detector that misfires shows up as wrong generated
 * output rather than as a directive that was misread.
 *
 * @packageDocumentation
 */

/** Represents the element type that the stream yields. */
export interface Value {
  /** Contains the key of the element. */
  key: string;
  /** Contains the body of the element. */
  body: string;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=streamreadertest/
export interface StreamReader {
  /**
   * Adds a value for the stream to yield.
   *
   * @remarks
   * Without a writer the subject stays empty forever, and every claim about
   * what a drain returns holds by vacancy — two drains of nothing agree, a
   * drain of nothing terminates, and no defect worn on list can make either
   * false.
   */
  add(v: Value): Promise<void>;

  /** Returns the values that add stored, as a stream. */
  list(): AsyncIterable<Value>;
}
