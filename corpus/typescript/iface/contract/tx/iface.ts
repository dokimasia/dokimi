// Copyright Dokimasia B.V. 2026
// SPDX-License-Identifier: Apache-2.0

/**
 * Declares the contract-axis fixture for the tx contract: a
 * begin/commit/rollback triple threading one handle.
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

/** Represents the handle that begin returns and both terminal operations consume. */
export interface Tx {
  /** Identifies the transaction. */
  id: number;
}

/** Reports a terminal operation on a transaction that the other terminal operation already closed. */
export const txClosed = new Error("tx: closed");

/** Reports a key that no transaction committed. */
export const notFound = new Error("tx: not found");

/** Represents the subject of the fixture. */
//dokimi:stub
//dokimi:suite
//dokimi:out path=txtest/
export interface Contract {
  /**
   * Opens a transaction, and returns the handle that its terminal pair
   * threads.
   *
   * @remarks
   * The handle is the shape a commit-XOR-rollback claim needs, because a
   * mutex over an unnamed transaction is unobservable.
   */
  //dokimi:contract tx role=begin closed=txClosed
  begin(): Promise<Tx>;

  /** Applies the transaction of the handle, terminally. */
  //dokimi:contract tx role=commit
  commit(tx: Tx): Promise<void>;

  /** Discards the transaction of the handle, terminally. */
  //dokimi:contract tx role=rollback
  rollback(tx: Tx): Promise<void>;

  /**
   * Stages a write under the open transaction of the handle.
   *
   * @remarks
   * On the interface rather than beside it, and that is a claim about claims:
   * an interface stating that its staging is invisible outside the
   * transaction is stating something about an operation, and an operation
   * nobody can call through the interface is one no conformance check can
   * reach. Held off it, the mid-transaction law bound a closure reaching past
   * to the concrete store, and no defect worn on the interface could make the
   * claim false — bound, green, and unfalsifiable.
   */
  put(tx: Tx, key: string, v: Value): Promise<void>;

  /**
   * Returns the committed value under the key.
   *
   * @remarks
   * What an open transaction staged must not have touched it, which is the
   * mid-transaction claim's outside read.
   */
  get(key: string): Promise<Value>;
}
