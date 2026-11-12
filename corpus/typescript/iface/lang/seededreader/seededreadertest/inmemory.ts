// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of {@link Catalog}.
 *
 * @packageDocumentation
 */

import type { Body, Catalog, Key } from "../iface";

/**
 * Implements {@link Catalog} in memory over the documents that it loads at its
 * construction. The corpus applies the suite of Catalog to it as a subject.
 */
//dokimi:subject Catalog new=newInMemory
export class InMemory implements Catalog {
  readonly #docs: ReadonlyMap<Key, Body>;

  /** Returns a catalog of a copy of the documents. */
  constructor(docs: ReadonlyMap<Key, Body>) {
    this.#docs = new Map(docs);
  }

  /** Returns the body of the key, and an empty body for a key without a document. */
  async lookup(key: Key, signal: AbortSignal): Promise<Body> {
    signal.throwIfAborted();
    return this.#docs.get(key) ?? "";
  }

  /** Returns the count of the documents. */
  async len(signal: AbortSignal): Promise<number> {
    signal.throwIfAborted();
    return this.#docs.size;
  }
}

/** Returns a new {@link InMemory} with the document of the defaults of Key and Body. */
export function newInMemory(): InMemory {
  return newInMemoryOf(new Map([["test-doc", "test-contents"]]));
}

/** Returns a new {@link InMemory} with the documents. */
export function newInMemoryOf(docs: ReadonlyMap<Key, Body>): InMemory {
  return new InMemory(docs);
}
