// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the subject of the suite of {@link Failure}.
 *
 * @packageDocumentation
 */

import type { Failure } from "../iface";

/**
 * Implements {@link Failure} as an Error that no retry can resolve. The corpus
 * applies the suite of Failure to it as a subject.
 */
//dokimi:subject Failure new=newInMemory
export class InMemory extends Error implements Failure {
  /** Returns a failure with its message. */
  constructor() {
    super("embeddedforeigntest: the call failed");
    this.name = "InMemory";
  }

  /** Reports that a retry of the failed call does not succeed. */
  retryable(): boolean {
    return false;
  }
}

/** Returns a new {@link InMemory}. */
export function newInMemory(): InMemory {
  return new InMemory();
}
