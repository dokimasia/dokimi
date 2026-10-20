// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the mixin-axis fixture for the errors mixin, which states that the
 * method reports misses through a declared sentinel.
 *
 * @remarks
 * The interface carries whatever methods the law needs to be stateable. A
 * mixin whose law spans two calls cannot be hosted by a single method: there
 * would be nothing to compare against, and the generated subtest would pass
 * by having nothing to check.
 *
 * A fixture without a declared sentinel leaves the generated assertion
 * nothing to compare against.
 *
 * There is no negated form here. eidos declares the mixin directive
 * DenyNegation, because a mixin is opt-in and deleting the directive is the
 * suppression.
 *
 * @packageDocumentation
 */

/** Reports a key that get misses, and is the sentinel the mixin declares. */
export const notFound = new Error("errors: not found");

/**
 * Reports a key that was removed.
 *
 * @remarks
 * A second sentinel exercises the list of the fault directive in a fixture
 * rather than only in a unit test.
 */
export const gone = new Error("errors: gone");

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=errorstest/
export interface Mixed {
  /**
   * Returns the value under the key.
   *
   * @remarks
   * get reports {@link notFound} and nothing else for a miss. The sentinel
   * has to be declared in the module for the assertion to reference it.
   *
   * The notfound mixin is what makes the miss derivable: nothing on this
   * interface writes, so without a declared sentinel no draw would be one
   * nothing supplied and the rule refuses. The fault directive is a different
   * key — it says what the DOUBLE can be told to return, and its list names
   * both sentinels on one line, each earning its own one-shot helper.
   */
  //dokimi:mixin errors
  //dokimi:mixin not-found sentinel=notFound
  //dokimi:fault errors=[notFound, gone]
  get(key: string): Promise<string>;
}
