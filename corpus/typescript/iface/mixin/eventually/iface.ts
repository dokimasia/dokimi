// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the eventually mixin, which states that
 * the effect becomes visible after a quiet period rather than immediately.
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

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=eventuallytest/
export interface Mixed {
  /**
   * Publishes the item, which is not observable on return.
   *
   * @remarks
   * The law is that it becomes so, which needs a settle step, an anti-entropy
   * exchange, and a read to observe after both — the members the directive
   * names. observe= names the read the effect is expected to reach. Without
   * it the directive states the property and no law can bind: settling is
   * getting to quiescence, and quiescence is not an observation.
   */
  //dokimi:mixin eventually settle=settle sync=sync observe=items
  publish(item: string): Promise<void>;

  /** Advances the quiet window that the law waits out. */
  settle(): Promise<void>;

  /**
   * Pulls the settled state of one peer into this replica.
   *
   * @remarks
   * sync is the pairwise anti-entropy step that the law composes into a full
   * round.
   */
  sync(peer: Mixed): Promise<void>;

  /**
   * Returns the eventual state in sorted order.
   *
   * @remarks
   * Two converged replicas return identical items.
   */
  items(): Promise<string[]>;
}
