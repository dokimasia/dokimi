// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the noduplicates mixin, which states
 * that a drain yields each element at most once.
 *
 * @remarks
 * The interface carries an append beside the drain: a claim about what a
 * drain yields needs something to have gone in, or the generated subtest
 * passes over an empty collection.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/** Represents the element that the collection contains. */
export interface Value {
  /** Contains the key of the element. */
  key: string;
  /** Contains the body of the element. */
  body: string;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=noduplicatestest/
export interface Mixed {
  /** Adds an element for the drain to yield. */
  add(v: Value): Promise<void>;

  /** Drains the collection, where no element appears twice in one drain. */
  //dokimi:mixin no-duplicates
  items(): Promise<Value[]>;
}
