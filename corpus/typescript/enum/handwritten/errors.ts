// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Reports text that names no weekday.
 *
 * @remarks
 * The author's own parse error, named the way the generator would have named
 * it.
 */
export class UnknownWeekdayError extends Error {
  /** Contains the text that names no weekday. */
  readonly text: string;

  constructor(text: string) {
    super(`handwritten: unknown weekday ${JSON.stringify(text)}`);
    this.name = "UnknownWeekdayError";
    this.text = text;
  }
}
