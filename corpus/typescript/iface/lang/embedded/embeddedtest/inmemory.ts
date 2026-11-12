// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of {@link Composed}.
 *
 * @packageDocumentation
 */

import type { Base, Closer, Composed } from "../iface";

/** Reports an identity without a value. */
export const notFound = new Error("embeddedtest: not found");

/**
 * Implements {@link Composed} in memory. The corpus applies the suite of
 * Composed to it as a subject.
 */
//dokimi:subject Composed new=newInMemory
export class InMemory implements Base, Closer, Composed {
  readonly #items = new Map<string, string>();

  /** Returns the value under the key, and rejects with {@link notFound} for a key without one. */
  async get(key: string, signal: AbortSignal): Promise<string> {
    signal.throwIfAborted();
    const value = this.#items.get(key);
    if (value === undefined) {
      throw notFound;
    }
    return value;
  }

  /** Stores the value under the key. A test seeds the store through it, because the interface declares no writer. */
  put(key: string, value: string): void {
    this.#items.set(key, value);
  }

  /** Checks that the subject is reachable. */
  async ping(signal: AbortSignal): Promise<void> {
    signal.throwIfAborted();
  }

  /** Releases the subject. A second call returns as well. */
  async close(signal: AbortSignal): Promise<void> {
    signal.throwIfAborted();
  }
}

/** Returns a new {@link InMemory}. */
export function newInMemory(): InMemory {
  return new InMemory();
}
