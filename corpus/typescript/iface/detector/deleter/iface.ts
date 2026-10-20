// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the deleter shape: a removal named
 * for it — the key of the entity to remove in, nothing out.
 *
 * @remarks
 * Structurally this is the writer shape, and the name is the discrimination:
 * nothing in the signature separates a delete from a put, and the two license
 * opposite expectations — a law deriving write-then-read-back against a
 * delete asserts the reverse of correct behaviour. The name gate is what
 * keeps that from happening.
 *
 * No classifying directive appears here: the classification comes from the
 * signature and the name alone, so a detector that misfires shows up as
 * wrong generated output rather than as a directive that was misread.
 *
 * @packageDocumentation
 */

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=deletertest/
export interface Deleter {
  /** Removes the entity under the key. */
  delete(key: string): Promise<void>;
}
