// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the transaction contract: a callable
 * running inside a transactional scope.
 *
 * @remarks
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

/** Reports a key that no transaction committed. */
export const notFound = new Error("transaction: not found");

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=transactiontest/
export interface Contract {
  /**
   * Runs the body inside a transactional scope.
   *
   * @remarks
   * run is the transaction contract's fn role. Its directive names the write
   * and the read that check a rollback. It accepts the body it scopes — the
   * shape a rollback claim needs, because the claim is about what a failing
   * body leaves behind, and a run that takes no body cannot be made to fail
   * on demand.
   */
  //dokimi:contract transaction role=fn not-found=notFound write=put read=get
  run(body: () => Promise<void>): Promise<void>;

  /**
   * Stores the value under the key, as the write that a body performs inside
   * the scope that run opened.
   *
   * @remarks
   * put is the reason the rollback claim is checkable at all: a claim about
   * what a failing body leaves behind needs the body to leave something
   * behind. Called outside a run it writes through, which is the ordinary
   * store behaviour a keyed writer promises; called inside one it stages, and
   * a failing body discards the staging whole.
   */
  put(key: string, v: Value): Promise<void>;

  /** Returns the committed value under the key, which a rolled-back run must not have touched. */
  get(key: string): Promise<Value>;
}
