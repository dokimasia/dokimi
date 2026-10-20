// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the fixture for an enumeration that already has the surface.
 *
 * @remarks
 * An author who wrote their own weekdayToString meant to keep it, so it must
 * not be regenerated — a second render beside the enum would contradict the
 * one the author wrote.
 *
 * parseWeekday goes with it. The parse function and the list of the variants
 * are module-level functions rather than members of the enum, so the
 * conservative reading is that a type keeping its own render keeps its own
 * parse, which this fixture declares to prove the pair travels together. What
 * remains is isWeekday, which depends on neither.
 *
 * @packageDocumentation
 */

import { UnknownWeekdayError } from "./errors";

/** Represents a day of the week. */
//dokimi:enum
export enum Weekday {
  Monday = 1,
  Tuesday,
  Wednesday,
}

/**
 * Returns the textual form of a weekday.
 *
 * @remarks
 * The render is deliberately hand-written, and deliberately not what the
 * generator would have derived: lowercase, where the derivation would give
 * `Monday`. A fixture whose hand-written form matched the derived one could
 * not tell "the author's was kept" from "the generator's was emitted".
 */
export function weekdayToString(w: Weekday): string {
  switch (w) {
    case Weekday.Monday:
      return "mon";
    case Weekday.Tuesday:
      return "tue";
    case Weekday.Wednesday:
      return "wed";
  }
  return "unknown";
}

/**
 * Returns the weekday whose textual form is the text, as the author's own
 * parse that matches their render.
 *
 * @throws {@link UnknownWeekdayError} for text that names no weekday.
 */
export function parseWeekday(text: string): Weekday {
  for (const w of [Weekday.Monday, Weekday.Tuesday, Weekday.Wednesday]) {
    if (weekdayToString(w) === text) {
      return w;
    }
  }
  throw new UnknownWeekdayError(text);
}
