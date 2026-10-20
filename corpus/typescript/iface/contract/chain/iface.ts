// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the chain contract: entries appended
 * to a log, replayed in order, and verified against the digest that links
 * them.
 *
 * @remarks
 * Every role the contract declares is present, because a contract is a
 * multi-callable protocol and a fixture missing a partner exercises the
 * validator's failure path rather than the contract itself.
 *
 * @packageDocumentation
 */

/** Represents one link in the chain. */
export interface Entry {
  /** Contains the key of the link. */
  key: string;
  /** Contains the body of the link. */
  body: string;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=chaintest/
export interface Contract {
  /**
   * Appends the entry, in the append role of the chain contract.
   *
   * @remarks
   * The contract requires a replay beside every append: a log nothing can
   * read back states nothing.
   */
  //dokimi:contract chain role=append
  append(e: Entry): Promise<void>;

  /** Returns the entries in order, in the replay role of the chain contract. */
  //dokimi:contract chain role=replay
  replay(): Promise<Entry[]>;

  /**
   * Recomputes the links, in the verify role of the chain contract.
   *
   * @remarks
   * verify reports a break rather than serving a broken chain.
   */
  //dokimi:contract chain role=verify
  verify(): Promise<void>;
}
