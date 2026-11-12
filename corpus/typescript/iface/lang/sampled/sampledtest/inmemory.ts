// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of {@link Labeller}.
 *
 * @packageDocumentation
 */

import { type Labeller, type Level, newTestTag, type Tag, unknownTag } from "../iface";

/**
 * Implements {@link Labeller} in memory over the vocabulary of the tag that
 * {@link newTestTag} returns. The corpus applies the suite of Labeller to it
 * as a subject.
 */
//dokimi:subject Labeller new=newInMemory
export class InMemory implements Labeller {
  readonly #known = new Set<Tag>([newTestTag()]);
  readonly #applied = new Map<Tag, Level>();

  /** Applies the tag at the level, and rejects with {@link unknownTag} for a tag outside the vocabulary. */
  async apply(tag: Tag, level: Level, signal: AbortSignal): Promise<void> {
    signal.throwIfAborted();
    if (!this.#known.has(tag)) {
      throw unknownTag;
    }
    this.#applied.set(tag, level);
  }

  /** Returns the level that apply set for the tag. */
  levelOf(tag: Tag): Level | undefined {
    return this.#applied.get(tag);
  }
}

/** Returns a new {@link InMemory}. */
export function newInMemory(): InMemory {
  return new InMemory();
}
