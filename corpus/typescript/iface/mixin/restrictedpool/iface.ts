// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for pool provenance: a declaration whose
 * input domain a run can narrow, and a model tier that has to respect the
 * narrowing.
 *
 * @remarks
 * Every other fixture draws its sequences from a pool the generator derived,
 * and widening a derived pool costs nothing — nobody said what the
 * implementation accepts, so arbitrary draws are fair. This one says it. The
 * roles below carry pools a run replaces by passing a config, and a run that
 * replaced them has stated what the subject takes. A tier that blended
 * arbitrary or hostile draws past that statement would red correct code
 * against inputs its owner ruled out, which is worse than not probing: the
 * red is real, the defect is not, and the fix is to weaken a check that was
 * right.
 *
 * The pair is the smallest shape that carries both halves. Roles on named
 * types are what produce the config a run overrides; the keyed read and
 * write are what give the model tier sequences to draw for.
 *
 * @packageDocumentation
 */

/** Reports a key that nothing wrote. */
export const notFound = new Error("restrictedpool: not found");

/**
 * Identifies one record.
 *
 * @remarks
 * Named rather than a bare string because a role is written on a
 * declaration, and a bare parameter has none. Both Key and Body declare their
 * default through the directive, because a named type has no field to carry
 * one.
 */
//dokimi:default "test-key"
export type Key = string;

/** Represents what a record contains. */
//dokimi:default "test-body"
export type Body = string;

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=restrictedpooltest/
export interface Store {
  /** Writes a body under a key, which draws both roles. */
  put(key: Key, body: Body): Promise<void>;

  /**
   * Returns the body under the key.
   *
   * @remarks
   * get reports a miss with the declared sentinel {@link notFound}.
   */
  get(key: Key): Promise<Body>;
}
