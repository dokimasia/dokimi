// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of {@link Ranked}.
 *
 * @packageDocumentation
 */

import type { Ordered, Ranked, Score } from "../iface";

/** Reports a key without a rank. */
export const notFound = new Error("genericboundtest: not found");

/**
 * Implements {@link Ranked} in memory. The corpus applies the suite of Ranked
 * to the ranks of scores under numbers that {@link newInMemory} returns.
 */
//dokimi:subject Ranked new=newInMemory
export class InMemory<K extends Ordered, V> implements Ranked<K, V> {
  readonly #ranks = new Map<K, V>();

  /** Returns the value of the rank of the key, and rejects with {@link notFound} for a key without one. */
  async rank(key: K, signal: AbortSignal): Promise<V> {
    signal.throwIfAborted();
    if (!this.#ranks.has(key)) {
      throw notFound;
    }
    return this.#ranks.get(key) as V;
  }

  /** Clears the ranks. */
  async reset(signal: AbortSignal): Promise<void> {
    signal.throwIfAborted();
    this.#ranks.clear();
  }

  /** Sets the rank of the key. A test seeds the ranks through it, because the interface declares no writer. */
  set(key: K, value: V): void {
    this.#ranks.set(key, value);
  }
}

/** Returns a new {@link InMemory} of scores under numbers. */
export function newInMemory(): InMemory<number, Score> {
  return new InMemory<number, Score>();
}
