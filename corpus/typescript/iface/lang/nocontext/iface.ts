// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for synchronous methods, which every
 * assertion about a promise has to skip rather than mis-apply.
 *
 * @remarks
 * This axis varies the type system rather than the classification, so these
 * break generators independently of any directive.
 *
 * @packageDocumentation
 */

/** Reports a zero divisor. */
export const divideByZero = new Error("nocontext: divide by zero");

/**
 * Represents a calculator that returns no promise anywhere.
 *
 * @remarks
 * Shapes are defined with the promise optional, so the same detectors must
 * fire here as on the asynchronous forms.
 */
//dokimi:stub
//dokimi:suite
//dokimi:out path=nocontexttest/
export interface Calculator {
  /** Returns the sum of a and b. */
  add(a: number, b: number): number;

  /**
   * Returns the quotient of a and b.
   *
   * @throws {@link divideByZero} for a zero divisor.
   */
  divide(a: number, b: number): number;

  /** Clears the state of the calculator. */
  reset(): void;
}
