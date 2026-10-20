// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the detector-axis fixture for the closer shape: the bare teardown
 * — a callable taking nothing and returning nothing, named for release.
 *
 * @remarks
 * This is the catalogue's one name-gated detector, and the fixture holds
 * that: `close(): void` and `stop(): void` of the voidlifecycle fixture are
 * the same signature in TypeScript, so the name is what separates a release
 * from any other call without arguments. A teardown is expected to be allowed
 * to answer differently the second time, which is why no repeat claim is
 * stated at the floor.
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
//dokimi:out path=closertest/
export interface Closer {
  /** Releases the resources of the subject. */
  close(): void;
}
