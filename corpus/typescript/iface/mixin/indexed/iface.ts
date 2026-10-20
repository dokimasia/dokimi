// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the indexed mixin, which states that a
 * number parameter addresses a position in the collection another method
 * sizes.
 *
 * @remarks
 * The shape the fixture axis had none of, and the one that made the gap
 * visible: a derived value for a bare `number` is a number invented from the
 * type, and an index invented from the type is out of range for every
 * collection smaller than it. Nothing in the signature says the parameter is
 * a position rather than a magnitude, so nothing could have known.
 *
 * `by=len` names the sizing method as a sibling callable, so a misspelling
 * fails at the directive. It cannot supply the bound: the size is a fact
 * about the seeded subject at run time, not about the declaration.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/** Represents the element that the positions address. */
export interface Value {
  /** Contains the key of the element. */
  key: string;
  /** Contains the body of the element. */
  body: string;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=indexedtest/
export interface Ranked {
  /** Appends an element, so the positions have something to address. */
  add(v: Value): Promise<void>;

  /**
   * Returns the count of the elements that the positions address.
   *
   * @remarks
   * The count is the bound that the mixin names, and an ordinary read that
   * the sequences still drive.
   */
  len(): Promise<number>;

  /**
   * Returns the element at a position.
   *
   * @remarks
   * at reports a miss for a position that the collection does not contain.
   */
  //dokimi:mixin indexed by=len
  at(i: number): Promise<Value>;
}
