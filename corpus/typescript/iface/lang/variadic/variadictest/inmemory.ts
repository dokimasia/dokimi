// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of {@link Finder}.
 *
 * @packageDocumentation
 */

import type { Finder } from "../iface";

/** Reports a find without keys. */
export const noKeys = new Error("variadictest: no keys requested");

/**
 * Implements {@link Finder} in memory. The corpus applies the suite of Finder
 * to it as a subject.
 */
//dokimi:subject Finder new=newInMemory
export class InMemory implements Finder {
  readonly #items = new Map<string, string>();

  /** Stores the value under the key. A test seeds the finder through it, because the interface declares no writer. */
  put(key: string, value: string): void {
    this.#items.set(key, value);
  }

  /** Returns the values of the keys that the finder knows, and rejects with {@link noKeys} for no key. */
  async find(signal: AbortSignal, ...keys: string[]): Promise<string[]> {
    signal.throwIfAborted();
    if (keys.length === 0) {
      throw noKeys;
    }
    return keys.flatMap((k) => {
      const value = this.#items.get(k);
      return value === undefined ? [] : [value];
    });
  }

  /** Returns at most limit of the values that find returns for the keys. A negative limit returns them all. */
  async findWithLimit(limit: number, signal: AbortSignal, ...keys: string[]): Promise<string[]> {
    const found = await this.find(signal, ...keys);
    return limit >= 0 && found.length > limit ? found.slice(0, limit) : found;
  }
}

/** Returns a new {@link InMemory}. */
export function newInMemory(): InMemory {
  return new InMemory();
}
