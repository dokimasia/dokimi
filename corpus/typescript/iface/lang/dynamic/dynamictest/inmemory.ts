// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of {@link Diary}.
 *
 * @packageDocumentation
 */

import type { Diary } from "../iface";

/**
 * Implements {@link Diary} in memory. The corpus applies the suite of Diary
 * to it as a subject.
 */
//dokimi:subject Diary new=newInMemory
export class InMemory implements Diary {
  readonly #records: string[] = [];
  #zone = "UTC";
  #scale = 1n;

  /** Records the message. */
  print(msg: string, ..._values: unknown[]): void {
    this.#records.push(msg);
  }

  /** Records the name of the field. */
  field(name: string, _value: unknown): void {
    this.#records.push(name);
  }

  /** Moves the diary into the time zone. */
  in(zone: string): void {
    this.#zone = zone;
  }

  /** Sets the scale of the diary to n. */
  scale(n: bigint): void {
    this.#scale = n;
  }

  /** Returns the number of records of the diary. */
  count(): number {
    return this.#records.length;
  }

  /** Returns the zone and the scale of the diary. */
  settings(): [string, bigint] {
    return [this.#zone, this.#scale];
  }
}

/** Returns an empty {@link InMemory} in the zone UTC at the scale 1. */
export function newInMemory(): InMemory {
  return new InMemory();
}
