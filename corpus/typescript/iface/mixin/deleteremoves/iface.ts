// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the deleteremoves mixin, which states
 * that a delete makes the key unreadable afterwards.
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
 * Reports a read of a key that nothing contains.
 *
 * @remarks
 * The error is the contract's own sentinel, which the delete law compares
 * against.
 */
export const gone = new Error("deleteremoves: not found");

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=deleteremovestest/
export interface Mixed {
  /**
   * Removes the key, which makes it unreadable.
   *
   * @remarks
   * Without read the law has no observation to make, and without put there
   * is nothing to delete.
   */
  //dokimi:mixin delete-removes read=read sentinel=gone
  delete(key: string): Promise<void>;

  /** Stores the value under the key, which establishes the key that the delete removes. */
  put(key: string, value: string): Promise<void>;

  /** Returns the value under the key, which is how the law observes the removal. */
  read(key: string): Promise<string>;
}
