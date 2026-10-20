// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the associative mixin, which states
 * that combining in any grouping yields the same value.
 *
 * @remarks
 * The interface carries a total beside the step: an algebraic claim about a
 * fold needs somewhere to read the fold's result, or the generated subtest
 * passes over a state nothing can see.
 *
 * The step takes a BATCH, and that is what makes the claim witnessable at
 * all. Grouping has to be something a caller can vary: with one delta per
 * call the only grouping is the call sequence itself, every implementation
 * satisfies the claim by having no second grouping to disagree with, and a
 * check would pass over a subject that folds a batch inconsistently. Two
 * calls carrying the same deltas split differently — [a,b] then [c] against
 * [a] then [b,c] — are two groupings, and the totals must agree.
 *
 * This is the axis on which it differs from commutative, which varies the
 * ORDER and needs no batch to do it.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/** Represents one contribution to the fold. */
export interface Delta {
  /** Contains the key of the contribution. */
  key: string;
  /** Contains the amount of the contribution. */
  amount: number;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=associativetest/
export interface Mixed {
  /**
   * Folds a batch of deltas in.
   *
   * @remarks
   * Splitting the same deltas across calls differently does not change the
   * total.
   */
  //dokimi:mixin associative
  apply(ds: Delta[]): Promise<void>;

  /** Returns the fold, which is what makes the claim about apply checkable at all. */
  total(): Promise<number>;
}
