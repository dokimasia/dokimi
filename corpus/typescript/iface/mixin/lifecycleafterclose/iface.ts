// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the lifecycleafterclose mixin, which
 * states that every call after teardown reports the closed sentinel.
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

/** Reports an operation after teardown. */
export const closed = new Error("lifecycleafterclose: closed");

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=lifecycleafterclosetest/
export interface Mixed {
  /** Tears the subject down, which is where the claim measures from. */
  close(): Promise<void>;

  /**
   * Does the work of the subject.
   *
   * @remarks
   * work must report {@link closed} once close has run — the mixin rides the
   * operation, because the claim is about what the operation answers, and
   * its parameters name the teardown and the sentinel the refusal carries.
   */
  //dokimi:mixin lifecycle-after-close close=close sentinel=closed
  work(): Promise<void>;
}
