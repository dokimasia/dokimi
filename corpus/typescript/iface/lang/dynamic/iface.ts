// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the language-axis fixture for parameters of the type unknown and
 * for the values that a shape states without a declaration of the module.
 *
 * @packageDocumentation
 */

/** Represents a diary that records messages with values of any type, in a zone and at a scale. */
//dokimi:suite
//dokimi:out path=dynamictest/
export interface Diary {
  /** Records the message with its values. */
  print(msg: string, ...values: unknown[]): void;
  /** Records the value of the field. */
  field(name: string, value: unknown): void;
  /** Moves the diary into the time zone with the identifier zone, such as Europe/Amsterdam. */
  in(zone: string): void;
  /** Sets the scale of the diary to n. */
  scale(n: bigint): void;
  /** Returns the number of records of the diary. */
  count(): number;
}
