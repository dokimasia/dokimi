// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the lifecycle shape: a failable
 * teardown.
 *
 * @remarks
 * close must be idempotent and every later operation must report the closed
 * sentinel, which is the law the shape carries.
 *
 * No classifying directive appears here: the classification comes from the
 * signature alone, so a detector that misfires shows up as wrong generated
 * output rather than as a directive that was misread.
 *
 * @packageDocumentation
 */

/** Reports an operation after teardown. */
export const closed = new Error("lifecycle: closed");

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=lifecycletest/
export interface Lifecycle {
  /**
   * Tears the subject down, and declares the teardown safe to repeat.
   *
   * @remarks
   * The idempotence is what a double close relies on, which the mixin states
   * because no signature can.
   */
  //dokimi:mixin idempotent
  close(): Promise<void>;
}
