// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of {@link Store}.
 *
 * @packageDocumentation
 */

import type { Rec, Store } from "../iface";

/**
 * Implements {@link Store} in memory. The corpus applies the suite of Store
 * to it as a subject.
 */
//dokimi:subject Store new=newInMemory
export class InMemory implements Store {
  readonly #weights = new Map<string, number>();
  readonly #at = new Map<string, number>();
  readonly #buckets = new Map<number, number>();
  readonly #recs = new Map<string, Rec>();
  #rev = 0;

  /** Adds the weight to the key. */
  async touch(key: string, weight: number, signal: AbortSignal): Promise<void> {
    signal.throwIfAborted();
    this.#weights.set(key, (this.#weights.get(key) ?? 0) + weight);
  }

  /** Returns the weight of the key. */
  async seen(k: string, signal: AbortSignal): Promise<number> {
    signal.throwIfAborted();
    return this.#weights.get(k) ?? 0;
  }

  /** Moves one unit from one identifier to the other. */
  async move(from: string, to: string, signal: AbortSignal): Promise<void> {
    signal.throwIfAborted();
    this.#at.set(from, (this.#at.get(from) ?? 0) - 1);
    this.#at.set(to, (this.#at.get(to) ?? 0) + 1);
  }

  /** Returns the units at the identifier. */
  async at(where: string, signal: AbortSignal): Promise<number> {
    signal.throwIfAborted();
    return this.#at.get(where) ?? 0;
  }

  /** Records an event in the bucket of the length of the id. */
  async emit(id: string, signal: AbortSignal): Promise<void> {
    signal.throwIfAborted();
    this.#buckets.set(id.length, (this.#buckets.get(id.length) ?? 0) + 1);
  }

  /** Returns how many events the bucket contains. */
  async count(bucket: number, signal: AbortSignal): Promise<number> {
    signal.throwIfAborted();
    return this.#buckets.get(bucket) ?? 0;
  }

  /** Returns the record under the key, and an empty record for a key without one. */
  async version(key: string, signal: AbortSignal): Promise<Rec> {
    signal.throwIfAborted();
    return this.#recs.get(key) ?? { key: "", rev: 0 };
  }

  /** Stores the record and stamps its revision. */
  async put(r: Rec, signal: AbortSignal): Promise<void> {
    signal.throwIfAborted();
    this.#store(r);
  }

  /** Stores the record in place of the record under its key and stamps its revision. */
  async replace(r: Rec, signal: AbortSignal): Promise<void> {
    signal.throwIfAborted();
    this.#store(r);
  }

  /** Stores a copy of the record with the next revision. */
  #store(r: Rec): void {
    this.#rev += 1;
    this.#recs.set(r.key, { ...r, rev: this.#rev });
  }
}

/** Returns a new {@link InMemory}. */
export function newInMemory(): InMemory {
  return new InMemory();
}
