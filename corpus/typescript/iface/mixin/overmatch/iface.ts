// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the overmatch mixin, which states that
 * a drain yields at least the elements required of it, and may yield more.
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
//dokimi:out path=overmatchtest/
export interface Mixed {
  /** Adds an element for the drain to yield. */
  add(v: Value): Promise<void>;

  /** Drains the collection, and yields a superset of what was required. */
  //dokimi:mixin over-match
  items(): Promise<Value[]>;
}
