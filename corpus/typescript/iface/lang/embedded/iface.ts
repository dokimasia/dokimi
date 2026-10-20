// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for an interface that extends others, so
 * the member set is larger than the declaration.
 *
 * @remarks
 * This axis varies the type system rather than the classification, so these
 * break generators independently of any directive.
 *
 * Every base here is declared in this module. The embeddedforeign fixture
 * covers the base from outside the tree, and the two must agree — which is
 * the property having both is for.
 *
 * @packageDocumentation
 */

/**
 * Reports a subject that ping cannot reach.
 *
 * @remarks
 * The sentinel is declared here, so a fault helper generated onto a double of
 * an extending interface has something to import.
 */
export const unreachable = new Error("embedded: unreachable");

/** Represents the base that {@link Composed} extends. */
//dokimi:stub
//dokimi:out path=embeddedtest/
export interface Base {
  /**
   * Checks that the subject is reachable.
   *
   * @remarks
   * An extending interface's double implements ping without declaring it. A
   * contributor reading only the declarations would silently skip it.
   */
  ping(): Promise<void>;
}

/** Represents the base that {@link Composed} extends beside {@link Base}. */
//dokimi:stub
//dokimi:out path=embeddedtest/
export interface Closer {
  /** Releases the subject. */
  close(): Promise<void>;
}

/**
 * Represents an interface that extends two interfaces and adds one method.
 *
 * @remarks
 * A generator that reads only the declared methods emits an incomplete stub,
 * which fails the check that the double implements the interface.
 */
//dokimi:stub
//dokimi:suite
//dokimi:out path=embeddedtest/
export interface Composed extends Base, Closer {
  /** Returns the value under the key. */
  get(key: string): Promise<string>;
}
