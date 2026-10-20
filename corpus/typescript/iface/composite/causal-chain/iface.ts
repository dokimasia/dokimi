// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the composite fixture that pairs the chain contract with the
 * causal mixin: an append-only log whose admission policy is causality.
 *
 * @remarks
 * An entry lands only after everything it depends on. The pairing exists
 * because the replay-causality rule needs both stamps on one method: the
 * chain says the log replays, the causal claim says the replay respects the
 * dependency order, and neither alone states the rule.
 *
 * @packageDocumentation
 */

/**
 * Returns the declared Entry of the fixture.
 *
 * @remarks
 * Declared rather than derived, because dependsOn is an admission
 * PRECONDITION and not a payload. A derived entry fills every member it can
 * and so names a dependency that never landed, which this log refuses whole:
 * the append fails, the replay finds nothing, and the chain check skips
 * having asserted anything. An entry depending on nothing is the one that
 * always lands, which is what a claim about the log's length needs.
 */
export function rootEntry(): Entry {
  return { id: "root", dependsOn: [], body: "held", seq: 0 };
}

/**
 * Represents the element of the log: an identifier, the entries it depends
 * on, the payload they order, and the position the log stamped it at.
 *
 * @remarks
 * seq is the store's, not the caller's. causal's version= names the member
 * carrying the ordering stamp, and a stamp the caller chose would make the
 * replay-order claim a question about the check's own arithmetic rather than
 * about the log.
 */
//dokimi:sample value=rootEntry()
export interface Entry {
  /** Identifies the entry. */
  id: string;
  /** Lists the entries that this entry depends on. */
  dependsOn: string[];
  /** Contains the payload of the entry. */
  body: string;
  /** Contains the position that the log stamped the entry at. */
  seq: number;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=causalchaintest/
export interface Log {
  /**
   * Admits an entry only after everything it depends on, and stamps it with
   * the position it took.
   *
   * @remarks
   * Both stamps together are the claim: chain says the log replays, causal
   * says which member orders it, and the replay coming back in that order is
   * what neither states alone.
   */
  //dokimi:contract chain role=append
  //dokimi:mixin causal version=seq
  append(e: Entry): Promise<void>;

  /** Returns the entries in order, in the replay role of the chain contract. */
  //dokimi:contract chain role=replay
  replay(): Promise<Entry[]>;
}
