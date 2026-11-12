// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of {@link Counter}.
 *
 * @packageDocumentation
 */

import type { Counter, Node, Step } from "../iface";

/**
 * Implements {@link Counter}. It walks each value to its end, so each value of
 * a recursive type that the suite draws has to end. The corpus applies the
 * suite of Counter to it as a subject.
 */
//dokimi:subject Counter new=newInMemory
export class InMemory implements Counter {
  /** Counts the nodes of the tree under root. */
  nodes(root: Node): number {
    return 1 + root.children().reduce((n, child) => n + this.nodes(child), 0);
  }

  /** Counts the steps of a run of the machine from start. */
  steps(start: Step): number {
    let n = 0;
    for (let s: Step | undefined = start; s !== undefined; s = s("next")) {
      n += 1;
    }
    return n;
  }
}

/** Returns a new {@link InMemory}. */
export function newInMemory(): InMemory {
  return new InMemory();
}
