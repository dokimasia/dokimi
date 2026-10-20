// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the conservative mixin, which states
 * that the operation preserves a quantity the interface names.
 *
 * @remarks
 * The interface carries a total beside the step: an algebraic claim about a
 * fold needs somewhere to read the fold's result, or the generated subtest
 * passes over a state nothing can see.
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
  /** Contains the amount of the contribution, which the operation preserves. */
  amount: number;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=conservativetest/
export interface Mixed {
  /** Moves a quantity without creating or destroying any of it. */
  //dokimi:mixin conservative field=amount
  apply(d: Delta): Promise<void>;

  /** Returns the fold, which is what makes the claim about apply checkable at all. */
  total(): Promise<number>;
}
