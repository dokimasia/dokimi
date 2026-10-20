// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the voidlifecycle shape: a teardown
 * that cannot fail.
 *
 * @remarks
 * The idempotence law is all that remains: calling it twice must be safe.
 * TypeScript declares no error result, so this signature is the closer's
 * under another name.
 *
 * No classifying directive appears here: the classification comes from the
 * signature alone, so a detector that misfires shows up as wrong generated
 * output rather than as a directive that was misread.
 *
 * @packageDocumentation
 */

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=voidlifecycletest/
export interface VoidLifecycle {
  /** Stops the subject. */
  stop(): void;
}
