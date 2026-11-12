// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of {@link Store}.
 *
 * @packageDocumentation
 */

import type { Key, Payload, Store } from "../iface";

/** Reports a key without a payload. */
export const notFound = new Error("roledtypestest: not found");

/**
 * Implements {@link Store} in memory. The corpus applies the suite of Store
 * to it as a subject.
 */
//dokimi:subject Store new=newInMemory
export class InMemory implements Store {
  readonly #payloads = new Map<Key, Payload>();

  /** Stores the payload under the key. */
  async put(key: Key, payload: Payload, signal: AbortSignal): Promise<void> {
    signal.throwIfAborted();
    this.#payloads.set(key, payload);
  }

  /** Returns the payload under the key, and rejects with {@link notFound} for a key without one. */
  async get(key: Key, signal: AbortSignal): Promise<Payload> {
    signal.throwIfAborted();
    const payload = this.#payloads.get(key);
    if (payload === undefined) {
      throw notFound;
    }
    return payload;
  }
}

/** Returns a new {@link InMemory}. */
export function newInMemory(): InMemory {
  return new InMemory();
}
