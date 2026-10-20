// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the saga contract: two steps, each
 * paired with the compensation that undoes it.
 *
 * @remarks
 * Two pairs because compensation is an ORDER. A saga that fails partway
 * unwinds the steps before the failure, last first, and leaves the failed
 * one alone — none of which a single-step fixture can state.
 *
 * Every role the contract declares is present, because a contract is a
 * multi-callable protocol and a fixture missing a partner exercises the
 * validator's failure path rather than the contract itself.
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
//dokimi:out path=sagatest/
export interface Contract {
  /** Applies the first step, in the step role of the saga contract. */
  //dokimi:contract saga role=step id=step
  step(v: Value): Promise<void>;

  /** Undoes the first step, in the compensate role of the saga contract. */
  //dokimi:contract saga role=compensate id=step
  compensate(v: Value): Promise<void>;

  /**
   * Applies a SECOND step, with its own compensator.
   *
   * @remarks
   * Two pairs rather than one, and the ordering is the reason. A saga
   * compensates in reverse from the step that failed, so a single-step
   * fixture leaves that whole claim unstated: there is nothing before the
   * failure to undo. With two, a failure in the second has to unwind the
   * first and only the first.
   *
   * The id ship separates this pair from the pair of step.
   */
  //dokimi:contract saga role=step id=ship
  ship(v: Value): Promise<void>;

  /** Undoes what ship applied. */
  //dokimi:contract saga role=compensate id=ship
  undoShip(v: Value): Promise<void>;

  /**
   * Returns the fingerprint of the applied steps, in order.
   *
   * @remarks
   * The fingerprint is the observation a compensation claim compares before
   * and after a failed run. A string rather than a listing because the
   * comparison is equality, and one composed value cannot half-match.
   */
  state(): Promise<string>;
}
