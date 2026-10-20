// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the orderafter mixin, which states that
 * the call is only valid after a named predecessor has run.
 *
 * @remarks
 * The interface carries whatever methods the law needs to be stateable. A
 * mixin whose law spans two calls cannot be hosted by a single method: there
 * would be nothing to compare against, and the generated subtest would pass
 * by having nothing to check.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/**
 * Reports a commit before prepare has run.
 *
 * @remarks
 * The error is the ordering's own sentinel, declared beside the shape so the
 * directive can name it.
 */
export const notReady = new Error("orderafter: not ready");

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=orderaftertest/
export interface Mixed {
  /**
   * Commits the work, which is valid only after prepare.
   *
   * @remarks
   * The fn parameter names the predecessor, and the resolver resolves it to a
   * sibling — so prepare has to exist in this interface or the reference
   * dangles. The unready key is what turns "refused" into "refused for being
   * early": without it, an implementation failing on a missing map passes the
   * ordering check.
   */
  //dokimi:mixin order-after fn=prepare unready=notReady
  commit(): Promise<void>;

  /** Prepares the commit, as the predecessor that fn names. */
  prepare(): Promise<void>;
}
