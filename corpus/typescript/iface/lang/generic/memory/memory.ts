// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of the generic store.
 *
 * @packageDocumentation
 */

import type { Store as Contract } from "../iface";

/** Reports a key without a value. */
export const notFound = new Error("memory: the store has no value under the key");

/**
 * Implements the generic store in memory for keys of the type K and values
 * of the type V. The corpus applies the suite of the store to the store of
 * numbers that {@link newInts} returns.
 */
//dokimi:subject Contract new=newInts
export class Store<K, V> implements Contract<K, V> {
  readonly #values = new Map<K, V>();

  /** Returns the value under the key, and rejects with {@link notFound} for a key without one. */
  async get(key: K, signal: AbortSignal): Promise<V> {
    signal.throwIfAborted();
    if (!this.#values.has(key)) {
      throw notFound;
    }
    return this.#values.get(key) as V;
  }

  /** Stores the value under the key. */
  async put(key: K, value: V, signal: AbortSignal): Promise<void> {
    signal.throwIfAborted();
    this.#values.set(key, value);
  }
}

/** Returns a new store of numbers under keys of numbers. */
export function newInts(): Store<number, number> {
  return new Store<number, number>();
}
