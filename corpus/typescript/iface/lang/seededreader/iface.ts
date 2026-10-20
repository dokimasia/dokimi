// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for the SEED SEAM: an interface nothing
 * can write to.
 *
 * @remarks
 * A run populates a subject through the surface under test, and this surface
 * offers no way in. A read against an empty subject cannot tell a miss from a
 * bug, so the harness receives the corpus instead of a bare constructor — the
 * suite decides what is seeded, from its own pools, and the harness decides
 * how loading happens.
 *
 * Both roles are stamped because the corpus is zipped from them: a key pool
 * with no payload pool names half a map.
 *
 * @packageDocumentation
 */

/**
 * Identifies one document.
 *
 * @remarks
 * Both Key and Body declare their default on the type, because a named type
 * has no field to carry one.
 */
//dokimi:default "test-doc"
export type Key = string;

/** Represents what a document contains. */
//dokimi:default "test-contents"
export type Body = string;

/** Represents a catalog that reads documents and writes none. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=seededreadertest/
export interface Catalog {
  /** Returns the body of a seeded key, and reports a miss for any other. */
  lookup(key: Key): Promise<Body>;

  /** Returns the count of the documents that the harness loaded. */
  len(): Promise<number>;
}
