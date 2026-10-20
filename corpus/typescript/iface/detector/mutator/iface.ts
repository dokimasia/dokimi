// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the mutator shape: state change
 * with nothing returned.
 *
 * @remarks
 * The only observable effect is on a later read. TypeScript declares no
 * error result, so this signature is the writer's with a key in place of a
 * value.
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
//dokimi:out path=mutatortest/
export interface Mutator {
  /** Marks the entry under the key as used. */
  touch(key: string): Promise<void>;
}
