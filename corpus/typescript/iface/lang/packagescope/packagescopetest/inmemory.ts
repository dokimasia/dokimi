// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of {@link Slugger}.
 *
 * @packageDocumentation
 */

import type { Slugger } from "../iface";

/**
 * Implements {@link Slugger} with the methods of a string. The corpus applies
 * the suite of Slugger to it as a subject.
 */
//dokimi:subject Slugger new=newInMemory
export class InMemory implements Slugger {
  /** Returns the words of title in lower case, joined by hyphens. */
  slug(title: string): string {
    return title.toLowerCase().split(/\s+/).filter((word) => word !== "").join("-");
  }
}

/** Returns a new {@link InMemory}. */
export function newInMemory(): InMemory {
  return new InMemory();
}
