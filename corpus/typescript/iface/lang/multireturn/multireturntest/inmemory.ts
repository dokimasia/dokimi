// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of {@link Wide}.
 *
 * @packageDocumentation
 */

import type { Wide } from "../iface";

/** Reports an id without a name. */
export const notFound = new Error("multireturntest: not found");

/**
 * Implements {@link Wide} in memory. The corpus applies the suite of Wide to
 * it as a subject.
 */
//dokimi:subject Wide new=newInMemory
export class InMemory implements Wide {
  readonly #items = new Map<string, string>();

  /** Stores the name under the id. A test seeds the names through it, because the interface declares no writer. */
  put(id: string, name: string): void {
    this.#items.set(id, name);
  }

  /** Returns the name, its length and true for the id, and rejects with {@link notFound} for an id without a name. */
  async quad(id: string, signal: AbortSignal): Promise<[string, number, boolean]> {
    signal.throwIfAborted();
    const name = this.#items.get(id);
    if (name === undefined) {
      throw notFound;
    }
    return [name, name.length, true];
  }

  /** Returns the name, its length and whether the id has a name. */
  async triple(id: string, signal: AbortSignal): Promise<[string, number, boolean]> {
    signal.throwIfAborted();
    const name = this.#items.get(id);
    return [name ?? "", name?.length ?? 0, name !== undefined];
  }

  /** Returns the name and its length for the id. */
  async noError(id: string, signal: AbortSignal): Promise<[string, number]> {
    signal.throwIfAborted();
    const name = this.#items.get(id) ?? "";
    return [name, name.length];
  }
}

/** Returns a new {@link InMemory}. */
export function newInMemory(): InMemory {
  return new InMemory();
}
