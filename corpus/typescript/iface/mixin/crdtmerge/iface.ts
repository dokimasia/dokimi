// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the crdtmerge mixin, which states that
 * replicas converge to the same state regardless of merge order.
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
 * Represents the peer that merge folds in.
 *
 * @remarks
 * The law compares two replicas merged in opposite orders, so the type has
 * to be nameable in the signature.
 */
//dokimi:stub
//dokimi:out path=crdtmergetest/
export interface Replica {
  /** Returns the items of the replica. */
  items(): Promise<string[]>;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=crdtmergetest/
export interface Mixed {
  /**
   * Folds another replica in.
   *
   * @remarks
   * Convergence is a statement about two merges in opposite orders, so the
   * method has to take a peer rather than a value.
   */
  //dokimi:mixin crdt-merge
  merge(peer: Replica): Promise<void>;

  /** Adds an item, which introduces divergence for the merge to reconcile. */
  add(item: string): Promise<void>;

  /** Returns the items, which is how the law observes convergence. */
  items(): Promise<string[]>;
}
