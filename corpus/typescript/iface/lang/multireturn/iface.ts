// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for more than two returned values, past
 * the point where a generator can assume a single value.
 *
 * @remarks
 * A TypeScript method returns one value, so the values are the elements of a
 * tuple. This axis varies the type system rather than the classification, so
 * these break generators independently of any directive.
 *
 * @packageDocumentation
 */

/**
 * Represents an interface whose methods return tuples.
 *
 * @remarks
 * Anything that assumes the result is a single scalar produces wrong samples
 * and wrong types here.
 */
//dokimi:stub
//dokimi:suite
//dokimi:out path=multireturntest/
export interface Wide {
  /** Returns the name, the count and the flag of the id, or rejects. */
  quad(id: string): Promise<[string, number, boolean]>;
  /** Returns the name, the count and the flag of the id. */
  triple(id: string): Promise<[string, number, boolean]>;
  /** Returns the name and the count of the id. */
  noError(id: string): Promise<[string, number]>;
}
