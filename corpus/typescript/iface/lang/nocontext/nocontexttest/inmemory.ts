// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of {@link Calculator}.
 *
 * @packageDocumentation
 */

import { type Calculator, divideByZero } from "../iface";

/**
 * Implements {@link Calculator} in memory, and counts its calls. The corpus
 * applies the suite of Calculator to it as a subject.
 */
//dokimi:subject Calculator new=newInMemory
export class InMemory implements Calculator {
  #calls = 0;

  /** Returns the sum of a and b. */
  add(a: number, b: number): number {
    this.#calls += 1;
    return a + b;
  }

  /**
   * Returns the quotient of a and b.
   *
   * @throws {@link divideByZero} for a zero divisor.
   */
  divide(a: number, b: number): number {
    this.#calls += 1;
    if (b === 0) {
      throw divideByZero;
    }
    return a / b;
  }

  /** Clears the count of the calls. */
  reset(): void {
    this.#calls = 0;
  }

  /** Returns the count of the calls since the last reset. */
  calls(): number {
    return this.#calls;
  }
}

/** Returns a new {@link InMemory}. */
export function newInMemory(): InMemory {
  return new InMemory();
}
