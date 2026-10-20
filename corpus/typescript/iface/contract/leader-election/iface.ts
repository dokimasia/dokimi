// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the leader-election contract: a
 * leadership protocol.
 *
 * @remarks
 * Every role the contract declares is present, because a contract is a
 * multi-callable protocol and a fixture missing a partner exercises the
 * validator's failure path rather than the contract itself.
 *
 * The model tier rides the twin floor. A lease tracker is the near miss: it
 * models exclusion, but it is keyed by what acquire takes and campaign takes
 * nothing, so the store has no key to keep the leadership under.
 *
 * @packageDocumentation
 */

/** Represents the payload that the roles of the contract exchange. */
export interface Value {
  /** Contains the key of the payload. */
  key: string;
  /** Contains the body of the payload. */
  body: string;
}

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=leaderelectiontest/
export interface Contract {
  /** Campaigns for the leadership, in the campaign role of the leader-election contract. */
  //dokimi:contract leader-election role=campaign
  campaign(): Promise<void>;

  /** Gives the leadership up, in the resign role of the leader-election contract. */
  //dokimi:contract leader-election role=resign
  resign(): Promise<void>;

  /** Reports whether the subject leads, in the is-leader role of the leader-election contract. */
  //dokimi:contract leader-election role=is-leader
  isLeader(): Promise<boolean>;
}
